package calendar

import (
	"fmt"

	"google.golang.org/api/calendar/v3"
)

func GetCalendars(service *calendar.Service) *calendar.CalendarList {
	calendarList, err := service.CalendarList.List().ShowDeleted(false).Do()
	if err != nil {
		fmt.Printf("Unable to retrieve calendar list: %v", err)
	}

	return calendarList
}
