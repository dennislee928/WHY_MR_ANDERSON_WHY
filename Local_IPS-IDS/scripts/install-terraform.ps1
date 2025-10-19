# Simple Terraform Installer for Windows
$ErrorActionPreference = "Stop"

Write-Host "`n=== Terraform Installer for Windows ===`n" -ForegroundColor Green

# Check if Terraform is already installed
$terraformCmd = Get-Command terraform -ErrorAction SilentlyContinue
if ($terraformCmd) {
    Write-Host "[OK] Terraform is already installed at: $($terraformCmd.Source)" -ForegroundColor Green
    & terraform version
    exit 0
}

Write-Host "[INFO] Downloading Terraform..." -ForegroundColor Cyan

# Configuration
$version = "1.6.6"
$url = "https://releases.hashicorp.com/terraform/$version/terraform_${version}_windows_amd64.zip"
$downloadPath = "$env:TEMP\terraform.zip"
$installPath = "$env:USERPROFILE\terraform"

# Download
Invoke-WebRequest -Uri $url -OutFile $downloadPath -UseBasicParsing
Write-Host "[OK] Download complete" -ForegroundColor Green

# Extract
Write-Host "[INFO] Extracting..." -ForegroundColor Cyan
if (Test-Path $installPath) {
    Remove-Item $installPath -Recurse -Force
}
New-Item -ItemType Directory -Path $installPath -Force | Out-Null
Expand-Archive -Path $downloadPath -DestinationPath $installPath -Force
Write-Host "[OK] Extraction complete" -ForegroundColor Green

# Add to PATH
Write-Host "[INFO] Adding to PATH..." -ForegroundColor Cyan
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if (-not $userPath.Contains($installPath)) {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installPath", "User")
    $env:Path = "$env:Path;$installPath"
    Write-Host "[OK] Added to PATH" -ForegroundColor Green
}

# Cleanup
Remove-Item $downloadPath -Force

Write-Host "`n[SUCCESS] Terraform installed successfully!`n" -ForegroundColor Green
Write-Host "Installation path: $installPath" -ForegroundColor Yellow
Write-Host "`nVerifying installation..." -ForegroundColor Cyan
& "$installPath\terraform.exe" version

Write-Host "`nNext steps:" -ForegroundColor Cyan
Write-Host "  1. Restart your terminal (or run this command):" -ForegroundColor White
Write-Host "     `$env:Path = [System.Environment]::GetEnvironmentVariable('Path','User')" -ForegroundColor Gray
Write-Host "  2. Verify: terraform version" -ForegroundColor White
Write-Host "  3. Go to terraform directory: cd terraform" -ForegroundColor White
Write-Host "  4. Initialize: terraform init`n" -ForegroundColor White

