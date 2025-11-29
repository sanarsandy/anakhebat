-- Script untuk membuat user admin pertama
-- Password: admin123 (bcrypt hash)

-- Cek apakah user sudah ada
DO $$
DECLARE
    user_exists boolean;
BEGIN
    SELECT EXISTS(SELECT 1 FROM users WHERE email = 'admin@jurnalsikecil.com') INTO user_exists;
    
    IF user_exists THEN
        -- Update existing user menjadi admin
        UPDATE users 
        SET role = 'admin',
            password_hash = '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', -- admin123
            full_name = COALESCE(full_name, 'Admin Jurnal Si Kecil'),
            auth_provider = COALESCE(auth_provider, 'email')
        WHERE email = 'admin@jurnalsikecil.com';
        
        RAISE NOTICE 'User updated to admin';
    ELSE
        -- Buat user admin baru
        INSERT INTO users (email, password_hash, full_name, role, auth_provider)
        VALUES (
            'admin@jurnalsikecil.com',
            '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', -- admin123
            'Admin Jurnal Si Kecil',
            'admin',
            'email'
        );
        
        RAISE NOTICE 'Admin user created';
    END IF;
END $$;

-- Verify admin user created
SELECT id, email, full_name, role, created_at 
FROM users 
WHERE role = 'admin';
