#!/bin/bash
# Database initialization script
# This script creates the initial database and runs migrations

set -e

# Configuration
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-velora}"
DB_PASSWORD="${DB_PASSWORD:-}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Starting database initialization...${NC}"

# Check if PostgreSQL is running
echo "Checking PostgreSQL connection..."
if ! pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" 2>/dev/null; then
    echo -e "${RED}Error: PostgreSQL is not running or not accessible at $DB_HOST:$DB_PORT${NC}"
    exit 1
fi

echo -e "${GREEN}✓ PostgreSQL is accessible${NC}"

# Create database if it doesn't exist
echo "Creating database '$DB_NAME' if it doesn't exist..."
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -tc \
    "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME'" | grep -q 1 || \
    PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c \
    "CREATE DATABASE $DB_NAME;"

echo -e "${GREEN}✓ Database created or already exists${NC}"

# Run migrations
echo "Running migrations..."
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
    -f migrations/001_create_users_table.sql \
    -f migrations/002_create_agents_table.sql \
    -f migrations/003_create_workflows_table.sql

echo -e "${GREEN}✓ Migrations completed${NC}"

# Verify tables were created
echo "Verifying tables..."
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" \
    -c "\dt"

echo -e "${GREEN}✓ Database initialization complete${NC}"
echo ""
echo "Database Details:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  User: $DB_USER"
echo "  Database: $DB_NAME"
