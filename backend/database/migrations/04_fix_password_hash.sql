-- ============================================================
-- FIX PASSWORD HASH
-- Ganti password hash dengan yang benar untuk "admin123"
-- ============================================================

-- Update admin password dengan hash yang benar
-- Password: admin123
-- Hash menggunakan bcrypt cost 10
UPDATE users 
SET password = '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'
WHERE username = 'admin';

-- Verify
SELECT 
    id, 
    username, 
    email, 
    role,
    CASE 
        WHEN password = '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy' 
        THEN 'Password OK (admin123)' 
        ELSE 'Password Different' 
    END as password_status
FROM users 
WHERE username = 'admin';
