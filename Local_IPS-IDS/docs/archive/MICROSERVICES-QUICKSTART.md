# 微服務快速啟動指南
## Pandora Box Console IDS-IPS

> 🚀 10 分鐘啟動完整的微服務架構

---

## 📋 前置需求

- Docker 和 Docker Compose
- Go 1.21+ (用於開發)
- Protocol Buffers 編譯器 (protoc)

---

## 🚀 快速啟動

### 1. 生成 gRPC 代碼

```bash
cd api/proto

# 安裝 protoc plugins
make install

# 生成所有服務的代碼
make generate

# 驗證生成的文件
ls -la *.pb.go
```

### 2. 啟動所有微服務

```bash
cd deployments/onpremise

# 啟動 RabbitMQ（消息隊列）
docker-compose up -d rabbitmq

# 等待 RabbitMQ 啟動
sleep 10

# 啟動所有微服務
docker-compose up -d device-service network-service control-service

# 檢查服務狀態
docker-compose ps
```

### 3. 驗證服務

```bash
# 檢查 Device Service
curl http://localhost:8081/health

# 檢查 Network Service
curl http://localhost:8082/health

# 檢查 Control Service
curl http://localhost:8083/health

# 預期輸出：
# {"status":"healthy","service":"Device Service","version":"1.0.0"}
```

---

## 🧪 測試微服務

### 使用 grpcurl 測試

```bash
# 安裝 grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# 列出 Device Service 的方法
grpcurl -plaintext localhost:50051 list pandora.device.DeviceService

# 連接設備
grpcurl -plaintext \
  -d '{"device_id":"usb-001","port":"/dev/ttyUSB0","baud_rate":115200}' \
  localhost:50051 \
  pandora.device.DeviceService/Connect

# 獲取設備狀態
grpcurl -plaintext \
  -d '{"device_id":"usb-001"}' \
  localhost:50051 \
  pandora.device.DeviceService/GetStatus

# 開始網路監控
grpcurl -plaintext \
  -d '{"interface_name":"eth0","promiscuous_mode":true}' \
  localhost:50052 \
  pandora.network.NetworkService/StartMonitoring

# 阻斷 IP
grpcurl -plaintext \
  -d '{"ip_address":"192.168.1.100","reason":"DDoS attack","duration_seconds":3600}' \
  localhost:50053 \
  pandora.control.ControlService/BlockIP

# 獲取阻斷列表
grpcurl -plaintext \
  -d '{}' \
  localhost:50053 \
  pandora.control.ControlService/GetBlockList
```

### 使用 Orchestrator 示例

```bash
cd examples/microservices

# 設置環境變數
export DEVICE_SERVICE_URL=localhost:50051
export NETWORK_SERVICE_URL=localhost:50052
export CONTROL_SERVICE_URL=localhost:50053

# 運行 Orchestrator
go run orchestrator.go
```

**預期輸出**:
```
INFO Starting Microservices Orchestrator Example
INFO All service clients connected

========== Scenario 1: Device Management ==========
INFO ✅ Device connected: Device connected successfully
INFO ✅ Device status: DEVICE_STATUS_CONNECTED (uptime: 5s)
INFO ✅ Total devices: 1
INFO    - usb-001: CH340 Serial (/dev/ttyUSB0)

========== Scenario 2: Network Monitoring ==========
INFO ✅ Monitoring started: session session_1728456789
INFO ✅ Network Statistics:
INFO    - Total Packets: 1000
INFO    - Total Bytes: 650000
INFO    - TCP: 800, UDP: 150, ICMP: 50
INFO    - Rate: 200.00 pps, 130000.00 Bps

========== Scenario 3: Threat Detection and Blocking ==========
INFO 🚨 DDoS attack detected from 192.168.1.100
INFO ✅ IP blocked: 192.168.1.100 (rule: block_ip_192.168.1.100_1728456789)
INFO ✅ Block List (Total: 1, Active: 1):
INFO    - BLOCK_TYPE_IP: 192.168.1.100 (reason: DDoS attack detected)
```

---

## 📊 監控微服務

### RabbitMQ 管理界面

訪問 http://localhost:15672
- 用戶名: `pandora`
- 密碼: `pandora123`

查看：
- **Connections**: 3 個連接（每個微服務一個）
- **Queues**: 查看事件隊列
- **Message rates**: 查看消息流量

### 服務日誌

```bash
# 查看所有服務日誌
docker-compose logs -f device-service network-service control-service

# 查看特定服務
docker-compose logs -f device-service
docker-compose logs -f network-service
docker-compose logs -f control-service
```

### 健康檢查

```bash
# 所有服務的健康狀態
for port in 8081 8082 8083; do
  echo "Port $port:"
  curl -s http://localhost:$port/health | jq
done
```

---

## 🔄 服務間通訊流程

### 完整的威脅響應流程

```
1. Device Service 檢測到異常輸入
   └─> 發布 device.data 事件到 RabbitMQ

2. Network Service 監控到攻擊流量
   └─> 發布 network.attack 事件到 RabbitMQ

3. Axiom Engine 接收事件並分析
   └─> 決定需要阻斷攻擊者 IP

4. Engine 調用 Control Service (gRPC)
   └─> BlockIP("192.168.1.100", "DDoS attack", 3600)

5. Control Service 執行阻斷
   ├─> 應用 iptables 規則
   └─> 發布 network.blocked 事件到 RabbitMQ

6. Engine 接收確認事件
   └─> 更新資料庫和儀表板
```

---

## 🧪 測試場景

### 場景 1: 設備連接和數據讀取

```bash
# 1. 連接設備
grpcurl -plaintext \
  -d '{"device_id":"usb-001","port":"/dev/ttyUSB0","baud_rate":115200}' \
  localhost:50051 \
  pandora.device.DeviceService/Connect

# 2. 讀取數據（串流）
grpcurl -plaintext \
  -d '{"device_id":"usb-001","buffer_size":1024}' \
  localhost:50051 \
  pandora.device.DeviceService/ReadData

# 3. 查看 RabbitMQ 中的 device.connected 事件
# 訪問 http://localhost:15672 → Queues → device_events
```

### 場景 2: 網路監控和異常檢測

```bash
# 1. 開始監控
SESSION_ID=$(grpcurl -plaintext \
  -d '{"interface_name":"eth0"}' \
  localhost:50052 \
  pandora.network.NetworkService/StartMonitoring | jq -r '.session_id')

# 2. 獲取統計
grpcurl -plaintext \
  -d "{\"session_id\":\"$SESSION_ID\"}" \
  localhost:50052 \
  pandora.network.NetworkService/GetStatistics

# 3. 檢測異常（串流）
grpcurl -plaintext \
  -d "{\"session_id\":\"$SESSION_ID\",\"threshold\":0.7}" \
  localhost:50052 \
  pandora.network.NetworkService/DetectAnomalies
```

### 場景 3: IP 阻斷和解除

```bash
# 1. 阻斷 IP
grpcurl -plaintext \
  -d '{"ip_address":"192.168.1.100","reason":"Test block","duration_seconds":300}' \
  localhost:50053 \
  pandora.control.ControlService/BlockIP

# 2. 查看阻斷列表
grpcurl -plaintext \
  -d '{}' \
  localhost:50053 \
  pandora.control.ControlService/GetBlockList

# 3. 解除阻斷
grpcurl -plaintext \
  -d '{"ip_address":"192.168.1.100","reason":"Test complete"}' \
  localhost:50053 \
  pandora.control.ControlService/UnblockIP
```

---

## 🔧 故障排除

### 服務無法啟動

```bash
# 檢查日誌
docker-compose logs device-service
docker-compose logs network-service
docker-compose logs control-service

# 檢查端口佔用
netstat -an | grep 50051
netstat -an | grep 50052
netstat -an | grep 50053

# 重啟服務
docker-compose restart device-service
```

### gRPC 連接失敗

```bash
# 檢查服務是否運行
docker-compose ps

# 測試端口連接
telnet localhost 50051
telnet localhost 50052
telnet localhost 50053

# 檢查防火牆規則
sudo iptables -L -n
```

### RabbitMQ 連接失敗

```bash
# 檢查 RabbitMQ 狀態
docker-compose ps rabbitmq

# 查看 RabbitMQ 日誌
docker-compose logs rabbitmq

# 重啟 RabbitMQ
docker-compose restart rabbitmq
```

---

## 📈 性能測試

### 使用 ghz 進行負載測試

```bash
# 安裝 ghz
go install github.com/bojand/ghz/cmd/ghz@latest

# 測試 Device Service
ghz --insecure \
  --proto api/proto/device.proto \
  --call pandora.device.DeviceService/GetStatus \
  -d '{"device_id":"usb-001"}' \
  -n 10000 \
  -c 100 \
  localhost:50051

# 測試 Network Service
ghz --insecure \
  --proto api/proto/network.proto \
  --call pandora.network.NetworkService/GetStatistics \
  -d '{"session_id":"test"}' \
  -n 10000 \
  -c 100 \
  localhost:50052

# 測試 Control Service
ghz --insecure \
  --proto api/proto/control.proto \
  --call pandora.control.ControlService/GetBlockList \
  -d '{}' \
  -n 10000 \
  -c 100 \
  localhost:50053
```

**預期結果**:
- 平均延遲: < 50ms
- 99th 百分位: < 100ms
- 吞吐量: > 1000 req/s

---

## 🎯 下一步

完成微服務啟動後：

1. ✅ **整合到 Axiom Engine** - 讓 Engine 使用微服務
2. ✅ **添加 mTLS** - 為 gRPC 通訊添加 TLS
3. ✅ **服務發現** - 整合 Consul 或 Kubernetes Service Discovery
4. ✅ **監控整合** - 添加 Prometheus 指標
5. ✅ **壓力測試** - 進行完整的性能測試

---

## 📚 相關文檔

- [微服務架構設計](architecture/microservices-design.md)
- [gRPC Proto 定義](../api/proto/README.md)
- [實施路線圖](IMPLEMENTATION-ROADMAP.md)
- [Docker Compose 配置](../deployments/onpremise/docker-compose.yml)

---

**最後更新**: 2025-10-09  
**維護者**: Pandora Box Team

