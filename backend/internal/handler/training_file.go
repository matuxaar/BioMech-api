package handler

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/service"
)

type TrainingFileHandler struct {
	fileService *service.TrainingFileService
}

func NewTrainingFileHandler(fileService *service.TrainingFileService) *TrainingFileHandler {
	return &TrainingFileHandler{fileService: fileService}
}

func (h *TrainingFileHandler) Upload(c *gin.Context) {
	userID := c.GetString("user_id")

	deviceID := c.PostForm("device_id")
	label := c.PostForm("label")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	tf, err := h.fileService.Upload(c.Request.Context(), userID, deviceID, label, header.Filename, file, header.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      tf.ID,
		"message": "file uploaded successfully",
	})
}

func (h *TrainingFileHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	files, err := h.fileService.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *TrainingFileHandler) Get(c *gin.Context) {
	file, err := h.fileService.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.JSON(http.StatusOK, file)
}

func (h *TrainingFileHandler) Download(c *gin.Context) {
	filePath, err := h.fileService.GetFilePath(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.FileAttachment(filePath, filepath.Base(filePath))
}

func (h *TrainingFileHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.fileService.Delete(c.Request.Context(), c.Param("id"), userID); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file deleted"})
}
