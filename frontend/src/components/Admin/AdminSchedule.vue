<template>
  <div class="placeholder-container">
    <div class="placeholder-content">
      <div class="placeholder-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
          <line x1="16" y1="2" x2="16" y2="6"></line>
          <line x1="8" y1="2" x2="8" y2="6"></line>
          <line x1="3" y1="10" x2="21" y2="10"></line>
        </svg>
      </div>
      <h1 class="placeholder-title">Task Schedule</h1>
      <p class="placeholder-text">This feature is coming soon...</p>
      <p class="placeholder-subtext">Schedule and manage production tasks</p>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.placeholder-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: #f5f7fa;
}

.placeholder-content {
  text-align: center;
  max-width: 500px;
}

.placeholder-icon {
  color: #cbd5e1;
  margin-bottom: 1.5rem;
}

.placeholder-title {
  font-size: 2rem;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 0.5rem;
}

.placeholder-text {
  font-size: 1.125rem;
  color: #64748b;
  margin-bottom: 0.5rem;
}

.placeholder-subtext {
  font-size: 0.875rem;
  color: #94a3b8;
}
</style>
