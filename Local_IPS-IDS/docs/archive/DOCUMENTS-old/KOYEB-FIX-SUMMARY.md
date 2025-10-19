# Koyeb 部署問題修復摘要

## 🐛 問題描述

部署到 Koyeb 時遇到錯誤：

```
error: failed to solve: failed to read dockerfile: 
open ./Dockerfile.agent.koyeb: no such file or directory
Build failed ❌
```

## 🔍 問題分析

### 根本原因

Koyeb 在從 GitHub 建置時，無法正確找到 `Dockerfile.agent.koyeb` 檔案。可能原因：

1. **配置檔路徑錯誤**: `.koyeb/config.yaml` 中的 dockerfile 路徑設定不正確
2. **Koyeb 配置格式**: Koyeb 可能有特定的配置檔格式要求
3. **UI 設定不完整**: 在 Koyeb Dashboard 中未明確指定 Dockerfile 路徑

## ✅ 解決方案

### 已實施的修復

1. **更新 `.koyeb/config.yaml`**
   - 從 `dockerfile: Dockerfile.agent` 改為 `dockerfile: Dockerfile.agent.koyeb`
   - 確保路徑與實際檔案名稱一致

2. **建立多種 Koyeb 配置格式**
   - `.koyeb.yml` - 簡化配置
   - `koyeb.yaml` - 官方完整配置格式
   - 提供多種選擇以應對不同情況

3. **建立詳細部署指南**
   - `KOYEB-DEPLOYMENT-GUIDE.md` - 完整故障排除指南
   - `KOYEB-QUICK-START.md` - 5 分鐘快速參考卡

4. **更新主要文件**
   - `README-PAAS-DEPLOYMENT.md` 加入 Koyeb 特別注意事項

## 📋 建議的部署步驟

### 方法 A: 使用 Koyeb Web Dashboard（最推薦）

這是最可靠的方法，因為可以在 UI 中明確指定 Dockerfile 路徑：

1. 登入 https://app.koyeb.com
2. Create App → 選擇 GitHub → 選擇 Repository
3. **關鍵步驟**: 在 Builder 設定中
   - **Builder**: Docker
   - **Dockerfile path**: `Dockerfile.agent.koyeb` （不要加 `./` 前綴）
   - **Build context**: `/`
4. 設定其他配置（Region, Instance Type, Port, Env）
5. Deploy

### 方法 B: 先建置 Docker 映像再部署

這個方法繞過 Dockerfile 路徑問題：

```bash
# 1. 本地建置映像
docker build -f Dockerfile.agent.koyeb -t YOUR_USERNAME/pandora-agent:latest .

# 2. 推送到 Docker Hub
docker push YOUR_USERNAME/pandora-agent:latest

# 3. 在 Koyeb Dashboard 中選擇 "Docker" 部署方式
# 4. 輸入映像: YOUR_USERNAME/pandora-agent:latest
```

### 方法 C: 使用標準 Dockerfile 名稱（臨時方案）

如果上述方法都不行：

```bash
# 複製為標準名稱
cp Dockerfile.agent.koyeb Dockerfile

# 在 Koyeb Dashboard 中
# Dockerfile path: Dockerfile
```

## 📁 新增的檔案

```
.
├── .koyeb/
│   └── config.yaml          (已更新)
├── .koyeb.yml              (新增 - 簡化配置)
├── koyeb.yaml              (新增 - 官方配置)
├── KOYEB-DEPLOYMENT-GUIDE.md   (新增 - 詳細指南)
├── KOYEB-QUICK-START.md        (新增 - 快速參考)
├── KOYEB-FIX-SUMMARY.md        (本檔案)
└── README-PAAS-DEPLOYMENT.md   (已更新)
```

## 🎯 關鍵配置對照

### ✅ 正確配置

| 項目 | 正確值 |
|------|--------|
| Dockerfile path | `Dockerfile.agent.koyeb` |
| Build context | `/` |
| 不要加前綴 | ❌ `./Dockerfile.agent.koyeb` |

### 📋 完整 Koyeb 配置範例

```yaml
# koyeb.yaml
app:
  name: pandora-agent

services:
  - name: pandora-agent
    build:
      type: dockerfile
      dockerfile: Dockerfile.agent.koyeb
      context: .
    
    regions:
      - fra
    
    instance_type: nano
    
    scaling:
      min: 1
      max: 2
    
    ports:
      - port: 8080
        protocol: http
    
    routes:
      - path: /
        port: 8080
    
    health_checks:
      - type: http
        port: 8080
        path: /health
        interval: 30s
```

## 🔍 驗證部署

部署成功後，執行以下驗證：

```bash
# 設定 URL（替換為實際的 Koyeb URL）
KOYEB_URL="https://pandora-agent-xxx.koyeb.app"

# 1. 健康檢查
curl $KOYEB_URL/health
# 預期: {"status":"ok"}

# 2. API 狀態
curl $KOYEB_URL/api/v1/status
# 預期: 返回系統狀態 JSON

# 3. Metrics
curl $KOYEB_URL/metrics
# 預期: Prometheus 格式的指標

# 4. 檢查日誌
koyeb service logs pandora-agent/pandora-agent --follow
```

## 📊 Koyeb 免費方案限制

- **CPU**: 0.1 vCPU per instance
- **Memory**: 512 MB per instance
- **Instances**: 2 個 Nano (永不休眠)
- **Storage**: Ephemeral (不持久化)
- **Build time**: 無限制
- **Bandwidth**: 無限制

**Agent 預估資源使用**:
- Memory: 150-200 MB
- CPU: < 5% (閒置時)
- Disk: < 100 MB (容器大小)

## 🚨 監控建議

1. **記憶體監控**
   - 設定告警: Memory > 400MB (80%)
   - 避免 OOM Kill

2. **健康檢查**
   - 確保 `/health` 端點正常回應
   - 回應時間 < 1s

3. **日誌監控**
   - 定期檢查錯誤日誌
   - Promtail 會將日誌推送到 Loki

## 📚 參考文件

- [Koyeb 官方文件](https://www.koyeb.com/docs)
- [Koyeb Docker 部署指南](https://www.koyeb.com/docs/build-and-deploy/build-from-dockerfile)
- [Koyeb 配置檔參考](https://www.koyeb.com/docs/reference/koyeb-config-file)

## 💡 下次部署提醒

1. ✅ 確認 Dockerfile 存在於根目錄
2. ✅ 在 Dashboard 中明確指定 Dockerfile 路徑
3. ✅ 使用 `/` 作為 Build context
4. ✅ 設定所有必要的環境變數
5. ✅ 健康檢查路徑為 `/health`
6. ✅ 選擇 Frankfurt (fra) region

## 🎉 預期結果

部署成功後，您將看到：

- ✅ 2 個 pandora-agent instances 運行中
- ✅ 健康檢查通過
- ✅ `/health` 端點回應正常
- ✅ 日誌顯示 Agent 和 Promtail 都在運行
- ✅ Metrics 可以被 Prometheus 抓取

---

**修復日期**: 2024-12-19
**狀態**: ✅ 已提供多種解決方案
**建議**: 優先使用方法 A (Web Dashboard)

