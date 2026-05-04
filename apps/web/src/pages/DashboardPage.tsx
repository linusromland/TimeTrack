import React from 'react'
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Avatar,
  useTheme,
} from '@mui/material'
import {
  TrendingUp as TrendingUpIcon,
  Schedule as ScheduleIcon,
  Assessment as AssessmentIcon,
} from '@mui/icons-material'
import { LineChart } from '@mui/x-charts/LineChart'
import { useAuth } from '@/contexts/AuthContext'

const StatCard = ({ 
  title, 
  value, 
  change, 
  icon, 
  color,
  gradient 
}: {
  title: string
  value: string
  change?: string
  icon: React.ReactNode
  color: string
  gradient: string
}) => {
  const theme = useTheme()
  
  return (
    <Card sx={{ 
      height: '100%',
      position: 'relative',
      overflow: 'hidden',
    }}>
      <Box
        sx={{
          position: 'absolute',
          top: 0,
          right: 0,
          width: 80,
          height: 80,
          background: gradient,
          borderRadius: '50%',
          transform: 'translate(30%, -30%)',
          opacity: 0.1,
        }}
      />
      <CardContent sx={{ position: 'relative', zIndex: 1 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
          <Typography variant="body2" sx={{ color: '#94A3B8', fontWeight: 500 }}>
            {title}
          </Typography>
          <Avatar
            sx={{
              width: 48,
              height: 48,
              background: gradient,
              boxShadow: `0 4px 20px ${color}40`,
            }}
          >
            {icon}
          </Avatar>
        </Box>
        <Typography variant="h4" sx={{ fontWeight: 700, color: '#F8FAFC', mb: 1 }}>
          {value}
        </Typography>
        {change && (
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
            <TrendingUpIcon sx={{ fontSize: 16, color: '#10B981' }} />
            <Typography variant="body2" sx={{ color: '#10B981', fontWeight: 600 }}>
              {change}
            </Typography>
            <Typography variant="body2" sx={{ color: '#64748B' }}>
              from last week
            </Typography>
          </Box>
        )}
      </CardContent>
    </Card>
  )
}

export const DashboardPage: React.FC = () => {
  const { user } = useAuth()
  const theme = useTheme()

  const userName = user?.email?.split('@')[0] || 'User'

  return (
    <Box>
      {/* Welcome Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" sx={{ 
          mb: 1,
          background: 'linear-gradient(135deg, #3B82F6 0%, #8B5CF6 100%)',
          backgroundClip: 'text',
          WebkitBackgroundClip: 'text',
          WebkitTextFillColor: 'transparent',
          fontWeight: 700,
        }}>
          Welcome back, {userName}! 👋
        </Typography>
        <Typography variant="body1" sx={{ color: '#94A3B8' }}>
          Here's what's happening with your projects today.
        </Typography>
      </Box>

      {/* Stats Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={6}>
          <StatCard
            title="Today's Time"
            value="4h 32m"
            change="+2h 15m"
            icon={<ScheduleIcon />}
            color="#3B82F6"
            gradient="linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%)"
          />
        </Grid>
        <Grid item xs={12} sm={6}>
          <StatCard
            title="This Week"
            value="28h 45m"
            change="+8h 20m"
            icon={<AssessmentIcon />}
            color="#10B981"
            gradient="linear-gradient(135deg, #10B981 0%, #34D399 100%)"
          />
        </Grid>
      </Grid>

      <Grid container spacing={3}>
        {/* Weekly Time Tracking Chart */}
        <Grid item xs={12}>
          <Card sx={{ height: 400 }}>
            <CardContent>
              <Typography variant="h6" sx={{ mb: 3, color: '#F8FAFC', fontWeight: 600 }}>
                Weekly Time Overview
              </Typography>
              <Box sx={{ height: 300, width: '100%' }}>
                <LineChart
                  xAxis={[{
                    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
                    scaleType: 'point',
                  }]}
                  series={[{
                    data: [8.2, 6.5, 7.8, 9.1, 8.9, 4.2, 2.1],
                    label: 'Hours Worked',
                    color: '#3B82F6',
                  }]}
                  width={undefined}
                  height={300}
                  sx={{
                    '& .MuiChartsAxis-line': {
                      stroke: '#475569',
                    },
                    '& .MuiChartsAxis-tick': {
                      stroke: '#475569',
                    },
                    '& .MuiChartsAxis-tickLabel': {
                      fill: '#94A3B8',
                    },
                    '& .MuiChartsLegend-label': {
                      fill: '#F8FAFC',
                    },
                  }}
                />
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  )
}