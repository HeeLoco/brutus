package configs

import (
	"os"
)

type Config struct {
	Port           string
	SubscriptionID string
	TenantID       string
	ClientID       string
	ClientSecret   string
}

func Load() *Config {
	return &Config{
		Port:           getEnvOrDefault("PORT", "8080"),
		SubscriptionID: getEnvOrDefault("AZURE_SUBSCRIPTION_ID", ""),
		TenantID:       getEnvOrDefault("AZURE_TENANT_ID", ""),
		ClientID:       getEnvOrDefault("AZURE_CLIENT_ID", ""),
		ClientSecret:   getEnvOrDefault("AZURE_CLIENT_SECRET", ""),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}