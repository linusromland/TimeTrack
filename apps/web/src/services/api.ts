import {
  User,
  Project,
  TimeEntry,
  TimeEntryStatistics,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  CreateTimeEntryInput,
  UpdateTimeEntryInput,
} from '@/types'

const API_BASE_URL = '/api/v1'

class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
    public response?: any
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let errorMessage = `HTTP error! status: ${response.status}`
    try {
      const errorData = await response.json()
      errorMessage = errorData.error || errorMessage
    } catch {
      // If we can't parse JSON, use the default message
    }
    throw new ApiError(response.status, errorMessage)
  }
  
  const contentType = response.headers.get('content-type')
  if (contentType && contentType.includes('application/json')) {
    return response.json()
  }
  
  return response.text() as T
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`
  
  const config: RequestInit = {
    credentials: 'include', // Include cookies for httpOnly JWT
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  }
  
  const response = await fetch(url, config)
  return handleResponse<T>(response)
}

export const authApi = {
  async login(data: LoginRequest): Promise<AuthResponse> {
    return apiRequest<AuthResponse>('/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async register(data: RegisterRequest): Promise<AuthResponse> {
    return apiRequest<AuthResponse>('/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async logout(): Promise<void> {
    return apiRequest<void>('/logout', {
      method: 'POST',
    })
  },

  async getCurrentUser(): Promise<User> {
    return apiRequest<User>('/user/')
  },

  async getAtlassianOAuthUrl(): Promise<{ url: string }> {
    return apiRequest<{ url: string }>('/user/oauth/atlassian')
  },
}

export const projectApi = {
  async getProjects(params?: {
    name?: string
    ids?: string[]
    skip?: number
    limit?: number
  }): Promise<Project[]> {
    const searchParams = new URLSearchParams()
    if (params?.name) searchParams.append('name', params.name)
    if (params?.ids) searchParams.append('ids', params.ids.join(','))
    if (params?.skip) searchParams.append('skip', params.skip.toString())
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    
    const query = searchParams.toString()
    return apiRequest<Project[]>(`/projects${query ? `?${query}` : ''}`)
  },

  async createProject(data: { name: string; integration?: any }): Promise<Project> {
    return apiRequest<Project>('/projects', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async updateProject(id: string, data: { name?: string; integration?: any }): Promise<Project> {
    return apiRequest<Project>(`/projects/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },

  async deleteProject(id: string): Promise<void> {
    return apiRequest<void>(`/projects/${id}`, {
      method: 'DELETE',
    })
  },
}

export const timeEntryApi = {
  async getTimeEntries(params?: {
    from?: string
    to?: string
    skip?: number
    limit?: number
  }): Promise<TimeEntry[]> {
    const searchParams = new URLSearchParams()
    if (params?.from) searchParams.append('from', params.from)
    if (params?.to) searchParams.append('to', params.to)
    if (params?.skip) searchParams.append('skip', params.skip.toString())
    if (params?.limit) searchParams.append('limit', params.limit.toString())
    
    const query = searchParams.toString()
    return apiRequest<TimeEntry[]>(`/time-entries${query ? `?${query}` : ''}`)
  },

  async createTimeEntry(data: CreateTimeEntryInput): Promise<TimeEntry> {
    return apiRequest<TimeEntry>('/time-entries', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  async updateTimeEntry(id: string, data: UpdateTimeEntryInput): Promise<TimeEntry> {
    return apiRequest<TimeEntry>(`/time-entries/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },

  async deleteTimeEntry(id: string): Promise<void> {
    return apiRequest<void>(`/time-entries/${id}`, {
      method: 'DELETE',
    })
  },

  async getStatistics(params?: {
    from?: string
    to?: string
    format?: 'd' | 'w' | 'm'
  }): Promise<TimeEntryStatistics> {
    const searchParams = new URLSearchParams()
    if (params?.from) searchParams.append('from', params.from)
    if (params?.to) searchParams.append('to', params.to)
    if (params?.format) searchParams.append('format', params.format)
    
    const query = searchParams.toString()
    return apiRequest<TimeEntryStatistics>(`/time-entries/statistics${query ? `?${query}` : ''}`)
  },
}

export const healthApi = {
  async getHealth(): Promise<{ status: string }> {
    return apiRequest<{ status: string }>('/health')
  },
}

export { ApiError }