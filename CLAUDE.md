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

**TypeScript API:**
- `GET /health` - Health check with detailed Azure connectivity
- `GET /api/v1/resource-groups` - List resource groups with structured response
- `POST /api/v1/resource-groups` - Create resource group with validation
- `GET /api/v1/resource-groups/{name}` - Get specific resource group
- `PUT /api/v1/resource-groups/{name}` - Update resource group with validation
- `DELETE /api/v1/resource-groups/{name}` - Delete resource group with confirmation

### Docker Configuration
- **Dockerfile.server.az**: Python 3.12.3 based container for Azure server
- **compose.yml**: Defines backend-az service on port 8000