package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"stori-gateway/internal/config"

	"github.com/gin-gonic/gin"
)

func Handler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		service := c.Param("service")
		path := c.Param("path")

		config.Lock.RLock()
		defer config.Lock.RUnlock()

		svc, ok := cfg.Services[service]
		if !ok {
			c.JSON(http.StatusBadGateway, gin.H{"error": "unknown service"})
			return
		}

		target, err := url.Parse(svc.URL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid target url"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)

		c.Request.URL.Path = path
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
