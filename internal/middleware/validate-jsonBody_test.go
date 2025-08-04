package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"stori-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouterWithValidateJSON() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ValidateJSONMiddleware())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "valid"})
	})
	return r
}

func TestValidateJSON_ValidJSON(t *testing.T) {
	r := setupRouterWithValidateJSON()

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"key":"value"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "valid")
}

func TestValidateJSON_InvalidJSON(t *testing.T) {
	r := setupRouterWithValidateJSON()

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`{"key":`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid JSON")
}

func TestValidateJSON_UnsupportedContentType(t *testing.T) {
	r := setupRouterWithValidateJSON()

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(`<xml></xml>`))
	req.Header.Set("Content-Type", "application/xml")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestValidateJSON_NonModifyingMethod(t *testing.T) {
	r := setupRouterWithValidateJSON()

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// no json to validate, still hits route
	assert.Equal(t, 200, w.Code)
}
