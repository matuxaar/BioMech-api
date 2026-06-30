package handler

import (
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) SyncUser(c *gin.Context) {
	userID := c.GetString("user_id")

	emailRaw, exists := c.Get("email")
	var emailStr string
	if exists {
		var ok bool
		emailStr, ok = emailRaw.(string)
		if !ok {
			log.Warn().Str("uid", userID).Msg("sync user: email claim is not a string")
		}
	} else {
		log.Warn().Str("uid", userID).Msg("sync user: email claim missing from token")
	}

	user, err := h.authService.SyncUser(c.Request.Context(), userID, emailStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
