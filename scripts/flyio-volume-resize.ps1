# Fly.io Volume 調整腳本 (PowerShell 版本)
# 用於降低 pandora-monitoring 的 volume 大小以減少費用

Write-Host "🔍 檢查當前 Fly.io volumes..." -ForegroundColor Cyan

# 檢查 flyctl 是否安裝
try {
    flyctl version | Out-Null
} catch {
    Write-Host "❌ 錯誤: flyctl 未安裝或不在 PATH 中" -ForegroundColor Red
    Write-Host "請先安裝 Fly.io CLI:" -ForegroundColor Yellow
    Write-Host "Windows: winget install Fly.Flyctl" -ForegroundColor White
    Write-Host "或下載: https://fly.io/docs/getting-started/installing-flyctl/" -ForegroundColor White
    exit 1
}

# 檢查當前 volumes
Write-Host "📋 當前 volumes 列表：" -ForegroundColor Green
flyctl volumes list --app pandora-monitoring

Write-Host ""
Write-Host "⚠️  WARNING: 調整 volume 大小需要停機時間！" -ForegroundColor Yellow
Write-Host "📌 建議的兩種方式：" -ForegroundColor White
Write-Host ""
Write-Host "方式 1 (推薦 - 簡單)：重新部署應用" -ForegroundColor Green
Write-Host "1. 使用新的 3GB volume 配置重新部署" -ForegroundColor White
Write-Host "2. Fly.io 會自動創建新的 3GB volume" -ForegroundColor White
Write-Host "3. 手動刪除舊的 18GB volume 以停止計費" -ForegroundColor White
Write-Host ""
Write-Host "方式 2 (複雜)：手動遷移數據" -ForegroundColor Yellow
Write-Host "1. 創建新的 3GB volume" -ForegroundColor White
Write-Host "2. 遷移數據" -ForegroundColor White
Write-Host "3. 更新配置並刪除舊 volume" -ForegroundColor White

Write-Host ""
$choice = Read-Host "選擇執行方式 (1=重新部署/2=手動遷移/N=取消)"

switch ($choice) {
    "1" {
        Write-Host "🚀 執行重新部署方式..." -ForegroundColor Green
        Write-Host "正在重新部署 pandora-monitoring..." -ForegroundColor Cyan
        
        try {
            flyctl deploy --app pandora-monitoring --config deployments/paas/flyio/fly-monitoring.toml --dockerfile build/docker/monitoring.dockerfile
            Write-Host "✅ 重新部署完成！" -ForegroundColor Green
            Write-Host ""
            Write-Host "📋 檢查新的 volumes：" -ForegroundColor Cyan
            flyctl volumes list --app pandora-monitoring
            Write-Host ""
            Write-Host "⚠️  重要：請手動刪除舊的 18GB volume 以停止計費！" -ForegroundColor Red
            Write-Host "使用命令: flyctl volumes delete <OLD_VOLUME_ID>" -ForegroundColor Yellow
        } catch {
            Write-Host "❌ 部署失敗: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    "2" {
        Write-Host "🔧 執行手動遷移方式..." -ForegroundColor Yellow
        Write-Host "📦 創建新的 3GB volume..." -ForegroundColor Cyan
        
        try {
            flyctl volumes create monitoring_data_new --app pandora-monitoring --region nrt --size 3
            Write-Host "✅ 新 volume 創建完成！" -ForegroundColor Green
            Write-Host ""
            Write-Host "📝 下一步手動操作：" -ForegroundColor White
            Write-Host "1. flyctl ssh console --app pandora-monitoring" -ForegroundColor Gray
            Write-Host "2. 複製重要資料到新 volume" -ForegroundColor Gray
            Write-Host "3. 更新應用配置" -ForegroundColor Gray
            Write-Host "4. 刪除舊 volume" -ForegroundColor Gray
        } catch {
            Write-Host "❌ 創建 volume 失敗: $($_.Exception.Message)" -ForegroundColor Red
        }
    }
    default {
        Write-Host "❌ 取消操作" -ForegroundColor Red
        exit 0
    }
}

Write-Host ""
Write-Host "💰 記住：刪除舊的大 volume 才能停止計費！" -ForegroundColor Yellow