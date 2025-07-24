<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth'
import DemoLoginForm from './components/DemoLoginForm.vue'
import ResourceGroupList from './components/ResourceGroupList.vue'

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
  <div id="app">
    <DemoLoginForm v-if="!authStore.state.isAuthenticated" />
    <ResourceGroupList v-else />
  </div>
</template>

<style>
* {
  box-sizing: border-box;
}

html, body {
  margin: 0;
  padding: 0;
  min-height: 100vh;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background-color: #f8fafc;
}

#app {
  min-height: 100vh;
}
</style>
