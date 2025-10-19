# Docker 容器錯誤分析與修復

**分析日期**: 2025-10-15  
**容器數量**: 14 個  
**問題數量**: 6 個

---

## 📊 容器健康狀態總覽

| 容器 | 狀態 | 問題 |
|------|------|------|
| cyber-ai-quantum | ✅ healthy | 無 |
| axiom-be | ✅ healthy | 無 |
| grafana | ✅ healthy | 無 |
| loki | ✅ healthy | 無 |
| prometheus | ✅ healthy | 無 |
| alertmanager | ✅ healthy | ⚠️ 告警失敗 |
| postgres | ✅ healthy | ⚠️ 連接警告 |
| rabbitmq | ✅ healthy | 無 |
| redis | ✅ healthy | ⚠️ 誤報 |
| **nginx** | ❌ **unhealthy** | ⚠️ healthcheck 失敗 |
| **portainer** | ❌ **unhealthy** | ⚠️ healthcheck 失敗 |
| pandora-agent | ✅ healthy | ⚠️ mTLS 警告 |
| promtail | ⚠️ running | ⚠️ 檔案寫入錯誤 |
| node-exporter | ⚠️ running | ⚠️ nfsd 收集錯誤 |

**健康狀態**: 10/14 healthy (71%)

---

## 🔍 詳細錯誤分析

### 1. ❌ Nginx - Unhealthy

**錯誤**:
```
nginx: Configuration complete; ready for start up
healthcheck 失敗
```

**原因**:
- Healthcheck 測試 `http://localhost/health` 端點
- 該端點可能不存在或配置錯誤

**影響**: ⚠️ 中等 - Nginx 仍在運行，但健康檢查失敗

**修復方案**:

#### 方案 A: 修改 healthcheck 端點

```yaml
# docker-compose.yml
nginx:
  healthcheck:
    test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:80"]
    # 或測試根路徑
    # test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost/"]
```

#### 方案 B: 添加 /health 端點到 nginx 配置

```nginx
# configs/nginx/default-paas.conf
location /health {
    access_log off;
    return 200 "healthy\n";
    add_header Content-Type text/plain;
}
```

---

### 2. ⚠️ Pandora-Agent - mTLS 錯誤

**錯誤**:
```
level=error msg="初始化代理程式失敗: mTLS客戶端初始化失敗: 
載入客戶端憑證失敗: open : no such file or directory"
```

**原因**:
- mTLS 憑證路徑為空或不存在
- `MTLS_CERT_PATH` 和 `MTLS_KEY_PATH` 配置問題

**影響**: ⚠️ 低 - Agent 仍在運行（healthy），但 mTLS 功能不可用

**修復方案**:

```yaml
# docker-compose.yml
pandora-agent:
  environment:
    - MTLS_CERT_PATH=/certs/client.crt
    - MTLS_KEY_PATH=/certs/client.key
  volumes:
    - ./certs:/certs:ro  # 確保憑證目錄存在
```

**快速修復**:
```bash
# 生成測試憑證
cd configs
mkdir -p certs
openssl req -x509 -newkey rsa:2048 -keyout certs/client.key -out certs/client.crt -days 365 -nodes -subj "/CN=pandora-agent"
```

---

### 3. ⚠️ Promtail - Read-only File System

**錯誤**:
```
level=error msg="error writing positions file" 
error="open /app/data/.positions.yaml: read-only file system"
```

**原因**:
- `/app/data` 目錄掛載為唯讀
- Promtail 需要寫入 positions 檔案來追蹤日誌位置

**影響**: ⚠️ 低 - Promtail 仍可讀取日誌，但無法保存進度

**修復方案**:

```yaml
# docker-compose.yml
promtail:
  volumes:
    - promtail-data:/app/data  # 使用 volume 而不是唯讀掛載
    - /var/log:/var/log:ro
```

---

### 4. ⚠️ Alertmanager - Cannot Resolve axiom-ui

**錯誤**:
```
msg="Notify attempt failed" 
err="dial tcp: lookup axiom-ui on 127.0.0.11:53: no such host"
```

**原因**:
- Alertmanager 配置中使用 `axiom-ui` 主機名
- 但容器名稱是 `axiom-be`

**影響**: ⚠️ 中等 - 告警通知失敗，但不影響監控

**修復方案**:

```yaml
# configs/alertmanager.yml
receivers:
  - name: 'default-receiver'
    webhook_configs:
      - url: 'http://axiom-be:3001/api/v1/alerts'  # 改為正確的容器名
```

---

### 5. ⚠️ Redis - Security Attack Detected

**錯誤**:
```
Possible SECURITY ATTACK detected. It looks like somebody is sending 
POST or Host: commands to Redis.
```

**原因**:
- 這是 **誤報**
- 來自健康檢查或其他服務的連接嘗試
- Redis 檢測到非標準命令

**影響**: ✅ 無 - Redis 仍正常運行

**解決方案**: 
- ✅ 可以忽略（正常的連接測試）
- 或在 Redis 配置中禁用此警告

---

### 6. ⚠️ Postgres - Invalid Startup Packet

**錯誤**:
```
LOG: invalid length of startup packet
```

**原因**:
- 健康檢查或其他服務發送了格式錯誤的連接包
- 通常是工具發送的探測包

**影響**: ✅ 無 - Postgres 正常運行（healthy）

**解決方案**: ✅ 可以忽略

---

### 7. ❌ Portainer - Unhealthy

**狀態**: Unhealthy

**原因**: 類似 nginx，healthcheck 可能配置不當

**影響**: ⚠️ 低 - Web UI 可能仍可訪問

**測試**:
```bash
# 訪問 Portainer
http://localhost:9000
```

---

## 🔧 優先修復順序

### 🔴 高優先級（影響功能）

#### 1. Nginx Unhealthy

<function_calls>
<invoke name="search_replace">
<parameter name="file_path">Application/docker-compose.yml
