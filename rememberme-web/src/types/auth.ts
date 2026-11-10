export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data?: T
  total?: number
  page?: number
  size?: number
}

export interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  email: string
}

export interface LoginReq {
  username: string
  password: string
}

export interface LoginResp {
  access_token: string
  refresh_token: string
  user: UserInfo
}

export interface RegisterReq {
  username: string
  email: string
  password: string
}

export interface RefreshReq {
  refresh_token: string
}

export interface ChangePasswordReq {
  oldPassword: string
  newPassword: string
}
