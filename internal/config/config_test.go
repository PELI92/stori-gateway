package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAPIKey(t *testing.T) {
	cfg := &Config{
		apiKey: "real-api-key",
	}
	assert.Equal(t, "real-api-key", cfg.GetAPIKey())
}

func TestGetServiceURL(t *testing.T) {
	cfg := &Config{
		services: map[string]Service{
			"auth": {URL: "http://localhost:9001"},
		},
	}

	url, ok := cfg.GetServiceURL("auth")
	assert.True(t, ok)
	assert.Equal(t, "http://localhost:9001", url)

	_, ok = cfg.GetServiceURL("missing")
	assert.False(t, ok)
}

func TestGetAllServices(t *testing.T) {
	cfg := &Config{
		services: map[string]Service{
			"user": {URL: "http://localhost:9002"},
		},
	}

	all := cfg.GetAllServices()
	assert.Len(t, all, 1)
	assert.Equal(t, "http://localhost:9002", all["user"].URL)
}
