<template>
  <div class="pem-wrapper">
    <header class="main-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="logo">IMETRAX</h1>
          <span class="divider">|</span>
          <h2 class="page-title">PEM - Operation Plans</h2>
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

    <main class="pem-body">
      <div class="tabs">
        <button
          :class="['tab-btn', { active: activeTab === 'schedules' }]"
          @click="activeTab = 'schedules'"
        >
          PPIC Schedules
        </button>
        <button
          :class="['tab-btn', { active: activeTab === 'plans' }]"
          @click="activeTab = 'plans'"
        >
          My Operation Plans
        </button>
        <button
          :class="['tab-btn', { active: activeTab === 'approvals' }]"
          @click="activeTab = 'approvals'"
        >
          Pending Approvals
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
                  <button @click="createOperationPlan(schedule)" class="btn-create">
                    Create OP
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- My Operation Plans Tab -->
      <div v-if="activeTab === 'plans'" class="tab-content">
        <OperationPlanList @edit="handleEditPlan" />
      </div>

      <!-- Pending Approvals Tab -->
      <div v-if="activeTab === 'approvals'" class="tab-content">
        <h3>Plans Pending Your Approval</h3>
        <div v-if="loadingApprovals" class="loading">Loading approvals...</div>
        <div v-else-if="pendingApprovals.length === 0" class="no-data">No pending approvals</div>

        <div v-else class="approvals-grid">
          <div v-for="plan in pendingApprovals" :key="plan.id" class="approval-card">
            <div class="card-header">
              <h4>{{ plan.form_number }}</h4>
              <span class="badge status-pending_approval">Pending</span>
            </div>
            <div class="card-body">
              <p><strong>Part Name:</strong> {{ plan.part_name }}</p>
              <p><strong>Material:</strong> {{ plan.material }}</p>
              <p><strong>Created By:</strong> {{ plan.creator?.username }}</p>
              <p><strong>Created:</strong> {{ formatDate(plan.created_at) }}</p>
            </div>
            <div class="card-actions">
              <button @click="viewPlan(plan.id)" class="btn-view">View Details</button>
              <button @click="approvePlan(plan.id)" class="btn-approve">Approve</button>
              <button @click="rejectPlan(plan.id)" class="btn-reject">Reject</button>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Operation Plan Form Modal -->
    <OperationPlanForm
      v-if="showForm"
      :schedule="selectedSchedule"
      :editPlan="editingPlan"
      @close="closeForm"
      @saved="handlePlanSaved"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../services/api.js';
import OperationPlanForm from './OperationPlanForm.vue';
import OperationPlanList from './OperationPlanList.vue';

const router = useRouter();
const activeTab = ref('schedules');
const schedules = ref([]);
const loading = ref(false);
const error = ref(null);
const searchQuery = ref('');
const showForm = ref(false);
const selectedSchedule = ref(null);
const editingPlan = ref(null);
const pendingApprovals = ref([]);
const loadingApprovals = ref(false);

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

async function loadPendingApprovals() {
  loadingApprovals.value = true;

  try {
    const response = await api.getPendingPEMApprovals();
    if (response.success && response.data) {
      pendingApprovals.value = response.data;
    }
  } catch (err) {
    console.error('Error loading pending approvals:', err);
  } finally {
    loadingApprovals.value = false;
  }
}

function createOperationPlan(schedule) {
  selectedSchedule.value = schedule;
  editingPlan.value = null;
  showForm.value = true;
}

function handleEditPlan(plan) {
  editingPlan.value = plan;
  selectedSchedule.value = null;
  showForm.value = true;
}

function closeForm() {
  showForm.value = false;
  selectedSchedule.value = null;
  editingPlan.value = null;
}

function handlePlanSaved() {
  closeForm();
  activeTab.value = 'plans';
}

async function viewPlan(planId) {
  router.push(`/pem/plan/${planId}`);
}

function getUserRoleForPlan(plan) {
  // Get current user
  const currentUser = api.getCurrentUser();
  if (!currentUser) return null;

  // Find which approval role(s) this user has for this plan
  const userApprovals = plan.approvals?.filter(
    approval => approval.approver_id === currentUser.id && approval.status === 'pending'
  );

  if (!userApprovals || userApprovals.length === 0) return null;

  // If user has multiple roles, ask them to choose
  if (userApprovals.length > 1) {
    const roleOptions = userApprovals.map(a => a.approver_role).join(', ');
    const selectedRole = prompt(`You have multiple approval roles for this plan (${roleOptions}). Enter the role you want to approve as:`);
    return selectedRole;
  }

  return userApprovals[0].approver_role;
}

async function approvePlan(planId) {
  const plan = pendingApprovals.value.find(p => p.id === planId);
  if (!plan) return;

  const role = getUserRoleForPlan(plan);
  if (!role) {
    alert('Unable to determine your approval role for this plan');
    return;
  }

  const comments = prompt(`Approving as ${role}. Enter approval comments (optional):`);
  if (comments === null) return; // User cancelled

  try {
    await api.approvePEMOperationPlan(planId, role, comments);
    alert('Plan approved successfully!');
    loadPendingApprovals();
  } catch (err) {
    alert('Failed to approve plan: ' + err.message);
  }
}

async function rejectPlan(planId) {
  const plan = pendingApprovals.value.find(p => p.id === planId);
  if (!plan) return;

  const role = getUserRoleForPlan(plan);
  if (!role) {
    alert('Unable to determine your approval role for this plan');
    return;
  }

  const comments = prompt(`Rejecting as ${role}. Enter rejection reason:`);
  if (!comments) {
    alert('Rejection reason is required');
    return;
  }

  try {
    await api.rejectPEMOperationPlan(planId, role, comments);
    alert('Plan rejected');
    loadPendingApprovals();
  } catch (err) {
    alert('Failed to reject plan: ' + err.message);
  }
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
  loadPendingApprovals();
});
</script>

<style scoped>
.pem-wrapper {
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

.pem-body {
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
.status-pending_approval { background: #fed7aa; color: #9a3412; }

.btn-create {
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

.btn-create:hover {
  background: #5568d3;
}

.approvals-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.approval-card {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1.5rem;
  background: white;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.card-header h4 {
  margin: 0;
  color: #1f2937;
}

.card-body p {
  margin: 0.5rem 0;
  color: #6b7280;
}

.card-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
}

.btn-view, .btn-approve, .btn-reject {
  flex: 1;
  padding: 0.6rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-view {
  background: #e5e7eb;
  color: #374151;
}

.btn-view:hover {
  background: #d1d5db;
}

.btn-approve {
  background: #10b981;
  color: white;
}

.btn-approve:hover {
  background: #059669;
}

.btn-reject {
  background: #ef4444;
  color: white;
}

.btn-reject:hover {
  background: #dc2626;
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
