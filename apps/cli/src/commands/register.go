package commands

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/auth"
	"TimeTrack-cli/src/database"

	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

func getRegisterCommand(ctx *app.AppContext) *cli.Command {
	return &cli.Command{
		Name:  "register",
		Usage: "User registration",
		Action: func(c *cli.Context) error {
			appUI := tview.NewApplication()
			serverURL := ctx.DB.Get(database.ServerURLKey)

			var registerForm *tview.Form
			registerForm = auth.CreateRegisterForm(appUI, serverURL, func(email, password string) {
				err := ctx.API.Register(email, password)
				if err != nil {
					auth.ShowModal(appUI, "Registration failed: "+err.Error(), func() {
						appUI.SetRoot(registerForm, true)
					})
					return
				}

				err = ctx.API.Login(email, password)
				if err != nil {
					auth.ShowModal(appUI, "Login after registration failed: "+err.Error(), func() {
						appUI.SetRoot(registerForm, true)
					})
					return
				}

				auth.ShowModal(appUI, "Registration successful!", func() {
					appUI.Stop()
				})
			})

			return appUI.SetRoot(registerForm, true).Run()
		},
	}
}
