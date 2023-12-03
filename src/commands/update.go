package commands

import (
	"TimeTrack/src/utils"
	"fmt"

	"github.com/urfave/cli/v2"
)

var UpdateCommand = &cli.Command{
	Name:    "update",
	Aliases: []string{"u"},
	Usage:   "update the application",
	Action: func(c *cli.Context) error {
		version := c.App.Version
		newVersion, err := utils.CheckForUpdate(version, true)

		if err != nil {
			return err
		} else if !newVersion {
			fmt.Println("You are already running the latest version. The current version is:", version)
		}

		return nil
	},
}
