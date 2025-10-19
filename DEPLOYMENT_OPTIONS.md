# Security Platform - Cloudflare Deployment Guide

## 🚨 **當前狀態**

由於 Node.js 權限問題 (`EPERM: operation not permitted, lstat 'C:\Users\pclee'`)，我們提供以下部署方案：

## 📋 **部署選項**

### **選項 1: 使用 Docker Compose (推薦 - 本地測試)**

這個方案可以立即在本地運行所有容器服務：

```powershell
# 在 Cloudflare 目錄中執行
cd C:\Users\USER\Desktop\WHY_MR_ANDERSON_WHY\infrastructure\cloud-configs\cloudflare

# 啟動所有容器
docker-compose up -d

# 檢查容器狀態
docker-compose ps

# 查看日誌
docker-compose logs -f

# 測試服務
curl http://localhost:3000/health      # Backend API
curl http://localhost:8000/health      # AI/Quantum
curl http://localhost:8080/health      # Security Tools
curl http://localhost:9090/-/healthy   # Monitoring
```

**優點**：
- ✅ 立即可用，無需解決 Node.js 權限問題
- ✅ 完整的容器化環境
- ✅ 所有服務本地運行
- ✅ 易於測試和除錯

**缺點**：
- ❌ 只在本地運行，不是 Cloudflare Workers

---

### **選項 2: 使用 Cloudflare Dashboard (推薦 - 正式部署)**

透過 Cloudflare Web 介面直接部署 Worker：

#### **步驟 1: 準備 Worker 代碼**
```powershell
# 確保 dist 目錄存在
cd C:\Users\USER\Desktop\WHY_MR_ANDERSON_WHY\infrastructure\cloud-configs\cloudflare

# 手動構建（如果有 dist/index.js）
# 或者直接使用 src/index.js
```

#### **步驟 2: 登入 Cloudflare Dashboard**
1. 前往 https://dash.cloudflare.com/
2. 選擇您的帳戶
3. 點擊 "Workers & Pages"

#### **步驟 3: 創建新 Worker**
1. 點擊 "Create application"
2. 選擇 "Create Worker"
3. 命名為 `security-platform-worker`
4. 點擊 "Deploy"

#### **步驟 4: 上傳代碼**
1. 在 Worker 編輯器中
2. 複製 `src/index.js` 的內容
3. 貼上到編輯器
4. 點擊 "Save and Deploy"

#### **步驟 5: 配置 Bindings**
1. 點擊 "Settings" > "Variables"
2. 根據需要添加：
   - D1 Database
   - KV Namespaces
   - Durable Objects
   - 等等

**優點**：
- ✅ 正式部署到 Cloudflare Workers
- ✅ 全球 CDN 加速
- ✅ 無需本地環境
- ✅ 簡單直觀

**缺點**：
- ❌ 需要手動操作
- ❌ 不適合自動化 CI/CD

---

### **選項 3: 使用 WSL (Windows Subsystem for Linux)**

在 WSL 中運行 Wrangler，避免 Windows 權限問題：

```bash
# 安裝 WSL (如果尚未安裝)
wsl --install

# 在 WSL 中操作
wsl

# 導航到專案目錄
cd /mnt/c/Users/USER/Desktop/WHY_MR_ANDERSON_WHY/infrastructure/cloud-configs/cloudflare

# 安裝 Node.js 和 npm
sudo apt update
sudo apt install nodejs npm -y

# 安裝 Wrangler
npm install -g wrangler

# 登入 Cloudflare
wrangler login

# 部署
wrangler deploy
```

**優點**：
- ✅ 完整的 Linux 環境
- ✅ 沒有 Windows 權限問題
- ✅ 支援 CI/CD
- ✅ 官方支援的部署方式

**缺點**：
- ❌ 需要安裝 WSL
- ❌ 需要額外設定

---

### **選項 4: 使用 Docker 構建 Wrangler 容器**

創建一個 Docker 容器來運行 Wrangler：

```dockerfile
# Dockerfile.wrangler
FROM node:18-alpine

WORKDIR /app

RUN npm install -g wrangler

COPY . .

CMD ["wrangler", "deploy"]
```

```powershell
# 構建 Wrangler 容器
docker build -f Dockerfile.wrangler -t wrangler-deploy .

# 使用容器部署
docker run -it --rm -v ${PWD}:/app wrangler-deploy
```

**優點**：
- ✅ 隔離的環境
- ✅ 可重複使用
- ✅ 適合 CI/CD

**缺點**：
- ❌ 需要額外配置
- ❌ 互動式登入較複雜

---

## 🎯 **推薦部署流程**

### **階段 1: 本地測試 (立即可用)**
```powershell
# 使用 Docker Compose
cd C:\Users\USER\Desktop\WHY_MR_ANDERSON_WHY\infrastructure\cloud-configs\cloudflare
docker-compose up -d
```

### **階段 2: Cloudflare 部署 (生產環境)**
1. 使用 **Cloudflare Dashboard** 手動部署 Worker
2. 或安裝 **WSL** 使用 Wrangler CLI

---

## 📊 **各選項比較表**

| 選項 | 難度 | 速度 | 適用場景 | 推薦度 |
|------|------|------|----------|--------|
| Docker Compose | ⭐ | ⭐⭐⭐ | 本地測試 | ⭐⭐⭐⭐⭐ |
| Cloudflare Dashboard | ⭐⭐ | ⭐⭐ | 正式部署 | ⭐⭐⭐⭐ |
| WSL | ⭐⭐⭐ | ⭐⭐ | CI/CD | ⭐⭐⭐⭐ |
| Docker Wrangler | ⭐⭐⭐⭐ | ⭐ | 進階使用 | ⭐⭐⭐ |

---

## ❓ **需要我協助哪個選項？**

請告訴我您想要：
1. **立即本地測試** → 使用 Docker Compose
2. **部署到 Cloudflare** → 使用 Dashboard 或 WSL
3. **解決 Node.js 問題** → 重新安裝 Node.js 或使用 nvm-windows
4. **其他方案** → 請說明您的需求

---

## 🔧 **快速啟動命令**

### **本地測試 (Docker Compose)**
```powershell
cd C:\Users\USER\Desktop\WHY_MR_ANDERSON_WHY\infrastructure\cloud-configs\cloudflare
docker-compose up -d
docker-compose logs -f
```

### **測試 API**
```powershell
# 測試各個服務
curl http://localhost:3000/health
curl http://localhost:8000/health
curl http://localhost:8080/health
curl http://localhost:5432  # PostgreSQL
curl http://localhost:9090/-/healthy  # Prometheus
```

### **停止服務**
```powershell
docker-compose down
```

---

讓我知道您想要使用哪個選項，我會協助您完成部署！ 🚀
