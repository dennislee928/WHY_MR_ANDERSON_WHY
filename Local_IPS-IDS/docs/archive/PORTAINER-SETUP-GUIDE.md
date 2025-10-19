# Portainer 設置與使用指南

## 📋 概述

Portainer 是一個輕量級的容器管理平台，提供直觀的 Web UI 來管理 Docker 容器、映像、卷和網路。

在 Pandora Box Console 系統中，Portainer 作為**主要的容器管理工具**，讓您可以：
- 📦 集中查看所有 14 個容器的狀態
- 📋 統一查看和搜索日誌
- 📊 監控資源使用情況
- 🔧 快速執行管理操作

---

## 🚀 快速開始

### 1. 啟動 Portainer

Portainer 已包含在 `docker-compose.yml` 中，隨系統一起啟動：

```bash
cd Application
./docker-start.sh
```

或單獨啟動 Portainer：

```bash
docker-compose up -d portainer
```

### 2. 初次設置

1. **訪問 Portainer**：
   - HTTP: http://localhost:9000
   - HTTPS: https://localhost:9443

2. **創建管理員帳號**：
   ```
   用戶名: admin
   密碼: (至少12個字元，建議使用 pandora123456)
   ```

3. **選擇環境**：
   - 選擇 "Get Started"
   - Portainer 會自動連接到本地 Docker 環境

4. **完成**：
   - 您現在可以看到 Dashboard

---

## 🎯 核心功能使用指南

### 功能 1: 容器管理

#### 查看所有容器

1. 點擊左側菜單 **"Containers"**
2. 查看所有 14 個 Pandora 容器：

| 容器名稱 | 狀態 | 用途 |
|---------|------|------|
| `portainer` | 運行中 | 容器管理平台 |
| `axiom-be` | 運行中 | 後端 API 服務 |
| `pandora-agent` | 運行中 | 核心 Agent |
| `cyber-ai-quantum` | 運行中 | AI/量子服務 |
| `prometheus` | 運行中 | 指標收集 |
| `grafana` | 運行中 | 監控儀表板 |
| `loki` | 運行中 | 日誌聚合 |
| `alertmanager` | 運行中 | 告警管理 |
| `rabbitmq` | 運行中 | 消息隊列 |
| `postgres` | 運行中 | 資料庫 |
| `redis` | 運行中 | 快取系統 |
| `node-exporter` | 運行中 | 系統指標 |
| `promtail` | 運行中 | 日誌收集 |
| `nginx` | 運行中 | 反向代理 |

#### 快速操作

- **啟動容器**: 選擇容器 → 點擊 "Start"
- **停止容器**: 選擇容器 → 點擊 "Stop"
- **重啟容器**: 選擇容器 → 點擊 "Restart"
- **刪除容器**: 選擇容器 → 點擊 "Remove"

#### 批量操作

1. 勾選多個容器
2. 點擊頂部的批量操作按鈕
3. 選擇操作（Start/Stop/Restart/Remove）

---

### 功能 2: 日誌查看

#### 查看單個容器日誌

1. 點擊 **容器名稱**（如 `cyber-ai-quantum`）
2. 點擊 **"Logs"** 標籤
3. 功能：
   - **即時串流**: 自動更新最新日誌
   - **搜索**: 輸入關鍵字（如 "error", "quantum"）
   - **時間戳**: 顯示完整時間
   - **複製**: 選擇並複製日誌
   - **下載**: 點擊 "Download" 下載完整日誌

#### 常用日誌搜索

**查找錯誤**：
```
搜索: error
結果: 所有包含 "error" 的日誌行
```

**查找量子作業**：
```
容器: cyber-ai-quantum
搜索: quantum job
```

**查找告警**：
```
容器: alertmanager
搜索: Notify for alerts failed
```

---

### 功能 3: 資源監控

#### 查看容器資源使用

1. 點擊容器名稱
2. 點擊 **"Stats"** 標籤
3. 查看：
   - **CPU 使用率** (%)
   - **記憶體使用** (MB / GB)
   - **網路 I/O** (上傳/下載)
   - **磁碟 I/O** (讀/寫)

#### Dashboard 總覽

在主 Dashboard 可以看到：
- 總容器數量：14
- 運行中：14
- 停止：0
- 映像總數
- 卷總數
- 網路總數

---

### 功能 4: 容器終端訪問

#### 進入容器 Shell

1. 點擊容器名稱
2. 點擊 **"Console"** 標籤
3. 選擇 Shell：
   - `/bin/sh` (Alpine 容器)
   - `/bin/bash` (Debian/Ubuntu 容器)
4. 點擊 "Connect"

#### 常用診斷命令

**檢查 Cyber AI/Quantum 容器**：
```bash
# 進入容器
# 選擇 cyber-ai-quantum → Console → /bin/sh

# 查看 Python 進程
ps aux | grep python

# 測試 API
curl http://localhost:8000/health

# 查看環境變數
env | grep QUANTUM

# 查看已安裝的包
pip list | grep qiskit
```

**檢查 Axiom BE 容器**：
```bash
# 進入容器
# 選擇 axiom-be → Console → /bin/sh

# 查看進程
ps aux

# 測試 API
curl http://localhost:3001/api/v1/health

# 查看配置
cat /app/configs/ui-config.yaml
```

---

### 功能 5: 映像管理

#### 查看映像

1. 點擊左側 **"Images"**
2. 查看所有 Docker 映像：
   - `application-axiom-be`
   - `application-cyber-ai-quantum`
   - `prom/prometheus`
   - `grafana/grafana`
   - 等等

#### 映像操作

- **拉取新映像**: Import → Pull Image
- **刪除映像**: 選擇 → Remove
- **查看層級**: 點擊映像名稱 → Layers

---

### 功能 6: 卷管理

#### 查看所有卷

1. 點擊左側 **"Volumes"**
2. 查看 Pandora 系統的卷：
   - `prometheus-data` (30.5 MB)
   - `grafana-data` (12.3 MB)
   - `postgres-data` (145 MB)
   - `redis-data` (2.1 MB)
   - `rabbitmq-data` (8.5 MB)
   - `ai-models` (500 MB)
   - `portainer-data`
   - 等等

#### 卷操作

- **瀏覽內容**: 點擊卷名稱 → Browse
- **備份卷**: 點擊 → 複製卷名 → 手動 `docker cp`
- **刪除卷**: 選擇 → Remove (⚠️ 謹慎操作)

---

### 功能 7: Stack 管理

#### 查看 Docker Compose Stack

1. 點擊左側 **"Stacks"**
2. 查看 `application` stack
3. 功能：
   - 查看所有服務
   - 編輯 docker-compose.yml
   - 重新部署
   - 停止整個 stack
   - 刪除 stack

#### 更新 Stack

1. 點擊 stack 名稱
2. 點擊 **"Editor"**
3. 編輯 `docker-compose.yml`
4. 點擊 **"Update the stack"**

---

## 🔍 故障排除場景

### 場景 1: Cyber AI/Quantum 容器無法啟動

**步驟**：

1. Portainer → Containers → 找到 `cyber-ai-quantum`
2. 查看狀態（可能是 "Exited" 或 "Restarting"）
3. 點擊容器名稱 → **Logs**
4. 搜索 "error" 或 "failed"
5. 找到錯誤原因：
   ```
   ModuleNotFoundError: No module named 'qiskit'
   ```
6. 解決方案：重新構建映像
   ```bash
   docker-compose build --no-cache cyber-ai-quantum
   docker-compose up -d cyber-ai-quantum
   ```

### 場景 2: AlertManager Webhook 404

**步驟**：

1. Portainer → Containers → `alertmanager` → Logs
2. 搜索 "404"
3. 發現：
   ```
   unexpected status code 404: http://axiom-be:3001/api/v1/alerts/webhook
   ```
4. 原因：端點未實現
5. 解決：參考 `docs/ERROR-ANALYSIS-AND-SOLUTIONS.md`

### 場景 3: 資源耗盡

**步驟**：

1. Portainer → Dashboard
2. 查看總資源使用
3. Containers → 排序 by CPU 或 Memory
4. 找到資源消耗最高的容器
5. 點擊容器 → Stats → 查看詳細趨勢
6. 決定是否需要擴展資源或優化

---

## 📊 Portainer 監控指標

### Dashboard 關鍵指標

```
┌─────────────────────────────────────────────┐
│         Portainer Dashboard                 │
├─────────────────────────────────────────────┤
│ Stacks:     1                               │
│ Containers: 14 (14 running, 0 stopped)      │
│ Images:     12                              │
│ Volumes:    14                              │
│ Networks:   1                               │
└─────────────────────────────────────────────┘

容器健康狀態：
✅ portainer          - healthy
✅ axiom-be           - healthy
✅ pandora-agent      - healthy
✅ cyber-ai-quantum   - healthy
✅ prometheus         - healthy
✅ grafana            - healthy
✅ loki               - healthy
✅ alertmanager       - healthy
✅ rabbitmq           - healthy
✅ postgres           - healthy
✅ redis              - healthy
✅ node-exporter      - up
✅ promtail           - healthy
✅ nginx              - healthy
```

---

## 🛠️ 進階配置

### 1. 啟用認證

預設情況下，Portainer 需要登入。建議設置：

```
用戶名: admin
密碼: pandora_portainer_2025!
```

### 2. 配置 Teams 和 RBAC

如果多人使用，可以配置：
- 創建團隊
- 設置角色權限
- 限制操作範圍

### 3. Webhook 通知

設置 Portainer webhook 來接收容器事件：

1. Settings → Notifications
2. 添加 Webhook URL
3. 選擇事件類型（Container stopped, Image removed）

---

## 🎯 Portainer vs 其他工具

| 功能 | Portainer | k9s | Docker Desktop | Grafana |
|------|-----------|-----|----------------|---------|
| **容器管理** | ✅ Web UI | ✅ CLI | ✅ GUI | ❌ |
| **日誌查看** | ✅ 即時 | ✅ 即時 | ✅ 基礎 | ⚠️ 需 Loki |
| **資源監控** | ✅ 圖表 | ✅ 即時 | ✅ 基礎 | ✅ 專業 |
| **終端訪問** | ✅ Web | ✅ CLI | ✅ GUI | ❌ |
| **學習曲線** | ⭐ 簡單 | ⭐⭐ 中等 | ⭐ 簡單 | ⭐⭐⭐ 複雜 |
| **適用場景** | Docker Compose | Kubernetes | 本地開發 | 生產監控 |

**推薦**：
- **Docker Compose 環境**: Portainer (最佳選擇)
- **Kubernetes 環境**: k9s
- **本地開發**: Docker Desktop
- **生產監控**: Grafana + Prometheus

---

## 📚 常用操作速查

### 查看容器日誌
```
Containers → 點擊容器 → Logs → 搜索過濾
```

### 重啟服務
```
Containers → 選擇容器 → Restart
```

### 查看資源使用
```
Containers → 點擊容器 → Stats
```

### 進入容器終端
```
Containers → 點擊容器 → Console → Connect
```

### 查看容器配置
```
Containers → 點擊容器 → Inspect
```

### 更新容器映像
```
Containers → 點擊容器 → Duplicate/Edit → Pull latest image
```

### 管理 Stack
```
Stacks → application → Editor → Update
```

### 清理資源
```
Images → Unused images → Remove
Volumes → Unused volumes → Remove
```

---

## 🔒 安全建議

### 1. 使用 HTTPS

預設 Portainer 同時提供 HTTP (9000) 和 HTTPS (9443)。

**生產環境建議**：
- 僅使用 HTTPS (9443)
- 關閉 HTTP 端口

### 2. 強密碼

初次設置時使用強密碼（至少 12 字元）：
```
pandora_portainer_2025!
```

### 3. 限制訪問

如果暴露到網路，使用防火牆限制：
```bash
# 僅允許本地訪問
sudo ufw allow from 127.0.0.1 to any port 9000
sudo ufw allow from 127.0.0.1 to any port 9443
```

### 4. 定期更新

保持 Portainer 更新到最新版本：
```bash
docker-compose pull portainer
docker-compose up -d portainer
```

---

## 📈 整合到 Pandora 監控體系

### 三層監控架構

```
┌─────────────────────────────────────────────┐
│  Layer 1: 容器層級管理 (Portainer)          │
│  • 容器狀態、日誌、終端                      │
│  • 快速故障排除                              │
│  • 資源即時監控                              │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────┴──────────────────────────┐
│  Layer 2: 應用層級監控 (Grafana)            │
│  • 業務指標、性能指標                        │
│  • 歷史趨勢分析                              │
│  • 告警規則管理                              │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────┴──────────────────────────┐
│  Layer 3: 數據收集 (Prometheus/Loki)        │
│  • 指標持久化                                │
│  • 日誌聚合                                  │
│  • 長期儲存                                  │
└─────────────────────────────────────────────┘
```

### 推薦監控工作流程

1. **日常檢查** (每天):
   - 訪問 Portainer Dashboard
   - 確認所有容器運行中
   - 快速掃描 CPU/Memory 使用率

2. **詳細分析** (每週):
   - 訪問 Grafana
   - 查看業務指標趨勢
   - 檢查告警歷史

3. **故障排除** (發生問題時):
   - Portainer → 找到異常容器
   - 查看即時日誌
   - 必要時進入終端診斷
   - 參考 Grafana 歷史數據

4. **性能優化** (每月):
   - Portainer → Stats → 找到資源瓶頸
   - Grafana → 長期趨勢分析
   - 決定擴展或優化策略

---

## 🎓 Portainer 最佳實踐

### 1. 定期備份

**備份 Portainer 數據**：
```bash
docker run --rm \
  -v portainer-data:/data \
  -v $(pwd):/backup \
  alpine tar czf /backup/portainer-backup-$(date +%Y%m%d).tar.gz /data
```

**恢復備份**：
```bash
docker run --rm \
  -v portainer-data:/data \
  -v $(pwd):/backup \
  alpine sh -c "cd /data && tar xzf /backup/portainer-backup-20250114.tar.gz --strip 1"
```

### 2. 使用標籤

為容器添加標籤，方便分類：
```yaml
# docker-compose.yml
labels:
  - "category=monitoring"
  - "tier=frontend"
  - "criticality=high"
```

### 3. 設置告警

在 Portainer 中設置基本告警：
- 容器停止 → 發送通知
- 資源使用超過閾值 → 告警

---

## 📝 故障排除

### 問題 1: 無法訪問 Portainer

**解決方案**：
```bash
# 檢查容器狀態
docker ps | grep portainer

# 查看日誌
docker logs portainer

# 重啟容器
docker-compose restart portainer
```

### 問題 2: 忘記密碼

**解決方案**：
```bash
# 重置管理員密碼
docker stop portainer
docker run --rm -v portainer-data:/data alpine sh -c "rm -rf /data/portainer.db"
docker-compose up -d portainer
# 重新訪問並設置新密碼
```

### 問題 3: 容器列表不完整

**解決方案**：
- 確認 Docker socket 正確掛載
- 檢查 `/var/run/docker.sock` 權限
- 重啟 Portainer

---

## 🎉 快速開始檢查清單

- [ ] 訪問 http://localhost:9000
- [ ] 創建管理員帳號
- [ ] 選擇 "Get Started"
- [ ] 確認看到 14 個容器
- [ ] 點擊 `cyber-ai-quantum` 查看日誌
- [ ] 搜索 "quantum" 關鍵字
- [ ] 查看 `axiom-be` 的 Stats
- [ ] 進入 `redis` 容器終端
- [ ] 執行 `redis-cli ping`
- [ ] 查看 Images 列表
- [ ] 查看 Volumes 使用情況

---

## 📖 相關文檔

- 📚 [Portainer 官方文檔](https://docs.portainer.io/)
- 📚 [Quick-Start.md](../Quick-Start.md) - 系統快速啟動
- 📚 [ERROR-ANALYSIS-AND-SOLUTIONS.md](ERROR-ANALYSIS-AND-SOLUTIONS.md) - 錯誤分析
- 📚 [README.md](../README.md) - 專案主文檔

---

**版本**: Portainer CE 2.19.4  
**維護者**: Pandora Security Team  
**最後更新**: 2025-01-14  
**用途**: Pandora Box Console 容器集中管理

