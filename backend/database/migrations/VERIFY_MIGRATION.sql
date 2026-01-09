-- ============================================================
-- VERIFY MIGRATION SCRIPT
-- ============================================================
-- Script untuk verifikasi bahwa migration berhasil
-- ============================================================
--
-- CARA PENGGUNAAN:
-- psql -U postgres -d ganttpro_db -f VERIFY_MIGRATION.sql
--
-- ============================================================

\echo '============================================================'
\echo 'VERIFIKASI DATABASE MIGRATION - GANTTPRO ERP'
\echo '============================================================'
\echo ''

-- Set display format
\pset border 2
\pset format wrapped

-- ============================================================
-- 1. CEK SEMUA TABEL
-- ============================================================

\echo '1. DAFTAR SEMUA TABEL'
\echo '------------------------------------------------------------'

SELECT
    table_name,
    (SELECT COUNT(*)
     FROM information_schema.columns
     WHERE columns.table_name = tables.table_name) as jumlah_kolom
FROM information_schema.tables
WHERE table_schema = 'public'
  AND table_type = 'BASE TABLE'
ORDER BY table_name;

\echo ''

-- ============================================================
-- 2. CEK JUMLAH DATA
-- ============================================================

\echo '2. JUMLAH DATA PER TABEL'
\echo '------------------------------------------------------------'

SELECT 'users' AS tabel, COUNT(*) AS jumlah_record FROM users
UNION ALL
SELECT 'token_blacklist', COUNT(*) FROM token_blacklist
UNION ALL
SELECT 'machines', COUNT(*) FROM machines
UNION ALL
SELECT 'job_orders', COUNT(*) FROM job_orders
UNION ALL
SELECT 'toolpather_files', COUNT(*) FROM toolpather_files
UNION ALL
SELECT 'g_code_files', COUNT(*) FROM g_code_files
UNION ALL
SELECT 'pem_operation_plans', COUNT(*) FROM pem_operation_plans
UNION ALL
SELECT 'operation_plans', COUNT(*) FROM operation_plans
UNION ALL
SELECT 'operation_plan_approvals', COUNT(*) FROM operation_plan_approvals
UNION ALL
SELECT 'ppic_schedules', COUNT(*) FROM ppic_schedules
UNION ALL
SELECT 'machine_assignments', COUNT(*) FROM machine_assignments
UNION ALL
SELECT 'ppic_links', COUNT(*) FROM ppic_links
ORDER BY tabel;

\echo ''

-- ============================================================
-- 3. CEK DEFAULT USERS
-- ============================================================

\echo '3. DEFAULT USERS'
\echo '------------------------------------------------------------'

SELECT
    id,
    username,
    email,
    role,
    full_name,
    CASE
        WHEN deleted_at IS NULL THEN 'Active'
        ELSE 'Deleted'
    END as status
FROM users
ORDER BY role, username;

\echo ''

-- ============================================================
-- 4. CEK MACHINES
-- ============================================================

\echo '4. DAFTAR MACHINES'
\echo '------------------------------------------------------------'

SELECT
    id,
    machine_code,
    machine_name,
    status,
    CASE
        WHEN deleted_at IS NULL THEN 'Active'
        ELSE 'Deleted'
    END as record_status
FROM machines
ORDER BY machine_code;

\echo ''

-- ============================================================
-- 5. CEK INDEXES
-- ============================================================

\echo '5. DATABASE INDEXES'
\echo '------------------------------------------------------------'

SELECT
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE schemaname = 'public'
ORDER BY tablename, indexname;

\echo ''

-- ============================================================
-- 6. CEK FOREIGN KEYS
-- ============================================================

\echo '6. FOREIGN KEY CONSTRAINTS'
\echo '------------------------------------------------------------'

SELECT
    tc.table_name as tabel,
    kcu.column_name as kolom,
    ccu.table_name AS foreign_tabel,
    ccu.column_name AS foreign_kolom,
    tc.constraint_name as constraint_name
FROM information_schema.table_constraints AS tc
JOIN information_schema.key_column_usage AS kcu
    ON tc.constraint_name = kcu.constraint_name
    AND tc.table_schema = kcu.table_schema
JOIN information_schema.constraint_column_usage AS ccu
    ON ccu.constraint_name = tc.constraint_name
    AND ccu.table_schema = tc.table_schema
WHERE tc.constraint_type = 'FOREIGN KEY'
  AND tc.table_schema = 'public'
ORDER BY tc.table_name, kcu.column_name;

\echo ''

-- ============================================================
-- 7. CEK SEQUENCES
-- ============================================================

\echo '7. DATABASE SEQUENCES'
\echo '------------------------------------------------------------'

SELECT
    sequence_name,
    data_type,
    start_value,
    minimum_value,
    maximum_value,
    increment
FROM information_schema.sequences
WHERE sequence_schema = 'public'
ORDER BY sequence_name;

\echo ''

-- ============================================================
-- 8. CEK PPIC SCHEDULES (jika ada)
-- ============================================================

\echo '8. PPIC SCHEDULES (Sample Data)'
\echo '------------------------------------------------------------'

SELECT
    id,
    njo,
    part_name,
    start_date,
    finish_date,
    priority,
    material_status,
    status,
    progress
FROM ppic_schedules
WHERE deleted_at IS NULL
ORDER BY start_date DESC
LIMIT 10;

\echo ''

-- ============================================================
-- 9. STORAGE SIZE
-- ============================================================

\echo '9. DATABASE STORAGE SIZE'
\echo '------------------------------------------------------------'

SELECT
    pg_database.datname as database_name,
    pg_size_pretty(pg_database_size(pg_database.datname)) AS size
FROM pg_database
WHERE datname = current_database();

\echo ''

SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size,
    pg_size_pretty(pg_relation_size(schemaname||'.'||tablename)) AS table_size,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename) - pg_relation_size(schemaname||'.'||tablename)) AS index_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

\echo ''

-- ============================================================
-- 10. SUMMARY VERIFICATION
-- ============================================================

\echo '10. MIGRATION SUMMARY'
\echo '============================================================'

DO $$
DECLARE
    v_table_count INTEGER;
    v_user_count INTEGER;
    v_machine_count INTEGER;
    v_index_count INTEGER;
    v_fk_count INTEGER;
    v_is_valid BOOLEAN := TRUE;
    v_errors TEXT := '';
BEGIN
    -- Count tables
    SELECT COUNT(*) INTO v_table_count
    FROM information_schema.tables
    WHERE table_schema = 'public' AND table_type = 'BASE TABLE';

    -- Count users
    SELECT COUNT(*) INTO v_user_count FROM users;

    -- Count machines
    SELECT COUNT(*) INTO v_machine_count FROM machines;

    -- Count indexes
    SELECT COUNT(*) INTO v_index_count
    FROM pg_indexes WHERE schemaname = 'public';

    -- Count foreign keys
    SELECT COUNT(*) INTO v_fk_count
    FROM information_schema.table_constraints
    WHERE constraint_type = 'FOREIGN KEY' AND table_schema = 'public';

    -- Validation
    IF v_table_count < 12 THEN
        v_is_valid := FALSE;
        v_errors := v_errors || '- Tabel kurang dari 12 (expected: 12, found: ' || v_table_count || ')' || E'\n';
    END IF;

    IF v_user_count < 4 THEN
        v_is_valid := FALSE;
        v_errors := v_errors || '- User kurang dari 4 (expected: 4+, found: ' || v_user_count || ')' || E'\n';
    END IF;

    IF v_machine_count < 8 THEN
        v_is_valid := FALSE;
        v_errors := v_errors || '- Machine kurang dari 8 (expected: 8+, found: ' || v_machine_count || ')' || E'\n';
    END IF;

    -- Display results
    RAISE NOTICE '============================================================';

    IF v_is_valid THEN
        RAISE NOTICE '✓ MIGRATION BERHASIL!';
        RAISE NOTICE '============================================================';
        RAISE NOTICE 'Total Tables      : % (✓ OK)', v_table_count;
        RAISE NOTICE 'Total Users       : % (✓ OK)', v_user_count;
        RAISE NOTICE 'Total Machines    : % (✓ OK)', v_machine_count;
        RAISE NOTICE 'Total Indexes     : % (✓ OK)', v_index_count;
        RAISE NOTICE 'Total Foreign Keys: % (✓ OK)', v_fk_count;
        RAISE NOTICE '============================================================';
        RAISE NOTICE '';
        RAISE NOTICE 'Database siap digunakan!';
        RAISE NOTICE 'Default credentials:';
        RAISE NOTICE '  Username: admin | Password: admin123';
        RAISE NOTICE '  Username: ppic  | Password: ppic123';
        RAISE NOTICE '  Username: pem   | Password: pem123';
        RAISE NOTICE '  Username: operator | Password: operator123';
        RAISE NOTICE '';
        RAISE NOTICE '⚠️  PENTING: Ganti password default setelah login!';
    ELSE
        RAISE NOTICE '✗ MIGRATION INCOMPLETE!';
        RAISE NOTICE '============================================================';
        RAISE NOTICE 'Total Tables      : % (Expected: 12)', v_table_count;
        RAISE NOTICE 'Total Users       : % (Expected: 4+)', v_user_count;
        RAISE NOTICE 'Total Machines    : % (Expected: 8+)', v_machine_count;
        RAISE NOTICE '============================================================';
        RAISE NOTICE 'Errors:';
        RAISE NOTICE '%', v_errors;
        RAISE NOTICE '============================================================';
        RAISE NOTICE 'Silakan jalankan ulang MIGRATION_COMPLETE.sql';
    END IF;

    RAISE NOTICE '============================================================';
END $$;

\echo ''
\echo 'Verifikasi selesai!'
\echo ''
\echo 'Untuk test koneksi dari backend:'
\echo '  cd backend'
\echo '  go run main.go'
\echo ''
\echo 'Untuk test API:'
\echo '  curl http://localhost:8080/health'
\echo ''
