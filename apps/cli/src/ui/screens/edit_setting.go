package screens

import (
	"TimeTrack-cli/src/app"
	"TimeTrack-cli/src/database"
	"TimeTrack-cli/src/ui"
	"TimeTrack-cli/src/ui/components"

	"github.com/rivo/tview"
)

// TODO: make this dynamic for other field other than server URL
func EditSettingModal(nav *ui.Navigator, ctx *app.AppContext) tview.Primitive {
	current := getServerURL(ctx)

	form := components.StyledForm("Edit Server URL")
	var newValue string = current

	form.AddInputField("Server URL", current, 40, nil, func(text string) { newValue = text })
	form.AddButton("Save", func() {
		ctx.DB.Set(database.ServerURLKey, newValue)
		nav.Show(DashboardScreen(nav, ctx))
	})
	form.AddButton("Cancel", func() {
		nav.Show(DashboardScreen(nav, ctx))
	})

	return form
}
