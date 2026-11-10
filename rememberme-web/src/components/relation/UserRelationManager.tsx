import { useEffect, useState } from 'react'
import { Table, Tabs, Button, Tag, message, Popconfirm, Space, Empty } from 'antd'
import { CheckOutlined, CloseOutlined } from '@ant-design/icons'
import * as kinshipApi from '../../api/kinship'
import type { KinshipRelationResp } from '../../types/kinship'
import { RELATION_LABELS, RELATION_STATUS_MAP } from '../../constants/enums'

interface Props {
  graveId: number
  refreshKey: number
}

export default function UserRelationManager({ graveId, refreshKey }: Props) {
  const [allRelations, setAllRelations] = useState<KinshipRelationResp[]>([])
  const [loading, setLoading] = useState(true)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await kinshipApi.getKinshipRelations(graveId)
      if (res.code === 0 && res.data) {
        setAllRelations(res.data.items)
      }
    } catch {
      message.error('获取认亲数据失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchData() }, [graveId, refreshKey])

  const handleConfirm = async (id: number) => {
    try {
      const res = await kinshipApi.confirmKinshipRelation(id)
      if (res.code === 0) {
        message.success('已确认')
        fetchData()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('操作失败')
    }
  }

  const handleReject = async (id: number) => {
    try {
      const res = await kinshipApi.rejectKinshipRelation(id)
      if (res.code === 0) {
        message.success('已拒绝')
        fetchData()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('操作失败')
    }
  }

  const columns = (showActions: boolean) => [
    {
      title: '用户', key: 'user', dataIndex: 'userName',
      render: (v: string) => v || '-',
    },
    {
      title: '逝者', key: 'deceased', dataIndex: 'deceasedName',
    },
    {
      title: '关系', dataIndex: 'relationType', key: 'relationType',
      render: (v: number) => <Tag>{RELATION_LABELS[v] || `未知(${v})`}</Tag>,
    },
    {
      title: '状态', dataIndex: 'confirmationStatus', key: 'status',
      render: (v: number) => {
        const s = RELATION_STATUS_MAP[v] || { label: '未知', color: 'default' }
        return <Tag color={s.color}>{s.label}</Tag>
      },
    },
    { title: '备注', dataIndex: 'note', key: 'note', ellipsis: true },
    ...(showActions ? [{
      title: '操作', key: 'actions',
      render: (_: unknown, r: KinshipRelationResp) => (
        <Space>
          <Popconfirm title="确认通过此认亲关系？" onConfirm={() => handleConfirm(r.id)}>
            <Button type="link" icon={<CheckOutlined />} style={{ color: '#52c41a' }}>确认</Button>
          </Popconfirm>
          <Popconfirm title="确定拒绝此认亲？" onConfirm={() => handleReject(r.id)}>
            <Button type="link" danger icon={<CloseOutlined />}>拒绝</Button>
          </Popconfirm>
        </Space>
      ),
    }] : []),
  ]

  const pending = allRelations.filter((r) => r.confirmationStatus === 0)
  const confirmed = allRelations.filter((r) => r.confirmationStatus === 1)
  const rejected = allRelations.filter((r) => r.confirmationStatus === 2)

  const tabItems = [
    {
      key: 'pending',
      label: `待确认 (${pending.length})`,
      children: (
        <Table
          columns={columns(true)}
          dataSource={pending}
          rowKey="id"
          loading={loading}
          pagination={false}
          locale={{ emptyText: <Empty description="暂无待确认的认亲" /> }}
        />
      ),
    },
    {
      key: 'confirmed',
      label: `已确认 (${confirmed.length})`,
      children: (
        <Table
          columns={columns(false)}
          dataSource={confirmed}
          rowKey="id"
          loading={loading}
          pagination={false}
          locale={{ emptyText: <Empty description="暂无已确认的认亲" /> }}
        />
      ),
    },
    {
      key: 'rejected',
      label: `已拒绝 (${rejected.length})`,
      children: (
        <Table
          columns={columns(false)}
          dataSource={rejected}
          rowKey="id"
          loading={loading}
          pagination={false}
          locale={{ emptyText: <Empty description="暂无已拒绝的认亲" /> }}
        />
      ),
    },
  ]

  return <Tabs items={tabItems} />
}
