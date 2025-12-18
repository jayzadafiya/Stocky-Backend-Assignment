.PHONY: help migrate run dev watch build clean test

help: 
	@echo "Available commands:"
	@echo "  make migrate    - Run database migrations and seed data"
	@echo "  make run        - Start the server in production mode"
	@echo "  make dev        - Start the server in development mode"
	@echo "  make watch      - Start server with hot-reload (requires air)"
	@echo "  make build      - Build the application"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make test       - Run tests"

migrate: 
	go run cmd/migrate/main.go

run: 
	go run main.go

dev: 
	GIN_MODE=debug LOG_LEVEL=debug go run main.go

watch:
	@command -v air > /dev/null 2>&1 || { echo "air is not installed. Install with: go install github.com/air-verse/air@latest"; exit 1; }
	@echo "Starting server with hot-reload..."
	GIN_MODE=debug LOG_LEVEL=debug air

build: 
	go build -o bin/stocky-backend main.go
	go build -o bin/migrate cmd/migrate/main.go

clean: 
	rm -rf bin/
	go clean

test: 
	go test -v ./...
