package utils

import (
	"time"
)

func GetTaskTime(start time.Time, end time.Time) (hours int, minutes int) {
	startHour := start.Hour()
	startMinute := start.Minute()
	endHour := end.Hour()
	endMinute := end.Minute()

	if endMinute < startMinute {
		endMinute += 60
		endHour--
	}

	if endHour < startHour {
		endHour += 24
	}

	hours = endHour - startHour
	minutes = endMinute - startMinute

	return hours, minutes
}
