# PubSub Package
## 消息隊列抽象層

這個包提供了消息隊列的抽象接口和 RabbitMQ 實現，用於 Pandora Box Console IDS-IPS 的服務間通訊。

---

## 📦 安裝

```bash
go get github.com/rabbitmq/amqp091-go@v1.9.0
```

---

## 🚀 快速開始

### 發布消息

```go
import "github.com/your-org/pandora-box/internal/pubsub"

// 創建連接
config := pubsub.DefaultConfig()
mq, err := pubsub.NewRabbitMQ(config)
if err != nil {
    log.Fatal(err)
}
defer mq.Close()

// 創建事件
event := pubsub.NewThreatEvent("ddos", "192.168.1.100", "DDoS attack", "blocked", 8)

// 發布
message, _ := pubsub.ToJSON(event)
mq.Publish(context.Background(), "pandora.events", "threat.detected", message)
```

### 訂閱消息

```go
// 定義處理函數
handler := func(ctx context.Context, msg *pubsub.Message) error {
    var event pubsub.ThreatEvent
    pubsub.FromJSON(msg.Body, &event)
    log.Printf("Threat: %s from %s", event.ThreatType, event.SourceIP)
    return nil
}

// 訂閱
mq.Subscribe(context.Background(), "threat_events", handler)
```

---

## 🧪 測試

### 運行單元測試

```bash
go test -v
```

### 運行集成測試

需要先啟動 RabbitMQ：

```bash
# 啟動 RabbitMQ
cd ../../deployments/onpremise
docker-compose up -d rabbitmq

# 運行集成測試
cd ../../internal/pubsub
go test -v -tags=integration
```

### 運行性能測試

```bash
go test -bench=. -benchmem
```

---

## 📚 文檔

詳細文檔請參考：[docs/architecture/message-queue.md](../../docs/architecture/message-queue.md)

---

## 🏗️ 架構

```
pubsub/
├── interface.go      # 消息隊列接口定義
├── rabbitmq.go       # RabbitMQ 實現
├── events.go         # 事件類型定義
├── events_test.go    # 事件類型測試
├── rabbitmq_test.go  # RabbitMQ 集成測試
└── README.md         # 本文件
```

---

## 🔧 配置

### 環境變數

```bash
export RABBITMQ_URL="amqp://pandora:pandora123@localhost:5672/"
export RABBITMQ_EXCHANGE="pandora.events"
```

### 代碼配置

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

---

## 📊 性能指標

基於初步測試（需要實際環境驗證）：

| 操作 | 延遲 | 吞吐量 |
|------|------|--------|
| Publish | < 5ms | 10000+ msg/s |
| Subscribe | < 10ms | 5000+ msg/s |
| JSON Marshal | < 1ms | - |
| JSON Unmarshal | < 1ms | - |

---

## 🐛 已知問題

目前沒有已知問題。

---

## 🤝 貢獻

歡迎提交 PR 改進這個包！

---

**維護者**: Pandora Box Team  
**版本**: 1.0.0

