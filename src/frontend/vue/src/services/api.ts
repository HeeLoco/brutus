import { ResourceManagementClient } from '@azure/arm-resources';
import type { TokenCredential } from '@azure/core-auth';

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

// Custom token credential using MSAL token
class MSALTokenCredential implements TokenCredential {
  constructor(private getTokenFn: () => Promise<string | null>) {}

  async getToken(): Promise<{ token: string; expiresOnTimestamp: number } | null> {
    const token = await this.getTokenFn();
    if (!token) return null;

    // Parse JWT to get expiration (basic implementation)
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return {
        token,
        expiresOnTimestamp: (payload.exp || Date.now() / 1000 + 3600) * 1000
      };
    } catch {
      // Default to 1 hour expiration if parsing fails
      return {
        token,
        expiresOnTimestamp: Date.now() + 3600000
      };
    }
  }
}

export class AzureApiService {
  private resourceClient: ResourceManagementClient | null = null;
  private subscriptionId: string;

  constructor(subscriptionId: string) {
    this.subscriptionId = subscriptionId;
  }

  initialize(getTokenFn: () => Promise<string | null>) {
    const credential = new MSALTokenCredential(getTokenFn);
    this.resourceClient = new ResourceManagementClient(credential, this.subscriptionId);
  }

  async getResourceGroups(): Promise<ResourceGroup[]> {
    if (!this.resourceClient) {
      throw new Error('API service not initialized. Please login first.');
    }

    try {
      const resourceGroups: ResourceGroup[] = [];
      
      for await (const resourceGroup of this.resourceClient.resourceGroups.list()) {
        resourceGroups.push({
          id: resourceGroup.id,
          name: resourceGroup.name,
          location: resourceGroup.location,
          tags: resourceGroup.tags,
          type: resourceGroup.type,
          properties: resourceGroup.properties
        });
      }

      return resourceGroups;
    } catch (error) {
      console.error('Failed to fetch resource groups:', error);
      throw new Error(`Failed to fetch resource groups: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  async createResourceGroup(data: CreateResourceGroupRequest): Promise<ResourceGroup> {
    if (!this.resourceClient) {
      throw new Error('API service not initialized. Please login first.');
    }

    try {
      const result = await this.resourceClient.resourceGroups.createOrUpdate(data.name, {
        location: data.location,
        tags: data.tags
      });

      return {
        id: result.id,
        name: result.name,
        location: result.location,
        tags: result.tags,
        type: result.type,
        properties: result.properties
      };
    } catch (error) {
      console.error('Failed to create resource group:', error);
      throw new Error(`Failed to create resource group: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  async deleteResourceGroup(name: string): Promise<void> {
    if (!this.resourceClient) {
      throw new Error('API service not initialized. Please login first.');
    }

    try {
      // Start the delete operation (it's async in Azure)
      await this.resourceClient.resourceGroups.beginDeleteAndWait(name);
    } catch (error) {
      console.error('Failed to delete resource group:', error);
      throw new Error(`Failed to delete resource group: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  async getResourceGroup(name: string): Promise<ResourceGroup> {
    if (!this.resourceClient) {
      throw new Error('API service not initialized. Please login first.');
    }

    try {
      const result = await this.resourceClient.resourceGroups.get(name);
      return {
        id: result.id,
        name: result.name,
        location: result.location,
        tags: result.tags,
        type: result.type,
        properties: result.properties
      };
    } catch (error) {
      console.error('Failed to get resource group:', error);
      throw new Error(`Failed to get resource group: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }
}