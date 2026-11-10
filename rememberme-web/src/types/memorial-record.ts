export interface MemorialRecordResp {
  id: number
  graveID: number
  deceasedID?: number
  userID: number
  type: number
  message: string
  createdAt: string
  updatedAt: string
}

export interface MemorialRecordListResp {
  items: MemorialRecordResp[]
}
