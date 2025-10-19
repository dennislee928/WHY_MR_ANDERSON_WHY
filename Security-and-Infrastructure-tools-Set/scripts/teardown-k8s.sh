#!/bin/bash

# ============================================
# Kubernetes 服務清理腳本
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

# 確認清理操作
confirm_cleanup() {
    log_warning "此操作將刪除所有 Kubernetes 服務和資源！"
    read -p "確定要繼續嗎？(y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "操作已取消"
        exit 0
    fi
}

# 清理 Parser 服務
cleanup_parsers() {
    log_info "清理 Parser 服務..."
    kubectl delete -f k8s/parser-nuclei.yaml --ignore-not-found=true
    kubectl delete -f k8s/parser-amass.yaml --ignore-not-found=true
    log_success "Parser 服務清理完成"
}

# 清理 SecureCodeBox Operator
cleanup_operator() {
    log_info "清理 SecureCodeBox Operator..."
    kubectl delete -f k8s/securecodebox-operator.yaml --ignore-not-found=true
    kubectl delete -f k8s/securecodebox-rbac.yaml --ignore-not-found=true
    log_success "SecureCodeBox Operator 清理完成"
}

# 清理 ArgoCD
cleanup_argocd() {
    log_info "清理 ArgoCD..."
    kubectl delete -f k8s/argocd-service.yaml --ignore-not-found=true
    kubectl delete namespace argocd --ignore-not-found=true
    log_success "ArgoCD 清理完成"
}

# 清理配置
cleanup_configs() {
    log_info "清理配置..."
    kubectl delete -f k8s/configmap.yaml --ignore-not-found=true
    kubectl delete -f k8s/secrets.yaml --ignore-not-found=true
    log_success "配置清理完成"
}

# 清理命名空間
cleanup_namespaces() {
    log_info "清理命名空間..."
    kubectl delete -f k8s/namespaces.yaml --ignore-not-found=true
    log_success "命名空間清理完成"
}

# 清理 CRDs
cleanup_crds() {
    log_info "清理 SecureCodeBox CRDs..."
    kubectl delete -f https://raw.githubusercontent.com/secureCodeBox/secureCodeBox/main/operator/crds/cascading-rule.yaml --ignore-not-found=true
    kubectl delete -f https://raw.githubusercontent.com/secureCodeBox/secureCodeBox/main/operator/crds/scan-type.yaml --ignore-not-found=true
    kubectl delete -f https://raw.githubusercontent.com/secureCodeBox/secureCodeBox/main/operator/crds/scan.yaml --ignore-not-found=true
    log_success "CRDs 清理完成"
}

# 驗證清理
verify_cleanup() {
    log_info "驗證清理狀態..."
    
    echo ""
    log_info "=== 剩餘 Pod ==="
    kubectl get pods -n security-tools 2>/dev/null || log_info "security-tools 命名空間不存在"
    kubectl get pods -n argocd 2>/dev/null || log_info "argocd 命名空間不存在"
    
    echo ""
    log_info "=== 剩餘 Service ==="
    kubectl get services -n security-tools 2>/dev/null || log_info "security-tools 命名空間不存在"
    kubectl get services -n argocd 2>/dev/null || log_info "argocd 命名空間不存在"
    
    log_success "清理驗證完成"
}

# 主函數
main() {
    log_info "開始清理 Kubernetes 服務..."
    
    confirm_cleanup
    cleanup_parsers
    cleanup_operator
    cleanup_argocd
    cleanup_configs
    cleanup_namespaces
    cleanup_crds
    verify_cleanup
    
    log_success "所有服務清理完成！"
}

# 執行主函數
main "$@"

