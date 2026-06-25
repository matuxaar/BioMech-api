package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort        string
	DatabaseURL       string
	MLServiceURL      string
	FirebaseCredsFile string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/desertacia?sslmode=disable"),
		MLServiceURL:      getEnv("ML_SERVICE_URL", "http://ml:8000"),
		FirebaseCredsFile: getEnv("FIREBASE_CREDENTIALS", "/secrets/firebase-service-account.json"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
