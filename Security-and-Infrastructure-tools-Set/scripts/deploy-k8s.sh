#!/bin/bash

# ============================================
# Kubernetes 服務部署腳本
# Dennis Security And Infra Toolkit
# ============================================

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日誌函數
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 檢查 kubectl 是否可用
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl 未安裝或不在 PATH 中"
        exit 1
    fi
    log_success "kubectl 已安裝"
}

# 檢查 Kubernetes 集群連接
check_cluster() {
    log_info "檢查 Kubernetes 集群連接..."
    if kubectl cluster-info &> /dev/null; then
        log_success "Kubernetes 集群連接正常"
        kubectl cluster-info
    else
        log_error "無法連接到 Kubernetes 集群"
        log_warning "請確保 Docker Desktop Kubernetes 已啟用"
        log_warning "或運行: kubectl config use-context docker-desktop"
        exit 1
    fi
}

# 部署命名空間
deploy_namespaces() {
    log_info "部署命名空間..."
    kubectl apply -f k8s/namespaces.yaml
    log_success "命名空間部署完成"
}

# 部署配置
deploy_configs() {
    log_info "部署 ConfigMap 和 Secret..."
    kubectl apply -f k8s/configmap.yaml
    kubectl apply -f k8s/secrets.yaml
    log_success "配置部署完成"
}

# 部署 ArgoCD
deploy_argocd() {
    log_info "部署 ArgoCD..."
    
    # 創建 ArgoCD 命名空間
    kubectl create namespace argocd --dry-run=client -o yaml | kubectl apply -f -
    
    # 安裝 ArgoCD
    kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
    
    # 等待 ArgoCD 就緒
    log_info "等待 ArgoCD 就緒..."
    kubectl wait --for=condition=available --timeout=300s deployment/argocd-server -n argocd
    
    # 部署 NodePort Service
    kubectl apply -f k8s/argocd-service.yaml
    
    log_success "ArgoCD 部署完成"
    
    # 獲取初始密碼
    log_info "獲取 ArgoCD 初始密碼..."
    ARGOCD_PASSWORD=$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)
    log_success "ArgoCD 初始密碼: $ARGOCD_PASSWORD"
    log_info "ArgoCD UI: http://localhost:30081"
    log_info "用戶名: admin"
}

# 部署 SecureCodeBox Operator
deploy_operator() {
    log_info "部署 SecureCodeBox Operator..."
    
    # 安裝 CRDs
    kubectl apply -f https://raw.githubusercontent.com/secureCodeBox/secureCodeBox/main/operator/crds/cascading-rule.yaml
    kubectl apply -f https://raw.githubusercontent.com/secureCodeBox/secureCodeBox/main/operator/crds/scan-type.yaml
    kubectl apply -f https://raw.githubusercontent.com/secureCodeBox/secureCodeBox/main/operator/crds/scan.yaml
    
    # 部署 RBAC
    kubectl apply -f k8s/securecodebox-rbac.yaml
    
    # 部署 Operator
    kubectl apply -f k8s/securecodebox-operator.yaml
    
    # 等待 Operator 就緒
    log_info "等待 SecureCodeBox Operator 就緒..."
    kubectl wait --for=condition=available --timeout=300s deployment/securecodebox-operator -n security-tools
    
    log_success "SecureCodeBox Operator 部署完成"
}

# 部署 Parser 服務
deploy_parsers() {
    log_info "部署 Parser 服務..."
    
    # 部署 Nuclei Parser
    kubectl apply -f k8s/parser-nuclei.yaml
    
    # 部署 AMASS Parser
    kubectl apply -f k8s/parser-amass.yaml
    
    # 等待 Parser 服務就緒
    log_info "等待 Parser 服務就緒..."
    kubectl wait --for=condition=available --timeout=300s deployment/parser-nuclei -n security-tools
    kubectl wait --for=condition=available --timeout=300s deployment/parser-amass -n security-tools
    
    log_success "Parser 服務部署完成"
}

# 驗證部署
verify_deployment() {
    log_info "驗證部署狀態..."
    
    echo ""
    log_info "=== Pod 狀態 ==="
    kubectl get pods -n security-tools
    kubectl get pods -n argocd
    
    echo ""
    log_info "=== Service 狀態 ==="
    kubectl get services -n security-tools
    kubectl get services -n argocd
    
    echo ""
    log_info "=== 部署狀態 ==="
    kubectl get deployments -n security-tools
    kubectl get deployments -n argocd
    
    log_success "部署驗證完成"
}

# 顯示訪問信息
show_access_info() {
    echo ""
    log_success "=== 部署完成！==="
    echo ""
    log_info "服務訪問信息:"
    echo "  ArgoCD UI: http://localhost:30081"
    echo "  ArgoCD CLI: kubectl port-forward -n argocd svc/argocd-server 8080:443"
    echo ""
    log_info "Parser 服務測試:"
    echo "  kubectl port-forward -n security-tools svc/parser-nuclei 8080:8080"
    echo "  kubectl port-forward -n security-tools svc/parser-amass 8081:8080"
    echo ""
    log_info "查看日誌:"
    echo "  kubectl logs -n security-tools -l app=securecodebox-operator"
    echo "  kubectl logs -n security-tools -l app=parser-nuclei"
    echo "  kubectl logs -n security-tools -l app=parser-amass"
}

# 主函數
main() {
    log_info "開始部署 Kubernetes 服務..."
    
    check_kubectl
    check_cluster
    deploy_namespaces
    deploy_configs
    deploy_argocd
    deploy_operator
    deploy_parsers
    verify_deployment
    show_access_info
    
    log_success "所有服務部署完成！"
}

# 執行主函數
main "$@"

