// Date and time utility functions for consistent formatting

/**
 * Format date as YYYY-MM-DD
 */
export const formatDate = (date: Date | string): string => {
  const d = new Date(date)
  return d.toISOString().split('T')[0]
}

/**
 * Format time as HH:MM
 */
export const formatTime = (date: Date | string): string => {
  const d = new Date(date)
  return d.toTimeString().slice(0, 5)
}

/**
 * Format date and time as YYYY-MM-DD HH:MM
 */
export const formatDateTime = (date: Date | string): string => {
  const d = new Date(date)
  return `${formatDate(d)} ${formatTime(d)}`
}

/**
 * Format duration from seconds to human readable format
 */
export const formatDuration = (seconds: number): string => {
  const hours = Math.floor(seconds / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${hours}h ${minutes}m`
}

/**
 * Parse duration string (like "1h 30m" or "2.5h") to seconds
 */
export const parseDuration = (durationStr: string): number => {
  const normalized = durationStr.toLowerCase().trim()
  
  // Handle formats like "1h 30m"
  const hoursMinutesMatch = normalized.match(/(\d+)h\s*(\d+)m/)
  if (hoursMinutesMatch) {
    const hours = parseInt(hoursMinutesMatch[1])
    const minutes = parseInt(hoursMinutesMatch[2])
    return hours * 3600 + minutes * 60
  }
  
  // Handle formats like "1h"
  const hoursOnlyMatch = normalized.match(/(\d+(?:\.\d+)?)h/)
  if (hoursOnlyMatch) {
    return Math.round(parseFloat(hoursOnlyMatch[1]) * 3600)
  }
  
  // Handle formats like "30m"
  const minutesOnlyMatch = normalized.match(/(\d+)m/)
  if (minutesOnlyMatch) {
    return parseInt(minutesOnlyMatch[1]) * 60
  }
  
  // Handle decimal hours like "1.5"
  const decimalMatch = normalized.match(/^(\d+(?:\.\d+)?)$/)
  if (decimalMatch) {
    return Math.round(parseFloat(decimalMatch[1]) * 3600)
  }
  
  return 0
}

/**
 * Get today's date as YYYY-MM-DD
 */
export const getToday = (): string => {
  return formatDate(new Date())
}

/**
 * Get current time as HH:MM
 */
export const getCurrentTime = (): string => {
  return formatTime(new Date())
}

/**
 * Combine date and time strings into ISO datetime
 */
export const combineDateAndTime = (date: string, time: string): string => {
  return `${date}T${time}:00`
}

/**
 * Add duration to a datetime string
 */
export const addDuration = (startDateTime: string, durationSeconds: number): string => {
  const start = new Date(startDateTime)
  const end = new Date(start.getTime() + durationSeconds * 1000)
  return end.toISOString()
}

/**
 * Get start and end of day for a given date
 */
export const getDayRange = (date: string): { start: string; end: string } => {
  return {
    start: `${date}T00:00:00`,
    end: `${date}T23:59:59`
  }
}

/**
 * Format date for display (e.g., "May 4, 2026")
 */
export const formatDateDisplay = (date: Date | string): string => {
  const d = new Date(date)
  return d.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })
}

/**
 * Format time for display (e.g., "9:00 AM")
 */
export const formatTimeDisplay = (date: Date | string): string => {
  const d = new Date(date)
  return d.toLocaleTimeString('en-US', { 
    hour: 'numeric', 
    minute: '2-digit',
    hour12: true
  })
}