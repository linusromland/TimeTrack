package apiService

import (
	"encoding/json"
	"fmt"
	"net/http"

	"TimeTrack-cli/src/dtos"
	"TimeTrack-cli/src/models"
)

func (api *APIService) CreateTimeEntry(entry *dtos.CreateTimeEntryInput) (*models.TimeEntry, error) {
	reqURL := fmt.Sprintf("%s/time-entries", api.baseURL)

	body, err := json.Marshal(entry)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal time entry: %w", err)
	}

	req, err := api.newAuthRequest("POST", reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create time entry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create time entry: %s", resp.Status)
	}

	var createdEntry models.TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&createdEntry); err != nil {
		return nil, fmt.Errorf("failed to parse created time entry response: %w", err)
	}

	return &createdEntry, nil
}
