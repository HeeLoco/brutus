import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'root',
      component: () => import('../components/MainApp.vue')
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../components/AppLayout.vue'),
      children: [
        {
          path: '',
          component: () => import('../components/Dashboard.vue'),
        }
      ]
    },
    {
      path: '/resource-groups',
      name: 'resource-groups',
      component: () => import('../components/AppLayout.vue'),
      children: [
        {
          path: '',
          component: () => import('../components/ResourceGroupList.vue'),
        }
      ]
    },
    {
      path: '/storage-accounts',
      name: 'storage-accounts',
      component: () => import('../components/AppLayout.vue'),
      children: [
        {
          path: '',
          component: () => import('../components/StorageAccountList.vue'),
        }
      ]
    },
    {
      path: '/settings',
      name: 'settings',
      component: () => import('../components/AppLayout.vue'),
      children: [
        {
          path: '',
          component: () => import('../components/Settings.vue'),
        }
      ]
    },
    {
      path: '/debug',
      name: 'debug',
      component: () => import('../components/AppLayout.vue'),
      children: [
        {
          path: '',
          component: () => import('../components/EnvironmentDebug.vue'),
        }
      ]
    }
  ],
})

// Route guard to check authentication
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // Routes that require authentication
  const protectedRoutes = ['/dashboard', '/resource-groups', '/storage-accounts', '/settings', '/debug']
  
  if (protectedRoutes.includes(to.path)) {
    if (!authStore.state.isAuthenticated) {
      // Redirect to login if not authenticated
      next('/')
      return
    }
  }
  
  // If going to root and already authenticated, redirect to dashboard
  if (to.path === '/' && authStore.state.isAuthenticated) {
    next('/dashboard')
    return
  }
  
  next()
})

export default router
