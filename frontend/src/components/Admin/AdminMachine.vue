<template>
  <div class="machine-container">
    <!-- Header -->
    <header class="machine-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="page-title">Machine Management</h1>
          <p class="page-subtitle">Manage production machines</p>
        </div>
        <div class="header-right">
          <button @click="openAddModal" class="btn-add-machine">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="12" y1="8" x2="12" y2="16"></line>
              <line x1="8" y1="12" x2="16" y2="12"></line>
            </svg>
            Add New Machine
          </button>
        </div>
      </div>
    </header>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>Loading machines...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-container">
      <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="12" y1="8" x2="12" y2="12"></line>
        <line x1="12" y1="16" x2="12.01" y2="16"></line>
      </svg>
      <p class="error-message">{{ error }}</p>
      <button @click="loadMachines" class="retry-button">Retry</button>
    </div>

    <!-- Machines Grid -->
    <main v-else class="machine-main">
      <div class="machines-info">
        <h2>Total Machines: <span class="count">{{ machines.length }}</span></h2>
      </div>

      <div class="machines-grid">
        <div v-for="machine in machines" :key="machine.id" class="machine-card">
          <div class="machine-card-header">
            <div class="machine-icon">
              <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="2" y="7" width="20" height="14" rx="2" ry="2"></rect>
                <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"></path>
              </svg>
            </div>
            <span class="status-badge" :class="'status-' + machine.status">{{ machine.status }}</span>
          </div>
          <div class="machine-card-body">
            <h3 class="machine-name">{{ machine.machine_name }}</h3>
            <p class="machine-code">{{ machine.machine_code }}</p>
            <div class="machine-details">
              <div class="detail-item">
                <span class="detail-label">Type:</span>
                <span class="detail-value">{{ machine.machine_type || '-' }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">Location:</span>
                <span class="detail-value">{{ machine.location || '-' }}</span>
              </div>
            </div>
          </div>
          <div class="machine-card-footer">
            <button @click="openEditModal(machine)" class="btn-edit">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
                <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
              </svg>
              Edit
            </button>
            <button @click="openDeleteModal(machine)" class="btn-delete">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
              </svg>
              Delete
            </button>
          </div>
        </div>

        <!-- Empty State -->
        <div v-if="machines.length === 0" class="empty-state">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="2" y="7" width="20" height="14" rx="2" ry="2"></rect>
            <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"></path>
          </svg>
          <h3>No Machines Found</h3>
          <p>Click "Add New Machine" to create your first machine</p>
        </div>
      </div>
    </main>

    <!-- Add Machine Modal -->
    <div v-if="showAddModal" class="modal-overlay" @click="closeAddModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h2>Add New Machine</h2>
          <button @click="closeAddModal" class="close-button">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label for="new-machine-code">Machine Code <span class="required">*</span></label>
            <input 
              v-model="addForm.machine_code" 
              type="text" 
              id="new-machine-code" 
              placeholder="e.g., YSD01"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label for="new-machine-name">Machine Name <span class="required">*</span></label>
            <input 
              v-model="addForm.machine_name" 
              type="text" 
              id="new-machine-name" 
              placeholder="e.g., Yasda CNC"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label for="new-machine-type">Machine Type</label>
            <input 
              v-model="addForm.machine_type" 
              type="text" 
              id="new-machine-type" 
              placeholder="e.g., CNC"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label for="new-location">Location</label>
            <input 
              v-model="addForm.location" 
              type="text" 
              id="new-location" 
              placeholder="e.g., Workshop A"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label for="new-status">Status</label>
            <select v-model="addForm.status" id="new-status" class="form-select">
              <option value="active">Active</option>
              <option value="maintenance">Maintenance</option>
              <option value="inactive">Inactive</option>
            </select>
          </div>

          <div v-if="addError" class="error-alert">
            {{ addError }}
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeAddModal" class="btn-cancel">Cancel</button>
          <button @click="addMachine" :disabled="adding" class="btn-save">
            <span v-if="adding">Creating...</span>
            <span v-else>Create Machine</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Machine Modal -->
    <div v-if="showEditModal" class="modal-overlay" @click="closeEditModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h2>Edit Machine</h2>
          <button @click="closeEditModal" class="close-button">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>Machine Code</label>
            <input type="text" :value="editingMachine.machine_code" disabled class="input-disabled" />
          </div>

          <div class="form-group">
            <label for="edit-machine-name">Machine Name</label>
            <input 
              v-model="editForm.machine_name" 
              type="text" 
              id="edit-machine-name" 
              placeholder="Enter machine name"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label for="edit-machine-type">Machine Type</label>
            <input 
              v-model="editForm.machine_type" 
              type="text" 
              id="edit-machine-type" 
              placeholder="Enter machine type"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label for="edit-location">Location</label>
            <input 
              v-model="editForm.location" 
              type="text" 
              id="edit-location" 
              placeholder="Enter location"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label for="edit-status">Status</label>
            <select v-model="editForm.status" id="edit-status" class="form-select">
              <option value="active">Active</option>
              <option value="maintenance">Maintenance</option>
              <option value="inactive">Inactive</option>
            </select>
          </div>

          <div v-if="editError" class="error-alert">
            {{ editError }}
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeEditModal" class="btn-cancel">Cancel</button>
          <button @click="updateMachine" :disabled="updating" class="btn-save">
            <span v-if="updating">Updating...</span>
            <span v-else>Save Changes</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Machine Modal -->
    <div v-if="showDeleteModal" class="modal-overlay" @click="closeDeleteModal">
      <div class="modal-content modal-small" @click.stop>
        <div class="modal-header">
          <h2>Delete Machine</h2>
          <button @click="closeDeleteModal" class="close-button">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <div class="delete-warning">
            <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="8" x2="12" y2="12"></line>
              <line x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <h3>Are you sure?</h3>
            <p>Do you want to delete machine <strong>{{ deletingMachine?.machine_name }}</strong>?</p>
            <p class="warning-text">This action cannot be undone.</p>
          </div>

          <div v-if="deleteError" class="error-alert">
            {{ deleteError }}
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeDeleteModal" class="btn-cancel">Cancel</button>
          <button @click="deleteMachine" :disabled="deleting" class="btn-delete-confirm">
            <span v-if="deleting">Deleting...</span>
            <span v-else>Delete Machine</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../../services/api';

// State
const machines = ref([]);
const loading = ref(false);
const error = ref('');

// Add Modal
const showAddModal = ref(false);
const addForm = ref({
  machine_code: '',
  machine_name: '',
  machine_type: '',
  location: '',
  status: 'active'
});
const addError = ref('');
const adding = ref(false);

// Edit Modal
const showEditModal = ref(false);
const editingMachine = ref(null);
const editForm = ref({
  machine_name: '',
  machine_type: '',
  location: '',
  status: 'active'
});
const editError = ref('');
const updating = ref(false);

// Delete Modal
const showDeleteModal = ref(false);
const deletingMachine = ref(null);
const deleteError = ref('');
const deleting = ref(false);

// Load machines
const loadMachines = async () => {
  try {
    loading.value = true;
    error.value = '';
    const response = await api.getAllMachines();
    machines.value = response.machines || [];
  } catch (err) {
    console.error('Error loading machines:', err);
    error.value = err.message || 'Failed to load machines';
  } finally {
    loading.value = false;
  }
};

// Add machine methods
const openAddModal = () => {
  addForm.value = {
    machine_code: '',
    machine_name: '',
    machine_type: '',
    location: '',
    status: 'active'
  };
  addError.value = '';
  showAddModal.value = true;
};

const closeAddModal = () => {
  showAddModal.value = false;
  addForm.value = {
    machine_code: '',
    machine_name: '',
    machine_type: '',
    location: '',
    status: 'active'
  };
  addError.value = '';
};

const addMachine = async () => {
  try {
    // Validation
    if (!addForm.value.machine_code || !addForm.value.machine_name) {
      addError.value = 'Machine code and name are required';
      return;
    }

    adding.value = true;
    addError.value = '';

    await api.createMachine(addForm.value);
    
    // Reload machines
    await loadMachines();
    closeAddModal();
  } catch (err) {
    console.error('Error adding machine:', err);
    addError.value = err.message || 'Failed to create machine';
  } finally {
    adding.value = false;
  }
};

// Edit machine methods
const openEditModal = (machine) => {
  editingMachine.value = machine;
  editForm.value = {
    machine_name: machine.machine_name,
    machine_type: machine.machine_type,
    location: machine.location,
    status: machine.status
  };
  editError.value = '';
  showEditModal.value = true;
};

const closeEditModal = () => {
  showEditModal.value = false;
  editingMachine.value = null;
  editForm.value = {
    machine_name: '',
    machine_type: '',
    location: '',
    status: 'active'
  };
  editError.value = '';
};

const updateMachine = async () => {
  try {
    updating.value = true;
    editError.value = '';

    await api.updateMachine(editingMachine.value.id, editForm.value);
    
    // Reload machines
    await loadMachines();
    closeEditModal();
  } catch (err) {
    console.error('Error updating machine:', err);
    editError.value = err.message || 'Failed to update machine';
  } finally {
    updating.value = false;
  }
};

// Delete machine methods
const openDeleteModal = (machine) => {
  deletingMachine.value = machine;
  deleteError.value = '';
  showDeleteModal.value = true;
};

const closeDeleteModal = () => {
  showDeleteModal.value = false;
  deletingMachine.value = null;
  deleteError.value = '';
};

const deleteMachine = async () => {
  try {
    deleting.value = true;
    deleteError.value = '';

    await api.deleteMachine(deletingMachine.value.id);
    
    // Reload machines
    await loadMachines();
    closeDeleteModal();
  } catch (err) {
    console.error('Error deleting machine:', err);
    deleteError.value = err.message || 'Failed to delete machine';
  } finally {
    deleting.value = false;
  }
};

// Load machines on mount
onMounted(() => {
  loadMachines();
});
</script>

<style src="./AdminMachine.css"></style>