package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/model"
)

type EMGHandler struct {
	emgService EMGService
}

func NewEMGHandler(emgService EMGService) *EMGHandler {
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
	userID := c.GetString("user_id")
	if err := h.emgService.EndSession(c.Request.Context(), userID, c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session ended"})
}

func (h *EMGHandler) ListSessions(c *gin.Context) {
	userID := c.GetString("user_id")
	p := model.ParsePagination(c)

	result, err := h.emgService.ListSessions(c.Request.Context(), userID, p.Page, p.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *EMGHandler) GetSession(c *gin.Context) {
	userID := c.GetString("user_id")
	session, err := h.emgService.GetSession(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (h *EMGHandler) AddSample(c *gin.Context) {
	userID := c.GetString("user_id")
	sessionID := c.Param("id")

	var req model.AddSampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sample, err := h.emgService.AddSample(c.Request.Context(), userID, sessionID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sample)
}

func (h *EMGHandler) AddSamplesBatch(c *gin.Context) {
	userID := c.GetString("user_id")
	sessionID := c.Param("id")

	var req struct {
		Samples []model.AddSampleRequest `json:"samples" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emgService.AddSamplesBatch(c.Request.Context(), userID, sessionID, req.Samples); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "samples added", "count": len(req.Samples)})
}

func (h *EMGHandler) GetSamples(c *gin.Context) {
	userID := c.GetString("user_id")
	p := model.ParsePagination(c)

	result, err := h.emgService.GetSamples(c.Request.Context(), userID, c.Param("id"), p.Page, p.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
