-- ============================================================
-- CLEANUP OLD DATA
-- Hapus data lama atau testing
-- ============================================================

-- Hapus PPIC links
DELETE FROM ppic_links;

-- Hapus machine assignments
DELETE FROM machine_assignments;

-- Hapus PPIC schedules
DELETE FROM ppic_schedules;

-- Hapus job orders
DELETE FROM job_orders;

-- Hapus operation plan approvals
DELETE FROM operation_plan_approvals;

-- Hapus operation plans
DELETE FROM operation_plans;

-- Hapus pem operation plans
DELETE FROM pem_operation_plans;

-- Hapus g-code files
DELETE FROM g_code_files;

-- Hapus toolpather files
DELETE FROM toolpather_files;

-- Hapus token blacklist yang sudah expired
DELETE FROM token_blacklist WHERE expires_at < NOW();

-- Reset sequences (optional)
ALTER SEQUENCE ppic_schedules_id_seq RESTART WITH 1;
ALTER SEQUENCE machine_assignments_id_seq RESTART WITH 1;
ALTER SEQUENCE ppic_links_id_seq RESTART WITH 1;
ALTER SEQUENCE job_orders_id_seq RESTART WITH 1;

-- Verify cleanup
SELECT 
    'ppic_schedules' as table_name, COUNT(*) as remaining_rows FROM ppic_schedules
UNION ALL SELECT 'machine_assignments', COUNT(*) FROM machine_assignments
UNION ALL SELECT 'ppic_links', COUNT(*) FROM ppic_links
UNION ALL SELECT 'job_orders', COUNT(*) FROM job_orders;

SELECT 'Cleanup completed!' as status;
