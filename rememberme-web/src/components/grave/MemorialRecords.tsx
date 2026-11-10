import { useEffect, useState } from 'react'
import { Button, Select, message, Popconfirm, Tag, Modal, Input, Timeline } from 'antd'
import { DeleteOutlined } from '@ant-design/icons'
import * as recordApi from '../../api/memorial-record'
import * as deceasedApi from '../../api/deceased'
import type { MemorialRecordResp } from '../../types/memorial-record'
import type { DeceasedResp } from '../../types/deceased'
import '../../styles/page-enhance.css'
import { RECORD_TYPES } from '../../constants/enums'

interface Props {
  graveId: number
  refreshKey: number
}

export default function MemorialRecords({ graveId, refreshKey }: Props) {
  const [records, setRecords] = useState<MemorialRecordResp[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [deceasedList, setDeceasedList] = useState<DeceasedResp[]>([])
  const [selectedDeceased, setSelectedDeceased] = useState<number | undefined>(undefined)

  const [msgOpen, setMsgOpen] = useState(false)
  const [msgText, setMsgText] = useState('')
  const [sending, setSending] = useState(false)

  const fetchRecords = async (p = 1) => {
    setLoading(true)
    try {
      const res = await recordApi.getRecordsByGrave(graveId, p, 10)
      if (res.code === 0) {
        const newRecords = res.data || []
        if (p === 1) {
          setRecords(newRecords)
        } else {
          setRecords((prev) => [...prev, ...newRecords])
        }
        setTotal(res.total || 0)
        setPage(res.page || 1)
      }
    } catch {
      message.error('获取祭祀记录失败')
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

  useEffect(() => { fetchRecords(1); fetchDeceased() }, [graveId, refreshKey])

  const handleQuickAction = async (type: number) => {
    setSending(true)
    try {
      const res = await recordApi.createRecord(graveId, { type, deceasedID: selectedDeceased })
      if (res.code === 0) {
        message.success(`${RECORD_TYPES[type].icon} ${RECORD_TYPES[type].label}成功`)
        fetchRecords(1)
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('操作失败')
    } finally {
      setSending(false)
    }
  }

  const handleSendMessage = async () => {
    if (!msgText.trim()) return
    setSending(true)
    try {
      const res = await recordApi.createRecord(graveId, {
        type: 5,
        deceasedID: selectedDeceased,
        message: msgText.trim(),
      })
      if (res.code === 0) {
        message.success('留言成功')
        setMsgOpen(false)
        setMsgText('')
        fetchRecords(1)
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('留言失败')
    } finally {
      setSending(false)
    }
  }

  const handleDelete = async (id: number) => {
    try {
      const res = await recordApi.deleteRecord(id)
      if (res.code === 0) {
        message.success('已删除')
        fetchRecords(1)
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('删除失败')
    }
  }

  const hasMore = records.length < total

  return (
    <div>
      <div style={{ marginBottom: 16, display: 'flex', alignItems: 'center', gap: 8, flexWrap: 'wrap' }}>
        <Select
          placeholder="全部逝者"
          allowClear
          style={{ width: 140 }}
          value={selectedDeceased}
          onChange={(v) => setSelectedDeceased(v)}
          options={deceasedList.map((d) => ({ label: d.name, value: d.id }))}
        />
        {[1, 2, 3, 4].map((type) => {
          const t = RECORD_TYPES[type]
          return (
            <Button key={type} loading={sending} onClick={() => handleQuickAction(type)}>
              {t.icon} {t.label}
            </Button>
          )
        })}
        <Button type="primary" onClick={() => setMsgOpen(true)}>
          💬 留言
        </Button>
      </div>

      <div className="record-timeline">
        {records.length === 0 && !loading ? (
          <div style={{ textAlign: 'center', padding: 40, color: '#aaa' }}>暂无祭祀记录</div>
        ) : (
          <Timeline
            items={records.map((r) => {
              const t = RECORD_TYPES[r.type]
              const deceasedName = r.deceasedID
                ? (deceasedList.find((x) => x.id === r.deceasedID)?.name || '-')
                : '墓园'
              return {
                dot: <span style={{ fontSize: 18 }}>{t?.icon}</span>,
                color: t?.color,
                children: (
                  <div className="memorial-record-card">
                    <div className="record-header">
                      <span>
                        <strong>{t?.label}</strong>
                        {r.deceasedID ? (
                          <Tag color="blue" style={{ marginLeft: 8 }}>{deceasedName}</Tag>
                        ) : (
                          <Tag style={{ marginLeft: 8 }}>墓园</Tag>
                        )}
                      </span>
                      <Popconfirm title="确定删除此记录？" onConfirm={() => handleDelete(r.id)}>
                        <Button type="link" danger icon={<DeleteOutlined />} size="small" />
                      </Popconfirm>
                    </div>
                    {r.message && (
                      <div className="record-message">{r.message}</div>
                    )}
                    <div className="record-time">{r.createdAt}</div>
                  </div>
                ),
              }
            })}
          />
        )}
        {hasMore && (
          <div style={{ textAlign: 'center', marginTop: 16 }}>
            <Button loading={loading} onClick={() => fetchRecords(page + 1)}>
              加载更多
            </Button>
          </div>
        )}
      </div>

      <Modal
        title="留言"
        open={msgOpen}
        onCancel={() => { setMsgOpen(false); setMsgText('') }}
        onOk={handleSendMessage}
        confirmLoading={sending}
      >
        <Input.TextArea
          rows={4}
          value={msgText}
          onChange={(e) => setMsgText(e.target.value)}
          placeholder="写下你想说的话..."
          maxLength={500}
        />
      </Modal>
    </div>
  )
}
