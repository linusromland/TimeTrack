package settings

import (
	"TimeTrack-cli/src/database"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RenderSettingsUI(app *tview.Application, db *database.DBWrapper, goBack func()) *tview.List {
	settingsList := GetAllSettings(db) // always get fresh values

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(" Settings ")
	list.SetTitleAlign(tview.AlignLeft)
	list.SetBorderPadding(1, 1, 2, 2)
	list.SetSecondaryTextColor(tcell.ColorGrey)

	for _, s := range settingsList {
		if s.Type() == "static" {
			list.AddItem(s.Label(), s.Get(), 0, nil)
		} else {
			list.AddItem(s.Label(), "Current value: "+s.Get(), 0, nil)
		}
	}

	list.AddItem("Back", "", 'b', goBack)

	list.SetSelectedFunc(func(index int, mainText, _ string, _ rune) {
		if mainText == "Back" {
			goBack()
			return
		}
		if index < len(settingsList) {
			if settingsList[index].Type() == "static" {
				// Re-render with fresh data
				app.SetRoot(RenderSettingsUI(app, db, goBack), true)
			} else {
				editSetting(app, db, settingsList, index, goBack)
			}
		}
	})

	return list
}

func editSetting(app *tview.Application, db *database.DBWrapper, settingsList []Setting, index int, goBack func()) {
	setting := settingsList[index]
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle(" Edit " + setting.Label() + " ")
	form.SetTitleAlign(tview.AlignLeft)
	form.SetBorderPadding(1, 1, 2, 2)

	currentValue := setting.Get()

	switch setting.Type() {
	case "text":
		form.AddInputField(setting.Label(), currentValue, 0, nil, func(text string) {
			currentValue = text
		})
	case "number":
		form.AddInputField(setting.Label(), currentValue, 0, func(textToCheck string, lastChar rune) bool {
			if textToCheck == "" {
				return true
			}
			_, err := strconv.Atoi(textToCheck)
			return err == nil
		}, func(text string) {
			currentValue = text
		})
	case "checkbox":
		boolVal := currentValue == "true"
		form.AddCheckbox(setting.Label(), boolVal, func(checked bool) {
			currentValue = strconv.FormatBool(checked)
		})
	default:
		panic("unknown setting type: " + setting.Type())
	}

	form.AddButton("Save", func() {
		setting.Set(currentValue)
		app.SetRoot(RenderSettingsUI(app, db, goBack), true)
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(RenderSettingsUI(app, db, goBack), true)
	})

	app.SetRoot(form, true)
}
