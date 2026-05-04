import client from './client'

export interface TimePeriod {
  started: string
  ended: string
  duration: number
}

export interface TimeEntry {
  id: string
  project_id: string
  owner_id: string
  period: TimePeriod
  note: string
  created_at: string
  updated_at: string
}

export interface TimeEntryStatistics {
  total_entries: number
  total_time: number
  format: string
  entries_per_date: Array<{ timeframe: string; total_time: number }>
  entries_per_project: Array<{ project_id: string; total_time: number }>
}

export interface CreateTimeEntryInput {
  project_id: string
  period: { start: string; end: string }
  note?: string
}

export async function getTimeEntries(params?: {
  from?: string
  to?: string
  skip?: number
  limit?: number
}): Promise<TimeEntry[]> {
  const res = await client.get<TimeEntry[]>('/time-entries', { params })
  return res.data
}

export async function createTimeEntry(data: CreateTimeEntryInput): Promise<TimeEntry> {
  const res = await client.post<TimeEntry>('/time-entries', data)
  return res.data
}

export async function updateTimeEntry(id: string, data: Partial<CreateTimeEntryInput>): Promise<TimeEntry> {
  const res = await client.put<TimeEntry>(`/time-entries/${id}`, data)
  return res.data
}

export async function deleteTimeEntry(id: string): Promise<void> {
  await client.delete(`/time-entries/${id}`)
}

export async function getStatistics(params?: {
  from?: string
  to?: string
  format?: string
}): Promise<TimeEntryStatistics> {
  const res = await client.get<TimeEntryStatistics>('/time-entries/statistics', { params })
  return res.data
}
