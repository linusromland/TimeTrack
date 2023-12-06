package calendar

import (
	"TimeTrack/core/oauth"
	"context"
	"fmt"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetCalendarService() *calendar.Service {
	client := oauth.GetClient()
	ctx := context.Background()
	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve calendar Client %v", err)
	}
	return service
}
