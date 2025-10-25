package screens

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/ui"
	"TimeTrack-cli/src/ui/components"

	"github.com/rivo/tview"
)

func RegisterModal(nav *ui.Navigator, ctx *app.AppContext, exitOnCancel bool) tview.Primitive {
	form := components.StyledForm("Register")
	serverURL := getServerURL(ctx)
	var email, confirmEmail, password, confirmPassword string

	form.AddTextView("Register at", serverURL, 40, 1, false, false)
	form.AddInputField("Email", "", 40, nil, func(text string) { email = text })
	form.AddInputField("Confirm Email", "", 40, nil, func(text string) { confirmEmail = text })
	form.AddPasswordField("Password", "", 40, '*', func(text string) { password = text })
	form.AddPasswordField("Confirm Password", "", 40, '*', func(text string) { confirmPassword = text })

	form.AddButton("Register", func() {
		if email != confirmEmail {
			nav.Show(components.StyledModal("Emails do not match", func() { nav.Show(RegisterModal(nav, ctx, exitOnCancel)) }))
			return
		}
		if password != confirmPassword {
			nav.Show(components.StyledModal("Passwords do not match", func() { nav.Show(RegisterModal(nav, ctx, exitOnCancel)) }))
			return
		}
		err := ctx.API.Register(email, password)
		if err != nil {
			nav.Show(components.StyledModal("Registration failed: "+err.Error(), func() { nav.Show(RegisterModal(nav, ctx, exitOnCancel)) }))
			return
		}
		nav.Show(components.StyledModal("Registration successful!", func() { nav.Show(DashboardScreen(nav, ctx)) }))
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
