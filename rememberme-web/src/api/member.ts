import client from './client'
import type { ApiResponse } from '../types/auth'
import type {
  AddMemberReq,
  GraveMemberResp,
  GraveMemberListResp,
  UpdateRoleReq,
  UpdateRelationDescReq,
  TransferOwnershipReq,
  TransferOwnershipResp,
} from '../types/member'

export async function getMembersByGrave(graveId: number): Promise<ApiResponse<GraveMemberListResp>> {
  return client.get<ApiResponse<GraveMemberListResp>>(`/graves/${graveId}/members`)
}

export async function addMember(graveId: number, data: AddMemberReq): Promise<ApiResponse<GraveMemberResp>> {
  return client.post<ApiResponse<GraveMemberResp>>(`/graves/${graveId}/members`, data)
}

export async function updateMemberRole(memberId: number, data: UpdateRoleReq): Promise<ApiResponse<GraveMemberResp>> {
  return client.put<ApiResponse<GraveMemberResp>>(`/members/${memberId}/role`, data)
}

export async function updateMemberRelationDesc(memberId: number, data: UpdateRelationDescReq): Promise<ApiResponse<GraveMemberResp>> {
  return client.put<ApiResponse<GraveMemberResp>>(`/members/${memberId}/relation-desc`, data)
}

export async function removeMember(memberId: number): Promise<ApiResponse> {
  return client.delete<ApiResponse>(`/members/${memberId}`)
}

export async function transferOwnership(data: TransferOwnershipReq): Promise<ApiResponse<TransferOwnershipResp>> {
  return client.post<ApiResponse<TransferOwnershipResp>>('/members/transfer', data)
}
