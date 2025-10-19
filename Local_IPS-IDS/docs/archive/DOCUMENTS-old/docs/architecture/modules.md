# Pandora Box Console - 新功能實作文件

## 📋 概述

本文件說明 Pandora Box Console IDS-IPS 系統中新增的四大核心功能：

1. **MQTT 協定支援** - 物聯網設備通訊
2. **Load Balancer** - 高可用性負載平衡
3. **Pub/Sub 系統** - 事件驅動架構
4. **Brute Force Protection** - 暴力攻擊防護

## 🏗️ 系統架構

```
┌─────────────────────────────────────────────────────────────┐
│                   Pandora Box Console                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Rate Limit  │  │     MQTT     │  │   Pub/Sub    │      │
│  │  Middleware  │  │    Broker    │  │   System     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │     Gin      │  │     Auth     │  │    Metrics   │      │
│  │    Router    │  │   Handler    │  │  Collector   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│                                                               │
│  ┌─────────────────────────────────────────────────┐        │
│  │           Load Balancer (Optional)              │        │
│  │  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐        │        │
│  │  │Agent1│  │Agent2│  │Agent3│  │Agent4│        │        │
│  │  └──────┘  └──────┘  └──────┘  └──────┘        │        │
│  └─────────────────────────────────────────────────┘        │
└─────────────────────────────────────────────────────────────┘
         │                  │                │
         ▼                  ▼                ▼
    ┌────────┐         ┌────────┐      ┌────────┐
    │ Redis  │         │ MQTT   │      │Postgres│
    │(Pub/Sub)         │ Broker │      │   DB   │
    └────────┘         └────────┘      └────────┘
```

## 1️⃣ MQTT 協定支援

### 功能特性

- ✅ 完整的 MQTT v3.1.1 協定支援
- ✅ QoS 0/1/2 訊息可靠性保證
- ✅ TLS/SSL 加密通訊
- ✅ 自動重連機制
- ✅ 主題訂閱/發布管理
- ✅ Clean Session 支援
- ✅ Last Will & Testament (LWT)

### 使用範例

#### 基本配置

```yaml
# configs/console-config.yaml
mqtt:
  enabled: true
  broker: "mqtt.example.com"
  port: 1883
  client_id: "pandora-console"
  username: "console_user"
  password: "${MQTT_PASSWORD}"
  default_qos: 1
  auto_reconnect: true
```

#### 訂閱主題

```go
// 訂閱設備認證請求
mqttBroker.Subscribe("device/+/auth/request", func(topic string, payload []byte) error {
    var authReq AuthRequest
    if err := json.Unmarshal(payload, &authReq); err != nil {
        return err
    }
    
    // 處理認證請求
    result := processAuth(authReq)
    
    // 回應結果
    responseTopic := fmt.Sprintf("device/%s/auth/response", authReq.DeviceID)
    return mqttBroker.Publish(responseTopic, result, 1, false)
})
```

#### 發布訊息

```go
// 發布認證成功事件
event := AuthEvent{
    Type:      "auth_success",
    DeviceID:  "device-001",
    Timestamp: time.Now(),
}

data, _ := json.Marshal(event)
mqttBroker.Publish("console/events/auth", data, 1, false)
```

### MQTT 主題架構

```
device/
├── {device_id}/
│   ├── auth/
│   │   ├── request          # 設備發送認證請求
│   │   └── response         # Console 回應認證結果
│   ├── status               # 設備狀態更新
│   ├── events               # 設備事件
│   └── commands             # Console 發送的命令

agent/
├── {agent_id}/
│   ├── logs                 # Agent 日誌
│   ├── metrics              # Agent 指標
│   └── alerts               # Agent 告警

console/
├── broadcast                # 廣播訊息
├── events/
│   ├── auth                 # 認證事件
│   ├── security             # 安全事件
│   └── system               # 系統事件
```

## 2️⃣ Load Balancer

### 功能特性

- ✅ Round Robin 策略
- ✅ Least Connections 策略
- ✅ IP Hash 策略
- ✅ 自動健康檢查
- ✅ 故障轉移
- ✅ 連接池管理
- ✅ 統計資訊收集

### 配置範例

```yaml
loadbalancer:
  enabled: true
  backends:
    - "http://agent-1:8080"
    - "http://agent-2:8080"
    - "http://agent-3:8080"
  
  strategy: "least_connections"
  
  health_check_enabled: true
  health_check_interval: "10s"
  health_check_timeout: "5s"
  health_check_path: "/health"
  
  max_retries: 3
  retry_delay: "1s"
```

### 使用範例

```go
// 初始化 Load Balancer
lb, err := loadbalancer.NewLoadBalancer(config, logger)
if err != nil {
    log.Fatal(err)
}

// 獲取後端服務器
backend, err := lb.GetNextBackend()
if err != nil {
    return err
}

// 代理請求
backend.ReverseProxy.ServeHTTP(w, r)

// 查看統計資訊
stats := lb.GetStats()
for _, stat := range stats {
    fmt.Printf("Backend: %s, Alive: %v, Connections: %d\n",
        stat["url"], stat["alive"], stat["active_connections"])
}
```

### API 端點

```bash
# 查看 Load Balancer 統計
GET /api/v1/admin/loadbalancer/stats

# 回應範例
{
  "backends": [
    {
      "url": "http://agent-1:8080",
      "alive": true,
      "active_connections": 5,
      "total_requests": 1234,
      "failed_requests": 2,
      "last_check": "2024-01-01T12:00:00Z"
    }
  ]
}
```

## 3️⃣ Pub/Sub 事件系統

### 功能特性

- ✅ Redis Pub/Sub 支援
- ✅ In-Memory Pub/Sub（單機部署）
- ✅ 多訂閱者支援
- ✅ 事件廣播
- ✅ 異步處理
- ✅ 緩衝區管理

### 配置範例

```yaml
pubsub:
  enabled: true
  type: "redis"  # 或 "memory"
  redis_addr: "redis:6379"
  redis_password: "${REDIS_PASSWORD}"
  redis_db: 0
  buffer_size: 100
  
  channels:
    - "auth.events"
    - "security.events"
    - "device.events"
    - "system.events"
```

### 使用範例

#### 訂閱事件

```go
// 訂閱認證事件
pubsubInstance.Subscribe(ctx, "auth.events", func(topic string, message []byte) error {
    var event AuthEvent
    if err := json.Unmarshal(message, &event); err != nil {
        return err
    }
    
    logger.Infof("收到認證事件: %+v", event)
    
    // 處理事件
    handleAuthEvent(event)
    return nil
})
```

#### 發布事件

```go
// 發布安全事件
event := SecurityEvent{
    Type:      "brute_force_detected",
    Source:    clientIP,
    Timestamp: time.Now(),
    Data: map[string]interface{}{
        "attempts": 5,
        "blocked":  true,
    },
}

pubsubInstance.Publish(ctx, "security.events", event)
```

### 事件類型

#### 1. 認證事件 (auth.events)

```json
{
  "type": "auth_success | auth_failure | tpm_verify",
  "source": "device_id | pc_identifier",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "method": "psk | tpm",
    "reason": "...",
    "metadata": {}
  }
}
```

#### 2. 安全事件 (security.events)

```json
{
  "type": "brute_force | ip_blocked | suspicious_activity",
  "source": "ip_address",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "severity": "low | medium | high | critical",
    "details": "...",
    "action_taken": "..."
  }
}
```

## 4️⃣ Brute Force Protection

### 功能特性

- ✅ Token Bucket 演算法
- ✅ Sliding Window 支援
- ✅ IP 封鎖機制
- ✅ 帳號鎖定
- ✅ 動態速率限制
- ✅ 白名單支援
- ✅ 自動清理過期資料

### 配置範例

```yaml
ratelimit:
  enabled: true
  
  # Token Bucket 設定
  rate: 10              # 每秒 10 個請求
  burst: 20             # 突發容量 20
  window_size: "1s"
  
  # 暴力攻擊防護
  max_attempts: 5       # 最大失敗次數
  lockout_time: "10m"   # 鎖定 10 分鐘
  
  # IP 封鎖
  block_enabled: true
  block_time: "1h"      # 封鎖 1 小時
  
  cleanup_interval: "5m"
  
  # 白名單
  whitelist_ips:
    - "127.0.0.1"
    - "10.0.0.0/8"
```

### 使用範例

#### 中間件整合

```go
// 全域速率限制
router.Use(rateLimitMiddleware.Handler())

// 特定路由的暴力攻擊防護
authGroup := v1.Group("/auth")
authGroup.Use(rateLimitMiddleware.BruteForceProtection())
{
    authGroup.POST("/login", handleLogin)
    authGroup.POST("/verify", handleVerify)
}
```

#### 手動檢查

```go
// 檢查是否允許請求
allowed, err := rateLimiter.Allow(clientIP)
if !allowed {
    return http.StatusTooManyRequests, "Too many requests"
}

// 記錄失敗嘗試
if authFailed {
    rateLimiter.RecordFailedAttempt(clientIP)
}

// 檢查是否被鎖定
if rateLimiter.IsLocked(clientIP) {
    return http.StatusForbidden, "Account locked"
}

// 檢查是否被封鎖
if rateLimiter.IsBlocked(clientIP) {
    return http.StatusForbidden, "IP blocked"
}

// 獲取狀態
status := rateLimiter.GetStatus(clientIP)
fmt.Printf("Remaining tokens: %d\n", status["tokens"])
fmt.Printf("Locked: %v\n", status["locked"])
fmt.Printf("Blocked: %v\n", status["blocked"])
```

### API 端點

```bash
# 查看 Rate Limit 狀態
GET /api/v1/admin/ratelimit/status/:key

# 回應範例
{
  "key": "192.168.1.100",
  "exists": true,
  "tokens": 5,
  "attempts": 2,
  "locked": false,
  "blocked": false,
  "locked_until": null,
  "blocked_until": null
}
```

### HTTP Headers

請求回應會包含速率限制資訊：

```http
X-RateLimit-Limit: 10
X-RateLimit-Remaining: 5
Retry-After: 60
```

## 🚀 部署指南

### 1. Docker Compose 部署

```yaml
version: '3.8'

services:
  console:
    image: pandora-console:latest
    ports:
      - "3001:3001"
    environment:
      - MQTT_BROKER=mosquitto
      - REDIS_HOST=redis
      - DATABASE_HOST=postgres
    depends_on:
      - redis
      - mosquitto
      - postgres
    volumes:
      - ./configs:/app/configs

  mosquitto:
    image: eclipse-mosquitto:2
    ports:
      - "1883:1883"
      - "9001:9001"
    volumes:
      - ./mosquitto/config:/mosquitto/config

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --requirepass ${REDIS_PASSWORD}

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=pandora
      - POSTGRES_PASSWORD=${DATABASE_PASSWORD}
      - POSTGRES_DB=pandora
```

### 2. Kubernetes 部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pandora-console
spec:
  replicas: 3
  selector:
    matchLabels:
      app: pandora-console
  template:
    metadata:
      labels:
        app: pandora-console
    spec:
      containers:
      - name: console
        image: pandora-console:latest
        ports:
        - containerPort: 3001
        env:
        - name: MQTT_BROKER
          value: "mqtt-service"
        - name: REDIS_HOST
          value: "redis-service"
        - name: RATELIMIT_ENABLED
          value: "true"
        - name: LOADBALANCER_ENABLED
          value: "true"
```

## 🧪 測試

### 單元測試

```bash
# 執行所有測試
go test ./...

# 測試特定模組
go test ./internal/mqtt/...
go test ./internal/ratelimit/...
go test ./internal/pubsub/...
go test ./internal/loadbalancer/...

# 帶覆蓋率
go test -cover ./...
```

### 整合測試

```bash
# MQTT 測試
mosquitto_pub -h localhost -t "device/test/auth/request" -m '{"device_id":"test","token":"123"}'

# Rate Limit 測試
for i in {1..15}; do curl http://localhost:3001/api/v1/verify/pc; done

# Pub/Sub 測試
redis-cli PUBLISH auth.events '{"type":"test","source":"test"}'

# Load Balancer 測試
curl http://localhost:3001/api/v1/admin/loadbalancer/stats
```

## 📊 監控與觀察

### Prometheus 指標

```prometheus
# Rate Limit 指標
pandora_ratelimit_requests_total{status="allowed|blocked"}
pandora_ratelimit_locks_total
pandora_ratelimit_blocks_total

# MQTT 指標
pandora_mqtt_messages_published_total
pandora_mqtt_messages_received_total
pandora_mqtt_connection_status

# Pub/Sub 指標
pandora_pubsub_events_published_total{topic="..."}
pandora_pubsub_events_consumed_total{topic="..."}

# Load Balancer 指標
pandora_loadbalancer_backend_status{backend="..."}
pandora_loadbalancer_requests_total{backend="..."}
pandora_loadbalancer_active_connections{backend="..."}
```

### Grafana Dashboard

請參考 `configs/grafana/dashboards/pandora-console.json`

## 🔒 安全建議

1. **MQTT**:
   - 始終啟用 TLS 加密
   - 使用強密碼
   - 定期輪換憑證
   - 限制主題存取權限

2. **Rate Limiting**:
   - 根據業務需求調整速率
   - 維護白名單
   - 監控異常封鎖
   - 定期審查鎖定日誌

3. **Pub/Sub**:
   - 驗證訊息來源
   - 加密敏感數據
   - 限制訂閱權限
   - 監控訊息流量

4. **Load Balancer**:
   - 啟用健康檢查
   - 設定適當的逾時
   - 監控後端狀態
   - 實作熔斷機制

## 📚 參考資源

- [MQTT 規範](https://mqtt.org/)
- [Redis Pub/Sub](https://redis.io/topics/pubsub)
- [Token Bucket Algorithm](https://en.wikipedia.org/wiki/Token_bucket)
- [Load Balancing Algorithms](https://en.wikipedia.org/wiki/Load_balancing_(computing))

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

## 📄 授權

MIT License

