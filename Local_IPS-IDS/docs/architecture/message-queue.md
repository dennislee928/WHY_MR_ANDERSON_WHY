# RabbitMQ 消息隊列整合文檔
## Pandora Box Console IDS-IPS

> 📖 基於專家反饋實現的消息隊列架構，用於解耦 Agent 和 Engine

---

## 📋 概述

本文檔描述如何在 Pandora Box Console IDS-IPS 中使用 RabbitMQ 消息隊列實現服務間的非同步通訊。

### 為什麼需要 RabbitMQ？

根據 `newspec.md` 專家分析：
- **降低耦合度**：Agent 作為中心樞紐，太多依賴導致單點故障風險
- **提高可靠性**：消息持久化，系統重啟後不丟失
- **支援擴展**：為未來的微服務架構奠定基礎
- **非同步處理**：提升系統響應速度和吞吐量

---

## 🏗️ 架構設計

### 消息流向

```
┌─────────────┐     Publish      ┌──────────────┐     Route      ┌─────────────┐
│ Pandora     │ ───────────────> │  RabbitMQ    │ ─────────────> │  Axiom      │
│ Agent       │   (Events)       │  Exchange    │   (Queues)     │  Engine     │
└─────────────┘                  └──────────────┘                └─────────────┘
      │                                 │                               │
      │                                 │                               │
  Threat Event                    pandora.events                 threat_events
  Network Event                   (Topic Exchange)               network_events
  System Event                                                   system_events
  Device Event                                                   device_events
```

### 交換機和隊列

| 交換機 | 類型 | 說明 |
|--------|------|------|
| `pandora.events` | Topic | 主要事件交換機 |
| `pandora.dlx` | Fanout | 死信交換機（處理失敗的消息） |

| 隊列 | 路由鍵 | 說明 |
|------|--------|------|
| `threat_events` | `threat.*` | 威脅事件隊列 |
| `network_events` | `network.*` | 網路事件隊列 |
| `system_events` | `system.*` | 系統事件隊列 |
| `device_events` | `device.*` | 設備事件隊列 |

---

## 🚀 快速開始

### 1. 啟動 RabbitMQ

```bash
# 使用 Docker Compose 啟動
cd deployments/onpremise
docker-compose up -d rabbitmq

# 檢查狀態
docker-compose ps rabbitmq

# 訪問管理界面
open http://localhost:15672
# 用戶名: pandora
# 密碼: pandora123
```

### 2. 發布事件（Publisher）

```go
package main

import (
    "context"
    "log"
    "github.com/your-org/pandora-box/internal/pubsub"
)

func main() {
    // 創建 RabbitMQ 連接
    config := pubsub.DefaultConfig()
    mq, err := pubsub.NewRabbitMQ(config)
    if err != nil {
        log.Fatal(err)
    }
    defer mq.Close()

    // 創建威脅事件
    event := pubsub.NewThreatEvent(
        "ddos",              // 威脅類型
        "192.168.1.100",     // 來源 IP
        "DDoS attack detected", // 描述
        "blocked",           // 動作
        8,                   // 威脅等級 (1-10)
    )

    // 設置額外資訊
    event.TargetIP = "10.0.0.1"
    event.TargetPort = 80
    event.Protocol = "tcp"

    // 轉換為 JSON
    message, err := pubsub.ToJSON(event)
    if err != nil {
        log.Fatal(err)
    }

    // 發布到 RabbitMQ
    ctx := context.Background()
    err = mq.Publish(
        ctx,
        "pandora.events",        // 交換機
        "threat.detected",       // 路由鍵
        message,                 // 消息內容
    )
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Event published successfully")
}
```

### 3. 訂閱事件（Subscriber）

```go
package main

import (
    "context"
    "log"
    "github.com/your-org/pandora-box/internal/pubsub"
)

func main() {
    // 創建 RabbitMQ 連接
    config := pubsub.DefaultConfig()
    mq, err := pubsub.NewRabbitMQ(config)
    if err != nil {
        log.Fatal(err)
    }
    defer mq.Close()

    // 定義消息處理函數
    handler := func(ctx context.Context, msg *pubsub.Message) error {
        log.Printf("Received message: %s", msg.RoutingKey)

        // 解析威脅事件
        var event pubsub.ThreatEvent
        if err := pubsub.FromJSON(msg.Body, &event); err != nil {
            return err
        }

        // 處理事件
        log.Printf("Threat detected: %s from %s (level: %d)",
            event.ThreatType,
            event.SourceIP,
            event.ThreatLevel,
        )

        // 這裡添加您的業務邏輯
        // 例如：分析威脅、更新資料庫、發送告警等

        return nil
    }

    // 訂閱威脅事件隊列
    ctx := context.Background()
    err = mq.Subscribe(ctx, "threat_events", handler)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Subscribed to threat_events, waiting for messages...")

    // 保持運行
    select {}
}
```

---

## 📦 事件類型

### 1. 威脅事件 (ThreatEvent)

用於安全威脅檢測和響應。

**路由鍵**:
- `threat.detected` - 威脅被檢測到
- `threat.blocked` - 威脅被阻斷
- `threat.analyzed` - 威脅分析完成
- `threat.resolved` - 威脅已解決

**範例**:
```go
event := pubsub.NewThreatEvent("ddos", "192.168.1.100", "DDoS attack", "blocked", 8)
event.TargetIP = "10.0.0.1"
event.TargetPort = 80
event.Protocol = "tcp"
event.Evidence["packet_count"] = 10000
event.Evidence["duration"] = "5m"
```

### 2. 網路事件 (NetworkEvent)

用於網路流量監控和分析。

**路由鍵**:
- `network.attack` - 網路攻擊
- `network.scan` - 端口掃描
- `network.anomaly` - 異常流量
- `network.blocked` - 流量被阻斷

**範例**:
```go
event := pubsub.NewNetworkEvent("port_scan", "192.168.1.100", "10.0.0.1", "tcp")
event.SourcePort = 54321
event.DestPort = 22
event.PacketCount = 100
event.Flags = []string{"SYN"}
```

### 3. 系統事件 (SystemEvent)

用於系統狀態監控和告警。

**路由鍵**:
- `system.started` - 系統啟動
- `system.stopped` - 系統停止
- `system.error` - 系統錯誤
- `system.healthy` - 健康檢查通過

**範例**:
```go
event := pubsub.NewSystemEvent("pandora-agent", "running", "Agent started successfully")
event.Metrics["cpu_usage"] = 25.5
event.Metrics["memory_usage"] = 512
event.Metrics["goroutines"] = 42
```

### 4. 設備事件 (DeviceEvent)

用於 IoT 設備監控。

**路由鍵**:
- `device.connected` - 設備連接
- `device.disconnected` - 設備斷開
- `device.data` - 設備數據
- `device.error` - 設備錯誤

**範例**:
```go
event := pubsub.NewDeviceEvent("usb-001", "usb-serial", "connected")
event.Port = "/dev/ttyUSB0"
event.DeviceName = "CH340 Serial"
event.Data["baud_rate"] = 115200
event.Data["timeout"] = "30s"
```

---

## ⚙️ 配置選項

### RabbitMQ 配置

```go
config := &pubsub.Config{
    URL:                  "amqp://pandora:pandora123@localhost:5672/",
    Exchange:             "pandora.events",
    ConnectionTimeout:    30 * time.Second,
    HeartbeatInterval:    60 * time.Second,
    ReconnectDelay:       5 * time.Second,
    MaxReconnectAttempts: 10,
}
```

### 發布選項

```go
opts := &pubsub.PublishOptions{
    ContentType: "application/json",
    Priority:    5,              // 0-9, 數字越大優先級越高
    Expiration:  "60000",        // 消息過期時間（毫秒）
    Persistent:  true,           // 消息持久化
    Headers: map[string]interface{}{
        "source": "pandora-agent",
        "version": "1.0.0",
    },
}
```

### 訂閱選項

```go
opts := &pubsub.SubscribeOptions{
    AutoAck:       false,        // 手動確認（推薦）
    Exclusive:     false,        // 非獨占消費者
    PrefetchCount: 10,           // 預取消息數量
    RetryPolicy: &pubsub.RetryPolicy{
        MaxRetries:      3,
        InitialInterval: 1 * time.Second,
        MaxInterval:     30 * time.Second,
        Multiplier:      2.0,    // 指數退避
    },
}
```

---

## 🔧 最佳實踐

### 1. 消息確認

**總是使用手動確認**，確保消息不丟失：

```go
handler := func(ctx context.Context, msg *pubsub.Message) error {
    // 處理消息
    if err := processMessage(msg); err != nil {
        // 返回錯誤，消息會被 Nack 並重試
        return err
    }
    // 返回 nil，消息會被 Ack
    return nil
}
```

### 2. 錯誤處理

實現重試邏輯和死信隊列：

```go
handler := func(ctx context.Context, msg *pubsub.Message) error {
    // 檢查重試次數
    if msg.Redelivered && msg.Headers["retry_count"].(int) >= 3 {
        log.Printf("Message failed after 3 retries, sending to DLX")
        // 消息會被發送到死信交換機
        return nil
    }

    // 處理消息
    if err := processMessage(msg); err != nil {
        // 增加重試計數
        retryCount := 0
        if count, ok := msg.Headers["retry_count"].(int); ok {
            retryCount = count
        }
        msg.Headers["retry_count"] = retryCount + 1
        return err
    }

    return nil
}
```

### 3. 監控和告警

使用 Prometheus 監控 RabbitMQ：

```go
// 在 Prometheus 配置中添加 RabbitMQ 指標
scrape_configs:
  - job_name: 'rabbitmq'
    static_configs:
      - targets: ['rabbitmq:15692']
```

關鍵指標：
- `rabbitmq_queue_messages` - 隊列中的消息數
- `rabbitmq_queue_messages_ready` - 待處理消息數
- `rabbitmq_queue_messages_unacknowledged` - 未確認消息數
- `rabbitmq_channel_consumers` - 消費者數量

### 4. 性能優化

**批量發布**：
```go
// 使用事務或 publisher confirms
for _, event := range events {
    message, _ := pubsub.ToJSON(event)
    mq.Publish(ctx, "pandora.events", event.Type, message)
}
```

**並發消費**：
```go
// 增加預取數量和消費者數量
opts := &pubsub.SubscribeOptions{
    PrefetchCount: 50,  // 增加預取數量
}

// 啟動多個消費者
for i := 0; i < 5; i++ {
    go mq.Subscribe(ctx, "threat_events", handler)
}
```

---

## 🧪 測試

### 單元測試

```bash
cd internal/pubsub
go test -v
```

### 集成測試

```bash
# 啟動 RabbitMQ
docker-compose up -d rabbitmq

# 運行集成測試
go test -v -tags=integration ./internal/pubsub/...
```

### 手動測試

使用 RabbitMQ 管理界面發布測試消息：

1. 訪問 http://localhost:15672
2. 登入（pandora / pandora123）
3. 進入 "Exchanges" → "pandora.events"
4. 點擊 "Publish message"
5. 設置 Routing key: `threat.detected`
6. 設置 Payload:
```json
{
  "id": "evt_test_001",
  "type": "threat.detected",
  "timestamp": "2025-10-09T12:00:00Z",
  "source": "test",
  "severity": "high",
  "threat_type": "ddos",
  "threat_level": 8,
  "source_ip": "192.168.1.100",
  "description": "Test DDoS attack",
  "action": "blocked"
}
```

---

## 🔍 故障排除

### 問題 1: 連接失敗

**症狀**: `failed to connect to RabbitMQ: dial tcp: connection refused`

**解決方法**:
```bash
# 檢查 RabbitMQ 是否運行
docker-compose ps rabbitmq

# 查看日誌
docker-compose logs rabbitmq

# 重啟 RabbitMQ
docker-compose restart rabbitmq
```

### 問題 2: 消息未被消費

**症狀**: 消息堆積在隊列中

**解決方法**:
1. 檢查消費者是否運行
2. 檢查消費者日誌是否有錯誤
3. 增加消費者數量或預取數量
4. 檢查消息處理邏輯是否有阻塞

### 問題 3: 消息丟失

**症狀**: 消息發布後找不到

**解決方法**:
1. 確保消息持久化：`Persistent: true`
2. 確保隊列持久化：`durable: true`
3. 使用 publisher confirms
4. 檢查死信隊列

---

## 📚 相關文檔

- [實施路線圖](../IMPLEMENTATION-ROADMAP.md)
- [系統架構分析](../../newspec.md)
- [RabbitMQ 官方文檔](https://www.rabbitmq.com/documentation.html)
- [AMQP 協議規範](https://www.amqp.org/specification/0-9-1/amqp-org-download)

---

## 🎯 下一步

完成 RabbitMQ 整合後，繼續實施：

1. ✅ **微服務拆分** - 將 Agent 拆分為獨立的微服務
2. ✅ **gRPC 通訊** - 實現服務間的 gRPC 通訊
3. ✅ **服務發現** - 整合 Consul 或 Kubernetes Service Discovery
4. ✅ **監控告警** - 添加 RabbitMQ 監控指標到 Grafana

參考：[TODO.md](../../TODO.md) 和 [IMPLEMENTATION-ROADMAP.md](../IMPLEMENTATION-ROADMAP.md)

---

**最後更新**: 2025-10-09  
**維護者**: Pandora Box Team  
**版本**: 1.0.0

