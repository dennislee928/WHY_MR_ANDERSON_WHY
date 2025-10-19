# ============================================
# ArgoCD 訪問解決方案
# Dennis Security And Infra Toolkit
# ============================================

Write-Host "ArgoCD 訪問解決方案" -ForegroundColor Cyan
Write-Host "===================" -ForegroundColor Cyan
Write-Host ""

Write-Host "方法 1: 使用 Port Forward (推薦)" -ForegroundColor Green
Write-Host "kubectl port-forward -n argocd svc/argocd-server 30081:80" -ForegroundColor Yellow
Write-Host "然後訪問: http://localhost:30081" -ForegroundColor Yellow
Write-Host ""

Write-Host "方法 2: 使用 NodePort" -ForegroundColor Green
Write-Host "訪問: http://localhost:30081" -ForegroundColor Yellow
Write-Host "注意: Docker Desktop Kubernetes 可能需要額外配置" -ForegroundColor Red
Write-Host ""

Write-Host "方法 3: 使用 LoadBalancer" -ForegroundColor Green
Write-Host "等待 EXTERNAL-IP 分配後訪問" -ForegroundColor Yellow
Write-Host ""

Write-Host "登錄信息:" -ForegroundColor Cyan
Write-Host "用戶名: admin" -ForegroundColor Yellow
Write-Host "密碼: nyyTFntVy68o5v2u" -ForegroundColor Yellow
Write-Host ""

Write-Host "如果仍然無法訪問，請嘗試:" -ForegroundColor Red
Write-Host "1. 檢查 Docker Desktop Kubernetes 是否啟用" -ForegroundColor Yellow
Write-Host "2. 重啟 Docker Desktop" -ForegroundColor Yellow
Write-Host "3. 使用 port-forward 方法" -ForegroundColor Yellow
