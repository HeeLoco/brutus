# Brutus - Multi-Cloud API Backend

A production-ready, multi-language API backend project for managing cloud resources across different providers (Azure, AWS, GCP) with three complementary language implementations.

## 🏗️ Architecture Overview

This project implements a clean **multi-cloud, multi-language architecture** designed to leverage the strengths of different programming languages:

```
src/
└── backend/
    ├── go/         # High-performance, concurrent operations
    ├── python/     # Data processing, ML workflows, rapid prototyping  
    ├── typescript/ # Complex business logic, enterprise features
    └── [future]    # aws/, gcp/ for other cloud providers
```

Each language implementation provides the same Azure Resource Management API with different strengths and use cases.

## 🚀 Language Implementations

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

## 📋 Comparison Matrix

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

All implementations use the same environment variables:

```bash
# Required
AZURE_SUBSCRIPTION_ID="your-subscription-id"

# Optional (uses DefaultAzureCredential if not provided)
AZURE_TENANT_ID="your-tenant-id"
AZURE_CLIENT_ID="your-client-id" 
AZURE_CLIENT_SECRET="your-client-secret"

# Server configuration (optional)
PORT=8000  # or 8080 for Go
HOST=0.0.0.0
DEBUG=false
```

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

All implementations support multiple Azure authentication methods:

1. **Service Principal** (recommended for production)
   - Set `AZURE_TENANT_ID`, `AZURE_CLIENT_ID`, `AZURE_CLIENT_SECRET`
2. **Managed Identity** (for Azure-hosted applications)
3. **Azure CLI** (for development)
4. **Visual Studio Code** (for development)

If service principal credentials are not provided, all implementations automatically fall back to Azure DefaultAzureCredential.

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

### **Choose Go when:**
- Maximum performance is critical
- Building microservices
- Need excellent concurrency
- Memory usage is a concern
- Building CLI tools

### **Choose Python when:**
- Rapid prototyping and development
- Data processing and analysis
- ML/AI integration needed
- Working with scientific computing
- Team has strong Python experience

### **Choose TypeScript when:**
- Complex business logic
- Enterprise-grade features needed
- Type safety is important
- Real-time features required
- Frontend and backend team overlap

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

- **AWS Support:** `src/backend/*/aws/` implementations
- **Google Cloud Platform:** `src/backend/*/gcp/` implementations
- **Additional Languages:** Rust, Java, C# implementations
- **Advanced Features:** GraphQL APIs, WebSocket support, event streaming
- **DevOps:** Kubernetes manifests, Terraform modules, CI/CD pipelines

---

**Choose your language, build amazing cloud APIs!** 🚀