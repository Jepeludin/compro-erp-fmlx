# Quick Start Guide

## Backend Setup

1. **Masuk ke folder backend:**
   ```powershell
   cd backend
   ```

2. **Copy dan edit file .env:**
   ```powershell
   Copy-Item .env.example .env
   # Edit .env dengan text editor
   ```

3. **Buat database (di MySQL/PostgreSQL):**
   ```sql
   CREATE DATABASE ganttpro_db;
   ```

4. **Install dependencies dan jalankan:**
   ```powershell
   go mod download
   go run main.go
   ```

5. **Buat user admin (opsional):**
   ```powershell
   # Jalankan script SQL di scripts/create_admin.sql
   # Atau gunakan endpoint /api/v1/auth/register
   ```

## Frontend Setup

1. **Masuk ke folder frontend:**
   ```powershell
   cd frontend
   ```

2. **Install dependencies:**
   ```powershell
   npm install
   ```

3. **Jalankan dev server:**
   ```powershell
   npm run dev
   ```

4. **Buka browser:**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080

## Test Login

Gunakan credentials default (setelah run script SQL):
- **Username:** admin
- **Password:** admin123

Atau register user baru melalui API:
```powershell
# PowerShell
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"userid":"testuser","password":"password123","name":"Test User","email":"test@example.com"}'
```

## API Endpoints

- **POST** `/api/v1/auth/login` - Login
- **POST** `/api/v1/auth/register` - Register
- **GET** `/api/v1/auth/profile` - Get profile (requires token)
- **GET** `/health` - Health check

## Struktur Database

```
users
- id (PK)
- user_id (unique)
- password (hashed)
- name
- email (unique)
- role (admin/user)
- is_active
- created_at
- updated_at
- deleted_at
```

## File Structure

```
backend/
├── config/          # Configuration
├── database/        # DB connection & migration
├── handlers/        # HTTP handlers
├── middleware/      # Middleware (CORS, Auth)
├── models/          # Database models
├── repository/      # Data access layer
├── routes/          # Route definitions
├── services/        # Business logic
├── utils/           # Utilities
└── main.go          # Entry point

frontend/
├── src/
│   ├── components/  # Vue components
│   ├── router/      # Vue Router
│   ├── services/    # API services
│   └── main.js
└── package.json
```
