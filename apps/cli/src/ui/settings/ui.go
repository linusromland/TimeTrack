package settingsui

import (
	"TimeTrack-cli/src/database"
	"TimeTrack-cli/src/settings"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

func RenderSettingsUI(app *tview.Application, db *database.DBWrapper, goBack func()) *tview.List {
	settingsList := settings.GetAllSettings(db)

	grouped := make(map[string][]settings.Setting)
	var order []string
	for _, s := range settingsList {
		cat := s.Category()
		if _, exists := grouped[cat]; !exists {
			order = append(order, cat)
		}
		grouped[cat] = append(grouped[cat], s)
	}

	sort.Strings(order)
	for cat := range grouped {
		sort.Slice(grouped[cat], func(i, j int) bool {
			return grouped[cat][i].Label() < grouped[cat][j].Label()
		})
	}

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(" Settings ")
	list.SetTitleAlign(tview.AlignLeft)
	list.SetBorderPadding(1, 1, 2, 2)

	nonSelectable := make(map[int]bool)

	for _, cat := range order {
		headerIndex := list.GetItemCount()
		list.AddItem(fmt.Sprintf("[::b]%s[::-]", cat), "", 0, nil)
		nonSelectable[headerIndex] = true

		for _, s := range grouped[cat] {
			var secondary string
			if s.Type() == "static" {
				secondary = s.Get()
			} else {
				secondary = "Current: " + s.Get()
			}
			list.AddItem(fmt.Sprintf("  %s", s.Label()), secondary, 0, nil)
		}
	}

	list.AddItem("Back", "", 'b', goBack)

	list.SetSelectedFunc(func(index int, mainText, _ string, _ rune) {
		if nonSelectable[index] {
			return
		}

		label := strings.TrimSpace(stripFormatting(mainText))
		if label == "Back" {
			goBack()
			return
		}

		for _, s := range settingsList {
			if s.Label() == label {
				if s.Type() != "static" {
					editSetting(app, s, func() {
						app.SetRoot(RenderSettingsUI(app, db, goBack), true)
					})
				} else if s.Action() != nil {
					msg, err := s.Action()()

					// show either the message or error, if message show it and then exit
					if err != nil {
						showError(app, err.Error(), goBack)
						return
					}

					if msg != "" {
						modal := tview.NewModal().
							SetText(msg).
							AddButtons([]string{"OK"}).
							SetDoneFunc(func(_ int, _ string) {
								// exit cli fully
								app.Stop()
							})
						app.SetRoot(modal, true)
					}
				} else {
					showError(app, "No action defined for this setting.", goBack)
				}

				return
			}
		}
	})

	return list
}

func editSetting(app *tview.Application, setting settings.Setting, goBack func()) {
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
		form.AddInputField(setting.Label(), currentValue, 0, func(textToCheck string, _ rune) bool {
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
		showError(app, "Unknown setting type: "+setting.Type(), goBack)
		return
	}

	form.AddButton("Save", func() {
		setting.Set(currentValue)
		goBack()
	})
	form.AddButton("Cancel", func() {
		goBack()
	})

	app.SetRoot(form, true)
}

func showError(app *tview.Application, message string, goBack func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			goBack()
		})
	app.SetRoot(modal, true)
}

func stripFormatting(s string) string {
	var out strings.Builder
	skip := false
	for _, r := range s {
		if r == '[' {
			skip = true
			continue
		}
		if r == ']' {
			skip = false
			continue
		}
		if !skip {
			out.WriteRune(r)
		}
	}
	return out.String()
}
