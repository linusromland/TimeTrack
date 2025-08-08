package main

import (
	"TimeTrack-cli/src/commands"

	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var (
	version = "dev"
)

func main() {

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
