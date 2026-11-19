# Architecture Decision Records

## ADR-004: Transactional Outbox over Dual Writes

**Status**: Accepted  
**Context**: We must guarantee that agent state changes and domain events are consistent.  
**Decision**: Write both in one Postgres transaction. A separate process publishes from the outbox.  
**Consequences**: Eventual consistency for consumers; operational need to monitor outbox lag; excellent reliability.  
**Alternatives considered**: CDC (Debezium) — more operational complexity for this showcase.

## ADR-005: Redis Streams as Event Bus

**Status**: Accepted  
**Context**: Need a lightweight, production-capable bus for the demo that supports consumer groups.  
**Decision**: Redis Streams.  
**Consequences**: Simple to run locally; can be swapped for NATS/Kafka without domain changes.

## ADR-006: OpenTelemetry + Prometheus

**Status**: Accepted  
**Context**: Production systems must be observable.  
**Decision**: OTel for traces, Prometheus client for metrics.  
**Consequences**: Vendor-neutral; integrates with existing Grafana/Jaeger stack.

## ADR-007: Optimistic Concurrency via Version Column

**Status**: Accepted  
**Context**: Multiple requests may try to transition the same agent.  
**Decision**: Increment `version` on every write; reject if version mismatch (future enhancement).  
**Consequences**: Simple and effective for moderate contention.
