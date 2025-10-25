package apiService

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (api *APIService) GetAtlassianAuthURL() (string, error) {
	reqURL := fmt.Sprintf("%s/user/oauth/atlassian", api.baseURL)

	req, err := api.newAuthRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, _ := api.client.Do(req)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get Atlassian auth URL, status code: %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode Atlassian auth URL response: %w", err)
	}
	oauthURL, exists := result["oauth_url"]
	if !exists {
		return "", fmt.Errorf("oauth_url not found in response")
	}

	return oauthURL, nil
}
