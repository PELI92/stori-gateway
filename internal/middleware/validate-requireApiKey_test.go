package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"stori-gateway/internal/config"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"stori-gateway/internal/middleware"
)

// mock config
type mockConfig struct{}

func (m *mockConfig) GetAPIKey() string {
	return "supersecretkey"
}

func (m *mockConfig) GetServiceURL(string) (string, bool) {
	return "", false
}

func (m *mockConfig) GetAllServices() map[string]config.Service {
	return nil
}

func TestRequireAPIKey_ValidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequireAPIKey(&mockConfig{}))
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-api-key", "supersecretkey")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestRequireAPIKey_MissingKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequireAPIKey(&mockConfig{}))
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "should not reach"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "x-api-key header is required")
}

func TestRequireAPIKey_InvalidKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RequireAPIKey(&mockConfig{}))
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "should not reach"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("x-api-key", "wrongkey")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid API key")
}
