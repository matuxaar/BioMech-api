package middleware

import (
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

var allowedOrigins []string

func InitCORS(origins string) {
	if origins == "" || origins == "*" {
		allowedOrigins = []string{"*"}
	} else {
		for _, o := range strings.Split(origins, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				allowedOrigins = append(allowedOrigins, o)
			}
		}
		if len(allowedOrigins) == 0 {
			allowedOrigins = []string{"*"}
		}
	}
}

func CORS() gin.HandlerFunc {
	allWildcard := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" && !allWildcard && !slices.Contains(allowedOrigins, origin) {
			c.AbortWithStatus(403)
			return
		}
		if allWildcard {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if slices.Contains(allowedOrigins, origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		}
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
