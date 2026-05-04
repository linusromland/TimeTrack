import React, { useState } from 'react'
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
  CircularProgress,
  Alert,
  Chip,
  useTheme,
} from '@mui/material'
import {
  Add as AddIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  FolderOpen as ProjectIcon,
  IntegrationInstructions as IntegrationIcon,
} from '@mui/icons-material'
import {
  useProjects,
  useCreateProject,
  useUpdateProject,
  useDeleteProject,
} from '@/hooks/api'
import type { Project } from '@/types'
import { formatDateDisplay } from '@/utils/dateUtils'

interface ProjectFormData {
  name: string
}

export const ProjectsPage: React.FC = () => {
  const theme = useTheme()
  const [open, setOpen] = useState(false)
  const [editingProject, setEditingProject] = useState<Project | null>(null)
  const [formData, setFormData] = useState<ProjectFormData>({
    name: '',
  })

  const { data: projects, isLoading: projectsLoading, error: projectsError } = useProjects()
  const createProject = useCreateProject()
  const updateProject = useUpdateProject()
  const deleteProject = useDeleteProject()

  const formatDate = (dateString: string) => {
    return formatDateDisplay(dateString)
  }

  const getIntegrationType = (project: Project) => {
    if (!project.integration) return 'None'
    return project.integration.type.charAt(0).toUpperCase() + project.integration.type.slice(1)
  }

  const handleOpen = (project?: Project) => {
    if (project) {
      setEditingProject(project)
      setFormData({
        name: project.name,
      })
    } else {
      setEditingProject(null)
      setFormData({
        name: '',
      })
    }
    setOpen(true)
  }

  const handleClose = () => {
    setOpen(false)
    setEditingProject(null)
  }

  const handleSubmit = async () => {
    if (!formData.name.trim()) return

    const data = {
      name: formData.name.trim(),
    }

    try {
      if (editingProject) {
        await updateProject.mutateAsync({ id: editingProject.id, data })
      } else {
        await createProject.mutateAsync(data)
      }
      handleClose()
    } catch (error) {
      console.error('Failed to save project:', error)
    }
  }

  const handleDelete = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this project? This will also delete all associated time entries.')) {
      try {
        await deleteProject.mutateAsync(id)
      } catch (error) {
        console.error('Failed to delete project:', error)
      }
    }
  }

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" sx={{ color: '#F8FAFC', fontWeight: 600 }}>
          Projects
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => handleOpen()}
          sx={{
            background: 'linear-gradient(135deg, #10B981 0%, #34D399 100%)',
            boxShadow: '0 8px 32px rgba(16, 185, 129, 0.3)',
          }}
        >
          Add Project
        </Button>
      </Box>

      {projectsError && (
        <Alert 
          severity="error" 
          sx={{ 
            mb: 3,
            background: 'rgba(239, 68, 68, 0.1)',
            border: '1px solid rgba(239, 68, 68, 0.2)',
            color: '#F87171',
          }}
        >
          Failed to load projects. Please check your connection.
        </Alert>
      )}

      <Paper 
        sx={{ 
          p: 0,
          background: 'rgba(15, 23, 42, 0.8)',
          backdropFilter: 'blur(10px)',
          border: '1px solid rgba(148, 163, 184, 0.1)',
          borderRadius: 2,
        }}
      >
        <TableContainer>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell sx={{ color: '#94A3B8', fontWeight: 600 }}>Name</TableCell>
                <TableCell sx={{ color: '#94A3B8', fontWeight: 600 }}>Integration</TableCell>
                <TableCell sx={{ color: '#94A3B8', fontWeight: 600 }}>Created</TableCell>
                <TableCell sx={{ color: '#94A3B8', fontWeight: 600 }}>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {projectsLoading ? (
                <TableRow>
                  <TableCell colSpan={4} align="center">
                    <Box sx={{ py: 4 }}>
                      <CircularProgress sx={{ color: '#3B82F6' }} />
                    </Box>
                  </TableCell>
                </TableRow>
              ) : !projects || projects.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={4} align="center">
                    <Box sx={{ py: 4, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                      <ProjectIcon sx={{ fontSize: 48, color: '#475569', mb: 2 }} />
                      <Typography color="textSecondary">
                        No projects found. Click "Add Project" to create your first project.
                      </Typography>
                    </Box>
                  </TableCell>
                </TableRow>
              ) : (
                projects.map((project) => (
                  <TableRow key={project.id} sx={{ '&:hover': { bgcolor: 'rgba(16, 185, 129, 0.1)' } }}>
                    <TableCell>
                      <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                        <ProjectIcon sx={{ color: '#10B981' }} />
                        <Typography sx={{ color: '#F8FAFC', fontWeight: 600 }}>
                          {project.name}
                        </Typography>
                      </Box>
                    </TableCell>
                    <TableCell>
                      <Chip
                        icon={<IntegrationIcon />}
                        label={getIntegrationType(project)}
                        size="small"
                        sx={{
                          background: project.integration 
                            ? 'linear-gradient(135deg, #8B5CF6 0%, #A78BFA 100%)'
                            : 'rgba(71, 85, 105, 0.3)',
                          color: 'white',
                          border: 'none',
                        }}
                      />
                    </TableCell>
                    <TableCell sx={{ color: '#E2E8F0' }}>
                      {formatDate(project.created_at)}
                    </TableCell>
                    <TableCell>
                      <IconButton
                        size="small"
                        onClick={() => handleOpen(project)}
                        sx={{ color: '#3B82F6', mr: 1 }}
                      >
                        <EditIcon fontSize="small" />
                      </IconButton>
                      <IconButton
                        size="small"
                        onClick={() => handleDelete(project.id)}
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

      {/* Add/Edit Dialog */}
      <Dialog 
        open={open} 
        onClose={handleClose}
        maxWidth="sm"
        fullWidth
        PaperProps={{
          sx: {
            background: 'rgba(15, 23, 42, 0.95)',
            backdropFilter: 'blur(10px)',
            border: '1px solid rgba(148, 163, 184, 0.2)',
          },
        }}
      >
        <DialogTitle sx={{ color: '#F8FAFC' }}>
          {editingProject ? 'Edit Project' : 'Add Project'}
        </DialogTitle>
        <DialogContent>
          <Box sx={{ pt: 2 }}>
            <TextField
              label="Project Name"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              fullWidth
              required
              helperText="Enter a descriptive name for your project"
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
            disabled={!formData.name.trim()}
            sx={{
              background: 'linear-gradient(135deg, #10B981 0%, #34D399 100%)',
            }}
          >
            {editingProject ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  )
}