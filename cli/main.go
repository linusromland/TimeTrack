package main

import (
	"fmt"
	"os"

	"TimeTrack/cli/commands"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
)

func main() {
	if version == "dev" {
		fmt.Println("Running in development mode.")
		godotenv.Load()
	}

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
