# Axiom Backend V2 - 擴展 API 功能規格

> **版本** : 2.1.0
>
> **日期** : 2025-10-16
>
> **擴展** : 實用/實驗/組合 APIs

---

## 1. 實用功能 APIs (Practical Function APIs)

### 1.1 pandora-agent 擴展

#### 資產發現與清點

```
GET  /api/v2/agent/{agentId}/assets/discovery
POST /api/v2/agent/{agentId}/assets/scan
GET  /api/v2/agent/{agentId}/assets/inventory
```

* 自動發現網路資產
* 軟體清單掃描
* 漏洞狀態清點

#### 合規性檢查

```
POST /api/v2/agent/{agentId}/compliance/check
GET  /api/v2/agent/{agentId}/compliance/report
POST /api/v2/agent/{agentId}/compliance/remediate
```

* CIS Benchmark 檢查
* NIST 框架合規驗證
* 自動修復建議

#### 遠端指令執行 (受控)

```
POST /api/v2/agent/{agentId}/exec/command
GET  /api/v2/agent/{agentId}/exec/history
POST /api/v2/agent/{agentId}/exec/script
```

* 安全的遠端命令執行
* PowerShell/Bash 腳本執行
* 執行歷史審計

### 1.2 prometheus 擴展

#### 智能基線與異常檢測

```
POST /api/v2/metrics/baseline/create
GET  /api/v2/metrics/baseline/{metricName}
POST /api/v2/metrics/anomaly/detect
GET  /api/v2/metrics/anomaly/history
```

* 自動建立指標基線
* 統計異常檢測
* 異常事件追蹤

#### 容量規劃

```
GET  /api/v2/metrics/capacity/forecast
POST /api/v2/metrics/capacity/analysis
GET  /api/v2/metrics/capacity/recommendations
```

* 資源使用預測
* 容量增長趨勢分析
* 擴容建議

#### 自定義指標聚合

```
POST /api/v2/metrics/aggregation/custom
GET  /api/v2/metrics/aggregation/{aggId}
DELETE /api/v2/metrics/aggregation/{aggId}
```

* 創建自定義聚合規則
* 多維度數據透視

### 1.3 loki 擴展

#### 日誌模式挖掘

```
POST /api/v2/logs/patterns/extract
GET  /api/v2/logs/patterns/list
GET  /api/v2/logs/patterns/{patternId}/occurrences
```

* 自動識別日誌模式
* 頻繁模式挖掘
* 異常模式檢測

#### 日誌關聯分析

```
POST /api/v2/logs/correlation/analyze
GET  /api/v2/logs/correlation/timeline
POST /api/v2/logs/correlation/traces
```

* 跨服務日誌關聯
* 事件時間線重建
* 分佈式追蹤整合

#### 智能日誌解析

```
POST /api/v2/logs/parse/auto
POST /api/v2/logs/parse/template
GET  /api/v2/logs/parse/fields/{logType}
```

* 自動日誌格式識別
* 欄位提取模板
* 結構化日誌轉換

### 1.4 alertmanager 擴展

#### 告警聚類與去重

```
POST /api/v2/alerts/clustering/analyze
GET  /api/v2/alerts/clustering/groups
POST /api/v2/alerts/deduplication/rules
```

* 相似告警自動聚類
* 智能去重規則
* 告警風暴抑制

#### 告警優先級管理

```
POST /api/v2/alerts/priority/calculate
PUT  /api/v2/alerts/{alertId}/priority
GET  /api/v2/alerts/priority/matrix
```

* 動態優先級計算
* 業務影響評估
* 告警優先級矩陣

#### 告警根因分析

```
POST /api/v2/alerts/rca/analyze
GET  /api/v2/alerts/rca/{incidentId}
POST /api/v2/alerts/rca/suggest-actions
```

* 自動根因推理
* 依賴關係分析
* 修復建議生成

### 1.5 cyber-ai-quantum 擴展

#### 威脅情報整合

```
POST /api/v2/quantum/threat-intel/ingest
GET  /api/v2/quantum/threat-intel/search
POST /api/v2/quantum/threat-intel/correlate
```

* 外部威脅情報源整合
* IoC 查詢與匹配
* 威脅關聯分析

#### 量子安全評估

```
POST /api/v2/quantum/security/assess
GET  /api/v2/quantum/security/score
POST /api/v2/quantum/security/recommend
```

* 量子抗性評估
* 密碼學安全評分
* 遷移建議

#### ML 模型管理

```
POST /api/v2/quantum/ml/models/train
GET  /api/v2/quantum/ml/models/list
POST /api/v2/quantum/ml/models/{modelId}/predict
PUT  /api/v2/quantum/ml/models/{modelId}/retrain
```

* 模型訓練與部署
* 模型版本管理
* 在線預測服務

### 1.6 n8n 擴展

#### 工作流模板市場

```
GET  /api/v2/workflows/templates/list
POST /api/v2/workflows/templates/{templateId}/instantiate
POST /api/v2/workflows/templates/publish
```

* 預建工作流模板
* 一鍵部署
* 社區模板分享

#### 工作流測試與驗證

```
POST /api/v2/workflows/{workflowId}/test
POST /api/v2/workflows/{workflowId}/validate
GET  /api/v2/workflows/{workflowId}/coverage
```

* 工作流邏輯測試
* 輸入輸出驗證
* 測試覆蓋率報告

---

## 2. 實驗性 APIs (Experimental APIs)

### 2.1 量子增強功能

#### 量子隨機數生成器 (QRNG)

```
GET  /api/v2/experimental/quantum/random/generate
POST /api/v2/experimental/quantum/random/stream
GET  /api/v2/experimental/quantum/random/entropy-pool
```

* 真量子隨機數
* 高熵密鑰生成
* 熵池管理

#### 量子機器學習 (QML)

```
POST /api/v2/experimental/quantum/qml/classify
POST /api/v2/experimental/quantum/qml/cluster
POST /api/v2/experimental/quantum/qml/optimize
```

* 量子分類器
* 量子聚類算法
* 量子優化求解器

#### 量子區塊鏈整合

```
POST /api/v2/experimental/quantum/blockchain/sign
POST /api/v2/experimental/quantum/blockchain/verify
GET  /api/v2/experimental/quantum/blockchain/audit-trail
```

* 量子簽名
* 抗量子區塊鏈
* 不可變審計日誌

### 2.2 AI 驅動自動化

#### 自然語言查詢 (NLQ)

```
POST /api/v2/experimental/ai/nlq/query
POST /api/v2/experimental/ai/nlq/translate
GET  /api/v2/experimental/ai/nlq/suggestions
```

* 自然語言轉 LogQL/PromQL
* 智能查詢建議
* 上下文感知搜索

#### 自動化運維決策 (AIOps)

```
POST /api/v2/experimental/ai/aiops/incident-predict
POST /api/v2/experimental/ai/aiops/auto-remediate
GET  /api/v2/experimental/ai/aiops/playbook-recommend
```

* 故障預測
* 自動修復執行
* Playbook 推薦

#### 行為分析與異常檢測

```
POST /api/v2/experimental/ai/behavior/profile
POST /api/v2/experimental/ai/behavior/detect-anomaly
GET  /api/v2/experimental/ai/behavior/{entityId}/timeline
```

* 用戶/實體行為畫像
* 異常行為檢測
* 行為時間線分析

### 2.3 邊緣計算與分佈式處理

#### 邊緣節點管理

```
POST /api/v2/experimental/edge/nodes/register
GET  /api/v2/experimental/edge/nodes/list
POST /api/v2/experimental/edge/nodes/{nodeId}/deploy-workload
```

* 邊緣節點註冊
* 工作負載分發
* 邊緣-雲協同

#### 分佈式查詢引擎

```
POST /api/v2/experimental/distributed/query/submit
GET  /api/v2/experimental/distributed/query/{queryId}/status
GET  /api/v2/experimental/distributed/query/{queryId}/results
```

* 分佈式日誌查詢
* 跨集群數據聚合
* 查詢優化

### 2.4 混沌工程

#### 故障注入

```
POST /api/v2/experimental/chaos/inject/latency
POST /api/v2/experimental/chaos/inject/failure
POST /api/v2/experimental/chaos/inject/resource-pressure
DELETE /api/v2/experimental/chaos/experiments/{expId}
```

* 網路延遲注入
* 服務故障模擬
* 資源壓力測試

#### 彈性測試

```
POST /api/v2/experimental/chaos/resilience/test
GET  /api/v2/experimental/chaos/resilience/report
POST /api/v2/experimental/chaos/resilience/game-day
```

* 自動化彈性測試
* 韌性評分報告
* Game Day 演練

---

## 3. 組合實例 APIs (Combination Instance APIs)

### 3.1 安全事件響應工作流

#### 一鍵事件調查

```
POST /api/v2/combined/incident/investigate
```

 **組合服務** : Loki + Prometheus + AlertManager + pandora-agent + cyber-ai-quantum

 **流程** :

1. 從 AlertManager 獲取告警詳情
2. 從 Loki 查詢相關日誌上下文
3. 從 Prometheus 獲取指標異常
4. 通過 Agent 收集主機取證數據
5. AI 服務進行威脅情報關聯和根因分析
6. 生成調查報告

 **回應範例** :

```json
{
  "incidentId": "INC-2025-001",
  "timeline": [...],
  "affectedAssets": [...],
  "threatIntel": {...},
  "rootCause": {...},
  "recommendations": [...]
}
```

#### 自動化威脅狩獵

```
POST /api/v2/combined/threat-hunting/campaign
GET  /api/v2/combined/threat-hunting/{campaignId}/results
```

 **組合服務** : cyber-ai-quantum + Loki + prometheus + pandora-agent

 **流程** :

1. 定義狩獵假設和 IoC
2. 跨所有日誌源搜索匹配
3. 關聯時序指標數據
4. Agent 深入主機調查
5. AI 評估威脅可信度
6. 生成狩獵報告

### 3.2 性能優化建議引擎

#### 全棧性能分析

```
POST /api/v2/combined/performance/analyze
GET  /api/v2/combined/performance/bottlenecks
POST /api/v2/combined/performance/optimize
```

 **組合服務** : Prometheus + Loki + Grafana + postgres + redis

 **流程** :

1. 從 Prometheus 獲取所有服務指標
2. 從 Loki 分析慢查詢日誌
3. 檢查資料庫查詢性能
4. 分析快取命中率
5. 生成優化建議
6. 自動調整配置參數

### 3.3 智能容量管理

#### 預測性擴容

```
POST /api/v2/combined/capacity/forecast-and-scale
GET  /api/v2/combined/capacity/predictions
POST /api/v2/combined/capacity/auto-scale
```

 **組合服務** : Prometheus + cyber-ai-quantum + Portainer + RabbitMQ

 **流程** :

1. 收集歷史資源使用數據
2. AI 預測未來資源需求
3. 計算最優擴容時機
4. 通過 Portainer 自動擴容
5. 調整 RabbitMQ 隊列配置
6. 監控擴容效果

### 3.4 合規性自動化

#### 端到端合規檢查

```
POST /api/v2/combined/compliance/full-audit
GET  /api/v2/combined/compliance/dashboard
POST /api/v2/combined/compliance/remediate-all
```

 **組合服務** : pandora-agent + Loki + postgres + n8n

 **流程** :

1. Agent 執行系統合規掃描
2. 檢查日誌審計合規性
3. 驗證資料庫訪問控制
4. 觸發 n8n 修復工作流
5. 生成合規報告
6. 自動化證據收集

### 3.5 災難恢復演練

#### 全系統 DR 測試

```
POST /api/v2/combined/dr/test/initiate
GET  /api/v2/combined/dr/test/{testId}/status
POST /api/v2/combined/dr/test/{testId}/failover
POST /api/v2/combined/dr/test/{testId}/rollback
```

 **組合服務** : Portainer + postgres + redis + RabbitMQ + n8n + Prometheus

 **流程** :

1. 創建當前狀態快照
2. 模擬服務故障
3. 執行自動化故障轉移
4. 驗證數據一致性
5. 測試服務恢復
6. 生成 DR 報告

### 3.6 多維度可觀測性

#### 統一可觀測性儀表板

```
POST /api/v2/combined/observability/dashboard/create
GET  /api/v2/combined/observability/dashboard/unified
POST /api/v2/combined/observability/correlate-events
```

 **組合服務** : Prometheus + Loki + Grafana + AlertManager + node-exporter

 **流程** :

1. 整合指標、日誌、追蹤
2. 自動創建關聯視圖
3. 智能異常標註
4. 跨維度根因分析
5. 生成統一儀表板

### 3.7 量子增強安全流水線

#### 端到端量子加密通道

```
POST /api/v2/combined/quantum-security/establish-channel
POST /api/v2/combined/quantum-security/secure-transfer
GET  /api/v2/combined/quantum-security/audit
```

 **組合服務** : cyber-ai-quantum + Redis + RabbitMQ + postgres

 **流程** :

1. 建立 QKD 加密通道
2. 量子加密敏感數據
3. 通過 RabbitMQ 安全傳輸
4. Redis 快取加密會話
5. Postgres 記錄審計日誌
6. 定期密鑰輪換

### 3.8 智能告警降噪

#### AI 驅動告警聚合

```
POST /api/v2/combined/alerts/intelligent-grouping
GET  /api/v2/combined/alerts/noise-reduction-report
POST /api/v2/combined/alerts/auto-suppress
```

 **組合服務** : AlertManager + cyber-ai-quantum + Loki + Prometheus

 **流程** :

1. 收集所有告警事件
2. AI 分析告警模式
3. 識別告警風暴
4. 關聯日誌和指標
5. 自動生成抑制規則
6. 僅推送根告警

### 3.9 服務依賴地圖

#### 動態拓撲發現

```
POST /api/v2/combined/topology/discover
GET  /api/v2/combined/topology/map
POST /api/v2/combined/topology/impact-analysis
```

 **組合服務** : Prometheus + Loki + Grafana + RabbitMQ + nginx

 **流程** :

1. 從 Prometheus 服務發現
2. 分析 Nginx 訪問日誌
3. 追蹤 RabbitMQ 消息流
4. 構建服務依賴圖
5. 影響範圍分析
6. Grafana 可視化

### 3.10 成本優化引擎

#### 資源成本分析

```
POST /api/v2/combined/cost/analyze
GET  /api/v2/combined/cost/recommendations
POST /api/v2/combined/cost/optimize
```

 **組合服務** : Prometheus + Portainer + postgres + redis + cyber-ai-quantum

 **流程** :

1. 收集資源使用指標
2. 分析容器資源分配
3. 識別閒置資源
4. AI 優化建議
5. 自動資源調整
6. 成本節省報告

---

## 4. 跨服務事件總線 (Event Bus APIs)

### 4.1 事件發布訂閱

```
POST /api/v2/events/publish
POST /api/v2/events/subscribe
GET  /api/v2/events/stream
DELETE /api/v2/events/subscriptions/{subId}
```

### 4.2 事件重放

```
POST /api/v2/events/replay
GET  /api/v2/events/history
POST /api/v2/events/filter
```

---

## 5. API 認證與授權

### 5.1 量子安全認證

```
POST /api/v2/auth/quantum/challenge
POST /api/v2/auth/quantum/response
POST /api/v2/auth/quantum/refresh
```

### 5.2 多因素驗證

```
POST /api/v2/auth/mfa/setup
POST /api/v2/auth/mfa/verify
GET  /api/v2/auth/mfa/backup-codes
```

### 5.3 API 密鑰管理

```
POST /api/v2/auth/api-keys/generate
GET  /api/v2/auth/api-keys/list
DELETE /api/v2/auth/api-keys/{keyId}
PUT  /api/v2/auth/api-keys/{keyId}/rotate
```

---

## 6. 系統健康與維護

### 6.1 系統級健康檢查

```
GET  /api/v2/system/health/comprehensive
GET  /api/v2/system/health/dependencies
POST /api/v2/system/health/test-all
```

### 6.2 自動化維護任務

```
POST /api/v2/system/maintenance/schedule
POST /api/v2/system/maintenance/execute
GET  /api/v2/system/maintenance/history
```

### 6.3 配置備份與恢復

```
POST /api/v2/system/backup/create
GET  /api/v2/system/backup/list
POST /api/v2/system/backup/{backupId}/restore
```

---

## 實施優先級建議

### 🔴 高優先級 (P0)

1. 統一可觀測性儀表板
2. 一鍵事件調查
3. 智能告警降噪
4. 全棧性能分析

### 🟡 中優先級 (P1)

1. 自動化威脅狩獵
2. 預測性擴容
3. 端到端合規檢查
4. 日誌模式挖掘

### 🟢 低優先級 (P2)

1. 量子機器學習
2. 混沌工程
3. 邊緣計算管理
4. 成本優化引擎

### 🔵 實驗性 (P3)

1. 量子區塊鏈整合
2. 自然語言查詢
3. 分佈式查詢引擎
4. 量子隨機數生成器

---

## API 響應格式標準

所有 API 遵循統一響應格式：

```json
{
  "success": true,
  "timestamp": "2025-10-16T10:30:00Z",
  "requestId": "req-uuid-here",
  "data": { ... },
  "metadata": {
    "executionTime": "1.23s",
    "servicesInvolved": ["loki", "prometheus"],
    "cacheHit": false
  },
  "errors": [],
  "warnings": []
}
```

---

## 版本控制策略

* API 版本通過 URL 路徑指定 (`/api/v2/...`)
* 實驗性 API 在 `/api/v2/experimental/` 下
* 向後兼容至少維護兩個大版本
* 廢棄 API 提供至少 6 個月遷移期
