package settings

import (
	"TimeTrack-cli/src/config"
	"TimeTrack-cli/src/database"
)

type Setting interface {
	ID() string
	Label() string
	Get() string
	Set(string) error
	Type() string
}

type TextSetting struct {
	id    string
	label string
	db    *database.DBWrapper
	key   string
	def   string
}

func (s TextSetting) ID() string    { return s.id }
func (s TextSetting) Label() string { return s.label }
func (s TextSetting) Type() string  { return "text" }
func (s TextSetting) Get() string {
	val := s.db.Get(s.key)
	if val == "" {
		return s.def
	}
	return val
}
func (s TextSetting) Set(v string) error { return s.db.Set(s.key, v) }

type SettingCategory struct {
	ID       string
	Label    string
	Settings []Setting
}

func GetAllSettings(db *database.DBWrapper) []Setting {
	return []Setting{
		TextSetting{
			id:    "url",
			label: "Server URL",
			db:    db,
			key:   database.ServerURLKey,
			def:   config.DefaultServerURL,
		},
	}
}
