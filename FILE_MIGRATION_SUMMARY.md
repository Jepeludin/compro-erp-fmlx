# 📦 Summary - File Migration Database

File-file yang telah dibuat untuk memudahkan migrasi database ke PC baru.

---

## ✅ File yang Berhasil Dibuat

### 📁 Root Directory

```
compro-erp-fmlx/
├── PANDUAN_MIGRASI_DATABASE.md        ✅ Panduan lengkap migrasi (Bahasa Indonesia)
├── QUICK_MIGRATION_GUIDE.md           ✅ Quick reference 5 menit
└── FILE_MIGRATION_SUMMARY.md          ✅ File ini
```

### 📁 Backend Directory

```
backend/
├── setup_database_migration.ps1       ✅ Script otomatis PowerShell (UTAMA)
├── setup_database_migration.bat       ✅ Wrapper untuk Command Prompt
├── MIGRATION_CHECKLIST.md             ✅ Checklist step-by-step
└── .env.example                       ✅ Template konfigurasi (sudah ada)
```

### 📁 Backend/Database/Migrations Directory

```
backend/database/migrations/
├── MIGRATION_COMPLETE.sql             ✅ Setup database lengkap (UTAMA)
├── VERIFY_MIGRATION.sql               ✅ Verifikasi hasil migration
├── EXPORT_DATA.sql                    ✅ Export data dari database lama
├── BACKUP_DATABASE.bat                ✅ Script backup otomatis
├── RESTORE_DATABASE.bat               ✅ Script restore otomatis
└── README.md                          ✅ Dokumentasi migration files
```

---

## 🎯 Cara Menggunakan (Quick Start)

### Option 1: Script Otomatis (RECOMMENDED) ⭐

```powershell
# 1. Buka PowerShell di folder backend
cd backend

# 2. Jalankan script
.\setup_database_migration.ps1

# 3. Masukkan password PostgreSQL
# 4. Done! Database siap digunakan
```

**Yang Dilakukan Script:**
- ✅ Cek prerequisites (PostgreSQL, Go)
- ✅ Test koneksi database
- ✅ Buat database `ganttpro_db`
- ✅ Jalankan migration lengkap (12 tabel)
- ✅ Insert default data (4 users, 8 machines)
- ✅ Buat/update file `.env`
- ✅ Verifikasi instalasi

### Option 2: Manual

```bash
# 1. Buat database
psql -U postgres
CREATE DATABASE ganttpro_db;
\q

# 2. Jalankan migration
cd backend/database/migrations
psql -U postgres -d ganttpro_db -f MIGRATION_COMPLETE.sql

# 3. Verifikasi
psql -U postgres -d ganttpro_db -f VERIFY_MIGRATION.sql

# 4. Edit .env
cd ../..
notepad .env
# Update DB_PASSWORD
```

---

## 📋 Struktur Database yang Dibuat

### 12 Tabel Utama

1. **users** - User authentication (4 default users)
2. **token_blacklist** - JWT token management
3. **machines** - Daftar mesin produksi (8 default machines)
4. **job_orders** - Job order pelanggan
5. **toolpather_files** - File CAD/CAM
6. **g_code_files** - G-Code untuk CNC
7. **pem_operation_plans** - PEM operation plans
8. **operation_plans** - General operation plans
9. **operation_plan_approvals** - Approval workflow
10. **ppic_schedules** - PPIC Gantt schedules
11. **machine_assignments** - Machine assignments
12. **ppic_links** - Schedule dependencies

### Default Data

**4 Users:**
| Username | Password | Role |
|----------|----------|------|
| admin | admin123 | admin |
| ppic | ppic123 | ppic |
| pem | pem123 | pem |
| operator | operator123 | operator |

**8 Machines:**
- M01-CNC: CNC Milling Machine 1
- M02-CNC: CNC Milling Machine 2
- M03-LATHE: CNC Lathe 1
- M04-GRIND: Surface Grinder 1
- M05-DRILL: Drilling Machine 1
- M06-EDM: EDM Machine 1
- M07-MILL: Manual Milling 1
- M08-WELD: Welding Station 1

---

## 🔄 Workflow Berbagai Skenario

### Skenario 1: Setup PC Baru (Fresh Install)

```
Step 1: Install PostgreSQL + Go
   ↓
Step 2: Clone repository
   ↓
Step 3: .\setup_database_migration.ps1
   ↓
Step 4: go run main.go
   ↓
Done! ✅
```

### Skenario 2: Pindah PC dengan Data Lama

**Di PC Lama:**
```
Step 1: BACKUP_DATABASE.bat
   ↓
Step 2: Copy backup.sql ke USB/Cloud
```

**Di PC Baru:**
```
Step 3: .\setup_database_migration.ps1 (fresh database)
   ↓
Step 4: RESTORE_DATABASE.bat
   ↓
Step 5: Pilih file backup.sql
   ↓
Done! ✅
```

### Skenario 3: Reset Database

```
Step 1: BACKUP_DATABASE.bat (safety!)
   ↓
Step 2: .\setup_database_migration.ps1
   ↓
Step 3: Pilih "y" untuk drop & recreate
   ↓
Done! ✅
```

---

## 📊 File Reference Chart

| Kebutuhan | File yang Digunakan | Lokasi |
|-----------|---------------------|--------|
| Setup database baru | setup_database_migration.ps1 | backend/ |
| Setup manual | MIGRATION_COMPLETE.sql | backend/database/migrations/ |
| Verifikasi setup | VERIFY_MIGRATION.sql | backend/database/migrations/ |
| Backup database | BACKUP_DATABASE.bat | backend/database/migrations/ |
| Restore database | RESTORE_DATABASE.bat | backend/database/migrations/ |
| Export data saja | EXPORT_DATA.sql | backend/database/migrations/ |
| Panduan lengkap | PANDUAN_MIGRASI_DATABASE.md | root/ |
| Quick reference | QUICK_MIGRATION_GUIDE.md | root/ |
| Checklist | MIGRATION_CHECKLIST.md | backend/ |

---

## ✅ Verification Checklist

Setelah migration, pastikan:

- [ ] Script `setup_database_migration.ps1` selesai tanpa error
- [ ] Output menunjukkan "MIGRATION COMPLETED SUCCESSFULLY"
- [ ] Total Tables: 12
- [ ] Total Users: 4
- [ ] Total Machines: 8
- [ ] File `.env` ada di folder `backend/`
- [ ] Backend bisa jalan: `go run main.go`
- [ ] Health check OK: `curl http://localhost:8080/health`
- [ ] Login berhasil dengan `admin` / `admin123`

---

## 🛠️ Tools & Scripts Summary

### PowerShell Scripts

| Script | Fungsi | Auto/Manual |
|--------|--------|-------------|
| setup_database_migration.ps1 | Setup lengkap otomatis | Auto |

### Batch Scripts

| Script | Fungsi | Auto/Manual |
|--------|--------|-------------|
| setup_database_migration.bat | Wrapper PS1 script | Auto |
| BACKUP_DATABASE.bat | Backup database | Interactive |
| RESTORE_DATABASE.bat | Restore database | Interactive |

### SQL Scripts

| Script | Fungsi | Output |
|--------|--------|--------|
| MIGRATION_COMPLETE.sql | Create all tables + data | 12 tables, 4 users, 8 machines |
| VERIFY_MIGRATION.sql | Verification report | Detailed report |
| EXPORT_DATA.sql | Export to SQL | SQL INSERT statements |

---

## 📝 Configuration Files

### backend/.env

Template:
```env
PORT=8080
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=ganttpro_db
JWT_SECRET=change-this-secret
```

**PENTING:**
- Ganti `DB_PASSWORD` dengan password PostgreSQL Anda
- Ganti `JWT_SECRET` dengan random string

---

## 🔐 Security Notes

1. **Password Default:** Ganti semua password default setelah login!
2. **JWT Secret:** Gunakan string random yang kuat untuk production
3. **File .env:** Jangan commit ke git (sudah ada di .gitignore)
4. **Backup:** Simpan backup di tempat aman
5. **Database Access:** Batasi akses remote jika tidak diperlukan

---

## 🆘 Troubleshooting

### Error: "psql: command not found"
**Solusi:** Tambahkan PostgreSQL ke PATH
```
C:\Program Files\PostgreSQL\15\bin
```

### Error: "password authentication failed"
**Solusi:**
- Cek password PostgreSQL benar
- Test manual: `psql -U postgres`

### Error: "port 8080 already in use"
**Solusi:**
- Ganti PORT di .env: `PORT=8081`
- Atau kill process: `taskkill /F /PID <pid>`

### Error: "database already exists"
**Solusi:**
- Pilih "y" saat script tanya drop database
- Atau manual: `DROP DATABASE ganttpro_db;`

---

## 📞 Support & Documentation

### Dokumentasi Lengkap

1. **[PANDUAN_MIGRASI_DATABASE.md](PANDUAN_MIGRASI_DATABASE.md)**
   - Panduan lengkap Bahasa Indonesia
   - Step-by-step dengan screenshots concept
   - Troubleshooting detail

2. **[QUICK_MIGRATION_GUIDE.md](QUICK_MIGRATION_GUIDE.md)**
   - Quick reference 5 menit
   - Command cheat sheet
   - Fast troubleshooting

3. **[backend/MIGRATION_CHECKLIST.md](backend/MIGRATION_CHECKLIST.md)**
   - Checklist lengkap
   - Verification steps
   - Security checklist

4. **[backend/database/migrations/README.md](backend/database/migrations/README.md)**
   - File overview
   - Usage matrix
   - Workflow diagrams

---

## 🎉 Success Criteria

Migration berhasil jika:

✅ Database `ganttpro_db` terbuat
✅ 12 tabel terinstall dengan benar
✅ 4 default users ada
✅ 8 default machines ada
✅ Backend berjalan tanpa error
✅ API health check response OK
✅ Login berhasil dengan user admin
✅ Frontend bisa connect ke backend

---

## 📈 Next Steps

Setelah migration berhasil:

1. **Test Backend:**
   ```bash
   cd backend
   go run main.go
   ```

2. **Test API:**
   ```bash
   curl http://localhost:8080/health
   ```

3. **Setup Frontend:**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

4. **Test Login:**
   - Buka http://localhost:5173
   - Login: admin / admin123
   - Ganti password default

5. **Production Ready:**
   - Ganti semua password
   - Update JWT_SECRET
   - Setup backup regular
   - Configure firewall

---

## 📅 Maintenance

### Backup Regular

```bash
# Setiap minggu/bulan
cd backend/database/migrations
BACKUP_DATABASE.bat
```

### Update Schema

Jika ada perubahan schema:
1. Backup database dulu
2. Buat file migration baru: `08_new_feature.sql`
3. Test di development
4. Apply di production

---

## 🏆 Summary

**Total File Created:** 9 files

**Dokumentasi:** 4 files
- PANDUAN_MIGRASI_DATABASE.md
- QUICK_MIGRATION_GUIDE.md
- MIGRATION_CHECKLIST.md
- migrations/README.md

**Scripts:** 5 files
- setup_database_migration.ps1
- setup_database_migration.bat
- BACKUP_DATABASE.bat
- RESTORE_DATABASE.bat
- MIGRATION_COMPLETE.sql
- VERIFY_MIGRATION.sql
- EXPORT_DATA.sql

**Waktu Setup:** ~5 menit (dengan script otomatis)

**Kompleksitas:** ⭐⭐☆☆☆ (Easy-Medium)

---

**✨ Database siap digunakan di PC baru! Selamat bekerja dengan GanttPro ERP! 🎉**

---

**Last Updated:** 2026-01-09
**Version:** 1.0
**Status:** ✅ Production Ready
