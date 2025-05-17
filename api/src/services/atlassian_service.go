package services

import (
	"TimeTrack-api/src/config"
	"TimeTrack-api/src/models"
	"encoding/json"
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

func (s *AtlassianService) CheckIfJiraTicketExists(userId string, ticketId string) error {
	log.Println("Fetching ticket info from Atlassian")

	atlassianIntegration, err := s.userService.GetAtlassianIntegration(userId)
	if err != nil {
		log.Println("Error fetching Atlassian integration:", err)
		return err
	}
	if !atlassianIntegration.Enabled {
		log.Println("Atlassian integration is not enabled for user:", userId)
		return nil
	}
	if atlassianIntegration.AccessToken == "" {
		log.Println("Access token is empty for user:", userId)
		return nil
	}

	jiraUrl := "https://api.atlassian.com/ex/jira/" + s.config.Audience + "/rest/api/2/issue/" + ticketId
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
		return err
	}

	var ticketInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ticketInfo); err != nil {
		log.Println("Error decoding response:", err)
		return err
	}

	// Check if the ticket exists
	if _, ok := ticketInfo["key"]; !ok {
		log.Println("Ticket not found:", ticketId)
		return err
	}
	log.Println("Ticket found:", ticketId)
	return nil
}
