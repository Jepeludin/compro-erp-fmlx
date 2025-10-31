# 📝 SIGNUP FEATURE - GANTTPRO

## 🎯 Overview
Fitur Sign Up memungkinkan pengguna baru untuk mendaftar akun dengan role **"Guest"** secara otomatis.

---

## ✨ Fitur Utama

### 1. **Form Sign Up**
- **Username** - Nama pengguna (e.g., John, Sarah)
- **User ID** - ID unik pengguna (e.g., PI0824.5001)
- **Password** - Password minimal 4 karakter
- **Role** - Otomatis diset sebagai **"Guest"**

### 2. **Toggle Login/Signup**
- User bisa beralih antara mode Login dan Sign Up
- Link "Sign Up" di halaman Login
- Link "Login" di halaman Sign Up

### 3. **Validasi**
- Semua field wajib diisi
- Password minimal 4 karakter
- Username harus unik
- User ID harus unik

### 4. **Auto-Switch ke Login**
- Setelah signup berhasil, akan menampilkan pesan sukses
- Otomatis beralih ke mode Login setelah 3.5 detik
- User bisa langsung login dengan akun baru

---

## 🚀 Cara Menggunakan

### **Sign Up - Mendaftar Akun Baru**

1. **Buka halaman Login**
   - URL: `http://localhost:5173` (atau port Vite Anda)

2. **Klik "Sign Up"**
   - Link ada di bawah tombol Login

3. **Isi Form Sign Up**
   ```
   Username  : JohnDoe
   User ID   : PI0824.5001
   Password  : mypassword123
   ```

4. **Klik tombol "Sign Up"**
   - Sistem akan membuat akun dengan role "Guest"
   - Menampilkan pesan sukses
   - Otomatis pindah ke mode Login

5. **Login dengan Akun Baru**
   ```
   User ID   : PI0824.5001
   Password  : mypassword123
   ```

---

## 🔧 Implementasi Teknis

### **Frontend - LoginPage.vue**

#### **Data Reactive**
```javascript
const isSignupMode = ref(false);  // Toggle mode
const username = ref('');         // Field username (signup only)
const userId = ref('');           // User ID
const password = ref('');         // Password
const successMessage = ref('');   // Success message
```

#### **Functions**
```javascript
toggleMode()     // Toggle antara Login dan Signup
handleLogin()    // Handler untuk login
handleSignup()   // Handler untuk signup (role: Guest)
```

#### **API Call - Signup**
```javascript
const response = await api.register({
  username: username.value,
  user_id: userId.value,
  password: password.value,
  role: 'Guest',        // Otomatis Guest
  operator: null
});
```

### **Backend - API Endpoint**

#### **POST /api/v1/auth/register**
```json
Request Body:
{
  "username": "JohnDoe",
  "user_id": "PI0824.5001",
  "password": "mypassword123",
  "role": "Guest",
  "operator": null
}

Response (201 Created):
{
  "message": "Registration successful",
  "data": {
    "id": 19,
    "username": "JohnDoe",
    "user_id": "PI0824.5001",
    "role": "Guest",
    "operator": null,
    "is_active": true,
    "created_at": "2025-10-26T00:30:00Z"
  }
}
```

### **Database Schema**

Role "Guest" akan disimpan di tabel `users`:

```sql
INSERT INTO public.users (username, user_id, password, role, operator, is_active)
VALUES ('JohnDoe', 'PI0824.5001', '$2a$10$...', 'Guest', NULL, TRUE);
```

---

## 🎨 UI/UX Features

### **Success Message**
```css
Green background (#d4edda)
Smooth slide-down animation
Auto-dismiss after switching to login mode
```

### **Toggle Link**
```css
Blue link with underline animation on hover
Smooth transition between modes
Form fields cleared on mode switch
```

### **Form Animations**
- Slide up animation on load
- Shake animation for errors
- Smooth transitions for all interactions

---

## 🔐 Security

### **Password Handling**
- Password di-hash dengan bcrypt (cost 10)
- Never stored in plain text
- Minimal 4 karakter

### **Validation**
- Backend validates all inputs
- Checks for duplicate username/user_id
- Returns appropriate error messages

---

## 📊 User Roles

| Role | Description | Created By |
|------|-------------|------------|
| Admin | Full access | Manual (SQL/Admin) |
| PPIC | Production Planning | Manual |
| Supervisor | Supervisor access | Manual |
| Operator | Operator access | Manual |
| Warehouse | Warehouse management | Manual |
| **Guest** | Limited access | **Sign Up (Auto)** |

---

## 🧪 Testing

### **Test Case 1: Successful Signup**
```
Input:
  Username: TestUser1
  User ID: PI0824.9001
  Password: test1234

Expected:
  ✅ Success message shown
  ✅ User created in database with role "Guest"
  ✅ Auto-switch to login mode
```

### **Test Case 2: Duplicate Username**
```
Input:
  Username: Jeremy (already exists)
  User ID: PI0824.9002
  Password: test1234

Expected:
  ❌ Error: "username already exists"
```

### **Test Case 3: Duplicate User ID**
```
Input:
  Username: NewUser
  User ID: PI0824.0001 (already exists)
  Password: test1234

Expected:
  ❌ Error: "user_id already exists"
```

### **Test Case 4: Short Password**
```
Input:
  Username: TestUser2
  User ID: PI0824.9003
  Password: abc (3 chars)

Expected:
  ❌ Error: "Password minimal 4 karakter"
```

---

## 📝 Example SQL Query

### **Check Guest Users**
```sql
SELECT 
    id,
    username,
    user_id,
    role,
    is_active,
    created_at
FROM public.users
WHERE role = 'Guest'
ORDER BY created_at DESC;
```

### **Count Users by Role**
```sql
SELECT 
    role,
    COUNT(*) as total
FROM public.users
WHERE is_active = true
GROUP BY role
ORDER BY total DESC;
```

---

## 🎯 Future Enhancements

### **Possible Improvements**
1. Email verification
2. Password strength indicator
3. CAPTCHA for bot prevention
4. Email/phone field
5. Profile picture upload
6. Role upgrade request
7. Terms & conditions checkbox
8. Password confirmation field

### **Role Management**
- Admin can upgrade Guest to other roles
- Guest users have limited permissions
- Automatic role expiry after X days

---

## 📚 Related Files

### **Frontend**
- `frontend/src/components/LoginPage.vue` - Main UI component
- `frontend/src/services/api.js` - API service

### **Backend**
- `backend/handlers/auth_handler.go` - Auth handlers
- `backend/services/auth_service.go` - Auth business logic
- `backend/repository/user_repository.go` - Database operations
- `backend/models/user.go` - User model

### **Database**
- `backend/scripts/pgadmin_complete_setup.sql` - Initial setup
- `backend/scripts/fix_passwords.sql` - Password fix script

---

## ✅ Checklist

- [x] Frontend UI - Login/Signup toggle
- [x] Frontend form validation
- [x] Frontend API integration
- [x] Backend register endpoint
- [x] Backend validation
- [x] Database schema support
- [x] Role "Guest" auto-assignment
- [x] Success/Error messages
- [x] Smooth UI transitions
- [x] Documentation

---

## 🆘 Troubleshooting

### **Problem: "username already exists"**
**Solution:** Choose a different username

### **Problem: "user_id already exists"**
**Solution:** Choose a different user_id (increment the number)

### **Problem: Signup button not responding**
**Solution:** Check browser console for errors, ensure backend is running

### **Problem: Password too short**
**Solution:** Use at least 4 characters

---

## 📞 Support

Jika ada masalah atau pertanyaan:
1. Check browser console for errors
2. Check backend terminal for logs
3. Verify database connection
4. Review this documentation

---

**Created:** October 26, 2025  
**Version:** 1.0  
**Author:** GitHub Copilot & Team
