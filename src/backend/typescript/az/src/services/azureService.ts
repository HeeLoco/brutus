import { ResourceManagementClient } from '@azure/arm-resources';
import { DefaultAzureCredential, ClientSecretCredential, TokenCredential } from '@azure/identity';
import { 
  ResourceGroup, 
  CreateResourceGroupRequest, 
  UpdateResourceGroupRequest,
  AzureErrorResponse 
} from '@/types/azure';
import config from '@/config/environment';
import { createLogger } from '@/config/logger';

export class AzureResourceError extends Error {
  constructor(
    message: string,
    public readonly code: string = 'AZURE_ERROR',
    public readonly statusCode: number = 500,
    public readonly azureError?: AzureErrorResponse
  ) {
    super(message);
    this.name = 'AzureResourceError';
  }
}

export class AzureService {
  private resourceClient: ResourceManagementClient;
  private logger = createLogger();

  constructor(correlationId?: string) {
    this.logger = createLogger(correlationId);
    this.resourceClient = this.createResourceClient();
  }

  private createResourceClient(): ResourceManagementClient {
    try {
      let credential: TokenCredential;

      // Try service principal authentication first
      if (config.AZURE_CLIENT_ID && config.AZURE_CLIENT_SECRET && config.AZURE_TENANT_ID) {
        this.logger.info('Using service principal authentication');
        credential = new ClientSecretCredential(
          config.AZURE_TENANT_ID,
          config.AZURE_CLIENT_ID,
          config.AZURE_CLIENT_SECRET
        );
      } else {
        this.logger.info('Using default Azure credential');
        credential = new DefaultAzureCredential();
      }

      const client = new ResourceManagementClient(
        credential,
        config.AZURE_SUBSCRIPTION_ID
      );

      this.logger.info('Successfully created Azure resource client');
      return client;

    } catch (error) {
      this.logger.error('Failed to create Azure resource client', { error });
      throw new AzureResourceError(
        'Failed to initialize Azure connection',
        'INITIALIZATION_ERROR',
        503
      );
    }
  }

  async testConnection(): Promise<boolean> {
    try {
      this.logger.info('Testing Azure connection');
      const startTime = Date.now();
      
      // Try to list resource groups with a limit to test connectivity
      const iterator = this.resourceClient.resourceGroups.list({ top: 1 });
      await iterator.next();
      
      const responseTime = Date.now() - startTime;
      this.logger.info('Azure connectivity test successful', { responseTime });
      return true;
      
    } catch (error) {
      this.logger.error('Azure connectivity test failed', { error });
      return false;
    }
  }

  async listResourceGroups(): Promise<ResourceGroup[]> {
    try {
      this.logger.info('Listing resource groups');
      const resourceGroups: ResourceGroup[] = [];

      for await (const resourceGroup of this.resourceClient.resourceGroups.list()) {
        resourceGroups.push({
          id: resourceGroup.id ?? '',
          name: resourceGroup.name ?? '',
          location: resourceGroup.location ?? '',
          type: resourceGroup.type ?? 'Microsoft.Resources/resourceGroups',
          tags: resourceGroup.tags ?? {},
          managedBy: resourceGroup.managedBy,
          provisioningState: resourceGroup.properties?.provisioningState
        });
      }

      this.logger.info(`Found ${resourceGroups.length} resource groups`);
      return resourceGroups;

    } catch (error) {
      this.logger.error('Failed to list resource groups', { error });
      this.handleAzureError(error);
      throw error; // This won't be reached due to handleAzureError throwing
    }
  }

  async getResourceGroup(name: string): Promise<ResourceGroup> {
    try {
      this.logger.info(`Getting resource group: ${name}`);
      
      const resourceGroup = await this.resourceClient.resourceGroups.get(name);
      
      const result: ResourceGroup = {
        id: resourceGroup.id ?? '',
        name: resourceGroup.name ?? '',
        location: resourceGroup.location ?? '',
        type: resourceGroup.type ?? 'Microsoft.Resources/resourceGroups',
        tags: resourceGroup.tags ?? {},
        managedBy: resourceGroup.managedBy,
        provisioningState: resourceGroup.properties?.provisioningState
      };

      this.logger.info(`Successfully retrieved resource group: ${name}`);
      return result;

    } catch (error) {
      this.logger.error(`Failed to get resource group ${name}`, { error });
      this.handleAzureError(error, name);
      throw error;
    }
  }

  async createResourceGroup(request: CreateResourceGroupRequest): Promise<ResourceGroup> {
    try {
      this.logger.info(`Creating resource group: ${request.name}`);
      
      const parameters = {
        location: request.location,
        tags: request.tags ?? {}
      };

      const result = await this.resourceClient.resourceGroups.createOrUpdate(
        request.name,
        parameters
      );

      const resourceGroup: ResourceGroup = {
        id: result.id ?? '',
        name: result.name ?? '',
        location: result.location ?? '',
        type: result.type ?? 'Microsoft.Resources/resourceGroups',
        tags: result.tags ?? {},
        managedBy: result.managedBy,
        provisioningState: result.properties?.provisioningState
      };

      this.logger.info(`Successfully created resource group: ${request.name}`);
      return resourceGroup;

    } catch (error) {
      this.logger.error(`Failed to create resource group ${request.name}`, { error });
      this.handleAzureError(error, request.name);
      throw error;
    }
  }

  async updateResourceGroup(name: string, request: UpdateResourceGroupRequest): Promise<ResourceGroup> {
    try {
      this.logger.info(`Updating resource group: ${name}`);
      
      // First check if resource group exists
      await this.getResourceGroup(name);
      
      const parameters = {
        tags: request.tags ?? {}
      };

      const result = await this.resourceClient.resourceGroups.update(name, parameters);

      const resourceGroup: ResourceGroup = {
        id: result.id ?? '',
        name: result.name ?? '',
        location: result.location ?? '',
        type: result.type ?? 'Microsoft.Resources/resourceGroups',
        tags: result.tags ?? {},
        managedBy: result.managedBy,
        provisioningState: result.properties?.provisioningState
      };

      this.logger.info(`Successfully updated resource group: ${name}`);
      return resourceGroup;

    } catch (error) {
      this.logger.error(`Failed to update resource group ${name}`, { error });
      this.handleAzureError(error, name);
      throw error;
    }
  }

  async deleteResourceGroup(name: string): Promise<void> {
    try {
      this.logger.info(`Deleting resource group: ${name}`);
      
      // First check if resource group exists
      await this.getResourceGroup(name);
      
      // Start the deletion (this is an async operation in Azure)
      const deletePoller = await this.resourceClient.resourceGroups.beginDelete(name);
      
      this.logger.info(`Started deletion of resource group: ${name}`);
      // Note: We don't wait for completion as it can take a long time
      
    } catch (error) {
      this.logger.error(`Failed to delete resource group ${name}`, { error });
      this.handleAzureError(error, name);
      throw error;
    }
  }

  private handleAzureError(error: any, resourceName?: string): never {
    // Handle different types of Azure errors
    if (error?.code === 'ResourceGroupNotFound' || error?.statusCode === 404) {
      throw new AzureResourceError(
        `Resource group${resourceName ? ` '${resourceName}'` : ''} not found`,
        'RESOURCE_NOT_FOUND',
        404
      );
    }
    
    if (error?.code === 'ResourceGroupBeingDeleted') {
      throw new AzureResourceError(
        `Resource group '${resourceName}' is being deleted`,
        'RESOURCE_BEING_DELETED',
        409
      );
    }

    if (error?.statusCode === 403 || error?.code?.includes('Authorization')) {
      throw new AzureResourceError(
        'Insufficient permissions to perform this operation',
        'AUTHORIZATION_ERROR',
        403
      );
    }

    if (error?.statusCode === 409 || error?.code?.includes('Conflict')) {
      throw new AzureResourceError(
        `Resource group '${resourceName}' already exists or has a conflict`,
        'CONFLICT_ERROR',
        409
      );
    }

    // Handle rate limiting
    if (error?.statusCode === 429) {
      throw new AzureResourceError(
        'Too many requests to Azure API, please try again later',
        'RATE_LIMIT_ERROR',
        429
      );
    }

    // Generic Azure error
    const message = error?.message || error?.body?.message || 'An Azure operation failed';
    const code = error?.code || 'AZURE_ERROR';
    const statusCode = error?.statusCode || 500;

    throw new AzureResourceError(message, code, statusCode);
  }
}