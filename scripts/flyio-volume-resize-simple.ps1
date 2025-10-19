# Fly.io Volume Resize Script - Simplified Version
# Reduces pandora-monitoring volume from 18GB to 3GB

Write-Host "Checking Fly.io volumes..." -ForegroundColor Cyan

# Check if flyctl is installed
try {
    $null = flyctl version
} catch {
    Write-Host "ERROR: flyctl not found in PATH" -ForegroundColor Red
    Write-Host "Install: winget install Fly.Flyctl" -ForegroundColor Yellow
    exit 1
}

# List current volumes
Write-Host "`nCurrent volumes:" -ForegroundColor Green
flyctl volumes list --app pandora-monitoring

Write-Host "`nWARNING: This will require downtime!" -ForegroundColor Yellow
Write-Host "`nOptions:" -ForegroundColor White
Write-Host "1. Redeploy with new 3GB volume (Recommended)" -ForegroundColor Green
Write-Host "2. Manual migration" -ForegroundColor Yellow
Write-Host "N. Cancel" -ForegroundColor Red

$choice = Read-Host "`nSelect option (1/2/N)"

if ($choice -eq "1") {
    Write-Host "`nRedeploying pandora-monitoring..." -ForegroundColor Green
    
    try {
        flyctl deploy --app pandora-monitoring --config deployments/paas/flyio/fly-monitoring.toml --dockerfile build/docker/monitoring.dockerfile
        
        Write-Host "`nDeployment complete!" -ForegroundColor Green
        Write-Host "`nNew volumes:" -ForegroundColor Cyan
        flyctl volumes list --app pandora-monitoring
        
        Write-Host "`nIMPORTANT: Delete old 18GB volumes to stop billing!" -ForegroundColor Red
        Write-Host "Command: flyctl volumes delete <OLD_VOLUME_ID>" -ForegroundColor Yellow
    } catch {
        Write-Host "`nDeployment failed: $_" -ForegroundColor Red
        exit 1
    }
} elseif ($choice -eq "2") {
    Write-Host "`nCreating new 3GB volume..." -ForegroundColor Yellow
    
    try {
        flyctl volumes create monitoring_data_new --app pandora-monitoring --region nrt --size 3
        
        Write-Host "`nVolume created!" -ForegroundColor Green
        Write-Host "`nNext steps:" -ForegroundColor White
        Write-Host "1. flyctl ssh console --app pandora-monitoring" -ForegroundColor Gray
        Write-Host "2. Copy important data to new volume" -ForegroundColor Gray
        Write-Host "3. Update app configuration" -ForegroundColor Gray
        Write-Host "4. Delete old volumes" -ForegroundColor Gray
    } catch {
        Write-Host "`nVolume creation failed: $_" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "`nOperation cancelled" -ForegroundColor Red
    exit 0
}

Write-Host "`nRemember: Delete old volumes to stop billing!" -ForegroundColor Yellow
