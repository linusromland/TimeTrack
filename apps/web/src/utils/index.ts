export * from './dateUtils'

/**
 * Get the start and end of this week in ISO format
 */
export function getThisWeekRange(): { start: string; end: string } {
  const now = new Date()
  const dayOfWeek = now.getDay()
  const start = new Date(now.getFullYear(), now.getMonth(), now.getDate() - dayOfWeek)
  const end = new Date(now.getFullYear(), now.getMonth(), now.getDate() + (7 - dayOfWeek))
  
  return {
    start: start.toISOString(),
    end: end.toISOString()
  }
}

/**
 * Get the start and end of this month in ISO format
 */
export function getThisMonthRange(): { start: string; end: string } {
  const now = new Date()
  const start = new Date(now.getFullYear(), now.getMonth(), 1)
  const end = new Date(now.getFullYear(), now.getMonth() + 1, 1)
  
  return {
    start: start.toISOString(),
    end: end.toISOString()
  }
}

/**
 * Validate email format
 */
export function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * Calculate duration between two date strings in seconds
 */
export function calculateDuration(startDate: string, endDate: string): number {
  const start = new Date(startDate)
  const end = new Date(endDate)
  return Math.max(0, Math.floor((end.getTime() - start.getTime()) / 1000))
}

/**
 * Convert local date and time inputs to ISO string
 */
export function combineDateTime(date: string, time: string): string {
  const dateTime = new Date(`${date}T${time}`)
  return dateTime.toISOString()
}

/**
 * Extract date part from ISO string for date input
 */
export function getDatePart(isoString: string): string {
  return isoString.split('T')[0]
}

/**
 * Extract time part from ISO string for time input  
 */
export function getTimePart(isoString: string): string {
  const date = new Date(isoString)
  return date.toTimeString().slice(0, 5) // HH:MM
}

/**
 * Debounce function calls
 */
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: NodeJS.Timeout
  
  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => func(...args), delay)
  }
}