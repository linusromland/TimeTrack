package calendar

import (
	"TimeTrack/src/oauth"
	"context"
	"log"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetCalendarService() *calendar.Service {
	client := oauth.GetClient()
	ctx := context.Background()
	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve calendar Client %v", err)
	}
	return service
}
