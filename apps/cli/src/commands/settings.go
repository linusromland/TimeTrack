package commands

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/ui"

	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

func getSettingsCommand(ctx *app.AppContext) *cli.Command {
	return &cli.Command{
		Name:  "settings",
		Usage: "Manage application settings",
		Action: func(c *cli.Context) error {
			appUI := tview.NewApplication()

			mainUI := ui.RenderSettingsUI(appUI, ctx.DB, func() {
				appUI.Stop()
			})

			if err := appUI.SetRoot(mainUI, true).Run(); err != nil {
				return err
			}
			return nil
		},
	}
}
