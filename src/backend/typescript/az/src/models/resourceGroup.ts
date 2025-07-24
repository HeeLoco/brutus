import Joi from 'joi';
import { CreateResourceGroupRequest, UpdateResourceGroupRequest } from '@/types/azure';

// Validation schemas
export const createResourceGroupSchema = Joi.object<CreateResourceGroupRequest>({
  name: Joi.string()
    .min(1)
    .max(90)
    .pattern(/^[-\w\._\(\)]+$/)
    .required()
    .messages({
      'string.pattern.base': 'Resource group name can only contain alphanumeric characters, periods, underscores, hyphens, and parentheses',
      'string.max': 'Resource group name must be 90 characters or less',
      'any.required': 'Resource group name is required'
    }),
  
  location: Joi.string()
    .min(1)
    .required()
    .messages({
      'any.required': 'Location is required'
    }),
  
  tags: Joi.object()
    .pattern(
      Joi.string().max(512),
      Joi.string().max(256)
    )
    .max(50)
    .optional()
    .messages({
      'object.max': 'Maximum of 50 tags allowed',
      'string.max': 'Tag keys must be 512 characters or less, tag values must be 256 characters or less'
    })
});

export const updateResourceGroupSchema = Joi.object<UpdateResourceGroupRequest>({
  tags: Joi.object()
    .pattern(
      Joi.string().max(512),
      Joi.string().max(256)
    )
    .max(50)
    .optional()
    .messages({
      'object.max': 'Maximum of 50 tags allowed',
      'string.max': 'Tag keys must be 512 characters or less, tag values must be 256 characters or less'
    })
});

export const resourceGroupNameSchema = Joi.string()
  .min(1)
  .max(90)
  .pattern(/^[-\w\._\(\)]+$/)
  .required()
  .messages({
    'string.pattern.base': 'Resource group name can only contain alphanumeric characters, periods, underscores, hyphens, and parentheses',
    'string.max': 'Resource group name must be 90 characters or less',
    'any.required': 'Resource group name is required'
  });

// Azure location constants
export const AZURE_LOCATIONS = [
  'eastus',
  'eastus2',
  'southcentralus',
  'westus2',
  'westus3',
  'australiaeast',
  'southeastasia',
  'northeurope',
  'swedencentral',
  'uksouth',
  'westeurope',
  'centralus',
  'southafricanorth',
  'centralindia',
  'eastasia',
  'japaneast',
  'koreacentral',
  'canadacentral',
  'francecentral',
  'germanywestcentral',
  'norwayeast',
  'brazilsouth',
  'uaenorth'
] as const;

export type AzureLocation = typeof AZURE_LOCATIONS[number];

// Helper functions for validation
export function validateCreateResourceGroup(data: unknown): { error?: Joi.ValidationError; value?: CreateResourceGroupRequest } {
  return createResourceGroupSchema.validate(data);
}

export function validateUpdateResourceGroup(data: unknown): { error?: Joi.ValidationError; value?: UpdateResourceGroupRequest } {
  return updateResourceGroupSchema.validate(data);
}

export function validateResourceGroupName(name: unknown): { error?: Joi.ValidationError; value?: string } {
  return resourceGroupNameSchema.validate(name);
}