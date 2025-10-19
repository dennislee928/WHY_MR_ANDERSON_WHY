# RabbitMQ 整合示例
## Pandora Box Console IDS-IPS

這個目錄包含了 RabbitMQ 整合的示例代碼，展示如何在 Agent 和 Engine 中使用消息隊列。

---

## 🚀 快速開始

### 1. 啟動 RabbitMQ

```bash
cd ../../deployments/onpremise
docker-compose up -d rabbitmq

# 等待 RabbitMQ 啟動
sleep 10

# 檢查狀態
docker-compose ps rabbitmq
```

### 2. 運行 Agent 示例（發布者）

在一個終端中運行：

```bash
cd examples/rabbitmq-integration

# 設置環境變數
export RABBITMQ_URL="amqp://pandora:pandora123@localhost:5672/"
export RABBITMQ_EXCHANGE="pandora.events"

# 運行 Agent 示例
go run agent_example.go
```

**輸出示例**：
```
INFO[0000] Event publisher initialized
INFO[0000] Agent running, press Ctrl+C to stop...
INFO[0010] Published threat detection event
INFO[0015] Published network attack event
INFO[0020] Published threat detection event
```

### 3. 運行 Engine 示例（訂閱者）

在另一個終端中運行：

```bash
cd examples/rabbitmq-integration

# 設置環境變數
export RABBITMQ_URL="amqp://pandora:pandora123@localhost:5672/"
export RABBITMQ_EXCHANGE="pandora.events"

# 運行 Engine 示例
go run engine_example.go
```

**輸出示例**：
```
INFO[0000] Event subscriber initialized
INFO[0000] Engine is now listening for events...
INFO[0000] Subscribed to:
INFO[0000]   - threat_events (威脅事件)
INFO[0000]   - network_events (網路事件)
INFO[0000]   - system_events (系統事件)
INFO[0000]   - device_events (設備事件)
INFO[0010] Received threat event: ddos from 192.168.1.100 (level: 8)
INFO[0015] Received network event: port_scan from 192.168.1.101 to 10.0.0.1
```

---

## 📊 觀察消息流

### 使用 RabbitMQ 管理界面

1. 訪問 http://localhost:15672
2. 登入（用戶名: `pandora`, 密碼: `pandora123`）
3. 查看：
   - **Exchanges** → `pandora.events` - 查看消息路由
   - **Queues** → 查看各個隊列的消息數量
   - **Connections** → 查看 Agent 和 Engine 的連接

### 使用命令行工具

```bash
# 查看隊列狀態
docker-compose exec rabbitmq rabbitmqctl list_queues name messages consumers

# 查看交換機
docker-compose exec rabbitmq rabbitmqctl list_exchanges name type

# 查看綁定關係
docker-compose exec rabbitmq rabbitmqctl list_bindings
```

---

## 🧪 測試場景

### 場景 1: 威脅檢測和分析

1. Agent 檢測到 DDoS 攻擊
2. Agent 發布 `threat.detected` 事件
3. Engine 接收事件並進行分析
4. Engine 更新威脅資料庫

### 場景 2: 網路攻擊響應

1. Agent 檢測到端口掃描
2. Agent 發布 `network.attack` 事件
3. Engine 接收事件並分析模式
4. Engine 觸發自動阻斷（如果配置）

### 場景 3: 系統健康監控

1. Agent 定期發布健康檢查事件
2. Engine 接收並記錄系統狀態
3. 如果健康檢查失敗，觸發告警

### 場景 4: 設備管理

1. USB 設備連接
2. Agent 發布 `device.connected` 事件
3. Engine 更新設備狀態
4. 設備斷開時發布 `device.disconnected` 事件

---

## 📈 性能測試

### 測試消息延遲

```bash
# 在 Agent 示例中添加時間戳
# 在 Engine 示例中計算延遲

# 預期結果：
# - 發布延遲: < 5ms
# - 端到端延遲: < 100ms
# - 吞吐量: > 1000 msg/s
```

### 測試高負載

```bash
# 修改 Agent 示例，增加發布頻率
# 觀察 RabbitMQ 管理界面的指標
# 確保沒有消息積壓
```

---

## 🔧 配置選項

### 環境變數

```bash
# RabbitMQ 連接
export RABBITMQ_URL="amqp://pandora:pandora123@localhost:5672/"
export RABBITMQ_EXCHANGE="pandora.events"

# 日誌等級
export LOG_LEVEL="debug"
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

## 🐛 故障排除

### Agent 無法連接 RabbitMQ

```bash
# 檢查 RabbitMQ 是否運行
docker-compose ps rabbitmq

# 檢查網路連接
telnet localhost 5672

# 查看 RabbitMQ 日誌
docker-compose logs rabbitmq
```

### 消息未被 Engine 接收

1. 檢查 Engine 是否正在運行
2. 檢查隊列綁定是否正確
3. 查看 RabbitMQ 管理界面的消息統計
4. 檢查 Engine 日誌是否有錯誤

### 消息積壓

1. 增加 Engine 實例數量
2. 增加預取數量（PrefetchCount）
3. 優化消息處理邏輯
4. 檢查是否有阻塞操作

---

## 📚 相關文檔

- [消息隊列架構文檔](../../docs/architecture/message-queue.md)
- [RabbitMQ 快速啟動](../../docs/QUICKSTART-RABBITMQ.md)
- [實施路線圖](../../docs/IMPLEMENTATION-ROADMAP.md)

---

## 🎯 下一步

完成示例測試後，繼續：

1. 整合到實際的 Agent 代碼
2. 整合到實際的 Engine 代碼
3. 添加完整的錯誤處理
4. 添加 Prometheus 監控指標
5. 進行性能測試和優化

---

**維護者**: Pandora Box Team  
**最後更新**: 2025-10-09

