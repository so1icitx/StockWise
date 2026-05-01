package config

import "os"

// Config contains runtime settings for the StockWise API.
type Config struct {
	AppName     string
	AppEnv      string
	ServerHost  string
	ServerPort  string
	DatabaseURL string
}

// Load reads application configuration from environment variables.
func Load() Config {
	return Config{
		AppName:     getEnv("APP_NAME", "StockWise"),
		AppEnv:      getEnv("APP_ENV", "development"),
		ServerHost:  getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/stockwise?sslmode=disable"),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
