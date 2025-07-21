import { Request } from 'express';

export interface RequestWithCorrelation extends Request {
  correlationId: string;
}

export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
    details?: any;
  };
  correlationId: string;
  timestamp: string;
}

export interface PaginatedResponse<T = any> extends ApiResponse<T> {
  pagination?: {
    page: number;
    limit: number;
    total: number;
    hasNext: boolean;
    hasPrev: boolean;
  };
}

export interface ValidationError {
  field: string;
  message: string;
  value?: any;
}

export interface ErrorResponse extends ApiResponse {
  success: false;
  error: {
    code: string;
    message: string;
    details?: ValidationError[];
  };
}