//go:build !e2e

package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/middleware"
)

func TestDevModeAuthBypass(t *testing.T) {
	middleware.InitAuth(true)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.AuthRequired(nil))
	r.GET("/test", func(c *gin.Context) {
		uid := c.GetString("user_id")
		if uid != "dev-user-id" {
			t.Errorf("expected dev-user-id, got %s", uid)
		}
		email := c.GetString("email")
		if email != "dev@biomech.app" {
			t.Errorf("expected dev@biomech.app, got %s", email)
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestDevModeDisabledRequiresAuth(t *testing.T) {
	middleware.InitAuth(false)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.AuthRequired(nil))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestInternalAuthMiddleware(t *testing.T) {
	t.Setenv("INTERNAL_API_KEY", "test-key-123")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/callback", middleware.InternalAuth(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	t.Run("valid key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/callback", nil)
		req.Header.Set("X-API-Key", "test-key-123")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("invalid key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/callback", nil)
		req.Header.Set("X-API-Key", "wrong-key")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	t.Run("missing key", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/callback", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})
}

func TestInternalAuthNotConfigured(t *testing.T) {
	t.Setenv("INTERNAL_API_KEY", "")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/callback", middleware.InternalAuth(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/callback", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}
