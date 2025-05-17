package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI          string
	Port              string
	JWTSecret         string
	IntegrationSecret string
}

var AppConfig *Config

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
	if cfg.IntegrationSecret == "" {
		log.Fatal("INTEGRATION_SECRET environment variable is not defined")
	}
}
