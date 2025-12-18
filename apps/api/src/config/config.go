package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AtlassianConfig struct {
	Audience     string
	ClientId     string
	ClientSecret string
	Scope        string
	CallbackUrl  string
}

type Config struct {
	APIVersion      string
	MongoURI        string
	Port            string
	JWTSecret       string
	AtlassianConfig AtlassianConfig
}

var AppConfig *Config

// Version can be set via build flags (-ldflags "-X 'TimeTrack-api/src/config.Version=...'")
var Version = "dev"

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Could not load .env file, proceeding with environment variables")
	}

	cfg := &Config{
		APIVersion: Version,
		MongoURI:   os.Getenv("MONGO_URI"),
		Port:       os.Getenv("PORT"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		AtlassianConfig: AtlassianConfig{
			Audience:     os.Getenv("ATLASSIAN_AUDIENCE"),
			ClientId:     os.Getenv("ATLASSIAN_CLIENT_ID"),
			ClientSecret: os.Getenv("ATLASSIAN_CLIENT_SECRET"),
			Scope:        os.Getenv("ATLASSIAN_SCOPE"),
			CallbackUrl:  os.Getenv("ATLASSIAN_CALLBACK_URL"),
		},
	}
	AppConfig = cfg
	return cfg
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
	if cfg.AtlassianConfig.Audience == "" {
		log.Fatal("ATLASSIAN_AUDIENCE environment variable is not defined")
	}
	if cfg.AtlassianConfig.ClientId == "" {
		log.Fatal("ATLASSIAN_CLIENT_ID environment variable is not defined")
	}
	if cfg.AtlassianConfig.ClientSecret == "" {
		log.Fatal("ATLASSIAN_CLIENT_SECRET environment variable is not defined")
	}
	if cfg.AtlassianConfig.Scope == "" {
		log.Fatal("ATLASSIAN_SCOPE environment variable is not defined")
	}
	if cfg.AtlassianConfig.CallbackUrl == "" {
		log.Fatal("ATLASSIAN_CALLBACK_URL environment variable is not defined")
	}
}
