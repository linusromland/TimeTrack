package apiService

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"TimeTrack-shared/dtos"
	"TimeTrack-shared/models"
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
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("error closing response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create time entry: %s", resp.Status)
	}

	var createdEntry models.TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&createdEntry); err != nil {
		return nil, fmt.Errorf("failed to parse created time entry response: %w", err)
	}

	return &createdEntry, nil
}

func (api *APIService) GetTimeEntries(startDate, endDate string, page int) ([]*models.TimeEntry, error) {
	limit := 25
	skip := (page - 1) * limit

	reqURL := fmt.Sprintf("%s/time-entries?from=%s&to=%s&skip=%d&limit=%d", api.baseURL, startDate, endDate, skip, limit)

	req, err := api.newAuthRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get time entries: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("error closing response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get time entries: %s", resp.Status)
	}

	var entries []*models.TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, fmt.Errorf("failed to parse time entries response: %w", err)
	}

	return entries, nil
}

func (api *APIService) GetTimeEntryStatistics(startDate, endDate string) (*models.TimeEntryStatistics, error) {
	reqURL := fmt.Sprintf("%s/time-entries/statistics?from=%s&to=%s", api.baseURL, startDate, endDate)

	req, err := api.newAuthRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get time entries: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("error closing response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get time entries: %s", resp.Status)
	}

	var stats *models.TimeEntryStatistics
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to parse time entries response: %w", err)
	}

	return stats, nil
}

func (api *APIService) DeleteTimeEntry(timeEntryId string) error {
	reqURL := fmt.Sprintf("%s/time-entries/%s", api.baseURL, timeEntryId)

	req, err := api.newAuthRequest("DELETE", reqURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete time entry: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("error closing response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete time entry: %s", resp.Status)
	}

	return nil
}
