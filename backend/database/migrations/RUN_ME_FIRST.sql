-- ============================================================
-- COMPLETE DATABASE SETUP - ALL IN ONE
-- Jalankan file ini untuk setup database dari nol
-- ============================================================
-- 
-- Cara pakai:
-- psql -U postgres -p 3241 -d ganttpro_db -f database/migrations/RUN_ME_FIRST.sql
--
-- ============================================================

\echo ''
\echo '========================================='
\echo '  Starting Complete Database Setup'
\echo '========================================='
\echo ''

-- ============================================================
-- STEP 1: CREATE ALL TABLES
-- ============================================================
\echo 'Step 1: Creating all tables...'

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

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- 2. TOKEN BLACKLIST TABLE
CREATE TABLE IF NOT EXISTS token_blacklist (
    id BIGSERIAL PRIMARY KEY,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_token_blacklist_token ON token_blacklist(token);
CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires_at ON token_blacklist(expires_at);

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
        REFERENCES ppic_schedules (id) ON DELETE CASCADE,
    CONSTRAINT ppic_links_target_schedule_id_fkey FOREIGN KEY (target_schedule_id)
        REFERENCES ppic_schedules (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_ppic_links_source ON ppic_links(source_schedule_id);
CREATE INDEX IF NOT EXISTS idx_ppic_links_target ON ppic_links(target_schedule_id);

\echo '✓ All tables created'
\echo ''

-- ============================================================
-- STEP 2: INSERT DEFAULT DATA
-- ============================================================
\echo 'Step 2: Inserting default data...'

-- Insert admin user with correct password hash for "admin123"
INSERT INTO users (email, username, password, full_name, role)
VALUES ('admin@example.com', 'admin', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Administrator', 'admin')
ON CONFLICT (email) DO NOTHING;

-- Insert sample machines
INSERT INTO machines (machine_code, machine_name, description, status)
VALUES 
    ('CNC-001', 'CNC Machine 1', 'Primary CNC machining center', 'active'),
    ('CNC-002', 'CNC Machine 2', 'Secondary CNC machining center', 'active'),
    ('MILL-001', 'Milling Machine 1', 'Vertical milling machine', 'active'),
    ('LATHE-001', 'Lathe Machine 1', 'CNC Lathe', 'active'),
    ('GRIND-001', 'Grinding Machine 1', 'Surface grinding', 'active')
ON CONFLICT (machine_code) DO NOTHING;

\echo '✓ Default data inserted'
\echo ''

-- ============================================================
-- STEP 3: INSERT SAMPLE DATA (OPTIONAL)
-- ============================================================
\echo 'Step 3: Inserting sample data for testing...'

-- Additional users
INSERT INTO users (email, username, password, full_name, role)
VALUES 
    ('ppic@example.com', 'ppic', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'PPIC User', 'ppic'),
    ('operator@example.com', 'operator', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'Operator User', 'operator')
ON CONFLICT (email) DO NOTHING;

-- Sample job orders
INSERT INTO job_orders (order_number, customer_name, product_name, quantity, status, start_date, due_date, notes, created_by)
VALUES 
    ('JO-2026-001', 'PT. ABC Indonesia', 'Gear Box Housing', 100, 'in_progress', '2026-01-01', '2026-01-31', 'High priority order', 1),
    ('JO-2026-002', 'PT. XYZ Manufacturing', 'Shaft Connector', 200, 'pending', '2026-01-10', '2026-02-15', 'Standard order', 1)
ON CONFLICT (order_number) DO NOTHING;

-- Sample PPIC schedules
INSERT INTO ppic_schedules (njo, part_name, start_date, finish_date, priority, priority_alpha, material_status, ppic_notes, status, progress, created_by)
VALUES 
    ('NJO-2026-001', 'Gear Box Housing - Roughing', '2026-01-06', '2026-01-10', 'Urgent', 'A', 'Ready', 'Material sudah tersedia', 'in_progress', 30, 1),
    ('NJO-2026-002', 'Gear Box Housing - Finishing', '2026-01-11', '2026-01-15', 'Urgent', 'A', 'Ready', 'Lanjutan dari roughing', 'pending', 0, 1)
ON CONFLICT (njo) DO NOTHING;

-- Machine assignments
INSERT INTO machine_assignments (schedule_id, machine_id, sequence, target_hours, scheduled_start, scheduled_end, status)
SELECT 
    ps.id, 1, 1, 8.0,
    '2026-01-06 08:00:00', '2026-01-06 16:00:00', 'in_progress'
FROM ppic_schedules ps WHERE ps.njo = 'NJO-2026-001'
ON CONFLICT DO NOTHING;

INSERT INTO machine_assignments (schedule_id, machine_id, sequence, target_hours, scheduled_start, scheduled_end, status)
SELECT 
    ps.id, 5, 1, 16.0,
    '2026-01-11 08:00:00', '2026-01-13 16:00:00', 'pending'
FROM ppic_schedules ps WHERE ps.njo = 'NJO-2026-002'
ON CONFLICT DO NOTHING;

-- PPIC links
INSERT INTO ppic_links (source_schedule_id, target_schedule_id, link_type)
SELECT 
    (SELECT id FROM ppic_schedules WHERE njo = 'NJO-2026-001'),
    (SELECT id FROM ppic_schedules WHERE njo = 'NJO-2026-002'),
    '0'
WHERE EXISTS (SELECT 1 FROM ppic_schedules WHERE njo = 'NJO-2026-001')
  AND EXISTS (SELECT 1 FROM ppic_schedules WHERE njo = 'NJO-2026-002');

\echo '✓ Sample data inserted'
\echo ''

-- ============================================================
-- STEP 4: VERIFICATION
-- ============================================================
\echo 'Step 4: Verifying setup...'
\echo ''

SELECT 'Table Count:' as info, COUNT(*) as total 
FROM information_schema.tables 
WHERE table_schema = 'public';

SELECT 'Users:' as info, COUNT(*) as total FROM users;
SELECT 'Machines:' as info, COUNT(*) as total FROM machines;
SELECT 'PPIC Schedules:' as info, COUNT(*) as total FROM ppic_schedules;
SELECT 'Machine Assignments:' as info, COUNT(*) as total FROM machine_assignments;

\echo ''
\echo '========================================='
\echo '  Database Setup BERHASIL!'
\echo '========================================='
\echo ''
\echo 'Default Login:'
\echo '  Email: admin@example.com'
\echo '  Username: admin'
\echo '  Password: admin123'
\echo ''
\echo 'Next steps:'
\echo '  1. Update .env file dengan:'
\echo '     DB_HOST=localhost'
\echo '     DB_PORT=3241'
\echo '     DB_USER=postgres'
\echo '     DB_NAME=ganttpro_db'
\echo ''
\echo '  2. Run backend: go run main.go'
\echo ''
\echo '  3. Test API: http://localhost:8080/api/ppic/gantt-data'
\echo ''
\echo 'PENTING: Ganti password default setelah login!'
\echo ''
