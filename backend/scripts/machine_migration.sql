-- ============================================
-- GANTTPRO - MACHINE TRACKING SYSTEM
-- ============================================
-- Database schema untuk machine tracking system
-- Centralized structure untuk unlimited machines
-- ============================================

-- ============================================
-- STEP 1: CREATE NEW TABLES
-- ============================================

-- 1. Tabel Machines (Master data mesin)
CREATE TABLE IF NOT EXISTS machines (
    id SERIAL PRIMARY KEY,
    machine_code VARCHAR(50) UNIQUE NOT NULL,
    machine_name VARCHAR(100) NOT NULL,
    machine_type VARCHAR(50),
    location VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 2. Tabel Job Orders (Order/Project per mesin)
CREATE TABLE IF NOT EXISTS job_orders (
    id SERIAL PRIMARY KEY,
    machine_id INTEGER REFERENCES machines(id) ON DELETE CASCADE,
    njo VARCHAR(50),
    project TEXT,
    item TEXT,
    note TEXT,
    deadline VARCHAR(50),
    operator_id INTEGER REFERENCES users(id),
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 3. Tabel Process Stages (4 stages: setting, proses, cmm, kalibrasi)
CREATE TABLE IF NOT EXISTS process_stages (
    id SERIAL PRIMARY KEY,
    job_order_id INTEGER REFERENCES job_orders(id) ON DELETE CASCADE,
    stage_name VARCHAR(50) NOT NULL,
    start_time TIMESTAMP,
    finish_time TIMESTAMP,
    duration_minutes INTEGER,
    operator_id INTEGER REFERENCES users(id),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- STEP 2: CREATE INDEXES
-- ============================================

CREATE INDEX IF NOT EXISTS idx_machines_code ON machines(machine_code);
CREATE INDEX IF NOT EXISTS idx_machines_status ON machines(status);
CREATE INDEX IF NOT EXISTS idx_job_orders_machine ON job_orders(machine_id);
CREATE INDEX IF NOT EXISTS idx_job_orders_status ON job_orders(status);
CREATE INDEX IF NOT EXISTS idx_job_orders_njo ON job_orders(njo);
CREATE INDEX IF NOT EXISTS idx_process_stages_job ON process_stages(job_order_id);
CREATE INDEX IF NOT EXISTS idx_process_stages_stage ON process_stages(stage_name);

-- ============================================
-- STEP 3: CREATE TRIGGER FOR AUTO-CALCULATE DURATION
-- ============================================

-- Function untuk auto calculate duration
CREATE OR REPLACE FUNCTION calculate_duration()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.start_time IS NOT NULL AND NEW.finish_time IS NOT NULL THEN
        NEW.duration_minutes := EXTRACT(EPOCH FROM (NEW.finish_time - NEW.start_time)) / 60;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger untuk auto calculate duration
DROP TRIGGER IF EXISTS trigger_calculate_duration ON process_stages;
CREATE TRIGGER trigger_calculate_duration
    BEFORE INSERT OR UPDATE ON process_stages
    FOR EACH ROW
    EXECUTE FUNCTION calculate_duration();

-- ============================================
-- STEP 4: INSERT SAMPLE MACHINES
-- ============================================

INSERT INTO machines (machine_code, machine_name, machine_type, location, status) VALUES
('YSD01', 'Yasda', 'CNC', 'Workshop A', 'active'),
('SDC01', 'Sodick', 'EDM', 'Workshop A', 'active'),
('WRC01', 'Wire Cut', 'Wire EDM', 'Workshop B', 'active')
ON CONFLICT (machine_code) DO NOTHING;

-- ============================================
-- STEP 5: VERIFICATION QUERIES
-- ============================================

-- Check machines
SELECT * FROM machines;

-- Check job_orders count
SELECT m.machine_name, COUNT(jo.id) as job_count
FROM machines m
LEFT JOIN job_orders jo ON jo.machine_id = m.id
GROUP BY m.id, m.machine_name;

-- Check process_stages count
SELECT ps.stage_name, COUNT(*) as stage_count
FROM process_stages ps
GROUP BY ps.stage_name;

-- Check duration calculation
SELECT 
    jo.njo,
    ps.stage_name,
    ps.start_time,
    ps.finish_time,
    ps.duration_minutes,
    ROUND(ps.duration_minutes / 60.0, 2) as duration_hours
FROM process_stages ps
JOIN job_orders jo ON jo.id = ps.job_order_id
ORDER BY jo.id, ps.stage_name
LIMIT 10;

-- ============================================
-- STEP 6: GRANT PERMISSIONS
-- ============================================

ALTER TABLE machines OWNER TO postgres;
ALTER TABLE job_orders OWNER TO postgres;
ALTER TABLE process_stages OWNER TO postgres;

-- ============================================
-- NOTES
-- ============================================
-- Duration otomatis dihitung dalam MENIT dari finish_time - start_time
-- Trigger akan auto-calculate duration saat INSERT atau UPDATE
-- Struktur ini scalable untuk unlimited machines
-- Tidak perlu create table baru untuk setiap mesin
-- ============================================
