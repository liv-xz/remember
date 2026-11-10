import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Button, Typography, Result, Spin, Space, Modal, Select, Input, Form } from 'antd'
import { ArrowRightOutlined } from '@ant-design/icons'
import { useAuth } from '../../contexts/AuthContext'
import * as shareLinkApi from '../../api/share-link'
import type { ShareLinkValidResp, DeceasedBrief } from '../../types/share-link'
import { RELATION_LABELS } from '../../constants/enums'

const { Title, Paragraph } = Typography

export default function ShareJoinPage() {
  const { shareCode } = useParams<{ shareCode: string }>()
  const navigate = useNavigate()
  const { isAuthenticated, loading: authLoading } = useAuth()

  const [validateResult, setValidateResult] = useState<ShareLinkValidResp | null>(null)
  const [validating, setValidating] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [joined, setJoined] = useState(false)

  const [joinOpen, setJoinOpen] = useState(false)
  const [joining, setJoining] = useState(false)
  const [selectedDeceased, setSelectedDeceased] = useState<number | null>(null)
  const [selectedRelation, setSelectedRelation] = useState<number | null>(null)
  const [note, setNote] = useState('')

  useEffect(() => {
    if (!shareCode) return
    setValidating(true)
    shareLinkApi.validateShareLink(shareCode)
      .then((res) => {
        if (res.code === 0 && res.data) {
          setValidateResult(res.data)
        } else {
          setError(res.message || '分享链接无效')
        }
      })
      .catch(() => {
        setError('验证分享链接失败')
      })
      .finally(() => setValidating(false))
  }, [shareCode])

  const handleJoin = async () => {
    if (!shareCode || !selectedDeceased || !selectedRelation) return
    setJoining(true)
    try {
      const res = await shareLinkApi.joinShareLink(shareCode, {
        deceasedID: selectedDeceased,
        relationType: selectedRelation,
        note,
      })
      if (res.code === 0) {
        setJoined(true)
        setJoinOpen(false)
      } else {
        setError(res.message || '加入失败')
      }
    } catch {
      setError('加入失败，请重试')
    } finally {
      setJoining(false)
    }
  }

  const openJoinModal = () => {
    setSelectedDeceased(null)
    setSelectedRelation(null)
    setNote('')
    setJoinOpen(true)
  }

  const handleGoToGrave = () => {
    if (validateResult) {
      navigate(`/grave/${validateResult.graveID}`)
    }
  }

  if (validating || authLoading) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', background: '#f0f2f5' }}>
        <Spin size="large" tip="验证分享链接..." />
      </div>
    )
  }

  if (error) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', background: '#f0f2f5' }}>
        <Result
          status="error"
          title="分享链接已失效"
          subTitle={error}
        />
      </div>
    )
  }

  if (!validateResult) return null

  const redirectUrl = `/join/${shareCode}`

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', background: '#f0f2f5' }}>
      <Card style={{ width: 480, boxShadow: '0 2px 8px rgba(0,0,0,0.1)' }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <Title level={2} style={{ marginBottom: 8 }}>🕯️ RememberMe</Title>
          <Title level={3} style={{ marginBottom: 4, fontWeight: 'normal' }}>{validateResult.graveName}</Title>
          {validateResult.graveDesc && (
            <Paragraph type="secondary">{validateResult.graveDesc}</Paragraph>
          )}
        </div>

        {!isAuthenticated && (
          <div style={{ textAlign: 'center' }}>
            <Paragraph type="secondary">你被邀请加入这个纪念园</Paragraph>
            <Space size="middle">
              <Button type="primary" size="large" onClick={() => navigate(`/login?redirect=${encodeURIComponent(redirectUrl)}`)}>
                登录
              </Button>
              <Button size="large" onClick={() => navigate(`/register?redirect=${encodeURIComponent(redirectUrl)}`)}>
                注册
              </Button>
            </Space>
          </div>
        )}

        {isAuthenticated && !joined && (
          <div style={{ textAlign: 'center' }}>
            <Button
              type="primary"
              size="large"
              icon={<ArrowRightOutlined />}
              onClick={openJoinModal}
            >
              加入墓园
            </Button>
          </div>
        )}

        {isAuthenticated && joined && (
          <div style={{ textAlign: 'center' }}>
            <Result
              status="success"
              title="加入成功"
              subTitle={`你已成为「${validateResult.graveName}」的成员`}
              extra={
                <Button type="primary" size="large" onClick={handleGoToGrave}>
                  前往墓园
                </Button>
              }
            />
          </div>
        )}
      </Card>

      <Modal
        title="加入墓园"
        open={joinOpen}
        onCancel={() => setJoinOpen(false)}
        onOk={handleJoin}
        confirmLoading={joining}
        okText="确认加入"
        okButtonProps={{ disabled: !selectedDeceased || !selectedRelation }}
      >
        <Form layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item label="选择逝者" required>
            <Select
              placeholder="请选择一位逝者"
              value={selectedDeceased}
              onChange={(v) => setSelectedDeceased(v)}
              options={(validateResult.deceased || []).map((d: DeceasedBrief) => ({
                label: `${d.name}${d.gender === 1 ? ' (男)' : d.gender === 2 ? ' (女)' : ''}`,
                value: d.id,
              }))}
            />
          </Form.Item>

          <Form.Item label="你与该逝者的关系" required>
            <Select
              placeholder="请选择关系"
              value={selectedRelation}
              onChange={(v) => setSelectedRelation(v)}
              options={Object.entries(RELATION_LABELS).map(([k, v]) => ({
                label: v,
                value: Number(k),
              }))}
            />
          </Form.Item>

          <Form.Item label="留言（可选）">
            <Input.TextArea
              rows={2}
              value={note}
              onChange={(e) => setNote(e.target.value)}
              placeholder="写给墓园管理员的话..."
            />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
