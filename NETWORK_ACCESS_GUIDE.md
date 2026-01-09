# Panduan Akses Aplikasi dari PC Lain di Jaringan Lokal

## Konfigurasi yang Telah Diterapkan

### 1. Backend (Go)
- Server sekarang bind ke `0.0.0.0:8080` (bukan hanya `localhost`)
- Ini memungkinkan backend menerima koneksi dari semua network interfaces

### 2. Frontend (Vite)
- Dev server dikonfigurasi dengan `host: '0.0.0.0'`
- Port: 5173
- Dapat diakses dari komputer lain di jaringan yang sama

## Cara Mengakses dari PC Lain

### Langkah 1: Cek IP Address PC Server
Di PC yang menjalankan aplikasi, jalankan perintah:

**Windows (Command Prompt/PowerShell):**
```cmd
ipconfig
```

**Windows (Git Bash/Terminal):**
```bash
hostname -I
```

Cari IP address yang diawali dengan:
- `192.168.x.x` (jaringan rumah/kantor umum)
- `10.x.x.x` (jaringan kantor)
- `172.16.x.x` - `172.31.x.x` (jaringan kantor)

Contoh: `192.168.1.100`

### Langkah 2: Jalankan Aplikasi

**Backend:**
```bash
cd backend
go run main.go
```
Server akan berjalan di: `http://0.0.0.0:8080` (dapat diakses dari semua network interfaces)

**Frontend:**
```bash
cd frontend
npm run dev
```
Server akan menampilkan:
```
VITE v5.x.x  ready in xxx ms

➜  Local:   http://localhost:5173/
➜  Network: http://192.168.1.100:5173/
```

### Langkah 3: Akses dari PC Lain

Dari PC lain di jaringan yang sama:

**Frontend:**
```
http://192.168.1.100:5173
```

**Backend API:**
```
http://192.168.1.100:8080/api
```

**Contoh request API:**
```
http://192.168.1.100:8080/api/auth/login
http://192.168.1.100:8080/api/machines
```

## Konfigurasi CORS

Jika mengakses dari IP berbeda, Anda mungkin perlu menambahkan allowed origin di environment variables:

### Cara 1: Update Environment Variable (Recommended)

Tambahkan IP address PC Anda ke `ALLOWED_ORIGINS` di file `.env` (buat jika belum ada):

**File: `backend/.env`**
```env
# Multiple origins dipisah dengan koma
ALLOWED_ORIGINS=http://localhost:5173,http://192.168.1.100:5173,http://192.168.1.101:5173

# Frontend URL untuk email links
FRONTEND_URL=http://192.168.1.100:5173
```

### Cara 2: Allow All Origins (Development Only - TIDAK AMAN untuk Production)

Di [backend/config/config.go](backend/config/config.go), origin default sudah di-set. Anda bisa menambahkan wildcard untuk development:

```go
// Untuk development, bisa tambahkan "*" ke allowed origins
AllowedOrigins: []string{"*"}  // WARNING: Jangan gunakan di production!
```

## Troubleshooting

### 1. Tidak Bisa Akses dari PC Lain

**Penyebab:** Windows Firewall memblokir koneksi

**Solusi:**
```powershell
# Jalankan sebagai Administrator di PowerShell

# Allow port 8080 (Backend)
New-NetFirewallRule -DisplayName "Gantt Pro Backend" -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow

# Allow port 5173 (Frontend)
New-NetFirewallRule -DisplayName "Gantt Pro Frontend" -Direction Inbound -LocalPort 5173 -Protocol TCP -Action Allow
```

Atau secara manual:
1. Buka **Windows Defender Firewall**
2. Pilih **Advanced Settings**
3. Pilih **Inbound Rules** → **New Rule**
4. Pilih **Port** → Next
5. Pilih **TCP** → Specific local ports: `8080,5173`
6. Pilih **Allow the connection**
7. Apply untuk semua profiles (Domain, Private, Public)
8. Beri nama: "Gantt Pro Application"

### 2. CORS Error

**Error:** `Access to XMLHttpRequest has been blocked by CORS policy`

**Solusi:**
- Tambahkan IP address frontend ke `ALLOWED_ORIGINS` di `.env`
- Restart backend setelah update `.env`

### 3. Frontend Tidak Load Asset/API

**Penyebab:** Frontend masih menggunakan hardcoded `localhost` untuk API URL

**Solusi:** 
- Update file konfigurasi frontend untuk menggunakan IP dinamis atau environment variable
- Contoh di `src/config.js` atau `.env` di frontend:
  ```javascript
  // Auto-detect atau gunakan environment variable
  const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://192.168.1.100:8080'
  ```

### 4. PC Tidak Saling Terhubung

**Checklist:**
- ✅ Kedua PC terhubung ke WiFi/network yang sama
- ✅ PC tidak dalam mode "Public Network" (ubah ke "Private Network")
- ✅ Network discovery enabled di Windows Settings
- ✅ Coba ping dari PC lain: `ping 192.168.1.100`

## Catatan Keamanan

⚠️ **PENTING untuk Production:**

1. **JANGAN** bind ke `0.0.0.0` di production public server
2. **JANGAN** gunakan `ALLOWED_ORIGINS=*` di production
3. **SELALU** gunakan HTTPS di production
4. **BATASI** allowed origins ke domain spesifik
5. **GUNAKAN** reverse proxy (nginx/caddy) dengan proper security headers

## Testing Koneksi

### Test Backend
```bash
# Dari PC lain, test health check
curl http://192.168.1.100:8080/health
```

Harusnya return:
```json
{"status":"ok"}
```

### Test Frontend
Buka browser di PC lain dan akses:
```
http://192.168.1.100:5173
```

## Konfigurasi untuk Development Team

Jika bekerja dalam tim, setiap developer bisa:

1. Share IP address mereka
2. Tambahkan semua IP ke `ALLOWED_ORIGINS`
3. Test API menggunakan IP rekan

**Contoh `.env` untuk tim:**
```env
ALLOWED_ORIGINS=http://localhost:5173,http://192.168.1.100:5173,http://192.168.1.101:5173,http://192.168.1.102:5173
```

---

**Last Updated:** January 6, 2026
