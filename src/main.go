package main

import (
	"log"
	"os"

	"TimeTrack/src/commands"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "TimeTrack",
		Usage:    "Easy time tracking from the command line. With built-in integration for Google Calendar.",
		Commands: commands.AllCommands,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
