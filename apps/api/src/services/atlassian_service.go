package services

import (
	"TimeTrack-api/src/config"
	"TimeTrack-shared/models"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"strings"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type AtlassianAPIErr struct {
	Error         string   `json:"error"`
	ErrorMessage  string   `json:"errorMessage"`
	ErrorMessages []string `json:"errorMessages"`
}

type AtlassianService struct {
	config      config.AtlassianConfig
	userService *UserService
	httpClient  *http.Client
}

func NewAtlassianService(c config.AtlassianConfig, us UserService) *AtlassianService {
	return &AtlassianService{
		config:      c,
		userService: &us,
		httpClient:  &http.Client{},
	}
}

func (s *AtlassianService) makeAtlassianRequest(method, reqURL string, accessToken string, reqBody interface{}, respTarget interface{}) error {
	var bodyReader io.Reader
	if reqBody != nil {
		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			log.Printf("Error marshaling request body: %v", err)
			return errors.New("failed to marshal request body")
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	} else {
		bodyReader = nil
	}

	req, err := http.NewRequest(method, reqURL, bodyReader)
	if err != nil {
		log.Printf("Error creating HTTP request for %s %s: %v", method, reqURL, err)
		return errors.New("failed to create HTTP request")
	}

	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	req.Header.Set("Accept", "application/json")
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("Error sending HTTP request to %s: %v", reqURL, err)
		return errors.New("failed to send HTTP request")
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("error closing response body: %v", cerr)
		}
	}()

	if resp.StatusCode >= 400 {
		var apiErr AtlassianAPIErr
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Atlassian API error response: Status %d, Body: %s", resp.StatusCode, string(bodyBytes))

		if err := json.Unmarshal(bodyBytes, &apiErr); err == nil {
			if apiErr.Error != "" {
				return errors.New("Atlassian API error: " + apiErr.Error + " (Status: " + resp.Status + ")")
			}
			if apiErr.ErrorMessage != "" {
				return errors.New("Atlassian API error: " + apiErr.ErrorMessage + " (Status: " + resp.Status + ")")
			}
			if len(apiErr.ErrorMessages) > 0 {
				return errors.New("Atlassian API errors: " + strings.Join(apiErr.ErrorMessages, "; ") + " (Status: " + resp.Status + ")")
			}
		}
		return errors.New("Atlassian API returned status: " + resp.Status)
	}

	if respTarget != nil {
		if err := json.NewDecoder(resp.Body).Decode(respTarget); err != nil {
			log.Printf("Error decoding Atlassian API response from %s: %v", reqURL, err)
			return errors.New("failed to decode Atlassian API response")
		}
	}
	return nil
}

func (s *AtlassianService) GetOAuthURL(c *gin.Context) {
	log.Println("Generating OAuth URL for Atlassian")

	userId := c.GetString("user_id")

	oauthURL := "https://auth.atlassian.com/authorize?" +
		"audience=" + s.config.Audience +
		"&client_id=" + s.config.ClientId +
		"&scope=" + s.config.Scope +
		"&redirect_uri=" + s.config.CallbackUrl +
		"&state=" + userId +
		"&response_type=code" +
		"&prompt=consent"

	log.Println("Generated OAuth URL:", oauthURL)

	c.JSON(http.StatusOK, gin.H{
		"oauth_url": oauthURL,
	})
}

func (s *AtlassianService) HandleOAuthCallback(c *gin.Context) {
	log.Println("Handling OAuth callback from Atlassian")

	code := c.Query("code")
	userId := c.Query("state")

	if code == "" || userId == "" {
		log.Println("Missing code or state in OAuth callback")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code or state"})
		return
	}

	_, err := s.userService.GetUserByID(c, userId)
	if err != nil {
		log.Printf("User not found for OAuth callback (userId: %s): %v", userId, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	oauthTokenUrl := "https://auth.atlassian.com/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", s.config.ClientId)
	data.Set("client_secret", s.config.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", s.config.CallbackUrl)

	req, err := http.NewRequest(http.MethodPost, oauthTokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		log.Printf("Error creating token exchange request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token request"})
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("Error sending token exchange request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send token request"})
		return
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("error closing response body: %v", cerr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Error response from Atlassian token endpoint: Status %d, Body: %s", resp.StatusCode, string(bodyBytes))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		log.Printf("Error parsing token response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse token response"})
		return
	}

	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		log.Println("Access token not found in response or not a string type assertion failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Access token not found in response"})
		return
	}

	err = s.userService.UpdateIntegration(c, userId, "atlassian", models.UserIntegration{
		Atlassian: models.AtlassianIntegration{
			Enabled:     true,
			AccessToken: accessToken,
		},
	})
	if err != nil {
		log.Printf("Error updating Atlassian integration for user %s: %v", userId, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save Atlassian integration"})
		return
	}

	log.Printf("Authentication with Atlassian successful for user: %s", userId)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Authentication with Atlassian successful. You can now safely close this window.",
	})
}

func (s *AtlassianService) GetCloudId(userId string) (string, error) {
	log.Printf("Fetching cloud ID from Atlassian for user: %s", userId)

	atlassianIntegration, err := s.userService.GetAtlassianIntegration(userId)
	if err != nil {
		log.Printf("Error fetching Atlassian integration for user %s: %v", userId, err)
		return "", err
	}
	if !atlassianIntegration.Enabled {
		return "", errors.New("atlassian integration not enabled for user: " + userId)
	}
	if atlassianIntegration.AccessToken == "" {
		return "", errors.New("atlassian access token is empty for user: " + userId)
	}

	cloudIdUrl := "https://api.atlassian.com/oauth/token/accessible-resources"
	var resources []map[string]interface{}
	err = s.makeAtlassianRequest(http.MethodGet, cloudIdUrl, atlassianIntegration.AccessToken, nil, &resources)
	if err != nil {
		log.Printf("Error making request to get accessible resources for user %s: %v", userId, err)
		return "", err
	}

	if len(resources) == 0 {
		return "", errors.New("no accessible Atlassian resources found for user: " + userId)
	}

	log.Printf("Accessible resources found for user %s: %+v", userId, resources)

	cloudId, ok := resources[0]["id"].(string)
	if !ok {
		return "", errors.New("cloud ID not found in the first accessible resource or not a string")
	}

	log.Printf("Cloud ID found for user %s: %s", userId, cloudId)
	return cloudId, nil
}

func (s *AtlassianService) CheckIfJiraTicketExists(userId string, ticketId string) error {
	log.Printf("Checking Jira ticket existence for user: %s, ticket: %s", userId, ticketId)

	atlassianIntegration, err := s.userService.GetAtlassianIntegration(userId)
	if err != nil {
		log.Printf("Error fetching Atlassian integration for user %s: %v", userId, err)
		return err
	}
	if !atlassianIntegration.Enabled {
		return errors.New("atlassian integration is not enabled for user: " + userId)
	}
	if atlassianIntegration.AccessToken == "" {
		return errors.New("access token is empty for user: " + userId)
	}

	cloudId, err := s.GetCloudId(userId)
	if err != nil {
		log.Printf("Error fetching cloud ID for user %s: %v", userId, err)
		return err
	}
	if cloudId == "" {
		return errors.New("could not determine Atlassian cloud ID for user: " + userId)
	}

	jiraUrl := "https://api.atlassian.com/ex/jira/" + cloudId + "/rest/api/3/issue/" + ticketId
	var ticketInfo map[string]interface{}
	err = s.makeAtlassianRequest(http.MethodGet, jiraUrl, atlassianIntegration.AccessToken, nil, &ticketInfo)
	if err != nil {
		log.Printf("Error checking Jira ticket %s for user %s: %v", ticketId, userId, err)
		if strings.Contains(err.Error(), "Status: 404 Not Found") {
			return errors.New("jira ticket not found: " + ticketId)
		}
		return err
	}

	if _, ok := ticketInfo["key"]; !ok {
		log.Printf("Ticket key not found in response for Jira ticket: %s", ticketId)
		return errors.New("jira ticket not found: " + ticketId)
	}

	log.Printf("Jira ticket found: %s for user: %s", ticketId, userId)
	return nil
}

func (s *AtlassianService) AddTimeEntryToJira(entry *models.TimeEntry, ticketId string) (string, error) {
	log.Printf("Adding time entry to Jira ticket: %s for owner: %s", ticketId, entry.OwnerID)

	userId := entry.OwnerID

	atlassianIntegration, err := s.userService.GetAtlassianIntegration(userId)
	if err != nil {
		log.Printf("Error fetching Atlassian integration for user %s: %v", userId, err)
		return "", err
	}
	if !atlassianIntegration.Enabled {
		return "", errors.New("atlassian integration is not enabled for user: " + userId)
	}
	if atlassianIntegration.AccessToken == "" {
		return "", errors.New("access token is empty for user: " + userId)
	}

	cloudId, err := s.GetCloudId(userId)
	if err != nil {
		log.Printf("Error fetching cloud ID for user %s: %v", userId, err)
		return "", err
	}
	if cloudId == "" {
		return "", errors.New("could not determine Atlassian cloud ID for user: " + userId)
	}

	jiraUrl := "https://api.atlassian.com/ex/jira/" + cloudId + "/rest/api/2/issue/" + ticketId + "/worklog"

	reqBody := map[string]interface{}{
		"comment":          entry.Note,
		"timeSpentSeconds": entry.Period.Duration,
	}

	var worklogResponse map[string]interface{}
	err = s.makeAtlassianRequest(http.MethodPost, jiraUrl, atlassianIntegration.AccessToken, reqBody, &worklogResponse)
	if err != nil {
		log.Printf("Error adding time entry to Jira ticket %s for user %s: %v", ticketId, userId, err)
		return "", errors.New("failed to add time entry to Jira: " + err.Error())
	}

	worklogId, ok := worklogResponse["id"].(string)
	if !ok {
		log.Println("Worklog ID not found in response or not a string type assertion failed")
		return "", errors.New("worklog ID not found in Jira response")
	}

	log.Printf("Worklog ID %s added to Jira ticket %s for user %s", worklogId, ticketId, userId)
	return worklogId, nil
}

func (s *AtlassianService) UpdateTimeEntryInJira(entry *models.TimeEntry, ticketId string, worklogId string) (string, error) {
	log.Printf("Updating time entry %s in Jira ticket: %s for owner: %s", worklogId, ticketId, entry.OwnerID)

	userId := entry.OwnerID

	atlassianIntegration, err := s.userService.GetAtlassianIntegration(userId)
	if err != nil {
		log.Printf("Error fetching Atlassian integration for user %s: %v", userId, err)
		return "", err
	}
	if !atlassianIntegration.Enabled {
		return "", errors.New("atlassian integration is not enabled for user: " + userId)
	}
	if atlassianIntegration.AccessToken == "" {
		return "", errors.New("access token is empty for user: " + userId)
	}

	cloudId, err := s.GetCloudId(userId)
	if err != nil {
		log.Printf("Error fetching cloud ID for user %s: %v", userId, err)
		return "", err
	}
	if cloudId == "" {
		return "", errors.New("could not determine Atlassian cloud ID for user: " + userId)
	}
	jiraUrl := "https://api.atlassian.com/ex/jira/" + cloudId + "/rest/api/2/issue/" + ticketId + "/worklog/" + worklogId

	reqBody := map[string]interface{}{
		"comment":          entry.Note,
		"timeSpentSeconds": entry.Period.Duration,
	}

	var worklogResponse map[string]interface{}
	err = s.makeAtlassianRequest(http.MethodPut, jiraUrl, atlassianIntegration.AccessToken, reqBody, &worklogResponse)
	if err != nil {
		log.Printf("Error updating time entry %s in Jira ticket %s: %v", worklogId, ticketId, err)
		return "", errors.New("failed to update time entry in Jira: " + err.Error())
	}

	updatedWorklogId, ok := worklogResponse["id"].(string)
	if !ok {
		log.Println("Updated worklog ID not found in response or type assertion failed")
		return "", errors.New("updated worklog ID not found in Jira response")
	}

	log.Printf("Worklog ID %s successfully updated in Jira ticket %s", updatedWorklogId, ticketId)
	return updatedWorklogId, nil
}

func (s *AtlassianService) RemoveTimeEntryFromJira(userId string, ticketId string, worklogId string) error {
	log.Printf("Removing time entry %s from Jira ticket: %s for owner: %s", worklogId, ticketId, userId)

	atlassianIntegration, err := s.userService.GetAtlassianIntegration(userId)
	if err != nil {
		log.Printf("Error fetching Atlassian integration for user %s: %v", userId, err)
		return err
	}
	if !atlassianIntegration.Enabled {
		return errors.New("atlassian integration is not enabled for user: " + userId)
	}
	if atlassianIntegration.AccessToken == "" {
		return errors.New("access token is empty for user: " + userId)
	}

	cloudId, err := s.GetCloudId(userId)
	if err != nil {
		log.Printf("Error fetching cloud ID for user %s: %v", userId, err)
		return err
	}
	if cloudId == "" {
		return errors.New("could not determine Atlassian cloud ID for user: " + userId)
	}

	jiraUrl := "https://api.atlassian.com/ex/jira/" + cloudId + "/rest/api/2/issue/" + ticketId + "/worklog/" + worklogId

	err = s.makeAtlassianRequest(http.MethodDelete, jiraUrl, atlassianIntegration.AccessToken, nil, nil)
	if err != nil {
		log.Printf("Error removing time entry %s from Jira ticket %s: %v", worklogId, ticketId, err)
		return errors.New("failed to remove time entry from Jira: " + err.Error())
	}

	log.Printf("Worklog ID %s successfully removed from Jira ticket %s", worklogId, ticketId)
	return nil
}
