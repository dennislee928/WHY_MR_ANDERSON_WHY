#!/bin/bash

# Pandora Box Console IDS-IPS OCI 部署腳本
# 此腳本會將所有後端服務部署到 OCI Kubernetes Cluster

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置變數
NAMESPACE="pandora-box"
OCI_REGION="${OCI_REGION:-us-ashburn-1}"
OCI_NAMESPACE="${OCI_NAMESPACE:-your-oci-namespace}"
OCI_REGISTRY="${OCI_REGION}.ocir.io/${OCI_NAMESPACE}"
CLUSTER_NAME="${CLUSTER_NAME:-pandora-cluster}"

echo -e "${BLUE}🚀 開始部署 Pandora Box Console IDS-IPS 到 OCI...${NC}"

# 檢查必要的工具
check_prerequisites() {
    echo -e "${YELLOW}📋 檢查必要工具...${NC}"
    
    local missing_tools=()
    
    if ! command -v kubectl &> /dev/null; then
        missing_tools+=("kubectl")
    fi
    
    if ! command -v oci &> /dev/null; then
        missing_tools+=("oci-cli")
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

# 配置 OCI CLI
configure_oci() {
    echo -e "${YELLOW}🔧 配置 OCI CLI...${NC}"
    
    if [ -z "$OCI_USER" ] || [ -z "$OCI_TENANCY" ] || [ -z "$OCI_REGION" ]; then
        echo -e "${RED}❌ 請設定以下環境變數:${NC}"
        echo "  export OCI_USER=your-user-id"
        echo "  export OCI_TENANCY=your-tenancy-id"
        echo "  export OCI_REGION=us-ashburn-1"
        echo "  export OCI_FINGERPRINT=your-api-key-fingerprint"
        echo "  export OCI_KEY_FILE=path/to/your/private-key.pem"
        exit 1
    fi
    
    # 設定 OCI CLI 配置
    oci setup config \
        --user $OCI_USER \
        --tenancy $OCI_TENANCY \
        --region $OCI_REGION \
        --fingerprint $OCI_FINGERPRINT \
        --key-file $OCI_KEY_FILE
    
    echo -e "${GREEN}✅ OCI CLI 配置完成${NC}"
}

# 登入 OCI Container Registry
login_oci_registry() {
    echo -e "${YELLOW}🔐 登入 OCI Container Registry...${NC}"
    
    # 獲取 auth token
    OCI_TOKEN=$(oci iam auth-token get \
        --user-id $OCI_USER \
        --description "Docker login token" \
        --query 'data.token' \
        --raw-output 2>/dev/null || echo "")
    
    if [ -z "$OCI_TOKEN" ]; then
        echo -e "${YELLOW}⚠️  未找到現有的 auth token，正在創建新的...${NC}"
        OCI_TOKEN=$(oci iam auth-token create \
            --user-id $OCI_USER \
            --description "Docker login token" \
            --query 'data.token' \
            --raw-output)
    fi
    
    # 登入 Docker Registry
    echo "$OCI_TOKEN" | docker login $OCI_REGISTRY \
        --username="${OCI_NAMESPACE}/${OCI_USER}" \
        --password-stdin
    
    echo -e "${GREEN}✅ 已登入 OCI Container Registry${NC}"
}

# 建置並推送 Docker 映像
build_and_push_images() {
    echo -e "${YELLOW}🐳 建置並推送 Docker 映像...${NC}"
    
    # 建置 Agent 映像
    echo -e "${BLUE}建置 pandora-agent 映像...${NC}"
    docker build -f Dockerfile.agent -t ${OCI_REGISTRY}/pandora-agent:latest .
    docker push ${OCI_REGISTRY}/pandora-agent:latest
    
    # 建置 Console 映像
    echo -e "${BLUE}建置 pandora-console 映像...${NC}"
    docker build -f Dockerfile -t ${OCI_REGISTRY}/pandora-console:latest .
    docker push ${OCI_REGISTRY}/pandora-console:latest
    
    echo -e "${GREEN}✅ Docker 映像推送完成${NC}"
}

# 配置 kubectl 連接到 OCI Kubernetes Cluster
configure_kubectl() {
    echo -e "${YELLOW}⚙️  配置 kubectl 連接到 OCI Kubernetes Cluster...${NC}"
    
    # 獲取 cluster OCID (需要從 OCI Console 獲取)
    if [ -z "$CLUSTER_OCID" ]; then
        echo -e "${RED}❌ 請設定 CLUSTER_OCID 環境變數${NC}"
        echo "  export CLUSTER_OCID=ocid1.cluster.oc1.xxxxx"
        exit 1
    fi
    
    # 下載 kubeconfig
    oci ce cluster create-kubeconfig \
        --cluster-id $CLUSTER_OCID \
        --file ~/.kube/config \
        --region $OCI_REGION \
        --token-version 2.0.0
    
    echo -e "${GREEN}✅ kubectl 配置完成${NC}"
}

# 更新 Kubernetes manifests 中的映像 URL
update_image_urls() {
    echo -e "${YELLOW}🔄 更新 Kubernetes manifests 中的映像 URL...${NC}"
    
    # 替換 manifests 中的映像 URL
    find k8s/ -name "*.yaml" -exec sed -i "s|iad.ocir.io/YOUR_NAMESPACE|${OCI_REGISTRY}|g" {} \;
    
    echo -e "${GREEN}✅ 映像 URL 更新完成${NC}"
}

# 創建 Kubernetes 密鑰
create_secrets() {
    echo -e "${YELLOW}🔑 創建 Kubernetes 密鑰...${NC}"
    
    # 創建 Docker Registry 密鑰
    kubectl create secret docker-registry oci-registry-secret \
        --docker-server=$OCI_REGISTRY \
        --docker-username="${OCI_NAMESPACE}/${OCI_USER}" \
        --docker-password="$OCI_TOKEN" \
        --namespace=$NAMESPACE \
        --dry-run=client -o yaml | kubectl apply -f -
    
    echo -e "${GREEN}✅ Kubernetes 密鑰創建完成${NC}"
}

# 部署到 Kubernetes
deploy_to_kubernetes() {
    echo -e "${YELLOW}🚀 部署到 Kubernetes...${NC}"
    
    # 使用 kustomize 部署
    kubectl apply -k k8s/
    
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
    echo -e "${GREEN}🎉 部署完成！${NC}"
    echo -e "${BLUE}📋 訪問資訊:${NC}"
    echo ""
    echo -e "${YELLOW}Grafana Dashboard:${NC}"
    echo "  URL: https://pandora.yourdomain.com/grafana"
    echo "  用戶名: admin"
    echo "  密碼: pandora123"
    echo ""
    echo -e "${YELLOW}Prometheus Metrics:${NC}"
    echo "  URL: https://pandora.yourdomain.com/prometheus"
    echo ""
    echo -e "${YELLOW}Pandora Console API:${NC}"
    echo "  URL: https://pandora.yourdomain.com/api/v1/health"
    echo ""
    echo -e "${YELLOW}Agent API:${NC}"
    echo "  URL: https://pandora.yourdomain.com/agent/health"
    echo ""
    echo -e "${YELLOW}查看 Pod 狀態:${NC}"
    echo "  kubectl get pods -n $NAMESPACE"
    echo ""
    echo -e "${YELLOW}查看服務狀態:${NC}"
    echo "  kubectl get services -n $NAMESPACE"
}

# 主執行流程
main() {
    check_prerequisites
    configure_oci
    login_oci_registry
    build_and_push_images
    configure_kubectl
    update_image_urls
    create_secrets
    deploy_to_kubernetes
    show_status
    show_access_info
}

# 執行主函數
main "$@"
