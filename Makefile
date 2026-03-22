.PHONY: dev-api dev-worker dev-web dev-infra migrate-up migrate-down seed sqlc test build clean

# Development
dev-api:
	go run ./cmd/api

dev-worker:
	go run ./cmd/worker

dev-web:
	cd web && npm run dev

dev-infra:
	docker compose up -d

dev-infra-down:
	docker compose down

# Database
migrate-up:
	migrate -path migrations -database "$${DATABASE_URL:-postgres://curator:curator@localhost:5432/curator?sslmode=disable}" up

migrate-down:
	migrate -path migrations -database "$${DATABASE_URL:-postgres://curator:curator@localhost:5432/curator?sslmode=disable}" down 1

seed:
	go run ./cmd/seed

collect:
	go run ./cmd/collect

discover:
	go run ./cmd/discover

# Code generation
sqlc:
	sqlc generate

# Testing
test:
	go test $(shell go list ./... | grep -v /web/) -v -count=1

test-cover:
	go test $(shell go list ./... | grep -v /web/) -coverprofile=coverage.out && go tool cover -html=coverage.out

# Build
build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker
	go build -o bin/seed ./cmd/seed
	cd web && npm run build

clean:
	rm -rf bin/

# Lint
lint:
	golangci-lint run ./...
