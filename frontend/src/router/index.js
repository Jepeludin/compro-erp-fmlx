import { createRouter, createWebHistory } from 'vue-router';
import LoginPage from '../components/LoginPage.vue';
import Dashboard from '../components/Dashboard.vue';
import AdminLayout from '../components/AdminLayout.vue';
import AdminPage from '../components/AdminPage.vue';
import AdminMachine from '../components/AdminMachine.vue';
import AdminSchedule from '../components/AdminSchedule.vue';
import api from '../services/api.js';

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

// Definisikan routes
const routes = [
  {
    path: '/',
    // Mengarahkan halaman utama (root) langsung ke /login
    redirect: '/login', 
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginPage,
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  // Placeholder routes untuk menu dashboard
  {
    path: '/ppic',
    name: 'PPIC',
    component: () => import('../components/Dashboard.vue'),
    meta: { requiresAuth: true, allowedRoles: ['Admin', 'PPIC'] }
  },
  {
    path: '/toolpather',
    name: 'Toolpather',
    component: () => import('../components/Dashboard.vue'),
    meta: { requiresAuth: true, allowedRoles: ['Admin', 'Toolpather'] }
  },
  {
    path: '/pem',
    name: 'PEM',
    component: () => import('../components/Dashboard.vue'),
    meta: { requiresAuth: true, allowedRoles: ['Admin', 'Toolpather', 'PEM', 'QC', 'Engineering'] }
  },
  {
    path: '/qc',
    name: 'QC',
    component: () => import('../components/Dashboard.vue'),
    meta: { requiresAuth: true, allowedRoles: ['Admin', 'QC'] }
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true, requiresAdmin: true },
    redirect: '/admin/users',
    children: [
      {
        path: 'users',
        name: 'AdminUsers',
        component: AdminPage,
      },
      {
        path: 'machines',
        name: 'AdminMachines',
        component: AdminMachine,
      },
      {
        path: 'schedule',
        name: 'AdminSchedule',
        component: AdminSchedule,
      }
    ]
  },
  {
    path: '/database',
    name: 'Database',
    component: () => import('../components/Dashboard.vue'),
    meta: { requiresAuth: true, allowedRoles: ['Admin', 'PPIC', 'Toolpather', 'PEM', 'QC', 'Engineering', 'Guest'] }
  },
  {
    path: '/timetrack',
    name: 'TimeTrack',
    component: () => import('../components/Dashboard.vue'),
    meta: { requiresAuth: true, allowedRoles: ['Admin', 'PPIC', 'Toolpather', 'PEM', 'QC', 'Engineering', 'Guest'] }
  },
  {
    path: '/reporttrack',
    name: 'ReportTrack',
    component: () => import('../components/Dashboard.vue'),
    meta: { requiresAuth: true, allowedRoles: ['Admin', 'PPIC', 'Toolpather', 'PEM', 'QC', 'Engineering', 'Guest'] }
  }
];

// Buat instance router
const router = createRouter({
  history: createWebHistory(),
  routes, // sama dengan `routes: routes`
});

// Navigation guard untuk protect routes yang memerlukan auth dan role
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token');
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth);
  const requiresAdmin = to.matched.some(record => record.meta.requiresAdmin);
  const allowedRoles = to.matched.find(record => record.meta.allowedRoles)?.meta.allowedRoles;

  if (requiresAuth && !token) {
    // Jika route memerlukan auth tapi tidak ada token, redirect ke login
    next('/login');
  } else if (requiresAdmin) {
    // Check if user is admin
    const user = api.getCurrentUser();
    if (!user || user.role !== 'Admin') {
      // Not admin, redirect to dashboard with error message
      alert('Access Denied! Admin access required.');
      next('/dashboard');
    } else {
      next();
    }
  } else if (allowedRoles) {
    // Check if user has required role
    const user = api.getCurrentUser();
    if (!user) {
      next('/login');
    } else if (!allowedRoles.includes(user.role)) {
      // User doesn't have permission
      alert(`Access Denied! Your role (${user.role}) doesn't have permission to access this page.`);
      next('/dashboard');
    } else {
      next();
    }
  } else if (to.path === '/login' && token) {
    // Jika sudah login dan mencoba akses halaman login, redirect ke dashboard
    next('/dashboard');
  } else {
    next();
  }
});

export default router;