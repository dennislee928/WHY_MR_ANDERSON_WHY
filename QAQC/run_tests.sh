#!/bin/bash

# Security Platform API Tests - Robot Framework Execution Script
# This script runs comprehensive API tests for the Cloudflare Workers deployment

set -e

# Configuration
PROJECT_NAME="security-platform-api-tests"
BASE_URL="https://security-platform-worker.workers.dev"
API_VERSION="v1"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
RESULTS_DIR="results_${TIMESTAMP}"
LOGS_DIR="logs_${TIMESTAMP}"

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
    
    if ! command -v robot &> /dev/null; then
        print_error "Robot Framework not found. Installing..."
        pip install robotframework robotframework-requests
    else
        print_success "Robot Framework found"
    fi
    
    if ! command -v python3 &> /dev/null; then
        print_error "Python3 not found. Please install Python 3.7+"
        exit 1
    else
        print_success "Python3 found"
    fi
}

# Create directories
setup_directories() {
    print_header "Setting Up Directories"
    
    mkdir -p "${RESULTS_DIR}"
    mkdir -p "${LOGS_DIR}"
    mkdir -p "test_data"
    
    print_success "Directories created"
}

# Test API connectivity
test_connectivity() {
    print_header "Testing API Connectivity"
    
    echo "Testing endpoint: ${BASE_URL}/api/${API_VERSION}/health"
    
    if curl -s -f "${BASE_URL}/api/${API_VERSION}/health" > /dev/null; then
        print_success "API is accessible"
    else
        print_error "API is not accessible. Please check the URL and deployment status."
        exit 1
    fi
}

# Run smoke tests
run_smoke_tests() {
    print_header "Running Smoke Tests"
    
    robot \
        --outputdir "${RESULTS_DIR}" \
        --logdir "${LOGS_DIR}" \
        --include smoke \
        --variable BASE_URL:${BASE_URL} \
        --variable API_VERSION:${API_VERSION} \
        --name "Smoke Tests" \
        api_tests.robot
    
    if [ $? -eq 0 ]; then
        print_success "Smoke tests passed"
    else
        print_error "Smoke tests failed"
        return 1
    fi
}

# Run regression tests
run_regression_tests() {
    print_header "Running Regression Tests"
    
    robot \
        --outputdir "${RESULTS_DIR}" \
        --logdir "${LOGS_DIR}" \
        --include regression \
        --variable BASE_URL:${BASE_URL} \
        --variable API_VERSION:${API_VERSION} \
        --name "Regression Tests" \
        api_tests.robot
    
    if [ $? -eq 0 ]; then
        print_success "Regression tests passed"
    else
        print_error "Regression tests failed"
        return 1
    fi
}

# Run Cloudflare Workers specific tests
run_cloudflare_tests() {
    print_header "Running Cloudflare Workers Tests"
    
    robot \
        --outputdir "${RESULTS_DIR}" \
        --logdir "${LOGS_DIR}" \
        --variable BASE_URL:${BASE_URL} \
        --variable API_VERSION:${API_VERSION} \
        --name "Cloudflare Workers Tests" \
        cloudflare_workers_tests.robot
    
    if [ $? -eq 0 ]; then
        print_success "Cloudflare Workers tests passed"
    else
        print_error "Cloudflare Workers tests failed"
        return 1
    fi
}

# Run performance tests
run_performance_tests() {
    print_header "Running Performance Tests"
    
    robot \
        --outputdir "${RESULTS_DIR}" \
        --logdir "${LOGS_DIR}" \
        --include performance \
        --variable BASE_URL:${BASE_URL} \
        --variable API_VERSION:${API_VERSION} \
        --name "Performance Tests" \
        api_tests.robot
    
    if [ $? -eq 0 ]; then
        print_success "Performance tests passed"
    else
        print_error "Performance tests failed"
        return 1
    fi
}

# Run integration tests
run_integration_tests() {
    print_header "Running Integration Tests"
    
    robot \
        --outputdir "${RESULTS_DIR}" \
        --logdir "${LOGS_DIR}" \
        --include integration \
        --variable BASE_URL:${BASE_URL} \
        --variable API_VERSION:${API_VERSION} \
        --name "Integration Tests" \
        test_suite_config.robot
    
    if [ $? -eq 0 ]; then
        print_success "Integration tests passed"
    else
        print_error "Integration tests failed"
        return 1
    fi
}

# Run all tests
run_all_tests() {
    print_header "Running All Tests"
    
    robot \
        --outputdir "${RESULTS_DIR}" \
        --logdir "${LOGS_DIR}" \
        --variable BASE_URL:${BASE_URL} \
        --variable API_VERSION:${API_VERSION} \
        --name "Complete Test Suite" \
        *.robot
    
    if [ $? -eq 0 ]; then
        print_success "All tests passed"
    else
        print_error "Some tests failed"
        return 1
    fi
}

# Generate test report
generate_report() {
    print_header "Generating Test Report"
    
    if [ -f "${RESULTS_DIR}/output.xml" ]; then
        # Generate HTML report
        rebot \
            --outputdir "${RESULTS_DIR}" \
            --name "Security Platform API Test Report" \
            --report "${RESULTS_DIR}/report.html" \
            --log "${RESULTS_DIR}/log.html" \
            "${RESULTS_DIR}/output.xml"
        
        print_success "Test report generated: ${RESULTS_DIR}/report.html"
    else
        print_warning "No output.xml found. Cannot generate report."
    fi
}

# Cleanup function
cleanup() {
    print_header "Cleanup"
    
    # Keep results and logs for analysis
    print_success "Test results saved in: ${RESULTS_DIR}"
    print_success "Test logs saved in: ${LOGS_DIR}"
}

# Main execution
main() {
    print_header "Security Platform API Test Suite"
    echo "Project: ${PROJECT_NAME}"
    echo "Base URL: ${BASE_URL}"
    echo "API Version: ${API_VERSION}"
    echo "Timestamp: ${TIMESTAMP}"
    echo ""
    
    # Check command line arguments
    case "${1:-all}" in
        "smoke")
            check_dependencies
            setup_directories
            test_connectivity
            run_smoke_tests
            generate_report
            cleanup
            ;;
        "regression")
            check_dependencies
            setup_directories
            test_connectivity
            run_regression_tests
            generate_report
            cleanup
            ;;
        "cloudflare")
            check_dependencies
            setup_directories
            test_connectivity
            run_cloudflare_tests
            generate_report
            cleanup
            ;;
        "performance")
            check_dependencies
            setup_directories
            test_connectivity
            run_performance_tests
            generate_report
            cleanup
            ;;
        "integration")
            check_dependencies
            setup_directories
            test_connectivity
            run_integration_tests
            generate_report
            cleanup
            ;;
        "all"|"")
            check_dependencies
            setup_directories
            test_connectivity
            run_all_tests
            generate_report
            cleanup
            ;;
        "help"|"-h"|"--help")
            echo "Usage: $0 [test_type]"
            echo ""
            echo "Test types:"
            echo "  smoke       - Run smoke tests only"
            echo "  regression  - Run regression tests only"
            echo "  cloudflare  - Run Cloudflare Workers specific tests"
            echo "  performance - Run performance tests only"
            echo "  integration - Run integration tests only"
            echo "  all         - Run all tests (default)"
            echo "  help        - Show this help message"
            ;;
        *)
            print_error "Unknown test type: $1"
            echo "Use '$0 help' for usage information"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
