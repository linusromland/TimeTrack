package calendar

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetCalendarService(client *http.Client) *calendar.Service {
	ctx := context.Background()
	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Unable to retrieve calendar Client %v", err)
	}
	return service
}
