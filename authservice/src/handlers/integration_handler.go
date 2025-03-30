package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"authservice/src/services" 
)

type IntegrationHandler struct {
	userService *services.UserService
}

func NewIntegrationHandler(us *services.UserService) *IntegrationHandler {
	return &IntegrationHandler{
		userService: us,
	}
}

func (h *IntegrationHandler) ValidateIntegrationToken(c *gin.Context) {
	// Integration token is already validated by the middleware
	userToken := c.DefaultQuery("user_token", "")
	if userToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User token is required in query parameters"})
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	parsedUserToken, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !parsedUserToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired user token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Integration token and user token are valid"})
}

func (h *IntegrationHandler) GetUserForIntegration(c *gin.Context) {
	userToken := c.DefaultQuery("user_token", "")
	if userToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User token is required in query parameters"})
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	parsedUserToken, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !parsedUserToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired user token"})
		return
	}

	claims, ok := parsedUserToken.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user token claims"})
		return
	}

	userID := claims["user_id"].(string)
	user, err := h.userService.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}