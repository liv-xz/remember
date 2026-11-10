import type { ApiResponse } from './auth'

export interface ShareLinkResp {
  id: number
  graveID: number
  graveName: string
  shareCode: string
  shareURL: string
  expireAt?: string
  maxUseCount: number
  usedCount: number
  isActive: boolean
  createdAt: string
}

export interface ShareLinkListResp {
  items: ShareLinkResp[]
}

export interface DeceasedBrief {
  id: number
  name: string
  gender: number
}

export interface ShareLinkValidResp {
  graveID: number
  graveName: string
  graveDesc: string
  createdBy: number
  deceased: DeceasedBrief[]
}

export interface CreateShareLinkReq {
  expireDays?: number
}

export interface JoinShareLinkReq {
  deceasedID: number
  relationType: number
  note?: string
}