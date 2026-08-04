package config

import (
	"os"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	GitHubClientID string
	GitHubSecret   string
	FrontendURL    string
}

func Load() *Config {
	
	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        log.Fatal("DATABASE_URL environment variable is required but empty")
    }
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		GitHubClientID: getEnv("GITHUB_CLIENT_ID", ""),
		GitHubSecret:   getEnv("GITHUB_CLIENT_SECRET", ""),
		FrontendURL:    getEnv("FRONTEND_URL", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}