#!/bin/bash
# 驗證所有 GitHub Actions workflows 的語法

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}======================================${NC}"
echo -e "${CYAN}  GitHub Actions Workflow 驗證${NC}"
echo -e "${CYAN}======================================${NC}"
echo ""

WORKFLOWS_DIR=".github/workflows"
PASSED=0
FAILED=0

# 檢查是否安裝 actionlint
if ! command -v actionlint &> /dev/null; then
    echo -e "${YELLOW}警告: actionlint 未安裝${NC}"
    echo "安裝指令: "
    echo "  # macOS: brew install actionlint"
    echo "  # Linux: go install github.com/rhysd/actionlint/cmd/actionlint@latest"
    echo ""
    echo "跳過 YAML 語法驗證，僅檢查檔案存在性..."
    echo ""
    
    # 簡單檢查檔案存在
    for file in "$WORKFLOWS_DIR"/*.yml "$WORKFLOWS_DIR"/*.yaml; do
        if [ -f "$file" ]; then
            echo -e "${GREEN}✓${NC} 檔案存在: $(basename "$file")"
            ((PASSED++))
        fi
    done
else
    # 使用 actionlint 驗證
    echo "使用 actionlint 驗證 workflows..."
    echo ""
    
    for file in "$WORKFLOWS_DIR"/*.yml "$WORKFLOWS_DIR"/*.yaml; do
        if [ -f "$file" ]; then
            filename=$(basename "$file")
            echo -n "驗證: $filename ... "
            
            if actionlint "$file" > /dev/null 2>&1; then
                echo -e "${GREEN}✓ 通過${NC}"
                ((PASSED++))
            else
                echo -e "${RED}✗ 失敗${NC}"
                actionlint "$file"
                ((FAILED++))
            fi
        fi
    done
fi

echo ""
echo -e "${CYAN}======================================${NC}"
echo -e "${CYAN}  驗證結果${NC}"
echo -e "${CYAN}======================================${NC}"
echo -e "通過: ${GREEN}$PASSED${NC}"
echo -e "失敗: ${RED}$FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ 所有 workflows 驗證通過！${NC}"
    exit 0
else
    echo -e "${RED}✗ 有 $FAILED 個 workflows 驗證失敗${NC}"
    exit 1
fi

