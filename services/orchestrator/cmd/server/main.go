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
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		otlptracehttp.WithEndpoint(getEnv("OTEL_ENDPOINT", "jaeger:4318")),
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
		logger.Info().Msg("OpenTelemetry tracing enabled")
	}

	// Postgres
	dbURL := getEnv("DATABASE_URL", "postgres://synapse:synapse@postgres:5432/synapse?sslmode=disable")
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatal("postgres: ", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("postgres ping: ", err)
	}

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
			processed_at TIMESTAMPTZ,
			retry_count INT DEFAULT 0,
			next_retry_at TIMESTAMPTZ
		);
	`)
	if err != nil {
		log.Fatal("migration: ", err)
	}

	// Redis
	redisAddr := getEnv("REDIS_URL", "redis:6379")
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("redis: ", err)
	}

	repo := postgres.NewAgentRepository(pool)
	eventBus := redisadapter.NewEventBus(rdb)
	service := application.NewOrchestratorService(repo, eventBus, logger)

	outbox := postgres.NewOutboxProcessor(pool, rdb, logger)
	outboxCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go outbox.Start(outboxCtx)

	handler := http.NewHandler(service, logger)
	mux := handler.Router()
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		logger.Info().Msg("Synapse Orchestrator (staff-level) ready")
		logger.Info().Msg("Demo UI  → http://localhost:8080")
		logger.Info().Msg("Metrics  → http://localhost:8080/metrics")
		logger.Info().Msg("Jaeger   → http://localhost:16686")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info().Msg("graceful shutdown")
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	_ = server.Shutdown(shutdownCtx)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
