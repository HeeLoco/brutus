# Python Azure API

Production-ready FastAPI application for managing Azure resources with clean architecture, comprehensive error handling, and modern Python best practices.

## Features

- **Clean Architecture**: Layered structure with routers, services, models, and middleware
- **Async/Await Support**: Full async support for better performance
- **Azure Integration**: Comprehensive Azure Resource Management with proper authentication
- **Configuration Management**: Environment-based configuration with validation
- **Structured Logging**: Correlation ID tracking and structured logging
- **Error Handling**: Comprehensive exception handling with proper HTTP status codes
- **API Documentation**: Auto-generated OpenAPI/Swagger documentation
- **Input Validation**: Pydantic models for request/response validation
- **CORS Support**: Cross-origin requests enabled
- **Health Checks**: Built-in health monitoring with Azure connectivity tests

## Project Structure

```
src/backend/python/az/
├── api/
│   ├── core/
│   │   ├── config.py           # Settings and configuration
│   │   ├── logging.py          # Logging configuration
│   │   └── exceptions.py       # Custom exception classes
│   ├── models/
│   │   └── resource_group.py   # Pydantic models
│   ├── services/
│   │   └── azure_service.py    # Azure SDK integration
│   ├── routers/
│   │   ├── health.py           # Health check endpoints
│   │   └── resource_groups.py  # Resource group endpoints
│   └── middleware/
│       ├── correlation.py      # Request correlation tracking
│       └── error_handler.py    # Global error handling
├── main.py                     # FastAPI application entry point
├── requirements.txt            # Python dependencies
└── .env.example               # Environment variables template
```

## Environment Variables

Create a `.env` file from `.env.example`:

```bash
# Required
AZURE_SUBSCRIPTION_ID=your-subscription-id

# Optional (uses DefaultAzureCredential if not provided)
AZURE_TENANT_ID=your-tenant-id
AZURE_CLIENT_ID=your-client-id
AZURE_CLIENT_SECRET=your-client-secret

# Server configuration (optional)
HOST=0.0.0.0
PORT=8000
DEBUG=false
LOG_LEVEL=INFO
```

## Authentication

The API supports multiple Azure authentication methods:

1. **Service Principal** (recommended for production)
2. **Managed Identity** (for Azure-hosted applications)
3. **Azure CLI** (for development)
4. **Visual Studio Code** (for development)

Set the appropriate environment variables or rely on Azure DefaultAzureCredential.

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | API information |
| GET | `/health` | Health check with Azure connectivity |
| GET | `/docs` | Swagger UI documentation |
| GET | `/redoc` | ReDoc documentation |
| GET | `/api/v1/resource-groups/` | List all resource groups |
| POST | `/api/v1/resource-groups/` | Create resource group |
| GET | `/api/v1/resource-groups/{name}` | Get specific resource group |
| PUT | `/api/v1/resource-groups/{name}` | Update resource group |
| DELETE | `/api/v1/resource-groups/{name}` | Delete resource group |

## Request/Response Examples

### Create Resource Group
```bash
POST /api/v1/resource-groups/
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
  },
  "type": "Microsoft.Resources/resourceGroups"
}
```

### Error Response
```json
{
  "error": "ResourceNotFoundError",
  "message": "Resource group 'non-existent' not found",
  "correlation_id": "a1b2c3d4",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Development

### Setup
```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

### Run Development Server
```bash
export AZURE_SUBSCRIPTION_ID="your-subscription-id"
uvicorn main:app --host 0.0.0.0 --port 8000 --reload
```

### Testing
Visit the API documentation:
- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc
- Health check: http://localhost:8000/health

## Production Deployment

For production deployment, consider:

1. Set `DEBUG=false`
2. Use proper Azure authentication (service principal or managed identity)
3. Configure appropriate CORS origins
4. Set up proper monitoring and logging
5. Use a production WSGI server configuration