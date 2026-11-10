import { Modal, Form, Input, InputNumber, Select, message } from 'antd'
import * as memberApi from '../../api/member'

interface Props {
  open: boolean
  graveId: number
  onCancel: () => void
  onSuccess: () => void
}

export default function MemberForm({ open, graveId, onCancel, onSuccess }: Props) {
  const [form] = Form.useForm()

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      const res = await memberApi.addMember(graveId, values)
      if (res.code === 0) {
        message.success('添加成功')
        form.resetFields()
        onSuccess()
      } else {
        message.error(res.message)
      }
    } catch {
      // validation failed
    }
  }

  return (
    <Modal
      title="添加成员"
      open={open}
      onOk={handleSubmit}
      onCancel={() => { form.resetFields(); onCancel() }}
      destroyOnClose
    >
      <Form form={form} layout="vertical" initialValues={{ role: 30 }}>
        <Form.Item name="userID" label="用户ID" rules={[{ required: true, message: '请输入用户ID' }]}>
          <InputNumber style={{ width: '100%' }} placeholder="输入目标用户的ID" min={1} />
        </Form.Item>
        <Form.Item name="role" label="角色" rules={[{ required: true }]}>
          <Select
            options={[
              { value: 20, label: '管理员' },
              { value: 30, label: '成员' },
            ]}
          />
        </Form.Item>
        <Form.Item name="relationDesc" label="关系描述">
          <Input placeholder="如：侄孙、长孙" maxLength={50} />
        </Form.Item>
      </Form>
    </Modal>
  )
}
