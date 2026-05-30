# Database Schema

## Overview
The Velora database uses PostgreSQL 16 with GORM for ORM. All tables use UUID as primary keys and include soft deletes with `deleted_at` field.

## Tables

### users
Stores user account information and authentication details.

| Column | Type | Constraints | Description |
|--------|------|-----------|-------------|
| id | VARCHAR(36) | PRIMARY KEY | Unique identifier (UUID) |
| email | VARCHAR(255) | UNIQUE, NOT NULL | User email address |
| password | VARCHAR(255) | NOT NULL | Bcrypt hashed password |
| full_name | VARCHAR(255) | | User's full name |
| avatar | VARCHAR(500) | | Avatar URL |
| role | VARCHAR(50) | DEFAULT 'user' | User role (admin, user, etc) |
| is_active | BOOLEAN | DEFAULT true | Account active status |
| created_at | TIMESTAMP | DEFAULT NOW() | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT NOW() | Last update timestamp |
| deleted_at | TIMESTAMP | | Soft delete timestamp |

**Indexes:**
- `email` - For email lookups during login
- `is_active` - For filtering active users
- `deleted_at` - For soft delete queries

---

### agents
Stores AI agent configurations and metadata.

| Column | Type | Constraints | Description |
|--------|------|-----------|-------------|
| id | VARCHAR(36) | PRIMARY KEY | Unique identifier (UUID) |
| name | VARCHAR(255) | NOT NULL | Agent name |
| description | TEXT | | Agent description |
| user_id | VARCHAR(36) | NOT NULL, FK | Owner user ID |
| config | JSONB | DEFAULT '{}' | Agent configuration (JSON) |
| status | VARCHAR(50) | DEFAULT 'active' | Status (active, inactive, etc) |
| created_at | TIMESTAMP | DEFAULT NOW() | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT NOW() | Last update timestamp |
| deleted_at | TIMESTAMP | | Soft delete timestamp |

**Indexes:**
- `user_id` - For querying user's agents
- `status` - For filtering by status
- `deleted_at` - For soft delete queries

**Foreign Keys:**
- `user_id` → `users.id` (ON DELETE CASCADE)

---

### workflows
Stores workflow definitions and automation rules.

| Column | Type | Constraints | Description |
|--------|------|-----------|-------------|
| id | VARCHAR(36) | PRIMARY KEY | Unique identifier (UUID) |
| name | VARCHAR(255) | NOT NULL | Workflow name |
| description | TEXT | | Workflow description |
| agent_id | VARCHAR(36) | NOT NULL, FK | Parent agent ID |
| definition | JSONB | DEFAULT '{}' | Workflow definition (JSON) |
| is_active | BOOLEAN | DEFAULT true | Activation status |
| created_at | TIMESTAMP | DEFAULT NOW() | Creation timestamp |
| updated_at | TIMESTAMP | DEFAULT NOW() | Last update timestamp |
| deleted_at | TIMESTAMP | | Soft delete timestamp |

**Indexes:**
- `agent_id` - For querying agent's workflows
- `is_active` - For filtering active workflows
- `deleted_at` - For soft delete queries

**Foreign Keys:**
- `agent_id` → `agents.id` (ON DELETE CASCADE)

---

## Relationships

```
users (1) ──── (*) agents
              │
              └──── (*) workflows
```

- One user can have multiple agents
- One agent can have multiple workflows
- Deleting a user cascades delete to agents and workflows
- Deleting an agent cascades delete to its workflows

## Migrations

Migrations are stored in `migrations/` directory:

- `001_create_users_table.sql` - Users table with indexes
- `002_create_agents_table.sql` - Agents table with foreign key
- `003_create_workflows_table.sql` - Workflows table with foreign key

### Running Migrations

**Automatic (via Go app):**
```bash
go run cmd/server/main.go  # Runs migrations on startup
```

**Manual (SQL):**
```bash
make db-init  # Run all SQL migration files
```

## GORM Models

All models are defined in `internal/domain/models/models.go`:

```go
// User model with BeforeCreate hook for UUID generation
type User struct {
    ID        string
    Email     string    // Unique
    Password  string
    FullName  string
    Role      string    // Default: "user"
    IsActive  bool      // Default: true
    ...
}

// Agent model
type Agent struct {
    ID          string
    Name        string
    UserID      string    // Foreign key
    Config      string    // JSONB
    Status      string    // Default: "active"
    ...
}

// Workflow model
type Workflow struct {
    ID          string
    Name        string
    AgentID     string    // Foreign key
    Definition  string    // JSONB
    IsActive    bool      // Default: true
    ...
}
```

## Development Database Commands

```bash
make db-up              # Start PostgreSQL service
make db-down            # Stop PostgreSQL service
make db-init            # Initialize database
make db-seed            # Seed sample data
make db-reset           # Reset database (down, up, init, seed)
```

## Timestamps and Soft Deletes

All tables have:
- `created_at` - Automatically set on insert
- `updated_at` - Automatically updated on any change
- `deleted_at` - NULL by default, set on logical delete (soft delete)

PostgreSQL triggers ensure `updated_at` is automatically updated.

## Sample Data

Development sample data includes:
- 3 test users (admin, user, agent-demo)
- 3 sample agents (email assistant, data analyzer, content generator)
- 3 sample workflows

See `scripts/seed-db.sh` for details.

## Indexes Strategy

Indexes are created for:
- **Foreign keys** - For JOIN operations
- **Status fields** - For filtering queries
- **Unique fields** - For constraint checking
- **Soft delete** - For excluding deleted records efficiently

## Backup & Recovery

### Backup Database
```bash
pg_dump -h localhost -U postgres -d velora > backup.sql
```

### Restore Database
```bash
psql -h localhost -U postgres -d velora < backup.sql
```

## Monitoring

Check database connection and tables:
```bash
psql -h localhost -U postgres -d velora -c "\dt+"
```

Check indexes:
```bash
psql -h localhost -U postgres -d velora -c "\di+"
```

Check table sizes:
```bash
psql -h localhost -U postgres -d velora -c "SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) FROM pg_tables WHERE schemaname='public' ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"
```
