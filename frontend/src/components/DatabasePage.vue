<template>
  <div class="database-wrapper">
    <header class="main-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="logo">IMETRAX</h1>
          <span class="divider">|</span>
          <h2 class="page-title">Database - Google Sheets Data</h2>
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

    <main class="database-body">
      <div class="database-container">
        <!-- Search and Filter Section -->
        <div class="controls-section">
          <div class="search-box">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"></circle>
              <path d="M21 21l-4.35-4.35"></path>
            </svg>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search in all columns..."
              class="search-input"
            />
          </div>
          <button @click="loadData" class="btn-refresh" :disabled="loading">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"/>
            </svg>
            Refresh Data
          </button>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>Loading Google Sheets data...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="error-state">
          <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="8" x2="12" y2="12"></line>
            <line x1="12" y1="16" x2="12.01" y2="16"></line>
          </svg>
          <h3>Error Loading Data</h3>
          <p>{{ error }}</p>
          <button @click="loadData" class="btn-retry">Try Again</button>
        </div>

        <!-- Data Table -->
        <div v-else-if="tableData.rows.length > 0" class="table-section">
          <div class="table-info">
            <p>
              Showing {{ filteredRows.length }} of {{ tableData.rows.length }} records
            </p>
          </div>

          <div class="table-wrapper">
            <table class="data-table">
              <thead>
                <tr>
                  <th style="width: 60px; text-align: center;">Detail</th>
                  <th v-for="(header, index) in tableData.headers" :key="index">
                    {{ header }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(row, rowIndex) in paginatedRows" :key="rowIndex">
                  <td style="text-align: center;">
                    <button
                      @click="openDetailSidebar(row, rowIndex)"
                      class="btn-detail"
                      title="View Details"
                    >
                      <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18"
                           viewBox="0 0 24 24" fill="none" stroke="currentColor"
                           stroke-width="2">
                        <circle cx="11" cy="11" r="8"></circle>
                        <path d="M21 21l-4.35-4.35"></path>
                      </svg>
                    </button>
                  </td>
                  <td v-for="(cell, cellIndex) in row" :key="cellIndex">
                    {{ cell || '-' }}
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Pagination -->
          <div class="pagination" v-if="totalPages > 1">
            <button
              @click="currentPage = 1"
              :disabled="currentPage === 1"
              class="btn-page"
            >
              First
            </button>
            <button
              @click="currentPage--"
              :disabled="currentPage === 1"
              class="btn-page"
            >
              Previous
            </button>

            <span class="page-info">
              Page {{ currentPage }} of {{ totalPages }}
            </span>

            <button
              @click="currentPage++"
              :disabled="currentPage === totalPages"
              class="btn-page"
            >
              Next
            </button>
            <button
              @click="currentPage = totalPages"
              :disabled="currentPage === totalPages"
              class="btn-page"
            >
              Last
            </button>

            <select v-model.number="rowsPerPage" class="rows-per-page">
              <option :value="25">25 per page</option>
              <option :value="50">50 per page</option>
              <option :value="100">100 per page</option>
              <option :value="200">200 per page</option>
            </select>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="empty-state">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
            <polyline points="14 2 14 8 20 8"></polyline>
            <line x1="9" y1="15" x2="15" y2="15"></line>
          </svg>
          <h3>No Data Available</h3>
          <p>The Google Sheets data is empty or could not be loaded.</p>
        </div>
      </div>
    </main>

    <!-- Detail Sidebar Overlay -->
    <Transition name="sidebar">
      <div v-if="sidebarOpen" class="sidebar-overlay" @click="closeSidebar">
        <div class="sidebar-panel" @click.stop>
          <div class="sidebar-header">
            <h3>Order Details</h3>
            <button @click="closeSidebar" class="btn-close">
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24"
                   viewBox="0 0 24 24" fill="none" stroke="currentColor"
                   stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <div class="sidebar-body">
            <!-- Operation Plans Section -->
            <div class="op-section">
              <div class="op-header">
                <h4>PEM Operation Plans</h4>
              </div>

              <!-- Loading state -->
              <div v-if="opLoading" class="op-loading">
                <div class="spinner-small"></div>
                <span>Checking operation plans...</span>
              </div>

              <!-- No PPIC Schedule found -->
              <div v-else-if="!opSchedule" class="op-empty">
                <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32"
                     viewBox="0 0 24 24" fill="none" stroke="currentColor"
                     stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="12" y1="8" x2="12" y2="12"></line>
                  <line x1="12" y1="16" x2="12.01" y2="16"></line>
                </svg>
                <p>No PPIC Schedule found for this order</p>
                <small>NJO: {{ getNoOrderFromRow() }}</small>
              </div>

              <!-- Operation Plans found -->
              <div v-else class="op-content">
                <!-- Status Info -->
                <div class="op-status">
                  <div class="op-status-item">
                    <span class="op-status-label">PPIC Schedule:</span>
                    <span class="op-status-value">{{ opSchedule.njo }}</span>
                  </div>
                  <div class="op-status-item">
                    <span class="op-status-label">Part Name:</span>
                    <span class="op-status-value">{{ opSchedule.part_name }}</span>
                  </div>
                  <div class="op-status-item">
                    <span class="op-status-label">Operation Plans:</span>
                    <span class="op-status-badge" :class="opPlans.length > 0 ? 'badge-success' : 'badge-warning'">
                      {{ opPlans.length }} Plan(s)
                    </span>
                  </div>
                </div>

                <!-- Plans List -->
                <div v-if="opPlans.length > 0" class="op-list">
                  <div
                    v-for="plan in opPlans"
                    :key="plan.id"
                    class="op-card"
                  >
                    <div class="op-card-header">
                      <span class="op-form-number">{{ plan.form_number }}</span>
                      <span
                        class="op-status-badge"
                        :class="{
                          'badge-draft': plan.status === 'draft',
                          'badge-pending': plan.status === 'pending_approval',
                          'badge-success': plan.status === 'approved',
                          'badge-danger': plan.status === 'rejected'
                        }"
                      >
                        {{ plan.status }}
                      </span>
                    </div>
                    <div class="op-card-body">
                      <div class="op-info-row">
                        <span class="op-info-label">Material:</span>
                        <span class="op-info-value">{{ plan.material || '-' }}</span>
                      </div>
                      <div class="op-info-row">
                        <span class="op-info-label">Steps:</span>
                        <span class="op-info-value">{{ plan.steps?.length || 0 }} step(s)</span>
                      </div>
                    </div>
                    <button
                      @click="viewOperationPlanDetail(plan)"
                      class="btn-op-detail"
                    >
                      View Full Detail
                    </button>
                  </div>
                </div>

                <!-- Create Button -->
                <div class="op-actions">
                  <button
                    @click="openCreateModal"
                    class="btn-create-op"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18"
                         viewBox="0 0 24 24" fill="none" stroke="currentColor"
                         stroke-width="2">
                      <line x1="12" y1="5" x2="12" y2="19"></line>
                      <line x1="5" y1="12" x2="19" y2="12"></line>
                    </svg>
                    Create New Operation Plan
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Operation Plan Form Modal -->
    <OperationPlanForm
      v-if="showCreateModal"
      :schedule="opSchedule"
      @close="closeCreateModal"
      @saved="handlePlanCreated"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api.js';
import OperationPlanForm from './PEM/OperationPlanForm.vue';

const router = useRouter();
const loading = ref(false);
const error = ref(null);
const searchQuery = ref('');
const currentPage = ref(1);
const rowsPerPage = ref(50);

const tableData = ref({
  headers: [],
  rows: [],
  total: 0
});

// Sidebar state
const sidebarOpen = ref(false);
const selectedRow = ref(null);
const selectedRowIndex = ref(null);

// Operation Plan state
const opLoading = ref(false);
const opSchedule = ref(null);
const opPlans = ref([]);
const showCreateModal = ref(false);

// Filtered rows based on search
const filteredRows = computed(() => {
  if (!searchQuery.value) {
    return tableData.value.rows;
  }

  const query = searchQuery.value.toLowerCase();
  return tableData.value.rows.filter(row => {
    return row.some(cell => {
      return String(cell).toLowerCase().includes(query);
    });
  });
});

// Paginated rows
const paginatedRows = computed(() => {
  const start = (currentPage.value - 1) * rowsPerPage.value;
  const end = start + rowsPerPage.value;
  return filteredRows.value.slice(start, end);
});

// Total pages
const totalPages = computed(() => {
  return Math.ceil(filteredRows.value.length / rowsPerPage.value);
});

// Reset to first page when search query changes
watch(searchQuery, () => {
  currentPage.value = 1;
});

// Reset to first page when rows per page changes
watch(rowsPerPage, () => {
  currentPage.value = 1;
});

async function loadData() {
  loading.value = true;
  error.value = null;

  try {
    const response = await api.getAllGoogleSheetsData();

    if (response.success && response.data) {
      tableData.value = {
        headers: response.data.headers || [],
        rows: response.data.rows || [],
        total: response.data.total || 0
      };
    } else {
      throw new Error('Invalid response format');
    }
  } catch (err) {
    error.value = err.message || 'Failed to load Google Sheets data';
    console.error('Error loading Google Sheets data:', err);
  } finally {
    loading.value = false;
  }
}

function goToDashboard() {
  router.push('/dashboard');
}

async function openDetailSidebar(row, rowIndex) {
  selectedRow.value = row;
  selectedRowIndex.value = rowIndex;
  sidebarOpen.value = true;

  // Load operation plans for this row
  const noOrder = getNoOrderFromRow();
  if (noOrder) {
    await loadOperationPlans(noOrder);
  }
}

function closeSidebar() {
  sidebarOpen.value = false;
  selectedRow.value = null;
  selectedRowIndex.value = null;
  opSchedule.value = null;
  opPlans.value = [];
}

async function loadOperationPlans(noOrder) {
  if (!noOrder) return;

  opLoading.value = true;
  try {
    // 1. Find PPIC Schedule by NJO
    const schedulesResponse = await api.getAllPPICSchedules();
    const schedules = schedulesResponse.data || [];
    const schedule = schedules.find(s => s.njo === noOrder);

    if (schedule) {
      opSchedule.value = schedule;

      // 2. Get PEM Operation Plans for this schedule
      const plansResponse = await api.getPEMPlansByPPICSchedule(schedule.id);
      opPlans.value = plansResponse.data || [];
    } else {
      opSchedule.value = null;
      opPlans.value = [];
    }
  } catch (err) {
    console.error('Error loading operation plans:', err);
    opSchedule.value = null;
    opPlans.value = [];
  } finally {
    opLoading.value = false;
  }
}

function viewOperationPlanDetail(plan) {
  // Navigate to PEM page with plan detail
  router.push(`/pem?plan_id=${plan.id}`);
}

function openCreateModal() {
  showCreateModal.value = true;
}

function closeCreateModal() {
  showCreateModal.value = false;
}

function handlePlanCreated() {
  showCreateModal.value = false;
  // Reload operation plans
  const noOrder = getNoOrderFromRow();
  if (noOrder) {
    loadOperationPlans(noOrder);
  }
}

function getNoOrderFromRow() {
  if (!selectedRow.value || !tableData.value.headers.length) return null;
  const noOrderIndex = tableData.value.headers.findIndex(h =>
    h.toLowerCase().includes('no order') || h.toLowerCase().includes('no. order')
  );
  return noOrderIndex >= 0 ? selectedRow.value[noOrderIndex] : null;
}

function getPartNameFromRow() {
  if (!selectedRow.value || !tableData.value.headers.length) return null;
  const partNameIndex = tableData.value.headers.findIndex(h =>
    h.toLowerCase().includes('part name') || h.toLowerCase().includes('part')
  );
  return partNameIndex >= 0 ? selectedRow.value[partNameIndex] : null;
}

onMounted(() => {
  loadData();
});
</script>

<style scoped>
.database-wrapper {
  min-height: 100vh;
  background: #f3f4f6;
  display: flex;
  flex-direction: column;
}

.main-header {
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  padding: 1rem 2rem;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1800px;
  margin: 0 auto;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.logo {
  color: rgb(0, 0, 0);
  font-size: 1.75rem;
  font-weight: 700;
  margin: 0;
}

.divider {
  color: rgba(255, 255, 255, 0.5);
  font-size: 1.5rem;
}

.page-title {
  color: rgb(0, 0, 0);
  font-size: 1.25rem;
  font-weight: 500;
  margin: 0;
}

.btn-back {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1.2rem;
  background: rgba(38, 32, 32, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 8px;
  color: rgb(0, 0, 0);
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-back:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.database-body {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
}

.database-container {
  max-width: 1800px;
  margin: 0 auto;
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  padding: 2rem;
}

.controls-section {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
}

.search-box {
  flex: 1;
  min-width: 300px;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}

.search-box svg {
  color: #6b7280;
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  border: none;
  background: none;
  font-size: 0.95rem;
  outline: none;
}

.btn-refresh {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-refresh:hover:not(:disabled) {
  background: #5568d3;
  transform: translateY(-2px);
}

.btn-refresh:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-state,
.error-state,
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #6b7280;
}

.loading-state .spinner {
  width: 48px;
  height: 48px;
  margin: 0 auto 1rem;
  border: 4px solid #e5e7eb;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error-state svg,
.empty-state svg {
  color: #9ca3af;
  margin-bottom: 1rem;
}

.error-state h3,
.empty-state h3 {
  color: #1f2937;
  margin-bottom: 0.5rem;
}

.btn-retry {
  margin-top: 1rem;
  padding: 0.75rem 1.5rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-retry:hover {
  background: #5568d3;
}

.table-section {
  margin-top: 1rem;
}

.table-info {
  margin-bottom: 1rem;
  color: #6b7280;
  font-size: 0.9rem;
}

.table-wrapper {
  overflow-x: auto;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.data-table thead {
  background: #f9fafb;
  position: sticky;
  top: 0;
  z-index: 10;
}

.data-table th {
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #374151;
  border-bottom: 2px solid #e5e7eb;
  white-space: nowrap;
}

.data-table td {
  padding: 0.875rem 1rem;
  border-bottom: 1px solid #f3f4f6;
  color: #6b7280;
}

.data-table tbody tr:hover {
  background: #f9fafb;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  margin-top: 2rem;
  flex-wrap: wrap;
}

.btn-page {
  padding: 0.5rem 1rem;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  color: #374151;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-page:hover:not(:disabled) {
  background: #f9fafb;
  border-color: #667eea;
}

.btn-page:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.page-info {
  padding: 0.5rem 1rem;
  color: #374151;
  font-size: 0.875rem;
  font-weight: 500;
}

.rows-per-page {
  padding: 0.5rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 0.875rem;
  background: white;
  cursor: pointer;
}

.rows-per-page:focus {
  outline: none;
  border-color: #667eea;
}

/* Detail Button */
.btn-detail {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-detail:hover {
  background: #667eea;
  color: white;
  transform: scale(1.1);
}

.btn-detail svg {
  display: block;
}

/* Sidebar Overlay */
.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: flex-end;
  align-items: stretch;
  z-index: 1000;
  backdrop-filter: blur(2px);
}

.sidebar-panel {
  width: 500px;
  max-width: 90vw;
  background: white;
  box-shadow: -4px 0 20px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid #e5e7eb;
  background: #f9fafb;
}

.sidebar-header h3 {
  margin: 0;
  color: #1f2937;
  font-size: 1.25rem;
  font-weight: 600;
}

.btn-close {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.5rem;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-close:hover {
  background: #e5e7eb;
  color: #1f2937;
}

.sidebar-body {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem 2rem;
}

.detail-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #f3f4f6;
}

.detail-item:last-child {
  border-bottom: none;
}

.detail-label {
  font-size: 0.75rem;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.detail-value {
  font-size: 0.95rem;
  color: #1f2937;
  word-wrap: break-word;
}

/* Sidebar Transition */
.sidebar-enter-active,
.sidebar-leave-active {
  transition: all 0.3s ease;
}

.sidebar-enter-from,
.sidebar-leave-to {
  opacity: 0;
}

.sidebar-enter-from .sidebar-panel,
.sidebar-leave-to .sidebar-panel {
  transform: translateX(100%);
}

.sidebar-enter-active .sidebar-panel,
.sidebar-leave-active .sidebar-panel {
  transition: transform 0.3s ease;
}

/* Responsive */
@media (max-width: 768px) {
  .sidebar-panel {
    width: 100%;
    max-width: 100vw;
  }
}

/* Operation Plans Section */
.op-section {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 2px solid #e5e7eb;
}

.op-header {
  margin-bottom: 1rem;
}

.op-header h4 {
  margin: 0;
  color: #1f2937;
  font-size: 1rem;
  font-weight: 600;
}

.op-loading {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1.5rem;
  background: #f9fafb;
  border-radius: 8px;
  color: #6b7280;
  font-size: 0.875rem;
}

.spinner-small {
  width: 20px;
  height: 20px;
  border: 2px solid #e5e7eb;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.op-empty {
  text-align: center;
  padding: 2rem 1rem;
  background: #fef3c7;
  border: 1px solid #fde68a;
  border-radius: 8px;
}

.op-empty svg {
  color: #f59e0b;
  margin-bottom: 0.75rem;
}

.op-empty p {
  margin: 0.5rem 0;
  color: #92400e;
  font-size: 0.9rem;
  font-weight: 500;
}

.op-empty small {
  color: #b45309;
  font-size: 0.75rem;
}

.op-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.op-status {
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 8px;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.op-status-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
}

.op-status-label {
  color: #0369a1;
  font-weight: 500;
}

.op-status-value {
  color: #075985;
  font-weight: 600;
}

.op-status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.badge-success {
  background: #d1fae5;
  color: #065f46;
}

.badge-warning {
  background: #fef3c7;
  color: #92400e;
}

.badge-draft {
  background: #e5e7eb;
  color: #374151;
}

.badge-pending {
  background: #dbeafe;
  color: #1e40af;
}

.badge-danger {
  background: #fee2e2;
  color: #991b1b;
}

.op-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.op-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1rem;
  transition: all 0.2s;
}

.op-card:hover {
  border-color: #667eea;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
}

.op-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #f3f4f6;
}

.op-form-number {
  font-size: 0.875rem;
  font-weight: 600;
  color: #667eea;
}

.op-card-body {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.op-info-row {
  display: flex;
  justify-content: space-between;
  font-size: 0.8rem;
}

.op-info-label {
  color: #6b7280;
}

.op-info-value {
  color: #1f2937;
  font-weight: 500;
}

.btn-op-detail {
  width: 100%;
  padding: 0.5rem;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  color: #667eea;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-op-detail:hover {
  background: #667eea;
  color: white;
  border-color: #667eea;
}

.op-actions {
  margin-top: 0.5rem;
}

.btn-create-op {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.75rem;
  background: #667eea;
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-create-op:hover {
  background: #5568d3;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-create-op svg {
  flex-shrink: 0;
}
</style>
