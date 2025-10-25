package commands

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/ui"
	"TimeTrack-cli/src/ui/screens"

	"github.com/urfave/cli/v2"
)

func getSettingsCommand(ctx *app.AppContext) *cli.Command {
	return &cli.Command{
		Name:  "settings",
		Usage: "Show dashboard and manage settings",
		Action: func(c *cli.Context) error {
			nav := ui.NewNavigator()
			return nav.Run(screens.DashboardScreen(nav, ctx))
		},
	}
}
