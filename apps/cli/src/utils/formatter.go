package utils

import "time"

func FormatDate(dateStr string, inputFormat string) string {
	date, _ := time.Parse(inputFormat, dateStr)
	return date.Local().Format("2006-01-02 15:04")
}

func FormatTime(t time.Time, format string) string {
	return t.Local().Format(format)
}
