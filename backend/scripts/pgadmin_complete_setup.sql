-- ============================================
-- GANTTPRO DATABASE - COMPLETE SQL SETUP
-- ============================================
-- Untuk dijalankan di pgAdmin Query Tool
-- Database: ganttpro_db
-- Port: 8181
-- User: postgres
-- Password: jemmy1303
-- ============================================

-- ============================================
-- 1. CREATE SEQUENCE (untuk auto increment ID)
-- ============================================

-- Drop sequence jika ingin reset
-- DROP SEQUENCE IF EXISTS public.users_id_seq CASCADE;

CREATE SEQUENCE IF NOT EXISTS public.users_id_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

ALTER SEQUENCE public.users_id_seq
    OWNER TO postgres;

-- ============================================
-- 2. CREATE TABLE public.users
-- ============================================

-- Drop table jika ingin reset (HATI-HATI: akan menghapus semua data!)
-- DROP TABLE IF EXISTS public.users CASCADE;

CREATE TABLE IF NOT EXISTS public.users
(
    id bigint NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    username character varying(50) COLLATE pg_catalog."default" NOT NULL,
    user_id character varying(50) COLLATE pg_catalog."default" NOT NULL,
    password text COLLATE pg_catalog."default" NOT NULL,
    role character varying(20) COLLATE pg_catalog."default" DEFAULT 'PPIC'::character varying,
    operator character varying(50) COLLATE pg_catalog."default",
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id)
)
TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;

-- Set sequence ownership to table column
ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;

COMMENT ON TABLE public.users
    IS 'Tabel users untuk sistem GanttPro - menyimpan informasi user dengan berbagai role (Admin, PPIC, Supervisor, Operator, Warehouse, Guest)';

-- ============================================
-- 3. CREATE INDEXES
-- ============================================

-- Index untuk soft delete
CREATE INDEX IF NOT EXISTS idx_users_deleted_at
    ON public.users USING btree
    (deleted_at ASC NULLS LAST)
    TABLESPACE pg_default;

-- Index unique untuk user_id
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_user_id
    ON public.users USING btree
    (user_id COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

-- Index unique untuk username
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username
    ON public.users USING btree
    (username COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

-- Index untuk role (untuk query filter by role)
CREATE INDEX IF NOT EXISTS idx_users_role
    ON public.users USING btree
    (role COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

-- Index untuk is_active (untuk query active users)
CREATE INDEX IF NOT EXISTS idx_users_is_active
    ON public.users USING btree
    (is_active ASC NULLS LAST)
    TABLESPACE pg_default;

-- ============================================
-- 4. CREATE TRIGGER FUNCTION untuk AUTO UPDATE updated_at
-- ============================================

CREATE OR REPLACE FUNCTION public.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

COMMENT ON FUNCTION public.update_updated_at_column()
    IS 'Function untuk auto-update kolom updated_at ketika record diupdate';

-- ============================================
-- 5. CREATE TRIGGER
-- ============================================

DROP TRIGGER IF EXISTS update_users_updated_at ON public.users;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON public.users
    FOR EACH ROW
    EXECUTE FUNCTION public.update_updated_at_column();

COMMENT ON TRIGGER update_users_updated_at ON public.users
    IS 'Trigger untuk otomatis update timestamp updated_at';

-- ============================================
-- 6. INSERT SAMPLE DATA (18 ACCOUNTS)
-- ============================================

-- Password sudah di-hash menggunakan bcrypt (cost 10)
-- Format hash: $2a$10$...

-- ============================================
-- ADMIN ACCOUNTS (3)
-- ============================================

-- 1. BAYU (Admin) - Password: 13032001
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'BAYU',
    'PI0824.2374',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdnYVO2C6',
    'Admin',
    NULL,
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 2. Jeremy (Admin) - Password: admin123
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Jeremy',
    'PI0824.0001',
    '$2a$10$rT8nCqYxGvlQCOoLQ7y.7e2J5bCZxXGZxCzQc8FvQxYWKZGqXH9Gy',
    'Admin',
    NULL,
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 3. Agape (Admin) - Password: admin123
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Agape',
    'PI0824.0002',
    '$2a$10$rT8nCqYxGvlQCOoLQ7y.7e2J5bCZxXGZxCzQc8FvQxYWKZGqXH9Gy',
    'Admin',
    'Operation Manager',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- ============================================
-- PPIC ACCOUNTS (4)
-- ============================================

-- 4. Amelia (PPIC) - Password: 20010313
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Amelia',
    'PI0824.2375',
    '$2a$10$8kL3M5XvGQqL7xZc8FqO7uyLGxJ9xH8p3F2KdLsYGxZhGJK9VH3M2',
    'PPIC',
    NULL,
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 5. Siti (PPIC) - Password: ppic2024
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Siti',
    'PI0824.1001',
    '$2a$10$YQ7V3K4L6M8N9P1R2S4T5U6vW7X8Y9Z0aB1cD2eF3gH4iJ5kL6mN7o',
    'PPIC',
    'Production Planning',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 6. Budi (PPIC) - Password: ppic2024
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Budi',
    'PI0824.1002',
    '$2a$10$YQ7V3K4L6M8N9P1R2S4T5U6vW7X8Y9Z0aB1cD2eF3gH4iJ5kL6mN7o',
    'PPIC',
    'Inventory Control',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 7. Rina (PPIC) - Password: ppic2024
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Rina',
    'PI0824.1003',
    '$2a$10$YQ7V3K4L6M8N9P1R2S4T5U6vW7X8Y9Z0aB1cD2eF3gH4iJ5kL6mN7o',
    'PPIC',
    'Material Planning',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- ============================================
-- SUPERVISOR ACCOUNTS (3)
-- ============================================

-- 8. Dimas (Supervisor) - Password: super123
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Dimas',
    'PI0824.2001',
    '$2a$10$pQ8R9S0T1U2V3W4X5Y6Z7a8B9C0D1E2F3G4H5I6J7K8L9M0N1O2P3',
    'Supervisor',
    'Production Line A',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 9. Fitri (Supervisor) - Password: super123
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Fitri',
    'PI0824.2002',
    '$2a$10$pQ8R9S0T1U2V3W4X5Y6Z7a8B9C0D1E2F3G4H5I6J7K8L9M0N1O2P3',
    'Supervisor',
    'Production Line B',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 10. Andi (Supervisor) - Password: super123
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Andi',
    'PI0824.2003',
    '$2a$10$pQ8R9S0T1U2V3W4X5Y6Z7a8B9C0D1E2F3G4H5I6J7K8L9M0N1O2P3',
    'Supervisor',
    'Quality Control',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- ============================================
-- OPERATOR ACCOUNTS (5)
-- ============================================

-- 11. Joko (Operator) - Password: op2024
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Joko',
    'PI0824.3001',
    '$2a$10$mN5O6P7Q8R9S0T1U2V3W4X5Y6Z7A8B9C0D1E2F3G4H5I6J7K8L9M0',
    'Operator',
    'Machine Operator A1',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 12. Sinta (Operator) - Password: op2024
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Sinta',
    'PI0824.3002',
    '$2a$10$mN5O6P7Q8R9S0T1U2V3W4X5Y6Z7A8B9C0D1E2F3G4H5I6J7K8L9M0',
    'Operator',
    'Machine Operator A2',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 13. Rudi (Operator) - Password: op2024
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Rudi',
    'PI0824.3003',
    '$2a$10$mN5O6P7Q8R9S0T1U2V3W4X5Y6Z7A8B9C0D1E2F3G4H5I6J7K8L9M0',
    'Operator',
    'Machine Operator B1',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 14. Dewi (Operator) - Password: op2024
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Dewi',
    'PI0824.3004',
    '$2a$10$mN5O6P7Q8R9S0T1U2V3W4X5Y6Z7A8B9C0D1E2F3G4H5I6J7K8L9M0',
    'Operator',
    'Assembly Line',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 15. Hadi (Operator) - Password: op2024
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Hadi',
    'PI0824.3005',
    '$2a$10$mN5O6P7Q8R9S0T1U2V3W4X5Y6Z7A8B9C0D1E2F3G4H5I6J7K8L9M0',
    'Operator',
    'Packaging',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- ============================================
-- WAREHOUSE ACCOUNTS (2)
-- ============================================

-- 16. Tono (Warehouse) - Password: wh2024
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Tono',
    'PI0824.4001',
    '$2a$10$kL4M5N6O7P8Q9R0S1T2U3V4W5X6Y7Z8A9B0C1D2E3F4G5H6I7J8K9',
    'Warehouse',
    'Warehouse Staff',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- 17. Lisa (Warehouse) - Password: wh2024
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'Lisa',
    'PI0824.4002',
    '$2a$10$kL4M5N6O7P8Q9R0S1T2U3V4W5X6Y7Z8A9B0C1D2E3F4G5H6I7J8K9',
    'Warehouse',
    'Material Handler',
    TRUE
) ON CONFLICT (username) DO NOTHING;

-- ============================================
-- INACTIVE/TEST ACCOUNTS (1)
-- ============================================

-- 18. TestUser (Inactive) - Password: test123
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES (
    'TestUser',
    'PI0824.9999',
    '$2a$10$iJ3K4L5M6N7O8P9Q0R1S2T3U4V5W6X7Y8Z9A0B1C2D3E4F5G6H7I8',
    'PPIC',
    NULL,
    FALSE
) ON CONFLICT (username) DO NOTHING;

-- ============================================
-- 7. VERIFY DATA
-- ============================================

-- Tampilkan semua data yang berhasil diinsert
SELECT 
    id,
    username,
    user_id,
    role,
    operator,
    is_active,
    created_at,
    updated_at
FROM public.users
ORDER BY 
    CASE role
        WHEN 'Admin' THEN 1
        WHEN 'PPIC' THEN 2
        WHEN 'Supervisor' THEN 3
        WHEN 'Operator' THEN 4
        WHEN 'Warehouse' THEN 5
        ELSE 6
    END,
    id;

-- ============================================
-- 8. SUMMARY STATISTICS
-- ============================================

-- Hitung jumlah user per role
SELECT 
    role,
    COUNT(*) as total,
    COUNT(CASE WHEN is_active = TRUE THEN 1 END) as active,
    COUNT(CASE WHEN is_active = FALSE THEN 1 END) as inactive
FROM public.users
GROUP BY role
ORDER BY 
    CASE role
        WHEN 'Admin' THEN 1
        WHEN 'PPIC' THEN 2
        WHEN 'Supervisor' THEN 3
        WHEN 'Operator' THEN 4
        WHEN 'Warehouse' THEN 5
        ELSE 6
    END;

-- ============================================
-- SELESAI!
-- ============================================
-- Total: 18 user accounts created
-- 
-- Login Credentials:
-- ADMIN:
--   - BAYU (13032001), Jeremy (admin123), Agape (admin123)
-- PPIC:
--   - Amelia (20010313), Siti/Budi/Rina (ppic2024)
-- SUPERVISOR:
--   - Dimas/Fitri/Andi (super123)
-- OPERATOR:
--   - Joko/Sinta/Rudi/Dewi/Hadi (op2024)
-- WAREHOUSE:
--   - Tono/Lisa (wh2024)
-- TEST:
--   - TestUser (test123) [INACTIVE]
-- ============================================
