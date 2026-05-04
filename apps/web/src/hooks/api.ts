import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { projectApi, timeEntryApi } from '@/services/api'
import type { Project, TimeEntry, TimeEntryStatistics } from '@/types'

// Project hooks
export const useProjects = (params?: {
  name?: string
  ids?: string[]
  skip?: number
  limit?: number
}) => {
  return useQuery({
    queryKey: ['projects', params],
    queryFn: () => projectApi.getProjects(params),
  })
}

export const useCreateProject = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: projectApi.createProject,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
  })
}

export const useUpdateProject = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) =>
      projectApi.updateProject(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
  })
}

export const useDeleteProject = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: projectApi.deleteProject,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] })
      queryClient.invalidateQueries({ queryKey: ['timeEntries'] })
    },
  })
}

// Time entry hooks
export const useTimeEntries = (params?: {
  from?: string
  to?: string
  skip?: number
  limit?: number
}) => {
  return useQuery({
    queryKey: ['timeEntries', params],
    queryFn: () => timeEntryApi.getTimeEntries(params),
  })
}

export const useCreateTimeEntry = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: timeEntryApi.createTimeEntry,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['timeEntries'] })
      queryClient.invalidateQueries({ queryKey: ['statistics'] })
    },
  })
}

export const useUpdateTimeEntry = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) =>
      timeEntryApi.updateTimeEntry(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['timeEntries'] })
      queryClient.invalidateQueries({ queryKey: ['statistics'] })
    },
  })
}

export const useDeleteTimeEntry = () => {
  const queryClient = useQueryClient()
  
  return useMutation({
    mutationFn: timeEntryApi.deleteTimeEntry,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['timeEntries'] })
      queryClient.invalidateQueries({ queryKey: ['statistics'] })
    },
  })
}

// Statistics hooks
export const useTimeEntryStatistics = (params?: {
  from?: string
  to?: string
  format?: 'd' | 'w' | 'm'
}) => {
  return useQuery({
    queryKey: ['statistics', params],
    queryFn: () => timeEntryApi.getStatistics(params),
  })
}