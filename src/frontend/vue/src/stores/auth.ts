import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { AccountInfo, AuthenticationResult } from '@azure/msal-browser'
import { msalInstance, loginRequest, tokenRequest } from '../config/auth'

export interface AuthState {
  isAuthenticated: boolean
  account: AccountInfo | null
  accessToken: string | null
  subscriptionId: string | null
  loading: boolean
  error: string | null
}

export const useAuthStore = defineStore('auth', () => {
  const state = ref<AuthState>({
    isAuthenticated: false,
    account: null,
    accessToken: null,
    subscriptionId: import.meta.env.VITE_AZURE_SUBSCRIPTION_ID || null,
    loading: false,
    error: null
  })

  const login = async () => {
    try {
      state.value.loading = true
      state.value.error = null

      // Perform redirect login
      await msalInstance.loginRedirect(loginRequest)
    } catch (error) {
      console.error('Login failed:', error)
      state.value.error = error instanceof Error ? error.message : 'Login failed'
      state.value.loading = false
    }
  }

  const handleRedirectResult = async () => {
    try {
      state.value.loading = true
      console.log('Handling redirect result...')
      const response = await msalInstance.handleRedirectPromise()
      
      if (response) {
        // Login was successful
        console.log('Redirect response received:', response)
        setAuthenticatedState(response)
      } else {
        // Check if user is already logged in
        console.log('No redirect response, checking for existing accounts...')
        const accounts = msalInstance.getAllAccounts()
        console.log('Found accounts:', accounts.length)
        if (accounts.length > 0) {
          state.value.account = accounts[0]
          state.value.isAuthenticated = true
          await acquireTokenSilent()
        }
      }
    } catch (error) {
      console.error('Handle redirect failed:', error)
      state.value.error = error instanceof Error ? error.message : 'Authentication failed'
    } finally {
      state.value.loading = false
    }
  }

  const acquireTokenSilent = async (): Promise<string | null> => {
    try {
      if (!state.value.account) {
        throw new Error('No account found')
      }

      const response = await msalInstance.acquireTokenSilent({
        ...tokenRequest,
        account: state.value.account
      })

      state.value.accessToken = response.accessToken
      return response.accessToken
    } catch (error) {
      console.error('Silent token acquisition failed:', error)
      // If silent token acquisition fails, try interactive login
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
      
      // Clear state
      state.value.isAuthenticated = false
      state.value.account = null
      state.value.accessToken = null
      state.value.error = null

      // Perform logout
      await msalInstance.logoutRedirect({
        account: state.value.account
      })
    } catch (error) {
      console.error('Logout failed:', error)
      state.value.error = error instanceof Error ? error.message : 'Logout failed'
    } finally {
      state.value.loading = false
    }
  }

  const setAuthenticatedState = (response: AuthenticationResult) => {
    state.value.isAuthenticated = true
    state.value.account = response.account
    state.value.accessToken = response.accessToken
    state.value.error = null
  }

  const initializeAuth = async () => {
    try {
      await msalInstance.initialize()
      await handleRedirectResult()
    } catch (error) {
      console.error('Failed to initialize auth:', error)
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
    handleRedirectResult
  }
})