import { useState, useEffect, useCallback, type ReactNode, type MouseEvent } from 'react'
import {
  Box,
  Card,
  CardContent,
  Typography,
  Grid,
  ToggleButton,
  ToggleButtonGroup,
  Skeleton,
} from '@mui/material'
import AccessTimeIcon from '@mui/icons-material/AccessTime'
import ListAltIcon from '@mui/icons-material/ListAlt'
import FolderIcon from '@mui/icons-material/Folder'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  Legend,
} from 'recharts'
import { getStatistics, type TimeEntryStatistics } from '../api/timeEntries'
import { getProjects, type Project } from '../api/projects'

const COLORS = ['#0891b2', '#06b6d4', '#67e8f9', '#a5f3fc', '#e0f2fe', '#0e7490', '#155e75']

function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  if (h === 0) return `${m}m`
  return `${h}h ${m}m`
}

interface StatCardProps {
  icon: ReactNode
  label: string
  value: string
  color: string
  loading: boolean
}

function StatCard({ icon, label, value, color, loading }: StatCardProps) {
  return (
    <Card sx={{ borderRadius: 3, height: '100%' }}>
      <CardContent sx={{ p: 3 }}>
        <Box display="flex" alignItems="center" gap={2}>
          <Box
            sx={{
              width: 48,
              height: 48,
              borderRadius: 2,
              bgcolor: `${color}20`,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: color,
            }}
          >
            {icon}
          </Box>
          <Box>
            <Typography variant="body2" color="text.secondary">
              {label}
            </Typography>
            {loading ? (
              <Skeleton width={80} height={32} />
            ) : (
              <Typography variant="h5" fontWeight={700}>
                {value}
              </Typography>
            )}
          </Box>
        </Box>
      </CardContent>
    </Card>
  )
}

export default function DashboardPage() {
  const [range, setRange] = useState('30')
  const [stats, setStats] = useState<TimeEntryStatistics | null>(null)
  const [projects, setProjects] = useState<Project[]>([])
  const [loading, setLoading] = useState(true)

  const loadData = useCallback(async () => {
    setLoading(true)
    const to = new Date()
    const from = new Date()
    from.setDate(from.getDate() - parseInt(range))
    try {
      const [statsData, projectsData] = await Promise.all([
        getStatistics({
          from: from.toISOString(),
          to: to.toISOString(),
          format: range === '7' ? 'day' : range === '30' ? 'day' : 'week',
        }),
        getProjects(),
      ])
      setStats(statsData)
      setProjects(projectsData)
    } catch {
      // ignore
    } finally {
      setLoading(false)
    }
  }, [range])

  useEffect(() => {
    loadData()
  }, [loadData])

  const projectMap = Object.fromEntries(projects.map((p) => [p.id, p.name]))

  const barData = (stats?.entries_per_date ?? []).map((d) => ({
    date: d.timeframe,
    hours: +(d.total_time / 3600).toFixed(2),
  }))

  const pieData = (stats?.entries_per_project ?? []).map((d) => ({
    name: projectMap[d.project_id] ?? d.project_id,
    value: +(d.total_time / 3600).toFixed(2),
  }))

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h5" fontWeight={700}>
          Dashboard
        </Typography>
        <ToggleButtonGroup
          value={range}
          exclusive
          onChange={(_e: MouseEvent, v: string | null) => v && setRange(v)}
          size="small"
        >
          <ToggleButton value="7">7d</ToggleButton>
          <ToggleButton value="30">30d</ToggleButton>
          <ToggleButton value="90">90d</ToggleButton>
        </ToggleButtonGroup>
      </Box>

      <Grid container spacing={3} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={4}>
          <StatCard
            icon={<ListAltIcon />}
            label="Total Entries"
            value={String(stats?.total_entries ?? 0)}
            color="#0891b2"
            loading={loading}
          />
        </Grid>
        <Grid item xs={12} sm={4}>
          <StatCard
            icon={<AccessTimeIcon />}
            label="Total Time"
            value={formatDuration(stats?.total_time ?? 0)}
            color="#0e7490"
            loading={loading}
          />
        </Grid>
        <Grid item xs={12} sm={4}>
          <StatCard
            icon={<FolderIcon />}
            label="Active Projects"
            value={String(projects.length)}
            color="#155e75"
            loading={loading}
          />
        </Grid>
      </Grid>

      <Grid container spacing={3}>
        <Grid item xs={12} md={8}>
          <Card sx={{ borderRadius: 3 }}>
            <CardContent sx={{ p: 3 }}>
              <Typography variant="h6" fontWeight={600} mb={2}>
                Time per Day (hours)
              </Typography>
              {loading ? (
                <Skeleton variant="rectangular" height={260} />
              ) : (
                <ResponsiveContainer width="100%" height={260}>
                  <BarChart data={barData}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
                    <XAxis dataKey="date" tick={{ fontSize: 12 }} />
                    <YAxis tick={{ fontSize: 12 }} />
                    <Tooltip formatter={(v: number) => [`${v}h`, 'Hours']} />
                    <Bar dataKey="hours" fill="#0891b2" radius={[4, 4, 0, 0]} />
                  </BarChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} md={4}>
          <Card sx={{ borderRadius: 3, height: '100%' }}>
            <CardContent sx={{ p: 3 }}>
              <Typography variant="h6" fontWeight={600} mb={2}>
                Time per Project
              </Typography>
              {loading ? (
                <Skeleton variant="rectangular" height={260} />
              ) : pieData.length === 0 ? (
                <Box display="flex" alignItems="center" justifyContent="center" height={260}>
                  <Typography color="text.secondary">No data</Typography>
                </Box>
              ) : (
                <ResponsiveContainer width="100%" height={260}>
                  <PieChart>
                    <Pie
                      data={pieData}
                      cx="50%"
                      cy="45%"
                      outerRadius={80}
                      dataKey="value"
                      label={({ name, percent }: { name: string; percent: number }) =>
                        `${name} ${(percent * 100).toFixed(0)}%`
                      }
                      labelLine={false}
                    >
                      {pieData.map((_entry, i) => (
                        <Cell key={i} fill={COLORS[i % COLORS.length]} />
                      ))}
                    </Pie>
                    <Legend />
                    <Tooltip formatter={(v: number) => [`${v}h`, 'Hours']} />
                  </PieChart>
                </ResponsiveContainer>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  )
}
