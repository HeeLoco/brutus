import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import { CSRFService } from './services/csrf'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// Initialize CSRF protection
CSRFService.initialize().catch(console.error)

app.mount('#app')
