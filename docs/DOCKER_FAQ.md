# 🤔 Docker FAQ - Common Questions & Issues

## 🚨 Most Common Issues

### ❌ "Docker command not found"
**What it means**: Docker is not installed or not in your system PATH.

**Fix**:
- **Windows/Mac**: Download and install [Docker Desktop](https://www.docker.com/products/docker-desktop)
- **Linux**: `sudo apt install docker.io docker-compose`
- Restart your terminal after installation

### ❌ "Cannot connect to the Docker daemon"
**What it means**: Docker Desktop is not running.

**Fix**:
- **Windows/Mac**: Start Docker Desktop from your applications
- **Linux**: `sudo systemctl start docker`
- Wait for Docker Desktop to show "Running" status

### ❌ "Permission denied while trying to connect to Docker"
**What it means**: Your user account doesn't have Docker permissions (Linux only).

**Fix**:
```bash
sudo usermod -aG docker $USER
# Log out and back in, then test:
docker run hello-world
```

### ❌ "Port is already allocated" or "Address already in use"
**What it means**: Another service is using the same port.

**Fix**:
```bash
# Find what's using port 8080
lsof -i :8080

# Kill the process (replace PID with actual number)
kill -9 <PID>

# Or change ports in docker-compose.yml
```

### ❌ "No space left on device"
**What it means**: Docker is using too much disk space.

**Fix**:
```bash
# Clean up unused Docker resources
docker system prune -a

# Remove unused volumes
docker volume prune

# Check disk usage
docker system df
```

### ❌ "Stop" command shows warnings and containers still running
**What it means**: Multiple docker-compose files exist and containers from different setups
**Fix**:
```bash
# Use our improved script that handles both files
./scripts/docker-scripts.sh stop

# For stubborn containers, use nuclear option
./scripts/docker-scripts.sh force-stop

# Check what's still running
./scripts/docker-scripts.sh status
```

### ❌ Container exits immediately
**What it means**: Application crashes on startup
**Fix**:
```bash
# Check logs for error messages
./scripts/docker-scripts.sh logs <service-name>

# Run container interactively for debugging
docker run -it <image-name> /bin/bash
```

## 🤷 Frequently Asked Questions

### Q: Do I need to install Go/Python/Node.js?
**A**: No! Docker containers include everything. You only need Docker.

### Q: How much disk space does this use?
**A**: Initial setup: ~2-3GB. After development: ~5-10GB depending on usage.

### Q: Can I run this on Windows?
**A**: Yes! Docker Desktop works on Windows 10/11. WSL2 is recommended.

### Q: Which services should I start?
**A**: For development: `./scripts/docker-scripts.sh dev` (Vue frontend + Go backend)

### Q: How do I see what's running?
**A**: Use `./scripts/docker-scripts.sh status` or `docker-compose ps`

### Q: The containers start but I can't access the website
**A**: Check:
1. Docker Desktop is running
2. No other services on same ports
3. Firewall isn't blocking ports
4. Wait 30-60 seconds for services to fully start

### Q: How do I update to latest version?
**A**: 
```bash
git pull
./scripts/docker-scripts.sh stop
./scripts/docker-scripts.sh build
./scripts/docker-scripts.sh dev
```

### Q: Can I develop without Docker?
**A**: Yes, but you'll need to install all languages (Go, Python, Node.js) manually. Docker is much easier.

### Q: How do I see error messages?
**A**: Check logs: `./scripts/docker-scripts.sh logs [service-name]`

### Q: The build takes forever
**A**: First build downloads images (~1-2GB). Subsequent builds are much faster.

### Q: Can I run just the frontend?
**A**: Yes, but you'll need a backend. Start with: `./scripts/docker-scripts.sh dev`

## 🛠️ Development Workflow

### Starting Development
```bash
# Start everything
./scripts/docker-scripts.sh dev

# Check it's working
./scripts/docker-scripts.sh status

# Open browser to http://localhost:5173
```

### Making Changes
1. Edit code in your editor
2. Changes auto-reload in development mode
3. Check logs if something breaks: `./scripts/docker-scripts.sh logs`

### Stopping Work
```bash
# Stop all services
./scripts/docker-scripts.sh stop

# Or use Ctrl+C if running in foreground
```

## 🔍 Debugging Tips

### Check Service Status
```bash
# Overview of all services
./scripts/docker-scripts.sh status

# Detailed view
docker-compose ps
```

### View Logs
```bash
# All services
./scripts/docker-scripts.sh logs

# Specific service
./scripts/docker-scripts.sh logs frontend-dev
./scripts/docker-scripts.sh logs backend-go
```

### Test Connections
```bash
# Test all endpoints
./scripts/docker-scripts.sh test

# Manual testing
curl http://localhost:8080/health
curl http://localhost:5173
```

### Interactive Container Access
```bash
# Get shell access to a container
docker exec -it brutus-backend-go-1 /bin/bash

# Or with sh if bash not available
docker exec -it brutus-backend-go-1 /bin/sh
```

## 🧹 Cleanup & Maintenance

### Regular Cleanup
```bash
# Remove unused containers and images
docker system prune

# More aggressive cleanup
docker system prune -a
```

### Complete Reset
```bash
# Nuclear option - removes everything
./scripts/docker-scripts.sh clean

# Rebuild from scratch
./scripts/docker-scripts.sh build
./scripts/docker-scripts.sh dev
```

### Check Resource Usage
```bash
# See what's using space
docker system df

# Monitor resource usage
docker stats
```

## 🆘 Still Stuck?

If nothing here helps:

1. **Check our guides**:
   - [Docker Quick Start](DOCKER_QUICK_START.md)
   - [Complete Docker Guide](DOCKER_GUIDE.md)
   - [Troubleshooting Guide](TROUBLESHOOTING.md)

2. **Run validation**:
   ```bash
   ./scripts/validate-docker.sh
   ```

3. **Create an issue** with:
   - Your operating system
   - Docker version: `docker --version`
   - Error messages from logs
   - Steps you tried

4. **Include this information**:
   ```bash
   # System info
   docker version
   docker-compose version
   
   # What's running
   docker ps -a
   
   # Recent logs
   ./scripts/docker-scripts.sh logs | tail -50
   ```

Remember: Every Docker expert was once a beginner! Don't give up! 🚀