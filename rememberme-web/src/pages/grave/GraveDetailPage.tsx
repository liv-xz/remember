import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Typography, Tabs, Button, Card, Descriptions, Tag, Spin, message, Space } from 'antd'
import { ArrowLeftOutlined, EditOutlined, TeamOutlined, InfoCircleOutlined, UserOutlined, ShareAltOutlined, HistoryOutlined, CalendarOutlined, HeartOutlined, ApartmentOutlined, NodeIndexOutlined } from '@ant-design/icons'
import '../../styles/page-enhance.css'
import * as graveApi from '../../api/grave'
import type { GraveResp } from '../../types/grave'
import type { DeceasedResp } from '../../types/deceased'
import DeceasedList from '../../components/deceased/DeceasedList'
import GraveDetail from '../../components/grave/GraveDetail'
import GraveForm from '../../components/grave/GraveForm'
import RelationGraph from '../../components/relation-graph/RelationGraph'
import DeceasedRelationManager from '../../components/relation/DeceasedRelationManager'
import UserRelationManager from '../../components/relation/UserRelationManager'
import MemorialRelations from '../../components/grave/MemorialRelations'
import ShareLinkManager from '../../components/grave/ShareLinkManager'
import MemorialRecords from '../../components/grave/MemorialRecords'
import MemorialEvents from '../../components/grave/MemorialEvents'

const { Title } = Typography

export default function GraveDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const graveId = Number(id)

  const [grave, setGrave] = useState<GraveResp | null>(null)
  const [loading, setLoading] = useState(true)
  const [editOpen, setEditOpen] = useState(false)
  const [deceasedRefreshKey, setDeceasedRefreshKey] = useState(0)

  const fetchGrave = async () => {
    setLoading(true)
    try {
      const res = await graveApi.getGraveById(graveId)
      if (res.code === 0 && res.data) {
        setGrave(res.data)
      } else {
        message.error(res.message || '墓地不存在')
        navigate('/')
      }
    } catch {
      message.error('获取墓地详情失败')
      navigate('/')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { if (graveId) fetchGrave() }, [graveId])

  if (loading) {
    return <div style={{ textAlign: 'center', padding: 80 }}><Spin size="large" /></div>
  }

  if (!grave) return null

  const refreshTrigger = () => setDeceasedRefreshKey((k) => k + 1)

  const commonTabs = [
    {
      key: 'info',
      label: '墓地信息',
      icon: <InfoCircleOutlined />,
      children: <GraveDetail grave={grave} />,
    },
    {
      key: 'deceased',
      label: '逝者列表',
      icon: <UserOutlined />,
      children: (
        <DeceasedList
          graveId={graveId}
          graveType={grave.graveType}
          graveMode={grave.graveMode}
          refreshKey={deceasedRefreshKey}
          onRefresh={refreshTrigger}
        />
      ),
    },
    {
      key: 'share-links',
      label: '分享链接',
      icon: <ShareAltOutlined />,
      children: <ShareLinkManager graveId={graveId} refreshKey={deceasedRefreshKey} />,
    },
    {
      key: 'records',
      label: '祭祀记录',
      icon: <HistoryOutlined />,
      children: <MemorialRecords graveId={graveId} refreshKey={deceasedRefreshKey} />,
    },
    {
      key: 'events',
      label: '纪念事件',
      icon: <CalendarOutlined />,
      children: <MemorialEvents graveId={graveId} refreshKey={deceasedRefreshKey} />,
    },
  ]

  const memorialTabs = [
    {
      key: 'memorial-relations',
      label: '纪念关系',
      icon: <HeartOutlined />,
      children: <MemorialRelations graveId={graveId} refreshKey={deceasedRefreshKey} />,
    },
  ]

  const familyTreeTabs = [
    {
      key: 'relations',
      label: '逝者间关系',
      icon: <ApartmentOutlined />,
      children: <DeceasedRelationManager graveId={graveId} refreshKey={deceasedRefreshKey} />,
    },
    {
      key: 'relation-graph',
      label: '关系图',
      icon: <NodeIndexOutlined />,
      children: <RelationGraph graveId={graveId} />,
    },
    {
      key: 'user-relations',
      label: '认亲管理',
      icon: <TeamOutlined />,
      children: <UserRelationManager graveId={graveId} refreshKey={deceasedRefreshKey} />,
    },
  ]

  const tabItems = grave.graveMode === 2
    ? [...commonTabs, ...familyTreeTabs]
    : [...commonTabs, ...memorialTabs]

  return (
    <div>
      <Card style={{ marginBottom: 16 }}>
        <div className="grave-detail-header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Space>
            <Button icon={<ArrowLeftOutlined />} onClick={() => navigate('/')}>返回</Button>
            <Title level={3} style={{ margin: 0 }}>{grave.name}</Title>
          </Space>
          <Space>
            <Button icon={<TeamOutlined />} onClick={() => navigate(`/grave/${graveId}/members`)}>
              成员管理
            </Button>
            <Button type="primary" icon={<EditOutlined />} onClick={() => setEditOpen(true)}>
              编辑
            </Button>
          </Space>
        </div>
      </Card>

      <Tabs items={tabItems} />

      <GraveForm
        open={editOpen}
        editingGrave={grave}
        onCancel={() => setEditOpen(false)}
        onSuccess={() => { setEditOpen(false); fetchGrave() }}
      />
    </div>
  )
}
