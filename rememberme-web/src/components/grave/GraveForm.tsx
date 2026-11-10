import { useEffect, useState } from 'react'
import { Modal, Form, Input, Select, Upload, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import type { RcFile } from 'antd/es/upload/interface'
import * as graveApi from '../../api/grave'
import type { GraveResp, CreateGraveReq, LatLng } from '../../types/grave'
import MapPicker from './MapPicker'
import { GRAVE_TYPE_OPTIONS, LAYOUT_OPTIONS, GRAVE_MODE_OPTIONS, GRAVE_STATUS_OPTIONS } from '../../constants/enums'
import { validateImageFile, handleUploadWithState } from '../../utils/upload'

interface Props {
  open: boolean
  editingGrave: GraveResp | null
  onCancel: () => void
  onSuccess: (graveId?: number) => void
}

function graveToLatLng(grave: GraveResp): LatLng | undefined {
  if (grave.latitude != null && grave.longitude != null) {
    return {
      latitude: grave.latitude,
      longitude: grave.longitude,
      address: grave.address || '',
    }
  }
  return undefined
}

export default function GraveForm({ open, editingGrave, onCancel, onSuccess }: Props) {
  const [form] = Form.useForm()
  const [location, setLocation] = useState<LatLng | null>(null)
  const [coverUrl, setCoverUrl] = useState('')
  const [coverUploading, setCoverUploading] = useState(false)
  const [saving, setSaving] = useState(false)
  const isEdit = !!editingGrave
  const watchedGraveType: number = Form.useWatch('graveType', form)

  useEffect(() => {
    if (watchedGraveType) {
      form.setFieldValue('layoutStyle', watchedGraveType)
    }
  }, [watchedGraveType, form])

  useEffect(() => {
    if (open) {
      if (editingGrave) {
        form.setFieldsValue(editingGrave)
        setLocation(graveToLatLng(editingGrave) ?? null)
        setCoverUrl(editingGrave.coverImage || '')
      } else {
        form.resetFields()
        setLocation(null)
        setCoverUrl('')
      }
    }
  }, [open, editingGrave, form])

  const handleCoverUpload = async (file: RcFile) => {
    return handleUploadWithState(file, setCoverUploading, 'cover').then((url) => {
      setCoverUrl(url)
      return url
    })
  }

  const beforeUpload = validateImageFile

  const handleSubmit = async () => {
    setSaving(true)
    try {
      const values = await form.validateFields()

      if (isEdit) {
        const payload: Record<string, unknown> = { ...values }
        if (coverUrl) {
          payload.coverImage = coverUrl
        }
        if (location && location.latitude !== 0 && location.longitude !== 0) {
          payload.latitude = location.latitude
          payload.longitude = location.longitude
          payload.address = location.address
        } else {
          payload.clearLocation = true
        }
        const res = await graveApi.updateGrave(editingGrave!.id, payload as never)
        if (res.code === 0) {
          message.success('更新成功')
          onSuccess()
        } else {
          message.error(res.message)
        }
      } else {
        const payload: CreateGraveReq = {
          name: values.name,
          graveType: values.graveType,
          graveMode: values.graveMode,
          layoutStyle: values.layoutStyle,
          description: values.description || '',
        }
        if (coverUrl) {
          payload.coverImage = coverUrl
        }
        if (location && location.latitude !== 0 && location.longitude !== 0) {
          payload.latitude = location.latitude
          payload.longitude = location.longitude
          payload.address = location.address
        }
        const res = await graveApi.createGrave(payload)
        if (res.code === 0 && res.data) {
          message.success('创建成功')
          onSuccess(res.data.id)
        } else {
          message.error(res.message)
        }
      }
    } catch (err) {
      console.error('Grave form submit failed:', err)
      message.error('操作失败，请稍后重试')
    } finally {
      setSaving(false)
    }
  }

  return (
    <Modal
      title={isEdit ? '编辑墓园' : '创建墓园'}
      open={open}
      onOk={handleSubmit}
      onCancel={onCancel}
      destroyOnClose
      confirmLoading={saving}
      width={640}
    >
      <Form form={form} layout="vertical" initialValues={{ graveType: 1, layoutStyle: 1 }}>
        <Form.Item name="name" label="墓园名称" rules={[{ required: true, message: '请输入墓园名称' }]}>
          <Input placeholder="例如：先祖陵" maxLength={50} />
        </Form.Item>
        <Form.Item name="graveType" label="墓地类型" rules={[{ required: true }]}>
          <Select options={GRAVE_TYPE_OPTIONS} />
        </Form.Item>
        <Form.Item name="layoutStyle" label="墓碑布局" rules={[{ required: true }]}>
          <Select options={LAYOUT_OPTIONS} />
        </Form.Item>
        <Form.Item name="graveMode" label="墓园模式" rules={[{ required: true }]}>
          <Select disabled={isEdit} options={GRAVE_MODE_OPTIONS} />
        </Form.Item>
        <Form.Item label="封面图片">
          <Upload
            listType="picture-card"
            maxCount={1}
            showUploadList={false}
            beforeUpload={beforeUpload}
            customRequest={async ({ file }) => {
              await handleCoverUpload(file as RcFile)
            }}
          >
            {coverUrl ? (
              <img src={coverUrl} alt="封面" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
            ) : (
              <div>
                {coverUploading ? '上传中...' : <PlusOutlined />}
                <div style={{ marginTop: 8 }}>封面图</div>
              </div>
            )}
          </Upload>
        </Form.Item>
        <Form.Item name="description" label="墓园描述">
          <Input.TextArea rows={3} placeholder="描述一下这座墓园（可选）" maxLength={500} />
        </Form.Item>
        <Form.Item label="墓园位置">
          <MapPicker value={location ?? undefined} onChange={setLocation} />
        </Form.Item>
        {isEdit && (
          <Form.Item name="status" label="状态">
            <Select options={GRAVE_STATUS_OPTIONS} />
          </Form.Item>
        )}
      </Form>
    </Modal>
  )
}
