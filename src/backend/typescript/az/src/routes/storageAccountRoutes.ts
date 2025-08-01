import { Router } from 'express';
import { storageAccountController } from '@/controllers/storageAccountController';
import { validateResourceGroupName } from '@/models/resourceGroup';
import { validateStorageAccountName } from '@/models/storageAccount';

const router = Router();

/**
 * @route   GET /api/v1/storage-accounts
 * @desc    List all storage accounts in the subscription
 * @access  Private
 */
router.get('/', storageAccountController.listStorageAccounts);

/**
 * @route   GET /api/v1/storage-accounts/:resourceGroup/:name
 * @desc    Get a specific storage account
 * @access  Private
 */
router.get('/:resourceGroup/:name', storageAccountController.getStorageAccount);

/**
 * @route   POST /api/v1/storage-accounts
 * @desc    Create a new storage account
 * @access  Private
 */
router.post('/', storageAccountController.createStorageAccount);

/**
 * @route   PUT /api/v1/storage-accounts/:resourceGroup/:name
 * @desc    Update an existing storage account
 * @access  Private
 */
router.put('/:resourceGroup/:name', storageAccountController.updateStorageAccount);

/**
 * @route   DELETE /api/v1/storage-accounts/:resourceGroup/:name
 * @desc    Delete a storage account
 * @access  Private
 */
router.delete('/:resourceGroup/:name', storageAccountController.deleteStorageAccount);

export default router;