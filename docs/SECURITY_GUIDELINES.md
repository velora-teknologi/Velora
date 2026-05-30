# Security Guidelines

Requirements

- JWT
- Refresh Tokens
- RBAC
- Rate Limiting
- Input Validation
- Audit Logging

Never

- Store plaintext passwords
- Hardcode secrets
- Disable TLS

Passwords

bcrypt minimum cost 12

API Keys

hashed before storage
