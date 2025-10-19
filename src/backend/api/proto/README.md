# gRPC Proto Definitions
## Pandora Box Console IDS-IPS Microservices

這個目錄包含所有微服務的 gRPC proto 定義文件。

---

## 📦 Proto 文件

| 文件 | 服務 | 說明 |
|------|------|------|
| `common.proto` | - | 共享類型和枚舉 |
| `device.proto` | DeviceService | 設備管理服務 |
| `network.proto` | NetworkService | 網路監控服務 |
| `control.proto` | ControlService | 網路控制服務 |

---

## 🚀 快速開始

### 1. 安裝工具

```bash
# 安裝 protoc（Protocol Buffer 編譯器）
# macOS
brew install protobuf

# Ubuntu/Debian
sudo apt-get install -y protobuf-compiler

# Windows
choco install protoc

# 安裝 Go plugins
make install
```

### 2. 生成代碼

```bash
# 生成所有服務的代碼
make generate

# 或者單獨生成
make device
make network
make control
```

### 3. 驗證

```bash
# 驗證 proto 文件語法
make validate
```

---

## 📚 服務說明

### Device Service

**端口**: 50051  
**職責**: 管理 USB-SERIAL CH340 設備

**主要 RPC**:
- `Connect` - 連接設備
- `ReadData` - 讀取設備數據（串流）
- `GetStatus` - 獲取設備狀態
- `ListDevices` - 列出所有設備

### Network Service

**端口**: 50052  
**職責**: 監控網路流量和檢測異常

**主要 RPC**:
- `StartMonitoring` - 開始監控
- `GetStatistics` - 獲取統計數據
- `AnalyzeTraffic` - 分析流量（串流）
- `DetectAnomalies` - 檢測異常（串流）

### Control Service

**端口**: 50053  
**職責**: 執行網路控制和阻斷

**主要 RPC**:
- `BlockIP` - 阻斷 IP
- `BlockPort` - 阻斷端口
- `ApplyFirewallRule` - 應用防火牆規則
- `GetBlockList` - 獲取阻斷列表

---

## 🔧 使用範例

### 調用 Device Service

```go
import (
    pb "pandora_box_console_ids_ips/api/proto/device"
    "google.golang.org/grpc"
)

// 創建連接
conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// 創建客戶端
client := pb.NewDeviceServiceClient(conn)

// 調用 RPC
req := &pb.ConnectRequest{
    DeviceId: "usb-001",
    Port: "/dev/ttyUSB0",
    BaudRate: 115200,
}
resp, err := client.Connect(context.Background(), req)
```

### 調用 Network Service

```go
import pb "pandora_box_console_ids_ips/api/proto/network"

client := pb.NewNetworkServiceClient(conn)

// 開始監控
req := &pb.MonitorRequest{
    InterfaceName: "eth0",
    PromiscuousMode: true,
}
resp, err := client.StartMonitoring(context.Background(), req)
```

### 調用 Control Service

```go
import pb "pandora_box_console_ids_ips/api/proto/control"

client := pb.NewControlServiceClient(conn)

// 阻斷 IP
req := &pb.BlockIPRequest{
    IpAddress: "192.168.1.100",
    Reason: "DDoS attack detected",
    DurationSeconds: 3600,
    Action: pb.BlockAction_BLOCK_ACTION_DROP,
}
resp, err := client.BlockIP(context.Background(), req)
```

---

## 📝 開發指南

### 添加新的 RPC

1. 編輯對應的 `.proto` 文件
2. 添加新的 RPC 定義
3. 運行 `make generate`
4. 實現服務端邏輯
5. 添加測試

### 修改現有 RPC

1. 更新 `.proto` 文件
2. 運行 `make generate`
3. 更新服務端實現
4. 更新客戶端調用
5. 更新測試

### 版本管理

使用 proto 包版本控制：

```protobuf
package pandora.device.v1;
package pandora.device.v2;
```

---

## 🧪 測試

### 使用 grpcurl 測試

```bash
# 安裝 grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# 列出服務
grpcurl -plaintext localhost:50051 list

# 列出方法
grpcurl -plaintext localhost:50051 list pandora.device.DeviceService

# 調用方法
grpcurl -plaintext \
  -d '{"device_id":"usb-001"}' \
  localhost:50051 \
  pandora.device.DeviceService/GetStatus
```

### 使用 ghz 性能測試

```bash
# 安裝 ghz
go install github.com/bojand/ghz/cmd/ghz@latest

# 性能測試
ghz --insecure \
  --proto device.proto \
  --call pandora.device.DeviceService/GetStatus \
  -d '{"device_id":"usb-001"}' \
  -n 10000 \
  -c 100 \
  localhost:50051
```

---

## 📚 相關資源

- [gRPC 官方文檔](https://grpc.io/docs/)
- [Protocol Buffers 文檔](https://protobuf.dev/)
- [微服務架構設計](../architecture/microservices-design.md)
- [實施路線圖](../IMPLEMENTATION-ROADMAP.md)

---

**維護者**: Pandora Box Team  
**最後更新**: 2025-10-09  
**版本**: 1.0.0

