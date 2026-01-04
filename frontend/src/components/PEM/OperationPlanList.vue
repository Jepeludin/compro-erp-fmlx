<template>
  <div class="plan-list">
    <div class="list-header">
      <h3>My Operation Plans</h3>
      <div class="filters">
        <select v-model="statusFilter" class="filter-select">
          <option value="">All Statuses</option>
          <option value="draft">Draft</option>
          <option value="pending_approval">Pending Approval</option>
          <option value="approved">Approved</option>
          <option value="rejected">Rejected</option>
        </select>
      </div>
    </div>

    <div v-if="loading" class="loading">Loading plans...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="filteredPlans.length === 0" class="no-data">
      No operation plans found. Create one from the PPIC Schedules tab.
    </div>

    <div v-else class="plans-grid">
      <div v-for="plan in filteredPlans" :key="plan.id" class="plan-card">
        <div class="card-header">
          <div>
            <h4>{{ plan.form_number }}</h4>
            <p class="part-name">{{ plan.part_name }}</p>
          </div>
          <span :class="['badge', 'status-' + plan.status]">
            {{ formatStatus(plan.status) }}
          </span>
        </div>

        <div class="card-body">
          <div class="info-row">
            <span class="label">Material:</span>
            <span>{{ plan.material || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="label">Dial Size:</span>
            <span>{{ plan.dial_size || '-' }}</span>
          </div>
          <div class="info-row">
            <span class="label">Quantity:</span>
            <span>{{ plan.quantity || 0 }}</span>
          </div>
          <div class="info-row">
            <span class="label">Steps:</span>
            <span>{{ plan.steps?.length || 0 }} step(s)</span>
          </div>
          <div class="info-row">
            <span class="label">Created:</span>
            <span>{{ formatDate(plan.created_at) }}</span>
          </div>
          <div v-if="plan.approval && plan.approval.approver" class="info-row">
            <span class="label">Approver:</span>
            <span>{{ plan.approval.approver.username }}</span>
          </div>
        </div>

        <div class="card-actions">
          <button @click="viewPlan(plan)" class="btn-view">View</button>
          <button
            v-if="plan.status === 'draft'"
            @click="editPlan(plan)"
            class="btn-edit"
          >
            Edit
          </button>
          <button
            v-if="plan.status === 'draft'"
            @click="deletePlan(plan)"
            class="btn-delete"
          >
            Delete
          </button>
          <button
            v-if="plan.status === 'rejected'"
            @click="resubmitPlan(plan)"
            class="btn-resubmit"
          >
            Resubmit
          </button>
        </div>
      </div>
    </div>

    <!-- View Plan Modal -->
    <div v-if="viewingPlan" class="modal-overlay" @click.self="viewingPlan = null">
      <div class="view-modal">
        <div class="modal-header">
          <h2>{{ viewingPlan.form_number }}</h2>
          <button @click="viewingPlan = null" class="btn-close">&times;</button>
        </div>
        <div class="modal-body">
          <div class="view-section">
            <h3>Part Information</h3>
            <div class="info-grid">
              <div><strong>Part Name:</strong> {{ viewingPlan.part_name }}</div>
              <div><strong>Material:</strong> {{ viewingPlan.material || '-' }}</div>
              <div><strong>Dial Size:</strong> {{ viewingPlan.dial_size || '-' }}</div>
              <div><strong>Quantity:</strong> {{ viewingPlan.quantity || 0 }}</div>
              <div><strong>Revision:</strong> {{ viewingPlan.revision || '-' }}</div>
              <div><strong>No. WP:</strong> {{ viewingPlan.no_wp || '-' }}</div>
              <div><strong>Page:</strong> {{ viewingPlan.page || '-' }}</div>
              <div><strong>Status:</strong> <span :class="['badge', 'status-' + viewingPlan.status]">{{ formatStatus(viewingPlan.status) }}</span></div>
            </div>
          </div>

          <div class="view-section">
            <h3>Process Steps</h3>
            <div v-if="!viewingPlan.steps || viewingPlan.steps.length === 0" class="no-data">
              No steps defined
            </div>
            <div v-else class="steps-list">
              <div v-for="step in viewingPlan.steps" :key="step.id" class="step-view">
                <h4>Step {{ step.step_number }}</h4>
                <div v-if="step.picture_url" class="step-image">
                  <img :src="getImageUrl(step.picture_url)" :alt="'Step ' + step.step_number" />
                </div>
                <div class="step-details">
                  <div v-if="step.clamping_system">
                    <strong>Clamping System:</strong> {{ step.clamping_system }}
                  </div>
                  <div v-if="step.raw_material">
                    <strong>Raw Material:</strong> {{ step.raw_material }}
                  </div>
                  <div v-if="step.setting">
                    <strong>Setting:</strong> {{ step.setting }}
                  </div>
                  <div v-if="step.process">
                    <strong>Process:</strong> {{ step.process }}
                  </div>
                  <div v-if="step.note">
                    <strong>Note:</strong> {{ step.note }}
                  </div>
                  <div v-if="step.checking_method">
                    <strong>Checking Method:</strong> {{ step.checking_method }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="viewingPlan.approval" class="view-section">
            <h3>Approval Information</h3>
            <div class="info-grid">
              <div v-if="viewingPlan.approval.approver">
                <strong>Approver:</strong> {{ viewingPlan.approval.approver.username }}
              </div>
              <div>
                <strong>Status:</strong> <span :class="['badge', 'status-' + viewingPlan.approval.status]">{{ formatStatus(viewingPlan.approval.status) }}</span>
              </div>
              <div v-if="viewingPlan.approval.approved_at">
                <strong>Approved/Rejected At:</strong> {{ formatDate(viewingPlan.approval.approved_at) }}
              </div>
              <div v-if="viewingPlan.approval.comments" class="full-width">
                <strong>Comments:</strong> {{ viewingPlan.approval.comments }}
              </div>
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

const emit = defineEmits(['edit']);

const plans = ref([]);
const loading = ref(false);
const error = ref(null);
const statusFilter = ref('');
const viewingPlan = ref(null);

const filteredPlans = computed(() => {
  if (!statusFilter.value) return plans.value;
  return plans.value.filter(plan => plan.status === statusFilter.value);
});

async function loadPlans() {
  loading.value = true;
  error.value = null;

  try {
    const response = await api.getAllPEMOperationPlans();
    if (response.success && response.data) {
      // Get current user from API service
      const currentUser = api.getCurrentUser();
      // Filter to show only plans created by current user
      plans.value = response.data.filter(plan => plan.created_by === currentUser.id);
    }
  } catch (err) {
    error.value = 'Failed to load operation plans: ' + err.message;
    console.error('Error loading plans:', err);
  } finally {
    loading.value = false;
  }
}

async function viewPlan(plan) {
  try {
    const response = await api.getPEMOperationPlan(plan.id);
    if (response.success && response.data) {
      viewingPlan.value = response.data;
    }
  } catch (err) {
    alert('Failed to load plan details: ' + err.message);
  }
}

async function editPlan(plan) {
  // Load full plan details before editing
  try {
    const response = await api.getPEMOperationPlan(plan.id);
    if (response.success && response.data) {
      emit('edit', response.data);
    }
  } catch (err) {
    alert('Failed to load plan for editing: ' + err.message);
  }
}

async function deletePlan(plan) {
  if (!confirm(`Are you sure you want to delete ${plan.form_number}?`)) {
    return;
  }

  try {
    await api.deletePEMOperationPlan(plan.id);
    alert('Plan deleted successfully');
    loadPlans();
  } catch (err) {
    alert('Failed to delete plan: ' + err.message);
  }
}

function resubmitPlan(plan) {
  // TODO: Implement resubmit functionality
  alert('Resubmit functionality coming soon! You can edit the plan and submit again.');
}

function formatDate(dateString) {
  if (!dateString) return '-';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' });
}

function formatStatus(status) {
  const statusMap = {
    'draft': 'Draft',
    'pending_approval': 'Pending Approval',
    'approved': 'Approved',
    'rejected': 'Rejected'
  };
  return statusMap[status] || status;
}

function getImageUrl(relPath) {
  // Remove /api/v1 from baseURL to get the root server URL
  const serverUrl = api.baseURL.replace('/api/v1', '');
  return `${serverUrl}/uploads/operation-plan-images/${relPath}`;
}

onMounted(() => {
  loadPlans();
});
</script>

<style scoped>
.plan-list {
  padding: 0;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.list-header h3 {
  margin: 0;
  color: #1f2937;
  font-size: 1.5rem;
}

.filter-select {
  padding: 0.6rem 1rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 0.95rem;
  background: white;
  cursor: pointer;
}

.plans-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.plan-card {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
  background: white;
  transition: all 0.3s;
}

.plan-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  background: #f9fafb;
  padding: 1rem 1.5rem;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  border-bottom: 1px solid #e5e7eb;
}

.card-header h4 {
  margin: 0 0 0.25rem 0;
  color: #1f2937;
  font-size: 1.1rem;
}

.part-name {
  margin: 0;
  color: #6b7280;
  font-size: 0.9rem;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 500;
}

.status-draft { background: #e5e7eb; color: #374151; }
.status-pending_approval { background: #fed7aa; color: #9a3412; }
.status-approved { background: #d1fae5; color: #065f46; }
.status-rejected { background: #fee2e2; color: #991b1b; }
.status-pending { background: #fed7aa; color: #9a3412; }

.card-body {
  padding: 1.5rem;
}

.info-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.75rem;
  color: #6b7280;
  font-size: 0.9rem;
}

.info-row .label {
  font-weight: 500;
  color: #374151;
}

.card-actions {
  padding: 1rem 1.5rem;
  background: #f9fafb;
  border-top: 1px solid #e5e7eb;
  display: flex;
  gap: 0.5rem;
}

.btn-view, .btn-edit, .btn-delete, .btn-resubmit {
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
  background: #667eea;
  color: white;
}

.btn-view:hover {
  background: #5568d3;
}

.btn-edit {
  background: #3b82f6;
  color: white;
}

.btn-edit:hover {
  background: #2563eb;
}

.btn-delete {
  background: #ef4444;
  color: white;
}

.btn-delete:hover {
  background: #dc2626;
}

.btn-resubmit {
  background: #f59e0b;
  color: white;
}

.btn-resubmit:hover {
  background: #d97706;
}

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
  padding: 2rem;
}

.view-modal {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 900px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-header {
  position: sticky;
  top: 0;
  background: white;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
  z-index: 1;
}

.modal-header h2 {
  margin: 0;
  color: #1f2937;
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
  padding: 2rem;
}

.view-section {
  margin-bottom: 2rem;
}

.view-section h3 {
  margin: 0 0 1rem 0;
  color: #1f2937;
  border-bottom: 2px solid #e5e7eb;
  padding-bottom: 0.5rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.info-grid .full-width {
  grid-column: 1 / -1;
}

.steps-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.step-view {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1.5rem;
  background: #f9fafb;
}

.step-view h4 {
  margin: 0 0 1rem 0;
  color: #1f2937;
}

.step-image {
  margin-bottom: 1rem;
}

.step-image img {
  max-width: 100%;
  max-height: 400px;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
}

.step-details div {
  margin-bottom: 0.75rem;
  color: #6b7280;
}

.step-details strong {
  color: #374151;
  display: inline-block;
  min-width: 150px;
}

.loading, .error, .no-data {
  text-align: center;
  padding: 3rem;
  color: #6b7280;
}

.error {
  color: #dc2626;
}
</style>
