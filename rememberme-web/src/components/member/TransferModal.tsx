import { useState } from 'react'
import { Modal, Select, message } from 'antd'
import * as memberApi from '../../api/member'
import type { GraveMemberResp } from '../../types/member'

interface Props {
  open: boolean
  members: GraveMemberResp[]
  onCancel: () => void
  onSuccess: () => void
}

export default function TransferModal({ open, members, onCancel, onSuccess }: Props) {
  const [targetId, setTargetId] = useState<number | undefined>()
  const [loading, setLoading] = useState(false)

  const adminMembers = members.filter((m) => m.role === 20)

  const handleTransfer = async () => {
    if (!targetId) {
      message.warning('请选择接收方')
      return
    }
    setLoading(true)
    try {
      const res = await memberApi.transferOwnership({ targetMemberID: targetId })
      if (res.code === 0) {
        message.success('所有权转让成功')
        onSuccess()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('转让失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <Modal
      title="转让所有权"
      open={open}
      onOk={handleTransfer}
      onCancel={onCancel}
      confirmLoading={loading}
    >
      <p style={{ marginBottom: 16, color: '#666' }}>
        转让后您将降级为管理员，接收方成为新的所有者。只能转让给管理员。
      </p>
      {adminMembers.length === 0 ? (
        <p style={{ color: '#ff4d4f' }}>没有可接收的管理员，请先将成员升级为管理员。</p>
      ) : (
        <Select
          style={{ width: '100%' }}
          placeholder="选择接收方"
          value={targetId}
          onChange={setTargetId}
          options={adminMembers.map((m) => ({
            value: m.id,
            label: `${m.nickname || m.username} (${m.relationDesc || '无关系描述'})`,
          }))}
        />
      )}
    </Modal>
  )
}
