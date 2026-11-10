import client from './client'
import type { ApiResponse } from '../types/auth'
import type { CreateDeceasedReq, DeceasedResp, DeceasedListResp, UpdateDeceasedReq } from '../types/deceased'

export async function getDeceasedByGrave(graveId: number): Promise<ApiResponse<DeceasedListResp>> {
  return client.get<ApiResponse<DeceasedListResp>>(`/graves/${graveId}/deceased`)
}

export async function createDeceased(graveId: number, data: CreateDeceasedReq): Promise<ApiResponse<DeceasedResp>> {
  return client.post<ApiResponse<DeceasedResp>>(`/graves/${graveId}/deceased`, data)
}

export async function updateDeceased(id: number, data: UpdateDeceasedReq): Promise<ApiResponse<DeceasedResp>> {
  return client.put<ApiResponse<DeceasedResp>>(`/deceased/${id}`, data)
}

export async function deleteDeceased(id: number): Promise<ApiResponse> {
  return client.delete<ApiResponse>(`/deceased/${id}`)
}
