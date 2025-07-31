// Input validation and sanitization utilities

export class InputValidator {
  // UUID/GUID validation for Azure subscription IDs
  static validateSubscriptionId(id: string): boolean {
    const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[4-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i
    return uuidRegex.test(id.trim())
  }

  // Resource group name validation (Azure rules)
  static validateResourceGroupName(name: string): boolean {
    // Azure resource group naming rules:
    // - 1-90 characters
    // - Alphanumeric, underscore, parentheses, hyphen, period
    // - Cannot end with period
    const rgNameRegex = /^[a-zA-Z0-9_\-\(\)\.]{1,90}$/
    const trimmedName = name.trim()
    
    if (!rgNameRegex.test(trimmedName) || trimmedName.endsWith('.')) {
      return false
    }
    
    return true
  }

  // Azure location validation
  static validateAzureLocation(location: string): boolean {
    const validLocations = [
      'eastus', 'eastus2', 'southcentralus', 'westus2', 'westus3', 'australiaeast',
      'southeastasia', 'northeurope', 'swedencentral', 'uksouth', 'westeurope',
      'centralus', 'southafricanorth', 'centralindia', 'eastasia', 'japaneast',
      'koreacentral', 'canadacentral', 'francecentral', 'germanywestcentral',
      'norwayeast', 'switzerlandnorth', 'uaenorth', 'brazilsouth', 'qatarcentral'
    ]
    
    return validLocations.includes(location.toLowerCase().trim())
  }

  // General string sanitization
  static sanitizeString(input: string, maxLength: number = 255): string {
    return input
      .trim()
      .replace(/[<>'"&]/g, '') // Remove potentially dangerous characters
      .substring(0, maxLength)
  }

  // Tag key validation
  static validateTagKey(key: string): boolean {
    // Azure tag key rules: 1-512 characters, no < > % & \ ? /
    const tagKeyRegex = /^[^<>%&\\?\/]{1,512}$/
    return tagKeyRegex.test(key.trim())
  }

  // Tag value validation  
  static validateTagValue(value: string): boolean {
    // Azure tag value rules: 0-256 characters, no < > % & \ ? /
    const tagValueRegex = /^[^<>%&\\?\/]{0,256}$/
    return tagValueRegex.test(value.trim())
  }

  // Validate entire tag object
  static validateTags(tags: Record<string, string>): { valid: boolean; errors: string[] } {
    const errors: string[] = []
    const entries = Object.entries(tags)

    // Azure allows up to 50 tags per resource
    if (entries.length > 50) {
      errors.push('Maximum 50 tags allowed')
    }

    entries.forEach(([key, value]) => {
      if (!this.validateTagKey(key)) {
        errors.push(`Invalid tag key: ${key}`)
      }
      if (!this.validateTagValue(value)) {
        errors.push(`Invalid tag value for key ${key}: ${value}`)
      }
    })

    return {
      valid: errors.length === 0,
      errors
    }
  }
}

// User-friendly error messages
export class ValidationMessages {
  static getSubscriptionIdError(): string {
    return 'Please enter a valid Azure subscription ID (GUID format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)'
  }

  static getResourceGroupNameError(): string {
    return 'Resource group name must be 1-90 characters, contain only alphanumeric characters, underscores, parentheses, hyphens, and periods, and cannot end with a period'
  }

  static getLocationError(): string {
    return 'Please select a valid Azure region'
  }

  static getTagError(errors: string[]): string {
    return `Tag validation failed: ${errors.join(', ')}`
  }
}

// Form validation composable
export function useFormValidation() {
  const errors = ref<Record<string, string>>({})

  const validateField = (field: string, value: string, validator: (val: string) => boolean, errorMessage: string) => {
    if (!validator(value)) {
      errors.value[field] = errorMessage
      return false
    } else {
      delete errors.value[field]
      return true
    }
  }

  const clearErrors = () => {
    errors.value = {}
  }

  const hasErrors = computed(() => Object.keys(errors.value).length > 0)

  return {
    errors: readonly(errors),
    validateField,
    clearErrors,
    hasErrors
  }
}

// Import ref, computed, readonly for composable
import { ref, computed, readonly } from 'vue'