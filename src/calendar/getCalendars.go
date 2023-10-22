package calendar

import (
	"log"

	"google.golang.org/api/calendar/v3"
)

func GetCalendars() *calendar.CalendarList {
	service := GetCalendarService()
	calendarList, err := service.CalendarList.List().ShowDeleted(false).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve calendar list: %v", err)
	}

	return calendarList
}
