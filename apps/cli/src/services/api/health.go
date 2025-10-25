package apiService

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type HealthResponse struct {
	OK      bool   `json:"ok"`
	Error   string `json:"error,omitempty"`
	Version string `json:"version"`
}

func (api *APIService) HealthCheck() (*HealthResponse, error) {
	healthURL := fmt.Sprintf("%s/health", api.baseURL)

	resp, err := api.client.Get(healthURL)
	if err != nil {
		return nil, fmt.Errorf("failed to reach API health endpoint: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("error closing response body: %v", cerr)
		}
	}()

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
