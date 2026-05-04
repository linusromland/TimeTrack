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
  FormControlLabel,
  Switch,
  Grid,
  Divider,
} from '@mui/material'
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  AccessTime as TimeIcon,
  FilterList as FilterIcon,
  Today as TodayIcon,
} from '@mui/icons-material'
import {
  useTimeEntries,
  useProjects,
  useCreateTimeEntry,
  useUpdateTimeEntry,
  useDeleteTimeEntry,
  useCreateProject,
} from '@/hooks/api'
import type { TimeEntry, Project } from '@/types'
import {
  formatDate,
  formatTime,
  formatDateTime,
  formatDuration,
  parseDuration,
  getToday,
  getCurrentTime,
  combineDateAndTime,
  addDuration,
  getDayRange,
  formatDateDisplay,
  formatTimeDisplay,
} from '@/utils/dateUtils'

interface TimeEntryFormData {
  project_id: string
  date: string
  start_time: string
  end_time: string
  duration_input: string
  note: string
  input_mode: 'datetime' | 'duration' // datetime+datetime vs datetime+duration
  new_project_name: string
}

export const TimeEntriesPage: React.FC = () => {
  const [open, setOpen] = useState(false)
  const [editingEntry, setEditingEntry] = useState<TimeEntry | null>(null)
  const [filterDate, setFilterDate] = useState<string>(getToday())
  const [showDateFilter, setShowDateFilter] = useState(true)
  const [showCreateProject, setShowCreateProject] = useState(false)
  
  const [formData, setFormData] = useState<TimeEntryFormData>({
    project_id: '',
    date: getToday(),
    start_time: getCurrentTime(),
    end_time: '',
    duration_input: '1h',
    note: '',
    input_mode: 'duration',
    new_project_name: '',
  })

  // Filter time entries by selected date
  const timeEntriesParams = useMemo(() => {
    if (showDateFilter && filterDate) {
      const { start, end } = getDayRange(filterDate)
      return { from: start, to: end }
    }
    return {}
  }, [showDateFilter, filterDate])

  const { data: timeEntries, isLoading: entriesLoading, error: entriesError } = useTimeEntries(timeEntriesParams)
  const { data: projects, refetch: refetchProjects } = useProjects()
  const createTimeEntry = useCreateTimeEntry()
  const updateTimeEntry = useUpdateTimeEntry()
  const deleteTimeEntry = useDeleteTimeEntry()
  const createProject = useCreateProject()

  const projectsMap = useMemo(() => {
    return (projects || []).reduce((acc, project) => {
      acc[project.id] = project
      return acc
    }, {} as Record<string, Project>)
  }, [projects])

  const handleOpen = (entry?: TimeEntry) => {
    if (entry) {
      setEditingEntry(entry)
      const startDate = new Date(entry.period.started)
      const endDate = new Date(entry.period.ended)
      
      setFormData({
        project_id: entry.project_id,
        date: formatDate(startDate),
        start_time: formatTime(startDate),
        end_time: formatTime(endDate),
        duration_input: formatDuration(entry.period.duration),
        note: entry.note,
        input_mode: 'datetime',
        new_project_name: '',
      })
    } else {
      setEditingEntry(null)
      setFormData({
        project_id: '',
        date: getToday(),
        start_time: getCurrentTime(),
        end_time: '',
        duration_input: '1h',
        note: '',
        input_mode: 'duration',
        new_project_name: '',
      })
    }
    setShowCreateProject(false)
    setOpen(true)
  }

  const handleClose = () => {
    setOpen(false)
    setEditingEntry(null)
    setShowCreateProject(false)
  }

  const handleInputModeChange = (mode: 'datetime' | 'duration') => {
    setFormData(prev => ({
      ...prev,
      input_mode: mode,
      end_time: mode === 'duration' ? '' : prev.end_time,
    }))
  }

  const calculateEndTime = (): string => {
    if (formData.input_mode === 'duration' && formData.duration_input) {
      const startDateTime = combineDateAndTime(formData.date, formData.start_time)
      const durationSeconds = parseDuration(formData.duration_input)
      const endDateTime = addDuration(startDateTime, durationSeconds)
      return formatTime(endDateTime)
    }
    return formData.end_time
  }

  const handleSubmit = async () => {
    // Create new project if needed
    let projectId = formData.project_id
    if (showCreateProject && formData.new_project_name.trim()) {
      try {
        const newProject = await createProject.mutateAsync({
          name: formData.new_project_name.trim()
        })
        projectId = newProject.id
        await refetchProjects()
      } catch (error) {
        console.error('Failed to create project:', error)
        return
      }
    }

    if (!projectId || !formData.date || !formData.start_time) return

    const startDateTime = combineDateAndTime(formData.date, formData.start_time)
    let endDateTime: string

    if (formData.input_mode === 'duration') {
      if (!formData.duration_input) return
      const durationSeconds = parseDuration(formData.duration_input)
      endDateTime = addDuration(startDateTime, durationSeconds)
    } else {
      if (!formData.end_time) return
      endDateTime = combineDateAndTime(formData.date, formData.end_time)
    }

    const data = {
      project_id: projectId,
      period: {
        start: new Date(startDateTime).toISOString(),
        end: new Date(endDateTime).toISOString(),
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

  const setToday = () => {
    setFilterDate(getToday())
  }

  const isFormValid = () => {
    if (showCreateProject) {
      return formData.new_project_name.trim() && formData.date && formData.start_time &&
        (formData.input_mode === 'duration' ? formData.duration_input : formData.end_time)
    }
    return formData.project_id && formData.date && formData.start_time &&
      (formData.input_mode === 'duration' ? formData.duration_input : formData.end_time)
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

      {/* Date Filter Controls */}
      <Paper sx={{ ...paperStyles, p: 2, mb: 3 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap' }}>
          <FilterIcon sx={{ color: '#94A3B8' }} />
          <Typography variant="h6" sx={{ color: '#F8FAFC', mr: 2 }}>
            Filter by Date
          </Typography>
          
          <FormControlLabel
            control={
              <Switch
                checked={showDateFilter}
                onChange={(e) => setShowDateFilter(e.target.checked)}
                sx={{
                  '& .MuiSwitch-switchBase.Mui-checked': { color: '#3B82F6' },
                  '& .MuiSwitch-switchBase.Mui-checked + .MuiSwitch-track': { backgroundColor: '#3B82F6' },
                }}
              />
            }
            label={<Typography sx={{ color: '#94A3B8' }}>Show specific date only</Typography>}
          />

          {showDateFilter && (
            <>
              <TextField
                type="date"
                value={filterDate}
                onChange={(e) => setFilterDate(e.target.value)}
                size="small"
                sx={{ minWidth: 150 }}
                InputLabelProps={{ shrink: true }}
              />
              <Button
                variant="outlined"
                startIcon={<TodayIcon />}
                onClick={setToday}
                size="small"
                sx={{
                  borderColor: '#3B82F6',
                  color: '#3B82F6',
                  '&:hover': { borderColor: '#60A5FA', backgroundColor: 'rgba(59, 130, 246, 0.1)' },
                }}
              >
                Today
              </Button>
              <Typography variant="body2" sx={{ color: '#94A3B8', fontStyle: 'italic' }}>
                Showing entries for {formatDateDisplay(filterDate)}
              </Typography>
            </>
          )}
        </Box>
      </Paper>

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
                <TableCell sx={{ color: '#94A3B8', fontWeight: 600 }}>Date</TableCell>
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
                  <TableCell colSpan={7} align="center">
                    <Box sx={{ py: 4 }}>
                      <CircularProgress sx={{ color: '#3B82F6' }} />
                    </Box>
                  </TableCell>
                </TableRow>
              ) : !timeEntries || timeEntries.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} align="center">
                    <Box sx={{ py: 4, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                      <TimeIcon sx={{ fontSize: 48, color: '#475569', mb: 2 }} />
                      <Typography color="textSecondary">
                        {showDateFilter 
                          ? `No time entries found for ${formatDateDisplay(filterDate)}.`
                          : 'No time entries found.'
                        } Click "Add Entry" to create your first time entry.
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
                    <TableCell sx={{ color: '#E2E8F0', fontWeight: 600 }}>
                      {formatDate(entry.period.started)}
                    </TableCell>
                    <TableCell sx={{ color: '#E2E8F0' }}>
                      {formatTimeDisplay(entry.period.started)}
                    </TableCell>
                    <TableCell sx={{ color: '#E2E8F0' }}>
                      {formatTimeDisplay(entry.period.ended)}
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
            
            {/* Project Selection or Creation */}
            <Box>
              <FormControlLabel
                control={
                  <Switch
                    checked={showCreateProject}
                    onChange={(e) => setShowCreateProject(e.target.checked)}
                    sx={{
                      '& .MuiSwitch-switchBase.Mui-checked': { color: '#10B981' },
                      '& .MuiSwitch-switchBase.Mui-checked + .MuiSwitch-track': { backgroundColor: '#10B981' },
                    }}
                  />
                }
                label={<Typography sx={{ color: '#E2E8F0' }}>Create new project</Typography>}
              />
              
              {showCreateProject ? (
                <TextField
                  label="New Project Name"
                  value={formData.new_project_name}
                  onChange={(e) => setFormData({ ...formData, new_project_name: e.target.value })}
                  fullWidth
                  required
                  sx={{ mt: 2 }}
                  helperText="Enter a name for the new project"
                />
              ) : (
                <TextField
                  select
                  label="Project"
                  value={formData.project_id}
                  onChange={(e) => setFormData({ ...formData, project_id: e.target.value })}
                  fullWidth
                  required
                  sx={{ mt: 2 }}
                >
                  {(projects || []).map((project) => (
                    <MenuItem key={project.id} value={project.id}>
                      {project.name}
                    </MenuItem>
                  ))}
                </TextField>
              )}
            </Box>

            <Divider sx={{ borderColor: 'rgba(148, 163, 184, 0.2)' }} />

            {/* Date */}
            <TextField
              label="Date"
              type="date"
              value={formData.date}
              onChange={(e) => setFormData({ ...formData, date: e.target.value })}
              fullWidth
              required
              InputLabelProps={{ shrink: true }}
            />

            {/* Time Input Mode Selection */}
            <Box>
              <Typography variant="subtitle2" sx={{ color: '#E2E8F0', mb: 1 }}>
                Time Input Method
              </Typography>
              <Box sx={{ display: 'flex', gap: 2 }}>
                <Button
                  variant={formData.input_mode === 'duration' ? 'contained' : 'outlined'}
                  onClick={() => handleInputModeChange('duration')}
                  sx={{
                    borderColor: '#8B5CF6',
                    color: formData.input_mode === 'duration' ? 'white' : '#8B5CF6',
                    backgroundColor: formData.input_mode === 'duration' ? '#8B5CF6' : 'transparent',
                  }}
                >
                  Start Time + Duration
                </Button>
                <Button
                  variant={formData.input_mode === 'datetime' ? 'contained' : 'outlined'}
                  onClick={() => handleInputModeChange('datetime')}
                  sx={{
                    borderColor: '#8B5CF6',
                    color: formData.input_mode === 'datetime' ? 'white' : '#8B5CF6',
                    backgroundColor: formData.input_mode === 'datetime' ? '#8B5CF6' : 'transparent',
                  }}
                >
                  Start Time + End Time
                </Button>
              </Box>
            </Box>

            {/* Time Inputs */}
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <TextField
                  label="Start Time"
                  type="time"
                  value={formData.start_time}
                  onChange={(e) => setFormData({ ...formData, start_time: e.target.value })}
                  fullWidth
                  required
                  InputLabelProps={{ shrink: true }}
                />
              </Grid>
              <Grid item xs={6}>
                {formData.input_mode === 'duration' ? (
                  <TextField
                    label="Duration"
                    value={formData.duration_input}
                    onChange={(e) => setFormData({ ...formData, duration_input: e.target.value })}
                    fullWidth
                    required
                    placeholder="e.g., 1h, 30m, 1h 30m, 1.5h"
                    helperText={`End time: ${calculateEndTime() || 'Invalid duration'}`}
                  />
                ) : (
                  <TextField
                    label="End Time"
                    type="time"
                    value={formData.end_time}
                    onChange={(e) => setFormData({ ...formData, end_time: e.target.value })}
                    fullWidth
                    required
                    InputLabelProps={{ shrink: true }}
                  />
                )}
              </Grid>
            </Grid>

            <TextField
              label="Note"
              value={formData.note}
              onChange={(e) => setFormData({ ...formData, note: e.target.value })}
              multiline
              rows={3}
              fullWidth
              placeholder="Optional description of work done..."
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
            disabled={!isFormValid()}
            sx={{ background: 'linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%)' }}
          >
            {editingEntry ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}