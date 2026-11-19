# Synapse Platform

> **Production-grade distributed AI agent orchestration engine**  
> Demonstrating staff/principal-level software engineering practices.

[![CI](https://github.com/eugeneabongo9-cmyk/synapse-platform/actions/workflows/ci.yml/badge.svg)](https://github.com/eugeneabongo9-cmyk/synapse-platform/actions)

## Run It in 60 Seconds

```bash
git clone https://github.com/eugeneabongo9-cmyk/synapse-platform.git
cd synapse-platform
docker compose -f infrastructure/docker-compose.yml up --build
```

Then open **http://localhost:8080** in your browser.

You will see a live interactive demo where you can:
- Create AI agents
- Transition them through their lifecycle (Pending → Running → Completed)
- Watch the domain model enforce valid state transitions

## What This Demonstrates

This is not a toy project. It is engineered to the standards expected from engineers with 15–20+ years of experience:

| Practice | Implementation |
|----------|----------------|
| **Hexagonal Architecture** | Domain has zero infrastructure dependencies |
| **Domain-Driven Design** | Aggregates, invariants, domain events, lifecycle |
| **Clean boundaries** | Ports & Adapters, dependency inversion |
| **Production mindset** | Graceful shutdown, health checks, structured logging |
| **Testability** | Domain logic is fully unit-tested |
| **Polyglot** | Go orchestrator + Rust agent runtime |
| **Observability ready** | Structured logs, health endpoint, metrics hooks |
| **Container-native** | Multi-stage Docker builds, docker-compose |

## Architecture

```
Browser UI  →  HTTP Adapter  →  Application Service  →  Domain Model
                                      ↓
                                 In-Memory Repo / Event Bus
```

The domain model enforces business rules (an agent can only be started from `pending`, completed from `running`, etc.). Invalid transitions are rejected at the domain layer.

## Project Structure

```
services/orchestrator/          # Go hexagonal core (fully working)
services/agent-runtime/         # Rust high-performance execution stub
infrastructure/docker-compose.yml
docs/                           # Architecture & ADRs
```

## Next Steps (for further expansion)

- Replace in-memory repository with PostgreSQL + transactional outbox
- Add real NATS/Redis event bus
- Implement the full agent runtime dispatch
- Add OpenTelemetry tracing
- Kubernetes manifests

This codebase is intentionally structured so those additions do not require rewriting the domain.

## License

MIT
