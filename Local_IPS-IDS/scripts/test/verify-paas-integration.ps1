# PaaS Service Integration Verification Script (PowerShell)
# Checks if all PaaS platform services are running and properly integrated

# Service URLs
$KOYEB_AGENT_URL = "https://dizzy-sher-mitake-7f13854a.koyeb.app:8080"
$FLYIO_MONITORING_URL = "https://pandora-monitoring.fly.dev"
$RENDER_NGINX_URL = "https://nginx-stable-perl-boqt.onrender.com"

# Counters
$TotalTests = 0
$PassedTests = 0
$FailedTests = 0

# Test function
function Test-Service {
    param(
        [string]$Name,
        [string]$Url,
        [int]$ExpectedCode = 200
    )
    
    $script:TotalTests++
    Write-Host -NoNewline "Testing $Name... "
    
    try {
        $response = Invoke-WebRequest -Uri $Url -Method Get -TimeoutSec 30 -UseBasicParsing -ErrorAction Stop
        $statusCode = $response.StatusCode
        
        if ($statusCode -eq $ExpectedCode -or $statusCode -eq 200 -or $statusCode -eq 302) {
            Write-Host "PASS" -ForegroundColor Green -NoNewline
            Write-Host " (HTTP $statusCode)"
            $script:PassedTests++
            return $true
        } else {
            Write-Host "FAIL" -ForegroundColor Red -NoNewline
            Write-Host " (HTTP $statusCode)"
            $script:FailedTests++
            return $false
        }
    } catch {
        Write-Host "FAIL" -ForegroundColor Red -NoNewline
        Write-Host " (Connection failed)"
        $script:FailedTests++
        return $false
    }
}

# Banner
Write-Host ""
Write-Host "╔═══════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║     Pandora Box Console - PaaS Integration Test         ║" -ForegroundColor Cyan
Write-Host "╚═══════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

# 1. Koyeb Agent Tests
Write-Host "`n[1/4] Testing Koyeb Agent..." -ForegroundColor Yellow
Test-Service -Name "Koyeb Agent Health" -Url "$KOYEB_AGENT_URL/health"
Test-Service -Name "Koyeb Agent Metrics" -Url "$KOYEB_AGENT_URL/metrics"

# 2. Fly.io Monitoring Stack Tests
Write-Host "`n[2/4] Testing Fly.io Monitoring Stack..." -ForegroundColor Yellow
Test-Service -Name "Prometheus Health" -Url "$FLYIO_MONITORING_URL:9090/-/healthy"
Test-Service -Name "Prometheus Targets" -Url "$FLYIO_MONITORING_URL:9090/api/v1/targets"
Test-Service -Name "Grafana Health" -Url "$FLYIO_MONITORING_URL:3000/api/health"
Test-Service -Name "Loki Ready" -Url "$FLYIO_MONITORING_URL:3100/ready"
Test-Service -Name "AlertManager Health" -Url "$FLYIO_MONITORING_URL:9093/-/healthy"

# 3. Render Services Tests
Write-Host "`n[3/4] Testing Render Services..." -ForegroundColor Yellow
Test-Service -Name "Nginx Proxy" -Url "$RENDER_NGINX_URL/health"

# 4. Integration Tests
Write-Host "`n[4/4] Testing Service Integration..." -ForegroundColor Yellow

# Check if Prometheus is scraping Koyeb Agent
Write-Host -NoNewline "Testing Prometheus scraping Koyeb Agent... "
try {
    $targets = Invoke-RestMethod -Uri "$FLYIO_MONITORING_URL:9090/api/v1/targets" -Method Get -TimeoutSec 30
    if ($targets -match "dizzy-sher-mitake") {
        Write-Host "PASS" -ForegroundColor Green
        $PassedTests++
    } else {
        Write-Host "FAIL (Target not found)" -ForegroundColor Red
        $FailedTests++
    }
} catch {
    Write-Host "FAIL (Connection failed)" -ForegroundColor Red
    $FailedTests++
}
$TotalTests++

# Results Summary
Write-Host "`n═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "Test Summary" -ForegroundColor Cyan
Write-Host "═══════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "Total Tests:  $TotalTests"
Write-Host "Passed:       $PassedTests" -ForegroundColor Green
Write-Host "Failed:       $FailedTests" -ForegroundColor Red
Write-Host "Success Rate: $([math]::Round($PassedTests * 100 / $TotalTests, 2))%"

# Service Status Table
Write-Host "`nService Status:" -ForegroundColor Cyan
Write-Host "┌─────────────────────────┬────────────────────────────────────────────────┐"
Write-Host "│ Service                 │ URL                                            │"
Write-Host "├─────────────────────────┼────────────────────────────────────────────────┤"
Write-Host "│ Koyeb Agent             │ $KOYEB_AGENT_URL                               │"
Write-Host "│ Fly.io Grafana          │ $FLYIO_MONITORING_URL`:3000                    │"
Write-Host "│ Fly.io Prometheus       │ $FLYIO_MONITORING_URL`:9090                    │"
Write-Host "│ Fly.io Loki             │ $FLYIO_MONITORING_URL`:3100                    │"
Write-Host "│ Fly.io AlertManager     │ $FLYIO_MONITORING_URL`:9093                    │"
Write-Host "│ Render Nginx            │ $RENDER_NGINX_URL                              │"
Write-Host "└─────────────────────────┴────────────────────────────────────────────────┘"

# Recommendations
if ($FailedTests -gt 0) {
    Write-Host "`nRecommendations:" -ForegroundColor Yellow
    Write-Host "1. Check service logs for errors"
    Write-Host "2. Verify environment variables are set correctly"
    Write-Host "3. Ensure all services are deployed and running"
    Write-Host "4. Review the integration guide: docs/deployment/paas-integration-guide.md"
    exit 1
} else {
    Write-Host "`nAll services are running and integrated successfully!" -ForegroundColor Green
    exit 0
}
