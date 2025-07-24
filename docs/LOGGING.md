# Frontend Logging System

This document describes the comprehensive logging system implemented for the Brutus multi-cloud API project. The logging system captures all frontend events and sends them to container logs for centralized monitoring and debugging.

## Overview

The logging system consists of three main components:

1. **Frontend Logger Service** - Centralized logging service for Vue.js frontend
2. **Backend Log Handler** - Go API endpoints that receive and process frontend logs
3. **Container Log Output** - Structured JSON logs that appear in Docker/Kubernetes container logs

## Architecture

```
Frontend Components
       ↓
Logger Service (logger.ts)
       ↓
HTTP API (/api/v1/logs)
       ↓
Backend Log Handler (logs.go)
       ↓
Container stdout/stderr
       ↓
Docker/Kubernetes Logs
```

## Frontend Logger Service

### Location
`src/frontend/vue/src/services/logger.ts`

### Features
- **Singleton Pattern**: Single instance shared across the application
- **Structured Logging**: Consistent log format with metadata
- **Async Transmission**: Non-blocking log transmission to backend
- **Console Fallback**: Always logs to browser console for immediate debugging
- **Session Tracking**: Unique session ID for tracing user sessions
- **Context Enrichment**: Automatic addition of URL, user agent, and timestamp
- **Error Handling**: Graceful degradation when backend is unavailable

### Usage

```typescript
import { logger } from '../services/logger'

// Basic logging
logger.info('User logged in successfully')
logger.warn('API response was slow')
logger.error('Failed to load resource groups')
logger.debug('Token validation passed')

// Logging with context
logger.info('Resource group created', { 
  name: 'my-rg', 
  location: 'eastus',
  tags: { environment: 'dev' }
})

// Error logging with full context
try {
  await apiCall()
} catch (error) {
  logger.logError(error, 'API call failed', { 
    endpoint: '/resource-groups',
    method: 'POST'
  })
}
```

### Log Levels

| Level | Purpose | Example Use Case |
|-------|---------|------------------|
| `debug` | Detailed diagnostic information | Token validation, API parsing details |
| `info` | General informational messages | Successful operations, state changes |
| `warn` | Warning messages for recoverable issues | Slow API responses, fallback mechanisms |
| `error` | Error messages for failures | API failures, authentication errors |

### Log Entry Structure

Each log entry contains the following fields:

```typescript
interface LogEntry {
  timestamp: string        // ISO 8601 timestamp
  level: LogLevel         // debug, info, warn, error
  message: string         // Human-readable message
  context?: object        // Additional metadata
  source: 'frontend'      // Always 'frontend'
  sessionId: string       // Unique session identifier
  userAgent: string       // Browser user agent string
  url: string            // Current page URL
}
```

## Backend Log Handler

### Location
`src/backend/go/az/internal/handlers/logs.go`

### API Endpoints

#### POST /api/v1/logs
Receives a single log entry from the frontend.

**Request Body:**
```json
{
  "timestamp": "2024-01-24T10:30:45.123Z",
  "level": "info",
  "message": "Resource groups fetched successfully",
  "context": { "count": 3 },
  "source": "frontend",
  "sessionId": "frontend-1706089845123-abc123def",
  "userAgent": "Mozilla/5.0...",
  "url": "http://localhost:8080"
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

#### POST /api/v1/logs/batch
Receives multiple log entries in a single request.

**Request Body:**
```json
{
  "logs": [
    { /* LogEntry 1 */ },
    { /* LogEntry 2 */ },
    { /* LogEntry 3 */ }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "message": "Log batch received",
  "count": 3
}
```

### Container Log Format

Frontend logs appear in container logs with the following format:

```
[FRONTEND-LOG] {"timestamp":"2024-01-24T10:30:45.123Z","level":"info","source":"frontend","sessionId":"frontend-1706089845123-abc123def","message":"Resource groups fetched successfully","url":"http://localhost:8080","userAgent":"Mozilla/5.0...","context":{"count":3}}
```

## Integration Examples

### Authentication Store (auth.ts)

```typescript
// Before
console.log('Handling redirect result...')
console.error('Login failed:', error)

// After
logger.info('Handling redirect result...')
logger.logError(error, 'Login failed')
```

### Resource Group Management (ResourceGroupList.vue)

```typescript
// Before
console.log('Access token acquired, reinitializing API service...')
console.error('Failed to delete resource group:', err)

// After
logger.info('Access token acquired, reinitializing API service...')
logger.logError(err, 'Failed to delete resource group', { resourceGroupName: name })
```

### API Service (simpleApi.ts)

```typescript
// Before
console.log('API returned empty response - this is normal for DELETE operations')
console.error('Failed to fetch resource groups:', error)

// After
logger.debug('Azure API returned empty response - this is normal for DELETE operations')
logger.logError(error, 'Failed to fetch resource groups', { mode: this.mode })
```

## Monitoring and Debugging

### Viewing Logs in Development

1. **Browser Console**: All logs still appear in browser dev tools
2. **Backend Container**: View logs using Docker commands:
   ```bash
   # View Go backend logs
   docker-compose logs -f backend-az
   
   # Filter for frontend logs only
   docker-compose logs backend-az | grep "FRONTEND-LOG"
   ```

### Log Analysis

Frontend logs can be analyzed using standard log processing tools:

```bash
# Extract all frontend errors
docker logs brutus-backend-az 2>&1 | grep "FRONTEND-LOG" | jq 'select(.level == "error")'

# Count logs by level
docker logs brutus-backend-az 2>&1 | grep "FRONTEND-LOG" | jq -r '.level' | sort | uniq -c

# Find logs for specific session
docker logs brutus-backend-az 2>&1 | grep "FRONTEND-LOG" | jq 'select(.sessionId == "frontend-1706089845123-abc123def")'
```

### Production Monitoring

In production environments, these logs can be:
- Collected by log aggregation systems (ELK Stack, Splunk, etc.)
- Monitored for error patterns and alerts
- Used for user session tracking and debugging
- Analyzed for performance and usage patterns

## Configuration

### Environment Variables

The logging system respects the following environment variables:

| Variable | Purpose | Default |
|----------|---------|---------|
| `VITE_API_BASE_URL` | Backend URL for log transmission | `http://localhost:8081/api/v1` |

### Customization

To customize logging behavior, modify the `LoggerService` constructor in `logger.ts`:

```typescript
constructor() {
  this.sessionId = this.generateSessionId();
  this.backendUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8081/api/v1';
  this.fallbackToConsole = true; // Set to false to disable console logging
}
```

## Performance Considerations

1. **Async Processing**: Log transmission doesn't block UI operations
2. **Batching**: Multiple logs are sent in batches to reduce HTTP overhead
3. **Queue Management**: Failed log transmissions don't cause memory leaks
4. **Graceful Degradation**: Application continues working even if logging fails

## Security Considerations

1. **No Sensitive Data**: Avoid logging passwords, tokens, or personal information
2. **Context Filtering**: Be careful what context data is included in logs
3. **Log Retention**: Consider log retention policies for GDPR compliance
4. **Access Control**: Secure access to container logs in production

## Troubleshooting

### Common Issues

1. **Logs not appearing in container**: Check that backend is running and accessible
2. **High log volume**: Adjust log levels or add filtering
3. **Performance impact**: Monitor for excessive logging in tight loops

### Debug Mode

To enable more verbose logging, use debug level:

```typescript
logger.debug('Detailed diagnostic information', { 
  state: currentState,
  config: configuration 
})
```

## Future Enhancements

Potential improvements to the logging system:

1. **Log Levels Configuration**: Runtime configuration of log levels
2. **Sampling**: Reduce log volume in high-traffic scenarios
3. **Metrics Integration**: Automatic generation of metrics from logs
4. **Client-side Buffering**: Improved offline log handling
5. **Structured Query Support**: Enhanced log querying capabilities