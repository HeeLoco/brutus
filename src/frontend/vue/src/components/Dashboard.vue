<template>
  <div class="dashboard">
    <!-- Dashboard Header -->
    <div class="dashboard-header">
      <h1>Dashboard</h1>
      <p class="dashboard-subtitle">Azure Resource Management Overview</p>
    </div>

    <!-- Quick Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📁</div>
        <div class="stat-content">
          <div class="stat-number">{{ stats.resourceGroups }}</div>
          <div class="stat-label">Resource Groups</div>
        </div>
        <router-link to="/resource-groups" class="stat-link">View All →</router-link>
      </div>

      <div class="stat-card">
        <div class="stat-icon">💾</div>
        <div class="stat-content">
          <div class="stat-number">{{ stats.storageAccounts }}</div>
          <div class="stat-label">Storage Accounts</div>
        </div>
        <router-link to="/storage-accounts" class="stat-link">View All →</router-link>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🌐</div>
        <div class="stat-content">
          <div class="stat-number">{{ getModeDisplay(authStore.state.apiMode) }}</div>
          <div class="stat-label">API Mode</div>
        </div>
        <router-link to="/settings" class="stat-link">Configure →</router-link>
      </div>

      <div class="stat-card">
        <div class="stat-icon">✅</div>
        <div class="stat-content">
          <div class="stat-number" :class="{ 'status-healthy': isHealthy, 'status-unhealthy': !isHealthy }">
            {{ isHealthy ? 'Healthy' : 'Issues' }}
          </div>
          <div class="stat-label">System Status</div>
        </div>
        <router-link to="/debug" class="stat-link">Details →</router-link>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="content-grid">
      <!-- Quick Actions -->
      <div class="content-card">
        <div class="card-header">
          <h3>🚀 Quick Actions</h3>
        </div>
        <div class="card-content">
          <div class="action-buttons">
            <router-link to="/resource-groups" class="action-btn primary">
              <span class="btn-icon">📁</span>
              <div class="btn-content">
                <div class="btn-title">Create Resource Group</div>
                <div class="btn-subtitle">Organize your Azure resources</div>
              </div>
            </router-link>

            <router-link to="/storage-accounts" class="action-btn primary">
              <span class="btn-icon">💾</span>
              <div class="btn-content">
                <div class="btn-title">Create Storage Account</div>
                <div class="btn-subtitle">Set up blob, file, table storage</div>
              </div>
            </router-link>

            <router-link to="/settings" class="action-btn secondary">
              <span class="btn-icon">⚙️</span>
              <div class="btn-content">
                <div class="btn-title">Configure Settings</div>
                <div class="btn-subtitle">Update preferences and API mode</div>
              </div>
            </router-link>

            <button @click="refreshData" class="action-btn secondary" :disabled="loading">
              <span class="btn-icon">{{ loading ? '⟳' : '🔄' }}</span>
              <div class="btn-content">
                <div class="btn-title">{{ loading ? 'Refreshing...' : 'Refresh Data' }}</div>
                <div class="btn-subtitle">Update resource information</div>
              </div>
            </button>
          </div>
        </div>
      </div>

      <!-- System Information -->
      <div class="content-card">
        <div class="card-header">
          <h3>ℹ️ System Information</h3>
        </div>
        <div class="card-content">
          <div class="info-list">
            <div class="info-item">
              <span class="info-label">Subscription ID:</span>
              <span class="info-value subscription-id">{{ authStore.state.subscriptionId || 'Not Set' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Authentication Mode:</span>
              <span class="info-value" :class="getModeClass(authStore.state.apiMode)">
                {{ getModeDisplay(authStore.state.apiMode) }}
              </span>
            </div>
            <div class="info-item">
              <span class="info-label">User Account:</span>
              <span class="info-value">{{ authStore.state.account?.name || 'Demo User' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Last Updated:</span>
              <span class="info-value">{{ formatDate(lastUpdated) }}</span>
            </div>
          </div>
        </div>
      </div>


      <!-- Mode-Specific Information -->
      <div class="content-card">
        <div class="card-header">
          <h3>🔧 Current Mode: {{ getModeDisplay(authStore.state.apiMode) }}</h3>
        </div>
        <div class="card-content">
          <div class="mode-info">
            <div v-if="authStore.state.apiMode === 'demo'" class="mode-details">
              <p><strong>Demo Mode Active</strong></p>
              <p>You're using mock data for testing. No real Azure resources are affected.</p>
              <ul>
                <li>✅ Safe testing environment</li>
                <li>✅ All features available</li>
                <li>✅ No Azure credentials required</li>
                <li>⚠️ Changes are not persistent</li>
              </ul>
            </div>
            <div v-else-if="authStore.state.apiMode === 'backend'" class="mode-details">
              <p><strong>Backend Mode Active</strong></p>
              <p>Connected to backend API server for resource management.</p>
              <ul>
                <li>✅ Real backend integration</li>
                <li>✅ API validation</li>
                <li>✅ Server-side processing</li>
                <li>ℹ️ Requires backend server running</li>
              </ul>
            </div>
            <div v-else-if="authStore.state.apiMode === 'azure'" class="mode-details">
              <p><strong>Azure Direct Mode Active</strong></p>
              <p>Direct connection to Azure Management APIs.</p>
              <ul>
                <li>✅ Real Azure resources</li>
                <li>✅ Full Azure integration</li>
                <li>✅ Live data and operations</li>
                <li>⚠️ Changes affect real Azure resources</li>
              </ul>
            </div>
            <div v-else class="mode-details">
              <p><strong>Unknown Mode</strong></p>
              <p>Please check your configuration or restart the application.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { SimpleApiService } from '../services/simpleApi'
import { logger } from '../services/logger'

const authStore = useAuthStore()
const loading = ref(false)
const lastUpdated = ref(new Date())

const stats = ref({
  resourceGroups: 0,
  storageAccounts: 0
})

const isHealthy = computed(() => {
  return authStore.state.isAuthenticated && !loading.value
})

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

const formatDate = (date: Date): string => {
  return date.toLocaleString()
}

const refreshData = async () => {
  loading.value = true
  try {
    logger.info('Refreshing dashboard data')
    
    const apiService = new SimpleApiService(
      authStore.state.subscriptionId || 'demo-subscription',
      authStore.state.apiMode || 'demo'
    )
    
    if (authStore.state.accessToken) {
      apiService.setAccessToken(authStore.state.accessToken)
    }

    // Fetch resource counts
    const [resourceGroups, storageAccounts] = await Promise.all([
      apiService.getResourceGroups().catch(() => []),
      apiService.getStorageAccounts().catch(() => [])
    ])

    stats.value.resourceGroups = resourceGroups.length
    stats.value.storageAccounts = storageAccounts.length
    lastUpdated.value = new Date()
    
    logger.info('Dashboard data refreshed', { 
      resourceGroups: stats.value.resourceGroups,
      storageAccounts: stats.value.storageAccounts
    })
  } catch (error) {
    logger.logError(error, 'Failed to refresh dashboard data')
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await refreshData()
})
</script>

<style scoped>
.dashboard {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.dashboard-header {
  margin-bottom: 2rem;
}

.dashboard-header h1 {
  font-size: 2.5rem;
  font-weight: 700;
  color: #1e293b;
  margin: 0;
}

.dashboard-subtitle {
  color: #64748b;
  font-size: 1.125rem;
  margin: 0.5rem 0 0 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  position: relative;
  transition: all 0.2s;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  font-size: 2rem;
  margin-right: 1rem;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 1.875rem;
  font-weight: 700;
  color: #1e293b;
  margin-bottom: 0.25rem;
}

.stat-label {
  color: #64748b;
  font-size: 0.875rem;
  font-weight: 500;
}

.stat-link {
  position: absolute;
  top: 0.75rem;
  right: 1rem;
  color: #0369a1;
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 500;
}

.stat-link:hover {
  color: #0284c7;
}

.status-healthy {
  color: #059669;
}

.status-unhealthy {
  color: #dc2626;
}

.content-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1.5rem;
}

.content-card {
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

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.action-btn {
  display: flex;
  align-items: center;
  padding: 1rem;
  border-radius: 8px;
  text-decoration: none;
  transition: all 0.2s;
  border: none;
  cursor: pointer;
  width: 100%;
  text-align: left;
}

.action-btn.primary {
  background: #f0f9ff;
  color: #0369a1;
  border: 1px solid #bae6fd;
}

.action-btn.primary:hover {
  background: #e0f2fe;
  border-color: #7dd3fc;
}

.action-btn.secondary {
  background: #f8fafc;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.action-btn.secondary:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

.action-btn:disabled {
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

.info-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-label {
  font-weight: 500;
  color: #374151;
  font-size: 0.875rem;
}

.info-value {
  color: #1e293b;
  font-size: 0.875rem;
}

.subscription-id {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  background: #f8fafc;
  padding: 0.25rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
}

.mode-demo { color: #0369a1; }
.mode-backend { color: #059669; }
.mode-azure { color: #dc2626; }
.mode-unknown { color: #6b7280; }

.activity-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.activity-item {
  display: flex;
  align-items: flex-start;
}

.activity-icon {
  font-size: 1.25rem;
  margin-right: 1rem;
  margin-top: 0.125rem;
}

.activity-content {
  flex: 1;
}

.activity-title {
  font-weight: 600;
  color: #1e293b;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
}

.activity-subtitle {
  color: #64748b;
  font-size: 0.75rem;
}

.mode-info {
  color: #374151;
  font-size: 0.875rem;
}

.mode-details p {
  margin: 0 0 0.75rem 0;
}

.mode-details ul {
  margin: 0.75rem 0 0 0;
  padding-left: 1.25rem;
}

.mode-details li {
  margin-bottom: 0.375rem;
  font-size: 0.8125rem;
}

@media (max-width: 768px) {
  .dashboard {
    padding: 1rem;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .content-grid {
    grid-template-columns: 1fr;
  }
  
  .dashboard-header h1 {
    font-size: 2rem;
  }
  
  .stat-card {
    flex-direction: column;
    text-align: center;
  }
  
  .stat-icon {
    margin-right: 0;
    margin-bottom: 0.5rem;
  }
  
  .stat-link {
    position: static;
    margin-top: 0.75rem;
  }
}
</style>