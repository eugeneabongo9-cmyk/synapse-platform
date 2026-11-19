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
	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/domain"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Dependency injection following hexagonal architecture
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
		logger.Info().Msg("Starting Synapse Orchestrator on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited")
}
