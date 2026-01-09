# Quick Fix: Network Access Error

## Error yang Terjadi
❌ **"Failed to fetch"** saat diakses dari PC lain

## Penyebab
1. Frontend masih menggunakan `localhost` di konfigurasi
2. CORS backend belum dikonfigurasi untuk menerima request dari IP network
3. Windows Firewall mungkin memblokir port 8080 dan 5173

## Solusi yang Sudah Diterapkan

### ✅ 1. Backend CORS Configuration
File: `backend/.env`
```env
ALLOWED_ORIGINS=http://localhost:5173,http://10.250.23.241:5173,http://localhost:3000,http://10.250.23.241:3000
FRONTEND_URL=http://10.250.23.241:5173
```

### ✅ 2. Frontend Configuration
Dibuat 2 file konfigurasi:

**File: `frontend/.env`** (untuk akses localhost)
```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

**File: `frontend/.env.network`** (untuk akses dari PC lain)
```env
VITE_API_BASE_URL=http://10.250.23.241:8080/api/v1
```

### ✅ 3. Helper Scripts

Dibuat 2 script untuk memudahkan:
- `start_network.bat` - Menjalankan aplikasi dengan konfigurasi network
- `setup_firewall.bat` - Membuka Windows Firewall (run as Administrator)

## Cara Menggunakan

### Opsi 1: Gunakan Helper Script (RECOMMENDED)

#### Step 1: Buka Firewall (Sekali saja)
```cmd
# Klik kanan -> Run as Administrator
setup_firewall.bat
```

#### Step 2: Jalankan Aplikasi
```cmd
# Double click atau run dari terminal
start_network.bat
```

### Opsi 2: Manual

#### Step 1: Konfigurasi Frontend
Pilih salah satu:

**A. Untuk akses dari PC lain:**
```cmd
cd frontend
copy .env.network .env
npm run dev
```

**B. Untuk akses localhost saja:**
```cmd
cd frontend
# Pastikan .env menggunakan localhost
npm run dev
```

#### Step 2: Jalankan Backend
```cmd
cd backend
go run main.go
```

#### Step 3: Buka Firewall (jika belum)
Run sebagai Administrator:
```cmd
netsh advfirewall firewall add rule name="GanttPro Backend" dir=in action=allow protocol=TCP localport=8080
netsh advfirewall firewall add rule name="GanttPro Frontend" dir=in action=allow protocol=TCP localport=5173
```

## URL untuk Akses

### Dari PC Server (localhost):
- Frontend: http://localhost:5173
- Backend: http://localhost:8080

### Dari PC Lain di Network:
- Frontend: http://10.250.23.241:5173
- Backend: http://10.250.23.241:8080

## Testing

### 1. Test Backend dari PC lain
```bash
curl http://10.250.23.241:8080/health
```
Hasil yang diharapkan:
```json
{"status":"ok"}
```

### 2. Test Frontend dari PC lain
Buka browser dan akses:
```
http://10.250.23.241:5173
```

## Troubleshooting

### ❌ Masih "Failed to fetch"?

**Checklist:**

1. ✅ **Backend sudah jalan?**
   ```cmd
   # Di PC server, cek apakah ada output log backend
   ```

2. ✅ **Frontend sudah pakai konfigurasi network?**
   ```cmd
   cd frontend
   type .env
   # Harus menunjukkan: VITE_API_BASE_URL=http://10.250.23.241:8080/api/v1
   ```

3. ✅ **Firewall sudah dibuka?**
   ```cmd
   # Run as Administrator
   netsh advfirewall firewall show rule name="GanttPro Backend"
   netsh advfirewall firewall show rule name="GanttPro Frontend"
   ```

4. ✅ **Kedua PC di network yang sama?**
   ```cmd
   # Dari PC lain, test ping
   ping 10.250.23.241
   ```

5. ✅ **Restart browser di PC client**
   - Clear cache browser (Ctrl + Shift + Delete)
   - Buka incognito/private window

### ❌ IP Address Berubah?

Jika IP address PC server berubah (misalnya setelah restart):

1. Cek IP baru:
   ```cmd
   ipconfig
   ```

2. Update file-file berikut:
   - `backend/.env` → `ALLOWED_ORIGINS` dan `FRONTEND_URL`
   - `frontend/.env.network` → `VITE_API_BASE_URL`
   - `start_network.bat` → IP di echo messages

3. Restart backend dan frontend

## Network Details

**IP Address Server:** 10.250.23.241  
**Subnet:** 10.250.0.0/16  
**Network Type:** Corporate Network (formulatrix.internal)

**Catatan:** Karena ini jaringan kantor, pastikan tidak ada proxy atau VPN yang memblokir koneksi.

---

**Last Updated:** January 6, 2026
