# RabbitMQ 快速啟動指南
## 5 分鐘開始使用 Pandora Box 消息隊列

---

## 🚀 快速啟動

### 1. 啟動 RabbitMQ

```bash
cd deployments/onpremise
docker-compose up -d rabbitmq

# 等待 RabbitMQ 啟動（約 10 秒）
docker-compose logs -f rabbitmq
```

### 2. 驗證安裝

訪問管理界面：http://localhost:15672

- **用戶名**: `pandora`
- **密碼**: `pandora123`

檢查以下內容：
- ✅ Exchanges: `pandora.events` (topic)
- ✅ Queues: `threat_events`, `network_events`, `system_events`, `device_events`
- ✅ Bindings: 4 個綁定關係

### 3. 測試發布消息

創建測試文件 `test_publisher.go`:

```go
package main

import (
	"context"
	"log"
	"github.com/your-org/pandora-box/internal/pubsub"
)

func main() {
	// 連接 RabbitMQ
	config := pubsub.DefaultConfig()
	mq, err := pubsub.NewRabbitMQ(config)
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close()

	// 創建測試事件
	event := pubsub.NewThreatEvent(
		"test_attack",
		"192.168.1.100",
		"Test threat event",
		"logged",
		5,
	)

	// 發布事件
	message, _ := pubsub.ToJSON(event)
	err = mq.Publish(context.Background(), "pandora.events", "threat.detected", message)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Event published successfully!")
}
```

運行：
```bash
go run test_publisher.go
```

### 4. 測試訂閱消息

創建測試文件 `test_subscriber.go`:

```go
package main

import (
	"context"
	"log"
	"github.com/your-org/pandora-box/internal/pubsub"
)

func main() {
	// 連接 RabbitMQ
	config := pubsub.DefaultConfig()
	mq, err := pubsub.NewRabbitMQ(config)
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close()

	// 定義處理函數
	handler := func(ctx context.Context, msg *pubsub.Message) error {
		var event pubsub.ThreatEvent
		if err := pubsub.FromJSON(msg.Body, &event); err != nil {
			return err
		}

		log.Printf("✅ Received threat: %s from %s (level: %d)",
			event.ThreatType,
			event.SourceIP,
			event.ThreatLevel,
		)
		return nil
	}

	// 訂閱事件
	err = mq.Subscribe(context.Background(), "threat_events", handler)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("📡 Listening for threat events...")
	select {} // 保持運行
}
```

運行：
```bash
go run test_subscriber.go
```

---

## 📦 事件類型速查

### 威脅事件

```go
event := pubsub.NewThreatEvent("ddos", "192.168.1.100", "DDoS attack", "blocked", 8)
mq.Publish(ctx, "pandora.events", "threat.detected", message)
```

### 網路事件

```go
event := pubsub.NewNetworkEvent("port_scan", "192.168.1.100", "10.0.0.1", "tcp")
mq.Publish(ctx, "pandora.events", "network.scan", message)
```

### 系統事件

```go
event := pubsub.NewSystemEvent("pandora-agent", "running", "Agent started")
mq.Publish(ctx, "pandora.events", "system.started", message)
```

### 設備事件

```go
event := pubsub.NewDeviceEvent("usb-001", "usb-serial", "connected")
mq.Publish(ctx, "pandora.events", "device.connected", message)
```

---

## 🔍 故障排除

### RabbitMQ 無法啟動

```bash
# 檢查端口是否被佔用
netstat -an | grep 5672
netstat -an | grep 15672

# 查看日誌
docker-compose logs rabbitmq

# 重啟服務
docker-compose restart rabbitmq
```

### 連接被拒絕

```bash
# 檢查 RabbitMQ 狀態
docker-compose ps rabbitmq

# 檢查網路連接
docker-compose exec rabbitmq rabbitmq-diagnostics ping
```

### 消息未被消費

1. 檢查隊列中的消息數量（管理界面）
2. 確認消費者正在運行
3. 查看消費者日誌

---

## 📚 下一步

- 📖 閱讀完整文檔：[message-queue.md](architecture/message-queue.md)
- 🔧 查看實施路線圖：[IMPLEMENTATION-ROADMAP.md](IMPLEMENTATION-ROADMAP.md)
- ✅ 查看 TODO 清單：[TODO.md](../TODO.md)

---

**需要幫助？** 查看 [故障排除指南](architecture/message-queue.md#故障排除)

