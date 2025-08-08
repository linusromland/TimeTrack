package settings

import (
	"TimeTrack-cli/src/database"

	badger "github.com/dgraph-io/badger/v4"
)

const (
	SERVER_URL = "https://timetrack.linusromland.com"
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
	"Server": "server",
}

func getSettings(db *badger.DB) []SettingCategory {
	settings := []SettingCategory{
		{
			Id:    "server",
			Label: "Server",
			Settings: []Setting{
				{
					Id:    "url",
					Type:  "text",
					Label: "Server URL",
					GetValue: func(_ string) (string, string) {
						value := database.GetData(db, database.SERVER_URL)

						if value == "" {
							value = SERVER_URL // Default value if not set
						}

						return value, value
					},
					SetValue: func(value string) {
						database.SetData(db, database.SERVER_URL, value)
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
	settingsCategoryId := categorySettings[category]

	settingCategory := getSettingCategory(db, settingsCategoryId)

	return settingCategory.Settings
}
