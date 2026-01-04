<template>
  <div class="file-list-container">
    <div class="list-header">
      <h3>{{ showMyFilesOnly ? 'My Uploaded Files' : 'All Files' }}</h3>
      <div class="header-actions">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by order number or file name..."
          class="search-input"
        />
        <button @click="loadFiles" class="btn-refresh" :disabled="loading">
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="23 4 23 10 17 10"></polyline>
            <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path>
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading">Loading files...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="filteredFiles.length === 0" class="no-data">
      {{ showMyFilesOnly ? 'You have not uploaded any files yet' : 'No files found' }}
    </div>

    <div v-else class="files-grid">
      <div v-for="file in filteredFiles" :key="file.id" class="file-card">
        <div class="card-header">
          <div class="file-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
              <polyline points="13 2 13 9 20 9"></polyline>
            </svg>
          </div>
          <div class="file-info">
            <h4 class="file-name">{{ file.file_name }}</h4>
            <p class="file-size">{{ formatFileSize(file.file_size) }}</p>
          </div>
        </div>

        <div class="card-body">
          <div class="info-row">
            <span class="label">Order Number:</span>
            <span class="value">{{ file.order_number }}</span>
          </div>
          <div v-if="file.part_name" class="info-row">
            <span class="label">Part Name:</span>
            <span class="value">{{ file.part_name }}</span>
          </div>
          <div class="info-row">
            <span class="label">Uploaded By:</span>
            <span class="value">{{ file.uploader?.username || 'Unknown' }}</span>
          </div>
          <div class="info-row">
            <span class="label">Upload Date:</span>
            <span class="value">{{ formatDate(file.created_at) }}</span>
          </div>
          <div v-if="file.notes" class="info-row">
            <span class="label">Notes:</span>
            <span class="value">{{ file.notes }}</span>
          </div>
        </div>

        <div class="card-actions">
          <button @click="downloadFile(file)" class="btn-download">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
              <polyline points="7 10 12 15 17 10"></polyline>
              <line x1="12" y1="15" x2="12" y2="3"></line>
            </svg>
            Download
          </button>
          <button v-if="canDelete(file)" @click="deleteFile(file)" class="btn-delete">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import api from '../../services/api.js';

const props = defineProps({
  showMyFilesOnly: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['refresh']);

const files = ref([]);
const loading = ref(false);
const error = ref(null);
const searchQuery = ref('');
const currentUser = ref(null);

const filteredFiles = computed(() => {
  if (!searchQuery.value) return files.value;

  const query = searchQuery.value.toLowerCase();
  return files.value.filter(file =>
    file.order_number?.toLowerCase().includes(query) ||
    file.file_name?.toLowerCase().includes(query) ||
    file.part_name?.toLowerCase().includes(query)
  );
});

async function loadFiles() {
  loading.value = true;
  error.value = null;

  try {
    let response;

    if (props.showMyFilesOnly) {
      response = await api.getMyToolpatherFiles();
    } else {
      response = await api.getAllToolpatherFiles();
    }

    if (response.success && response.data) {
      files.value = response.data;
    }
  } catch (err) {
    error.value = 'Failed to load files: ' + err.message;
    console.error('Error loading files:', err);
  } finally {
    loading.value = false;
  }
}

function downloadFile(file) {
  const downloadUrl = api.getToolpatherFileDownloadUrl(file.id);

  // Create temporary link and trigger download
  const link = document.createElement('a');
  link.href = downloadUrl;
  link.download = file.file_name;

  // Add auth token to request
  const token = api.getToken();
  if (token) {
    link.setAttribute('data-token', token);
  }

  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

async function deleteFile(file) {
  const confirmDelete = confirm(
    `Are you sure you want to delete "${file.file_name}"?\nThis action cannot be undone.`
  );

  if (!confirmDelete) return;

  try {
    const response = await api.deleteToolpatherFile(file.id);

    if (response.success) {
      alert('File deleted successfully');
      loadFiles();
      emit('refresh');
    }
  } catch (err) {
    alert('Failed to delete file: ' + err.message);
    console.error('Delete error:', err);
  }
}

function canDelete(file) {
  if (!currentUser.value) return false;

  // Admin can delete any file
  if (currentUser.value.role === 'Admin') return true;

  // Uploader can delete their own files
  return file.uploaded_by === currentUser.value.id;
}

function formatFileSize(bytes) {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

function formatDate(dateString) {
  if (!dateString) return '-';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
}

onMounted(() => {
  currentUser.value = api.getCurrentUser();
  loadFiles();
});
</script>

<style scoped>
.file-list-container {
  width: 100%;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.list-header h3 {
  margin: 0;
  color: #1f2937;
  font-size: 1.5rem;
}

.header-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.search-input {
  padding: 0.6rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  width: 300px;
  font-size: 0.95rem;
}

.btn-refresh {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1.2rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.95rem;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-refresh:hover:not(:disabled) {
  background: #5568d3;
}

.btn-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.loading, .error, .no-data {
  text-align: center;
  padding: 3rem;
  color: #6b7280;
  font-size: 1.1rem;
}

.error {
  color: #dc2626;
}

.files-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.file-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s;
}

.file-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: #667eea;
}

.card-header {
  display: flex;
  gap: 1rem;
  padding: 1.5rem;
  background: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
}

.file-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.file-icon svg {
  color: white;
}

.file-info {
  flex: 1;
  min-width: 0;
}

.file-name {
  margin: 0 0 0.25rem;
  color: #1f2937;
  font-size: 1rem;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  margin: 0;
  color: #6b7280;
  font-size: 0.875rem;
}

.card-body {
  padding: 1.5rem;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 0.5rem 0;
  border-bottom: 1px solid #f3f4f6;
}

.info-row:last-child {
  border-bottom: none;
}

.info-row .label {
  color: #6b7280;
  font-size: 0.875rem;
  font-weight: 500;
}

.info-row .value {
  color: #1f2937;
  font-size: 0.875rem;
  text-align: right;
  max-width: 60%;
  overflow: hidden;
  text-overflow: ellipsis;
}

.card-actions {
  display: flex;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  background: #f9fafb;
  border-top: 1px solid #e5e7eb;
}

.btn-download, .btn-delete {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.6rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-download {
  background: #667eea;
  color: white;
}

.btn-download:hover {
  background: #5568d3;
}

.btn-delete {
  background: #fee2e2;
  color: #dc2626;
}

.btn-delete:hover {
  background: #fecaca;
}
</style>
