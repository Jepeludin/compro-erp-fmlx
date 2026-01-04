<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-container">
      <div class="modal-header">
        <h2>{{ editPlan ? 'Edit Operation Plan' : 'Create Operation Plan' }}</h2>
        <button @click="$emit('close')" class="btn-close">&times;</button>
      </div>

      <div class="modal-body">
        <form @submit.prevent="savePlan">
          <!-- Header Information -->
          <div class="form-section">
            <h3>Part Information</h3>
            <div class="form-grid">
              <div class="form-group">
                <label>Part Name *</label>
                <input
                  v-model="formData.part_name"
                  type="text"
                  required
                  placeholder="Enter part name"
                />
              </div>

              <div class="form-group">
                <label>Material</label>
                <input
                  v-model="formData.material"
                  type="text"
                  placeholder="e.g., AL 6061"
                />
              </div>

              <div class="form-group">
                <label>Dial Size</label>
                <input
                  v-model="formData.dial_size"
                  type="text"
                  placeholder="e.g., 16 x 135 mm"
                />
              </div>

              <div class="form-group">
                <label>Quantity</label>
                <input
                  v-model.number="formData.quantity"
                  type="number"
                  min="1"
                  placeholder="0"
                />
              </div>

              <div class="form-group">
                <label>Revision</label>
                <input
                  v-model="formData.revision"
                  type="text"
                  placeholder="e.g., Rev. A"
                />
              </div>

              <div class="form-group">
                <label>No. WP</label>
                <input
                  v-model="formData.no_wp"
                  type="text"
                  placeholder="Work Package number"
                />
              </div>

              <div class="form-group">
                <label>Page</label>
                <input
                  v-model="formData.page"
                  type="text"
                  placeholder="e.g., 1/2"
                />
              </div>
            </div>
          </div>

          <!-- Process Steps -->
          <div class="form-section">
            <div class="section-header">
              <h3>Process Steps</h3>
              <button type="button" @click="addStep" class="btn-add-step">
                + Add Step
              </button>
            </div>

            <div v-if="formData.steps.length === 0" class="no-steps">
              No steps added yet. Click "Add Step" to begin.
            </div>

            <div v-for="(step, index) in formData.steps" :key="index" class="step-card">
              <div class="step-header">
                <h4>Step {{ index + 1 }}</h4>
                <button type="button" @click="removeStep(index)" class="btn-remove-step">
                  Remove
                </button>
              </div>

              <div class="step-grid">
                <div class="form-group full-width">
                  <label>Picture (PNG, JPG, PDF - Max 5MB)</label>
                  <input
                    type="file"
                    accept=".png,.jpg,.jpeg,.pdf"
                    @change="handleImageUpload(index, $event)"
                    class="file-input"
                  />
                  <div v-if="step.imagePreview" class="image-preview">
                    <img v-if="isImage(step.imageFile)" :src="step.imagePreview" alt="Preview" />
                    <div v-else class="pdf-preview">
                      <span>📄 {{ step.imageFile?.name }}</span>
                    </div>
                    <button type="button" @click="removeImage(index)" class="btn-remove-image">
                      &times;
                    </button>
                  </div>
                </div>

                <div class="form-group">
                  <label>Clamping System</label>
                  <textarea
                    v-model="step.clamping_system"
                    rows="2"
                    placeholder="Describe clamping method"
                  ></textarea>
                </div>

                <div class="form-group">
                  <label>Raw Material</label>
                  <textarea
                    v-model="step.raw_material"
                    rows="2"
                    placeholder="Material specifications"
                  ></textarea>
                </div>

                <div class="form-group">
                  <label>Setting</label>
                  <textarea
                    v-model="step.setting"
                    rows="2"
                    placeholder="Machine settings"
                  ></textarea>
                </div>

                <div class="form-group">
                  <label>Process</label>
                  <textarea
                    v-model="step.process"
                    rows="2"
                    placeholder="Process description"
                  ></textarea>
                </div>

                <div class="form-group">
                  <label>Note</label>
                  <textarea
                    v-model="step.note"
                    rows="2"
                    placeholder="Additional notes"
                  ></textarea>
                </div>

                <div class="form-group">
                  <label>Checking Method</label>
                  <textarea
                    v-model="step.checking_method"
                    rows="2"
                    placeholder="Quality check method"
                  ></textarea>
                </div>
              </div>
            </div>
          </div>

          <!-- Approver Selection -->
          <div class="form-section">
            <h3>Approvers (5 Required)</h3>
            <div class="approvers-grid">
              <div v-for="role in approverRoles" :key="role.key" class="approver-group">
                <label>{{ role.label }} *</label>
                <button
                  type="button"
                  @click="openApproverModal(role.key)"
                  class="btn-select-approver"
                >
                  {{ selectedApprovers[role.key] ? selectedApprovers[role.key].username : `Select ${role.label}` }}
                </button>
                <p v-if="selectedApprovers[role.key]" class="approver-info">
                  {{ selectedApprovers[role.key].email }} ({{ selectedApprovers[role.key].role }})
                </p>
              </div>
            </div>
          </div>

          <!-- Form Actions -->
          <div class="form-actions">
            <button type="button" @click="$emit('close')" class="btn-cancel">
              Cancel
            </button>
            <button type="submit" :disabled="saving" class="btn-save">
              {{ saving ? 'Saving...' : 'Save as Draft' }}
            </button>
            <button
              type="button"
              @click="saveAndSubmit"
              :disabled="saving || !areAllApproversSelected()"
              class="btn-submit"
            >
              {{ saving ? 'Submitting...' : 'Save & Submit for Approval' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Approver Selection Modal -->
    <ApproverModal
      v-if="showApproverModal"
      @close="showApproverModal = false"
      @select="selectApprover"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue';
import api from '../../services/api.js';
import ApproverModal from './ApproverModal.vue';

const props = defineProps({
  schedule: Object,
  editPlan: Object
});

const emit = defineEmits(['close', 'saved']);

const formData = reactive({
  ppic_schedule_id: props.editPlan?.ppic_schedule_id || props.schedule?.id || null,
  part_name: props.editPlan?.part_name || props.schedule?.part_name || '',
  material: props.editPlan?.material || '',
  dial_size: props.editPlan?.dial_size || '',
  quantity: props.editPlan?.quantity || 0,
  revision: props.editPlan?.revision || '',
  no_wp: props.editPlan?.no_wp || '',
  page: props.editPlan?.page || '',
  steps: props.editPlan?.steps?.map(step => ({
    id: step.id,
    step_number: step.step_number,
    clamping_system: step.clamping_system || '',
    raw_material: step.raw_material || '',
    setting: step.setting || '',
    process: step.process || '',
    note: step.note || '',
    checking_method: step.checking_method || '',
    picture_url: step.picture_url || null,
    imageFile: null,
    imagePreview: step.picture_url ? getImageUrl(step.picture_url) : null
  })) || []
});

const saving = ref(false);
const showApproverModal = ref(false);

// Initialize approvers from editPlan if editing
const initializeApprovers = () => {
  console.log('Initializing approvers, editPlan:', props.editPlan);
  if (props.editPlan?.approvals) {
    console.log('Approvals found:', props.editPlan.approvals);
    const approvers = {};
    props.editPlan.approvals.forEach(approval => {
      console.log('Processing approval:', approval);
      if (approval.approver) {
        approvers[approval.approver_role] = approval.approver;
      } else {
        // Initialize with null if approver not assigned yet
        approvers[approval.approver_role] = null;
      }
    });
    console.log('Initialized approvers:', approvers);
    return approvers;
  }
  console.log('No editPlan or approvals, returning empty');
  return {
    PEM: null,
    Toolpather: null,
    QC: null,
    Custom1: null,
    Custom2: null
  };
};

const selectedApprovers = ref(initializeApprovers());

function getImageUrl(relPath) {
  // Remove /api/v1 from baseURL to get the root server URL
  const serverUrl = api.baseURL.replace('/api/v1', '');
  return `${serverUrl}/uploads/operation-plan-images/${relPath}`;
}
const currentSelectingRole = ref(null);
const createdPlanId = ref(null);

const approverRoles = [
  { key: 'PEM', label: 'PEM' },
  { key: 'Toolpather', label: 'Toolpather' },
  { key: 'QC', label: 'QC' },
  { key: 'Custom1', label: 'Custom Role 1' },
  { key: 'Custom2', label: 'Custom Role 2' }
];

function addStep() {
  formData.steps.push({
    step_number: formData.steps.length + 1,
    clamping_system: '',
    raw_material: '',
    setting: '',
    process: '',
    note: '',
    checking_method: '',
    imageFile: null,
    imagePreview: null
  });
}

function removeStep(index) {
  formData.steps.splice(index, 1);
  // Renumber steps
  formData.steps.forEach((step, idx) => {
    step.step_number = idx + 1;
  });
}

function handleImageUpload(index, event) {
  const file = event.target.files[0];
  if (!file) return;

  // Validate file size (5MB)
  if (file.size > 5 * 1024 * 1024) {
    alert('File size must be less than 5MB');
    event.target.value = '';
    return;
  }

  // Validate file type
  const validTypes = ['image/png', 'image/jpeg', 'image/jpg', 'application/pdf'];
  if (!validTypes.includes(file.type)) {
    alert('Only PNG, JPG, and PDF files are allowed');
    event.target.value = '';
    return;
  }

  formData.steps[index].imageFile = file;

  // Create preview for images
  if (file.type.startsWith('image/')) {
    const reader = new FileReader();
    reader.onload = (e) => {
      formData.steps[index].imagePreview = e.target.result;
    };
    reader.readAsDataURL(file);
  } else {
    formData.steps[index].imagePreview = 'pdf';
  }
}

function removeImage(index) {
  formData.steps[index].imageFile = null;
  formData.steps[index].imagePreview = null;
}

function isImage(file) {
  return file?.type.startsWith('image/');
}

function openApproverModal(roleKey) {
  currentSelectingRole.value = roleKey;
  showApproverModal.value = true;
}

function selectApprover(approver) {
  if (currentSelectingRole.value) {
    selectedApprovers.value[currentSelectingRole.value] = approver;
  }
  showApproverModal.value = false;
  currentSelectingRole.value = null;
}

async function savePlan() {
  saving.value = true;

  try {
    const planData = {
      ppic_schedule_id: formData.ppic_schedule_id ? parseInt(formData.ppic_schedule_id) : null,
      part_name: formData.part_name,
      material: formData.material,
      dial_size: formData.dial_size,
      quantity: formData.quantity,
      revision: formData.revision,
      no_wp: formData.no_wp,
      page: formData.page,
      steps: formData.steps.map(step => ({
        step_number: step.step_number,
        clamping_system: step.clamping_system,
        raw_material: step.raw_material,
        setting: step.setting,
        process: step.process,
        note: step.note,
        checking_method: step.checking_method
      }))
    };

    console.log('Sending plan data:', JSON.stringify(planData, null, 2));

    let response;
    if (props.editPlan) {
      // Update existing plan
      response = await api.updatePEMOperationPlan(props.editPlan.id, planData);
    } else {
      // Create new plan
      response = await api.createPEMOperationPlan(planData);
    }

    if (response.success && response.data) {
      createdPlanId.value = response.data.id;

      // Upload images for steps that have them
      await uploadStepImages(response.data);

      alert(`Operation plan ${props.editPlan ? 'updated' : 'saved'} successfully!`);
      emit('saved');
    }
  } catch (err) {
    alert(`Failed to ${props.editPlan ? 'update' : 'save'} plan: ` + err.message);
    console.error('Error saving plan:', err);
  } finally {
    saving.value = false;
  }
}

async function uploadStepImages(plan) {
  // Find steps with images
  const stepsWithImages = formData.steps
    .map((step, index) => ({ ...step, originalIndex: index }))
    .filter(step => step.imageFile);

  for (const step of stepsWithImages) {
    // Find the corresponding step ID from the saved plan
    const savedStep = plan.steps?.find(s => s.step_number === step.step_number);
    if (!savedStep) continue;

    try {
      const formData = new FormData();
      formData.append('image', step.imageFile);
      await api.uploadStepImage(savedStep.id, formData);
    } catch (err) {
      console.error(`Failed to upload image for step ${step.step_number}:`, err);
    }
  }
}

function areAllApproversSelected() {
  return Object.values(selectedApprovers.value).every(approver => approver !== null);
}

async function saveAndSubmit() {
  if (!areAllApproversSelected()) {
    alert('Please select all 5 approvers before submitting');
    return;
  }

  saving.value = true;

  try {
    // First save the plan
    await savePlan();

    // Get plan ID - use editPlan.id if editing, otherwise use createdPlanId
    const planId = props.editPlan?.id || createdPlanId.value;

    if (!planId) {
      throw new Error('Failed to save plan');
    }

    // Assign all 5 approvers (map role to user ID)
    const approvers = {};
    Object.keys(selectedApprovers.value).forEach(role => {
      approvers[role] = selectedApprovers.value[role].id;
    });

    console.log('Assigning approvers to plan ID:', planId, 'Approvers:', approvers);
    await api.assignPEMApprovers(planId, approvers);

    // Submit for approval
    console.log('Submitting plan ID:', planId, 'for approval');
    await api.submitPEMOperationPlan(planId);

    alert('Operation plan submitted for approval successfully!');
    emit('saved');
  } catch (err) {
    alert('Failed to submit plan: ' + err.message);
    console.error('Error submitting plan:', err);
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  // Add one initial step
  addStep();
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
  z-index: 1000;
  overflow-y: auto;
  padding: 2rem;
}

.modal-container {
  background: white;
  border-radius: 12px;
  width: 100%;
  max-width: 1200px;
  max-height: 90vh;
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
  padding: 2rem;
}

.form-section {
  margin-bottom: 2rem;
}

.form-section h3 {
  margin: 0 0 1rem 0;
  color: #1f2937;
  font-size: 1.25rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
}

.form-group.full-width {
  grid-column: 1 / -1;
}

.form-group label {
  margin-bottom: 0.5rem;
  color: #374151;
  font-weight: 500;
  font-size: 0.9rem;
}

.form-group input,
.form-group textarea {
  padding: 0.6rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 0.95rem;
  font-family: inherit;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
}

.btn-add-step {
  padding: 0.6rem 1.2rem;
  background: #667eea;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
}

.btn-add-step:hover {
  background: #5568d3;
}

.no-steps {
  text-align: center;
  padding: 3rem;
  color: #6b7280;
  background: #f9fafb;
  border-radius: 8px;
}

.step-card {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1rem;
  background: #f9fafb;
}

.step-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.step-header h4 {
  margin: 0;
  color: #1f2937;
}

.btn-remove-step {
  padding: 0.4rem 0.8rem;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85rem;
}

.btn-remove-step:hover {
  background: #dc2626;
}

.step-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.file-input {
  padding: 0.5rem !important;
}

.image-preview {
  position: relative;
  margin-top: 0.5rem;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  overflow: hidden;
}

.image-preview img {
  width: 100%;
  max-height: 200px;
  object-fit: contain;
  display: block;
}

.pdf-preview {
  padding: 2rem;
  background: #f3f4f6;
  text-align: center;
  color: #6b7280;
}

.btn-remove-image {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  width: 30px;
  height: 30px;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  font-size: 1.2rem;
  line-height: 1;
}

.btn-remove-image:hover {
  background: #dc2626;
}

.approvers-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.approver-group {
  display: flex;
  flex-direction: column;
}

.approver-group label {
  margin-bottom: 0.5rem;
  color: #374151;
  font-weight: 500;
  font-size: 0.9rem;
}

.btn-select-approver {
  width: 100%;
  padding: 0.8rem;
  background: white;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  text-align: left;
  cursor: pointer;
  font-size: 0.95rem;
}

.btn-select-approver:hover {
  border-color: #667eea;
}

.approver-info {
  margin-top: 0.5rem;
  color: #6b7280;
  font-size: 0.85rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  padding-top: 2rem;
  border-top: 1px solid #e5e7eb;
}

.btn-cancel,
.btn-save,
.btn-submit {
  padding: 0.8rem 1.5rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.95rem;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-cancel {
  background: #e5e7eb;
  color: #374151;
}

.btn-cancel:hover {
  background: #d1d5db;
}

.btn-save {
  background: #667eea;
  color: white;
}

.btn-save:hover:not(:disabled) {
  background: #5568d3;
}

.btn-submit {
  background: #10b981;
  color: white;
}

.btn-submit:hover:not(:disabled) {
  background: #059669;
}

button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
