package commands

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/ui"
	"TimeTrack-cli/src/ui/screens"

	"github.com/urfave/cli/v2"
)

func getRegisterCommand(ctx *app.AppContext) *cli.Command {
	return &cli.Command{
		Name:  "register",
		Usage: "User registration",
		Action: func(c *cli.Context) error {
			nav := ui.NewNavigator()
			return nav.Run(screens.RegisterModal(nav, ctx, true))
		},
	}
}
