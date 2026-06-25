CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id VARCHAR(128) PRIMARY KEY,
    email VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(128) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('prosthetic', 'sensor')),
    name VARCHAR(255) NOT NULL,
    hw_version VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_devices_user_id ON devices(user_id);

CREATE TABLE emg_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(128) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    label VARCHAR(255),
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_emg_sessions_user_id ON emg_sessions(user_id);
CREATE INDEX idx_emg_sessions_device_id ON emg_sessions(device_id);

CREATE TABLE emg_samples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES emg_sessions(id) ON DELETE CASCADE,
    timestamp TIMESTAMPTZ NOT NULL,
    channel_1 DOUBLE PRECISION NOT NULL DEFAULT 0,
    channel_2 DOUBLE PRECISION NOT NULL DEFAULT 0,
    channel_3 DOUBLE PRECISION NOT NULL DEFAULT 0,
    channel_4 DOUBLE PRECISION NOT NULL DEFAULT 0,
    channel_5 DOUBLE PRECISION NOT NULL DEFAULT 0,
    channel_6 DOUBLE PRECISION NOT NULL DEFAULT 0,
    channel_7 DOUBLE PRECISION NOT NULL DEFAULT 0,
    channel_8 DOUBLE PRECISION NOT NULL DEFAULT 0,
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX idx_emg_samples_session_id ON emg_samples(session_id);
CREATE INDEX idx_emg_samples_timestamp ON emg_samples(timestamp);

CREATE TABLE training_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(128) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_ids UUID[] NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'running', 'completed', 'failed')),
    model_path VARCHAR(500),
    accuracy DOUBLE PRECISION DEFAULT 0,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_training_jobs_user_id ON training_jobs(user_id);
CREATE INDEX idx_training_jobs_status ON training_jobs(status);
