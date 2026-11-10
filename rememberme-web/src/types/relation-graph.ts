import type { ApiResponse } from './auth'

export interface GraphNode {
  id: string
  type: 'deceased' | 'user'
  name: string
  avatar: string
  gender: number
  generation: number
  isDead: boolean
  birthDate?: string
  deathDate?: string
}

export interface GraphEdge {
  from: string
  to: string
  relationType: string
  confirmed: boolean
  direction?: number
  confirmationStatus?: number
  confirmedBy?: number
}

export interface RelationGraphResp {
  graveID: number
  graveName: string
  nodes: GraphNode[]
  edges: GraphEdge[]
}