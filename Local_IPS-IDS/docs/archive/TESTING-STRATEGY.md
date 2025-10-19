# 測試策略與實施計劃
## Phase 4: 從 0.7% 到 80% 測試覆蓋率

> 📅 **創建日期**: 2025-10-09  
> 🎯 **目標**: 80% 測試覆蓋率  
> 📊 **當前狀態**: 0.7% (180 行 / 25,653 行)  
> 🔴 **優先級**: P0 - Critical

---

## 🚨 當前問題

### 測試覆蓋不足

```
當前狀態:
- 總代碼: 25,653 行
- 測試代碼: 180 行
- 覆蓋率: 0.7%

這是嚴重的生產風險！
```

### 專家反饋

根據 `newspec.md` 第 181-204 行：

> "Current status: Only 180 lines of test code for 25,653 lines!  
> This is approximately 0.7% test coverage - CRITICAL GAP"

**優先測試領域**:
- ✅ gRPC service integration tests
- ✅ RabbitMQ message flow tests
- ✅ ML model prediction tests
- ✅ Rate limiting under load
- ✅ Circuit breaker failure scenarios
- ✅ Database transaction rollback tests

---

## 📋 測試策略

### 測試金字塔

```
        ┌─────────────┐
        │  E2E Tests  │  10% (關鍵流程)
        │   (100+)    │
        ├─────────────┤
        │Integration  │  30% (服務間)
        │   Tests     │
        │   (300+)    │
        ├─────────────┤
        │  Unit Tests │  60% (單元邏輯)
        │   (600+)    │
        └─────────────┘
```

**總計**: 1000+ 測試案例

---

## 🧪 測試類型詳細計劃

### 1. 單元測試（600+ 測試）

**目標**: 60% 覆蓋率

#### internal/pubsub/ (50 tests)
```go
// tests/unit/pubsub/rabbitmq_test.go
func TestRabbitMQPublish(t *testing.T)
func TestRabbitMQSubscribe(t *testing.T)
func TestRabbitMQReconnect(t *testing.T)
func TestRabbitMQMessagePersistence(t *testing.T)
func TestEventSerialization(t *testing.T)
// ... 45 more tests
```

#### internal/services/ (150 tests)
```go
// tests/unit/services/device/service_test.go
func TestDeviceServiceGetStatus(t *testing.T)
func TestDeviceServiceListDevices(t *testing.T)
func TestSerialCommunication(t *testing.T)
// ... 50 tests per service × 3 services
```

#### internal/ml/ (100 tests)
```go
// tests/unit/ml/bot_detector_test.go
func TestBotDetectorPredict(t *testing.T)
func TestFeatureExtraction(t *testing.T)
func TestScoreCalculation(t *testing.T)

// tests/unit/ml/deep_learning_test.go
func TestNeuralNetworkForward(t *testing.T)
func TestModelTraining(t *testing.T)
func TestFeatureNormalization(t *testing.T)

// tests/unit/ml/behavior_baseline_test.go
func TestAnomalyDetection(t *testing.T)
func TestBaselineComputation(t *testing.T)
func TestDeviationCalculation(t *testing.T)
```

#### internal/security/ (80 tests)
```go
// tests/unit/security/waf_test.go
func TestWAFInspectRequest(t *testing.T)
func TestSQLInjectionDetection(t *testing.T)
func TestXSSDetection(t *testing.T)

// tests/unit/security/tls_fingerprint_test.go
func TestJA3Generation(t *testing.T)
func TestBotIdentification(t *testing.T)
func TestMalwareDetection(t *testing.T)
```

#### internal/ratelimit/ (40 tests)
```go
// tests/unit/ratelimit/token_bucket_test.go
func TestTokenBucketAllow(t *testing.T)
func TestTokenBucketRefill(t *testing.T)
func TestMultiLevelRateLimit(t *testing.T)
```

#### internal/resilience/ (60 tests)
```go
// tests/unit/resilience/retry_test.go
func TestExponentialBackoff(t *testing.T)
func TestMaxRetries(t *testing.T)

// tests/unit/resilience/circuit_breaker_test.go
func TestCircuitBreakerOpen(t *testing.T)
func TestCircuitBreakerHalfOpen(t *testing.T)
func TestCircuitBreakerClose(t *testing.T)
```

#### 其他模組 (120 tests)
- internal/cache/ (30 tests)
- internal/tracing/ (20 tests)
- internal/multitenant/ (30 tests)
- internal/automation/ (40 tests)

### 2. 集成測試（300+ 測試）

**目標**: 30% 覆蓋率

#### gRPC 服務集成 (60 tests)
```go
// tests/integration/grpc/device_service_test.go
func TestDeviceServiceIntegration(t *testing.T) {
    // 啟動真實 gRPC 服務
    // 測試所有 6 個 RPCs
    // 驗證 mTLS 連接
}

// tests/integration/grpc/network_service_test.go
func TestNetworkServiceIntegration(t *testing.T) {
    // 測試所有 7 個 RPCs
}

// tests/integration/grpc/control_service_test.go
func TestControlServiceIntegration(t *testing.T) {
    // 測試所有 9 個 RPCs
}
```

#### RabbitMQ 消息流 (40 tests)
```go
// tests/integration/rabbitmq/message_flow_test.go
func TestThreatEventFlow(t *testing.T)
func TestNetworkEventFlow(t *testing.T)
func TestSystemEventFlow(t *testing.T)
func TestDeviceEventFlow(t *testing.T)
func TestMessagePersistence(t *testing.T)
func TestMessageOrdering(t *testing.T)
```

#### 資料庫集成 (60 tests)
```go
// tests/integration/database/postgresql_test.go
func TestEventStorage(t *testing.T)
func TestTransactionRollback(t *testing.T)
func TestConcurrentWrites(t *testing.T)
func TestQueryPerformance(t *testing.T)
```

#### Redis 集成 (40 tests)
```go
// tests/integration/redis/cache_test.go
func TestCacheHitRate(t *testing.T)
func TestCacheEviction(t *testing.T)
func TestDistributedLock(t *testing.T)
```

#### 端到端流程 (100 tests)
```go
// tests/integration/e2e/threat_detection_test.go
func TestDDoSDetectionFlow(t *testing.T) {
    // 1. 模擬 DDoS 流量
    // 2. Network Service 檢測
    // 3. RabbitMQ 傳遞事件
    // 4. ML 模型分析
    // 5. SOAR 自動響應
    // 6. 驗證 IP 被阻斷
}
```

### 3. 端到端測試（100+ 測試）

**目標**: 10% 覆蓋率

#### 關鍵業務流程 (50 tests)
```go
// tests/e2e/business_flows_test.go
func TestCompleteSecurityIncidentFlow(t *testing.T)
func TestUserOnboardingFlow(t *testing.T)
func TestThreatResponsePlaybookFlow(t *testing.T)
func TestMultiTenantIsolationFlow(t *testing.T)
```

#### 用戶場景 (50 tests)
```go
// tests/e2e/user_scenarios_test.go
func TestAdminBlocksIPScenario(t *testing.T)
func TestAnalystInvestigatesThreatScenario(t *testing.T)
func TestSOCRespondsToAlertScenario(t *testing.T)
```

---

## 🚀 性能測試計劃

### 負載測試（k6）

```javascript
// tests/load/throughput_test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 },    // Ramp up
    { duration: '5m', target: 100000 }, // Sustained load
    { duration: '2m', target: 0 },      // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(99)<2'],     // 99% < 2ms
    http_req_failed: ['rate<0.01'],     // < 1% errors
  },
};

export default function () {
  let res = http.get('http://localhost:8081/health');
  check(res, { 'status is 200': (r) => r.status === 200 });
  sleep(0.01);
}
```

### 壓力測試（Apache JMeter）

```xml
<!-- tests/load/stress_test.jmx -->
<jmeterTestPlan>
  <hashTree>
    <ThreadGroup>
      <stringProp name="ThreadGroup.num_threads">10000</stringProp>
      <stringProp name="ThreadGroup.ramp_time">60</stringProp>
      <stringProp name="ThreadGroup.duration">600</stringProp>
    </ThreadGroup>
  </hashTree>
</jmeterTestPlan>
```

### 基準測試（Go Benchmark）

```go
// tests/benchmarks/microservices_bench_test.go
func BenchmarkDeviceServiceGetStatus(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Call gRPC
    }
}

func BenchmarkMLPrediction(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Run prediction
    }
}

func BenchmarkCacheGet(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Cache lookup
    }
}
```

---

## 🔥 混沌工程測試

### Chaos Mesh 配置

```yaml
# tests/chaos/pod-failure.yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: PodChaos
metadata:
  name: device-service-failure
  namespace: pandora-system
spec:
  action: pod-failure
  mode: one
  duration: "30s"
  selector:
    namespaces:
      - pandora-system
    labelSelectors:
      app: device-service
```

```yaml
# tests/chaos/network-delay.yaml
apiVersion: chaos-mesh.org/v1alpha1
kind: NetworkChaos
metadata:
  name: network-delay
spec:
  action: delay
  mode: all
  selector:
    namespaces:
      - pandora-system
  delay:
    latency: "100ms"
    correlation: "25"
    jitter: "10ms"
  duration: "5m"
```

### 測試場景

| 場景 | 預期結果 | 驗證方法 |
|------|----------|----------|
| Pod 隨機終止 | 自動重啟，服務不中斷 | 監控 Prometheus |
| 網路延遲 100ms | 重試機制生效 | 檢查日誌 |
| CPU 100% | HPA 擴展到 10 副本 | kubectl get hpa |
| 記憶體洩漏 | OOM Kill + 重啟 | kubectl describe pod |
| PostgreSQL 斷線 | 斷路器開啟 | 檢查 circuit breaker 狀態 |
| RabbitMQ 故障 | 消息緩存 + 重連 | 檢查消息隊列 |

---

## 🛡️ 安全測試計劃

### 滲透測試（OWASP ZAP）

```bash
# tests/security/penetration_test.sh

# 1. SQL 注入測試
zap-cli quick-scan --spider http://localhost:3001
zap-cli active-scan http://localhost:3001/api/v1/*

# 2. XSS 測試
zap-cli alerts --alert-type XSS

# 3. CSRF 測試
zap-cli alerts --alert-type CSRF

# 4. 認證測試
zap-cli alerts --alert-type "Broken Authentication"
```

### 漏洞掃描（Trivy）

```bash
# tests/security/vulnerability_scan.sh

# 掃描 Docker 鏡像
trivy image pandora-box/device-service:latest
trivy image pandora-box/network-service:latest
trivy image pandora-box/control-service:latest

# 掃描依賴
trivy fs --security-checks vuln,config .

# 生成報告
trivy image --format json --output trivy-report.json pandora-box/device-service:latest
```

### mTLS 驗證（testssl.sh）

```bash
# tests/security/mtls_validation.sh

# 測試 TLS 配置
testssl.sh --full localhost:50051
testssl.sh --full localhost:50052
testssl.sh --full localhost:50053

# 驗證證書鏈
openssl s_client -connect localhost:50051 -CAfile certs/ca.crt

# 測試證書輪換
./scripts/rotate-certs.sh --dry-run
```

---

## 📊 測試覆蓋率目標

### 按模組

| 模組 | 當前 | 目標 | 測試數 |
|------|------|------|--------|
| internal/pubsub/ | 10% | 85% | 50 |
| internal/services/ | 0% | 80% | 150 |
| internal/ml/ | 0% | 75% | 100 |
| internal/security/ | 0% | 90% | 80 |
| internal/grpc/ | 0% | 85% | 40 |
| internal/resilience/ | 0% | 90% | 60 |
| internal/ratelimit/ | 0% | 85% | 40 |
| internal/cache/ | 0% | 80% | 30 |
| internal/tracing/ | 0% | 70% | 20 |
| internal/multitenant/ | 0% | 75% | 30 |
| internal/automation/ | 0% | 70% | 40 |
| cmd/*-service/ | 0% | 60% | 30 |
| **總計** | **0.7%** | **80%** | **670** |

### 按測試類型

| 類型 | 測試數 | 覆蓋率 |
|------|--------|--------|
| 單元測試 | 600 | 60% |
| 集成測試 | 300 | 30% |
| E2E 測試 | 100 | 10% |
| **總計** | **1000** | **100%** |

---

## 🎯 實施計劃

### Week 1: 基礎測試框架

**Day 1-2**: 設置測試環境
```bash
# 安裝測試工具
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install github.com/onsi/gomega/...@latest
go install github.com/golang/mock/mockgen@latest

# 創建測試目錄結構
mkdir -p tests/{unit,integration,e2e,load,chaos,security}
```

**Day 3-4**: 創建測試輔助工具
```go
// tests/testutil/helpers.go
package testutil

func SetupTestDB(t *testing.T) *sql.DB
func SetupTestRedis(t *testing.T) *redis.Client
func SetupTestRabbitMQ(t *testing.T) *amqp.Connection
func MockGRPCServer(t *testing.T) *grpc.Server
```

**Day 5**: 配置 CI 測試流程
```yaml
# .github/workflows/test-suite.yml
name: Comprehensive Test Suite
on: [push, pull_request]
jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go test -cover -coverprofile=coverage.out ./...
      - run: go tool cover -html=coverage.out -o coverage.html
      - uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: coverage.html
```

### Week 2: 單元測試實施

**目標**: 完成 600 個單元測試

**每日目標**: 60 個測試/天 × 10 天

**Day 1**: internal/pubsub/ (50 tests)  
**Day 2**: internal/services/device/ (50 tests)  
**Day 3**: internal/services/network/ (50 tests)  
**Day 4**: internal/services/control/ (50 tests)  
**Day 5**: internal/ml/bot_detector (50 tests)  
**Day 6**: internal/ml/deep_learning (50 tests)  
**Day 7**: internal/security/ (80 tests)  
**Day 8**: internal/resilience/ + ratelimit/ (100 tests)  
**Day 9**: internal/cache/ + tracing/ (50 tests)  
**Day 10**: internal/multitenant/ + automation/ (70 tests)

### Week 3: 集成測試實施

**目標**: 完成 300 個集成測試

**Day 1-2**: gRPC 服務集成 (60 tests)  
**Day 3**: RabbitMQ 消息流 (40 tests)  
**Day 4-5**: 資料庫集成 (60 tests)  
**Day 6**: Redis 集成 (40 tests)  
**Day 7-10**: 端到端流程 (100 tests)

### Week 4: E2E 和性能測試

**Day 1-2**: E2E 測試 (100 tests)  
**Day 3-4**: 負載測試（k6）  
**Day 5**: 壓力測試（JMeter）  
**Day 6-7**: 基準測試（Go Benchmark）  
**Day 8-10**: 性能報告生成

---

## 🔧 測試工具清單

### Go 測試工具

```bash
# 安裝所有測試工具
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install github.com/onsi/gomega/...@latest
go install github.com/golang/mock/mockgen@latest
go install github.com/stretchr/testify@latest
go install gotest.tools/gotestsum@latest
go install github.com/axw/gocov/gocov@latest
go install github.com/AlekSi/gocov-xml@latest
```

### 性能測試工具

```bash
# k6
brew install k6

# Apache JMeter
brew install jmeter

# Gatling
brew install gatling
```

### 安全測試工具

```bash
# OWASP ZAP
brew install zaproxy

# Trivy
brew install aquasecurity/trivy/trivy

# gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest
```

### 混沌工程工具

```bash
# Chaos Mesh
kubectl apply -f https://mirrors.chaos-mesh.org/latest/crd.yaml
kubectl apply -f https://mirrors.chaos-mesh.org/latest/chaos-mesh.yaml

# Litmus
kubectl apply -f https://litmuschaos.github.io/litmus/litmus-operator-latest.yaml
```

---

## 📈 測試報告

### 生成覆蓋率報告

```bash
# 運行所有測試並生成報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 生成 JSON 報告
gocov convert coverage.out | gocov-xml > coverage.xml

# 上傳到 Codecov
bash <(curl -s https://codecov.io/bash)
```

### CI 集成

```yaml
# .github/workflows/test-suite.yml 片段
- name: Run tests with coverage
  run: |
    go test -v -coverprofile=coverage.out -covermode=atomic ./...
    
- name: Upload coverage to Codecov
  uses: codecov/codecov-action@v3
  with:
    files: ./coverage.out
    flags: unittests
    name: codecov-umbrella
```

---

## ✅ 驗收標準

### 必須達成（P0）

- [ ] 單元測試覆蓋率 ≥ 60%
- [ ] 集成測試覆蓋率 ≥ 30%
- [ ] E2E 測試覆蓋率 ≥ 10%
- [ ] 總體覆蓋率 ≥ 80%
- [ ] 所有測試通過
- [ ] 無高危安全漏洞
- [ ] 性能指標驗證通過

### 應該達成（P1）

- [ ] 測試執行時間 < 10 分鐘
- [ ] 測試穩定性 > 99%
- [ ] 測試文檔完整
- [ ] CI/CD 集成完成

### 可以達成（P2）

- [ ] 測試覆蓋率 ≥ 90%
- [ ] 變異測試（Mutation Testing）
- [ ] 屬性測試（Property-Based Testing）
- [ ] 模糊測試（Fuzzing）

---

## 📚 參考資源

- [Go Testing Best Practices](https://golang.org/doc/code.html#Testing)
- [k6 Load Testing](https://k6.io/docs/)
- [Chaos Mesh Documentation](https://chaos-mesh.org/docs/)
- [OWASP Testing Guide](https://owasp.org/www-project-web-security-testing-guide/)

---

**狀態**: 📅 計劃中  
**優先級**: 🔴 P0 - Critical  
**預計完成**: Week 1-4（測試實施）  
**負責人**: QA Team + DevOps

**🚀 讓我們把測試覆蓋率提升到世界級標準！**

