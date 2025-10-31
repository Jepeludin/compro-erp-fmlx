<template>
  <div class="dashboard-container">
    <!-- Header -->
    <header class="dashboard-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="dashboard-title">IMETRAX</h1>
          <p class="dashboard-subtitle">Welcome back, {{ userName }}</p>
        </div>
        <div class="header-right">
          <div class="user-info">
            <div class="user-avatar">{{ userInitial }}</div>
            <div class="user-details">
              <span class="user-name">{{ userName }}</span>
              <span class="user-role">{{ userRole }}</span>
            </div>
          </div>
          <button @click="handleLogout" class="logout-button">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
              <polyline points="16 17 21 12 16 7"></polyline>
              <line x1="21" y1="12" x2="9" y2="12"></line>
            </svg>
            Logout
          </button>
        </div>
      </div>
    </header>

    <!-- Dashboard Cards -->
    <main class="dashboard-main">
      <div class="cards-grid">
        <div 
          v-for="card in visibleCards" 
          :key="card.id"
          @click="handleCardClick(card)"
          class="dashboard-card"
        >
          <div class="card-icon">
            <component :is="'div'" v-html="card.icon"></component>
          </div>
          <h3 class="card-title">{{ card.title }}</h3>
          <p class="card-description">{{ card.description }}</p>
          <div class="card-arrow">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="5" y1="12" x2="19" y2="12"></line>
              <polyline points="12 5 19 12 12 19"></polyline>
            </svg>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import api from '../services/api.js';

const router = useRouter();

// User data
const userName = ref('');
const userRole = ref('');
const userInitial = computed(() => {
  return userName.value ? userName.value.charAt(0).toUpperCase() : 'U';
});

// Role-based access control configuration
const rolePermissions = {
  'Admin': ['ppic', 'toolpather', 'pem', 'qc', 'admin', 'database', 'timetrack', 'reporttrack'],
  'PPIC': ['ppic', 'database', 'timetrack', 'reporttrack'],
  'Toolpather': ['toolpather', 'pem', 'database', 'timetrack', 'reporttrack'],
  'PEM': ['pem', 'database', 'timetrack', 'reporttrack'],
  'QC': ['qc', 'pem', 'database', 'timetrack', 'reporttrack'],
  'Engineering': ['pem', 'database', 'timetrack', 'reporttrack'],
  'Guest': ['database', 'timetrack', 'reporttrack']
};

// Dashboard cards configuration
const dashboardCards = ref([
  {
    id: 'ppic',
    title: 'PPIC',
    description: 'Production Planning & Inventory Control',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <rect x="3" y="3" width="18" height="18" rx="2" ry="2"></rect>
      <line x1="9" y1="9" x2="15" y2="9"></line>
      <line x1="9" y1="15" x2="15" y2="15"></line>
    </svg>`,
    route: '/ppic'
  },
  {
    id: 'toolpather',
    title: 'Toolpather',
    description: 'Tool Path Management System',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"></path>
    </svg>`,
    route: '/toolpather'
  },
  {
    id: 'pem',
    title: 'PEM',
    description: 'Process Engineering Mold',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="12" cy="12" r="10"></circle>
      <polyline points="12 6 12 12 16 14"></polyline>
    </svg>`,
    route: '/pem'
  },
  {
    id: 'qc',
    title: 'QC',
    description: 'Quality Control & Assurance',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <polyline points="20 6 9 17 4 12"></polyline>
    </svg>`,
    route: '/qc'
  },
  {
    id: 'admin',
    title: 'Admin',
    description: 'System Administration',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path d="M12 2L2 7l10 5 10-5-10-5z"></path>
      <path d="M2 17l10 5 10-5"></path>
      <path d="M2 12l10 5 10-5"></path>
    </svg>`,
    route: '/admin'
  },
  {
    id: 'database',
    title: 'Database',
    description: 'Data Management System',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <ellipse cx="12" cy="5" rx="9" ry="3"></ellipse>
      <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"></path>
      <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"></path>
    </svg>`,
    route: '/database'
  },
  {
    id: 'timetrack',
    title: 'Time Track',
    description: 'Time Tracking & Monitoring',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <circle cx="12" cy="12" r="10"></circle>
      <polyline points="12 6 12 12 16 14"></polyline>
    </svg>`,
    route: '/timetrack'
  },
  {
    id: 'reporttrack',
    title: 'Report Track',
    description: 'Reporting & Analytics',
    icon: `<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <line x1="18" y1="20" x2="18" y2="10"></line>
      <line x1="12" y1="20" x2="12" y2="4"></line>
      <line x1="6" y1="20" x2="6" y2="14"></line>
    </svg>`,
    route: '/reporttrack'
  }
]);

// Computed property to filter cards based on user role
const visibleCards = computed(() => {
  if (!userRole.value) return [];
  
  const allowedPages = rolePermissions[userRole.value] || [];
  return dashboardCards.value.filter(card => allowedPages.includes(card.id));
});

// Load user data
onMounted(() => {
  const user = api.getCurrentUser();
  if (user) {
    userName.value = user.username;
    userRole.value = user.role;
  } else {
    // Jika tidak ada user data, redirect ke login
    router.push('/');
  }
});

// Handle card click
const handleCardClick = (card) => {
  console.log('Card clicked:', card.title);
  
  const user = api.getCurrentUser();
  if (!user) {
    router.push('/');
    return;
  }

  // Check permission
  const allowedPages = rolePermissions[user.role] || [];
  if (!allowedPages.includes(card.id)) {
    alert(`Access Denied!\n\nYour role (${user.role}) doesn't have permission to access ${card.title}.`);
    return;
  }

  // Navigate to route
  if (card.route === '/admin') {
    if (user.role === 'Admin') {
      router.push(card.route);
    } else {
      alert('Access Denied!\n\nOnly Admin users can access this page.');
    }
  } else {
    // For other routes, show coming soon alert or navigate
    alert(`Navigating to ${card.title}...\nRoute: ${card.route}\n\nFitur ini akan segera tersedia!`);
    
    // Uncomment untuk navigasi real
    // router.push(card.route);
  }
};

// Handle logout
const handleLogout = () => {
  if (confirm('Are you sure you want to logout?')) {
    api.clearAuth();
    router.push('/');
  }
};
</script>

<style scoped>
@import './Dashboard.css';
</style>
