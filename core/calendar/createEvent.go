package calendar

import (
	"time"

	"google.golang.org/api/calendar/v3"
)

func CreateEvent(calendarId string, title string, description string, startTime string, endTime string) *calendar.Event {
	service := GetCalendarService()

	// Removes the timezone from the time string
	parsedStartTime := startTime[:len(startTime)-6]
	parsedEndTime := endTime[:len(endTime)-6]

	timeZone := time.Now().Format("Z07:00")
	event := &calendar.Event{
		Summary:     title,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: parsedStartTime,
			TimeZone: timeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: parsedEndTime,
			TimeZone: timeZone,
		},
	}

	event, err := service.Events.Insert(calendarId, event).Do()
	if err != nil {
		panic(err)
	}

	return event
}
