import { useEffect, useState } from 'react'
import { Card, Button, Tag, message, Popconfirm, Row, Col, Avatar } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined, UserOutlined } from '@ant-design/icons'
import * as deceasedApi from '../../api/deceased'
import type { DeceasedResp } from '../../types/deceased'
import DeceasedForm from './DeceasedForm'
import { GENDER_MAP } from '../../constants/enums'
import '../../styles/page-enhance.css'

interface Props {
  graveId: number
  graveType: number
  graveMode: 1 | 2
  refreshKey: number
  onRefresh: () => void
}

export default function DeceasedList({ graveId, graveType, graveMode, refreshKey, onRefresh }: Props) {
  const [deceasedList, setDeceasedList] = useState<DeceasedResp[]>([])
  const [loading, setLoading] = useState(true)
  const [formOpen, setFormOpen] = useState(false)
  const [editingDeceased, setEditingDeceased] = useState<DeceasedResp | null>(null)

  const fetchList = async () => {
    setLoading(true)
    try {
      const res = await deceasedApi.getDeceasedByGrave(graveId)
      if (res.code === 0 && res.data) {
        setDeceasedList(res.data.items)
      }
    } catch {
      message.error('获取逝者列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchList() }, [graveId, refreshKey])

  const handleDelete = async (id: number) => {
    try {
      const res = await deceasedApi.deleteDeceased(id)
      if (res.code === 0) {
        message.success('删除成功')
        fetchList()
        onRefresh()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('删除失败')
    }
  }

  const capacityMsg = graveType === 1 ? '1位' : graveType === 2 ? '2位' : '不限'

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <span>逝者数量：{deceasedList.length}，容量上限：{capacityMsg}</span>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => { setEditingDeceased(null); setFormOpen(true) }}
          disabled={graveType === 1 && deceasedList.length >= 1 || graveType === 2 && deceasedList.length >= 2}
        >
          添加逝者
        </Button>
      </div>

      <Row gutter={[16, 16]}>
        {deceasedList.map((record) => (
          <Col xs={24} sm={12} lg={8} key={record.id}>
            <Card className="deceased-card" hoverable>
              <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 12 }}>
                <Avatar
                  className="deceased-avatar"
                  src={record.portraitUrl || undefined}
                  icon={!record.portraitUrl ? <UserOutlined /> : undefined}
                  size={60}
                  style={!record.portraitUrl ? { backgroundColor: '#d9d9d9', color: '#999' } : undefined}
                />
                <div>
                  <div style={{ fontSize: 16, fontWeight: 600 }}>{record.name}</div>
                  <Tag>{GENDER_MAP[record.gender] || '未知'}</Tag>
                </div>
              </div>
              {(record.birthDate || record.deathDate) && (
                <div style={{ color: '#888', fontSize: 13, marginBottom: 8 }}>
                  {record.birthDate ? record.birthDate.substring(0, 4) : '?'} — {record.deathDate ? record.deathDate.substring(0, 4) : '?'}
                </div>
              )}
              {record.epitaph && (
                <div className="epitaph" style={{
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  display: '-webkit-box',
                  WebkitLineClamp: 2,
                  WebkitBoxOrient: 'vertical',
                }}>
                  "{record.epitaph}"
                </div>
              )}
              <div style={{ marginTop: 12, display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
                <Button
                  size="small"
                  icon={<EditOutlined />}
                  onClick={() => { setEditingDeceased(record); setFormOpen(true) }}
                >
                  编辑
                </Button>
                <Popconfirm title="确定删除此逝者？" onConfirm={() => handleDelete(record.id)}>
                  <Button size="small" danger icon={<DeleteOutlined />}>删除</Button>
                </Popconfirm>
              </div>
            </Card>
          </Col>
        ))}
      </Row>

      <DeceasedForm
        open={formOpen}
        graveId={graveId}
        graveType={graveType}
        graveMode={graveMode}
        deceasedCount={deceasedList.length}
        editingDeceased={editingDeceased}
        onCancel={() => { setFormOpen(false); setEditingDeceased(null) }}
        onSuccess={() => {
          setFormOpen(false)
          setEditingDeceased(null)
          fetchList()
          onRefresh()
        }}
      />
    </div>
  )
}
