# Quick Start - Setup Database di Laptop Baru

## Masalah yang Dihadapi
Error: `kolom ms.scheduled_start_belum ada (SQLSTATE 42703)`

Ini terjadi karena database di laptop belum di-setup dengan benar.

## Solusi Cepat (Recommended)

### Opsi 1: Gunakan PowerShell Script (Paling Mudah!)

1. Buka PowerShell sebagai Administrator
2. Navigate ke folder backend:
   ```powershell
   cd C:\Jemmy\compro\compro-erp-fmlx\backend
   ```

3. Jalankan script setup:
   ```powershell
   .\setup_database.ps1
   ```

4. Ikuti instruksi:
   - Masukkan username PostgreSQL (default: postgres)
   - Masukkan password PostgreSQL
   - Masukkan nama database (default: ganttpro_db)

5. Script akan otomatis:
   - ✓ Test koneksi ke PostgreSQL
   - ✓ Buat database baru
   - ✓ Jalankan semua migration
   - ✓ Insert data default (admin user & machines)
   - ✓ Update file .env

6. Selesai! Jalankan backend:
   ```powershell
   go run main.go
   ```

### Opsi 2: Manual Setup

#### Step 1: Buat Database

```bash
# Login ke PostgreSQL
psql -U postgres

# Buat database
CREATE DATABASE ganttpro_db;

# Keluar
\q
```

#### Step 2: Jalankan Migration

```bash
# Jalankan migration script
psql -U postgres -d ganttpro_db -f database/migrations/00_complete_setup.sql
```

#### Step 3: Update .env

Edit file `.env` di folder backend:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ganttpro_db
DB_SSLMODE=disable

PORT=8080
JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRY=24h
ENVIRONMENT=development
```

#### Step 4: Verifikasi

```bash
# Cek apakah setup berhasil
psql -U postgres -d ganttpro_db -f database/migrations/verify_database.sql
```

#### Step 5: Jalankan Backend

```bash
go run main.go
```

## Verifikasi Setup Berhasil

Buka browser dan test:
```
http://localhost:8080/api/ppic/gantt-data
```

Harusnya return JSON kosong atau data gantt:
```json
{
  "success": true,
  "data": {
    "sections": [],
    "links": []
  }
}
```

**BUKAN error 500!**

## Default Login

Setelah setup selesai, gunakan kredensial ini untuk login:

- **Email**: admin@example.com
- **Username**: admin
- **Password**: admin123

## Troubleshooting

### Error: "psql is not recognized"

PostgreSQL belum di-install atau belum di PATH.

**Solusi:**
1. Install PostgreSQL dari https://www.postgresql.org/download/
2. Atau tambahkan ke PATH: `C:\Program Files\PostgreSQL\15\bin`

### Error: "password authentication failed"

Password salah.

**Solusi:**
- Pastikan password yang dimasukkan sesuai dengan password PostgreSQL Anda
- Atau reset password PostgreSQL

### Error: "database ganttpro_db already exists"

Database sudah ada tapi struktur salah.

**Solusi:**
```sql
-- Drop database lama
psql -U postgres
DROP DATABASE ganttpro_db;
CREATE DATABASE ganttpro_db;
\q

-- Jalankan ulang migration
psql -U postgres -d ganttpro_db -f database/migrations/00_complete_setup.sql
```

### Error: "relation ppic_schedules does not exist"

Migration belum dijalankan.

**Solusi:**
```bash
psql -U postgres -d ganttpro_db -f database/migrations/00_complete_setup.sql
```

### Backend masih error 500 setelah setup

**Solusi:**
1. Restart backend (Ctrl+C lalu `go run main.go` lagi)
2. Cek koneksi database di console output
3. Verifikasi database dengan script verify:
   ```bash
   psql -U postgres -d ganttpro_db -f database/migrations/verify_database.sql
   ```

## File-File Penting

- `database/migrations/00_complete_setup.sql` - Script migration lengkap
- `database/migrations/verify_database.sql` - Script untuk cek database
- `setup_database.ps1` - Script otomatis setup (Windows)
- `SETUP_DATABASE.md` - Dokumentasi lengkap
- `.env` - File konfigurasi (jangan di-commit ke git!)

## Struktur Database

Database akan punya 12 tables:
1. `users` - User accounts
2. `token_blacklist` - JWT blacklist
3. `machines` - Mesin produksi
4. `job_orders` - Job orders
5. `toolpather_files` - File toolpather
6. `g_code_files` - File G-code
7. `pem_operation_plans` - PEM operation plans
8. `operation_plans` - Operation plans
9. `operation_plan_approvals` - Approvals
10. `ppic_schedules` - **PPIC schedules (untuk Gantt)**
11. `machine_assignments` - **Assignment mesin (untuk Gantt)**
12. `ppic_links` - **Dependencies antar task (untuk Gantt)**

## Tips

1. **Selalu backup database** sebelum testing:
   ```bash
   pg_dump -U postgres ganttpro_db > backup.sql
   ```

2. **Restore dari backup**:
   ```bash
   psql -U postgres -d ganttpro_db < backup.sql
   ```

3. **Reset database** ke kondisi awal:
   ```bash
   psql -U postgres -d ganttpro_db -f database/migrations/00_complete_setup.sql
   ```

## Kontak Support

Jika masih ada error, screenshot error message dan kirim ke tim development.

**Penting:** Jangan share file `.env` atau password database!
