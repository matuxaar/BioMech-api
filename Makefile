.PHONY: help build-backend build-ml up down logs lint-backend test-backend clean

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build-backend: ## Build Go backend binary
	cd backend && go build -o bin/api ./cmd/api

build-ml: ## Build ML service Docker image
	docker compose build ml

up: ## Start full stack (db + api + ml)
	docker compose up --build -d

up-backend: ## Start backend stack (db + api, no ML)
	docker compose -f backend/docker-compose.yml up --build -d

down: ## Stop all services
	docker compose down

logs: ## Tail logs from all services
	docker compose logs -f

logs-api: ## Tail API logs
	docker compose logs -f api

lint-backend: ## Run go vet on backend
	cd backend && go vet ./...

test-backend: ## Run Go tests
	cd backend && go test ./...

seed: ## Seed synthetic test data via scripts/seed_data.py
	python scripts/seed_data.py

clean: ## Clean build artifacts
	rm -rf backend/bin/
	find . -name '__pycache__' -type d -exec rm -rf {} + 2>/dev/null || true
