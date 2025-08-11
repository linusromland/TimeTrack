package screens

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/ui"
	"TimeTrack-cli/src/ui/components"

	"github.com/rivo/tview"
)

func LoginModal(nav *ui.Navigator, ctx *app.AppContext, exitOnCancel bool) tview.Primitive {
	form := components.StyledForm("Login")
	serverURL := getServerURL(ctx)
	var email, password string

	form.AddTextView("Login to", serverURL, 40, 1, false, false)
	form.AddInputField("Email", "", 40, nil, func(text string) { email = text })
	form.AddPasswordField("Password", "", 40, '*', func(text string) { password = text })
	form.AddButton("Login", func() {
		if email == "" || password == "" {
			nav.Show(components.StyledModal("Email and password required", func() { nav.Show(LoginModal(nav, ctx, exitOnCancel)) }))
			return
		}
		err := ctx.API.Login(email, password)
		if err != nil {
			nav.Show(components.StyledModal("Login failed: "+err.Error(), func() { nav.Show(LoginModal(nav, ctx, exitOnCancel)) }))
			return
		}
		nav.Show(components.StyledModal("Login successful!", func() { nav.Show(DashboardScreen(nav, ctx)) }))
	})
	form.AddButton("Cancel", func() {
		if exitOnCancel {
			nav.Stop()
			return
		}

		nav.Show(DashboardScreen(nav, ctx))
	})

	return form
}
