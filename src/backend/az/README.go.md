# Go Azure API

Production-ready REST API server built with Gin framework for managing Azure resources.

## Features

- **Azure Resource Group Management**: Complete CRUD operations
- **Authentication**: Azure DefaultAzureCredential with service principal support
- **Clean Architecture**: Organized with handlers, services, models, and middleware
- **CORS Support**: Cross-origin requests enabled
- **Health Checks**: Built-in health monitoring endpoint
- **Error Handling**: Structured error responses

## Project Structure

```
src/backend/go/az/
├── cmd/server/          # Application entry point
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── services/        # Business logic layer
│   ├── models/          # Data models
│   └── middleware/      # HTTP middleware
├── configs/             # Configuration management
└── bin/                 # Compiled binaries
```

## Environment Variables

```bash
# Required
AZURE_SUBSCRIPTION_ID="your-subscription-id"

# Optional (uses DefaultAzureCredential if not provided)
AZURE_TENANT_ID="your-tenant-id"
AZURE_CLIENT_ID="your-client-id" 
AZURE_CLIENT_SECRET="your-client-secret"

# Server configuration
PORT="8080"  # Default port
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/` | Service information |
| GET | `/api/v1/resource-groups` | List all resource groups |
| POST | `/api/v1/resource-groups` | Create resource group |
| GET | `/api/v1/resource-groups/{name}` | Get specific resource group |
| PUT | `/api/v1/resource-groups/{name}` | Update resource group |
| DELETE | `/api/v1/resource-groups/{name}` | Delete resource group |

## Request/Response Examples

### Create Resource Group
```bash
POST /api/v1/resource-groups
Content-Type: application/json

{
  "name": "my-resource-group",
  "location": "eastus",
  "tags": {
    "environment": "dev",
    "project": "brutus"
  }
}
```

### Response
```json
{
  "id": "/subscriptions/.../resourceGroups/my-resource-group",
  "name": "my-resource-group", 
  "location": "eastus",
  "tags": {
    "environment": "dev",
    "project": "brutus"
  }
}
```