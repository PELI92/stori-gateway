package docs

import (
	_ "encoding/json"
	"net/http"
	"stori-gateway/internal/config"

	"github.com/gin-gonic/gin"
)

func Handler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfgData := getSwaggerSpec(cfg)
		c.JSON(http.StatusOK, cfgData)
	}
}

func getSwaggerSpec(cfg *config.Config) map[string]interface{} {
	spec := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]string{
			"title":   "API Gateway",
			"version": "1.0.0",
		},
		"paths": map[string]interface{}{},
	}

	cfgServices := cfg.GetAllServices()

	for name := range cfgServices {
		path := "/api/" + name + "/{any}"
		spec["paths"].(map[string]interface{})[path] = map[string]interface{}{
			"get": map[string]interface{}{
				"summary": "Proxy to " + name,
				"responses": map[string]interface{}{
					"200": map[string]string{"description": "OK"},
				},
			},
		}
	}

	return spec
}
