package commands

import (
	"TimeTrack/calendar"
	"TimeTrack/database"
	"TimeTrack/utils"

	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

var ChangeCommand = &cli.Command{
	Name:    "change",
	Aliases: []string{"c"},
	Usage:   "End the current task and start a new one",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "end",
			Aliases: []string{"e"},
			Value:   time.Now().Format("15:04"),
			Usage:   "End time of current task and start time of new task (format: HH:mm)",
		},
		&cli.StringFlag{
			Name:     "name",
			Aliases:  []string{"n"},
			Required: true,
			Usage:    "Name of new task",
		},
		&cli.StringFlag{
			Name:    "description",
			Aliases: []string{"desc", "D"},
			Value:   "",
			Usage:   "Description of the task that's ending",
		},
	},
	Action: func(c *cli.Context) error {
		db, err := database.OpenDB()
		if err != nil {
			fmt.Println(err)
			return nil
		}
		calendarId := database.GetData(db, database.CALENDAR_ID)
		if calendarId == "" {
			fmt.Println("No calendar selected. Please select a calendar with the selectCalendar command.")
			return nil
		}

		if !utils.IsValidTime(c.String("end")) {
			return cli.Exit("Invalid end/start time. Please use the following format: HH:mm", 1)
		}

		// End current task
		currentTask := database.GetData(db, database.CURRENT_TASK)
		if currentTask == "" {
			return cli.Exit("No task is currently running. Please start a task before changing it.", 1)
		}

		startTime := database.GetData(db, database.CURRENT_TASK_START_TIME)
		if startTime == "" {
			return cli.Exit("No start time found for current task. Please start a task before changing it.", 1)
		}

		endTime := fmt.Sprintf("%sT%s:00+02:00", time.Now().Format("2006-01-02"), c.String("end"))
		parsedEndTime, _ := time.Parse(time.RFC3339, endTime)
		parsedStartTime, _ := time.Parse(time.RFC3339, startTime)

		event := calendar.CreateEvent(calendarId, currentTask, c.String("description"), startTime, endTime)

		hours, minutes := utils.GetTaskTime(parsedStartTime, parsedEndTime)

		fmt.Printf("Ended task '%s' at %s\n", currentTask, parsedEndTime.Format("2006-01-02 15:04"))
		if c.String("description") != "" {
			fmt.Printf("Description: %s\n", c.String("description"))
		}
		fmt.Printf("Duration: %dh %dm\n", hours, minutes)
		fmt.Printf("Event created: %s\n", event.HtmlLink)

		// Start new task
		database.SetData(db, database.CURRENT_TASK, c.String("name"))
		database.SetData(db, database.CURRENT_TASK_START_TIME, endTime)

		fmt.Printf("Started new task '%s' at %s\n", c.String("name"), parsedEndTime.Format("2006-01-02 15:04"))
		return nil
	},
}
