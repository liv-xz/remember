import client from './client'
import type { ApiResponse } from '../types/auth'
import type { MemorialRelationListResp, CreateMemorialRelationReq } from '../types/memorial-relation'

export async function createMemorialRelation(
  graveId: number,
  data: CreateMemorialRelationReq
): Promise<ApiResponse> {
  return client.post<ApiResponse>(`/graves/${graveId}/memorial-relations`, data)
}

export async function getMemorialRelations(graveId: number): Promise<ApiResponse<MemorialRelationListResp>> {
  return client.get<ApiResponse<MemorialRelationListResp>>(`/graves/${graveId}/memorial-relations`)
}

export async function confirmMemorialRelation(relationId: number): Promise<ApiResponse> {
  return client.put<ApiResponse>(`/memorial-relations/${relationId}/confirm`)
}

export async function rejectMemorialRelation(relationId: number): Promise<ApiResponse> {
  return client.put<ApiResponse>(`/memorial-relations/${relationId}/reject`)
}
