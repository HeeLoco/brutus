import { Request, Response, NextFunction } from 'express';
import { v4 as uuidv4 } from 'uuid';
import { RequestWithCorrelation } from '@/types/express';

export function correlationMiddleware(req: Request, res: Response, next: NextFunction): void {
  const correlationId = req.headers['x-correlation-id'] as string || uuidv4().substring(0, 8);
  
  // Add correlation ID to request object
  (req as RequestWithCorrelation).correlationId = correlationId;
  
  // Add correlation ID to response headers
  res.setHeader('X-Correlation-ID', correlationId);
  
  next();
}