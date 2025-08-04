package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateJSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			if c.GetHeader("Content-Type") == "application/json" {
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
					c.Abort()
					return
				}

				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				
				var tmp interface{}
				if err := json.Unmarshal(bodyBytes, &tmp); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
