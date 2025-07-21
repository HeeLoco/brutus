import { Request, Response, NextFunction } from 'express';
import { AzureService } from '@/services/azureService';
import { RequestWithCorrelation, ApiResponse, PaginatedResponse } from '@/types/express';
import { ResourceGroup, CreateResourceGroupRequest, UpdateResourceGroupRequest } from '@/types/azure';
import { createLogger } from '@/config/logger';

export class ResourceGroupController {
  
  async listResourceGroups(req: Request, res: Response, next: NextFunction): Promise<void> {
    const correlationId = (req as RequestWithCorrelation).correlationId;
    const logger = createLogger(correlationId);
    
    try {
      logger.info('Listing resource groups');
      
      const azureService = new AzureService(correlationId);
      const resourceGroups = await azureService.listResourceGroups();
      
      const response: ApiResponse<{ resourceGroups: ResourceGroup[]; count: number }> = {
        success: true,
        data: {
          resourceGroups,
          count: resourceGroups.length
        },
        correlationId,
        timestamp: new Date().toISOString()
      };
      
      logger.info(`Successfully listed ${resourceGroups.length} resource groups`);
      res.status(200).json(response);
      
    } catch (error) {
      logger.error('Failed to list resource groups', { error });
      next(error);
    }
  }
  
  async getResourceGroup(req: Request, res: Response, next: NextFunction): Promise<void> {
    const correlationId = (req as RequestWithCorrelation).correlationId;
    const logger = createLogger(correlationId);
    const { name } = req.params;
    
    try {
      logger.info(`Getting resource group: ${name}`);
      
      const azureService = new AzureService(correlationId);
      const resourceGroup = await azureService.getResourceGroup(name);
      
      const response: ApiResponse<ResourceGroup> = {
        success: true,
        data: resourceGroup,
        correlationId,
        timestamp: new Date().toISOString()
      };
      
      logger.info(`Successfully retrieved resource group: ${name}`);
      res.status(200).json(response);
      
    } catch (error) {
      logger.error(`Failed to get resource group ${name}`, { error });
      next(error);
    }
  }
  
  async createResourceGroup(req: Request, res: Response, next: NextFunction): Promise<void> {
    const correlationId = (req as RequestWithCorrelation).correlationId;
    const logger = createLogger(correlationId);
    const createRequest = req.body as CreateResourceGroupRequest;
    
    try {
      logger.info(`Creating resource group: ${createRequest.name}`);
      
      const azureService = new AzureService(correlationId);
      const resourceGroup = await azureService.createResourceGroup(createRequest);
      
      const response: ApiResponse<ResourceGroup> = {
        success: true,
        data: resourceGroup,
        correlationId,
        timestamp: new Date().toISOString()
      };
      
      logger.info(`Successfully created resource group: ${createRequest.name}`);
      res.status(201).json(response);
      
    } catch (error) {
      logger.error(`Failed to create resource group ${createRequest.name}`, { error });
      next(error);
    }
  }
  
  async updateResourceGroup(req: Request, res: Response, next: NextFunction): Promise<void> {
    const correlationId = (req as RequestWithCorrelation).correlationId;
    const logger = createLogger(correlationId);
    const { name } = req.params;
    const updateRequest = req.body as UpdateResourceGroupRequest;
    
    try {
      logger.info(`Updating resource group: ${name}`);
      
      const azureService = new AzureService(correlationId);
      const resourceGroup = await azureService.updateResourceGroup(name, updateRequest);
      
      const response: ApiResponse<ResourceGroup> = {
        success: true,
        data: resourceGroup,
        correlationId,
        timestamp: new Date().toISOString()
      };
      
      logger.info(`Successfully updated resource group: ${name}`);
      res.status(200).json(response);
      
    } catch (error) {
      logger.error(`Failed to update resource group ${name}`, { error });
      next(error);
    }
  }
  
  async deleteResourceGroup(req: Request, res: Response, next: NextFunction): Promise<void> {
    const correlationId = (req as RequestWithCorrelation).correlationId;
    const logger = createLogger(correlationId);
    const { name } = req.params;
    
    try {
      logger.info(`Deleting resource group: ${name}`);
      
      const azureService = new AzureService(correlationId);
      await azureService.deleteResourceGroup(name);
      
      const response: ApiResponse<{ message: string; name: string }> = {
        success: true,
        data: {
          message: `Resource group '${name}' deletion initiated`,
          name
        },
        correlationId,
        timestamp: new Date().toISOString()
      };
      
      logger.info(`Successfully initiated deletion of resource group: ${name}`);
      res.status(202).json(response);
      
    } catch (error) {
      logger.error(`Failed to delete resource group ${name}`, { error });
      next(error);
    }
  }
}