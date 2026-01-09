-- ============================================================
-- QUICK CHECK - Verifikasi Cepat Database
-- ============================================================

-- Check current database
SELECT current_database() as database_name;

-- Count all tables
SELECT COUNT(*) as total_tables 
FROM information_schema.tables 
WHERE table_schema = 'public';

-- Count data in each table
SELECT 
    'users' as table_name, COUNT(*) as row_count FROM users
UNION ALL SELECT 'machines', COUNT(*) FROM machines
UNION ALL SELECT 'job_orders', COUNT(*) FROM job_orders
UNION ALL SELECT 'ppic_schedules', COUNT(*) FROM ppic_schedules
UNION ALL SELECT 'machine_assignments', COUNT(*) FROM machine_assignments
UNION ALL SELECT 'ppic_links', COUNT(*) FROM ppic_links
UNION ALL SELECT 'operation_plans', COUNT(*) FROM operation_plans
UNION ALL SELECT 'pem_operation_plans', COUNT(*) FROM pem_operation_plans
UNION ALL SELECT 'g_code_files', COUNT(*) FROM g_code_files
UNION ALL SELECT 'toolpather_files', COUNT(*) FROM toolpather_files
UNION ALL SELECT 'token_blacklist', COUNT(*) FROM token_blacklist;

-- Check users
SELECT id, username, email, role, created_at 
FROM users 
ORDER BY id;

-- Check machines
SELECT id, machine_code, machine_name, status 
FROM machines 
ORDER BY id;

-- Check PPIC schedules with assignments
SELECT 
    ps.njo,
    ps.part_name,
    ps.priority,
    ps.status,
    ps.progress,
    COUNT(ma.id) as machine_count
FROM ppic_schedules ps
LEFT JOIN machine_assignments ma ON ps.id = ma.schedule_id
GROUP BY ps.id, ps.njo, ps.part_name, ps.priority, ps.status, ps.progress
ORDER BY ps.id;

-- Check machine_assignments table structure
SELECT 
    column_name,
    data_type,
    is_nullable
FROM information_schema.columns
WHERE table_name = 'machine_assignments'
ORDER BY ordinal_position;
