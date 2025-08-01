<template>
  <div class="app-layout">
    <!-- Left Navigation Sidebar -->
    <nav class="sidebar">
      <div class="sidebar-header">
        <div class="logo">
          <h2>🌩️ Brutus</h2>
          <p class="subtitle">Azure Resource Manager</p>
        </div>
      </div>

      <div class="sidebar-content">
        <!-- User Info -->
        <div class="user-info">
          <div class="user-avatar">
            {{ getInitials(authStore.state.account?.name || 'User') }}
          </div>
          <div class="user-details">
            <div class="user-name">{{ authStore.state.account?.name || 'Demo User' }}</div>
            <div class="user-mode">{{ getModeDisplay(authStore.state.apiMode) }}</div>
          </div>
        </div>

        <!-- Navigation Menu -->
        <ul class="nav-menu">
          <li>
            <router-link to="/dashboard" class="nav-item">
              <span class="nav-icon">🏠</span>
              <span class="nav-text">Dashboard</span>
            </router-link>
          </li>
          <li>
            <router-link to="/resource-groups" class="nav-item">
              <span class="nav-icon">📁</span>
              <span class="nav-text">Resource Groups</span>
            </router-link>
          </li>
          <li>
            <router-link to="/storage-accounts" class="nav-item">
              <span class="nav-icon">💾</span>
              <span class="nav-text">Storage Accounts</span>
            </router-link>
          </li>
          <li>
            <router-link to="/settings" class="nav-item">
              <span class="nav-icon">⚙️</span>
              <span class="nav-text">Settings</span>
            </router-link>
          </li>
          <li>
            <router-link to="/debug" class="nav-item">
              <span class="nav-icon">🔧</span>
              <span class="nav-text">Debug</span>
            </router-link>
          </li>
        </ul>
      </div>

      <!-- Sidebar Footer -->
      <div class="sidebar-footer">
        <div class="subscription-info">
          <div class="subscription-label">Subscription</div>
          <div class="subscription-id">{{ formatSubscriptionId(authStore.state.subscriptionId) }}</div>
        </div>
        <button @click="handleLogout" class="logout-btn">
          <span class="nav-icon">🚪</span>
          <span class="nav-text">Logout</span>
        </button>
      </div>
    </nav>

    <!-- Main Content Area -->
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

const getInitials = (name: string): string => {
  return name
    .split(' ')
    .map(word => word.charAt(0).toUpperCase())
    .slice(0, 2)
    .join('')
}

const getModeDisplay = (mode: string | null): string => {
  switch (mode) {
    case 'demo': return 'Demo Mode'
    case 'backend': return 'Backend Mode'
    case 'azure': return 'Azure Direct'
    default: return 'Unknown Mode'
  }
}

const formatSubscriptionId = (id: string | null): string => {
  if (!id) return 'Not Set'
  return id.length > 20 ? `${id.substring(0, 8)}...${id.substring(id.length - 8)}` : id
}

const handleLogout = async () => {
  try {
    await authStore.logout()
    // Force reload to ensure clean state
    window.location.href = '/'
  } catch (error) {
    console.error('Logout failed:', error)
    // Force reload anyway
    window.location.href = '/'
  }
}
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
  background: #f8fafc;
}

.sidebar {
  width: 280px;
  background: white;
  border-right: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.sidebar-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e2e8f0;
}

.logo h2 {
  margin: 0;
  color: #1e293b;
  font-size: 1.5rem;
  font-weight: 700;
}

.subtitle {
  margin: 0.25rem 0 0 0;
  color: #64748b;
  font-size: 0.875rem;
}

.sidebar-content {
  flex: 1;
  padding: 1.5rem 0;
}

.user-info {
  display: flex;
  align-items: center;
  padding: 0 1.5rem 1.5rem 1.5rem;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid #f1f5f9;
}

.user-avatar {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.875rem;
  margin-right: 0.75rem;
}

.user-details {
  flex: 1;
}

.user-name {
  font-weight: 600;
  color: #1e293b;
  font-size: 0.875rem;
}

.user-mode {
  color: #64748b;
  font-size: 0.75rem;
  margin-top: 0.125rem;
}

.nav-menu {
  list-style: none;
  margin: 0;
  padding: 0;
}

.nav-menu li {
  margin-bottom: 0.25rem;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 0.75rem 1.5rem;
  color: #64748b;
  text-decoration: none;
  transition: all 0.2s;
  border-right: 3px solid transparent;
}

.nav-item:hover {
  background: #f8fafc;
  color: #1e293b;
}

.nav-item.router-link-active {
  background: #f0f9ff;
  color: #0369a1;
  border-right-color: #0369a1;
}

.nav-icon {
  font-size: 1.125rem;
  margin-right: 0.75rem;
  width: 20px;
  text-align: center;
}

.nav-text {
  font-weight: 500;
  font-size: 0.875rem;
}

.sidebar-footer {
  padding: 1.5rem;
  border-top: 1px solid #e2e8f0;
}

.subscription-info {
  margin-bottom: 1rem;
}

.subscription-label {
  font-size: 0.75rem;
  color: #64748b;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.25rem;
}

.subscription-id {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 0.75rem;
  color: #1e293b;
  background: #f8fafc;
  padding: 0.375rem 0.5rem;
  border-radius: 0.375rem;
  border: 1px solid #e2e8f0;
}

.logout-btn {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 0.75rem;
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.875rem;
}

.logout-btn:hover {
  background: #fee2e2;
  border-color: #fca5a5;
}

.main-content {
  flex: 1;
  overflow-x: auto;
}

/* Responsive Design */
@media (max-width: 768px) {
  .sidebar {
    width: 240px;
  }
  
  .user-info {
    flex-direction: column;
    align-items: flex-start;
    text-align: left;
  }
  
  .user-avatar {
    margin-right: 0;
    margin-bottom: 0.5rem;
  }
}

@media (max-width: 640px) {
  .app-layout {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    height: auto;
  }
  
  .nav-menu {
    display: flex;
    overflow-x: auto;
    padding: 0 1rem;
  }
  
  .nav-menu li {
    margin-bottom: 0;
    margin-right: 0.5rem;
    flex-shrink: 0;
  }
  
  .nav-item {
    padding: 0.5rem 1rem;
    border-radius: 0.5rem;
    border-right: none;
    white-space: nowrap;
  }
  
  .nav-item.router-link-active {
    background: #0369a1;
    color: white;
  }
}
</style>