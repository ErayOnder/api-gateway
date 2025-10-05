package config

import (
	"os"
)

// Config holds all configuration for the API Gateway
type Config struct {
	ServerPort     string
	LLMServiceURL  string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		LLMServiceURL: getEnv("LLM_SERVICE_URL", "http://localhost:3000"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
