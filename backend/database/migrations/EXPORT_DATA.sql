-- ============================================================
-- EXPORT DATA SCRIPT
-- ============================================================
-- Script ini untuk export data dari database lama
-- Jalankan di PC lama sebelum migrasi
-- ============================================================
--
-- CARA PENGGUNAAN:
-- 1. Di PC lama, jalankan: psql -U postgres -d ganttpro_db -f EXPORT_DATA.sql > data_export.sql
-- 2. Copy file data_export.sql ke PC baru
-- 3. Di PC baru, import dengan: psql -U postgres -d ganttpro_db -f data_export.sql
--
-- ============================================================

\echo '============================================================'
\echo 'EXPORT DATA - GANTTPRO ERP'
\echo '============================================================'
\echo ''

-- Set output format
\pset tuples_only on
\pset format unaligned

-- ============================================================
-- EXPORT USERS (tanpa password hash untuk keamanan)
-- ============================================================

\echo '-- ============================================================'
\echo '-- USERS DATA'
\echo '-- ============================================================'
\echo ''

SELECT 'INSERT INTO users (email, username, password, full_name, role, created_at, updated_at) VALUES '
    || string_agg(
        format('(%L, %L, %L, %L, %L, %L, %L)',
            email,
            username,
            password,
            full_name,
            role,
            created_at,
            updated_at
        ),
        ', '
    )
    || ' ON CONFLICT (email) DO UPDATE SET '
    || 'username = EXCLUDED.username, '
    || 'full_name = EXCLUDED.full_name, '
    || 'role = EXCLUDED.role;'
FROM users
WHERE deleted_at IS NULL;

\echo ''

-- ============================================================
-- EXPORT MACHINES
-- ============================================================

\echo '-- ============================================================'
\echo '-- MACHINES DATA'
\echo '-- ============================================================'
\echo ''

SELECT 'INSERT INTO machines (machine_code, machine_name, description, status, created_at, updated_at) VALUES '
    || string_agg(
        format('(%L, %L, %L, %L, %L, %L)',
            machine_code,
            machine_name,
            description,
            status,
            created_at,
            updated_at
        ),
        ', '
    )
    || ' ON CONFLICT (machine_code) DO UPDATE SET '
    || 'machine_name = EXCLUDED.machine_name, '
    || 'description = EXCLUDED.description, '
    || 'status = EXCLUDED.status;'
FROM machines
WHERE deleted_at IS NULL;

\echo ''

-- ============================================================
-- EXPORT JOB ORDERS
-- ============================================================

\echo '-- ============================================================'
\echo '-- JOB ORDERS DATA'
\echo '-- ============================================================'
\echo ''

SELECT 'INSERT INTO job_orders (order_number, customer_name, product_name, quantity, status, start_date, due_date, notes, created_at, updated_at) VALUES '
    || string_agg(
        format('(%L, %L, %L, %s, %L, %L, %L, %L, %L, %L)',
            order_number,
            customer_name,
            product_name,
            quantity,
            status,
            start_date,
            due_date,
            notes,
            created_at,
            updated_at
        ),
        ', '
    )
    || ' ON CONFLICT (order_number) DO NOTHING;'
FROM job_orders
WHERE deleted_at IS NULL
    AND EXISTS (SELECT 1 FROM job_orders WHERE deleted_at IS NULL);

\echo ''

-- ============================================================
-- EXPORT PPIC SCHEDULES
-- ============================================================

\echo '-- ============================================================'
\echo '-- PPIC SCHEDULES DATA'
\echo '-- ============================================================'
\echo ''

SELECT 'INSERT INTO ppic_schedules (njo, part_name, start_date, finish_date, priority, priority_alpha, material_status, ppic_notes, status, progress, created_at, updated_at) VALUES '
    || string_agg(
        format('(%L, %L, %L, %L, %L, %L, %L, %L, %L, %s, %L, %L)',
            njo,
            part_name,
            start_date,
            finish_date,
            priority,
            priority_alpha,
            material_status,
            ppic_notes,
            status,
            progress,
            created_at,
            updated_at
        ),
        ', '
    )
    || ' ON CONFLICT (njo) DO UPDATE SET '
    || 'part_name = EXCLUDED.part_name, '
    || 'status = EXCLUDED.status, '
    || 'progress = EXCLUDED.progress;'
FROM ppic_schedules
WHERE deleted_at IS NULL
    AND EXISTS (SELECT 1 FROM ppic_schedules WHERE deleted_at IS NULL);

\echo ''

-- ============================================================
-- EXPORT MACHINE ASSIGNMENTS
-- ============================================================

\echo '-- ============================================================'
\echo '-- MACHINE ASSIGNMENTS DATA'
\echo '-- ============================================================'
\echo ''

-- Note: Requires schedule_id and machine_id mapping
SELECT 'INSERT INTO machine_assignments (schedule_id, machine_id, target_hours, scheduled_start, scheduled_end, actual_start, actual_end, status, sequence, created_at, updated_at) '
    || 'SELECT '
    || '(SELECT id FROM ppic_schedules WHERE njo = ' || quote_literal(ps.njo) || '), '
    || '(SELECT id FROM machines WHERE machine_code = ' || quote_literal(m.machine_code) || '), '
    || string_agg(
        format('%s, %L, %L, %L, %L, %L, %s, %L, %L',
            ma.target_hours,
            ma.scheduled_start,
            ma.scheduled_end,
            ma.actual_start,
            ma.actual_end,
            ma.status,
            ma.sequence,
            ma.created_at,
            ma.updated_at
        ),
        ' UNION ALL SELECT '
    )
    || ';'
FROM machine_assignments ma
JOIN ppic_schedules ps ON ma.schedule_id = ps.id
JOIN machines m ON ma.machine_id = m.id
WHERE ps.deleted_at IS NULL
    AND m.deleted_at IS NULL
    AND EXISTS (SELECT 1 FROM machine_assignments)
GROUP BY ps.njo, m.machine_code;

\echo ''

-- ============================================================
-- EXPORT PPIC LINKS
-- ============================================================

\echo '-- ============================================================'
\echo '-- PPIC LINKS DATA'
\echo '-- ============================================================'
\echo ''

SELECT 'INSERT INTO ppic_links (source_schedule_id, target_schedule_id, link_type, created_at, updated_at) '
    || 'SELECT '
    || '(SELECT id FROM ppic_schedules WHERE njo = ' || quote_literal(source.njo) || '), '
    || '(SELECT id FROM ppic_schedules WHERE njo = ' || quote_literal(target.njo) || '), '
    || string_agg(
        format('%L, %L, %L',
            pl.link_type,
            pl.created_at,
            pl.updated_at
        ),
        ' UNION ALL SELECT '
    )
    || ';'
FROM ppic_links pl
JOIN ppic_schedules source ON pl.source_schedule_id = source.id
JOIN ppic_schedules target ON pl.target_schedule_id = target.id
WHERE EXISTS (SELECT 1 FROM ppic_links)
GROUP BY source.njo, target.njo;

\echo ''

-- ============================================================
-- EXPORT PEM OPERATION PLANS
-- ============================================================

\echo '-- ============================================================'
\echo '-- PEM OPERATION PLANS DATA'
\echo '-- ============================================================'
\echo ''

SELECT 'INSERT INTO pem_operation_plans (operation_plan_code, part_name, part_number, machine_id, material, quantity, finish_dimension, description, created_at, updated_at) '
    || 'VALUES '
    || string_agg(
        format('(%L, %L, %L, (SELECT id FROM machines WHERE machine_code = %L), %L, %s, %L, %L, %L, %L)',
            operation_plan_code,
            part_name,
            part_number,
            (SELECT machine_code FROM machines WHERE id = pem_operation_plans.machine_id),
            material,
            quantity,
            finish_dimension,
            description,
            created_at,
            updated_at
        ),
        ', '
    )
    || ' ON CONFLICT (operation_plan_code) DO NOTHING;'
FROM pem_operation_plans
WHERE deleted_at IS NULL
    AND EXISTS (SELECT 1 FROM pem_operation_plans WHERE deleted_at IS NULL);

\echo ''

-- ============================================================
-- EXPORT SUMMARY
-- ============================================================

\echo ''
\echo '-- ============================================================'
\echo '-- EXPORT SUMMARY'
\echo '-- ============================================================'

\pset tuples_only off
\pset format aligned

SELECT 'users' AS table_name, COUNT(*) AS record_count FROM users WHERE deleted_at IS NULL
UNION ALL
SELECT 'machines', COUNT(*) FROM machines WHERE deleted_at IS NULL
UNION ALL
SELECT 'job_orders', COUNT(*) FROM job_orders WHERE deleted_at IS NULL
UNION ALL
SELECT 'ppic_schedules', COUNT(*) FROM ppic_schedules WHERE deleted_at IS NULL
UNION ALL
SELECT 'machine_assignments', COUNT(*) FROM machine_assignments
UNION ALL
SELECT 'ppic_links', COUNT(*) FROM ppic_links
UNION ALL
SELECT 'pem_operation_plans', COUNT(*) FROM pem_operation_plans WHERE deleted_at IS NULL;

\echo ''
\echo '-- Export completed!'
\echo '-- Save this output to a file and run it on the new database'
