-- ============================================================
-- RESET DATABASE - DROP ALL TABLES
-- HATI-HATI: Ini akan menghapus semua data!
-- ============================================================

-- Drop tables in correct order (reverse of dependencies)
DROP TABLE IF EXISTS ppic_links CASCADE;
DROP TABLE IF EXISTS machine_assignments CASCADE;
DROP TABLE IF EXISTS ppic_schedules CASCADE;
DROP TABLE IF EXISTS operation_plan_approvals CASCADE;
DROP TABLE IF EXISTS operation_plans CASCADE;
DROP TABLE IF EXISTS pem_operation_plans CASCADE;
DROP TABLE IF EXISTS g_code_files CASCADE;
DROP TABLE IF EXISTS toolpather_files CASCADE;
DROP TABLE IF EXISTS job_orders CASCADE;
DROP TABLE IF EXISTS machines CASCADE;
DROP TABLE IF EXISTS token_blacklist CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Drop sequences
DROP SEQUENCE IF EXISTS ppic_links_id_seq CASCADE;

-- Verify all tables dropped
SELECT 'All tables dropped successfully!' AS status;
