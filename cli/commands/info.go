package commands

import (
	"TimeTrack/core/database"
	"TimeTrack/core/utils"

	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

var InfoCommand = &cli.Command{
	Name:    "info",
	Aliases: []string{"i"},
	Usage:   "Show info about current task",
	Action: func(c *cli.Context) error {
		db, err := database.OpenDB()
		if err != nil {
			fmt.Println(err)
			return nil
		}

		currentTask := database.GetData(db, "currentTask")
		if currentTask == "" {
			return cli.Exit("No task is currently running. Please start a task to see info about it.", 1)
		}

		startTime := database.GetData(db, "currentTaskStartTime")
		if startTime == "" {
			return cli.Exit("No start time found for current task. Please start a task to see info about it.", 1)
		}

		parsedStartTime, _ := time.Parse(time.RFC3339, startTime)

		currentTime := time.Now()

		hours, minutes := utils.GetTaskTime(parsedStartTime, currentTime)

		fmt.Printf("Current task '%s' started at %s\n", currentTask, parsedStartTime.Format("15:04"))
		fmt.Printf("Current task has been running for %d hours and %d minutes\n", hours, minutes)
		return nil
	},
}
