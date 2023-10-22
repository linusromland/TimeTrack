package main

import (
	"fmt"
	"os"

	"TimeTrack/src/commands"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	//Load .env file
	godotenv.Load()

	app := &cli.App{
		Name:     "TimeTrack",
		Usage:    "Easy time tracking from the command line. With built-in integration for Google Calendar.",
		Commands: commands.AllCommands,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
