package main

import (
	"fmt"
	"os"

	"TimeTrack/src/commands"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
)

func main() {
	if version == "dev" {
		godotenv.Load()
	}

	if os.Getenv("enviroment") == "development" {
		fmt.Println("Running in development mode")
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
