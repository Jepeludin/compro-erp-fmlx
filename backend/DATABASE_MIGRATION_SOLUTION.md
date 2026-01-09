# Perubahan dan Solusi untuk Error Database

## Problem
Error yang terjadi di laptop:
```
ERROR: kolom ms.scheduled_start_belum ada (SQLSTATE 42703)
```

Error ini terjadi karena database di laptop belum di-setup dengan benar setelah pindah dari PC.

## Root Cause
- Database di laptop tidak memiliki table `machine_assignments` dengan kolom `scheduled_start`
- Migration belum dijalankan di laptop baru
- Struktur database tidak lengkap

## Solusi yang Sudah Dibuat

### 1. Complete Migration Script
**File:** `database/migrations/00_complete_setup.sql`

Script SQL lengkap yang membuat semua table dari nol:
- ✅ 12 tables (users, machines, ppic_schedules, machine_assignments, dll)
- ✅ All indexes dan foreign keys
- ✅ Default admin user (admin@example.com / admin123)
- ✅ 5 sample machines (CNC-001, CNC-002, MILL-001, LATHE-001, GRIND-001)
- ✅ Proper column definitions termasuk `scheduled_start`, `scheduled_end`, dll

### 2. PowerShell Setup Script
**File:** `setup_database.ps1`

Script otomatis untuk Windows yang:
- ✅ Test koneksi PostgreSQL
- ✅ Buat database baru (atau drop & recreate)
- ✅ Jalankan migration otomatis
- ✅ Verifikasi setup berhasil
- ✅ Generate file .env otomatis
- ✅ User-friendly dengan progress indicator

**Cara pakai:**
```powershell
cd C:\Jemmy\compro\compro-erp-fmlx\backend
.\setup_database.ps1
```

### 3. Database Verification Script
**File:** `database/migrations/verify_database.sql`

Script untuk verifikasi database sudah setup dengan benar:
- ✅ Check semua table ada
- ✅ Check struktur table `machine_assignments`
- ✅ Verify critical columns (scheduled_start, scheduled_end, dll)
- ✅ Check indexes dan foreign keys
- ✅ Test query yang sama seperti aplikasi
- ✅ Count data (users, machines, schedules)

**Cara pakai:**
```bash
psql -U postgres -d ganttpro_db -f database/migrations/verify_database.sql
```

### 4. Quick Start Guide
**File:** `QUICK_START.md`

Panduan singkat dalam Bahasa Indonesia untuk:
- Setup database di laptop baru (< 5 menit)
- Troubleshooting common errors
- Default credentials
- Tips & tricks

### 5. Complete Setup Documentation
**File:** `SETUP_DATABASE.md`

Dokumentasi lengkap tentang:
- Step-by-step manual setup
- Cara menggunakan psql dan pgAdmin
- Penjelasan struktur database
- Troubleshooting lengkap
- Reset database procedures

### 6. Updated README
**File:** `README.md`

Update README dengan:
- Link ke quick start guide
- Penjelasan struktur project
- Setup instructions yang jelas

## Cara Menggunakan

### Opsi 1: Setup Otomatis (RECOMMENDED)

1. Buka PowerShell
2. Navigate ke backend folder:
   ```powershell
   cd C:\Jemmy\compro\compro-erp-fmlx\backend
   ```
3. Jalankan setup script:
   ```powershell
   .\setup_database.ps1
   ```
4. Ikuti instruksi (masukkan username & password PostgreSQL)
5. Selesai! Jalankan backend:
   ```powershell
   go run main.go
   ```

### Opsi 2: Setup Manual

1. Create database:
   ```sql
   psql -U postgres
   CREATE DATABASE ganttpro_db;
   \q
   ```

2. Run migration:
   ```bash
   psql -U postgres -d ganttpro_db -f database/migrations/00_complete_setup.sql
   ```

3. Verify:
   ```bash
   psql -U postgres -d ganttpro_db -f database/migrations/verify_database.sql
   ```

4. Update .env dengan credentials yang benar

5. Run backend:
   ```bash
   go run main.go
   ```

## Verifikasi Setup Berhasil

1. Backend harus start tanpa error:
   ```
   Connected to database successfully!
   Server is running on port 8080
   ```

2. Test API harus return data (bukan error 500):
   ```
   http://localhost:8080/api/ppic/gantt-data
   ```
   
   Response:
   ```json
   {
     "success": true,
     "data": {
       "sections": [],
       "links": []
     }
   }
   ```

## Default Credentials

Setelah setup selesai:
- **Email:** admin@example.com
- **Username:** admin
- **Password:** admin123

⚠️ **PENTING:** Ganti password default setelah login pertama kali!

## File Structure

```
backend/
├── database/
│   └── migrations/
│       ├── 00_complete_setup.sql        # ← MAIN MIGRATION (RUN THIS)
│       ├── verify_database.sql          # ← VERIFICATION SCRIPT
│       ├── 001_create_ppic_tables.sql   # (old, not used)
│       ├── migrationppic.sql            # (old, not used)
│       └── ppicjob.sql                  # (old, not used)
├── setup_database.ps1                   # ← AUTO SETUP (WINDOWS)
├── QUICK_START.md                       # ← READ THIS FIRST!
├── SETUP_DATABASE.md                    # ← COMPLETE DOCS
└── README.md                            # ← UPDATED
```

## Troubleshooting

### Error: psql not found
- Install PostgreSQL atau tambahkan ke PATH

### Error: password authentication failed
- Check password di .env atau saat run setup script

### Error: database already exists
- Pilih "yes" saat setup script tanya untuk drop & recreate
- Atau manual: `DROP DATABASE ganttpro_db; CREATE DATABASE ganttpro_db;`

### Error 500 setelah setup
- Restart backend
- Verify database dengan verify_database.sql
- Check koneksi di .env

## Testing

Setelah setup, test endpoints ini:

1. **Health check:**
   ```
   GET http://localhost:8080/health
   ```

2. **Gantt data:**
   ```
   GET http://localhost:8080/api/ppic/gantt-data
   ```

3. **Login:**
   ```
   POST http://localhost:8080/api/v1/auth/login
   {
     "username": "admin",
     "password": "admin123"
   }
   ```

## Migration Strategy

Script `00_complete_setup.sql` menggunakan `IF NOT EXISTS`, jadi:
- ✅ Aman dijalankan berkali-kali
- ✅ Tidak akan error jika table sudah ada
- ✅ Bisa dipakai untuk fresh install atau repair

## Next Steps

1. **Setup selesai** → Run backend → Test API
2. **Ganti password** admin default
3. **Backup database** secara berkala:
   ```bash
   pg_dump -U postgres ganttpro_db > backup_$(date +%Y%m%d).sql
   ```
4. **Deploy** ke production (update .env untuk prod)

## Notes

- Semua file dokumentasi dalam Bahasa Indonesia untuk kemudahan
- Script setup dibuat user-friendly dengan progress indicator
- Verification script memberikan output yang jelas (✓ atau ✗)
- Default data (admin & machines) otomatis di-insert

## Support

Jika masih ada masalah:
1. Check QUICK_START.md untuk troubleshooting
2. Run verify_database.sql untuk diagnostic
3. Screenshot error message
4. Hubungi tim development

---

**Created:** January 5, 2026
**Purpose:** Fix database migration issue saat pindah dari PC ke laptop
**Status:** ✅ Ready to use
