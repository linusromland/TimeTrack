package settings

import (
	"TimeTrack/core/calendar"
	"TimeTrack/core/database"
	"fmt"

	badger "github.com/dgraph-io/badger/v4"
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

func getSettings(db *badger.DB) []SettingCategory {
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
					Label:   "Select Calendar",
					Options: availableCalendars,
					GetValue: func(value string) (string, string) {
						if value == "" {
							value = database.GetData(db, database.CALENDAR_ID)
						}

						for _, calendar := range availableCalendars {
							if calendar.Value == value {
								return calendar.Value, calendar.Label
							}
						}

						return "", ""
					},
					SetValue: func(value string) {
						database.SetData(db, database.CALENDAR_ID, value)
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
						value := database.GetData(db, database.CLOUD_SYNC_ENABLED)

						label := "Disabled"
						if value == "true" {
							label = "Enabled"
						}

						return value, label
					},
					SetValue: func(value string) {
						database.SetData(db, database.CLOUD_SYNC_ENABLED, value)
					},
				},
				{
					Id:    "url",
					Type:  "text",
					Label: "Cloud Sync URL",
					GetValue: func(_ string) (string, string) {
						value := database.GetData(db, database.CLOUD_SYNC_URL)

						return value, value
					},
					SetValue: func(value string) {
						database.SetData(db, database.CLOUD_SYNC_URL, value)
					},
				},
				{
					Id:    "interval",
					Type:  "number",
					Label: "Sync Interval",
					GetValue: func(_ string) (string, string) {
						value := database.GetData(db, database.CLOUD_SYNC_INTERVAL)

						if value == "" {
							return "", ""
						}

						label := "Every time"
						if value != "0" {
							label = fmt.Sprintf("Every %s hour(s)", value)
						}

						return value, label
					},
					SetValue: func(value string) {
						database.SetData(db, database.CLOUD_SYNC_INTERVAL, value)
					},
				},
			},
		},
	}

	return settings
}

func getSettingCategory(db *badger.DB, id string) SettingCategory {
	settings := getSettings(db)

	for _, settingCategory := range settings {
		if settingCategory.Id == id {
			return settingCategory
		}
	}

	return SettingCategory{}
}

func GetSettingsByCategory(db *badger.DB, category string) []Setting {
	// get the settingsCategory from the map
	settingsCategoryId := categorySettings[category]

	settingCategory := getSettingCategory(db, settingsCategoryId)

	return settingCategory.Settings
}
