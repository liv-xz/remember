import { useEffect, useState } from 'react'
import { Table, Button, Select, DatePicker, Input, Switch, message, Popconfirm, Modal, Form, Space, Tag } from 'antd'
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons'
import * as eventApi from '../../api/memorial-event'
import * as deceasedApi from '../../api/deceased'
import type { MemorialEventResp } from '../../types/memorial-event'
import type { DeceasedResp } from '../../types/deceased'
import dayjs from 'dayjs'

const EVENT_TYPES: Record<number, { label: string; icon: string }> = {
  1: { label: '诞辰', icon: '🎂' },
  2: { label: '忌日', icon: '🕯️' },
  3: { label: '合葬日', icon: '💑' },
  4: { label: '自定义', icon: '📅' },
}

interface Props {
  graveId: number
  refreshKey: number
}

export default function MemorialEvents({ graveId, refreshKey }: Props) {
  const [events, setEvents] = useState<MemorialEventResp[]>([])
  const [loading, setLoading] = useState(true)
  const [deceasedList, setDeceasedList] = useState<DeceasedResp[]>([])

  const [modalOpen, setModalOpen] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [saving, setSaving] = useState(false)

  const [formDeceasedID, setFormDeceasedID] = useState<number | undefined>(undefined)
  const [formEventType, setFormEventType] = useState<number>(1)
  const [formEventDate, setFormEventDate] = useState<string>('')
  const [formTitle, setFormTitle] = useState('')
  const [formLunar, setFormLunar] = useState(0)

  const fetchEvents = async () => {
    setLoading(true)
    try {
      const res = await eventApi.getEventsByGrave(graveId)
      if (res.code === 0 && res.data) {
        setEvents(res.data.items)
      }
    } catch {
      message.error('获取纪念事件失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchDeceased = async () => {
    try {
      const res = await deceasedApi.getDeceasedByGrave(graveId)
      if (res.code === 0 && res.data) {
        setDeceasedList(res.data.items)
      }
    } catch { /* ignore */ }
  }

  useEffect(() => { fetchEvents(); fetchDeceased() }, [graveId, refreshKey])

  const resetForm = () => {
    setFormDeceasedID(undefined)
    setFormEventType(1)
    setFormEventDate('')
    setFormTitle('')
    setFormLunar(0)
    setEditingId(null)
  }

  const openCreate = () => {
    resetForm()
    setModalOpen(true)
  }

  const openEdit = (e: MemorialEventResp) => {
    setEditingId(e.id)
    setFormDeceasedID(e.deceasedID)
    setFormEventType(e.eventType)
    setFormEventDate(e.eventDate)
    setFormTitle(e.title)
    setFormLunar(e.lunar)
    setModalOpen(true)
  }

  const handleSave = async () => {
    if (!formDeceasedID || !formEventDate.trim() || !formTitle.trim()) {
      message.warning('请填写必填项')
      return
    }
    setSaving(true)
    try {
      if (editingId) {
        const res = await eventApi.updateEvent(editingId, {
          deceasedID: formDeceasedID,
          eventType: formEventType,
          eventDate: formEventDate,
          title: formTitle.trim(),
          lunar: formLunar,
        })
        if (res.code === 0) {
          message.success('更新成功')
          setModalOpen(false)
          resetForm()
          fetchEvents()
        } else {
          message.error(res.message)
        }
      } else {
        const res = await eventApi.createEvent(graveId, {
          deceasedID: formDeceasedID,
          eventType: formEventType,
          eventDate: formEventDate,
          title: formTitle.trim(),
          lunar: formLunar,
        })
        if (res.code === 0) {
          message.success('创建成功')
          setModalOpen(false)
          resetForm()
          fetchEvents()
        } else {
          message.error(res.message)
        }
      }
    } catch {
      message.error('操作失败')
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async (id: number) => {
    try {
      const res = await eventApi.deleteEvent(id)
      if (res.code === 0) {
        message.success('已删除')
        fetchEvents()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('删除失败')
    }
  }

  const columns = [
    {
      title: '逝者', key: 'deceased',
      render: (_: unknown, r: MemorialEventResp) => {
        const d = deceasedList.find((x) => x.id === r.deceasedID)
        return d?.name || '-'
      },
    },
    {
      title: '类型', dataIndex: 'eventType', key: 'eventType', width: 90,
      render: (v: number) => {
        const t = EVENT_TYPES[v]
        return <span>{t?.icon} {t?.label}</span>
      },
    },
    {
      title: '日期', dataIndex: 'eventDate', key: 'eventDate', width: 120,
    },
    {
      title: '标题', dataIndex: 'title', key: 'title', ellipsis: true,
    },
    {
      title: '历法', dataIndex: 'lunar', key: 'lunar', width: 60,
      render: (v: number) => v === 1 ? <Tag color="purple">农历</Tag> : <Tag>公历</Tag>,
    },
    {
      title: '操作', key: 'actions', width: 120,
      render: (_: unknown, r: MemorialEventResp) => (
        <Space>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => openEdit(r)} />
          <Popconfirm title="确定删除此事件？" onConfirm={() => handleDelete(r.id)}>
            <Button type="link" danger size="small" icon={<DeleteOutlined />} />
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          创建纪念事件
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={events}
        rowKey="id"
        loading={loading}
        pagination={false}
        locale={{ emptyText: '暂无纪念事件' }}
      />

      <Modal
        title={editingId ? '编辑纪念事件' : '创建纪念事件'}
        open={modalOpen}
        onCancel={() => { setModalOpen(false); resetForm() }}
        onOk={handleSave}
        confirmLoading={saving}
        destroyOnClose
      >
        <Form layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item label="逝者" required>
            <Select
              placeholder="选择逝者"
              value={formDeceasedID}
              onChange={(v) => setFormDeceasedID(v)}
              options={deceasedList.map((d) => ({ label: d.name, value: d.id }))}
            />
          </Form.Item>

          <Form.Item label="事件类型" required>
            <Select
              value={formEventType}
              onChange={(v) => setFormEventType(v)}
              options={Object.entries(EVENT_TYPES).map(([k, v]) => ({
                label: `${v.icon} ${v.label}`,
                value: Number(k),
              }))}
            />
          </Form.Item>

          <Form.Item label="日期" required>
            <DatePicker
              style={{ width: '100%' }}
              value={formEventDate ? dayjs(formEventDate) : null}
              onChange={(d) => setFormEventDate(d ? d.format('YYYY-MM-DD') : '')}
            />
          </Form.Item>

          <Form.Item label="标题" required>
            <Input
              value={formTitle}
              onChange={(e) => setFormTitle(e.target.value)}
              placeholder="如：张三诞辰纪念"
              maxLength={100}
            />
          </Form.Item>

          <Form.Item label="农历">
            <Switch
              checked={formLunar === 1}
              onChange={(v) => setFormLunar(v ? 1 : 0)}
              checkedChildren="农历"
              unCheckedChildren="公历"
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
