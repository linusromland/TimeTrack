package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "TimeTrack-api/src/docs" // Import the generated docs

	"TimeTrack-api/src/config"
	"TimeTrack-api/src/database"
	"TimeTrack-api/src/handlers"
	"TimeTrack-api/src/middleware"
	"TimeTrack-api/src/services"
)

// @title TimeTrack API
// @version 1.0
// @description Time tracking API for managing projects and time entries.
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load and check configuration
	cfg := config.LoadConfig()
	config.CheckRequiredVariables(cfg)

	// Connect to MongoDB
	database.ConnectDB(cfg.MongoURI)
	defer database.DisconnectDB()

	// Initialize services
	userService := services.NewUserService(database.Database)
	tokenService := services.NewTokenService(database.Database, cfg.JWTSecret)
	atlassianService := services.NewAtlassianService(cfg.AtlassianConfig, *userService)
	projectService := services.NewProjectService(database.Database, atlassianService)
	timeEntryService := services.NewTimeEntryService(database.Database, projectService, atlassianService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, tokenService)
	projectHandler := handlers.NewProjectHandler(projectService)
	timeEntryHandler := handlers.NewTimeEntryHandler(timeEntryService, projectService)
	healthHandler := handlers.NewHealthHandler(database.Database, cfg.APIVersion)

	// Setup Gin router
	r := gin.Default()

	// Serve Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Versioned API path `/api/v1`
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/register", userHandler.RegisterUser)
		apiV1.POST("/login", userHandler.LoginUser)

		// Health check endpoint
		apiV1.GET("/health", healthHandler.CheckHealth)

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
			authGroup.GET("/time-entries/statistics", timeEntryHandler.Statistics)
		}
	}

	// Start server
	port := cfg.Port
	log.Printf("Server listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
