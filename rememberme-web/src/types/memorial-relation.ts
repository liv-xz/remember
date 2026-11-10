export interface MemorialRelationResp {
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

export interface MemorialRelationListResp {
  items: MemorialRelationResp[]
}

export interface CreateMemorialRelationReq {
  deceasedID: number
  relationType: number
  note?: string
}
