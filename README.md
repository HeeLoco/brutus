# Brutus - Multi-Cloud API Backend

A production-ready, multi-language API backend project for managing cloud resources across different providers (Azure, AWS, GCP) with three complementary language implementations.

## 🏗️ Architecture Overview

This project implements a clean **multi-cloud, multi-language architecture** designed to leverage the strengths of different programming languages:

```
src/
├── backend/
│   ├── go/         # High-performance, concurrent operations
│   ├── python/     # Data processing, ML workflows, rapid prototyping  
│   ├── typescript/ # Complex business logic, enterprise features
│   └── [future]    # aws/, gcp/ for other cloud providers
└── frontend/
    └── vue/        # Vue 3 + Vite frontend with TypeScript
```

Each backend language implementation provides the same Azure Resource Management API with different strengths and use cases. The Vue frontend provides a modern, containerized interface for consuming these APIs.

## 🚀 Backend Implementations

### **Go Implementation** (`src/backend/go/az/`)
**Best for:** High-performance operations, concurrent processing, microservices
- **Framework:** Gin HTTP framework
- **Port:** 8080 (default)
- **Strengths:** Blazing fast, low memory usage, excellent concurrency
- **Authentication:** Azure DefaultAzureCredential + Service Principal
- **Features:** Clean architecture, CORS support, structured error handling

### **Python Implementation** (`src/backend/python/az/`)
**Best for:** Data processing, ML workflows, rapid development
- **Framework:** FastAPI with async support
- **Port:** 8000 (default)  
- **Strengths:** Rapid development, rich ecosystem, excellent for data science
- **Authentication:** Azure DefaultAzureCredential + Service Principal
- **Features:** Auto-generated docs, Pydantic validation, correlation ID tracking

### **TypeScript Implementation** (`src/backend/typescript/az/`)
**Best for:** Complex business logic, enterprise features, real-time applications
- **Framework:** Express.js with full TypeScript support
- **Port:** 8000 (default)
- **Strengths:** Type safety, enterprise features, excellent developer experience
- **Authentication:** Azure DefaultAzureCredential + Service Principal
- **Features:** Enterprise middleware, Joi validation, Winston logging, rate limiting

## 🎨 Frontend Implementation

### **Vue 3 + Vite** (`src/frontend/vue/`)
**Modern web interface for Azure Resource Management**
- **Framework:** Vue 3 with Composition API + Vite build tool
- **Port:** 5173 (dev), 80 (container)
- **Authentication:** Azure AD integration with MSAL.js
- **Features:** TypeScript, Vue Router, Pinia state management, Vitest testing
- **Container:** Multi-stage Docker build with Nginx (~20MB final image)
- **API Integration:** Direct Azure SDK calls and backend API support

## 📋 Comparison Matrix

### Backend Comparison
| Feature | **Go** | **Python** | **TypeScript** |
|---------|---------|------------|----------------|
| **Performance** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Development Speed** | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Type Safety** | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Enterprise Features** | ⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Memory Usage** | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐ |
| **Ecosystem** | ⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |


## 🛠️ Quick Start

### Prerequisites
- **Go:** >= 1.24
- **Python:** >= 3.11
- **Node.js:** >= 18.0
- **Azure Subscription** with appropriate permissions

### Environment Setup

For **frontend authentication** (recommended), **no Azure environment variables are needed**:

```bash
# Server configuration (optional)
PORT=8000  # or 8080 for Go
HOST=0.0.0.0
DEBUG=false
```

The backends are configured for **frontend authentication** where:
- Users authenticate directly with Azure AD through the frontend (MSAL.js)
- Users enter their subscription ID in the frontend interface
- Frontend sends user access tokens AND subscription ID to backend APIs
- Backend uses these tokens to make Azure API calls on behalf of the user
- **No Azure credentials or configuration needed in backend!**

### 🐹 Go Setup

```bash
cd src/backend/go/az
export AZURE_SUBSCRIPTION_ID="your-subscription-id"

# Run development
go run cmd/server/main.go

# Build for production
go build -o bin/azure-api cmd/server/main.go
./bin/azure-api
```

**Available at:** http://localhost:8080

### 🐍 Python Setup

```bash
cd src/backend/python/az

# Setup virtual environment
python3 -m venv .venv
source .venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Setup environment
cp .env.example .env
# Edit .env with your Azure credentials

# Run development server
python main.py
# Or with uvicorn
uvicorn main:app --host 0.0.0.0 --port 8000 --reload
```

**Available at:** http://localhost:8000
- **Swagger UI:** http://localhost:8000/docs
- **ReDoc:** http://localhost:8000/redoc

### 📘 TypeScript Setup

```bash
cd src/backend/typescript/az

# Install dependencies
npm install

# Setup environment
cp .env.example .env
# Edit .env with your Azure credentials

# Run development server (with hot reload)
npm run dev

# Build and run production
npm run build && npm start
```

**Available at:** http://localhost:8000

### 🎨 Frontend Setup

#### Vue 3 + Vite Frontend

```bash
cd src/frontend/vue

# Install dependencies
npm install

# Setup environment
cp .env.example .env
# Edit .env to point to your preferred backend

# Run development server
npm run dev
```

**Available at:** http://localhost:5173


### 🐳 Docker Deployment (Recommended)

**New to Docker?** Check our beginner guides:
- 📚 **[Docker Quick Start](docs/DOCKER_QUICK_START.md)** - Get running in 3 minutes
- 📖 **[Complete Docker Guide](docs/DOCKER_GUIDE.md)** - Comprehensive tutorial

#### Quick Start with Docker
```bash
# Start development environment (Vue + Go backend)
./scripts/docker-scripts.sh dev

# Access applications:
# - Vue Frontend: http://localhost:5173
# - Go Backend: http://localhost:8080

# View logs
./scripts/docker-scripts.sh logs

# Stop everything
./scripts/docker-scripts.sh stop
```

#### Available Docker Commands
```bash
./scripts/docker-scripts.sh dev         # Development mode (Vue + Go backend)
./scripts/docker-scripts.sh prod        # Production mode  
./scripts/docker-scripts.sh build       # Build images only
./scripts/docker-scripts.sh status      # Show running services
./scripts/docker-scripts.sh help        # Show all commands
```

#### Manual Docker Build (Advanced)
```bash
# Vue frontend
cd src/frontend/vue
docker build -t brutus-vue-frontend .
docker run -p 80:80 brutus-vue-frontend
```

## 🔗 API Endpoints

All implementations provide the same REST API endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/health` | Health check with Azure connectivity test |
| `GET` | `/` | API information and status |
| `GET` | `/api/v1/resource-groups` | List all resource groups |
| `POST` | `/api/v1/resource-groups` | Create resource group |
| `GET` | `/api/v1/resource-groups/{name}` | Get specific resource group |
| `PUT` | `/api/v1/resource-groups/{name}` | Update resource group |
| `DELETE` | `/api/v1/resource-groups/{name}` | Delete resource group |

### Example API Request

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

## 🔒 Authentication

The application uses **frontend authentication** for better security and user experience:

### **Frontend Authentication (Current Implementation)**
- Users log in directly to Azure AD through the frontend (MSAL.js)
- Frontend obtains user access tokens and sends them to backend APIs
- Backend uses user tokens to access Azure resources **on behalf of the user**
- Only requires `AZURE_SUBSCRIPTION_ID` in backend configuration
- Users see their own Azure resources based on their Azure AD permissions

### **Benefits of Frontend Authentication:**
- ✅ No service principal secrets stored in backend
- ✅ Users authenticate with their own Azure AD credentials  
- ✅ Fine-grained access control based on user permissions
- ✅ Better security posture - no shared credentials
- ✅ Users can only access resources they have permission to

## 🏢 Production Deployment

### Using Docker

Python implementation includes Docker support:

```bash
docker-compose up --build backend-az
```

### Manual Deployment

For production deployment:

1. Set `NODE_ENV=production` (TypeScript) or equivalent
2. Use proper Azure authentication (service principal or managed identity)
3. Configure appropriate CORS origins
4. Set up monitoring and log aggregation
5. Use process managers (PM2, systemd, Docker)
6. Implement health check monitoring

## 📊 When to Use Which Implementation

### **Choose Go Backend when:**
- Maximum performance is critical
- Building microservices
- Need excellent concurrency
- Memory usage is a concern
- Building CLI tools

### **Choose Python Backend when:**
- Rapid prototyping and development
- Data processing and analysis
- ML/AI integration needed
- Working with scientific computing
- Team has strong Python experience

### **Choose TypeScript Backend when:**
- Complex business logic
- Enterprise-grade features needed
- Type safety is important
- Real-time features required
- Frontend and backend team overlap

### **Choose Vue 3 + Vite Frontend when:**
- Building component-based applications
- Need rapid development with great DX
- Want comprehensive ecosystem and community
- Working with complex state management
- Team has Vue/React experience

### **Choose SvelteKit Frontend when:**
- Performance is critical (smallest bundles)
- Building content-focused applications
- Want compile-time optimization
- Prefer simpler, less verbose syntax
- Need excellent SEO and static generation

## 🔄 Multi-Language Benefits

This architecture provides several advantages:

1. **Performance Optimization:** Use the right tool for each job
2. **Team Flexibility:** Leverage existing team expertise
3. **Technology Evolution:** Easily adopt new technologies
4. **Risk Mitigation:** Reduce dependency on a single language
5. **Learning Opportunities:** Cross-pollinate ideas between implementations

## 🧪 Testing

Each implementation includes comprehensive testing setup:

```bash
# Go
go test ./...

# Python  
pytest

# TypeScript
npm test
```

## 📚 Documentation

Detailed documentation for each implementation:

- **Go:** `src/backend/go/az/README.go.md`
- **Python:** `src/backend/python/az/README.python.md`
- **TypeScript:** `src/backend/typescript/az/README.typescript.md`
- **Architecture:** `CLAUDE.md`

## 🤝 Contributing

1. Follow each language's best practices and conventions
2. Maintain API compatibility across all implementations
3. Update all relevant documentation
4. Add tests for new features
5. Ensure proper error handling and logging

## 🔮 Future Roadmap

### Backend Expansion
- **AWS Support:** `src/backend/*/aws/` implementations
- **Google Cloud Platform:** `src/backend/*/gcp/` implementations
- **Additional Languages:** Rust, Java, C# implementations
- **Advanced Features:** GraphQL APIs, WebSocket support, event streaming

### Frontend Enhancements
- **Additional Frameworks:** React, Angular, Solid.js implementations
- **Mobile Apps:** React Native, Flutter companion apps
- **Advanced Features:** Real-time updates, offline support, PWA capabilities

### DevOps & Infrastructure
- **Container Orchestration:** Kubernetes manifests, Helm charts
- **Infrastructure as Code:** Terraform modules, Pulumi templates
- **CI/CD Pipelines:** GitHub Actions, Azure DevOps workflows

---

**Choose your language, build amazing cloud APIs!** 🚀