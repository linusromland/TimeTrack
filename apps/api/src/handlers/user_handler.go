package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"TimeTrack-api/src/services"
	"TimeTrack-shared/models"
)

type UserHandler struct {
	userService  *services.UserService
	tokenService *services.TokenService
}

func NewUserHandler(us *services.UserService, ts *services.TokenService) *UserHandler {
	return &UserHandler{
		userService:  us,
		tokenService: ts,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	existingUser, err := h.userService.GetUserByEmail(c, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking existing user"})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	if err := h.userService.RegisterUser(c, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var loginData models.User
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// log email and password for debugging purposes
	fmt.Printf("Login attempt with email: %s, password: %s\n", loginData.Email, loginData.Password)

	user, err := h.userService.LoginUser(c, &loginData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenString, err := h.tokenService.GenerateAuthToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing the token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.userService.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user"})
		return
	}
	c.JSON(http.StatusOK, user)
}
