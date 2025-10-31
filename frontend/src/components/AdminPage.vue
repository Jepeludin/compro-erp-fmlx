<template>
  <div class="admin-container">
    <!-- Header -->
    <header class="admin-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="page-title">User Management</h1>
          <p class="page-subtitle">Manage system users and permissions</p>
        </div>
        <div class="header-right">
          <div class="search-box">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"></circle>
              <path d="m21 21-4.35-4.35"></path>
            </svg>
            <input 
              v-model="searchQuery" 
              type="text" 
              placeholder="Search users..."
              class="search-input"
            />
          </div>
        </div>
      </div>
    </header>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>Loading users...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-container">
      <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="12" y1="8" x2="12" y2="12"></line>
        <line x1="12" y1="16" x2="12.01" y2="16"></line>
      </svg>
      <p class="error-message">{{ error }}</p>
      <button @click="loadUsers" class="retry-button">Retry</button>
    </div>

    <!-- Users Table -->
    <main v-else class="admin-main">
      <div class="table-header">
        <div class="table-info">
          <h2>Total Users: <span class="count">{{ filteredUsers.length }}</span></h2>
        </div>
        <div class="header-actions">
          <div class="filter-group">
            <label>Role Filter:</label>
            <select v-model="roleFilter" class="role-filter">
              <option value="">All Roles</option>
              <option value="Admin">Admin</option>
              <option value="PPIC">PPIC</option>
              <option value="Toolpather">Toolpather</option>
              <option value="PEM">PEM</option>
              <option value="QC">QC</option>
              <option value="Engineering">Engineering</option>
              <option value="Guest">Guest</option>
              <option value="Operator">Operator</option>
            </select>
          </div>
          <button @click="openAddModal" class="btn-add-user">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
              <circle cx="8.5" cy="7" r="4"></circle>
              <line x1="20" y1="8" x2="20" y2="14"></line>
              <line x1="23" y1="11" x2="17" y2="11"></line>
            </svg>
            Add New User
          </button>
        </div>
      </div>

      <div class="table-container">
        <table class="users-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Username</th>
              <th>User ID</th>
              <th>Role</th>
              <th>Operator</th>
              <th>Status</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in filteredUsers" :key="user.id" :class="{ 'inactive': !user.is_active }">
              <td>{{ user.id }}</td>
              <td class="username">{{ user.username }}</td>
              <td>{{ user.user_id }}</td>
              <td>
                <span class="role-badge" :class="'role-' + user.role.toLowerCase()">
                  {{ user.role }}
                </span>
              </td>
              <td>{{ user.operator || '-' }}</td>
              <td>
                <span class="status-badge" :class="user.is_active ? 'active' : 'inactive'">
                  {{ user.is_active ? 'Active' : 'Inactive' }}
                </span>
              </td>
              <td>{{ formatDate(user.created_at) }}</td>
              <td class="actions">
                <button @click="openEditModal(user)" class="btn-edit" title="Edit">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
                  </svg>
                </button>
                <button @click="confirmDelete(user)" class="btn-delete" title="Delete">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="filteredUsers.length === 0" class="no-results">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <path d="M16 16s-1.5-2-4-2-4 2-4 2"></path>
            <line x1="9" y1="9" x2="9.01" y2="9"></line>
            <line x1="15" y1="9" x2="15.01" y2="9"></line>
          </svg>
          <p>No users found</p>
        </div>
      </div>
    </main>

    <!-- Add User Modal -->
    <div v-if="showAddModal" class="modal-overlay" @click="closeAddModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h2>Add New User</h2>
          <button @click="closeAddModal" class="close-button">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label for="new-username">Username <span class="required">*</span></label>
            <input 
              v-model="addForm.username" 
              type="text" 
              id="new-username" 
              placeholder="Enter username"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label for="new-userid">User ID <span class="required">*</span></label>
            <input 
              v-model="addForm.user_id" 
              type="text" 
              id="new-userid" 
              placeholder="e.g., PI0824.1001"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label for="new-password">Password <span class="required">*</span></label>
            <input 
              v-model="addForm.password" 
              type="password" 
              id="new-password" 
              placeholder="Enter password"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label for="new-role">Role <span class="required">*</span></label>
            <select v-model="addForm.role" id="new-role" class="form-select">
              <option value="">Select Role</option>
              <option value="Admin">Admin</option>
              <option value="QC">QC</option>
              <option value="Toolpather">Toolpather</option>
              <option value="PEM">PEM</option>
              <option value="PPIC">PPIC</option>
              <option value="Engineering">Engineering</option>
              <option value="Guest">Guest</option>
              <option value="Operator">Operator</option>
            </select>
          </div>

          <div v-if="addForm.role === 'Operator'" class="form-group">
            <label for="new-operator">Machine <span class="required">*</span></label>
            <select v-model="addForm.operator" id="new-operator" class="form-select">
              <option value="">Select Machine</option>
              <option v-for="machine in machines" :key="machine.id" :value="machine.machine_name">
                {{ machine.machine_name }} ({{ machine.machine_code }})
              </option>
            </select>
            <p class="form-hint">Assign this operator to a specific machine</p>
          </div>

          <div v-if="addError" class="error-alert">
            {{ addError }}
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeAddModal" class="btn-cancel">Cancel</button>
          <button @click="addUser" :disabled="adding" class="btn-save">
            <span v-if="adding">Creating...</span>
            <span v-else>Create User</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="showEditModal" class="modal-overlay" @click="closeEditModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h2>Edit User</h2>
          <button @click="closeEditModal" class="close-button">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>Username</label>
            <input type="text" :value="editingUser.username" disabled class="input-disabled" />
          </div>

          <div class="form-group">
            <label>User ID</label>
            <input type="text" :value="editingUser.user_id" disabled class="input-disabled" />
          </div>

          <div class="form-group">
            <label for="role">Role</label>
            <select v-model="editForm.role" id="role" class="form-select">
              <option value="Admin">Admin</option>
              <option value="QC">QC</option>
              <option value="Toolpather">Toolpather</option>
              <option value="PEM">PEM</option>
              <option value="PPIC">PPIC</option>
              <option value="Engineering">Engineering</option>
              <option value="Guest">Guest</option>
              <option value="Operator">Operator</option>
            </select>
          </div>

          <div v-if="editForm.role === 'Operator'" class="form-group">
            <label for="edit-operator">Machine <span class="required">*</span></label>
            <select v-model="editForm.operator" id="edit-operator" class="form-select">
              <option value="">Select Machine</option>
              <option v-for="machine in machines" :key="machine.id" :value="machine.machine_name">
                {{ machine.machine_name }} ({{ machine.machine_code }})
              </option>
            </select>
            <p class="form-hint">Assign this operator to a specific machine</p>
          </div>

          <div class="form-group">
            <label for="password">New Password (leave empty to keep current)</label>
            <input 
              v-model="editForm.password" 
              type="password" 
              id="password" 
              placeholder="Enter new password"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input v-model="editForm.is_active" type="checkbox" class="form-checkbox" />
              <span>Active</span>
            </label>
          </div>

          <div v-if="editError" class="error-alert">
            {{ editError }}
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeEditModal" class="btn-cancel">Cancel</button>
          <button @click="updateUser" :disabled="updating" class="btn-save">
            <span v-if="updating">Updating...</span>
            <span v-else>Save Changes</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="modal-overlay" @click="closeDeleteModal">
      <div class="modal-content modal-small" @click.stop>
        <div class="modal-header">
          <h2>Confirm Delete</h2>
          <button @click="closeDeleteModal" class="close-button">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <p class="delete-warning">
            Are you sure you want to delete user <strong>{{ deletingUser?.username }}</strong>?
          </p>
          <p class="delete-warning-note">This action cannot be undone.</p>

          <div v-if="deleteError" class="error-alert">
            {{ deleteError }}
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeDeleteModal" class="btn-cancel">Cancel</button>
          <button @click="deleteUser" :disabled="deleting" class="btn-delete-confirm">
            <span v-if="deleting">Deleting...</span>
            <span v-else>Delete User</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api.js';

const router = useRouter();

// State
const users = ref([]);
const loading = ref(true);
const error = ref('');
const searchQuery = ref('');
const roleFilter = ref('');

// Machines data
const machines = ref([]);
const loadingMachines = ref(false);

// Add Modal
const showAddModal = ref(false);
const addForm = ref({
  username: '',
  user_id: '',
  password: '',
  role: '',
  operator: ''
});
const addError = ref('');
const adding = ref(false);

// Edit Modal
const showEditModal = ref(false);
const editingUser = ref(null);
const editForm = ref({
  role: '',
  password: '',
  operator: '',
  is_active: true
});
const editError = ref('');
const updating = ref(false);

// Delete Modal
const showDeleteModal = ref(false);
const deletingUser = ref(null);
const deleteError = ref('');
const deleting = ref(false);

// Computed
const filteredUsers = computed(() => {
  let filtered = users.value;

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(user => 
      user.username.toLowerCase().includes(query) ||
      user.user_id.toLowerCase().includes(query) ||
      (user.operator && user.operator.toLowerCase().includes(query))
    );
  }

  // Filter by role
  if (roleFilter.value) {
    filtered = filtered.filter(user => user.role === roleFilter.value);
  }

  return filtered;
});

// Methods
const loadUsers = async () => {
  try {
    loading.value = true;
    error.value = '';
    const response = await api.getAllUsers();
    users.value = response.users || [];
  } catch (err) {
    console.error('Error loading users:', err);
    error.value = err.response?.data?.error || 'Failed to load users';
  } finally {
    loading.value = false;
  }
};

const loadMachines = async () => {
  try {
    loadingMachines.value = true;
    const response = await api.getAllMachines();
    machines.value = response.machines || [];
  } catch (err) {
    console.error('Error loading machines:', err);
  } finally {
    loadingMachines.value = false;
  }
};

const openAddModal = () => {
  addForm.value = {
    username: '',
    user_id: '',
    password: '',
    role: '',
    operator: ''
  };
  addError.value = '';
  showAddModal.value = true;
};

const closeAddModal = () => {
  showAddModal.value = false;
  addForm.value = {
    username: '',
    user_id: '',
    password: '',
    role: '',
    operator: ''
  };
  addError.value = '';
};

const addUser = async () => {
  try {
    // Validation
    if (!addForm.value.username || !addForm.value.user_id || !addForm.value.password || !addForm.value.role) {
      addError.value = 'Please fill all required fields';
      return;
    }

    // Validate operator if role is Operator
    if (addForm.value.role === 'Operator' && !addForm.value.operator) {
      addError.value = 'Please select a machine for the operator';
      return;
    }

    adding.value = true;
    addError.value = '';

    const payload = {
      username: addForm.value.username,
      user_id: addForm.value.user_id,
      password: addForm.value.password,
      role: addForm.value.role,
      operator: addForm.value.operator || '-'
    };

    await api.register(payload);
    
    // Reload users
    await loadUsers();
    closeAddModal();
  } catch (err) {
    console.error('Error adding user:', err);
    addError.value = err.response?.data?.error || err.message || 'Failed to create user';
  } finally {
    adding.value = false;
  }
};

const openEditModal = (user) => {
  editingUser.value = user;
  editForm.value = {
    role: user.role,
    password: '',
    operator: user.operator || '',
    is_active: user.is_active
  };
  editError.value = '';
  showEditModal.value = true;
};

const closeEditModal = () => {
  showEditModal.value = false;
  editingUser.value = null;
  editForm.value = {
    role: '',
    password: '',
    operator: '',
    is_active: true
  };
  editError.value = '';
};

const updateUser = async () => {
  try {
    // Validate operator if role is Operator
    if (editForm.value.role === 'Operator' && !editForm.value.operator) {
      editError.value = 'Please select a machine for the operator';
      return;
    }

    updating.value = true;
    editError.value = '';

    const payload = {
      role: editForm.value.role,
      is_active: editForm.value.is_active,
      operator: editForm.value.operator || '-'
    };

    // Only include password if provided
    if (editForm.value.password) {
      payload.password = editForm.value.password;
    }

    await api.updateUser(editingUser.value.id, payload);
    
    // Reload users
    await loadUsers();
    closeEditModal();
  } catch (err) {
    console.error('Error updating user:', err);
    editError.value = err.response?.data?.error || 'Failed to update user';
  } finally {
    updating.value = false;
  }
};

const confirmDelete = (user) => {
  deletingUser.value = user;
  deleteError.value = '';
  showDeleteModal.value = true;
};

const closeDeleteModal = () => {
  showDeleteModal.value = false;
  deletingUser.value = null;
  deleteError.value = '';
};

const deleteUser = async () => {
  try {
    deleting.value = true;
    deleteError.value = '';

    await api.deleteUser(deletingUser.value.id);
    
    // Reload users
    await loadUsers();
    closeDeleteModal();
  } catch (err) {
    console.error('Error deleting user:', err);
    deleteError.value = err.response?.data?.error || 'Failed to delete user';
  } finally {
    deleting.value = false;
  }
};

const formatDate = (dateString) => {
  if (!dateString) return '-';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

// Load users on mount
onMounted(() => {
  const user = api.getCurrentUser();
  if (!user || user.role !== 'Admin') {
    router.push('/dashboard');
    return;
  }
  loadUsers();
  loadMachines(); // Load machines for operator dropdown
});
</script>

<style scoped>
@import './AdminPage.css';
</style>
