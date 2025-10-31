# Struktur User Database - GanttPro

## 📋 Struktur Tabel Users

Berdasarkan contoh data yang diberikan, berikut adalah struktur lengkap tabel `users`:

| Field      | Type          | Description                                    | Example           |
|------------|---------------|------------------------------------------------|-------------------|
| id         | BIGINT        | Primary key (auto increment)                   | 1, 2, 3           |
| username   | VARCHAR(50)   | Nama user (unique)                             | BAYU, Amelia      |
| user_id    | VARCHAR(50)   | ID format PI (unique)                          | PI0824.2374       |
| password   | VARCHAR(255)  | Password ter-hash (bcrypt)                     | $2a$10$...       |
| role       | VARCHAR(20)   | Role/jabatan user                              | Admin, PPIC       |
| operator   | VARCHAR(50)   | Field operator (nullable)                      | NULL              |
| is_active  | BOOLEAN       | Status aktif user                              | true, false       |
| created_at | TIMESTAMP     | Waktu pembuatan record                         | 2024-10-24 10:00  |
| updated_at | TIMESTAMP     | Waktu update terakhir                          | 2024-10-24 10:00  |
| deleted_at | TIMESTAMP     | Soft delete (nullable)                         | NULL              |

## 👥 Contoh Data User

### User 1: BAYU (Admin)
```json
{
  "username": "BAYU",
  "user_id": "PI0824.2374",
  "password": "13032001",  // Plain (akan di-hash)
  "role": "Admin",
  "operator": null
}
```

### User 2: Amelia (PPIC)
```json
{
  "username": "Amelia",
  "user_id": "PI0824.2375",
  "password": "20010313",  // Plain (akan di-hash)
  "role": "PPIC",
  "operator": null
}
```

## 🔐 Hashing Password

Password disimpan menggunakan **bcrypt** dengan cost 10:

```go
// Contoh hashing di Go
import "golang.org/x/crypto/bcrypt"

hashedPassword, err := bcrypt.GenerateFromPassword([]byte("13032001"), bcrypt.DefaultCost)
```

### Password Hash Examples:
- Password `13032001` → `$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdnYVO2C6`
- Password `20010313` → `$2a$10$8kL3M5XvGQqL7xZc8FqO7uyLGxJ9xH8p3F2KdLsYGxZhGJK9VH3M2`

## 🔑 API Login Request

### Request Format:
```json
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "BAYU",
  "password": "13032001"
}
```

### Response Format (Success):
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
      "operator": null,
      "is_active": true,
      "created_at": "2024-10-24T10:00:00Z"
    }
  }
}
```

### Response Format (Error):
```json
{
  "error": "invalid credentials"
}
```

## 📝 API Register Request

### Request Format:
```json
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "NewUser",
  "user_id": "PI0824.2376",
  "password": "password123",
  "role": "PPIC",
  "operator": ""
}
```

### Response Format (Success):
```json
{
  "message": "Registration successful",
  "data": {
    "id": 3,
    "username": "NewUser",
    "user_id": "PI0824.2376",
    "role": "PPIC",
    "operator": "",
    "is_active": true,
    "created_at": "2024-10-24T10:30:00Z"
  }
}
```

## 🗄️ SQL Setup Database

### Membuat Tabel Users:
```sql
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    user_id VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'PPIC',
    operator VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_users_username (username),
    INDEX idx_users_user_id (user_id),
    INDEX idx_users_role (role)
);
```

### Insert Sample Data:
```sql
-- BAYU (Admin)
INSERT INTO users (username, user_id, password, role, operator, is_active)
VALUES (
  'BAYU',
  'PI0824.2374',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdnYVO2C6',
  'Admin',
  NULL,
  TRUE
);

-- Amelia (PPIC)
INSERT INTO users (username, user_id, password, role, operator, is_active)
VALUES (
  'Amelia',
  'PI0824.2375',
  '$2a$10$8kL3M5XvGQqL7xZc8FqO7uyLGxJ9xH8p3F2KdLsYGxZhGJK9VH3M2',
  'PPIC',
  NULL,
  TRUE
);
```

## 🧪 Testing Login

### Menggunakan curl:
```bash
# Login sebagai BAYU
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"BAYU","password":"13032001"}'

# Login sebagai Amelia
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"Amelia","password":"20010313"}'
```

### Menggunakan Postman:
1. Buat request POST ke `http://localhost:8080/api/v1/auth/login`
2. Set Header: `Content-Type: application/json`
3. Set Body (raw JSON):
   ```json
   {
     "username": "BAYU",
     "password": "13032001"
   }
   ```
4. Klik Send

## 📱 Frontend Integration

### Login Component (Vue.js):
```vue
<script setup>
import { ref } from 'vue';
import api from '@/services/api';

const username = ref('');
const password = ref('');

const handleLogin = async () => {
  try {
    const response = await api.login(username.value, password.value);
    api.saveAuth(response.data.token, response.data.user);
    alert(`Welcome, ${response.data.user.username}!`);
  } catch (error) {
    alert('Login failed: ' + error.message);
  }
};
</script>

<template>
  <form @submit.prevent="handleLogin">
    <input v-model="username" placeholder="Username" required />
    <input v-model="password" type="password" placeholder="Password" required />
    <button type="submit">Login</button>
  </form>
</template>
```

## 🔒 JWT Token Claims

Token JWT berisi informasi berikut:
```json
{
  "user_id": 1,
  "username": "BAYU",
  "user_id_string": "PI0824.2374",
  "role": "Admin",
  "operator": null,
  "exp": 1729868400,
  "iat": 1729782000
}
```

## ✅ Checklist Implementasi

- [x] Model User dengan field username, user_id, role, operator
- [x] Repository dengan FindByUsername dan FindByUserIDString
- [x] Auth Service untuk login dengan username
- [x] Auth Handler untuk endpoint /auth/login
- [x] Frontend LoginPage dengan input username
- [x] API service dengan method login(username, password)
- [x] SQL scripts untuk setup database
- [x] Sample users: BAYU (Admin) dan Amelia (PPIC)
- [x] Password hashing dengan bcrypt
- [x] JWT token generation dengan claims lengkap

## 📚 Referensi File

- Backend Model: `backend/models/user.go`
- Backend Repository: `backend/repository/user_repository.go`
- Backend Service: `backend/services/auth_service.go`
- Backend Handler: `backend/handlers/auth_handler.go`
- Frontend Component: `frontend/src/components/LoginPage.vue`
- Frontend API: `frontend/src/services/api.js`
- SQL Setup: `backend/scripts/setup_database.sql`
- Sample Users: `backend/scripts/create_sample_users.sql`
