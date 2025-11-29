-- Migration: Add Admin RBAC Support
-- This migration adds tables and constraints for admin role-based access control

-- Ensure role can be 'admin' in users table
-- Drop constraint if exists first, then add it
DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'check_role') THEN
        ALTER TABLE users DROP CONSTRAINT check_role;
    END IF;
END $$;

ALTER TABLE users 
ADD CONSTRAINT check_role CHECK (role IN ('parent', 'admin'));

-- Create Audit Logs Table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL, -- 'create', 'update', 'delete', 'view'
    resource_type VARCHAR(50) NOT NULL, -- 'user', 'child', 'milestone', etc.
    resource_id UUID,
    before_data JSONB,
    after_data JSONB,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for audit logs
CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);

-- Create System Settings Table
CREATE TABLE IF NOT EXISTS system_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT,
    type VARCHAR(50) DEFAULT 'string', -- 'string', 'number', 'boolean', 'json'
    category VARCHAR(50) DEFAULT 'general', -- 'general', 'email', 'notifications', 'features', etc.
    description TEXT,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for system settings
CREATE INDEX IF NOT EXISTS idx_system_settings_key ON system_settings(key);
CREATE INDEX IF NOT EXISTS idx_system_settings_category ON system_settings(category);

-- Insert default system settings
INSERT INTO system_settings (key, value, type, category, description) VALUES
('app_name', 'Jurnal Si Kecil', 'string', 'general', 'Application name'),
('maintenance_mode', 'false', 'boolean', 'general', 'Enable/disable maintenance mode'),
('max_otp_attempts', '3', 'number', 'security', 'Maximum OTP verification attempts'),
('otp_expiry_minutes', '10', 'number', 'security', 'OTP code expiry time in minutes'),
('enable_whatsapp_notifications', 'true', 'boolean', 'notifications', 'Enable WhatsApp notifications'),
('enable_email_notifications', 'true', 'boolean', 'notifications', 'Enable email notifications')
ON CONFLICT (key) DO NOTHING;

-- Add comments
COMMENT ON TABLE audit_logs IS 'Audit trail for all admin actions and system changes';
COMMENT ON TABLE system_settings IS 'System configuration settings managed by admin';
COMMENT ON COLUMN audit_logs.before_data IS 'JSON data before the action (for update/delete)';
COMMENT ON COLUMN audit_logs.after_data IS 'JSON data after the action (for create/update)';

