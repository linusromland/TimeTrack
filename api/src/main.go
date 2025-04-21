package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"TimeTrack-api/src/config"
	"TimeTrack-api/src/database"
	"TimeTrack-api/src/handlers"
	"TimeTrack-api/src/middleware"
	"TimeTrack-api/src/services"
)

func main() {
	// Load and check configuration
	cfg := config.LoadConfig()
	config.CheckRequiredVariables(cfg)

	// Setup Google OAuth config
	config.SetupGoogleOAuthConfig(cfg)

	// Connect to MongoDB
	database.ConnectDB(cfg.MongoURI)
	defer database.DisconnectDB()

	// Initialize services
	userService := services.NewUserService(database.AuthDatabase)
	tokenService := services.NewTokenService(database.AuthDatabase, cfg.JWTSecret)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, tokenService)

	// Setup Gin router
	r := gin.Default()

	// Versioned API path `/api/v1`
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/register", userHandler.RegisterUser)
		apiV1.POST("/login", userHandler.LoginUser)

		// Authentication routes
		authGroup := apiV1.Group("/auth")
		authGroup.Use(middleware.AuthMiddleware())
		{
			authGroup.GET("/user", userHandler.GetUser)
		}
	}

	// Start server
	port := cfg.Port
	log.Printf("Server listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
