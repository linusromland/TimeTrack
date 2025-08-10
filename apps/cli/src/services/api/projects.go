package apiService

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"TimeTrack-shared/dtos"
	"TimeTrack-shared/models"
)

func (api *APIService) GetProjectByName(name string) (*models.Project, error) {
	reqURL := fmt.Sprintf("%s/projects?name=%s", api.baseURL, url.QueryEscape(name))

	req, err := api.newAuthRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get project by name: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get project by name: %s", resp.Status)
	}

	var projects []models.Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, fmt.Errorf("failed to parse projects response: %w", err)
	}
	if len(projects) == 0 {
		return nil, fmt.Errorf("no project found with name: %s", name)
	}
	return &projects[0], nil
}

func (api *APIService) CreateProject(project *dtos.CreateProjectInput) (*models.Project, error) {
	reqURL := fmt.Sprintf("%s/projects", api.baseURL)

	body, err := json.Marshal(project)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal project: %w", err)
	}

	req, err := api.newAuthRequest("POST", reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to create project: %s", resp.Status)
	}

	var createdProject models.Project
	if err := json.NewDecoder(resp.Body).Decode(&createdProject); err != nil {
		return nil, fmt.Errorf("failed to parse created project response: %w", err)
	}

	return &createdProject, nil
}
