package calendar

import (
	"fmt"

	"google.golang.org/api/calendar/v3"
)

func GetCalendars() *calendar.CalendarList {
	service := GetCalendarService()
	calendarList, err := service.CalendarList.List().ShowDeleted(false).Do()
	if err != nil {
		fmt.Printf("Unable to retrieve calendar list: %v", err)
	}

	return calendarList
}
