# ✅ Checklist Migrasi Database ke PC Baru

Gunakan checklist ini untuk memastikan semua langkah sudah dilakukan dengan benar.

---

## 📋 Pre-Installation

- [ ] PostgreSQL terinstall (versi 12+)
- [ ] Go terinstall (versi 1.21+)
- [ ] Git terinstall
- [ ] Repository sudah di-clone ke PC baru
- [ ] Catat password PostgreSQL yang akan digunakan

---

## 🚀 Database Setup

- [ ] Buka PowerShell/Command Prompt sebagai Administrator
- [ ] Masuk ke folder backend: `cd backend`
- [ ] Jalankan script setup:
  - [ ] **Windows PowerShell**: `.\setup_database_migration.ps1`
  - [ ] **Command Prompt**: `setup_database_migration.bat`
- [ ] Masukkan password PostgreSQL saat diminta
- [ ] Verifikasi output menunjukkan:
  - [ ] "DATABASE MIGRATION COMPLETED SUCCESSFULLY!"
  - [ ] Total Tables: 12
  - [ ] Total Users: 4
  - [ ] Total Machines: 8

**Atau Manual:**

- [ ] Buat database: `CREATE DATABASE ganttpro_db;`
- [ ] Jalankan migration: `psql -U postgres -d ganttpro_db -f database\migrations\MIGRATION_COMPLETE.sql`
- [ ] Verifikasi dengan query: `SELECT COUNT(*) FROM users;`

---

## ⚙️ Backend Configuration

- [ ] File `.env` sudah dibuat di folder `backend/`
- [ ] Update konfigurasi di `.env`:
  - [ ] `DB_PASSWORD` sesuai dengan PostgreSQL
  - [ ] `DB_HOST=localhost` (atau sesuai kebutuhan)
  - [ ] `DB_PORT=5432` (default PostgreSQL)
  - [ ] `DB_NAME=ganttpro_db`
  - [ ] `JWT_SECRET` diganti dengan random string
  - [ ] `PORT=8080` (atau port yang tersedia)

---

## 🧪 Testing Backend

- [ ] Download dependencies: `go mod download`
- [ ] Build backend (opsional): `go build -o ganttpro-backend.exe`
- [ ] Jalankan backend: `go run main.go`
- [ ] Verifikasi output:
  - [ ] "Server running on port 8080"
  - [ ] "Database connected successfully!"
  - [ ] Tidak ada error koneksi database
- [ ] Test health check: Buka http://localhost:8080/health
  - [ ] Response: `{"status":"ok"}`

---

## 🔐 Test Authentication

- [ ] Test login dengan curl atau Postman
- [ ] Test dengan user admin:
  ```bash
  curl -X POST http://localhost:8080/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"admin\",\"password\":\"admin123\"}"
  ```
- [ ] Response berisi token JWT
- [ ] Response berisi data user (id, username, role)

**Default Credentials:**

- [ ] Admin: `admin` / `admin123` ✓
- [ ] PPIC: `ppic` / `ppic123` ✓
- [ ] PEM: `pem` / `pem123` ✓
- [ ] Operator: `operator` / `operator123` ✓

---

## 🎨 Frontend Setup

- [ ] Masuk ke folder frontend: `cd ../frontend`
- [ ] Install dependencies: `npm install`
- [ ] Update file `.env` di frontend:
  - [ ] `VITE_API_BASE_URL=http://localhost:8080/api/v1`
  - [ ] Atau sesuaikan dengan IP jika network mode
- [ ] Jalankan frontend: `npm run dev`
- [ ] Buka browser: http://localhost:5173
- [ ] Test login dengan user admin

---

## 🔧 Verifikasi Database

Di PostgreSQL (psql atau pgAdmin):

- [ ] Koneksi ke database: `\c ganttpro_db`
- [ ] Lihat semua tabel: `\dt`
- [ ] Cek data users:
  ```sql
  SELECT id, username, email, role FROM users;
  ```
  Harusnya ada 4 users

- [ ] Cek data machines:
  ```sql
  SELECT id, machine_code, machine_name, status FROM machines;
  ```
  Harusnya ada 8 machines

- [ ] Cek struktur tabel ppic_schedules:
  ```sql
  \d ppic_schedules
  ```

---

## 🛡️ Security Checklist

- [ ] Password PostgreSQL cukup kuat
- [ ] JWT_SECRET sudah diganti dari default
- [ ] File `.env` TIDAK di-commit ke git (ada di `.gitignore`)
- [ ] Ganti password default user setelah login pertama:
  - [ ] Admin
  - [ ] PPIC
  - [ ] PEM
  - [ ] Operator

---

## 📁 File Structure Check

Pastikan file-file berikut ada:

```
compro-erp-fmlx/
├── backend/
│   ├── .env                                    ✓
│   ├── go.mod                                  ✓
│   ├── main.go                                 ✓
│   ├── setup_database_migration.ps1            ✓
│   ├── setup_database_migration.bat            ✓
│   ├── MIGRATION_CHECKLIST.md                  ✓
│   └── database/
│       └── migrations/
│           └── MIGRATION_COMPLETE.sql          ✓
├── frontend/
│   ├── .env                                    ✓
│   ├── package.json                            ✓
│   └── vite.config.js                          ✓
└── PANDUAN_MIGRASI_DATABASE.md                 ✓
```

---

## 🔄 Network/Remote Access (Opsional)

Jika ingin akses dari PC lain dalam jaringan:

- [ ] Update `DB_HOST` di backend `.env` ke IP server database
- [ ] Update PostgreSQL `pg_hba.conf` untuk allow remote connections
- [ ] Update PostgreSQL `postgresql.conf`: `listen_addresses = '*'`
- [ ] Restart PostgreSQL service
- [ ] Buka firewall untuk port 5432 (PostgreSQL) dan 8080 (Backend)
- [ ] Update frontend `.env`: `VITE_API_BASE_URL=http://<backend-ip>:8080/api/v1`

---

## ✅ Final Verification

- [ ] Backend berjalan tanpa error
- [ ] Frontend berjalan tanpa error
- [ ] Login berhasil dengan user admin
- [ ] Bisa akses halaman dashboard
- [ ] Bisa lihat data PPIC schedules (jika ada)
- [ ] Bisa lihat daftar machines
- [ ] API response time normal (< 1 detik)

---

## 🆘 Troubleshooting

Jika ada masalah, cek:

- [ ] PostgreSQL service sedang running
  - Windows: `net start postgresql-x64-15`
  - Services: Cari "postgresql" dan pastikan "Running"

- [ ] Port tidak bentrok:
  - PostgreSQL (5432): `netstat -ano | findstr :5432`
  - Backend (8080): `netstat -ano | findstr :8080`

- [ ] Firewall tidak blocking:
  - Windows Firewall → Allow PostgreSQL & Go

- [ ] Credentials benar:
  - Test manual: `psql -U postgres -h localhost -d ganttpro_db`

- [ ] File `.env` terbaca:
  - Print env di backend untuk debug

---

## 📝 Notes

- **Backup Database**: Sebelum migration production, backup database lama:
  ```bash
  pg_dump -U postgres -d ganttpro_db > backup_$(date +%Y%m%d).sql
  ```

- **Restore dari Backup**:
  ```bash
  psql -U postgres -d ganttpro_db < backup_20260109.sql
  ```

- **Reset Database** (hati-hati!):
  ```sql
  DROP DATABASE ganttpro_db;
  CREATE DATABASE ganttpro_db;
  ```
  Lalu jalankan ulang migration.

---

## ✨ Success!

Jika semua checklist tercentang, database berhasil di-migrate ke PC baru!

**Selamat menggunakan GanttPro ERP! 🎉**

---

**Last Updated**: 2026-01-09
**Migration Version**: 1.0
