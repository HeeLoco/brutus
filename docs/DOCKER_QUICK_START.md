# 🚀 Docker Quick Start - Brutus Project

**For complete beginners who just want to get started quickly!**

## ⚡ Super Quick Start (3 minutes)

### 1. Install Docker
- **Windows/Mac**: Download [Docker Desktop](https://www.docker.com/products/docker-desktop)
- **Linux**: `sudo apt install docker.io docker-compose`

### 2. Start Docker
- Launch Docker Desktop (Windows/Mac)
- Or start service: `sudo systemctl start docker` (Linux)

### 3. Clone and Run
```bash
git clone <your-repo-url>
cd brutus
./scripts/docker-scripts.sh dev
```

### 4. Open Your Browser
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080

That's it! 🎉

## 🔧 Essential Commands

### Using Our Helper Script (Recommended for Beginners)
```bash
# Start development environment
./scripts/docker-scripts.sh dev

# View what's running
./scripts/docker-scripts.sh status

# See application logs
./scripts/docker-scripts.sh logs

# Stop everything
./scripts/docker-scripts.sh stop

# Get help
./scripts/docker-scripts.sh help
```

### Direct Docker Compose Commands
```bash
# Start all services
docker-compose up

# Start in background
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs

# Rebuild and start
docker-compose up --build
```

## 🆘 Common Issues & Quick Fixes

| Problem | Quick Fix |
|---------|-----------|
| "Docker command not found" | Install Docker Desktop |
| "Docker is not running" | Start Docker Desktop |
| "Port already in use" | Stop other services or change ports |
| "Permission denied" (Linux) | `sudo usermod -aG docker $USER` then logout/login |
| Container won't start | Check logs: `./scripts/docker-scripts.sh logs` |

## 📍 Important URLs

Once running, access these URLs:

| Service | URL | Description |
|---------|-----|-------------|
| Vue Frontend | http://localhost:5173 | Main web interface |
| SvelteKit Frontend | http://localhost:5174 | Alternative interface |
| Go Backend | http://localhost:8080 | API server |
| Health Check | http://localhost:8080/health | Server status |

## 🎯 What Each Command Does

### Development Mode (`./scripts/docker-scripts.sh dev`)
- Starts Vue frontend with hot reload
- Starts Go backend API
- Perfect for coding and testing

### Production Mode (`./scripts/docker-scripts.sh prod`)
- Builds optimized versions
- Runs like a real deployment
- Use for final testing

### Build Only (`./scripts/docker-scripts.sh build`)
- Creates Docker images
- Doesn't start anything
- Good for preparing deployment

## 🛑 Before You Start

Make sure you have:
- [ ] Docker installed and running
- [ ] At least 4GB free RAM
- [ ] Internet connection (for first-time setup)
- [ ] Ports 5173, 8080 available

## 💡 Pro Tips for Beginners

1. **Always check if Docker is running** before running commands
2. **Use our helper script** - it's designed for beginners
3. **Check logs when things break**: `./scripts/docker-scripts.sh logs`
4. **Stop services when done** to free up resources
5. **Don't panic** - most issues have simple fixes

## 🆔 Need More Help?

- **Full guide**: Read `docs/DOCKER_GUIDE.md`
- **Troubleshooting**: Check `docs/TROUBLESHOOTING.md`
- **Validation**: Run `./scripts/validate-docker.sh`
- **Stuck?**: Create a GitHub issue with error details

Remember: Docker might seem complex, but with our scripts, it's just a few commands! 🚀