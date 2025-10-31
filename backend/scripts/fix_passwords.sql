-- ============================================
-- FIX PASSWORD HASHES
-- ============================================
-- Script untuk memperbaiki password hash yang tidak cocok
-- Jalankan script ini di pgAdmin Query Tool

-- Update Admin accounts
UPDATE public.users 
SET password = '$2a$10$obyfJlWKMNx/eDQ.p5FZpevnU4OLoc2nGcmusurRjJe37/NBvnTl6'
WHERE user_id = 'PI0824.2374'; -- BAYU password: 13032001

UPDATE public.users 
SET password = '$2a$10$emd0QZ4S6YhjnroHATGF5uCvrn09Ib76nr1o1butTG6lYCQHFBrjG'
WHERE user_id = 'PI0824.0001'; -- Jeremy password: admin123

UPDATE public.users 
SET password = '$2a$10$emd0QZ4S6YhjnroHATGF5uCvrn09Ib76nr1o1butTG6lYCQHFBrjG'
WHERE user_id = 'PI0824.0002'; -- Agape password: admin123

-- Update PPIC accounts  
UPDATE public.users 
SET password = '$2a$10$IHPFXWljbjx4UK3.wPkJP.dRbjcqjQS0FMD4RSjOVzsVV/ILPwvee'
WHERE user_id = 'PI0824.2375'; -- Amelia password: 20010313

UPDATE public.users 
SET password = '$2a$10$VKlCwWaZt4w7KiPtDkRDGOlyIDhClUPM9w2U3AjdQBxUPR.rJHKv.'
WHERE user_id IN ('PI0824.1001', 'PI0824.1002', 'PI0824.1003'); -- Siti, Budi, Rina password: ppic2024

-- Update Supervisor accounts
UPDATE public.users 
SET password = '$2a$10$ik449BK/DRm4UcExUwED6OZ7Uo6nsZ5HivKkLozJymbvnVBKATinK'
WHERE user_id IN ('PI0824.2001', 'PI0824.2002', 'PI0824.2003'); -- Dimas, Fitri, Andi password: super123

-- Update Operator accounts
UPDATE public.users 
SET password = '$2a$10$2XPFYcBFWU6ELF0ifVifxeNjm9f8OY2IqWvmb20efh042v4P1gbNW'
WHERE user_id IN ('PI0824.3001', 'PI0824.3002', 'PI0824.3003', 'PI0824.3004', 'PI0824.3005'); -- All operators password: op2024

-- Update Warehouse accounts
UPDATE public.users 
SET password = '$2a$10$qvQr1kDDp3R8xoUQeXv4p.Qcek7D5gfZgmtzu6.UqAxaM5plCELHC'
WHERE user_id IN ('PI0824.4001', 'PI0824.4002'); -- Tono, Lisa password: wh2024

-- Update Test account
UPDATE public.users 
SET password = '$2a$10$FM4RpLsmgSOkOb.fm1sypOuVTdZHEDjba5UQXEtrBNGt5Y3sxFUae'
WHERE user_id = 'PI0824.9999'; -- TestUser password: test123

-- Verify updates
SELECT 
    user_id,
    username,
    role,
    CASE 
        WHEN user_id = 'PI0824.2374' THEN '13032001'
        WHEN user_id IN ('PI0824.0001', 'PI0824.0002') THEN 'admin123'
        WHEN user_id = 'PI0824.2375' THEN '20010313'
        WHEN user_id IN ('PI0824.1001', 'PI0824.1002', 'PI0824.1003') THEN 'ppic2024'
        WHEN user_id IN ('PI0824.2001', 'PI0824.2002', 'PI0824.2003') THEN 'super123'
        WHEN user_id IN ('PI0824.3001', 'PI0824.3002', 'PI0824.3003', 'PI0824.3004', 'PI0824.3005') THEN 'op2024'
        WHEN user_id IN ('PI0824.4001', 'PI0824.4002') THEN 'wh2024'
        WHEN user_id = 'PI0824.9999' THEN 'test123'
    END as password_plaintext,
    LEFT(password, 20) || '...' as password_hash_preview,
    updated_at
FROM public.users
ORDER BY 
    CASE role
        WHEN 'Admin' THEN 1
        WHEN 'PPIC' THEN 2
        WHEN 'Supervisor' THEN 3
        WHEN 'Operator' THEN 4
        WHEN 'Warehouse' THEN 5
        ELSE 6
    END,
    id;

-- Done!
SELECT 'Password hashes updated successfully!' as status;
