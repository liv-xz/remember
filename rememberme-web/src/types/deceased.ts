export interface DeceasedResp {
  id: number
  graveId: number
  name: string
  gender: number
  birthDate: string
  deathDate: string
  biography: string
  portraitUrl: string
  epitaph: string
  burialOrder: number
  sortOrder: number
  originGraveID?: number
  createdAt: string
  updatedAt: string
}

export interface DeceasedListResp {
  items: DeceasedResp[]
}

export interface CreateDeceasedReq {
  name: string
  gender: number
  birthDate?: string
  deathDate?: string
  biography?: string
  portraitUrl?: string
  epitaph?: string
  burialOrder?: number
  sortOrder?: number
  originGraveID?: number
  userRelation?: {
    relationType: number
    note?: string
  }
}

export interface UpdateDeceasedReq {
  name?: string
  gender?: number
  birthDate?: string
  deathDate?: string
  biography?: string
  portraitUrl?: string
  epitaph?: string
  burialOrder?: number
  sortOrder?: number
}
