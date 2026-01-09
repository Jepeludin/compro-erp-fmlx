# 📚 Database Migration Files - GanttPro ERP

Koleksi lengkap file migration dan tools untuk setup database PostgreSQL.

---

## 📁 File Overview

### 🎯 Migration Scripts (SQL)

| File | Deskripsi | Kapan Digunakan |
|------|-----------|-----------------|
| **MIGRATION_COMPLETE.sql** | ⭐ Setup database lengkap dari nol | **PC baru / fresh install** |
| VERIFY_MIGRATION.sql | Verifikasi bahwa migration berhasil | Setelah migration untuk cek |
| EXPORT_DATA.sql | Export data dari database lama | Sebelum pindah PC (backup data) |
| verify_database.sql | Quick verification | Legacy, pakai VERIFY_MIGRATION.sql |

### 🔧 Automation Scripts (Windows)

| File | Deskripsi | Kapan Digunakan |
|------|-----------|-----------------|
| **setup_database_migration.ps1** | ⭐ Script setup otomatis PowerShell | **Cara tercepat setup database** |
| setup_database_migration.bat | Wrapper untuk PS1 script | Untuk Command Prompt |
| BACKUP_DATABASE.bat | Backup database otomatis | Sebelum update/migrasi |
| RESTORE_DATABASE.bat | Restore dari backup | Recovery atau pindah PC |

### 📖 Documentation

| File | Deskripsi |
|------|-----------|
| MIGRATION_CHECKLIST.md | Checklist lengkap step-by-step |
| README.md | File ini |

### 🗂️ Legacy Files

| File | Status | Note |
|------|--------|------|
| 00_complete_setup.sql | Legacy | Gunakan MIGRATION_COMPLETE.sql |
| 01-07_*.sql | Legacy | Script lama, sudah digabung |
| RUN_ME_FIRST.sql | Legacy | Gunakan setup_database_migration.ps1 |

---

## 🚀 Quick Start

### Cara Paling Cepat (Recommended)

```powershell
# 1. Buka PowerShell di folder backend
cd C:\path\to\compro-erp-fmlx\backend

# 2. Jalankan script otomatis
.\setup_database_migration.ps1

# 3. Ikuti instruksi di layar
# Done! Database ready.
```

### Manual Way

```bash
# 1. Buat database
psql -U postgres
CREATE DATABASE ganttpro_db;
\q

# 2. Jalankan migration
cd database/migrations
psql -U postgres -d ganttpro_db -f MIGRATION_COMPLETE.sql

# 3. Verifikasi
psql -U postgres -d ganttpro_db -f VERIFY_MIGRATION.sql
```

---

## 📋 Workflow Migrasi

### Scenario 1: Fresh Install (PC Baru, Database Kosong)

```
1. setup_database_migration.ps1
   ↓
2. MIGRATION_COMPLETE.sql (auto)
   ↓
3. VERIFY_MIGRATION.sql (optional)
   ↓
4. Done! ✅
```

### Scenario 2: Pindah PC dengan Data Lama

**Di PC Lama:**
```
1. BACKUP_DATABASE.bat
   ↓
2. Copy file backup.sql ke PC baru
```

**Di PC Baru:**
```
3. setup_database_migration.ps1 (skip jika database sudah ada)
   ↓
4. RESTORE_DATABASE.bat
   ↓
5. Pilih file backup.sql
   ↓
6. VERIFY_MIGRATION.sql
   ↓
7. Done! ✅
```

### Scenario 3: Update/Refresh Database

```
1. BACKUP_DATABASE.bat (safety first!)
   ↓
2. setup_database_migration.ps1
   ↓
3. Pilih "y" untuk drop & recreate
   ↓
4. Done! ✅
```

---

## 🗄️ Database Schema

Migration akan membuat 12 tabel:

```
ganttpro_db
├── users                      # User & authentication
├── token_blacklist            # JWT token management
├── machines                   # Daftar mesin produksi
├── job_orders                 # Job order pelanggan
├── toolpather_files           # CAD/CAM files
├── g_code_files               # G-Code untuk CNC
├── pem_operation_plans        # PEM operation plans
├── operation_plans            # General operation plans
├── operation_plan_approvals   # Approval workflow
├── ppic_schedules             # PPIC Gantt schedules
├── machine_assignments        # Machine-schedule assignments
└── ppic_links                 # Schedule dependencies
```

### Default Data

**Users (4):**
- admin / admin123 (role: admin)
- ppic / ppic123 (role: ppic)
- pem / pem123 (role: pem)
- operator / operator123 (role: operator)

**Machines (8):**
- M01-CNC: CNC Milling Machine 1
- M02-CNC: CNC Milling Machine 2
- M03-LATHE: CNC Lathe 1
- M04-GRIND: Surface Grinder 1
- M05-DRILL: Drilling Machine 1
- M06-EDM: EDM Machine 1
- M07-MILL: Manual Milling 1
- M08-WELD: Welding Station 1

---

## 🔍 Verification

Setelah migration, jalankan verifikasi:

```bash
# SQL Verification
psql -U postgres -d ganttpro_db -f VERIFY_MIGRATION.sql

# Backend Verification
cd ../..
go run main.go
# Output: "Server running on port 8080"

# API Verification
curl http://localhost:8080/health
# Output: {"status":"ok"}
```

---

## 🛠️ Troubleshooting

### Problem: "psql: command not found"

**Solution:**
```
Tambahkan PostgreSQL ke PATH:
C:\Program Files\PostgreSQL\15\bin
```

### Problem: "password authentication failed"

**Solution:**
```
1. Cek password PostgreSQL
2. Edit pg_hba.conf jika perlu
3. Restart PostgreSQL service
```

### Problem: "database already exists"

**Solution:**
```
# Option 1: Drop & recreate
psql -U postgres
DROP DATABASE ganttpro_db;
CREATE DATABASE ganttpro_db;

# Option 2: Use script
.\setup_database_migration.ps1
# Pilih "y" saat ditanya drop database
```

### Problem: Migration errors

**Solution:**
```bash
# Cek log detail
psql -U postgres -d ganttpro_db -f MIGRATION_COMPLETE.sql 2>&1 | tee migration.log

# Cek tabel yang sudah dibuat
psql -U postgres -d ganttpro_db
\dt

# Drop semua tabel dan coba lagi
# Uncomment DROP TABLE di MIGRATION_COMPLETE.sql
```

---

## 📊 File Usage Matrix

| Use Case | Primary File | Secondary File |
|----------|--------------|----------------|
| Setup PC baru | setup_database_migration.ps1 | MIGRATION_COMPLETE.sql |
| Verifikasi setup | VERIFY_MIGRATION.sql | - |
| Backup data | BACKUP_DATABASE.bat | - |
| Restore data | RESTORE_DATABASE.bat | backup.sql |
| Export data saja | EXPORT_DATA.sql | - |
| Manual setup | MIGRATION_COMPLETE.sql | VERIFY_MIGRATION.sql |

---

## 🔗 Related Documentation

- [PANDUAN_MIGRASI_DATABASE.md](../../../PANDUAN_MIGRASI_DATABASE.md) - Panduan lengkap Bahasa Indonesia
- [QUICK_MIGRATION_GUIDE.md](../../../QUICK_MIGRATION_GUIDE.md) - Quick reference
- [MIGRATION_CHECKLIST.md](../../MIGRATION_CHECKLIST.md) - Step-by-step checklist
- [backend/README.md](../../README.md) - Backend documentation

---

## 📝 Notes

1. **Password Security**: Ganti semua password default setelah login!
2. **Backup Regular**: Jalankan BACKUP_DATABASE.bat secara berkala
3. **Version Control**: Jangan commit file `.env` atau `backup.sql` ke git
4. **Production**: Gunakan password & JWT secret yang kuat
5. **Testing**: Selalu test di environment development dulu

---

## 📞 Support

Jika masih ada masalah:
1. Cek dokumentasi lengkap di root folder
2. Review error logs PostgreSQL
3. Pastikan semua prerequisites terinstall

---

**Last Updated**: 2026-01-09
**Version**: 1.0
**PostgreSQL**: 12+
**Status**: ✅ Production Ready
