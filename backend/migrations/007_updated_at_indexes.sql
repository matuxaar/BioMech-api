ALTER TABLE devices
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE emg_sessions
  ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_devices_updated_at ON devices(updated_at);
CREATE INDEX IF NOT EXISTS idx_emg_sessions_updated_at ON emg_sessions(updated_at);
