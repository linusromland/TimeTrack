package oauth

import (
	coreCalendar "TimeTrack/core/calendar"
	"TimeTrack/core/oauth"

	"google.golang.org/api/calendar/v3"
)

var (
	PRODUCTION_CLIENT_ID     string
	PRODUCTION_CLIENT_SECRET string
)

func GetCalendarService() *calendar.Service {
	client := oauth.GetClient(PRODUCTION_CLIENT_ID, PRODUCTION_CLIENT_SECRET)
	return coreCalendar.GetCalendarService(client)
}
