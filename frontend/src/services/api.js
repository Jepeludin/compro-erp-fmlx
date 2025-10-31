// api.js - API Service untuk Frontend
// Letakkan file ini di frontend/src/services/api.js

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

class ApiService {
  constructor() {
    this.baseURL = API_BASE_URL;
  }

  // Get auth token from localStorage
  getToken() {
    return localStorage.getItem('token');
  }

  // Get auth headers
  getHeaders(includeAuth = false) {
    const headers = {
      'Content-Type': 'application/json',
    };

    if (includeAuth) {
      const token = this.getToken();
      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    return headers;
  }

  // Generic request method
  async request(endpoint, options = {}) {
    const url = `${this.baseURL}${endpoint}`;
    
    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          ...this.getHeaders(options.auth),
          ...options.headers,
        },
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Request failed');
      }

      return data;
    } catch (error) {
      console.error('API Error:', error);
      throw error;
    }
  }

  // Auth endpoints
  async login(userId, password) {
    return this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ user_id: userId, password }),
    });
  }

  async register(userData) {
    // userData should include: username, user_id, password, role, operator
    return this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  }

  async getProfile() {
    return this.request('/auth/profile', {
      method: 'GET',
      auth: true,
    });
  }

  // Save user data to localStorage
  saveAuth(token, user) {
    localStorage.setItem('token', token);
    localStorage.setItem('user', JSON.stringify(user));
  }

  // Clear auth data
  clearAuth() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }

  // Get current user
  getCurrentUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
  }

  // Check if user is authenticated
  // Check if user is authenticated
  isAuthenticated() {
    return !!this.getToken();
  }

  // Admin endpoints
  async getAllUsers() {
    return this.request('/admin/users', {
      method: 'GET',
      auth: true,
    });
  }

  async updateUser(userId, userData) {
    // userData can include: password, role, is_active
    return this.request(`/admin/users/${userId}`, {
      method: 'PUT',
      auth: true,
      body: JSON.stringify(userData),
    });
  }

  async deleteUser(userId) {
    return this.request(`/admin/users/${userId}`, {
      method: 'DELETE',
      auth: true,
    });
  }

  // Machine endpoints
  async getAllMachines() {
    return this.request('/machines', {
      method: 'GET',
      auth: true,
    });
  }

  async getMachine(machineId) {
    return this.request(`/machines/${machineId}`, {
      method: 'GET',
      auth: true,
    });
  }

  async createMachine(machineData) {
    // machineData: { machine_code, machine_name, machine_type, location, status }
    return this.request('/admin/machines', {
      method: 'POST',
      auth: true,
      body: JSON.stringify(machineData),
    });
  }

  async updateMachine(machineId, machineData) {
    // machineData: { machine_name, machine_type, location, status }
    return this.request(`/admin/machines/${machineId}`, {
      method: 'PUT',
      auth: true,
      body: JSON.stringify(machineData),
    });
  }

  async deleteMachine(machineId) {
    return this.request(`/admin/machines/${machineId}`, {
      method: 'DELETE',
      auth: true,
    });
  }

  // Job Order endpoints
  async getAllJobOrders() {
    return this.request('/job-orders', {
      method: 'GET',
      auth: true,
    });
  }

  async getJobOrder(jobOrderId) {
    return this.request(`/job-orders/${jobOrderId}`, {
      method: 'GET',
      auth: true,
    });
  }

  async getJobOrdersByMachine(machineId) {
    return this.request(`/job-orders/machine/${machineId}`, {
      method: 'GET',
      auth: true,
    });
  }

  async createJobOrder(jobOrderData) {
    // jobOrderData: { machine_id, njo, project, item, note, deadline, operator_id }
    return this.request('/job-orders', {
      method: 'POST',
      auth: true,
      body: JSON.stringify(jobOrderData),
    });
  }

  async updateJobOrder(jobOrderId, jobOrderData) {
    // jobOrderData: { project, item, note, deadline, operator_id, status }
    return this.request(`/job-orders/${jobOrderId}`, {
      method: 'PUT',
      auth: true,
      body: JSON.stringify(jobOrderData),
    });
  }

  async deleteJobOrder(jobOrderId) {
    return this.request(`/job-orders/${jobOrderId}`, {
      method: 'DELETE',
      auth: true,
    });
  }

  // Process Stage endpoints
  async updateProcessStage(stageId, stageData) {
    // stageData: { start_time, finish_time, operator_id, notes }
    return this.request(`/process-stages/${stageId}`, {
      method: 'PUT',
      auth: true,
      body: JSON.stringify(stageData),
    });
  }
}

export default new ApiService();
