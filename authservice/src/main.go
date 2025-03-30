package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"authservice/src/config"
	"authservice/src/database"
	"authservice/src/handlers"
	"authservice/src/middleware"
	"authservice/src/services"
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
	authHandler := handlers.NewAuthHandler(tokenService)
	integrationHandler := handlers.NewIntegrationHandler(userService)

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
			authGroup.POST("/token/generate", authHandler.GenerateAPIToken)
			authGroup.GET("/token/list", authHandler.ListUserTokens)
			authGroup.DELETE("/token/revoke/:id", authHandler.RevokeToken)
			authGroup.GET("/user", userHandler.GetUser)
		}

		integrationGroup := apiV1.Group("/integration")
		integrationGroup.Use(middleware.IntegrationAuthMiddleware())
		{
			integrationGroup.POST("/validate", integrationHandler.ValidateIntegrationToken)
			integrationGroup.GET("/user", integrationHandler.GetUserForIntegration)
		}
	}

	// Start server
	port := cfg.Port
	log.Printf("Server listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
