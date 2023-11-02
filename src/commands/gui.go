package commands

import (
	"TimeTrack/src/gui"

	"github.com/urfave/cli/v2"
)

var GuiCommand = &cli.Command{
	Name:  "gui",
	Usage: "start the Graphical User Interface(GUI) for TimeTrack",
	Action: func(c *cli.Context) error {
		gui.Launch()
		return nil
	},
}
