ALTER TABLE users ADD COLUMN IF NOT EXISTS nickname VARCHAR(50) UNIQUE;
CREATE INDEX IF NOT EXISTS idx_users_nickname ON users(nickname);
