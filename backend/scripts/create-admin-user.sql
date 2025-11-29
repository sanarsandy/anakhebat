-- Script untuk membuat user admin pertama
-- Usage: 
--   Via Docker: docker exec -i tukem-db psql -U tukem_user -d tukem_db < create-admin-user.sql
--   Via psql: psql -h localhost -U tukem_user -d tukem_db -f create-admin-user.sql

-- Option 1: Update existing user menjadi admin
-- Ganti email dengan email user yang ingin dijadikan admin
-- UPDATE users SET role = 'admin' WHERE email = 'admin@example.com';

-- Option 2: Buat user admin baru
-- Ganti email, password, dan full_name sesuai kebutuhan
INSERT INTO users (email, password_hash, full_name, role, auth_provider)
VALUES (
    'admin@jurnalsikecil.com',
    '$2a$10$rK8X8X8X8X8X8X8X8X8Xe8X8X8X8X8X8X8X8X8X8X8X8X8X8X8X8X', -- Password: admin123 (harus di-hash dengan bcrypt)
    'Admin Jurnal Si Kecil',
    'admin',
    'email'
)
ON CONFLICT (email) DO UPDATE SET role = 'admin';

-- Note: Password hash di atas adalah placeholder
-- Untuk generate password hash yang benar, gunakan bcrypt dengan cost 10
-- Atau gunakan script Go di bawah untuk generate hash


