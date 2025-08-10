package settings

import (
	"TimeTrack-cli/src/config"
	"TimeTrack-cli/src/database"
	services "TimeTrack-cli/src/services/api"
	"fmt"
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

type StaticSetting struct {
	id     string
	label  string
	getter func() string
}

func (s StaticSetting) ID() string    { return s.id }
func (s StaticSetting) Label() string { return s.label }
func (s StaticSetting) Type() string  { return "static" }
func (s StaticSetting) Get() string   { return s.getter() }
func (s StaticSetting) Set(_ string) error {
	// Static values can't be set
	return nil
}

type SettingCategory struct {
	ID       string
	Label    string
	Settings []Setting
}

func GetAllSettings(db *database.DBWrapper) []Setting {
	api := services.NewAPIService(db)

	return []Setting{
		StaticSetting{
			id:    "server_health",
			label: "Server Health",
			getter: func() string {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Recovered in getter: %v\n", r)
					}
				}()

				health, err := api.HealthCheck()
				if err != nil {
					return fmt.Sprintf("Unhealthy - %s", err.Error())
				}
				return fmt.Sprintf("Healthy - Version: %s", health.Version)
			},
		},
		StaticSetting{
			id:    "user",
			label: "User",
			getter: func() string {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Recovered in getter: %v\n", r)
					}
				}()

				user, err := api.GetCurrentUser()
				if err != nil {
					return "Unauthorized / Not logged in"
				}
				return fmt.Sprintf("Logged in as: %s", user.Email)
			},
		},
		TextSetting{
			id:    "url",
			label: "Server URL",
			db:    db,
			key:   database.ServerURLKey,
			def:   config.DefaultServerURL,
		},
	}
}
