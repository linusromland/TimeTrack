import React, { useState } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import {
  AppBar,
  Box,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography,
  useMediaQuery,
  useTheme,
  Avatar,
  Menu,
  MenuItem,
  Divider,
  Button,
} from '@mui/material'
import {
  Menu as MenuIcon,
  Dashboard as DashboardIcon,
  AccessTime as TimeIcon,
  Folder as ProjectIcon,
  Settings as SettingsIcon,
  AccountCircle as AccountIcon,
  Refresh as RefreshIcon,
  LogoutRounded as LogoutIcon,
  Schedule as ScheduleIcon,
} from '@mui/icons-material'
import { useAuth } from '@/contexts/AuthContext'

const drawerWidth = 280

const navigationItems = [
  { label: 'Dashboard', icon: <DashboardIcon />, path: '/dashboard' },
  { label: 'Time Entries', icon: <TimeIcon />, path: '/time-entries' },
  { label: 'Projects', icon: <ProjectIcon />, path: '/projects' },
  { label: 'Settings', icon: <SettingsIcon />, path: '/settings' },
]

export const Layout: React.FC = () => {
  const theme = useTheme()
  const isMobile = useMediaQuery(theme.breakpoints.down('md'))
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuth()

  const [mobileOpen, setMobileOpen] = useState(false)
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen)
  }

  const handleMenuClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget)
  }

  const handleMenuClose = () => {
    setAnchorEl(null)
  }

  const handleLogout = async () => {
    handleMenuClose()
    try {
      await logout()
      navigate('/login')
    } catch (error) {
      console.error('Logout failed:', error)
    }
  }

  const handleRefresh = () => {
    window.location.reload()
  }

  const drawer = (
    <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Toolbar sx={{ 
        px: 3, 
        py: 3,
        borderBottom: '1px solid rgba(59, 130, 246, 0.1)'
      }}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Box
            sx={{
              width: 40,
              height: 40,
              borderRadius: '12px',
              background: 'linear-gradient(135deg, #3B82F6 0%, #8B5CF6 100%)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              boxShadow: '0 4px 20px rgba(59, 130, 246, 0.3)',
            }}
          >
            <ScheduleIcon sx={{ color: 'white', fontSize: 24 }} />
          </Box>
          <Typography 
            variant="h6" 
            sx={{ 
              fontWeight: 700,
              background: 'linear-gradient(135deg, #3B82F6 0%, #8B5CF6 100%)',
              backgroundClip: 'text',
              WebkitBackgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
            }}
          >
            TimeTrack
          </Typography>
        </Box>
      </Toolbar>
      
      <Box sx={{ flex: 1, px: 2, py: 1 }}>
        <List sx={{ '& .MuiListItem-root': { px: 0 } }}>
          {navigationItems.map((item) => (
            <ListItem key={item.path} disablePadding sx={{ mb: 0.5 }}>
              <ListItemButton
                selected={location.pathname === item.path}
                onClick={() => {
                  navigate(item.path)
                  if (isMobile) {
                    setMobileOpen(false)
                  }
                }}
                sx={{
                  py: 1.5,
                  px: 2,
                  borderRadius: '12px',
                  mx: 1,
                }}
              >
                <ListItemIcon sx={{ 
                  color: location.pathname === item.path ? '#3B82F6' : '#94A3B8',
                  minWidth: 44,
                }}>
                  {item.icon}
                </ListItemIcon>
                <ListItemText 
                  primary={item.label} 
                  primaryTypographyProps={{
                    fontWeight: location.pathname === item.path ? 600 : 500,
                    color: location.pathname === item.path ? '#F8FAFC' : '#CBD5E1',
                  }}
                />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      </Box>
    </Box>
  )

  return (
    <Box sx={{ display: 'flex' }}>
      <AppBar
        position="fixed"
        elevation={0}
        sx={{
          width: { md: `calc(100% - ${drawerWidth}px)` },
          ml: { md: `${drawerWidth}px` },
        }}
      >
        <Toolbar sx={{ px: 3, py: 1 }}>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ 
              mr: 2, 
              display: { md: 'none' },
              background: 'rgba(59, 130, 246, 0.1)',
              backdropFilter: 'blur(10px)',
              '&:hover': {
                background: 'rgba(59, 130, 246, 0.2)',
              }
            }}
          >
            <MenuIcon />
          </IconButton>
          
          <Typography variant="h5" component="div" sx={{ flexGrow: 1, fontWeight: 600 }}>
            {navigationItems.find(item => item.path === location.pathname)?.label || 'TimeTrack'}
          </Typography>

          <Button
            variant="outlined"
            startIcon={<RefreshIcon />}
            onClick={handleRefresh}
            sx={{ 
              mr: 2,
              borderRadius: '10px',
            }}
          >
            Refresh
          </Button>

          <IconButton
            size="large"
            aria-label="account menu"
            aria-controls="account-menu"
            aria-haspopup="true"
            onClick={handleMenuClick}
            sx={{
              background: 'rgba(59, 130, 246, 0.1)',
              backdropFilter: 'blur(10px)',
              '&:hover': {
                background: 'rgba(59, 130, 246, 0.2)',
              }
            }}
          >
            <Avatar 
              sx={{ 
                width: 32, 
                height: 32,
                background: 'linear-gradient(135deg, #3B82F6 0%, #8B5CF6 100%)',
              }}
            >
              <AccountIcon />
            </Avatar>
          </IconButton>
          
          <Menu
            id="account-menu"
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
            onClick={handleMenuClose}
            PaperProps={{
              elevation: 0,
              sx: {
                overflow: 'visible',
                background: theme.glassMorphism.background,
                backdropFilter: theme.glassMorphism.backdropFilter,
                border: theme.glassMorphism.border,
                borderRadius: '12px',
                mt: 1.5,
                minWidth: 200,
                '& .MuiMenuItem-root': {
                  borderRadius: '8px',
                  mx: 1,
                  my: 0.5,
                  '&:hover': {
                    background: 'rgba(59, 130, 246, 0.1)',
                  },
                },
              },
            }}
            transformOrigin={{ horizontal: 'right', vertical: 'top' }}
            anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
          >
            <MenuItem disabled sx={{ opacity: 1 }}>
              <Typography variant="body2" sx={{ color: '#94A3B8', fontWeight: 500 }}>
                {user?.email}
              </Typography>
            </MenuItem>
            <Divider sx={{ my: 1, borderColor: 'rgba(59, 130, 246, 0.1)' }} />
            <MenuItem onClick={() => navigate('/settings')}>
              <ListItemIcon>
                <SettingsIcon fontSize="small" sx={{ color: '#CBD5E1' }} />
              </ListItemIcon>
              <Typography sx={{ color: '#F8FAFC' }}>Settings</Typography>
            </MenuItem>
            <MenuItem onClick={handleLogout}>
              <ListItemIcon>
                <LogoutIcon fontSize="small" sx={{ color: '#EF4444' }} />
              </ListItemIcon>
              <Typography sx={{ color: '#EF4444' }}>Logout</Typography>
            </MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      <Box
        component="nav"
        sx={{ width: { md: drawerWidth }, flexShrink: { md: 0 } }}
        aria-label="navigation"
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true,
          }}
          sx={{
            display: { xs: 'block', md: 'none' },
            '& .MuiDrawer-paper': { 
              boxSizing: 'border-box', 
              width: drawerWidth,
            },
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: 'none', md: 'block' },
            '& .MuiDrawer-paper': { 
              boxSizing: 'border-box', 
              width: drawerWidth,
            },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>

      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { md: `calc(100% - ${drawerWidth}px)` },
          minHeight: '100vh',
        }}
      >
        <Toolbar />
        <Outlet />
      </Box>
    </Box>
  )
}