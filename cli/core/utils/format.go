package utils

import "time"

func FormatDate(dateStr string, inputFormat string) string {
	date, _ := time.Parse(inputFormat, dateStr)
	return date.Format("2006-01-02 15:04")
}
