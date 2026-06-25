package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/service"
)

type TrainingHandler struct {
	trainingService *service.TrainingService
}

func NewTrainingHandler(trainingService *service.TrainingService) *TrainingHandler {
	return &TrainingHandler{trainingService: trainingService}
}

func (h *TrainingHandler) CreateJob(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.CreateTrainingJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.trainingService.CreateJob(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.trainingService.StartTraining(c.Request.Context(), job.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, job)
}

func (h *TrainingHandler) ListJobs(c *gin.Context) {
	userID := c.GetString("user_id")

	jobs, err := h.trainingService.ListJobs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func (h *TrainingHandler) GetJob(c *gin.Context) {
	job, err := h.trainingService.GetJob(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}
