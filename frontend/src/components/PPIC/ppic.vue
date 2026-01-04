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
import api from '../../services/api.js';
import 'dhtmlx-gantt/codebase/dhtmlxgantt.css';
import './ppic.css';

const ganttContainer = ref(null);
const router = useRouter();
const machines = ref([]);

async function fetchGanttData() {
  try {
    const response = await api.getGanttChart();

    // Backend returns: { success: true, data: { sections: [...], machines: [...], summary: {...} } }
    if (response.success && response.data) {
      return response.data;
    }

    // Fallback for other response shapes
    if (response.tasks || response.links || response.sections) return response;
    if (response.data && (response.data.tasks || response.data.links || response.data.sections)) return response.data;

    return response;
  } catch (err) {
    return null;
  }
}

function transformTasksData(tasks) {
  if (!tasks) return tasks;
  return tasks.map(task => {
    // Extract machine name from machines array
    let machineName = '-';
    if (task.machines && task.machines.length > 0) {
      const machineNames = task.machines.map(m => m.machine_name || m.machine_code).filter(Boolean);
      machineName = machineNames.join(', ');
    }

    // Extract schedule ID from task_id (format: "task-5" -> 5)
    let scheduleId = null;
    const taskId = task.task_id || task.id;
    if (taskId && taskId.toString().startsWith('task-')) {
      scheduleId = parseInt(taskId.toString().replace('task-', ''));
    }

    // Map backend GanttTask fields to dhtmlx-gantt format
    const transformedTask = {
      id: taskId,                                            // task_id -> id
      ppic_schedule_id: scheduleId,                          // Store for updates
      text: task.part_name || task.task_name || task.text,   // part_name/task_name -> text
      order_number: task.njo || task.order_number,           // njo -> order_number
      start_date: task.start || task.start_date,             // start -> start_date
      end_date: task.end || task.end_date,                   // end -> end_date
      priority: task.priority,
      material: task.material_status || task.material,       // material_status -> material
      machine: machineName,
      ppic_notes: task.ppic_notes,
      status: task.status,
      progress: task.progress || 0,
      color: task.color
    };

    return transformedTask;
  });
}

async function reloadGanttData() {
  try {
    const remote = await fetchGanttData();

    if (remote) {
      gantt.clearAll();

      let tasksToLoad = [];
      if (remote.sections && Array.isArray(remote.sections)) {
        // Backend returns sections with tasks
        remote.sections.forEach(section => {
          if (section.tasks && Array.isArray(section.tasks)) {
            tasksToLoad = tasksToLoad.concat(transformTasksData(section.tasks));
          }
        });
      } else if (remote.tasks) {
        if (remote.tasks.data) {
          tasksToLoad = transformTasksData(remote.tasks.data);
        } else if (Array.isArray(remote.tasks)) {
          tasksToLoad = transformTasksData(remote.tasks);
        }
      } else if (Array.isArray(remote)) {
        tasksToLoad = transformTasksData(remote);
      }
      gantt.parse({ data: tasksToLoad, links: remote.links?.data || remote.links || [] });
    } 
  } catch (err) {
  }
}

async function fetchMachines() {
  try {
    const response = await api.getAllMachines();

    // Handle different response structures
    const machinesData = response?.machines || response?.data?.machines || response?.data || [];

    if (machinesData && machinesData.length > 0) {
      // Store full machine data for later lookup
      window.machinesData = machinesData;
      machines.value = [
        { key: "", label: "No Machine", machine_id: null },
        ...machinesData.map(m => ({
          key: `${m.machine_name} (${m.machine_code})`,
          label: `${m.machine_name} (${m.machine_code})`,
          machine_id: m.id
        }))
      ];


      // Update gantt lightbox sections after fetching machines
      if (gantt && gantt.config && gantt.config.lightbox) {
        const machineSection = gantt.config.lightbox.sections.find(s => s.name === 'Mesin');
        if (machineSection) {
          machineSection.options = machines.value;
        }
      }
    }
  } catch (err) {
    window.machinesData = [];
    machines.value = [
      { key: "", label: "No Machine", machine_id: null },
      { key: "Machine 1", label: "Machine 1", machine_id: 1 },
      { key: "Machine 2", label: "Machine 2", machine_id: 2 },
      { key: "Machine 3", label: "Machine 3", machine_id: 3 }
    ];
  }
}

onMounted(async () => {
  // ========================================
  // GRID COLUMNS CONFIGURATION
  // ========================================
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
    { name: "machine", label: "Mesin", align: "center", width: 120,
      template: (obj) => {
        return obj.machine || "-";
      }
    },
    { name: "priority", label: "Priority", align: "center", width: 100,
      template: (obj) => {
        const colors = { "Top Urgent": "red", "Urgent": "orange", "Medium": "blue", "Low": "gray" };
        return `<span style="color:${colors[obj.priority]}">${obj.priority}</span>`;
      }
    },
    { name: "add", label: "", width: 44 },
    { name: "delete", label: "", width: 44, template: (obj) => {
        return '<div class="gantt_delete"><i class="fa fa-trash"></i></div>';
      }
    }
  ];

  // ========================================
  // DATE FORMAT
  // ========================================
  gantt.config.date_format = "%Y-%m-%d";

  // ========================================
  // TIMELINE WIDTH CONFIGURATION - PERBAIKAN UTAMA
  // ========================================
  gantt.config.column_width = 10;       // Width per day column (compact view)
  gantt.config.min_column_width = 5;   // Minimum width per column
  gantt.config.autosize = false;        // Jangan auto-resize
  gantt.config.fit_tasks = false;       // Jangan fit ke container

  // ========================================
  // SCALE HEIGHT CONFIGURATION
  // ========================================
  gantt.config.scale_height = 50;       // Total height untuk 2 baris scale (month + day)

  // ========================================
  // ROW & BAR CONFIGURATION
  // ========================================
  gantt.config.row_height = 36;         // Tinggi setiap row
  gantt.config.bar_height = 36;         // Tinggi task bar
  gantt.config.show_task_cells = true;  // Tampilkan grid cells
  gantt.config.static_background = true; // Better performance

  // ========================================
  // SCALES CONFIGURATION - PERBAIKAN UTAMA
  // Ini yang memperbaiki masalah bulan terpotong
  // ========================================
  gantt.config.scales = [
    { 
      unit: "month", 
      step: 1, 
      format: "%F %Y",  // Full month name: "December 2025"
      // Alternatif format:
      // "%M %Y" = "Dec 2025" (abbreviated)
      // "%F" = "December" (month only)
      css: function(date) {
        return "scale_month";
      }
    },
    { 
      unit: "day", 
      step: 1, 
      format: "%d",     // Just day number: "28", "29", "30", "31"
      css: function(date) {
        // Highlight weekend
        if (date.getDay() === 0 || date.getDay() === 6) {
          return "scale_day weekend";
        }
        return "scale_day";
      }
    }
  ];

  // ========================================
  // TEMPLATE UNTUK SCALE STYLING
  // ========================================
  gantt.templates.scale_cell_class = function(date, scale) {
    if (scale.unit === "month") {
      return "month_scale_cell";
    }
    if (scale.unit === "day") {
      // Weekend highlighting
      if (date.getDay() === 0 || date.getDay() === 6) {
        return "day_scale_cell weekend_cell";
      }
      return "day_scale_cell";
    }
    return "";
  };

  // Template untuk task cell (background grid) - weekend highlighting
  gantt.templates.timeline_cell_class = function(_item, date) {
    if (date.getDay() === 0 || date.getDay() === 6) {
      return "weekend_cell";
    }
    return "";
  };

  // ========================================
  // CUSTOM LIGHTBOX SECTION FOR ORDER NUMBER
  // ========================================
  gantt.form_blocks["order_input"] = {
    render: function(sns) {
      return "<div class='gantt_cal_ltext'><input type='text' name='" + sns.name + "' style='width: 100%; padding: 8px; border: 1px solid #ccc; border-radius: 4px;'></div>";
    },
    set_value: function(node, value, task, section) {
      node.querySelector("input").value = value || "";
    },
    get_value: function(node, task, section) {
      return node.querySelector("input").value;
    },
    focus: function(node) {
      const input = node.querySelector("input");
      input.focus();
    }
  };

  // ========================================
  // LIGHTBOX CONFIGURATION
  // ========================================
  gantt.config.lightbox.sections = [
    { name: "Order Number", height: 38, map_to: "order_number", type: "order_input", focus: true },
    { name: "Part Name", height: 38, map_to: "text", type: "template", template: (obj) => {
        return `<div class="gantt_cal_ltext" style="padding: 8px; background: #f5f5f5; border-radius: 4px;">${obj.text || 'New task'}</div>`;
    }},
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
    { name: "Mesin", height: 38, map_to: "machine", type: "select", options: machines.value }
  ];

  // ========================================
  // EVENT: BEFORE LIGHTBOX CLOSE - VALIDATE PART NAME
  // ========================================
  gantt.attachEvent("onBeforeLightbox", (taskId) => {
    // Store the original task text
    const task = gantt.getTask(taskId);
    if (!task._originalText) {
      task._originalText = task.text;
    }
    return true;
  });

  // Prevent saving if Part Name is still "New task"
  gantt.attachEvent("onLightboxSave", (id, task, isNew) => {
    if (task.text === "New task" || !task.text || task.text.trim() === "") {
      alert("Please enter a valid Order Number and wait for the Part Name to be fetched from Google Sheets.");
      return false; // Prevent save
    }
    return true; // Allow save
  });

  // ========================================
  // INITIALIZE GANTT
  // ========================================
  gantt.init(ganttContainer.value);

  // Force remove overflow hidden from gantt elements after init
  setTimeout(() => {
    const elementsToFix = [
      '.gantt_data_area',
      '.gantt_task_area',
      '.gantt_bars_area',
      '.gantt_task',
      '.gantt_layout',
      '.gantt_layout_content',
      '.gantt_task_bg',
      '.gantt-container'
    ];

    elementsToFix.forEach(selector => {
      const elements = document.querySelectorAll(selector);
      elements.forEach(el => {
        el.style.overflow = 'visible';
        el.style.overflowX = 'visible';
        el.style.overflowY = 'visible';
      });
    });
  }, 100);

  // ========================================
  // EVENT: LIGHTBOX CHANGE - AUTO-POPULATE PART NAME FROM BACKEND
  // ========================================
  gantt.attachEvent("onLightbox", (taskId) => {
    // Add event listener to Order Number input
    setTimeout(() => {
      const orderInput = document.querySelector('input[name="Order Number"]');
      if (orderInput) {
        // Function to fetch part name from backend
        const fetchPartName = async () => {
          const orderNumber = orderInput.value.trim();

          if (orderNumber) {
            try {
              // Call backend API to get Part Name
              const response = await api.getPartNameByOrderNumber(orderNumber);

              if (response && response.success && response.data && response.data.part_name) {
                // Update the task object
                const task = gantt.getTask(taskId);
                task.text = response.data.part_name;

                // Update the Part Name display in the lightbox
                // Find all divs with class gantt_cal_ltext, the second one is Part Name (first is Order Number input wrapper)
                const allDivs = document.querySelectorAll('.gantt_cal_ltext');
                if (allDivs.length >= 2) {
                  // allDivs[0] is Order Number input wrapper
                  // allDivs[1] is Part Name display
                  allDivs[1].textContent = response.data.part_name;
                }
              }
            } catch (err) {
              console.error('Error fetching Part Name:', err);
              // Optionally show error to user
              // alert('Failed to fetch Part Name for Order Number: ' + orderNumber);
            }
          }
        };

        // Prevent Enter key from closing/submitting the lightbox
        orderInput.addEventListener('keydown', (e) => {
          if (e.key === 'Enter') {
            e.preventDefault();
            e.stopPropagation();
            // Trigger fetch when Enter is pressed
            fetchPartName();
            return false;
          }
        });

        // Add blur event to fetch Part Name from backend
        orderInput.addEventListener('blur', fetchPartName);
      }
    }, 100);
  });

  // Fetch machines AFTER gantt init, before loading data
  await fetchMachines();

  const remote = await fetchGanttData();
  if (remote) {
    let tasksToLoad = [];
    if (remote.sections && Array.isArray(remote.sections)) {
      // Backend returns sections with tasks
      remote.sections.forEach(section => {
        if (section.tasks && Array.isArray(section.tasks)) {
          tasksToLoad = tasksToLoad.concat(transformTasksData(section.tasks));
        }
      });
    } else if (remote.tasks) {
      if (remote.tasks.data) {
        tasksToLoad = transformTasksData(remote.tasks.data);
      } else if (Array.isArray(remote.tasks)) {
        tasksToLoad = transformTasksData(remote.tasks);
      }
    } else if (Array.isArray(remote)) {
      tasksToLoad = transformTasksData(remote);
    }

    gantt.parse({ data: tasksToLoad, links: remote.links?.data || remote.links || [] });
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

  // ========================================
  // EVENT: AFTER TASK ADD
  // ========================================
  gantt.attachEvent("onAfterTaskAdd", async (id, item) => {
    const fmt = (d) => {
      if (!d) return null;
      if (typeof d === 'string') return d.slice(0,10);
      return d.toISOString().slice(0,10);
    };

    const parseMachineId = (machineName) => {
      if (!machineName || machineName === "" || machineName === "No Machine") return null;
      // Find machine by name from machines array
      const machine = machines.value.find(m => m.key === machineName || m.label === machineName);
      if (machine && machine.machine_id) {
        return machine.machine_id;
      }
      // Fallback: try to find from window.machinesData
      if (window.machinesData) {
        const machineData = window.machinesData.find(m =>
          `${m.machine_name} (${m.machine_code})` === machineName ||
          m.machine_name === machineName ||
          m.machine_code === machineName
        );
        if (machineData) return machineData.id;
      }
      return null; // Return null if no machine found
    };

    const machineId = parseMachineId(item.machine);
    const payload = {
      njo: item.order_number || '',
      part_name: item.text || '',
      start_date: fmt(item.start_date) || fmt(new Date()),
      finish_date: fmt(item.end_date || (item.start_date && new Date(item.start_date.getTime() + (item.duration||0)*24*60*60*1000))) || fmt(new Date()),
      priority: item.priority || 'Low',
      priority_alpha: '',
      material_status: item.material || 'Ready',
      ppic_notes: item.ppic_notes || '',
      machine_assignments: machineId ? [
        {
          machine_id: machineId,
          target_hours: (item.duration || 1) * 8,
          sequence: 1
        }
      ] : []
    };


    try {
      const response = await api.createPPICSchedule(payload);

      if (response && response.data && response.data.id) {
        // Reload data from database to show the saved schedule
        await reloadGanttData();
        return true;
      }
      return true;
    } catch (err) {

      // Extract error message from response
      let errorMessage = 'Failed to create PPIC schedule';
      if (err.response && err.response.data) {
        errorMessage = err.response.data.error || err.response.data.message || errorMessage;
      } else if (err.message) {
        errorMessage = err.message;
      }

      // Show error to user
      if (errorMessage.includes('NJO already exists') || errorMessage.includes('already exists')) {
        alert(`Error: Order Number (NJO) "${payload.njo}" sudah ada di database. Gunakan nomor order yang berbeda.`);
      } else {
        alert(`Error: ${errorMessage}`);
      }

      // Delete the task from gantt since it failed to save
      gantt.deleteTask(id);
      return false;
    }
  });

  // ========================================
  // EVENT: AFTER TASK UPDATE
  // ========================================
  gantt.attachEvent("onAfterTaskUpdate", async (id, item) => {
    const fmt = (d) => {
      if (!d) return null;
      if (typeof d === 'string') return d.slice(0,10);
      return d.toISOString().slice(0,10);
    };

    const parseMachineId = (machineName) => {
      if (!machineName || machineName === "" || machineName === "No Machine" || machineName === "-") return null;
      // Find machine by name from machines array
      const machine = machines.value.find(m => m.key === machineName || m.label === machineName);
      if (machine && machine.machine_id) {
        return machine.machine_id;
      }
      // Fallback: try to find from window.machinesData
      if (window.machinesData) {
        const machineData = window.machinesData.find(m =>
          `${m.machine_name} (${m.machine_code})` === machineName ||
          m.machine_name === machineName ||
          m.machine_code === machineName
        );
        if (machineData) return machineData.id;
      }
      return null; // Return null if no machine found
    };

    // Skip if this is a new task without ppic_schedule_id (will be handled by onAfterTaskAdd)
    if (!item.ppic_schedule_id && !item.id.toString().startsWith('task-')) {
      console.log('Skipping update for new task:', id);
      return true;
    }

    const machineId = parseMachineId(item.machine);
    const payload = {
      part_name: item.text || '',
      priority: item.priority || 'Low',
      priority_alpha: '',
      material_status: item.material || 'Ready',
      start_date: fmt(item.start_date) || fmt(new Date()),
      finish_date: fmt(item.end_date || (item.start_date && new Date(item.start_date.getTime() + (item.duration||0)*24*60*60*1000))) || fmt(new Date()),
      ppic_notes: item.ppic_notes || '',
      machine_assignments: machineId ? [
        {
          machine_id: machineId,
          target_hours: (item.duration || 1) * 8,
          sequence: 1
        }
      ] : []
    };

    console.log('Updating PPIC schedule:', id, JSON.stringify(payload, null, 2));

    try {
      // Extract numeric ID from task_id format (e.g., "task-5" -> 5)
      let scheduleId = item.ppic_schedule_id;
      if (!scheduleId && item.id.toString().startsWith('task-')) {
        scheduleId = parseInt(item.id.toString().replace('task-', ''));
      }

      if (!scheduleId) {
        console.error('Cannot update: No schedule ID found');
        return false;
      }

      const response = await api.updatePPICSchedule(scheduleId, payload);

      if (response && response.success) {
        console.log('PPIC schedule updated successfully:', scheduleId);
        // Reload data from database to show the updated schedule and any cascaded changes
        await reloadGanttData();
        return true;
      }
      return true;
    } catch (err) {
      console.error('Error updating PPIC schedule:', err);
      console.error('Failed payload was:', JSON.stringify(payload, null, 2));

      // Extract error message from response
      let errorMessage = 'Failed to update PPIC schedule';
      if (err.response && err.response.data) {
        errorMessage = err.response.data.error || err.response.data.message || errorMessage;
      } else if (err.message) {
        errorMessage = err.message;
      }

      console.error('Backend error:', errorMessage);
      alert(`Error updating schedule: ${errorMessage}`);

      // Reload to revert changes
      await reloadGanttData();
      return false;
    }
  });

  // ========================================
  // EVENT: AFTER LINK ADD
  // ========================================
  gantt.attachEvent("onAfterLinkAdd", async (id, link) => {
    try {
      // Extract schedule IDs from task IDs
      const getScheduleId = (taskId) => {
        if (taskId && taskId.toString().startsWith('task-')) {
          return parseInt(taskId.toString().replace('task-', ''));
        }
        return null;
      };

      const sourceScheduleId = getScheduleId(link.source);
      const targetScheduleId = getScheduleId(link.target);

      if (!sourceScheduleId || !targetScheduleId) {
        console.error('Cannot create link: Invalid task IDs');
        alert('Error: Invalid task IDs for linking');
        gantt.deleteLink(id);
        return false;
      }

      const payload = {
        source_schedule_id: sourceScheduleId,
        target_schedule_id: targetScheduleId,
        link_type: link.type || '0'
      };

      console.log('Creating PPIC link:', JSON.stringify(payload, null, 2));

      const response = await api.createPPICLink(payload);

      if (response && response.data && response.data.id) {
        // Store the backend link ID on the gantt link object
        const ganttLink = gantt.getLink(id);
        ganttLink.ppic_link_id = response.data.id;
        console.log('PPIC link created successfully:', response.data.id);

        // Reload gantt data to show any auto-rescheduled tasks
        await reloadGanttData();
      }

      return true;
    } catch (err) {
      console.error('Error creating PPIC link:', err);

      // Extract error message from response
      let errorMessage = 'Failed to create link';
      if (err.response && err.response.data) {
        errorMessage = err.response.data.error || err.response.data.message || errorMessage;
      } else if (err.message) {
        errorMessage = err.message;
      }

      alert(`Error: ${errorMessage}`);

      // Delete the link from gantt since it failed to save
      gantt.deleteLink(id);
      return false;
    }
  });

  // ========================================
  // EVENT: BEFORE LINK DELETE
  // ========================================
  gantt.attachEvent("onBeforeLinkDelete", async (id, link) => {
    try {
      // Get the backend link ID
      let linkId = link.ppic_link_id;

      if (linkId) {
        console.log('Deleting PPIC link:', linkId);
        await api.deletePPICLink(linkId);
        console.log('PPIC link deleted successfully');
      }

      return true;
    } catch (err) {
      console.error('Error deleting PPIC link:', err);
      return true; // Allow deletion in gantt even if backend fails
    }
  });

  // ========================================
  // EVENT: BEFORE TASK DELETE
  // ========================================
  gantt.attachEvent("onBeforeTaskDelete", async (id, item) => {
    try {
      // Skip if this is a new task without ppic_schedule_id
      if (!item.ppic_schedule_id && !item.id.toString().startsWith('task-')) {
        console.log('Skipping delete for new task:', id);
        return true;
      }

      // Show confirmation dialog
      const confirmed = confirm(`Apakah Anda yakin ingin menghapus schedule "${item.text}"?`);
      if (!confirmed) {
        return false; // Cancel deletion
      }

      // Extract numeric ID from task_id format (e.g., "task-5" -> 5)
      let scheduleId = item.ppic_schedule_id;
      if (!scheduleId && item.id.toString().startsWith('task-')) {
        scheduleId = parseInt(item.id.toString().replace('task-', ''));
      }

      if (!scheduleId) {
        console.error('Cannot delete: No schedule ID found');
        return false;
      }

      console.log('Deleting PPIC schedule:', scheduleId);
      const response = await api.deletePPICSchedule(scheduleId);

      if (response && response.success) {
        console.log('PPIC schedule deleted successfully:', scheduleId);
        return true; // Allow deletion in gantt
      }

      return true;
    } catch (err) {
      console.error('Error deleting PPIC schedule:', err);

      // Extract error message from response
      let errorMessage = 'Failed to delete PPIC schedule';
      if (err.response && err.response.data) {
        errorMessage = err.response.data.error || err.response.data.message || errorMessage;
      } else if (err.message) {
        errorMessage = err.message;
      }

      alert(`Error deleting schedule: ${errorMessage}`);
      return false; // Prevent deletion in gantt if backend fails
    }
  });

  // ========================================
  // EVENT: ON GANTT RENDER (Delete button handler)
  // ========================================
  gantt.attachEvent("onGanttRender", () => {
    const deleteButtons = document.querySelectorAll('.gantt_delete');
    deleteButtons.forEach(button => {
      button.onclick = (e) => {
        e.stopPropagation();
        const row = e.target.closest('.gantt_row');
        if (row) {
          const taskId = row.getAttribute('task_id');
          if (taskId) {
            const task = gantt.getTask(taskId);
            if (task) {
              gantt.deleteTask(taskId);
            }
          }
        }
      };
    });
  });
});

const goToDashboard = () => {
  router.push('/dashboard');
};
</script>