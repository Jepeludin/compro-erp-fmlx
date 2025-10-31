# PostgreSQL Setup Guide - GanttPro

## 📋 Konfigurasi PostgreSQL

**Database Configuration:**
- Host: `localhost`
- Port: `8181`
- User: `postgres`
- Password: `jemmy1303`
- Database: `ganttpro_db`

## 🚀 Setup Database (Otomatis)

### Cara Tercepat - Gunakan PowerShell Script

```powershell
# 1. Masuk ke folder scripts
cd "d:\Jeremy Agape\Ganttpro\Project\backend\scripts"

# 2. Jalankan setup script
.\setup_database.ps1
```

Script ini akan otomatis:
- ✅ Create database `ganttpro_db`
- ✅ Create tabel `users` dengan struktur lengkap
- ✅ Create indexes untuk performa
- ✅ Create trigger auto-update `updated_at`
- ✅ Insert sample users (BAYU & Amelia)

## 🔧 Setup Database (Manual)

Jika ingin setup manual, ikuti step berikut:

### Step 1: Set Password Environment

```powershell
$env:PGPASSWORD = "jemmy1303"
```

### Step 2: Create Database

```powershell
psql -h localhost -p 8181 -U postgres -c "CREATE DATABASE ganttpro_db;"
```

### Step 3: Run SQL Script

```powershell
cd "d:\Jeremy Agape\Ganttpro\Project\backend\scripts"
psql -h localhost -p 8181 -U postgres -d ganttpro_db -f setup_postgres.sql
```

### Step 4: Verify Setup

```powershell
psql -h localhost -p 8181 -U postgres -d ganttpro_db -c "SELECT id, username, user_id, role FROM users;"
```

Expected output:
```
 id | username | user_id     | role  
----+----------+-------------+-------
  1 | BAYU     | PI0824.2374 | Admin
  2 | Amelia   | PI0824.2375 | PPIC
```

## 📊 Database Schema

### Tabel: users

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    user_id VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'PPIC',
    operator VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

**Indexes:**
- `idx_users_username` - untuk query by username
- `idx_users_user_id` - untuk query by user_id
- `idx_users_role` - untuk query by role
- `idx_users_is_active` - untuk filter active users
- `idx_users_deleted_at` - untuk soft delete

**Trigger:**
- `update_users_updated_at` - auto-update `updated_at` pada setiap UPDATE

## 👥 Sample Data

### User 1: BAYU (Admin)
```json
{
  "username": "BAYU",
  "user_id": "PI0824.2374",
  "password": "13032001",  // Plain (tersimpan ter-hash)
  "role": "Admin"
}
```

### User 2: Amelia (PPIC)
```json
{
  "username": "Amelia",
  "user_id": "PI0824.2375",
  "password": "20010313",  // Plain (tersimpan ter-hash)
  "role": "PPIC"
}
```

## 🏃‍♂️ Menjalankan Backend

### 1. Check File .env

File `backend/.env` sudah dikonfigurasi:
```bash
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=8181
DB_USER=postgres
DB_PASSWORD=jemmy1303
DB_NAME=ganttpro_db
DB_SSLMODE=disable
```

### 2. Install Dependencies

```powershell
cd "d:\Jeremy Agape\Ganttpro\Project\backend"
go mod tidy
```

### 3. Run Backend Server

```powershell
go run main.go
```

Output yang diharapkan:
```
2024/10/25 10:00:00 Connecting to database...
2024/10/25 10:00:00 Database connected successfully
2024/10/25 10:00:00 Starting server on :8080
```

## 🧪 Testing

### Test Connection ke Database

```powershell
psql -h localhost -p 8181 -U postgres -d ganttpro_db
```

Di dalam psql:
```sql
-- List all tables
\dt

-- Check users table
SELECT * FROM users;

-- Check table structure
\d users

-- Exit
\q
```

### Test Backend API - Login

```powershell
# Test login BAYU
curl -X POST http://localhost:8080/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{\"username\":\"BAYU\",\"password\":\"13032001\"}'
```

Success response:
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "BAYU",
      "user_id": "PI0824.2374",
      "role": "Admin",
      "is_active": true
    }
  }
}
```

## ⚠️ Troubleshooting

### Problem: `psql: command not found`

**Solusi:**
1. Cari PostgreSQL bin folder (contoh: `C:\Program Files\PostgreSQL\16\bin`)
2. Tambahkan ke System Environment Variable PATH
3. Restart PowerShell dan coba lagi

### Problem: `Connection refused` atau `could not connect to server`

**Solusi:**
1. Pastikan PostgreSQL service berjalan:
   ```powershell
   # Check service
   Get-Service -Name postgresql*
   
   # Start service jika stopped
   Start-Service postgresql-x64-16  # sesuaikan nama service
   ```

2. Check port 8181:
   ```powershell
   netstat -an | findstr 8181
   ```

3. Check `postgresql.conf`:
   - Pastikan `port = 8181`
   - Restart service setelah perubahan

### Problem: `password authentication failed`

**Solusi:**
1. Pastikan password benar: `jemmy1303`
2. Check `pg_hba.conf`:
   ```
   # Tambahkan/ubah line berikut:
   host    all    all    127.0.0.1/32    md5
   ```
3. Restart PostgreSQL service

### Problem: `database "ganttpro_db" already exists`

**Solusi:**
Ini normal. Script akan melanjutkan ke create table.

Jika ingin reset database:
```powershell
# Drop database
psql -h localhost -p 8181 -U postgres -c "DROP DATABASE IF EXISTS ganttpro_db;"

# Jalankan ulang setup script
.\setup_database.ps1
```

### Problem: Backend error `failed to connect to database`

**Solusi:**
1. Check `.env` file ada dan isinya benar
2. Check PostgreSQL berjalan di port 8181
3. Check kredensial: user `postgres`, password `jemmy1303`
4. Check database `ganttpro_db` sudah dibuat

## 📁 File Penting

```
backend/
├── .env                          # Environment config (JANGAN commit!)
├── .env.example                  # Template environment
├── config/
│   └── config.go                # Config loader (sudah diupdate untuk postgres)
├── scripts/
│   ├── setup_database.ps1       # PowerShell setup script (RUN THIS!)
│   ├── setup_postgres.sql       # SQL setup script
│   └── create_sample_users.sql  # Sample users insert
└── database/
    └── db.go                    # Database connection
```

## ✅ Checklist Setup

- [ ] PostgreSQL installed dan berjalan di port 8181
- [ ] Password postgres diset ke `jemmy1303`
- [ ] psql command tersedia di PATH
- [ ] Database `ganttpro_db` dibuat
- [ ] Tabel `users` dibuat dengan sample data
- [ ] File `backend/.env` dikonfigurasi
- [ ] Dependencies Go installed (`go mod tidy`)
- [ ] Backend berjalan (`go run main.go`)
- [ ] Login API tested dan berhasil

## 🎯 Next Steps

Setelah database setup berhasil:

1. **Jalankan Backend:**
   ```powershell
   cd "d:\Jeremy Agape\Ganttpro\Project\backend"
   go run main.go
   ```

2. **Jalankan Frontend:**
   ```powershell
   cd "d:\Jeremy Agape\Ganttpro\Project\frontend"
   npm install
   npm run dev
   ```

3. **Test Login:**
   - Buka http://localhost:5173
   - Login dengan: Username `BAYU`, Password `13032001`

## 📚 Dokumentasi Lainnya

- `README.md` - Overview project
- `API_TESTING.md` - API testing guide
- `JWT_CONFIGURATION.md` - JWT configuration
- `STRUKTUR_USER.md` - User structure details

---

**Database**: PostgreSQL  
**Port**: 8181  
**Setup Date**: 25 Oktober 2024
