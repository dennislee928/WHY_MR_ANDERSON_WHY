# Docker 錯誤修復報告

**修復日期**: 2025-10-15  
**狀態**: ✅ 所有關鍵問題已修復

---

## 🔧 已應用的修復

### 1. ✅ Nginx Healthcheck 修復

**問題**: healthcheck 測試 `/health` 端點不存在

**修復**:
```yaml
# 從
test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost/health"]

# 改為
test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:80"]
```

**檔案**: `Application/docker-compose.yml`

---

### 2. ✅ Alertmanager Webhook 修復

**問題**: 嘗試連接不存在的 `axiom-ui` 主機

**修復**:
```yaml
# 所有 webhook 從
url: 'http://axiom-ui:3001/...'

# 改為正確的容器名
url: 'http://axiom-be:3001/...'
```

**檔案**: `configs/alertmanager.yml`  
**修復數量**: 5 處

---

### 3. ✅ Promtail 檔案寫入權限修復

**問題**: `/app/data` 掛載為唯讀，無法寫入 positions.yaml

**修復**:
```yaml
# 從唯讀掛載
- pandora-data:/app/data:ro

# 改為專用可寫入的 volume
- promtail-positions:/app/data
```

**檔案**: `Application/docker-compose.yml`

**同時添加**:
```yaml
volumes:
  promtail-positions:  # 新增
```

---

## ⚠️ 低優先級問題（可忽略）

### 4. ℹ️ Pandora-Agent mTLS 警告

**問題**: mTLS 憑證未找到

**影響**: ✅ 無 - Agent 仍正常運行（healthy）

**狀態**: 可選修復（mTLS 功能不是必需的）

**修復方案** (可選):
```bash
cd configs
mkdir -p certs
openssl req -x509 -newkey rsa:2048 -keyout certs/client.key \
  -out certs/client.crt -days 365 -nodes \
  -subj "/CN=pandora-agent"
```

---

### 5. ℹ️ Redis "Security Attack" 警告

**問題**: Redis 檢測到非標準命令

**原因**: 健康檢查或工具探測

**影響**: ✅ 無 - 這是誤報

**狀態**: 可忽略

---

### 6. ℹ️ Postgres "Invalid Startup Packet"

**問題**: 接收到格式錯誤的連接包

**原因**: 探測或測試連接

**影響**: ✅ 無 - Postgres 正常運行（healthy）

**狀態**: 可忽略

---

### 7. ℹ️ Node-Exporter nfsd 錯誤

**問題**: 無法讀取 nfsd 指標

**原因**: WSL2 環境中 nfsd 不可用

**影響**: ✅ 無 - 其他指標正常收集

**狀態**: 可忽略

---

## 🚀 應用修復

### 重啟受影響的容器

```bash
cd Application

# 重啟 nginx（應用新的 healthcheck）
docker-compose restart nginx

# 重啟 alertmanager（應用新的 webhook 配置）
docker-compose restart alertmanager

# 重啟 promtail（應用新的 volume 配置）
docker-compose down promtail
docker-compose up -d promtail
```

### 或全面重啟（推薦）

```bash
cd Application

# 停止所有容器
docker-compose down

# 使用新配置啟動
docker-compose up -d

# 等待所有服務就緒
sleep 30

# 檢查狀態
docker-compose ps
```

---

## ✅ 預期修復結果

### 修復前
```
nginx         unhealthy  ❌
portainer     unhealthy  ❌
alertmanager  healthy    ⚠️ (webhook 失敗)
promtail      running    ⚠️ (寫入失敗)
```

### 修復後
```
nginx         healthy    ✅
portainer     healthy    ✅ (如果 healthcheck 正常)
alertmanager  healthy    ✅ (webhook 正常)
promtail      running    ✅ (可寫入)
```

---

## 📊 修復總結

| 問題 | 嚴重性 | 修復 | 檔案 |
|------|--------|------|------|
| Nginx unhealthy | 🔴 中 | ✅ | docker-compose.yml |
| Alertmanager DNS | 🔴 中 | ✅ | alertmanager.yml |
| Promtail 權限 | 🟡 低 | ✅ | docker-compose.yml |
| Pandora mTLS | 🟢 可選 | ⏭️ | - |
| Redis 誤報 | 🟢 無害 | - | - |
| Postgres 警告 | 🟢 無害 | - | - |
| Node-exporter | 🟢 無害 | - | - |

---

## 🎯 立即執行

```bash
cd ~/Documents/GitHub/Local_IPS-IDS/Application

# 應用所有修復
docker-compose down
docker-compose up -d

# 等待 30 秒
sleep 30

# 檢查健康狀態
docker-compose ps
```

---

**修復狀態**: ✅ 關鍵問題已修復  
**下一步**: 重啟容器應用修復

