package commands

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/utils"
	"TimeTrack-shared/dtos"
	"TimeTrack-shared/models"

	"fmt"
	"time"

	"github.com/urfave/cli/v2"
)

func getAddTimeEntryCommand(ctx *app.AppContext) *cli.Command {
	return &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "Add a new time entry",
		Flags:   addTimeEntryFlags(),
		Action: func(c *cli.Context) error {
			// Ensure user is authenticated
			if _, err := ctx.API.GetCurrentUser(); err != nil {
				return cli.Exit("Unauthorized or not logged in. Please login or register first.", 1)
			}

			// Validate input formats
			if err := validateTimeEntryInputs(c); err != nil {
				return err
			}

			// Determine start/end times
			endDate := c.String("endDate")
			if endDate == "" {
				endDate = c.String("date")
			}

			startTimeStr := fmt.Sprintf("%sT%s:00+02:00", c.String("date"), c.String("start"))
			endTimeStr := fmt.Sprintf("%sT%s:00+02:00", endDate, c.String("end"))

			// Retrieve or create project
			project, err := getOrCreateProject(ctx, c.String("name"))
			if err != nil {
				return err
			}

			// Parse start/end times and validate order
			startTimeParsed, _ := time.Parse(time.RFC3339, startTimeStr)
			endTimeParsed, _ := time.Parse(time.RFC3339, endTimeStr)
			if endTimeParsed.Before(startTimeParsed) {
				return cli.Exit("End time is before start time.", 1)
			}

			// Create time entry object
			entry := &dtos.CreateTimeEntryInput{
				ProjectID: project.ID,
				Note:      c.String("description"),
				Period: dtos.TimePeriod{
					Start: startTimeParsed,
					End:   endTimeParsed,
				},
			}

			// Confirm and create
			if c.Bool("skipConfirmation") || utils.Confirm(
				"This will create a new time entry with the following details:\n\n"+
					getTimeEntryInformationString(project, entry)+
					"\n\nDo you want to proceed?",
			) {
				if _, err := ctx.API.CreateTimeEntry(entry); err != nil {
					return cli.Exit("Failed to create time entry: "+err.Error(), 1)
				}
				fmt.Println("Time Entry created with the following details:")
				fmt.Println(getTimeEntryInformationString(project, entry))
			} else {
				fmt.Println("Time Entry not created.")
			}

			return nil
		},
	}
}

func addTimeEntryFlags() []cli.Flag {
	return []cli.Flag{
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
			Usage:   "Description of time entry",
		},
		&cli.StringFlag{
			Name:     "start",
			Aliases:  []string{"s"},
			Required: true,
			Usage:    "Start time of time entry (format: HH:mm)",
		},
		&cli.StringFlag{
			Name:     "end",
			Aliases:  []string{"e"},
			Required: true,
			Usage:    "End time of time entry (format: HH:mm)",
		},
		&cli.StringFlag{
			Name:    "date",
			Value:   time.Now().Format("2006-01-02"),
			Aliases: []string{"d"},
			Usage:   "Date of time entry. (format: YYYY-MM-DD)",
		},
		&cli.StringFlag{
			Name:    "endDate",
			Aliases: []string{"E"},
			Usage:   "End date of time entry. If not provided, the date flag will be used. (format: YYYY-MM-DD)",
		},
		&cli.BoolFlag{
			Name:    "skipConfirmation",
			Aliases: []string{"yes", "y"},
			Usage:   "Skip confirmation",
		},
	}
}

func validateTimeEntryInputs(c *cli.Context) error {
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
	return nil
}

func getOrCreateProject(ctx *app.AppContext, name string) (*models.Project, error) {
	project, err := ctx.API.GetProjectByName(name)
	if project == nil || err != nil {
		if !utils.Confirm("Project not found. Do you want to create a new project with the name '" + name + "'?") {
			return nil, cli.Exit("Project creation aborted.", 1)
		}

		projectInput := &dtos.CreateProjectInput{Name: name}
		project, err = ctx.API.CreateProject(projectInput)
		if err != nil {
			return nil, cli.Exit("Failed to create project: "+err.Error(), 1)
		}
		fmt.Printf("Project created: %s\n", project.Name)
	}

	if project.ID == "" {
		return nil, cli.Exit("Failed to create or retrieve project.", 1)
	}
	return project, nil
}

func getTimeEntryInformationString(project *models.Project, entry *dtos.CreateTimeEntryInput) string {
	note := entry.Note
	if note == "" {
		note = "(no description provided)"
	}

	return fmt.Sprintf(
		"Project: %s\nDescription: %s\nStart: %s\nEnd: %s",
		project.Name,
		note,
		utils.FormatDate(entry.Period.Start.Format(time.RFC3339), time.RFC3339),
		utils.FormatDate(entry.Period.End.Format(time.RFC3339), time.RFC3339),
	)
}
