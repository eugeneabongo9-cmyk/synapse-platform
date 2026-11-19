package memory

import (
	"context"

	"github.com/eugeneabongo9-cmyk/synapse-platform/services/orchestrator/internal/domain"
)

// InMemoryEventBus is a simple in-memory event bus for demonstration.
type InMemoryEventBus struct {
	handlers map[string][]func(domain.DomainEvent)
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		handlers: make(map[string][]func(domain.DomainEvent)),
	}
}

func (b *InMemoryEventBus) Publish(ctx context.Context, events []domain.DomainEvent) error {
	for _, event := range events {
		if handlers, ok := b.handlers[event.EventType()]; ok {
			for _, handler := range handlers {
				handler(event)
			}
		}
	}
	return nil
}

func (b *InMemoryEventBus) Subscribe(eventType string, handler func(domain.DomainEvent)) {
	b.handlers[eventType] = append(b.handlers[eventType], handler)
}
