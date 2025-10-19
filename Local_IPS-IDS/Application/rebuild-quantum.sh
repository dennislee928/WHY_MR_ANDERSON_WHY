#!/bin/bash
# ============================================================================
# Pandora Cyber AI/Quantum - 重新建構與部署腳本 (Linux/macOS)
# ============================================================================
# 用途: 重新建構 cyber-ai-quantum 服務並重新啟動
# 執行方式: ./rebuild-quantum.sh
# ============================================================================

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 參數處理
NO_BUILD=false
CLEAN=false
IBM_TOKEN="${IBM_QUANTUM_TOKEN:-}"

while [[ $# -gt 0 ]]; do
    case $1 in
        --no-build)
            NO_BUILD=true
            shift
            ;;
        --clean)
            CLEAN=true
            shift
            ;;
        --token)
            IBM_TOKEN="$2"
            shift 2
            ;;
        *)
            echo -e "${RED}[ERROR] 未知參數: $1${NC}"
            exit 1
            ;;
    esac
done

echo -e "${CYAN}============================================================================${NC}"
echo -e "${CYAN}  Pandora Cyber AI/Quantum - 重新建構與部署${NC}"
echo -e "${CYAN}============================================================================${NC}"
echo ""

# 切換到腳本所在目錄
cd "$(dirname "$0")"

# 檢查 Docker
echo -e "${YELLOW}[1/6] 檢查 Docker 環境...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}[ERROR] Docker 未安裝${NC}"
    exit 1
fi

if ! docker info &> /dev/null; then
    echo -e "${RED}[ERROR] Docker daemon 未運行${NC}"
    exit 1
fi

DOCKER_VERSION=$(docker version --format '{{.Server.Version}}')
echo -e "${GREEN}[OK] Docker 版本: ${DOCKER_VERSION}${NC}"

# 設定環境變數
if [ -n "$IBM_TOKEN" ]; then
    echo -e "\n${YELLOW}[2/6] 設定 IBM Quantum Token...${NC}"
    export IBM_QUANTUM_TOKEN="$IBM_TOKEN"
    echo -e "${GREEN}[OK] Token 已設定 (${IBM_TOKEN:0:8}...)${NC}"
else
    echo -e "\n${YELLOW}[2/6] 警告: 未設定 IBM Quantum Token${NC}"
    echo -e "${CYAN}[INFO] 可使用 --token 參數或設定環境變數 IBM_QUANTUM_TOKEN${NC}"
fi

# 清理
if [ "$CLEAN" = true ]; then
    echo -e "\n${YELLOW}[3/6] 清理舊容器和映像...${NC}"
    
    echo "  停止容器..."
    docker-compose stop cyber-ai-quantum 2>/dev/null || true
    docker-compose rm -f cyber-ai-quantum 2>/dev/null || true
    
    echo "  刪除舊映像..."
    docker rmi application-cyber-ai-quantum 2>/dev/null || true
    
    echo "  清理 dangling 映像..."
    docker image prune -f
    
    echo -e "${GREEN}[OK] 清理完成${NC}"
else
    echo -e "\n${YELLOW}[3/6] 跳過清理 (使用 --clean 參數強制清理)${NC}"
fi

# 建構映像
if [ "$NO_BUILD" = false ]; then
    echo -e "\n${YELLOW}[4/6] 重新建構 cyber-ai-quantum 映像...${NC}"
    echo -e "${CYAN}[INFO] 這可能需要幾分鐘時間...${NC}"
    
    docker-compose build --no-cache cyber-ai-quantum
    
    echo -e "${GREEN}[OK] 映像建構成功${NC}"
else
    echo -e "\n${YELLOW}[4/6] 跳過建構 (使用現有映像)${NC}"
fi

# 停止舊容器
echo -e "\n${YELLOW}[5/6] 停止舊容器...${NC}"
docker-compose stop cyber-ai-quantum
echo -e "${GREEN}[OK] 容器已停止${NC}"

# 啟動新容器
echo -e "\n${YELLOW}[6/6] 啟動新容器...${NC}"
docker-compose up -d cyber-ai-quantum

echo -e "${GREEN}[OK] 容器已啟動${NC}"

# 等待服務就緒
echo -e "\n${YELLOW}[等待] 等待服務就緒...${NC}"
sleep 5

# 健康檢查
echo -e "\n${YELLOW}[健康檢查] 檢查服務狀態...${NC}"
max_attempts=10
attempt=0
healthy=false

while [ $attempt -lt $max_attempts ] && [ "$healthy" = false ]; do
    ((attempt++))
    if curl -s -f http://localhost:8000/health > /dev/null 2>&1; then
        healthy=true
        echo -e "${GREEN}[OK] 服務健康檢查通過！${NC}"
    else
        echo "  嘗試 $attempt/$max_attempts : 服務尚未就緒..."
        sleep 3
    fi
done

if [ "$healthy" = false ]; then
    echo -e "${YELLOW}[WARNING] 健康檢查失敗，但容器可能仍在啟動中${NC}"
    echo -e "${CYAN}[INFO] 請稍後手動檢查: http://localhost:8000/health${NC}"
fi

# 顯示容器狀態
echo -e "\n${YELLOW}[狀態] 容器資訊:${NC}"
docker ps --filter "name=cyber-ai-quantum" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# 顯示日誌
echo -e "\n${YELLOW}[日誌] 最近的日誌 (最後 20 行):${NC}"
echo -e "${CYAN}============================================================================${NC}"
docker-compose logs --tail=20 cyber-ai-quantum
echo -e "${CYAN}============================================================================${NC}"

# 管理指令提示
echo -e "\n${CYAN}============================================================================${NC}"
echo -e "${GREEN}  部署完成！${NC}"
echo -e "${CYAN}============================================================================${NC}"
echo -e "\n${YELLOW}常用指令:${NC}"
echo "  查看日誌:       docker-compose logs -f cyber-ai-quantum"
echo "  進入容器:       docker exec -it cyber-ai-quantum /bin/bash"
echo "  重新啟動:       docker-compose restart cyber-ai-quantum"
echo "  停止服務:       docker-compose stop cyber-ai-quantum"
echo "  查看狀態:       docker-compose ps cyber-ai-quantum"
echo ""
echo -e "${YELLOW}API 端點:${NC}"
echo "  健康檢查:       http://localhost:8000/health"
echo "  API 文檔:       http://localhost:8000/docs"
echo "  接收日誌:       POST http://localhost:8000/api/v1/agent/log"
echo "  Zero Trust:     POST http://localhost:8000/api/v1/zerotrust/predict"
echo ""

echo -e "${GREEN}[完成] 所有操作完成！${NC}"
echo ""

