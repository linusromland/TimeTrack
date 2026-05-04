import client from './client'

export interface User {
  id: string
  email: string
  created_at: string
  updated_at: string
}

export async function login(email: string, password: string): Promise<string> {
  const res = await client.post<{ token: string }>('/login', { email, password })
  return res.data.token
}

export async function register(email: string, password: string): Promise<void> {
  await client.post('/register', { email, password })
}

export async function getUser(): Promise<User> {
  const res = await client.get<User>('/user/')
  return res.data
}
