package proxy

import (
	"bytes"
	"github.com/rs/zerolog/log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"stori-gateway/internal/config"

	"github.com/gin-gonic/gin"
)

const maxLogBodySize = 1000

type ReverseProxy struct {
	cfg config.Provider
}

func NewReverseProxy(cfg config.Provider) *ReverseProxy {
	return &ReverseProxy{cfg: cfg}
}

func (p *ReverseProxy) Handle(c *gin.Context) {
	service := c.Param("service")
	path := c.Param("path")

	serviceURL, ok := p.cfg.GetServiceURL(service)
	if !ok {
		c.JSON(http.StatusBadGateway, gin.H{"error": "unknown service"})
		return
	}

	target, err := url.Parse(serviceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid target url"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// This is not explicitly requested in the assigment, but there is no point in propagating the api key downstream once the request is validated. All other headers will be propagated.
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Del("x-api-key")
	}

	c.Request.URL.Path = path

	rec := &ResponseRecorder{
		ResponseWriter: c.Writer,
		Body:           &bytes.Buffer{},
		StatusCode:     200, // default if WriteHeader never called
	}

	proxy.ServeHTTP(rec, c.Request)

	go func() {
		responseBody := rec.Body.String()
		if len(responseBody) > maxLogBodySize {
			responseBody = responseBody[:maxLogBodySize] + "..."
		}

		log.Info().
			Int("downstream_status", rec.StatusCode).
			Str("response_body", rec.Body.String()).
			Msg("downstream response")
	}()
}
