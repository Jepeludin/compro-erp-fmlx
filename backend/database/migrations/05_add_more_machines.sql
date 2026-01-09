-- ============================================================
-- ADD MORE MACHINES
-- Tambah mesin-mesin tambahan
-- ============================================================

INSERT INTO machines (machine_code, machine_name, description, status)
VALUES 
    ('CNC-003', 'CNC Machine 3', '3-axis CNC machining center', 'active'),
    ('CNC-004', 'CNC Machine 4', '5-axis CNC machining center', 'active'),
    ('MILL-002', 'Milling Machine 2', 'Horizontal milling machine', 'active'),
    ('LATHE-002', 'Lathe Machine 2', 'Manual Lathe', 'active'),
    ('DRILL-001', 'Drilling Machine 1', 'Radial drilling machine', 'active'),
    ('WELD-001', 'Welding Station 1', 'TIG/MIG welding', 'active'),
    ('QC-001', 'Quality Control Station', 'Measurement and inspection', 'active')
ON CONFLICT (machine_code) DO NOTHING;

-- Verify
SELECT COUNT(*) as total_machines FROM machines;
SELECT machine_code, machine_name, status FROM machines ORDER BY id;
