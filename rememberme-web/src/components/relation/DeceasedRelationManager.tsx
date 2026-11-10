import { useEffect, useState } from 'react'
import { Table, Button, Select, Form, Tag, message, Popconfirm, Space, Card, Empty } from 'antd'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
import * as relationApi from '../../api/relation'
import * as deceasedApi from '../../api/deceased'
import type { DeceasedRelationResp } from '../../types/relation'
import type { DeceasedResp } from '../../types/deceased'
import { DECEASED_RELATION_LABELS, DIR_LABELS } from '../../constants/enums'

interface Props {
  graveId: number
  refreshKey: number
}

export default function DeceasedRelationManager({ graveId, refreshKey }: Props) {
  const [relations, setRelations] = useState<DeceasedRelationResp[]>([])
  const [deceasedList, setDeceasedList] = useState<DeceasedResp[]>([])
  const [loading, setLoading] = useState(true)
  const [formShow, setFormShow] = useState(false)
  const [aId, setAId] = useState<number | null>(null)
  const [bId, setBId] = useState<number | null>(null)
  const [relType, setRelType] = useState<number>(1)
  const [dir, setDir] = useState<number>(0)

  const fetchData = async () => {
    setLoading(true)
    try {
      const [relRes, decRes] = await Promise.all([
        relationApi.getDeceasedRelations(graveId),
        deceasedApi.getDeceasedByGrave(graveId),
      ])
      if (relRes.code === 0 && relRes.data) setRelations(relRes.data.items)
      if (decRes.code === 0 && decRes.data) setDeceasedList(decRes.data.items)
    } catch {
      message.error('获取数据失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [graveId, refreshKey])

  const deceasedOptions = deceasedList.map((d) => ({
    value: d.id,
    label: `${d.name}${d.gender === 1 ? '(男)' : d.gender === 2 ? '(女)' : ''}`,
  }))

  const deceasedMap = Object.fromEntries(deceasedList.map((d) => [d.id, d.name]))

  const handleCreate = async () => {
    if (!aId || !bId || aId === bId) {
      message.warning('请选择两位不同的逝者')
      return
    }
    try {
      const res = await relationApi.createDeceasedRelation(graveId, {
        deceasedID1: aId,
        deceasedID2: bId,
        relationType: relType,
        direction: dir,
      })
      if (res.code === 0) {
        message.success('关系创建成功')
        setFormShow(false)
        setAId(null)
        setBId(null)
        fetchData()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('创建失败')
    }
  }

  const handleDelete = async (id: number) => {
    try {
      const res = await relationApi.deleteDeceasedRelation(id)
      if (res.code === 0) {
        message.success('删除成功')
        fetchData()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('删除失败')
    }
  }

  const columns = [
    {
      title: '逝者 A', key: 'a',
      render: (_: unknown, r: DeceasedRelationResp) => deceasedMap[r.deceasedID1] || `ID:${r.deceasedID1}`,
    },
    {
      title: '逝者 B', key: 'b',
      render: (_: unknown, r: DeceasedRelationResp) => deceasedMap[r.deceasedID2] || `ID:${r.deceasedID2}`,
    },
    {
      title: '关系', dataIndex: 'relationType', key: 'relationType',
      render: (v: number) => <Tag color={v === 1 ? 'pink' : 'blue'}>{DECEASED_RELATION_LABELS[v] || v}</Tag>,
    },
    {
      title: '方向', dataIndex: 'direction', key: 'direction',
      render: (v: number, r: DeceasedRelationResp) => {
        if (r.relationType === 1) return <Tag>无向</Tag>
        if (v === 1) return <Tag color="orange">{deceasedMap[r.deceasedID1]} 是 {deceasedMap[r.deceasedID2]} 的长辈</Tag>
        if (v === 2) return <Tag color="orange">{deceasedMap[r.deceasedID2]} 是 {deceasedMap[r.deceasedID1]} 的长辈</Tag>
        return <Tag>{DIR_LABELS[v]}</Tag>
      },
    },
    {
      title: '操作', key: 'actions',
      render: (_: unknown, r: DeceasedRelationResp) => (
        <Popconfirm title="确定删除此关系？" onConfirm={() => handleDelete(r.id)}>
          <Button type="link" danger icon={<DeleteOutlined />}>删除</Button>
        </Popconfirm>
      ),
    },
  ]

  return (
    <Card
      title="逝者间关系"
      extra={
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setFormShow(true)}>
          添加关系
        </Button>
      }
    >
      {deceasedList.length < 2 ? (
        <Empty description="至少需要两位逝者才能建立关系" />
      ) : (
        <>
          <Table
            columns={columns}
            dataSource={relations}
            rowKey="id"
            loading={loading}
            pagination={false}
            locale={{ emptyText: '暂无关系，点击"添加关系"开始建立家族树' }}
            style={{ marginBottom: 16 }}
          />

          {formShow && (
            <Card size="small" title="添加关系" style={{ marginTop: 16 }} bordered>
              <Space wrap>
                <Select
                  style={{ width: 180 }}
                  placeholder="选择逝者 A"
                  value={aId}
                  onChange={(v) => setAId(v)}
                  options={deceasedOptions}
                />
                <Select
                  style={{ width: 140 }}
                  value={relType}
                  onChange={(v) => { setRelType(v); setDir(0) }}
                  options={[
                    { value: 1, label: '夫妻' },
                    { value: 2, label: '亲子' },
                  ]}
                />
                <Select
                  style={{ width: 180 }}
                  placeholder="选择逝者 B"
                  value={bId}
                  onChange={(v) => setBId(v)}
                  options={deceasedOptions}
                />
                {relType === 2 && (
                  <Select
                    style={{ width: 160 }}
                    value={dir}
                    onChange={(v) => setDir(v)}
                    options={[
                      { value: 1, label: 'A 是 B 的父/母' },
                      { value: 2, label: 'B 是 A 的父/母' },
                    ]}
                    placeholder="选择方向"
                  />
                )}
                <Button type="primary" onClick={handleCreate}>确定添加</Button>
                <Button onClick={() => setFormShow(false)}>取消</Button>
              </Space>
            </Card>
          )}
        </>
      )}
    </Card>
  )
}
