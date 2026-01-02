<template>
  <div class="timetrack-wrapper">
    <!-- Header -->
    <header class="timetrack-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="logo">IMETRAX</h1>
          <span class="divider">|</span>
          <h2 class="page-title">Machine Monitoring (IME)</h2>
        </div>
        <div class="header-right">
          <div class="update-time">Last Update: {{ lastUpdate }}</div>
          <button @click="goBack" class="btn-back">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            Back to Dashboard
          </button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="timetrack-body">
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading machine data...</p>
      </div>

      <div v-else-if="error" class="error-state">
        <p>{{ error }}</p>
        <button @click="fetchData" class="retry-button">Retry</button>
      </div>

      <div v-else class="machines-grid">
        <!-- Machine Card -->
        <div
          v-for="machine in machinesWithSchedules"
          :key="machine.machine_id"
          class="machine-card"
        >
          <div class="machine-header">
            <h3 class="machine-name">{{ machine.machine_name }}</h3>
            <span class="job-count">{{ machine.schedules.length }} Jobs</span>
          </div>

          <div class="jobs-list">
            <div
              v-for="schedule in machine.schedules"
              :key="schedule.id"
              class="job-item"
              :class="getStatusClass(schedule)"
            >
              <div class="job-header">
                <span class="job-priority">{{ schedule.priority || 'FAST' }}</span>
                <span class="job-number">{{ schedule.njo }}</span>
              </div>
              <div class="job-details">
                <p class="job-part">{{ schedule.part_name }}</p>
                <div class="job-meta">
                  <span class="job-duration">Dur: {{ calculateDuration(schedule.start_date, schedule.finish_date) }}</span>
                </div>
              </div>
              <div class="job-footer">
                <span class="job-date">{{ formatDate(schedule.start_date) }}</span>
                <div class="status-indicator" :class="getStatusIndicatorClass(schedule)">
                  <span class="status-dot"></span>
                </div>
              </div>
            </div>

            <!-- Empty state for machines with no jobs -->
            <div v-if="machine.schedules.length === 0" class="empty-jobs">
              <p>No jobs scheduled</p>
            </div>
          </div>
        </div>

        <!-- Empty state for no machines -->
        <div v-if="machinesWithSchedules.length === 0" class="no-machines">
          <p>No machines found</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import api from '../../services/api.js';

const router = useRouter();

// State
const loading = ref(true);
const error = ref(null);
const machines = ref([]);
const schedules = ref([]);
const lastUpdate = ref('');

// Computed
const machinesWithSchedules = computed(() => {
  return machines.value.map(machine => {
    // Filter schedules that have this machine assigned
    // Use machine.id instead of machine.machine_id (API returns 'id' not 'machine_id')
    const machineSchedules = schedules.value.filter(schedule => {
      return schedule.machine_assignments &&
             schedule.machine_assignments.some(ma => ma.machine_id === machine.id);
    });

    // Sort by start date (earliest first)
    machineSchedules.sort((a, b) => {
      return new Date(a.start_date) - new Date(b.start_date);
    });

    return {
      ...machine,
      machine_id: machine.id, // Add machine_id for compatibility
      schedules: machineSchedules
    };
  });
});

// Methods
const fetchData = async () => {
  loading.value = true;
  error.value = null;

  try {
    // Fetch machines and schedules in parallel
    const [machinesData, schedulesData] = await Promise.all([
      api.getAllMachines(),
      api.getAllPPICSchedules()
    ]);

    machines.value = machinesData.machines || [];
    schedules.value = schedulesData.data || [];

    // Update last update time
    updateLastUpdateTime();
  } catch (err) {
    console.error('Error fetching data:', err);
    error.value = 'Failed to load machine data. Please try again.';
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

const calculateDuration = (startDate, finishDate) => {
  if (!startDate || !finishDate) return '-';

  const start = new Date(startDate);
  const finish = new Date(finishDate);
  const diffTime = Math.abs(finish - start);
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

  return `${diffDays}`;
};

const formatDate = (dateString) => {
  if (!dateString) return '-';

  const date = new Date(dateString);
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const year = date.getFullYear();

  return `${month}/${day}/${year}`;
};

const getStatusClass = (schedule) => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const startDate = new Date(schedule.start_date);
  startDate.setHours(0, 0, 0, 0);

  const finishDate = new Date(schedule.finish_date);
  finishDate.setHours(0, 0, 0, 0);

  if (schedule.status === 'Completed') return 'status-completed';
  if (today >= startDate && today <= finishDate) return 'status-ongoing';
  if (today < startDate) return 'status-pending';
  if (today > finishDate) return 'status-overdue';

  return '';
};

const getStatusIndicatorClass = (schedule) => {
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const startDate = new Date(schedule.start_date);
  startDate.setHours(0, 0, 0, 0);

  const finishDate = new Date(schedule.finish_date);
  finishDate.setHours(0, 0, 0, 0);

  if (schedule.status === 'Completed') return 'indicator-completed';
  if (today >= startDate && today <= finishDate) return 'indicator-ongoing';
  if (today < startDate) return 'indicator-pending';
  if (today > finishDate) return 'indicator-overdue';

  return 'indicator-pending';
};

const goBack = () => {
  router.push('/dashboard');
};

// Lifecycle
onMounted(() => {
  fetchData();

  // Auto-refresh every 30 seconds
  const refreshInterval = setInterval(() => {
    fetchData();
  }, 30000);

  // Cleanup on unmount
  return () => {
    clearInterval(refreshInterval);
  };
});
</script>

<style scoped>
@import './timetrack.css';
</style>
