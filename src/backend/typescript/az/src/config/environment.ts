import Joi from 'joi';

interface EnvironmentConfig {
  NODE_ENV: 'development' | 'production' | 'test';
  PORT: number;
  HOST: string;
  
  // Azure Configuration
  AZURE_SUBSCRIPTION_ID: string;
  AZURE_TENANT_ID?: string;
  AZURE_CLIENT_ID?: string;
  AZURE_CLIENT_SECRET?: string;
  
  // API Configuration
  API_PREFIX: string;
  CORS_ORIGINS: string[];
  
  // Logging Configuration
  LOG_LEVEL: 'error' | 'warn' | 'info' | 'debug';
  LOG_FORMAT: 'json' | 'simple';
}

const envSchema = Joi.object({
  NODE_ENV: Joi.string()
    .valid('development', 'production', 'test')
    .default('development'),
  PORT: Joi.number().port().default(8000),
  HOST: Joi.string().default('0.0.0.0'),
  
  // Azure Configuration (required)
  AZURE_SUBSCRIPTION_ID: Joi.string().required(),
  AZURE_TENANT_ID: Joi.string().optional(),
  AZURE_CLIENT_ID: Joi.string().optional(),
  AZURE_CLIENT_SECRET: Joi.string().optional(),
  
  // API Configuration
  API_PREFIX: Joi.string().default('/api/v1'),
  CORS_ORIGINS: Joi.string().default('*'),
  
  // Logging Configuration
  LOG_LEVEL: Joi.string()
    .valid('error', 'warn', 'info', 'debug')
    .default('info'),
  LOG_FORMAT: Joi.string()
    .valid('json', 'simple')
    .default('simple')
}).unknown();

function validateEnvironment(): EnvironmentConfig {
  const { error, value } = envSchema.validate(process.env);
  
  if (error) {
    throw new Error(`Environment validation error: ${error.message}`);
  }
  
  // Parse CORS origins
  const corsOrigins = value.CORS_ORIGINS === '*' 
    ? ['*'] 
    : (value.CORS_ORIGINS as string).split(',').map((origin: string) => origin.trim());
  
  return {
    ...value,
    CORS_ORIGINS: corsOrigins
  } as EnvironmentConfig;
}

export const config = validateEnvironment();

export default config;