package commands

import (
	"TimeTrack/src/calendar"
	"TimeTrack/src/database"
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

var SelectDatabaseCommand = &cli.Command{
	Name:    "selectDatabase",
	Aliases: []string{"sd"},
	Usage:   "Select database",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "calendarId",
			Aliases: []string{"c"},
			Usage:   "Calendar ID",
		},
	},
	Action: func(c *cli.Context) error {
		//Open database
		db, err := database.OpenDB()
		if err != nil {
			fmt.Println(err)
		}

		calendarId := c.String("calendarId")
		if calendarId == "" {
			calendars := calendar.GetCalendars()

			app := tview.NewApplication()

			list := tview.NewList().ShowSecondaryText(false).
				AddItem("Select a calendar", "", 'a', nil)

			for _, cal := range calendars.Items {
				list.AddItem(cal.Summary, cal.Id, 0, func() {
					app.Stop()
				})
			}

			list.AddItem("Exit", "", 'e', func() {
				app.Stop()
			})

			// log all items in list
			list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
				calendarId = secondaryText
			})

			if err := app.SetRoot(list, true).Run(); err != nil {
				panic(err)
			}

		} else {
			updateCalendarId(db, calendarId)
		}

		if calendarId != "" {
			updateCalendarId(db, calendarId)
		} else {
			fmt.Println("No calendar selected.")
		}

		// Close database
		err = database.CloseDB(db)
		if err != nil {
			fmt.Println(err)
		}

		return nil
	},
}

func updateCalendarId(db *badger.DB, calendarId string) error {
	calendar := calendar.GetCalendar(calendarId)

	// Check if calendar exists
	if calendar == nil {
		fmt.Println("Calendar not found.")
		return nil
	}

	// Save calendar ID to database
	err := database.InsertData(db, "calendarId", calendarId)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Calendar '" + calendar.Summary + "' selected.")
	}

	return nil
}
