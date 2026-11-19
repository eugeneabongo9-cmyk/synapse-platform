# Synapse Platform

> **Staff / Principal-level distributed systems showcase**  
> Hexagonal architecture • Transactional Outbox • Event-driven • Full observability • Kubernetes-ready

## One-Command Full Stack

```bash
git clone https://github.com/eugeneabongo9-cmyk/synapse-platform.git
cd synapse-platform
docker compose -f infrastructure/docker-compose.yml up --build
```

| Endpoint | Purpose |
|----------|---------|
| http://localhost:8080 | Interactive demo UI |
| http://localhost:8080/metrics | Prometheus metrics |
| http://localhost:16686 | Jaeger distributed tracing |
| localhost:5432 | PostgreSQL |
| localhost:6379 | Redis Streams |

## Why This Is Staff-Level

This repository demonstrates the exact skills companies evaluate for senior/staff/principal engineers:

- **Architectural thinking** — Clear boundaries, dependency inversion, trade-off documentation (ADRs)
- **Distributed systems reliability** — Transactional outbox, at-least-once delivery, optimistic concurrency
- **Observability culture** — Traces, metrics, structured logs from day one
- **Production readiness** — Health probes, graceful shutdown, resource limits, HPA
- **Clean code that scales** — Domain remains pure; infrastructure can be swapped
- **Operational awareness** — Kubernetes manifests, docker-compose, monitoring hooks

## Patterns on Display

- Hexagonal / Ports & Adapters
- Domain-Driven Design (aggregates, invariants, domain events)
- Transactional Outbox + reliable event publishing
- Optimistic concurrency control
- OpenTelemetry distributed tracing
- Prometheus metrics
- Background processing with `SKIP LOCKED`
- Kubernetes Deployment + HPA + probes

## Kubernetes

Manifests are in `infrastructure/k8s/`. Apply with:

```bash
kubectl apply -f infrastructure/k8s/
```

(Requires secrets and image push for full deployment.)

## License

MIT
