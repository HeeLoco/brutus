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