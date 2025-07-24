<template>
  <div class="login-container">
    <div class="login-card">
      <h2>Azure Resource Manager</h2>
      <p class="subtitle">Sign in with your Microsoft account</p>
      
      <div v-if="authStore.state.loading" class="loading-state">
        <div class="spinner"></div>
        <p>Authenticating with Azure...</p>
      </div>

      <div v-else-if="authStore.state.error" class="error-message">
        <h4>Authentication Error</h4>
        <p>{{ authStore.state.error }}</p>
        <button @click="retryLogin" class="retry-button">Try Again</button>
      </div>

      <div v-else class="login-form">
        <div class="form-group">
          <label for="subscriptionId">Azure Subscription ID:</label>
          <input
            id="subscriptionId"
            v-model="subscriptionId"
            type="text"
            placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
            class="form-input"
          />
          <small>Optional: Leave empty to use environment variable</small>
        </div>

        <button 
          @click="handleLogin"
          :disabled="authStore.state.loading"
          class="login-button"
        >
          <span class="login-icon">🔐</span>
          Sign in with Microsoft
        </button>

        <div class="info-box">
          <h4>Azure AD Authentication</h4>
          <p>This application uses Microsoft Authentication Library (MSAL) for secure Azure AD login.</p>
          <ul>
            <li>You'll be redirected to Microsoft's login page</li>
            <li>Sign in with your Azure account credentials</li>
            <li>Grant permissions for Azure Resource Management</li>
            <li>You'll be redirected back to manage your resources</li>
          </ul>
        </div>

        <div class="setup-info">
          <h4>Setup Required</h4>
          <p>To use this application, you need:</p>
          <ul>
            <li>An Azure AD App Registration with client ID configured</li>
            <li>Redirect URI configured: <code>http://localhost:5173</code></li>
            <li>API permissions for Azure Resource Manager</li>
          </ul>
          <p>Configure these values in your <code>.env</code> file.</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const subscriptionId = ref(import.meta.env.VITE_AZURE_SUBSCRIPTION_ID || '')

const handleLogin = async () => {
  // Update subscription ID if provided
  if (subscriptionId.value.trim()) {
    authStore.state.subscriptionId = subscriptionId.value.trim()
  }
  
  await authStore.login()
}

const retryLogin = () => {
  authStore.state.error = null
  handleLogin()
}

onMounted(async () => {
  await authStore.initializeAuth()
})
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
  max-width: 480px;
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

.login-form {
  margin-bottom: 2rem;
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

.error-message {
  color: #dc2626;
  background: #fef2f2;
  padding: 0.75rem;
  border-radius: 6px;
  border: 1px solid #fecaca;
  margin-top: 1rem;
  font-size: 0.875rem;
}

.info-box {
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 6px;
  padding: 1rem;
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

.loading-state {
  text-align: center;
  padding: 2rem;
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

.retry-button {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  margin-top: 1rem;
}

.login-icon {
  margin-right: 0.5rem;
}

.form-group small {
  color: #6b7280;
  font-size: 0.75rem;
  margin-top: 0.25rem;
  display: block;
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

.setup-info ul {
  color: #92400e;
  margin: 0;
  padding-left: 1rem;
  font-size: 0.75rem;
}

.setup-info li {
  margin-bottom: 0.25rem;
}

.setup-info code {
  background: #fef3c7;
  padding: 0.125rem 0.25rem;
  border-radius: 3px;
  font-family: monospace;
}
</style>