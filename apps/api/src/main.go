package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

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

	// Serve static files for React SPA
	staticPath := "./static"
	if _, err := os.Stat(staticPath); err == nil {
		// Serve static assets (CSS, JS, images, etc.)
		r.Static("/assets", filepath.Join(staticPath, "assets"))
		
		// Serve other static files
		r.StaticFS("/static", http.Dir(staticPath))
		
		// Health check and Swagger before SPA fallback
		r.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}

	// Serve Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Versioned API path `/api/v1`
	apiV1 := r.Group("/api/v1")
	{
		apiV1.POST("/register", userHandler.RegisterUser)
		apiV1.POST("/login", userHandler.LoginUser)
		apiV1.POST("/logout", userHandler.LogoutUser)

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

	// SPA fallback - serve index.html for all non-API routes
	r.NoRoute(func(c *gin.Context) {
		// Check if this is an API request
		if c.Request.URL.Path[:4] == "/api" || c.Request.URL.Path[:8] == "/swagger" {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}
		
		// Serve React SPA for all other routes
		staticPath := "./static"
		indexPath := filepath.Join(staticPath, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			c.File(indexPath)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		}
	})

	// Start server
	port := cfg.Port
	log.Printf("Server listening on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
