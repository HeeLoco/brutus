# Troubleshooting Guide

This guide helps you troubleshoot common issues with the Brutus multi-cloud resource management platform.

## Table of Contents

1. [Frontend Logging Issues](#frontend-logging-issues)
2. [Authentication Problems](#authentication-problems)
3. [API Connection Issues](#api-connection-issues)
4. [Resource Group Management](#resource-group-management)
5. [Docker and Container Issues](#docker-and-container-issues)
6. [Development Environment](#development-environment)

## Frontend Logging Issues

### Logs Not Appearing in Container

**Symptoms:**
- Frontend actions occur but no logs appear in container output
- Only browser console shows logs

**Diagnosis:**
```bash
# Check if backend is running
docker-compose ps

# Check backend logs for errors
docker-compose logs backend-az

# Test log endpoint directly
curl -X POST "http://localhost:8081/api/v1/logs" \
  -H "Content-Type: application/json" \
  -d '{"timestamp":"2024-01-24T10:30:45.123Z","level":"info","message":"Test","source":"frontend","sessionId":"test-123"}'
```

**Solutions:**
1. **Backend Not Running:**
   ```bash
   docker-compose up backend-az
   ```

2. **Wrong Backend URL:**
   Check `VITE_API_BASE_URL` in frontend environment:
   ```bash
   # In src/frontend/vue/.env
   VITE_API_BASE_URL=http://localhost:8081/api/v1
   ```

3. **CORS Issues:**
   Verify CORS settings in backend allow frontend origin

### High Log Volume / Performance Issues

**Symptoms:**
- Application becomes slow
- Too many logs in container output

**Solutions:**
1. **Adjust Log Levels:**
   ```typescript
   // Temporarily disable debug logs
   if (process.env.NODE_ENV === 'production') {
     // Override debug method
     logger.debug = () => {};
   }
   ```

2. **Add Filtering:**
   ```typescript
   // Add conditional logging
   if (shouldLog(level, context)) {
     logger.info(message, context);
   }
   ```

### Log Format Issues

**Symptoms:**
- Logs appear malformed in container output
- JSON parsing errors in log analysis

**Diagnosis:**
```bash
# Check log format
docker-compose logs backend-az | grep "FRONTEND-LOG" | head -1 | jq '.'
```

**Solutions:**
1. **Validate Log Structure:**
   ```typescript
   // Ensure all required fields are present
   const logEntry: LogEntry = {
     timestamp: new Date().toISOString(),
     level: 'info',
     message: 'Test message',
     source: 'frontend',
     sessionId: 'session-123'
   };
   ```

## Authentication Problems

### Azure AD Login Fails

**Symptoms:**
- Redirect loop during login
- "Invalid client" errors
- Token acquisition failures

**Diagnosis:**
1. **Check App Registration:**
   ```bash
   # Use the verification script
   ./verify-azure-app.sh
   ```

2. **Verify Redirect URIs:**
   - Check Azure Portal → App Registration → Authentication
   - Ensure `http://localhost:8080` is configured as SPA redirect URI

**Solutions:**
1. **Fix Redirect URIs:**
   ```bash
   # Add missing redirect URIs
   az ad app update --id <client-id> --spa-redirect-uris "http://localhost:8080" "http://localhost:5173"
   ```

2. **Enable Public Client Flows:**
   ```bash
   az ad app update --id <client-id> --is-fallback-public-client
   ```

3. **Check Permissions:**
   - Azure Service Management → user_impersonation
   - Microsoft Graph → User.Read

### Token Expiration Issues

**Symptoms:**
- Authentication works initially but fails after some time
- "Token expired" errors in logs

**Solutions:**
1. **Check Token Refresh:**
   ```typescript
   // Ensure silent token refresh is working
   logger.info('Token refresh attempt', { hasAccount: !!account });
   ```

2. **Handle Token Renewal:**
   ```typescript
   // Add token renewal logic
   if (tokenExpired) {
     await authStore.getAccessToken(); // Triggers silent renewal
   }
   ```

### Demo Mode Stuck

**Symptoms:**
- Application shows demo data despite authentication
- "Using Demo mode (fallback)" in logs

**Diagnosis:**
```typescript
// Check authentication state
logger.info('Auth state debug', {
  hasToken: !!authStore.state.accessToken,
  hasAccount: !!authStore.state.account,
  isAuthenticated: authStore.state.isAuthenticated
});
```

**Solutions:**
1. **Wait for Token:**
   Authentication timing issue - token may still be loading

2. **Force Refresh:**
   ```typescript
   // Manually trigger token acquisition
   await authStore.getAccessToken();
   ```

## API Connection Issues

### Backend Unreachable

**Symptoms:**
- "Failed to fetch" errors
- Connection refused messages

**Diagnosis:**
```bash
# Check if backend is running
curl http://localhost:8081/health

# Check Docker containers
docker-compose ps

# Check port bindings
docker-compose port backend-az 8080
```

**Solutions:**
1. **Start Backend:**
   ```bash
   docker-compose up backend-az
   ```

2. **Check Port Configuration:**
   ```yaml
   # In compose.yml
   services:
     backend-az:
       ports:
         - "8081:8080"  # host:container
   ```

### API Authentication Fails

**Symptoms:**
- 401 Unauthorized errors
- "Invalid token" messages

**Diagnosis:**
```bash
# Test with valid token
TOKEN="your-azure-token-here"
curl -H "Authorization: Bearer $TOKEN" http://localhost:8081/api/v1/resource-groups
```

**Solutions:**
1. **Verify Token:**
   ```typescript
   // Log token validity
   logger.debug('Token check', { 
     hasToken: !!token,
     tokenLength: token?.length,
     tokenPrefix: token?.substring(0, 10)
   });
   ```

2. **Check Token Scopes:**
   Token must include Azure Resource Manager scopes

### CORS Errors

**Symptoms:**
- "CORS policy" errors in browser console
- Preflight request failures

**Solutions:**
1. **Update CORS Configuration:**
   ```go
   // In Go backend middleware
   c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
   c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
   c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
   ```

## Resource Group Management

### Resource Groups Not Loading

**Symptoms:**
- Empty resource group list
- Loading state persists

**Diagnosis:**
```typescript
// Check API mode and responses
logger.info('API diagnosis', {
  mode: apiService.mode,
  subscriptionId: authStore.state.subscriptionId,
  hasToken: !!authStore.state.accessToken
});
```

**Solutions:**
1. **Verify Subscription Access:**
   ```bash
   # Check Azure CLI access
   az account show
   az group list
   ```

2. **Check Subscription ID:**
   ```bash
   # Verify subscription ID in environment
   echo $AZURE_SUBSCRIPTION_ID
   ```

### Delete Operations Fail

**Symptoms:**
- Delete button doesn't work
- Resource groups remain after deletion

**Diagnosis:**
```typescript
// Add debug logging to delete function
logger.info('Delete attempt', { name, mode: apiService.mode });
```

**Solutions:**
1. **Check Permissions:**
   User needs Contributor role on resource group

2. **Verify API Mode:**
   Ensure not in demo mode for actual deletions

## Docker and Container Issues

### Container Won't Start

**Symptoms:**
- "docker-compose up" fails
- Container exits immediately

**Diagnosis:**
```bash
# Check container logs
docker-compose logs backend-az

# Check image build
docker-compose build backend-az

# Verify Dockerfile
cat src/backend/go/az/Dockerfile
```

**Solutions:**
1. **Rebuild Container:**
   ```bash
   docker-compose build --no-cache backend-az
   docker-compose up backend-az
   ```

2. **Check Dependencies:**
   ```bash
   # Verify Go modules
   cd src/backend/go/az
   go mod tidy
   go mod verify
   ```

### Port Conflicts

**Symptoms:**
- "Port already in use" errors
- Cannot bind to port

**Solutions:**
1. **Check Port Usage:**
   ```bash
   # Find process using port 8081
   lsof -i :8081
   
   # Kill process if needed
   sudo kill -9 <pid>
   ```

2. **Use Different Port:**
   ```yaml
   # In compose.yml
   services:
     backend-az:
       ports:
         - "8082:8080"  # Use different host port
   ```

### Volume Mount Issues

**Symptoms:**
- Code changes not reflected in container
- File not found errors

**Solutions:**
1. **Fix Volume Paths:**
   ```yaml
   # In compose.yml
   volumes:
     - ./src/backend/go/az:/app
   ```

2. **Check File Permissions:**
   ```bash
   # Ensure files are readable
   chmod -R 755 src/backend/go/az
   ```

## Development Environment

### Frontend Build Fails

**Symptoms:**
- Vite build errors
- TypeScript compilation issues

**Solutions:**
1. **Clear Cache:**
   ```bash
   cd src/frontend/vue
   rm -rf node_modules package-lock.json
   npm install
   ```

2. **Check Node Version:**
   ```bash
   node --version  # Should be >= 16
   npm --version
   ```

### Environment Variables Not Loading

**Symptoms:**
- Configuration values are undefined
- Default values being used

**Solutions:**
1. **Check .env File:**
   ```bash
   # Verify .env file exists and has correct format
   cat src/frontend/vue/.env
   
   # Variables must start with VITE_
   VITE_AZURE_CLIENT_ID=your-client-id
   ```

2. **Restart Development Server:**
   ```bash
   # Environment changes require restart
   npm run dev
   ```

### Hot Reload Not Working

**Symptoms:**
- Changes require manual refresh
- Development server doesn't detect file changes

**Solutions:**
1. **Check File Watchers:**
   ```bash
   # Increase file watcher limit (Linux)
   echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf
   sudo sysctl -p
   ```

2. **Use Polling:**
   ```javascript
   // In vite.config.ts
   export default defineConfig({
     server: {
       watch: {
         usePolling: true
       }
     }
   })
   ```

## Debugging Commands

### Useful Docker Commands

```bash
# View all logs with timestamps
docker-compose logs -f -t

# Filter specific service logs
docker-compose logs -f backend-az

# Execute commands in running container
docker-compose exec backend-az /bin/sh

# View container resource usage
docker stats $(docker-compose ps -q)
```

### Log Analysis Commands

```bash
# Count log entries by level
docker-compose logs backend-az | grep "FRONTEND-LOG" | jq -r '.level' | sort | uniq -c

# Find recent errors
docker-compose logs --since=10m backend-az | grep -i error

# Extract specific session logs
docker-compose logs backend-az | grep "FRONTEND-LOG" | jq 'select(.sessionId == "session-id")'

# Monitor logs in real-time
docker-compose logs -f backend-az | grep "FRONTEND-LOG" | jq '.'
```

### Health Check Commands

```bash
# Test all endpoints
curl http://localhost:8081/health
curl http://localhost:8081/api/v1/resource-groups
curl -X POST http://localhost:8081/api/v1/logs -d '{"test":"data"}'

# Monitor endpoint availability
watch -n 5 'curl -s http://localhost:8081/health | jq .'
```

## Getting Help

If you're still experiencing issues:

1. **Check Logs:** Review container logs for error messages
2. **Verify Configuration:** Ensure all environment variables are set correctly
3. **Test Components Individually:** Isolate the problem to frontend, backend, or Azure
4. **Review Documentation:** Check `docs/LOGGING.md` and `docs/API.md`
5. **Create Minimal Reproduction:** Create a simple test case that demonstrates the issue

For Azure-specific issues, also check:
- Azure Portal for resource states
- Azure AD App Registration configuration
- Azure subscription permissions and limits