package commands

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/urfave/cli/v2"
)

type Setting struct {
	Id       string
	Type     string
	Label    string
	Options  []string                            // For selects
	GetValue func(value string) (string, string) // Returns value and label. value is used for selects, for all others, leave blank
	SetValue func(string)                        // Sets value
}

type SettingCategory struct {
	Id       string
	Label    string
	Settings []Setting
}

var settings = []SettingCategory{
	{
		Id:    "calendar",
		Label: "Calendar",
		Settings: []Setting{
			{
				Id:      "calendar",
				Type:    "select",
				Label:   "Selected Calendar",
				Options: []string{"123", "234"},
				GetValue: func(value string) (string, string) {
					if value == "123" {
						return "123Value0", "123Label0"
					} else {
						return "234Value0", "234Label0"
					}
				},
				SetValue: func(value string) {
					fmt.Printf("Set calendar to %s\n", value)
				},
			},
		},
	},
	{
		Id:    "cloudSync",
		Label: "Cloud Sync",
		Settings: []Setting{
			{
				Id:    "enabled",
				Type:  "checkbox",
				Label: "Cloud Sync Enabled",
				GetValue: func(_ string) (string, string) {
					return "true", "Enabled"
				},
				SetValue: func(value string) {
					fmt.Printf("Set enabled to %s\n", value)
				},
			},
			{
				Id:    "url",
				Type:  "text",
				Label: "Cloud Sync URL",
				GetValue: func(_ string) (string, string) {
					return "https://example.com", "https://example.com"
				},
				SetValue: func(value string) {
					fmt.Printf("Set url to %s\n", value)
				},
			},
			{
				Id:    "interval",
				Type:  "number",
				Label: "Sync Interval",
				GetValue: func(_ string) (string, string) {
					return "0", "Every time"
				},
				SetValue: func(value string) {
					fmt.Printf("Set interval to %s\n", value)
				},
			},
		},
	},
}

var categorySettings = map[string]string{
	"Calendar":   "calendar",
	"Cloud Sync": "cloudSync",
}

var SettingsCommand = &cli.Command{
	Name:  "settings",
	Usage: "Manage application settings",
	Action: func(c *cli.Context) error {
		app := tview.NewApplication()
		mainList := createMainList(app)

		if err := app.SetRoot(mainList, true).Run(); err != nil {
			panic(err)
		}

		return nil
	},
}

func getSettingCategory(id string) SettingCategory {
	for _, settingCategory := range settings {
		if settingCategory.Id == id {
			return settingCategory
		}
	}

	return SettingCategory{}
}

func GetSettingsByCategory(category string) []Setting {
	// get the settingsCategory from the map
	settingsCategoryId := categorySettings[category]

	settingCategory := getSettingCategory(settingsCategoryId)

	return settingCategory.Settings
}

func createMainList(app *tview.Application) *tview.List {
	mainList := tview.NewList()

	mainList.SetBorderPadding(1, 1, 2, 2).SetBorder(true).SetTitle(" Settings ").SetTitleAlign(tview.AlignLeft)

	mainList.AddItem("Calendar", "", 0, nil)
	mainList.AddItem("Cloud Sync", "", 0, nil)
	mainList.AddItem("Exit", "", 'e', func() {
		app.Stop()
	})

	mainList.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if mainText == "Exit" {
			app.Stop()
			return
		}
		showSettingsForCategory(app, mainList, mainText)
	})

	return mainList
}

func showSettingsForCategory(app *tview.Application, mainList *tview.List, category string) {
	settingsList := tview.NewList()

	settingsList.SetBorderPadding(1, 1, 2, 2).SetBorder(true).SetTitle(" " + category + " ").SetTitleAlign(tview.AlignLeft)

	settings := GetSettingsByCategory(category)

	for _, setting := range settings {
		_, settingLabel := setting.GetValue("")

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

func editSetting(app *tview.Application, mainList *tview.List, setting Setting) {
	settingType := setting.Type // text, number, checkbox, select

	if settingType == "select" {
		editSelectSetting(app, mainList, setting)
		return
	}

	form := tview.NewForm()

	form.SetBorderPadding(1, 1, 2, 2).SetBorder(true).SetTitle(" " + setting.Label + " ").SetTitleAlign(tview.AlignLeft)

	switch settingType {
	case "text":
		form.AddInputField(setting.Label, "", 0, nil, nil)
	case "number":
		form.AddInputField(setting.Label, "", 0, nil, nil)

	case "checkbox":
		form.AddCheckbox(setting.Label, false, nil)
	default:
		panic("Unknown setting type: " + settingType)
	}

	form.AddButton("Save", func() {
		app.SetRoot(mainList, true)
	})

	form.AddButton("Cancel", func() {
		app.SetRoot(mainList, true)
	})

	app.SetRoot(form, true)
}

func editSelectSetting(app *tview.Application, mainList *tview.List, setting Setting) {
	selectList := tview.NewList()

	if setting.Type != "select" {
		panic("editSelectSetting called with non-select setting")
	}

	if len(setting.Options) == 0 {
		panic("editSelectSetting called with setting with no options")
	}

	selectList.SetBorderPadding(1, 1, 2, 2).SetBorder(true).SetTitle(" " + setting.Label + " ").SetTitleAlign(tview.AlignLeft)

	for _, option := range setting.Options {
		value, label := setting.GetValue(option)

		selectList.AddItem(label, "", 0, func() {
			setting.SetValue(value)
			app.SetRoot(mainList, true)
		})
	}

	selectList.AddItem("Back", "", 'b', func() {
		app.SetRoot(mainList, true)
	})

	app.SetRoot(selectList, true)
}
