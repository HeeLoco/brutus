# TypeScript Azure API

Production-ready Express.js/TypeScript API server for managing Azure resources with enterprise-grade features, comprehensive error handling, and modern development practices.

## Features

- **Enterprise Architecture**: Clean layered architecture with controllers, services, models, and middleware
- **Type Safety**: Full TypeScript implementation with strict type checking
- **Azure Integration**: Comprehensive Azure Resource Management with multiple authentication methods
- **Configuration Management**: Environment-based configuration with Joi validation
- **Structured Logging**: Winston-based logging with correlation ID tracking
- **Error Handling**: Comprehensive error handling with proper HTTP status codes
- **Security**: Helmet, rate limiting, CORS, and input validation
- **Validation**: Joi-based request/response validation
- **Performance**: Async/await throughout, efficient middleware chain
- **Developer Experience**: Hot reload, ESLint, path aliases, comprehensive npm scripts

## Project Structure

```
src/backend/typescript/az/
├── src/
│   ├── config/
│   │   ├── environment.ts      # Environment configuration with validation
│   │   └── logger.ts           # Winston logging configuration
│   ├── controllers/
│   │   ├── healthController.ts  # Health check controller
│   │   └── resourceGroupController.ts # Resource group operations
│   ├── middleware/
│   │   ├── correlation.ts      # Request correlation tracking
│   │   ├── errorHandler.ts     # Global error handling
│   │   └── validation.ts       # Request validation middleware
│   ├── models/
│   │   └── resourceGroup.ts    # Joi validation schemas
│   ├── routes/
│   │   ├── healthRoutes.ts     # Health check routes
│   │   └── resourceGroupRoutes.ts # Resource group routes
│   ├── services/
│   │   └── azureService.ts     # Azure SDK integration
│   ├── types/
│   │   ├── azure.ts           # Azure-specific type definitions
│   │   └── express.ts         # Express-specific type definitions
│   ├── utils/
│   │   └── asyncHandler.ts    # Async error handling utility
│   ├── app.ts                 # Express application setup
│   └── server.ts              # Server entry point
├── dist/                      # Compiled JavaScript output
├── package.json
├── tsconfig.json
├── .eslintrc.js
├── nodemon.json
├── .env.example
└── README.typescript.md
```

## Environment Variables

Create a `.env` file from `.env.example`:

```bash
# Environment
NODE_ENV=development

# Server Configuration
HOST=0.0.0.0
PORT=8000

# Azure Configuration (Required)
AZURE_SUBSCRIPTION_ID=your-subscription-id

# Azure Authentication (Optional)
AZURE_TENANT_ID=your-tenant-id
AZURE_CLIENT_ID=your-client-id
AZURE_CLIENT_SECRET=your-client-secret

# API Configuration
API_PREFIX=/api/v1
CORS_ORIGINS=*

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=simple
```

## Development Setup

### Prerequisites
- Node.js >= 18.0.0
- npm >= 9.0.0

### Installation
```bash
cd src/backend/typescript/az
npm install
```

### Development Scripts
```bash
# Development with hot reload
npm run dev

# Build TypeScript
npm run build

# Build and watch for changes
npm run build:watch

# Start production server
npm start

# Type checking only
npm run type-check

# Linting
npm run lint
npm run lint:fix

# Testing
npm test
npm run test:watch
npm run test:coverage

# Clean build directory
npm run clean
```

## Authentication

The API supports multiple Azure authentication methods:

1. **Service Principal** (recommended for production)
   - Set `AZURE_TENANT_ID`, `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`
2. **Managed Identity** (for Azure-hosted applications)
3. **Azure CLI** (for development)
4. **Visual Studio Code** (for development)

If service principal credentials are not provided, the API automatically falls back to Azure DefaultAzureCredential.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | API information and status |
| GET | `/health` | Health check with Azure connectivity |
| GET | `/api/v1` | API version information |
| GET | `/api/v1/resource-groups` | List all resource groups |
| POST | `/api/v1/resource-groups` | Create resource group |
| GET | `/api/v1/resource-groups/{name}` | Get specific resource group |
| PUT | `/api/v1/resource-groups/{name}` | Update resource group |
| DELETE | `/api/v1/resource-groups/{name}` | Delete resource group |

## API Response Format

All API responses follow a consistent format:

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "correlationId": "a1b2c3d4",
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Request validation failed",
    "details": [
      {
        "field": "name",
        "message": "Resource group name is required",
        "value": null
      }
    ]
  },
  "correlationId": "a1b2c3d4",
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

## Request/Response Examples

### Create Resource Group
```bash
POST /api/v1/resource-groups
Content-Type: application/json
X-Correlation-ID: my-request-id

{
  "name": "my-resource-group",
  "location": "eastus",
  "tags": {
    "environment": "dev",
    "project": "brutus"
  }
}
```

### Health Check Response
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-15T10:30:00.000Z",
    "version": "1.0.0",
    "azureConnection": true,
    "uptime": 12345678,
    "checks": {
      "azure": {
        "status": "pass",
        "responseTime": 150
      }
    }
  },
  "correlationId": "a1b2c3d4",
  "timestamp": "2024-01-15T10:30:00.000Z"
}
```

## Security Features

- **Helmet**: Security headers
- **Rate Limiting**: 100 requests per 15 minutes (configurable)
- **CORS**: Configurable origins
- **Input Validation**: Joi schema validation
- **Error Handling**: No sensitive information leakage

## Development Features

- **Hot Reload**: Automatic server restart on code changes
- **TypeScript**: Full type safety with strict configuration
- **ESLint**: Code quality and consistency
- **Path Aliases**: Clean imports with `@/` prefix
- **Correlation IDs**: Request tracking across the application
- **Structured Logging**: JSON or simple format logging
- **Async/Await**: Modern JavaScript patterns throughout

## Production Deployment

For production deployment:

1. Set `NODE_ENV=production`
2. Use proper Azure authentication (service principal or managed identity)
3. Configure appropriate CORS origins
4. Set up monitoring and log aggregation
5. Use PM2 or similar process manager
6. Set up health check monitoring

### Build for Production
```bash
npm run build
NODE_ENV=production npm start
```

## Testing

The project is set up for testing with Jest:

```bash
# Run tests
npm test

# Watch mode
npm run test:watch

# Coverage report
npm run test:coverage
```

## Contributing

1. Follow TypeScript strict mode guidelines
2. Use ESLint configuration provided
3. Write tests for new features
4. Follow the existing architecture patterns
5. Use correlation IDs in logging