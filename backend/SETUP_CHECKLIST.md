# ✅ Database Setup Checklist

Ikuti checklist ini untuk setup database di laptop baru.

## 📋 Pre-Setup

- [ ] PostgreSQL sudah terinstall
- [ ] PostgreSQL service sedang running
- [ ] Sudah tahu username & password PostgreSQL
- [ ] Sudah masuk ke folder backend: `cd C:\Jemmy\compro\compro-erp-fmlx\backend`

## 🚀 Setup Database (Pilih salah satu)

### Opsi A: Setup Otomatis (Recommended) ⭐

- [ ] Buka PowerShell
- [ ] Run: `.\setup_database.ps1`
- [ ] Masukkan username PostgreSQL (default: postgres)
- [ ] Masukkan password PostgreSQL
- [ ] Masukkan nama database (default: ganttpro_db)
- [ ] Tunggu sampai script selesai
- [ ] Lihat output: "Setup Database BERHASIL!"

### Opsi B: Setup Manual

- [ ] Create database:
  ```bash
  psql -U postgres
  CREATE DATABASE ganttpro_db;
  \q
  ```
- [ ] Run migration:
  ```bash
  psql -U postgres -d ganttpro_db -f database/migrations/00_complete_setup.sql
  ```
- [ ] Update file .env dengan kredensial yang benar
- [ ] Verify setup:
  ```bash
  psql -U postgres -d ganttpro_db -f database/migrations/verify_database.sql
  ```

## ✓ Verifikasi

- [ ] Check file .env ada dan isinya benar:
  ```
  DB_HOST=localhost
  DB_PORT=5432
  DB_USER=postgres
  DB_PASSWORD=your_password
  DB_NAME=ganttpro_db
  ```

- [ ] Check database punya 12 tables:
  ```bash
  psql -U postgres -d ganttpro_db -c "\dt"
  ```
  Harus ada: users, machines, ppic_schedules, machine_assignments, ppic_links, dll

- [ ] Check admin user sudah ada:
  ```bash
  psql -U postgres -d ganttpro_db -c "SELECT email, username, role FROM users;"
  ```
  Harus ada: admin@example.com | admin | admin

- [ ] Check machines sudah ada:
  ```bash
  psql -U postgres -d ganttpro_db -c "SELECT machine_code, machine_name FROM machines;"
  ```
  Harus ada 5 machines: CNC-001, CNC-002, MILL-001, LATHE-001, GRIND-001

- [ ] Check kolom scheduled_start ada:
  ```bash
  psql -U postgres -d ganttpro_db -c "\d machine_assignments"
  ```
  Harus ada kolom: scheduled_start, scheduled_end, actual_start, actual_end

## 🏃 Run Backend

- [ ] Install dependencies (jika belum):
  ```bash
  go mod download
  ```

- [ ] Run backend:
  ```bash
  go run main.go
  ```

- [ ] Check output console:
  ```
  Connected to database successfully!
  Server is running on port 8080
  ```
  ✅ Jika ada pesan ini, setup BERHASIL!
  ❌ Jika ada error, lihat Troubleshooting

## 🧪 Test API

- [ ] Open browser atau Postman

- [ ] Test health endpoint:
  ```
  GET http://localhost:8080/health
  ```
  Expected: Status 200 OK

- [ ] Test gantt-data endpoint:
  ```
  GET http://localhost:8080/api/ppic/gantt-data
  ```
  Expected: 
  ```json
  {
    "success": true,
    "data": {
      "sections": [],
      "links": []
    }
  }
  ```
  ✅ Jika dapat response ini, API WORKS!
  ❌ Jika dapat error 500, ada masalah database

- [ ] Test login:
  ```
  POST http://localhost:8080/api/v1/auth/login
  Body:
  {
    "username": "admin",
    "password": "admin123"
  }
  ```
  Expected: JWT token dalam response

## 🔒 Security

- [ ] Ganti password admin default:
  ```bash
  psql -U postgres -d ganttpro_db
  UPDATE users SET password = '$2a$10$NEW_HASH_HERE' WHERE username = 'admin';
  ```
  Atau login ke aplikasi dan ganti lewat UI

- [ ] Update JWT_SECRET di .env dengan value yang strong
  ```
  JWT_SECRET=generate-random-secret-here-min-32-chars
  ```

- [ ] Pastikan .env ada di .gitignore (jangan commit password!)

## 💾 Backup

- [ ] Buat backup database pertama:
  ```bash
  pg_dump -U postgres ganttpro_db > backup_initial_setup.sql
  ```

- [ ] Simpan backup di tempat yang aman

## 📝 Documentation

- [ ] Baca QUICK_START.md untuk troubleshooting
- [ ] Baca SETUP_DATABASE.md untuk dokumentasi lengkap
- [ ] Bookmark API.md untuk referensi API endpoints

## ⚠️ Troubleshooting

Jika ada error, cek:

### Error: "psql is not recognized"
- [ ] Install PostgreSQL
- [ ] Tambahkan PostgreSQL ke PATH
- [ ] Restart terminal/PowerShell

### Error: "password authentication failed"
- [ ] Cek password di .env benar
- [ ] Test login manual: `psql -U postgres`
- [ ] Reset password PostgreSQL jika lupa

### Error: "database ganttpro_db already exists"
- [ ] Drop database: `psql -U postgres -c "DROP DATABASE ganttpro_db;"`
- [ ] Buat ulang: `psql -U postgres -c "CREATE DATABASE ganttpro_db;"`
- [ ] Run migration lagi

### Error 500 dari API
- [ ] Check backend console untuk error message
- [ ] Verify database: `psql -U postgres -d ganttpro_db -f database/migrations/verify_database.sql`
- [ ] Restart backend
- [ ] Check .env connection string

### Backend tidak bisa connect
- [ ] Cek PostgreSQL service running
- [ ] Cek port 5432 available: `netstat -ano | findstr :5432`
- [ ] Cek firewall tidak block PostgreSQL
- [ ] Test manual connection: `psql -U postgres -d ganttpro_db -h localhost -p 5432`

## 🎉 Done!

Jika semua checklist ✅, selamat! Database sudah ready dan backend berjalan dengan baik.

Next steps:
1. Mulai develop/test features
2. Buat backup database secara berkala
3. Update documentation jika ada perubahan

---

**Need help?** Check:
- QUICK_START.md - Panduan cepat
- SETUP_DATABASE.md - Dokumentasi lengkap
- DATABASE_MIGRATION_SOLUTION.md - Penjelasan solusi

**Support:** Tim Development
