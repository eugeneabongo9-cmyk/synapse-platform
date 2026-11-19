package redisadapter

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

// Publish is a no-op here because events are published reliably via the transactional outbox.
// This satisfies the application port.
func (b *EventBus) Publish(ctx context.Context, events []domain.DomainEvent) error {
	return nil
}
