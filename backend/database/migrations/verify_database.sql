-- ============================================================
-- DATABASE VERIFICATION SCRIPT
-- Run this to check if database is properly set up
-- ============================================================

\echo ''
\echo '========================================='
\echo '  Database Verification'
\echo '========================================='
\echo ''

-- Check database exists
\echo 'Current Database:'
SELECT current_database();
\echo ''

-- List all tables
\echo 'Tables in database:'
SELECT 
    table_name,
    (SELECT COUNT(*) FROM information_schema.columns WHERE table_name = t.table_name) as column_count
FROM information_schema.tables t
WHERE table_schema = 'public' 
  AND table_type = 'BASE TABLE'
ORDER BY table_name;
\echo ''

-- Check critical tables exist
\echo 'Critical Tables Check:'
SELECT 
    CASE 
        WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'users') 
        THEN '✓ users' 
        ELSE '✗ users MISSING!' 
    END as users_table,
    CASE 
        WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'machines') 
        THEN '✓ machines' 
        ELSE '✗ machines MISSING!' 
    END as machines_table,
    CASE 
        WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'ppic_schedules') 
        THEN '✓ ppic_schedules' 
        ELSE '✗ ppic_schedules MISSING!' 
    END as ppic_schedules_table,
    CASE 
        WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'machine_assignments') 
        THEN '✓ machine_assignments' 
        ELSE '✗ machine_assignments MISSING!' 
    END as machine_assignments_table,
    CASE 
        WHEN EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'ppic_links') 
        THEN '✓ ppic_links' 
        ELSE '✗ ppic_links MISSING!' 
    END as ppic_links_table;
\echo ''

-- Check machine_assignments structure (the problematic table)
\echo 'machine_assignments Table Structure:'
SELECT 
    column_name,
    data_type,
    character_maximum_length,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_name = 'machine_assignments'
ORDER BY ordinal_position;
\echo ''

-- Verify critical columns exist in machine_assignments
\echo 'machine_assignments Critical Columns:'
SELECT 
    CASE 
        WHEN EXISTS (
            SELECT 1 FROM information_schema.columns 
            WHERE table_name = 'machine_assignments' AND column_name = 'scheduled_start'
        ) 
        THEN '✓ scheduled_start column exists' 
        ELSE '✗ scheduled_start column MISSING!' 
    END as scheduled_start_check,
    CASE 
        WHEN EXISTS (
            SELECT 1 FROM information_schema.columns 
            WHERE table_name = 'machine_assignments' AND column_name = 'scheduled_end'
        ) 
        THEN '✓ scheduled_end column exists' 
        ELSE '✗ scheduled_end column MISSING!' 
    END as scheduled_end_check,
    CASE 
        WHEN EXISTS (
            SELECT 1 FROM information_schema.columns 
            WHERE table_name = 'machine_assignments' AND column_name = 'actual_start'
        ) 
        THEN '✓ actual_start column exists' 
        ELSE '✗ actual_start column MISSING!' 
    END as actual_start_check,
    CASE 
        WHEN EXISTS (
            SELECT 1 FROM information_schema.columns 
            WHERE table_name = 'machine_assignments' AND column_name = 'actual_end'
        ) 
        THEN '✓ actual_end column exists' 
        ELSE '✗ actual_end column MISSING!' 
    END as actual_end_check;
\echo ''

-- Check ppic_schedules structure
\echo 'ppic_schedules Table Structure:'
SELECT 
    column_name,
    data_type,
    is_nullable
FROM information_schema.columns
WHERE table_name = 'ppic_schedules'
ORDER BY ordinal_position;
\echo ''

-- Check indexes
\echo 'Indexes on machine_assignments:'
SELECT 
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'machine_assignments';
\echo ''

-- Check foreign keys
\echo 'Foreign Keys on machine_assignments:'
SELECT
    tc.constraint_name,
    tc.table_name,
    kcu.column_name,
    ccu.table_name AS foreign_table_name,
    ccu.column_name AS foreign_column_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
WHERE tc.constraint_type = 'FOREIGN KEY'
    AND tc.table_name = 'machine_assignments';
\echo ''

-- Check data
\echo 'Data Count:'
SELECT 
    (SELECT COUNT(*) FROM users) as users_count,
    (SELECT COUNT(*) FROM machines) as machines_count,
    (SELECT COUNT(*) FROM ppic_schedules) as schedules_count,
    (SELECT COUNT(*) FROM machine_assignments) as assignments_count,
    (SELECT COUNT(*) FROM ppic_links) as links_count;
\echo ''

-- Check users
\echo 'Users in system:'
SELECT id, email, username, role, created_at 
FROM users 
ORDER BY id;
\echo ''

-- Check machines
\echo 'Machines in system:'
SELECT id, machine_code, machine_name, status 
FROM machines 
ORDER BY id;
\echo ''

-- Test query similar to what the application uses
\echo 'Test Query (similar to application):'
SELECT 
    ma.id, 
    ma.schedule_id, 
    ma.machine_id, 
    m.machine_name, 
    m.machine_code,
    ma.sequence, 
    ma.target_hours, 
    ma.scheduled_start, 
    ma.scheduled_end,
    ma.actual_start, 
    ma.actual_end, 
    ma.status
FROM machine_assignments ma
JOIN machines m ON ma.machine_id = m.id
LIMIT 5;
\echo ''

\echo '========================================='
\echo '  Verification Complete'
\echo '========================================='
\echo ''
\echo 'If all checks show ✓, database is properly set up!'
\echo 'If any check shows ✗, please run the migration script again.'
\echo ''
