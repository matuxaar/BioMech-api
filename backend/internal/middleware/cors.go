package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/matuxaar/BioMech-api/internal/config"
)

func CORS() gin.HandlerFunc {
	origins := config.Load().CORSOrigins
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origins)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
