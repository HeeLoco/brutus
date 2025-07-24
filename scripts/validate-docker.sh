#!/bin/bash

# Docker Configuration Validation Script

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

print_check() {
    echo -e "${BLUE}🔍 Checking: $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  Docker Configuration Validation${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_header

# Check if Docker is available
print_check "Docker installation"
if command -v docker >/dev/null 2>&1; then
    print_success "Docker is installed"
    docker --version
else
    print_error "Docker is not installed"
    exit 1
fi

# Check if Docker daemon is running
print_check "Docker daemon"
if docker info >/dev/null 2>&1; then
    print_success "Docker daemon is running"
else
    print_warning "Docker daemon is not running - start Docker Desktop"
fi

# Check Docker Compose
print_check "Docker Compose"
if docker-compose --version >/dev/null 2>&1; then
    print_success "Docker Compose is available"
    docker-compose --version
else
    print_error "Docker Compose is not available"
fi

# Validate Dockerfile syntax
print_check "Vue Dockerfile (production)"
if [ -f "src/frontend/vue/Dockerfile" ]; then
    # Basic syntax validation
    if docker build --dry-run -f src/frontend/vue/Dockerfile src/frontend/vue >/dev/null 2>&1; then
        print_success "Vue production Dockerfile syntax is valid"
    else
        print_warning "Vue production Dockerfile has issues (Docker not running?)"
    fi
else
    print_error "Vue production Dockerfile not found"
fi

print_check "Vue Dockerfile (development)"
if [ -f "src/frontend/vue/Dockerfile.dev" ]; then
    print_success "Vue development Dockerfile found"
else
    print_error "Vue development Dockerfile not found"
fi

print_check "SvelteKit Dockerfile (production)"
if [ -f "src/frontend/svelte/Dockerfile" ]; then
    print_success "SvelteKit production Dockerfile found"
else
    print_error "SvelteKit production Dockerfile not found"
fi

print_check "SvelteKit Dockerfile (development)"
if [ -f "src/frontend/svelte/Dockerfile.dev" ]; then
    print_success "SvelteKit development Dockerfile found"
else
    print_error "SvelteKit development Dockerfile not found"
fi

# Validate docker-compose.yml
print_check "Docker Compose configuration"
if [ -f "docker-compose.yml" ]; then
    if docker-compose config >/dev/null 2>&1; then
        print_success "docker-compose.yml syntax is valid"
    else
        print_warning "docker-compose.yml has syntax issues"
        docker-compose config 2>&1 | head -5
    fi
else
    print_error "docker-compose.yml not found"
fi

# Check required files
print_check "Required configuration files"
files=(
    "src/frontend/vue/package.json"
    "src/frontend/vue/nginx.conf"
    "src/frontend/svelte/package.json"
    "src/frontend/svelte/nginx.conf"
    ".env.docker.example"
)

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        print_success "$file exists"
    else
        print_error "$file missing"
    fi
done

# Check frontend dependencies
print_check "Vue frontend dependencies"
if [ -f "src/frontend/vue/package.json" ]; then
    cd src/frontend/vue
    if [ -f "package-lock.json" ] || [ -f "yarn.lock" ]; then
        print_success "Vue dependencies are locked"
    else
        print_warning "Vue dependencies not locked - run npm install"
    fi
    cd ../../..
fi

print_check "SvelteKit frontend dependencies"
if [ -f "src/frontend/svelte/package.json" ]; then
    cd src/frontend/svelte
    if [ -f "package-lock.json" ] || [ -f "yarn.lock" ]; then
        print_success "SvelteKit dependencies are locked"
    else
        print_warning "SvelteKit dependencies not locked - run npm install"
    fi
    cd ../../..
fi

# Validate script permissions
print_check "Script permissions"
if [ -x "docker-scripts.sh" ]; then
    print_success "docker-scripts.sh is executable"
else
    print_warning "docker-scripts.sh is not executable - run: chmod +x docker-scripts.sh"
fi

# Summary
echo
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Validation Summary${NC}"
echo -e "${BLUE}========================================${NC}"

if docker info >/dev/null 2>&1; then
    echo -e "${GREEN}🚀 Ready to run Docker commands!${NC}"
    echo
    echo "Quick start:"
    echo "  ./docker-scripts.sh dev     # Start development"
    echo "  ./docker-scripts.sh prod    # Start production"
    echo "  ./docker-scripts.sh help    # Show all commands"
else
    echo -e "${YELLOW}⚠️  Start Docker Desktop first, then run:${NC}"
    echo "  ./docker-scripts.sh dev     # Start development"
fi

echo
echo "For detailed instructions, see: DOCKER.md"