# Developer Guide

This guide provides comprehensive information for developers working with the Brutus multi-cloud API backend project.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Development vs Production Environments](#development-vs-production-environments)
3. [Docker Management](#docker-management)
4. [Azure AD Configuration](#azure-ad-configuration)
5. [Frontend Development](#frontend-development)
6. [Troubleshooting](#troubleshooting)
7. [Environment Configuration](#environment-configuration)

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Node.js 20+ (for local development)
- Go 1.24+ (for local development)
- Azure subscription and Azure AD app registration

### Getting Started

1. **Clone and navigate to the project:**
   ```bash
   git clone <repository-url>
   cd brutus
   ```

2. **Set up environment variables:**
   ```bash
   cd src/frontend/vue
   cp .env.example .env.development
   cp .env.example .env.production
   # Edit both files with your Azure credentials
   ```

3. **Start development environment:**
   ```bash
   ./scripts/docker-scripts.sh dev
   ```

4. **Access the application:**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - Health Check: http://localhost:8080/health

## Development vs Production Environments

### Overview

The project supports two distinct environments with the **same source code** but different builds and configurations:

| Aspect | Development | Production |
|--------|-------------|------------|
| **URL** | `http://localhost:5173` | `http://localhost` |
| **Server** | Vite dev server | Nginx static files |
| **Build** | Unminified, hot reload | Minified, optimized |
| **Environment File** | `.env.development` | `.env.production` |
| **Docker Command** | `docker-scripts.sh dev` | `docker-scripts.sh prod` |
| **Debugging** | Vue DevTools, source maps | Production optimized |

### Key Differences

#### Development Mode
- **Hot Module Replacement**: Changes reflect immediately
- **Source Maps**: Full debugging support
- **Vue DevTools**: Available in browser
- **Unminified Code**: Readable in browser DevTools
- **Error Overlay**: Vite shows build errors in browser

#### Production Mode  
- **Static Files**: Served by Nginx for performance
- **Minified Bundle**: Smaller file sizes
- **Tree Shaking**: Unused code removed
- **No Debugging Tools**: Production optimized
- **Gzip Compression**: Faster loading

### Environment Variables

Both environments use the same Azure credentials but different redirect URIs:

**Development (`.env.development`):**
```bash
VITE_REDIRECT_URI=http://localhost:5173
VITE_POST_LOGOUT_REDIRECT_URI=http://localhost:5173
```

**Production (`.env.production`):**
```bash
VITE_REDIRECT_URI=http://localhost
VITE_POST_LOGOUT_REDIRECT_URI=http://localhost
```

## Docker Management

### Docker Scripts

The project includes a comprehensive Docker management script at `./scripts/docker-scripts.sh`:

```bash
# Start development environment
./scripts/docker-scripts.sh dev

# Start production environment  
./scripts/docker-scripts.sh prod

# Build all images
./scripts/docker-scripts.sh build

# Build specific services
./scripts/docker-scripts.sh build backend          # Go backend only
./scripts/docker-scripts.sh build frontend-dev    # Development frontend only  
./scripts/docker-scripts.sh build frontend-prod   # Production frontend only

# Force rebuild (no cache)
./scripts/docker-scripts.sh build --force         # Rebuild all from scratch
./scripts/docker-scripts.sh build backend --force # Force rebuild backend only

# Stop all services
./scripts/docker-scripts.sh stop

# Check environment configuration
./scripts/docker-scripts.sh check-env

# View container status
./scripts/docker-scripts.sh status

# View logs
./scripts/docker-scripts.sh logs

# Test all endpoints
./scripts/docker-scripts.sh test

# Clean up everything
./scripts/docker-scripts.sh clean
```

### Container Architecture

#### Development Stack
- **frontend-dev**: Vue.js with Vite dev server (port 5173)
- **backend-go**: Go API server (port 8080)

#### Production Stack  
- **frontend-prod**: Nginx serving static files (port 80)
- **backend-go**: Go API server (port 8080)

### Docker Files

- **`Dockerfile.dev`**: Single-stage development container with hot reload
- **`Dockerfile`**: Multi-stage production build with Nginx
- **`docker-compose.yml`**: Service definitions with profiles (dev/prod)

### Build Management

The enhanced build system supports granular control over which services to build:

#### Build All Services
```bash
./scripts/docker-scripts.sh build
```
Builds all three services: backend, frontend-dev, and frontend-prod with their respective profiles.

#### Build Individual Services
```bash
# Backend service (Go API)
./scripts/docker-scripts.sh build backend

# Development frontend (Vite dev server)
./scripts/docker-scripts.sh build frontend-dev

# Production frontend (Nginx static files)  
./scripts/docker-scripts.sh build frontend-prod

# Alternative service names
./scripts/docker-scripts.sh build dev    # Same as frontend-dev
./scripts/docker-scripts.sh build prod   # Same as frontend-prod
```

#### Force Rebuild (No Cache)
```bash
# Force rebuild all services
./scripts/docker-scripts.sh build --force

# Force rebuild specific service
./scripts/docker-scripts.sh build backend --force
./scripts/docker-scripts.sh build frontend-dev --force
```

#### Build Output
The build command shows:
- Build progress for each service
- Success/failure status with error handling
- List of built images with sizes and timestamps
- Clear indication of which services were built

## Azure AD Configuration

### App Registration Requirements

Your Azure AD app registration must include **both** redirect URIs:

1. **Authentication → Redirect URIs:**
   - `http://localhost:5173` (for development)
   - `http://localhost` (for production)

2. **API Permissions:**
   - `https://management.azure.com/user_impersonation`
   - `openid`
   - `profile` 
   - `offline_access`

### Environment Variables

Required in both `.env.development` and `.env.production`:

```bash
# Azure App Registration
VITE_AZURE_CLIENT_ID=your-client-id
VITE_AZURE_TENANT_ID=your-tenant-id-or-common
VITE_AZURE_SUBSCRIPTION_ID=your-subscription-id

# Backend Configuration
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_BACKEND_TYPE=go
```

### Multi-Tenant Support

- Use `VITE_AZURE_TENANT_ID=common` for multi-tenant applications
- Use specific tenant ID for single-tenant applications

## Frontend Development

### Architecture

The Vue.js frontend uses:

- **Vue 3** with Composition API and TypeScript
- **Pinia** for state management  
- **Vue Router** for navigation
- **MSAL.js** for Azure AD authentication
- **Vite** for build tooling

### Key Components

- **`MainApp.vue`**: Main authentication wrapper
- **`DemoLoginForm.vue`**: Login interface with multiple modes
- **`ResourceGroupList.vue`**: Resource management interface
- **`EnvironmentDebug.vue`**: Environment debugging tool

### Authentication Modes

The frontend supports multiple API modes:

1. **Demo Mode**: Mock data for testing
2. **Backend Mode**: Routes through Go/Python/TypeScript backends
3. **Azure Direct**: Direct Azure Management API calls

### Development Workflow

1. **Start development environment:**
   ```bash
   ./scripts/docker-scripts.sh dev
   ```

2. **Access development tools:**
   - Application: http://localhost:5173
   - Vue DevTools: Available in browser
   - Environment Debug: http://localhost:5173/debug

3. **Hot reload**: Changes reflect immediately
4. **Debugging**: Full source maps and Vue DevTools support

### Build Process

**Development:**
```bash
npm run dev  # Starts Vite dev server
```

**Production:**
```bash
npm run build      # Creates optimized build
npm run preview    # Preview production build locally
```

## Troubleshooting

### Common Issues

#### 1. Login Fails with "interaction_in_progress" Error

**Symptoms:**
```
[ERROR] Login failed 
BrowserAuthError: interaction_in_progress: Interaction is currently in progress
```

**Solutions (in order):**

1. **Use Clear Login State button:**
   - Click "Clear Login State" in the login form
   - Wait for success message
   - Try logging in again

2. **Refresh the page:**
   - Click "Refresh Page" button
   - Or manually refresh browser (F5)

3. **Clear browser storage manually:**
   ```javascript
   // In browser console - Clear cookies and session storage
   document.cookie.split(";").forEach(c => {
     document.cookie = c.replace(/^ +/, "").replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/");
   });
   sessionStorage.clear()
   location.reload()
   ```

4. **Use incognito/private browsing mode**

#### 2. Different Behavior Between Dev and Prod

**This is normal!** Dev and prod environments have different:
- URLs (localhost:5173 vs localhost)
- Build optimization
- Debugging tools
- Server technology (Vite vs Nginx)

**Solution:** Ensure both redirect URIs are registered in Azure AD.

#### 3. Container Build Failures

**Check environment configuration:**
```bash
./scripts/docker-scripts.sh check-env
```

**View container logs:**
```bash
./scripts/docker-scripts.sh logs [service-name]
```

**Clean and rebuild:**
```bash
./scripts/docker-scripts.sh clean
./scripts/docker-scripts.sh build

# Or force rebuild specific service
./scripts/docker-scripts.sh build backend --force
./scripts/docker-scripts.sh build frontend-dev --force
```

#### 4. Port Conflicts

**Default ports:**
- Development frontend: 5173
- Production frontend: 80  
- Backend API: 8080

**Check for conflicts:**
```bash
lsof -i :5173
lsof -i :80
lsof -i :8080
```

### Debug Tools

#### Environment Debug Page

Access `http://localhost:5173/debug` to view:
- Environment variables (with sensitive values masked)
- Runtime information
- Authentication status
- Azure AD configuration requirements
- Backend connectivity testing

#### Log Analysis

**Frontend logs appear in backend container:**
```bash
# View all frontend logs
docker logs brutus-backend-go-1 | grep "FRONTEND-LOG"

# Parse as JSON
docker logs brutus-backend-go-1 | grep "FRONTEND-LOG" | jq '.'

# Filter by log level
docker logs brutus-backend-go-1 | grep "FRONTEND-LOG" | jq 'select(.level == "error")'
```

#### Health Checks

```bash
# Backend health
curl http://localhost:8080/health

# Test all endpoints
./scripts/docker-scripts.sh test
```

## Environment Configuration

### File Structure

```
src/frontend/vue/
├── .env.example          # Template with all variables
├── .env.development      # Development-specific values
├── .env.production       # Production-specific values
└── .env                  # Fallback (optional)
```

### Environment Loading Order

1. `.env.development` (in dev mode)
2. `.env.production` (in production mode)  
3. `.env` (fallback)
4. `.env.example` (documentation only)

### Security Notes

- Never commit actual credentials to git
- Use `.env.example` for documentation
- Sensitive values are masked in debug displays
- Frontend logs mask credentials automatically

### Validation

Use the environment check command:
```bash
./scripts/docker-scripts.sh check-env
```

This validates:
- Required environment files exist
- Key variables are set
- Container configuration
- Azure AD redirect URI setup

---

## Quick Reference

### Essential Commands

```bash
# Development
./scripts/docker-scripts.sh dev
# → http://localhost:5173

# Production  
./scripts/docker-scripts.sh prod
# → http://localhost

# Environment debug
# → http://localhost:5173/debug

# Logs and status
./scripts/docker-scripts.sh logs
./scripts/docker-scripts.sh status
./scripts/docker-scripts.sh check-env
```

### Troubleshooting Checklist

1. ✅ Both redirect URIs registered in Azure AD?
2. ✅ Environment files configured with correct values?
3. ✅ Docker containers running and healthy?
4. ✅ No port conflicts?
5. ✅ Browser storage cleared if authentication fails?

For additional help, check the environment debug page or container logs.