package services

import (
	"TimeTrack-api/src/config"
	"TimeTrack-api/src/models"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type AtlassianService struct {
	config      config.AtlassianConfig
	userService *UserService
}

func NewAtlassianService(c config.AtlassianConfig, us UserService) *AtlassianService {
	return &AtlassianService{
		config:      c,
		userService: &us,
	}
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

	c.JSON(200, gin.H{
		"oauth_url": oauthURL,
	})
}

func (s *AtlassianService) HandleOAuthCallback(c *gin.Context) {
	log.Println("Handling OAuth callback from Atlassian")

	code := c.Query("code")
	userId := c.Query("state")

	if code == "" || userId == "" {
		c.JSON(400, gin.H{"error": "Missing code or state"})
		return
	}

	_, err := s.userService.GetUserByID(c, userId)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	oauthTokenUrl := "https://auth.atlassian.com/oauth/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", s.config.ClientId)
	data.Set("client_secret", s.config.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", s.config.CallbackUrl)
	req, err := http.NewRequest("POST", oauthTokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to send request"})
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.JSON(500, gin.H{"error": "Failed to exchange code for token"})
		return
	}

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse response"})
		return
	}
	accessToken, ok := tokenResponse["access_token"].(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Access token not found in response"})
		return
	}

	s.userService.UpdateIntegration(c, userId, "atlassian", models.UserIntegration{
		Atlassian: models.AtlassianIntegration{
			Enabled:     true,
			AccessToken: accessToken,
		},
	})

	c.JSON(200, gin.H{
		"success": true,
		"message": "Authentication with Atlassian successful. You can now safely close this window.",
	})

}

// TODO: Add support fo specifiying what cloud the user wants to connect to, and also save the cloud id in the database
func (s *AtlassianService) GetCloudId(userId string) (string, error) {
	log.Println("Fetching cloud ID from Atlassian")

	atlassianIntegration, err := s.userService.GetAtlassianIntegration(userId)
	if err != nil {
		log.Println("Error fetching Atlassian integration:", err)
		return "", err
	}
	if !atlassianIntegration.Enabled {
		log.Println("Atlassian integration is not enabled for user:", userId)
		return "", nil
	}
	if atlassianIntegration.AccessToken == "" {
		log.Println("Access token is empty for user:", userId)
		return "", nil
	}

	// Make request to get the cloud ID
	cloudIdUrl := "https://api.atlassian.com/oauth/token/accessible-resources"
	req, err := http.NewRequest("GET", cloudIdUrl, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+atlassianIntegration.AccessToken)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("Error response from Atlassian:", resp.StatusCode)
		return "", err
	}
	var resources []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&resources); err != nil {
		log.Println("Error decoding response:", err)
		return "", err
	}
	// Extract the cloud ID from the response
	if len(resources) == 0 {
		log.Println("No resources found in response")
		return "", nil
	}

	log.Println("Resources found:", resources)

	cloudId, ok := resources[0]["id"].(string)
	if !ok {
		log.Println("Cloud ID not found in response")
		return "", nil
	}
	log.Println("Cloud ID found:", cloudId)
	return cloudId, nil
}

func (s *AtlassianService) CheckIfJiraTicketExists(userId string, ticketId string) error {
	log.Println("Fetching ticket info from Atlassian")

	atlassianIntegration, err := s.userService.GetAtlassianIntegration(userId)
	if err != nil {
		log.Println("Error fetching Atlassian integration:", err)
		return err
	}
	if !atlassianIntegration.Enabled {
		log.Println("Atlassian integration is not enabled for user:", userId)
		return errors.New("atlassian integration is not enabled")
	}
	if atlassianIntegration.AccessToken == "" {
		log.Println("Access token is empty for user:", userId)
		return errors.New("access token is empty")
	}

	cloudId, err := s.GetCloudId(userId)
	if err != nil {
		log.Println("Error fetching cloud ID:", err)
		return err
	}

	jiraUrl := "https://api.atlassian.com/ex/jira/" + cloudId + "/rest/api/3/issue/" + ticketId
	req, err := http.NewRequest("GET", jiraUrl, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return err
	}
	req.Header.Set("Authorization", "Bearer "+atlassianIntegration.AccessToken)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("Error response from Atlassian:", resp.StatusCode)
		return errors.New("ticket not found")
	}

	var ticketInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ticketInfo); err != nil {
		log.Println("Error decoding response:", err)
		return err
	}

	if _, ok := ticketInfo["key"]; !ok {
		log.Println("Ticket not found:", ticketId)
		return errors.New("ticket not found")
	}
	log.Println("Ticket found:", ticketId)
	return nil
}

func (s *AtlassianService) AddTimeEntryToJira(entry *models.TimeEntry, ticketId string) (string, error) {
	log.Println("Adding time entry to Jira")

	userId := entry.OwnerID

	atlassianIntegration, err := s.userService.GetAtlassianIntegration(userId)
	if err != nil {
		log.Println("Error fetching Atlassian integration:", err)
		return "", err
	}
	if !atlassianIntegration.Enabled {
		log.Println("Atlassian integration is not enabled for user:", userId)
		return "", errors.New("atlassian integration is not enabled")
	}
	if atlassianIntegration.AccessToken == "" {
		log.Println("Access token is empty for user:", userId)
		return "", errors.New("access token is empty")
	}

	cloudId, err := s.GetCloudId(userId)
	if err != nil {
		log.Println("Error fetching cloud ID:", err)
		return "", err
	}

	jiraUrl := "https://api.atlassian.com/ex/jira/" + cloudId + "/rest/api/2/issue/" + ticketId + "/worklog"
	reqBody := map[string]interface{}{
		"comment":          entry.Note,
		"timeSpentSeconds": entry.Period.Duration,
	}
	reqBodyJson, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", jiraUrl, strings.NewReader(string(reqBodyJson)))
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+atlassianIntegration.AccessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Println("Error response from Atlassian:", resp.StatusCode)
		return "", errors.New("failed to add time entry to Jira")
	}

	var worklogResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&worklogResponse); err != nil {
		log.Println("Error decoding response:", err)
		return "", err
	}

	worklogId, ok := worklogResponse["id"].(string)
	if !ok {
		log.Println("Worklog ID not found in response")
		return "", errors.New("worklog ID not found")
	}

	log.Println("Worklog ID found:",
		worklogId)
	return worklogId, nil
}
