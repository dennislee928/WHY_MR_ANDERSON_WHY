# Nginx 在 Pandora Box Console 的角色與功能

**版本**: v3.4.1  
**配置**: `configs/nginx/default-paas.conf`

---

## 🎯 核心角色

### 1. **統一入口點（API Gateway）**

Nginx 作為整個 Pandora Box Console 系統的**唯一對外入口**，統一管理所有服務的訪問。

```
外部請求 → Nginx (Port 80/443) → 內部微服務
```

---

## 🔧 主要功能

### 1. ✅ 反向代理（Reverse Proxy）

**功能**: 將外部請求路由到對應的內部服務

#### 配置的上游服務 (5 個)

```nginx
upstream axiom_ui {
    server axiom-be:3001;       # API 後端服務
    keepalive 32;
}

upstream grafana {
    server grafana:3000;        # 視覺化儀表板
    keepalive 32;
}

upstream prometheus {
    server prometheus:9090;     # 監控指標
    keepalive 32;
}

upstream loki {
    server loki:3100;           # 日誌聚合
    keepalive 32;
}

upstream alertmanager {
    server alertmanager:9093;   # 告警管理
    keepalive 32;
}
```

---

### 2. ✅ 路由管理（URL Routing）

**功能**: 根據 URL 路徑分發到不同服務

| URL 路徑 | 轉發目標 | 用途 |
|---------|---------|------|
| `/api/` | axiom-be:3001 | RESTful API |
| `/ws` | axiom-be:3001 | WebSocket 即時通訊 |
| `/grafana/` | grafana:3000 | Grafana 儀表板 |
| `/prometheus/` | prometheus:9090 | Prometheus 監控 |
| `/loki/` | loki:3100 | Loki 日誌查詢 |
| `/alertmanager/` | alertmanager:9093 | AlertManager 管理 |
| `/health` | 內建端點 | 健康檢查 |
| `/` | 靜態檔案 | 前端 UI |

#### 範例

```bash
# 外部訪問
curl http://localhost/api/v1/health

# Nginx 轉發到
→ http://axiom-be:3001/api/v1/health
```

---

### 3. ✅ WebSocket 支援

**功能**: 支援長連接的即時通訊

```nginx
location /ws {
    proxy_pass http://axiom_ui;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    
    # WebSocket 超時設定（7天）
    proxy_connect_timeout 7d;
    proxy_send_timeout 7d;
    proxy_read_timeout 7d;
}
```

**用途**:
- 即時監控數據推送
- 告警即時通知
- 雙向通訊

---

### 4. ✅ 安全防護

#### A. 安全標頭（Security Headers）

```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "no-referrer-when-downgrade" always;
```

**防護效果**:
- ✅ **X-Frame-Options**: 防止 Clickjacking 攻擊
- ✅ **X-Content-Type-Options**: 防止 MIME 類型混淆
- ✅ **X-XSS-Protection**: 防止跨站腳本攻擊
- ✅ **Referrer-Policy**: 控制 Referer 資訊洩露

#### B. 連接資訊轉發

```nginx
proxy_set_header X-Real-IP $remote_addr;
proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
proxy_set_header X-Forwarded-Proto $scheme;
```

**功能**:
- 保留真實客戶端 IP
- 支援 IP 白名單/黑名單
- 審計日誌追蹤

---

### 5. ✅ 效能優化

#### A. Gzip 壓縮

```nginx
gzip  on;
gzip_comp_level 6;
gzip_types text/plain text/css application/json ...
```

**效果**:
- 減少傳輸量 60-80%
- 加快頁面載入速度
- 節省頻寬成本

#### B. 連接池（Keep-Alive）

```nginx
upstream axiom_ui {
    server axiom-be:3001;
    keepalive 32;  # 保持 32 個持久連接
}
```

**效果**:
- 減少 TCP 握手次數
- 降低延遲
- 提升吞吐量

#### C. 靜態檔案快取

```nginx
location / {
    expires 1h;
    add_header Cache-Control "public, immutable";
}
```

**效果**:
- 瀏覽器快取 1 小時
- 減少重複請求
- 加快頁面載入

---

### 6. ✅ 負載均衡（預留）

**當前**: 單一後端

**擴展性**: 可輕鬆添加多個後端

```nginx
# 未來可擴展為
upstream axiom_ui {
    server axiom-be-1:3001 weight=5;
    server axiom-be-2:3001 weight=3;
    server axiom-be-3:3001 backup;
    keepalive 32;
}
```

---

### 7. ✅ 超時管理

**API 請求超時**:
```nginx
proxy_connect_timeout 60s;
proxy_send_timeout 60s;
proxy_read_timeout 60s;
```

**WebSocket 超時**:
```nginx
proxy_connect_timeout 7d;  # 7 天長連接
```

**保護效果**:
- 防止慢速攻擊
- 避免資源耗盡
- 優雅處理超時

---

### 8. ✅ 健康檢查端點

```nginx
location /health {
    access_log off;        # 不記錄健康檢查日誌
    return 200 "OK\n";     # 直接返回 200
    add_header Content-Type text/plain;
}
```

**用途**:
- Docker healthcheck
- 負載均衡器檢查
- 監控系統探測

---

### 9. ✅ 靜態檔案服務

```nginx
location / {
    root /usr/share/nginx/html;
    index index.html;
    try_files $uri $uri/ /index.html;  # SPA 路由支援
}
```

**功能**:
- 提供前端 UI 檔案
- 支援 React/Vue SPA 路由
- 靜態資源快取

---

### 10. ✅ 日誌記錄

```nginx
access_log  /var/log/nginx/access.log  main;
error_log  /var/log/nginx/error.log warn;
```

**記錄內容**:
- 客戶端 IP
- 請求時間
- 請求路徑
- HTTP 狀態碼
- User-Agent
- 回應時間

---

## 🏗️ 系統架構圖

```
                           ┌─────────────────┐
                           │   外部訪問       │
                           │  (Port 80/443)  │
                           └────────┬────────┘
                                    │
                        ┌───────────▼──────────┐
                        │      Nginx           │
                        │   反向代理 + 閘道     │
                        │   (統一入口點)       │
                        └───────────┬──────────┘
                                    │
                ┌───────────────────┼───────────────────┐
                │                   │                   │
        ┌───────▼───────┐   ┌──────▼──────┐   ┌──────▼──────┐
        │  /api/        │   │  /grafana/  │   │  /ws        │
        │  axiom-be     │   │  grafana    │   │  WebSocket  │
        │  :3001        │   │  :3000      │   │  即時通訊    │
        └───────────────┘   └─────────────┘   └─────────────┘
                │
        ┌───────┴───────┐
        │               │
    ┌───▼───┐   ┌──────▼──────┐   ┌───────────┐
    │Postgres│   │cyber-ai-    │   │RabbitMQ   │
    │:5432  │   │quantum:8000 │   │:5672      │
    └───────┘   └─────────────┘   └───────────┘
```

---

## 📊 Nginx 在 Pandora 的具體功能清單

### ✅ 核心功能

| 功能 | 說明 | 狀態 |
|------|------|------|
| **API 閘道** | 統一入口 | ✅ |
| **反向代理** | 5 個上游服務 | ✅ |
| **WebSocket** | 長連接支援 | ✅ |
| **安全標頭** | 4 種防護 | ✅ |
| **Gzip 壓縮** | 60-80% 壓縮率 | ✅ |
| **健康檢查** | /health 端點 | ✅ |
| **靜態檔案** | 前端 UI | ✅ |
| **日誌記錄** | 訪問和錯誤日誌 | ✅ |
| **IPv4/IPv6** | 雙協議支援 | ✅ |
| **Keep-Alive** | 連接池優化 | ✅ |

---

## 🔐 安全性優勢

### 1. **單一入口控制**
- 所有外部請求必須經過 Nginx
- 內部服務不直接暴露
- 統一的安全策略

### 2. **DDoS 緩解**
```nginx
worker_connections  1024;  # 限制並發連接
proxy_read_timeout 60s;    # 超時保護
```

### 3. **隱藏後端資訊**
- 不暴露內部服務 IP 和端口
- 統一的錯誤頁面
- 隱藏技術棧細節

---

## 🚀 效能優勢

### 1. **連接復用**
- Keep-Alive: 32 個持久連接
- 減少 TCP 握手開銷

### 2. **壓縮傳輸**
- Gzip 壓縮: 節省 60-80% 頻寬
- 加快載入速度

### 3. **靜態快取**
- 瀏覽器快取: 1 小時
- 減少伺服器負載

---

## 🔄 實際流量範例

### 範例 1: API 請求

```
客戶端請求:
  → http://localhost/api/v1/agents

Nginx 處理:
  1. 接收請求
  2. 添加安全標頭
  3. 記錄訪問日誌
  4. 轉發到 axiom-be:3001/api/v1/agents
  5. 接收回應
  6. Gzip 壓縮
  7. 返回給客戶端
```

### 範例 2: Grafana 訪問

```
客戶端請求:
  → http://localhost/grafana/d/dashboard

Nginx 處理:
  1. 接收請求
  2. URL 重寫: /grafana/d/dashboard → /d/dashboard
  3. 轉發到 grafana:3000/d/dashboard
  4. 返回儀表板頁面
```

### 範例 3: WebSocket 連接

```
客戶端:
  → ws://localhost/ws

Nginx 處理:
  1. 檢測 Upgrade 標頭
  2. 建立 WebSocket 連接到 axiom-be
  3. 維持長連接（最長 7 天）
  4. 雙向數據轉發
```

---

## 📊 Nginx 與其他服務的關係

### 前端架構

```
┌─────────────────────────────────────────────┐
│            瀏覽器/客戶端                     │
└───────────────────┬─────────────────────────┘
                    │ HTTP/HTTPS
          ┌─────────▼─────────┐
          │      Nginx        │ ◄─── 你在這裡
          │   (Port 80/443)   │
          └─────────┬─────────┘
                    │
    ┌───────────────┼───────────────┐
    │               │               │
┌───▼────┐    ┌────▼────┐    ┌────▼────┐
│Axiom-BE│    │ Grafana │    │Prometheus│
│:3001   │    │:3000    │    │:9090     │
└────────┘    └─────────┘    └──────────┘
```

### 作為中間層

- **前端**: 瀏覽器只需訪問一個地址（Port 80）
- **後端**: 各服務在內部網路通訊，不暴露
- **Nginx**: 橋接前後端，提供統一介面

---

## 🎯 為什麼需要 Nginx？

### 優點

#### 1. **簡化訪問**
```bash
# 不使用 Nginx（需要記住多個端口）
http://localhost:3001/api/...      # axiom-be
http://localhost:3000/...          # grafana  
http://localhost:9090/...          # prometheus
http://localhost:8000/...          # cyber-ai-quantum

# 使用 Nginx（統一入口）
http://localhost/api/...           # → axiom-be
http://localhost/grafana/...       # → grafana
http://localhost/prometheus/...    # → prometheus
```

#### 2. **安全性**
- ✅ 內部服務不直接暴露
- ✅ 統一的安全標頭
- ✅ 集中的訪問控制
- ✅ 防止直接攻擊後端

#### 3. **可維護性**
- ✅ 單點配置管理
- ✅ 容易添加/移除服務
- ✅ 統一的日誌格式
- ✅ 版本升級不影響前端

#### 4. **效能**
- ✅ Gzip 壓縮
- ✅ 連接池
- ✅ 靜態檔案快取
- ✅ HTTP/2 支援（可選）

#### 5. **可擴展性**
- ✅ 容易添加負載均衡
- ✅ 支援 SSL/TLS 終止
- ✅ 可整合 WAF
- ✅ 支援 CDN

---

## 🔧 當前配置總結

### 監聽

```
Port 80   (HTTP)  - IPv4 + IPv6 ✅
Port 443  (HTTPS) - 預留（需配置 SSL）
```

### 轉發的服務 (5 個)

| 服務 | 內部地址 | 外部路徑 |
|------|---------|---------|
| Axiom BE | axiom-be:3001 | /api/, /ws |
| Grafana | grafana:3000 | /grafana/ |
| Prometheus | prometheus:9090 | /prometheus/ |
| Loki | loki:3100 | /loki/ |
| AlertManager | alertmanager:9093 | /alertmanager/ |

### 特殊端點 (2 個)

| 端點 | 功能 | 訪問 |
|------|------|------|
| /health | 健康檢查 | 公開 |
| / | 靜態 UI | 公開 |

---

## 📈 效能指標

### 當前配置

```nginx
worker_processes  auto;         # CPU 核心數
worker_connections  1024;       # 每個 worker 可處理 1024 連接
keepalive 32;                   # 每個上游 32 個持久連接
```

### 理論容量

- **並發連接**: ~1024 × CPU核心數
- **每秒請求**: ~5,000 - 10,000 RPS（取決於後端）
- **帶寬節省**: 60-80%（Gzip）

---

## 🛠️ 可擴展功能（未來）

### 1. SSL/TLS 終止

```nginx
server {
    listen 443 ssl http2;
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    
    # SSL 優化
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
}
```

### 2. 速率限制（Rate Limiting）

```nginx
limit_req_zone $binary_remote_addr zone=api:10m rate=100r/s;

location /api/ {
    limit_req zone=api burst=20;
    proxy_pass http://axiom_ui;
}
```

### 3. IP 白名單

```nginx
location /admin/ {
    allow 192.168.1.0/24;
    deny all;
    proxy_pass http://axiom_ui;
}
```

### 4. WAF 整合

```nginx
# 整合 ModSecurity
load_module modules/ngx_http_modsecurity_module.so;
modsecurity on;
modsecurity_rules_file /etc/nginx/modsecurity.conf;
```

---

## 📊 在 Pandora 系統中的重要性

| 方面 | 重要性 | 說明 |
|------|--------|------|
| **統一入口** | 🔴 關鍵 | 所有外部訪問的必經之路 |
| **安全防護** | 🔴 關鍵 | 第一道防線 |
| **路由管理** | 🟡 重要 | 簡化前端開發 |
| **效能優化** | 🟡 重要 | 提升用戶體驗 |
| **可觀測性** | 🟢 輔助 | 統一日誌收集 |

---

## ✅ 總結

### Nginx 在 Pandora Box Console 的角色

**主要角色**:
1. 🎯 **API Gateway** - 統一入口點
2. 🔀 **反向代理** - 路由分發
3. 🛡️ **安全閘道** - 第一道防線
4. ⚡ **效能優化** - 壓縮、快取、連接池
5. 📊 **可觀測性** - 統一日誌

### 具備的功能

- ✅ HTTP/HTTPS 反向代理
- ✅ WebSocket 長連接支援
- ✅ 靜態檔案服務
- ✅ Gzip 壓縮
- ✅ 安全標頭注入
- ✅ 連接池管理
- ✅ 超時控制
- ✅ 健康檢查
- ✅ 錯誤處理
- ✅ 訪問日誌

### 對系統的價值

**沒有 Nginx**:
- ❌ 每個服務需要單獨訪問
- ❌ 內部服務直接暴露
- ❌ 無統一的安全策略
- ❌ 前端需要處理多個 API 端點

**有了 Nginx**:
- ✅ 單一入口，簡化訪問
- ✅ 內部服務受保護
- ✅ 統一安全防護
- ✅ 前端只需一個地址

---

**總結**: Nginx 是 Pandora Box Console 的**統一入口閘道**，扮演**API Gateway + 反向代理 + 安全防護**的三重角色，對系統的可用性、安全性和可維護性至關重要。

---

**版本**: v3.4.1  
**狀態**: ✅ Healthy  
**配置**: `configs/nginx/default-paas.conf`


