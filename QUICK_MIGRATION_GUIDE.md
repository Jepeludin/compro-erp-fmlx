# ⚡ Quick Migration Guide - GanttPro ERP

Panduan super cepat untuk migrasi database ke PC baru dalam 5 menit!

---

## 🎯 Untuk PC Baru (Fresh Install)

### Option 1: Script Otomatis (RECOMMENDED)

```bash
# 1. Clone repository
git clone <repository-url>
cd compro-erp-fmlx/backend

# 2. Run setup script
.\setup_database_migration.ps1

# 3. Jalankan backend
go run main.go
```

**Done! Database siap digunakan dengan:**
- ✅ 12 tabel database
- ✅ 4 default users (admin, ppic, pem, operator)
- ✅ 8 sample machines

### Option 2: Manual (5 Langkah)

```bash
# 1. Buat database
psql -U postgres
CREATE DATABASE ganttpro_db;
\q

# 2. Jalankan migration
cd backend/database/migrations
psql -U postgres -d ganttpro_db -f MIGRATION_COMPLETE.sql

# 3. Edit .env
cd ../..
copy .env.example .env
# Edit DB_PASSWORD di .env

# 4. Test backend
go run main.go

# 5. Login
Username: admin
Password: admin123
```

---

## 🔄 Untuk Migrasi Data dari PC Lama

### Di PC Lama:

```bash
# Backup database
cd backend/database/migrations
BACKUP_DATABASE.bat

# Atau manual:
pg_dump -U postgres ganttpro_db > backup.sql
```

### Di PC Baru:

```bash
# 1. Setup database kosong (pakai script di atas)

# 2. Restore backup
cd backend/database/migrations
RESTORE_DATABASE.bat
# Pilih file backup

# Atau manual:
psql -U postgres -d ganttpro_db < backup.sql
```

---

## 📁 File-file Penting

| File | Fungsi |
|------|--------|
| `MIGRATION_COMPLETE.sql` | Setup database lengkap dari nol |
| `setup_database_migration.ps1` | Script otomatis Windows |
| `BACKUP_DATABASE.bat` | Backup database |
| `RESTORE_DATABASE.bat` | Restore dari backup |
| `MIGRATION_CHECKLIST.md` | Checklist lengkap |
| `PANDUAN_MIGRASI_DATABASE.md` | Dokumentasi detail |

---

## 🔐 Default Credentials

| Role | Username | Password |
|------|----------|----------|
| Admin | admin | admin123 |
| PPIC | ppic | ppic123 |
| PEM | pem | pem123 |
| Operator | operator | operator123 |

**⚠️ Ganti password setelah login pertama!**

---

## 🆘 Troubleshooting Cepat

**Error: "password authentication failed"**
```bash
# Cek password PostgreSQL benar
psql -U postgres
```

**Error: "database does not exist"**
```bash
# Buat database dulu
psql -U postgres
CREATE DATABASE ganttpro_db;
```

**Error: "port 8080 already in use"**
```bash
# Ganti PORT di .env
PORT=8081
```

**Backend tidak konek ke database**
```bash
# Cek .env:
# - DB_PASSWORD benar
# - DB_PORT=5432 (atau sesuai PostgreSQL Anda)
# - DB_HOST=localhost
```

---

## ✅ Verifikasi

```bash
# 1. Backend running
go run main.go
# Output: "Server running on port 8080"

# 2. Health check
curl http://localhost:8080/health
# Output: {"status":"ok"}

# 3. Test login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
# Output: token JWT
```

---

## 📞 Butuh Bantuan?

- 📖 Dokumentasi lengkap: [PANDUAN_MIGRASI_DATABASE.md](PANDUAN_MIGRASI_DATABASE.md)
- ✅ Checklist detail: [backend/MIGRATION_CHECKLIST.md](backend/MIGRATION_CHECKLIST.md)
- 🔧 Backend README: [backend/README.md](backend/README.md)

---

**✨ Selamat! Database siap digunakan!**

Waktu setup: ~5 menit | Kompleksitas: ⭐⭐☆☆☆
