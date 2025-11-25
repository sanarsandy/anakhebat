-- Migration: Add phone-based OTP authentication
-- This allows users to login with WhatsApp number + OTP

-- Add phone_number column to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS phone_number VARCHAR(20),
ADD COLUMN IF NOT EXISTS phone_verified BOOLEAN DEFAULT false,
ADD COLUMN IF NOT EXISTS phone_verified_at TIMESTAMP;

-- Make email optional (for Google OAuth or backup)
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

-- Update auth_provider to support multiple methods
-- 'phone' - phone auth only
-- 'google' - Google OAuth only  
-- 'phone_google' - both methods
-- 'email' - deprecated (existing users)

-- Create unique index for phone_number
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone_number ON users(phone_number) WHERE phone_number IS NOT NULL;

-- Create OTP codes table
CREATE TABLE IF NOT EXISTS otp_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) NOT NULL,
    otp_code VARCHAR(6) NOT NULL,
    purpose VARCHAR(50) NOT NULL DEFAULT 'login', -- 'login', 'register', 'verify_phone'
    
    -- Expiration
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    is_used BOOLEAN DEFAULT false,
    
    -- Rate limiting
    attempt_count INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    
    -- Metadata
    ip_address VARCHAR(45),
    user_agent TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for OTP codes
CREATE INDEX IF NOT EXISTS idx_otp_codes_phone ON otp_codes(phone_number);
CREATE INDEX IF NOT EXISTS idx_otp_codes_code ON otp_codes(otp_code);
CREATE INDEX IF NOT EXISTS idx_otp_codes_expires ON otp_codes(expires_at);
CREATE INDEX IF NOT EXISTS idx_otp_codes_used ON otp_codes(is_used);

-- Create OTP rate limits table
CREATE TABLE IF NOT EXISTS otp_rate_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number VARCHAR(20) NOT NULL,
    request_count INTEGER DEFAULT 1,
    window_start TIMESTAMP NOT NULL,
    window_end TIMESTAMP NOT NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(phone_number, window_start)
);

-- Create indexes for rate limits
CREATE INDEX IF NOT EXISTS idx_otp_rate_limits_phone ON otp_rate_limits(phone_number);
CREATE INDEX IF NOT EXISTS idx_otp_rate_limits_window ON otp_rate_limits(window_start, window_end);

-- Add comments
COMMENT ON COLUMN users.phone_number IS 'Phone number for OTP authentication (WhatsApp format)';
COMMENT ON COLUMN users.phone_verified IS 'Whether phone number has been verified via OTP';
COMMENT ON COLUMN users.phone_verified_at IS 'Timestamp when phone number was verified';
COMMENT ON TABLE otp_codes IS 'Stores OTP codes for phone authentication';
COMMENT ON TABLE otp_rate_limits IS 'Rate limiting for OTP requests per phone number';

