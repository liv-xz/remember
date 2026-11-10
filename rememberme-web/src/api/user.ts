import client from './client'
import type { ApiResponse } from '../types/auth'
import type { UserProfileResp, UpdateProfileReq } from '../types/user'

export async function getProfile(): Promise<ApiResponse<UserProfileResp>> {
  return client.get<ApiResponse<UserProfileResp>>('/user/profile')
}

export async function updateProfile(data: UpdateProfileReq): Promise<ApiResponse<UserProfileResp>> {
  return client.put<ApiResponse<UserProfileResp>>('/user/profile', data)
}
