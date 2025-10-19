# Simple PaaS Service Integration Test

$services = @(
    @{Name="Koyeb Agent Health"; Url="https://dizzy-sher-mitake-7f13854a.koyeb.app:8080/health"},
    @{Name="Koyeb Agent Metrics"; Url="https://dizzy-sher-mitake-7f13854a.koyeb.app:8080/metrics"},
    @{Name="Prometheus Health"; Url="https://pandora-monitoring.fly.dev:9090/-/healthy"},
    @{Name="Grafana Health"; Url="https://pandora-monitoring.fly.dev:3000/api/health"},
    @{Name="Loki Ready"; Url="https://pandora-monitoring.fly.dev:3100/ready"},
    @{Name="AlertManager Health"; Url="https://pandora-monitoring.fly.dev:9093/-/healthy"},
    @{Name="Nginx Proxy"; Url="https://nginx-stable-perl-boqt.onrender.com/health"}
)

Write-Host "`nPandora Box Console - PaaS Integration Test`n" -ForegroundColor Cyan

$passed = 0
$failed = 0

foreach ($service in $services) {
    Write-Host -NoNewline "Testing $($service.Name)... "
    try {
        $response = Invoke-WebRequest -Uri $service.Url -Method Get -TimeoutSec 30 -UseBasicParsing -ErrorAction Stop
        Write-Host "PASS (HTTP $($response.StatusCode))" -ForegroundColor Green
        $passed++
    } catch {
        Write-Host "FAIL" -ForegroundColor Red
        $failed++
    }
}

Write-Host "`nResults:" -ForegroundColor Cyan
Write-Host "Passed: $passed" -ForegroundColor Green
Write-Host "Failed: $failed" -ForegroundColor Red

if ($failed -eq 0) {
    Write-Host "`nAll services are running!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "`nSome services failed. Check the integration guide." -ForegroundColor Yellow
    exit 1
}
