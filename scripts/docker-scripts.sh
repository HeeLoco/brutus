#!/bin/bash

# Docker Management Scripts for Brutus Project

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}========================================${NC}"
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

# Check if Docker is running
check_docker() {
    if ! docker info >/dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
}

# Development Environment
dev() {
    print_header "Starting Development Environment"
    check_docker
    
    print_warning "This will start the Vue frontend in development mode with hot reload"
    docker-compose -f docker-compose.yml --profile dev up --build -d
    
    print_success "Development environment started!"
    echo -e "🌐 Vue Frontend: ${BLUE}http://localhost:5173${NC}"
    echo -e "🚀 Go Backend: ${BLUE}http://localhost:8080${NC}"
    echo -e "📊 Health Check: ${BLUE}http://localhost:8080/health${NC}"
}

# Production Environment
prod() {
    print_header "Starting Production Environment"
    check_docker
    
    print_warning "This will build and start the production environment"
    docker-compose -f docker-compose.yml --profile prod up --build -d
    
    print_success "Production environment started!"
    echo -e "🌐 Frontend: ${BLUE}http://localhost${NC}"
    echo -e "🚀 Backend: ${BLUE}http://localhost:8080${NC}"
}

# Build only (no start)
build() {
    print_header "Building All Images"
    check_docker
    
    print_warning "Building all Docker images..."
    docker-compose -f docker-compose.yml build
    
    print_success "All images built successfully!"
}

# Stop all services
stop() {
    print_header "Stopping All Services"
    
    # Stop services from both compose files
    print_warning "Stopping services from docker-compose.yml..."
    docker-compose -f docker-compose.yml down 2>/dev/null || echo "No services from docker-compose.yml to stop"
    
    print_warning "Stopping services from compose.yml..."
    docker-compose -f compose.yml down 2>/dev/null || echo "No services from compose.yml to stop"
    
    # Stop any remaining brutus containers
    print_warning "Stopping any remaining brutus containers..."
    docker ps --filter "name=brutus" --format "{{.Names}}" | xargs -r docker stop 2>/dev/null || echo "No brutus containers to stop"
    
    print_success "All services stopped!"
}

# Force stop all containers (nuclear option)
force_stop() {
    print_header "Force Stopping All Containers"
    
    print_warning "This will forcefully stop and remove ALL running containers"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Force stopping all containers..."
        docker ps -q | xargs -r docker stop
        docker ps -a -q | xargs -r docker rm
        print_success "All containers forcefully stopped and removed!"
    else
        print_warning "Force stop cancelled."
    fi
}

# Clean up (remove containers, images, volumes)
clean() {
    print_header "Cleaning Up Docker Resources"
    
    read -p "This will remove all containers, images, and volumes. Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Cleaning up docker-compose.yml services..."
        docker-compose -f docker-compose.yml down -v --rmi all --remove-orphans 2>/dev/null || echo "No docker-compose.yml services to clean"
        
        print_warning "Cleaning up compose.yml services..."
        docker-compose -f compose.yml down -v --rmi all --remove-orphans 2>/dev/null || echo "No compose.yml services to clean"
        
        print_warning "Removing any remaining brutus containers and images..."
        docker ps -a --filter "name=brutus" --format "{{.Names}}" | xargs -r docker rm -f 2>/dev/null || echo "No brutus containers to remove"
        docker images --filter "reference=brutus*" --format "{{.Repository}}:{{.Tag}}" | xargs -r docker rmi -f 2>/dev/null || echo "No brutus images to remove"
        
        print_warning "General system cleanup..."
        docker system prune -f
        
        print_success "Cleanup completed!"
    else
        print_warning "Cleanup cancelled."
    fi
}

# Show logs
logs() {
    service=${1:-}
    if [ -z "$service" ]; then
        print_header "Showing All Logs"
        docker-compose -f docker-compose.yml logs -f
    else
        print_header "Showing Logs for $service"
        docker-compose -f docker-compose.yml logs -f "$service"
    fi
}

# Show status
status() {
    print_header "Docker Services Status"
    echo "Services from docker-compose.yml:"
    docker-compose -f docker-compose.yml ps
    echo
    echo "Services from compose.yml:"
    docker-compose -f compose.yml ps 2>/dev/null || echo "No services in compose.yml"
    echo
    print_header "All Running Containers"
    docker ps
    echo
    print_header "Docker Images"
    docker images | grep brutus || echo "No brutus images found"
}

# Test endpoints
test() {
    print_header "Testing Endpoints"
    
    echo "Testing backend health..."
    curl -f http://localhost:8080/health 2>/dev/null && print_success "Backend health check passed" || print_error "Backend health check failed"
    
    echo "Testing frontend..."
    curl -f http://localhost:5173 2>/dev/null && print_success "Frontend accessible" || print_error "Frontend not accessible"
    
    echo "Testing production frontend..."
    curl -f http://localhost 2>/dev/null && print_success "Production frontend accessible" || print_warning "Production frontend not running"
}

# Show help
help() {
    print_header "Brutus Docker Management"
    echo "Usage: $0 <command>"
    echo
    echo "Commands:"
    echo "  dev          - Start development environment (Vue frontend + Go backend)"
    echo "  prod         - Start production environment"
    echo "  build        - Build all Docker images"
    echo "  stop         - Stop all services (handles both compose files)"
    echo "  force-stop   - Force stop ALL containers (nuclear option)"
    echo "  clean        - Remove all containers, images, and volumes"
    echo "  logs [service] - Show logs (all services or specific service)"
    echo "  status       - Show services and images status"
    echo "  test         - Test all endpoints"
    echo "  help         - Show this help message"
    echo
    echo "Examples:"
    echo "  $0 dev                 # Start development"
    echo "  $0 logs frontend-dev   # Show frontend logs"
    echo "  $0 stop               # Stop all services"
    echo "  $0 force-stop         # Force stop everything"
    echo "  $0 clean              # Clean everything"
    echo
    echo "Note: This script handles both docker-compose.yml and compose.yml files"
}

# Main script logic
case "${1:-help}" in
    dev)
        dev
        ;;
    prod)
        prod
        ;;
    build)
        build
        ;;
    stop)
        stop
        ;;
    force-stop)
        force_stop
        ;;
    clean)
        clean
        ;;
    logs)
        logs "$2"
        ;;
    status)
        status
        ;;
    test)
        test
        ;;
    help|--help|-h)
        help
        ;;
    *)
        print_error "Unknown command: $1"
        help
        exit 1
        ;;
esac