<template>
  <div class="login-container">
    <div class="login-card">
      <h2>Azure Resource Manager</h2>
      <p class="subtitle">Demo Mode - Test without Azure credentials</p>
      
      <div class="demo-form">
        <div class="form-group">
          <label for="subscriptionId">Azure Subscription ID:</label>
          <input
            id="subscriptionId"
            v-model="subscriptionId"
            type="text"
            placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
            class="form-input"
          />
          <small>Enter any valid GUID format for demo</small>
        </div>

        <div class="form-group">
          <label for="backendType">API Mode:</label>
          <select
            id="backendType"
            v-model="apiMode"
            class="form-select"
          >
            <option value="demo">Demo Mode (Mock Data)</option>
            <option value="backend-go">Go Backend (Port 8080) - Currently Running</option>
            <option value="backend-python" disabled>Python Backend (Port 8000) - Not Available</option>
            <option value="backend-typescript" disabled>TypeScript Backend (Port 3000) - Not Available</option>
            <option value="azure-direct">Azure Direct (Real Azure API)</option>
          </select>
        </div>

        <button 
          @click="handleLogin"
          :disabled="loading"
          class="login-button"
        >
          {{ loading ? 'Connecting...' : 'Connect' }}
        </button>

        <button 
          @click="clearLoginState"
          :disabled="loading"
          class="clear-button"
          title="Clear cached login state if authentication is stuck"
        >
          Clear Login State
        </button>

        <button 
          @click="refreshPage"
          :disabled="loading"
          class="refresh-button"
          title="Refresh the page for a complete reset"
        >
          Refresh Page
        </button>

        <div v-if="error" class="error-message">
          {{ error }}
          <div v-if="error.includes('interaction_in_progress') || error.includes('already in progress')" class="error-suggestion">
            💡 Try clicking "Clear Login State" button above and then try logging in again.
          </div>
        </div>

        <div class="info-box">
          <h4>Available Modes</h4>
          <p>Choose your preferred way to interact with Azure resources:</p>
          <ul>
            <li><strong>Demo Mode:</strong> Uses mock data, perfect for testing (no setup required)</li>
            <li><strong>Go Backend:</strong> Connects to the running Go API server (currently available)</li>
            <li><strong>Azure Direct:</strong> Real Azure authentication with your credentials</li>
          </ul>
          <p><em>Note: Python and TypeScript backends are implemented but not currently running in Docker.</em></p>
        </div>

        <div v-if="apiMode === 'azure-direct'" class="setup-info">
          <h4>Azure Setup Required</h4>
          <p>For Azure Direct mode, configure your <code>.env</code> file:</p>
          <pre><code>VITE_AZURE_CLIENT_ID=your-app-registration-client-id
VITE_AZURE_TENANT_ID=your-tenant-id
VITE_AZURE_SUBSCRIPTION_ID=your-subscription-id</code></pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { CookieService } from '../services/cookies'

const authStore = useAuthStore()
const subscriptionId = ref('7b1c880d-e26e-4916-b60d-5e475ea49dca')
const apiMode = ref('azure-direct')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  if (!subscriptionId.value.trim()) {
    error.value = 'Please enter a subscription ID'
    return
  }

  loading.value = true
  error.value = ''

  try {
    console.log('Starting login with mode:', apiMode.value)
    // Set up the auth state based on selected mode
    authStore.state.subscriptionId = subscriptionId.value.trim()
    
    if (apiMode.value === 'azure-direct') {
      console.log('Using MSAL authentication')
      authStore.state.apiMode = 'azure-direct'
      // Save apiMode to cookies immediately to persist through MSAL redirect
      CookieService.setApiMode('azure-direct')
      // Use MSAL authentication - don't set isAuthenticated here, let MSAL handle it
      await authStore.login()
    } else {
      console.log('Using demo/backend mode')
      // Direct login for demo/backend modes - keep the specific mode
      authStore.state.apiMode = apiMode.value as 'demo' | 'backend' | 'backend-go' | 'backend-python' | 'backend-typescript' | 'azure' | 'azure-direct'
      // Save apiMode to cookies
      CookieService.setApiMode(apiMode.value)
      authStore.state.isAuthenticated = true
      authStore.state.account = {
        homeAccountId: 'demo-account',
        environment: 'demo',
        tenantId: 'demo-tenant',
        username: apiMode.value === 'demo' ? 'demo@demo.com' : `${apiMode.value}@demo.com`,
        localAccountId: 'demo-local',
        name: apiMode.value === 'demo' ? 'Demo User' : `${apiMode.value.replace('-', ' ').replace('backend', 'Backend')} User`
      }
    }
    
  } catch (err) {
    console.error('Login failed:', err)
    error.value = err instanceof Error ? err.message : 'Login failed'
    authStore.state.isAuthenticated = false
  } finally {
    // Only set loading to false for non-azure-direct modes
    // Azure direct will handle this in the redirect
    if (apiMode.value !== 'azure-direct') {
      loading.value = false
    }
  }
}

const clearLoginState = async () => {
  try {
    loading.value = true
    error.value = ''
    
    console.log('Clearing login state...')
    await authStore.clearInteractionState()
    
    // Small delay to let MSAL reset
    await new Promise(resolve => setTimeout(resolve, 500))
    
    console.log('Login state cleared successfully')
    error.value = '✅ Login state cleared. You can now try logging in again, or refresh the page for a complete reset.'
  } catch (err) {
    console.error('Failed to clear login state:', err)
    error.value = 'Failed to clear login state'
  } finally {
    loading.value = false
  }
}

const refreshPage = () => {
  window.location.reload()
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 1rem;
}

.login-card {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  width: 100%;
  max-width: 520px;
}

h2 {
  color: #1f2937;
  font-size: 1.875rem;
  font-weight: bold;
  text-align: center;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: #6b7280;
  text-align: center;
  margin-bottom: 2rem;
}

.demo-form {
  margin-bottom: 1rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  color: #374151;
  font-weight: 500;
  margin-bottom: 0.5rem;
}

.form-input,
.form-select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 1rem;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-input:focus,
.form-select:focus {
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

.login-button {
  width: 100%;
  background: #3b82f6;
  color: white;
  border: none;
  padding: 0.875rem;
  border-radius: 6px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
}

.login-button:hover:not(:disabled) {
  background: #2563eb;
}

.login-button:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.clear-button {
  width: 100%;
  background: #f59e0b;
  color: white;
  border: none;
  padding: 0.875rem;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-top: 0.5rem;
}

.clear-button:hover:not(:disabled) {
  background: #d97706;
}

.clear-button:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.refresh-button {
  width: 100%;
  background: #6b7280;
  color: white;
  border: none;
  padding: 0.875rem;
  border-radius: 6px;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s;
  margin-top: 0.5rem;
}

.refresh-button:hover:not(:disabled) {
  background: #4b5563;
}

.refresh-button:disabled {
  background: #9ca3af;
  cursor: not-allowed;
}

.error-message {
  color: #dc2626;
  background: #fef2f2;
  padding: 0.75rem;
  border-radius: 6px;
  border: 1px solid #fecaca;
  margin-top: 1rem;
  font-size: 0.875rem;
}

.error-suggestion {
  margin-top: 0.5rem;
  padding: 0.5rem;
  background: #fef3c7;
  color: #92400e;
  border-radius: 4px;
  font-size: 0.8rem;
  border-left: 3px solid #f59e0b;
}

.info-box {
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 6px;
  padding: 1rem;
  margin-top: 1rem;
}

.info-box h4 {
  color: #0c4a6e;
  margin: 0 0 0.5rem 0;
  font-size: 0.875rem;
  font-weight: 600;
}

.info-box p {
  color: #0c4a6e;
  margin: 0 0 0.5rem 0;
  font-size: 0.75rem;
}

.info-box ul {
  color: #0c4a6e;
  margin: 0;
  padding-left: 1rem;
  font-size: 0.75rem;
}

.info-box li {
  margin-bottom: 0.25rem;
}

.setup-info {
  background: #fefce8;
  border: 1px solid #fef08a;
  border-radius: 6px;
  padding: 1rem;
  margin-top: 1rem;
}

.setup-info h4 {
  color: #92400e;
  margin: 0 0 0.5rem 0;
  font-size: 0.875rem;
  font-weight: 600;
}

.setup-info p {
  color: #92400e;
  margin: 0 0 0.5rem 0;
  font-size: 0.75rem;
}

.setup-info pre {
  background: #fef3c7;
  padding: 0.5rem;
  border-radius: 4px;
  font-size: 0.6rem;
  overflow-x: auto;
  margin: 0.5rem 0 0 0;
}

.setup-info code {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  color: #92400e;
}
</style>