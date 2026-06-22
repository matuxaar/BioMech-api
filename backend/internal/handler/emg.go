package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/motvii/desertacia/internal/model"
	"github.com/motvii/desertacia/internal/service"
)

type EMGHandler struct {
	emgService *service.EMGService
}

func NewEMGHandler(emgService *service.EMGService) *EMGHandler {
	return &EMGHandler{emgService: emgService}
}

func (h *EMGHandler) StartSession(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.CreateEMGSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.emgService.StartSession(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *EMGHandler) EndSession(c *gin.Context) {
	if err := h.emgService.EndSession(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session ended"})
}

func (h *EMGHandler) ListSessions(c *gin.Context) {
	userID := c.GetString("user_id")

	sessions, err := h.emgService.ListSessions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *EMGHandler) GetSession(c *gin.Context) {
	session, err := h.emgService.GetSession(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *EMGHandler) AddSample(c *gin.Context) {
	sessionID := c.Param("id")

	var req model.AddSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sample, err := h.emgService.AddSample(c.Request.Context(), sessionID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sample)
}

func (h *EMGHandler) AddSamplesBatch(c *gin.Context) {
	sessionID := c.Param("id")

	var req struct {
		Samples []model.AddSampleRequest `json:"samples" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emgService.AddSamplesBatch(c.Request.Context(), sessionID, req.Samples); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "samples added", "count": len(req.Samples)})
}

func (h *EMGHandler) GetSamples(c *gin.Context) {
	samples, err := h.emgService.GetSamples(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, samples)
}
