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
        // Include details from backend if available
        const errorMessage = data.details
          ? `${data.error}: ${data.details}`
          : (data.error || 'Request failed');
        throw new Error(errorMessage);
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

  // Get users for approver selection (accessible by all authenticated users)
  async getUsers() {
    return this.request('/users', {
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

  // PPIC Schedule endpoints
  async getAllPPICSchedules() {
    return this.request('/ppic-schedules', {
      method: 'GET',
      auth: true,
    });
  }

  async getPPICSchedule(scheduleId) {
    return this.request(`/ppic-schedules/${scheduleId}`, {
      method: 'GET',
      auth: true,
    });
  }

  async createPPICSchedule(scheduleData) {
    // scheduleData: { njo, part_name, priority, priority_alpha, material_status,
    //                 start_date, finish_date, ppic_notes, machine_assignments }
    return this.request('/ppic-schedules', {
      method: 'POST',
      auth: true,
      body: JSON.stringify(scheduleData),
    });
  }

  async updatePPICSchedule(scheduleId, scheduleData) {
    // scheduleData: { part_name, priority, material_status, status, progress,
    //                 start_date, finish_date, ppic_notes, machine_assignments }
    return this.request(`/ppic-schedules/${scheduleId}`, {
      method: 'PUT',
      auth: true,
      body: JSON.stringify(scheduleData),
    });
  }

  async deletePPICSchedule(scheduleId) {
    return this.request(`/ppic-schedules/${scheduleId}`, {
      method: 'DELETE',
      auth: true,
    });
  }

  async getSchedulesByMachine(machineId) {
    return this.request(`/ppic-schedules/machine/${machineId}`, {
      method: 'GET',
      auth: true,
    });
  }

  // Gantt Chart endpoint
  async getGanttChart(filters = {}) {
    const params = new URLSearchParams(filters).toString();
    const endpoint = params ? `/gantt-chart?${params}` : '/gantt-chart';
    return this.request(endpoint, {
      method: 'GET',
      auth: true,
    });
  }

  // PPIC Links endpoints
  async createPPICLink(linkData) {
    // linkData: { source_schedule_id, target_schedule_id, link_type }
    return this.request('/ppic-links', {
      method: 'POST',
      auth: true,
      body: JSON.stringify(linkData),
    });
  }

  async deletePPICLink(linkId) {
    return this.request(`/ppic-links/${linkId}`, {
      method: 'DELETE',
      auth: true,
    });
  }

  async getAllPPICLinks() {
    return this.request('/ppic-links', {
      method: 'GET',
      auth: true,
    });
  }

  // Google Sheets endpoints
  async getPartNameByOrderNumber(orderNumber) {
    return this.request(`/google-sheets/part-name/${encodeURIComponent(orderNumber)}`, {
      method: 'GET',
      auth: true,
    });
  }

  async getAllGoogleSheetsData() {
    return this.request('/google-sheets/all-data', {
      method: 'GET',
      auth: true,
    });
  }

  // PEM Operation Plan endpoints
  async createPEMOperationPlan(planData) {
    return this.request('/pem-operation-plans', {
      method: 'POST',
      auth: true,
      body: JSON.stringify(planData),
    });
  }

  async getAllPEMOperationPlans(filters = {}) {
    const params = new URLSearchParams(filters).toString();
    const endpoint = params ? `/pem-operation-plans?${params}` : '/pem-operation-plans';
    return this.request(endpoint, {
      method: 'GET',
      auth: true,
    });
  }

  async getPEMOperationPlan(planId) {
    return this.request(`/pem-operation-plans/${planId}`, {
      method: 'GET',
      auth: true,
    });
  }

  async updatePEMOperationPlan(planId, updates) {
    return this.request(`/pem-operation-plans/${planId}`, {
      method: 'PUT',
      auth: true,
      body: JSON.stringify(updates),
    });
  }

  async deletePEMOperationPlan(planId) {
    return this.request(`/pem-operation-plans/${planId}`, {
      method: 'DELETE',
      auth: true,
    });
  }

  async addOperationPlanStep(planId, stepData) {
    return this.request(`/pem-operation-plans/${planId}/steps`, {
      method: 'POST',
      auth: true,
      body: JSON.stringify(stepData),
    });
  }

  async updateOperationPlanStep(stepId, updates) {
    return this.request(`/pem-operation-plans/steps/${stepId}`, {
      method: 'PUT',
      auth: true,
      body: JSON.stringify(updates),
    });
  }

  async deleteOperationPlanStep(stepId) {
    return this.request(`/pem-operation-plans/steps/${stepId}`, {
      method: 'DELETE',
      auth: true,
    });
  }

  async uploadStepImage(stepId, formData) {
    const url = `${this.baseURL}/pem-operation-plans/steps/${stepId}/image`;
    const token = this.getToken();

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        // Don't set Content-Type, browser will set it with boundary for multipart
      },
      body: formData,
    });

    const data = await response.json();
    if (!response.ok) {
      throw new Error(data.error || 'Upload failed');
    }
    return data;
  }

  async deleteStepImage(stepId) {
    return this.request(`/pem-operation-plans/steps/${stepId}/image`, {
      method: 'DELETE',
      auth: true,
    });
  }

  async assignPEMApprovers(planId, approvers) {
    return this.request(`/pem-operation-plans/${planId}/assign-approvers`, {
      method: 'POST',
      auth: true,
      body: JSON.stringify({ approvers }),
    });
  }

  async submitPEMOperationPlan(planId) {
    return this.request(`/pem-operation-plans/${planId}/submit`, {
      method: 'POST',
      auth: true,
    });
  }

  async approvePEMOperationPlan(planId, role, comments = '') {
    return this.request(`/pem-operation-plans/${planId}/approve?role=${encodeURIComponent(role)}`, {
      method: 'POST',
      auth: true,
      body: JSON.stringify({ comments }),
    });
  }

  async rejectPEMOperationPlan(planId, role, comments = '') {
    return this.request(`/pem-operation-plans/${planId}/reject?role=${encodeURIComponent(role)}`, {
      method: 'POST',
      auth: true,
      body: JSON.stringify({ comments }),
    });
  }

  async getPEMPlansByPPICSchedule(scheduleId) {
    return this.request(`/pem-operation-plans/ppic-schedule/${scheduleId}`, {
      method: 'GET',
      auth: true,
    });
  }

  async getPendingPEMApprovals() {
    return this.request('/pem-operation-plans/pending-approvals', {
      method: 'GET',
      auth: true,
    });
  }

  // Toolpather File Upload endpoints
  async uploadToolpatherFiles(formData) {
    const url = `${this.baseURL}/toolpather-files/upload`;
    const token = this.getToken();

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        // Don't set Content-Type, browser will set it with boundary for multipart
      },
      body: formData,
    });

    const data = await response.json();
    if (!response.ok) {
      throw new Error(data.error || 'Upload failed');
    }
    return data;
  }

  async getAllToolpatherFiles(filters = {}) {
    const params = new URLSearchParams(filters).toString();
    const endpoint = params ? `/toolpather-files?${params}` : '/toolpather-files';
    return this.request(endpoint, {
      method: 'GET',
      auth: true,
    });
  }

  async getMyToolpatherFiles() {
    return this.request('/toolpather-files/my-files', {
      method: 'GET',
      auth: true,
    });
  }

  async getToolpatherFile(id) {
    return this.request(`/toolpather-files/${id}`, {
      method: 'GET',
      auth: true,
    });
  }

  async getFilesByOrderNumber(orderNumber) {
    return this.request(`/toolpather-files/order/${encodeURIComponent(orderNumber)}`, {
      method: 'GET',
      auth: true,
    });
  }

  getToolpatherFileDownloadUrl(id) {
    return `${this.baseURL}/toolpather-files/${id}/download`;
  }

  async deleteToolpatherFile(id) {
    return this.request(`/toolpather-files/${id}`, {
      method: 'DELETE',
      auth: true,
    });
  }
}

export default new ApiService();
