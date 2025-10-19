#!/bin/bash
# ============================================================================
# 推送 Pandora 映像到 Docker Hub
# ============================================================================

set -e

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}============================================================================${NC}"
echo -e "${BLUE}  推送 Pandora 映像到 Docker Hub${NC}"
echo -e "${BLUE}============================================================================${NC}"
echo ""

# Docker Hub 帳號（請修改為您的帳號）
DOCKERHUB_USERNAME="${DOCKERHUB_USERNAME:-}"
VERSION="${VERSION:-v3.4.1}"

# 如果未設定，詢問用戶
if [ -z "$DOCKERHUB_USERNAME" ]; then
    read -p "請輸入您的 Docker Hub 帳號: " DOCKERHUB_USERNAME
fi

if [ -z "$DOCKERHUB_USERNAME" ]; then
    echo -e "${RED}❌ Docker Hub 帳號未設定${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Docker Hub 帳號: $DOCKERHUB_USERNAME${NC}"
echo -e "${BLUE}ℹ️  版本標籤: $VERSION${NC}"
echo ""

# 要推送的映像列表
declare -a IMAGES=(
    "application-axiom-be"
    "application-axiom-ui"
    "application-pandora-agent"
    "application-cyber-ai-quantum"
)

# 步驟 1: 登入 Docker Hub
echo -e "${YELLOW}[1/3] 登入 Docker Hub...${NC}"
echo ""

if docker login; then
    echo -e "${GREEN}✅ 登入成功${NC}"
else
    echo -e "${RED}❌ 登入失敗${NC}"
    exit 1
fi

echo ""

# 步驟 2: 標記映像
echo -e "${YELLOW}[2/3] 標記映像...${NC}"
echo ""

for IMAGE in "${IMAGES[@]}"; do
    LOCAL_IMAGE="${IMAGE}:latest"
    REMOTE_IMAGE="$DOCKERHUB_USERNAME/${IMAGE#application-}:$VERSION"
    REMOTE_LATEST="$DOCKERHUB_USERNAME/${IMAGE#application-}:latest"
    
    echo -e "${BLUE}標記: $LOCAL_IMAGE${NC}"
    echo -e "  → $REMOTE_IMAGE"
    echo -e "  → $REMOTE_LATEST"
    
    if docker tag "$LOCAL_IMAGE" "$REMOTE_IMAGE"; then
        echo -e "${GREEN}  ✅ 版本標籤成功${NC}"
    else
        echo -e "${RED}  ❌ 標記失敗${NC}"
        continue
    fi
    
    if docker tag "$LOCAL_IMAGE" "$REMOTE_LATEST"; then
        echo -e "${GREEN}  ✅ latest 標籤成功${NC}"
    else
        echo -e "${RED}  ❌ 標記失敗${NC}"
    fi
    
    echo ""
done

# 步驟 3: 推送映像
echo -e "${YELLOW}[3/3] 推送映像到 Docker Hub...${NC}"
echo ""

PUSHED_COUNT=0
FAILED_COUNT=0

for IMAGE in "${IMAGES[@]}"; do
    SHORT_NAME="${IMAGE#application-}"
    REMOTE_IMAGE="$DOCKERHUB_USERNAME/$SHORT_NAME:$VERSION"
    REMOTE_LATEST="$DOCKERHUB_USERNAME/$SHORT_NAME:latest"
    
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    echo -e "${BLUE}推送: $SHORT_NAME${NC}"
    echo -e "${BLUE}════════════════════════════════════════${NC}"
    
    # 推送版本標籤
    echo -e "\n推送 $VERSION 標籤..."
    if docker push "$REMOTE_IMAGE"; then
        echo -e "${GREEN}✅ $REMOTE_IMAGE 推送成功${NC}"
    else
        echo -e "${RED}❌ $REMOTE_IMAGE 推送失敗${NC}"
        FAILED_COUNT=$((FAILED_COUNT + 1))
        continue
    fi
    
    # 推送 latest 標籤
    echo -e "\n推送 latest 標籤..."
    if docker push "$REMOTE_LATEST"; then
        echo -e "${GREEN}✅ $REMOTE_LATEST 推送成功${NC}"
        PUSHED_COUNT=$((PUSHED_COUNT + 1))
    else
        echo -e "${RED}❌ $REMOTE_LATEST 推送失敗${NC}"
        FAILED_COUNT=$((FAILED_COUNT + 1))
    fi
    
    echo ""
done

# 總結
echo ""
echo -e "${BLUE}============================================================================${NC}"
echo -e "${BLUE}  推送完成！${NC}"
echo -e "${BLUE}============================================================================${NC}"
echo ""
echo -e "${GREEN}✅ 成功推送: $PUSHED_COUNT 個映像${NC}"
if [ $FAILED_COUNT -gt 0 ]; then
    echo -e "${RED}❌ 失敗: $FAILED_COUNT 個映像${NC}"
fi

echo ""
echo "您的映像現已可用於："
echo ""

for IMAGE in "${IMAGES[@]}"; do
    SHORT_NAME="${IMAGE#application-}"
    echo -e "  docker pull $DOCKERHUB_USERNAME/$SHORT_NAME:$VERSION"
    echo -e "  docker pull $DOCKERHUB_USERNAME/$SHORT_NAME:latest"
    echo ""
done

echo -e "${BLUE}============================================================================${NC}"
echo -e "${GREEN}完成！映像已推送到 Docker Hub${NC}"
echo -e "${BLUE}============================================================================${NC}"

