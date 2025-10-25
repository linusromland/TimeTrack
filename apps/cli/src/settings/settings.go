package settings

import (
	"TimeTrack-cli/src/config"
	"TimeTrack-cli/src/database"
	"TimeTrack-cli/src/services"
	"TimeTrack-cli/src/utils"
	"fmt"
)

type Setting interface {
	ID() string
	Label() string
	Get() string
	Set(string) error
	Type() string
	Category() string
	Action() func() (string, error)
}

type TextSetting struct {
	id       string
	label    string
	category string
	db       *database.DBWrapper
	key      string
	def      string
}

func (s TextSetting) ID() string       { return s.id }
func (s TextSetting) Label() string    { return s.label }
func (s TextSetting) Type() string     { return "text" }
func (s TextSetting) Category() string { return s.category }
func (s TextSetting) Get() string {
	val := s.db.Get(s.key)
	if val == "" {
		return s.def
	}
	return val
}
func (s TextSetting) Set(v string) error { return s.db.Set(s.key, v) }
func (s TextSetting) Action() func() (string, error) {
	return nil // No action for text settings
}

type StaticSetting struct {
	id       string
	label    string
	category string
	getter   func() string
	action   func() (string, error)
}

func (s StaticSetting) ID() string       { return s.id }
func (s StaticSetting) Label() string    { return s.label }
func (s StaticSetting) Type() string     { return "static" }
func (s StaticSetting) Category() string { return s.category }
func (s StaticSetting) Get() string      { return s.getter() }
func (s StaticSetting) Set(_ string) error {
	return nil
}
func (s StaticSetting) Action() func() (string, error) {
	if s.action != nil {
		return s.action
	}
	return nil
}

func GetAllSettings(db *database.DBWrapper) []Setting {
	api := services.NewAPIService(db)

	return []Setting{
		StaticSetting{
			id:       "server_health",
			label:    "Server Health",
			category: "Server Info",
			getter: func() string {
				health, err := api.HealthCheck()
				if err != nil {
					return fmt.Sprintf("Unhealthy - %s", err.Error())
				}
				return fmt.Sprintf("Healthy - Version: %s", health.Version)
			},
		},
		StaticSetting{
			id:       "user",
			label:    "User",
			category: "Server Info",
			getter: func() string {
				user, err := api.GetCurrentUser()
				if err != nil {
					return "Unauthorized / Not logged in"
				}
				return fmt.Sprintf("Logged in as: %s", user.Email)
			},
		},
		TextSetting{
			id:       "url",
			label:    "Server URL",
			category: "General",
			db:       db,
			key:      database.ServerURLKey,
			def:      config.DefaultServerURL,
		},
		StaticSetting{
			id:       "atlassian_integration_status",
			label:    "Atlassian Status",
			category: "Integrations",
			getter: func() string {
				user, err := api.GetCurrentUser()
				if err != nil {
					return "Disabled"
				}
				if user.Integration.Atlassian.Enabled {
					return "Enabled"
				}
				return "Disabled"
			},
		},
		StaticSetting{
			id:       "atlassian_authenticate",
			label:    "Authenticate Atlassian",
			category: "Integrations",
			getter: func() string {
				return "Click to start authentication"
			},
			action: func() (string, error) {
				url, err := api.GetAtlassianAuthURL()
				if err != nil {
					return "", err
				}
				err = utils.OpenBrowser(url)
				if err != nil {
					return "", err
				}
				return "Follow the instructions in your browser to complete authentication.", nil
			},
		},
	}
}
