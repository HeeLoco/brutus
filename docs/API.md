# API Documentation

This document provides comprehensive API documentation for the Brutus multi-cloud resource management platform.

## Base URLs

| Environment | URL | Description |
|-------------|-----|-------------|
| Go Backend | `http://localhost:8081/api/v1` | Primary Go-based API server |
| Python Backend | `http://localhost:8000/api/v1` | Python FastAPI server |
| TypeScript Backend | `http://localhost:3000/api/v1` | Node.js Express server |

## Authentication

All resource management endpoints require Azure AD authentication. The API supports user token-based authentication:

```http
Authorization: Bearer <azure-ad-user-token>
```

## Common Response Formats

### Success Response
```json
{
  "success": true,
  "data": { /* response data */ },
  "message": "Operation completed successfully"
}
```

### Error Response
```json
{
  "error": "Error description",
  "code": 400,
  "message": "Detailed error message"
}
```

## Health Check Endpoints

### GET /health

Check the health status of the API server.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-24T10:30:45.123Z",
  "service": "Azure API",
  "version": "1.0.0"
}
```

## Resource Group Management

### GET /api/v1/resource-groups

List all resource groups in the subscription.

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:**
```json
{
  "resource_groups": [
    {
      "id": "/subscriptions/sub-id/resourceGroups/my-rg",
      "name": "my-rg",
      "location": "eastus",
      "type": "Microsoft.Resources/resourceGroups",
      "tags": {
        "environment": "dev",
        "project": "brutus"
      },
      "properties": {
        "provisioningState": "Succeeded"
      }
    }
  ],
  "count": 1
}
```

### GET /api/v1/resource-groups/{name}

Get details of a specific resource group.

**Parameters:**
- `name` (path) - Resource group name

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": "/subscriptions/sub-id/resourceGroups/my-rg",
  "name": "my-rg",
  "location": "eastus",
  "type": "Microsoft.Resources/resourceGroups",
  "tags": {
    "environment": "dev"
  },
  "properties": {
    "provisioningState": "Succeeded"
  }
}
```

**Error Responses:**
- `404 Not Found`: Resource group does not exist
- `401 Unauthorized`: Invalid or expired token
- `403 Forbidden`: Insufficient permissions

### POST /api/v1/resource-groups

Create a new resource group.

**Headers:**
```http
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "my-new-rg",
  "location": "eastus",
  "tags": {
    "environment": "production",
    "team": "devops"
  }
}
```

**Response:**
```json
{
  "id": "/subscriptions/sub-id/resourceGroups/my-new-rg",
  "name": "my-new-rg",
  "location": "eastus",
  "type": "Microsoft.Resources/resourceGroups",
  "tags": {
    "environment": "production",
    "team": "devops"
  },
  "properties": {
    "provisioningState": "Succeeded"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or missing required fields
- `409 Conflict`: Resource group already exists
- `401 Unauthorized`: Invalid or expired token

### PUT /api/v1/resource-groups/{name}

Update an existing resource group.

**Parameters:**
- `name` (path) - Resource group name

**Headers:**
```http
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "tags": {
    "environment": "staging",
    "updated": "2024-01-24"
  }
}
```

**Response:**
```json
{
  "id": "/subscriptions/sub-id/resourceGroups/my-rg",
  "name": "my-rg",
  "location": "eastus",
  "type": "Microsoft.Resources/resourceGroups",
  "tags": {
    "environment": "staging",
    "updated": "2024-01-24"
  },
  "properties": {
    "provisioningState": "Succeeded"
  }
}
```

### DELETE /api/v1/resource-groups/{name}

Delete a resource group and all its resources.

**Parameters:**
- `name` (path) - Resource group name

**Headers:**
```http
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "Resource group deletion initiated",
  "name": "my-rg"
}
```

**Error Responses:**
- `404 Not Found`: Resource group does not exist
- `401 Unauthorized`: Invalid or expired token
- `403 Forbidden`: Insufficient permissions

## Frontend Logging Endpoints

### POST /api/v1/logs

Receive a single log entry from the frontend application.

**Headers:**
```http
Content-Type: application/json
```

**Request Body:**
```json
{
  "timestamp": "2024-01-24T10:30:45.123Z",
  "level": "info",
  "message": "User logged in successfully",
  "context": {
    "username": "user@example.com",
    "loginMethod": "azure-ad"
  },
  "source": "frontend",
  "sessionId": "frontend-1706089845123-abc123def",
  "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
  "url": "http://localhost:8080/resource-groups"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Log entry received",
  "count": 1
}
```

### POST /api/v1/logs/batch

Receive multiple log entries in a single request for improved performance.

**Headers:**
```http
Content-Type: application/json
```

**Request Body:**
```json
{
  "logs": [
    {
      "timestamp": "2024-01-24T10:30:45.123Z",
      "level": "info",
      "message": "Component mounted",
      "source": "frontend",
      "sessionId": "frontend-1706089845123-abc123def",
      "userAgent": "Mozilla/5.0...",
      "url": "http://localhost:8080"
    },
    {
      "timestamp": "2024-01-24T10:30:46.456Z",
      "level": "debug",
      "message": "API call initiated",
      "context": {
        "endpoint": "/resource-groups",
        "method": "GET"
      },
      "source": "frontend",
      "sessionId": "frontend-1706089845123-abc123def",
      "userAgent": "Mozilla/5.0...",
      "url": "http://localhost:8080"
    }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Log batch received",
  "count": 2
}
```

## Log Entry Schema

### LogEntry Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `timestamp` | string | Yes | ISO 8601 timestamp |
| `level` | string | Yes | Log level: `debug`, `info`, `warn`, `error` |
| `message` | string | Yes | Human-readable log message |
| `context` | object | No | Additional structured data |
| `source` | string | Yes | Always `"frontend"` for frontend logs |  
| `sessionId` | string | Yes | Unique session identifier |
| `userAgent` | string | No | Browser user agent string |
| `url` | string | No | Current page URL |

### Log Levels

| Level | Purpose | Example Use Case |
|-------|---------|------------------|
| `debug` | Detailed diagnostic information | API response parsing, token validation |
| `info` | General informational messages | User actions, successful operations |
| `warn` | Warning conditions | Slow responses, fallback mechanisms |
| `error` | Error conditions | API failures, authentication errors |

## Error Codes

| Code | Description | Common Causes |
|------|-------------|---------------|
| 400 | Bad Request | Invalid JSON, missing required fields |
| 401 | Unauthorized | Missing or invalid authorization token |
| 403 | Forbidden | Token valid but insufficient permissions |
| 404 | Not Found | Resource group does not exist |
| 409 | Conflict | Resource group already exists |
| 500 | Internal Server Error | Server-side errors, Azure API failures |

## Rate Limiting

The API implements rate limiting to prevent abuse:

- **Resource Groups**: 100 requests per minute per user
- **Logging**: 1000 requests per minute per session

Rate limit headers are included in responses:
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1706089905
```

## CORS Support

The API supports Cross-Origin Resource Sharing (CORS) for web applications:

**Allowed Origins:**
- `http://localhost:8080` (Frontend development)
- `http://localhost:5173` (Vite development server)
- `http://localhost:3000` (Alternative frontend port)

**Allowed Methods:**
- GET, POST, PUT, DELETE, OPTIONS

**Allowed Headers:**
- Authorization, Content-Type, X-Requested-With

## Examples

### Frontend Integration

```typescript
// Resource Group API Client Example
class ResourceGroupClient {
  constructor(private baseUrl: string, private token: string) {}

  async listResourceGroups(): Promise<ResourceGroup[]> {
    const response = await fetch(`${this.baseUrl}/resource-groups`, {
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      }
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    const data = await response.json();
    return data.resource_groups;
  }

  async createResourceGroup(rg: CreateResourceGroupRequest): Promise<ResourceGroup> {
    const response = await fetch(`${this.baseUrl}/resource-groups`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(rg)
    });

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }

    return response.json();
  }
}
```

### Logging Integration

```typescript
// Frontend Logger Integration Example
import { logger } from '../services/logger';

// Log user action
logger.info('Resource group created', {
  name: 'my-rg',
  location: 'eastus',
  duration: 1234
});

// Log error with context
try {
  await createResourceGroup(data);
} catch (error) {
  logger.logError(error, 'Failed to create resource group', {
    resourceGroupName: data.name,
    location: data.location
  });
}
```

### cURL Examples

```bash
# List resource groups
curl -X GET "http://localhost:8081/api/v1/resource-groups" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json"

# Create resource group
curl -X POST "http://localhost:8081/api/v1/resource-groups" \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-rg",
    "location": "eastus",
    "tags": {"environment": "test"}
  }'

# Send log entry
curl -X POST "http://localhost:8081/api/v1/logs" \
  -H "Content-Type: application/json" \
  -d '{
    "timestamp": "2024-01-24T10:30:45.123Z",
    "level": "info",
    "message": "Test log entry",
    "source": "frontend",
    "sessionId": "test-session-123"
  }'
```

## Monitoring and Observability

### Container Logs

All API requests and frontend logs are written to container stdout/stderr:

```bash
# View all logs
docker-compose logs -f backend-az

# Filter frontend logs
docker-compose logs backend-az | grep "FRONTEND-LOG"

# Filter API access logs
docker-compose logs backend-az | grep "API-REQUEST"
```

### Health Monitoring

Monitor API health using the health endpoint:

```bash
# Check API health
curl http://localhost:8081/health

# Monitor continuously
watch -n 5 'curl -s http://localhost:8081/health | jq .'
```

### Metrics Collection

The API provides metrics suitable for monitoring systems:

- Request count and duration
- Error rates by endpoint
- Authentication success/failure rates
- Frontend log volume by level

For production deployments, consider integrating with:
- Prometheus for metrics collection
- Grafana for visualization
- ELK Stack for log analysis
- Azure Monitor for cloud-native monitoring