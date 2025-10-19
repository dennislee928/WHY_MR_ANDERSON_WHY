# Cloudflare Containers Build and Deploy Script - PowerShell
# This script builds and deploys all containers for the security platform

param(
    [Parameter(Position=0)]
    [string]$Version = "latest",
    [Parameter(Position=1)]
    [string]$Environment = "production",
    [Parameter(Position=2)]
    [ValidateSet("build", "push", "deploy", "test", "all", "help")]
    [string]$Action = "all"
)

# Configuration
$ProjectName = "security-platform-containers"
$Registry = "ghcr.io/dennislee928"

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
    
    # Check Docker
    try {
        $dockerVersion = docker --version 2>&1
        Write-Success "Docker found: $dockerVersion"
    }
    catch {
        Write-Error "Docker not found. Please install Docker Desktop."
        exit 1
    }
    
    # Check Docker Compose
    try {
        $composeVersion = docker-compose --version 2>&1
        Write-Success "Docker Compose found: $composeVersion"
    }
    catch {
        Write-Error "Docker Compose not found. Please install Docker Compose."
        exit 1
    }
    
    # Check Wrangler
    try {
        $wranglerVersion = wrangler --version 2>&1
        Write-Success "Wrangler found: $wranglerVersion"
    }
    catch {
        Write-Warning "Wrangler not found. Installing..."
        npm install -g wrangler
        Write-Success "Wrangler installed"
    }
}

# Build containers
function Build-Containers {
    Write-Header "Building Containers"
    
    # Build Backend API
    Write-Success "Building Backend API container..."
    docker build -t "${Registry}/${ProjectName}-backend-api:${Version}" ./containers/backend-api/
    
    # Build AI/Quantum
    Write-Success "Building AI/Quantum container..."
    docker build -t "${Registry}/${ProjectName}-ai-quantum:${Version}" ./containers/ai-quantum/
    
    # Build Security Tools
    Write-Success "Building Security Tools container..."
    docker build -t "${Registry}/${ProjectName}-security-tools:${Version}" ./containers/security-tools/
    
    # Build Database
    Write-Success "Building Database container..."
    docker build -t "${Registry}/${ProjectName}-database:${Version}" ./containers/database/
    
    # Build Monitoring
    Write-Success "Building Monitoring container..."
    docker build -t "${Registry}/${ProjectName}-monitoring:${Version}" ./containers/monitoring/
    
    Write-Success "All containers built successfully"
}

# Push containers to registry
function Push-Containers {
    Write-Header "Pushing Containers to Registry"
    
    # Login to registry (if needed)
    if ($Registry -ne "localhost") {
        Write-Success "Logging in to registry..."
        $env:GITHUB_TOKEN | docker login ghcr.io -u $env:GITHUB_USERNAME --password-stdin
    }
    
    # Push all containers
    docker push "${Registry}/${ProjectName}-backend-api:${Version}"
    docker push "${Registry}/${ProjectName}-ai-quantum:${Version}"
    docker push "${Registry}/${ProjectName}-security-tools:${Version}"
    docker push "${Registry}/${ProjectName}-database:${Version}"
    docker push "${Registry}/${ProjectName}-monitoring:${Version}"
    
    Write-Success "All containers pushed successfully"
}

# Deploy to Cloudflare Workers
function Deploy-Workers {
    Write-Header "Deploying to Cloudflare Workers"
    
    # Update wrangler.toml with container images
    $wranglerConfig = Get-Content "wrangler-containers.toml" -Raw
    $wranglerConfig = $wranglerConfig -replace "image = `"security-platform/backend-api`"", "image = `"${Registry}/${ProjectName}-backend-api:${Version}`""
    $wranglerConfig = $wranglerConfig -replace "image = `"security-platform/ai-quantum`"", "image = `"${Registry}/${ProjectName}-ai-quantum:${Version}`""
    $wranglerConfig = $wranglerConfig -replace "image = `"security-platform/security-tools`"", "image = `"${Registry}/${ProjectName}-security-tools:${Version}`""
    $wranglerConfig = $wranglerConfig -replace "image = `"security-platform/database`"", "image = `"${Registry}/${ProjectName}-database:${Version}`""
    $wranglerConfig = $wranglerConfig -replace "image = `"security-platform/monitoring`"", "image = `"${Registry}/${ProjectName}-monitoring:${Version}`""
    Set-Content "wrangler-containers.toml" $wranglerConfig
    
    # Deploy to Cloudflare Workers
    wrangler deploy --config wrangler-containers.toml --env $Environment
    
    Write-Success "Deployed to Cloudflare Workers successfully"
}

# Test deployment
function Test-Deployment {
    Write-Header "Testing Deployment"
    
    # Get deployment URL
    $deploymentOutput = wrangler deployments list --config wrangler-containers.toml --env $Environment
    $deploymentUrl = ($deploymentOutput | Select-String "https://" | Select-Object -First 1).Line.Trim()
    
    if (-not $deploymentUrl) {
        Write-Error "Could not get deployment URL"
        return $false
    }
    
    Write-Success "Testing deployment at: $deploymentUrl"
    
    # Test health endpoint
    try {
        $healthResponse = Invoke-WebRequest -Uri "${deploymentUrl}/api/v1/containers/health" -Method GET -TimeoutSec 30
        if ($healthResponse.StatusCode -eq 200) {
            Write-Success "Health check passed"
        }
        else {
            Write-Error "Health check failed with status: $($healthResponse.StatusCode)"
            return $false
        }
    }
    catch {
        Write-Error "Health check failed: $($_.Exception.Message)"
        return $false
    }
    
    # Test services endpoint
    try {
        $servicesResponse = Invoke-WebRequest -Uri "${deploymentUrl}/api/v1/services" -Method GET -TimeoutSec 30
        if ($servicesResponse.StatusCode -eq 200) {
            Write-Success "Services endpoint working"
        }
        else {
            Write-Error "Services endpoint failed with status: $($servicesResponse.StatusCode)"
            return $false
        }
    }
    catch {
        Write-Error "Services endpoint failed: $($_.Exception.Message)"
        return $false
    }
    
    Write-Success "All tests passed"
    return $true
}

# Cleanup
function Invoke-Cleanup {
    Write-Header "Cleanup"
    
    # Remove local images to save space
    docker rmi "${Registry}/${ProjectName}-backend-api:${Version}" 2>$null
    docker rmi "${Registry}/${ProjectName}-ai-quantum:${Version}" 2>$null
    docker rmi "${Registry}/${ProjectName}-security-tools:${Version}" 2>$null
    docker rmi "${Registry}/${ProjectName}-database:${Version}" 2>$null
    docker rmi "${Registry}/${ProjectName}-monitoring:${Version}" 2>$null
    
    Write-Success "Cleanup completed"
}

# Show help
function Show-Help {
    Write-Host "Cloudflare Containers Deployment Script" -ForegroundColor Blue
    Write-Host ""
    Write-Host "Usage: .\deploy-containers.ps1 [version] [environment] [action]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Parameters:" -ForegroundColor Cyan
    Write-Host "  version     - Container version (default: latest)"
    Write-Host "  environment - Deployment environment (default: production)"
    Write-Host "  action      - Action to perform (default: all)"
    Write-Host ""
    Write-Host "Actions:" -ForegroundColor Cyan
    Write-Host "  build       - Build containers only"
    Write-Host "  push        - Push containers to registry"
    Write-Host "  deploy      - Deploy to Cloudflare Workers"
    Write-Host "  test        - Test deployment"
    Write-Host "  all         - Build, push, deploy, and test (default)"
    Write-Host "  help        - Show this help message"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Cyan
    Write-Host "  .\deploy-containers.ps1 v1.0.0 production build"
    Write-Host "  .\deploy-containers.ps1 latest staging deploy"
    Write-Host "  .\deploy-containers.ps1 v2.0.0 production all"
}

# Main execution
function Main {
    Write-Header "Cloudflare Containers Deployment"
    Write-Host "Project: $ProjectName"
    Write-Host "Registry: $Registry"
    Write-Host "Version: $Version"
    Write-Host "Environment: $Environment"
    Write-Host ""
    
    switch ($Action) {
        "build" {
            Test-Dependencies
            Build-Containers
        }
        "push" {
            Test-Dependencies
            Push-Containers
        }
        "deploy" {
            Test-Dependencies
            Deploy-Workers
        }
        "test" {
            Test-Deployment
        }
        "all" {
            Test-Dependencies
            Build-Containers
            Push-Containers
            Deploy-Workers
            Test-Deployment
            Invoke-Cleanup
        }
        "help" {
            Show-Help
        }
        default {
            Write-Error "Unknown action: $Action"
            Show-Help
            exit 1
        }
    }
}

# Run main function
Main
