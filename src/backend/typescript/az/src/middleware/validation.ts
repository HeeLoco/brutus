import { Request, Response, NextFunction } from 'express';
import Joi from 'joi';
import { RequestWithCorrelation } from '@/types/express';

export function validateBody(schema: Joi.ObjectSchema) {
  return (req: Request, res: Response, next: NextFunction): void => {
    const { error, value } = schema.validate(req.body);
    
    if (error) {
      return next(error);
    }
    
    req.body = value;
    next();
  };
}

export function validateParams(schema: Joi.ObjectSchema) {
  return (req: Request, res: Response, next: NextFunction): void => {
    const { error, value } = schema.validate(req.params);
    
    if (error) {
      return next(error);
    }
    
    req.params = value;
    next();
  };
}

export function validateQuery(schema: Joi.ObjectSchema) {
  return (req: Request, res: Response, next: NextFunction): void => {
    const { error, value } = schema.validate(req.query);
    
    if (error) {
      return next(error);
    }
    
    req.query = value;
    next();
  };
}