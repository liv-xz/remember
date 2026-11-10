import client from './client'
import type { ApiResponse } from '../types/auth'
import type { ShareLinkListResp, CreateShareLinkReq, ShareLinkValidResp, JoinShareLinkReq } from '../types/share-link'

export async function createShareLink(graveId: number, data?: CreateShareLinkReq): Promise<ApiResponse<ShareLinkListResp['items'][0]>> {
  return client.post<ApiResponse<ShareLinkListResp['items'][0]>>(`/graves/${graveId}/share-links`, data || {})
}

export async function getShareLinks(graveId: number): Promise<ApiResponse<ShareLinkListResp>> {
  return client.get<ApiResponse<ShareLinkListResp>>(`/graves/${graveId}/share-links`)
}

export async function deactivateShareLink(graveId: number, linkId: number): Promise<ApiResponse> {
  return client.put<ApiResponse>(`/graves/${graveId}/share-links/${linkId}/deactivate`)
}

export async function validateShareLink(shareCode: string): Promise<ApiResponse<ShareLinkValidResp>> {
  return client.get<ApiResponse<ShareLinkValidResp>>(`/share-links/${shareCode}`)
}

export async function joinShareLink(shareCode: string, data: JoinShareLinkReq): Promise<ApiResponse<ShareLinkValidResp>> {
  return client.post<ApiResponse<ShareLinkValidResp>>(`/share-links/${shareCode}/join`, data)
}