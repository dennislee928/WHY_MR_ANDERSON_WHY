#!/bin/bash
# Pandora Box Console - 本地構建腳本（Linux/macOS）
# 用於在本地環境構建地端部署版本

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 參數
TARGET="${1:-all}"
VERSION="${2:-dev}"
SKIP_FRONTEND="${SKIP_FRONTEND:-false}"
SKIP_BACKEND="${SKIP_BACKEND:-false}"
CLEAN="${CLEAN:-false}"

echo -e "${CYAN}=====================================${NC}"
echo -e "${CYAN}  Pandora Box Console 本地構建工具   ${NC}"
echo -e "${CYAN}=====================================${NC}"
echo ""

# 目錄設定
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
BACKEND_DIR="$SCRIPT_DIR/be"
FRONTEND_DIR="$SCRIPT_DIR/Fe"
DIST_DIR="$SCRIPT_DIR/dist"
BUILD_DATE=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${YELLOW}📋 構建資訊:${NC}"
echo "   版本: $VERSION"
echo "   構建日期: $BUILD_DATE"
echo "   Git Commit: $GIT_COMMIT"
echo "   目標: $TARGET"
echo ""

# 清理
if [ "$CLEAN" = "true" ]; then
    echo -e "${YELLOW}🧹 清理舊的構建產物...${NC}"
    rm -rf "$DIST_DIR"
    echo -e "${GREEN}✅ 清理完成${NC}"
    echo ""
fi

# 創建輸出目錄
mkdir -p "$DIST_DIR"/{backend,frontend}

# 構建後端
if [ "$SKIP_BACKEND" != "true" ]; then
    echo -e "${YELLOW}🔨 構建後端...${NC}"
    
    # 檢查 Go 是否安裝
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ 錯誤: 未找到 Go。請安裝 Go 1.24 或更高版本。${NC}"
        exit 1
    fi
    
    GO_VERSION=$(go version)
    echo "   使用 Go: $GO_VERSION"
    
    # 設定環境變數
    export CGO_ENABLED=0
    export GOOS=${GOOS:-linux}
    export GOARCH=${GOARCH:-amd64}
    
    LDFLAGS="-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_DATE -X main.GitCommit=$GIT_COMMIT"
    
    cd "$ROOT_DIR"
    
    # 下載依賴
    echo "   正在下載 Go 依賴..."
    go mod download
    
    # 構建 Agent
    echo "   正在構建 Agent..."
    go build -ldflags "$LDFLAGS" -o "$DIST_DIR/backend/pandora-agent" ./cmd/agent/main.go
    
    # 構建 Console
    echo "   正在構建 Console..."
    go build -ldflags "$LDFLAGS" -o "$DIST_DIR/backend/pandora-console" ./cmd/console/main.go
    
    # 構建 UI Server
    echo "   正在構建 UI Server..."
    go build -ldflags "$LDFLAGS" -o "$DIST_DIR/backend/axiom-ui" ./cmd/ui/main.go
    
    # 設定執行權限
    chmod +x "$DIST_DIR/backend"/*
    
    # 複製配置檔案
    echo "   正在複製配置檔案..."
    cp -r "$ROOT_DIR/configs" "$DIST_DIR/backend/"
    
    echo -e "${GREEN}✅ 後端構建完成${NC}"
    echo ""
fi

# 構建前端
if [ "$SKIP_FRONTEND" != "true" ]; then
    echo -e "${YELLOW}🎨 構建前端...${NC}"
    
    # 檢查 Node.js 是否安裝
    if ! command -v node &> /dev/null; then
        echo -e "${RED}❌ 錯誤: 未找到 Node.js。請安裝 Node.js 18 或更高版本。${NC}"
        exit 1
    fi
    
    NODE_VERSION=$(node --version)
    echo "   使用 Node.js: $NODE_VERSION"
    
    cd "$FRONTEND_DIR"
    
    # 安裝依賴
    if [ ! -d "node_modules" ]; then
        echo "   正在安裝依賴..."
        npm install
    fi
    
    # 構建前端
    echo "   正在構建前端應用程式..."
    export NEXT_PUBLIC_APP_VERSION="$VERSION"
    export NODE_ENV="production"
    npm run build
    
    # 複製構建產物
    echo "   正在複製構建產物..."
    [ -d ".next/standalone" ] && cp -r .next/standalone/* "$DIST_DIR/frontend/" || echo "   警告: 未找到 standalone 輸出"
    [ -d ".next/static" ] && cp -r .next/static "$DIST_DIR/frontend/.next/" || echo "   警告: 未找到 static 輸出"
    [ -d "public" ] && cp -r public "$DIST_DIR/frontend/" 2>/dev/null || true
    
    echo -e "${GREEN}✅ 前端構建完成${NC}"
    echo ""
fi

# 創建啟動腳本
echo -e "${YELLOW}📝 創建啟動腳本...${NC}"

cat > "$DIST_DIR/start.sh" <<'EOF'
#!/bin/bash

echo "====================================="
echo "   Pandora Box Console IDS-IPS"
echo "   版本: VERSION_PLACEHOLDER"
echo "====================================="
echo ""

# 設定環境變數
export LOG_LEVEL=info
export DEVICE_PORT=/dev/ttyUSB0
export CONFIG_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/backend/configs" && pwd)"

echo "正在啟動服務..."
echo ""

# 啟動後端服務
cd "$(dirname "${BASH_SOURCE[0]}")/backend"

nohup ./pandora-agent --config "$CONFIG_DIR/agent-config.yaml" > logs/agent.log 2>&1 &
echo "✓ Pandora Agent 已啟動 (PID: $!)"
sleep 2

nohup ./pandora-console --config "$CONFIG_DIR/console-config.yaml" > logs/console.log 2>&1 &
echo "✓ Pandora Console 已啟動 (PID: $!)"
sleep 2

nohup ./axiom-ui --config "$CONFIG_DIR/ui-config.yaml" > logs/ui.log 2>&1 &
echo "✓ Axiom UI 已啟動 (PID: $!)"

echo ""
echo "====================================="
echo "   所有服務已啟動！"
echo "====================================="
echo ""
echo "訪問 Web 介面: http://localhost:3001"
echo "訪問 Grafana: http://localhost:3000"
echo "訪問 Prometheus: http://localhost:9090"
echo ""
echo "查看日誌: tail -f backend/logs/*.log"
echo ""
EOF

sed -i "s/VERSION_PLACEHOLDER/$VERSION/g" "$DIST_DIR/start.sh"
chmod +x "$DIST_DIR/start.sh"

cat > "$DIST_DIR/stop.sh" <<'EOF'
#!/bin/bash

echo "正在停止 Pandora Box Console 服務..."

pkill -f pandora-agent
pkill -f pandora-console
pkill -f axiom-ui

echo "所有服務已停止。"
EOF

chmod +x "$DIST_DIR/stop.sh"

echo -e "${GREEN}✅ 啟動腳本已創建${NC}"
echo ""

# 創建 README
cat > "$DIST_DIR/README.txt" <<EOF
Pandora Box Console IDS-IPS v$VERSION
=====================================

構建資訊
--------
版本: $VERSION
構建日期: $BUILD_DATE
Git Commit: $GIT_COMMIT

快速開始
--------

1. 確保已安裝必要的依賴：
   - PostgreSQL 14+
   - Redis 7+

2. 編輯配置檔案（位於 backend/configs/）

3. 執行 ./start.sh 啟動所有服務

4. 訪問 http://localhost:3001 使用 Web 介面

停止服務
--------
執行 ./stop.sh 停止所有服務

服務端口
--------
- Axiom UI: 3001
- Grafana: 3000
- Prometheus: 9090
- Agent API: 8080

日誌位置
--------
- Agent: backend/logs/agent.log
- Console: backend/logs/console.log
- UI: backend/logs/ui.log

技術支援
--------
問題回報: https://github.com/your-org/pandora_box_console_IDS-IPS/issues
電子郵件: support@pandora-ids.com

授權條款
--------
MIT License - 詳見 LICENSE 檔案
EOF

echo ""
echo -e "${GREEN}=====================================${NC}"
echo -e "${GREEN}  ✅ 構建完成！${NC}"
echo -e "${GREEN}=====================================${NC}"
echo ""
echo -e "${CYAN}構建產物位於: $DIST_DIR${NC}"
echo ""
echo -e "${YELLOW}目錄結構:${NC}"
echo "  backend/          - 後端程式"
echo "  frontend/         - 前端程式"
echo "  start.sh          - 啟動所有服務"
echo "  stop.sh           - 停止所有服務"
echo "  README.txt        - 說明文件"
echo ""
echo -e "${YELLOW}下一步:${NC}"
echo "  1. cd $DIST_DIR"
echo "  2. 編輯 backend/configs/ 中的配置檔案"
echo "  3. 執行 ./start.sh 啟動服務"
echo ""

