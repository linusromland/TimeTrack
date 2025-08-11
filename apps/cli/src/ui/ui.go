package ui

import (
	"TimeTrack-cli/src/database"
	settingsUi "TimeTrack-cli/src/ui/settings"

	"github.com/rivo/tview"
)

func RenderSettingsUI(app *tview.Application, db *database.DBWrapper, goBack func()) *tview.List {
	return settingsUi.RenderSettingsUI(app, db, goBack)
}
