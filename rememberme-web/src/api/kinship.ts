import client from './client'
import type { ApiResponse } from '../types/auth'
import type { KinshipRelationListResp, CreateKinshipRelationReq } from '../types/kinship'

export async function createKinshipRelation(
  graveId: number,
  data: CreateKinshipRelationReq
): Promise<ApiResponse> {
  return client.post<ApiResponse>(`/graves/${graveId}/kinship-relations`, data)
}

export async function getKinshipRelations(graveId: number): Promise<ApiResponse<KinshipRelationListResp>> {
  return client.get<ApiResponse<KinshipRelationListResp>>(`/graves/${graveId}/kinship-relations`)
}

export async function confirmKinshipRelation(relationId: number): Promise<ApiResponse> {
  return client.put<ApiResponse>(`/kinship-relations/${relationId}/confirm`)
}

export async function rejectKinshipRelation(relationId: number): Promise<ApiResponse> {
  return client.put<ApiResponse>(`/kinship-relations/${relationId}/reject`)
}
