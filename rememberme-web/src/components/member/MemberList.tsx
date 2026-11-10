import { useState } from 'react'
import { Table, Button, Tag, message, Popconfirm, Space, Modal, Typography, Input } from 'antd'
import { UserOutlined, SwapOutlined } from '@ant-design/icons'
import * as memberApi from '../../api/member'
import type { GraveMemberResp } from '../../types/member'
import TransferModal from './TransferModal'
import { ROLE_MAP } from '../../constants/enums'

const { Text } = Typography
interface Props {
  graveId: number
  members: GraveMemberResp[]
  onRefresh: () => void
}

export default function MemberList({ graveId, members, onRefresh }: Props) {
  const [transferOpen, setTransferOpen] = useState(false)
  const [descModalOpen, setDescModalOpen] = useState(false)
  const [editingMember, setEditingMember] = useState<GraveMemberResp | null>(null)
  const [newDesc, setNewDesc] = useState('')

  const handleRoleToggle = async (member: GraveMemberResp) => {
    const newRole = member.role === 20 ? 30 : 20
    const action = newRole === 20 ? '升级为管理员' : '降级为成员'
    try {
      const res = await memberApi.updateMemberRole(member.id, { role: newRole })
      if (res.code === 0) {
        message.success(`已${action}`)
        onRefresh()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('操作失败')
    }
  }

  const handleRemove = async (memberId: number) => {
    try {
      const res = await memberApi.removeMember(memberId)
      if (res.code === 0) {
        message.success('已移除成员')
        onRefresh()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('移除失败')
    }
  }

  const handleUpdateDesc = async () => {
    if (!editingMember) return
    try {
      const res = await memberApi.updateMemberRelationDesc(editingMember.id, { relationDesc: newDesc })
      if (res.code === 0) {
        message.success('已更新关系描述')
        setDescModalOpen(false)
        onRefresh()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('更新失败')
    }
  }

  const columns = [
    {
      title: '用户', key: 'user',
      render: (_: unknown, r: GraveMemberResp) => (
        <Space>
          <UserOutlined />
          <span>{r.nickname || r.username}</span>
        </Space>
      ),
    },
    {
      title: '角色', dataIndex: 'role', key: 'role',
      render: (v: number) => {
        const info = ROLE_MAP[v] || { label: '未知', color: 'default' }
        return <Tag color={info.color}>{info.label}</Tag>
      },
    },
    {
      title: '关系', dataIndex: 'relationDesc', key: 'relationDesc',
      render: (v: string, r: GraveMemberResp) => (
        <Button
          type="link"
          onClick={() => { setEditingMember(r); setNewDesc(r.relationDesc); setDescModalOpen(true) }}
        >
          {v || '点击设置'}
        </Button>
      ),
    },
    {
      title: '加入时间', dataIndex: 'createdAt', key: 'createdAt',
    },
    {
      title: '操作', key: 'actions',
      render: (_: unknown, record: GraveMemberResp) => (
        <Space>
          {record.role !== 10 && (
            <>
              <Button type="link" onClick={() => handleRoleToggle(record)}>
                {record.role === 20 ? '降为成员' : '升为管理员'}
              </Button>
              <Popconfirm title="确定移除该成员？" onConfirm={() => handleRemove(record.id)}>
                <Button type="link" danger>移除</Button>
              </Popconfirm>
            </>
          )}
          {record.role === 10 && (
            <Button
              type="link"
              icon={<SwapOutlined />}
              onClick={() => setTransferOpen(true)}
            >
              转让所有权
            </Button>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div>
      <Table columns={columns} dataSource={members} rowKey="id" pagination={false} />

      <Modal
        title="修改关系描述"
        open={descModalOpen}
        onOk={handleUpdateDesc}
        onCancel={() => { setDescModalOpen(false); setEditingMember(null) }}
      >
        <Text>用户：{editingMember?.nickname || editingMember?.username}</Text>
        <Input.TextArea
          style={{ marginTop: 12 }}
          value={newDesc}
          onChange={(e) => setNewDesc(e.target.value)}
          placeholder="与墓主的关系，如：侄孙"
          maxLength={50}
          rows={2}
        />
      </Modal>

      <TransferModal
        open={transferOpen}
        members={members}
        onCancel={() => setTransferOpen(false)}
        onSuccess={() => { setTransferOpen(false); onRefresh() }}
      />
    </div>
  )
}
