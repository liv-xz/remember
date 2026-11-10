import { useEffect, useState } from 'react'
import { Typography, Form, Input, Button, message, Card, Divider, Modal, Space } from 'antd'
import { useAuth } from '../../contexts/AuthContext'
import * as userApi from '../../api/user'
import * as authApi from '../../api/auth'
import type { UserProfileResp } from '../../types/user'

const { Title } = Typography

export default function ProfilePage() {
  const { profile, updateProfile, logout } = useAuth()
  const [form] = Form.useForm()
  const [saving, setSaving] = useState(false)
  const [pwdModalOpen, setPwdModalOpen] = useState(false)
  const [pwdForm] = Form.useForm()
  const [pwdSaving, setPwdSaving] = useState(false)

  useEffect(() => {
    if (profile) {
      form.setFieldsValue({
        nickname: profile.nickname,
        avatar: profile.avatar,
        abstract: profile.abstract,
      })
    }
  }, [profile, form])

  const handleSaveProfile = async (values: { nickname: string; avatar: string; abstract: string }) => {
    setSaving(true)
    try {
      const res = await userApi.updateProfile(values)
      if (res.code === 0 && res.data) {
        message.success('保存成功')
        updateProfile(res.data)
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('保存失败')
    } finally {
      setSaving(false)
    }
  }

  const handleChangePassword = async (values: { oldPassword: string; newPassword: string }) => {
    setPwdSaving(true)
    try {
      const res = await authApi.changePassword(values)
      if (res.code === 0) {
        message.success('密码修改成功，请重新登录')
        setPwdModalOpen(false)
        pwdForm.resetFields()
        await logout()
      } else {
        message.error(res.message)
      }
    } catch {
      message.error('修改密码失败')
    } finally {
      setPwdSaving(false)
    }
  }

  return (
    <div className="profile-container" style={{ maxWidth: 640, margin: '0 auto' }}>
      <Title level={3}>个人资料</Title>

      <Card>
        <Form layout="vertical" form={form} onFinish={handleSaveProfile}>
          <Form.Item label="用户名">
            <Input value={profile?.username} disabled />
          </Form.Item>
          <Form.Item label="邮箱">
            <Input value={profile?.email} disabled />
          </Form.Item>
          <Form.Item name="nickname" label="昵称">
            <Input placeholder="设置昵称" />
          </Form.Item>
          <Form.Item name="avatar" label="头像URL">
            <Input placeholder="输入头像图片地址" />
          </Form.Item>
          <Form.Item name="abstract" label="个人简介">
            <Input.TextArea rows={3} placeholder="介绍一下自己" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={saving}>
              保存修改
            </Button>
          </Form.Item>
        </Form>
      </Card>

      <Divider />

      <Card>
        <Title level={5}>账号安全</Title>
        <Button onClick={() => setPwdModalOpen(true)}>修改密码</Button>
      </Card>

      <Modal
        title="修改密码"
        open={pwdModalOpen}
        onCancel={() => { setPwdModalOpen(false); pwdForm.resetFields() }}
        footer={null}
      >
        <Form form={pwdForm} layout="vertical" onFinish={handleChangePassword}>
          <Form.Item name="oldPassword" label="当前密码" rules={[{ required: true, message: '请输入当前密码' }]}>
            <Input.Password />
          </Form.Item>
          <Form.Item name="newPassword" label="新密码" rules={[
            { required: true, message: '请输入新密码' },
            { min: 8, message: '密码至少8位' },
            { pattern: /^(?=.*[a-zA-Z])(?=.*\d)/, message: '密码需包含字母和数字' },
          ]}>
            <Input.Password />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" loading={pwdSaving}>确认修改</Button>
              <Button onClick={() => { setPwdModalOpen(false); pwdForm.resetFields() }}>取消</Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}
