import type { RcFile, UploadFile } from 'antd/es/upload'
import { message } from 'antd'
import * as uploadApi from '../api/upload'

export function validateImageFile(file: RcFile): boolean {
  if (!file.type.startsWith('image/')) {
    message.error('只能上传图片文件')
    return false
  }
  if (file.size / 1024 / 1024 > 5) {
    message.error('图片大小不能超过 5MB')
    return false
  }
  return true
}

export async function uploadCoverFile(file: RcFile): Promise<string> {
  const url = await uploadApi.uploadFile(file, 'cover')
  return url
}

export async function uploadPortraitFile(file: RcFile): Promise<string> {
  const url = await uploadApi.uploadFile(file, 'portrait')
  return url
}

export async function handleUploadWithState(
  file: RcFile,
  setUploading: (v: boolean) => void,
  folder: string,
): Promise<string> {
  setUploading(true)
  try {
    const url = await uploadApi.uploadFile(file, folder)
    return url
  } catch {
    message.error('图片上传失败')
    throw new Error('upload failed')
  } finally {
    setUploading(false)
  }
}
