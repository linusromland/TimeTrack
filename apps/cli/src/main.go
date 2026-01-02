package main

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/commands"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var version = "dev"

func main() {
	ctx := app.NewAppContext(version)

	appCLI := &cli.App{
		Name:     "TimeTrack",
		Usage:    "Easy time tracking from the command line.",
		Version:  version,
		Before:   ctx.Startup,
		Commands: commands.GetAllCommands(ctx),
	}

	if err := appCLI.Run(os.Args); err != nil {
		fmt.Println("Error:", err)
	}
}
