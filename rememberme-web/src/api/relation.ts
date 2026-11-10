import client from './client'
import type { ApiResponse } from '../types/auth'
import type { DeceasedRelationResp, DeceasedRelationListResp, CreateDeceasedRelationReq } from '../types/relation'

export async function getDeceasedRelations(graveId: number): Promise<ApiResponse<DeceasedRelationListResp>> {
  return client.get<ApiResponse<DeceasedRelationListResp>>(`/graves/${graveId}/kinship-deceased-relations`)
}

export async function createDeceasedRelation(
  graveId: number,
  data: CreateDeceasedRelationReq
): Promise<ApiResponse<DeceasedRelationResp>> {
  return client.post<ApiResponse<DeceasedRelationResp>>(`/graves/${graveId}/kinship-deceased-relations`, data)
}

export async function deleteDeceasedRelation(id: number): Promise<ApiResponse> {
  return client.delete<ApiResponse>(`/kinship-deceased-relations/${id}`)
}
