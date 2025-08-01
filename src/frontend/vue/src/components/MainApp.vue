<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import DemoLoginForm from './DemoLoginForm.vue'

const authStore = useAuthStore()
const router = useRouter()

onMounted(async () => {
  try {
    await authStore.initializeAuth()
    // If authenticated after initialization, redirect to dashboard
    if (authStore.state.isAuthenticated) {
      router.push('/dashboard')
    }
  } catch (error) {
    console.error('Failed to initialize auth:', error)
    // Continue without auth for demo mode
  }
})

// Watch for authentication changes and redirect accordingly
watch(() => authStore.state.isAuthenticated, (isAuthenticated) => {
  if (isAuthenticated) {
    router.push('/dashboard')
  }
})
</script>

<template>
  <div class="main-app">
    <DemoLoginForm />
  </div>
</template>

<style scoped>
.main-app {
  min-height: 100vh;
}
</style>