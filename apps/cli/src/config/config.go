package config

import (
	"TimeTrack-cli/src/database"
	"fmt"
)

const DefaultServerURL = "https://timetrack.linusromland.com"

func InitializeDB(db *database.DBWrapper) error {
	if db == nil || db.DB == nil {
		return nil
	}

	if db.Get(database.ServerURLKey) == "" {
		// Set the default server URL
		err := db.Set(database.ServerURLKey, DefaultServerURL)
		if err != nil {
			return fmt.Errorf("failed to set default server URL: %w", err)
		}
	}

	return nil
}
