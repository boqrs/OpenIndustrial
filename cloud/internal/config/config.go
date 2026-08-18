package config

import (
	"github.com/OpenIndustrial/cloud/internal/kernel/security/provider"

	"os"
)

// Config holds all configuration for the application.
type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	RedisAddr   string
	PKI provider.ProviderConfig `mapstructure:"provider"`
}

// Load reads configuration from environment variables or uses sensible defaults.
func Load() (*Config, error) {
	// For development, we can use defaults if environment variables aren't set.
	// In production, you would set these environment variables.
	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/open_industrial?sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", "a-very-secret-and-strong-key-that-should-be-changed"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
	}
	return cfg, nil
}

// getEnv is a helper function to read an environment variable or return a fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}