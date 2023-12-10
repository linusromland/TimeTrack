package main

import (
	"TimeTrack/cli/commands"
	cliUtils "TimeTrack/cli/utils"
	"TimeTrack/core/database"
	"TimeTrack/core/utils"

	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
)

func main() {
	if version == "dev" {
		fmt.Println("Running in development mode.")
		godotenv.Load("../.env")
	}

	checkForUpdate()

	app := &cli.App{
		Name:     "TimeTrack",
		Usage:    "Easy time tracking from the command line. With built-in integration for Google Calendar.",
		Version:  version,
		Commands: commands.AllCommands,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func checkForUpdate() {
	updateAvailable, err := utils.CheckForUpdate(version)
	if err != nil {
		fmt.Println(err)
	}
	if updateAvailable != "" {
		db, err := database.OpenDB()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("There is a new version of TimeTrack available: %s\n", updateAvailable)

		// Check if the user has skipped this update.
		skipUpdate := database.GetData(db, database.SKIP_UPDATE)
		if skipUpdate == updateAvailable {
			if cliUtils.Confirm("Do you want to update?") {
				err = utils.UpdateVersion(updateAvailable)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Skipping this update, you can update later by running 'timetrack update'.")
				database.SetData(db, database.SKIP_UPDATE, updateAvailable)
			}
		}
		database.CloseDB(db)
	}
}
