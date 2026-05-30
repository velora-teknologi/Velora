.PHONY: help backend frontend docker-up docker-down docker-build migrate-up migrate-down test lint clean

help:
	@echo "Velora Development Commands"
	@echo ""
	@echo "Backend:"
	@echo "  make backend              - Run backend server"
	@echo "  make backend-build        - Build backend binary"
	@echo "  make backend-test         - Run backend tests"
	@echo ""
	@echo "Frontend:"
	@echo "  make frontend             - Run frontend server"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build         - Build Docker images"
	@echo "  make docker-up            - Start Docker containers"
	@echo "  make docker-down          - Stop Docker containers"
	@echo ""
	@echo "Database:"
	@echo "  make migrate-up           - Run database migrations"
	@echo "  make migrate-down         - Rollback database migrations"
	@echo ""
	@echo "Tools:"
	@echo "  make lint                 - Run linter"
	@echo "  make test                 - Run all tests"
	@echo "  make clean                - Clean build artifacts"

# Backend commands
backend:
	@echo "Starting backend server..."
	cd backend && go run cmd/server/main.go

backend-build:
	@echo "Building backend..."
	cd backend && go build -o bin/server cmd/server/main.go

backend-test:
	@echo "Running backend tests..."
	cd backend && go test ./...

# Frontend commands
frontend:
	@echo "Starting frontend server..."
	cd frontend && npm run dev

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs:
	docker-compose logs -f

# Database commands
db-up:
	@echo "Ensuring database is running..."
	sudo systemctl start postgresql

db-down:
	@echo "Stopping database..."
	sudo systemctl stop postgresql

db-init:
	@echo "Initializing database..."
	cd backend && bash scripts/init-db.sh

db-seed:
	@echo "Seeding development data..."
	cd backend && bash scripts/seed-db.sh

db-reset: db-down db-up db-init db-seed
	@echo "Database reset complete!"

redis-up:
	@echo "Starting Redis..."
	sudo systemctl start redis-server

redis-down:
	@echo "Stopping Redis..."
	sudo systemctl stop redis-server

# Development commands
dev: docker-up
	@echo "Development environment is ready!"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:3000"

lint:
	@echo "Running linter..."
	cd backend && golangci-lint run ./...

test: backend-test
	@echo "All tests passed!"

clean:
	@echo "Cleaning build artifacts..."
	cd backend && rm -rf bin/
	find . -type f -name "*.test" -delete
	find . -type f -name "coverage.out" -delete
