package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/service"
)

type DeviceHandler struct {
	deviceService DeviceService
}

func NewDeviceHandler(deviceService DeviceService) *DeviceHandler {
	return &DeviceHandler{deviceService: deviceService}
}

func (h *DeviceHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := h.deviceService.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, device)
}

func (h *DeviceHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	p := model.ParsePagination(c)

	result, err := h.deviceService.ListByUser(c.Request.Context(), userID, p.Page, p.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *DeviceHandler) GetByID(c *gin.Context) {
	userID := c.GetString("user_id")
	device, err := h.deviceService.GetByID(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}

func (h *DeviceHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := h.deviceService.Update(c.Request.Context(), userID, c.Param("id"), &req)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrDeviceNotFound):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrAccessDenied):
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device)
}

func (h *DeviceHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := h.deviceService.Delete(c.Request.Context(), userID, c.Param("id")); err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrDeviceNotFound):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrAccessDenied):
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "device deleted"})
}

func (h *DeviceHandler) GetActions(c *gin.Context) {
	userID := c.GetString("user_id")

	actions, err := h.deviceService.GetActions(c.Request.Context(), userID, c.Param("id"))
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, service.ErrDeviceNotFound):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrAccessDenied):
			status = http.StatusForbidden
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, actions)
}
