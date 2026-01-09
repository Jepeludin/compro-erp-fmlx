# GanttPro Backend - Golang API

Backend API untuk aplikasi GanttPro ERP yang dibangun dengan Golang menggunakan Gin framework.

## 🚀 Fitur

- ✅ Authentication dengan JWT
- ✅ Password hashing dengan bcrypt
- ✅ Database PostgreSQL dengan migration support
- ✅ CORS middleware
- ✅ Clean architecture dengan separation of concerns
- ✅ Environment-based configuration
- ✅ PPIC Gantt Chart scheduling
- ✅ Machine assignments & dependencies
- ✅ Role-based access control

## 📁 Struktur Project

```
backend/
├── config/              # Konfigurasi aplikasi
├── database/            # Database connection & migration
│   └── migrations/      # SQL migration files
├── handlers/            # HTTP request handlers
├── middleware/          # Middleware (CORS, Auth, Role, dll)
├── models/              # Database models & structs
├── repository/          # Database operations layer
├── routes/              # Route definitions
├── services/            # Business logic layer
├── utils/               # Utility functions
├── testing/             # Unit tests
├── uploads/             # Upload directory
├── .env                 # Environment variables (jangan commit!)
├── setup_database.ps1   # Auto setup script (Windows)
├── QUICK_START.md       # Panduan cepat setup database
├── SETUP_DATABASE.md    # Dokumentasi lengkap database setup
├── go.mod               # Go module dependencies
└── main.go              # Entry point aplikasi
```

## ⚡ Quick Start (Setup di Laptop Baru)

**Jika Anda baru pindah ke laptop baru dan database belum di-setup:**

### Cara Tercepat (Windows):

```powershell
# 1. Masuk ke folder backend
cd C:\Jemmy\compro\compro-erp-fmlx\backend

# 2. Jalankan script setup otomatis
.\setup_database.ps1

# 3. Jalankan backend
go run main.go
```

Script akan otomatis:
- ✓ Test koneksi PostgreSQL
- ✓ Buat database baru
- ✓ Jalankan semua migration
- ✓ Insert data default (admin user & machines)
- ✓ Update file .env

**Lihat [QUICK_START.md](QUICK_START.md) untuk detail lengkap!**

## 🛠️ Setup Manual

### Prerequisites

- Go 1.21 atau lebih tinggi
- PostgreSQL 12 atau lebih tinggi
- Git

### Langkah Instalasi

1. **Masuk ke direktori backend:**
   ```bash
   cd backend
   ```

2. **Setup Database:**
   
   Lihat dokumentasi lengkap di [SETUP_DATABASE.md](SETUP_DATABASE.md)
   
   Atau jalankan quick setup:
   ```bash
   # Windows PowerShell
   Copy-Item .env.example .env
   
   # Atau manual copy paste
   ```

3. **Edit file `.env` sesuai konfigurasi Anda:**
   ```env
   PORT=8080
   DB_DRIVER=mysql
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=your_password
   DB_NAME=ganttpro_db
   JWT_SECRET=your-secret-key-change-this
   ```

4. **Buat database:**
   ```sql
   -- MySQL
   CREATE DATABASE ganttpro_db;
   
   -- PostgreSQL
   CREATE DATABASE ganttpro_db;
   ```

5. **Install dependencies:**
   ```bash
   go mod download
   ```

6. **Jalankan aplikasi:**
   ```bash
   go run main.go
   ```

Server akan berjalan di `http://localhost:8080`

## 📡 API Endpoints

### Authentication

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "userid": "user123",
  "password": "password123"
}
```

Response:
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "userid": "user123",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "user",
      "is_active": true
    }
  }
}
```

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "userid": "newuser",
  "password": "password123",
  "name": "New User",
  "email": "newuser@example.com"
}
```

#### Get Profile (Protected)
```http
GET /api/v1/auth/profile
Authorization: Bearer <your-jwt-token>
```

### Health Check
```http
GET /health
```

## 🔐 Membuat User Pertama

Untuk membuat user pertama, gunakan endpoint register atau jalankan query SQL:

```sql
-- MySQL/PostgreSQL
INSERT INTO users (user_id, password, name, email, role, is_active, created_at, updated_at)
VALUES (
  'admin',
  '$2a$10$YourBcryptHashedPasswordHere',
  'Administrator',
  'admin@example.com',
  'admin',
  true,
  NOW(),
  NOW()
);
```

Atau gunakan Go untuk generate password hash:
```bash
go run -c 'package main; import ("fmt"; "golang.org/x/crypto/bcrypt"); func main() { hash, _ := bcrypt.GenerateFromPassword([]byte("yourpassword"), bcrypt.DefaultCost); fmt.Println(string(hash)) }'
```

## 🔧 Development

### Build aplikasi:
```bash
go build -o ganttpro-backend.exe
```

### Run aplikasi yang sudah di-build:
```bash
./ganttpro-backend.exe
```

### Run dengan hot reload (gunakan air):
```bash
go install github.com/cosmtrek/air@latest
air
```

## 🌐 Integrasi dengan Frontend

Update file `.env` frontend Anda untuk mengarah ke backend:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

Contoh penggunaan di frontend (Vue.js):

```javascript
// Login request
async function handleLogin() {
  try {
    const response = await fetch('http://localhost:8080/api/v1/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        userid: userid.value,
        password: password.value
      })
    });

    const data = await response.json();
    
    if (response.ok) {
      // Simpan token
      localStorage.setItem('token', data.data.token);
      localStorage.setItem('user', JSON.stringify(data.data.user));
      
      // Redirect ke dashboard
      router.push('/dashboard');
    } else {
      alert(data.error);
    }
  } catch (error) {
    console.error('Login error:', error);
  }
}
```

## 📝 Database Schema

### Users Table
```sql
CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    role VARCHAR(20) DEFAULT 'user',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

## 🤝 Contributing

Silakan buat branch baru untuk setiap fitur atau perbaikan:

```bash
git checkout -b feature/nama-fitur
```

## 📄 License

MIT License
