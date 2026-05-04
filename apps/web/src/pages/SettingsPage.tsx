import React from 'react'
import {
  Box,
  Typography,
  Paper,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Switch,
  Divider,
  Button,
} from '@mui/material'
import { useAuth } from '@/contexts/AuthContext'

export const SettingsPage: React.FC = () => {
  const { user } = useAuth()

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Settings
      </Typography>

      <Paper sx={{ mb: 3 }}>
        <List>
          <ListItem>
            <ListItemText 
              primary="Account Information"
              secondary={`Logged in as: ${user?.email}`}
            />
          </ListItem>
        </List>
      </Paper>

      <Paper sx={{ mb: 3 }}>
        <List>
          <ListItem>
            <ListItemText
              primary="Atlassian Integration"
              secondary="Connect your Jira account for automatic time logging"
            />
            <ListItemSecondaryAction>
              <Switch
                edge="end"
                checked={user?.integration?.atlassian?.enabled || false}
                disabled
              />
            </ListItemSecondaryAction>
          </ListItem>
          <Divider />
          <ListItem>
            <ListItemText
              primary="Connect to Jira"
              secondary="Set up integration with your Jira workspace"
            />
            <ListItemSecondaryAction>
              <Button variant="outlined" disabled>
                Connect
              </Button>
            </ListItemSecondaryAction>
          </ListItem>
        </List>
      </Paper>

      <Paper>
        <List>
          <ListItem>
            <ListItemText
              primary="Notifications"
              secondary="Receive notifications for time tracking reminders"
            />
            <ListItemSecondaryAction>
              <Switch
                edge="end"
                defaultChecked={false}
                disabled
              />
            </ListItemSecondaryAction>
          </ListItem>
          <Divider />
          <ListItem>
            <ListItemText
              primary="Auto-save"
              secondary="Automatically save time entries as you work"
            />
            <ListItemSecondaryAction>
              <Switch
                edge="end"
                defaultChecked={true}
                disabled
              />
            </ListItemSecondaryAction>
          </ListItem>
        </List>
      </Paper>
    </Box>
  )
}