# ============================================
# Security Stack PowerShell Commands
# Windows 替代 Makefile 的腳本
# ============================================

param(
    [Parameter(Position=0)]
    [string]$Command = "help",
    
    [Parameter(Position=1)]
    [string]$Target = ""
)

$COMPOSE_DIR = "..\Docker\compose"
$SCRIPT_DIR = $PSScriptRoot

function Show-Help {
    Write-Host "Security Stack Commands:" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage: .\make.ps1 <command> [options]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Commands:" -ForegroundColor Green
    Write-Host "  up              - Start all services"
    Write-Host "  down            - Stop all services"
    Write-Host "  restart         - Restart all services"
    Write-Host "  ps              - View service status"
    Write-Host "  health          - Check health status"
    Write-Host "  logs            - View logs (use -f for follow mode)"
    Write-Host "  scan-nuclei     - Run Nuclei scan (requires -Target)"
    Write-Host "  scan-nmap       - Run Nmap scan (requires -Target)"
    Write-Host "  scan-amass      - Run AMASS scan (requires -Target)"
    Write-Host "  scan-burp       - Run Burp Suite scan (requires -Target)"
    Write-Host "  scan-intelowl   - Run IntelOwl Nuclei scan (requires -Target)"
    Write-Host "  backup          - Backup databases"
    Write-Host "  clean           - Remove all volumes"
    Write-Host "  deploy-k8s      - Deploy Kubernetes services"
    Write-Host "  teardown-k8s    - Remove Kubernetes services"
    Write-Host "  k8s-status      - Check Kubernetes services status"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  .\make.ps1 up"
    Write-Host "  .\make.ps1 scan-nuclei -Target https://example.com"
    Write-Host "  .\make.ps1 scan-nmap -Target 192.168.1.1"
}

function Invoke-Up {
    Write-Host "Starting all services..." -ForegroundColor Green
    Push-Location $COMPOSE_DIR
    docker-compose up -d
    Write-Host "Waiting for services to be healthy..." -ForegroundColor Yellow
    Start-Sleep -Seconds 10
    docker-compose ps
    Pop-Location
}

function Invoke-Down {
    Write-Host "Stopping all services..." -ForegroundColor Yellow
    Push-Location $COMPOSE_DIR
    docker-compose down
    Pop-Location
}

function Invoke-Restart {
    Write-Host "Restarting all services..." -ForegroundColor Yellow
    Push-Location $COMPOSE_DIR
    docker-compose restart
    Pop-Location
}

function Invoke-Ps {
    Push-Location $COMPOSE_DIR
    docker-compose ps
    Pop-Location
}

function Invoke-Health {
    Write-Host "Checking service health..." -ForegroundColor Cyan
    Push-Location $COMPOSE_DIR
    docker-compose ps --format "table {{.Service}}`t{{.Status}}"
    Pop-Location
}

function Invoke-Logs {
    Push-Location $COMPOSE_DIR
    docker-compose logs -f
    Pop-Location
}

function Invoke-ScanNuclei {
    if (-not $Target) {
        Write-Host "Error: TARGET is required" -ForegroundColor Red
        Write-Host "Usage: .\make.ps1 scan-nuclei -Target https://example.com"
        exit 1
    }
    
    $timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
    Write-Host "Running Nuclei scan on $Target..." -ForegroundColor Green
    
    Push-Location $COMPOSE_DIR
    docker-compose run --rm scanner-nuclei nuclei -u $Target -o "/results/nuclei-$timestamp.json"
    Pop-Location
}

function Invoke-ScanNmap {
    if (-not $Target) {
        Write-Host "Error: TARGET is required" -ForegroundColor Red
        Write-Host "Usage: .\make.ps1 scan-nmap -Target 192.168.1.1"
        exit 1
    }
    
    $timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
    Write-Host "Running Nmap scan on $Target..." -ForegroundColor Green
    
    Push-Location $COMPOSE_DIR
    docker-compose run --rm nmap nmap $Target -oX "/results/nmap-$timestamp.xml"
    Pop-Location
}

function Invoke-Backup {
    $timestamp = Get-Date -Format "yyyyMMdd-HHmmss"
    $backupDir = "..\backups"
    
    if (-not (Test-Path $backupDir)) {
        New-Item -ItemType Directory -Path $backupDir | Out-Null
    }
    
    Write-Host "Creating backup..." -ForegroundColor Green
    Push-Location $COMPOSE_DIR
    docker-compose exec -T postgres pg_dump -U sectools security > "$backupDir\db-$timestamp.sql"
    Pop-Location
    
    Write-Host "Backup created in backups\" -ForegroundColor Green
}

function Invoke-Clean {
    Write-Host "WARNING: This will remove all volumes and data!" -ForegroundColor Red
    $confirmation = Read-Host "Are you sure? (y/N)"
    
    if ($confirmation -eq 'y' -or $confirmation -eq 'Y') {
        Push-Location $COMPOSE_DIR
        docker-compose down -v
        Write-Host "All volumes removed" -ForegroundColor Yellow
        Pop-Location
    } else {
        Write-Host "Cancelled" -ForegroundColor Green
    }
}

function Invoke-ScanAmass {
    if (-not $Target) {
        Write-Host "Error: -Target parameter required for AMASS scan" -ForegroundColor Red
        Write-Host "Example: .\make.ps1 scan-amass -Target example.com" -ForegroundColor Yellow
        return
    }
    
    Write-Host "Running AMASS scan on $Target..." -ForegroundColor Green
    Push-Location $COMPOSE_DIR
    docker-compose run --rm scanner-amass amass enum -d $Target -o /results/amass-$Target.txt
    Pop-Location
    Write-Host "AMASS scan completed. Results saved to scan_results volume." -ForegroundColor Green
}

function Invoke-ScanBurp {
    if (-not $Target) {
        Write-Host "Error: -Target parameter required for Burp Suite scan" -ForegroundColor Red
        Write-Host "Example: .\make.ps1 scan-burp -Target https://example.com" -ForegroundColor Yellow
        return
    }
    
    Write-Host "Running Burp Suite scan on $Target..." -ForegroundColor Green
    Write-Host "Note: Burp Suite requires GUI. This is a basic containerized version." -ForegroundColor Yellow
    Push-Location $COMPOSE_DIR
    docker-compose run --rm scanner-burpsuite java -jar /opt/burpsuite/burpsuite.jar --help
    Pop-Location
}

function Invoke-ScanIntelOwl {
    if (-not $Target) {
        Write-Host "Error: -Target parameter required for IntelOwl scan" -ForegroundColor Red
        Write-Host "Example: .\make.ps1 scan-intelowl -Target https://example.com" -ForegroundColor Yellow
        return
    }
    
    Write-Host "Running IntelOwl Nuclei scan on $Target..." -ForegroundColor Green
    Push-Location $COMPOSE_DIR
    docker-compose run --rm intelowl-nuclei nuclei -u $Target -o /results/intelowl-$Target.json
    Pop-Location
    Write-Host "IntelOwl scan completed. Results saved to scan_results volume." -ForegroundColor Green
}

function Invoke-DeployK8s {
    Write-Host "Deploying Kubernetes services..." -ForegroundColor Green
    
    # 檢查 kubectl
    if (-not (Get-Command kubectl -ErrorAction SilentlyContinue)) {
        Write-Host "Error: kubectl not found. Please install kubectl first." -ForegroundColor Red
        return
    }
    
    # 檢查 Kubernetes 集群
    try {
        kubectl cluster-info | Out-Null
        Write-Host "Kubernetes cluster connection verified" -ForegroundColor Green
    }
    catch {
        Write-Host "Error: Cannot connect to Kubernetes cluster" -ForegroundColor Red
        Write-Host "Please ensure Docker Desktop Kubernetes is enabled" -ForegroundColor Yellow
        return
    }
    
    # 執行部署腳本
    $scriptPath = Join-Path $PSScriptRoot "..\scripts\deploy-k8s.ps1"
    if (Test-Path $scriptPath) {
        Write-Host "Running PowerShell deployment script..." -ForegroundColor Yellow
        PowerShell -ExecutionPolicy Bypass -File $scriptPath
    }
    else {
        Write-Host "Error: Deployment script not found at $scriptPath" -ForegroundColor Red
        Write-Host "Available files in scripts directory:" -ForegroundColor Yellow
        Get-ChildItem (Join-Path $PSScriptRoot "..\scripts") -ErrorAction SilentlyContinue | ForEach-Object { Write-Host "  $($_.Name)" }
    }
}

function Invoke-TeardownK8s {
    Write-Host "Removing Kubernetes services..." -ForegroundColor Yellow
    
    # 檢查 kubectl
    if (-not (Get-Command kubectl -ErrorAction SilentlyContinue)) {
        Write-Host "Error: kubectl not found. Please install kubectl first." -ForegroundColor Red
        return
    }
    
    # 執行清理腳本
    $scriptPath = Join-Path $PSScriptRoot "..\scripts\teardown-k8s.sh"
    if (Test-Path $scriptPath) {
        Write-Host "Running teardown script..." -ForegroundColor Yellow
        bash $scriptPath
    }
    else {
        Write-Host "Error: Teardown script not found at $scriptPath" -ForegroundColor Red
    }
}

function Invoke-K8sStatus {
    Write-Host "Checking Kubernetes services status..." -ForegroundColor Green
    
    # 檢查 kubectl
    if (-not (Get-Command kubectl -ErrorAction SilentlyContinue)) {
        Write-Host "Error: kubectl not found. Please install kubectl first." -ForegroundColor Red
        return
    }
    
    try {
        Write-Host "=== Cluster Info ===" -ForegroundColor Cyan
        kubectl cluster-info
        
        Write-Host "`n=== Pods Status ===" -ForegroundColor Cyan
        kubectl get pods -n security-tools 2>$null
        kubectl get pods -n argocd 2>$null
        
        Write-Host "`n=== Services Status ===" -ForegroundColor Cyan
        kubectl get services -n security-tools 2>$null
        kubectl get services -n argocd 2>$null
        
        Write-Host "`n=== Deployments Status ===" -ForegroundColor Cyan
        kubectl get deployments -n security-tools 2>$null
        kubectl get deployments -n argocd 2>$null
    }
    catch {
        Write-Host "Error: Cannot connect to Kubernetes cluster" -ForegroundColor Red
        Write-Host "Please ensure Docker Desktop Kubernetes is enabled" -ForegroundColor Yellow
    }
}

# Main execution
switch ($Command.ToLower()) {
    "help" { Show-Help }
    "up" { Invoke-Up }
    "down" { Invoke-Down }
    "restart" { Invoke-Restart }
    "ps" { Invoke-Ps }
    "health" { Invoke-Health }
    "logs" { Invoke-Logs }
    "scan-nuclei" { Invoke-ScanNuclei }
    "scan-nmap" { Invoke-ScanNmap }
    "scan-amass" { Invoke-ScanAmass }
    "scan-burp" { Invoke-ScanBurp }
    "scan-intelowl" { Invoke-ScanIntelOwl }
    "backup" { Invoke-Backup }
    "clean" { Invoke-Clean }
    "deploy-k8s" { Invoke-DeployK8s }
    "teardown-k8s" { Invoke-TeardownK8s }
    "k8s-status" { Invoke-K8sStatus }
    default {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help
        exit 1
    }
}

