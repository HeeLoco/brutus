import { Request, Response, NextFunction } from 'express';
import Joi from 'joi';
import { AzureResourceError } from '@/services/azureService';
import { RequestWithCorrelation, ErrorResponse } from '@/types/express';
import { createLogger } from '@/config/logger';

export function errorHandler(
  error: Error,
  req: Request,
  res: Response,
  next: NextFunction
): void {
  const correlationId = (req as RequestWithCorrelation).correlationId || 'unknown';
  const logger = createLogger(correlationId);
  
  let statusCode = 500;
  let code = 'INTERNAL_SERVER_ERROR';
  let message = 'An unexpected error occurred';
  let details: any = undefined;

  // Handle Joi validation errors
  if (Joi.isError(error)) {
    statusCode = 422;
    code = 'VALIDATION_ERROR';
    message = 'Request validation failed';
    details = error.details.map(detail => ({
      field: detail.path.join('.'),
      message: detail.message,
      value: detail.context?.value
    }));
    
    logger.warn('Validation error', { error: error.message, details });
  }
  // Handle Azure-specific errors
  else if (error instanceof AzureResourceError) {
    statusCode = error.statusCode;
    code = error.code;
    message = error.message;
    
    logger.error('Azure resource error', { 
      code: error.code,
      statusCode: error.statusCode,
      message: error.message
    });
  }
  // Handle other known error types
  else if (error.name === 'UnauthorizedError') {
    statusCode = 401;
    code = 'UNAUTHORIZED';
    message = 'Authentication required';
    
    logger.warn('Unauthorized access attempt');
  }
  else if (error.name === 'ForbiddenError') {
    statusCode = 403;
    code = 'FORBIDDEN';
    message = 'Access forbidden';
    
    logger.warn('Forbidden access attempt');
  }
  // Handle unexpected errors
  else {
    logger.error('Unexpected error', { 
      error: error.message,
      stack: error.stack,
      name: error.name
    });
  }

  const errorResponse: ErrorResponse = {
    success: false,
    error: {
      code,
      message,
      details
    },
    correlationId,
    timestamp: new Date().toISOString()
  };

  res.status(statusCode).json(errorResponse);
}

export function notFoundHandler(req: Request, res: Response): void {
  const correlationId = (req as RequestWithCorrelation).correlationId || 'unknown';
  
  const errorResponse: ErrorResponse = {
    success: false,
    error: {
      code: 'NOT_FOUND',
      message: `Route ${req.method} ${req.path} not found`
    },
    correlationId,
    timestamp: new Date().toISOString()
  };

  res.status(404).json(errorResponse);
}