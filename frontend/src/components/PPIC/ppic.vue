<template>
  <div class="ppic-wrapper">
    <header class="main-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="logo">IMETRAX</h1>
          <span class="divider">|</span>
          <h2 class="page-title">PPIC - Production Planning</h2>
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

    <main class="ppic-body">
      <div class="gantt-container">
        <div ref="ganttContainer"></div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { gantt } from 'dhtmlx-gantt'; 
import { useRouter } from 'vue-router';
import 'dhtmlx-gantt/codebase/dhtmlxgantt.css'; 
import './ppic.css'; 

const ganttContainer = ref(null);
const router = useRouter();

async function fetchGanttData() {
  try {
    const token = localStorage.getItem('token');
    const headers = token ? { 'Authorization': `Bearer ${token}` } : {};
    const res = await fetch('/api/v1/gantt-chart', { headers });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = await res.json();
    // Support common response shapes: { tasks, links } or { data: { tasks, links } }
    if (data.tasks || data.links) return data;
    if (data.data && (data.data.tasks || data.data.links)) return data.data;
    return data;
  } catch (err) {
    console.error('Failed fetching gantt data', err);
    return null;
  }
}

onMounted(async () => {
  gantt.config.columns = [
    { name: "text", label: "Part Name", tree: true, width: 150 }, 
    { name: "order_number", label: "Order #", align: "center", width: 80 }, 
    { name: "start_date", label: "Start Date", align: "center", width: 100, 
      template: (obj) => {
        return gantt.templates.date_grid(obj.start_date);
      }
    },
    { name: "finish_date", label: "Finish Date", align: "center", width: 100, 
      template: (obj) => {
        return gantt.templates.date_grid(obj.end_date);
      }
    },
    { name: "priority", label: "Priority", align: "center", width: 100, 
      template: (obj) => {
        const colors = { "Top Urgent": "red", "Urgent": "orange", "Medium": "blue", "Low": "gray" };
        return `<span style="color:${colors[obj.priority]}">${obj.priority}</span>`;
      }
    },
    { name: "add", label: "", width: 44 }
  ];

  gantt.config.date_format = "%Y-%m-%d";
  
  gantt.config.lightbox.sections = [
    { name: "Order Number", height: 38, map_to: "order_number", type: "textarea", focus: true },
    { name: "Part Name", height: 38, map_to: "text", type: "textarea" },
    { name: "Time Period", type: "time", map_to: "auto", time_format: ["%d", "%m", "%Y"] },
    { name: "Priority", height: 38, map_to: "priority", type: "select", options: [
        { key: "Low", label: "Low" },
        { key: "Medium", label: "Medium" },
        { key: "Urgent", label: "Urgent" },
        { key: "Top Urgent", label: "Top Urgent" }
    ]},
    { name: "Material Status", height: 38, map_to: "material", type: "select", options: [
        { key: "Ready", label: "Ready" },
        { key: "Not Ready", label: "Not Ready" }
    ]},
    { name: "Mesin", height: 38, map_to: "machine", type: "select", options: [
        { key: "Machine 1", label: "Machine 1" },
        { key: "Machine 2", label: "Machine 2" },
        { key: "Machine 3", label: "Machine 3" }
    ]}
  ];

  gantt.init(ganttContainer.value);

  const remote = await fetchGanttData();
  if (remote) {
    if (remote.tasks && remote.tasks.data) {
      gantt.parse({ tasks: remote.tasks.data, links: (remote.links && remote.links.data) || remote.links || [] });
    } else if (remote.tasks) {
      gantt.parse({ tasks: remote.tasks, links: remote.links || [] });
    } else if (Array.isArray(remote)) {
      gantt.parse({ tasks: remote, links: [] });
    } else {
      gantt.parse(remote);
    }
  } else {
    // Keep local sample if fetch failed
    // gantt.parse({
    //   tasks: [
    //     { id: 10, text: "Section 1", open: true, type: "project" },
    //     { id: 1, text: "Part A-101", order_number: "5252", priority: "Top Urgent", material: "Ready", start_date: "2025-12-17", duration: 10, parent: 10 },
    //     { id: 2, text: "Part B-202", order_number: "5253", priority: "Medium", material: "Ready", start_date: "2025-12-18", duration: 5, parent: 10 }
    //   ],
    //   links: []
    // });
  }

  gantt.attachEvent("onAfterTaskAdd", async (id, item) => {
    try {
      const token = localStorage.getItem('token');
      const headers = { 'Content-Type': 'application/json' };
      if (token) headers['Authorization'] = `Bearer ${token}`;

      const fmt = (d) => {
        if (!d) return null;
        if (typeof d === 'string') return d.slice(0,10);
        return d.toISOString().slice(0,10);
      };

      const parseMachineId = (m) => {
        if (!m) return 1;
        const parts = String(m).trim().split(/\s+/);
        const last = parts[parts.length-1];
        const n = parseInt(last, 10);
        return isFinite(n) && n > 0 ? n : 1;
      };

      const payload = {
        njo: item.order_number || '',
        part_name: item.text || '',
        start_date: fmt(item.start_date) || fmt(new Date()),
        finish_date: fmt(item.end_date || (item.start_date && new Date(item.start_date.getTime() + (item.duration||0)*24*60*60*1000))) || fmt(new Date()),
        priority: item.priority || 'Low',
        priority_alpha: '',
        material_status: item.material || 'Ready',
        ppic_notes: item.ppic_notes || '',
        machine_assignments: [
          {
            machine_id: parseMachineId(item.machine),
            target_hours: (item.duration || 1) * 8,
            sequence: 1
          }
        ]
      };

      const res = await fetch('/api/v1/ppic-schedules', {
        method: 'POST',
        headers,
        body: JSON.stringify(payload)
      });

      const body = await res.json().catch(() => ({}));
      if (!res.ok) {
        console.error('Failed creating PPIC schedule', res.status, body);
        return true;
      }

      if (body && body.data && body.data.id) {
        gantt.changeTask(id, { ppic_schedule_id: body.data.id });
      }
      return true;
    } catch (err) {
      console.error('Error in onAfterTaskAdd', err);
      return true;
    }
  });
});

const goToDashboard = () => {
  router.push('/dashboard');
};
</script>