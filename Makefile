.PHONY: build test lint up down clean

build:
	cd services/orchestrator && go build -o ../../bin/orchestrator ./cmd/server
	cd services/agent-runtime && cargo build --release

test:
	cd services/orchestrator && go test -v -race ./...
	cd services/agent-runtime && cargo test

lint:
	cd services/orchestrator && golangci-lint run
	cd services/agent-runtime && cargo clippy -- -D warnings

up:
	docker compose -f infrastructure/docker-compose.yml up --build

down:
	docker compose -f infrastructure/docker-compose.yml down -v

clean:
	rm -rf bin/
	cd services/agent-runtime && cargo clean
