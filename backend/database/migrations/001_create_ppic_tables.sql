-- Migration: Create PPIC Schedule tables for Gantt Chart feature


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

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_njo ON ppic_schedules(njo);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_start_date ON ppic_schedules(start_date);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_finish_date ON ppic_schedules(finish_date);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_priority ON ppic_schedules(priority);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_status ON ppic_schedules(status);
CREATE INDEX IF NOT EXISTS idx_ppic_schedules_deleted_at ON ppic_schedules(deleted_at);

CREATE INDEX IF NOT EXISTS idx_machine_assignments_schedule_id ON machine_assignments(schedule_id);
CREATE INDEX IF NOT EXISTS idx_machine_assignments_machine_id ON machine_assignments(machine_id);
CREATE INDEX IF NOT EXISTS idx_machine_assignments_status ON machine_assignments(status);

-- Add comments to tables
COMMENT ON TABLE ppic_schedules IS 'PPIC scheduling data for Gantt chart display';
COMMENT ON TABLE machine_assignments IS 'Machine assignments for each PPIC schedule entry';

-- Column comments for ppic_schedules
COMMENT ON COLUMN ppic_schedules.njo IS 'Order Number (unique identifier)';
COMMENT ON COLUMN ppic_schedules.part_name IS 'Name of the part being produced';
COMMENT ON COLUMN ppic_schedules.priority IS 'Priority level: Low, Medium, Urgent, Top Urgent';
COMMENT ON COLUMN ppic_schedules.priority_alpha IS 'Alphabetic priority code (A, B, C, etc.)';
COMMENT ON COLUMN ppic_schedules.material_status IS 'Material availability: Ready, Pending, Ordered, Not Ready';
COMMENT ON COLUMN ppic_schedules.status IS 'Schedule status: pending, in_progress, completed';
COMMENT ON COLUMN ppic_schedules.progress IS 'Completion percentage (0-100)';

-- Column comments for machine_assignments
COMMENT ON COLUMN machine_assignments.target_hours IS 'Estimated duration in hours for this machine';
COMMENT ON COLUMN machine_assignments.sequence IS 'Order of machine in production process (1-5)';
COMMENT ON COLUMN machine_assignments.scheduled_start IS 'Planned start time for this machine';
COMMENT ON COLUMN machine_assignments.scheduled_end IS 'Planned end time for this machine';
COMMENT ON COLUMN machine_assignments.actual_start IS 'Actual start time';
COMMENT ON COLUMN machine_assignments.actual_end IS 'Actual end time';



