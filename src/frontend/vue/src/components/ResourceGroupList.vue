<template>
  <div class="resource-groups-container">
    <div class="header">
      <div class="header-content">
        <h1>Azure Resource Groups</h1>
        <div class="header-info">
          <span class="subscription-id">{{ authStore.state.subscriptionId }}</span>
          <span class="user-info">{{ authStore.state.account?.username || authStore.state.account?.name || 'Azure User' }}</span>
        </div>
      </div>
      <div class="header-actions">
        <button @click="refreshResourceGroups" :disabled="loading" class="refresh-button">
          <span v-if="loading">⟳</span>
          <span v-else>↻</span>
          Refresh
        </button>
        <button @click="showCreateModal = true" class="create-button">
          + Create Resource Group
        </button>
        <button @click="authStore.logout()" class="logout-button">
          Logout
        </button>
      </div>
    </div>

    <div v-if="loading && resourceGroups.length === 0" class="loading-state">
      <div class="spinner"></div>
      <p>Loading resource groups...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>Failed to load resource groups</h3>
      <p>{{ error }}</p>
      <button @click="refreshResourceGroups" class="retry-button">Try Again</button>
    </div>

    <div v-else-if="resourceGroups.length === 0" class="empty-state">
      <div class="empty-icon">📁</div>
      <h3>No resource groups found</h3>
      <p>Create your first resource group to get started with Azure resources.</p>
      <button @click="showCreateModal = true" class="create-button">
        Create Resource Group
      </button>
    </div>

    <div v-else class="resource-groups-grid">
      <div
        v-for="rg in resourceGroups"
        :key="rg.id"
        class="resource-group-card"
      >
        <div class="card-header">
          <h3>{{ rg.name }}</h3>
          <div class="card-actions">
            <button
              @click="rg.name && deleteResourceGroup(rg.name)"
              :disabled="loading"
              class="delete-button"
              title="Delete Resource Group"
            >
              🗑️
            </button>
          </div>
        </div>
        <div class="card-content">
          <div class="location">
            <span class="label">Location:</span>
            <span class="value">{{ rg.location }}</span>
          </div>
          <div v-if="rg.tags && Object.keys(rg.tags).length > 0" class="tags">
            <span class="label">Tags:</span>
            <div class="tag-list">
              <span
                v-for="(value, key) in rg.tags"
                :key="key"
                class="tag"
              >
                {{ key }}: {{ value }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Resource Group Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click="closeCreateModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h3>Create Resource Group</h3>
          <button @click="closeCreateModal" class="close-button">×</button>
        </div>
        <form @submit.prevent="createResourceGroup" class="modal-form">
          <div class="form-group">
            <label for="rgName">Name:</label>
            <input
              id="rgName"
              v-model="createForm.name"
              type="text"
              required
              class="form-input"
              placeholder="my-resource-group"
            />
          </div>
          <div class="form-group">
            <label for="rgLocation">Location:</label>
            <select
              id="rgLocation"
              v-model="createForm.location"
              required
              class="form-select"
            >
              <option value="">Select a location</option>
              <option value="eastus">East US</option>
              <option value="westus2">West US 2</option>
              <option value="westeurope">West Europe</option>
              <option value="northeurope">North Europe</option>
              <option value="southeastasia">Southeast Asia</option>
            </select>
          </div>
          <div class="form-group">
            <label for="rgTags">Tags (optional):</label>
            <input
              id="rgTags"
              v-model="createForm.tagsInput"
              type="text"
              class="form-input"
              placeholder="environment:dev,project:brutus"
            />
            <small>Format: key:value,key:value</small>
          </div>
          <div class="modal-actions">
            <button type="button" @click="closeCreateModal" class="cancel-button">
              Cancel
            </button>
            <button type="submit" :disabled="creating" class="submit-button">
              {{ creating ? 'Creating...' : 'Create' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { SimpleApiService, type ResourceGroup, type CreateResourceGroupRequest } from '../services/simpleApi'

const authStore = useAuthStore()

const resourceGroups = ref<ResourceGroup[]>([])
const loading = ref(false)
const error = ref('')
const showCreateModal = ref(false)
const creating = ref(false)
let apiService: SimpleApiService | null = null

const createForm = ref({
  name: '',
  location: '',
  tagsInput: ''
})

const initializeApiService = async () => {
  if (!authStore.state.subscriptionId) {
    throw new Error('No subscription ID available')
  }
  
  // Determine API mode based on authentication state and user choice
  let mode: 'demo' | 'backend' | 'azure' = 'demo'
  
  if (authStore.state.accessToken) {
    mode = 'azure'
  } else if (authStore.state.account?.username?.includes('backend')) {
    mode = 'backend'
  }
  
  // Select backend URL based on the authenticated mode
  let backendUrl = 'http://localhost:8081/api/v1' // Default to Go backend
  if (authStore.state.account?.username?.includes('python')) {
    backendUrl = 'http://localhost:8000/api/v1'
  } else if (authStore.state.account?.username?.includes('typescript')) {
    backendUrl = 'http://localhost:3000/api/v1'
  }
  
  // Override with environment variable if set
  backendUrl = import.meta.env.VITE_API_BASE_URL || backendUrl
  
  apiService = new SimpleApiService(authStore.state.subscriptionId, mode, backendUrl)
  
  if (authStore.state.accessToken) {
    apiService.setAccessToken(authStore.state.accessToken)
  }
}

const refreshResourceGroups = async () => {
  loading.value = true
  error.value = ''
  
  try {
    if (!apiService) {
      await initializeApiService()
    }
    
    if (!apiService) {
      throw new Error('Failed to initialize API service')
    }
    
    resourceGroups.value = await apiService.getResourceGroups()
  } catch (err) {
    console.error('Failed to fetch resource groups:', err)
    error.value = err instanceof Error ? err.message : 'Unknown error occurred'
  } finally {
    loading.value = false
  }
}

const createResourceGroup = async () => {
  creating.value = true
  
  try {
    if (!apiService) {
      await initializeApiService()
    }
    
    if (!apiService) {
      throw new Error('Failed to initialize API service')
    }
    
    const tags: Record<string, string> = {}
    if (createForm.value.tagsInput.trim()) {
      createForm.value.tagsInput.split(',').forEach(tag => {
        const [key, value] = tag.split(':').map(s => s.trim())
        if (key && value) {
          tags[key] = value
        }
      })
    }

    const request: CreateResourceGroupRequest = {
      name: createForm.value.name,
      location: createForm.value.location,
      tags: Object.keys(tags).length > 0 ? tags : undefined
    }

    await apiService.createResourceGroup(request)

    closeCreateModal()
    await refreshResourceGroups()
  } catch (err) {
    console.error('Failed to create resource group:', err)
    error.value = err instanceof Error ? err.message : 'Failed to create resource group'
  } finally {
    creating.value = false
  }
}

const deleteResourceGroup = async (name: string) => {
  if (!name) {
    error.value = 'Resource group name is required'
    return
  }
  
  if (!confirm(`Are you sure you want to delete resource group "${name}"? This action cannot be undone.`)) {
    return
  }

  try {
    if (!apiService) {
      await initializeApiService()
    }
    
    if (!apiService) {
      throw new Error('Failed to initialize API service')
    }
    
    await apiService.deleteResourceGroup(name)
    await refreshResourceGroups()
  } catch (err) {
    console.error('Failed to delete resource group:', err)
    error.value = err instanceof Error ? err.message : 'Failed to delete resource group'
  }
}

const closeCreateModal = () => {
  showCreateModal.value = false
  createForm.value = {
    name: '',
    location: '',
    tagsInput: ''
  }
}

onMounted(async () => {
  try {
    await initializeApiService()
    await refreshResourceGroups()
  } catch (err) {
    console.error('Failed to initialize:', err)
    error.value = err instanceof Error ? err.message : 'Failed to initialize application'
  }
})
</script>

<style scoped>
.resource-groups-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.header-content h1 {
  margin: 0;
  color: #1f2937;
  font-size: 2rem;
  font-weight: bold;
}

.header-info {
  display: flex;
  gap: 1rem;
  align-items: center;
  margin-top: 0.5rem;
}

.subscription-id {
  font-family: monospace;
  background: #f3f4f6;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
  color: #6b7280;
}

.user-info {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: bold;
  background: #f0f9ff;
  color: #0c4a6e;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  flex-wrap: wrap;
}

.refresh-button, .create-button, .logout-button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.refresh-button {
  background: #f3f4f6;
  color: #374151;
}

.refresh-button:hover:not(:disabled) {
  background: #e5e7eb;
}

.create-button {
  background: #10b981;
  color: white;
}

.create-button:hover {
  background: #059669;
}

.logout-button {
  background: #ef4444;
  color: white;
}

.logout-button:hover {
  background: #dc2626;
}

.loading-state, .error-state, .empty-state {
  text-align: center;
  padding: 3rem;
  color: #6b7280;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f4f6;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-icon, .empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.retry-button {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  margin-top: 1rem;
}

.resource-groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.resource-group-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: box-shadow 0.2s;
}

.resource-group-card:hover {
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.card-header h3 {
  margin: 0;
  color: #1f2937;
  font-size: 1.125rem;
}

.delete-button {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25rem;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.delete-button:hover:not(:disabled) {
  background: #fee2e2;
}

.card-content .label {
  font-weight: 500;
  color: #374151;
  margin-right: 0.5rem;
}

.card-content .value {
  color: #6b7280;
}

.location {
  margin-bottom: 0.75rem;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.tag {
  background: #f3f4f6;
  color: #374151;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  margin: 0;
  color: #1f2937;
}

.close-button {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #6b7280;
}

.modal-form {
  padding: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  font-weight: 500;
  color: #374151;
  margin-bottom: 0.5rem;
}

.form-input, .form-select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 1rem;
  box-sizing: border-box;
}

.form-input:focus, .form-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-group small {
  color: #6b7280;
  font-size: 0.75rem;
  margin-top: 0.25rem;
  display: block;
}

.modal-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

.cancel-button, .submit-button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
}

.cancel-button {
  background: #f3f4f6;
  color: #374151;
}

.submit-button {
  background: #3b82f6;
  color: white;
}

.submit-button:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}
</style>