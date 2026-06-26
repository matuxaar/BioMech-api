package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/model"
	"github.com/matuxaar/BioMech-api/internal/service"
)

type UserHandler struct {
	authService *service.AuthService
}

func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")

	profile, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *UserHandler) Update(c *gin.Context) {
	userID := c.GetString("user_id")

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.UpdateProfile(c.Request.Context(), userID, &req)
	if err != nil {
		if err.Error() == "nickname already taken" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
