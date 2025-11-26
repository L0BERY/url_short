package config

import (
	"os"

	"shorturl.com/pkg/logger"
)

type Config struct {
	ServerAddress string
	BaseURL       string
	DatabaseURL   string
	LogLevel      logger.Level
}

func LoadConfig() *Config {
	return &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		BaseURL:       getEnv("BASE_URL", "http://localhost:8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/url_shortener?sslmode=disable"),
		LogLevel:      logger.LevelInfo,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
