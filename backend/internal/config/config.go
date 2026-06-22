package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	DatabaseURL   string
	JWTSecret     string
	JWTExpiresIn  time.Duration
	RefreshExpiry time.Duration
	MLServiceURL  string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/desertacia?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		JWTExpiresIn:  15 * time.Minute,
		RefreshExpiry: 7 * 24 * time.Hour,
		MLServiceURL:  getEnv("ML_SERVICE_URL", "http://ml:8000"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
