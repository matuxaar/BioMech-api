package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/matuxaar/BioMech-api/internal/middleware"
)

func TestRequestIDMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("generates ID when not provided", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.RequestID())
		r.GET("/test", func(c *gin.Context) {
			id := middleware.GetRequestID(c)
			if id == "" {
				t.Error("expected non-empty request_id")
			}
			_, err := uuid.Parse(id)
			if err != nil {
				t.Errorf("expected valid UUID, got %s: %v", id, err)
			}
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		respID := w.Header().Get("X-Request-ID")
		if respID == "" {
			t.Error("expected X-Request-ID response header")
		}
	})

	t.Run("preserves provided ID", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.RequestID())
		r.GET("/test", func(c *gin.Context) {
			id := middleware.GetRequestID(c)
			if id != "client-provided-id" {
				t.Errorf("expected client-provided-id, got %s", id)
			}
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Request-ID", "client-provided-id")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
		if w.Header().Get("X-Request-ID") != "client-provided-id" {
			t.Errorf("expected X-Request-ID to be echoed back")
		}
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.RateLimit(100, 200))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("allows requests within limit", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Errorf("request %d: expected 200, got %d", i, w.Code)
			}
		}
	})

	t.Run("different IP is not rate limited", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "10.0.0.1:12345"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200 for different IP, got %d", w.Code)
		}
	})
}

func TestRateLimitExceedsBurst(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.RateLimit(10, 5))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	burst := 0
	for i := 0; i < 50; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code == http.StatusTooManyRequests {
			burst++
		}
	}
	if burst == 0 {
		t.Error("expected at least one 429 response when exceeding burst")
	}
}
