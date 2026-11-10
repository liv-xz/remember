import { useEffect, useState, useCallback } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Typography, Button, message, Space, Spin } from 'antd'
import { ArrowLeftOutlined } from '@ant-design/icons'
import * as memberApi from '../../api/member'
import type { GraveMemberResp } from '../../types/member'
import MemberList from '../../components/member/MemberList'
import MemberForm from '../../components/member/MemberForm'

const { Title } = Typography

export default function MemberPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const graveId = Number(id)

  const [members, setMembers] = useState<GraveMemberResp[]>([])
  const [loading, setLoading] = useState(true)
  const [addOpen, setAddOpen] = useState(false)

  const fetchMembers = useCallback(async () => {
    setLoading(true)
    try {
      const res = await memberApi.getMembersByGrave(graveId)
      if (res.code === 0 && res.data) {
        setMembers(res.data.items)
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('获取成员列表失败')
    } finally {
      setLoading(false)
    }
  }, [graveId])

  useEffect(() => { if (graveId) fetchMembers() }, [graveId, fetchMembers])

  if (loading) {
    return <div style={{ textAlign: 'center', padding: 80 }}><Spin size="large" /></div>
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Space>
          <Button icon={<ArrowLeftOutlined />} onClick={() => navigate(`/grave/${graveId}`)}>返回墓地</Button>
          <Title level={3} style={{ margin: 0 }}>成员管理</Title>
        </Space>
        <Button type="primary" onClick={() => setAddOpen(true)}>添加成员</Button>
      </div>

      <MemberList
        graveId={graveId}
        members={members}
        onRefresh={fetchMembers}
      />

      <MemberForm
        open={addOpen}
        graveId={graveId}
        onCancel={() => setAddOpen(false)}
        onSuccess={() => { setAddOpen(false); fetchMembers() }}
      />
    </div>
  )
}
