package commands

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/ui"
	"TimeTrack-cli/src/ui/screens"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

func getListTimeEntriesCommand(ctx *app.AppContext) *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List time entries",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "start",
				Value:   time.Now().Format("2006-01-02"),
				Aliases: []string{"s"},
				Usage:   "Start date of time entries to show. (format: YYYY-MM-DD)",
			},
			&cli.StringFlag{
				Name:    "end",
				Aliases: []string{"e"},
				Usage:   "End date of time entries to show. If not set, the start date will be used. (format: YYYY-MM-DD)",
			},
			&cli.StringFlag{
				Name:    "next",
				Aliases: []string{"n"},
				Usage:   "Show time entries from the next x units. This will start from the start date. (format: <number><unit>, unit options: d (days), w (week), m (month))",
			},
			&cli.StringFlag{
				Name:    "last",
				Aliases: []string{"l"},
				Usage:   "Show time entries from the last x units. This will start from the start date and go backwards. (format: <number><unit>, unit options: d (days), w (week), m (month))",
			},
		},
		Action: func(c *cli.Context) error {
			startDate, err := time.Parse("2006-01-02", c.String("start"))
			if err != nil {
				return fmt.Errorf("invalid start date: %v", err)
			}

			endDate := startDate
			if c.String("end") != "" {
				endDate, err = time.Parse("2006-01-02", c.String("end"))
				if err != nil {
					return fmt.Errorf("invalid end date: %v", err)
				}
			}

			if next := c.String("next"); next != "" {
				dur, err := parseRelative(next)
				if err != nil {
					return err
				}
				endDate = startDate.Add(dur)
			}

			if last := c.String("last"); last != "" {
				dur, err := parseRelative(last)
				if err != nil {
					return err
				}
				startDate = startDate.Add(-dur)
			}

			startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
			endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), endDate.Location())

			nav := ui.NewNavigator()
			return nav.Run(screens.TimeEntriesScreen(nav, ctx, startDate, endDate))
		},
	}
}

func parseRelative(s string) (time.Duration, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid relative format: %s", s)
	}
	numPart := s[:len(s)-1]
	unit := s[len(s)-1]
	val, err := time.ParseDuration(fmt.Sprintf("%sh", numPart))
	if err != nil {
		var num int
		_, err = fmt.Sscanf(numPart, "%d", &num)
		if err != nil {
			return 0, fmt.Errorf("invalid number in: %s", s)
		}
		switch unit {
		case 'd':
			return time.Duration(num) * 24 * time.Hour, nil
		case 'w':
			return time.Duration(num*7) * 24 * time.Hour, nil
		case 'm':
			return time.Duration(num*30) * 24 * time.Hour, nil
		default:
			return 0, fmt.Errorf("unknown unit: %c", unit)
		}
	}
	return val, nil
}
