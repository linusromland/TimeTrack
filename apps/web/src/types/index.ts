export interface User {
  id: string
  email: string
  created_at: string
  updated_at: string
  integration?: {
    atlassian?: {
      enabled: boolean
    }
  }
}

export interface Project {
  id: string
  name: string
  owner_id: string
  integration?: {
    type: string
    key: string
    external_id: string
  }
  created_at: string
  updated_at: string
}

export interface TimePeriod {
  started: string
  ended: string
  duration: number
}

export interface ReportStatus {
  done: boolean
  integration: string
  external_id: string
  reported_at?: string
  updated_at?: string
}

export interface TimeEntry {
  id: string
  project_id: string
  owner_id: string
  period: TimePeriod
  note: string
  reported?: ReportStatus
  created_at: string
  updated_at: string
}

export interface CreateTimeEntryInput {
  project_id: string
  period: {
    start: string
    end: string
  }
  note?: string
}

export interface UpdateTimeEntryInput {
  project_id?: string
  period?: {
    start?: string
    end?: string
  }
  note?: string
}

export interface TimeEntryStatPerDate {
  timeframe: string
  total_time: number
}

export interface TimeEntryPerProject {
  project_id: string
  total_time: number
}

export interface TimeEntryStatistics {
  total_entries: number
  total_time: number
  format: string
  entries_per_date: TimeEntryStatPerDate[]
  entries_per_project: TimeEntryPerProject[]
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
}

export interface AuthResponse {
  user: User
}