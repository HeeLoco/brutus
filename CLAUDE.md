# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a multi-cloud, multi-language API backend project for controlling cloud resources across different providers:

**Directory Structure:**
- `src/` - Source code for all components  
- `src/backend/` - Backend API implementations
- `src/backend/go/` - Go implementations
- `src/backend/python/` - Python implementations
- `src/backend/typescript/` - TypeScript/Node.js implementations
- `az/` - Azure-specific APIs (current: `src/backend/go/az/`, `src/backend/python/az/`, `src/backend/typescript/az/`)
- Future cloud providers: `aws/` (Amazon Web Services), `gcp/` (Google Cloud Platform), and others

**Current Implementation:**
- **Python Azure API**: FastAPI-based Azure resource management server in `src/backend/python/az/`
- **Go Azure API**: Production-ready Gin-based REST API server in `src/backend/go/az/` for Azure resource management
- **TypeScript Azure API**: Express.js-based enterprise API server in `src/backend/typescript/az/` for Azure resource management

## Development Commands

### Python (Azure Server)

Setup virtual environment:
```bash
python3 -m venv .venv
source .venv/bin/activate  # bash/zsh
source .venv/bin/activate.fish  # fish
```

Install dependencies:
```bash
pip install -r src/backend/python/az/requirements.txt
```

Set up environment variables:
```bash
cd src/backend/python/az
cp .env.example .env
# Edit .env with your Azure credentials
export AZURE_SUBSCRIPTION_ID="your-subscription-id"
```

Run the server:
```bash
cd src/backend/python/az
python main.py
# Or with uvicorn directly:
uvicorn main:app --host 0.0.0.0 --port 8000 --reload
```

View API documentation:
- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc

### TypeScript (Azure APIs)

Install dependencies:
```bash
cd src/backend/typescript/az
npm install
```

Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your Azure credentials
export AZURE_SUBSCRIPTION_ID="your-subscription-id"
```

Run development server:
```bash
npm run dev
# Or production build:
npm run build && npm start
```

API endpoints available at:
- Health check: http://localhost:8000/health
- API: http://localhost:8000/api/v1
- Resource groups: http://localhost:8000/api/v1/resource-groups

### Go (Azure APIs)

Run the Go Azure API server:
```bash
cd src/backend/go/az
# Set required environment variables first
export AZURE_SUBSCRIPTION_ID="your-subscription-id"
export AZURE_TENANT_ID="your-tenant-id"  # Optional with DefaultAzureCredential
export AZURE_CLIENT_ID="your-client-id"  # Optional with DefaultAzureCredential
export AZURE_CLIENT_SECRET="your-client-secret"  # Optional with DefaultAzureCredential
export PORT="8080"  # Optional, defaults to 8080

# Run the server
go run cmd/server/main.go
```

Build the application:
```bash
cd src/backend/go/az
go build -o bin/azure-api cmd/server/main.go
./bin/azure-api
```

### Docker

Build and run the Azure server using Docker Compose:
```bash
docker-compose up --build backend-az
```

The server will be available on port 8000.

## Development Environment Management

### Docker Scripts

The project includes comprehensive Docker management scripts for easy development:

```bash
# Quick Start - Development Environment
./scripts/docker-scripts.sh dev
# → Frontend: http://localhost:5173 (Vue + Vite dev server)
# → Backend: http://localhost:8080 (Go API)

# Production Environment
./scripts/docker-scripts.sh prod  
# → Frontend: http://localhost (Nginx static files)
# → Backend: http://localhost:8080 (Go API)

# Build commands (supports individual services and force rebuild)
./scripts/docker-scripts.sh build                    # Build all images
./scripts/docker-scripts.sh build backend           # Build Go backend only
./scripts/docker-scripts.sh build frontend-dev      # Build dev frontend only
./scripts/docker-scripts.sh build frontend-prod     # Build prod frontend only
./scripts/docker-scripts.sh build --force           # Force rebuild all (no cache)

# Other useful commands
./scripts/docker-scripts.sh stop         # Stop all services
./scripts/docker-scripts.sh status       # Show container status
./scripts/docker-scripts.sh logs         # View logs
./scripts/docker-scripts.sh check-env    # Validate environment
./scripts/docker-scripts.sh test         # Test all endpoints
./scripts/docker-scripts.sh clean        # Remove everything
```

### Development vs Production Environments

The project supports two distinct environments with **identical source code** but different builds:

| Environment | URL | Server | Build Type | Use Case |
|-------------|-----|--------|------------|----------|
| **Development** | `http://localhost:5173` | Vite dev server | Unminified, hot reload | Development, debugging |
| **Production** | `http://localhost` | Nginx static files | Minified, optimized | Production testing |

**Key Differences:**
- **Development**: Hot module replacement, Vue DevTools, source maps, error overlay
- **Production**: Minified bundle, static file serving, performance optimized

**Environment Files:**
- `.env.development` → Used in development mode
- `.env.production` → Used in production mode
- Both require identical Azure credentials but different redirect URIs

### Frontend Environment Setup

1. **Copy environment templates:**
   ```bash
   cd src/frontend/vue
   cp .env.example .env.development
   cp .env.example .env.production
   ```

2. **Configure Azure AD credentials in both files:**
   ```bash
   # Required in both .env.development and .env.production
   VITE_AZURE_CLIENT_ID=your-client-id
   VITE_AZURE_TENANT_ID=your-tenant-id
   VITE_AZURE_SUBSCRIPTION_ID=your-subscription-id
   ```

3. **Different redirect URIs per environment:**
   ```bash
   # .env.development
   VITE_REDIRECT_URI=http://localhost:5173
   
   # .env.production  
   VITE_REDIRECT_URI=http://localhost
   ```

4. **Register BOTH redirect URIs in Azure AD App Registration:**
   - `http://localhost:5173` (development)
   - `http://localhost` (production)

### Debugging and Troubleshooting

**Environment Debug Page:**
- Access `http://localhost:5173/debug` to view configuration details
- Shows environment variables (with sensitive values masked)
- Tests backend connectivity
- Displays authentication status

**Login Issues:**
If authentication fails with "interaction_in_progress" error:
1. Click "Clear Login State" button in login form
2. Try "Refresh Page" button for complete reset  
3. Use browser incognito mode as fallback

**Log Analysis:**
```bash
# View frontend logs (appear in backend container)
docker logs brutus-backend-go-1 | grep "FRONTEND-LOG"

# Parse frontend logs as JSON
docker logs brutus-backend-go-1 | grep "FRONTEND-LOG" | jq '.'

# Environment validation
./scripts/docker-scripts.sh check-env
```

**Common Port Usage:**
- Development frontend: 5173
- Production frontend: 80
- Go backend: 8080
- Health checks: Available at `/health` endpoint

For comprehensive development guidance, see [`docs/DEVELOPER_GUIDE.md`](docs/DEVELOPER_GUIDE.md).

## Architecture

### Python Azure Server (`src/backend/python/az/`)
- **Framework**: FastAPI with async/await support
- **Architecture**: Clean layered architecture (routers, services, models, middleware)
- **Authentication**: Azure DefaultAzureCredential with service principal fallback
- **Configuration**: Pydantic settings with environment variable support
- **Logging**: Structured logging with correlation IDs
- **Error Handling**: Comprehensive exception handling with proper HTTP status codes
- **API Documentation**: Auto-generated OpenAPI/Swagger documentation
- **Middleware**: CORS, request correlation, and error handling
- **Dependencies**: FastAPI, Pydantic, Azure SDK, Uvicorn

### Go Application (`src/backend/go/az/`)
- **Framework**: Gin HTTP framework for REST API
- **Authentication**: Azure DefaultAzureCredential (supports managed identity, Azure CLI, service principal)
- **Structure**: Clean architecture with handlers, services, models, and middleware
- **API Endpoints**: Complete CRUD operations for Azure resource groups
- **Frontend Logging**: Centralized logging system that receives frontend logs via HTTP API
- **Dependencies**: Azure SDK for Go (azidentity, armresources)

### TypeScript Azure Server (`src/backend/typescript/az/`)
- **Framework**: Express.js with full TypeScript support
- **Architecture**: Enterprise-grade layered architecture (controllers, services, middleware)
- **Authentication**: Azure DefaultAzureCredential with service principal support
- **Security**: Helmet, rate limiting, CORS, input validation
- **Logging**: Winston with correlation ID tracking
- **Validation**: Joi schema validation for all inputs
- **Developer Experience**: Hot reload, ESLint, path aliases, comprehensive npm scripts
- **Error Handling**: Comprehensive error handling with proper HTTP status codes
- **Dependencies**: Express.js, Azure SDK, Winston, Joi, TypeScript

### Key API Endpoints

**Python API:**
- `GET /health` - Health check with Azure connectivity test
- `GET /api/v1/resource-groups/` - List all resource groups
- `POST /api/v1/resource-groups/` - Create resource group
- `GET /api/v1/resource-groups/{name}` - Get specific resource group
- `PUT /api/v1/resource-groups/{name}` - Update resource group
- `DELETE /api/v1/resource-groups/{name}` - Delete resource group

**Go API:**
- `GET /health` - Health check
- `GET /api/v1/resource-groups` - List resource groups
- `POST /api/v1/resource-groups` - Create resource group
- `GET /api/v1/resource-groups/{name}` - Get specific resource group
- `PUT /api/v1/resource-groups/{name}` - Update resource group
- `DELETE /api/v1/resource-groups/{name}` - Delete resource group
- `POST /api/v1/logs` - Receive single frontend log entry
- `POST /api/v1/logs/batch` - Receive batch of frontend log entries

**TypeScript API:**
- `GET /health` - Health check with detailed Azure connectivity
- `GET /api/v1/resource-groups` - List resource groups with structured response
- `POST /api/v1/resource-groups` - Create resource group with validation
- `GET /api/v1/resource-groups/{name}` - Get specific resource group
- `PUT /api/v1/resource-groups/{name}` - Update resource group with validation
- `DELETE /api/v1/resource-groups/{name}` - Delete resource group with confirmation

## Frontend Architecture

### Vue.js Application (`src/frontend/vue/`)
- **Framework**: Vue 3 with Composition API and TypeScript
- **Authentication**: MSAL (Microsoft Authentication Library) for Azure AD integration with secure cookie storage
- **State Management**: Pinia for centralized state management
- **API Integration**: Multi-mode API service (Demo, Backend, Azure Direct)
- **Logging**: Centralized logging service that sends logs to backend containers
- **Components**: Reactive components with auto-refresh capabilities
- **Cookie-Based Auth**: Secure authentication using HTTP cookies with proper security flags
- **Authentication Modes**:
  - **Demo Mode**: Uses mock data for testing
  - **Backend Mode**: Communicates through backend APIs
  - **Azure Direct Mode**: Direct calls to Azure Management APIs using user tokens

### Authentication Storage Architecture

The application uses **secure HTTP cookies** for authentication persistence instead of localStorage:

- **Primary Storage**: HTTP cookies with `Secure`, `SameSite=Strict` flags
- **Cookie Types**:
  - `auth_token` - Azure access token with automatic expiration
  - `account_info` - User account information (JSON)
  - `refresh_token` - Long-term refresh token (30-day expiry)
- **MSAL Cache**: Uses sessionStorage as fallback for MSAL's internal cache
- **Security Benefits**: 
  - CSRF protection via SameSite
  - Automatic expiration handling
  - SSR compatibility
  - Secure transport over HTTPS

### Frontend Logging System

The frontend implements a comprehensive logging system that captures all application events and sends them to container logs for centralized monitoring:

- **Logger Service** (`src/frontend/vue/src/services/logger.ts`): Centralized logging with structured JSON format
- **Session Tracking**: Unique session IDs for tracing user interactions
- **Context Enrichment**: Automatic addition of URL, user agent, timestamp, and custom context
- **Async Transmission**: Non-blocking log transmission to backend
- **Console Fallback**: Always maintains browser console output for immediate debugging
- **Error Handling**: Comprehensive error logging with stack traces and context

**Log Viewing:**
```bash
# View all container logs
docker-compose logs -f backend-az

# Filter for frontend logs only
docker-compose logs backend-az | grep "FRONTEND-LOG"

# Parse frontend logs as JSON
docker-compose logs backend-az | grep "FRONTEND-LOG" | jq '.'
```

**Frontend Dependencies**: Vue 3, TypeScript, MSAL, Pinia, Vite

### Docker Configuration
- **Dockerfile.server.az**: Python 3.12.3 based container for Azure server
- **compose.yml**: Defines backend-az service on port 8000
- **Frontend Logging**: All frontend logs appear in backend container logs with structured JSON format

## Logging and Monitoring

### Frontend Logging
All frontend application events are captured and sent to backend container logs:

- **Location**: `src/frontend/vue/src/services/logger.ts`
- **Format**: Structured JSON logs with session tracking
- **Transmission**: HTTP API to backend (`/api/v1/logs`)
- **Fallback**: Browser console for immediate debugging
- **Monitoring**: Container logs accessible via Docker/Kubernetes

### Log Analysis Examples
```bash
# View all frontend logs
docker-compose logs backend-az | grep "FRONTEND-LOG"

# Filter by log level
docker logs container_id | grep "FRONTEND-LOG" | jq 'select(.level == "error")'

# Track specific user session
docker logs container_id | grep "FRONTEND-LOG" | jq 'select(.sessionId == "session-id")'

# Count logs by level
docker logs container_id | grep "FRONTEND-LOG" | jq -r '.level' | sort | uniq -c
```

For detailed logging documentation, see: `docs/LOGGING.md`