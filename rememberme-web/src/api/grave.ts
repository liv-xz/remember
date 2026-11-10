import client from './client'
import type { ApiResponse } from '../types/auth'
import type { CreateGraveReq, CreateGraveWithDeceasedReq, GraveResp, GraveListResp, GraveWithDeceasedResp, UpdateGraveReq } from '../types/grave'

export async function getMyGraves(page = 1, size = 100): Promise<ApiResponse<GraveListResp>> {
  return client.get<ApiResponse<GraveListResp>>('/graves', { params: { page, size } })
}

export async function getGraveById(id: number): Promise<ApiResponse<GraveResp>> {
  return client.get<ApiResponse<GraveResp>>(`/graves/${id}`)
}

export async function createGrave(data: CreateGraveReq): Promise<ApiResponse<GraveResp>> {
  return client.post<ApiResponse<GraveResp>>('/graves', data)
}

export async function updateGrave(id: number, data: UpdateGraveReq): Promise<ApiResponse<GraveResp>> {
  return client.put<ApiResponse<GraveResp>>(`/graves/${id}`, data)
}

export async function deleteGrave(id: number): Promise<ApiResponse> {
  return client.delete<ApiResponse>(`/graves/${id}`)
}

export async function createGraveWithDeceased(data: CreateGraveWithDeceasedReq): Promise<ApiResponse<GraveWithDeceasedResp>> {
  return client.post<ApiResponse<GraveWithDeceasedResp>>('/graves/with-deceased', data)
}
