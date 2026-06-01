-- Migration: 007_create_invitations_table
-- Description: Create invitations table for tenant invitation workflow
-- Version: 007
-- Date: 2026-06-01

CREATE TABLE IF NOT EXISTS invitations (
    id VARCHAR(36) PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id VARCHAR(36) NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    role VARCHAR(50) DEFAULT 'member' NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' NOT NULL,
    invited_by VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    accepted_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_invitations_tenant_id ON invitations(tenant_id);
CREATE INDEX IF NOT EXISTS idx_invitations_email ON invitations(email);
CREATE INDEX IF NOT EXISTS idx_invitations_status ON invitations(status);
CREATE INDEX IF NOT EXISTS idx_invitations_deleted_at ON invitations(deleted_at);

CREATE OR REPLACE FUNCTION update_invitations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS invitations_update_updated_at ON invitations;
CREATE TRIGGER invitations_update_updated_at
    BEFORE UPDATE ON invitations
    FOR EACH ROW
    EXECUTE FUNCTION update_invitations_updated_at();
