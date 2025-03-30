package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"authservice/src/models" 
	"authservice/src/services"
)

type UserHandler struct {
	userService *services.UserService
	tokenService *services.TokenService
}

func NewUserHandler(us *services.UserService, ts *services.TokenService) *UserHandler {
	return &UserHandler{
		userService: us,
		tokenService: ts,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
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
	c.JSON(http.StatusOK, gin.H{"user": user})
}