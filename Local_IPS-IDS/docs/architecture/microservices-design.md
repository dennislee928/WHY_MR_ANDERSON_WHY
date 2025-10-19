# 微服務架構設計
## Pandora Box Console IDS-IPS

> 📖 基於專家反饋的微服務拆分方案  
> 🎯 目標：降低單點故障風險，提高系統可擴展性

---

## 📋 概述

根據 `newspec.md` 專家分析，當前 Pandora Agent 作為中心樞紐，集成了太多功能：
- USB-SERIAL CH340 設備管理
- 網路流量監控
- 網路控制和阻斷
- 與 Grafana/Prometheus 通訊

**問題**：
- 單點故障風險高
- 難以獨立擴展
- 維護複雜度高

**解決方案**：拆分為 3 個獨立的微服務

---

## 🏗️ 微服務架構

### 整體架構圖

```
┌─────────────────────────────────────────────────────────────────────┐
│                          Pandora Box Console                         │
└─────────────────────────────────────────────────────────────────────┘

┌──────────────────┐      ┌──────────────────┐      ┌──────────────────┐
│  Device Service  │      │ Network Service  │      │ Control Service  │
│                  │      │                  │      │                  │
│  • USB-SERIAL    │      │  • 流量監控      │      │  • 網路阻斷      │
│  • CH340 驅動    │      │  • 封包分析      │      │  • 防火牆規則    │
│  • 設備管理      │      │  • 異常檢測      │      │  • IP 黑名單     │
│  • 數據採集      │      │  • 統計分析      │      │  • 端口控制      │
└────────┬─────────┘      └────────┬─────────┘      └────────┬─────────┘
         │                         │                         │
         │ gRPC                    │ gRPC                    │ gRPC
         │                         │                         │
         └─────────────────────────┼─────────────────────────┘
                                   │
                                   ▼
                         ┌──────────────────┐
                         │    RabbitMQ      │
                         │  Message Queue   │
                         └────────┬─────────┘
                                  │
                                  ▼
                         ┌──────────────────┐
                         │  Axiom Engine    │
                         │  • 事件處理      │
                         │  • 威脅分析      │
                         │  • 告警管理      │
                         └──────────────────┘
```

---

## 🎯 微服務定義

### 1. Device Service

**職責**：
- 管理 USB-SERIAL CH340 設備
- 讀取硬體輸入數據
- 設備狀態監控
- 設備錯誤處理

**API**：
```protobuf
service DeviceService {
  rpc Connect(ConnectRequest) returns (ConnectResponse);
  rpc Disconnect(DisconnectRequest) returns (DisconnectResponse);
  rpc ReadData(ReadDataRequest) returns (stream DataResponse);
  rpc GetStatus(StatusRequest) returns (StatusResponse);
  rpc ListDevices(ListDevicesRequest) returns (ListDevicesResponse);
}
```

**事件發布**：
- `device.connected` - 設備連接
- `device.disconnected` - 設備斷開
- `device.data` - 設備數據
- `device.error` - 設備錯誤

**依賴**：
- RabbitMQ（發布事件）
- 無其他服務依賴

**端口**：
- gRPC: 50051
- HTTP Health: 8081

---

### 2. Network Service

**職責**：
- 監控網路流量
- 分析網路封包
- 檢測異常流量
- 統計網路指標

**API**：
```protobuf
service NetworkService {
  rpc StartMonitoring(MonitorRequest) returns (MonitorResponse);
  rpc StopMonitoring(StopRequest) returns (StopResponse);
  rpc GetStatistics(StatsRequest) returns (StatsResponse);
  rpc AnalyzeTraffic(AnalyzeRequest) returns (stream AnalysisResponse);
  rpc DetectAnomalies(AnomalyRequest) returns (stream AnomalyResponse);
}
```

**事件發布**：
- `network.attack` - 網路攻擊
- `network.scan` - 端口掃描
- `network.anomaly` - 異常流量
- `network.statistics` - 網路統計

**依賴**：
- RabbitMQ（發布事件）
- 無其他服務依賴

**端口**：
- gRPC: 50052
- HTTP Health: 8082

---

### 3. Control Service

**職責**：
- 執行網路阻斷
- 管理防火牆規則
- 維護 IP 黑名單
- 控制端口訪問

**API**：
```protobuf
service ControlService {
  rpc BlockIP(BlockIPRequest) returns (BlockIPResponse);
  rpc UnblockIP(UnblockIPRequest) returns (UnblockIPResponse);
  rpc BlockPort(BlockPortRequest) returns (BlockPortResponse);
  rpc UnblockPort(UnblockPortRequest) returns (UnblockPortResponse);
  rpc GetBlockList(BlockListRequest) returns (BlockListResponse);
  rpc ApplyFirewallRule(FirewallRuleRequest) returns (FirewallRuleResponse);
}
```

**事件發布**：
- `network.blocked` - IP/端口被阻斷
- `network.unblocked` - IP/端口解除阻斷
- `control.rule_applied` - 防火牆規則應用
- `control.error` - 控制錯誤

**依賴**：
- RabbitMQ（發布事件）
- Device Service（獲取設備狀態）
- Network Service（獲取流量統計）

**端口**：
- gRPC: 50053
- HTTP Health: 8083

---

## 🔄 服務間通訊

### 通訊模式

| 場景 | 通訊方式 | 說明 |
|------|----------|------|
| 事件通知 | RabbitMQ | 非同步事件發布/訂閱 |
| 服務調用 | gRPC | 同步 RPC 調用 |
| 健康檢查 | HTTP | RESTful API |

### 通訊流程

#### 場景 1: 威脅檢測和阻斷

```
1. Device Service 檢測到異常輸入
   └─> 發布 device.data 事件到 RabbitMQ

2. Network Service 監控到攻擊流量
   └─> 發布 network.attack 事件到 RabbitMQ

3. Axiom Engine 接收事件並分析
   └─> 決定需要阻斷

4. Engine 調用 Control Service (gRPC)
   └─> BlockIP("192.168.1.100")

5. Control Service 執行阻斷
   └─> 發布 network.blocked 事件到 RabbitMQ

6. Engine 接收確認事件
   └─> 更新資料庫和儀表板
```

#### 場景 2: 設備狀態查詢

```
1. UI 需要查詢設備狀態
   └─> 調用 Device Service (gRPC)

2. Device Service 返回設備列表
   └─> ListDevices() → [device1, device2, ...]

3. UI 顯示設備狀態
```

---

## 📦 服務部署

### Docker Compose 配置

```yaml
services:
  # Device Service
  device-service:
    build:
      context: .
      dockerfile: Dockerfile.device
    container_name: device-service
    restart: unless-stopped
    privileged: true
    ports:
      - "50051:50051"  # gRPC
      - "8081:8081"    # HTTP Health
    environment:
      - RABBITMQ_URL=amqp://pandora:pandora123@rabbitmq:5672/
      - DEVICE_PORT=/dev/ttyUSB0
    volumes:
      - /dev:/dev
    depends_on:
      - rabbitmq

  # Network Service
  network-service:
    build:
      context: .
      dockerfile: Dockerfile.network
    container_name: network-service
    restart: unless-stopped
    network_mode: host
    ports:
      - "50052:50052"  # gRPC
      - "8082:8082"    # HTTP Health
    environment:
      - RABBITMQ_URL=amqp://pandora:pandora123@rabbitmq:5672/
    depends_on:
      - rabbitmq

  # Control Service
  control-service:
    build:
      context: .
      dockerfile: Dockerfile.control
    container_name: control-service
    restart: unless-stopped
    privileged: true
    network_mode: host
    ports:
      - "50053:50053"  # gRPC
      - "8083:8083"    # HTTP Health
    environment:
      - RABBITMQ_URL=amqp://pandora:pandora123@rabbitmq:5672/
      - DEVICE_SERVICE_URL=device-service:50051
      - NETWORK_SERVICE_URL=network-service:50052
    depends_on:
      - rabbitmq
      - device-service
      - network-service
```

---

## 🔐 安全考量

### 服務間認證

所有 gRPC 通訊使用 mTLS：

```go
// 服務端配置
tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{serverCert},
    ClientAuth:   tls.RequireAndVerifyClientCert,
    ClientCAs:    clientCAPool,
}

// 客戶端配置
tlsConfig := &tls.Config{
    Certificates: []tls.Certificate{clientCert},
    RootCAs:      serverCAPool,
}
```

### 授權

使用 JWT Token 進行服務間授權：

```go
// 在 gRPC metadata 中傳遞 token
md := metadata.Pairs("authorization", "Bearer "+token)
ctx := metadata.NewOutgoingContext(context.Background(), md)
```

---

## 📊 監控和觀測

### 健康檢查

每個服務提供 HTTP 健康檢查端點：

```
GET /health
{
  "status": "healthy",
  "service": "device-service",
  "version": "1.0.0",
  "uptime": 3600,
  "dependencies": {
    "rabbitmq": "healthy",
    "device": "connected"
  }
}
```

### Prometheus 指標

每個服務暴露 Prometheus 指標：

```
# Device Service
device_connections_total
device_read_operations_total
device_errors_total
device_data_bytes_total

# Network Service
network_packets_total
network_bytes_total
network_attacks_detected_total
network_anomalies_total

# Control Service
control_blocks_total
control_unblocks_total
control_rules_applied_total
control_errors_total
```

---

## 🧪 測試策略

### 單元測試

每個服務的核心邏輯：
```bash
go test -v ./internal/services/device/...
go test -v ./internal/services/network/...
go test -v ./internal/services/control/...
```

### 集成測試

服務間通訊測試：
```bash
go test -v -tags=integration ./tests/integration/...
```

### 端到端測試

完整流程測試：
```bash
go test -v -tags=e2e ./tests/e2e/...
```

### 性能測試

壓力測試和負載測試：
```bash
# 使用 ghz 進行 gRPC 性能測試
ghz --insecure \
  --proto api/proto/device.proto \
  --call pandora.DeviceService/GetStatus \
  -d '{"device_id":"usb-001"}' \
  -n 10000 \
  -c 100 \
  localhost:50051
```

---

## 🔄 遷移策略

### 階段 1: 並行運行（Week 2-3）

```
┌─────────────────┐
│  舊 Agent       │ ← 保持運行
│  (Monolith)     │
└─────────────────┘

┌─────────────────┐
│  新微服務       │ ← 並行部署
│  (Microservices)│
└─────────────────┘
```

### 階段 2: 流量切換（Week 3-4）

```
舊 Agent: 90% 流量
新微服務: 10% 流量 ← 逐步增加
```

### 階段 3: 完全遷移（Week 4）

```
舊 Agent: 下線
新微服務: 100% 流量 ← 完全接管
```

---

## 📈 擴展性設計

### 水平擴展

每個服務都可以獨立擴展：

```yaml
# 擴展 Network Service 到 3 個實例
docker-compose up -d --scale network-service=3
```

### 負載均衡

使用 gRPC 客戶端負載均衡：

```go
// Round-robin 負載均衡
conn, err := grpc.Dial(
    "dns:///network-service:50052",
    grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
)
```

---

## 🎯 成功指標

| 指標 | 目標 | 測量方法 |
|------|------|----------|
| 服務獨立運行 | 100% | 單獨啟動每個服務 |
| 服務間通訊延遲 | < 50ms | gRPC 調用延遲 |
| 單服務故障隔離 | 100% | 停止一個服務，其他服務正常 |
| 水平擴展能力 | 支援 10+ 實例 | 壓力測試 |
| 部署時間 | < 5 分鐘 | Docker Compose 部署 |

---

## 📚 相關文檔

- [gRPC 服務定義](grpc-services.md)
- [服務發現機制](service-discovery.md)
- [監控和告警](monitoring.md)
- [部署指南](../deployment/microservices.md)

---

## 🔜 實施計劃

### Day 1-2: 定義 gRPC 接口
- [ ] 創建 proto 文件
- [ ] 生成 Go 代碼
- [ ] 定義服務接口

### Day 3-4: 實現 Device Service
- [ ] 創建服務框架
- [ ] 實現設備管理邏輯
- [ ] 整合 RabbitMQ
- [ ] 添加測試

### Day 5-6: 實現 Network Service
- [ ] 創建服務框架
- [ ] 實現流量監控邏輯
- [ ] 整合 RabbitMQ
- [ ] 添加測試

### Day 7-8: 實現 Control Service
- [ ] 創建服務框架
- [ ] 實現控制邏輯
- [ ] 整合 gRPC 客戶端
- [ ] 添加測試

### Day 9-10: 集成和測試
- [ ] 更新 Docker Compose
- [ ] 端到端測試
- [ ] 性能測試
- [ ] 文檔更新

---

**設計者**: AI Assistant  
**審核者**: @tech-lead  
**版本**: 1.0.0  
**最後更新**: 2025-10-09

