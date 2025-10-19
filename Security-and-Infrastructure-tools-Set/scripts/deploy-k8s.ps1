# ============================================
# Kubernetes 服務部署腳本 (PowerShell 版本)
# Dennis Security And Infra Toolkit
# ============================================

param(
    [switch]$SkipConfirmation
)

# 顏色定義
$Red = "Red"
$Green = "Green"
$Yellow = "Yellow"
$Blue = "Blue"

# 日誌函數
function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor $Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor $Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor $Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor $Red
}

# 檢查 kubectl 是否可用
function Test-Kubectl {
    if (-not (Get-Command kubectl -ErrorAction SilentlyContinue)) {
        Write-Error "kubectl 未安裝或不在 PATH 中"
        return $false
    }
    Write-Success "kubectl 已安裝"
    return $true
}

# 檢查 Kubernetes 集群連接
function Test-KubernetesCluster {
    Write-Info "檢查 Kubernetes 集群連接..."
    try {
        $null = kubectl cluster-info 2>$null
        Write-Success "Kubernetes 集群連接正常"
        kubectl cluster-info
        return $true
    }
    catch {
        Write-Error "無法連接到 Kubernetes 集群"
        Write-Warning "請確保 Docker Desktop Kubernetes 已啟用"
        Write-Warning "或運行: kubectl config use-context docker-desktop"
        return $false
    }
}

# 部署命名空間
function Deploy-Namespaces {
    Write-Info "部署命名空間..."
    kubectl apply -f k8s/namespaces.yaml
    Write-Success "命名空間部署完成"
}

# 部署配置
function Deploy-Configs {
    Write-Info "部署 ConfigMap 和 Secret..."
    kubectl apply -f k8s/configmap.yaml
    kubectl apply -f k8s/secrets.yaml
    Write-Success "配置部署完成"
}

# 部署 ArgoCD
function Deploy-ArgoCD {
    Write-Info "部署 ArgoCD..."
    
    # 創建 ArgoCD 命名空間
    kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
    
    # 安裝 ArgoCD
    kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
    
    # 等待 ArgoCD 就緒
    Write-Info "等待 ArgoCD 就緒..."
    kubectl wait --for=condition=available --timeout=300s deployment/argocd-server -n argocd
    
    # 部署 NodePort Service
    kubectl apply -f k8s/argocd-service.yaml
    
    Write-Success "ArgoCD 部署完成"
    
    # 獲取初始密碼
    Write-Info "獲取 ArgoCD 初始密碼..."
    try {
        $ARGOCD_PASSWORD = kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | ForEach-Object { [System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String($_)) }
        Write-Success "ArgoCD 初始密碼: $ARGOCD_PASSWORD"
        Write-Info "ArgoCD UI: http://localhost:30081"
        Write-Info "用戶名: admin"
    }
    catch {
        Write-Warning "無法獲取 ArgoCD 初始密碼，請手動檢查"
    }
}

# 部署 SecureCodeBox Operator
function Deploy-Operator {
    Write-Info "部署 SecureCodeBox Operator..."
    
    # 安裝 CRDs
    kubectl apply -f https://raw.githubusercontent.com/secureCodeBox/secureCodeBox/main/operator/crds/cascading-rule.yaml
    kubectl apply -f https://raw.githubusercontent.com/secureCodeBox/secureCodeBox/main/operator/crds/scan-type.yaml
    kubectl apply -f https://raw.githubusercontent.com/secureCodeBox/secureCodeBox/main/operator/crds/scan.yaml
    
    # 部署 RBAC
    kubectl apply -f k8s/securecodebox-rbac.yaml
    
    # 部署 Operator
    kubectl apply -f k8s/securecodebox-operator.yaml
    
    # 等待 Operator 就緒
    Write-Info "等待 SecureCodeBox Operator 就緒..."
    kubectl wait --for=condition=available --timeout=300s deployment/securecodebox-operator -n security-tools
    
    Write-Success "SecureCodeBox Operator 部署完成"
}

# 部署 Parser 服務
function Deploy-Parsers {
    Write-Info "部署 Parser 服務..."
    
    # 部署 Nuclei Parser
    kubectl apply -f k8s/parser-nuclei.yaml
    
    # 部署 AMASS Parser
    kubectl apply -f k8s/parser-amass.yaml
    
    # 等待 Parser 服務就緒
    Write-Info "等待 Parser 服務就緒..."
    kubectl wait --for=condition=available --timeout=300s deployment/parser-nuclei -n security-tools
    kubectl wait --for=condition=available --timeout=300s deployment/parser-amass -n security-tools
    
    Write-Success "Parser 服務部署完成"
}

# 驗證部署
function Test-Deployment {
    Write-Info "驗證部署狀態..."
    
    Write-Host ""
    Write-Info "=== Pod 狀態 ==="
    kubectl get pods -n security-tools
    kubectl get pods -n argocd
    
    Write-Host ""
    Write-Info "=== Service 狀態 ==="
    kubectl get services -n security-tools
    kubectl get services -n argocd
    
    Write-Host ""
    Write-Info "=== 部署狀態 ==="
    kubectl get deployments -n security-tools
    kubectl get deployments -n argocd
    
    Write-Success "部署驗證完成"
}

# 顯示訪問信息
function Show-AccessInfo {
    Write-Host ""
    Write-Success "=== 部署完成! ==="
    Write-Host ""
    Write-Info "服務訪問信息:"
    Write-Host "  ArgoCD UI: http://localhost:30081"
    Write-Host "  ArgoCD CLI: kubectl port-forward -n argocd svc/argocd-server 8080:443"
    Write-Host ""
    Write-Info "Parser 服務測試:"
    Write-Host "  kubectl port-forward -n security-tools svc/parser-nuclei 8080:8080"
    Write-Host "  kubectl port-forward -n security-tools svc/parser-amass 8081:8080"
    Write-Host ""
    Write-Info "查看日誌:"
    Write-Host "  kubectl logs -n security-tools -l app=securecodebox-operator"
    Write-Host "  kubectl logs -n security-tools -l app=parser-nuclei"
    Write-Host "  kubectl logs -n security-tools -l app=parser-amass"
}

# 主函數
function Start-Deployment {
    Write-Info "開始部署 Kubernetes 服務..."
    
    if (-not $SkipConfirmation) {
        Write-Warning "此操作將在 Kubernetes 集群中部署多個服務"
        $confirmation = Read-Host "確定要繼續嗎? (y/N)"
        if ($confirmation -ne 'y' -and $confirmation -ne 'Y') {
            Write-Info "操作已取消"
            return
        }
    }
    
    if (-not (Test-Kubectl)) { return }
    if (-not (Test-KubernetesCluster)) { return }
    
    Deploy-Namespaces
    Deploy-Configs
    Deploy-ArgoCD
    Deploy-Operator
    Deploy-Parsers
    Test-Deployment
    Show-AccessInfo
    
    Write-Success "所有服務部署完成!"
}

# 執行主函數
Start-Deployment