import { useEffect, useState } from 'react'
import { Modal, Form, Input, Select, DatePicker, Upload, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import type { RcFile } from 'antd/es/upload/interface'
import dayjs from 'dayjs'
import * as deceasedApi from '../../api/deceased'
import type { DeceasedResp, CreateDeceasedReq } from '../../types/deceased'
import { GENDER_OPTIONS, RELATION_OPTIONS } from '../../constants/enums'
import { validateImageFile, handleUploadWithState } from '../../utils/upload'

const { useWatch } = Form

interface Props {
  open: boolean
  graveId: number
  graveType: number
  graveMode: 1 | 2
  deceasedCount: number
  editingDeceased: DeceasedResp | null
  onCancel: () => void
  onSuccess: () => void
}

export default function DeceasedForm({ open, graveId, graveType, graveMode, deceasedCount, editingDeceased, onCancel, onSuccess }: Props) {
  const [form] = Form.useForm()
  const isEdit = !!editingDeceased
  const [portraitUrl, setPortraitUrl] = useState('')
  const [portraitUploading, setPortraitUploading] = useState(false)
  const [saving, setSaving] = useState(false)

  const relationRequired = !isEdit && graveMode === 2
  const watchedRelationType = useWatch('userRelationType', form)

  useEffect(() => {
    if (open) {
      if (editingDeceased) {
        form.setFieldsValue({
          ...editingDeceased,
          birthDate: editingDeceased.birthDate ? dayjs(editingDeceased.birthDate) : null,
          deathDate: editingDeceased.deathDate ? dayjs(editingDeceased.deathDate) : null,
        })
        setPortraitUrl(editingDeceased.portraitUrl || '')
      } else {
        form.resetFields()
        setPortraitUrl('')
      }
    }
  }, [open, editingDeceased, form])

  const beforeUpload = validateImageFile

  const handlePortraitUpload = async (file: RcFile) => {
    return handleUploadWithState(file, setPortraitUploading, 'portrait').then((url) => {
      setPortraitUrl(url)
      form.setFieldValue('portraitUrl', url)
      return url
    })
  }

  const handleSubmit = async () => {
    setSaving(true)
    try {
      const values = await form.validateFields()

      if (!isEdit) {
        const maxCapacity = graveType === 1 ? 1 : graveType === 2 ? 2 : Infinity
        if (deceasedCount >= maxCapacity) {
          message.error(`该墓园最多容纳 ${maxCapacity === 1 ? '1位' : '2位'} 逝者`)
          return
        }
      }

      const payload: CreateDeceasedReq & { userRelation?: { relationType: number; note?: string } } = {
        ...values,
        birthDate: values.birthDate ? values.birthDate.format('YYYY-MM-DD') : undefined,
        deathDate: values.deathDate ? values.deathDate.format('YYYY-MM-DD') : undefined,
      }

      if (!isEdit && values.userRelationType) {
        payload.userRelation = {
          relationType: values.userRelationType,
          note: values.userRelationNote || '',
        }
      }

      let res
      if (isEdit) {
        res = await deceasedApi.updateDeceased(editingDeceased!.id, payload)
      } else {
        res = await deceasedApi.createDeceased(graveId, payload as CreateDeceasedReq)
      }

      if (res.code === 0) {
        message.success(isEdit ? '更新成功' : '添加成功')
        onSuccess()
      } else {
        message.error(res.message)
      }
    } catch (err) {
      console.error('Deceased form submit failed:', err)
      message.error('操作失败，请稍后重试')
    } finally {
      setSaving(false)
    }
  }

  return (
    <Modal
      title={isEdit ? '编辑逝者' : '添加逝者'}
      open={open}
      onOk={handleSubmit}
      onCancel={onCancel}
      destroyOnClose
      confirmLoading={saving}
      width={560}
    >
      <Form form={form} layout="vertical" initialValues={{ gender: 1 }}>
        <Form.Item name="name" label="姓名" rules={[{ required: true, message: '请输入姓名' }]}>
          <Input placeholder="例如：张老丈" maxLength={50} />
        </Form.Item>
        <Form.Item name="gender" label="性别" rules={[{ required: true }]}>
          <Select options={GENDER_OPTIONS} />
        </Form.Item>
        <Form.Item name="birthDate" label="出生日期">
          <DatePicker style={{ width: '100%' }} placeholder="选择出生日期" />
        </Form.Item>
        <Form.Item
          name="deathDate"
          label="逝世日期"
          rules={[
            ({ getFieldValue }) => ({
              validator(_, value) {
                const birth = getFieldValue('birthDate')
                if (!value || !birth || value.isAfter(birth) || value.isSame(birth)) {
                  return Promise.resolve()
                }
                return Promise.reject(new Error('逝世日期不能早于出生日期'))
              },
            }),
          ]}
        >          <DatePicker style={{ width: '100%' }} placeholder="选择逝世日期" />
        </Form.Item>
        <Form.Item name="portraitUrl" label="遗像">
          <Upload
            listType="picture-card"
            maxCount={1}
            showUploadList={false}
            beforeUpload={beforeUpload}
            customRequest={async ({ file }) => {
              await handlePortraitUpload(file as RcFile)
            }}
          >
            {portraitUrl ? (
              <img src={portraitUrl} alt="遗像" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
            ) : (
              <div>
                {portraitUploading ? '上传中...' : <PlusOutlined />}
                <div style={{ marginTop: 8 }}>遗像</div>
              </div>
            )}
          </Upload>
        </Form.Item>
        <Form.Item name="epitaph" label="墓志铭">
          <Input placeholder="例如：慈父严师" maxLength={500} />
        </Form.Item>
        <Form.Item name="biography" label="生平简介">
          <Input.TextArea rows={3} placeholder="逝者生平简介" maxLength={2000} />
        </Form.Item>

        {!isEdit && (
          <>
            <Form.Item
              name="userRelationType"
              label="我与该逝者的关系"
              rules={relationRequired ? [{ required: true, message: '请选择关系' }] : []}
            >
              <Select
                allowClear
                placeholder={relationRequired ? '请选择你与逝者的关系' : '可选：你与逝者的关系'}
                options={RELATION_OPTIONS}
              />
            </Form.Item>
            <Form.Item name="userRelationNote" label="备注" hidden={!watchedRelationType}>
              <Input placeholder="关系备注（可选）" maxLength={255} />
            </Form.Item>
            {graveMode === 2 && (
              <Form.Item name="originGraveID" label="来源墓园ID">
                <Input placeholder="嫁入者的来源墓园ID（可选）" type="number" />
              </Form.Item>
            )}
          </>
        )}
      </Form>
    </Modal>
  )
}
