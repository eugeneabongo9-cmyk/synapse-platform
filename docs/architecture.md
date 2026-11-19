# Architecture Deep Dive — Staff Level

## Core Principles

1. **Domain purity** — The domain package has zero imports from infrastructure or frameworks.
2. **Transactional consistency** — Aggregate state + domain events are written atomically via the outbox.
3. **At-least-once delivery** — Outbox + idempotent consumers.
4. **Observability by default** — Traces (OTel), metrics (Prometheus), structured logs.
5. **Failure isolation** — Health checks, graceful degradation, future circuit breakers.

## Reliability Patterns Implemented

| Pattern | Implementation | Purpose |
|---------|----------------|---------|
| Transactional Outbox | Postgres `outbox` table + background processor | Prevent dual-write inconsistency |
| Optimistic concurrency | `version` column on agents | Detect concurrent modifications |
| Idempotent publishing | Outbox `processed_at` + consumer deduplication | Safe redelivery |
| Health & readiness | `/health` endpoint + Docker/K8s probes | Orchestration readiness |
| Structured logging | zerolog with correlation | Debuggability |
| Distributed tracing | OpenTelemetry → Jaeger | Request flow visibility |

## Future Patterns (ready to add)

- **Circuit Breaker** around external agent runtime calls
- **Saga orchestration** for multi-agent workflows
- **CQRS read models** for query optimization
- **Event sourcing** of full agent history (current model is state + events)
- **Rate limiting & bulkheads** per tenant

## Data Flow

```
Client → HTTP Adapter → Application Service → Domain Aggregate
                                      ↓
                               Postgres Transaction
                               ├─ agents (state)
                               └─ outbox (events)
                                      ↓
                            Outbox Processor (poll + SKIP LOCKED)
                                      ↓
                                 Redis Stream
                                      ↓
                           (Consumers / Agent Runtime)
```
