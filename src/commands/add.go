package commands

import (
	"TimeTrack/src/calendar"
	"TimeTrack/src/utils"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

var AddCommand = &cli.Command{
	Name:    "add",
	Aliases: []string{"a"},
	Usage:   "add a task to the list",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Required: true,
			Usage:    "Name of task",
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"desc", "D"},
			Value:   "",
			Usage:   "Description of task",
		},
		&cli.StringFlag{
			Name:     "start",
			Aliases:  []string{"s"},
			Required: true,
			Usage:    "Start time of task (format: HH:mm)",
		},
		&cli.StringFlag{
			Name:     "end",
			Aliases:  []string{"e"},
			Required: true,
			Usage:    "End time of task (format: HH:mm)",
		},
		&cli.StringFlag{
			Name:    "date",
			Value:   time.Now().Format("2006-01-02"),
			Aliases: []string{"d"},
			Usage:   "Date of task. (format: YYYY-MM-DD)",
		},
		&cli.StringFlag{
			Name:    "endDate",
			Aliases: []string{"E"},
			Usage:   "End date of task. If not provided, the date flag will be used. (format: YYYY-MM-DD)",
		},
		&cli.BoolFlag{
			Name:    "skipConfirmation",
			Aliases: []string{"yes", "y"},
			Usage:   "Skip confirmation",
		},
	},
	Action: func(c *cli.Context) error {
		if !utils.IsValidTime(c.String("start")) {
			return cli.Exit("Invalid start time. Please use the following format: HH:mm", 1)
		}
		if !utils.IsValidTime(c.String("end")) {
			return cli.Exit("Invalid end time. Please use the following format: HH:mm", 1)
		}
		if !utils.IsValidDate(c.String("date")) {
			return cli.Exit("Invalid date. Please use the following format: YYYY-MM-DD", 1)
		}
		if c.String("endDate") != "" && !utils.IsValidDate(c.String("endDate")) {
			return cli.Exit("Invalid end date. Please use the following format: YYYY-MM-DD", 1)
		}

		endDate := c.String("endDate")
		if endDate == "" {
			endDate = c.String("date")
		}

		startTime := fmt.Sprintf("%sT%s:00+02:00", c.String("date"), c.String("start"))
		endTime := fmt.Sprintf("%sT%s:00+02:00", endDate, c.String("end"))

		if c.Bool("skipConfirmation") || utils.Confirm("This will create a new task with the following information:\nName: "+c.String("name")+"\nStart: "+utils.FormatDate(startTime, time.RFC3339)+"\nEnd: "+utils.FormatDate(endTime, time.RFC3339)+"\n\nAre you sure?") {

			startTimeParsed, _ := time.Parse(time.RFC3339, startTime)
			endTimeParsed, _ := time.Parse(time.RFC3339, endTime)
			if endTimeParsed.Before(startTimeParsed) {
				return cli.Exit("End time is before start time.", 1)
			}

			event := calendar.CreateEvent(c.String("name"), c.String("description"), startTime, endTime)
			fmt.Printf("Event created with the following information:\nName: %s\nStart: %s\nEnd: %s\nLink: %s\n", event.Summary, utils.FormatDate(event.Start.DateTime, time.RFC3339), utils.FormatDate(event.End.DateTime, time.RFC3339), event.HtmlLink)
		} else {
			fmt.Println("Task not created.")
		}
		return nil
	},
}
