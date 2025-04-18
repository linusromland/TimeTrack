package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"authservice/src/services"
)

type AuthHandler struct {
	tokenService *services.TokenService
}

func NewAuthHandler(ts *services.TokenService) *AuthHandler {
	return &AuthHandler{
		tokenService: ts,
	}
}

func (h *AuthHandler) GenerateAPIToken(c *gin.Context) {
	var req struct {
		Expiry int `json:"expiry"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID := c.GetString("user_id")
	tokenString, _, err := h.tokenService.GenerateAPIToken(c, userID, req.Expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Token generated", "token": tokenString})
}

func (h *AuthHandler) ListUserTokens(c *gin.Context) {
	userID := c.GetString("user_id")
	tokens, err := h.tokenService.ListUserTokens(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tokens"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

func (h *AuthHandler) RevokeToken(c *gin.Context) {
	id := c.Param("id")
	err := h.tokenService.RevokeToken(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error revoking token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Token revoked"})
}
