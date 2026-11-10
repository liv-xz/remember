import client from './client'
import type { ApiResponse } from '../types/auth'
import type { MemorialRecordResp } from '../types/memorial-record'

interface RecordListData extends ApiResponse<MemorialRecordResp[]> {
  total?: number
  page?: number
  size?: number
}

export async function createRecord(
  graveId: number,
  data: { type: number; deceasedID?: number; message?: string }
): Promise<ApiResponse<MemorialRecordResp>> {
  return client.post<ApiResponse<MemorialRecordResp>>(`/graves/${graveId}/records`, data)
}

export async function getRecordsByGrave(
  graveId: number,
  page?: number,
  size?: number
): Promise<RecordListData> {
  return client.get<RecordListData>(`/graves/${graveId}/records`, { params: { page, size } })
}

export async function deleteRecord(recordId: number): Promise<ApiResponse> {
  return client.delete<ApiResponse>(`/records/${recordId}`)
}
