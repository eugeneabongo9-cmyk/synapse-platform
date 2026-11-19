package redis

import (
	"context"

	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/domain"
	"github.com/redis/go-redis/v9"
)

type EventBus struct {
	client *redis.Client
}

func NewEventBus(client *redis.Client) *EventBus {
	return &EventBus{client: client}
}

func (b *EventBus) Publish(ctx context.Context, events []domain.DomainEvent) error {
	// Events are published via the outbox processor, not directly.
	// This adapter exists to satisfy the port when needed for other use cases.
	return nil
}
