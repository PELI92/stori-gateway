package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-Request-ID header is required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
