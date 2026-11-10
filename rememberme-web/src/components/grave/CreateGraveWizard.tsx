import { useState, useEffect } from 'react'
import { Modal, Steps, Button, Form, Input, Select, DatePicker, Upload, message, Space, Card } from 'antd'
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons'
import type { RcFile } from 'antd/es/upload/interface'
import * as graveApi from '../../api/grave'
import * as uploadApi from '../../api/upload'
import type { CreateGraveWithDeceasedReq, DeceasedItem } from '../../types/grave'
import { GRAVE_TYPE_OPTIONS, LAYOUT_OPTIONS, GENDER_OPTIONS } from '../../constants/enums'
import { validateImageFile } from '../../utils/upload'
// TODO: 地图功能后续研究
// import MapPicker from './MapPicker'

const { Step } = Steps

interface Props {
  open: boolean
  graveMode: 1 | 2
  onCancel: () => void
  onSuccess: (graveId: number) => void
}

// TODO: 地图功能后续研究
/*
interface LatLng {
  latitude: number
  longitude: number
  address: string
}
*/

interface DeceasedFormData {
  key: number
  portraitUrl: string
  uploading: boolean
}

export default function CreateGraveWizard({ open, graveMode, onCancel, onSuccess }: Props) {
  const [step, setStep] = useState(0)
  const [submitting, setSubmitting] = useState(false)
  const [coverUrl, setCoverUrl] = useState('')
  const [coverUploading, setCoverUploading] = useState(false)
  // TODO: 地图功能后续研究
  // const [location, setLocation] = useState<LatLng | null>(null)
  const [deceasedForms, setDeceasedForms] = useState<DeceasedFormData[]>([
    { key: 0, portraitUrl: '', uploading: false },
  ])
  const [graveFormValues, setGraveFormValues] = useState<Record<string, unknown> | null>(null)

  const [graveForm] = Form.useForm()
  const [deceasedForm] = Form.useForm()

  const graveType = Form.useWatch('graveType', graveForm)
  const maxDeceased = graveType === 1 ? 1 : graveType === 2 ? 2 : 99

  useEffect(() => {
    if (graveType) {
      graveForm.setFieldValue('layoutStyle', graveType)
    }
  }, [graveType, graveForm])

  const handleCoverUpload = async (file: RcFile) => {
    setCoverUploading(true)
    try {
      const url = await uploadApi.uploadFile(file, 'cover')
      setCoverUrl(url)
      return url
    } catch {
      message.error('封面图上传失败')
      throw new Error('upload failed')
    } finally {
      setCoverUploading(false)
    }
  }

  const handlePortraitUpload = async (file: RcFile, key: number) => {
    setDeceasedForms((prev) =>
      prev.map((f) => (f.key === key ? { ...f, uploading: true } : f))
    )
    try {
      const url = await uploadApi.uploadFile(file, 'portrait')
      setDeceasedForms((prev) =>
        prev.map((f) => (f.key === key ? { ...f, portraitUrl: url, uploading: false } : f))
      )
      return url
    } catch {
      message.error('遗像上传失败')
      setDeceasedForms((prev) =>
        prev.map((f) => (f.key === key ? { ...f, uploading: false } : f))
      )
      throw new Error('upload failed')
    }
  }

  const beforeUpload = validateImageFile

  const addDeceased = () => {
    if (deceasedForms.length >= maxDeceased) {
      message.warning(`最多添加 ${maxDeceased} 位逝者`)
      return
    }
    const newKey = Date.now()
    setDeceasedForms((prev) => [
      ...prev,
      { key: newKey, portraitUrl: '', uploading: false },
    ])
  }

  const removeDeceased = (key: number) => {
    if (deceasedForms.length <= 1) return
    setDeceasedForms((prev) => prev.filter((f) => f.key !== key))
  }

  const handleNext = async () => {
    try {
      const values = await graveForm.validateFields()
      setGraveFormValues(values)
      setStep(1)
    } catch {
      // validation failed
    }
  }

  const handlePrev = () => {
    setGraveFormValues(null)
    setStep(0)
  }

  const handleSubmit = async () => {
    try {
      if (!graveFormValues) return
      const deceasedValues = await deceasedForm.validateFields()
      setSubmitting(true)

      const deceasedList: DeceasedItem[] = deceasedForms.map((df) => ({
        name: deceasedValues[`name_${df.key}`] || '',
        gender: deceasedValues[`gender_${df.key}`] ?? 1,
        birthDate: deceasedValues[`birthDate_${df.key}`]
          ? deceasedValues[`birthDate_${df.key}`].format('YYYY-MM-DD')
          : '',
        deathDate: deceasedValues[`deathDate_${df.key}`]
          ? deceasedValues[`deathDate_${df.key}`].format('YYYY-MM-DD')
          : '',
        biography: deceasedValues[`biography_${df.key}`] || '',
        portraitUrl: df.portraitUrl,
        epitaph: deceasedValues[`epitaph_${df.key}`] || '',
      }))

      const payload: CreateGraveWithDeceasedReq = {
        name: graveFormValues.name as string,
        graveType: graveFormValues.graveType as number,
        graveMode,
        layoutStyle: graveFormValues.layoutStyle as number,
        description: (graveFormValues.description as string) || '',
        coverImage: coverUrl,
        deceased: deceasedList,
      }

      const res = await graveApi.createGraveWithDeceased(payload)
      if (res.code === 0 && res.data) {
        message.success('创建成功')
        onSuccess(res.data.grave.id)
      } else {
        message.error(res.message || '创建失败')
      }
    } catch (err) {
      console.error('Create grave failed:', err)
      message.error('创建失败，请稍后重试')
    } finally {
      setSubmitting(false)
    }
  }

  const handleClose = () => {
    setStep(0)
    setGraveFormValues(null)
    setCoverUrl('')
    // setLocation(null)
    setDeceasedForms([{ key: 0, portraitUrl: '', uploading: false }])
    graveForm.resetFields()
    deceasedForm.resetFields()
    onCancel()
  }

  const coverUploadButton = (
    <div>
      {coverUploading ? '上传中...' : <PlusOutlined />}
      <div style={{ marginTop: 8 }}>封面图</div>
    </div>
  )

  return (
    <Modal
      title="创建墓园"
      open={open}
      onCancel={handleClose}
      footer={null}
      destroyOnClose
      width={680}
    >
      <Steps current={step} size="small" style={{ marginBottom: 32 }}>
        <Step title="墓园信息" />
        <Step title="逝者信息" />
      </Steps>

      {step === 0 && (
        <Form form={graveForm} layout="vertical" initialValues={{ graveType: 3, layoutStyle: 3 }}>
          <Form.Item
            name="name"
            label="墓园名称"
            rules={[{ required: true, message: '请输入墓园名称' }]}
          >
            <Input placeholder="例如：张家墓园" maxLength={100} />
          </Form.Item>

          <Form.Item
            name="graveType"
            label="墓地类型"
            rules={[{ required: true, message: '请选择墓地类型' }]}
          >
            <Select options={GRAVE_TYPE_OPTIONS} />
          </Form.Item>

          <Form.Item
            name="layoutStyle"
            label="墓碑布局"
            rules={[{ required: true, message: '请选择墓碑布局' }]}
          >
            <Select options={LAYOUT_OPTIONS} />
          </Form.Item>

          <Form.Item name="description" label="墓园描述">
            <Input.TextArea rows={2} placeholder="墓园简介（可选）" maxLength={500} />
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
                coverUploadButton
              )}
            </Upload>
          </Form.Item>

          {/*<Form.Item label="墓园位置">
            <MapPicker value={location ?? undefined} onChange={setLocation} />
          </Form.Item>*/}
        </Form>
      )}

      {step === 1 && (
        <Form form={deceasedForm} layout="vertical">
          {deceasedForms.map((df, index) => (
            <Card
              key={df.key}
              size="small"
              title={`逝者 ${index + 1}`}
              style={{ marginBottom: 16 }}
              extra={
                deceasedForms.length > 1 && (
                  <Button
                    type="link"
                    danger
                    icon={<MinusCircleOutlined />}
                    onClick={() => removeDeceased(df.key)}
                    size="small"
                  >
                    移除
                  </Button>
                )
              }
            >
              <Form.Item
                name={`name_${df.key}`}
                label="姓名"
                rules={[{ required: true, message: '请输入逝者姓名' }]}
              >
                <Input placeholder="例如：张太公" maxLength={50} />
              </Form.Item>

              <Form.Item
                name={`gender_${df.key}`}
                label="性别"
                rules={[{ required: true, message: '请选择性别' }]}
                initialValue={1}
              >
                <Select options={GENDER_OPTIONS} />
              </Form.Item>

              <Space style={{ width: '100%' }} size={12}>
                <Form.Item name={`birthDate_${df.key}`} label="出生日期" style={{ flex: 1 }}>
                  <DatePicker style={{ width: '100%' }} placeholder="选择出生日期" />
                </Form.Item>
                <Form.Item name={`deathDate_${df.key}`} label="逝世日期" style={{ flex: 1 }}>
                  <DatePicker style={{ width: '100%' }} placeholder="选择逝世日期" />
                </Form.Item>
              </Space>

              <Form.Item label="遗像">
                <Upload
                  listType="picture-card"
                  maxCount={1}
                  showUploadList={false}
                  beforeUpload={beforeUpload}
                  customRequest={async ({ file }) => {
                    await handlePortraitUpload(file as RcFile, df.key)
                  }}
                >
                  {df.portraitUrl ? (
                    <img src={df.portraitUrl} alt="遗像" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
                  ) : (
                    <div>
                      {df.uploading ? '上传中...' : <PlusOutlined />}
                      <div style={{ marginTop: 8 }}>遗像</div>
                    </div>
                  )}
                </Upload>
              </Form.Item>

              <Form.Item name={`epitaph_${df.key}`} label="墓志铭">
                <Input placeholder="例如：慈父严师" maxLength={500} />
              </Form.Item>

              <Form.Item name={`biography_${df.key}`} label="生平简介">
                <Input.TextArea rows={2} placeholder="逝者生平简介（可选）" maxLength={2000} />
              </Form.Item>
            </Card>
          ))}

          {deceasedForms.length < maxDeceased && (
            <Button type="dashed" onClick={addDeceased} icon={<PlusOutlined />} block>
              添加逝者
            </Button>
          )}
        </Form>
      )}

      <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: 24 }}>
        <div>
          {step === 1 && (
            <Button onClick={handlePrev}>上一步</Button>
          )}
        </div>
        <div>
          {step === 0 && (
            <Button type="primary" onClick={handleNext}>下一步</Button>
          )}
          {step === 1 && (
            <Button type="primary" onClick={handleSubmit} loading={submitting}>
              提交创建
            </Button>
          )}
        </div>
      </div>
    </Modal>
  )
}
