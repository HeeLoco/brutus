import { PublicClientApplication } from '@azure/msal-browser';
import type { Configuration, RedirectRequest, PopupRequest } from '@azure/msal-browser';

// MSAL configuration
export const msalConfig: Configuration = {
  auth: {
    clientId: import.meta.env.VITE_AZURE_CLIENT_ID || 'demo-client-id',
    authority: `https://login.microsoftonline.com/${import.meta.env.VITE_AZURE_TENANT_ID || 'common'}`,
    redirectUri: import.meta.env.VITE_REDIRECT_URI || (typeof window !== 'undefined' ? window.location.origin : 'http://localhost:5173'),
    postLogoutRedirectUri: import.meta.env.VITE_POST_LOGOUT_REDIRECT_URI || (typeof window !== 'undefined' ? window.location.origin : 'http://localhost:5173'),
  },
  cache: {
    cacheLocation: 'sessionStorage', // Use sessionStorage as fallback for MSAL's internal cache
    storeAuthStateInCookie: true, // Enable cookie storage for auth state
  },
};

// Scopes for Azure Resource Manager
export const ARM_SCOPES = [
  'https://management.azure.com/user_impersonation',
  'openid',
  'profile',
  'offline_access'
];

// Login request configuration
export const loginRequest: RedirectRequest = {
  scopes: ARM_SCOPES,
  prompt: 'select_account', // Forces account selection
};

// Silent token request configuration
export const tokenRequest = {
  scopes: ARM_SCOPES,
  forceRefresh: false, // Set to true to skip cache lookup
};

// Create MSAL instance
export const msalInstance = new PublicClientApplication(msalConfig);

// Initialize MSAL
export const initializeMsal = async () => {
  try {
    await msalInstance.initialize();
    console.log('MSAL initialized successfully');
  } catch (error) {
    console.error('Failed to initialize MSAL:', error);
    throw error;
  }
};