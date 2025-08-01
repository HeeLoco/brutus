<template>
  <div class="settings">
    <div class="settings-header">
      <h1>Settings</h1>
      <p class="settings-subtitle">Configure your Azure Resource Manager preferences</p>
    </div>

    <div class="settings-content">
      <!-- Current Configuration -->
      <div class="settings-card">
        <div class="card-header">
          <h3>🔧 Current Configuration</h3>
        </div>
        <div class="card-content">
          <div class="config-grid">
            <div class="config-item">
              <label>API Mode</label>
              <div class="config-value" :class="getModeClass(authStore.state.apiMode)">
                {{ getModeDisplay(authStore.state.apiMode) }}
              </div>
            </div>
            <div class="config-item">
              <label>Subscription ID</label>
              <div class="config-value subscription-id">
                {{ authStore.state.subscriptionId || 'Not Set' }}
              </div>
            </div>
            <div class="config-item">
              <label>User Account</label>
              <div class="config-value">
                {{ authStore.state.account?.name || 'Demo User' }}
              </div>
            </div>
            <div class="config-item">
              <label>Authentication Status</label>
              <div class="config-value" :class="{ 'status-active': authStore.state.isAuthenticated }">
                {{ authStore.state.isAuthenticated ? 'Authenticated' : 'Not Authenticated' }}
              </div>
            </div>
          </div>
        </div>
      </div>


      <!-- Actions -->
      <div class="settings-card">
        <div class="card-header">
          <h3>⚡ Actions</h3>
        </div>
        <div class="card-content">
          <div class="action-grid">
            <button @click="changeMode" class="action-button primary">
              <span class="btn-icon">🔄</span>
              <div class="btn-content">
                <div class="btn-title">Change Mode</div>
                <div class="btn-subtitle">Switch to a different authentication mode</div>
              </div>
            </button>

            <button @click="refreshConfiguration" class="action-button secondary" :disabled="loading">
              <span class="btn-icon">{{ loading ? '⟳' : '🔄' }}</span>
              <div class="btn-content">
                <div class="btn-title">{{ loading ? 'Refreshing...' : 'Refresh Configuration' }}</div>
                <div class="btn-subtitle">Reload current settings and status</div>
              </div>
            </button>

            <button @click="clearCache" class="action-button warning">
              <span class="btn-icon">🗑️</span>
              <div class="btn-content">
                <div class="btn-title">Clear Cache</div>
                <div class="btn-subtitle">Clear all cached data and tokens</div>
              </div>
            </button>

            <router-link to="/debug" class="action-button secondary">
              <span class="btn-icon">🔧</span>
              <div class="btn-content">
                <div class="btn-title">Debug Information</div>
                <div class="btn-subtitle">View detailed system information</div>
              </div>
            </router-link>
          </div>
        </div>
      </div>

      <!-- Environment Information -->
      <div class="settings-card">
        <div class="card-header">
          <h3>🌍 Environment Information</h3>
        </div>
        <div class="card-content">
          <div class="env-grid">
            <div class="env-item">
              <label>Frontend URL</label>
              <div class="env-value">{{ frontendUrl }}</div>
            </div>
            <div class="env-item">
              <label>Backend URL</label>
              <div class="env-value">{{ getBackendUrl() }}</div>
            </div>
            <div class="env-item">
              <label>Azure Client ID</label>
              <div class="env-value">{{ azureClientId || 'Not Set' }}</div>
            </div>
            <div class="env-item">
              <label>Azure Tenant ID</label>
              <div class="env-value">{{ azureTenantId || 'Not Set' }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import { logger } from '../services/logger'

const authStore = useAuthStore()
const router = useRouter()

const loading = ref(false)
const frontendUrl = window.location.origin

// Environment variables
const azureClientId = import.meta.env.VITE_AZURE_CLIENT_ID
const azureTenantId = import.meta.env.VITE_AZURE_TENANT_ID

const getModeDisplay = (mode: string | null): string => {
  switch (mode) {
    case 'demo': return 'Demo Mode'
    case 'backend': return 'Backend Mode'
    case 'backend-go': return 'Go Backend'
    case 'backend-python': return 'Python Backend'
    case 'backend-typescript': return 'TypeScript Backend'
    case 'azure': return 'Azure Direct'
    case 'azure-direct': return 'Azure Direct'
    default: return 'Unknown Mode'
  }
}

const getModeClass = (mode: string | null): string => {
  switch (mode) {
    case 'demo': return 'mode-demo'
    case 'backend': return 'mode-backend'
    case 'backend-go': return 'mode-backend'
    case 'backend-python': return 'mode-backend'
    case 'backend-typescript': return 'mode-backend'
    case 'azure': return 'mode-azure'
    case 'azure-direct': return 'mode-azure'
    default: return 'mode-unknown'
  }
}

const getBackendUrl = (): string => {
  const mode = authStore.state.apiMode
  if (mode === 'backend-go') {
    return 'http://localhost:8080'
  } else if (mode === 'backend-python') {
    return 'http://localhost:8000'
  } else if (mode === 'backend-typescript') {
    return 'http://localhost:3000'
  } else if (mode === 'backend') {
    // Legacy mode - check username for backward compatibility
    if (authStore.state.account?.username?.includes('python')) {
      return 'http://localhost:8000'
    } else if (authStore.state.account?.username?.includes('typescript')) {
      return 'http://localhost:3000'
    } else {
      return 'http://localhost:8080'
    }
  }
  return 'Not applicable'
}

const changeMode = async () => {
  try {
    logger.info('User requested mode change')
    await authStore.logout()
    // Force reload to ensure clean state
    window.location.href = '/'
  } catch (error) {
    logger.logError(error, 'Failed to change mode')
    // Force reload anyway
    window.location.href = '/'
  }
}

const refreshConfiguration = async () => {
  loading.value = true
  try {
    logger.info('Refreshing configuration')
    await authStore.initializeAuth()
    logger.info('Configuration refreshed successfully')
  } catch (error) {
    logger.logError(error, 'Failed to refresh configuration')
  } finally {
    loading.value = false
  }
}

const clearCache = async () => {
  if (confirm('Are you sure you want to clear all cached data? This will log you out.')) {
    try {
      logger.info('Clearing cache')
      await authStore.clearInteractionState()
      // Force reload to ensure clean state
      window.location.href = '/'
    } catch (error) {
      logger.logError(error, 'Failed to clear cache')
      // Force reload anyway
      window.location.href = '/'
    }
  }
}

</script>

<style scoped>
.settings {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.settings-header {
  margin-bottom: 2rem;
}

.settings-header h1 {
  font-size: 2.5rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.settings-subtitle {
  color: #64748b;
  font-size: 1.125rem;
  margin: 0.5rem 0 0 0;
}

.settings-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.settings-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.card-header {
  padding: 1.5rem 1.5rem 0 1.5rem;
}

.card-header h3 {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1e293b;
  margin: 0;
}

.card-content {
  padding: 1.5rem;
}

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.config-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.config-item label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
}

.config-value {
  padding: 0.75rem;
  background: #f8fafc;
  border-radius: 0.5rem;
  border: 1px solid #e2e8f0;
  font-size: 0.875rem;
  color: #1e293b;
}

.subscription-id {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.75rem;
}

.mode-demo { background: #f0f9ff; color: #0369a1; border-color: #bae6fd; }
.mode-backend { background: #f0fdf4; color: #059669; border-color: #bbf7d0; }
.mode-azure { background: #fef2f2; color: #dc2626; border-color: #fecaca; }
.mode-unknown { background: #f9fafb; color: #6b7280; border-color: #d1d5db; }

.status-active {
  background: #f0fdf4;
  color: #059669;
  border-color: #bbf7d0;
}


.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1rem;
}

.action-button {
  display: flex;
  align-items: center;
  padding: 1rem;
  border-radius: 0.5rem;
  text-decoration: none;
  transition: all 0.2s;
  border: none;
  cursor: pointer;
  width: 100%;
  text-align: left;
}

.action-button.primary {
  background: #f0f9ff;
  color: #0369a1;
  border: 1px solid #bae6fd;
}

.action-button.primary:hover {
  background: #e0f2fe;
  border-color: #7dd3fc;
}

.action-button.secondary {
  background: #f8fafc;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.action-button.secondary:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

.action-button.warning {
  background: #fefce8;
  color: #ca8a04;
  border: 1px solid #fde047;
}

.action-button.warning:hover {
  background: #fef3c7;
  border-color: #facc15;
}

.action-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-icon {
  font-size: 1.5rem;
  margin-right: 1rem;
}

.btn-content {
  flex: 1;
}

.btn-title {
  font-weight: 600;
  font-size: 0.875rem;
}

.btn-subtitle {
  font-size: 0.75rem;
  opacity: 0.8;
  margin-top: 0.25rem;
}

.env-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.env-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.env-item label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
}

.env-value {
  padding: 0.5rem;
  background: #f8fafc;
  border-radius: 0.375rem;
  border: 1px solid #e2e8f0;
  font-size: 0.75rem;
  color: #1e293b;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  word-break: break-all;
}

@media (max-width: 768px) {
  .settings {
    padding: 1rem;
  }
  
  .settings-header h1 {
    font-size: 2rem;
  }
  
  .config-grid, .action-grid, .env-grid {
    grid-template-columns: 1fr;
  }
}
</style>