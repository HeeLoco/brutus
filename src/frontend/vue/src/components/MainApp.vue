<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import DemoLoginForm from './DemoLoginForm.vue'
import ResourceGroupList from './ResourceGroupList.vue'

const authStore = useAuthStore()

onMounted(async () => {
  try {
    await authStore.initializeAuth()
  } catch (error) {
    console.error('Failed to initialize auth:', error)
    // Continue without auth for demo mode
  }
})
</script>

<template>
  <div class="main-app">
    <DemoLoginForm v-if="!authStore.state.isAuthenticated" />
    <ResourceGroupList v-else />
  </div>
</template>

<style scoped>
.main-app {
  min-height: 100vh;
}
</style>