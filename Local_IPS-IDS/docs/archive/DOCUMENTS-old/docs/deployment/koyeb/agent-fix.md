# Koyeb Pandora Agent 崩潰修復

## 🐛 問題描述

Koyeb 上的 `pandora-agent` 不斷崩潰重啟：

```
2025-10-07 08:06:06,722 INFO success: pandora-agent entered RUNNING state
2025-10-07 08:06:06,723 WARN exited: pandora-agent (exit status 1; not expected)
2025-10-07 08:07:03,465 WARN received SIGTERM indicating exit request
TCP health check failed on port 8080.
```

同時 `promtail` 也無法啟動：
```
2025-10-07 08:06:05,722 WARN exited: promtail (exit status 127; not expected)
```

## 🔍 根本原因

1. **pandora-agent 沒有 HTTP 服務器**
   - Agent 是純命令行應用程式
   - Koyeb 期望 HTTP 服務監聽 8080 端口
   - 健康檢查失敗導致容器被終止

2. **Promtail 二進制文件找不到**
   - Exit status 127 = "command not found"
   - 安裝過程可能有問題

3. **雲端環境沒有實體設備**
   - Agent 期望連接 USB-SERIAL 設備（`/dev/ttyUSB0`）
   - 初始化失敗導致程式退出

## ✅ 修復方案

已對以下文件進行修復：

### 1. `cmd/agent/main.go` - 添加 HTTP 健康檢查服務器

```go
// 新增功能：
// 1. HTTP 健康檢查端點 (/health)
// 2. 優雅關閉處理
// 3. 允許在無實體設備環境下運行
```

#### 主要變更：

- ✅ 添加 `/health` 端點（返回 JSON 狀態）
- ✅ 在 goroutine 中運行 HTTP 服務器
- ✅ 即使設備初始化失敗也繼續運行
- ✅ 支援信號處理（SIGTERM, SIGINT）
- ✅ 從環境變數讀取 PORT

### 2. `Dockerfile.agent.koyeb` - 修正 Promtail 安裝

```dockerfile
# 修正前：
RUN wget https://github.com/grafana/loki/releases/download/v${PROMTAIL_VERSION}/promtail-linux-amd64.zip && \
    unzip promtail-linux-amd64.zip && \
    ...

# 修正後：
RUN cd /tmp && \
    wget -q https://github.com/grafana/loki/releases/download/v${PROMTAIL_VERSION}/promtail-linux-amd64.zip && \
    unzip -q promtail-linux-amd64.zip && \
    mv promtail-linux-amd64 /usr/local/bin/promtail && \
    chmod +x /usr/local/bin/promtail && \
    rm -f promtail-linux-amd64.zip && \
    /usr/local/bin/promtail --version  # 驗證安裝
```

## 📦 重新部署

### 方法 1: 使用 Git 推送自動部署（推薦）

```bash
# 1. 提交變更
git add cmd/agent/main.go Dockerfile.agent.koyeb
git commit -m "fix: 添加 HTTP 健康檢查端點，修正 Promtail 安裝"
git push origin main

# 2. Koyeb 會自動檢測並重新部署
```

### 方法 2: 使用 Koyeb CLI 手動部署

```bash
# 1. 安裝 Koyeb CLI
curl -fsSL https://cli.koyeb.com/install.sh | sh

# 2. 登錄
koyeb login

# 3. 重新部署
koyeb service redeploy pandora-agent/pandora-agent
```

## 🔍 驗證修復

部署完成後，檢查以下內容：

### 1. 查看日誌
```bash
koyeb logs pandora-agent/pandora-agent --follow
```

**預期輸出**：
```
健康檢查服務器啟動於端口 8080
已載入配置檔案: /app/configs/agent-config.yaml
INFO success: pandora-agent entered RUNNING state (不再崩潰)
```

### 2. 測試健康檢查端點

```bash
curl https://pandora-agent-<your-app-id>.koyeb.app/health
```

**預期回應**：
```json
{
  "status": "healthy",
  "service": "pandora-agent",
  "timestamp": "2025-10-07T08:30:00Z"
}
```

### 3. 檢查服務狀態

在 [Koyeb Dashboard](https://app.koyeb.com/) 確認：
- ✅ Status: **Healthy**
- ✅ Instances: **1/1 running**
- ✅ Health Checks: **Passing**

## 🎯 功能說明

修復後的 Agent 現在可以：

1. **在雲端環境運行**
   - 即使沒有實體設備也能啟動
   - 提供 HTTP API 供監控

2. **健康檢查**
   - Koyeb 可以正確檢測服務狀態
   - 避免不必要的重啟

3. **優雅關閉**
   - 正確處理 SIGTERM 信號
   - 清理資源後退出

4. **日誌收集**（如果 Promtail 正常）
   - 自動發送日誌到 Loki
   - 支援中心化日誌管理

## 📊 新的 API 端點

### GET `/health`
健康檢查端點

**回應**：
```json
{
  "status": "healthy",
  "service": "pandora-agent",
  "timestamp": "2025-10-07T08:30:00Z"
}
```

### GET `/`
服務資訊

**回應**：
```json
{
  "service": "pandora-agent",
  "version": "1.0.0"
}
```

## ⚙️ 環境變數

確保在 Koyeb 設定以下環境變數：

| 變數名 | 說明 | 預設值 |
|--------|------|--------|
| `PORT` | HTTP 服務端口 | `8080` |
| `LOG_LEVEL` | 日誌等級 | `info` |
| `GIN_MODE` | Gin 模式 | `release` |

## 🐛 已知限制

### 在雲端環境：

1. **無法控制實體設備**
   - 沒有 USB-SERIAL 連接
   - 設備管理功能無法使用

2. **主要功能受限**
   - 網路阻斷功能不可用（需要實體設備）
   - PIN 碼系統無法使用（需要 IoT 設備）

3. **適合的用途**
   - ✅ 健康檢查和監控
   - ✅ 配置測試
   - ✅ API 端點測試
   - ✅ 日誌收集

## 💡 建議

### 對於生產環境：

1. **將 Agent 部署到本地**
   - 需要實體設備控制時
   - 使用 Docker Compose 或 Systemd

2. **雲端部署用於**
   - 中心化監控
   - 日誌聚合
   - API Gateway

3. **混合架構**
   ```
   ┌────────────────┐         ┌──────────────┐
   │  Koyeb Cloud   │◄────────│  Local Agent │
   │  (監控/日誌)    │         │  (設備控制)   │
   └────────────────┘         └──────────────┘
   ```

## 📚 相關文檔

- [Koyeb 健康檢查文檔](https://www.koyeb.com/docs/build-and-deploy/health-checks)
- [Docker Compose 本地部署](../docker-compose.yml)
- [Agent 配置說明](../configs/agent-config.yaml.template)

## 🎉 完成

修復後，Pandora Agent 應該能在 Koyeb 上穩定運行，並通過健康檢查。

---

**修復日期**: 2025-10-07  
**狀態**: ✅ 已修復  
**測試**: 等待用戶驗證

