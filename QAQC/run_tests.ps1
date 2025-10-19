# Security Platform API Tests - PowerShell Execution Script
# This script runs comprehensive API tests for the Cloudflare Workers deployment

param(
    [Parameter(Position=0)]
    [ValidateSet("smoke", "regression", "cloudflare", "performance", "integration", "all", "help")]
    [string]$TestType = "all"
)

# Configuration
$ProjectName = "security-platform-api-tests"
$BaseUrl = "https://security-platform-worker.workers.dev"
$ApiVersion = "v1"
$Timestamp = Get-Date -Format "yyyyMMdd_HHmmss"
$ResultsDir = "results_$Timestamp"
$LogsDir = "logs_$Timestamp"

# Functions
function Write-Header {
    param([string]$Message)
    Write-Host "========================================" -ForegroundColor Blue
    Write-Host $Message -ForegroundColor Blue
    Write-Host "========================================" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "⚠ $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor Red
}

# Check dependencies
function Test-Dependencies {
    Write-Header "Checking Dependencies"
    
    # Check Python
    try {
        $pythonVersion = python --version 2>&1
        Write-Success "Python found: $pythonVersion"
    }
    catch {
        Write-Error "Python not found. Please install Python 3.7+"
        exit 1
    }
    
    # Check Robot Framework
    try {
        $robotVersion = robot --version 2>&1
        Write-Success "Robot Framework found: $robotVersion"
    }
    catch {
        Write-Warning "Robot Framework not found. Installing..."
        pip install robotframework robotframework-requests
        Write-Success "Robot Framework installed"
    }
    
    # Check curl (for connectivity test)
    try {
        curl --version | Out-Null
        Write-Success "curl found"
    }
    catch {
        Write-Warning "curl not found. Using PowerShell Invoke-WebRequest instead"
    }
}

# Create directories
function New-TestDirectories {
    Write-Header "Setting Up Directories"
    
    New-Item -ItemType Directory -Path $ResultsDir -Force | Out-Null
    New-Item -ItemType Directory -Path $LogsDir -Force | Out-Null
    New-Item -ItemType Directory -Path "test_data" -Force | Out-Null
    
    Write-Success "Directories created"
}

# Test API connectivity
function Test-ApiConnectivity {
    Write-Header "Testing API Connectivity"
    
    $healthUrl = "$BaseUrl/api/$ApiVersion/health"
    Write-Host "Testing endpoint: $healthUrl"
    
    try {
        $response = Invoke-WebRequest -Uri $healthUrl -Method GET -TimeoutSec 30
        if ($response.StatusCode -eq 200) {
            Write-Success "API is accessible"
        }
        else {
            Write-Error "API returned status code: $($response.StatusCode)"
            exit 1
        }
    }
    catch {
        Write-Error "API is not accessible. Please check the URL and deployment status."
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
}

# Run smoke tests
function Invoke-SmokeTests {
    Write-Header "Running Smoke Tests"
    
    try {
        robot --outputdir $ResultsDir --logdir $LogsDir --include smoke --variable "BASE_URL:$BaseUrl" --variable "API_VERSION:$ApiVersion" --name "Smoke Tests" api_tests.robot
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Smoke tests passed"
        }
        else {
            Write-Error "Smoke tests failed"
            return $false
        }
    }
    catch {
        Write-Error "Error running smoke tests: $($_.Exception.Message)"
        return $false
    }
    return $true
}

# Run regression tests
function Invoke-RegressionTests {
    Write-Header "Running Regression Tests"
    
    try {
        robot --outputdir $ResultsDir --logdir $LogsDir --include regression --variable "BASE_URL:$BaseUrl" --variable "API_VERSION:$ApiVersion" --name "Regression Tests" api_tests.robot
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Regression tests passed"
        }
        else {
            Write-Error "Regression tests failed"
            return $false
        }
    }
    catch {
        Write-Error "Error running regression tests: $($_.Exception.Message)"
        return $false
    }
    return $true
}

# Run Cloudflare Workers specific tests
function Invoke-CloudflareTests {
    Write-Header "Running Cloudflare Workers Tests"
    
    try {
        robot --outputdir $ResultsDir --logdir $LogsDir --variable "BASE_URL:$BaseUrl" --variable "API_VERSION:$ApiVersion" --name "Cloudflare Workers Tests" cloudflare_workers_tests.robot
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Cloudflare Workers tests passed"
        }
        else {
            Write-Error "Cloudflare Workers tests failed"
            return $false
        }
    }
    catch {
        Write-Error "Error running Cloudflare tests: $($_.Exception.Message)"
        return $false
    }
    return $true
}

# Run performance tests
function Invoke-PerformanceTests {
    Write-Header "Running Performance Tests"
    
    try {
        robot --outputdir $ResultsDir --logdir $LogsDir --include performance --variable "BASE_URL:$BaseUrl" --variable "API_VERSION:$ApiVersion" --name "Performance Tests" api_tests.robot
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Performance tests passed"
        }
        else {
            Write-Error "Performance tests failed"
            return $false
        }
    }
    catch {
        Write-Error "Error running performance tests: $($_.Exception.Message)"
        return $false
    }
    return $true
}

# Run integration tests
function Invoke-IntegrationTests {
    Write-Header "Running Integration Tests"
    
    try {
        robot --outputdir $ResultsDir --logdir $LogsDir --include integration --variable "BASE_URL:$BaseUrl" --variable "API_VERSION:$ApiVersion" --name "Integration Tests" test_suite_config.robot
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "Integration tests passed"
        }
        else {
            Write-Error "Integration tests failed"
            return $false
        }
    }
    catch {
        Write-Error "Error running integration tests: $($_.Exception.Message)"
        return $false
    }
    return $true
}

# Run all tests
function Invoke-AllTests {
    Write-Header "Running All Tests"
    
    try {
        robot --outputdir $ResultsDir --logdir $LogsDir --variable "BASE_URL:$BaseUrl" --variable "API_VERSION:$ApiVersion" --name "Complete Test Suite" *.robot
        
        if ($LASTEXITCODE -eq 0) {
            Write-Success "All tests passed"
        }
        else {
            Write-Error "Some tests failed"
            return $false
        }
    }
    catch {
        Write-Error "Error running all tests: $($_.Exception.Message)"
        return $false
    }
    return $true
}

# Generate test report
function New-TestReport {
    Write-Header "Generating Test Report"
    
    $outputXml = Join-Path $ResultsDir "output.xml"
    if (Test-Path $outputXml) {
        try {
            rebot --outputdir $ResultsDir --name "Security Platform API Test Report" --report "$ResultsDir/report.html" --log "$ResultsDir/log.html" $outputXml
            Write-Success "Test report generated: $ResultsDir/report.html"
        }
        catch {
            Write-Warning "Error generating report: $($_.Exception.Message)"
        }
    }
    else {
        Write-Warning "No output.xml found. Cannot generate report."
    }
}

# Cleanup function
function Invoke-Cleanup {
    Write-Header "Cleanup"
    
    Write-Success "Test results saved in: $ResultsDir"
    Write-Success "Test logs saved in: $LogsDir"
    
    # Open results directory
    if (Test-Path $ResultsDir) {
        Write-Host "Opening results directory..." -ForegroundColor Cyan
        Start-Process explorer.exe -ArgumentList $ResultsDir
    }
}

# Show help
function Show-Help {
    Write-Host "Security Platform API Test Suite" -ForegroundColor Blue
    Write-Host ""
    Write-Host "Usage: .\run_tests.ps1 [test_type]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Test types:" -ForegroundColor Cyan
    Write-Host "  smoke       - Run smoke tests only"
    Write-Host "  regression  - Run regression tests only"
    Write-Host "  cloudflare  - Run Cloudflare Workers specific tests"
    Write-Host "  performance - Run performance tests only"
    Write-Host "  integration - Run integration tests only"
    Write-Host "  all         - Run all tests (default)"
    Write-Host "  help        - Show this help message"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Cyan
    Write-Host "  .\run_tests.ps1 smoke"
    Write-Host "  .\run_tests.ps1 cloudflare"
    Write-Host "  .\run_tests.ps1 all"
}

# Main execution
function Main {
    Write-Header "Security Platform API Test Suite"
    Write-Host "Project: $ProjectName"
    Write-Host "Base URL: $BaseUrl"
    Write-Host "API Version: $ApiVersion"
    Write-Host "Timestamp: $Timestamp"
    Write-Host ""
    
    switch ($TestType) {
        "smoke" {
            Test-Dependencies
            New-TestDirectories
            Test-ApiConnectivity
            Invoke-SmokeTests
            New-TestReport
            Invoke-Cleanup
        }
        "regression" {
            Test-Dependencies
            New-TestDirectories
            Test-ApiConnectivity
            Invoke-RegressionTests
            New-TestReport
            Invoke-Cleanup
        }
        "cloudflare" {
            Test-Dependencies
            New-TestDirectories
            Test-ApiConnectivity
            Invoke-CloudflareTests
            New-TestReport
            Invoke-Cleanup
        }
        "performance" {
            Test-Dependencies
            New-TestDirectories
            Test-ApiConnectivity
            Invoke-PerformanceTests
            New-TestReport
            Invoke-Cleanup
        }
        "integration" {
            Test-Dependencies
            New-TestDirectories
            Test-ApiConnectivity
            Invoke-IntegrationTests
            New-TestReport
            Invoke-Cleanup
        }
        "all" {
            Test-Dependencies
            New-TestDirectories
            Test-ApiConnectivity
            Invoke-AllTests
            New-TestReport
            Invoke-Cleanup
        }
        "help" {
            Show-Help
        }
        default {
            Write-Error "Unknown test type: $TestType"
            Show-Help
            exit 1
        }
    }
}

# Run main function
Main
