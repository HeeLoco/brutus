# 🐳 Docker Guide for Beginners

This guide will help you run the Brutus Azure Resource Manager application using Docker on your local machine, even if you're completely new to Docker.

## 📋 Table of Contents

- [What is Docker?](#what-is-docker)
- [Prerequisites](#prerequisites)
- [Docker Installation](#docker-installation)
- [Understanding Docker Concepts](#understanding-docker-concepts)
- [Project Docker Structure](#project-docker-structure)
- [Building Docker Images](#building-docker-images)
- [Using Docker Compose](#using-docker-compose)
- [Common Commands](#common-commands)
- [Troubleshooting](#troubleshooting)
- [Tips for Beginners](#tips-for-beginners)

## 🤔 What is Docker?

Docker is a platform that packages applications and their dependencies into lightweight, portable containers. Think of containers as:

- **Virtual boxes** that contain everything needed to run an application
- **Consistent environments** that work the same on any machine
- **Isolated spaces** that don't interfere with your local system

### Why Use Docker for This Project?

- ✅ **Easy setup**: No need to install Go, Python, Node.js manually
- ✅ **Consistent environment**: Works the same on Windows, Mac, Linux
- ✅ **Isolation**: Doesn't mess with your local system
- ✅ **Quick deployment**: Start all services with one command

## 🛠️ Prerequisites

Before starting, you need:

1. **A computer** running Windows 10/11, macOS, or Linux
2. **Administrator/root access** to install Docker
3. **At least 4GB RAM** available for Docker
4. **Internet connection** to download Docker images

## 📦 Docker Installation

### Windows 10/11

1. **Download Docker Desktop**:
   - Go to [docker.com/products/docker-desktop](https://www.docker.com/products/docker-desktop)
   - Click "Download for Windows"

2. **Install Docker Desktop**:
   - Run the downloaded installer
   - Follow the installation wizard
   - Restart your computer when prompted

3. **Verify Installation**:
   ```powershell
   docker --version
   docker-compose --version
   ```

### macOS

1. **Download Docker Desktop**:
   - Go to [docker.com/products/docker-desktop](https://www.docker.com/products/docker-desktop)
   - Click "Download for Mac"
   - Choose the right version (Intel or Apple Silicon)

2. **Install Docker Desktop**:
   - Open the downloaded `.dmg` file
   - Drag Docker to Applications folder
   - Launch Docker from Applications

3. **Verify Installation**:
   ```bash
   docker --version
   docker-compose --version
   ```

### Linux (Ubuntu/Debian)

1. **Update package index**:
   ```bash
   sudo apt update
   ```

2. **Install Docker**:
   ```bash
   sudo apt install docker.io docker-compose
   ```

3. **Add your user to docker group**:
   ```bash
   sudo usermod -aG docker $USER
   ```

4. **Log out and back in**, then verify:
   ```bash
   docker --version
   docker-compose --version
   ```

## 🧠 Understanding Docker Concepts

### Key Terms

- **Image**: A blueprint/template for creating containers (like a recipe)
- **Container**: A running instance of an image (like a cake made from the recipe)
- **Dockerfile**: Instructions for building an image (the recipe itself)
- **Docker Compose**: Tool for running multi-container applications (orchestrator)

### Visual Analogy

```
Dockerfile → Image → Container
(Recipe)   (Template) (Running App)
```

## 🏗️ Project Docker Structure

Our project has the following Docker-related files:

```
brutus/
├── docker-compose.yml          # Main orchestration file
├── compose.yml                 # Alternative compose file
├── Dockerfile.backend.python.az  # Python backend image
├── src/
│   ├── backend/
│   │   └── go/az/
│   │       └── Dockerfile      # Go backend image
│   └── frontend/
│       ├── vue/
│       │   ├── Dockerfile      # Vue production image
│       │   └── Dockerfile.dev  # Vue development image
│       └── svelte/
│           ├── Dockerfile      # Svelte production image
│           └── Dockerfile.dev  # Svelte development image
└── scripts/
    ├── docker-scripts.sh       # Helper scripts
    └── validate-docker.sh      # Validation script
```

## 🔨 Building Docker Images

### Method 1: Build Individual Images

Navigate to the project root directory:

```bash
cd /path/to/brutus
```

#### Build Python Backend
```bash
docker build -f Dockerfile.backend.python.az -t brutus-python-backend .
```

#### Build Go Backend
```bash
docker build -f src/backend/go/az/Dockerfile -t brutus-go-backend src/backend/go/az
```

#### Build Vue Frontend (Production)
```bash
docker build -f src/frontend/vue/Dockerfile -t brutus-vue-frontend src/frontend/vue
```

#### Build Vue Frontend (Development)
```bash
docker build -f src/frontend/vue/Dockerfile.dev -t brutus-vue-dev src/frontend/vue
```

### Method 2: Use Helper Script

We provide a helper script that builds all images:

```bash
# Make script executable
chmod +x scripts/docker-scripts.sh

# Build all images
./scripts/docker-scripts.sh build
```

### Understanding Build Commands

Let's break down a build command:
```bash
docker build -f Dockerfile.backend.python.az -t brutus-python-backend .
```

- `docker build`: Command to build an image
- `-f Dockerfile.backend.python.az`: Use this specific Dockerfile
- `-t brutus-python-backend`: Tag (name) the image
- `.`: Build context (current directory)

## 🐳 Using Docker Compose

Docker Compose allows you to run multiple containers together with one command.

### Understanding docker-compose.yml

Our `docker-compose.yml` file defines multiple services:

```yaml
services:
  backend-python:     # Python API server
  backend-go:         # Go API server  
  frontend-vue:       # Vue.js frontend
  frontend-svelte:    # SvelteKit frontend
```

### Starting Services

#### Start All Services
```bash
docker-compose up
```

#### Start in Background (Detached Mode)
```bash
docker-compose up -d
```

#### Start Specific Services
```bash
# Only Python backend
docker-compose up backend-python

# Backend + Vue frontend
docker-compose up backend-python frontend-vue
```

#### Start with Fresh Build
```bash
docker-compose up --build
```

### Stopping Services

#### Stop All Services
```bash
docker-compose down
```

#### Stop and Remove Volumes
```bash
docker-compose down -v
```

### Viewing Logs

#### All Services
```bash
docker-compose logs
```

#### Follow Logs (Live)
```bash
docker-compose logs -f
```

#### Specific Service
```bash
docker-compose logs backend-python
docker-compose logs frontend-vue
```

## 🚀 Quick Start Guide

### Step 1: Clone and Navigate
```bash
git clone <repository-url>
cd brutus
```

### Step 2: Set Up Environment Variables
```bash
# Copy example environment files
cp src/frontend/vue/.env.example src/frontend/vue/.env

# Edit the .env file with your Azure credentials
# (Use your favorite text editor)
nano src/frontend/vue/.env
```

### Step 3: Build and Start
```bash
# Build and start all services
docker-compose up --build

# Or use our helper script
./scripts/docker-scripts.sh start
```

### Step 4: Access Applications

Once running, you can access:

- **Vue Frontend**: http://localhost:5173
- **SvelteKit Frontend**: http://localhost:5174  
- **Python Backend**: http://localhost:8000
- **Go Backend**: http://localhost:8080

### Step 5: Stop When Done
```bash
# Stop all services
docker-compose down

# Or use Ctrl+C if running in foreground
```

## 🔧 Common Commands

### Docker Commands

```bash
# List running containers
docker ps

# List all containers (including stopped)
docker ps -a

# List images
docker images

# Remove container
docker rm <container-name>

# Remove image
docker rmi <image-name>

# View container logs
docker logs <container-name>

# Execute command in running container
docker exec -it <container-name> /bin/bash
```

### Docker Compose Commands

```bash
# Start services
docker-compose up

# Start in background
docker-compose up -d

# Stop services
docker-compose down

# View running services
docker-compose ps

# View logs
docker-compose logs

# Rebuild and start
docker-compose up --build

# Remove everything (containers, networks, volumes)
docker-compose down --volumes --remove-orphans
```

### Project-Specific Helper Commands

```bash
# Validate Docker setup
./scripts/validate-docker.sh

# Build all images
./scripts/docker-scripts.sh build

# Start all services
./scripts/docker-scripts.sh start

# Stop all services
./scripts/docker-scripts.sh stop

# Clean up everything
./scripts/docker-scripts.sh clean
```

## 🐛 Troubleshooting

### Common Issues and Solutions

#### 1. "Docker command not found"
**Problem**: Docker is not installed or not in PATH
**Solution**: 
- Reinstall Docker Desktop
- Restart your terminal
- On Linux, check if Docker service is running: `sudo systemctl start docker`

#### 2. "Permission denied"
**Problem**: User doesn't have Docker permissions (Linux)
**Solution**:
```bash
sudo usermod -aG docker $USER
# Log out and back in
```

#### 3. "Port already in use"
**Problem**: Another service is using the same port
**Solution**:
```bash
# Find what's using the port
lsof -i :8000  # or whatever port
# Kill the process or change port in docker-compose.yml
```

#### 4. "Cannot connect to Docker daemon"
**Problem**: Docker Desktop is not running
**Solution**: Start Docker Desktop application

#### 5. "Out of space" or "No space left on device"
**Problem**: Docker is using too much disk space
**Solution**:
```bash
# Clean up unused containers, images, networks
docker system prune -a

# Remove unused volumes
docker volume prune
```

#### 6. Container Exits Immediately
**Problem**: Application crashes on startup
**Solution**:
```bash
# Check logs
docker-compose logs <service-name>

# Run container interactively
docker run -it <image-name> /bin/bash
```

### Getting Help

If you encounter issues:

1. **Check logs**: `docker-compose logs <service-name>`
2. **Validate setup**: `./scripts/validate-docker.sh`
3. **Clean and rebuild**: 
   ```bash
   docker-compose down
   docker system prune -f
   docker-compose up --build
   ```

## 💡 Tips for Beginners

### Do's ✅

- **Start Docker Desktop** before running commands
- **Use helper scripts** - they're designed for beginners
- **Check logs** when something doesn't work
- **Stop services** when you're done to free resources
- **Keep Docker Desktop updated**

### Don'ts ❌

- **Don't panic** if something fails - check logs first
- **Don't run multiple instances** of the same service on same ports
- **Don't delete system containers** unless you know what you're doing
- **Don't ignore error messages** - they often contain the solution

### Development Workflow

1. **Start your development environment**:
   ```bash
   docker-compose up backend-python frontend-vue
   ```

2. **Make code changes** (Docker will auto-reload in dev mode)

3. **View logs** to debug:
   ```bash
   docker-compose logs -f backend-python
   ```

4. **Stop when done**:
   ```bash
   docker-compose down
   ```

### Resource Management

Docker can use significant resources. Monitor usage:

- **Docker Desktop**: Check resource usage in settings
- **Command line**: `docker stats` shows real-time usage
- **Clean up regularly**: `docker system prune` removes unused data

### Learning More

- **Official Docker Tutorial**: https://www.docker.com/101-tutorial
- **Docker Compose Documentation**: https://docs.docker.com/compose/
- **Project-specific help**: Check `docs/TROUBLESHOOTING.md`

## 🎉 Success Indicators

You know everything is working when:

- ✅ Docker Desktop shows green status
- ✅ `docker-compose ps` shows all services as "Up"
- ✅ You can access frontends at http://localhost:5173 and http://localhost:5174
- ✅ Backend APIs respond at http://localhost:8000 and http://localhost:8080
- ✅ No error messages in `docker-compose logs`

## 📞 Getting Help

If you're still stuck:

1. **Check existing documentation**:
   - `docs/TROUBLESHOOTING.md`
   - `docs/API.md`
   - Individual service README files

2. **Use validation script**:
   ```bash
   ./scripts/validate-docker.sh
   ```

3. **Create a GitHub issue** with:
   - Your operating system
   - Docker version (`docker --version`)
   - Error messages from logs
   - Steps you tried

Remember: Everyone was a Docker beginner once! Take your time, read error messages carefully, and don't hesitate to ask questions.

---

**Happy Dockerizing! 🐳✨**