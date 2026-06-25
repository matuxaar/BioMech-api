package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/motvii/desertacia/internal/service"
)

type UserHandler struct {
	authService *service.AuthService
}

func NewUserHandler(authService *service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

func (h *UserHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
