import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../components/MainApp.vue'),
    },
    {
      path: '/debug',
      name: 'debug',
      component: () => import('../components/EnvironmentDebug.vue'),
    },
  ],
})

export default router
