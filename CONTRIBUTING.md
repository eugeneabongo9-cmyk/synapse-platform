# Contributing to Synapse Platform

This repository is designed as a demonstration of advanced software engineering practices. Contributions that improve architecture, reliability, observability, or documentation are especially welcome.

## Development Setup

```bash
# Install dependencies
cd services/orchestrator && go mod tidy
cd ../agent-runtime && cargo fetch

# Run tests
make test

# Start the platform
make up
```

## Code Standards

- Domain logic must remain free of infrastructure dependencies
- All public APIs must include tests
- Follow the existing hexagonal architecture patterns
- Add architecture decision records for significant changes

## Pull Request Process

1. Ensure all tests pass
2. Update documentation if behavior changes
3. Add or update ADRs for architectural decisions
4. Request review
