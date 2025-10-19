#!/bin/bash
# Pandora Box Console - Docker 啟動腳本（Linux/macOS）

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  Pandora Box Console - Docker 啟動   ${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

# 檢查 Docker
echo -e "${YELLOW}檢查 Docker...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${RED}✗ Docker 未安裝${NC}"
    echo -e "${YELLOW}  請安裝 Docker: https://docs.docker.com/get-docker/${NC}"
    exit 1
fi

if ! docker ps &> /dev/null; then
    echo -e "${RED}✗ Docker 未運行${NC}"
    echo -e "${YELLOW}  請啟動 Docker daemon${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Docker 正在運行${NC}"

# 檢查 docker-compose
echo -e "${YELLOW}檢查 docker-compose...${NC}"
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}✗ docker-compose 未安裝${NC}"
    exit 1
fi

echo -e "${GREEN}✓ docker-compose 可用${NC}"
echo ""

# 檢查環境變數檔案
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}⚠️  未找到 .env 檔案${NC}"
    echo -e "${NC}   從 .env.example 複製...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ 已創建 .env 檔案${NC}"
    echo -e "${CYAN}   請編輯 .env 設定您的環境${NC}"
    echo ""
fi

# 啟動服務
echo -e "${YELLOW}啟動所有服務...${NC}"
echo ""

docker-compose up -d

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  ✓ 所有服務已啟動！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${CYAN}🌐 訪問以下 URL：${NC}"
    echo -e "  主介面:      ${GREEN}http://localhost:3001${NC}"
    echo -e "  Grafana:     ${GREEN}http://localhost:3000${NC}"
    echo -e "  Prometheus:  ${GREEN}http://localhost:9090${NC}"
    echo -e "  Loki:        ${GREEN}http://localhost:3100${NC}"
    echo -e "  AlertManager: ${GREEN}http://localhost:9093${NC}"
    echo -e "  Prometheus Node Exporter: ${GREEN}http://localhost:9100${NC}"
    echo -e "  Axiom UI: ${GREEN}http://localhost:3001${NC}"
    echo -e "  Cyber AI/Quantum: ${GREEN}http://localhost:8000${NC}"
    echo -e "  Portainer: ${GREEN}http://localhost:9000${NC}"
    echo -e "  RabbitMQ: ${GREEN}http://localhost:15672${NC}"
    echo -e "  PostgreSQL: ${GREEN}http://localhost:5432${NC}"
    echo -e "  Redis: ${GREEN}http://localhost:6379${NC}"
    echo -e "  Prometheus Node Exporter: ${GREEN}http://localhost:9100${NC}"
    echo -e "  Nginx: ${GREEN}http://localhost:443${NC}"
    echo -e "  Promtail: ${GREEN}http://localhost:8080${NC}"
    echo -e "  Node Exporter: ${GREEN}http://localhost:9100${NC}"
    echo ""
    echo -e "${CYAN}🔐 Grafana 預設帳號：${NC}"
    echo -e "  使用者名稱: ${GREEN}admin${NC}"
    echo -e "  密碼:       ${GREEN}pandora123${NC}"
    echo ""
    echo -e "${YELLOW}📊 查看服務狀態：${NC}"
    echo -e "  ${NC}docker-compose ps${NC}"
    echo ""
    echo -e "${YELLOW}📝 查看日誌：${NC}"
    echo -e "  ${NC}docker-compose logs -f${NC}"
    echo ""
    echo -e "${YELLOW}🛑 停止服務：${NC}"
    echo -e "  ${NC}docker-compose down${NC}"
    echo ""
else
    echo ""
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}  ✗ 啟動失敗${NC}"
    echo -e "${RED}========================================${NC}"
    echo ""
    echo -e "${YELLOW}請檢查錯誤訊息並修正${NC}"
    exit 1
fi

