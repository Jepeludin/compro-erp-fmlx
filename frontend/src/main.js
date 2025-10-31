import { createApp } from 'vue'
import './style.css' // CSS global
import App from './App.vue'
import router from './router' // <-- Impor router kita

const app = createApp(App)

app.use(router) // <-- Gunakan router

app.mount('#app')