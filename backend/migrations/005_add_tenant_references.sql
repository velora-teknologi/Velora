-- Migration: 005_add_tenant_references
-- Description: Add tenant references to users, agents, and workflows
-- Version: 005
-- Date: 2026-06-01

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(36),
    ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'user' NOT NULL;

ALTER TABLE agents
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(36);

ALTER TABLE workflows
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(36);

ALTER TABLE users
    ADD CONSTRAINT IF NOT EXISTS fk_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL;
ALTER TABLE agents
    ADD CONSTRAINT IF NOT EXISTS fk_agents_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL;
ALTER TABLE workflows
    ADD CONSTRAINT IF NOT EXISTS fk_workflows_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_users_tenant_id ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_agents_tenant_id ON agents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_workflows_tenant_id ON workflows(tenant_id);
