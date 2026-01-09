# 📦 Panduan Migrasi Database ke PC Baru

Panduan lengkap untuk setup database PostgreSQL di PC/laptop baru.

---

## 🎯 Prerequisites

Sebelum mulai, pastikan sudah terinstall:

1. **PostgreSQL** (versi 12 atau lebih baru)
   - Download: https://www.postgresql.org/download/windows/
   - Ingat password yang Anda set untuk user `postgres` saat instalasi!

2. **Go** (versi 1.21 atau lebih baru)
   - Download: https://go.dev/dl/

3. **Git** (untuk clone repository)
   - Download: https://git-scm.com/downloads

---

## 🚀 Cara Cepat (Otomatis)

### Windows PowerShell:

```powershell
# 1. Masuk ke folder backend
cd C:\Jemmy\compro\compro-erp-fmlx\backend

# 2. Jalankan script setup otomatis
.\setup_database_migration.ps1
```

Script akan otomatis:
- ✅ Test koneksi PostgreSQL
- ✅ Buat database `ganttpro_db`
- ✅ Jalankan migration lengkap
- ✅ Insert data default (4 user, 8 machines)
- ✅ Verifikasi instalasi

**Selesai! Langsung bisa jalankan backend.**

---

## 📝 Cara Manual

Jika script otomatis tidak berhasil, ikuti langkah manual berikut:

### 1. Buat Database Baru

Buka **pgAdmin** atau **Command Prompt** dan jalankan:

```bash
# Login ke PostgreSQL
psql -U postgres

# Di dalam psql prompt:
CREATE DATABASE ganttpro_db;

# Keluar dari psql
\q
```

### 2. Jalankan Migration Script

```bash
# Masuk ke folder migrations
cd C:\Jemmy\compro\compro-erp-fmlx\backend\database\migrations

# Jalankan migration complete
psql -U postgres -d ganttpro_db -f MIGRATION_COMPLETE.sql
```

Anda akan melihat output seperti:

```
NOTICE:  ============================================================
NOTICE:  DATABASE MIGRATION COMPLETED SUCCESSFULLY!
NOTICE:  ============================================================
NOTICE:  Total Tables Created: 12
NOTICE:  Total Users: 4
NOTICE:  Total Machines: 8
NOTICE:  Total Schedules: 3
NOTICE:  ============================================================
```

### 3. Update File `.env`

Edit file `backend/.env` (buat jika belum ada):

```env
# Server Configuration
PORT=8080

# Database Configuration
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password_anda_disini
DB_NAME=ganttpro_db

# JWT Secret (ganti dengan random string)
JWT_SECRET=your-very-secret-key-change-this-in-production

# Environment
ENV=development
```

**PENTING:** Ganti `DB_PASSWORD` dengan password PostgreSQL Anda!

### 4. Test Backend

```bash
# Masuk ke folder backend
cd C:\Jemmy\compro\compro-erp-fmlx\backend

# Download dependencies (jika belum)
go mod download

# Jalankan backend
go run main.go
```

Jika berhasil, Anda akan melihat:

```
Server running on port 8080
Database connected successfully!
```

---

## 🔐 Default User Credentials

Setelah migration, Anda bisa login dengan user berikut:

| Role     | Username  | Password     | Email                  |
|----------|-----------|--------------|------------------------|
| Admin    | admin     | admin123     | admin@ganttpro.com     |
| PPIC     | ppic      | ppic123      | ppic@ganttpro.com      |
| PEM      | pem       | pem123       | pem@ganttpro.com       |
| Operator | operator  | operator123  | operator@ganttpro.com  |

**PENTING:** Ganti password default setelah login pertama kali!

---

## 🗂️ Struktur Database

Migration akan membuat 12 tabel:

1. **users** - User authentication & roles
2. **token_blacklist** - JWT token management
3. **machines** - Daftar mesin produksi (8 sample machines)
4. **job_orders** - Job order pelanggan
5. **toolpather_files** - File CAD/CAM
6. **g_code_files** - G-Code untuk CNC
7. **pem_operation_plans** - Operation plans dari PEM
8. **operation_plans** - General operation plans
9. **operation_plan_approvals** - Workflow approval
10. **ppic_schedules** - PPIC Gantt chart schedules
11. **machine_assignments** - Assignment mesin ke schedule
12. **ppic_links** - Dependencies antar schedule

---

## ✅ Verifikasi Instalasi

### Cek di PostgreSQL:

```sql
-- Login ke database
psql -U postgres -d ganttpro_db

-- Lihat semua tabel
\dt

-- Cek jumlah users
SELECT COUNT(*) FROM users;  -- Harusnya: 4

-- Cek jumlah machines
SELECT COUNT(*) FROM machines;  -- Harusnya: 8

-- Lihat detail users
SELECT username, email, role FROM users;

-- Lihat daftar mesin
SELECT machine_code, machine_name, status FROM machines;
```

### Test API Login:

```bash
# Test dengan curl
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"admin\",\"password\":\"admin123\"}"
```

Atau buka browser:
- Frontend: http://localhost:5173
- Backend health check: http://localhost:8080/health

---

## 🛠️ Troubleshooting

### Problem 1: "psql: FATAL: password authentication failed"

**Solusi:**
```bash
# Reset password postgres
# 1. Edit pg_hba.conf (C:\Program Files\PostgreSQL\15\data\pg_hba.conf)
# 2. Ganti "md5" jadi "trust" untuk sementara
# 3. Restart PostgreSQL service
# 4. Reset password:
psql -U postgres
ALTER USER postgres PASSWORD 'password_baru';
# 5. Kembalikan pg_hba.conf ke "md5"
# 6. Restart PostgreSQL lagi
```

### Problem 2: "database ganttpro_db does not exist"

**Solusi:**
```bash
psql -U postgres
CREATE DATABASE ganttpro_db;
\q
```

### Problem 3: Backend error "failed to connect to database"

**Solusi:**
1. Cek apakah PostgreSQL service running:
   ```bash
   # Windows
   net start postgresql-x64-15
   ```

2. Cek file `.env`:
   - Pastikan `DB_PASSWORD` benar
   - Pastikan `DB_PORT=5432` (default PostgreSQL)
   - Pastikan `DB_HOST=localhost`

3. Test koneksi manual:
   ```bash
   psql -U postgres -h localhost -p 5432 -d ganttpro_db
   ```

### Problem 4: "listen tcp :8080: bind: address already in use"

**Solusi:**
```bash
# Cek proses yang pakai port 8080
netstat -ano | findstr :8080

# Kill proses tersebut
taskkill /PID <process_id> /F

# Atau ganti PORT di .env
PORT=8081
```

---

## 📋 Checklist Migrasi

Gunakan checklist ini saat setup di PC baru:

- [ ] Install PostgreSQL
- [ ] Install Go
- [ ] Clone repository
- [ ] Buat database `ganttpro_db`
- [ ] Jalankan `MIGRATION_COMPLETE.sql`
- [ ] Copy & edit file `.env`
- [ ] Update `DB_PASSWORD` di `.env`
- [ ] Test: `go run main.go`
- [ ] Verifikasi: Buka http://localhost:8080/health
- [ ] Test login dengan user admin
- [ ] Ganti password default

---

## 📞 Bantuan Lebih Lanjut

Jika masih ada masalah:

1. Cek log backend untuk error details
2. Cek PostgreSQL log di: `C:\Program Files\PostgreSQL\15\data\log`
3. Pastikan firewall tidak block port 8080 dan 5432

---

## 🔄 Reset Database

Jika ingin reset database dari awal:

```sql
-- HATI-HATI: Ini akan hapus semua data!
psql -U postgres

DROP DATABASE ganttpro_db;
CREATE DATABASE ganttpro_db;
\q

-- Jalankan ulang migration
psql -U postgres -d ganttpro_db -f MIGRATION_COMPLETE.sql
```

---

## 📄 File-file Penting

- `MIGRATION_COMPLETE.sql` - Script migration lengkap
- `setup_database_migration.ps1` - Script otomatis Windows
- `backend/.env` - Konfigurasi backend
- `backend/README.md` - Dokumentasi backend

---

**✨ Selamat! Database siap digunakan di PC baru Anda!**
