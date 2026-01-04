<template>
  <div class="toolpather-wrapper">
    <header class="main-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="logo">IMETRAX</h1>
          <span class="divider">|</span>
          <h2 class="page-title">Toolpather - File Management</h2>
        </div>

        <div class="header-right">
          <button @click="goToDashboard" class="btn-back">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            Dashboard
          </button>
        </div>
      </div>
    </header>

    <main class="toolpather-body">
      <div class="tabs">
        <button
          :class="['tab-btn', { active: activeTab === 'schedules' }]"
          @click="activeTab = 'schedules'"
        >
          PPIC Schedules
        </button>
        <button
          :class="['tab-btn', { active: activeTab === 'myFiles' }]"
          @click="activeTab = 'myFiles'"
        >
          My Files
        </button>
        <button
          :class="['tab-btn', { active: activeTab === 'allFiles' }]"
          @click="activeTab = 'allFiles'"
        >
          All Files
        </button>
      </div>

      <!-- PPIC Schedules Tab -->
      <div v-if="activeTab === 'schedules'" class="tab-content">
        <div class="section-header">
          <h3>PPIC Production Schedules</h3>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search by NJO or Part Name..."
            class="search-input"
          />
        </div>

        <div v-if="loading" class="loading">Loading schedules...</div>
        <div v-else-if="error" class="error">{{ error }}</div>
        <div v-else-if="filteredSchedules.length === 0" class="no-data">No schedules found</div>

        <div v-else class="table-container">
          <table class="schedules-table">
            <thead>
              <tr>
                <th>NJO</th>
                <th>Part Name</th>
                <th>Priority</th>
                <th>Material Status</th>
                <th>Start Date</th>
                <th>Finish Date</th>
                <th>Status</th>
                <th>Action</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="schedule in filteredSchedules" :key="schedule.id">
                <td>{{ schedule.njo }}</td>
                <td>{{ schedule.part_name }}</td>
                <td>
                  <span :class="['badge', 'priority-' + getPriorityClass(schedule.priority)]">
                    {{ schedule.priority }}
                  </span>
                </td>
                <td>
                  <span :class="['badge', 'material-' + getMaterialClass(schedule.material_status)]">
                    {{ schedule.material_status }}
                  </span>
                </td>
                <td>{{ formatDate(schedule.start_date) }}</td>
                <td>{{ formatDate(schedule.finish_date) }}</td>
                <td>
                  <span :class="['badge', 'status-' + schedule.status]">
                    {{ schedule.status }}
                  </span>
                </td>
                <td>
                  <button @click="uploadFiles(schedule)" class="btn-upload">
                    Upload Files
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- My Files Tab -->
      <div v-if="activeTab === 'myFiles'" class="tab-content">
        <FileList :showMyFilesOnly="true" @refresh="loadSchedules" />
      </div>

      <!-- All Files Tab -->
      <div v-if="activeTab === 'allFiles'" class="tab-content">
        <FileList :showMyFilesOnly="false" @refresh="loadSchedules" />
      </div>
    </main>

    <!-- File Upload Form Modal -->
    <FileUploadForm
      v-if="showUploadForm"
      :schedule="selectedSchedule"
      @close="closeUploadForm"
      @uploaded="handleFilesUploaded"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../services/api.js';
import FileUploadForm from './FileUploadForm.vue';
import FileList from './FileList.vue';

const router = useRouter();
const activeTab = ref('schedules');
const schedules = ref([]);
const loading = ref(false);
const error = ref(null);
const searchQuery = ref('');
const showUploadForm = ref(false);
const selectedSchedule = ref(null);

const filteredSchedules = computed(() => {
  if (!searchQuery.value) return schedules.value;

  const query = searchQuery.value.toLowerCase();
  return schedules.value.filter(schedule =>
    schedule.njo?.toLowerCase().includes(query) ||
    schedule.part_name?.toLowerCase().includes(query)
  );
});

async function loadSchedules() {
  loading.value = true;
  error.value = null;

  try {
    const response = await api.getGanttChart();

    if (response.success && response.data && response.data.sections) {
      // Extract all tasks from sections
      const allTasks = response.data.sections.flatMap(section => section.tasks || []);
      schedules.value = allTasks.map(task => ({
        id: task.task_id?.toString().replace('task-', ''),
        njo: task.njo,
        part_name: task.part_name || task.task_name,
        priority: task.priority,
        material_status: task.material_status,
        start_date: task.start,
        finish_date: task.end,
        status: task.status,
        ppic_notes: task.ppic_notes
      }));
    }
  } catch (err) {
    error.value = 'Failed to load PPIC schedules: ' + err.message;
    console.error('Error loading schedules:', err);
  } finally {
    loading.value = false;
  }
}

function uploadFiles(schedule) {
  selectedSchedule.value = schedule;
  showUploadForm.value = true;
}

function closeUploadForm() {
  showUploadForm.value = false;
  selectedSchedule.value = null;
}

function handleFilesUploaded() {
  closeUploadForm();
  activeTab.value = 'myFiles';
}

function formatDate(dateString) {
  if (!dateString) return '-';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
}

function getPriorityClass(priority) {
  const map = {
    'Top Urgent': 'top-urgent',
    'Urgent': 'urgent',
    'Medium': 'medium',
    'Low': 'low'
  };
  return map[priority] || 'low';
}

function getMaterialClass(status) {
  const map = {
    'Ready': 'ready',
    'Pending': 'pending',
    'Ordered': 'ordered',
    'Not Ready': 'not-ready'
  };
  return map[status] || 'pending';
}

function goToDashboard() {
  router.push('/dashboard');
}

onMounted(() => {
  loadSchedules();
});
</script>

<style scoped>
.toolpather-wrapper {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.main-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1.5rem 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1400px;
  margin: 0 auto;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.logo {
  font-size: 1.8rem;
  font-weight: 700;
  margin: 0;
}

.divider {
  color: rgba(255, 255, 255, 0.5);
  font-size: 1.5rem;
}

.page-title {
  font-size: 1.3rem;
  font-weight: 500;
  margin: 0;
}

.btn-back {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1.2rem;
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: white;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.95rem;
  transition: all 0.3s;
}

.btn-back:hover {
  background: rgba(255, 255, 255, 0.3);
}

.toolpather-body {
  max-width: 1400px;
  margin: 2rem auto;
  padding: 0 2rem;
}

.tabs {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  border-bottom: 2px solid #e5e7eb;
}

.tab-btn {
  padding: 1rem 2rem;
  background: none;
  border: none;
  color: #6b7280;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  border-bottom: 3px solid transparent;
  transition: all 0.3s;
}

.tab-btn.active {
  color: #667eea;
  border-bottom-color: #667eea;
}

.tab-btn:hover {
  color: #667eea;
}

.tab-content {
  background: white;
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-header h3 {
  margin: 0;
  color: #1f2937;
  font-size: 1.5rem;
}

.search-input {
  padding: 0.6rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  width: 300px;
  font-size: 0.95rem;
}

.table-container {
  overflow-x: auto;
}

.schedules-table {
  width: 100%;
  border-collapse: collapse;
}

.schedules-table th {
  background: #f9fafb;
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #374151;
  border-bottom: 2px solid #e5e7eb;
}

.schedules-table td {
  padding: 1rem;
  border-bottom: 1px solid #e5e7eb;
}

.schedules-table tbody tr:hover {
  background: #f9fafb;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 500;
}

.priority-top-urgent { background: #fee2e2; color: #991b1b; }
.priority-urgent { background: #fed7aa; color: #9a3412; }
.priority-medium { background: #fef3c7; color: #92400e; }
.priority-low { background: #d1fae5; color: #065f46; }

.material-ready { background: #d1fae5; color: #065f46; }
.material-pending { background: #fef3c7; color: #92400e; }
.material-ordered { background: #dbeafe; color: #1e40af; }
.material-not-ready { background: #fee2e2; color: #991b1b; }

.status-pending { background: #fef3c7; color: #92400e; }
.status-in_progress { background: #dbeafe; color: #1e40af; }
.status-completed { background: #d1fae5; color: #065f46; }

.btn-upload {
  padding: 0.5rem 1rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-upload:hover {
  background: #5568d3;
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
</style>
