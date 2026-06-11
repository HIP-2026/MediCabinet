package config

import "os"

// Config holds runtime configuration loaded from the environment.
type Config struct {
	Port        string
	DatabaseURL string
}

// Load reads configuration from the environment, applying defaults suitable
// for local development.
func Load() Config {
	return Config{
		Port:        getenv("PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://medicabinet:medicabinet@localhost:5432/medicabinet?sslmode=disable"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
