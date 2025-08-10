package app

import (
	"TimeTrack-cli/src/config"
	"TimeTrack-cli/src/database"
	services "TimeTrack-cli/src/services/api"
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

	// Initialize the database with default values
	if err := config.InitializeDB(a.DB); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	

	a.API = services.NewAPIService(a.DB)

	healthResponse, err := a.API.HealthCheck()
	if err != nil {
		fmt.Printf("Warning: API health check failed: %s\n", err)
	}

	if healthResponse != nil {
		if !healthResponse.OK {
			fmt.Printf("API is not healthy: %s\n", healthResponse.Error)
		}

		if healthResponse.Version != a.Version {
			fmt.Printf("Warning: API version mismatch. CLI version: %s, API version: %s\n", a.Version, healthResponse.Version)
		}
	}

	return nil
}

func (a *AppContext) Shutdown() {
	if a.DB != nil {
		a.DB.Close()
	}
}
