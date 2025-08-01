import Joi from 'joi';
import { CreateStorageAccountRequest, UpdateStorageAccountRequest } from '@/types/azure';

// Storage account kind constants
export const STORAGE_KINDS = [
  'Storage',
  'StorageV2',
  'BlobStorage',
  'FileStorage',
  'BlockBlobStorage'
] as const;

export type StorageKind = typeof STORAGE_KINDS[number];

// Storage account SKU constants
export const STORAGE_SKUS = [
  'Standard_LRS',
  'Standard_GRS',
  'Standard_RAGRS',
  'Standard_ZRS',
  'Premium_LRS',
  'Premium_ZRS',
  'Standard_GZRS',
  'Standard_RAGZRS'
] as const;

export type StorageSku = typeof STORAGE_SKUS[number];

// Access tier constants
export const ACCESS_TIERS = [
  'Hot',
  'Cool',
  'Archive'
] as const;

export type AccessTier = typeof ACCESS_TIERS[number];

// Validation schemas
export const createStorageAccountSchema = Joi.object<CreateStorageAccountRequest>({
  name: Joi.string()
    .min(3)
    .max(24)
    .pattern(/^[a-z0-9]+$/)
    .required()
    .messages({
      'string.pattern.base': 'Storage account name can only contain lowercase letters and numbers',
      'string.min': 'Storage account name must be at least 3 characters',
      'string.max': 'Storage account name must be 24 characters or less',
      'any.required': 'Storage account name is required'
    }),
  
  resourceGroup: Joi.string()
    .min(1)
    .max(90)
    .required()
    .messages({
      'any.required': 'Resource group is required'
    }),
  
  location: Joi.string()
    .min(1)
    .required()
    .messages({
      'any.required': 'Location is required'
    }),
  
  kind: Joi.string()
    .valid(...STORAGE_KINDS)
    .optional()
    .default('StorageV2')
    .messages({
      'any.only': `Storage kind must be one of: ${STORAGE_KINDS.join(', ')}`
    }),
  
  skuName: Joi.string()
    .valid(...STORAGE_SKUS)
    .optional()
    .default('Standard_LRS')
    .messages({
      'any.only': `SKU name must be one of: ${STORAGE_SKUS.join(', ')}`
    }),
  
  accessTier: Joi.string()
    .valid(...ACCESS_TIERS)
    .optional()
    .messages({
      'any.only': `Access tier must be one of: ${ACCESS_TIERS.join(', ')}`
    }),
  
  allowBlobPublicAccess: Joi.boolean()
    .optional()
    .default(false),
  
  allowSharedKeyAccess: Joi.boolean()
    .optional()
    .default(true),
  
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

export const updateStorageAccountSchema = Joi.object<UpdateStorageAccountRequest>({
  accessTier: Joi.string()
    .valid(...ACCESS_TIERS)
    .optional()
    .messages({
      'any.only': `Access tier must be one of: ${ACCESS_TIERS.join(', ')}`
    }),
  
  allowBlobPublicAccess: Joi.boolean()
    .optional(),
  
  allowSharedKeyAccess: Joi.boolean()
    .optional(),
  
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

export const storageAccountNameSchema = Joi.string()
  .min(3)
  .max(24)
  .pattern(/^[a-z0-9]+$/)
  .required()
  .messages({
    'string.pattern.base': 'Storage account name can only contain lowercase letters and numbers',
    'string.min': 'Storage account name must be at least 3 characters',
    'string.max': 'Storage account name must be 24 characters or less',
    'any.required': 'Storage account name is required'
  });

// Helper functions for validation
export function validateCreateStorageAccount(data: unknown): { error?: Joi.ValidationError; value?: CreateStorageAccountRequest } {
  return createStorageAccountSchema.validate(data);
}

export function validateUpdateStorageAccount(data: unknown): { error?: Joi.ValidationError; value?: UpdateStorageAccountRequest } {
  return updateStorageAccountSchema.validate(data);
}

export function validateStorageAccountName(name: unknown): { error?: Joi.ValidationError; value?: string } {
  return storageAccountNameSchema.validate(name);
}