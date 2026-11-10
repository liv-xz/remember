export interface MemorialEventResp {
  id: number
  graveID: number
  deceasedID: number
  eventType: number
  eventDate: string
  title: string
  lunar: number
  createdAt: string
  updatedAt: string
}

export interface MemorialEventListResp {
  items: MemorialEventResp[]
}

export interface CreateMemorialEventReq {
  deceasedID: number
  eventType: number
  eventDate: string
  title: string
  lunar?: number
}

export interface UpdateMemorialEventReq {
  deceasedID?: number
  eventType?: number
  eventDate?: string
  title?: string
  lunar?: number
}
