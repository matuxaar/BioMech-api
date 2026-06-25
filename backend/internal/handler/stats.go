package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/service"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: statsService}
}

func (h *StatsHandler) Dashboard(c *gin.Context) {
	userID := c.GetString("user_id")

	stats, err := h.statsService.GetDashboardStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
