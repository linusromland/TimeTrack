package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	MongoURI          string
	Port              string
	JWTSecret         string
	IntegrationSecret string
	GoogleClientID    string
	GoogleClientSecret string
	GoogleRedirectURL   string
}

var AppConfig *Config
var GoogleOAuthConfig *oauth2.Config

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	cfg := &Config{
		MongoURI:          os.Getenv("MONGO_URI"),
		Port:              os.Getenv("PORT"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		IntegrationSecret: os.Getenv("INTEGRATION_SECRET"),
		GoogleClientID:    os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:   os.Getenv("GOOGLE_REDIRECT_URL"),
	}
	AppConfig = cfg
	return cfg
}

func SetupGoogleOAuthConfig(cfg *Config) *oauth2.Config {
	oauthConf := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	GoogleOAuthConfig = oauthConf
	return oauthConf
}

func CheckRequiredVariables(cfg *Config) {
	if cfg.MongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not defined")
	}
	if cfg.Port == "" {
		log.Println("PORT environment variable is not defined, using default port 8080")
		cfg.Port = "8080"
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not defined")
	}
	if cfg.IntegrationSecret == "" {
		log.Fatal("INTEGRATION_SECRET environment variable is not defined")
	}
	if cfg.GoogleClientID == "" {
		log.Println("GOOGLE_CLIENT_ID environment variable is not defined (Google OAuth will not work)")
	}
	if cfg.GoogleClientSecret == "" {
		log.Println("GOOGLE_CLIENT_SECRET environment variable is not defined (Google OAuth will not work)")
	}
	if cfg.GoogleRedirectURL == "" {
		log.Println("GOOGLE_REDIRECT_URL environment variable is not defined (Google OAuth will not work)")
	}
}