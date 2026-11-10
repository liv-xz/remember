import type { ApiResponse } from './auth'

export interface KinshipRelationResp {
  id: number
  graveID: number
  userID: number
  userName: string
  userAvatar: string
  deceasedID: number
  deceasedName: string
  deceasedAvatar: string
  relationType: number
  confirmationStatus: number
  confirmedBy?: number
  note: string
  createdAt: string
  updatedAt: string
}

export interface KinshipRelationListResp {
  items: KinshipRelationResp[]
}

export interface CreateKinshipRelationReq {
  deceasedID: number
  relationType: number
  note?: string
}
