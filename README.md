# Synapse Platform

> Production-grade distributed AI agent orchestration engine built to demonstrate staff/principal-level software engineering.

[![CI](https://github.com/eugeneabongo9-cmyk/synapse-platform/actions/workflows/ci.yml/badge.svg)](https://github.com/eugeneabongo9-cmyk/synapse-platform/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## What This Project Demonstrates

This repository is intentionally engineered as a **showcase of advanced software engineering practices** expected from engineers with 15-20+ years of experience:

- **Hexagonal / Ports & Adapters Architecture** with clear domain boundaries
- **Event-driven microservices** with at-least-once delivery guarantees
- **Polyglot implementation** (Go for orchestration, Rust for high-performance agent runtime)
- **Domain-Driven Design** tactical patterns (Aggregates, Value Objects, Domain Events)
- **Observability-first design** (OpenTelemetry traces, metrics, structured logging)
- **Reliability engineering**: circuit breakers, retries with backoff, idempotency keys
- **Production-ready packaging**: multi-stage Docker builds, docker-compose, K8s manifests
- **Comprehensive CI/CD** with multi-language testing, linting, and security scanning
- **Clean architecture** that separates domain logic from infrastructure concerns

## Architecture Overview

```
┌───────────────────────────────────────────────┐
│                   API Gateway                        │
│              (REST + GraphQL + WebSocket)            │
└──────────────────────────────┬────────────────┘
                               │
         ┌───────────────────────┼─────────────────────┐
         │                    │                    │
┌──────────────────┐    ┌──────────────────┐    ┌──────────────────┐
│  Orchestrator  │    │   Event Bus    │    │ Agent Runtime │
│   (Go / Hex)    │────│  (Redis/NATS)  │────│   (Rust)       │
│               │    │               │    │              │
│ Domain Model  │    │ Domain Events │    │ High-perf     │
│ Aggregates    │    │ Outbox Pattern│    │ Execution     │
└───────────────────┘    └──────────────────┘    └──────────────────┘
```

## Quick Start

```bash
# Start the full platform
docker compose up --build

# The API will be available at http://localhost:8080
# Grafana dashboards at http://localhost:3000 (admin/admin)
```

## Key Design Principles

1. **Domain isolation** — Business logic has zero infrastructure dependencies
2. **Explicit contracts** — All boundaries use well-defined ports and adapters
3. **Observability by default** — Every critical path emits traces and metrics
4. **Fail gracefully** — Circuit breakers, bulkheads, and graceful degradation
5. **Testability** — Domain layer is 100% unit-testable without mocks of infrastructure

## Project Structure

```
synapse-platform/
├── services/
│   ├── orchestrator/     # Go hexagonal architecture core
│   ├── agent-runtime/    # Rust high-performance execution engine
│   └── api-gateway/      # GraphQL + REST + WebSocket gateway
├── infrastructure/
│   ├── docker-compose.yml
│   └── k8s/
├── docs/
│   ├── architecture.md
│   └── design-decisions.md
└── .github/workflows/
```

## Technology Choices & Rationale

| Component | Technology | Why |
|-----------|------------|-----|
| Orchestrator | Go | Excellent concurrency, strong typing, mature ecosystem for systems work |
| Agent Runtime | Rust | Zero-cost abstractions, memory safety, extreme performance for execution |
| Event Bus | Redis Streams / NATS | Lightweight yet production-capable for demonstration |
| Observability | OpenTelemetry | Vendor-neutral, industry standard |
| API | GraphQL + REST | Flexibility for clients + simplicity for internal services |

## What Hiring Managers Look For

This codebase demonstrates:

- Ability to design systems that scale and remain maintainable
- Deep understanding of distributed systems trade-offs
- Production mindset (observability, reliability, testing)
- Clean code that other engineers can extend
- Thoughtful technology selection with clear rationale

## Contributing

This is a showcase repository. Issues and PRs demonstrating improvements to architecture, reliability, or performance are welcome.

## License

MIT
