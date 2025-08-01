import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { AccountInfo, AuthenticationResult } from '@azure/msal-browser'
import { msalInstance, loginRequest, tokenRequest } from '../config/auth'
import { logger } from '../services/logger'
import { CookieService } from '../services/cookies'

export interface AuthState {
  isAuthenticated: boolean
  account: AccountInfo | null
  accessToken: string | null
  subscriptionId: string | null
  apiMode: 'demo' | 'backend' | 'backend-go' | 'backend-python' | 'backend-typescript' | 'azure' | 'azure-direct' | null
  loading: boolean
  error: string | null
}

export const useAuthStore = defineStore('auth', () => {
  const state = ref<AuthState>({
    isAuthenticated: false,
    account: null,
    accessToken: null,
    subscriptionId: import.meta.env.VITE_AZURE_SUBSCRIPTION_ID || null,
    apiMode: null,
    loading: false,
    error: null
  })

  const login = async () => {
    try {
      state.value.loading = true
      state.value.error = null

      // Check if there's already an active account (simple check)
      const activeAccount = msalInstance.getActiveAccount()
      if (activeAccount) {
        logger.warn('User already has an active account, skipping new login attempt')
        state.value.loading = false
        return
      }

      // Perform redirect login
      await msalInstance.loginRedirect(loginRequest)
    } catch (error) {
      logger.logError(error, 'Login failed')
      
      // Handle specific interaction_in_progress error
      if (error instanceof Error && error.message.includes('interaction_in_progress')) {
        state.value.error = 'Login is already in progress. Please wait or refresh the page to try again.'
        logger.warn('Interaction in progress detected, suggesting page refresh')
      } else {
        state.value.error = error instanceof Error ? error.message : 'Login failed'
      }
      
      state.value.loading = false
    }
  }

  const handleRedirectResult = async () => {
    try {
      state.value.loading = true
      logger.info('Handling redirect result...')
      const response = await msalInstance.handleRedirectPromise()
      
      if (response) {
        // Login was successful
        logger.info('Redirect response received', { accountId: response.account?.homeAccountId })
        setAuthenticatedState(response)
      } else {
        // Check if user is already logged in
        logger.info('No redirect response, checking for existing accounts...')
        const accounts = msalInstance.getAllAccounts()
        logger.info('Found accounts', { count: accounts.length })
        if (accounts.length > 0) {
          logger.info('Setting account and acquiring token...', { username: accounts[0].username })
          state.value.account = accounts[0]
          state.value.isAuthenticated = true
          try {
            const token = await acquireTokenSilent()
            logger.info('Token acquired successfully', { hasToken: !!token })
          } catch (tokenError) {
            logger.logError(tokenError, 'Token acquisition failed')
            // Reset auth state if token acquisition fails
            state.value.isAuthenticated = false
            state.value.account = null
          }
        }
      }
    } catch (error) {
      logger.logError(error, 'Handle redirect failed')
      state.value.error = error instanceof Error ? error.message : 'Authentication failed'
    } finally {
      state.value.loading = false
    }
  }

  const acquireTokenSilent = async (): Promise<string | null> => {
    try {
      logger.info('Attempting silent token acquisition...')
      if (!state.value.account) {
        throw new Error('No account found')
      }

      logger.info('Account found, requesting token for scopes', { scopes: tokenRequest.scopes, username: state.value.account.username })
      const response = await msalInstance.acquireTokenSilent({
        ...tokenRequest,
        account: state.value.account
      })

      logger.info('Silent token acquisition successful')
      state.value.accessToken = response.accessToken
      return response.accessToken
    } catch (error) {
      logger.logError(error, 'Silent token acquisition failed')
      // If silent token acquisition fails, try interactive login
      logger.info('Falling back to interactive login...')
      await login()
      return null
    }
  }

  const getAccessToken = async (): Promise<string | null> => {
    if (state.value.accessToken) {
      return state.value.accessToken
    }
    return await acquireTokenSilent()
  }

  const logout = async () => {
    try {
      state.value.loading = true
      
      // Check if this is a demo/backend mode (not using real MSAL authentication)
      const isDemoMode = state.value.apiMode === 'demo' || 
                         state.value.apiMode === 'backend' ||
                         state.value.apiMode === 'backend-go' ||
                         state.value.apiMode === 'backend-python' ||
                         state.value.apiMode === 'backend-typescript' ||
                         state.value.account?.username?.includes('demo') || 
                         state.value.account?.username?.includes('backend') ||
                         !state.value.accessToken
      
      logger.info('Logout initiated', { 
        isDemoMode, 
        apiMode: state.value.apiMode, 
        username: state.value.account?.username 
      })
      
      // Store account reference before clearing state
      const accountForLogout = state.value.account
      
      // Clear cookies first
      CookieService.clearAllAuthCookies()
      logger.info('Auth cookies cleared')
      
      // Force clear any session storage
      try {
        const msalKeys = Object.keys(sessionStorage).filter(key => key.includes('msal'))
        msalKeys.forEach(key => {
          sessionStorage.removeItem(key)
        })
        logger.info('MSAL session storage cleared')
      } catch (e) {
        logger.warn('Failed to clear session storage', e as Record<string, any>)
      }
      
      // Clear state completely
      const wasAuthenticated = state.value.isAuthenticated
      state.value.isAuthenticated = false
      state.value.account = null
      state.value.accessToken = null
      state.value.apiMode = null
      state.value.subscriptionId = null
      state.value.error = null
      state.value.loading = false

      logger.info('Auth state cleared', { wasAuthenticated, isDemoMode })

      // Only perform MSAL logout for real Azure authentication
      if (!isDemoMode && wasAuthenticated && accountForLogout) {
        logger.info('Performing MSAL logout for Azure authentication')
        await msalInstance.logoutRedirect({
          account: accountForLogout
        })
      } else {
        logger.info('Demo/backend mode logout completed - state cleared without MSAL redirect')
      }
    } catch (error) {
      logger.logError(error, 'Logout failed')
      // Force clear state even if there's an error
      state.value.isAuthenticated = false
      state.value.account = null
      state.value.accessToken = null
      state.value.apiMode = null
      state.value.subscriptionId = null
      state.value.loading = false
      state.value.error = error instanceof Error ? error.message : 'Logout failed'
    }
  }

  const setAuthenticatedState = (response: AuthenticationResult) => {
    state.value.isAuthenticated = true
    state.value.account = response.account
    state.value.accessToken = response.accessToken
    state.value.error = null
    
    // Preserve apiMode if it was already set (important for Azure Direct mode)
    if (!state.value.apiMode) {
      // If no apiMode was set, default to 'azure-direct' since we're using MSAL
      state.value.apiMode = 'azure-direct'
    }

    // Store auth data in cookies
    if (response.accessToken) {
      CookieService.setAuthToken(response.accessToken, response.expiresOn ? 
        Math.floor((response.expiresOn.getTime() - Date.now()) / 1000) : 3600)
    }
    
    if (response.account) {
      CookieService.setAccountInfo(response.account)
    }

    // Store apiMode in cookies
    if (state.value.apiMode) {
      CookieService.setApiMode(state.value.apiMode)
    }

    logger.info('Auth state saved to cookies')
  }

  const clearInteractionState = async () => {
    try {
      // Clear any cached interaction state
      logger.info('Clearing interaction state...')
      
      // Clear cookies (primary storage)
      CookieService.clearAllAuthCookies()
      
      // Clear MSAL's internal cache (sessionStorage only, as configured)
      const msalSessionKeys = Object.keys(sessionStorage).filter(key => key.includes('msal'))
      msalSessionKeys.forEach(key => {
        sessionStorage.removeItem(key)
        logger.info(`Removed sessionStorage key: ${key}`)
      })
      
      // Reset auth state
      state.value.isAuthenticated = false
      state.value.account = null
      state.value.accessToken = null
      state.value.apiMode = null
      state.value.loading = false
      state.value.error = null
      
      logger.info('Interaction state cleared successfully')
    } catch (error) {
      logger.logError(error, 'Failed to clear interaction state')
    }
  }

  const restoreAuthFromCookies = () => {
    try {
      const token = CookieService.getAuthToken()
      const accountInfo = CookieService.getAccountInfo()
      const apiMode = CookieService.getApiMode()
      
      if (token && accountInfo) {
        state.value.accessToken = token
        state.value.account = accountInfo
        state.value.isAuthenticated = true
        
        // Restore apiMode from cookies
        if (apiMode) {
          state.value.apiMode = apiMode as 'demo' | 'backend' | 'backend-go' | 'backend-python' | 'backend-typescript' | 'azure' | 'azure-direct'
        }
        
        logger.info('Auth state restored from cookies', { apiMode })
        return true
      }
      
      logger.info('No valid auth cookies found')
      return false
    } catch (error) {
      logger.logError(error, 'Failed to restore auth from cookies')
      CookieService.clearAllAuthCookies() // Clear invalid cookies
      return false
    }
  }

  const initializeAuth = async () => {
    try {
      await msalInstance.initialize()
      
      // First try to restore from cookies
      const restored = restoreAuthFromCookies()
      if (restored) {
        logger.info('Auth restored from cookies, skipping redirect handling')
        return
      }
      
      // If no valid cookies, handle redirect result
      await handleRedirectResult()
    } catch (error) {
      logger.logError(error, 'Failed to initialize auth')
      state.value.error = error instanceof Error ? error.message : 'Failed to initialize authentication'
      // Set loading to false so user can see error and retry
      state.value.loading = false
    }
  }

  return {
    state,
    login,
    logout,
    getAccessToken,
    initializeAuth,
    handleRedirectResult,
    clearInteractionState,
    restoreAuthFromCookies
  }
})