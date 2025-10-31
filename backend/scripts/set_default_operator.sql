-- ============================================
-- SET DEFAULT VALUE UNTUK KOLOM OPERATOR
-- ============================================
-- Script untuk mengatur default value pada kolom operator
-- agar tidak null saat user baru signup

-- 1. Update existing NULL values menjadi '-'
UPDATE public.users 
SET operator = '-' 
WHERE operator IS NULL;

-- 2. Set default value untuk kolom operator di schema
ALTER TABLE public.users 
ALTER COLUMN operator SET DEFAULT '-';

-- 3. Verify changes
SELECT 
    id,
    username,
    user_id,
    role,
    operator,
    is_active,
    created_at
FROM public.users
ORDER BY id;

-- ============================================
-- SELESAI!
-- ============================================
-- Sekarang semua user dengan operator NULL akan menjadi '-'
-- Dan user baru yang signup akan otomatis dapat operator = '-'
-- ============================================
