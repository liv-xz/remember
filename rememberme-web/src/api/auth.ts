import client from './client'
import type { ApiResponse, LoginReq, LoginResp, RegisterReq, RefreshReq, ChangePasswordReq } from '../types/auth'

export async function register(data: RegisterReq): Promise<ApiResponse> {
  return client.post<ApiResponse>('/auth/register', data)
}

export async function login(data: LoginReq): Promise<ApiResponse<LoginResp>> {
  return client.post<ApiResponse<LoginResp>>('/auth/login', data)
}

export async function refresh(data: RefreshReq): Promise<ApiResponse<LoginResp>> {
  return client.post<ApiResponse<LoginResp>>('/auth/refresh', data)
}

export async function logout(data: RefreshReq): Promise<ApiResponse> {
  return client.post<ApiResponse>('/auth/logout', data)
}

export async function changePassword(data: ChangePasswordReq): Promise<ApiResponse> {
  return client.put<ApiResponse>('/auth/password', data)
}
