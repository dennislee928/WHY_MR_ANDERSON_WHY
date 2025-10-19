#!/bin/bash
# Application/be/build.sh
# Linux/macOS 後端構建腳本

set -e

TARGET="${1:-all}"
VERSION="${2:-dev}"

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}========================================${NC}"
echo -e "${CYAN}  後端構建腳本 (Linux/macOS)${NC}"
echo -e "${CYAN}========================================${NC}"
echo ""

ROOT_DIR="$(cd ../.. && pwd)"

echo -e "${NC}專案根目錄: ${ROOT_DIR}${NC}"
echo -e "${NC}構建目標: ${TARGET}${NC}"
echo -e "${NC}版本: ${VERSION}${NC}"
echo ""

# 檢查 Go 是否安裝
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ 錯誤: 未找到 Go${NC}"
    echo -e "${YELLOW}  請安裝 Go 1.24+ 後再試${NC}"
    exit 1
fi

GO_VERSION=$(go version)
echo -e "${GREEN}✓ Go 環境: ${GO_VERSION}${NC}"

# 檢查根目錄是否正確
if [ ! -f "$ROOT_DIR/go.mod" ]; then
    echo -e "${RED}✗ 錯誤: 找不到 go.mod${NC}"
    echo -e "${YELLOW}  ROOT_DIR: ${ROOT_DIR}${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}開始構建...${NC}"
echo ""

# 使用 make
if [ "$TARGET" == "all" ]; then
    make all
else
    make "$TARGET"
fi

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}  ✓ 構建成功！${NC}"
    echo -e "${GREEN}========================================${NC}"
    echo ""
    echo -e "${CYAN}二進位檔案位於: ./bin/${NC}"
    echo ""
    if [ -d "./bin/" ]; then
        ls -lh ./bin/
    fi
else
    echo ""
    echo -e "${RED}========================================${NC}"
    echo -e "${RED}  ✗ 構建失敗${NC}"
    echo -e "${RED}========================================${NC}"
    exit 1
fi

