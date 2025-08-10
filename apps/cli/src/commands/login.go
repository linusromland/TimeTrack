package commands

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/auth"
	"TimeTrack-cli/src/database"

	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

func getLoginCommand(ctx *app.AppContext) *cli.Command {
	return &cli.Command{
		Name:  "login",
		Usage: "User login",
		Action: func(c *cli.Context) error {
			appUI := tview.NewApplication()
			serverURL := ctx.DB.Get(database.ServerURLKey)

			var loginForm *tview.Form
			loginForm = auth.CreateLoginForm(appUI, serverURL, func(email, password string) {
				err := ctx.API.Login(email, password)
				if err != nil {
					auth.ShowModal(appUI, "Login failed: "+err.Error(), func() {
						appUI.SetRoot(loginForm, true)
					})
					return
				}
				auth.ShowModal(appUI, "Login successful!", func() {
					appUI.Stop()
				})
			})

			return appUI.SetRoot(loginForm, true).Run()
		},
	}
}
