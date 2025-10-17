# ============================================
# ArgoCD 快速訪問腳本
# Dennis Security And Infra Toolkit
# ============================================

Write-Host "ArgoCD 訪問指南" -ForegroundColor Cyan
Write-Host "===================" -ForegroundColor Cyan
Write-Host ""

# 檢查 port-forward 狀態
Write-Host "檢查 port-forward 狀態..." -ForegroundColor Yellow
$portCheck = netstat -ano | findstr :8080
if ($portCheck) {
    Write-Host "Port-forward 正在運行" -ForegroundColor Green
    Write-Host ""
    Write-Host "訪問 ArgoCD:" -ForegroundColor Green
    Write-Host "   URL: http://localhost:8080" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "登錄信息:" -ForegroundColor Green
    Write-Host "   用戶名: admin" -ForegroundColor Yellow
    Write-Host "   密碼: nyyTFntVy68o5v2u" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "提示: 如果無法訪問，請確保 port-forward 正在運行" -ForegroundColor Cyan
} else {
    Write-Host "Port-forward 未運行" -ForegroundColor Red
    Write-Host ""
    Write-Host "解決方案:" -ForegroundColor Yellow
    Write-Host "   運行: kubectl port-forward -n argocd svc/argocd-server 8080:80" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "其他服務密碼:" -ForegroundColor Cyan
Write-Host "   PostgreSQL: changeme" -ForegroundColor Yellow
Write-Host "   Vault: root" -ForegroundColor Yellow
Write-Host "   ArgoCD: nyyTFntVy68o5v2u" -ForegroundColor Yellow