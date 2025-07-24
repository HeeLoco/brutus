# Docker Setup and Usage Guide

This guide covers how to use Docker for development and production with the Brutus Azure Resource Manager application.

## Prerequisites

1. **Docker Desktop** - Install from [docker.com](https://www.docker.com/products/docker-desktop/)
2. **Docker Compose** - Included with Docker Desktop

## Quick Start

### 1. Start Docker Desktop
Make sure Docker Desktop is running before proceeding.

### 2. Use the Docker Scripts
We've provided a convenient script to manage all Docker operations:

```bash
# Make the script executable
chmod +x docker-scripts.sh

# Show available commands
./docker-scripts.sh help

# Start development environment
./docker-scripts.sh dev

# Start production environment
./docker-scripts.sh prod
```

## Available Environments

### Development Environment
- **Vue.js Frontend**: Hot reload, source mapping, development server
- **Go Backend**: Live backend API
- **Port Mapping**: 
  - Frontend: http://localhost:5173
  - Backend: http://localhost:8080

```bash
# Start development
./docker-scripts.sh dev

# View logs
./docker-scripts.sh logs frontend-dev
```

### Production Environment
- **Optimized builds**: Minified, compressed assets
- **Nginx serving**: Production-ready web server
- **Health checks**: Built-in container health monitoring
- **Port Mapping**:
  - Frontend: http://localhost (port 80)
  - Backend: http://localhost:8080

```bash
# Start production
./docker-scripts.sh prod

# Check status
./docker-scripts.sh status
```

### SvelteKit Alternative
- **Both frontends**: Vue and SvelteKit running simultaneously
- **Port Mapping**:
  - Vue: http://localhost:5173
  - Svelte: http://localhost:5174

```bash
# Start with SvelteKit
./docker-scripts.sh dev-svelte
```

## Manual Docker Commands

If you prefer manual control:

### Frontend Development
```bash
# Build development image
docker build -f src/frontend/vue/Dockerfile.dev -t brutus-vue-dev src/frontend/vue

# Run development container
docker run -p 5173:5173 -v $(pwd)/src/frontend/vue:/app -v /app/node_modules brutus-vue-dev
```

### Frontend Production
```bash
# Build production image
docker build -f src/frontend/vue/Dockerfile -t brutus-vue-prod src/frontend/vue

# Run production container
docker run -p 80:80 brutus-vue-prod
```

### Using Docker Compose
```bash
# Development
docker-compose --profile dev up --build

# Production
docker-compose --profile prod up --build

# SvelteKit development
docker-compose --profile dev --profile svelte up --build

# Background mode (detached)
docker-compose --profile dev up --build -d
```

## Environment Configuration

### 1. Copy Environment File
```bash
cp .env.docker.example .env.docker
```

### 2. Configure Values
Edit `.env.docker` with your Azure credentials:

```bash
# Azure Authentication (Optional - for real Azure integration)
AZURE_SUBSCRIPTION_ID=your-subscription-id
AZURE_TENANT_ID=your-tenant-id
AZURE_CLIENT_ID=your-client-id
AZURE_CLIENT_SECRET=your-client-secret

# Frontend Environment Variables
VITE_AZURE_CLIENT_ID=your-app-registration-client-id
VITE_AZURE_TENANT_ID=your-tenant-id-or-common
VITE_AZURE_SUBSCRIPTION_ID=your-subscription-id
```

### 3. Load Environment
```bash
# Load environment variables
export $(cat .env.docker | xargs)

# Or use with docker-compose
docker-compose --env-file .env.docker --profile dev up
```

## Container Architecture

### Multi-stage Builds
Both frontends use optimized multi-stage Docker builds:

1. **Builder Stage**: Node.js environment for building the application
2. **Production Stage**: Lightweight Nginx serving optimized assets

### Image Sizes
- **Vue Development**: ~300MB (includes Node.js and source)
- **Vue Production**: ~25MB (Nginx + built assets)
- **SvelteKit Production**: ~20MB (Nginx + built assets)

### Health Checks
All production containers include health checks:
```bash
# Check container health
docker ps
# Look for "healthy" status

# Manual health check
curl http://localhost/health
```

## Networking

All services run on a custom Docker network `brutus-network`:
- **Internal communication**: Services can communicate using service names
- **External access**: Only exposed ports are accessible from host

### Service Communication
```yaml
# Frontend can call backend using:
http://backend-go:8080/api/v1/resource-groups

# Instead of:
http://localhost:8080/api/v1/resource-groups
```

## Troubleshooting

### Common Issues

1. **Port already in use**:
   ```bash
   # Stop conflicting services
   ./docker-scripts.sh stop
   
   # Or kill specific port
   lsof -ti:5173 | xargs kill -9
   ```

2. **Build fails**:
   ```bash
   # Clean build
   docker system prune -f
   ./docker-scripts.sh build
   ```

3. **Permission issues**:
   ```bash
   # Fix file permissions
   chmod +x docker-scripts.sh
   
   # Fix volume permissions (if needed)
   sudo chown -R $USER:$USER src/frontend/vue/node_modules
   ```

4. **Cannot connect to Docker daemon**:
   - Ensure Docker Desktop is running
   - Check Docker Desktop preferences
   - Restart Docker Desktop if needed

### Debugging

```bash
# View logs for all services
./docker-scripts.sh logs

# View logs for specific service
./docker-scripts.sh logs frontend-dev

# Interactive shell in container
docker exec -it <container-name> sh

# Inspect container
docker inspect <container-name>
```

### Performance

```bash
# View resource usage
docker stats

# Clean up unused resources
docker system prune -f

# Remove all brutus images
docker rmi $(docker images -q brutus*)
```

## Production Deployment

### Docker Swarm
```bash
# Initialize swarm
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.yml brutus
```

### Kubernetes
Convert the docker-compose.yml using tools like Kompose:
```bash
# Install kompose
curl -L https://github.com/kubernetes/kompose/releases/download/v1.26.0/kompose-linux-amd64 -o kompose

# Convert
./kompose convert -f docker-compose.yml
```

## Security Considerations

1. **Environment Variables**: Never commit real credentials to Git
2. **Network Security**: Use custom networks for isolation
3. **Image Security**: Regularly update base images
4. **Volume Permissions**: Set appropriate file permissions

## Monitoring

### Health Checks
```bash
# Test all endpoints
./docker-scripts.sh test

# Manual checks
curl http://localhost:5173    # Frontend dev
curl http://localhost/health  # Frontend prod
curl http://localhost:8080/health  # Backend
```

### Logs
```bash
# View logs
docker-compose logs -f

# Export logs
docker-compose logs > docker-logs.txt
```

This Docker setup provides a complete development and production environment that's easy to use and maintain!