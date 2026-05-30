#!/bin/bash
# Database seeding script for development
# This script adds sample data to the database

set -e

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-velora_dev}"
DB_PASSWORD="${DB_PASSWORD:-}"

echo "Seeding development data..."

# Create sample users (passwords are hashed with bcrypt)
PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<EOF
-- Insert sample users
INSERT INTO users (id, email, password, full_name, role, is_active) 
VALUES 
    ('550e8400-e29b-41d4-a716-446655440001', 'admin@velora.dev', '\$2a\$10\$QIvFIAJLWrWJF/H3R8VxvOKdnWBCkIslvJqt.jvpPYIUpADVzJQu2', 'Admin User', 'admin', true),
    ('550e8400-e29b-41d4-a716-446655440002', 'user@velora.dev', '\$2a\$10\$QIvFIAJLWrWJF/H3R8VxvOKdnWBCkIslvJqt.jvpPYIUpADVzJQu2', 'Test User', 'user', true),
    ('550e8400-e29b-41d4-a716-446655440003', 'agent-demo@velora.dev', '\$2a\$10\$QIvFIAJLWrWJF/H3R8VxvOKdnWBCkIslvJqt.jvpPYIUpADVzJQu2', 'Agent Demo', 'user', true)
ON CONFLICT (email) DO NOTHING;

-- Insert sample agents
INSERT INTO agents (id, name, description, user_id, config, status)
VALUES
    ('550e8400-e29b-41d4-a716-446655440101', 'Email Assistant', 'AI agent for email automation', '550e8400-e29b-41d4-a716-446655440001', '{"version": "1.0", "model": "gpt-4"}', 'active'),
    ('550e8400-e29b-41d4-a716-446655440102', 'Data Analyzer', 'AI agent for data analysis', '550e8400-e29b-41d4-a716-446655440001', '{"version": "1.0", "model": "gpt-4"}', 'active'),
    ('550e8400-e29b-41d4-a716-446655440103', 'Content Generator', 'AI agent for content generation', '550e8400-e29b-41d4-a716-446655440002', '{"version": "1.0", "model": "gpt-4"}', 'inactive')
ON CONFLICT (id) DO NOTHING;

-- Insert sample workflows
INSERT INTO workflows (id, name, description, agent_id, definition, is_active)
VALUES
    ('550e8400-e29b-41d4-a716-446655440201', 'Daily Email Summary', 'Summarize daily emails', '550e8400-e29b-41d4-a716-446655440101', '{"steps": [{"type": "fetch_emails", "params": {}}]}', true),
    ('550e8400-e29b-41d4-a716-446655440202', 'Auto Response', 'Automatic email response', '550e8400-e29b-41d4-a716-446655440101', '{"steps": [{"type": "respond_to_emails", "params": {}}]}', true),
    ('550e8400-e29b-41d4-a716-446655440203', 'Data Export', 'Export analyzed data', '550e8400-e29b-41d4-a716-446655440102', '{"steps": [{"type": "export_data", "params": {}}]}', false)
ON CONFLICT (id) DO NOTHING;

EOF

echo "✓ Development data seeded successfully"
echo ""
echo "Sample Credentials:"
echo "  Email: admin@velora.dev"
echo "  Email: user@velora.dev"
echo "  Email: agent-demo@velora.dev"
echo "  Password: password (hashed in database)"
