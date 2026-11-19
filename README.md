# Synapse Platform

> Staff / Principal-level distributed systems showcase  
> Hexagonal architecture • Transactional Outbox • Event-driven • Full observability • Kubernetes-ready

## Quickest Way to Run (No Docker Required)

You only need Go installed.

```bash
# Install Go if you don't have it (macOS)
brew install go

# Then run the orchestrator in in-memory mode
cd services/orchestrator
go mod tidy
go run ./cmd/server
```

Open **http://localhost:8080**

This mode uses an in-memory repository so you can demo the domain model, API, and UI immediately.

## Full Production Stack (Docker)

If you have Docker:

```bash
docker compose -f infrastructure/docker-compose.yml up --build
```

This starts Postgres + Redis + Jaeger + the orchestrator with the transactional outbox.

## Why This Is Staff-Level

- Hexagonal architecture (domain has zero infrastructure deps)
- Domain-Driven Design with invariants and domain events
- Transactional Outbox pattern
- OpenTelemetry tracing + Prometheus metrics
- Kubernetes manifests with HPA
- Graceful degradation (works with or without external services)

## License

MIT
