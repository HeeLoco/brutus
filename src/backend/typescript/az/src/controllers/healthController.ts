import { Request, Response, NextFunction } from 'express';
import { AzureService } from '@/services/azureService';
import { RequestWithCorrelation, ApiResponse } from '@/types/express';
import { HealthCheckResponse } from '@/types/azure';
import { createLogger } from '@/config/logger';
import config from '@/config/environment';

export class HealthController {
  private startTime = Date.now();

  async checkHealth(req: Request, res: Response, next: NextFunction): Promise<void> {
    const correlationId = (req as RequestWithCorrelation).correlationId;
    const logger = createLogger(correlationId);
    
    try {
      logger.info('Performing health check');
      
      // Test Azure connection
      const azureService = new AzureService(correlationId);
      const startTime = Date.now();
      const azureConnection = await azureService.testConnection();
      const azureResponseTime = Date.now() - startTime;
      
      const uptime = Date.now() - this.startTime;
      
      let status: 'healthy' | 'degraded' | 'unhealthy';
      if (azureConnection) {
        status = 'healthy';
      } else {
        status = 'unhealthy';
      }
      
      const healthCheck: HealthCheckResponse = {
        status,
        timestamp: new Date().toISOString(),
        version: '1.0.0',
        azureConnection,
        uptime,
        checks: {
          azure: {
            status: azureConnection ? 'pass' : 'fail',
            responseTime: azureResponseTime,
            error: azureConnection ? undefined : 'Unable to connect to Azure services'
          }
        }
      };
      
      const response: ApiResponse<HealthCheckResponse> = {
        success: true,
        data: healthCheck,
        correlationId,
        timestamp: new Date().toISOString()
      };
      
      const statusCode = status === 'healthy' ? 200 : (status === 'degraded' ? 200 : 503);
      logger.info(`Health check completed - Status: ${status}`, { azureConnection, responseTime: azureResponseTime });
      
      res.status(statusCode).json(response);
      
    } catch (error) {
      logger.error('Health check failed', { error });
      next(error);
    }
  }
}