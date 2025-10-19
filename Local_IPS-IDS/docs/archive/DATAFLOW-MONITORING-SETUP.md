# 資料流監控設定指南

> **版本**: 3.2.0  
> **服務**: Cyber AI/Quantum Security  
> **完成日期**: 2025-01-14

---

## 📋 概述

本文檔說明如何設定和使用 Pandora Box Console IDS-IPS 的 AI 資料流監控系統。

---

## 🏗️ 架構

```
網路流量 → 封包捕獲 → 特徵提取 → 異常檢測 → 告警 → 響應
    ↓           ↓           ↓           ↓         ↓       ↓
 Promtail    DataFlow    Analyzer    ML Model  RabbitMQ  Axiom
             Monitor                                      Engine
```

---

## 🚀 快速開始

### 啟動監控服務

```bash
# 啟動 Cyber AI/Quantum 容器
cd Application
docker-compose up -d cyber-ai-quantum

# 驗證服務
curl http://localhost:8000/health
```

### 查看資料流統計

```bash
curl http://localhost:8000/api/v1/dataflow/stats
```

**響應範例**:

```json
{
  "status": "success",
  "metrics": {
    "total_bytes": 1548900,
    "packet_count": 1234,
    "unique_sources": 45,
    "unique_destinations": 12,
    "protocols": {
      "TCP": 890,
      "UDP": 300,
      "ICMP": 44
    },
    "avg_packet_size": 1255.3,
    "flow_rate": 25815.0,
    "timestamp": "2025-01-14T12:00:00"
  }
}
```

---

## 📊 行為基線

### 建立基線

系統自動在7天內建立用戶和網路行為基線。

#### 基線特徵

- **流量率**: 平均流量和標準差
- **封包數**: 平均封包數和標準差
- **封包大小**: 平均大小和標準差
- **協議分布**: 各協議使用比例
- **時間模式**: 高峰和低峰時段

#### 基線公式

```python
baseline = {
    'mean': np.mean(historical_data),
    'std': np.std(historical_data),
    'min': np.min(historical_data),
    'max': np.max(historical_data),
    'median': np.median(historical_data)
}
```

### 查看基線

```bash
curl http://localhost:8000/api/v1/dataflow/baseline
```

**響應**:

```json
{
  "status": "success",
  "baseline": {
    "flow_rate": {
      "mean": 25000.5,
      "std": 5000.2
    },
    "packet_count": {
      "mean": 1200.3,
      "std": 300.8
    },
    "avg_packet_size": {
      "mean": 1250.0,
      "std": 250.5
    }
  },
  "last_updated": "2025-01-14T12:00:00"
}
```

---

## 🚨 異常檢測

### Z-Score 方法

#### 原理

```python
z_score = (current_value - baseline_mean) / baseline_std

if abs(z_score) > threshold:  # threshold = 3.0
    trigger_anomaly_alert()
```

#### 閾值設定

| Z-Score | 嚴重程度 | 建議操作 |
|---------|---------|---------|
| > 5.0 | CRITICAL | 立即阻斷並調查 |
| > 4.0 | HIGH | 告警並監控 |
| > 3.0 | MEDIUM | 記錄並觀察 |
| < 3.0 | - | 正常 |

### 異常類型

#### 流量異常

- **症狀**: 流量率顯著高於基線
- **可能原因**: DDoS、數據外洩
- **建議**: 分析來源 IP，考慮限流

#### 封包數異常

- **症狀**: 封包數激增
- **可能原因**: 掃描攻擊、Botnet
- **建議**: 檢查封包模式，可能阻斷

#### 封包大小異常

- **症狀**: 平均封包大小異常
- **可能原因**: 數據滲漏、異常協議
- **建議**: 深度封包檢測

### 查看異常

```bash
curl http://localhost:8000/api/v1/dataflow/anomalies
```

**響應**:

```json
{
  "status": "success",
  "anomalies": [
    {
      "alert_id": "anomaly_20250114120000",
      "anomaly_type": "flow_rate, packet_count",
      "severity": "high",
      "anomaly_score": 4.2,
      "timestamp": "2025-01-14T12:00:00"
    }
  ],
  "total": 15
}
```

---

## ⚙️ 配置

### 環境變數

```bash
# 監控窗口大小（秒）
DATAFLOW_WINDOW_SIZE=60

# 異常檢測閾值（Z-Score）
ANOMALY_THRESHOLD=3.0

# 基線更新頻率（小時）
BASELINE_UPDATE_INTERVAL=24

# 歷史數據保留（小時）
HISTORY_RETENTION=72
```

### 調整閾值

```python
# 調整異常檢測靈敏度
monitor.anomaly_detector.threshold = 2.5  # 更敏感
monitor.anomaly_detector.threshold = 4.0  # 更寬鬆
```

---

## 🔗 整合

### RabbitMQ 整合

```python
# 訂閱網路事件
await mq.subscribe(
    queue="network_events",
    handler=process_network_event
)

async def process_network_event(routing_key, message):
    event = json.loads(message)
    
    # 添加到資料流監控
    monitor.analyzer.process_packet(event)
    
    # 檢測異常
    metrics = monitor.analyzer.calculate_metrics()
    if metrics:
        anomaly = monitor.anomaly_detector.detect_anomaly(metrics)
        
        if anomaly:
            # 發布異常告警
            await mq.publish(
                exchange="pandora.events",
                routing_key="dataflow.anomaly",
                message=json.dumps(asdict(anomaly))
            )
```

### Grafana 可視化

```yaml
# Grafana Dashboard JSON
{
  "title": "資料流監控",
  "panels": [
    {
      "title": "流量趨勢",
      "targets": [{
        "expr": "dataflow_bytes_total"
      }]
    },
    {
      "title": "異常告警",
      "targets": [{
        "expr": "dataflow_anomalies_total"
      }]
    }
  ]
}
```

---

## 📈 性能優化

### 滑動窗口

```python
# 使用固定大小滑動窗口
window_size = 60  # 秒
cutoff_time = now - timedelta(seconds=window_size)

# 清理過期數據
flow_buffer = [p for p in flow_buffer if p['time'] > cutoff_time]
```

### 批次處理

```python
# 批次計算指標，減少 CPU 使用
batch_size = 100
for i in range(0, len(packets), batch_size):
    batch = packets[i:i+batch_size]
    process_batch(batch)
```

---

## 🧪 測試

### 正常流量測試

```python
# 生成正常流量
for i in range(100):
    packet = {
        'source_ip': f'192.168.1.{i % 20}',
        'size': random.randint(64, 1500),
        'protocol': 'TCP'
    }
    send_packet(packet)

# 驗證：應該無異常告警
```

### 異常流量測試

```python
# 生成 DDoS 流量
for i in range(5000):
    packet = {
        'source_ip': '192.168.1.100',
        'size': 64,
        'protocol': 'TCP'
    }
    send_packet(packet)

# 驗證：應該觸發異常告警
```

---

## 🎯 使用案例

### 案例 1: DDoS 檢測

**場景**: 檢測分散式拒絕服務攻擊

**指標異常**:
- flow_rate: 10x 基線
- packet_count: 15x 基線
- unique_sources: 100+ IPs

**系統響應**:
1. 觸發 CRITICAL 異常告警
2. 自動發送到 RabbitMQ
3. Axiom Engine 啟動阻斷
4. 記錄到 PostgreSQL

### 案例 2: 數據外洩檢測

**場景**: 檢測大量數據外傳

**指標異常**:
- flow_rate: 50x 基線（向外）
- avg_packet_size: 大封包
- 時間: 非工作時間

**系統響應**:
1. 觸發 HIGH 異常告警
2. 標記可疑 IP
3. 人工審核
4. 可能阻斷

### 案例 3: 內部威脅檢測

**場景**: 檢測內部人員異常行為

**指標異常**:
- 存取時間異常（深夜）
- 數據存取量異常（10x）
- 命令執行異常（sudo）

**系統響應**:
1. 觸發 MEDIUM 異常告警
2. 記錄行為模式
3. 通知安全團隊
4. 增強監控

---

## 📚 故障排除

### 問題：基線不準確

**原因**: 學習期數據不足或異常

**解決**:
```bash
# 重置基線
curl -X POST http://localhost:8000/api/v1/dataflow/baseline/reset

# 等待7天重新學習
```

### 問題：誤報率高

**原因**: 閾值設定過於敏感

**解決**:
```python
# 調整閾值
ANOMALY_THRESHOLD=4.0  # 從 3.0 提高到 4.0
```

### 問題：漏報威脅

**原因**: 閾值設定過於寬鬆

**解決**:
```python
# 調整閾值
ANOMALY_THRESHOLD=2.5  # 從 3.0 降低到 2.5

# 或增加特徵維度
```

---

**維護者**: Pandora DataFlow Team  
**最後更新**: 2025-01-14  
**版本**: 3.2.0

