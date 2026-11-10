export interface UserProfileResp {
  id: number
  username: string
  nickname: string
  avatar: string
  email: string
  abstract: string
}

export interface UpdateProfileReq {
  nickname: string
  avatar: string
  abstract: string
}
