package utils

import "time"

func IsValidTime(timeStr string) bool {
	var timeFormat = "15:04"
	_, err := time.Parse(timeFormat, timeStr)
	return err == nil
}

func IsValidDate(dateStr string) bool {
	var dateFormat = "2006-01-02"
	_, err := time.Parse(dateFormat, dateStr)
	return err == nil
}
