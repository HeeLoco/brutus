import express, { Application, Request, Response } from 'express';
import cors from 'cors';
import helmet from 'helmet';
import morgan from 'morgan';
import rateLimit from 'express-rate-limit';

import config from '@/config/environment';
import { loggerStream } from '@/config/logger';
import { correlationMiddleware } from '@/middleware/correlation';
import { errorHandler, notFoundHandler } from '@/middleware/errorHandler';
import { RequestWithCorrelation, ApiResponse } from '@/types/express';

import healthRoutes from '@/routes/healthRoutes';
import resourceGroupRoutes from '@/routes/resourceGroupRoutes';
import storageAccountRoutes from '@/routes/storageAccountRoutes';

class App {
  public app: Application;

  constructor() {
    this.app = express();
    this.initializeMiddlewares();
    this.initializeRoutes();
    this.initializeErrorHandling();
  }

  private initializeMiddlewares(): void {
    // Security middleware
    this.app.use(helmet());

    // Rate limiting
    const limiter = rateLimit({
      windowMs: 15 * 60 * 1000, // 15 minutes
      max: config.NODE_ENV === 'production' ? 100 : 1000, // Limit each IP
      message: {
        error: 'Too many requests from this IP, please try again later.',
      },
      standardHeaders: true,
      legacyHeaders: false,
    });
    this.app.use(limiter);

    // CORS
    this.app.use(cors({
      origin: config.CORS_ORIGINS,
      methods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
      allowedHeaders: ['Content-Type', 'Authorization', 'X-Correlation-ID'],
      exposedHeaders: ['X-Correlation-ID', 'X-RateLimit-Limit', 'X-RateLimit-Remaining'],
      credentials: true
    }));

    // Body parsing
    this.app.use(express.json({ limit: '10mb' }));
    this.app.use(express.urlencoded({ extended: true, limit: '10mb' }));

    // Logging
    const morganFormat = config.NODE_ENV === 'production' 
      ? 'combined' 
      : 'dev';
    this.app.use(morgan(morganFormat, { stream: loggerStream }));

    // Correlation ID middleware
    this.app.use(correlationMiddleware);
  }

  private initializeRoutes(): void {
    // Root endpoint
    this.app.get('/', (req: Request, res: Response) => {
      const correlationId = (req as RequestWithCorrelation).correlationId;
      
      const response: ApiResponse<{
        name: string;
        version: string;
        status: string;
        environment: string;
        endpoints: Record<string, string>;
      }> = {
        success: true,
        data: {
          name: 'Azure Resource Management API (TypeScript)',
          version: '1.0.0',
          status: 'running',
          environment: config.NODE_ENV,
          endpoints: {
            health: '/health',
            resourceGroups: `${config.API_PREFIX}/resource-groups`,
            storageAccounts: `${config.API_PREFIX}/storage-accounts`,
            docs: '/docs' // Future: OpenAPI documentation
          }
        },
        correlationId,
        timestamp: new Date().toISOString()
      };
      
      res.status(200).json(response);
    });

    // Health routes (not prefixed)
    this.app.use('/', healthRoutes);

    // API routes
    this.app.use(`${config.API_PREFIX}/resource-groups`, resourceGroupRoutes);
    this.app.use(`${config.API_PREFIX}/storage-accounts`, storageAccountRoutes);

    // API info endpoint
    this.app.get(`${config.API_PREFIX}`, (req: Request, res: Response) => {
      const correlationId = (req as RequestWithCorrelation).correlationId;
      
      const response: ApiResponse<{
        version: string;
        endpoints: string[];
      }> = {
        success: true,
        data: {
          version: 'v1',
          endpoints: [
            'GET /resource-groups - List resource groups',
            'POST /resource-groups - Create resource group',
            'GET /resource-groups/{name} - Get resource group',
            'PUT /resource-groups/{name} - Update resource group',
            'DELETE /resource-groups/{name} - Delete resource group',
            'GET /storage-accounts - List storage accounts',
            'POST /storage-accounts - Create storage account',
            'GET /storage-accounts/{resourceGroup}/{name} - Get storage account',
            'PUT /storage-accounts/{resourceGroup}/{name} - Update storage account',
            'DELETE /storage-accounts/{resourceGroup}/{name} - Delete storage account'
          ]
        },
        correlationId,
        timestamp: new Date().toISOString()
      };
      
      res.status(200).json(response);
    });
  }

  private initializeErrorHandling(): void {
    // 404 handler
    this.app.use('*', notFoundHandler);

    // Global error handler
    this.app.use(errorHandler);
  }

  public getApp(): Application {
    return this.app;
  }
}

export default App;