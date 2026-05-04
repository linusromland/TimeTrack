import client from './client'

export interface Project {
  id: string
  name: string
  integration: { type: string; key: string; external_id: string }
  owner_id: string
  created_at: string
  updated_at: string
}

export interface CreateProjectInput {
  name: string
  integration?: { type: string; key: string; external_id: string }
}

export async function getProjects(params?: { name?: string; skip?: number; limit?: number }): Promise<Project[]> {
  const res = await client.get<Project[]>('/projects', { params })
  return res.data
}

export async function createProject(data: CreateProjectInput): Promise<Project> {
  const res = await client.post<Project>('/projects', data)
  return res.data
}

export async function updateProject(id: string, data: Partial<CreateProjectInput>): Promise<Project> {
  const res = await client.put<Project>(`/projects/${id}`, data)
  return res.data
}

export async function deleteProject(id: string): Promise<void> {
  await client.delete(`/projects/${id}`)
}
