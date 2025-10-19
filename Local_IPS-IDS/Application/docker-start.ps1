# Pandora Box Console - Docker 啟動腳本（Windows）

$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Pandora Box Console - Docker 啟動   " -ForegroundColor Cyan  
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 檢查 Docker 是否運行
Write-Host "檢查 Docker..." -ForegroundColor Yellow
try {
    docker ps | Out-Null
    Write-Host "✓ Docker 正在運行" -ForegroundColor Green
} catch {
    Write-Host "✗ Docker 未運行或未安裝" -ForegroundColor Red
    Write-Host "  請啟動 Docker Desktop" -ForegroundColor Yellow
    exit 1
}

# 檢查 docker-compose
Write-Host "檢查 docker-compose..." -ForegroundColor Yellow
try {
    docker-compose version | Out-Null
    Write-Host "✓ docker-compose 可用" -ForegroundColor Green
} catch {
    Write-Host "✗ docker-compose 未安裝" -ForegroundColor Red
    exit 1
}

Write-Host ""

# 檢查環境變數檔案
if (-not (Test-Path ".env")) {
    Write-Host "⚠️  未找到 .env 檔案" -ForegroundColor Yellow
    Write-Host "   從 .env.example 複製..." -ForegroundColor Gray
    Copy-Item ".env.example" ".env"
    Write-Host "✓ 已創建 .env 檔案" -ForegroundColor Green
    Write-Host "   請編輯 .env 設定您的環境" -ForegroundColor Cyan
    Write-Host ""
}

# 啟動服務
Write-Host "啟動所有服務..." -ForegroundColor Yellow
Write-Host ""

docker-compose up -d

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  ✓ 所有服務已啟動！" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "🌐 訪問以下 URL：" -ForegroundColor Cyan
    Write-Host "  主介面:      http://localhost:3001" -ForegroundColor White
    Write-Host "  Grafana:     http://localhost:3000" -ForegroundColor White
    Write-Host "  Prometheus:  http://localhost:9090" -ForegroundColor White
    Write-Host "  Loki:        http://localhost:3100" -ForegroundColor White
    Write-Host "  AlertManager: http://localhost:9093" -ForegroundColor White
    Write-Host ""
    Write-Host "🔐 Grafana 預設帳號:" -ForegroundColor Cyan
    Write-Host "  使用者名稱: admin" -ForegroundColor White
    Write-Host "  密碼:       pandora123" -ForegroundColor White
    Write-Host ""
    Write-Host "📊 查看服務狀態：" -ForegroundColor Yellow
    Write-Host "  docker-compose ps" -ForegroundColor White
    Write-Host ""
    Write-Host "📝 查看日誌：" -ForegroundColor Yellow
    Write-Host "  docker-compose logs -f" -ForegroundColor White
    Write-Host ""
    Write-Host "🛑 停止服務：" -ForegroundColor Yellow
    Write-Host "  docker-compose down" -ForegroundColor White
    Write-Host ""
} else {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "  ✗ 啟動失敗" -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "請檢查錯誤訊息並修正" -ForegroundColor Yellow
    exit 1
}

