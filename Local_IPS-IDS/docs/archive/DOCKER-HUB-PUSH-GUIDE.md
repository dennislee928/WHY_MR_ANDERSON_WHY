# Docker Hub 推送指南

## 🚀 快速使用

### Git Bash / Linux

```bash
cd ~/Documents/GitHub/Local_IPS-IDS

# 方式 1: 使用腳本（推薦）
./scripts/push-to-dockerhub.sh

# 方式 2: 設定帳號後執行
export DOCKERHUB_USERNAME="你的Docker Hub帳號"
./scripts/push-to-dockerhub.sh

# 方式 3: 指定版本
VERSION=v3.4.1 ./scripts/push-to-dockerhub.sh
```

### Windows PowerShell

```powershell
cd C:\Users\dennis.lee\Documents\GitHub\Local_IPS-IDS

# 方式 1: 使用腳本（推薦）
.\scripts\push-to-dockerhub.ps1

# 方式 2: 指定參數
.\scripts\push-to-dockerhub.ps1 -Username "你的帳號" -Version "v3.4.1"

# 方式 3: 使用環境變數
$env:DOCKERHUB_USERNAME="你的帳號"
.\scripts\push-to-dockerhub.ps1
```

---

## 📝 手動命令

如果您想手動執行，以下是完整命令：

### 1. 登入 Docker Hub

```bash
docker login
# 輸入用戶名和密碼
```

### 2. 標記映像

```bash
# 設定變數
DOCKERHUB_USERNAME="你的Docker Hub帳號"
VERSION="v3.4.1"

# 標記 axiom-be
docker tag application-axiom-be:latest $DOCKERHUB_USERNAME/axiom-be:$VERSION
docker tag application-axiom-be:latest $DOCKERHUB_USERNAME/axiom-be:latest

# 標記 axiom-ui
docker tag application-axiom-ui:latest $DOCKERHUB_USERNAME/axiom-ui:$VERSION
docker tag application-axiom-ui:latest $DOCKERHUB_USERNAME/axiom-ui:latest

# 標記 pandora-agent
docker tag application-pandora-agent:latest $DOCKERHUB_USERNAME/pandora-agent:$VERSION
docker tag application-pandora-agent:latest $DOCKERHUB_USERNAME/pandora-agent:latest

# 標記 cyber-ai-quantum
docker tag application-cyber-ai-quantum:latest $DOCKERHUB_USERNAME/cyber-ai-quantum:$VERSION
docker tag application-cyber-ai-quantum:latest $DOCKERHUB_USERNAME/cyber-ai-quantum:latest
```

### 3. 推送映像

```bash
# 推送 axiom-be
docker push $DOCKERHUB_USERNAME/axiom-be:$VERSION
docker push $DOCKERHUB_USERNAME/axiom-be:latest

# 推送 axiom-ui
docker push $DOCKERHUB_USERNAME/axiom-ui:$VERSION
docker push $DOCKERHUB_USERNAME/axiom-ui:latest

# 推送 pandora-agent
docker push $DOCKERHUB_USERNAME/pandora-agent:$VERSION
docker push $DOCKERHUB_USERNAME/pandora-agent:latest

# 推送 cyber-ai-quantum
docker push $DOCKERHUB_USERNAME/cyber-ai-quantum:$VERSION
docker push $DOCKERHUB_USERNAME/cyber-ai-quantum:latest
```

---

## 📊 推送的映像

| 本地映像 | Docker Hub 映像 | 大小 |
|---------|----------------|------|
| application-axiom-be:latest | `你的帳號/axiom-be:v3.4.1` | ~50MB |
| application-axiom-ui:latest | `你的帳號/axiom-ui:v3.4.1` | ~50MB |
| application-pandora-agent:latest | `你的帳號/pandora-agent:v3.4.1` | ~50MB |
| application-cyber-ai-quantum:latest | `你的帳號/cyber-ai-quantum:v3.4.1` | ~2.2GB |

**總大小**: 約 2.4GB

---

## 🔐 登入方式

### 方式 1: 互動式登入

```bash
docker login
# Username: 你的帳號
# Password: 你的密碼或 Personal Access Token
```

### 方式 2: 使用 Token（推薦）

```bash
# 從標準輸入讀取密碼
echo "你的Token" | docker login --username 你的帳號 --password-stdin
```

### 方式 3: 環境變數

```bash
export DOCKER_USERNAME="你的帳號"
export DOCKER_PASSWORD="你的Token"

echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
```

---

## ⏱️ 預計時間

| 映像 | 大小 | 上傳時間（估計） |
|------|------|-----------------|
| axiom-be | ~50MB | ~2 分鐘 |
| axiom-ui | ~50MB | ~2 分鐘 |
| pandora-agent | ~50MB | ~2 分鐘 |
| cyber-ai-quantum | ~2.2GB | ~15-30 分鐘 |

**總計**: 約 20-40 分鐘（取決於網路速度）

---

## 📋 推送後的使用

### 拉取映像

```bash
# 其他人可以拉取您的映像
docker pull 你的帳號/cyber-ai-quantum:v3.4.1
docker pull 你的帳號/pandora-agent:latest
```

### 更新 docker-compose.yml

```yaml
services:
  cyber-ai-quantum:
    image: 你的帳號/cyber-ai-quantum:v3.4.1
    # 不需要 build 配置
```

---

## 🔄 更新和重新推送

```bash
# 重新建構映像
cd Application
docker-compose build cyber-ai-quantum

# 重新標記和推送
./scripts/push-to-dockerhub.sh
```

---

## 📊 檢查映像資訊

```bash
# 列出本地映像
docker images | grep application

# 檢查映像大小
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}" | grep application

# 查看映像歷史
docker history application-cyber-ai-quantum:latest
```

---

## 🎯 完整流程範例

```bash
# 1. 切換目錄
cd ~/Documents/GitHub/Local_IPS-IDS

# 2. 設定 Docker Hub 帳號
export DOCKERHUB_USERNAME="dennis-lee"  # 改為您的帳號

# 3. 執行推送腳本
./scripts/push-to-dockerhub.sh

# 4. 等待完成（約 20-40 分鐘）

# 5. 驗證
# 訪問 https://hub.docker.com/u/dennis-lee
```

---

## ✅ 快速命令（一鍵執行）

**Git Bash**:
```bash
cd ~/Documents/GitHub/Local_IPS-IDS && export DOCKERHUB_USERNAME="你的帳號" && ./scripts/push-to-dockerhub.sh
```

**PowerShell**:
```powershell
cd C:\Users\dennis.lee\Documents\GitHub\Local_IPS-IDS; .\scripts\push-to-dockerhub.ps1 -Username "你的帳號"
```

---

## 🎉 完成後

您的映像將可在以下位置訪問：

- https://hub.docker.com/r/你的帳號/axiom-be
- https://hub.docker.com/r/你的帳號/axiom-ui
- https://hub.docker.com/r/你的帳號/pandora-agent
- https://hub.docker.com/r/你的帳號/cyber-ai-quantum

---

**創建日期**: 2025-10-15  
**版本**: v3.4.1  
**狀態**: ✅ 就緒

