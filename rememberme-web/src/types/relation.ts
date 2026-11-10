export interface DeceasedRelationResp {
  id: number
  graveID: number
  deceasedID1: number
  deceasedID2: number
  relationType: number
  direction: number
  createdAt: string
  updatedAt: string
}

export interface DeceasedRelationListResp {
  items: DeceasedRelationResp[]
}

export interface CreateDeceasedRelationReq {
  deceasedID1: number
  deceasedID2: number
  relationType: number
  direction: number
}
