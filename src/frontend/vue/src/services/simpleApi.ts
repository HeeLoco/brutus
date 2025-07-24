// Fixed JSON parsing for DELETE operations - v2
import { logger } from './logger'

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

export class SimpleApiService {
  private accessToken: string | null = null;
  private subscriptionId: string;
  private mode: 'demo' | 'backend' | 'azure';
  private backendUrl: string;
  private demoData: ResourceGroup[];

  constructor(subscriptionId: string, mode: 'demo' | 'backend' | 'azure' = 'demo', backendUrl = 'http://localhost:8080/api/v1') {
    this.subscriptionId = subscriptionId;
    this.mode = mode;
    this.backendUrl = backendUrl;
    this.demoData = [...DEMO_RESOURCE_GROUPS]; // Copy to allow modifications
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
    
    const response = await fetch(url, {
      headers: {
        'Authorization': `Bearer ${this.accessToken}`,
        'Content-Type': 'application/json',
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
      } else if (this.mode === 'backend') {
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
      } else if (this.mode === 'backend') {
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
      } else if (this.mode === 'backend') {
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
      } else if (this.mode === 'backend') {
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
    } else if (this.mode === 'backend') {
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
}