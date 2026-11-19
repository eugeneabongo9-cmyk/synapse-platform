.PHONY: up down build logs

up:
	docker compose -f infrastructure/docker-compose.yml up --build

down:
	docker compose -f infrastructure/docker-compose.yml down -v

build:
	docker compose -f infrastructure/docker-compose.yml build

logs:
	docker compose -f infrastructure/docker-compose.yml logs -f orchestrator

test:
	cd services/orchestrator && go test -v ./...
