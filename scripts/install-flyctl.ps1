# Fly.io CLI 安裝腳本

Write-Host "🚀 安裝 Fly.io CLI (flyctl)..." -ForegroundColor Cyan

# 檢查是否已安裝
try {
    $version = flyctl version 2>$null
    if ($version) {
        Write-Host "✅ flyctl 已安裝: $version" -ForegroundColor Green
        exit 0
    }
} catch {
    # 繼續安裝
}

Write-Host "📦 開始安裝 flyctl..." -ForegroundColor Yellow

# 方法 1: 使用 winget (推薦)
try {
    Write-Host "嘗試使用 winget 安裝..." -ForegroundColor Cyan
    winget install Fly.Flyctl
    Write-Host "✅ flyctl 安裝完成！" -ForegroundColor Green
    Write-Host "請重新啟動 PowerShell 或命令提示字元" -ForegroundColor Yellow
    exit 0
} catch {
    Write-Host "⚠️  winget 安裝失敗，嘗試手動下載..." -ForegroundColor Yellow
}

# 方法 2: PowerShell 安裝腳本
try {
    Write-Host "使用 PowerShell 安裝腳本..." -ForegroundColor Cyan
    iwr https://fly.io/install.ps1 -useb | iex
    Write-Host "✅ flyctl 安裝完成！" -ForegroundColor Green
} catch {
    Write-Host "❌ 自動安裝失敗" -ForegroundColor Red
    Write-Host "請手動下載並安裝:" -ForegroundColor Yellow
    Write-Host "https://github.com/superfly/flyctl/releases" -ForegroundColor White
}

Write-Host ""
Write-Host "💡 安裝完成後，請執行以下命令登入："
Write-Host "flyctl auth login" -ForegroundColor Cyan

