import { useEffect, useState, useMemo, memo, useCallback, useRef } from 'react'
import { Spin, Empty, Card, Alert } from 'antd'
import * as relationGraphApi from '../../api/relation-graph'
import type { RelationGraphResp, GraphNode, GraphEdge } from '../../types/relation-graph'
import { GENDER_ICON } from '../../constants/enums'
import './RelationGraph.css'

const CARD_W = 128
const CARD_H = 64
const COL_STEP = CARD_W + 40
const ROW_STEP = CARD_H + 56
const PAD = 28

function labelOf(rt: string): string {
  const m: Record<string, string> = {
    spouse: '夫妻', parent: '亲子',
    son: '儿子', daughter: '女儿', brother: '兄弟', sister: '姐妹',
    grandchild: '孙辈', father: '父亲', mother: '母亲',
    grandfather: '祖父', grandmother: '祖母', elder: '长辈', follower: '其他',
  }
  return m[rt] || rt
}

interface PositionedNode {
  id: string
  node: GraphNode
  x: number
  y: number
}

interface ConnLine {
  fromX: number; fromY: number
  toX: number; toY: number
  fromId: string; toId: string
  relationType: string
  isSpouse: boolean
  isUserEdge: boolean
  confirmed: boolean
}

/* ========== 布局算法 ========== */

function buildData(nodes: GraphNode[], edges: GraphEdge[]) {
  const nodeMap = new Map<string, GraphNode>()
  nodes.forEach((n) => nodeMap.set(n.id, n))

  const spouseOf = new Map<string, string>()
  const childrenOf = new Map<string, string[]>()
  const parentOf = new Map<string, string>()

  for (const e of edges) {
    if (e.relationType === 'spouse' && e.confirmed) {
      spouseOf.set(e.from, e.to)
      spouseOf.set(e.to, e.from)
    }
    if (e.relationType === 'parent' && e.confirmed) {
      if (!childrenOf.has(e.from)) childrenOf.set(e.from, [])
      childrenOf.get(e.from)!.push(e.to)
      parentOf.set(e.to, e.from)
    }
  }

  // 确保 children 顺序稳定
  for (const [, kids] of childrenOf) kids.sort()

  const roots = nodes.filter((n) => !parentOf.has(n.id) && n.type === 'deceased').map((n) => n.id)
  roots.sort()

  return { nodeMap, spouseOf, childrenOf, parentOf, roots }
}

type ColMap = Map<string, number>

function layoutSubtree(
  nodeId: string,
  startCol: number,
  childrenOf: Map<string, string[]>,
  spouseOf: Map<string, string>,
  colMap: ColMap,
  visited: Set<string>
): number {
  if (visited.has(nodeId)) return startCol
  visited.add(nodeId)

  const kids = childrenOf.get(nodeId) ?? []
  const spouseId = spouseOf.get(nodeId)

  // 递归布局所有子女
  let childEnd = startCol
  for (const kid of kids) {
    childEnd = layoutSubtree(kid, childEnd, childrenOf, spouseOf, colMap, visited)
  }

  const hasChildren = kids.length > 0
  const hasSpouse = !!spouseId && !visited.has(spouseId)

  let myCol: number
  if (hasChildren) {
    myCol = Math.floor((startCol + childEnd - 2) / 2)
  } else {
    myCol = startCol
  }
  if (myCol < startCol) myCol = startCol

  colMap.set(nodeId, myCol)

  if (hasSpouse) {
    const spouseCol = myCol + 1
    // 确保不和下一个子树重叠
    const spouseOccupies = Math.max(spouseCol, childEnd > startCol ? childEnd - 1 : spouseCol)
    colMap.set(spouseId, spouseCol)
    visited.add(spouseId)
    return spouseOccupies + 1
  }

  return Math.max(myCol + 1, childEnd)
}

function computeLayout(nodes: GraphNode[], edges: GraphEdge[]) {
  const { nodeMap, spouseOf, childrenOf, parentOf, roots } = buildData(nodes, edges)
  const colMap: ColMap = new Map()
  const visited = new Set<string>()

  let nextCol = 0
  for (const root of roots) {
    nextCol = layoutSubtree(root, nextCol, childrenOf, spouseOf, colMap, visited)
    // 子树之间加一列间隔
    nextCol += 1
  }

  // 处理用户节点：放在关联逝者右侧
  const positioned: PositionedNode[] = []
  const placed = new Set<string>()

  for (const node of nodes) {
    const col = colMap.get(node.id)
    if (col !== undefined) {
      const row = Math.max(0, node.generation)
      positioned.push({
        id: node.id,
        node,
        x: PAD + col * COL_STEP + CARD_W / 2,
        y: PAD + row * ROW_STEP + CARD_H / 2,
      })
      placed.add(node.id)
    }
  }

  // 构建已占用位置网格（逝者节点），供用户节点碰撞避让
  const PAD_GAP = 8
  type Rect = { x1: number; x2: number }
  const rectsOverlap = (a: Rect, b: Rect) => a.x1 < b.x2 && a.x2 > b.x1
  const occupiedByRow = new Map<number, Rect[]>()
  for (const p of positioned) {
    const r = Math.round((p.y - PAD - CARD_H / 2) / ROW_STEP)
    if (!occupiedByRow.has(r)) occupiedByRow.set(r, [])
    occupiedByRow.get(r)!.push({ x1: p.x - CARD_W / 2 - PAD_GAP, x2: p.x + CARD_W / 2 + PAD_GAP })
  }
  const isSlotFree = (row: number, cx: number) => {
    const rects = occupiedByRow.get(row) || []
    const r: Rect = { x1: cx - CARD_W / 2 - PAD_GAP, x2: cx + CARD_W / 2 + PAD_GAP }
    return !rects.some((o) => rectsOverlap(o, r))
  }
  const occupySlot = (row: number, cx: number) => {
    if (!occupiedByRow.has(row)) occupiedByRow.set(row, [])
    occupiedByRow.get(row)!.push({ x1: cx - CARD_W / 2 - PAD_GAP, x2: cx + CARD_W / 2 + PAD_GAP })
  }
  const findFreeX = (row: number, prefX: number, step: number) => {
    let x = prefX
    let tries = 0
    while (!isSlotFree(row, x) && tries < 50) { x += step; tries++ }
    return x
  }

  // 用户节点：按用户分组，每个用户选最佳锚点逝者
  const userEdges = edges.filter((e) => {
    const a = nodeMap.get(e.from)
    const b = nodeMap.get(e.to)
    return a && b && !a.isDead !== !b.isDead
  })

  const userToDeceased = new Map<string, string[]>()
  for (const e of userEdges) {
    const userId = nodeMap.get(e.from)!.isDead ? e.to : e.from
    const decId = userId === e.from ? e.to : e.from
    if (!userToDeceased.has(userId)) userToDeceased.set(userId, [])
    if (!userToDeceased.get(userId)!.includes(decId)) userToDeceased.get(userId)!.push(decId)
  }

  const anchorStack = new Map<string, number>()
  const userAnchor = new Map<string, string>()

  for (const [userId, decIds] of userToDeceased) {
    if (placed.has(userId)) continue

    let anchorId = decIds[0]
    let bestPri = 999
    let bestRt = ''
    for (const did of decIds) {
      const dp = positioned.find((p) => p.id === did)
      if (!dp) continue
      const edge = userEdges.find((e) =>
        (e.from === userId && e.to === did) || (e.to === userId && e.from === did),
      )
      const rt = edge?.relationType || ''
      let pri = 50
      if (rt === 'brother' || rt === 'sister' || rt === 'spouse') pri = 0
      else if (rt === 'son' || rt === 'daughter') pri = 10
      else if (rt === 'grandchild') pri = 11
      else if (rt === 'father' || rt === 'mother') pri = 20
      else if (rt === 'grandfather' || rt === 'grandmother') pri = 30
      else if (rt === 'elder') pri = 15
      if (pri < bestPri) { bestPri = pri; anchorId = did; bestRt = rt }
    }

    userAnchor.set(userId, anchorId)

    const anchorPos = positioned.find((p) => p.id === anchorId)
    if (!anchorPos) continue
    const userNode = nodeMap.get(userId)
    if (!userNode) continue

    const anchorRow = Math.round((anchorPos.y - PAD - CARD_H / 2) / ROW_STEP)
    let targetRow: number
    let isLeft: boolean

    if (bestRt === 'son' || bestRt === 'daughter') {
      targetRow = anchorRow + 1; isLeft = false
    } else if (bestRt === 'grandchild') {
      targetRow = anchorRow + 2; isLeft = false
    } else if (bestRt === 'father' || bestRt === 'mother' || bestRt === 'elder') {
      targetRow = anchorRow - 1; isLeft = true
    } else if (bestRt === 'grandfather' || bestRt === 'grandmother') {
      targetRow = anchorRow - 2; isLeft = true
    } else {
      targetRow = anchorRow; isLeft = false
    }

    const stackIdx = anchorStack.get(anchorId) ?? 0
    anchorStack.set(anchorId, stackIdx + 1)

    const idealX = isLeft
      ? anchorPos.x - CARD_W - 28 - stackIdx * (CARD_W + 28)
      : anchorPos.x + CARD_W + 28 + stackIdx * (CARD_W + 28)
    const searchStep = isLeft ? -COL_STEP : COL_STEP

    const finalX = findFreeX(targetRow, idealX, searchStep)
    const finalY = PAD + targetRow * ROW_STEP + CARD_H / 2
    occupySlot(targetRow, finalX)

    positioned.push({
      id: userId,
      node: userNode,
      x: finalX,
      y: finalY,
    })
    placed.add(userId)
  }

  // 未放置节点的 fallback 布局：排放在图表右下方
  const unplaced = nodes.filter((n) => !placed.has(n.id))
  if (unplaced.length > 0) {
    let ux = 0
    let uy = 0
    for (const p of positioned) {
      if (p.x + CARD_W / 2 > ux) ux = p.x + CARD_W / 2 + COL_STEP
      if (p.y + CARD_H / 2 > uy) uy = p.y + CARD_H / 2
    }
    uy += ROW_STEP
    const perRow = Math.max(1, Math.floor((Math.max(...positioned.map((p) => p.x + CARD_W / 2), ux) - PAD * 2) / COL_STEP))
    unplaced.forEach((n, i) => {
      const col = i % perRow
      const row = Math.floor(i / perRow)
      positioned.push({
        id: n.id,
        node: n,
        x: ux + col * COL_STEP + CARD_W / 2,
        y: uy + row * ROW_STEP + CARD_H / 2,
      })
      placed.add(n.id)
    })
  }

  // 检查是否有节点超出左/上边界，整体平移
  const minX = Math.min(PAD, ...positioned.map((p) => p.x - CARD_W / 2))
  if (minX < PAD) {
    const shift = PAD - minX
    for (const p of positioned) p.x += shift
  }
  const minY = Math.min(PAD, ...positioned.map((p) => p.y - CARD_H / 2))
  if (minY < PAD) {
    const shift = PAD - minY
    for (const p of positioned) p.y += shift
  }

  // 构建连线
  const connections: ConnLine[] = []
  const posMap = new Map(positioned.map((p) => [p.id, p]))

  for (const e of edges) {
    const from = posMap.get(e.from)
    const to = posMap.get(e.to)
    if (!from || !to) continue

    const isUserEdge = !from.node.isDead !== !to.node.isDead
    if (isUserEdge) {
      const userId = from.node.isDead ? to.id : from.id
      const decId = from.node.isDead ? from.id : to.id
      if (userAnchor.get(userId) !== decId) continue
    }
    connections.push({
      fromX: from.x, fromY: from.y,
      toX: to.x, toY: to.y,
      fromId: e.from, toId: e.to,
      relationType: e.relationType,
      isSpouse: e.relationType === 'spouse',
      isUserEdge,
      confirmed: e.confirmed,
    })
  }

  // 计算 SVG 尺寸（基于所有节点实际位置）
  const allRight = positioned.map((p) => p.x + CARD_W / 2)
  const allBottom = positioned.map((p) => p.y + CARD_H / 2)
  const svgW = Math.max(...allRight) + PAD
  const svgH = Math.max(...allBottom) + PAD

  return { positioned, connections, svgW, svgH }
}

function formatLife(birth: string, death: string): string {
  if (!birth && !death) return ''
  const b = birth ? birth.substring(0, 4) : '?'
  const d = death ? death.substring(0, 4) : ''
  return `${b}-${d}`
}

/* ========== 组件 ========== */

const MemoConnection = memo(function MemoConnection({ c }: { c: ConnLine }) {
  return (
    <g>
      {c.isSpouse ? (
        <line
          x1={c.fromX} y1={c.fromY}
          x2={c.toX} y2={c.toY}
          stroke="#eb2f96" strokeWidth={2} strokeDasharray="6,3"
        />
      ) : c.isUserEdge ? (
        <>
          <line
            x1={c.fromX} y1={c.fromY}
            x2={c.toX} y2={c.toY}
            stroke="#52c41a" strokeWidth={1.5} strokeDasharray="4,4"
          />
          <text
            x={(c.fromX + c.toX) / 2}
            y={(c.fromY + c.toY) / 2 - 6}
            textAnchor="middle"
            fill="#52c41a"
            fontSize={11}
          >
            {labelOf(c.relationType)}
          </text>
        </>
      ) : (
        <>
          <line
            x1={c.fromX} y1={c.fromY + CARD_H / 2}
            x2={c.fromX} y2={c.toY - CARD_H / 2 - 10}
            stroke={c.confirmed ? '#1890ff' : '#ccc'} strokeWidth={1.5}
          />
          <line
            x1={c.fromX} y1={c.toY - CARD_H / 2 - 10}
            x2={c.toX} y2={c.toY - CARD_H / 2 - 10}
            stroke={c.confirmed ? '#1890ff' : '#ccc'} strokeWidth={1.5}
          />
          <polygon
            points={`${c.toX},${c.toY - CARD_H / 2 + 2} ${c.toX - 4},${c.toY - CARD_H / 2 - 6} ${c.toX + 4},${c.toY - CARD_H / 2 - 6}`}
            fill={c.confirmed ? '#1890ff' : '#ccc'}
          />
        </>
      )}
    </g>
  )
})

const MAX_NODES_WARN = 80

export default function RelationGraph({ graveId }: { graveId: number }) {
  const [graphData, setGraphData] = useState<RelationGraphResp | null>(null)
  const [loading, setLoading] = useState(true)
  const [scale, setScale] = useState(1)
  const [offsetX, setOffsetX] = useState(0)
  const [offsetY, setOffsetY] = useState(0)
  const touchRef = useRef<{ dist0: number; cx0: number; cy0: number; scale0: number; x0: number; y0: number } | null>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    fetchGraph()
  }, [graveId])

  const fetchGraph = async () => {
    setLoading(true)
    try {
      const res = await relationGraphApi.getRelationGraph(graveId)
      if (res.code === 0 && res.data) setGraphData(res.data)
    } catch {
      console.error('获取关系图失败')
    } finally {
      setLoading(false)
    }
  }

  const handleWheel = useCallback((e: WheelEvent) => {
    e.preventDefault()
    const delta = e.deltaY > 0 ? -0.1 : 0.1
    setScale((s) => Math.min(3, Math.max(0.3, s + delta)))
  }, [])

  const handleTouchStart = useCallback((e: React.TouchEvent) => {
    if (e.touches.length === 2) {
      const t0 = e.touches[0], t1 = e.touches[1]
      const dx = t1.clientX - t0.clientX, dy = t1.clientY - t0.clientY
      touchRef.current = {
        dist0: Math.sqrt(dx * dx + dy * dy),
        cx0: (t0.clientX + t1.clientX) / 2,
        cy0: (t0.clientY + t1.clientY) / 2,
        scale0: scale,
        x0: offsetX,
        y0: offsetY,
      }
    } else if (e.touches.length === 1) {
      touchRef.current = {
        dist0: 0,
        cx0: e.touches[0].clientX,
        cy0: e.touches[0].clientY,
        scale0: scale,
        x0: offsetX,
        y0: offsetY,
      }
    }
  }, [scale, offsetX, offsetY])

  const handleTouchMove = useCallback((e: React.TouchEvent) => {
    if (!touchRef.current) return
    if (e.touches.length === 2 && touchRef.current.dist0 > 0) {
      const t0 = e.touches[0], t1 = e.touches[1]
      const dx = t1.clientX - t0.clientX, dy = t1.clientY - t0.clientY
      const dist = Math.sqrt(dx * dx + dy * dy)
      const s = Math.min(3, Math.max(0.3, touchRef.current.scale0 * dist / touchRef.current.dist0))
      setScale(s)
    } else if (e.touches.length === 1) {
      setOffsetX(touchRef.current.x0 + e.touches[0].clientX - touchRef.current.cx0)
      setOffsetY(touchRef.current.y0 + e.touches[0].clientY - touchRef.current.cy0)
    }
  }, [])

  const handleTouchEnd = useCallback(() => {
    touchRef.current = null
  }, [])

  useEffect(() => {
    const el = containerRef.current
    if (el) {
      el.addEventListener('wheel', handleWheel, { passive: false })
      return () => el.removeEventListener('wheel', handleWheel)
    }
  }, [handleWheel])

  const layout = useMemo(() => {
    if (!graphData || graphData.nodes.length === 0) return null
    return computeLayout(graphData.nodes, graphData.edges)
  }, [graphData])

  if (loading) {
    return <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', padding: 80 }}><Spin size="large" /></div>
  }
  if (!graphData || graphData.nodes.length === 0) {
    return <div style={{ padding: 80 }}><Empty description="暂无关系数据" /></div>
  }
  if (!layout) return null

  return (
    <div style={{ padding: 20 }}>
      <Card title={`${graphData.graveName} - 关系图`}>
        {graphData.nodes.length >= MAX_NODES_WARN && (
          <Alert
            type="warning"
            showIcon
            message={`节点数量较多（${graphData.nodes.length}），渲染性能可能受影响`}
            style={{ marginBottom: 12 }}
          />
        )}
        <div
          className="rg-scroll"
          ref={containerRef}
          onTouchStart={handleTouchStart}
          onTouchMove={handleTouchMove}
          onTouchEnd={handleTouchEnd}
          style={{ touchAction: 'none', overflow: 'hidden' }}
        >
          <div style={{
            position: 'relative',
            width: layout.svgW,
            height: layout.svgH,
            transform: `scale(${scale}) translate(${offsetX / scale}px, ${offsetY / scale}px)`,
            transformOrigin: '0 0',
            willChange: 'transform',
          }}>
            {/* 连线层 */}
            <svg style={{ position: 'absolute', top: 0, left: 0, width: layout.svgW, height: layout.svgH, pointerEvents: 'none', zIndex: 1, willChange: 'transform' }}>
              {layout.connections.map((c, i) => (
                <MemoConnection key={i} c={c} />
              ))}
            </svg>

            {/* 节点卡片层 */}
            {layout.positioned.map((p) => {
              const isDead = p.node.isDead
              const borderColor = isDead ? '#999' : '#1890ff'
              const bg = isDead ? '#fafafa' : '#e6f7ff'
              return (
                <div
                  key={p.id}
                  className="rg-card"
                  style={{
                    left: p.x - CARD_W / 2,
                    top: p.y - CARD_H / 2,
                    width: CARD_W,
                    height: CARD_H,
                    borderColor,
                    backgroundColor: bg,
                    position: 'absolute',
                    zIndex: 2,
                  }}
                >
                  <div className="rg-card-name">
                    {!isDead && <span className="rg-user-icon">👤</span>}
                    <span>{p.node.name}</span>
                    {isDead && GENDER_ICON[p.node.gender] && (
                      <span style={{ marginLeft: 4, color: p.node.gender === 1 ? '#1890ff' : '#eb2f96', fontSize: 14 }}>
                        {GENDER_ICON[p.node.gender]}
                      </span>
                    )}
                  </div>
                  <div className="rg-card-life">
                    {isDead ? formatLife(p.node.birthDate || '', p.node.deathDate || '') : ''}
                  </div>
                </div>
              )
            })}
          </div>
        </div>

        {/* 图例 */}
        <div style={{ marginTop: 16, fontSize: 12, color: '#666' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 8 }}>
            <strong>图例说明：</strong>
            <span style={{ fontSize: 11, color: '#999' }}>
              缩放: {Math.round(scale * 100)}% | 滚轮缩放 · 双指缩放 · 拖动平移
            </span>
          </div>
          <div style={{ display: 'flex', gap: 18, flexWrap: 'wrap', marginTop: 6 }}>
            <span style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
              <span className="rg-legend-box" style={{ borderColor: '#999', backgroundColor: '#fafafa' }} />
              逝者节点
            </span>
            <span style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
              <span className="rg-legend-box" style={{ borderColor: '#1890ff', backgroundColor: '#e6f7ff' }} />
              用户节点
            </span>
            <span style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
              <span style={{ width: 36, height: 2, backgroundColor: '#1890ff', display: 'inline-block' }} />
              亲子
            </span>
            <span style={{ display: 'flex', alignItems: 'center', gap: 6 }}>
              <span style={{ width: 36, borderTop: '2px dashed #eb2f96', display: 'inline-block' }} />
              夫妻
            </span>
          </div>
        </div>
      </Card>
    </div>
  )
}
