package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/adapters/http"
	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/adapters/memory"
	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/application"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	repo := memory.NewAgentRepository()
	eventBus := memory.NewInMemoryEventBus()

	service := application.NewOrchestratorService(repo, eventBus, logger)
	handler := http.NewHandler(service, logger)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info().Msg("Synapse Orchestrator listening on :8080")
		logger.Info().Msg("Open http://localhost:8080 in your browser for the interactive demo")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("forced shutdown")
	}
	logger.Info().Msg("Orchestrator stopped")
}
