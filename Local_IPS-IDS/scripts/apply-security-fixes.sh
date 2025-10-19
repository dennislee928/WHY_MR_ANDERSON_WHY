#!/bin/bash
# Pandora Box Console - 安全修復應用腳本
# 自動應用 SAST 掃描發現的安全修復

set -e

echo "========================================="
echo "  🔒 Pandora 安全修復應用工具"
echo "========================================="

# 顏色定義
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 1. 更新 Go 依賴
echo -e "\n${YELLOW}📦 步驟 1/5: 更新 Go 依賴...${NC}"
go mod tidy
go mod download
echo -e "${GREEN}✅ Go 依賴已更新${NC}"

# 2. 更新 Python 依賴
echo -e "\n${YELLOW}📦 步驟 2/5: 更新 Python 依賴...${NC}"
cd Experimental/cyber-ai-quantum
pip install -r requirements.txt --upgrade --quiet
cd ../..
echo -e "${GREEN}✅ Python 依賴已更新${NC}"

# 3. 驗證 Dockerfile USER 指令
echo -e "\n${YELLOW}🔍 步驟 3/5: 驗證 Dockerfile 安全性...${NC}"
DOCKERFILES=(
    "Application/docker/agent.koyeb.dockerfile"
    "Application/docker/monitoring.dockerfile"
    "Application/docker/nginx.dockerfile"
    "Application/docker/test.dockerfile"
    "Application/docker/axiom-be.dockerfile"
)

for dockerfile in "${DOCKERFILES[@]}"; do
    if grep -q "^USER " "$dockerfile"; then
        echo -e "${GREEN}  ✅ $dockerfile - USER 指令已存在${NC}"
    else
        echo -e "${RED}  ❌ $dockerfile - 缺少 USER 指令${NC}"
    fi
done

# 4. 檢查 Alpine 版本
echo -e "\n${YELLOW}🔍 步驟 4/5: 檢查 Alpine 基礎映像版本...${NC}"
for dockerfile in Application/docker/*.dockerfile; do
    if grep -q "FROM alpine:3.21" "$dockerfile" || grep -q "FROM alpine:3.22" "$dockerfile"; then
        echo -e "${GREEN}  ✅ $(basename $dockerfile) - Alpine 版本安全${NC}"
    elif grep -q "FROM alpine:" "$dockerfile"; then
        echo -e "${YELLOW}  ⚠️  $(basename $dockerfile) - 建議更新到 Alpine 3.21+${NC}"
    fi
done

# 5. 重新構建關鍵服務
echo -e "\n${YELLOW}🔨 步驟 5/5: 重新構建 Docker 映像...${NC}"
cd Application

echo -e "${YELLOW}  構建 axiom-be...${NC}"
docker-compose build --no-cache axiom-be

echo -e "${YELLOW}  構建 cyber-ai-quantum...${NC}"
docker-compose build --no-cache cyber-ai-quantum

cd ..
echo -e "${GREEN}✅ Docker 映像已重新構建${NC}"

# 完成
echo -e "\n========================================="
echo -e "${GREEN}  ✅ 安全修復應用完成！${NC}"
echo -e "========================================="
echo -e "\n📋 下一步:"
echo -e "  1. 查看詳細報告: ${YELLOW}docs/SAST-SECURITY-FIXES.md${NC}"
echo -e "  2. 重啟服務: ${YELLOW}cd Application && docker-compose up -d${NC}"
echo -e "  3. 驗證服務: ${YELLOW}docker-compose ps${NC}"
echo -e "\n⚠️  需要手動處理的項目:"
echo -e "  - 配置 gRPC TLS 證書"
echo -e "  - 修復 exec.Command 輸入驗證"
echo -e "  - 修復 RWMutex 死鎖風險"
echo -e "  - 更新 GitHub Actions 配置\n"

