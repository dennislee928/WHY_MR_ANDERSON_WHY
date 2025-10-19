#!/bin/bash

# Cloudflare Containers Build and Deploy Script
# This script builds and deploys all containers for the security platform

set -e

# Configuration
PROJECT_NAME="security-platform-containers"
REGISTRY="ghcr.io/dennislee928"
VERSION=${1:-latest}
ENVIRONMENT=${2:-production}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Check dependencies
check_dependencies() {
    print_header "Checking Dependencies"
    
    if ! command -v docker &> /dev/null; then
        print_error "Docker not found. Please install Docker."
        exit 1
    else
        print_success "Docker found"
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose not found. Please install Docker Compose."
        exit 1
    else
        print_success "Docker Compose found"
    fi
    
    if ! command -v wrangler &> /dev/null; then
        print_warning "Wrangler not found. Installing..."
        npm install -g wrangler
    else
        print_success "Wrangler found"
    fi
}

# Build containers
build_containers() {
    print_header "Building Containers"
    
    # Build Backend API
    print_success "Building Backend API container..."
    docker build -t ${REGISTRY}/${PROJECT_NAME}-backend-api:${VERSION} ./containers/backend-api/
    
    # Build AI/Quantum
    print_success "Building AI/Quantum container..."
    docker build -t ${REGISTRY}/${PROJECT_NAME}-ai-quantum:${VERSION} ./containers/ai-quantum/
    
    # Build Security Tools
    print_success "Building Security Tools container..."
    docker build -t ${REGISTRY}/${PROJECT_NAME}-security-tools:${VERSION} ./containers/security-tools/
    
    # Build Database
    print_success "Building Database container..."
    docker build -t ${REGISTRY}/${PROJECT_NAME}-database:${VERSION} ./containers/database/
    
    # Build Monitoring
    print_success "Building Monitoring container..."
    docker build -t ${REGISTRY}/${PROJECT_NAME}-monitoring:${VERSION} ./containers/monitoring/
    
    print_success "All containers built successfully"
}

# Push containers to registry
push_containers() {
    print_header "Pushing Containers to Registry"
    
    # Login to registry (if needed)
    if [ "$REGISTRY" != "localhost" ]; then
        print_success "Logging in to registry..."
        echo $GITHUB_TOKEN | docker login ghcr.io -u $GITHUB_USERNAME --password-stdin
    fi
    
    # Push all containers
    docker push ${REGISTRY}/${PROJECT_NAME}-backend-api:${VERSION}
    docker push ${REGISTRY}/${PROJECT_NAME}-ai-quantum:${VERSION}
    docker push ${REGISTRY}/${PROJECT_NAME}-security-tools:${VERSION}
    docker push ${REGISTRY}/${PROJECT_NAME}-database:${VERSION}
    docker push ${REGISTRY}/${PROJECT_NAME}-monitoring:${VERSION}
    
    print_success "All containers pushed successfully"
}

# Deploy to Cloudflare Workers
deploy_workers() {
    print_header "Deploying to Cloudflare Workers"
    
    # Update wrangler.toml with container images
    sed -i "s|image = \"security-platform/backend-api\"|image = \"${REGISTRY}/${PROJECT_NAME}-backend-api:${VERSION}\"|g" wrangler-containers.toml
    sed -i "s|image = \"security-platform/ai-quantum\"|image = \"${REGISTRY}/${PROJECT_NAME}-ai-quantum:${VERSION}\"|g" wrangler-containers.toml
    sed -i "s|image = \"security-platform/security-tools\"|image = \"${REGISTRY}/${PROJECT_NAME}-security-tools:${VERSION}\"|g" wrangler-containers.toml
    sed -i "s|image = \"security-platform/database\"|image = \"${REGISTRY}/${PROJECT_NAME}-database:${VERSION}\"|g" wrangler-containers.toml
    sed -i "s|image = \"security-platform/monitoring\"|image = \"${REGISTRY}/${PROJECT_NAME}-monitoring:${VERSION}\"|g" wrangler-containers.toml
    
    # Deploy to Cloudflare Workers
    wrangler deploy --config wrangler-containers.toml --env ${ENVIRONMENT}
    
    print_success "Deployed to Cloudflare Workers successfully"
}

# Test deployment
test_deployment() {
    print_header "Testing Deployment"
    
    # Get deployment URL
    DEPLOYMENT_URL=$(wrangler deployments list --config wrangler-containers.toml --env ${ENVIRONMENT} | head -n 2 | tail -n 1 | awk '{print $2}')
    
    if [ -z "$DEPLOYMENT_URL" ]; then
        print_error "Could not get deployment URL"
        return 1
    fi
    
    print_success "Testing deployment at: $DEPLOYMENT_URL"
    
    # Test health endpoint
    if curl -f "${DEPLOYMENT_URL}/api/v1/containers/health" > /dev/null 2>&1; then
        print_success "Health check passed"
    else
        print_error "Health check failed"
        return 1
    fi
    
    # Test services endpoint
    if curl -f "${DEPLOYMENT_URL}/api/v1/services" > /dev/null 2>&1; then
        print_success "Services endpoint working"
    else
        print_error "Services endpoint failed"
        return 1
    fi
    
    print_success "All tests passed"
}

# Cleanup
cleanup() {
    print_header "Cleanup"
    
    # Remove local images to save space
    docker rmi ${REGISTRY}/${PROJECT_NAME}-backend-api:${VERSION} || true
    docker rmi ${REGISTRY}/${PROJECT_NAME}-ai-quantum:${VERSION} || true
    docker rmi ${REGISTRY}/${PROJECT_NAME}-security-tools:${VERSION} || true
    docker rmi ${REGISTRY}/${PROJECT_NAME}-database:${VERSION} || true
    docker rmi ${REGISTRY}/${PROJECT_NAME}-monitoring:${VERSION} || true
    
    print_success "Cleanup completed"
}

# Main execution
main() {
    print_header "Cloudflare Containers Deployment"
    echo "Project: $PROJECT_NAME"
    echo "Registry: $REGISTRY"
    echo "Version: $VERSION"
    echo "Environment: $ENVIRONMENT"
    echo ""
    
    case "${3:-all}" in
        "build")
            check_dependencies
            build_containers
            ;;
        "push")
            check_dependencies
            push_containers
            ;;
        "deploy")
            check_dependencies
            deploy_workers
            ;;
        "test")
            test_deployment
            ;;
        "all"|"")
            check_dependencies
            build_containers
            push_containers
            deploy_workers
            test_deployment
            cleanup
            ;;
        "help"|"-h"|"--help")
            echo "Usage: $0 [version] [environment] [action]"
            echo ""
            echo "Arguments:"
            echo "  version     - Container version (default: latest)"
            echo "  environment - Deployment environment (default: production)"
            echo "  action      - Action to perform (default: all)"
            echo ""
            echo "Actions:"
            echo "  build       - Build containers only"
            echo "  push        - Push containers to registry"
            echo "  deploy      - Deploy to Cloudflare Workers"
            echo "  test        - Test deployment"
            echo "  all         - Build, push, deploy, and test (default)"
            echo "  help        - Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0 v1.0.0 production build"
            echo "  $0 latest staging deploy"
            echo "  $0 v2.0.0 production all"
            ;;
        *)
            print_error "Unknown action: $3"
            echo "Use '$0 help' for usage information"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
