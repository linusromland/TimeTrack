package calendar

import "google.golang.org/api/calendar/v3"

func GetEvents(service *calendar.Service, calendarId string, startTime string, endTime string) *calendar.Events {
	events, err := service.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin(startTime).TimeMax(endTime).OrderBy("startTime").MaxResults(2500).Do()

	if err != nil {
		panic(err)
	}

	// Remove all events that are full day events
	for i := 0; i < len(events.Items); i++ {
		if events.Items[i].Start.DateTime == "" {
			events.Items = append(events.Items[:i], events.Items[i+1:]...)
			i--
		}
	}

	return events
}
