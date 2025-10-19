#!/bin/bash
# ============================================
# One-Click Setup Script
# ============================================
# This script automates the initial setup

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_banner() {
    echo -e "${BLUE}"
    echo "=========================================="
    echo "  Security & Infrastructure Tools Set"
    echo "           Setup Script v1.0"
    echo "=========================================="
    echo -e "${NC}"
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    local missing=0
    
    # Check Docker
    if command -v docker &>/dev/null; then
        local docker_version=$(docker --version | awk '{print $3}' | sed 's/,//')
        log_info "Docker: $docker_version ✓"
    else
        log_error "Docker is not installed"
        ((missing++))
    fi
    
    # Check Docker Compose
    if command -v docker-compose &>/dev/null; then
        local compose_version=$(docker-compose version --short)
        log_info "Docker Compose: $compose_version ✓"
    else
        log_error "Docker Compose is not installed"
        ((missing++))
    fi
    
    # Check Make
    if command -v make &>/dev/null; then
        log_info "Make: installed ✓"
    else
        log_warn "Make is not installed (optional)"
    fi
    
    if [ "$missing" -gt 0 ]; then
        log_error "Please install missing prerequisites"
        exit 1
    fi
    
    log_info "Prerequisites check passed ✓"
}

setup_environment() {
    log_info "Setting up environment..."
    
    cd "$PROJECT_ROOT" || exit 1
    
    # Check if .env exists
    if [ -f ".env" ]; then
        log_warn ".env file already exists"
        read -p "Do you want to overwrite it? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "Skipping .env creation"
            return 0
        fi
    fi
    
    # Create .env from template
    if [ -f ".env.template" ]; then
        cp .env.template .env
        log_info "Created .env from template"
        
        # Generate random passwords
        local db_pass=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-25)
        local vault_token=$(openssl rand -hex 16)
        
        # Update .env with generated passwords
        sed -i.bak "s/DB_PASSWORD=changeme/DB_PASSWORD=$db_pass/" .env
        sed -i.bak "s/VAULT_TOKEN=root/VAULT_TOKEN=$vault_token/" .env
        rm -f .env.bak
        
        log_info "Generated secure passwords ✓"
        log_warn "Please review .env file and adjust settings as needed"
    else
        log_error ".env.template not found"
        exit 1
    fi
}

create_directories() {
    log_info "Creating required directories..."
    
    cd "$PROJECT_ROOT" || exit 1
    
    mkdir -p backups
    mkdir -p docs/zh-TW docs/en docs/images
    mkdir -p scripts/init scripts/parsers
    mkdir -p examples
    
    log_info "Directories created ✓"
}

start_services() {
    log_info "Starting services..."
    
    cd "$PROJECT_ROOT/Docker/compose" || exit 1
    
    docker-compose up -d
    
    log_info "Waiting for services to start (30 seconds)..."
    sleep 30
    
    log_info "Services started ✓"
}

verify_deployment() {
    log_info "Verifying deployment..."
    
    cd "$PROJECT_ROOT" || exit 1
    
    if bash scripts/health-check.sh; then
        log_info "Deployment verification passed ✓"
        return 0
    else
        log_error "Deployment verification failed"
        return 1
    fi
}

print_summary() {
    echo ""
    echo -e "${GREEN}=========================================="
    echo "       Setup Completed Successfully!"
    echo "==========================================${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Review .env file: nano .env"
    echo "  2. Access services:"
    echo "     - Traefik Dashboard: http://localhost:8080"
    echo "     - Vault UI: http://localhost:8200"
    echo "     - ArgoCD UI: http://localhost:8081"
    echo "  3. Run your first scan:"
    echo "     cd Make_Files"
    echo "     make scan-nuclei TARGET=https://example.com"
    echo ""
    echo "Documentation: README.md"
    echo ""
}

# Main execution
main() {
    print_banner
    
    check_prerequisites
    setup_environment
    create_directories
    start_services
    verify_deployment
    
    print_summary
}

# Handle errors
trap 'log_error "Setup failed at line $LINENO"' ERR

# Run main
main "$@"

