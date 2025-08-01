export interface ResourceGroup {
  id: string;
  name: string;
  location: string;
  type: string;
  tags?: Record<string, string>;
  managedBy?: string;
  provisioningState?: string;
}

export interface CreateResourceGroupRequest {
  name: string;
  location: string;
  tags?: Record<string, string>;
}

export interface UpdateResourceGroupRequest {
  tags?: Record<string, string>;
}

export interface ResourceGroupList {
  value: ResourceGroup[];
  nextLink?: string;
}

export interface AzureErrorResponse {
  error: {
    code: string;
    message: string;
    target?: string;
    details?: Array<{
      code: string;
      message: string;
      target?: string;
    }>;
  };
}

export interface HealthCheckResponse {
  status: 'healthy' | 'degraded' | 'unhealthy';
  timestamp: string;
  version: string;
  azureConnection: boolean;
  uptime: number;
  checks: {
    azure: {
      status: 'pass' | 'fail';
      responseTime?: number;
      error?: string;
    };
  };
}

// Storage Account Types
export interface StorageEndpoints {
  blob?: string;
  queue?: string;
  table?: string;
  file?: string;
}

export interface StorageAccount {
  id: string;
  name: string;
  location: string;
  resourceGroup: string;
  kind: string;
  skuName: string;
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

export interface StorageAccountList {
  value: StorageAccount[];
  nextLink?: string;
}