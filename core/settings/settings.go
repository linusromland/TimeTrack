package settings

import (
	"TimeTrack/core/calendar"
	"fmt"
)

type Option struct {
	Value string
	Label string
}

type Setting struct {
	Id       string
	Type     string
	Label    string
	Options  []Option                            // For selects
	GetValue func(value string) (string, string) // Returns value and label. value is used for selects, for all others, leave blank
	SetValue func(string)                        // Sets value
}

type SettingCategory struct {
	Id       string
	Label    string
	Settings []Setting
}

var categorySettings = map[string]string{
	"Calendar":   "calendar",
	"Cloud Sync": "cloudSync",
}

func GetSettings() []SettingCategory {
	calendars := calendar.GetCalendars()
	availableCalendars := []Option{}

	for _, cal := range calendars.Items {
		availableCalendars = append(availableCalendars, Option{Value: cal.Id, Label: cal.Summary})
	}

	settings := []SettingCategory{
		{
			Id:    "calendar",
			Label: "Calendar",
			Settings: []Setting{
				{
					Id:      "calendar",
					Type:    "select",
					Label:   "Selected Calendar",
					Options: availableCalendars,
					GetValue: func(value string) (string, string) {
						for _, cal := range availableCalendars {
							if cal.Value == value {
								return cal.Value, cal.Label
							}
						}

						return "", ""
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

	return settings
}

func getSettingCategory(id string) SettingCategory {
	settings := GetSettings()

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
