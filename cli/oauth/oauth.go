package oauth

import (
	coreCalendar "TimeTrack/core/calendar"
	"TimeTrack/core/oauth"
	"fmt"

	"google.golang.org/api/calendar/v3"
)

var (
	PRODUCTION_CLIENT_ID     string
	PRODUCTION_CLIENT_SECRET string
)

func GetCalendarService() *calendar.Service {
	fmt.Println("Production client ID:", PRODUCTION_CLIENT_ID)
	fmt.Println("Production client secret:", PRODUCTION_CLIENT_SECRET)
	client := oauth.GetClient(PRODUCTION_CLIENT_ID, PRODUCTION_CLIENT_SECRET)
	return coreCalendar.GetCalendarService(client)
}
