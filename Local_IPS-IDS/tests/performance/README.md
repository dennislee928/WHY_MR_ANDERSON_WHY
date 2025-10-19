# 性能測試套件
## Pandora Box Console IDS-IPS Microservices

> 🚀 完整的性能測試和壓力測試套件

---

## 📋 測試類型

### 1. 基準測試（Benchmark）

測試單個 RPC 的性能：

```bash
# 運行所有基準測試
go test -bench=. -benchmem -benchtime=10s

# 運行特定服務的測試
go test -bench=BenchmarkDeviceService -benchmem
go test -bench=BenchmarkNetworkService -benchmem
go test -bench=BenchmarkControlService -benchmem
```

### 2. 負載測試（Load Test）

測試系統在高負載下的表現：

```bash
# 運行負載測試（30 秒，100 並發）
go test -v -tags=performance -run TestLoadTest

# 長時間負載測試（5 分鐘）
go test -v -tags=performance -run TestLoadTest -timeout 10m
```

### 3. 壓力測試（Stress Test）

使用 ghz 進行壓力測試：

```bash
# 安裝 ghz
go install github.com/bojand/ghz/cmd/ghz@latest

# Device Service 壓力測試
ghz --insecure \
  --proto ../../api/proto/device.proto \
  --call pandora.device.DeviceService/GetStatus \
  -d '{"device_id":"usb-001"}' \
  -n 100000 \
  -c 500 \
  --connections=50 \
  localhost:50051

# Network Service 壓力測試
ghz --insecure \
  --proto ../../api/proto/network.proto \
  --call pandora.network.NetworkService/GetStatistics \
  -d '{"session_id":"test"}' \
  -n 100000 \
  -c 500 \
  localhost:50052

# Control Service 壓力測試
ghz --insecure \
  --proto ../../api/proto/control.proto \
  --call pandora.control.ControlService/GetBlockList \
  -d '{}' \
  -n 100000 \
  -c 500 \
  localhost:50053
```

---

## 🎯 性能目標

| 指標 | 目標 | 說明 |
|------|------|------|
| 平均延遲 | < 50ms | 單個 RPC 調用 |
| P99 延遲 | < 100ms | 99th 百分位 |
| 吞吐量 | > 1000 req/s | 每個服務 |
| 並發連接 | > 1000 | 同時連接數 |
| 錯誤率 | < 0.1% | 失敗請求比例 |
| CPU 使用率 | < 70% | 正常負載下 |
| 內存使用 | < 200MB | 每個服務 |

---

## 📊 測試結果範例

### Device Service

```
BenchmarkDeviceServiceGetStatus-8        50000    25000 ns/op    1024 B/op    10 allocs/op

Summary:
  Count:        50000
  Total:        1.250 s
  Slowest:      45.23 ms
  Fastest:      2.15 ms
  Average:      25.00 ms
  Requests/sec: 40000.00

Latency distribution:
  10% in 15.00 ms
  25% in 20.00 ms
  50% in 25.00 ms
  75% in 30.00 ms
  90% in 35.00 ms
  95% in 40.00 ms
  99% in 45.00 ms

Status code distribution:
  [OK]   50000 responses
```

### Network Service

```
BenchmarkNetworkServiceGetStatistics-8   40000    30000 ns/op    2048 B/op    15 allocs/op

Summary:
  Count:        40000
  Total:        1.200 s
  Slowest:      55.12 ms
  Fastest:      5.23 ms
  Average:      30.00 ms
  Requests/sec: 33333.33
```

### Control Service

```
BenchmarkControlServiceGetBlockList-8    60000    20000 ns/op    512 B/op     8 allocs/op

Summary:
  Count:        60000
  Total:        1.200 s
  Slowest:      35.45 ms
  Fastest:      8.12 ms
  Average:      20.00 ms
  Requests/sec: 50000.00
```

---

## 🔧 性能優化建議

### 1. 連接池優化

```go
// 增加連接池大小
grpc.WithDefaultCallOptions(
    grpc.MaxCallRecvMsgSize(10 * 1024 * 1024),
    grpc.MaxCallSendMsgSize(10 * 1024 * 1024),
)

// 調整 keepalive 參數
grpc.WithKeepaliveParams(keepalive.ClientParameters{
    Time:                10 * time.Second,
    Timeout:             3 * time.Second,
    PermitWithoutStream: true,
})
```

### 2. 並發控制

```go
// 使用 worker pool 限制並發
type WorkerPool struct {
    workers   int
    taskQueue chan func()
}

func NewWorkerPool(workers int) *WorkerPool {
    wp := &WorkerPool{
        workers:   workers,
        taskQueue: make(chan func(), workers*2),
    }
    wp.start()
    return wp
}
```

### 3. 緩存優化

```go
// 使用 Redis 緩存頻繁查詢的數據
type Cache struct {
    redis *redis.Client
    ttl   time.Duration
}

func (c *Cache) Get(key string) ([]byte, error) {
    return c.redis.Get(context.Background(), key).Bytes()
}
```

### 4. 批量處理

```go
// 批量發布事件
func (p *EventPublisher) PublishBatch(ctx context.Context, events []interface{}) error {
    for _, event := range events {
        // 發布事件
    }
    return nil
}
```

---

## 📈 監控指標

### gRPC 指標

```go
import "github.com/grpc-ecosystem/go-grpc-prometheus"

// 添加 Prometheus 攔截器
grpcServer := grpc.NewServer(
    grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
    grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
)

// 註冊指標
grpc_prometheus.Register(grpcServer)
```

### 關鍵指標

- `grpc_server_handled_total` - 處理的請求總數
- `grpc_server_handling_seconds` - 請求處理時間
- `grpc_server_msg_received_total` - 接收的消息數
- `grpc_server_msg_sent_total` - 發送的消息數

---

## 🧪 運行測試

### 前置條件

```bash
# 啟動所有微服務
cd deployments/onpremise
docker-compose up -d rabbitmq device-service network-service control-service

# 等待服務就緒
sleep 10

# 驗證服務健康
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

### 運行基準測試

```bash
cd tests/performance

# 快速測試
go test -bench=. -benchtime=1s

# 完整測試
go test -bench=. -benchmem -benchtime=10s > bench_results.txt

# 比較結果
benchstat bench_results_old.txt bench_results.txt
```

### 運行負載測試

```bash
# 30 秒負載測試
go test -v -tags=performance -run TestLoadTest

# 5 分鐘壓力測試
go test -v -tags=performance -run TestLoadTest -timeout 10m
```

### 使用 ghz 壓力測試

```bash
# 創建測試腳本
cat > stress_test.sh <<'EOF'
#!/bin/bash

echo "Starting stress test..."

# Device Service
ghz --insecure \
  --proto ../../api/proto/device.proto \
  --call pandora.device.DeviceService/GetStatus \
  -d '{"device_id":"usb-001"}' \
  -n 100000 \
  -c 500 \
  --connections=50 \
  -o device_results.json \
  localhost:50051 &

# Network Service
ghz --insecure \
  --proto ../../api/proto/network.proto \
  --call pandora.network.NetworkService/GetStatistics \
  -d '{"session_id":"test"}' \
  -n 100000 \
  -c 500 \
  --connections=50 \
  -o network_results.json \
  localhost:50052 &

# Control Service
ghz --insecure \
  --proto ../../api/proto/control.proto \
  --call pandora.control.ControlService/GetBlockList \
  -d '{}' \
  -n 100000 \
  -c 500 \
  --connections=50 \
  -o control_results.json \
  localhost:50053 &

wait

echo "Stress test completed!"
echo "Results saved to *_results.json"
EOF

chmod +x stress_test.sh
./stress_test.sh
```

---

## 📊 分析結果

### 使用 ghz 結果

```bash
# 查看 JSON 結果
cat device_results.json | jq '.average, .fastest, .slowest, .rps'

# 生成報告
ghz --insecure \
  --proto ../../api/proto/device.proto \
  --call pandora.device.DeviceService/GetStatus \
  -d '{"device_id":"usb-001"}' \
  -n 10000 \
  -c 100 \
  --format html \
  -O device_report.html \
  localhost:50051
```

### 使用 pprof 分析

```go
import _ "net/http/pprof"

// 在服務中添加 pprof 端點
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

分析：
```bash
# CPU profile
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Memory profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine profile
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

---

## 🎯 性能調優檢查清單

- [ ] gRPC 連接池配置
- [ ] Keepalive 參數調整
- [ ] 消息大小限制
- [ ] 並發控制
- [ ] 緩存策略
- [ ] 數據庫連接池
- [ ] RabbitMQ 預取數量
- [ ] 日誌級別（生產環境使用 WARN）
- [ ] Goroutine 洩漏檢查
- [ ] 內存洩漏檢查

---

## 📚 相關文檔

- [微服務架構設計](../../docs/architecture/microservices-design.md)
- [實施路線圖](../../docs/IMPLEMENTATION-ROADMAP.md)
- [Week 2 完成報告](../../docs/PHASE1-WEEK2-COMPLETE.md)

---

**維護者**: Pandora Box Team  
**最後更新**: 2025-10-09

