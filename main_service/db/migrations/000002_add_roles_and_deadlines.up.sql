ALTER TABLE users ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user';

ALTER TABLE tasks ADD COLUMN IF NOT EXISTS deadline TIMESTAMPTZ;

-- Create index for role lookups
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
