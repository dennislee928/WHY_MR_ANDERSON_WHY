# Koyeb 部署指南 - Pandora Agent

## 問題診斷

如果遇到 `error: failed to solve: failed to read dockerfile: open ./Dockerfile.agent.koyeb: no such file or directory` 錯誤，請按照以下步驟操作。

## 解決方案

### 方案 1: 使用 Koyeb Web Dashboard (推薦)

1. **登入 Koyeb Dashboard**
   - 前往 https://app.koyeb.com
   - 登入您的帳號

2. **建立新應用**
   - 點擊 "Create App"
   - 選擇 "GitHub" 作為部署來源
   - 選擇 Repository: `pandora_box_console_IDS-IPS`
   - 選擇 Branch: `main`

3. **設定建置配置**
   - **Builder**: Docker
   - **Dockerfile path**: `Dockerfile.agent.koyeb`（重要！必須明確指定）
   - **Build context**: `/` (root)

4. **設定部署配置**
   - **App name**: `pandora-agent`
   - **Service name**: `pandora-agent`
   - **Region**: Europe (Frankfurt) - `fra`
   - **Instance type**: Nano (免費)
   - **Instances**: 2

5. **設定端口**
   - **Port**: `8080`
   - **Protocol**: HTTP
   - **Public**: Yes

6. **設定環境變數**
   ```
   LOG_LEVEL=info
   GIN_MODE=release
   PORT=8080
   DATABASE_URL=<your-railway-database-url>
   REDIS_URL=<your-render-redis-url>
   PROMETHEUS_URL=<your-prometheus-url>
   LOKI_URL=<your-loki-url>
   GRAFANA_URL=<your-grafana-url>
   JWT_SECRET=<your-jwt-secret>
   ENCRYPTION_KEY=<your-encryption-key>
   ```

7. **設定健康檢查**
   - **Type**: HTTP
   - **Path**: `/health`
   - **Port**: 8080
   - **Initial delay**: 30s
   - **Interval**: 30s
   - **Timeout**: 10s

8. **部署**
   - 點擊 "Deploy"
   - 等待建置完成

### 方案 2: 使用 Koyeb CLI

```bash
# 1. 安裝 Koyeb CLI
curl -fsSL https://cli.koyeb.com/install.sh | bash

# 2. 登入
koyeb login

# 3. 建立應用
koyeb app create pandora-agent

# 4. 建置並推送 Docker 映像到 Docker Hub
docker build -f Dockerfile.agent.koyeb -t YOUR_DOCKERHUB_USERNAME/pandora-agent:latest .
docker push YOUR_DOCKERHUB_USERNAME/pandora-agent:latest

# 5. 部署服務
koyeb service create pandora-agent \
  --app pandora-agent \
  --docker YOUR_DOCKERHUB_USERNAME/pandora-agent:latest \
  --ports 8080:http \
  --routes /:8080 \
  --regions fra \
  --instance-type nano \
  --scale 2 \
  --env LOG_LEVEL=info \
  --env GIN_MODE=release \
  --env PORT=8080 \
  --env DATABASE_URL=$RAILWAY_DATABASE_URL \
  --env REDIS_URL=$RENDER_REDIS_URL \
  --env PROMETHEUS_URL=$PROMETHEUS_URL \
  --env LOKI_URL=$LOKI_URL \
  --env GRAFANA_URL=$GRAFANA_URL \
  --env JWT_SECRET=$JWT_SECRET \
  --env ENCRYPTION_KEY=$ENCRYPTION_KEY
```

### 方案 3: 重新命名 Dockerfile（臨時解決方案）

如果 Koyeb 無法識別 `Dockerfile.agent.koyeb`，可以暫時重新命名：

```bash
# 備份原始檔案
cp Dockerfile.agent.koyeb Dockerfile.agent.koyeb.backup

# 複製為標準名稱
cp Dockerfile.agent.koyeb Dockerfile

# 在 Koyeb Dashboard 中使用 "Dockerfile" 作為路徑
```

然後在 Koyeb Dashboard 中：
- **Dockerfile path**: `Dockerfile`

## 常見問題排除

### 問題 1: Dockerfile 路徑錯誤

**錯誤訊息**:
```
error: failed to solve: failed to read dockerfile: open ./Dockerfile.agent.koyeb: no such file or directory
```

**解決方案**:
- 確認 Repository 根目錄有 `Dockerfile.agent.koyeb` 檔案
- 在 Koyeb Dashboard 的 "Dockerfile path" 欄位中明確指定路徑
- 確認路徑不包含 `./` 前綴，直接使用 `Dockerfile.agent.koyeb`

### 問題 2: 建置超時

**解決方案**:
- 使用方案 2 (Docker Hub) 先建置映像
- 或使用 GitHub Actions 自動建置並推送

### 問題 3: 環境變數未生效

**解決方案**:
- 確認在 Koyeb Dashboard 的 "Environment" 標籤中正確設定
- 敏感資訊使用 "Secret" 類型
- 重新部署服務以套用變更

## Dockerfile 說明

`Dockerfile.agent.koyeb` 特點：

1. **多階段建置**: 減少最終映像大小
2. **Supervisor 整合**: 同時運行 pandora-agent 和 promtail
3. **健康檢查**: 內建 `/health` 端點
4. **安全性**: 使用非 root 使用者運行

## 驗證部署

部署完成後，驗證服務：

```bash
# 檢查健康狀態
curl https://pandora-agent-xxx.koyeb.app/health

# 檢查 API 狀態
curl https://pandora-agent-xxx.koyeb.app/api/v1/status

# 檢查 metrics
curl https://pandora-agent-xxx.koyeb.app/metrics
```

## 監控與日誌

在 Koyeb Dashboard 中：

1. **Logs**: 查看即時日誌
2. **Metrics**: 監控 CPU, Memory, Network
3. **Events**: 查看部署歷史
4. **Settings**: 調整配置

## 成本估算

使用 Koyeb 免費方案：

- **2 個 Nano 實例**: 永久免費
- **Bandwidth**: 100GB/月
- **Build time**: 無限制
- **永不休眠**: ✅

## 支援

如有問題：

1. 查看 Koyeb 文件: https://www.koyeb.com/docs
2. Koyeb Discord: https://discord.gg/koyeb
3. GitHub Issues: https://github.com/your-org/pandora_box_console_IDS-IPS/issues

---

**提示**: 建議先使用方案 1 (Web Dashboard)，因為它最直觀且錯誤訊息最清晰。

