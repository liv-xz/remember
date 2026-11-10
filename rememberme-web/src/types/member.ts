export interface GraveMemberResp {
  id: number
  graveID: number
  userID: number
  username: string
  nickname: string
  avatar: string
  role: number
  relationDesc: string
  createdAt: string
  updatedAt: string
}

export interface GraveMemberListResp {
  items: GraveMemberResp[]
}

export interface AddMemberReq {
  userID: number
  role: number
  relationDesc?: string
}

export interface UpdateRoleReq {
  role: number
}

export interface UpdateRelationDescReq {
  relationDesc: string
}

export interface TransferOwnershipReq {
  targetMemberID: number
}

export interface TransferOwnershipResp {
  oldOwner: GraveMemberResp
  newOwner: GraveMemberResp
}
