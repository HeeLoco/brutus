package configs

import (
	"os"
)

type Config struct {
	Port string
}

func Load() *Config {
	return &Config{
		Port: getEnvOrDefault("PORT", "8080"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}