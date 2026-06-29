package integration

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/model"
)

func TestParsePagination(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		wantPage  int
		wantLimit int
	}{
		{"defaults", "", 1, 20},
		{"custom page", "?page=3", 3, 20},
		{"custom limit", "?limit=50", 1, 50},
		{"both", "?page=2&limit=10", 2, 10},
		{"page too small", "?page=0", 1, 20},
		{"limit too large", "?limit=200", 1, 20},
		{"limit too small", "?limit=0", 1, 20},
		{"invalid page", "?page=abc", 1, 20},
		{"max limit", "?limit=100", 1, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/test"+tt.query, nil)

			p := model.ParsePagination(c)
			if p.Page != tt.wantPage {
				t.Errorf("page: got %d, want %d", p.Page, tt.wantPage)
			}
			if p.Limit != tt.wantLimit {
				t.Errorf("limit: got %d, want %d", p.Limit, tt.wantLimit)
			}
		})
	}
}

func TestNewPaginatedResponse(t *testing.T) {
	data := []string{"a", "b", "c"}

	t.Run("single page", func(t *testing.T) {
		resp := model.NewPaginatedResponse(data, 3, 1, 20)
		if resp.Total != 3 {
			t.Errorf("total: got %d, want 3", resp.Total)
		}
		if resp.TotalPages != 1 {
			t.Errorf("total_pages: got %d, want 1", resp.TotalPages)
		}
		if len(resp.Data) != 3 {
			t.Errorf("data len: got %d, want 3", len(resp.Data))
		}
	})

	t.Run("multiple pages", func(t *testing.T) {
		resp := model.NewPaginatedResponse(data, 50, 2, 20)
		if resp.Total != 50 {
			t.Errorf("total: got %d, want 50", resp.Total)
		}
		if resp.TotalPages != 3 {
			t.Errorf("total_pages: got %d, want 3", resp.TotalPages)
		}
	})

	t.Run("exact fit", func(t *testing.T) {
		resp := model.NewPaginatedResponse(data, 40, 1, 20)
		if resp.TotalPages != 2 {
			t.Errorf("total_pages: got %d, want 2", resp.TotalPages)
		}
	})

	t.Run("empty data", func(t *testing.T) {
		resp := model.NewPaginatedResponse([]string{}, 0, 1, 20)
		if resp.Total != 0 {
			t.Errorf("total: got %d, want 0", resp.Total)
		}
		if resp.TotalPages != 0 {
			t.Errorf("total_pages: got %d, want 0", resp.TotalPages)
		}
	})
}
