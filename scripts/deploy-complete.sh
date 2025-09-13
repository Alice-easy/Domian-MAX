#!/bin/bash

# Domain MAX - Complete Deployment Script
# This script handles the complete deployment process

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if running on Windows with WSL
check_environment() {
    log_info "Checking environment..."
    
    if command -v docker >/dev/null 2>&1; then
        log_success "Docker is installed"
    else
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if command -v docker-compose >/dev/null 2>&1; then
        log_success "Docker Compose is installed"
    else
        log_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
}

# Setup environment
setup_environment() {
    log_info "Setting up environment..."
    
    cd "$(dirname "$0")/../deployments"
    
    # Copy environment file if it doesn't exist
    if [ ! -f .env ]; then
        if [ -f .env.example ]; then
            cp .env.example .env
            log_warning "Copied .env.example to .env. Please configure your settings!"
            log_warning "You must set secure passwords and secrets before continuing."
            echo ""
            echo "Required settings to configure in .env:"
            echo "- DB_PASSWORD (database password)"
            echo "- JWT_SECRET (JWT secret key)"
            echo "- ENCRYPTION_KEY (32-character encryption key)"
            echo ""
            read -p "Press Enter after configuring .env file..."
        else
            log_error ".env.example file not found"
            exit 1
        fi
    else
        log_success "Environment file (.env) exists"
    fi
    
    # Validate required environment variables
    source .env
    
    if [ -z "$DB_PASSWORD" ] || [ "$DB_PASSWORD" = "your_secure_database_password_here" ]; then
        log_error "Please set a secure DB_PASSWORD in .env file"
        exit 1
    fi
    
    if [ -z "$JWT_SECRET" ] || [ "$JWT_SECRET" = "your_jwt_secret_key_here_minimum_32_characters" ]; then
        log_error "Please set a secure JWT_SECRET in .env file"
        exit 1
    fi
    
    if [ -z "$ENCRYPTION_KEY" ] || [ "$ENCRYPTION_KEY" = "your_32_character_encryption_key_12" ]; then
        log_error "Please set a secure ENCRYPTION_KEY in .env file"
        exit 1
    fi
    
    log_success "Environment configuration validated"
}

# Generate SSL certificates
setup_ssl() {
    log_info "Setting up SSL certificates..."
    
    if [ ! -f ssl/nginx-selfsigned.crt ]; then
        log_info "Generating SSL certificates..."
        ../scripts/generate-ssl.sh
        log_success "SSL certificates generated"
    else
        log_success "SSL certificates already exist"
    fi
}

# Build and start services
deploy_services() {
    log_info "Building and deploying services..."
    
    # Build the application
    log_info "Building application image..."
    docker-compose build --no-cache
    log_success "Application image built successfully"
    
    # Start services
    log_info "Starting services..."
    docker-compose up -d
    
    # Wait for services to be ready
    log_info "Waiting for services to be ready..."
    sleep 30
    
    # Check service health
    check_services_health
}

# Check service health
check_services_health() {
    log_info "Checking service health..."
    
    # Check if containers are running
    if docker-compose ps | grep -q "Up"; then
        log_success "Services are running"
    else
        log_error "Some services failed to start"
        docker-compose logs --tail=50
        exit 1
    fi
    
    # Test database connection
    if docker-compose exec -T db pg_isready -U postgres >/dev/null 2>&1; then
        log_success "Database is ready"
    else
        log_warning "Database is not ready yet"
    fi
    
    # Test application health
    sleep 10
    if curl -sf http://localhost:8080/api/health >/dev/null 2>&1; then
        log_success "Application health check passed"
    else
        log_warning "Application health check failed, but container may still be starting"
    fi
    
    # Test nginx
    if curl -sf http://localhost/health >/dev/null 2>&1; then
        log_success "Nginx health check passed"
    else
        log_warning "Nginx health check failed"
    fi
}

# Show deployment information
show_deployment_info() {
    echo ""
    log_success "🎉 Domain MAX deployment completed!"
    echo ""
    echo "📋 Service Information:"
    echo "  🌐 Frontend: http://localhost (HTTP) / https://localhost (HTTPS)"
    echo "  🔌 API: http://localhost/api"
    echo "  🗄️  Database: localhost:5432"
    echo "  💾 Redis: localhost:6379"
    echo ""
    echo "🔧 Management Commands:"
    echo "  📊 View logs: docker-compose logs -f"
    echo "  🔄 Restart: docker-compose restart"
    echo "  🛑 Stop: docker-compose stop"
    echo "  🗑️  Clean up: docker-compose down -v"
    echo ""
    echo "🏥 Health Checks:"
    echo "  🔍 Application: curl http://localhost:8080/api/health"
    echo "  🔍 Nginx: curl http://localhost/health"
    echo "  🔍 Database: docker-compose exec db pg_isready -U postgres"
    echo ""
    log_warning "Note: If using HTTPS, you may need to accept the self-signed certificate"
}

# Main deployment process
main() {
    echo "🚀 Domain MAX Deployment Script"
    echo "================================="
    echo ""
    
    check_environment
    setup_environment
    setup_ssl
    deploy_services
    show_deployment_info
    
    log_success "Deployment completed successfully! 🎉"
}

# Handle script arguments
case "${1:-}" in
    --check-health)
        check_services_health
        ;;
    --setup-ssl)
        setup_ssl
        ;;
    --help)
        echo "Domain MAX Deployment Script"
        echo ""
        echo "Usage: $0 [option]"
        echo ""
        echo "Options:"
        echo "  (no option)     Full deployment"
        echo "  --check-health  Check service health"
        echo "  --setup-ssl     Generate SSL certificates only"
        echo "  --help          Show this help"
        ;;
    *)
        main
        ;;
esac