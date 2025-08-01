// Fixed JSON parsing for DELETE operations - v2
import { logger } from './logger'
import { CSRFService } from './csrf'

export interface ResourceGroup {
  id?: string;
  name?: string;
  location?: string;
  tags?: Record<string, string>;
  type?: string;
  properties?: {
    provisioningState?: string;
  };
}

export interface CreateResourceGroupRequest {
  name: string;
  location: string;
  tags?: Record<string, string>;
}

// Storage Account Interfaces
export interface StorageEndpoints {
  blob?: string;
  queue?: string;
  table?: string;
  file?: string;
}

export interface StorageAccount {
  id?: string;
  name?: string;
  location?: string;
  resourceGroup?: string;
  kind?: string;
  skuName?: string;
  skuTier?: string;
  accessTier?: string;
  allowBlobPublicAccess?: boolean;
  allowSharedKeyAccess?: boolean;
  tags?: Record<string, string>;
  creationTime?: string;
  primaryEndpoints?: StorageEndpoints;
}

export interface CreateStorageAccountRequest {
  name: string;
  resourceGroup: string;
  location: string;
  kind?: string;
  skuName?: string;
  accessTier?: string;
  allowBlobPublicAccess?: boolean;
  allowSharedKeyAccess?: boolean;
  tags?: Record<string, string>;
}

export interface UpdateStorageAccountRequest {
  accessTier?: string;
  allowBlobPublicAccess?: boolean;
  allowSharedKeyAccess?: boolean;
  tags?: Record<string, string>;
}

// Simple API service that works with both Azure direct calls and backend APIs
// Demo data for testing
const DEMO_RESOURCE_GROUPS: ResourceGroup[] = [
  {
    id: '/subscriptions/demo-sub/resourceGroups/demo-rg-1',
    name: 'demo-rg-1',
    location: 'eastus',
    type: 'Microsoft.Resources/resourceGroups',
    tags: { environment: 'demo', project: 'brutus' },
    properties: { provisioningState: 'Succeeded' }
  },
  {
    id: '/subscriptions/demo-sub/resourceGroups/demo-rg-2',
    name: 'demo-rg-2',
    location: 'westus2',
    type: 'Microsoft.Resources/resourceGroups',
    tags: { environment: 'test' },
    properties: { provisioningState: 'Succeeded' }
  },
  {
    id: '/subscriptions/demo-sub/resourceGroups/production-rg',
    name: 'production-rg',
    location: 'westeurope',
    type: 'Microsoft.Resources/resourceGroups',
    tags: { environment: 'production', team: 'devops' },
    properties: { provisioningState: 'Succeeded' }
  }
];

// Demo storage accounts for testing
const DEMO_STORAGE_ACCOUNTS: StorageAccount[] = [
  {
    id: '/subscriptions/demo-sub/resourceGroups/demo-rg-1/providers/Microsoft.Storage/storageAccounts/demostorage001',
    name: 'demostorage001',
    location: 'eastus',
    resourceGroup: 'demo-rg-1',
    kind: 'StorageV2',
    skuName: 'Standard_LRS',
    skuTier: 'Standard',
    accessTier: 'Hot',
    allowBlobPublicAccess: false,
    allowSharedKeyAccess: true,
    tags: { environment: 'demo', project: 'brutus' },
    creationTime: '2024-01-15T10:30:00Z',
    primaryEndpoints: {
      blob: 'https://demostorage001.blob.core.windows.net/',
      queue: 'https://demostorage001.queue.core.windows.net/',
      table: 'https://demostorage001.table.core.windows.net/',
      file: 'https://demostorage001.file.core.windows.net/'
    }
  },
  {
    id: '/subscriptions/demo-sub/resourceGroups/demo-rg-2/providers/Microsoft.Storage/storageAccounts/teststorage02',
    name: 'teststorage02',
    location: 'westus2',
    resourceGroup: 'demo-rg-2',
    kind: 'StorageV2',
    skuName: 'Standard_GRS',
    skuTier: 'Standard',
    accessTier: 'Cool',
    allowBlobPublicAccess: true,
    allowSharedKeyAccess: true,
    tags: { environment: 'test' },
    creationTime: '2024-02-10T14:45:00Z',
    primaryEndpoints: {
      blob: 'https://teststorage02.blob.core.windows.net/',
      queue: 'https://teststorage02.queue.core.windows.net/',
      table: 'https://teststorage02.table.core.windows.net/',
      file: 'https://teststorage02.file.core.windows.net/'
    }
  },
  {
    id: '/subscriptions/demo-sub/resourceGroups/production-rg/providers/Microsoft.Storage/storageAccounts/prodstorage',
    name: 'prodstorage',
    location: 'westeurope',
    resourceGroup: 'production-rg',
    kind: 'StorageV2',
    skuName: 'Premium_LRS',
    skuTier: 'Premium',
    accessTier: 'Hot',
    allowBlobPublicAccess: false,
    allowSharedKeyAccess: false,
    tags: { environment: 'production', team: 'devops', backup: 'enabled' },
    creationTime: '2023-12-01T09:15:00Z',
    primaryEndpoints: {
      blob: 'https://prodstorage.blob.core.windows.net/',
      queue: 'https://prodstorage.queue.core.windows.net/',
      table: 'https://prodstorage.table.core.windows.net/',
      file: 'https://prodstorage.file.core.windows.net/'
    }
  }
];

export class SimpleApiService {
  private accessToken: string | null = null;
  private subscriptionId: string;
  private mode: 'demo' | 'backend' | 'backend-go' | 'backend-python' | 'backend-typescript' | 'azure' | 'azure-direct';
  private backendUrl: string;
  private demoData: ResourceGroup[];
  private demoStorageData: StorageAccount[];

  constructor(subscriptionId: string, mode: 'demo' | 'backend' | 'backend-go' | 'backend-python' | 'backend-typescript' | 'azure' | 'azure-direct' = 'demo', backendUrl?: string) {
    this.subscriptionId = subscriptionId;
    this.mode = mode;
    
    // Set backend URL based on mode if not explicitly provided
    if (backendUrl) {
      this.backendUrl = backendUrl;
    } else {
      switch (mode) {
        case 'backend-go':
          this.backendUrl = 'http://localhost:8080/api/v1';
          break;
        case 'backend-python':
          this.backendUrl = 'http://localhost:8000/api/v1';
          break;
        case 'backend-typescript':
          this.backendUrl = 'http://localhost:3000/api/v1';
          break;
        default:
          this.backendUrl = 'http://localhost:8080/api/v1';
      }
    }
    
    this.demoData = [...DEMO_RESOURCE_GROUPS]; // Copy to allow modifications
    this.demoStorageData = [...DEMO_STORAGE_ACCOUNTS]; // Copy to allow modifications
  }

  private isBackendMode(): boolean {
    return this.mode === 'backend' || this.mode === 'backend-go' || this.mode === 'backend-python' || this.mode === 'backend-typescript';
  }

  private isAzureMode(): boolean {
    return this.mode === 'azure' || this.mode === 'azure-direct';
  }

  setAccessToken(token: string) {
    this.accessToken = token;
  }

  // Direct Azure REST API calls
  private async azureRequest(endpoint: string, options: RequestInit = {}): Promise<any> {
    if (!this.accessToken) {
      throw new Error('No access token available');
    }

    const url = `https://management.azure.com/subscriptions/${this.subscriptionId}${endpoint}?api-version=2021-04-01`;
    
    // Get CSRF headers for state-changing operations
    const csrfHeaders = options.method && options.method !== 'GET' ? await CSRFService.getHeaders() : {};
    
    const response = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
        'Content-Type': 'application/json',
        ...csrfHeaders,
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Azure API Error: ${response.status} ${response.statusText} - ${errorText}`);
    }

    // Check if response has content before parsing JSON
    const contentType = response.headers.get('content-type');
    const contentLength = response.headers.get('content-length');
    
    // If there's no content or it's not JSON, return empty object
    if (contentLength === '0' || !contentType?.includes('application/json')) {
      return {};
    }

    // Try to parse JSON, but handle empty responses gracefully
    const text = await response.text();
    if (!text || text.trim() === '') {
      logger.debug('Azure API returned empty response - this is normal for DELETE operations');
      return {};
    }
    
    try {
      return JSON.parse(text);
    } catch (error) {
      logger.warn('Failed to parse Azure API response as JSON', { responseText: text });
      return {};
    }
  }

  // Backend API calls (with user token)
  private async backendRequest(endpoint: string, options: RequestInit = {}): Promise<any> {
    const url = `${this.backendUrl}${endpoint}`;
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...(options.headers as Record<string, string> || {}),
    };

    // Pass user token to backend for user-context authentication
    if (this.accessToken) {
      headers['Authorization'] = `Bearer ${this.accessToken}`;
    }

    // Add subscription ID header for backend to know which subscription to use
    headers['X-Azure-Subscription-ID'] = this.subscriptionId;
    
    // Add CSRF protection for state-changing operations
    if (options.method && options.method !== 'GET') {
      const csrfHeaders = await CSRFService.getHeaders();
      Object.assign(headers, csrfHeaders);
    }
    
    const response = await fetch(url, {
      headers,
      ...options,
    });

    if (!response.ok) {
      throw new Error(`Backend API Error: ${response.status} ${response.statusText}`);
    }

    // Check if response has content before parsing JSON
    const contentType = response.headers.get('content-type');
    const contentLength = response.headers.get('content-length');
    
    // If there's no content or it's not JSON, return empty object
    if (contentLength === '0' || !contentType?.includes('application/json')) {
      return {};
    }

    // Try to parse JSON, but handle empty responses gracefully
    const text = await response.text();
    if (!text || text.trim() === '') {
      logger.debug('Backend API returned empty response - this is normal for DELETE operations');
      return {};
    }
    
    try {
      return JSON.parse(text);
    } catch (error) {
      logger.warn('Failed to parse backend API response as JSON', { responseText: text });
      return {};
    }
  }

  async getResourceGroups(): Promise<ResourceGroup[]> {
    try {
      if (this.mode === 'demo') {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 500));
        return [...this.demoData];
      } else if (this.isBackendMode()) {
        const response = await this.backendRequest('/resource-groups');
        return Array.isArray(response) ? response : response.data || [];
      } else {
        const response = await this.azureRequest('/resourcegroups');
        return response.value || [];
      }
    } catch (error) {
      logger.logError(error, 'Failed to fetch resource groups', { mode: this.mode });
      throw error;
    }
  }

  async createResourceGroup(data: CreateResourceGroupRequest): Promise<ResourceGroup> {
    try {
      if (this.mode === 'demo') {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 1000));
        
        const newRg: ResourceGroup = {
          id: `/subscriptions/${this.subscriptionId}/resourceGroups/${data.name}`,
          name: data.name,
          location: data.location,
          type: 'Microsoft.Resources/resourceGroups',
          tags: data.tags,
          properties: { provisioningState: 'Succeeded' }
        };
        
        this.demoData.push(newRg);
        return newRg;
      } else if (this.isBackendMode()) {
        const response = await this.backendRequest('/resource-groups', {
          method: 'POST',
          body: JSON.stringify(data),
        });
        return response.data || response;
      } else {
        const response = await this.azureRequest(`/resourcegroups/${data.name}`, {
          method: 'PUT',
          body: JSON.stringify({
            location: data.location,
            tags: data.tags,
          }),
        });
        return response;
      }
    } catch (error) {
      logger.logError(error, 'Failed to create resource group', { mode: this.mode, name: data.name });
      throw error;
    }
  }

  async deleteResourceGroup(name: string): Promise<void> {
    try {
      if (this.mode === 'demo') {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 800));
        
        const index = this.demoData.findIndex(rg => rg.name === name);
        if (index >= 0) {
          this.demoData.splice(index, 1);
        }
        return;
      } else if (this.isBackendMode()) {
        await this.backendRequest(`/resource-groups/${name}`, {
          method: 'DELETE',
        });
      } else {
        await this.azureRequest(`/resourcegroups/${name}`, {
          method: 'DELETE',
        });
      }
    } catch (error) {
      logger.logError(error, 'Failed to delete resource group', { mode: this.mode, name });
      throw error;
    }
  }

  async getResourceGroup(name: string): Promise<ResourceGroup> {
    try {
      if (this.mode === 'demo') {
        const rg = this.demoData.find(rg => rg.name === name);
        if (!rg) {
          throw new Error(`Resource group '${name}' not found`);
        }
        return rg;
      } else if (this.isBackendMode()) {
        const response = await this.backendRequest(`/resource-groups/${name}`);
        return response.data || response;
      } else {
        const response = await this.azureRequest(`/resourcegroups/${name}`);
        return response;
      }
    } catch (error) {
      logger.logError(error, 'Failed to get resource group', { mode: this.mode, name });
      throw error;
    }
  }

  async checkHealth(): Promise<{ status: string; message?: string }> {
    if (this.mode === 'demo') {
      return { status: 'healthy', message: 'Demo mode - using mock data' };
    } else if (this.isBackendMode()) {
      return this.backendRequest('/health');
    } else {
      // For Azure direct, just check if we can list subscriptions
      try {
        await this.azureRequest('');
        return { status: 'healthy', message: 'Azure API accessible' };
      } catch (error) {
        return { status: 'unhealthy', message: error instanceof Error ? error.message : 'Unknown error' };
      }
    }
  }

  // Storage Account Methods

  async getStorageAccounts(): Promise<StorageAccount[]> {
    try {
      if (this.mode === 'demo') {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 600));
        return [...this.demoStorageData];
      } else if (this.isBackendMode()) {
        const response = await this.backendRequest('/storage-accounts');
        return Array.isArray(response) ? response : response.storage_accounts || response.data || [];
      } else {
        // Azure direct API calls for storage accounts
        logger.info('Making Azure API call for storage accounts', { mode: this.mode, hasToken: !!this.accessToken, subscriptionId: this.subscriptionId });
        
        try {
          // First test if we can access the subscription (basic permission check)
          const testResponse = await this.azureRequest('');
          logger.info('Azure subscription access confirmed');
        } catch (subscriptionError) {
          logger.error('Failed to access Azure subscription', { 
            error: subscriptionError instanceof Error ? subscriptionError.message : 'Unknown error',
            subscriptionId: this.subscriptionId 
          });
          throw new Error(`Azure subscription access failed: ${subscriptionError instanceof Error ? subscriptionError.message : 'Unknown error'}`);
        }
        
        const response = await this.azureRequest('/providers/Microsoft.Storage/storageAccounts');
        logger.info('Azure storage accounts API response received', { count: response.value?.length || 0 });
        return response.value || [];
      }
    } catch (error) {
      logger.logError(error, 'Failed to fetch storage accounts', { mode: this.mode, hasToken: !!this.accessToken });
      // Don't throw error for Azure Direct mode - fall back to empty array to avoid breaking the UI
      if (this.isAzureMode()) {
        logger.warn('Azure storage accounts API failed, returning empty array', { error: error instanceof Error ? error.message : 'Unknown error' });
        return [];
      }
      throw error;
    }
  }

  async getStorageAccount(resourceGroup: string, name: string): Promise<StorageAccount> {
    try {
      if (this.mode === 'demo') {
        const sa = this.demoStorageData.find(sa => sa.name === name && sa.resourceGroup === resourceGroup);
        if (!sa) {
          throw new Error(`Storage account '${name}' not found in resource group '${resourceGroup}'`);
        }
        return sa;
      } else if (this.isBackendMode()) {
        const response = await this.backendRequest(`/storage-accounts/${resourceGroup}/${name}`);
        return response.data || response;
      } else {
        const response = await this.azureRequest(`/resourceGroups/${resourceGroup}/providers/Microsoft.Storage/storageAccounts/${name}`);
        return response;
      }
    } catch (error) {
      logger.logError(error, 'Failed to get storage account', { mode: this.mode, resourceGroup, name });
      throw error;
    }
  }

  async createStorageAccount(data: CreateStorageAccountRequest): Promise<StorageAccount> {
    try {
      if (this.mode === 'demo') {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        const newSa: StorageAccount = {
          id: `/subscriptions/${this.subscriptionId}/resourceGroups/${data.resourceGroup}/providers/Microsoft.Storage/storageAccounts/${data.name}`,
          name: data.name,
          location: data.location,
          resourceGroup: data.resourceGroup,
          kind: data.kind || 'StorageV2',
          skuName: data.skuName || 'Standard_LRS',
          skuTier: data.skuName?.includes('Premium') ? 'Premium' : 'Standard',
          accessTier: data.accessTier || 'Hot',
          allowBlobPublicAccess: data.allowBlobPublicAccess ?? false,
          allowSharedKeyAccess: data.allowSharedKeyAccess ?? true,
          tags: data.tags,
          creationTime: new Date().toISOString(),
          primaryEndpoints: {
            blob: `https://${data.name}.blob.core.windows.net/`,
            queue: `https://${data.name}.queue.core.windows.net/`,
            table: `https://${data.name}.table.core.windows.net/`,
            file: `https://${data.name}.file.core.windows.net/`
          }
        };
        
        this.demoStorageData.push(newSa);
        return newSa;
      } else if (this.isBackendMode()) {
        const response = await this.backendRequest('/storage-accounts', {
          method: 'POST',
          body: JSON.stringify(data),
        });
        return response.data || response;
      } else {
        const azureData = {
          location: data.location,
          kind: data.kind || 'StorageV2',
          sku: {
            name: data.skuName || 'Standard_LRS'
          },
          properties: {
            accessTier: data.accessTier,
            allowBlobPublicAccess: data.allowBlobPublicAccess ?? false,
            allowSharedKeyAccess: data.allowSharedKeyAccess ?? true
          },
          tags: data.tags
        };
        
        const response = await this.azureRequest(
          `/resourceGroups/${data.resourceGroup}/providers/Microsoft.Storage/storageAccounts/${data.name}`,
          {
            method: 'PUT',
            body: JSON.stringify(azureData),
          }
        );
        return response;
      }
    } catch (error) {
      logger.logError(error, 'Failed to create storage account', { mode: this.mode, name: data.name });
      throw error;
    }
  }

  async updateStorageAccount(resourceGroup: string, name: string, data: UpdateStorageAccountRequest): Promise<StorageAccount> {
    try {
      if (this.mode === 'demo') {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 1500));
        
        const index = this.demoStorageData.findIndex(sa => sa.name === name && sa.resourceGroup === resourceGroup);
        if (index === -1) {
          throw new Error(`Storage account '${name}' not found in resource group '${resourceGroup}'`);
        }
        
        // Update the demo data
        this.demoStorageData[index] = {
          ...this.demoStorageData[index],
          accessTier: data.accessTier || this.demoStorageData[index].accessTier,
          allowBlobPublicAccess: data.allowBlobPublicAccess ?? this.demoStorageData[index].allowBlobPublicAccess,
          allowSharedKeyAccess: data.allowSharedKeyAccess ?? this.demoStorageData[index].allowSharedKeyAccess,
          tags: data.tags || this.demoStorageData[index].tags
        };
        
        return this.demoStorageData[index];
      } else if (this.isBackendMode()) {
        const response = await this.backendRequest(`/storage-accounts/${resourceGroup}/${name}`, {
          method: 'PUT',
          body: JSON.stringify(data),
        });
        return response.data || response;
      } else {
        const azureData = {
          properties: {
            accessTier: data.accessTier,
            allowBlobPublicAccess: data.allowBlobPublicAccess,
            allowSharedKeyAccess: data.allowSharedKeyAccess
          },
          tags: data.tags
        };
        
        const response = await this.azureRequest(
          `/resourceGroups/${resourceGroup}/providers/Microsoft.Storage/storageAccounts/${name}`,
          {
            method: 'PATCH',
            body: JSON.stringify(azureData),
          }
        );
        return response;
      }
    } catch (error) {
      logger.logError(error, 'Failed to update storage account', { mode: this.mode, resourceGroup, name });
      throw error;
    }
  }

  async deleteStorageAccount(resourceGroup: string, name: string): Promise<void> {
    try {
      if (this.mode === 'demo') {
        // Simulate API delay
        await new Promise(resolve => setTimeout(resolve, 1200));
        
        const index = this.demoStorageData.findIndex(sa => sa.name === name && sa.resourceGroup === resourceGroup);
        if (index >= 0) {
          this.demoStorageData.splice(index, 1);
        }
        return;
      } else if (this.isBackendMode()) {
        await this.backendRequest(`/storage-accounts/${resourceGroup}/${name}`, {
          method: 'DELETE',
        });
      } else {
        await this.azureRequest(
          `/resourceGroups/${resourceGroup}/providers/Microsoft.Storage/storageAccounts/${name}`,
          {
            method: 'DELETE',
          }
        );
      }
    } catch (error) {
      logger.logError(error, 'Failed to delete storage account', { mode: this.mode, resourceGroup, name });
      throw error;
    }
  }
}