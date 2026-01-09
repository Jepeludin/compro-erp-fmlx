# Setup Database dari Nol

Panduan ini akan membantu Anda setup database PostgreSQL dari awal di laptop baru.

## Prerequisites

1. PostgreSQL terinstal (versi 12 atau lebih baru)
2. Akses ke PostgreSQL dengan user yang punya privilege CREATE DATABASE

## Langkah 1: Buat Database Baru

Buka terminal/command prompt dan login ke PostgreSQL:

```bash
psql -U postgres
```

Kemudian buat database baru:

```sql
CREATE DATABASE ganttpro_db;
```

Keluar dari psql:
```sql
\q
```

## Langkah 2: Update File .env

Pastikan file `.env` di root folder backend berisi konfigurasi yang benar:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=ganttpro_db
DB_SSLMODE=disable

# Server Configuration
PORT=8080

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRY=24h

# Environment
ENVIRONMENT=development
```

**Penting**: Ganti `your_password_here` dengan password PostgreSQL Anda!

## Langkah 3: Jalankan Migration Script

Ada 2 cara untuk menjalankan migration:

### Cara 1: Menggunakan psql (Recommended)

```bash
# Navigate ke folder backend
cd c:\Jemmy\compro\compro-erp-fmlx\backend

# Run migration script
psql -U postgres -d ganttpro_db -f database/migrations/00_complete_setup.sql
```

### Cara 2: Menggunakan pgAdmin

1. Buka pgAdmin
2. Connect ke server PostgreSQL Anda
3. Klik kanan pada database `ganttpro_db` → Query Tool
4. Buka file `database/migrations/00_complete_setup.sql`
5. Copy semua isinya ke Query Tool
6. Klik Execute (F5)

## Langkah 4: Verifikasi Database

Cek apakah semua table sudah dibuat:

```sql
-- Login ke database
psql -U postgres -d ganttpro_db

-- List semua tables
\dt

-- Cek struktur table machine_assignments
\d machine_assignments

-- Cek struktur table ppic_schedules
\d ppic_schedules

-- Cek apakah user admin sudah ada
SELECT id, email, username, role FROM users;

-- Cek apakah machines sudah ada
SELECT id, machine_code, machine_name, status FROM machines;
```

Output yang diharapkan:
- Minimal 12 tables: users, token_blacklist, machines, job_orders, toolpather_files, g_code_files, pem_operation_plans, operation_plans, operation_plan_approvals, ppic_schedules, machine_assignments, ppic_links
- 1 user admin dengan email `admin@example.com`
- 5 machines (CNC-001, CNC-002, MILL-001, LATHE-001, GRIND-001)

## Langkah 5: Test Backend

Jalankan backend:

```bash
cd c:\Jemmy\compro\compro-erp-fmlx\backend
go run main.go
```

Output yang diharapkan:
```
Connected to database successfully!
Server is running on port 8080
```

## Langkah 6: Test API

Buka browser dan test endpoint:

```
http://localhost:8080/api/ppic/gantt-data
```

Atau gunakan curl:

```bash
curl http://localhost:8080/api/ppic/gantt-data
```

## Troubleshooting

### Error: "column ms.scheduled_start does not exist"

Ini berarti migration belum jalan dengan benar. Solusi:

1. Drop database dan buat ulang:
```sql
psql -U postgres
DROP DATABASE ganttpro_db;
CREATE DATABASE ganttpro_db;
\q
```

2. Jalankan ulang migration (Langkah 3)

### Error: "database ganttpro_db does not exist"

Jalankan Langkah 1 untuk membuat database.

### Error: "password authentication failed"

Pastikan password di file `.env` sesuai dengan password PostgreSQL Anda.

### Error: "role postgres does not exist"

Ganti `postgres` dengan username PostgreSQL yang Anda gunakan di semua perintah di atas.

### Backend tidak bisa connect ke database

1. Cek apakah PostgreSQL service sedang running:
   - Windows: Buka Services → cari "PostgreSQL" → pastikan status "Running"
   - Linux: `sudo systemctl status postgresql`

2. Cek apakah port 5432 digunakan oleh PostgreSQL:
   - Windows: `netstat -ano | findstr :5432`
   - Linux: `netstat -tlnp | grep :5432`

3. Test koneksi manual:
```bash
psql -U postgres -d ganttpro_db -h localhost -p 5432
```

## Default Login Credentials

Setelah setup, Anda bisa login dengan:

- **Email**: `admin@example.com`
- **Username**: `admin`
- **Password**: `admin123`

**PENTING**: Segera ganti password default ini setelah login pertama kali!

## Reset Database (Jika Perlu)

Jika Anda ingin reset database ke kondisi awal:

```sql
-- Login ke PostgreSQL
psql -U postgres

-- Drop dan buat ulang database
DROP DATABASE ganttpro_db;
CREATE DATABASE ganttpro_db;
\q

-- Jalankan ulang migration
psql -U postgres -d ganttpro_db -f database/migrations/00_complete_setup.sql
```

## Struktur Table Penting

### ppic_schedules
Table untuk menyimpan jadwal produksi:
- `id`: Primary key
- `njo`: Nomor Job Order (unique)
- `part_name`: Nama part
- `start_date`: Tanggal mulai
- `finish_date`: Tanggal selesai
- `priority`: Low, Medium, Urgent, Top Urgent
- `status`: pending, in_progress, completed
- `progress`: 0-100

### machine_assignments
Table untuk menyimpan assignment mesin ke schedule:
- `id`: Primary key
- `schedule_id`: FK ke ppic_schedules
- `machine_id`: FK ke machines
- `target_hours`: Estimasi durasi (jam)
- `scheduled_start`: Jadwal mulai (TIMESTAMP)
- `scheduled_end`: Jadwal selesai (TIMESTAMP)
- `actual_start`: Actual mulai
- `actual_end`: Actual selesai
- `status`: pending, in_progress, completed
- `sequence`: Urutan mesin (1-5)

### ppic_links
Table untuk dependency antar schedule:
- `id`: Primary key
- `source_schedule_id`: Task predecessor
- `target_schedule_id`: Task successor
- `link_type`: "0" = finish-to-start

## Kontak

Jika masih ada masalah, hubungi tim development atau buka issue di repository.
