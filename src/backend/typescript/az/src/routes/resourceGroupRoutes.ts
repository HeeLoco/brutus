import { Router } from 'express';
import Joi from 'joi';
import { ResourceGroupController } from '@/controllers/resourceGroupController';
import { validateBody, validateParams } from '@/middleware/validation';
import { asyncHandler } from '@/utils/asyncHandler';
import { 
  createResourceGroupSchema, 
  updateResourceGroupSchema, 
  resourceGroupNameSchema 
} from '@/models/resourceGroup';

const router = Router();
const resourceGroupController = new ResourceGroupController();

// Parameter validation schema
const nameParamSchema = Joi.object({
  name: resourceGroupNameSchema
});

/**
 * @route   GET /api/v1/resource-groups
 * @desc    List all resource groups
 * @access  Public
 */
router.get(
  '/',
  asyncHandler(resourceGroupController.listResourceGroups.bind(resourceGroupController))
);

/**
 * @route   GET /api/v1/resource-groups/:name
 * @desc    Get a specific resource group
 * @access  Public
 */
router.get(
  '/:name',
  validateParams(nameParamSchema),
  asyncHandler(resourceGroupController.getResourceGroup.bind(resourceGroupController))
);

/**
 * @route   POST /api/v1/resource-groups
 * @desc    Create a new resource group
 * @access  Public
 */
router.post(
  '/',
  validateBody(createResourceGroupSchema),
  asyncHandler(resourceGroupController.createResourceGroup.bind(resourceGroupController))
);

/**
 * @route   PUT /api/v1/resource-groups/:name
 * @desc    Update a resource group
 * @access  Public
 */
router.put(
  '/:name',
  validateParams(nameParamSchema),
  validateBody(updateResourceGroupSchema),
  asyncHandler(resourceGroupController.updateResourceGroup.bind(resourceGroupController))
);

/**
 * @route   DELETE /api/v1/resource-groups/:name
 * @desc    Delete a resource group
 * @access  Public
 */
router.delete(
  '/:name',
  validateParams(nameParamSchema),
  asyncHandler(resourceGroupController.deleteResourceGroup.bind(resourceGroupController))
);

export default router;