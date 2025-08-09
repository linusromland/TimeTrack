package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"TimeTrack-cli/src/database"
)

type APIService struct {
	db      *database.DBWrapper
	baseURL string
	client  *http.Client
}

type HealthResponse struct {
	OK      bool   `json:"ok"`
	Error   string `json:"error,omitempty"`
	Version string `json:"version"`
}

func NewAPIService(db *database.DBWrapper) *APIService {
	if db == nil {
		panic("database wrapper cannot be nil")
	}

	rawURL := db.Get(database.ServerURLKey)
	if rawURL == "" {
		panic("server URL cannot be empty")
	}

	normalized := normalizeBaseURL(rawURL)

	return &APIService{
		db:      db,
		baseURL: normalized,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func normalizeBaseURL(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		panic(fmt.Sprintf("invalid server URL: %q", raw))
	}

	cleanPath := strings.TrimRight(parsed.Path, "/")
	parsed.Path = fmt.Sprintf("%s/api/v1", cleanPath)

	return parsed.String()
}

func (api *APIService) GetBaseURL() string {
	return api.baseURL
}

func (api *APIService) HealthCheck() (*HealthResponse, error) {
	healthURL := fmt.Sprintf("%s/health", api.baseURL)

	resp, err := api.client.Get(healthURL)
	if err != nil {
		return nil, fmt.Errorf("failed to reach API health endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("health check returned status %d", resp.StatusCode)
	}

	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("failed to parse health response: %w", err)
	}

	if !health.OK {
		return &health, errors.New(health.Error)
	}

	return &health, nil
}
