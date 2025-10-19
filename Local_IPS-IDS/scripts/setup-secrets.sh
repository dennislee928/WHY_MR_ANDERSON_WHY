#!/bin/bash

# Pandora Box Console IDS-IPS 密鑰設定腳本
# 此腳本協助設定 GitHub Actions 和 Kubernetes 密鑰

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔐 Pandora Box Console IDS-IPS 密鑰設定工具${NC}"
echo ""

# 檢查必要工具
check_prerequisites() {
    echo -e "${YELLOW}📋 檢查必要工具...${NC}"
    
    local missing_tools=()
    
    if ! command -v gh &> /dev/null; then
        missing_tools+=("gh (GitHub CLI)")
    fi
    
    if ! command -v kubectl &> /dev/null; then
        missing_tools+=("kubectl")
    fi
    
    if ! command -v base64 &> /dev/null; then
        missing_tools+=("base64")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        echo -e "${RED}❌ 缺少必要工具: ${missing_tools[*]}${NC}"
        echo "請安裝缺少的工具後重試"
        exit 1
    fi
    
    echo -e "${GREEN}✅ 所有必要工具已安裝${NC}"
}

# 設定 GitHub Actions 密鑰
setup_github_secrets() {
    echo -e "${YELLOW}🔧 設定 GitHub Actions 密鑰...${NC}"
    
    # 檢查是否已登入 GitHub CLI
    if ! gh auth status &> /dev/null; then
        echo -e "${YELLOW}⚠️  請先登入 GitHub CLI:${NC}"
        echo "  gh auth login"
        exit 1
    fi
    
    echo -e "${BLUE}請輸入以下 OCI 配置資訊:${NC}"
    
    # OCI 基本配置
    read -p "OCI User OCID: " OCI_USER
    read -p "OCI Tenancy OCID: " OCI_TENANCY
    read -p "OCI Region (預設: us-ashburn-1): " OCI_REGION
    OCI_REGION=${OCI_REGION:-us-ashburn-1}
    read -p "OCI API Key Fingerprint: " OCI_FINGERPRINT
    read -p "OCI Namespace: " OCI_NAMESPACE
    read -p "Kubernetes Cluster OCID: " CLUSTER_OCID
    
    # OCI Container Registry 認證
    read -p "OCI Registry Username (格式: namespace/username): " OCI_USERNAME
    read -s -p "OCI Registry Password/Auth Token: " OCI_PASSWORD
    echo ""
    
    # GCP 配置
    read -p "GCP Project ID: " GCP_PROJECT_ID
    read -p "GCP Cluster Name: " GCP_CLUSTER_NAME
    read -s -p "GCP Service Account Key (JSON): " GCP_SA_KEY
    echo ""
    
    # Vercel 配置
    read -p "Vercel Token: " VERCEL_TOKEN
    read -p "Vercel Organization ID: " VERCEL_ORG_ID
    read -p "Vercel Project ID (OCI): " VERCEL_PROJECT_ID_OCI
    read -p "Vercel Project ID (GCP): " VERCEL_PROJECT_ID_GCP
    read -p "API Base URL (OCI): " VERCEL_API_BASE_URL_OCI
    read -p "Grafana URL (OCI): " VERCEL_GRAFANA_URL_OCI
    read -p "Prometheus URL (OCI): " VERCEL_PROMETHEUS_URL_OCI
    read -p "API Base URL (GCP): " VERCEL_API_BASE_URL_GCP
    read -p "Grafana URL (GCP): " VERCEL_GRAFANA_URL_GCP
    read -p "Prometheus URL (GCP): " VERCEL_PROMETHEUS_URL_GCP
    
    echo -e "${YELLOW}正在設定 GitHub Actions 密鑰...${NC}"
    
    # 設定密鑰
    gh secret set OCI_USER --body "$OCI_USER"
    gh secret set OCI_TENANCY --body "$OCI_TENANCY"
    gh secret set OCI_REGION --body "$OCI_REGION"
    gh secret set OCI_FINGERPRINT --body "$OCI_FINGERPRINT"
    gh secret set OCI_NAMESPACE --body "$OCI_NAMESPACE"
    gh secret set CLUSTER_OCID --body "$CLUSTER_OCID"
    gh secret set OCI_USERNAME --body "$OCI_USERNAME"
    gh secret set OCI_PASSWORD --body "$OCI_PASSWORD"
    gh secret set GCP_PROJECT_ID --body "$GCP_PROJECT_ID"
    gh secret set GCP_CLUSTER_NAME --body "$GCP_CLUSTER_NAME"
    gh secret set GCP_SA_KEY --body "$GCP_SA_KEY"
    gh secret set VERCEL_TOKEN --body "$VERCEL_TOKEN"
    gh secret set VERCEL_ORG_ID --body "$VERCEL_ORG_ID"
    gh secret set VERCEL_PROJECT_ID_OCI --body "$VERCEL_PROJECT_ID_OCI"
    gh secret set VERCEL_PROJECT_ID_GCP --body "$VERCEL_PROJECT_ID_GCP"
    gh secret set VERCEL_API_BASE_URL_OCI --body "$VERCEL_API_BASE_URL_OCI"
    gh secret set VERCEL_GRAFANA_URL_OCI --body "$VERCEL_GRAFANA_URL_OCI"
    gh secret set VERCEL_PROMETHEUS_URL_OCI --body "$VERCEL_PROMETHEUS_URL_OCI"
    gh secret set VERCEL_API_BASE_URL_GCP --body "$VERCEL_API_BASE_URL_GCP"
    gh secret set VERCEL_GRAFANA_URL_GCP --body "$VERCEL_GRAFANA_URL_GCP"
    gh secret set VERCEL_PROMETHEUS_URL_GCP --body "$VERCEL_PROMETHEUS_URL_GCP"
    
    echo -e "${GREEN}✅ GitHub Actions 密鑰設定完成${NC}"
}

# 設定 Kubernetes 密鑰
setup_kubernetes_secrets() {
    echo -e "${YELLOW}⚙️  設定 Kubernetes 密鑰...${NC}"
    
    # 檢查 kubectl 連接
    if ! kubectl cluster-info &> /dev/null; then
        echo -e "${RED}❌ 無法連接到 Kubernetes 集群${NC}"
        echo "請確保 kubectl 已正確配置並能連接到目標集群"
        exit 1
    fi
    
    # 創建命名空間（如果不存在）
    kubectl create namespace pandora-box --dry-run=client -o yaml | kubectl apply -f -
    
    echo -e "${BLUE}請輸入資料庫密碼:${NC}"
    read -s -p "PostgreSQL 密碼 (預設: pandora123): " POSTGRES_PASSWORD
    POSTGRES_PASSWORD=${POSTGRES_PASSWORD:-pandora123}
    echo ""
    
    read -s -p "Redis 密碼 (預設: pandora123): " REDIS_PASSWORD
    REDIS_PASSWORD=${REDIS_PASSWORD:-pandora123}
    echo ""
    
    read -s -p "Grafana 管理員密碼 (預設: pandora123): " GRAFANA_PASSWORD
    GRAFANA_PASSWORD=${GRAFANA_PASSWORD:-pandora123}
    echo ""
    
    read -s -p "JWT 密鑰: " JWT_SECRET
    echo ""
    
    echo -e "${YELLOW}正在創建 Kubernetes 密鑰...${NC}"
    
    # 創建應用程式密鑰
    kubectl create secret generic pandora-secrets \
        --namespace=pandora-box \
        --from-literal=postgres-password="$(echo -n "$POSTGRES_PASSWORD" | base64)" \
        --from-literal=postgres-user="$(echo -n "pandora" | base64)" \
        --from-literal=redis-password="$(echo -n "$REDIS_PASSWORD" | base64)" \
        --from-literal=grafana-password="$(echo -n "$GRAFANA_PASSWORD" | base64)" \
        --from-literal=jwt-secret="$(echo -n "$JWT_SECRET" | base64)" \
        --dry-run=client -o yaml | kubectl apply -f -
    
    echo -e "${GREEN}✅ Kubernetes 密鑰創建完成${NC}"
}

# 生成 mTLS 憑證
generate_mtls_certs() {
    echo -e "${YELLOW}🔒 生成 mTLS 憑證...${NC}"
    
    if [ ! -d "certs" ]; then
        mkdir -p certs
    fi
    
    # 生成 CA 憑證
    openssl genrsa -out certs/ca.key 4096
    openssl req -new -x509 -days 365 -key certs/ca.key -out certs/ca.crt -subj "/C=TW/ST=Taiwan/L=Taipei/O=PandoraBox/CN=PandoraBox-CA"
    
    # 生成伺服器憑證
    openssl genrsa -out certs/server.key 4096
    openssl req -new -key certs/server.key -out certs/server.csr -subj "/C=TW/ST=Taiwan/L=Taipei/O=PandoraBox/CN=pandora-server"
    openssl x509 -req -days 365 -in certs/server.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/server.crt
    
    # 生成客戶端憑證
    openssl genrsa -out certs/client.key 4096
    openssl req -new -key certs/client.key -out certs/client.csr -subj "/C=TW/ST=Taiwan/L=Taipei/O=PandoraBox/CN=pandora-client"
    openssl x509 -req -days 365 -in certs/client.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/client.crt
    
    # 更新 Kubernetes 密鑰
    kubectl create secret tls pandora-mtls-certs \
        --namespace=pandora-box \
        --cert=certs/server.crt \
        --key=certs/server.key \
        --dry-run=client -o yaml | kubectl apply -f -
    
    kubectl create secret tls pandora-client-certs \
        --namespace=pandora-box \
        --cert=certs/client.crt \
        --key=certs/client.key \
        --dry-run=client -o yaml | kubectl apply -f -
    
    echo -e "${GREEN}✅ mTLS 憑證生成完成${NC}"
}

# 驗證設定
verify_setup() {
    echo -e "${YELLOW}🔍 驗證設定...${NC}"
    
    # 驗證 GitHub 密鑰
    echo -e "${BLUE}GitHub Actions 密鑰:${NC}"
    gh secret list | grep -E "(OCI_|VERCEL_)" || echo "未找到相關密鑰"
    
    # 驗證 Kubernetes 密鑰
    echo -e "${BLUE}Kubernetes 密鑰:${NC}"
    kubectl get secrets -n pandora-box | grep pandora || echo "未找到相關密鑰"
    
    # 驗證憑證檔案
    echo -e "${BLUE}mTLS 憑證檔案:${NC}"
    if [ -d "certs" ]; then
        ls -la certs/
    else
        echo "憑證目錄不存在"
    fi
    
    echo -e "${GREEN}✅ 驗證完成${NC}"
}

# 顯示設定摘要
show_summary() {
    echo -e "${GREEN}🎉 密鑰設定完成！${NC}"
    echo ""
    echo -e "${BLUE}📋 設定摘要:${NC}"
    echo ""
    echo -e "${YELLOW}GitHub Actions 密鑰:${NC}"
    echo "  - OCI 配置密鑰已設定"
    echo "  - Vercel 配置密鑰已設定"
    echo ""
    echo -e "${YELLOW}Kubernetes 密鑰:${NC}"
    echo "  - 資料庫密碼已設定"
    echo "  - JWT 密鑰已設定"
    echo "  - mTLS 憑證已生成"
    echo ""
    echo -e "${YELLOW}下一步:${NC}"
    echo "  1. 更新 k8s/manifests 中的映像 URL"
    echo "  2. 執行部署腳本: ./scripts/deploy-oci.sh"
    echo "  3. 設定 Vercel 專案"
    echo ""
    echo -e "${YELLOW}重要提醒:${NC}"
    echo "  - 請妥善保管生成的憑證檔案"
    echo "  - 定期更新密鑰和憑證"
    echo "  - 確保網路安全配置正確"
}

# 主選單
show_menu() {
    echo -e "${BLUE}請選擇操作:${NC}"
    echo "1) 設定 GitHub Actions 密鑰"
    echo "2) 設定 Kubernetes 密鑰"
    echo "3) 生成 mTLS 憑證"
    echo "4) 驗證設定"
    echo "5) 全部設定"
    echo "6) 退出"
    echo ""
    read -p "請輸入選項 (1-6): " choice
}

# 主執行流程
main() {
    check_prerequisites
    
    while true; do
        show_menu
        case $choice in
            1)
                setup_github_secrets
                ;;
            2)
                setup_kubernetes_secrets
                ;;
            3)
                generate_mtls_certs
                ;;
            4)
                verify_setup
                ;;
            5)
                setup_github_secrets
                setup_kubernetes_secrets
                generate_mtls_certs
                verify_setup
                show_summary
                break
                ;;
            6)
                echo -e "${GREEN}👋 再見！${NC}"
                exit 0
                ;;
            *)
                echo -e "${RED}❌ 無效選項，請重新選擇${NC}"
                ;;
        esac
        echo ""
        read -p "按 Enter 繼續..."
        echo ""
    done
}

# 執行主函數
main "$@"
