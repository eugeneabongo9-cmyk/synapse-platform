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
	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/adapters/postgres"
	redisadapter "github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/adapters/redis"
	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/application"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	ctx := context.Background()

	// OpenTelemetry
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("jaeger:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err == nil {
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exp),
			sdktrace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName("synapse-orchestrator"),
			)),
		)
		otel.SetTracerProvider(tp)
		defer func() { _ = tp.Shutdown(ctx) }()
		logger.Info().Msg("OpenTelemetry tracing enabled → Jaeger")
	} else {
		logger.Warn().Err(err).Msg("OTel exporter not available, continuing without tracing")
	}

	// Postgres
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed to connect to postgres: ", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("postgres ping failed: ", err)
	}

	// Simple migrations
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS agents (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			version INT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS outbox (
			id UUID PRIMARY KEY,
			aggregate_id UUID NOT NULL,
			event_type TEXT NOT NULL,
			payload JSONB NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			processed_at TIMESTAMPTZ
		);
	`)
	if err != nil {
		log.Fatal("migration failed: ", err)
	}

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("redis connection failed: ", err)
	}

	repo := postgres.NewAgentRepository(pool)
	eventBus := redisadapter.NewEventBus(rdb)
	service := application.NewOrchestratorService(repo, eventBus, logger)

	outbox := postgres.NewOutboxProcessor(pool, rdb, logger)
	outboxCtx, cancelOutbox := context.WithCancel(context.Background())
	defer cancelOutbox()
	go outbox.Start(outboxCtx)

	handler := http.NewHandler(service, logger)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info().Msg("Synapse Orchestrator ready")
		logger.Info().Msg("Demo UI     → http://localhost:8080")
		logger.Info().Msg("Jaeger UI   → http://localhost:16686")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("graceful shutdown started")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
}
