package apiService

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"TimeTrack-cli/src/database"
	"TimeTrack-shared/models"
)

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
	reqURL := fmt.Sprintf("%s/user", api.baseURL)

	req, err := api.newAuthRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

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
