# Architecture Decision Records (ADRs)

## ADR-001: Hexagonal Architecture for Orchestrator

**Status**: Accepted

**Context**: We need a maintainable core that can evolve as requirements change without rewriting infrastructure.

**Decision**: Adopt hexagonal architecture with strict dependency inversion.

**Consequences**:
- Domain logic is 100% unit testable
- Infrastructure can be swapped (e.g., Postgres → DynamoDB) without domain changes
- Slightly more upfront ceremony

## ADR-002: Polyglot Services

**Status**: Accepted

**Context**: Different parts of the system have different performance and safety requirements.

**Decision**:
- Go for orchestration and API layers (productivity + concurrency)
- Rust for agent runtime (performance + memory safety)

**Consequences**:
- Higher cognitive load for the team
- Demonstrates ability to choose the right tool for the job
- Clear service boundaries prevent language leakage

## ADR-003: Event-Driven with Outbox Pattern

**Status**: Accepted

**Context**: Need reliable inter-service communication without distributed transactions.

**Decision**: Use domain events + transactional outbox.

**Consequences**:
- Eventual consistency
- Guaranteed delivery
- Increased operational complexity (outbox processor)
