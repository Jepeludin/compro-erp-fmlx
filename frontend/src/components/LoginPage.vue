<template>
  <div class="login-container">
    <form @submit.prevent="isSignupMode ? handleSignup : handleLogin" class="login-form">
      <h2 class="login-title">{{ isSignupMode ? 'Sign Up' : 'Login' }}</h2>
      
      <!-- Success Message -->
      <div v-if="successMessage" class="success-message">
        {{ successMessage }}
      </div>
      
      <!-- Error Message -->
      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
      
      <!-- Username field (only for signup) -->
      <div v-if="isSignupMode" class="form-group username-group">
        <label for="username">Username</label>
        <input type="text" id="username" v-model="username" :disabled="isLoading" placeholder="e.g., John, Sarah" />
      </div>
      
      <div class="form-group userid-group">
        <label for="user_id">User ID</label>
        <input type="text" id="user_id" v-model="userId" :disabled="isLoading" :placeholder="isSignupMode ? 'e.g., PI0824.5001' : 'e.g., PI0824.0001, PI0824.2374'" />
      </div>
      <div class="form-group password-group">
        <label for="password">Password</label>
        <input type="password" id="password" v-model="password" :disabled="isLoading" :placeholder="isSignupMode ? 'Create a password' : ''" />
      </div>
      
      <button type="button" @click="handleFormSubmit" class="login-button" :class="{ loading: isLoading }" :disabled="isLoading">
        <span v-if="!isLoading">{{ isSignupMode ? 'SIGN UP' : 'LOGIN' }}</span>
        <span v-else>Loading...</span>
      </button>
      
      <!-- Toggle between Login and Signup -->
      <div class="toggle-mode">
        <span v-if="!isSignupMode">
          Don't have an account? 
          <a href="#" @click.prevent="toggleMode" class="toggle-link">Sign Up</a>
        </span>
        <span v-else>
          Already have an account? 
          <a href="#" @click.prevent="toggleMode" class="toggle-link">Login</a>
        </span>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api.js';

const router = useRouter();

// Data reaktif untuk menampung input pengguna
const userId = ref('');
const username = ref('');
const password = ref('');
const isLoading = ref(false);
const errorMessage = ref('');
const successMessage = ref('');
const isSignupMode = ref(false);

// Store original body styles
let originalBodyStyle = '';
let originalAppStyle = '';
let originalAppLayoutStyle = '';

// Apply fullscreen styles when component mounts
onMounted(() => {
  const body = document.body;
  const app = document.getElementById('app');
  const appLayout = document.getElementById('app-layout');
  
  // Store original styles
  originalBodyStyle = body.getAttribute('style') || '';
  if (app) originalAppStyle = app.getAttribute('style') || '';
  if (appLayout) originalAppLayoutStyle = appLayout.getAttribute('style') || '';
  
  // Apply fullscreen styles
  body.style.cssText = 'margin: 0 !important; padding: 0 !important; overflow: hidden !important; display: block !important;';
  if (app) {
    app.style.cssText = 'max-width: none !important; margin: 0 !important; padding: 0 !important; text-align: left !important; width: 100vw !important; height: 100vh !important; display: block !important; place-items: initial !important;';
  }
  if (appLayout) {
    appLayout.style.cssText = 'max-width: none !important; padding: 0 !important; margin: 0 !important; text-align: left !important; width: 100vw !important; height: 100vh !important; display: block !important;';
  }
});

// Restore original styles when component unmounts
onUnmounted(() => {
  const body = document.body;
  const app = document.getElementById('app');
  const appLayout = document.getElementById('app-layout');
  
  // Restore original styles
  body.setAttribute('style', originalBodyStyle);
  if (app) app.setAttribute('style', originalAppStyle);
  if (appLayout) appLayout.setAttribute('style', originalAppLayoutStyle);
});

// Toggle between Login and Signup mode
const toggleMode = () => {
  console.log('Toggle mode clicked');
  isSignupMode.value = !isSignupMode.value;
  errorMessage.value = '';
  successMessage.value = '';
  // Clear form fields when switching modes
  userId.value = '';
  username.value = '';
  password.value = '';
};

// Handle form submission
const handleFormSubmit = () => {
  console.log('Form submit clicked, isSignupMode:', isSignupMode.value);
  if (isSignupMode.value) {
    handleSignup();
  } else {
    handleLogin();
  }
};

// Fungsi yang dipanggil saat form login disubmit
const handleLogin = async () => {
  console.log('Login button clicked');
  console.log('User ID:', userId.value);
  console.log('Password:', password.value ? '***' : 'empty');
  
  // Reset messages
  errorMessage.value = '';
  successMessage.value = '';
  
  // Validasi input
  if (!userId.value || !password.value) {
    errorMessage.value = 'User ID dan Password harus diisi';
    return;
  }
  
  // Mulai loading
  isLoading.value = true;
  
  try {
    // Kirim request login ke backend
    const response = await api.login(userId.value, password.value);
    
    console.log('Login successful:', response);
    
    // Simpan token dan user data
    api.saveAuth(response.data.token, response.data.user);
    
    // Redirect ke dashboard
    router.push('/dashboard');
    
  } catch (error) {
    console.error('Login error:', error);
    errorMessage.value = error.message || 'Login gagal! Periksa User ID dan Password Anda.';
  } finally {
    // Selesai loading
    isLoading.value = false;
  }
};

// Fungsi yang dipanggil saat form signup disubmit
const handleSignup = async () => {
  console.log('Signup button clicked');
  console.log('Username:', username.value);
  console.log('User ID:', userId.value);
  console.log('Password:', password.value ? '***' : 'empty');
  
  // Reset messages
  errorMessage.value = '';
  successMessage.value = '';
  
  // Validasi input
  if (!username.value || !userId.value || !password.value) {
    errorMessage.value = 'Semua field harus diisi';
    return;
  }
  
  // Validasi panjang password
  if (password.value.length < 4) {
    errorMessage.value = 'Password minimal 4 karakter';
    return;
  }
  
  // Mulai loading
  isLoading.value = true;
  
  try {
    // Kirim request register ke backend dengan role "Guest"
    const response = await api.register({
      username: username.value,
      user_id: userId.value,
      password: password.value,
      role: 'Guest', // Role otomatis Guest
      operator: '[null]' // Set operator dengan tanda strip agar tidak null
    });
    
    console.log('Registration successful:', response);
    
    // Tampilkan pesan sukses
    successMessage.value = `Akun berhasil dibuat! Selamat datang, ${response.data.username}`;
    
    // Clear form
    username.value = '';
    userId.value = '';
    password.value = '';
    
    // Setelah 2 detik, pindah ke mode login
    setTimeout(() => {
      successMessage.value = 'Silakan login dengan akun baru Anda';
      setTimeout(() => {
        toggleMode();
      }, 1500);
    }, 2000);
    
  } catch (error) {
    console.error('Signup error:', error);
    errorMessage.value = error.message || 'Pendaftaran gagal! User ID atau Username mungkin sudah digunakan.';
  } finally {
    // Selesai loading
    isLoading.value = false;
  }
};
</script>

<style scoped>
/* Modern Login Page Design */
.login-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 50%, #7e22ce 100%);
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  padding: 20px;
  box-sizing: border-box;
  overflow: hidden;
}

/* Animated background elements */
.login-container::before {
  content: '';
  position: absolute;
  width: 400px;
  height: 400px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 50%;
  top: -200px;
  left: -200px;
  animation: float 20s infinite ease-in-out;
  pointer-events: none;
  z-index: 1;
}

.login-container::after {
  content: '';
  position: absolute;
  width: 300px;
  height: 300px;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 50%;
  bottom: -150px;
  right: -150px;
  animation: float 15s infinite ease-in-out reverse;
  pointer-events: none;
  z-index: 1;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(50px, 50px) scale(1.1);
  }
}

.login-form {
  position: relative;
  width: 100%;
  max-width: 440px;
  padding: 3.5rem 3rem;
  background: rgba(255, 255, 255, 0.98);
  border-radius: 24px;
  box-shadow: 0 30px 60px rgba(0, 0, 0, 0.3),
              0 0 0 1px rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  z-index: 10;
  animation: slideUp 0.6s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-title {
  text-align: center;
  color: #1e293b;
  margin-bottom: 0.5rem;
  font-size: 2.5rem;
  font-weight: 700;
  letter-spacing: -0.5px;
  position: relative;
}

.login-title::after {
  content: 'Welcome Back';
  display: block;
  font-size: 1rem;
  font-weight: 400;
  color: #64748b;
  margin-top: 0.5rem;
  letter-spacing: 0;
}

.form-group {
  margin-bottom: 1.8rem;
  position: relative;
}

.form-group label {
  display: block;
  margin-bottom: 0.6rem;
  color: #334155;
  font-weight: 600;
  font-size: 0.9rem;
  letter-spacing: 0.3px;
  text-transform: uppercase;
}

.form-group input {
  width: 100%;
  padding: 1rem 1.2rem;
  padding-left: 3.2rem;
  border: 2px solid #e2e8f0;
  border-radius: 14px;
  box-sizing: border-box;
  font-size: 1rem;
  background: #f8fafc;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  outline: none;
  color: #1e293b;
  font-weight: 500;
}

.form-group input:focus {
  border-color: #3b82f6;
  background: #ffffff;
  box-shadow: 0 0 0 4px rgba(59, 130, 246, 0.1);
  transform: translateY(-1px);
}

.form-group input::placeholder {
  color: #94a3b8;
}

/* Icon styling */
.form-group::before {
  content: '';
  position: absolute;
  left: 1.1rem;
  top: 2.8rem;
  width: 20px;
  height: 20px;
  background-size: contain;
  background-repeat: no-repeat;
  opacity: 0.5;
  transition: opacity 0.3s;
  z-index: 1;
}

.form-group:focus-within::before {
  opacity: 0.8;
}

.userid-group::before {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='%233b82f6' viewBox='0 0 24 24'%3E%3Cpath d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/%3E%3C/svg%3E");
}

.password-group::before {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='%233b82f6' viewBox='0 0 24 24'%3E%3Cpath d='M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zM12 17c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zM15.1 8H8.9V6c0-1.71 1.39-3.1 3.1-3.1s3.1 1.39 3.1 3.1v2z'/%3E%3C/svg%3E");
}

.username-group::before {
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='%233b82f6' viewBox='0 0 24 24'%3E%3Cpath d='M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z'/%3E%3C/svg%3E");
}

.login-button {
  width: 100%;
  padding: 1.1rem;
  border: none;
  border-radius: 14px;
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  color: white;
  font-size: 1.05rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  margin-top: 0.5rem;
  position: relative;
  overflow: hidden;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  box-shadow: 0 4px 15px rgba(59, 130, 246, 0.4);
  pointer-events: auto;
  z-index: 10;
}

.login-button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  transition: left 0.5s;
  pointer-events: none;
}

.login-button:hover::before {
  left: 100%;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(59, 130, 246, 0.5);
}

.login-button:active {
  transform: translateY(0);
  box-shadow: 0 4px 15px rgba(59, 130, 246, 0.4);
}

/* Loading state */
.login-button.loading {
  pointer-events: none;
  opacity: 0.8;
  background: linear-gradient(135deg, #64748b 0%, #475569 100%);
}

.login-button.loading::after {
  content: '';
  position: absolute;
  width: 20px;
  height: 20px;
  margin: auto;
  border: 3px solid transparent;
  border-top-color: #ffffff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

@keyframes spin {
  0% { transform: translate(-50%, -50%) rotate(0deg); }
  100% { transform: translate(-50%, -50%) rotate(360deg); }
}

/* Error message styling */
.error-message {
  background-color: #fee;
  border: 1px solid #fcc;
  color: #c33;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
  font-size: 0.9rem;
  text-align: center;
  animation: shake 0.3s ease-in-out;
}

/* Success message styling */
.success-message {
  background-color: #d4edda;
  border: 1px solid #c3e6cb;
  color: #155724;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  margin-bottom: 1rem;
  font-size: 0.9rem;
  text-align: center;
  animation: slideDown 0.4s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-10px); }
  75% { transform: translateX(10px); }
}

/* Responsive */
@media (max-width: 480px) {
  .login-form {
    padding: 2.5rem 2rem;
    max-width: 95%;
  }
  
  .login-title {
    font-size: 2rem;
  }
}

/* Accessibility - Focus visible */
.login-button:focus-visible {
  outline: 3px solid #3b82f6;
  outline-offset: 2px;
}

.form-group input:focus-visible {
  outline: none;
}

/* Toggle mode link */
.toggle-mode {
  text-align: center;
  margin-top: 1.5rem;
  font-size: 0.95rem;
  color: #64748b;
}

.toggle-link {
  color: #3b82f6;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s;
  position: relative;
}

.toggle-link::after {
  content: '';
  position: absolute;
  width: 0;
  height: 2px;
  bottom: -2px;
  left: 0;
  background-color: #3b82f6;
  transition: width 0.3s;
}

.toggle-link:hover::after {
  width: 100%;
}

.toggle-link:hover {
  color: #2563eb;
}
</style>