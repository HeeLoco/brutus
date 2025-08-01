import { Request, Response } from 'express';
import { AzureService, AzureResourceError } from '@/services/azureService';
import { 
  validateCreateStorageAccount, 
  validateUpdateStorageAccount, 
  validateStorageAccountName 
} from '@/models/storageAccount';
import { createLogger } from '@/config/logger';
import { asyncHandler } from '@/utils/asyncHandler';

const logger = createLogger();

export const storageAccountController = {
  /**
   * List all storage accounts
   */
  listStorageAccounts: asyncHandler(async (req: Request, res: Response): Promise<void> => {
    const correlationId = req.correlationId;
    logger.info('Controller: Listing storage accounts', { correlationId });

    const azureService = new AzureService(correlationId);
    const storageAccounts = await azureService.listStorageAccounts();

    res.status(200).json({
      success: true,
      data: {
        storageAccounts,
        count: storageAccounts.length
      },
      correlationId
    });
  }),

  /**
   * Get a specific storage account
   */
  getStorageAccount: asyncHandler(async (req: Request, res: Response): Promise<void> => {
    const correlationId = req.correlationId;
    const { resourceGroup, name } = req.params;

    logger.info(`Controller: Getting storage account ${name} in resource group ${resourceGroup}`, { 
      correlationId, 
      resourceGroup, 
      name 
    });

    // Validate storage account name
    const { error: nameError } = validateStorageAccountName(name);
    if (nameError) {
      res.status(400).json({
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Invalid storage account name',
          details: nameError.details.map(detail => ({
            field: detail.path?.join('.') || 'name',
            message: detail.message
          }))
        },
        correlationId
      });
      return;
    }

    const azureService = new AzureService(correlationId);
    const storageAccount = await azureService.getStorageAccount(resourceGroup, name);

    res.status(200).json({
      success: true,
      data: storageAccount,
      correlationId
    });
  }),

  /**
   * Create a new storage account
   */
  createStorageAccount: asyncHandler(async (req: Request, res: Response): Promise<void> => {
    const correlationId = req.correlationId;
    logger.info('Controller: Creating storage account', { correlationId, body: req.body });

    // Validate request body
    const { error, value } = validateCreateStorageAccount(req.body);
    if (error) {
      res.status(400).json({
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Invalid storage account configuration',
          details: error.details.map(detail => ({
            field: detail.path?.join('.') || 'unknown',
            message: detail.message
          }))
        },
        correlationId
      });
      return;
    }

    const azureService = new AzureService(correlationId);
    const storageAccount = await azureService.createStorageAccount(value!);

    logger.info(`Controller: Successfully created storage account ${value!.name}`, { 
      correlationId, 
      name: value!.name 
    });

    res.status(201).json({
      success: true,
      data: storageAccount,
      correlationId
    });
  }),

  /**
   * Update an existing storage account
   */
  updateStorageAccount: asyncHandler(async (req: Request, res: Response): Promise<void> => {
    const correlationId = req.correlationId;
    const { resourceGroup, name } = req.params;

    logger.info(`Controller: Updating storage account ${name} in resource group ${resourceGroup}`, { 
      correlationId, 
      resourceGroup, 
      name, 
      body: req.body 
    });

    // Validate storage account name
    const { error: nameError } = validateStorageAccountName(name);
    if (nameError) {
      res.status(400).json({
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Invalid storage account name',
          details: nameError.details.map(detail => ({
            field: detail.path?.join('.') || 'name',
            message: detail.message
          }))
        },
        correlationId
      });
      return;
    }

    // Validate request body
    const { error, value } = validateUpdateStorageAccount(req.body);
    if (error) {
      res.status(400).json({
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Invalid update configuration',
          details: error.details.map(detail => ({
            field: detail.path?.join('.') || 'unknown',
            message: detail.message
          }))
        },
        correlationId
      });
      return;
    }

    const azureService = new AzureService(correlationId);
    const storageAccount = await azureService.updateStorageAccount(resourceGroup, name, value!);

    logger.info(`Controller: Successfully updated storage account ${name}`, { 
      correlationId, 
      name 
    });

    res.status(200).json({
      success: true,
      data: storageAccount,
      correlationId
    });
  }),

  /**
   * Delete a storage account
   */
  deleteStorageAccount: asyncHandler(async (req: Request, res: Response): Promise<void> => {
    const correlationId = req.correlationId;
    const { resourceGroup, name } = req.params;

    logger.info(`Controller: Deleting storage account ${name} in resource group ${resourceGroup}`, { 
      correlationId, 
      resourceGroup, 
      name 
    });

    // Validate storage account name
    const { error: nameError } = validateStorageAccountName(name);
    if (nameError) {
      res.status(400).json({
        success: false,
        error: {
          code: 'VALIDATION_ERROR',
          message: 'Invalid storage account name',
          details: nameError.details.map(detail => ({
            field: detail.path?.join('.') || 'name',
            message: detail.message
          }))
        },
        correlationId
      });
      return;
    }

    const azureService = new AzureService(correlationId);
    await azureService.deleteStorageAccount(resourceGroup, name);

    logger.info(`Controller: Successfully deleted storage account ${name}`, { 
      correlationId, 
      name 
    });

    res.status(200).json({
      success: true,
      message: `Storage account '${name}' deleted successfully`,
      data: {
        resourceGroup,
        name
      },
      correlationId
    });
  })
};