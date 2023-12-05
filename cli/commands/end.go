package commands

import (
	"TimeTrack/core/calendar"
	"TimeTrack/core/database"
	"TimeTrack/core/utils"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

var EndCommand = &cli.Command{
	Name:    "end",
	Aliases: []string{"e"},
	Usage:   "End a task",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"desc", "D"},
			Value:   "",
			Usage:   "Description of task",
		},
		&cli.StringFlag{
			Name:    "end",
			Aliases: []string{"e"},
			Value:   time.Now().Format("15:04"),
			Usage:   "End time of task (format: HH:mm)",
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

		if !utils.IsValidTime(c.String("end")) {
			return cli.Exit("Invalid end time. Please use the following format: HH:mm", 1)
		}

		currentTask := database.GetData(db, "currentTask")
		if currentTask == "" {
			return cli.Exit("No task is currently running. Please start a task before ending it.", 1)
		}

		startTime := database.GetData(db, "currentTaskStartTime")
		if startTime == "" {
			return cli.Exit("No start time found for current task. Please start a task before ending it.", 1)
		}

		database.SetData(db, "currentTask", "")
		database.SetData(db, "currentTaskStartTime", "")

		parsedStartTime, _ := time.Parse(time.RFC3339, startTime)
		// if starttime is less than 1 minute ago, exit with error
		if parsedStartTime.Add(time.Minute * 1).After(time.Now()) {
			return cli.Exit("Task started less than 1 minute ago. Please wait at least 1 minute before ending the task.", 1)
		}

		endTime := fmt.Sprintf("%sT%s:00+02:00", time.Now().Format("2006-01-02"), c.String("end"))
		parsedEndTime, _ := time.Parse(time.RFC3339, endTime)

		event := calendar.CreateEvent(calendarId, currentTask, c.String("description"), startTime, endTime)

		hours, minutes := utils.GetTaskTime(parsedStartTime, parsedEndTime)

		fmt.Printf("Ended task '%s' at %s\n", currentTask, parsedEndTime.Format("2006-01-02 15:04"))
		if c.String("description") != "" {
			fmt.Printf("Description: %s\n", c.String("description"))
		}
		fmt.Printf("Duration: %dh %dm\n", hours, minutes)
		fmt.Printf("Event created: %s\n", event.HtmlLink)
		return nil
	},
}
