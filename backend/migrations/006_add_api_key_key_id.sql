-- 006_add_api_key_key_id.sql
ALTER TABLE api_keys
    ADD COLUMN IF NOT EXISTS key_id UUID NOT NULL DEFAULT uuid_generate_v4();

CREATE UNIQUE INDEX IF NOT EXISTS idx_api_keys_key_id ON api_keys(key_id);
