package commands

import (
	"TimeTrack/core/database"
	"TimeTrack/core/settings"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

var SettingsCommand = &cli.Command{
	Name:  "settings",
	Usage: "Manage application settings",
	Action: func(c *cli.Context) error {
		app := tview.NewApplication()

		db, err := database.OpenDB()
		if err != nil {
			panic(err)
		}

		mainList := createMainList(app, db)

		if err := app.SetRoot(mainList, true).Run(); err != nil {
			panic(err)
		}

		database.CloseDB(db)
		return nil
	},
}

func createMainList(app *tview.Application, db *badger.DB) *tview.List {
	mainList := tview.NewList()

	mainList.SetBorderPadding(1, 1, 2, 2).SetBorder(true).SetTitle(" Settings ").SetTitleAlign(tview.AlignLeft)

	mainList.AddItem("Calendar", "", 0, nil)
	mainList.AddItem("Cloud Sync", "", 0, nil)
	mainList.AddItem("Exit", "", 'e', func() {
		app.Stop()
	})
	mainList.ShowSecondaryText(false)

	mainList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if mainText == "Exit" {
			app.Stop()
			return
		}

		showSettingsForCategory(app, mainList, mainText, db)
	})

	return mainList
}

func showSettingsForCategory(app *tview.Application, mainList *tview.List, category string, db *badger.DB) {
	settingsList := tview.NewList()

	settingsList.SetBorderPadding(1, 1, 2, 2).SetBorder(true).SetTitle(" " + category + " ").SetTitleAlign(tview.AlignLeft)

	settings := settings.GetSettingsByCategory(db, category)

	for _, setting := range settings {
		_, settingLabel := setting.GetValue("")

		if settingLabel == "" {
			settingLabel = "Not set"
		}

		valueLabel := "Current value: " + settingLabel

		settingsList.AddItem(setting.Label, valueLabel, 0, nil)
		settingsList.SetSecondaryTextColor(tcell.ColorGrey)
	}

	settingsList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if mainText == "Back" {
			app.SetRoot(mainList, true)
			return
		}

		setting := settings[index]

		editSetting(app, mainList, setting)
	})

	settingsList.AddItem("Back", "", 'b', func() {
		app.SetRoot(mainList, true)
	})

	app.SetRoot(settingsList, true)
}

func editSetting(app *tview.Application, mainList *tview.List, setting settings.Setting) {
	settingType := setting.Type // text, number, checkbox, select

	if settingType == "select" {
		editSelectSetting(app, mainList, setting)
		return
	}

	form := tview.NewForm()

	form.SetBorderPadding(1, 1, 2, 2).SetBorder(true).SetTitle(" " + setting.Label + " ").SetTitleAlign(tview.AlignLeft)

	settingValue, _ := setting.GetValue("")

	textValue := ""
	numberValue := 0
	booleanValue := false

	switch settingType {
	case "text":
		textValue = settingValue
		form.AddInputField(setting.Label, settingValue, 0, nil, func(text string) {
			textValue = text
		})
	case "number":
		if settingValue != "" {
			numberValue = parseInt(settingValue)
		}
		form.AddInputField(setting.Label, settingValue, 0, func(textToCheck string, lastChar rune) bool {
			if textToCheck == "" {
				return true
			}

			_, err := strconv.Atoi(textToCheck)

			return err == nil
		}, func(text string) {
			if text == "" {
				numberValue = 0
				return
			}

			numberValue = parseInt(text)
		})
	case "checkbox":
		booleanValue = parseBool(settingValue)
		form.AddCheckbox(setting.Label, booleanValue, func(checked bool) {
			booleanValue = checked
		})
	default:
		panic("Unknown setting type: " + settingType)
	}

	form.AddButton("Save", func() {
		switch settingType {
		case "text":
			setting.SetValue(textValue)
		case "number":
			setting.SetValue(strconv.Itoa(numberValue))
		case "checkbox":
			setting.SetValue(strconv.FormatBool(booleanValue))
		default:
			panic("Unknown setting type: " + settingType)
		}
		app.SetRoot(mainList, true)
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(mainList, true)
	})

	app.SetRoot(form, true)
}

func parseInt(value string) int {
	intValue, err := strconv.Atoi(value)

	if err != nil {
		panic(err)
	} else {
		return intValue
	}
}

func parseBool(value string) bool {
	return value == "true"
}

func editSelectSetting(app *tview.Application, mainList *tview.List, setting settings.Setting) {
	selectList := tview.NewList()

	if setting.Type != "select" {
		panic("editSelectSetting called with non-select setting")
	}

	if len(setting.Options) == 0 {
		panic("editSelectSetting called with setting with no options")
	}

	selectList.SetBorderPadding(1, 1, 2, 2).SetBorder(true).SetTitle(" " + setting.Label + " ").SetTitleAlign(tview.AlignLeft)

	for _, option := range setting.Options {
		selectList.AddItem(option.Label, option.Value, 0, func() {
			setting.SetValue(option.Value)
			app.SetRoot(mainList, true)
		})
	}

	selectList.ShowSecondaryText(false)

	selectList.AddItem("Back", "", 'b', func() {
		app.SetRoot(mainList, true)
	})

	selectList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if mainText == "Back" {
			app.SetRoot(mainList, true)
			return
		}

		option := setting.Options[index]

		setting.SetValue(option.Value)
		app.SetRoot(mainList, true)
	})

	app.SetRoot(selectList, true)
}
