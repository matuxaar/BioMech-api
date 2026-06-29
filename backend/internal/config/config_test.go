package config

import (
	"os"
	"testing"
	"time"
)

func TestGetEnv(t *testing.T) {
	t.Run("uses value from env", func(t *testing.T) {
		t.Setenv("TEST_KEY", "actual")
		if got := getEnv("TEST_KEY", "fallback"); got != "actual" {
			t.Errorf("expected actual, got %s", got)
		}
	})

	t.Run("uses fallback when not set", func(t *testing.T) {
		os.Unsetenv("TEST_KEY_MISSING")
		if got := getEnv("TEST_KEY_MISSING", "fallback"); got != "fallback" {
			t.Errorf("expected fallback, got %s", got)
		}
	})

	t.Run("uses fallback when empty", func(t *testing.T) {
		t.Setenv("TEST_EMPTY", "")
		if got := getEnv("TEST_EMPTY", "fallback"); got != "fallback" {
			t.Errorf("expected fallback, got %s", got)
		}
	})
}

func TestIntEnv(t *testing.T) {
	t.Run("valid integer", func(t *testing.T) {
		t.Setenv("TEST_INT", "42")
		if got := intEnv("TEST_INT", 0); got != 42 {
			t.Errorf("expected 42, got %d", got)
		}
	})

	t.Run("invalid integer uses fallback", func(t *testing.T) {
		t.Setenv("TEST_INT_BAD", "not-a-number")
		if got := intEnv("TEST_INT_BAD", 10); got != 10 {
			t.Errorf("expected 10, got %d", got)
		}
	})

	t.Run("empty uses fallback", func(t *testing.T) {
		t.Setenv("TEST_INT_EMPTY", "")
		if got := intEnv("TEST_INT_EMPTY", 10); got != 10 {
			t.Errorf("expected 10, got %d", got)
		}
	})
}

func TestInt64Env(t *testing.T) {
	t.Run("valid int64", func(t *testing.T) {
		t.Setenv("TEST_INT64", "9223372036854775807")
		if got := int64Env("TEST_INT64", 0); got != 9223372036854775807 {
			t.Errorf("expected max int64, got %d", got)
		}
	})

	t.Run("invalid uses fallback", func(t *testing.T) {
		t.Setenv("TEST_INT64_BAD", "nan")
		if got := int64Env("TEST_INT64_BAD", 99); got != 99 {
			t.Errorf("expected 99, got %d", got)
		}
	})
}

func TestDurationEnv(t *testing.T) {
	t.Run("valid duration", func(t *testing.T) {
		t.Setenv("TEST_DUR", "30s")
		if got := durationEnv("TEST_DUR", 0); got != 30*time.Second {
			t.Errorf("expected 30s, got %v", got)
		}
	})

	t.Run("valid duration with minutes", func(t *testing.T) {
		t.Setenv("TEST_DUR_MIN", "5m")
		if got := durationEnv("TEST_DUR_MIN", 0); got != 5*time.Minute {
			t.Errorf("expected 5m, got %v", got)
		}
	})

	t.Run("invalid uses fallback", func(t *testing.T) {
		t.Setenv("TEST_DUR_BAD", "forever")
		if got := durationEnv("TEST_DUR_BAD", 10*time.Second); got != 10*time.Second {
			t.Errorf("expected 10s, got %v", got)
		}
	})
}

func TestLoadDefaults(t *testing.T) {
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("ML_SERVICE_URL")
	os.Unsetenv("DEV_MODE")
	os.Unsetenv("MAX_UPLOAD_SIZE_MB")
	os.Unsetenv("UPLOADS_DIR")
	os.Unsetenv("WS_BATCH_SIZE")
	os.Unsetenv("WS_BUFFER_INIT_CAP")

	cfg := Load()

	if cfg.ServerPort != "8080" {
		t.Errorf("expected 8080, got %s", cfg.ServerPort)
	}
	if cfg.MaxUploadSizeMB != 50 {
		t.Errorf("expected 50, got %d", cfg.MaxUploadSizeMB)
	}
	if cfg.DevMode {
		t.Error("expected DevMode false")
	}
	if cfg.WSBatchSize != 32 {
		t.Errorf("expected 32, got %d", cfg.WSBatchSize)
	}
	if cfg.ReadTimeout != 15*time.Second {
		t.Errorf("expected 15s, got %v", cfg.ReadTimeout)
	}
	if cfg.MLRequestTimeout != 30*time.Second {
		t.Errorf("expected 30s, got %v", cfg.MLRequestTimeout)
	}
}

func TestLoadWithEnvOverrides(t *testing.T) {
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://user:pass@remote/db")
	t.Setenv("ML_SERVICE_URL", "http://ml-prod:8000")
	t.Setenv("DEV_MODE", "true")
	t.Setenv("MAX_UPLOAD_SIZE_MB", "100")
	t.Setenv("UPLOADS_DIR", "/data/uploads")
	t.Setenv("WS_BATCH_SIZE", "64")
	t.Setenv("WS_BUFFER_INIT_CAP", "512")
	t.Setenv("READ_TIMEOUT", "45s")
	t.Setenv("ML_REQUEST_TIMEOUT", "2m")

	cfg := Load()

	if cfg.ServerPort != "9090" {
		t.Errorf("expected 9090, got %s", cfg.ServerPort)
	}
	if cfg.DatabaseURL != "postgres://user:pass@remote/db" {
		t.Errorf("expected remote db, got %s", cfg.DatabaseURL)
	}
	if cfg.MLServiceURL != "http://ml-prod:8000" {
		t.Errorf("expected ml-prod, got %s", cfg.MLServiceURL)
	}
	if !cfg.DevMode {
		t.Error("expected DevMode true")
	}
	if cfg.MaxUploadSizeMB != 100 {
		t.Errorf("expected 100, got %d", cfg.MaxUploadSizeMB)
	}
	if cfg.UploadsDir != "/data/uploads" {
		t.Errorf("expected /data/uploads, got %s", cfg.UploadsDir)
	}
	if cfg.WSBatchSize != 64 {
		t.Errorf("expected 64, got %d", cfg.WSBatchSize)
	}
	if cfg.WSBufferInitCap != 512 {
		t.Errorf("expected 512, got %d", cfg.WSBufferInitCap)
	}
	if cfg.ReadTimeout != 45*time.Second {
		t.Errorf("expected 45s, got %v", cfg.ReadTimeout)
	}
	if cfg.MLRequestTimeout != 2*time.Minute {
		t.Errorf("expected 2m, got %v", cfg.MLRequestTimeout)
	}
}

func TestLoadInvalidDuration(t *testing.T) {
	t.Setenv("READ_TIMEOUT", "invalid")
	cfg := Load()
	if cfg.ReadTimeout != 15*time.Second {
		t.Errorf("expected fallback 15s on invalid, got %v", cfg.ReadTimeout)
	}
}

func TestLoadInvalidInt(t *testing.T) {
	t.Setenv("WS_BATCH_SIZE", "not-a-number")
	cfg := Load()
	if cfg.WSBatchSize != 32 {
		t.Errorf("expected fallback 32 on invalid, got %d", cfg.WSBatchSize)
	}
}
