# Architecture Deep Dive

## Hexagonal Architecture in the Orchestrator

The orchestrator service follows strict hexagonal architecture:

```
                    ┌─────────────────────────────────────┐
                    │           Domain Layer              │
                    │  (Pure business logic, no deps)    │
                    │                                    │
                    │  • Aggregates (Agent, Workflow)  │
                    │  • Value Objects                 │
                    │  • Domain Events                 │
                    │  • Domain Services               │
                    └─────────────────────────────────────┘
                                      ▲
                                      │ Ports (interfaces)
                    ┌─────────────────────────────────────┐
                    │        Application Layer            │
                    │   (Use cases, orchestration)       │
                    └─────────────────────────────────────┘
                                      ▲
                                      │
         ┌──────────────────────────────────────────────┐
         │              Adapters                     │
         │  • HTTP/GraphQL   • Redis Event Bus     │
         │  • PostgreSQL     • OpenTelemetry       │
         │  • WebSocket      • Circuit Breaker     │
         └──────────────────────────────────────────────┘
```

## Event-Driven Communication

All inter-service communication uses domain events with the **Transactional Outbox** pattern to guarantee at-least-once delivery without dual-write problems.

### Event Flow

1. Domain aggregate raises domain event
2. Event is persisted in the same transaction as the aggregate state change
3. Background outbox processor publishes to the event bus
4. Consumers process events idempotently

## Reliability Patterns

- **Circuit Breaker**: Protects against cascading failures to external agent runtimes
- **Retry with Exponential Backoff + Jitter**: For transient failures
- **Idempotency Keys**: All commands carry unique identifiers
- **Bulkhead Isolation**: Separate thread pools for different concern areas
- **Graceful Shutdown**: Proper signal handling and in-flight request draining

## Observability

Every service emits:

- Distributed traces (OpenTelemetry)
- Business and technical metrics
- Structured logs with correlation IDs
- Health and readiness probes
