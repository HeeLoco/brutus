<template>
  <div class="environment-debug">
    <h2>🔧 Environment Debug Information</h2>
    
    <div class="debug-section">
      <h3>📋 Environment Variables</h3>
      <div class="env-grid">
        <div class="env-item">
          <strong>VITE_AZURE_CLIENT_ID:</strong>
          <span class="env-value">{{ maskedValue(envVars.VITE_AZURE_CLIENT_ID) }}</span>
        </div>
        <div class="env-item">
          <strong>VITE_AZURE_TENANT_ID:</strong>
          <span class="env-value">{{ maskedValue(envVars.VITE_AZURE_TENANT_ID) }}</span>
        </div>
        <div class="env-item">
          <strong>VITE_AZURE_SUBSCRIPTION_ID:</strong>
          <span class="env-value">{{ maskedValue(envVars.VITE_AZURE_SUBSCRIPTION_ID) }}</span>
        </div>
        <div class="env-item">
          <strong>VITE_REDIRECT_URI:</strong>
          <span class="env-value" :class="{ 'warning': !isCorrectRedirectUri }">
            {{ envVars.VITE_REDIRECT_URI }}
          </span>
          <span v-if="!isCorrectRedirectUri" class="warning-text">
            ⚠️ Should be {{ expectedRedirectUri }} for current environment
          </span>
        </div>
        <div class="env-item">
          <strong>VITE_POST_LOGOUT_REDIRECT_URI:</strong>
          <span class="env-value">{{ envVars.VITE_POST_LOGOUT_REDIRECT_URI }}</span>
        </div>
        <div class="env-item">
          <strong>VITE_API_BASE_URL:</strong>
          <span class="env-value">{{ envVars.VITE_API_BASE_URL }}</span>
        </div>
        <div class="env-item">
          <strong>NODE_ENV:</strong>
          <span class="env-value" :class="nodeEnvClass">{{ envVars.NODE_ENV }}</span>
        </div>
        <div class="env-item">
          <strong>MODE:</strong>
          <span class="env-value" :class="modeClass">{{ envVars.MODE }}</span>
        </div>
      </div>
    </div>

    <div class="debug-section">
      <h3>🌐 Runtime Information</h3>
      <div class="env-grid">
        <div class="env-item">
          <strong>Current URL:</strong>
          <span class="env-value">{{ currentUrl }}</span>
        </div>
        <div class="env-item">
          <strong>Origin:</strong>
          <span class="env-value">{{ runtimeInfo.origin }}</span>
        </div>
        <div class="env-item">
          <strong>Port:</strong>
          <span class="env-value">{{ runtimeInfo.port }}</span>
        </div>
        <div class="env-item">
          <strong>Protocol:</strong>
          <span class="env-value">{{ runtimeInfo.protocol }}</span>
        </div>
        <div class="env-item">
          <strong>User Agent:</strong>
          <span class="env-value small">{{ runtimeInfo.userAgent }}</span>
        </div>
      </div>
    </div>

    <div class="debug-section">
      <h3>🔐 Authentication Status</h3>
      <div class="env-grid">
        <div class="env-item">
          <strong>Logged In:</strong>
          <span class="env-value" :class="authStore.state.isAuthenticated ? 'success' : 'error'">
            {{ authStore.state.isAuthenticated ? '✅ Yes' : '❌ No' }}
          </span>
        </div>
        <div class="env-item" v-if="authStore.state.account">
          <strong>Username:</strong>
          <span class="env-value">{{ authStore.state.account.username }}</span>
        </div>
        <div class="env-item" v-if="authStore.state.accessToken">
          <strong>Has Access Token:</strong>
          <span class="env-value success">✅ Yes</span>
        </div>
        <div class="env-item" v-if="authStore.state.subscriptionId">
          <strong>Subscription ID:</strong>
          <span class="env-value">{{ maskedValue(authStore.state.subscriptionId) }}</span>
        </div>
      </div>
    </div>

    <div class="debug-section">
      <h3>🔗 Expected Azure AD Configuration</h3>
      <div class="config-info">
        <p><strong>Redirect URIs to register in Azure AD App:</strong></p>
        <ul>
          <li><code>http://localhost:5173</code> (for development)</li>
          <li><code>http://localhost</code> (for production)</li>
        </ul>
        <p><strong>Required Permissions:</strong></p>
        <ul>
          <li>https://management.azure.com/user_impersonation</li>
          <li>openid</li>
          <li>profile</li>
          <li>offline_access</li>
        </ul>
      </div>
    </div>

    <div class="debug-section">
      <h3>🧪 Test Backend Connection</h3>
      <div class="test-section">
        <button @click="testBackendHealth" class="test-button" :disabled="testing">
          {{ testing ? 'Testing...' : 'Test Backend Health' }}
        </button>
        <div v-if="backendTest.result" class="test-result" :class="backendTest.success ? 'success' : 'error'">
          <strong>{{ backendTest.success ? '✅ Success' : '❌ Error' }}:</strong>
          <pre>{{ backendTest.result }}</pre>
        </div>
      </div>
    </div>

    <div class="debug-actions">
      <button @click="copyToClipboard" class="copy-button">
        📋 Copy Debug Info to Clipboard
      </button>
      <button @click="refreshPage" class="refresh-button">
        🔄 Refresh Page
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const currentUrl = ref(window.location.href)
const testing = ref(false)
const backendTest = ref<{ result: string | null, success: boolean }>({ result: null, success: false })

// Environment variables (Vite exposes them on import.meta.env)
const envVars = computed(() => ({
  VITE_AZURE_CLIENT_ID: import.meta.env.VITE_AZURE_CLIENT_ID || 'Not set',
  VITE_AZURE_TENANT_ID: import.meta.env.VITE_AZURE_TENANT_ID || 'Not set',
  VITE_AZURE_SUBSCRIPTION_ID: import.meta.env.VITE_AZURE_SUBSCRIPTION_ID || 'Not set',
  VITE_REDIRECT_URI: import.meta.env.VITE_REDIRECT_URI || 'Not set',
  VITE_POST_LOGOUT_REDIRECT_URI: import.meta.env.VITE_POST_LOGOUT_REDIRECT_URI || 'Not set',
  VITE_API_BASE_URL: import.meta.env.VITE_API_BASE_URL || 'Not set',
  NODE_ENV: import.meta.env.NODE_ENV || 'Not set',
  MODE: import.meta.env.MODE || 'Not set',
}))

// Runtime information computed from window object
const runtimeInfo = computed(() => ({
  origin: window.location.origin,
  port: window.location.port || '80',
  protocol: window.location.protocol,
  userAgent: navigator.userAgent
}))

// Expected redirect URI based on current port
const expectedRedirectUri = computed(() => {
  const port = window.location.port
  if (port === '5173') return 'http://localhost:5173'
  if (port === '80' || port === '') return 'http://localhost'
  return `http://localhost:${port}`
})

// Check if redirect URI matches current environment
const isCorrectRedirectUri = computed(() => {
  return envVars.value.VITE_REDIRECT_URI === expectedRedirectUri.value
})

// Styling classes
const nodeEnvClass = computed(() => ({
  'success': envVars.value.NODE_ENV === 'development',
  'warning': envVars.value.NODE_ENV === 'production',
}))

const modeClass = computed(() => ({
  'success': envVars.value.MODE === 'development',
  'warning': envVars.value.MODE === 'production',
}))

// Utility functions
const maskedValue = (value: string): string => {
  if (!value || value === 'Not set') return value
  if (value.length > 8) {
    return `${value.substring(0, 8)}...${value.substring(value.length - 4)}`
  }
  return value
}

const testBackendHealth = async () => {
  testing.value = true
  try {
    const response = await fetch(`${envVars.value.VITE_API_BASE_URL.replace('/api/v1', '')}/health`)
    const data = await response.json()
    backendTest.value = {
      result: JSON.stringify(data, null, 2),
      success: response.ok
    }
  } catch (error) {
    backendTest.value = {
      result: `Error: ${error instanceof Error ? error.message : 'Unknown error'}`,
      success: false
    }
  } finally {
    testing.value = false
  }
}

const copyToClipboard = async () => {
  const debugInfo = {
    timestamp: new Date().toISOString(),
    url: currentUrl.value,
    environment: envVars.value,
    runtime: runtimeInfo.value,
    authentication: {
      isAuthenticated: authStore.state.isAuthenticated,
      hasAccessToken: !!authStore.state.accessToken,
      username: authStore.state.account?.username,
      subscriptionId: authStore.state.subscriptionId ? maskedValue(authStore.state.subscriptionId) : null,
    }
  }
  
  try {
    await navigator.clipboard.writeText(JSON.stringify(debugInfo, null, 2))
    alert('Debug information copied to clipboard!')
  } catch (error) {
    console.error('Failed to copy to clipboard:', error)
    alert('Failed to copy to clipboard. Check console for details.')
  }
}

const refreshPage = () => {
  window.location.reload()
}

onMounted(() => {
  // Update current URL in case it changes
  currentUrl.value = window.location.href
})
</script>

<style scoped>
.environment-debug {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

.environment-debug h2 {
  color: #2c3e50;
  border-bottom: 2px solid #3498db;
  padding-bottom: 10px;
  margin-bottom: 30px;
}

.debug-section {
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
}

.debug-section h3 {
  color: #495057;
  margin-top: 0;
  margin-bottom: 15px;
}

.env-grid {
  display: grid;
  gap: 10px;
}

.env-item {
  display: flex;
  flex-direction: column;
  padding: 10px;
  background: white;
  border-radius: 4px;
  border: 1px solid #e9ecef;
}

.env-item strong {
  color: #6c757d;
  font-size: 0.9em;
  margin-bottom: 5px;
}

.env-value {
  font-family: 'Courier New', monospace;
  font-size: 0.95em;
  word-break: break-all;
}

.env-value.small {
  font-size: 0.8em;
}

.env-value.success {
  color: #28a745;
  font-weight: bold;
}

.env-value.error {
  color: #dc3545;
  font-weight: bold;
}

.env-value.warning {
  color: #fd7e14;
  font-weight: bold;
}

.warning-text {
  color: #fd7e14;
  font-size: 0.8em;
  margin-top: 5px;
  font-style: italic;
}

.config-info {
  background: white;
  padding: 15px;
  border-radius: 4px;
  border: 1px solid #e9ecef;
}

.config-info ul {
  margin: 10px 0;
  padding-left: 20px;
}

.config-info code {
  background: #f1f3f4;
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
}

.test-section {
  background: white;
  padding: 15px;
  border-radius: 4px;
  border: 1px solid #e9ecef;
}

.test-button {
  background: #007bff;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.9em;
  margin-bottom: 15px;
}

.test-button:hover:not(:disabled) {
  background: #0056b3;
}

.test-button:disabled {
  background: #6c757d;
  cursor: not-allowed;
}

.test-result {
  padding: 10px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 0.8em;
}

.test-result.success {
  background: #d4edda;
  border: 1px solid #c3e6cb;
  color: #155724;
}

.test-result.error {
  background: #f8d7da;
  border: 1px solid #f5c6cb;
  color: #721c24;
}

.test-result pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.debug-actions {
  display: flex;
  gap: 15px;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 2px solid #dee2e6;
}

.copy-button, .refresh-button {
  padding: 12px 24px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9em;
  font-weight: bold;
  transition: all 0.2s ease;
}

.copy-button {
  background: #28a745;
  color: white;
}

.copy-button:hover {
  background: #218838;
}

.refresh-button {
  background: #6c757d;
  color: white;
}

.refresh-button:hover {
  background: #5a6268;
}

@media (max-width: 768px) {
  .environment-debug {
    padding: 10px;
  }
  
  .debug-actions {
    flex-direction: column;
  }
}
</style>