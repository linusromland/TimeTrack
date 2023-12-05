package commands

import (
	"TimeTrack/core/database"
	"TimeTrack/core/utils"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

var StartCommand = &cli.Command{
	Name:    "start",
	Aliases: []string{"s"},
	Usage:   "Start a task",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Required: true,
			Usage:    "Name of task",
		},
		&cli.StringFlag{
			Name:    "start",
			Aliases: []string{"s"},
			Value:   time.Now().Format("15:04"),
			Usage:   "Start time of task (format: HH:mm)",
		},
	},
	Action: func(c *cli.Context) error {
		db, err := database.OpenDB()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		calendarId := database.GetData(db, "calendarId")
		if calendarId == "" {
			fmt.Println("No calendar selected. Please select a calendar with the selectCalendar command.")
			return nil
		}

		if !utils.IsValidTime(c.String("start")) {
			return cli.Exit("Invalid start time. Please use the following format: HH:mm", 1)
		}

		currentTask := database.GetData(db, "currentTask")
		if currentTask != "" {
			return cli.Exit("Task '"+currentTask+"' is already running. Please stop the current task before starting a new one.", 1)
		}

		startTime := fmt.Sprintf("%sT%s:00+02:00", time.Now().Format("2006-01-02"), c.String("start"))
		parsedStartTime, _ := time.Parse(time.RFC3339, startTime)

		database.SetData(db, "currentTask", c.String("name"))
		database.SetData(db, "currentTaskStartTime", startTime)

		fmt.Printf("Started task '%s' at %s\n", c.String("name"), parsedStartTime.Format("2006-01-02 15:04"))
		return nil
	},
}
