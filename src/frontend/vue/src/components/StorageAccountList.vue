<template>
  <div class="storage-accounts-container">
    <div class="header">
      <div class="header-content">
        <h1>Storage Accounts</h1>
        <p class="page-subtitle">Manage your Azure storage accounts</p>
      </div>
      <div class="header-actions">
        <button @click="refreshStorageAccounts" :disabled="loading" class="refresh-button">
          <span v-if="loading">⟳</span>
          <span v-else>↻</span>
          Refresh
        </button>
        <button @click="showCreateModal = true" class="create-button">
          + Create Storage Account
        </button>
      </div>
    </div>

    <div v-if="loading && storageAccounts.length === 0" class="loading-state">
      <div class="spinner"></div>
      <p>Loading storage accounts...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>Failed to load storage accounts</h3>
      <p>{{ error }}</p>
      <button @click="refreshStorageAccounts" class="retry-button">Try Again</button>
    </div>

    <div v-else-if="storageAccounts.length === 0" class="empty-state">
      <div class="empty-icon">💾</div>
      <h3>No storage accounts found</h3>
      <p>Create your first storage account to store blobs, files, queues, and tables.</p>
      <button @click="showCreateModal = true" class="create-button">
        Create Storage Account
      </button>
    </div>

    <div v-else class="storage-accounts-grid">
      <div
        v-for="account in storageAccounts"
        :key="account.id"
        class="storage-account-card"
      >
        <div class="card-header">
          <div class="account-icon">💾</div>
          <div class="account-info">
            <h3 class="account-name">{{ account.name }}</h3>
            <p class="account-location">{{ account.location }}</p>
            <p class="account-resource-group">Resource Group: {{ account.resourceGroup }}</p>
          </div>
          <div class="account-actions">
            <button
              @click="deleteStorageAccount(account)"
              :disabled="loading"
              class="delete-button"
              title="Delete storage account"
            >
              🗑️
            </button>
          </div>
        </div>

        <div class="card-body">
          <div class="account-details">
            <div class="detail-row">
              <span class="detail-label">Kind:</span>
              <span class="detail-value">{{ account.kind || 'StorageV2' }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">SKU:</span>
              <span class="detail-value">{{ account.skuName || 'Standard_LRS' }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Access Tier:</span>
              <span class="detail-value access-tier" :class="getAccessTierClass(account.accessTier)">
                {{ account.accessTier || 'Hot' }}
              </span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Blob Public Access:</span>
              <span class="detail-value" :class="{ 'security-enabled': !account.allowBlobPublicAccess }">
                {{ account.allowBlobPublicAccess ? 'Allowed' : 'Disabled' }}
              </span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Shared Key Access:</span>
              <span class="detail-value" :class="{ 'security-enabled': !account.allowSharedKeyAccess }">
                {{ account.allowSharedKeyAccess ? 'Allowed' : 'Disabled' }}
              </span>
            </div>
            <div v-if="account.creationTime" class="detail-row">
              <span class="detail-label">Created:</span>
              <span class="detail-value">{{ formatDate(account.creationTime) }}</span>
            </div>
          </div>

          <div v-if="account.primaryEndpoints" class="endpoints-section">
            <h4>Primary Endpoints</h4>
            <div class="endpoints-grid">
              <div v-if="account.primaryEndpoints.blob" class="endpoint-item">
                <span class="endpoint-type">Blob:</span>
                <a :href="account.primaryEndpoints.blob" target="_blank" class="endpoint-url">
                  {{ account.primaryEndpoints.blob }}
                </a>
              </div>
              <div v-if="account.primaryEndpoints.queue" class="endpoint-item">
                <span class="endpoint-type">Queue:</span>
                <a :href="account.primaryEndpoints.queue" target="_blank" class="endpoint-url">
                  {{ account.primaryEndpoints.queue }}
                </a>
              </div>
              <div v-if="account.primaryEndpoints.table" class="endpoint-item">
                <span class="endpoint-type">Table:</span>
                <a :href="account.primaryEndpoints.table" target="_blank" class="endpoint-url">
                  {{ account.primaryEndpoints.table }}
                </a>
              </div>
              <div v-if="account.primaryEndpoints.file" class="endpoint-item">
                <span class="endpoint-type">File:</span>
                <a :href="account.primaryEndpoints.file" target="_blank" class="endpoint-url">
                  {{ account.primaryEndpoints.file }}
                </a>
              </div>
            </div>
          </div>

          <div v-if="account.tags && Object.keys(account.tags).length > 0" class="tags-section">
            <h4>Tags</h4>
            <div class="tags-grid">
              <span
                v-for="(value, key) in account.tags"
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

    <!-- Create Storage Account Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click="closeCreateModal">
      <div class="modal" @click.stop>
        <div class="modal-header">
          <h2>Create Storage Account</h2>
          <button @click="closeCreateModal" class="close-button">×</button>
        </div>
        
        <form @submit.prevent="handleCreateStorageAccount" class="modal-form">
          <div class="form-group">
            <label for="account-name">Storage Account Name*</label>
            <input
              id="account-name"
              v-model="createForm.name"
              type="text"
              required
              pattern="[a-z0-9]{3,24}"
              class="form-input"
              placeholder="Enter storage account name (3-24 lowercase letters/numbers)"
            />
            <small>Must be globally unique, 3-24 characters, lowercase letters and numbers only</small>
          </div>

          <div class="form-group">
            <label for="resource-group">Resource Group*</label>
            <select
              id="resource-group"
              v-model="createForm.resourceGroup"
              required
              class="form-select"
            >
              <option value="">Select resource group</option>
              <option v-for="rg in availableResourceGroups" :key="rg" :value="rg">
                {{ rg }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label for="location">Location*</label>
            <select
              id="location"
              v-model="createForm.location"
              required
              class="form-select"
            >
              <option value="">Select location</option>
              <option value="eastus">East US</option>
              <option value="eastus2">East US 2</option>
              <option value="westus">West US</option>
              <option value="westus2">West US 2</option>
              <option value="westeurope">West Europe</option>
              <option value="northeurope">North Europe</option>
              <option value="southeastasia">Southeast Asia</option>
              <option value="japaneast">Japan East</option>
            </select>
          </div>

          <div class="form-group">
            <label for="kind">Account Kind</label>
            <select
              id="kind"
              v-model="createForm.kind"
              class="form-select"
            >
              <option value="StorageV2">StorageV2 (General Purpose v2)</option>
              <option value="Storage">Storage (General Purpose v1)</option>
              <option value="BlobStorage">BlobStorage</option>
            </select>
          </div>

          <div class="form-group">
            <label for="sku">Performance Tier</label>
            <select
              id="sku"
              v-model="createForm.skuName"
              class="form-select"
            >
              <option value="Standard_LRS">Standard LRS (Locally Redundant)</option>
              <option value="Standard_GRS">Standard GRS (Geo Redundant)</option>
              <option value="Standard_RAGRS">Standard RA-GRS (Read-Access Geo Redundant)</option>
              <option value="Premium_LRS">Premium LRS (High Performance)</option>
            </select>
          </div>

          <div class="form-group">
            <label for="access-tier">Access Tier</label>
            <select
              id="access-tier"
              v-model="createForm.accessTier"
              class="form-select"
            >
              <option value="Hot">Hot (Frequently accessed data)</option>
              <option value="Cool">Cool (Infrequently accessed data)</option>
            </select>
          </div>

          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input
                type="checkbox"
                v-model="createForm.allowBlobPublicAccess"
              />
              Allow blob public access
            </label>
          </div>

          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input
                type="checkbox"
                v-model="createForm.allowSharedKeyAccess"
              />
              Allow shared key access
            </label>
          </div>

          <div class="modal-actions">
            <button type="button" @click="closeCreateModal" class="cancel-button">
              Cancel
            </button>
            <button type="submit" :disabled="creating" class="submit-button">
              {{ creating ? 'Creating...' : 'Create Storage Account' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { SimpleApiService, type StorageAccount, type CreateStorageAccountRequest } from '../services/simpleApi'
import { logger } from '../services/logger'

// Store
const authStore = useAuthStore()

// State
const storageAccounts = ref<StorageAccount[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const showCreateModal = ref(false)
const creating = ref(false)

// Create form
const createForm = ref<CreateStorageAccountRequest>({
  name: '',
  resourceGroup: '',
  location: '',
  kind: 'StorageV2',
  skuName: 'Standard_LRS',
  accessTier: 'Hot',
  allowBlobPublicAccess: false,
  allowSharedKeyAccess: true,
  tags: {}
})

// Available resource groups - loaded from API
const availableResourceGroups = ref<string[]>([])

// API Service
const apiService = computed(() => {
  const service = new SimpleApiService(
    authStore.state.subscriptionId || 'demo-subscription',
    authStore.state.apiMode || 'demo'
  )
  
  if (authStore.state.accessToken) {
    service.setAccessToken(authStore.state.accessToken)
  }
  
  return service
})

// Methods
async function loadResourceGroups() {
  try {
    logger.info('Loading resource groups for storage account creation')
    const resourceGroups = await apiService.value.getResourceGroups()
    availableResourceGroups.value = resourceGroups.map(rg => rg.name || '').filter(name => name)
    logger.info(`Loaded ${availableResourceGroups.value.length} resource groups`)
  } catch (err) {
    logger.logError(err, 'Failed to load resource groups')
    // Fallback to demo data only for demo mode
    if (authStore.state.apiMode === 'demo') {
      availableResourceGroups.value = ['demo-rg-1', 'demo-rg-2', 'production-rg']
    }
  }
}

async function loadStorageAccounts() {
  loading.value = true
  error.value = null
  
  try {
    logger.info('Loading storage accounts')
    const accounts = await apiService.value.getStorageAccounts()
    storageAccounts.value = accounts
    logger.info(`Loaded ${accounts.length} storage accounts`)
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred'
    error.value = errorMessage
    logger.logError(err, 'Failed to load storage accounts')
  } finally {
    loading.value = false
  }
}

async function refreshStorageAccounts() {
  await loadStorageAccounts()
}

async function deleteStorageAccount(account: StorageAccount) {
  if (!account.name || !account.resourceGroup) {
    logger.logError(new Error('Missing account name or resource group'), 'Cannot delete storage account')
    return
  }

  if (!confirm(`Are you sure you want to delete storage account "${account.name}"? This action cannot be undone.`)) {
    return
  }

  loading.value = true
  
  try {
    logger.info(`Deleting storage account: ${account.name}`)
    await apiService.value.deleteStorageAccount(account.resourceGroup, account.name)
    
    // Remove from local state
    storageAccounts.value = storageAccounts.value.filter(sa => sa.id !== account.id)
    logger.info(`Successfully deleted storage account: ${account.name}`)
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred'
    error.value = `Failed to delete storage account: ${errorMessage}`
    logger.logError(err, `Failed to delete storage account: ${account.name}`)
  } finally {
    loading.value = false
  }
}

function closeCreateModal() {
  showCreateModal.value = false
  // Reset form
  createForm.value = {
    name: '',
    resourceGroup: '',
    location: '',
    kind: 'StorageV2',
    skuName: 'Standard_LRS',
    accessTier: 'Hot',
    allowBlobPublicAccess: false,
    allowSharedKeyAccess: true,
    tags: {}
  }
}

async function handleCreateStorageAccount() {
  creating.value = true
  
  try {
    logger.info(`Creating storage account: ${createForm.value.name}`)
    const newAccount = await apiService.value.createStorageAccount(createForm.value)
    
    // Add to local state
    storageAccounts.value.push(newAccount)
    logger.info(`Successfully created storage account: ${createForm.value.name}`)
    
    closeCreateModal()
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred'
    error.value = `Failed to create storage account: ${errorMessage}`
    logger.logError(err, `Failed to create storage account: ${createForm.value.name}`)
  } finally {
    creating.value = false
  }
}

function formatDate(dateString: string): string {
  try {
    return new Date(dateString).toLocaleDateString()
  } catch {
    return dateString
  }
}

function getAccessTierClass(tier?: string): string {
  switch (tier?.toLowerCase()) {
    case 'hot': return 'tier-hot'
    case 'cool': return 'tier-cool'
    case 'archive': return 'tier-archive'
    default: return ''
  }
}

// Lifecycle
onMounted(async () => {
  // Debug: Log current authentication state
  logger.info('StorageAccountList mounted - Auth State Debug', {
    isAuthenticated: authStore.state.isAuthenticated,
    apiMode: authStore.state.apiMode,
    hasAccessToken: !!authStore.state.accessToken,
    subscriptionId: authStore.state.subscriptionId,
    accountName: authStore.state.account?.name
  })
  
  await Promise.all([
    loadResourceGroups(),
    loadStorageAccounts()
  ])
})
</script>

<style scoped>
.storage-accounts-container {
  padding: 2rem;
  background: #f8fafc;
  min-height: 100vh;
}

.header {
  background: white;
  border-radius: 12px;
  padding: 1.5rem 2rem;
  margin-bottom: 2rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
}

.header-content h1 {
  margin: 0;
  color: #1e293b;
  font-size: 2rem;
  font-weight: 700;
}

.page-subtitle {
  margin: 0.5rem 0 0 0;
  color: #64748b;
  font-size: 1rem;
}

.subscription-id, .user-info {
  font-size: 0.875rem;
  color: #6b7280;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
  align-items: center;
  flex-wrap: wrap;
}

.refresh-button, .create-button, .retry-button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.refresh-button {
  background: #f3f4f6;
  color: #374151;
}

.refresh-button:hover:not(:disabled) {
  background: #e5e7eb;
}

.refresh-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.create-button {
  background: #10b981;
  color: white;
}

.create-button:hover {
  background: #059669;
}

.nav-link, .debug-link {
  padding: 0.75rem 1.25rem;
  background: #6366f1;
  color: white;
  text-decoration: none;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.2s;
}

.nav-link:hover, .debug-link:hover {
  background: #4f46e5;
}

.logout-button {
  padding: 0.75rem 1.25rem;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.logout-button:hover {
  background: #dc2626;
}

.loading-state, .error-state, .empty-state {
  background: white;
  border-radius: 12px;
  padding: 3rem;
  text-align: center;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
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
  margin-top: 1rem;
}

.retry-button:hover {
  background: #2563eb;
}

.storage-accounts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 1.5rem;
}

.storage-account-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid #e2e8f0;
  overflow: hidden;
  transition: all 0.3s ease;
}

.storage-account-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: flex-start;
  gap: 1rem;
}

.account-icon {
  font-size: 2rem;
  flex-shrink: 0;
}

.account-info {
  flex: 1;
}

.account-name {
  margin: 0;
  color: #1f2937;
  font-size: 1.25rem;
  font-weight: 600;
}

.account-location, .account-resource-group {
  margin: 0.25rem 0 0 0;
  color: #6b7280;
  font-size: 0.875rem;
}

.account-actions {
  display: flex;
  gap: 0.5rem;
}

.delete-button {
  padding: 0.5rem;
  background: transparent;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.2s;
}

.delete-button:hover:not(:disabled) {
  background: #fca5a5;
  border-color: #f87171;
}

.delete-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.card-body {
  padding: 1.5rem;
}

.account-details {
  margin-bottom: 1.5rem;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
  border-bottom: 1px solid #f3f4f6;
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  font-weight: 500;
  color: #374151;
}

.detail-value {
  color: #6b7280;
}

.access-tier {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
}

.tier-hot {
  background: #fef3c7;
  color: #d97706;
}

.tier-cool {
  background: #dbeafe;
  color: #2563eb;
}

.tier-archive {
  background: #f3e8ff;
  color: #7c3aed;
}

.security-enabled {
  color: #059669;
  font-weight: 500;
}

.endpoints-section, .tags-section {
  margin-top: 1.5rem;
}

.endpoints-section h4, .tags-section h4 {
  margin: 0 0 0.75rem 0;
  color: #374151;
  font-size: 1rem;
  font-weight: 600;
}

.endpoints-grid {
  display: grid;
  gap: 0.5rem;
}

.endpoint-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
}

.endpoint-type {
  font-weight: 500;
  color: #374151;
  min-width: 60px;
}

.endpoint-url {
  color: #3b82f6;
  text-decoration: none;
  word-break: break-all;
}

.endpoint-url:hover {
  text-decoration: underline;
}

.tags-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  background: #f3f4f6;
  color: #374151;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
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
  padding: 1rem;
}

.modal {
  background: white;
  border-radius: 16px;
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-header {
  padding: 1.5rem 2rem;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h2 {
  margin: 0;
  color: #1f2937;
  font-size: 1.5rem;
  font-weight: 600;
}

.close-button {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #6b7280;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.close-button:hover {
  background: #f3f4f6;
}

.modal-form {
  padding: 2rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #374151;
  font-weight: 500;
}

.form-input, .form-select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-input:focus, .form-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-group small {
  display: block;
  margin-top: 0.25rem;
  color: #6b7280;
  font-size: 0.875rem;
}

.checkbox-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0;
  cursor: pointer;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

.cancel-button, .submit-button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
}

.cancel-button {
  background: #f3f4f6;
  color: #374151;
}

.cancel-button:hover {
  background: #e5e7eb;
}

.submit-button {
  background: #3b82f6;
  color: white;
}

.submit-button:hover:not(:disabled) {
  background: #2563eb;
}

.submit-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .storage-accounts-container {
    padding: 1rem;
  }
  
  .header {
    flex-direction: column;
    align-items: stretch;
  }
  
  .header-actions {
    justify-content: center;
  }
  
  .storage-accounts-grid {
    grid-template-columns: 1fr;
  }
  
  .modal {
    margin: 1rem;
    max-width: none;
  }
}
</style>