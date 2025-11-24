-- Migration: Add Google OAuth support to users table
-- This allows users to login with Google account

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS google_id VARCHAR(255) UNIQUE,
ADD COLUMN IF NOT EXISTS auth_provider VARCHAR(50) DEFAULT 'email', -- 'email' or 'google'
ALTER COLUMN password_hash DROP NOT NULL; -- Allow NULL for Google users

-- Add index for faster Google ID lookups
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id) WHERE google_id IS NOT NULL;

-- Add comment
COMMENT ON COLUMN users.google_id IS 'Google user ID for OAuth authentication';
COMMENT ON COLUMN users.auth_provider IS 'Authentication provider: email or google';

