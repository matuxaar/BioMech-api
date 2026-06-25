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

type uploadRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Label    string `json:"label"`
}

type predictRequest struct {
	Samples []model.AddSampleRequest `json:"samples" binding:"required"`
}

func (h *TrainingHandler) Predict(c *gin.Context) {
	var req predictRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	samples := make([]model.EMGSample, len(req.Samples))
	for i, s := range req.Samples {
		samples[i] = model.EMGSample{
			Timestamp: s.Timestamp,
			Channel1:  s.Channel1,
			Channel2:  s.Channel2,
			Channel3:  s.Channel3,
			Channel4:  s.Channel4,
			Channel5:  s.Channel5,
			Channel6:  s.Channel6,
			Channel7:  s.Channel7,
			Channel8:  s.Channel8,
			Metadata:  s.Metadata,
		}
	}

	result, err := h.trainingService.Predict(c.Request.Context(), samples)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *TrainingHandler) Upload(c *gin.Context) {
	userID := c.GetString("user_id")

	deviceID := c.PostForm("device_id")
	if deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id is required"})
		return
	}
	label := c.PostForm("label")

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	session, err := h.trainingService.ProcessUpload(c.Request.Context(), userID, deviceID, label, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"session_id": session.ID,
		"message":    "file uploaded and session created",
	})
}
