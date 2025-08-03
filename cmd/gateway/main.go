package main

import (
	"os"
	"stori-gateway/internal/config"
	"stori-gateway/internal/middleware"
	"stori-gateway/internal/proxy"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg := config.LoadConfig()

	r := gin.New()

	// This configuration is used to suppress warning. This is ok because we are not going to use a LB or reverse proxy before this service. (this service will be the reverse proxy)
	err := r.SetTrustedProxies(nil)
	if err != nil {
		return
	}

	// Zerolog configuration
	r.Use(middleware.ZerologMiddleware())

	// Used to handle and recover from panics, returning 500
	r.Use(gin.Recovery())

	r.Any("/api/:service/*path", proxy.Handler(cfg))

	if err := r.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}

}
