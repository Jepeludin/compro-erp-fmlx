<template>
  <div class="machine-detail-wrapper">
    <!-- Header -->
    <header class="machine-detail-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="logo">IMETRAX</h1>
          <span class="divider">|</span>
          <h2 class="page-title">{{ machineName }} - Time Tracking</h2>
        </div>
        <div class="header-right">
          <div class="update-time">Last Update: {{ lastUpdate }}</div>
          <button @click="goBack" class="btn-back">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            Back to Monitoring
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="machine-detail-body">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading machine schedules...</p>
      </div>

      <div v-else-if="error" class="error-state">
        <p>{{ error }}</p>
        <button @click="fetchData" class="retry-button">Retry</button>
      </div>

      <div v-else-if="schedules.length === 0" class="empty-state">
        <p class="empty-state-text">No work orders scheduled for this machine</p>
        <button @click="goBack" class="add-first-btn">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M19 12H5M12 19l-7-7 7-7"/>
          </svg>
          Back to Monitoring
        </button>
      </div>

      <div v-else class="work-orders-container">
        <!-- Work Order Card -->
        <div
          v-for="(schedule, index) in schedules"
          :key="schedule.id"
          class="work-order-card"
        >
          <div class="work-order-header" @click="toggleExpand(schedule.id)">
            <div class="order-badge">#{{ index + 1 }}</div>
            <div class="order-info-grid">
              <div class="info-item">
                <div class="info-label">Work Order</div>
                <div class="info-value">{{ schedule.njo || '-' }}</div>
              </div>
              <div class="info-item">
                <div class="info-label">Start Date</div>
                <div class="info-value">{{ formatDateTime(schedule.start_date) }}</div>
              </div>
              <!-- <div class="info-item">
                <div class="info-label">Project</div>
                <div class="info-value">{{ schedule.project_name || '-' }}</div>
              </div> -->
              <div class="info-item">
                <div class="info-label">Part Name</div>
                <div class="info-value">{{ schedule.part_name || '-' }}</div>
              </div>
              <div class="info-item">
                <div class="info-label">Status</div>
                <div>
                  <span class="status-badge" :class="getStatusClass(schedule)">
                    {{ schedule.status || 'Open' }}
                  </span>
                </div>
              </div>
              <div class="info-item">
                <div class="info-label">Operation Plan</div>
                <div>
                  <span class="status-badge not-ready">Not Ready</span>
                </div>
              </div>
              <div class="info-item">
                <div class="info-label">NC Sheet & G-Code</div>
                <div>
                  <span class="status-badge not-ready">Not Ready</span>
                </div>
              </div>
              <div class="info-item">
                <div class="info-label">Time Log</div>
                <div class="time-display">
                  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <polyline points="12 6 12 12 16 14"/>
                  </svg>
                  <span>{{ calculateDuration(schedule.start_date, schedule.finish_date) }}</span>
                </div>
              </div>
            </div>
            <svg 
              class="expand-icon" 
              :class="{ 'expanded': expandedOrders.includes(schedule.id) }"
              xmlns="http://www.w3.org/2000/svg" 
              width="24" 
              height="24" 
              viewBox="0 0 24 24" 
              fill="none" 
              stroke="currentColor" 
              stroke-width="2"
            >
              <polyline points="6 9 12 15 18 9"/>
            </svg>
          </div>

          <!-- Details Section (Expandable) -->
          <div class="details-section" :class="{ 'show': expandedOrders.includes(schedule.id) }">
            <template v-if="schedule.operations && schedule.operations.length > 0">
              <!-- Process Rows -->
              <div
                v-for="(operation, opIndex) in schedule.operations"
                :key="operation.id"
                class="process-row"
              >
                <div class="process-header-row">
                  <div class="info-item">
                    <div class="info-label">Sequence</div>
                    <div class="info-value">x{{ operation.sequence || (opIndex + 1) }}</div>
                  </div>
                  <div class="info-item">
                    <div class="info-label">WO</div>
                    <div class="info-value">{{ schedule.njo || '-' }}</div>
                  </div>
                  <div class="info-item">
                    <div class="info-label">Process</div>
                    <div class="info-value">{{ operation.operation_name }}</div>
                  </div>
                </div>

                <div class="process-content">
                  <div class="left-section">
                    <div>
                      <div class="info-label">Note</div>
                      <textarea class="input-field notes-area" placeholder="Add notes..."></textarea>
                    </div>
                    <div>
                      <div class="info-label">Operator</div>
                      <input type="text" class="input-field" placeholder="Operator name" />
                    </div>
                  </div>

                  <!-- Phases Container -->
                  <div class="phases-container">
                    <div
                      v-for="phase in phases"
                      :key="phase.id"
                      class="phase-item"
                      :class="{ 'running': activePhase === `${schedule.id}-${operation.id}-${phase.id}` }"
                    >
                      <div class="phase-header-inline">
                        <span class="phase-name">
                          <span 
                            class="status-indicator" 
                            :class="{ 'active': activePhase === `${schedule.id}-${operation.id}-${phase.id}` }"
                          ></span>
                          {{ phase.name }}
                        </span>
                      </div>

                      <div class="timer-text">
                        {{ formatTimerDisplay(`${schedule.id}-${operation.id}-${phase.id}`) }}
                      </div>

                      <div class="time-range">
                        <div class="time-range-row">
                          <span class="time-label">Start:</span>
                          <span class="time-value">{{ getPhaseStartTime(`${schedule.id}-${operation.id}-${phase.id}`) }}</span>
                        </div>
                        <div class="time-range-row">
                          <span class="time-label">End:</span>
                          <span class="time-value">{{ getPhaseEndTime(`${schedule.id}-${operation.id}-${phase.id}`) }}</span>
                        </div>
                      </div>

                      <button
                        v-if="activePhase !== `${schedule.id}-${operation.id}-${phase.id}`"
                        class="control-btn-small start"
                        @click="startPhase(`${schedule.id}-${operation.id}-${phase.id}`)"
                        :disabled="activePhase !== null"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                          <polygon points="5 3 19 12 5 21 5 3"/>
                        </svg>
                        Start
                      </button>

                      <button
                        v-else
                        class="control-btn-small stop"
                        @click="stopPhase(`${schedule.id}-${operation.id}-${phase.id}`)"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                          <rect x="6" y="6" width="12" height="12"/>
                        </svg>
                        Stop
                      </button>

                      <div class="time-meta">
                        Total: {{ getPhaseTotal(`${schedule.id}-${operation.id}-${phase.id}`) }}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              
              <button class="add-process-btn">
                <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 5v14M5 12h14"/>
                </svg>
                Add Process Row
              </button>
            </template>
            
            <template v-else>
              <div class="empty-operations">
                <p>No processes yet for this work order</p>
                <button class="add-first-btn">
                  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M12 5v14M5 12h14"/>
                  </svg>
                  Add First Process
                </button>
              </div>
            </template>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import api from '../../services/api.js';

const router = useRouter();
const route = useRoute();

// State
const loading = ref(true);
const error = ref(null);
const machineId = ref(null);
const machineName = ref('');
const schedules = ref([]);
const lastUpdate = ref('');
const expandedOrders = ref([]);
const activePhase = ref(null);
const timers = ref({});
const timerIntervals = ref({});

// Phase definitions (sesuai dengan template HTML)
const phases = [
  { id: 'proses', name: 'Proses' },
  { id: 'setting', name: 'Setting' },
  { id: 'cmm', name: 'CMM' },
  { id: 'multilevel', name: 'MultiLevel' }
];

// Methods
const fetchData = async () => {
  loading.value = true;
  error.value = null;

  try {
    machineId.value = route.params.id;
    
    // Fetch machine details and schedules
    const [machineData, schedulesData] = await Promise.all([
      api.getMachine(machineId.value),
      api.getAllPPICSchedules()
    ]);

    console.log('Machine Data:', machineData);
    console.log('Schedules Data:', schedulesData);

    // Handle different response structures
    machineName.value = machineData.machine?.machine_name || machineData.machine_name || 'Unknown Machine';

    // Filter schedules for this machine
    const allSchedules = schedulesData.data || [];
    schedules.value = allSchedules.filter(schedule => {
      return schedule.machine_assignments &&
             schedule.machine_assignments.some(ma => ma.machine_id === parseInt(machineId.value));
    }).map(schedule => {
      // If no operations, create default operation structure
      if (!schedule.operations || schedule.operations.length === 0) {
        schedule.operations = [
          {
            id: `default-${schedule.id}`,
            operation_name: schedule.part_name || 'Manufacturing Process',
            sequence: 1
          }
        ];
      }
      return schedule;
    });

    // Sort by start date
    schedules.value.sort((a, b) => new Date(a.start_date) - new Date(b.start_date));

    updateLastUpdateTime();
  } catch (err) {
    console.error('Error fetching data:', err);
    error.value = 'Failed to load machine schedules. Please try again.';
  } finally {
    loading.value = false;
  }
};

const updateLastUpdateTime = () => {
  const now = new Date();
  const year = now.getFullYear();
  const month = String(now.getMonth() + 1).padStart(2, '0');
  const day = String(now.getDate()).padStart(2, '0');
  const hours = String(now.getHours()).padStart(2, '0');
  const minutes = String(now.getMinutes()).padStart(2, '0');
  const seconds = String(now.getSeconds()).padStart(2, '0');
  lastUpdate.value = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

const formatDateTime = (dateString) => {
  if (!dateString) return '-';
  const date = new Date(dateString);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${year}-${month}-${day} ${hours}:${minutes}`;
};

const calculateDuration = (startDate, finishDate) => {
  if (!startDate || !finishDate) return '0h 0m';

  const start = new Date(startDate);
  const finish = new Date(finishDate);
  const diffTime = Math.abs(finish - start);
  const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
  const diffMinutes = Math.floor((diffTime % (1000 * 60 * 60)) / (1000 * 60));

  return `${diffHours}h ${diffMinutes}m`;
};

const getStatusClass = (schedule) => {
  const status = schedule.status?.toLowerCase() || 'open';
  if (status === 'completed') return 'status-completed';
  if (status === 'in progress' || status === 'in-progress') return 'in-progress';
  if (status === 'ready') return 'ready';
  return 'open';
};

const toggleExpand = (orderId) => {
  const index = expandedOrders.value.indexOf(orderId);
  if (index > -1) {
    expandedOrders.value.splice(index, 1);
  } else {
    expandedOrders.value.push(orderId);
  }
};

// Timer functions
const startPhase = (phaseId) => {
  if (activePhase.value) return;

  activePhase.value = phaseId;
  
  if (!timers.value[phaseId]) {
    timers.value[phaseId] = {
      startTime: Date.now(),
      endTime: null,
      elapsed: 0,
      isRunning: true
    };
  } else {
    timers.value[phaseId].startTime = Date.now() - timers.value[phaseId].elapsed;
    timers.value[phaseId].isRunning = true;
  }

  // Start interval for this timer
  timerIntervals.value[phaseId] = setInterval(() => {
    if (timers.value[phaseId]) {
      timers.value[phaseId].elapsed = Date.now() - timers.value[phaseId].startTime;
    }
  }, 1000);
};

const stopPhase = (phaseId) => {
  if (activePhase.value !== phaseId) return;

  if (timers.value[phaseId]) {
    timers.value[phaseId].endTime = Date.now();
    timers.value[phaseId].elapsed = timers.value[phaseId].endTime - timers.value[phaseId].startTime;
    timers.value[phaseId].isRunning = false;
  }

  // Clear interval
  if (timerIntervals.value[phaseId]) {
    clearInterval(timerIntervals.value[phaseId]);
    delete timerIntervals.value[phaseId];
  }

  activePhase.value = null;
};

const formatTimerDisplay = (phaseId) => {
  if (!timers.value[phaseId]) return '00:00:00';

  const elapsed = timers.value[phaseId].elapsed;
  const hours = Math.floor(elapsed / (1000 * 60 * 60));
  const minutes = Math.floor((elapsed % (1000 * 60 * 60)) / (1000 * 60));
  const seconds = Math.floor((elapsed % (1000 * 60)) / 1000);

  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
};

const getPhaseStartTime = (phaseId) => {
  if (!timers.value[phaseId] || !timers.value[phaseId].startTime) return '--:--';

  const date = new Date(timers.value[phaseId].startTime);
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${hours}:${minutes}`;
};

const getPhaseEndTime = (phaseId) => {
  if (!timers.value[phaseId] || !timers.value[phaseId].endTime) return '--:--';

  const date = new Date(timers.value[phaseId].endTime);
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${hours}:${minutes}`;
};

const getPhaseTotal = (phaseId) => {
  if (!timers.value[phaseId]) return '0m';

  const elapsed = timers.value[phaseId].elapsed;
  const minutes = Math.floor(elapsed / (1000 * 60));
  return `${minutes}m`;
};

const goBack = () => {
  router.push('/timetrack');
};

// Lifecycle
onMounted(() => {
  fetchData();

  // Auto-refresh every 60 seconds
  const refreshInterval = setInterval(() => {
    if (!activePhase.value) {
      fetchData();
    }
  }, 60000);

  // Cleanup on unmount
  onUnmounted(() => {
    clearInterval(refreshInterval);
    // Clear all timer intervals
    Object.keys(timerIntervals.value).forEach(key => {
      clearInterval(timerIntervals.value[key]);
    });
  });
});
</script>

<style scoped>
@import './machinedetail.css';
</style>
