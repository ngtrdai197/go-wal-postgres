package publicapi

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go-wal/config"
	"go-wal/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	config.Init()
	InitProvider()
}

func Launch() {
	r := gin.New()
	r.Use(
		middleware.GinLogger(),
		middleware.Recovery(),
	) // Add middleware here

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Register routes

	// Setup gin server with graceful shutdown
	server := &http.Server{
		Addr:    config.Config.APIInfo.PublicApiListen,
		Handler: r,
	}
	// Start the server in a separate goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Msgf("ListenAndServe: %v", err)
		}
	}()
	log.Info().Msgf("Server started on %s", config.Config.APIInfo.PublicApiListen)

	// Wait for an interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutting down server...")

	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Msgf("Server shutdown failed: %v", err)
	}

	log.Info().Msg("Server shutdown successfully")
}
