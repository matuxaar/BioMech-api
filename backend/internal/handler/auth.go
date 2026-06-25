package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SyncUser(c *gin.Context) {
	userID := c.GetString("user_id")
	email, _ := c.Get("email")

	emailStr, _ := email.(string)
	user, err := h.authService.SyncUser(c.Request.Context(), userID, emailStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
