Generate only production-ready code.

Requirements:

- Golang
- Fiber
- PostgreSQL
- Redis
- GORM
- JWT
- NATS

Architecture:

- Clean Architecture
- DDD
- SOLID

Every feature must include:

- handler.go
- service.go
- repository.go
- model.go
- dto.go
- routes.go
- tests

Always:

- Use dependency injection
- Add logging
- Add validation
- Add error handling
- Add unit tests

Never:

- Hardcode secrets
- Put business logic in handlers
- Use global mutable state

Target:

Enterprise-grade SaaS AI Automation Platform.
