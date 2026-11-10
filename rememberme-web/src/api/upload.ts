import client from './client'
import type { ApiResponse } from '../types/auth'

export async function uploadFile(file: File, subDir: string): Promise<string> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('subDir', subDir)

  const res = await client.post<ApiResponse<{ url: string }>>('/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })

  if (res.code === 0 && res.data) {
    return res.data.url
  }
  throw new Error(res.message || '上传失败')
}
