package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"stori-gateway/internal/config"
	"stori-gateway/internal/docs"
	"stori-gateway/internal/middleware"
	"stori-gateway/internal/proxy"
)

func main() {
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()

	cfg := config.LoadConfig()

	r := gin.New()

	// This configuration is used to suppress warning. This is ok because we are not going to use a LB or reverse proxy before this service. (this service will be the reverse proxy)
	err := r.SetTrustedProxies(nil)
	if err != nil {
		return
	}

	r.Use(
		middleware.ZeroLogMiddleware(), // Logger configuration // Validates API key
		gin.Recovery(),                 // Used to handle and recover from panics, returning 500

	)

	proxyHandler := proxy.NewReverseProxy(cfg)

	// Public endpoint will not require requestId nor api key
	public := r.Group("/")
	{
		public.GET("/swagger.json", docs.Handler(cfg))
		public.GET("/healthz", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})
	}

	// Protected endpoints will require requestId nor api key
	protected := r.Group("/")
	protected.Use(
		middleware.RequireRequestID(),
		middleware.RequireAPIKey(cfg),
		middleware.ValidateJSONMiddleware(),
	)
	{
		protected.Any("/api/:service/*path", proxyHandler.Handle)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}

}
