# 🚀 Quick Start - Sign Up Feature

## ✨ Fitur Baru: Sign Up dengan Role "Guest"

### 📋 Cara Menggunakan

#### **1. Sign Up (Daftar Akun Baru)**

1. Buka browser: `http://localhost:5173`
2. Klik link **"Sign Up"** di bawah tombol Login
3. Isi form:
   - **Username**: Nama Anda (contoh: JohnDoe)
   - **User ID**: ID unik (contoh: PI0824.5001)
   - **Password**: Minimal 4 karakter
4. Klik tombol **"Sign Up"**
5. Akan muncul pesan sukses ✅
6. Otomatis pindah ke mode Login setelah 3.5 detik

#### **2. Login dengan Akun Baru**

Setelah signup berhasil, login dengan:
- **User ID**: `PI0824.5001` (sesuai yang Anda daftarkan)
- **Password**: password yang Anda buat

---

## 🎯 Role Otomatis

Semua akun yang dibuat melalui Sign Up akan mendapat role **"Guest"** secara otomatis.

---

## 📝 Format User ID

Gunakan format: `PI0824.XXXX`

Contoh:
- `PI0824.5001` ✅
- `PI0824.5002` ✅
- `PI0824.6001` ✅

---

## ⚡ Contoh Cepat

```
Username  : Sarah
User ID   : PI0824.5001
Password  : sarah123
Role      : Guest (otomatis)
```

Setelah signup, login dengan:
```
User ID   : PI0824.5001
Password  : sarah123
```

---

## 🔍 Check Database

Untuk melihat user baru di database:

```sql
SELECT username, user_id, role, created_at 
FROM public.users 
WHERE role = 'Guest' 
ORDER BY created_at DESC;
```

---

## 📚 Dokumentasi Lengkap

Lihat file: `SIGNUP_FEATURE.md` untuk dokumentasi detail.

---

## ✅ Checklist Testing

- [ ] Buka `http://localhost:5173`
- [ ] Klik "Sign Up"
- [ ] Isi Username, User ID, Password
- [ ] Klik "Sign Up"
- [ ] Lihat pesan sukses
- [ ] Login dengan akun baru
- [ ] Check database untuk konfirmasi

---

**Happy Coding!** 🎉
