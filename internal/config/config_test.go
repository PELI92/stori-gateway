package config_test

import (
	"stori-gateway/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockConfig struct {
	APIKey   string
	Services map[string]config.Service
}

func (m *mockConfig) GetAPIKey() string {
	return m.APIKey
}

func (m *mockConfig) GetServiceURL(name string) (string, bool) {
	svc, ok := m.Services[name]
	return svc.URL, ok
}

func (m *mockConfig) GetAllServices() map[string]config.Service {
	return m.Services
}

func TestGetAPIKey(t *testing.T) {
	mock := &mockConfig{APIKey: "mocked-key"}
	assert.Equal(t, "mocked-key", mock.GetAPIKey())
}

func TestGetServiceURL(t *testing.T) {
	mock := &mockConfig{
		Services: map[string]config.Service{
			"auth": {URL: "http://localhost:9001"},
		},
	}

	url, ok := mock.GetServiceURL("auth")
	assert.True(t, ok)
	assert.Equal(t, "http://localhost:9001", url)

	_, ok = mock.GetServiceURL("notfound")
	assert.False(t, ok)
}

func TestGetAllServices(t *testing.T) {
	mock := &mockConfig{
		Services: map[string]config.Service{
			"user": {URL: "http://localhost:9002"},
		},
	}

	all := mock.GetAllServices()
	assert.Equal(t, 1, len(all))
	assert.Equal(t, "http://localhost:9002", all["user"].URL)
}
