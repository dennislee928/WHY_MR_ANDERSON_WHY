#!/bin/bash

# Axiom Backend V3 本地構建和測試腳本 (Bash)
# 版本: 3.1.0
# 使用方式：
#   ./build-axiom-be-v3.sh           # 正常構建
#   ./build-axiom-be-v3.sh --no-cache # 不使用快取構建

# 獲取腳本所在目錄
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR" || exit 1

echo "============================================"
echo "  Axiom Backend V3 本地構建和測試"
echo "  版本: 3.1.0"
echo "============================================"
echo ""
echo "工作目錄: $SCRIPT_DIR"
echo ""

# 設置變量
IMAGE_NAME="axiom-backend"
IMAGE_TAG="v3.1.0"
CONTAINER_NAME="axiom-be-v3-test"
PORT=3001

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Step 1: 清理舊容器和鏡像
echo -e "${YELLOW}[1/6] 清理舊容器...${NC}"

# 清理目標測試容器
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "  停止並刪除容器: $CONTAINER_NAME"
    docker stop $CONTAINER_NAME 2>/dev/null
    docker rm $CONTAINER_NAME 2>/dev/null
fi

# 清理所有佔用 3001 端口的容器
PORT_CONTAINERS=$(docker ps -a --filter "publish=${PORT}" --format "{{.Names}}" 2>/dev/null)
if [ ! -z "$PORT_CONTAINERS" ]; then
    echo "  發現佔用端口 ${PORT} 的容器："
    for container in $PORT_CONTAINERS; do
        echo "    - $container"
        docker stop $container 2>/dev/null
        docker rm $container 2>/dev/null
    done
fi

echo -e "${GREEN}✓ 完成${NC}"
echo ""

# Step 2: 構建 Docker 鏡像
echo -e "${YELLOW}[2/6] 構建 Docker 鏡像...${NC}"
echo -e "${CYAN}鏡像: $IMAGE_NAME:$IMAGE_TAG${NC}"

BUILD_START=$(date +%s)

# 檢查是否需要使用 --no-cache
NO_CACHE=""
if [ "$1" == "--no-cache" ]; then
    NO_CACHE="--no-cache"
    echo -e "${CYAN}使用 --no-cache 選項${NC}"
fi

docker build \
    $NO_CACHE \
    -f docker/axiom-be.dockerfile \
    -t ${IMAGE_NAME}:${IMAGE_TAG} \
    -t ${IMAGE_NAME}:latest \
    be

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 構建失敗！${NC}"
    exit 1
fi

BUILD_END=$(date +%s)
BUILD_TIME=$((BUILD_END - BUILD_START))
echo -e "${GREEN}✓ 構建成功！耗時: ${BUILD_TIME}秒${NC}"
echo ""

# Step 3: 檢查鏡像大小
echo -e "${YELLOW}[3/6] 檢查鏡像...${NC}"
docker images | grep $IMAGE_NAME | head -1
echo ""

# Step 4: 啟動容器（單機測試）
echo -e "${YELLOW}[4/6] 啟動測試容器...${NC}"
docker run -d \
    --name $CONTAINER_NAME \
    -p ${PORT}:3001 \
    -e POSTGRES_HOST=host.docker.internal \
    -e POSTGRES_PORT=5432 \
    -e POSTGRES_USER=pandora \
    -e POSTGRES_PASSWORD=pandora123 \
    -e POSTGRES_DB=pandora_db \
    -e REDIS_HOST=host.docker.internal \
    -e REDIS_PORT=6379 \
    -e REDIS_PASSWORD=pandora123 \
    -e PROMETHEUS_URL=http://host.docker.internal:9090 \
    -e LOKI_URL=http://host.docker.internal:3100 \
    -e QUANTUM_URL=http://host.docker.internal:8000 \
    ${IMAGE_NAME}:${IMAGE_TAG}

if [ $? -ne 0 ]; then
    echo -e "${RED}✗ 容器啟動失敗！${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 容器已啟動${NC}"
echo ""

# Step 5: 等待服務啟動
echo -e "${YELLOW}[5/6] 等待服務就緒...${NC}"
RETRY=0
MAX_RETRY=30

while [ $RETRY -lt $MAX_RETRY ]; do
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:${PORT}/health)
    if [ "$HTTP_CODE" == "200" ]; then
        echo -e "${GREEN}✓ 服務已就緒！${NC}"
        break
    fi
    RETRY=$((RETRY + 1))
    echo -e "  等待中... ($RETRY/$MAX_RETRY)"
    sleep 2
done

if [ $RETRY -eq $MAX_RETRY ]; then
    echo -e "${RED}✗ 服務啟動超時！${NC}"
    echo -e "${YELLOW}查看日誌:${NC}"
    docker logs $CONTAINER_NAME
    exit 1
fi
echo ""

# Step 6: 測試 API
echo -e "${YELLOW}[6/6] 測試 API 端點...${NC}"

# 測試健康檢查
echo "  測試健康檢查..."
HEALTH_RESPONSE=$(curl -s http://localhost:${PORT}/health)
if echo "$HEALTH_RESPONSE" | grep -q "healthy"; then
    echo -e "  ${GREEN}✓ 健康檢查通過${NC}"
else
    echo -e "  ${RED}✗ 健康檢查失敗${NC}"
fi

# 測試 Agent 註冊
echo "  測試 Agent 註冊..."
AGENT_RESPONSE=$(curl -s -X POST http://localhost:${PORT}/api/v2/agent/register \
    -H "Content-Type: application/json" \
    -d '{"mode":"internal","hostname":"test-server","ip_address":"127.0.0.1","capabilities":["windows_logs"]}')

if echo "$AGENT_RESPONSE" | grep -q "success"; then
    AGENT_ID=$(echo "$AGENT_RESPONSE" | grep -oP '(?<="agent_id":")[^"]*' | head -1)
    echo -e "  ${GREEN}✓ Agent 註冊成功: $AGENT_ID${NC}"
else
    echo -e "  ${RED}✗ Agent 註冊失敗${NC}"
fi

# 測試 PII 檢測
echo "  測試 PII 檢測..."
PII_RESPONSE=$(curl -s -X POST http://localhost:${PORT}/api/v2/compliance/pii/detect \
    -H "Content-Type: application/json" \
    -d '{"text":"Contact: test@example.com, Card: 4532-1234-5678-9010"}')

if echo "$PII_RESPONSE" | grep -q "pii_found"; then
    echo -e "  ${GREEN}✓ PII 檢測成功${NC}"
else
    echo -e "  ${RED}✗ PII 檢測失敗${NC}"
fi

# 測試儲存統計
echo "  測試儲存統計..."
STORAGE_RESPONSE=$(curl -s http://localhost:${PORT}/api/v2/storage/tiers/stats)
if echo "$STORAGE_RESPONSE" | grep -q "success"; then
    echo -e "  ${GREEN}✓ 儲存統計成功${NC}"
else
    echo -e "  ${RED}✗ 儲存統計失敗${NC}"
fi

echo ""
echo "============================================"
echo -e "${CYAN}  構建和測試完成！${NC}"
echo "============================================"
echo ""
echo -e "${YELLOW}容器信息:${NC}"
echo "  名稱: $CONTAINER_NAME"
echo "  端口: http://localhost:$PORT"
echo ""
echo -e "${YELLOW}常用指令:${NC}"
echo "  查看日誌: docker logs $CONTAINER_NAME"
echo "  查看日誌(實時): docker logs -f $CONTAINER_NAME"
echo "  停止容器: docker stop $CONTAINER_NAME"
echo "  刪除容器: docker rm $CONTAINER_NAME"
echo ""
echo -e "${YELLOW}API 端點:${NC}"
echo "  健康檢查: http://localhost:$PORT/health"
echo "  Agent 註冊: POST http://localhost:$PORT/api/v2/agent/register"
echo "  PII 檢測: POST http://localhost:$PORT/api/v2/compliance/pii/detect"
echo "  GDPR 刪除: POST http://localhost:$PORT/api/v2/compliance/gdpr/deletion-request"
echo ""
echo -e "${YELLOW}查看完整 API 文檔:${NC}"
echo "  ../docs/AXIOM-BACKEND-V3-API-DOCUMENTATION.md"
echo ""

