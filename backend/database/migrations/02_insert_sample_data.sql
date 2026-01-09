-- ============================================================
-- INSERT SAMPLE DATA FOR TESTING
-- ============================================================

-- Insert additional test users
INSERT INTO users (email, username, password, full_name, role)
VALUES 
    ('ppic@example.com', 'ppic', '$2a$10$XqsVLqLj2LXD0VZ5FY5g7.eMgZ9CQKZ9Yd3nE5h6h7P8u9w0x1y2z', 'PPIC User', 'ppic'),
    ('operator@example.com', 'operator', '$2a$10$XqsVLqLj2LXD0VZ5FY5g7.eMgZ9CQKZ9Yd3nE5h6h7P8u9w0x1y2z', 'Operator User', 'operator'),
    ('supervisor@example.com', 'supervisor', '$2a$10$XqsVLqLj2LXD0VZ5FY5g7.eMgZ9CQKZ9Yd3nE5h6h7P8u9w0x1y2z', 'Supervisor', 'supervisor')
ON CONFLICT (email) DO NOTHING;

-- Insert sample job orders
INSERT INTO job_orders (order_number, customer_name, product_name, quantity, status, start_date, due_date, notes, created_by)
VALUES 
    ('JO-2026-001', 'PT. ABC Indonesia', 'Gear Box Housing', 100, 'in_progress', '2026-01-01', '2026-01-31', 'High priority order', 1),
    ('JO-2026-002', 'PT. XYZ Manufacturing', 'Shaft Connector', 200, 'pending', '2026-01-10', '2026-02-15', 'Standard order', 1),
    ('JO-2026-003', 'CV. Jaya Makmur', 'Bracket Mount', 50, 'pending', '2026-01-15', '2026-02-28', 'Custom design', 1)
ON CONFLICT (order_number) DO NOTHING;

-- Insert sample PPIC schedules
INSERT INTO ppic_schedules (njo, part_name, start_date, finish_date, priority, priority_alpha, material_status, ppic_notes, status, progress, created_by)
VALUES 
    ('NJO-2026-001', 'Gear Box Housing - Roughing', '2026-01-06', '2026-01-10', 'Urgent', 'A', 'Ready', 'Material sudah tersedia', 'in_progress', 30, 1),
    ('NJO-2026-002', 'Gear Box Housing - Finishing', '2026-01-11', '2026-01-15', 'Urgent', 'A', 'Ready', 'Lanjutan dari roughing', 'pending', 0, 1),
    ('NJO-2026-003', 'Shaft Connector - Turning', '2026-01-08', '2026-01-12', 'Medium', 'B', 'Pending', 'Menunggu material', 'pending', 0, 1),
    ('NJO-2026-004', 'Bracket Mount - Milling', '2026-01-15', '2026-01-20', 'Low', 'C', 'Ordered', 'Material dalam pesanan', 'pending', 0, 1)
ON CONFLICT (njo) DO NOTHING;

-- Insert machine assignments for PPIC schedules
-- NJO-2026-001: Roughing (CNC-001 -> MILL-001)
INSERT INTO machine_assignments (schedule_id, machine_id, sequence, target_hours, scheduled_start, scheduled_end, status)
SELECT 
    ps.id,
    1, -- CNC-001
    1,
    8.0,
    '2026-01-06 08:00:00',
    '2026-01-06 16:00:00',
    'in_progress'
FROM ppic_schedules ps WHERE ps.njo = 'NJO-2026-001'
ON CONFLICT DO NOTHING;

INSERT INTO machine_assignments (schedule_id, machine_id, sequence, target_hours, scheduled_start, scheduled_end, status)
SELECT 
    ps.id,
    3, -- MILL-001
    2,
    12.0,
    '2026-01-07 08:00:00',
    '2026-01-08 12:00:00',
    'pending'
FROM ppic_schedules ps WHERE ps.njo = 'NJO-2026-001'
ON CONFLICT DO NOTHING;

-- NJO-2026-002: Finishing (GRIND-001)
INSERT INTO machine_assignments (schedule_id, machine_id, sequence, target_hours, scheduled_start, scheduled_end, status)
SELECT 
    ps.id,
    5, -- GRIND-001
    1,
    16.0,
    '2026-01-11 08:00:00',
    '2026-01-13 16:00:00',
    'pending'
FROM ppic_schedules ps WHERE ps.njo = 'NJO-2026-002'
ON CONFLICT DO NOTHING;

-- NJO-2026-003: Turning (LATHE-001)
INSERT INTO machine_assignments (schedule_id, machine_id, sequence, target_hours, scheduled_start, scheduled_end, status)
SELECT 
    ps.id,
    4, -- LATHE-001
    1,
    24.0,
    '2026-01-08 08:00:00',
    '2026-01-11 16:00:00',
    'pending'
FROM ppic_schedules ps WHERE ps.njo = 'NJO-2026-003'
ON CONFLICT DO NOTHING;

-- NJO-2026-004: Milling (MILL-001 -> CNC-002)
INSERT INTO machine_assignments (schedule_id, machine_id, sequence, target_hours, scheduled_start, scheduled_end, status)
SELECT 
    ps.id,
    3, -- MILL-001
    1,
    10.0,
    '2026-01-15 08:00:00',
    '2026-01-16 10:00:00',
    'pending'
FROM ppic_schedules ps WHERE ps.njo = 'NJO-2026-004'
ON CONFLICT DO NOTHING;

INSERT INTO machine_assignments (schedule_id, machine_id, sequence, target_hours, scheduled_start, scheduled_end, status)
SELECT 
    ps.id,
    2, -- CNC-002
    2,
    14.0,
    '2026-01-17 08:00:00',
    '2026-01-18 14:00:00',
    'pending'
FROM ppic_schedules ps WHERE ps.njo = 'NJO-2026-004'
ON CONFLICT DO NOTHING;

-- Insert PPIC links (dependencies)
-- NJO-2026-002 depends on NJO-2026-001 (Finishing after Roughing)
INSERT INTO ppic_links (source_schedule_id, target_schedule_id, link_type)
SELECT 
    (SELECT id FROM ppic_schedules WHERE njo = 'NJO-2026-001'),
    (SELECT id FROM ppic_schedules WHERE njo = 'NJO-2026-002'),
    '0' -- finish-to-start
WHERE EXISTS (SELECT 1 FROM ppic_schedules WHERE njo = 'NJO-2026-001')
  AND EXISTS (SELECT 1 FROM ppic_schedules WHERE njo = 'NJO-2026-002');

-- Verify data inserted
SELECT 'Sample data inserted successfully!' AS status;
SELECT COUNT(*) as user_count FROM users;
SELECT COUNT(*) as job_order_count FROM job_orders;
SELECT COUNT(*) as schedule_count FROM ppic_schedules;
SELECT COUNT(*) as assignment_count FROM machine_assignments;
SELECT COUNT(*) as link_count FROM ppic_links;
