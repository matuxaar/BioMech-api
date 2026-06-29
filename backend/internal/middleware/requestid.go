package middleware

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		c.Set("request_id", reqID)
		c.Header("X-Request-ID", reqID)
		slog.Debug("request", "method", c.Request.Method, "path", c.Request.URL.Path, "request_id", reqID)
		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get("request_id"); exists {
		return id.(string)
	}
	return ""
}
