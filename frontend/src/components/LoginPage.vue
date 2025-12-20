<template>
  <div class="login-container">
    <div class="animated-bg"></div>
    <form @submit.prevent="handleFormSubmit" class="login-form">
      <h2 class="login-title">{{ isSignupMode ? 'Sign Up' : 'Login' }}</h2>
      
      <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
      <div v-if="errorMessage" class="error-message">{{ errorMessage }}</div>
      
      <div v-if="isSignupMode" class="form-group username-group">
        <label for="username">Username</label>
        <input type="text" id="username" v-model="username" :disabled="isLoading" placeholder="e.g., John Doe" />
      </div>
      
      <div class="form-group userid-group">
        <label for="user_id">User ID</label>
        <input type="text" id="user_id" v-model="userId" :disabled="isLoading" :placeholder="isSignupMode ? 'e.g., PI0824.5001' : 'Enter your ID'" />
      </div>

      <div class="form-group password-group">
        <label for="password">Password</label>
        <input type="password" id="password" v-model="password" :disabled="isLoading" />
      </div>
      
      <button type="submit" class="login-button" :class="{ loading: isLoading }" :disabled="isLoading">
        <span v-if="!isLoading">{{ isSignupMode ? 'SIGN UP' : 'LOGIN' }}</span>
        <span v-else>Loading...</span>
      </button>
      
      <div class="toggle-mode">
        <span>
          {{ isSignupMode ? 'Already have an account?' : "Don't have an account?" }}
          <a href="#" @click.prevent="toggleMode" class="toggle-link">
            {{ isSignupMode ? 'Login' : 'Sign Up' }}
          </a>
        </span>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api.js';
// Impor CSS di sini jika tidak diimpor secara global
import './LoginPage.css';

const router = useRouter();
const userId = ref('');
const username = ref('');
const password = ref('');
const isLoading = ref(false);
const errorMessage = ref('');
const successMessage = ref('');
const isSignupMode = ref(false);

let originalStyles = { body: '', app: '' };

onMounted(() => {
  const body = document.body;
  const app = document.getElementById('app');
  
  originalStyles.body = body.style.cssText;
  if (app) originalStyles.app = app.style.cssText;
  
  body.classList.add('auth-page-active');
});

onUnmounted(() => {
  document.body.classList.remove('auth-page-active');
});

const toggleMode = () => {
  isSignupMode.value = !isSignupMode.value;
  errorMessage.value = '';
  successMessage.value = '';
  userId.value = '';
  username.value = '';
  password.value = '';
};

const handleFormSubmit = () => {
  isSignupMode.value ? handleSignup() : handleLogin();
};

const handleLogin = async () => {
  if (!userId.value || !password.value) {
    errorMessage.value = 'User ID dan Password harus diisi';
    return;
  }
  isLoading.value = true;
  try {
    const response = await api.login(userId.value, password.value);
    api.saveAuth(response.data.token, response.data.user);
    router.push('/dashboard');
  } catch (error) {
    errorMessage.value = error.message || 'Login gagal!';
  } finally {
    isLoading.value = false;
  }
};

const handleSignup = async () => {
  if (!username.value || !userId.value || !password.value) {
    errorMessage.value = 'Semua field harus diisi';
    return;
  }
  isLoading.value = true;
  try {
    const response = await api.register({
      username: username.value,
      user_id: userId.value,
      password: password.value,
      role: 'Guest',
      operator: '[null]'
    });
    successMessage.value = `Akun berhasil dibuat!`;
    setTimeout(() => toggleMode(), 2000);
  } catch (error) {
    errorMessage.value = error.message || 'Pendaftaran gagal!';
  } finally {
    isLoading.value = false;
  }
};
</script>