package commands

import (
	"TimeTrack/src/database"
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

var AbortCommand = &cli.Command{
	Name:    "abort",
	Aliases: []string{"a"},
	Usage:   "Abort a task",
	Action: func(c *cli.Context) error {
		db, err := database.OpenDB()
		if err != nil {
			fmt.Println(err)
			return nil
		}

		currentTask := database.GetData(db, "currentTask")
		if currentTask == "" {
			return cli.Exit("No task is currently running. Please start a task before aborting it.", 1)
		}

		startTime := database.GetData(db, "currentTaskStartTime")
		if startTime == "" {
			return cli.Exit("No start time found for current task. Please start a task before aborting it.", 1)
		}
		parsedStartTime, _ := time.Parse(time.RFC3339, startTime)

		database.SetData(db, "currentTask", "")
		database.SetData(db, "currentTaskStartTime", "")

		fmt.Printf("Aborted task '%s' at %s\n", currentTask, parsedStartTime.Format("2006-01-02 15:04"))
		return nil
	},
}
