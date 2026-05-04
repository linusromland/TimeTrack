import { createTheme } from '@mui/material/styles'

declare module '@mui/material/styles' {
  interface Theme {
    glassMorphism: {
      background: string
      border: string
      backdropFilter: string
    }
  }

  interface ThemeOptions {
    glassMorphism?: {
      background: string
      border: string
      backdropFilter: string
    }
  }
}

export const theme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#3B82F6', // Blue
      light: '#60A5FA',
      dark: '#1E40AF',
    },
    secondary: {
      main: '#8B5CF6', // Purple
      light: '#A78BFA',
      dark: '#6D28D9',
    },
    background: {
      default: 'linear-gradient(135deg, #0F172A 0%, #1E293B 50%, #334155 100%)',
      paper: 'rgba(30, 41, 59, 0.8)',
    },
    error: {
      main: '#EF4444',
      light: '#F87171',
      dark: '#DC2626',
    },
    success: {
      main: '#10B981',
      light: '#34D399',
      dark: '#059669',
    },
    warning: {
      main: '#F59E0B',
      light: '#FBBF24',
      dark: '#D97706',
    },
    text: {
      primary: '#F8FAFC',
      secondary: '#CBD5E1',
    },
  },
  typography: {
    fontFamily: '"Inter", "SF Pro Display", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    h1: {
      fontWeight: 700,
      fontSize: '2.5rem',
      background: 'linear-gradient(135deg, #3B82F6 0%, #8B5CF6 100%)',
      backgroundClip: 'text',
      WebkitBackgroundClip: 'text',
      WebkitTextFillColor: 'transparent',
    },
    h4: {
      fontWeight: 600,
      fontSize: '1.75rem',
      color: '#F8FAFC',
    },
    h5: {
      fontWeight: 600,
      fontSize: '1.5rem',
      color: '#F8FAFC',
    },
    h6: {
      fontWeight: 600,
      fontSize: '1.25rem',
      color: '#F8FAFC',
    },
    body1: {
      fontSize: '1rem',
      color: '#CBD5E1',
    },
    body2: {
      fontSize: '0.875rem',
      color: '#94A3B8',
    },
  },
  components: {
    MuiCssBaseline: {
      styleOverrides: {
        body: {
          background: 'linear-gradient(135deg, #0F172A 0%, #1E293B 50%, #334155 100%)',
          backgroundAttachment: 'fixed',
          minHeight: '100vh',
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          background: 'rgba(30, 41, 59, 0.8)',
          backdropFilter: 'blur(20px)',
          border: '1px solid rgba(59, 130, 246, 0.1)',
          borderRadius: '16px',
          boxShadow: '0 8px 32px rgba(0, 0, 0, 0.3)',
          transition: 'all 0.3s ease-in-out',
          '&:hover': {
            transform: 'translateY(-2px)',
            boxShadow: '0 12px 40px rgba(59, 130, 246, 0.2)',
            border: '1px solid rgba(59, 130, 246, 0.3)',
          },
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          background: 'rgba(30, 41, 59, 0.8)',
          backdropFilter: 'blur(20px)',
          border: '1px solid rgba(59, 130, 246, 0.1)',
          borderRadius: '16px',
          boxShadow: '0 8px 32px rgba(0, 0, 0, 0.3)',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: '12px',
          textTransform: 'none',
          fontWeight: 600,
          padding: '12px 24px',
          fontSize: '0.95rem',
          transition: 'all 0.3s ease-in-out',
        },
        contained: {
          background: 'linear-gradient(135deg, #3B82F6 0%, #8B5CF6 100%)',
          boxShadow: '0 4px 20px rgba(59, 130, 246, 0.3)',
          '&:hover': {
            transform: 'translateY(-1px)',
            boxShadow: '0 8px 30px rgba(59, 130, 246, 0.4)',
            background: 'linear-gradient(135deg, #60A5FA 0%, #A78BFA 100%)',
          },
        },
        outlined: {
          border: '2px solid rgba(59, 130, 246, 0.3)',
          color: '#3B82F6',
          '&:hover': {
            border: '2px solid #3B82F6',
            background: 'rgba(59, 130, 246, 0.1)',
          },
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            borderRadius: '12px',
            background: 'rgba(51, 65, 85, 0.5)',
            backdropFilter: 'blur(10px)',
            '& fieldset': {
              border: '1px solid rgba(59, 130, 246, 0.2)',
            },
            '&:hover fieldset': {
              border: '1px solid rgba(59, 130, 246, 0.4)',
            },
            '&.Mui-focused fieldset': {
              border: '2px solid #3B82F6',
            },
          },
          '& .MuiInputLabel-root': {
            color: '#CBD5E1',
          },
          '& .MuiOutlinedInput-input': {
            color: '#F8FAFC',
          },
        },
      },
    },
    MuiDrawer: {
      styleOverrides: {
        paper: {
          background: 'rgba(15, 23, 42, 0.9)',
          backdropFilter: 'blur(20px)',
          border: 'none',
          borderRight: '1px solid rgba(59, 130, 246, 0.1)',
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          background: 'rgba(15, 23, 42, 0.9)',
          backdropFilter: 'blur(20px)',
          boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
          borderBottom: '1px solid rgba(59, 130, 246, 0.1)',
        },
      },
    },
    MuiListItemButton: {
      styleOverrides: {
        root: {
          borderRadius: '12px',
          margin: '4px 8px',
          transition: 'all 0.3s ease-in-out',
          '&:hover': {
            background: 'rgba(59, 130, 246, 0.1)',
            transform: 'translateX(4px)',
          },
          '&.Mui-selected': {
            background: 'linear-gradient(135deg, rgba(59, 130, 246, 0.2) 0%, rgba(139, 92, 246, 0.2) 100%)',
            borderLeft: '3px solid #3B82F6',
            '&:hover': {
              background: 'linear-gradient(135deg, rgba(59, 130, 246, 0.3) 0%, rgba(139, 92, 246, 0.3) 100%)',
            },
          },
        },
      },
    },
  },
  glassMorphism: {
    background: 'rgba(30, 41, 59, 0.8)',
    border: '1px solid rgba(59, 130, 246, 0.1)',
    backdropFilter: 'blur(20px)',
  },
})