import client from './client'
import type { ApiResponse } from '../types/auth'
import type { RelationGraphResp } from '../types/relation-graph'

export async function getRelationGraph(graveId: number): Promise<ApiResponse<RelationGraphResp>> {
  return client.get<ApiResponse<RelationGraphResp>>(`/graves/${graveId}/relation-graph`)
}