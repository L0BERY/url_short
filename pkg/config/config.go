package config

import "os"

type Config struct {
	ServerAddress string
	BaseURL       string
	DatabaseURL   string
}

func LoadConfig() *Config {
	return &Config{
		ServerAddress: getEnv("SERVER_ADDRES", ":8080"),
		BaseURL:       getEnv("BASE_URL", "http://localhost:8080"),
		DatabaseURL:   getEnv("DTABASE_URL", "postgres://user:password@localhost:5432/url_shortener?sslmode=disable"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
