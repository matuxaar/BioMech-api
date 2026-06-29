package model

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewPaginatedResponse(t *testing.T) {
	t.Run("single page exact fit", func(t *testing.T) {
		data := []int{1, 2, 3}
		res := NewPaginatedResponse(data, 3, 1, 3)
		if res.TotalPages != 1 {
			t.Errorf("expected 1 page, got %d", res.TotalPages)
		}
		if res.Total != 3 {
			t.Errorf("expected total 3, got %d", res.Total)
		}
	})

	t.Run("multiple pages", func(t *testing.T) {
		data := []int{1, 2, 3}
		res := NewPaginatedResponse(data, 10, 1, 3)
		if res.TotalPages != 4 {
			t.Errorf("expected 4 pages, got %d", res.TotalPages)
		}
	})

	t.Run("empty data", func(t *testing.T) {
		data := []int{}
		res := NewPaginatedResponse(data, 0, 1, 20)
		if res.TotalPages != 0 {
			t.Errorf("expected 0 pages, got %d", res.TotalPages)
		}
	})

	t.Run("single item single page", func(t *testing.T) {
		data := []string{"a"}
		res := NewPaginatedResponse(data, 1, 1, 20)
		if res.TotalPages != 1 {
			t.Errorf("expected 1 page, got %d", res.TotalPages)
		}
	})
}

func TestParsePagination(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		c := newGinCtx(req)
		p := ParsePagination(c)
		if p.Page != 1 || p.Limit != 20 {
			t.Errorf("expected (1,20), got (%d,%d)", p.Page, p.Limit)
		}
	})

	t.Run("custom page and limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?page=3&limit=10", nil)
		c := newGinCtx(req)
		p := ParsePagination(c)
		if p.Page != 3 || p.Limit != 10 {
			t.Errorf("expected (3,10), got (%d,%d)", p.Page, p.Limit)
		}
	})

	t.Run("page too small", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?page=0&limit=20", nil)
		c := newGinCtx(req)
		p := ParsePagination(c)
		if p.Page != 1 {
			t.Errorf("expected page 1, got %d", p.Page)
		}
	})

	t.Run("limit too large", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?page=1&limit=200", nil)
		c := newGinCtx(req)
		p := ParsePagination(c)
		if p.Limit != 20 {
			t.Errorf("expected limit 20, got %d", p.Limit)
		}
	})

	t.Run("limit too small", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?page=1&limit=0", nil)
		c := newGinCtx(req)
		p := ParsePagination(c)
		if p.Limit != 20 {
			t.Errorf("expected limit 20, got %d", p.Limit)
		}
	})

	t.Run("invalid page", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?page=abc&limit=20", nil)
		c := newGinCtx(req)
		p := ParsePagination(c)
		if p.Page != 1 {
			t.Errorf("expected page 1, got %d", p.Page)
		}
	})

	t.Run("max limit", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/?page=1&limit=100", nil)
		c := newGinCtx(req)
		p := ParsePagination(c)
		if p.Limit != 100 {
			t.Errorf("expected limit 100, got %d", p.Limit)
		}
	})
}

func newGinCtx(req *http.Request) *gin.Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c
}
