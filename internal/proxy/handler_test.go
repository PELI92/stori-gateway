package proxy_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"stori-gateway/internal/config"
	"stori-gateway/internal/proxy"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockConfig struct {
	APIKey   string
	Services map[string]config.Service
}

func (m *mockConfig) GetAPIKey() string { return m.APIKey }
func (m *mockConfig) GetServiceURL(name string) (string, bool) {
	svc, ok := m.Services[name]
	return svc.URL, ok
}
func (m *mockConfig) GetAllServices() map[string]config.Service { return m.Services }

func TestHandle_RoutesToService(t *testing.T) {
	// backend mock
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Header().Set("X-Mock-Response", "true")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("echo: " + string(body)))
	}))
	defer backend.Close()

	cfg := &mockConfig{
		Services: map[string]config.Service{
			"auth": {URL: backend.URL},
		},
	}

	// proxy handler
	handler := proxy.NewReverseProxy(cfg)

	// gin test context
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Any("/api/:service/*path", handler.Handle)

	reqBody := `{"test":"data"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/echo", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "true", w.Header().Get("X-Mock-Response"))
	assert.Contains(t, w.Body.String(), "echo: "+reqBody)
}

func TestHandle_UnknownService(t *testing.T) {
	cfg := &mockConfig{
		Services: map[string]config.Service{}, // vacío
	}

	handler := proxy.NewReverseProxy(cfg)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Any("/api/:service/*path", handler.Handle)

	req := httptest.NewRequest(http.MethodGet, "/api/unknown/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadGateway, w.Code)
	assert.Contains(t, w.Body.String(), "unknown service")
}
