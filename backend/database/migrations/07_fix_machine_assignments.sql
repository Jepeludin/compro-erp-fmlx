-- ============================================================
-- FIX MACHINE ASSIGNMENTS TABLE
-- Tambahkan kolom yang hilang ke table machine_assignments
-- ============================================================

-- Tambahkan kolom scheduled_start
ALTER TABLE machine_assignments 
ADD COLUMN IF NOT EXISTS scheduled_start TIMESTAMP WITH TIME ZONE;

-- Tambahkan kolom scheduled_end
ALTER TABLE machine_assignments 
ADD COLUMN IF NOT EXISTS scheduled_end TIMESTAMP WITH TIME ZONE;

-- Tambahkan kolom actual_start (jika belum ada)
ALTER TABLE machine_assignments 
ADD COLUMN IF NOT EXISTS actual_start TIMESTAMP WITH TIME ZONE;

-- Tambahkan kolom actual_end (jika belum ada)
ALTER TABLE machine_assignments 
ADD COLUMN IF NOT EXISTS actual_end TIMESTAMP WITH TIME ZONE;

-- Verify struktur table
\d machine_assignments

-- Check apakah kolom sudah ada
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'machine_assignments' 
  AND column_name IN ('scheduled_start', 'scheduled_end', 'actual_start', 'actual_end')
ORDER BY column_name;

SELECT 'Machine assignments table fixed!' as status;
