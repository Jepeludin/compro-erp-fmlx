<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-container">
      <div class="modal-header">
        <h2>Select Approver</h2>
        <button @click="$emit('close')" class="btn-close">&times;</button>
      </div>

      <div class="modal-body">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by name or email..."
          class="search-input"
        />

        <div v-if="loading" class="loading">Loading users...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="filteredUsers.length === 0" class="no-data">No users found</div>

        <div v-else class="users-list">
          <div
            v-for="user in filteredUsers"
            :key="user.id"
            @click="selectUser(user)"
            class="user-card"
          >
            <div class="user-info">
              <div class="user-name">{{ user.username }}</div>
              <div class="user-email">{{ user.email }}</div>
            </div>
            <div class="user-role">
              <span class="role-badge">{{ user.role }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import api from '../../services/api.js';

const emit = defineEmits(['close', 'select']);

const users = ref([]);
const loading = ref(false);
const error = ref(null);
const searchQuery = ref('');

const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value;

  const query = searchQuery.value.toLowerCase();
  return users.value.filter(user =>
    user.username?.toLowerCase().includes(query) ||
    user.email?.toLowerCase().includes(query) ||
    user.role?.toLowerCase().includes(query)
  );
});

async function loadUsers() {
  loading.value = true;
  error.value = null;

  try {
    const response = await api.getUsers();
    if (response.success && response.users) {
      // Users are already filtered on backend (active users, excluding Guest)
      users.value = response.users;
    }
  } catch (err) {
    error.value = 'Failed to load users: ' + err.message;
    console.error('Error loading users:', err);
  } finally {
    loading.value = false;
  }
}

function selectUser(user) {
  emit('select', user);
}

onMounted(() => {
  loadUsers();
});
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1001;
}

.modal-container {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 600px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h2 {
  margin: 0;
  color: #1f2937;
  font-size: 1.5rem;
}

.btn-close {
  background: none;
  border: none;
  font-size: 2rem;
  color: #6b7280;
  cursor: pointer;
  line-height: 1;
}

.btn-close:hover {
  color: #1f2937;
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem 2rem;
}

.search-input {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 0.95rem;
  margin-bottom: 1rem;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
}

.users-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.user-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.user-card:hover {
  border-color: #667eea;
  background: #f9fafb;
}

.user-info {
  flex: 1;
}

.user-name {
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 0.25rem;
}

.user-email {
  color: #6b7280;
  font-size: 0.9rem;
}

.user-role {
  margin-left: 1rem;
}

.role-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  background: #dbeafe;
  color: #1e40af;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 500;
}

.loading,
.error,
.no-data {
  text-align: center;
  padding: 2rem;
  color: #6b7280;
}

.error {
  color: #dc2626;
}
</style>
