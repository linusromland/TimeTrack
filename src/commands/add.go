package commands

import (
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

		fmt.Println("Creating task with the following parameters:")
		fmt.Println("Name:", c.String("name"))
		if c.String("description") != "" {
			fmt.Println("Description:", c.String("description"))
		}
		fmt.Println("Start:", c.String("start"))
		fmt.Println("End:", c.String("end"))
		if c.String("date") != "" {
			fmt.Println("Date:", c.String("date"))
		}
		if c.Bool("skipConfirmation") || utils.Confirm("Are you sure you want to create this task?") {
			fmt.Println("Task created.")

		} else {
			fmt.Println("Task not created.")
		}
		return nil
	},
}
