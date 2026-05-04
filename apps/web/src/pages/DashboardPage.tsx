import React, { useMemo } from 'react'
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  Avatar,
  CircularProgress,
  Alert,
  useTheme,
} from '@mui/material'
import {
  TrendingUp as TrendingUpIcon,
  Schedule as ScheduleIcon,
  Assessment as AssessmentIcon,
} from '@mui/icons-material'
import { LineChart } from '@mui/x-charts/LineChart'
import { useAuth } from '@/contexts/AuthContext'
import { useTimeEntryStatistics } from '@/hooks/api'
import { formatDate, formatDuration } from '@/utils/dateUtils'

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

  // Get current date ranges
  const today = new Date()
  const todayStr = formatDate(today)
  
  const weekStart = new Date(today)
  weekStart.setDate(today.getDate() - today.getDay()) // Start of current week (Sunday)
  const weekStartStr = formatDate(weekStart)

  // Fetch today's statistics
  const { 
    data: todayStats, 
    isLoading: todayLoading, 
    error: todayError 
  } = useTimeEntryStatistics({
    from: todayStr,
    to: todayStr,
    format: 'd'
  })

  // Fetch this week's statistics  
  const { 
    data: weekStats, 
    isLoading: weekLoading, 
    error: weekError 
  } = useTimeEntryStatistics({
    from: weekStartStr,
    to: todayStr,
    format: 'd'
  })

  // Calculate display values
  const { todayTime, weekTime, chartData } = useMemo(() => {
    const todayTime = todayStats?.total_time || 0
    const weekTime = weekStats?.total_time || 0
    
    // Prepare chart data
    const daysOfWeek = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
    const chartData = daysOfWeek.map((day, index) => {
      const date = new Date(weekStart)
      date.setDate(weekStart.getDate() + index)
      const dateStr = formatDate(date)
      
      const dayData = weekStats?.entries_per_date?.find(
        entry => entry.timeframe === dateStr
      )
      
      return dayData ? Math.round((dayData.total_time / 3600) * 10) / 10 : 0 // Convert seconds to hours
    })

    return {
      todayTime: Math.round((todayTime / 3600) * 10) / 10, // Convert seconds to hours
      weekTime: Math.round((weekTime / 3600) * 10) / 10,
      chartData
    }
  }, [todayStats, weekStats, weekStart])

  const formatTime = (hours: number): string => {
    const totalSeconds = Math.round(hours * 3600)
    return formatDuration(totalSeconds)
  }

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
            value={todayLoading ? '...' : formatTime(todayTime)}
            change={todayError ? undefined : "+0h 0m"}
            icon={<ScheduleIcon />}
            color="#3B82F6"
            gradient="linear-gradient(135deg, #3B82F6 0%, #60A5FA 100%)"
          />
        </Grid>
        <Grid item xs={12} sm={6}>
          <StatCard
            title="This Week"
            value={weekLoading ? '...' : formatTime(weekTime)}
            change={weekError ? undefined : "+0h 0m"}
            icon={<AssessmentIcon />}
            color="#10B981"
            gradient="linear-gradient(135deg, #10B981 0%, #34D399 100%)"
          />
        </Grid>
      </Grid>

      <Grid container spacing={3}>
        {/* Error Display */}
        {(todayError || weekError) && (
          <Grid item xs={12}>
            <Alert 
              severity="error" 
              sx={{ 
                mb: 2,
                background: 'rgba(239, 68, 68, 0.1)',
                border: '1px solid rgba(239, 68, 68, 0.2)',
                color: '#F87171',
              }}
            >
              Failed to load time tracking data. Please check your connection.
            </Alert>
          </Grid>
        )}

        {/* Weekly Time Tracking Chart */}
        <Grid item xs={12}>
          <Card sx={{ height: 400 }}>
            <CardContent>
              <Typography variant="h6" sx={{ mb: 3, color: '#F8FAFC', fontWeight: 600 }}>
                Weekly Time Overview
              </Typography>
              {weekLoading ? (
                <Box 
                  sx={{ 
                    height: 300, 
                    display: 'flex', 
                    alignItems: 'center', 
                    justifyContent: 'center' 
                  }}
                >
                  <CircularProgress sx={{ color: '#3B82F6' }} />
                </Box>
              ) : (
                <Box sx={{ height: 300, width: '100%' }}>
                  <LineChart
                    xAxis={[{
                      data: ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'],
                      scaleType: 'point',
                    }]}
                    series={[{
                      data: chartData,
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
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  )
}