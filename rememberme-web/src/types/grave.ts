export interface GraveResp {
  id: number
  name: string
  graveType: number
  graveMode: 1 | 2
  layoutStyle: number
  description: string
  coverImage: string
  status: number
  latitude?: number
  longitude?: number
  address?: string
  createdBy: number
  creatorName: string
  createdAt: string
  updatedAt: string
}

export type GraveListResp = GraveResp[]

export interface CreateGraveReq {
  name: string
  graveType: number
  graveMode?: 1 | 2
  layoutStyle: number
  description?: string
  coverImage?: string
  latitude?: number
  longitude?: number
  address?: string
}

export interface UpdateGraveReq {
  name?: string
  graveType?: number
  graveMode?: 1 | 2
  layoutStyle?: number
  description?: string
  coverImage?: string
  status?: number
  latitude?: number
  longitude?: number
  address?: string
  clearLocation?: boolean
}

export interface CreateGraveWithDeceasedReq {
  name: string
  graveType: number
  graveMode?: 1 | 2
  layoutStyle: number
  description?: string
  coverImage?: string
  latitude?: number
  longitude?: number
  address?: string
  deceased: DeceasedItem[]
}

export interface DeceasedItem {
  name: string
  gender: number
  birthDate?: string
  deathDate?: string
  biography?: string
  portraitUrl?: string
  epitaph?: string
  burialOrder?: number
}

export interface GraveWithDeceasedResp {
  grave: GraveResp
  deceased: import('./deceased').DeceasedResp[]
}

export interface LatLng {
  latitude: number
  longitude: number
  address: string
}
