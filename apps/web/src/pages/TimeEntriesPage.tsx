import React, { useState, useMemo } from 'react'
import {
  Box,
  Typography,
  Paper,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  MenuItem,
  CircularProgress,
  Alert,
  Chip,
} from '@mui/material'
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  AccessTime as TimeIcon,
} from '@mui/icons-material'
import {
  useTimeEntries,
  useProjects,
  useCreateTimeEntry,
  useUpdateTimeEntry,
  useDeleteTimeEntry,
} from '@/hooks/api'
import type { TimeEntry, Project } from '@/types'

interface TimeEntryFormData {
  project_id: string
  start: string
  end: string
  note: string
}

export const TimeEntriesPage: React.FC = () => {
  const [open, setOpen] = useState(false)
  const [editingEntry, setEditingEntry] = useState<TimeEntry | null>(null)
  const [formData, setFormData] = useState<TimeEntryFormData>({
    project_id: '',
    start: '',
    end: '',
    note: '',
  })

  const { data: timeEntries, isLoading: entriesLoading, error: entriesError } = useTimeEntries()
  const { data: projects } = useProjects()
  const createTimeEntry = useCreateTimeEntry()
  const updateTimeEntry = useUpdateTimeEntry()
  const deleteTimeEntry = useDeleteTimeEntry()

  const projectsMap = useMemo(() => {
    return (projects || []).reduce((acc, project) => {
      acc[project.id] = project
      return acc
    }, {} as Record<string, Project>)
  }, [projects])

  const formatDateTime = (dateString: string) => {
    return new Date(dateString).toLocaleString()
  }

  const formatDuration = (duration: number) => {
    const hours = Math.floor(duration / 3600)
    const minutes = Math.floor((duration % 3600) / 60)
    return hours + 'h ' + minutes + 'm'
  }

  const handleOpen = (entry?: TimeEntry) => {
    if (entry) {
      setEditingEntry(entry)
      setFormData({
        project_id: entry.project_id,
        start: new Date(entry.period.started).toISOString().slice(0, 16),
        end: new Date(entry.period.ended).toISOString().slice(0, 16),
        note: entry.note,
      })
    } else {
      setEditingEntry(null)
      setFormData({
        project_id: '',
        start: '',
        end: '',
        note: '',
      })
    }
    setOpen(true)
  }

  const handleClose = () => {
    setOpen(false)
    setEditingEntry(null)
  }

  const handleSubmit = async () => {
    if (!formData.start || !formData.end || !formData.project_id) return

    const data = {
      project_id: formData.project_id,
      period: {
        start: new Date(formData.start).toISOString(),
        end: new Date(formData.end).toISOString(),
      },
      note: formData.note,
    }

    try {
      if (editingEntry) {
        await updateTimeEntry.mutateAsync({ id: editingEntry.id, data })
      } else {
        await createTimeEntry.mutateAsync(data)
      }
      handleClose()
    } catch (error) {
      console.error('Failed to save time entry:', error)
    }
  }

  const handleDelete = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this time entry?')) {
      try {
        await deleteTimeEntry.mutateAsync(id)
      } catch (error) {
        console.error('Failed to delete time entry:', error)
      }
    }
  }

  const paperStyles = {
    p: 0,
    background: 'rgba(15, 23, 42, 0.8)',
    backdropFilter: 'blur(10px)',
    border: '1px solid rgba(148, 163, 184, 0.1)',
    borderRadius: 2,
  }

  const dialogStyles = {
    background: 'rgba(15, 23, 42, 0.95)',
    backdropFilter: 'blur(10px)',
    border: '1px solid rgba(148, 163, 184, 0.2)',
  }

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" sx={{ color: '#F8FAFC', fontWeight: 600 }}>
          Time Entries
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => handleOpen()}
          sx={{
            background: 'linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%)',
            boxShadow: '0 8px 32px rgba(59, 130, 246, 0.3)',
          }}
        >
          Add Entry
        </Button>
      </Box>

      {entriesError && (
        <Alert 
          severity="error" 
          sx={{ 
            mb: 3,
            background: 'rgba(239, 68, 68, 0.1)',
            border: '1px solid rgba(239, 68, 68, 0.2)',
            color: '#F87171',
          }}
        >
          Failed to load time entries. Please check your connection.
        </Alert>
      )}

      <Paper sx={paperStyles}>
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell sx={{ color: '#94A3B8', fontWeight: 600 }}>Project</TableCell>
                <TableCell sx={{ color: '#94A3B8', fontWeight: 600 }}>Start</TableCell>
                <TableCell sx={{ color: '#94A3B8', fontWeight: 600 }}>End</TableCell>
                <TableCell sx={{ color: '#94A3B8', fontWeight: 600 }}>Duration</TableCell>
                <TableCell sx={{ color: '#94A3B8', fontWeight: 600 }}>Note</TableCell>
                <TableCell sx={{ color: '#94A3B8', fontWeight: 600 }}>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {entriesLoading ? (
                <TableRow>
                  <TableCell colSpan={6} align="center">
                    <Box sx={{ py: 4 }}>
                      <CircularProgress sx={{ color: '#3B82F6' }} />
                    </Box>
                  </TableCell>
                </TableRow>
              ) : !timeEntries || timeEntries.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} align="center">
                    <Box sx={{ py: 4, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                      <TimeIcon sx={{ fontSize: 48, color: '#475569', mb: 2 }} />
                      <Typography color="textSecondary">
                        No time entries found. Click "Add Entry" to create your first time entry.
                      </Typography>
                    </Box>
                  </TableCell>
                </TableRow>
              ) : (
                timeEntries.map((entry) => (
                  <TableRow key={entry.id} sx={{ '&:hover': { bgcolor: 'rgba(59, 130, 246, 0.1)' } }}>
                    <TableCell>
                      <Chip
                        label={projectsMap[entry.project_id]?.name || 'Unknown Project'}
                        sx={{
                          background: 'linear-gradient(135deg, #10B981 0%, #34D399 100%)',
                          color: 'white',
                          fontWeight: 500,
                        }}
                      />
                    </TableCell>
                    <TableCell sx={{ color: '#E2E8F0' }}>
                      {formatDateTime(entry.period.started)}
                    </TableCell>
                    <TableCell sx={{ color: '#E2E8F0' }}>
                      {formatDateTime(entry.period.ended)}
                    </TableCell>
                    <TableCell sx={{ color: '#E2E8F0', fontWeight: 600 }}>
                      {formatDuration(entry.period.duration)}
                    </TableCell>
                    <TableCell sx={{ color: '#E2E8F0' }}>
                      {entry.note || '-'}
                    </TableCell>
                    <TableCell>
                      <IconButton
                        size="small"
                        onClick={() => handleOpen(entry)}
                        sx={{ color: '#3B82F6', mr: 1 }}
                      >
                        <EditIcon fontSize="small" />
                      </IconButton>
                      <IconButton
                        size="small"
                        onClick={() => handleDelete(entry.id)}
                        sx={{ color: '#EF4444' }}
                      >
                        <DeleteIcon fontSize="small" />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Paper>

      <Dialog 
        open={open} 
        onClose={handleClose}
        maxWidth="md"
        fullWidth
        PaperProps={{ sx: dialogStyles }}
      >
        <DialogTitle sx={{ color: '#F8FAFC' }}>
          {editingEntry ? 'Edit Time Entry' : 'Add Time Entry'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2, display: 'flex', flexDirection: 'column', gap: 3 }}>
            <TextField
              select
              label="Project"
              value={formData.project_id}
              onChange={(e) => setFormData({ ...formData, project_id: e.target.value })}
              fullWidth
              required
            >
              {(projects || []).map((project) => (
                <MenuItem key={project.id} value={project.id}>
                  {project.name}
                </MenuItem>
              ))}
            </TextField>
            <TextField
              label="Start Time"
              type="datetime-local"
              value={formData.start}
              onChange={(e) => setFormData({ ...formData, start: e.target.value })}
              fullWidth
              required
              InputLabelProps={{ shrink: true }}
            />
            <TextField
              label="End Time"
              type="datetime-local"
              value={formData.end}
              onChange={(e) => setFormData({ ...formData, end: e.target.value })}
              fullWidth
              required
              InputLabelProps={{ shrink: true }}
            />
            <TextField
              label="Note"
              value={formData.note}
              onChange={(e) => setFormData({ ...formData, note: e.target.value })}
              multiline
              rows={3}
              fullWidth
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose} sx={{ color: '#94A3B8' }}>
            Cancel
          </Button>
          <Button
            onClick={handleSubmit}
            variant="contained"
            disabled={!formData.start || !formData.end || !formData.project_id}
            sx={{ background: 'linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%)' }}
          >
            {editingEntry ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}