package calendar

import (
	"time"

	"google.golang.org/api/calendar/v3"
)

func CreateEvent(calendarId string, title string, description string, startTime string, endTime string) *calendar.Event {
	service := GetCalendarService()

	timeZone := time.Now().Format("Z07:00")
	event := &calendar.Event{
		Summary:     title,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startTime,
			TimeZone: timeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: endTime,
			TimeZone: timeZone,
		},
	}

	event, err := service.Events.Insert(calendarId, event).Do()
	if err != nil {
		panic(err)
	}

	return event
}
