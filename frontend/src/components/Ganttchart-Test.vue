<template>
  <div class="gantt-container">
    <div ref="ganttContainer"></div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { gantt } from 'dhtmlx-gantt'; 
import 'dhtmlx-gantt/codebase/dhtmlxgantt.css'; 

const ganttContainer = ref(null);

const data = {
  tasks: [
    { 
      id: 1, 
      text: "Project Planning", 
      start_date: "2025-11-01", 
      duration: 5, 
      progress: 0.6, 
      open: true 
    },
    { 
      id: 2, 
      text: "Research", 
      start_date: "2025-11-02", 
      duration: 3, 
      parent: 1, 
      progress: 0.4
    },
    { 
      id: 3, 
      text: "Design Mockups", 
      start_date: "2025-11-05", 
      duration: 2, 
      parent: 1, 
      progress: 0.2
    },
    { 
      id: 4, 
      text: "Development", 
      start_date: "2025-11-07", 
      duration: 4, 
      progress: 0.0
    }
  ],
  links: [
    { id: 1, source: 3, target: 4, type: "0" } 
  ]
};

onMounted(() => {
  gantt.config.date_format = "%Y-%m-%d";

  gantt.init(ganttContainer.value);

  gantt.parse(data);
});

</script>

<style scoped>
.gantt-container {
  width: 100%;
  height: 500px;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding: 20px 0;
  box-sizing: border-box;
}

.gantt-container > div {
  width: 900px; 
  max-width: 95%; 
  height: 100%;
  box-shadow: 0 8px 24px rgba(0,0,0,0.08);
  border-radius: 8px;
  background: #fff;
  overflow: hidden;
}
</style>