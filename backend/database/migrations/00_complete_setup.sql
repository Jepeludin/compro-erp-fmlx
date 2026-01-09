-- ============================================================
-- COMPLETE DATABASE SETUP FROM SCRATCH
-- Run this script to create all tables from scratch
-- ============================================================

-- Drop existing tables if needed (use with caution!)
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

-- ============================================================
-- 1. USERS TABLE
-- ============================================================
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

-- ============================================================
-- 2. TOKEN BLACKLIST TABLE
-- ============================================================
CREATE TABLE IF NOT EXISTS token_blacklist (
    id BIGSERIAL PRIMARY KEY,
    token TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_token_blacklist_token ON token_blacklist(token);
CREATE INDEX IF NOT EXISTS idx_token_blacklist_expires_at ON token_blacklist(expires_at);

-- ============================================================
-- 3. MACHINES TABLE
-- ============================================================
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

-- ============================================================
-- 4. JOB ORDERS TABLE
-- ============================================================
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

-- ============================================================
-- 5. TOOLPATHER FILES TABLE
-- ============================================================
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

-- ============================================================
-- 6. G-CODE FILES TABLE
-- ============================================================
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

-- ============================================================
-- 7. PEM OPERATION PLANS TABLE
-- ============================================================
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

-- ============================================================
-- 8. OPERATION PLANS TABLE
-- ============================================================
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

-- ============================================================
-- 9. OPERATION PLAN APPROVALS TABLE
-- ============================================================
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

-- ============================================================
-- 10. PPIC SCHEDULES TABLE
-- ============================================================
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

-- Create indexes for ppic_schedules
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_njo ON ppic_schedules(njo);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_start_date ON ppic_schedules(start_date);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_finish_date ON ppic_schedules(finish_date);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_priority ON ppic_schedules(priority);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_status ON ppic_schedules(status);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_deleted_at ON ppic_schedules(deleted_at);

-- Add comments to ppic_schedules
COMMENT ON TABLE ppic_schedules IS 'PPIC scheduling data for Gantt chart display';
COMMENT ON COLUMN ppic_schedules.njo IS 'Order Number (unique identifier)';
COMMENT ON COLUMN ppic_schedules.part_name IS 'Name of the part being produced';
COMMENT ON COLUMN ppic_schedules.priority IS 'Priority level: Low, Medium, Urgent, Top Urgent';
COMMENT ON COLUMN ppic_schedules.priority_alpha IS 'Alphabetic priority code (A, B, C, etc.)';
COMMENT ON COLUMN ppic_schedules.material_status IS 'Material availability: Ready, Pending, Ordered, Not Ready';
COMMENT ON COLUMN ppic_schedules.status IS 'Schedule status: pending, in_progress, completed';
COMMENT ON COLUMN ppic_schedules.progress IS 'Completion percentage (0-100)';

-- ============================================================
-- 11. MACHINE ASSIGNMENTS TABLE
-- ============================================================
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

-- Create indexes for machine_assignments
CREATE INDEX IF NOT EXISTS idx_machine_assignments_schedule_id ON machine_assignments(schedule_id);
CREATE INDEX IF NOT EXISTS idx_machine_assignments_machine_id ON machine_assignments(machine_id);
CREATE INDEX IF NOT EXISTS idx_machine_assignments_status ON machine_assignments(status);

-- Add comments to machine_assignments
COMMENT ON TABLE machine_assignments IS 'Machine assignments for each PPIC schedule entry';
COMMENT ON COLUMN machine_assignments.target_hours IS 'Estimated duration in hours for this machine';
COMMENT ON COLUMN machine_assignments.sequence IS 'Order of machine in production process (1-5)';
COMMENT ON COLUMN machine_assignments.scheduled_start IS 'Planned start time for this machine';
COMMENT ON COLUMN machine_assignments.scheduled_end IS 'Planned end time for this machine';
COMMENT ON COLUMN machine_assignments.actual_start IS 'Actual start time';
COMMENT ON COLUMN machine_assignments.actual_end IS 'Actual end time';

-- ============================================================
-- 12. PPIC LINKS TABLE
-- ============================================================
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

-- ============================================================
-- INSERT SAMPLE DATA
-- ============================================================

-- Insert default admin user (password: admin123)
-- Password is hashed using bcrypt
INSERT INTO users (email, username, password, full_name, role)
VALUES ('admin@example.com', 'admin', '$2a$10$XqsVLqLj2LXD0VZ5FY5g7.eMgZ9CQKZ9Yd3nE5h6h7P8u9w0x1y2z', 'Administrator', 'admin')
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

-- ============================================================
-- VERIFICATION QUERIES
-- ============================================================

-- Run these queries to verify the setup:

-- Check all tables exist:
-- SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name;

-- Check users:
-- SELECT id, email, username, role FROM users;

-- Check machines:
-- SELECT id, machine_code, machine_name, status FROM machines;

-- Check table structure for machine_assignments:
-- SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = 'machine_assignments' ORDER BY ordinal_position;

-- Check table structure for ppic_schedules:
-- SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = 'ppic_schedules' ORDER BY ordinal_position;

-- ============================================================
-- DONE
-- ============================================================
