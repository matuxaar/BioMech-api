package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestGetRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("returns empty when not set", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		if id := GetRequestID(c); id != "" {
			t.Errorf("expected empty, got %s", id)
		}
	})

	t.Run("returns value when set", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("request_id", "test-id")
		if id := GetRequestID(c); id != "test-id" {
			t.Errorf("expected test-id, got %s", id)
		}
	})
}

func TestCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	InitCORS("*")

	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("sets CORS headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodOptions, "/test", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "GET")
		r.ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") == "" {
			t.Error("expected CORS headers")
		}
		if w.Code != http.StatusNoContent {
			t.Errorf("expected 204 for preflight, got %d", w.Code)
		}
	})

	t.Run("allows actual request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "http://example.com")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}

func TestCORSNoWildcard(t *testing.T) {
	gin.SetMode(gin.TestMode)
	InitCORS("http://app.example.com")

	r := gin.New()
	r.Use(CORS())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("allows configured origin", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "http://app.example.com")
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if w.Header().Get("Access-Control-Allow-Origin") != "http://app.example.com" {
			t.Errorf("expected specific origin")
		}
	})

	t.Run("blocks unknown origin", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Origin", "http://evil.com")
		r.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Errorf("expected 403, got %d", w.Code)
		}
	})
}

func TestAuthRequiredDevMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	InitAuth(true)

	r := gin.New()
	r.Use(AuthRequired(nil))
	r.GET("/test", func(c *gin.Context) {
		uid := c.GetString("user_id")
		if uid != "dev-user-id" {
			t.Errorf("expected dev-user-id, got %s", uid)
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthRequiredNoDevMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	InitAuth(false)

	r := gin.New()
	r.Use(AuthRequired(nil))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestRateLimitCleanup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := RateLimit(100, 10)
	r := gin.New()
	r.Use(handler)
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Send a request to initialize
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	// Can't directly test cleanup goroutine, but verify basic function
	StopRateLimiters()
}

func TestMetricsRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(Metrics())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestIncRateLimited(t *testing.T) {
	IncRateLimited()
}

func TestTrackWSOpenClose(t *testing.T) {
	TrackWSOpen()
	TrackWSClose()
}

func TestTrackWSMessage(t *testing.T) {
	TrackWSMessage("incoming")
	TrackWSMessage("outgoing")
}

func TestRequestIDSetsUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		id := GetRequestID(c)
		if _, err := uuid.Parse(id); err != nil {
			t.Errorf("expected valid UUID, got %s: %v", id, err)
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	respID := w.Header().Get("X-Request-ID")
	if respID == "" {
		t.Error("expected X-Request-ID header")
	}
}

func TestRequestIDPreservesClientID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		id := GetRequestID(c)
		if id != "client-id" {
			t.Errorf("expected client-id, got %s", id)
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "client-id")
	r.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") != "client-id" {
		t.Errorf("expected X-Request-ID to be echoed")
	}
}

func TestInternalAuth(t *testing.T) {
	t.Setenv("INTERNAL_API_KEY", "test-key")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/callback", InternalAuth(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("valid key", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/callback", nil)
		req.Header.Set("X-API-Key", "test-key")
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("invalid key", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/callback", nil)
		req.Header.Set("X-API-Key", "wrong")
		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})
}
