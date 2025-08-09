package services

import "TimeTrack-cli/src/database"

type APIService struct {
	db *database.DBWrapper
}

func NewAPIService(db *database.DBWrapper) *APIService {
	return &APIService{db: db}
}

// Placeholder for future API health check
func (api *APIService) HealthCheck() error {
	return nil
}
