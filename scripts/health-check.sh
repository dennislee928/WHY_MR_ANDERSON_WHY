#!/bin/bash
# ============================================
# Health Check Script
# ============================================
# This script checks the health status of all services

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
COMPOSE_DIR="$PROJECT_ROOT/Docker/compose"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Services to check
SERVICES=(
    "postgres:5432"
    "vault:8200"
    "traefik:80"
    "argocd:8080"
)

# Functions
print_header() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}    Security Stack Health Check${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo ""
}

check_docker() {
    echo -n "Checking Docker... "
    if docker info &>/dev/null; then
        echo -e "${GREEN}✓ Running${NC}"
        return 0
    else
        echo -e "${RED}✗ Not running${NC}"
        return 1
    fi
}

check_docker_compose() {
    echo -n "Checking Docker Compose... "
    if command -v docker-compose &>/dev/null; then
        local version=$(docker-compose version --short)
        echo -e "${GREEN}✓ Version $version${NC}"
        return 0
    else
        echo -e "${RED}✗ Not installed${NC}"
        return 1
    fi
}

check_container_status() {
    local container=$1
    echo -n "Checking $container... "
    
    cd "$COMPOSE_DIR" || exit 1
    
    # Check if container exists
    if ! docker-compose ps "$container" &>/dev/null; then
        echo -e "${RED}✗ Not found${NC}"
        return 1
    fi
    
    # Check if container is running
    local status=$(docker-compose ps -q "$container" | xargs docker inspect -f '{{.State.Status}}' 2>/dev/null)
    
    if [ "$status" = "running" ]; then
        # Check health status if available
        local health=$(docker-compose ps -q "$container" | xargs docker inspect -f '{{.State.Health.Status}}' 2>/dev/null)
        
        if [ "$health" = "healthy" ]; then
            echo -e "${GREEN}✓ Healthy${NC}"
            return 0
        elif [ "$health" = "unhealthy" ]; then
            echo -e "${RED}✗ Unhealthy${NC}"
            return 1
        else
            echo -e "${GREEN}✓ Running${NC}"
            return 0
        fi
    else
        echo -e "${RED}✗ $status${NC}"
        return 1
    fi
}

check_port() {
    local service=$1
    local port=$2
    echo -n "Checking $service port $port... "
    
    if nc -z localhost "$port" &>/dev/null; then
        echo -e "${GREEN}✓ Open${NC}"
        return 0
    else
        echo -e "${RED}✗ Closed${NC}"
        return 1
    fi
}

check_postgres() {
    echo -n "Checking PostgreSQL connection... "
    
    cd "$COMPOSE_DIR" || exit 1
    
    if docker-compose exec -T postgres psql -U sectools -d security -c "SELECT 1" &>/dev/null; then
        echo -e "${GREEN}✓ Connected${NC}"
        return 0
    else
        echo -e "${RED}✗ Cannot connect${NC}"
        return 1
    fi
}

check_vault() {
    echo -n "Checking Vault... "
    
    local status=$(curl -s http://localhost:8200/v1/sys/health | jq -r '.sealed' 2>/dev/null)
    
    if [ "$status" = "false" ]; then
        echo -e "${GREEN}✓ Unsealed${NC}"
        return 0
    elif [ "$status" = "true" ]; then
        echo -e "${YELLOW}⚠ Sealed${NC}"
        return 1
    else
        echo -e "${RED}✗ Not responding${NC}"
        return 1
    fi
}

check_traefik() {
    echo -n "Checking Traefik Dashboard... "
    
    local status=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/dashboard/)
    
    if [ "$status" = "200" ] || [ "$status" = "301" ]; then
        echo -e "${GREEN}✓ Accessible${NC}"
        return 0
    else
        echo -e "${RED}✗ HTTP $status${NC}"
        return 1
    fi
}

check_disk_space() {
    echo -n "Checking disk space... "
    
    local usage=$(df -h / | awk 'NR==2 {print $5}' | sed 's/%//')
    
    if [ "$usage" -lt 80 ]; then
        echo -e "${GREEN}✓ ${usage}% used${NC}"
        return 0
    elif [ "$usage" -lt 90 ]; then
        echo -e "${YELLOW}⚠ ${usage}% used${NC}"
        return 1
    else
        echo -e "${RED}✗ ${usage}% used${NC}"
        return 1
    fi
}

check_memory() {
    echo -n "Checking memory usage... "
    
    local usage=$(free | grep Mem | awk '{printf "%.0f", $3/$2 * 100}')
    
    if [ "$usage" -lt 80 ]; then
        echo -e "${GREEN}✓ ${usage}% used${NC}"
        return 0
    elif [ "$usage" -lt 90 ]; then
        echo -e "${YELLOW}⚠ ${usage}% used${NC}"
        return 1
    else
        echo -e "${RED}✗ ${usage}% used${NC}"
        return 1
    fi
}

show_summary() {
    local passed=$1
    local failed=$2
    local total=$((passed + failed))
    
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}              Summary${NC}"
    echo -e "${BLUE}========================================${NC}"
    echo -e "Total checks: $total"
    echo -e "${GREEN}Passed: $passed${NC}"
    
    if [ "$failed" -gt 0 ]; then
        echo -e "${RED}Failed: $failed${NC}"
    else
        echo -e "${GREEN}Failed: $failed${NC}"
    fi
    
    echo ""
    
    if [ "$failed" -eq 0 ]; then
        echo -e "${GREEN}✓ All checks passed!${NC}"
        return 0
    else
        echo -e "${RED}✗ Some checks failed!${NC}"
        return 1
    fi
}

# Main execution
main() {
    print_header
    
    local passed=0
    local failed=0
    
    # System checks
    echo -e "${BLUE}System Checks:${NC}"
    check_docker && ((passed++)) || ((failed++))
    check_docker_compose && ((passed++)) || ((failed++))
    check_disk_space && ((passed++)) || ((failed++))
    check_memory && ((passed++)) || ((failed++))
    
    echo ""
    
    # Service checks
    echo -e "${BLUE}Service Checks:${NC}"
    check_container_status postgres && ((passed++)) || ((failed++))
    check_container_status vault && ((passed++)) || ((failed++))
    check_container_status traefik && ((passed++)) || ((failed++))
    
    echo ""
    
    # Connectivity checks
    echo -e "${BLUE}Connectivity Checks:${NC}"
    check_postgres && ((passed++)) || ((failed++))
    check_vault && ((passed++)) || ((failed++))
    check_traefik && ((passed++)) || ((failed++))
    
    # Summary
    show_summary "$passed" "$failed"
}

# Run main
main "$@"

