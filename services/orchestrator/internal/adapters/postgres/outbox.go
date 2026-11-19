package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type OutboxProcessor struct {
	pool   *pgxpool.Pool
	redis  *redis.Client
	logger zerolog.Logger
}

func NewOutboxProcessor(pool *pgxpool.Pool, rdb *redis.Client, logger zerolog.Logger) *OutboxProcessor {
	return &OutboxProcessor{pool: pool, redis: rdb, logger: logger}
}

func (p *OutboxProcessor) Start(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	p.logger.Info().Msg("outbox processor started")

	for {
		select {
		case <-ctx.Done():
			p.logger.Info().Msg("outbox processor stopped")
			return
		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

func (p *OutboxProcessor) processBatch(ctx context.Context) {
	rows, err := p.pool.Query(ctx, `
		SELECT id, aggregate_id, event_type, payload
		FROM outbox
		WHERE processed_at IS NULL
		ORDER BY created_at
		LIMIT 50
		FOR UPDATE SKIP LOCKED
	`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, aggID uuid.UUID
		var eventType string
		var payload []byte

		if err := rows.Scan(&id, &aggID, &eventType, &payload); err != nil {
			continue
		}

		// Publish to Redis Stream
		err = p.redis.XAdd(ctx, &redis.XAddArgs{
			Stream: "synapse.events",
			Values: map[string]interface{}{
				"event_type": eventType,
				"aggregate_id": aggID.String(),
				"payload": string(payload),
			},
		}).Err()

		if err != nil {
			p.logger.Error().Err(err).Msg("failed to publish event")
			continue
		}

		// Mark as processed
		_, _ = p.pool.Exec(ctx, `UPDATE outbox SET processed_at = NOW() WHERE id = $1`, id)
		p.logger.Debug().Str("event", eventType).Msg("event published from outbox")
	}
}
