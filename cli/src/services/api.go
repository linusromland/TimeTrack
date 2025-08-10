package services

import (
	"TimeTrack-cli/src/database"
	apiPkg "TimeTrack-cli/src/services/api"
)

type APIService = apiPkg.APIService

func NewAPIService(db *database.DBWrapper) *APIService {
	return apiPkg.NewAPIService(db)
}
