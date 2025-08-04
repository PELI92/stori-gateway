package middleware

import (
	"net/http"
	"stori-gateway/internal/config"

	"github.com/gin-gonic/gin"
)

func RequireAPIKey(cfg config.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "x-api-key header is required"})
			c.Abort()
			return
		}

		if apiKey != cfg.GetAPIKey() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
