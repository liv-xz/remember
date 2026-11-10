import { useEffect, useState } from 'react'
import { Table, Button, Tag, message, Popconfirm, Modal, InputNumber, Space, Tooltip } from 'antd'
import { PlusOutlined, StopOutlined, CopyOutlined } from '@ant-design/icons'
import * as shareLinkApi from '../../api/share-link'
import type { ShareLinkResp } from '../../types/share-link'

interface Props {
  graveId: number
  refreshKey: number
}

export default function ShareLinkManager({ graveId, refreshKey }: Props) {
  const [links, setLinks] = useState<ShareLinkResp[]>([])
  const [loading, setLoading] = useState(true)
  const [createOpen, setCreateOpen] = useState(false)
  const [expireDays, setExpireDays] = useState<number | null>(0)
  const [creating, setCreating] = useState(false)

  const fetchLinks = async () => {
    setLoading(true)
    try {
      const res = await shareLinkApi.getShareLinks(graveId)
      if (res.code === 0 && res.data) {
        setLinks(res.data.items)
      }
    } catch {
      message.error('获取分享链接失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchLinks() }, [graveId, refreshKey])

  const handleCreate = async () => {
    setCreating(true)
    try {
      const res = await shareLinkApi.createShareLink(graveId, { expireDays: expireDays ?? 0 })
      if (res.code === 0) {
        message.success('分享链接创建成功')
        setCreateOpen(false)
        setExpireDays(0)
        fetchLinks()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('创建失败')
    } finally {
      setCreating(false)
    }
  }

  const handleDeactivate = async (linkId: number) => {
    try {
      const res = await shareLinkApi.deactivateShareLink(graveId, linkId)
      if (res.code === 0) {
        message.success('已停用')
        fetchLinks()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('操作失败')
    }
  }

  const handleCopy = (shareCode: string) => {
    const url = `${window.location.origin}/join/${shareCode}`
    navigator.clipboard.writeText(url).then(
      () => message.success('链接已复制到剪贴板'),
      () => message.error('复制失败，请手动复制')
    )
  }

  const formatExpire = (expireAt?: string, isActive?: boolean) => {
    if (!isActive) return <Tag color="red">已停用</Tag>
    if (!expireAt) return <Tag color="green">永久有效</Tag>
    const expire = new Date(expireAt)
    const now = new Date()
    if (expire < now) return <Tag color="red">已过期</Tag>
    const days = Math.ceil((expire.getTime() - now.getTime()) / (1000 * 60 * 60 * 24))
    return <Tag color="blue">剩余 {days} 天</Tag>
  }

  const columns = [
    {
      title: '分享码', dataIndex: 'shareCode', key: 'shareCode',
      render: (v: string) => <code>{v.slice(0, 8)}...</code>,
    },
    {
      title: '有效期', key: 'expire',
      render: (_: unknown, r: ShareLinkResp) => formatExpire(r.expireAt, r.isActive),
    },
    {
      title: '使用次数', key: 'usage',
      render: (_: unknown, r: ShareLinkResp) => (
        <span>{r.usedCount}{r.maxUseCount > 0 ? ` / ${r.maxUseCount}` : ' / ∞'}</span>
      ),
    },
    {
      title: '状态', dataIndex: 'isActive', key: 'status',
      render: (v: boolean) => v ? <Tag color="green">有效</Tag> : <Tag color="red">已停用</Tag>,
    },
    {
      title: '创建时间', dataIndex: 'createdAt', key: 'createdAt',
    },
    {
      title: '操作', key: 'actions',
      render: (_: unknown, r: ShareLinkResp) => (
        <Space>
          <Tooltip title="复制链接">
            <Button type="link" icon={<CopyOutlined />} onClick={() => handleCopy(r.shareCode)}>
              复制
            </Button>
          </Tooltip>
          {r.isActive && (
            <Popconfirm title="确定停用此分享链接？停用后该链接将立即失效" onConfirm={() => handleDeactivate(r.id)}>
              <Button type="link" danger icon={<StopOutlined />}>停用</Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setCreateOpen(true)}>
          创建分享链接
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={links}
        rowKey="id"
        loading={loading}
        pagination={false}
        locale={{ emptyText: '暂无分享链接，点击上方按钮创建' }}
      />

      <Modal
        title="创建分享链接"
        open={createOpen}
        onCancel={() => { setCreateOpen(false); setExpireDays(0) }}
        onOk={handleCreate}
        confirmLoading={creating}
      >
        <div style={{ marginBottom: 8 }}>设置有效期（天）</div>
        <InputNumber
          style={{ width: '100%' }}
          min={0}
          max={365}
          value={expireDays}
          onChange={(v) => setExpireDays(v)}
          placeholder="0 表示永久有效"
          addonAfter="天"
        />
        <div style={{ marginTop: 8, color: '#888' }}>
          {expireDays === 0 || expireDays === null ? '永久有效，不会过期' : `${expireDays} 天后自动失效`}
        </div>
      </Modal>
    </div>
  )
}
