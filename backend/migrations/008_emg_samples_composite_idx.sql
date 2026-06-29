CREATE INDEX IF NOT EXISTS idx_emg_samples_session_timestamp
  ON emg_samples(session_id, timestamp);

-- Update devices.updated_at on device modification
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'set_devices_updated_at') THEN
    CREATE TRIGGER set_devices_updated_at
      BEFORE UPDATE ON devices
      FOR EACH ROW
      EXECUTE FUNCTION update_updated_at_column();
  END IF;
END $$;
