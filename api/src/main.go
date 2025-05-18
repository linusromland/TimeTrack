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

	// Connect to MongoDB
	database.ConnectDB(cfg.MongoURI)
	defer database.DisconnectDB()

	// Initialize services
	userService := services.NewUserService(database.AuthDatabase)
	tokenService := services.NewTokenService(database.AuthDatabase, cfg.JWTSecret)
	atlassianService := services.NewAtlassianService(cfg.AtlassianConfig, *userService)
	projectService := services.NewProjectService(database.AuthDatabase, atlassianService)
	timeEntryService := services.NewTimeEntryService(database.AuthDatabase, projectService, atlassianService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, tokenService)
	projectHandler := handlers.NewProjectHandler(projectService)
	timeEntryHandler := handlers.NewTimeEntryHandler(timeEntryService, projectService)

	// Setup Gin router
	r := gin.Default()

	// Versioned API path `/api/v1`
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/register", userHandler.RegisterUser)
		apiV1.POST("/login", userHandler.LoginUser)

		// Authentication routes
		authGroup := apiV1.Group("/")

		// Oauth Callback routes
		authGroup.GET("/user/oauth/atlassian/callback", atlassianService.HandleOAuthCallback)


		authGroup.Use(middleware.AuthMiddleware())
		{
			userGroup := authGroup.Group("/user")
			{
				userGroup.GET("/", userHandler.GetUser)

				oauthGroup := userGroup.Group("/oauth")
				{
					oauthGroup.GET("/atlassian", atlassianService.GetOAuthURL)
				}
			}

			// Project routes
			authGroup.POST("/projects", projectHandler.Create)
			authGroup.PUT("/projects/:id", projectHandler.Update)
			authGroup.DELETE("/projects/:id", projectHandler.Delete)
			authGroup.GET("/projects", projectHandler.List)

			// Time Entry routes
			authGroup.POST("/time-entries", timeEntryHandler.Create)
			authGroup.PUT("/time-entries/:id", timeEntryHandler.Update)
			authGroup.DELETE("/time-entries/:id", timeEntryHandler.Delete)
			authGroup.GET("/time-entries", timeEntryHandler.List)
		}
	}

	// Start server
	port := cfg.Port
	log.Printf("Server listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
