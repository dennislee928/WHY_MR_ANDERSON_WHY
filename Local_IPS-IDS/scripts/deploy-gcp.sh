#!/bin/bash

# Pandora Box Console IDS-IPS GCP 部署腳本
# 此腳本會將所有後端服務部署到 Google Cloud Platform (GCP) Kubernetes Engine

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置變數
NAMESPACE="pandora-box"
GCP_PROJECT_ID="${GCP_PROJECT_ID:-your-project-id}"
GCP_REGION="${GCP_REGION:-us-central1}"
GCP_ZONE="${GCP_ZONE:-us-central1-a}"
GCP_CLUSTER_NAME="${GCP_CLUSTER_NAME:-pandora-cluster}"
GCP_REGISTRY="gcr.io/${GCP_PROJECT_ID}"

echo -e "${BLUE}🚀 開始部署 Pandora Box Console IDS-IPS 到 GCP...${NC}"

# 檢查必要的工具
check_prerequisites() {
    echo -e "${YELLOW}📋 檢查必要工具...${NC}"
    
    local missing_tools=()
    
    if ! command -v kubectl &> /dev/null; then
        missing_tools+=("kubectl")
    fi
    
    if ! command -v gcloud &> /dev/null; then
        missing_tools+=("gcloud")
    fi
    
    if ! command -v docker &> /dev/null; then
        missing_tools+=("docker")
    fi
    
    if [ ${#missing_tools[@]} -ne 0 ]; then
        echo -e "${RED}❌ 缺少必要工具: ${missing_tools[*]}${NC}"
        echo "請安裝缺少的工具後重試"
        exit 1
    fi
    
    echo -e "${GREEN}✅ 所有必要工具已安裝${NC}"
}

# 配置 GCP CLI
configure_gcp() {
    echo -e "${YELLOW}🔧 配置 GCP CLI...${NC}"
    
    if [ -z "$GCP_PROJECT_ID" ] || [ "$GCP_PROJECT_ID" = "your-project-id" ]; then
        echo -e "${RED}❌ 請設定 GCP_PROJECT_ID 環境變數:${NC}"
        echo "  export GCP_PROJECT_ID=your-actual-project-id"
        exit 1
    fi
    
    # 設定預設專案
    gcloud config set project $GCP_PROJECT_ID
    gcloud config set compute/region $GCP_REGION
    gcloud config set compute/zone $GCP_ZONE
    
    echo -e "${GREEN}✅ GCP CLI 配置完成${NC}"
}

# 啟用必要的 GCP APIs
enable_apis() {
    echo -e "${YELLOW}🔌 啟用必要的 GCP APIs...${NC}"
    
    gcloud services enable container.googleapis.com
    gcloud services enable containerregistry.googleapis.com
    gcloud services enable compute.googleapis.com
    gcloud services enable dns.googleapis.com
    
    echo -e "${GREEN}✅ GCP APIs 已啟用${NC}"
}

# 配置 Docker 使用 gcloud 認證
configure_docker() {
    echo -e "${YELLOW}🐳 配置 Docker 使用 gcloud 認證...${NC}"
    
    gcloud auth configure-docker
    
    echo -e "${GREEN}✅ Docker 認證配置完成${NC}"
}

# 建置並推送 Docker 映像
build_and_push_images() {
    echo -e "${YELLOW}🐳 建置並推送 Docker 映像到 GCR...${NC}"
    
    # 建置 Agent 映像
    echo -e "${BLUE}建置 pandora-agent 映像...${NC}"
    docker build -f Dockerfile.agent -t ${GCP_REGISTRY}/pandora-agent:latest .
    docker push ${GCP_REGISTRY}/pandora-agent:latest
    
    # 建置 Console 映像
    echo -e "${BLUE}建置 pandora-console 映像...${NC}"
    docker build -f Dockerfile -t ${GCP_REGISTRY}/pandora-console:latest .
    docker push ${GCP_REGISTRY}/pandora-console:latest
    
    echo -e "${GREEN}✅ Docker 映像推送完成${NC}"
}

# 配置 kubectl 連接到 GKE
configure_kubectl() {
    echo -e "${YELLOW}⚙️  配置 kubectl 連接到 GKE...${NC}"
    
    # 獲取 GKE 憑證
    gcloud container clusters get-credentials $GCP_CLUSTER_NAME \
        --zone $GCP_ZONE \
        --project $GCP_PROJECT_ID
    
    echo -e "${GREEN}✅ kubectl 配置完成${NC}"
}

# 更新 Kubernetes manifests 中的映像 URL
update_image_urls() {
    echo -e "${YELLOW}🔄 更新 Kubernetes manifests 中的映像 URL...${NC}"
    
    # 替換 manifests 中的映像 URL
    find k8s-gcp/ -name "*.yaml" -exec sed -i "s|gcr.io/YOUR_PROJECT_ID|${GCP_REGISTRY}|g" {} \;
    
    echo -e "${GREEN}✅ 映像 URL 更新完成${NC}"
}

# 創建 Kubernetes 密鑰
create_secrets() {
    echo -e "${YELLOW}🔑 創建 Kubernetes 密鑰...${NC}"
    
    # 創建 Docker Registry 密鑰 (使用 gcloud 認證)
    kubectl create secret docker-registry gcr-registry-secret \
        --docker-server=gcr.io \
        --docker-username=_json_key \
        --docker-password="$(gcloud auth print-access-token)" \
        --namespace=$NAMESPACE \
        --dry-run=client -o yaml | kubectl apply -f -
    
    echo -e "${GREEN}✅ Kubernetes 密鑰創建完成${NC}"
}

# 部署到 Kubernetes
deploy_to_kubernetes() {
    echo -e "${YELLOW}🚀 部署到 Kubernetes...${NC}"
    
    # 使用 kustomize 部署
    kubectl apply -k k8s-gcp/
    
    # 等待部署完成
    echo -e "${BLUE}⏳ 等待部署完成...${NC}"
    kubectl wait --for=condition=available --timeout=300s deployment/prometheus -n $NAMESPACE
    kubectl wait --for=condition=available --timeout=300s deployment/loki -n $NAMESPACE
    kubectl wait --for=condition=available --timeout=300s deployment/grafana -n $NAMESPACE
    kubectl wait --for=condition=available --timeout=300s deployment/postgres -n $NAMESPACE
    kubectl wait --for=condition=available --timeout=300s deployment/redis -n $NAMESPACE
    kubectl wait --for=condition=available --timeout=300s deployment/pandora-console -n $NAMESPACE
    
    echo -e "${GREEN}✅ 部署完成${NC}"
}

# 設定靜態 IP 和 SSL 憑證
setup_networking() {
    echo -e "${YELLOW}🌐 設定網路和 SSL...${NC}"
    
    # 創建靜態 IP
    gcloud compute addresses create pandora-gcp-ip \
        --global \
        --project=$GCP_PROJECT_ID || echo "靜態 IP 已存在"
    
    # 獲取靜態 IP
    STATIC_IP=$(gcloud compute addresses describe pandora-gcp-ip \
        --global \
        --project=$GCP_PROJECT_ID \
        --format="value(address)")
    
    echo -e "${BLUE}靜態 IP: ${STATIC_IP}${NC}"
    echo -e "${YELLOW}請將您的網域指向此 IP: ${STATIC_IP}${NC}"
    
    echo -e "${GREEN}✅ 網路設定完成${NC}"
}

# 顯示部署狀態
show_status() {
    echo -e "${BLUE}📊 部署狀態:${NC}"
    kubectl get pods -n $NAMESPACE
    echo ""
    kubectl get services -n $NAMESPACE
    echo ""
    kubectl get ingress -n $NAMESPACE
}

# 顯示訪問資訊
show_access_info() {
    echo -e "${GREEN}🎉 GCP 部署完成！${NC}"
    echo -e "${BLUE}📋 訪問資訊:${NC}"
    echo ""
    echo -e "${YELLOW}Grafana Dashboard:${NC}"
    echo "  URL: https://pandora-gcp.yourdomain.com/grafana"
    echo "  用戶名: admin"
    echo "  密碼: pandora123"
    echo ""
    echo -e "${YELLOW}Prometheus Metrics:${NC}"
    echo "  URL: https://pandora-gcp.yourdomain.com/prometheus"
    echo ""
    echo -e "${YELLOW}Pandora Console API:${NC}"
    echo "  URL: https://pandora-gcp.yourdomain.com/api/v1/health"
    echo ""
    echo -e "${YELLOW}Agent API:${NC}"
    echo "  URL: https://pandora-gcp.yourdomain.com/agent/health"
    echo ""
    echo -e "${YELLOW}查看 Pod 狀態:${NC}"
    echo "  kubectl get pods -n $NAMESPACE"
    echo ""
    echo -e "${YELLOW}查看服務狀態:${NC}"
    echo "  kubectl get services -n $NAMESPACE"
    echo ""
    echo -e "${YELLOW}查看 Ingress 狀態:${NC}"
    echo "  kubectl get ingress -n $NAMESPACE"
    echo ""
    echo -e "${YELLOW}查看 ManagedCertificate 狀態:${NC}"
    echo "  kubectl get managedcertificate -n $NAMESPACE"
}

# 清理函數
cleanup() {
    echo -e "${YELLOW}🧹 清理資源...${NC}"
    
    # 刪除部署
    kubectl delete -k k8s-gcp/ --ignore-not-found=true
    
    # 刪除靜態 IP
    gcloud compute addresses delete pandora-gcp-ip \
        --global \
        --project=$GCP_PROJECT_ID \
        --quiet || echo "靜態 IP 不存在"
    
    echo -e "${GREEN}✅ 清理完成${NC}"
}

# 主選單
show_menu() {
    echo -e "${BLUE}請選擇操作:${NC}"
    echo "1) 完整部署"
    echo "2) 僅建置和推送映像"
    echo "3) 僅部署到 Kubernetes"
    echo "4) 設定網路和 SSL"
    echo "5) 顯示狀態"
    echo "6) 清理資源"
    echo "7) 退出"
    echo ""
    read -p "請輸入選項 (1-7): " choice
}

# 主執行流程
main() {
    check_prerequisites
    configure_gcp
    enable_apis
    configure_docker
    
    while true; do
        show_menu
        case $choice in
            1)
                build_and_push_images
                configure_kubectl
                update_image_urls
                create_secrets
                deploy_to_kubernetes
                setup_networking
                show_status
                show_access_info
                break
                ;;
            2)
                build_and_push_images
                ;;
            3)
                configure_kubectl
                update_image_urls
                create_secrets
                deploy_to_kubernetes
                show_status
                ;;
            4)
                setup_networking
                ;;
            5)
                show_status
                ;;
            6)
                cleanup
                ;;
            7)
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
