# JWT Token Configuration - GanttPro Backend

## 📝 Overview

Backend GanttPro menggunakan **JWT (JSON Web Token)** untuk autentikasi dan otorisasi user. Token ini digunakan untuk mengamankan endpoint API yang memerlukan autentikasi.

## ⏰ Token Expiry (Kadaluwarsa)

### Default Configuration
- **Kadaluwarsa Token**: **168 jam** (7 hari / 7 x 24 jam)
- **Algoritma**: HS256 (HMAC with SHA-256)
- **Secret Key**: Dapat dikonfigurasi via environment variable

### Konfigurasi di Code

#### File: `backend/config/config.go`
```go
// Default JWT expiry: 7 days (7 * 24 hours = 168 hours)
jwtExpiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "168"))
```

#### File: `backend/services/auth_service.go`
```go
claims := jwt.MapClaims{
    "user_id":        user.ID,
    "username":       user.Username,
    "user_id_string": user.UserID,
    "role":           user.Role,
    "operator":       user.Operator,
    "exp":            time.Now().Add(time.Hour * time.Duration(s.config.JWTExpiryHours)).Unix(),
    "iat":            time.Now().Unix(),
}
```

## 🔧 Cara Mengubah Kadaluwarsa Token

### Opsi 1: Environment Variable (Recommended)
Buat file `.env` di folder `backend/`:

```bash
# Set token expiry ke 7 hari (default)
JWT_EXPIRY_HOURS=168

# Atau pilih nilai lain:
# 24 = 1 hari
# 72 = 3 hari
# 168 = 7 hari (recommended)
# 720 = 30 hari
```

### Opsi 2: Langsung di Code
Edit file `backend/config/config.go` baris 32:

```go
jwtExpiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "168")) // Ubah "168" ke nilai yang diinginkan
```

## 🔑 JWT Token Structure

### Token Claims (Payload)
```json
{
  "user_id": 1,                    // ID user (integer)
  "username": "BAYU",              // Username
  "user_id_string": "PI0824.2374", // User ID format PI
  "role": "Admin",                 // Role user
  "operator": null,                // Operator (jika ada)
  "iat": 1729782000,              // Issued At (waktu token dibuat)
  "exp": 1730386800               // Expiry (waktu token kadaluwarsa)
}
```

### Cara Decode Token
Token JWT terdiri dari 3 bagian yang dipisahkan tanda titik (.):
```
[Header].[Payload].[Signature]
eyJhbGc...  .  eyJ1c2Vy...  .  SflKxwRJ...
```

Anda dapat decode payload di website: https://jwt.io

## 📋 Token Lifecycle

### 1. Login → Generate Token
```
User login dengan username & password
   ↓
Backend verify credentials
   ↓
Generate JWT token (valid 7 hari)
   ↓
Return token ke frontend
```

### 2. Frontend Simpan Token
```javascript
// Simpan di localStorage
localStorage.setItem('token', token);
```

### 3. Request dengan Token
```javascript
// Setiap request, kirim token di header
headers: {
  'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
}
```

### 4. Backend Validate Token
```
Request masuk
   ↓
Ekstrak token dari header Authorization
   ↓
Verify signature & expiry
   ↓
Allow/Deny request
```

## 🔒 Security Best Practices

### 1. Secret Key
**SANGAT PENTING**: Ganti JWT secret key di production!

```bash
# .env file
JWT_SECRET=gunakan-secret-key-yang-kuat-dan-random-minimal-32-karakter
```

Generate secret key yang kuat:
```bash
# Menggunakan openssl
openssl rand -base64 32

# Output contoh:
# 5K8x9P2mN7qR4tW6vY8zA1bC3dE5fG7h
```

### 2. HTTPS Only
Selalu gunakan HTTPS di production untuk mencegah token dicuri (man-in-the-middle attack).

### 3. Token Storage
- ✅ **Recommended**: Simpan di `localStorage` atau `sessionStorage`
- ❌ **Not Recommended**: Simpan di cookies (jika tidak properly configured)

### 4. Token Refresh
Jika token expired, user harus login ulang. Implementasi refresh token bisa ditambahkan nanti jika diperlukan.

## 🧪 Testing JWT Token

### 1. Login dan Dapatkan Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "BAYU",
    "password": "13032001"
  }'
```

Response:
```json
{
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzAzODY4MDAsImlhdCI6MTcyOTc4MjAwMCwib3BlcmF0b3IiOm51bGwsInJvbGUiOiJBZG1pbiIsInVzZXJfaWQiOjEsInVzZXJfaWRfc3RyaW5nIjoiUEkwODI0LjIzNzQiLCJ1c2VybmFtZSI6IkJBWVUifQ.abcdef123456...",
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

### 2. Gunakan Token untuk Request
```bash
# Simpan token ke variable
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Request dengan token
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer $TOKEN"
```

### 3. Cek Kapan Token Expire
Decode token di https://jwt.io atau:

```javascript
// Di browser console
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...";
const payload = JSON.parse(atob(token.split('.')[1]));
console.log('Issued at:', new Date(payload.iat * 1000));
console.log('Expires at:', new Date(payload.exp * 1000));
console.log('Valid for:', (payload.exp - payload.iat) / 3600, 'hours');
```

## 📊 Token Expiry Examples

| Setting | Hours | Days | Use Case |
|---------|-------|------|----------|
| 1 | 1 | 0.04 | Testing only |
| 24 | 24 | 1 | High security apps |
| 72 | 72 | 3 | Medium security |
| **168** | **168** | **7** | **Default (Recommended)** |
| 720 | 720 | 30 | Low security / convenience |

## ⚠️ Troubleshooting

### Token Expired Error
```json
{
  "error": "Token is expired"
}
```
**Solusi**: User harus login ulang untuk mendapatkan token baru.

### Invalid Token Error
```json
{
  "error": "invalid token"
}
```
**Solusi**: 
1. Pastikan token dikirim dengan format: `Bearer <token>`
2. Pastikan token tidak corrupt/terpotong
3. Pastikan JWT secret key sama dengan saat generate token

### Token Not Found
```json
{
  "error": "Authorization header required"
}
```
**Solusi**: Pastikan header `Authorization` dikirim dalam request.

## 📚 Referensi

- JWT Official: https://jwt.io/
- RFC 7519: https://tools.ietf.org/html/rfc7519
- golang-jwt library: https://github.com/golang-jwt/jwt

---

**Update Terakhir**: 24 Oktober 2024  
**Version**: 1.0  
**Author**: GanttPro Development Team
