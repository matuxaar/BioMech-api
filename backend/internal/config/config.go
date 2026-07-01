package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort        string
	DatabaseURL       string
	DBMaxConns        int
	MLServiceURL      string
	FirebaseCredsFile string
	MigrationsDir     string
	CORSOrigins       string
	DevMode           bool
	MaxUploadSizeMB   int64
	UploadsDir        string
	AvatarsDir        string
	TrainingDir       string
	WSBatchSize       int
	WSBufferInitCap   int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
	MLRequestTimeout  time.Duration
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/desertacia?sslmode=disable"),
		DBMaxConns:        intEnv("DB_MAX_CONNS", 25),
		MLServiceURL:      getEnv("ML_SERVICE_URL", "http://ml:8000"),
		FirebaseCredsFile: getEnv("FIREBASE_CREDENTIALS", "/secrets/firebase-service-account.json"),
		MigrationsDir:     getEnv("MIGRATIONS_DIR", "migrations"),
		CORSOrigins:       getEnv("CORS_ORIGINS", "http://localhost:8080"),
		DevMode:           os.Getenv("DEV_MODE") == "true",
		MaxUploadSizeMB:   int64Env("MAX_UPLOAD_SIZE_MB", 50),
		UploadsDir:        getEnv("UPLOADS_DIR", "uploads"),
		AvatarsDir:        getEnv("AVATARS_DIR", getEnv("UPLOADS_DIR", "uploads")+"/avatars"),
		TrainingDir:       getEnv("TRAINING_DIR", "uploads/training"),
		WSBatchSize:       intEnv("WS_BATCH_SIZE", 32),
		WSBufferInitCap:   intEnv("WS_BUFFER_INIT_CAP", 256),
		ReadTimeout:       durationEnv("READ_TIMEOUT", 15*time.Second),
		WriteTimeout:      durationEnv("WRITE_TIMEOUT", 30*time.Second),
		ReadHeaderTimeout: durationEnv("READ_HEADER_TIMEOUT", 20*time.Second),
		IdleTimeout:       durationEnv("IDLE_TIMEOUT", 60*time.Second),
		MLRequestTimeout:  durationEnv("ML_REQUEST_TIMEOUT", 30*time.Second),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func intEnv(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}

func int64Env(key string, fallback int64) int64 {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.ParseInt(val, 10, 64); err == nil {
			return n
		}
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return fallback
}
