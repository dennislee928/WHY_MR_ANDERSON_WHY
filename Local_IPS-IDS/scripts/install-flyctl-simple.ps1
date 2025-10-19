Write-Host "Installing Fly.io CLI..." -ForegroundColor Green

try {
    winget install Fly.Flyctl
    Write-Host "Flyctl installed successfully!" -ForegroundColor Green
    Write-Host "Please restart PowerShell and run: flyctl auth login" -ForegroundColor Yellow
} catch {
    Write-Host "Winget failed, trying PowerShell install..." -ForegroundColor Yellow
    try {
        iwr https://fly.io/install.ps1 -useb | iex
        Write-Host "Flyctl installed successfully!" -ForegroundColor Green
    } catch {
        Write-Host "Auto install failed. Please manually download from:" -ForegroundColor Red
        Write-Host "https://github.com/superfly/flyctl/releases" -ForegroundColor White
    }
}

