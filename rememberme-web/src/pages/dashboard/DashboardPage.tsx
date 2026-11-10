import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Card, Button, Row, Col, Typography, Tag, message, Popconfirm, Space } from 'antd'
import { PlusOutlined, TeamOutlined, DeleteOutlined } from '@ant-design/icons'
import * as graveApi from '../../api/grave'
import type { GraveResp } from '../../types/grave'
import GraveForm from '../../components/grave/GraveForm'
import { GRAVE_TYPE_MAP, LAYOUT_MAP } from '../../constants/enums'
import { useAuth } from '../../contexts/AuthContext'
import '../../styles/page-enhance.css'

const { Title, Text, Paragraph } = Typography

export default function DashboardPage() {
  const [graves, setGraves] = useState<GraveResp[]>([])
  const [loading, setLoading] = useState(true)
  const [formOpen, setFormOpen] = useState(false)
  const navigate = useNavigate()
  const { user } = useAuth()

  const fetchGraves = async () => {
    setLoading(true)
    try {
      const res = await graveApi.getMyGraves(1, 100)
      if (res.code === 0 && res.data) {
        setGraves(res.data)
      }
    } catch {
      message.error('获取墓地列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchGraves() }, [])

  const handleDelete = async (id: number) => {
    try {
      const res = await graveApi.deleteGrave(id)
      if (res.code === 0) {
        message.success('删除成功')
        fetchGraves()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('删除失败')
    }
  }

  const handleCreateSuccess = (graveId?: number) => {
    setFormOpen(false)
    if (graveId) {
      navigate(`/grave/${graveId}`)
    } else {
      fetchGraves()
    }
  }

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title level={3} style={{ margin: 0 }}>我的墓园</Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setFormOpen(true)}>
          创建墓园
        </Button>
      </div>

      {graves.length === 0 && !loading ? (
        <Card>
          <div style={{ textAlign: 'center', padding: 48, color: '#999' }}>
            <Title level={4} type="secondary">还没有墓园</Title>
            <Paragraph type="secondary">点击"创建墓园"开始您的数字纪念馆</Paragraph>
          </div>
        </Card>
      ) : (
        <Row gutter={[16, 16]}>
          {graves.map((grave) => (
            <Col xs={24} sm={12} lg={8} key={grave.id}>
              <Card
                hoverable
                className="grave-card"
                cover={
                  <div className="grave-card-cover" onClick={() => navigate(`/grave/${grave.id}`)}>
                    {grave.coverImage ? (
                      <img src={grave.coverImage} alt={grave.name} />
                    ) : (
                      <div className="grave-card-placeholder" style={{
                        background: 'linear-gradient(135deg, #1a5632 0%, #2d8a4e 50%, #61b87a 100%)',
                      }}>
                        <span style={{ fontSize: 36, opacity: 0.5 }}>🪦</span>
                      </div>
                    )}
                    <span className="grave-card-title">{grave.name}</span>
                  </div>
                }
                actions={[
                  <TeamOutlined key="members" onClick={(e) => {
                    e.stopPropagation()
                    navigate(`/grave/${grave.id}/members`)
                  }} />,
                  ...(user && grave.createdBy === user.id ? [
                    <Popconfirm
                      key="delete"
                      title="确定删除此墓园？"
                      onConfirm={() => handleDelete(grave.id)}
                      onCancel={(e) => e?.stopPropagation()}
                    >
                      <DeleteOutlined onClick={(e) => {
                        e.stopPropagation()
                      }} />
                    </Popconfirm>,
                  ] : []),
                ]}
                extra={
                  <Space size={4}>
                    {user && grave.createdBy === user.id ? <Tag color="blue">我创建的</Tag> : <Tag color="orange">我加入的</Tag>}
                    {grave.status === 1 ? <Tag color="green">启用</Tag> : <Tag color="red">停用</Tag>}
                  </Space>
                }
                loading={loading}
              >
                <div className="grave-card-body">
                  <Paragraph type="secondary" ellipsis={{ rows: 2 }}>
                    {grave.description || '暂无描述'}
                  </Paragraph>
                  <div style={{ display: 'flex', gap: 8 }}>
                    <Tag>{GRAVE_TYPE_MAP[grave.graveType] || '未知'}</Tag>
                    <Tag>{LAYOUT_MAP[grave.layoutStyle] || '未知'}</Tag>
                  </div>
                  <Text type="secondary" style={{ fontSize: 12 }}>创建者：{grave.creatorName}</Text>
                </div>
              </Card>
            </Col>
          ))}
        </Row>
      )}

      <GraveForm
        open={formOpen}
        editingGrave={null}
        onCancel={() => setFormOpen(false)}
        onSuccess={handleCreateSuccess}
      />
    </div>
  )
}
