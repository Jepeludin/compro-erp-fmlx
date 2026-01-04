<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-container">
      <div class="modal-header">
        <h2>Upload Toolpather Files</h2>
        <button @click="$emit('close')" class="btn-close">&times;</button>
      </div>

      <div class="modal-body">
        <form @submit.prevent="handleUpload">
          <!-- Order Number -->
          <div class="form-group">
            <label>Order Number (NJO) <span class="required">*</span></label>
            <input
              v-model="formData.orderNumber"
              type="text"
              placeholder="Enter order number..."
              class="form-input"
              @blur="lookupPartName"
              required
            />
            <small v-if="lookingUpPartName" class="help-text">Looking up part name...</small>
          </div>

          <!-- Part Name (auto-filled) -->
          <div class="form-group">
            <label>Part Name</label>
            <input
              v-model="formData.partName"
              type="text"
              placeholder="Auto-filled from order number"
              class="form-input"
              readonly
            />
          </div>

          <!-- File Upload -->
          <div class="form-group">
            <label>Files (.txt only) <span class="required">*</span></label>
            <div class="file-upload-area" @click="triggerFileInput">
              <input
                ref="fileInput"
                type="file"
                multiple
                accept=".txt"
                @change="handleFileSelect"
                style="display: none"
              />
              <div v-if="selectedFiles.length === 0" class="upload-placeholder">
                <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                  <polyline points="17 8 12 3 7 8"></polyline>
                  <line x1="12" y1="3" x2="12" y2="15"></line>
                </svg>
                <p>Click to select .txt files or drag and drop</p>
                <small>You can select multiple files at once</small>
              </div>
              <div v-else class="selected-files">
                <div v-for="(file, index) in selectedFiles" :key="index" class="file-item">
                  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"></path>
                    <polyline points="13 2 13 9 20 9"></polyline>
                  </svg>
                  <span class="file-name">{{ file.name }}</span>
                  <span class="file-size">{{ formatFileSize(file.size) }}</span>
                  <button type="button" @click.stop="removeFile(index)" class="btn-remove">
                    &times;
                  </button>
                </div>
                <button type="button" @click.stop="triggerFileInput" class="btn-add-more">
                  + Add More Files
                </button>
              </div>
            </div>
            <small class="help-text">Only .txt files are allowed. Maximum 10MB per file.</small>
          </div>

          <!-- Notes (Optional) -->
          <div class="form-group">
            <label>Notes (Optional)</label>
            <textarea
              v-model="formData.notes"
              placeholder="Add any additional notes..."
              class="form-textarea"
              rows="3"
            ></textarea>
          </div>

          <!-- Error Message -->
          <div v-if="error" class="error-message">
            {{ error }}
          </div>

          <!-- Form Actions -->
          <div class="form-actions">
            <button type="button" @click="$emit('close')" class="btn-cancel">
              Cancel
            </button>
            <button type="submit" :disabled="uploading || selectedFiles.length === 0" class="btn-submit">
              {{ uploading ? 'Uploading...' : 'Upload Files' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../../services/api.js';

const props = defineProps({
  schedule: {
    type: Object,
    default: null
  }
});

const emit = defineEmits(['close', 'uploaded']);

const formData = ref({
  orderNumber: '',
  partName: '',
  notes: ''
});

const selectedFiles = ref([]);
const fileInput = ref(null);
const uploading = ref(false);
const error = ref(null);
const lookingUpPartName = ref(false);

onMounted(() => {
  if (props.schedule) {
    formData.value.orderNumber = props.schedule.njo || '';
    formData.value.partName = props.schedule.part_name || '';
  }
});

function triggerFileInput() {
  fileInput.value?.click();
}

function handleFileSelect(event) {
  const files = Array.from(event.target.files || []);

  // Validate files
  for (const file of files) {
    // Check file extension
    if (!file.name.toLowerCase().endsWith('.txt')) {
      error.value = `File "${file.name}" is not a .txt file. Only .txt files are allowed.`;
      return;
    }

    // Check file size (10MB max)
    if (file.size > 10 * 1024 * 1024) {
      error.value = `File "${file.name}" exceeds 10MB limit.`;
      return;
    }
  }

  // Add files to selected list (avoid duplicates)
  files.forEach(file => {
    if (!selectedFiles.value.some(f => f.name === file.name && f.size === file.size)) {
      selectedFiles.value.push(file);
    }
  });

  error.value = null;

  // Reset input
  if (fileInput.value) {
    fileInput.value.value = '';
  }
}

function removeFile(index) {
  selectedFiles.value.splice(index, 1);
}

function formatFileSize(bytes) {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
}

async function lookupPartName() {
  if (!formData.value.orderNumber) return;

  lookingUpPartName.value = true;
  error.value = null;

  try {
    const response = await api.getPartNameByOrderNumber(formData.value.orderNumber);
    if (response.success && response.data) {
      formData.value.partName = response.data.part_name;
    }
  } catch (err) {
    console.warn('Could not fetch part name:', err.message);
    // Don't show error, part name is optional
  } finally {
    lookingUpPartName.value = false;
  }
}

async function handleUpload() {
  if (selectedFiles.value.length === 0) {
    error.value = 'Please select at least one file';
    return;
  }

  if (!formData.value.orderNumber) {
    error.value = 'Order number is required';
    return;
  }

  uploading.value = true;
  error.value = null;

  try {
    const uploadFormData = new FormData();

    // Append files
    selectedFiles.value.forEach(file => {
      uploadFormData.append('files', file);
    });

    // Append form fields
    uploadFormData.append('order_number', formData.value.orderNumber);

    if (formData.value.partName) {
      uploadFormData.append('part_name', formData.value.partName);
    }

    if (formData.value.notes) {
      uploadFormData.append('notes', formData.value.notes);
    }

    if (props.schedule?.id) {
      uploadFormData.append('ppic_schedule_id', props.schedule.id);
    }

    const response = await api.uploadToolpatherFiles(uploadFormData);

    if (response.success) {
      alert(`Successfully uploaded ${selectedFiles.value.length} file(s)!`);
      emit('uploaded');
    }
  } catch (err) {
    error.value = 'Upload failed: ' + err.message;
    console.error('Upload error:', err);
  } finally {
    uploading.value = false;
  }
}
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
  z-index: 1000;
}

.modal-container {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid #e5e7eb;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.btn-close {
  background: none;
  border: none;
  color: white;
  font-size: 2rem;
  cursor: pointer;
  line-height: 1;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: background 0.2s;
}

.btn-close:hover {
  background: rgba(255, 255, 255, 0.2);
}

.modal-body {
  padding: 2rem;
  overflow-y: auto;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #374151;
}

.required {
  color: #ef4444;
}

.form-input, .form-textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-input:focus, .form-textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-input:read-only {
  background-color: #f9fafb;
  cursor: not-allowed;
}

.file-upload-area {
  border: 2px dashed #d1d5db;
  border-radius: 8px;
  padding: 2rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  background: #f9fafb;
}

.file-upload-area:hover {
  border-color: #667eea;
  background: #f5f7ff;
}

.upload-placeholder svg {
  color: #9ca3af;
  margin-bottom: 1rem;
}

.upload-placeholder p {
  margin: 0 0 0.5rem;
  color: #374151;
  font-weight: 500;
}

.upload-placeholder small {
  color: #6b7280;
}

.selected-files {
  text-align: left;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  margin-bottom: 0.5rem;
}

.file-item svg {
  color: #667eea;
  flex-shrink: 0;
}

.file-name {
  flex: 1;
  color: #374151;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-size {
  color: #6b7280;
  font-size: 0.85rem;
  flex-shrink: 0;
}

.btn-remove {
  background: #fee2e2;
  color: #dc2626;
  border: none;
  width: 24px;
  height: 24px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1.2rem;
  line-height: 1;
  flex-shrink: 0;
  transition: all 0.2s;
}

.btn-remove:hover {
  background: #fecaca;
}

.btn-add-more {
  width: 100%;
  padding: 0.75rem;
  background: white;
  border: 1px solid #667eea;
  color: #667eea;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-add-more:hover {
  background: #f5f7ff;
}

.help-text {
  display: block;
  margin-top: 0.25rem;
  color: #6b7280;
  font-size: 0.875rem;
}

.error-message {
  padding: 1rem;
  background: #fee2e2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  color: #dc2626;
  margin-bottom: 1rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #e5e7eb;
}

.btn-cancel, .btn-submit {
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 500;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel {
  background: #f3f4f6;
  color: #374151;
  border: 1px solid #d1d5db;
}

.btn-cancel:hover {
  background: #e5e7eb;
}

.btn-submit {
  background: #667eea;
  color: white;
  border: none;
}

.btn-submit:hover:not(:disabled) {
  background: #5568d3;
}

.btn-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
