package commands

import (
	cliUtils "TimeTrack/cli/utils"

	"TimeTrack/core/utils"

	"fmt"

	"github.com/urfave/cli/v2"
)

var UpdateCommand = &cli.Command{
	Name:    "update",
	Aliases: []string{"u"},
	Usage:   "update the application",
	Action: func(c *cli.Context) error {
		version := c.App.Version
		updateAvailable, err := utils.CheckForUpdate(version)

		if err != nil {
			return err
		}

		if updateAvailable == "" {
			fmt.Println("You are already running the latest version. The current version is:", version)
			return nil
		}

		if cliUtils.Confirm("Do you want to update?") {
			err = utils.UpdateVersion(updateAvailable)
			if err != nil {
				fmt.Println(err)
			}
		}

		return nil
	},
}
