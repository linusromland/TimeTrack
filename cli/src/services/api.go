package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"TimeTrack-cli/src/database"
	"TimeTrack-cli/src/models"
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

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (api *APIService) Register(email, password string) error {
	url := fmt.Sprintf("%s/register", api.baseURL)
	body, _ := json.Marshal(AuthPayload{Email: email, Password: password})

	resp, err := api.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed: %s", resp.Status)
	}
	return nil
}

func (api *APIService) Login(email, password string) error {
	url := fmt.Sprintf("%s/login", api.baseURL)
	body, _ := json.Marshal(AuthPayload{Email: email, Password: password})

	resp, err := api.client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	var response struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}
	if response.Token == "" {
		return errors.New("login response did not contain a token")
	}
	api.db.Set(database.AuthTokenKey, response.Token)

	return nil
}

func (api *APIService) GetCurrentUser() (*models.User, error) {
	url := fmt.Sprintf("%s/user", api.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	token := api.db.Get(database.AuthTokenKey)
	if token == "" {
		return nil, errors.New("no authentication token found")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get current user: %s", resp.Status)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	return &user, nil
}