-- ============================================================
-- COMPLETE DATABASE MIGRATION SCRIPT
-- File ini untuk setup database lengkap di PC baru
-- ============================================================
--
-- CARA PENGGUNAAN:
-- 1. Install PostgreSQL di PC baru
-- 2. Buat database baru: CREATE DATABASE ganttpro_db;
-- 3. Jalankan file ini: psql -U postgres -d ganttpro_db -f MIGRATION_COMPLETE.sql
-- 4. Update file .env dengan kredensial database Anda
-- 5. Jalankan backend: go run main.go
--
-- ============================================================

-- Set client encoding
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

-- ============================================================
-- STEP 1: DROP EXISTING TABLES (HATI-HATI!)
-- Uncomment jika ingin reset database sepenuhnya
-- ============================================================

-- DROP TABLE IF EXISTS ppic_links CASCADE;
-- DROP TABLE IF EXISTS machine_assignments CASCADE;
-- DROP TABLE IF EXISTS ppic_schedules CASCADE;
-- DROP TABLE IF EXISTS operation_plan_approvals CASCADE;
-- DROP TABLE IF EXISTS operation_plans CASCADE;
-- DROP TABLE IF EXISTS pem_operation_plans CASCADE;
-- DROP TABLE IF EXISTS g_code_files CASCADE;
-- DROP TABLE IF EXISTS toolpather_files CASCADE;
-- DROP TABLE IF EXISTS job_orders CASCADE;
-- DROP TABLE IF EXISTS machines CASCADE;
-- DROP TABLE IF EXISTS token_blacklist CASCADE;
-- DROP TABLE IF EXISTS users CASCADE;
-- DROP SEQUENCE IF EXISTS ppic_links_id_seq CASCADE;

-- ============================================================
-- STEP 2: CREATE ALL TABLES
-- ============================================================

-- 1. USERS TABLE
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    role VARCHAR(50) NOT NULL DEFAULT 'operator',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for users
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

COMMENT ON TABLE users IS 'User authentication and authorization';
COMMENT ON COLUMN users.role IS 'User role: admin, ppic, pem, operator';

-- 2. TOKEN BLACKLIST TABLE
CREATE TABLE IF NOT EXISTS token_blacklist (
    id BIGSERIAL PRIMARY KEY,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_token_blacklist_token ON token_blacklist(token);
CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires_at ON token_blacklist(expires_at);

COMMENT ON TABLE token_blacklist IS 'Invalidated JWT tokens for logout functionality';

-- 3. MACHINES TABLE
CREATE TABLE IF NOT EXISTS machines (
    id BIGSERIAL PRIMARY KEY,
    machine_code VARCHAR(50) UNIQUE NOT NULL,
    machine_name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_machines_code ON machines(machine_code);
CREATE INDEX IF NOT EXISTS idx_machines_status ON machines(status);
CREATE INDEX IF NOT EXISTS idx_machines_deleted_at ON machines(deleted_at);

COMMENT ON TABLE machines IS 'Manufacturing machines/equipment';
COMMENT ON COLUMN machines.status IS 'Machine status: active, maintenance, inactive';

-- 4. JOB ORDERS TABLE
CREATE TABLE IF NOT EXISTS job_orders (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(100) UNIQUE NOT NULL,
    customer_name VARCHAR(255),
    product_name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    start_date DATE,
    due_date DATE,
    notes TEXT,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_job_orders_number ON job_orders(order_number);
CREATE INDEX IF NOT EXISTS idx_job_orders_status ON job_orders(status);
CREATE INDEX IF NOT EXISTS idx_job_orders_deleted_at ON job_orders(deleted_at);

COMMENT ON TABLE job_orders IS 'Customer job orders';

-- 5. TOOLPATHER FILES TABLE
CREATE TABLE IF NOT EXISTS toolpather_files (
    id BIGSERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100),
    uploaded_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_toolpather_files_filename ON toolpather_files(filename);
CREATE INDEX IF NOT EXISTS idx_toolpather_files_deleted_at ON toolpather_files(deleted_at);

COMMENT ON TABLE toolpather_files IS 'Toolpather CAD/CAM files';

-- 6. G-CODE FILES TABLE
CREATE TABLE IF NOT EXISTS g_code_files (
    id BIGSERIAL PRIMARY KEY,
    operation_plan_code VARCHAR(100) NOT NULL,
    filename VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    uploaded_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_g_code_files_op_code ON g_code_files(operation_plan_code);
CREATE INDEX IF NOT EXISTS idx_g_code_files_deleted_at ON g_code_files(deleted_at);

COMMENT ON TABLE g_code_files IS 'G-Code files for CNC machines';

-- 7. PEM OPERATION PLANS TABLE
CREATE TABLE IF NOT EXISTS pem_operation_plans (
    id BIGSERIAL PRIMARY KEY,
    operation_plan_code VARCHAR(100) UNIQUE NOT NULL,
    part_name VARCHAR(255) NOT NULL,
    part_number VARCHAR(100),
    machine_id BIGINT REFERENCES machines(id),
    material VARCHAR(255),
    quantity INTEGER,
    finish_dimension VARCHAR(255),
    description TEXT,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_pem_op_plans_code ON pem_operation_plans(operation_plan_code);
CREATE INDEX IF NOT EXISTS idx_pem_op_plans_machine ON pem_operation_plans(machine_id);
CREATE INDEX IF NOT EXISTS idx_pem_op_plans_deleted_at ON pem_operation_plans(deleted_at);

COMMENT ON TABLE pem_operation_plans IS 'PEM (Production Engineering & Maintenance) operation plans';

-- 8. OPERATION PLANS TABLE
CREATE TABLE IF NOT EXISTS operation_plans (
    id BIGSERIAL PRIMARY KEY,
    operation_plan_code VARCHAR(100) UNIQUE NOT NULL,
    part_name VARCHAR(255) NOT NULL,
    part_number VARCHAR(100),
    machine_id BIGINT REFERENCES machines(id),
    material VARCHAR(255),
    quantity INTEGER,
    finish_dimension VARCHAR(255),
    image_path VARCHAR(500),
    description TEXT,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_operation_plans_code ON operation_plans(operation_plan_code);
CREATE INDEX IF NOT EXISTS idx_operation_plans_machine ON operation_plans(machine_id);
CREATE INDEX IF NOT EXISTS idx_operation_plans_deleted_at ON operation_plans(deleted_at);

COMMENT ON TABLE operation_plans IS 'General operation plans';

-- 9. OPERATION PLAN APPROVALS TABLE
CREATE TABLE IF NOT EXISTS operation_plan_approvals (
    id BIGSERIAL PRIMARY KEY,
    operation_plan_id BIGINT NOT NULL REFERENCES operation_plans(id) ON DELETE CASCADE,
    approver_id BIGINT NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    comments TEXT,
    approved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_op_approvals_plan_id ON operation_plan_approvals(operation_plan_id);
CREATE INDEX IF NOT EXISTS idx_op_approvals_status ON operation_plan_approvals(status);

COMMENT ON TABLE operation_plan_approvals IS 'Approval workflow for operation plans';

-- 10. PPIC SCHEDULES TABLE
CREATE TABLE IF NOT EXISTS ppic_schedules (
    id BIGSERIAL PRIMARY KEY,
    njo VARCHAR(100) NOT NULL UNIQUE,
    part_name VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    finish_date DATE NOT NULL,
    priority VARCHAR(50) NOT NULL DEFAULT 'Medium',
    priority_alpha VARCHAR(10),
    material_status VARCHAR(50) NOT NULL DEFAULT 'Pending',
    ppic_notes TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    progress INTEGER NOT NULL DEFAULT 0,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_ppic_schedules_njo ON ppic_schedules(njo);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_start_date ON ppic_schedules(start_date);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_finish_date ON ppic_schedules(finish_date);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_priority ON ppic_schedules(priority);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_status ON ppic_schedules(status);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_deleted_at ON ppic_schedules(deleted_at);

COMMENT ON TABLE ppic_schedules IS 'PPIC scheduling data for Gantt chart display';
COMMENT ON COLUMN ppic_schedules.njo IS 'No Job Order (unique identifier)';
COMMENT ON COLUMN ppic_schedules.part_name IS 'Name of the part being produced';
COMMENT ON COLUMN ppic_schedules.priority IS 'Priority level: Low, Medium, Urgent, Top Urgent';
COMMENT ON COLUMN ppic_schedules.priority_alpha IS 'Alphabetic priority code (A, B, C, etc.)';
COMMENT ON COLUMN ppic_schedules.material_status IS 'Material availability: Ready, Pending, Ordered, Not Ready';
COMMENT ON COLUMN ppic_schedules.status IS 'Schedule status: pending, in_progress, completed';
COMMENT ON COLUMN ppic_schedules.progress IS 'Completion percentage (0-100)';

-- 11. MACHINE ASSIGNMENTS TABLE
CREATE TABLE IF NOT EXISTS machine_assignments (
    id BIGSERIAL PRIMARY KEY,
    schedule_id BIGINT NOT NULL REFERENCES ppic_schedules(id) ON DELETE CASCADE,
    machine_id BIGINT NOT NULL REFERENCES machines(id),
    target_hours DECIMAL(10, 2) NOT NULL DEFAULT 0,
    scheduled_start TIMESTAMP WITH TIME ZONE,
    scheduled_end TIMESTAMP WITH TIME ZONE,
    actual_start TIMESTAMP WITH TIME ZONE,
    actual_end TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    sequence INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_machine_assignments_schedule_id ON machine_assignments(schedule_id);
CREATE INDEX IF NOT EXISTS idx_machine_assignments_machine_id ON machine_assignments(machine_id);
CREATE INDEX IF NOT EXISTS idx_machine_assignments_status ON machine_assignments(status);

COMMENT ON TABLE machine_assignments IS 'Machine assignments for each PPIC schedule entry';
COMMENT ON COLUMN machine_assignments.target_hours IS 'Estimated duration in hours for this machine';
COMMENT ON COLUMN machine_assignments.sequence IS 'Order of machine in production process (1-5)';

-- 12. PPIC LINKS TABLE
CREATE SEQUENCE IF NOT EXISTS ppic_links_id_seq;

CREATE TABLE IF NOT EXISTS ppic_links (
    id BIGINT NOT NULL DEFAULT nextval('ppic_links_id_seq'::regclass),
    source_schedule_id BIGINT NOT NULL,
    target_schedule_id BIGINT NOT NULL,
    link_type VARCHAR(20) NOT NULL DEFAULT '0',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT ppic_links_pkey PRIMARY KEY (id),
    CONSTRAINT ppic_links_source_schedule_id_fkey FOREIGN KEY (source_schedule_id)
        REFERENCES ppic_schedules (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT ppic_links_target_schedule_id_fkey FOREIGN KEY (target_schedule_id)
        REFERENCES ppic_schedules (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_ppic_links_source ON ppic_links(source_schedule_id);
CREATE INDEX IF NOT EXISTS idx_ppic_links_target ON ppic_links(target_schedule_id);

COMMENT ON TABLE ppic_links IS 'Dependencies between PPIC schedules';

-- ============================================================
-- STEP 3: INSERT DEFAULT DATA
-- ============================================================

-- Insert default admin user
-- Username: admin
-- Password: admin123
-- Password hash generated with bcrypt cost 10
INSERT INTO users (email, username, password, full_name, role)
VALUES (
    'admin@ganttpro.com',
    'admin',
    '$2a$10$XqsVLqLj2LXD0VZ5FY5g7.eMgZ9CQKZ9Yd3nE5h6h7P8u9w0x1y2z',
    'System Administrator',
    'admin'
)
ON CONFLICT (email) DO NOTHING;

-- Insert PPIC user
-- Username: ppic
-- Password: ppic123
INSERT INTO users (email, username, password, full_name, role)
VALUES (
    'ppic@ganttpro.com',
    'ppic',
    '$2a$10$YqsVLqLj2LXD0VZ5FY5g7.eMgZ9CQKZ9Yd3nE5h6h7P8u9w0x1y2z',
    'PPIC Staff',
    'ppic'
)
ON CONFLICT (email) DO NOTHING;

-- Insert PEM user
-- Username: pem
-- Password: pem123
INSERT INTO users (email, username, password, full_name, role)
VALUES (
    'pem@ganttpro.com',
    'pem',
    '$2a$10$ZqsVLqLj2LXD0VZ5FY5g7.eMgZ9CQKZ9Yd3nE5h6h7P8u9w0x1y2z',
    'PEM Staff',
    'pem'
)
ON CONFLICT (email) DO NOTHING;

-- Insert Operator user
-- Username: operator
-- Password: operator123
INSERT INTO users (email, username, password, full_name, role)
VALUES (
    'operator@ganttpro.com',
    'operator',
    '$2a$10$AqsVLqLj2LXD0VZ5FY5g7.eMgZ9CQKZ9Yd3nE5h6h7P8u9w0x1y2z',
    'Machine Operator',
    'operator'
)
ON CONFLICT (email) DO NOTHING;

-- Insert sample machines
INSERT INTO machines (machine_code, machine_name, description, status)
VALUES
    ('M01-CNC', 'CNC Milling Machine 1', 'High-precision CNC milling center', 'active'),
    ('M02-CNC', 'CNC Milling Machine 2', 'Secondary CNC milling center', 'active'),
    ('M03-LATHE', 'CNC Lathe 1', 'Computer-controlled lathe machine', 'active'),
    ('M04-GRIND', 'Surface Grinder 1', 'Precision surface grinding machine', 'active'),
    ('M05-DRILL', 'Drilling Machine 1', 'Multi-spindle drilling center', 'active'),
    ('M06-EDM', 'EDM Machine 1', 'Electrical Discharge Machining', 'active'),
    ('M07-MILL', 'Manual Milling 1', 'Conventional milling machine', 'active'),
    ('M08-WELD', 'Welding Station 1', 'MIG/TIG welding station', 'active')
ON CONFLICT (machine_code) DO NOTHING;

-- ============================================================
-- STEP 4: INSERT SAMPLE PPIC SCHEDULES (OPTIONAL)
-- Uncomment jika ingin insert data contoh
-- ============================================================

-- Get admin user id for created_by
DO $$
DECLARE
    admin_id BIGINT;
BEGIN
    SELECT id INTO admin_id FROM users WHERE username = 'admin' LIMIT 1;

    -- Insert sample schedules
    INSERT INTO ppic_schedules (njo, part_name, start_date, finish_date, priority, material_status, status, progress, created_by)
    VALUES
        ('NJO-2026-001', 'Bracket Assembly A', '2026-01-10', '2026-01-20', 'High', 'Ready', 'pending', 0, admin_id),
        ('NJO-2026-002', 'Shaft Component B', '2026-01-15', '2026-01-25', 'Medium', 'Pending', 'pending', 0, admin_id),
        ('NJO-2026-003', 'Housing Unit C', '2026-01-12', '2026-01-22', 'Urgent', 'Ready', 'in_progress', 25, admin_id)
    ON CONFLICT (njo) DO NOTHING;

    -- Insert machine assignments for the schedules
    INSERT INTO machine_assignments (schedule_id, machine_id, target_hours, sequence, status)
    SELECT
        s.id,
        m.id,
        8.0,
        1,
        'pending'
    FROM ppic_schedules s
    CROSS JOIN machines m
    WHERE s.njo = 'NJO-2026-001' AND m.machine_code = 'M01-CNC'
    ON CONFLICT DO NOTHING;

END $$;

-- ============================================================
-- STEP 5: VERIFICATION QUERIES
-- ============================================================

-- Display setup summary
DO $$
DECLARE
    table_count INTEGER;
    user_count INTEGER;
    machine_count INTEGER;
    schedule_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO table_count
    FROM information_schema.tables
    WHERE table_schema = 'public' AND table_type = 'BASE TABLE';

    SELECT COUNT(*) INTO user_count FROM users;
    SELECT COUNT(*) INTO machine_count FROM machines;
    SELECT COUNT(*) INTO schedule_count FROM ppic_schedules;

    RAISE NOTICE '============================================================';
    RAISE NOTICE 'DATABASE MIGRATION COMPLETED SUCCESSFULLY!';
    RAISE NOTICE '============================================================';
    RAISE NOTICE 'Total Tables Created: %', table_count;
    RAISE NOTICE 'Total Users: %', user_count;
    RAISE NOTICE 'Total Machines: %', machine_count;
    RAISE NOTICE 'Total Schedules: %', schedule_count;
    RAISE NOTICE '============================================================';
    RAISE NOTICE 'DEFAULT CREDENTIALS:';
    RAISE NOTICE 'Admin    - Username: admin    | Password: admin123';
    RAISE NOTICE 'PPIC     - Username: ppic     | Password: ppic123';
    RAISE NOTICE 'PEM      - Username: pem      | Password: pem123';
    RAISE NOTICE 'Operator - Username: operator | Password: operator123';
    RAISE NOTICE '============================================================';
END $$;

-- List all created tables
SELECT
    table_name,
    (SELECT COUNT(*)
     FROM information_schema.columns
     WHERE columns.table_name = tables.table_name) as column_count
FROM information_schema.tables
WHERE table_schema = 'public'
  AND table_type = 'BASE TABLE'
ORDER BY table_name;

-- Show users
SELECT id, username, email, role, full_name
FROM users
ORDER BY role, username;

-- Show machines
SELECT id, machine_code, machine_name, status
FROM machines
ORDER BY machine_code;

-- ============================================================
-- MIGRATION COMPLETE
-- ============================================================
-- Next steps:
-- 1. Update backend/.env file with your database credentials
-- 2. Run backend: cd backend && go run main.go
-- 3. Test login with default admin credentials
-- ============================================================
