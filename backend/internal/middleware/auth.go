package middleware

import (
	"context"
	"net/http"
	"strings"

		firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
)

func AuthRequired(firebaseApp *firebase.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		client, err := firebaseApp.Auth(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "auth service unavailable"})
			return
		}

		token, err := client.VerifyIDToken(context.Background(), parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user_id", token.UID)
		if email, ok := token.Claims["email"]; ok {
			c.Set("email", email)
		}
		c.Next()
	}
}
