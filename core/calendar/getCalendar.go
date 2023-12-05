package calendar

import (
	"fmt"

	"google.golang.org/api/calendar/v3"
)

func GetCalendar(service *calendar.Service, calendarId string) *calendar.CalendarListEntry {
	calendar, err := service.CalendarList.Get(calendarId).Do()
	if err != nil {
		fmt.Printf("Unable to retrieve calendar: %v", err)
	}

	return calendar
}
