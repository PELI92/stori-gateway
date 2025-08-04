package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ZeroLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		var bodyCopy []byte
		if c.Request.Body != nil {
			bodyCopy, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))
		}

		c.Next()

		service := c.Param("service")
		method := c.Request.Method
		path := c.Request.URL.Path
		status := c.Writer.Status()
		duration := time.Since(start)
		userAgent := c.Request.UserAgent()
		requestID := c.GetHeader("X-Request-ID")
		errors := c.Errors

		go func() {
			event := log.Info()
			if status >= 400 {
				event = log.Error()
			} else if log.Logger.GetLevel() == zerolog.DebugLevel {
				event = log.Debug()
			}

			event = event.
				Str("service", service).
				Str("method", method).
				Str("path", path).
				Int("status", status).
				Dur("duration", duration).
				Str("user-agent", userAgent).
				Str("request-id", requestID)

			if log.Logger.GetLevel() == zerolog.DebugLevel && len(bodyCopy) > 0 {
				var tmp interface{}
				if json.Unmarshal(bodyCopy, &tmp) == nil {
					event = event.RawJSON("body", bodyCopy)
				} else {
					event = event.Str("raw_body", string(bodyCopy))
				}
			}

			event.Msg("handled request")

			for _, err := range errors {
				log.Error().Err(err).Msg("handler error")
			}
		}()
	}
}
