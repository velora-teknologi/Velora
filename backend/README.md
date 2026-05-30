# Velora Backend

Enterprise-grade AI Automation Platform Backend built with Go, Fiber, and PostgreSQL.

## Architecture

This project follows **Clean Architecture** and **Domain-Driven Design (DDD)** principles:

```
backend/
├── cmd/server/           # Application entry point
├── internal/
│   ├── domain/           # Domain layer (models, repositories interfaces)
│   ├── application/      # Application layer (services, DTOs)
│   └── infrastructure/   # Infrastructure layer (handlers, config, database)
├── pkg/                  # Shared packages (logger, middleware, errors)
└── tests/                # Unit and integration tests
```

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL 16+
- Redis 7+
- NATS 2.8+

### Setup

1. **Install dependencies:**
   ```bash
   cd backend
   go mod download
   ```

2. **Configure environment:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Run the server:**
   ```bash
   go run cmd/server/main.go
   ```

### Docker Setup

```bash
# Build all services
docker-compose build

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

## Project Structure

### Domain Layer (`internal/domain/`)

Contains the core business logic and domain entities:

- `models/`: Domain entities (User, Agent, Workflow, etc.)
- `repositories/`: Repository interfaces (contracts for data access)

### Application Layer (`internal/application/`)

Implements business use cases:

- `services/`: Business logic and orchestration
- `dtos/`: Data Transfer Objects for API requests/responses

### Infrastructure Layer (`internal/infrastructure/`)

Implements technical concerns:

- `handlers/`: HTTP request handlers
- `persistence/`: Repository implementations using GORM
- `config/`: Configuration and external service initialization
- `routes/`: API route definitions

### Shared Packages (`pkg/`)

- `logger/`: Structured logging with Zap
- `middleware/`: HTTP middleware (JWT, error handling)
- `errors/`: Custom error types and handling

## API Endpoints

### Health Check
- `GET /health` - Server health check

### Users
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `GET /api/v1/users` - List users

### Authentication
- `POST /api/v1/auth/login` - Login user

## Development

### Running Tests

```bash
go test ./...

# With coverage
go test -cover ./...

# With detailed output
go test -v ./...
```

### Linting

```bash
golangci-lint run ./...
```

### Building

```bash
go build -o bin/server cmd/server/main.go
```

## Key Features

✅ Clean Architecture with DDD
✅ Dependency Injection
✅ Structured Logging with Zap
✅ Error Handling
✅ JWT Authentication
✅ PostgreSQL with GORM
✅ Redis Caching
✅ NATS Event Bus
✅ Docker & Docker Compose
✅ Unit Tests

## Configuration

See `.env.example` for all available configuration options.

## Database Migrations

Migrations are managed through GORM models. Create new migrations by:

1. Define models in `internal/domain/models/`
2. Run migrations in application startup

## Contributing

Follow these guidelines:

1. Use dependency injection
2. Add logging for important operations
3. Add error handling and validation
4. Write unit tests for services
5. Follow Go conventions and style guide
6. Never hardcode secrets
7. Don't put business logic in handlers

## License

Proprietary - Velora Technologies
