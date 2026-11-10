import client from './client'
import type { ApiResponse } from '../types/auth'
import type {
  MemorialEventResp,
  MemorialEventListResp,
  CreateMemorialEventReq,
  UpdateMemorialEventReq,
} from '../types/memorial-event'

export async function createEvent(
  graveId: number,
  data: CreateMemorialEventReq
): Promise<ApiResponse<MemorialEventResp>> {
  return client.post<ApiResponse<MemorialEventResp>>(`/graves/${graveId}/events`, data)
}

export async function getEventsByGrave(
  graveId: number
): Promise<ApiResponse<MemorialEventListResp>> {
  return client.get<ApiResponse<MemorialEventListResp>>(`/graves/${graveId}/events`)
}

export async function updateEvent(
  eventId: number,
  data: UpdateMemorialEventReq
): Promise<ApiResponse<MemorialEventResp>> {
  return client.put<ApiResponse<MemorialEventResp>>(`/events/${eventId}`, data)
}

export async function deleteEvent(eventId: number): Promise<ApiResponse> {
  return client.delete<ApiResponse>(`/events/${eventId}`)
}
