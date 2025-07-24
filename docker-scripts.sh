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
    docker-compose --profile dev up --build -d
    
    print_success "Development environment started!"
    echo -e "🌐 Vue Frontend: ${BLUE}http://localhost:5173${NC}"
    echo -e "🚀 Go Backend: ${BLUE}http://localhost:8080${NC}"
    echo -e "📊 Health Check: ${BLUE}http://localhost:8080/health${NC}"
}

# Development with SvelteKit
dev_svelte() {
    print_header "Starting SvelteKit Development Environment"
    check_docker
    
    print_warning "This will start the SvelteKit frontend in development mode"
    docker-compose --profile dev --profile svelte up --build -d
    
    print_success "SvelteKit development environment started!"
    echo -e "🌐 Vue Frontend: ${BLUE}http://localhost:5173${NC}"
    echo -e "🔥 Svelte Frontend: ${BLUE}http://localhost:5174${NC}"
    echo -e "🚀 Go Backend: ${BLUE}http://localhost:8080${NC}"
}

# Production Environment
prod() {
    print_header "Starting Production Environment"
    check_docker
    
    print_warning "This will build and start the production environment"
    docker-compose --profile prod up --build -d
    
    print_success "Production environment started!"
    echo -e "🌐 Frontend: ${BLUE}http://localhost${NC}"
    echo -e "🚀 Backend: ${BLUE}http://localhost:8080${NC}"
}

# Build only (no start)
build() {
    print_header "Building All Images"
    check_docker
    
    print_warning "Building all Docker images..."
    docker-compose build
    
    print_success "All images built successfully!"
}

# Stop all services
stop() {
    print_header "Stopping All Services"
    docker-compose down
    print_success "All services stopped!"
}

# Clean up (remove containers, images, volumes)
clean() {
    print_header "Cleaning Up Docker Resources"
    
    read -p "This will remove all containers, images, and volumes. Are you sure? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        docker-compose down -v --rmi all --remove-orphans
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
        docker-compose logs -f
    else
        print_header "Showing Logs for $service"
        docker-compose logs -f "$service"
    fi
}

# Show status
status() {
    print_header "Docker Services Status"
    docker-compose ps
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
    echo "  dev          - Start development environment (Vue + Go backend)"
    echo "  dev-svelte   - Start development with SvelteKit"
    echo "  prod         - Start production environment"
    echo "  build        - Build all Docker images"
    echo "  stop         - Stop all services"
    echo "  clean        - Remove all containers, images, and volumes"
    echo "  logs [service] - Show logs (all services or specific service)"
    echo "  status       - Show services and images status"
    echo "  test         - Test all endpoints"
    echo "  help         - Show this help message"
    echo
    echo "Examples:"
    echo "  $0 dev                 # Start development"
    echo "  $0 logs frontend-dev   # Show frontend logs"
    echo "  $0 prod               # Start production"
    echo "  $0 clean              # Clean everything"
}

# Main script logic
case "${1:-help}" in
    dev)
        dev
        ;;
    dev-svelte)
        dev_svelte
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