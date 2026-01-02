-- Create ppic_schedules table
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

-- Create machine_assignments table
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

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_njo ON ppic_schedules(njo);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_start_date ON ppic_schedules(start_date);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_finish_date ON ppic_schedules(finish_date);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_priority ON ppic_schedules(priority);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_status ON ppic_schedules(status);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_deleted_at ON ppic_schedules(deleted_at);
CREATE INDEX IF NOT EXISTS idx_machine_assignments_schedule_id ON machine_assignments(schedule_id);
CREATE INDEX IF NOT EXISTS idx_machine_assignments_machine_id ON machine_assignments(machine_id);
CREATE INDEX IF NOT EXISTS idx_machine_assignments_status ON machine_assignments(status);