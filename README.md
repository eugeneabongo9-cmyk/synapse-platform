# Synapse Platform

> **Production-grade distributed AI agent orchestration engine**  
> Demonstrating staff/principal-level software engineering practices.

## Run the Full Stack

```bash
git clone https://github.com/eugeneabongo9-cmyk/synapse-platform.git
cd synapse-platform
docker compose -f infrastructure/docker-compose.yml up --build
```

| Service | URL |
|---------|-----|
| **Interactive Demo** | http://localhost:8080 |
| **Jaeger Tracing UI** | http://localhost:16686 |
| **Postgres** | localhost:5432 |
| **Redis** | localhost:6379 |

## What This Version Demonstrates

This is a significantly more advanced system:

- **PostgreSQL persistence** with proper schema
- **Transactional Outbox Pattern** — domain events are written in the same transaction as the aggregate, then reliably published
- **Redis Streams** as the event bus
- **OpenTelemetry distributed tracing** (view traces in Jaeger)
- **Hexagonal architecture** still fully intact — domain has zero knowledge of Postgres, Redis, or OTel
- **Background outbox processor** with `SKIP LOCKED` for safe concurrent processing
- **Graceful startup dependencies** via Docker healthchecks

### Architecture Flow

```
HTTP Request
    ↓
Application Service
    ↓
Domain Aggregate (raises events)
    ↓
Postgres Transaction
    ├── Update agents table
    └── Insert into outbox table
            ↓
Outbox Processor (background)
    ↓
Redis Stream (synapse.events)
    ↓
(Future consumers / agent runtime)
```

The domain model remains pure. All infrastructure concerns live in adapters.

## Key Patterns on Display

| Pattern | Where |
|---------|-------|
| Hexagonal Architecture | Entire orchestrator service |
| Transactional Outbox | `postgres/repository.go` + `outbox.go` |
| Idempotent event publishing | Outbox `processed_at` + `SKIP LOCKED` |
| Distributed Tracing | OpenTelemetry + Jaeger |
| Domain invariants | `domain/agent.go` |
| Clean dependency inversion | Ports defined in application layer |

## License

MIT
