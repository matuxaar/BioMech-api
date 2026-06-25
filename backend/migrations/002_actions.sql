ALTER TABLE users ADD COLUMN IF NOT EXISTS display_name VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS photo_url VARCHAR(500) NOT NULL DEFAULT '';

CREATE TABLE IF NOT EXISTS device_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    emoji VARCHAR(10) NOT NULL DEFAULT '',
    accuracy DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(device_id, name)
);

CREATE INDEX IF NOT EXISTS idx_device_actions_device_id ON device_actions(device_id);
