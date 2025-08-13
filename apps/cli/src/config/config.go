package config

import (
	"TimeTrack-cli/src/database"
)

const DefaultServerURL = "https://timetrack.linusromland.com"

func InitializeDB(db *database.DBWrapper) error {
	if db == nil || db.DB == nil {
		return nil
	}

	if db.Get(database.ServerURLKey) == "" {
		// Set the default server URL
		db.Set(database.ServerURLKey, DefaultServerURL)
	}

	return nil
}
