package app

import (
	"TimeTrack-cli/src/database"
	"TimeTrack-cli/src/services"
	"fmt"

	"github.com/urfave/cli/v2"
)

type AppContext struct {
	Version string
	DB      *database.DBWrapper
	API     *services.APIService
}

func NewAppContext(version string) *AppContext {
	return &AppContext{
		Version: version,
	}
}

func (a *AppContext) Startup(c *cli.Context) error {
	db, err := database.OpenDB()
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	a.DB = db

	a.API = services.NewAPIService(a.DB)

	if err := a.API.HealthCheck(); err != nil {
		return fmt.Errorf("API health check failed: %w", err)
	}

	return nil
}

func (a *AppContext) Shutdown() {
	if a.DB != nil {
		a.DB.Close()
	}
}
