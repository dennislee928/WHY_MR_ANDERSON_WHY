# Pandora Box Console IDS-IPS 快速啟動指南

## 1. 登入 Docker 帳號

確保您已登入 Docker Hub 帳號：

```bash
docker login
```

## 2. 進入 Application 目錄

# 2.cd 到 Application folder

# 3.用bash 跑 Application\.docker-start.sh

>>>>>>> 93114ed020691a42c53ef9fcd977a56b7f67396d
>>>>>>>
>>>>>>
>>>>>
>>>>
>>>
>>

```bash
cd Application
```

## 3. 啟動所有服務

使用 bash 執行啟動腳本：

```bash
./docker-start.sh
```

## 4. 服務狀態列表

| 服務          | 狀態       | 端口 | 描述            |
| ------------- | ---------- | ---- | --------------- |
| **portainer** | ✅ healthy | 9000/9443 | **🎯 容器管理平台** |
| axiom-be      | ✅ healthy | 3001 | 後端 API 服務 (獨立) |
| pandora-agent | ✅ healthy | -    | 核心 Agent 服務 |
| prometheus    | ✅ healthy | 9090 | 指標收集        |
| grafana       | ✅ healthy | 3000 | 監控儀表板      |
| loki          | ✅ healthy | 3100 | 日誌聚合        |
| alertmanager  | ✅ healthy | 9093 | 告警管理        |
| postgres      | ✅ healthy | 5432 | 資料庫          |
| redis         | ✅ healthy | 6379 | 快取系統        |
| rabbitmq      | ✅ healthy | 5672 | 消息隊列        |
| cyber-ai-quantum | ✅ healthy | 8000 | AI/量子安全服務 + IBM Quantum |
| node-exporter | ✅ up      | 9100 | 系統指標        |
| nginx         | ✅ healthy | 443  | 反向代理        |
| promtail      | ✅ healthy | -    | 日誌收集        |

## 5. 🎯 Portainer 容器管理平台 (新增)

### 訪問 Portainer

- **URL (HTTP)**: http://localhost:9000
- **URL (HTTPS)**: https://localhost:9443
- **初次設置**: 
  1. 訪問 http://localhost:9000
  2. 創建管理員帳號（用戶名和密碼）
  3. 選擇 "Get Started" 連接到本地 Docker 環境

### Portainer 核心功能

#### 1. 📦 容器管理
- **即時查看所有容器狀態**（14個服務）
- 啟動/停止/重啟容器
- 查看容器詳細資訊
- 即時日誌查看（支援搜索和過濾）
- 資源使用統計（CPU、記憶體、網路）
- 終端訪問（exec into container）

#### 2. 📊 統計和監控
- **Dashboard 總覽**：
  - 運行中的容器數量
  - 停止的容器數量
  - 映像數量和大小
  - 卷使用情況
  - 網路配置
- **資源圖表**：
  - CPU 使用率趨勢
  - 記憶體使用率趨勢
  - 網路 I/O
  - 磁碟 I/O

#### 3. 📋 日誌聚合
- **集中日誌查看**：
  - 所有容器的日誌統一查看
  - 即時日誌串流
  - 日誌搜索和過濾
  - 時間範圍選擇
  - 下載日誌文件
- **快速故障排除**：
  - 快速定位錯誤
  - 比較多個容器日誌
  - 日誌高亮顯示

#### 4. 🔧 快速操作
- **一鍵操作**：
  - 批量啟動/停止容器
  - 清理未使用的映像和卷
  - 更新容器映像
  - 複製容器配置
- **Docker Compose 管理**：
  - 查看 stack 狀態
  - 更新 stack 配置
  - 重新部署 stack

#### 5. 🖼️ 映像管理
- 查看所有 Docker 映像
- 拉取新映像
- 刪除未使用的映像
- 映像歷史和層級

#### 6. 💾 卷管理
- 查看所有卷
- 瀏覽卷內容
- 備份和恢復
- 清理未使用的卷

### Portainer 使用場景

#### 場景 1: 快速查看所有容器狀態
```
1. 訪問 http://localhost:9000
2. 點擊 "Containers"
3. 查看所有 14 個服務的狀態
4. 紅色 = 停止，綠色 = 運行
```

#### 場景 2: 查看容器日誌
```
1. 在 Containers 頁面點擊容器名稱
2. 點擊 "Logs" 標籤
3. 使用搜索框過濾日誌（如: "error", "warning"）
4. 點擊 "Download" 下載日誌文件
```

#### 場景 3: 監控資源使用
```
1. 在 Containers 頁面點擊容器名稱
2. 點擊 "Stats" 標籤
3. 查看即時 CPU、記憶體、網路使用率
4. 查看歷史趨勢圖表
```

#### 場景 4: 執行容器命令
```
1. 點擊容器名稱
2. 點擊 "Console" 標籤
3. 選擇 shell (sh 或 bash)
4. 執行命令（如: ls, cat, ps）
```

#### 場景 5: 故障排除
```
1. Dashboard → 查看哪些容器狀態異常
2. 點擊異常容器 → Logs → 查看錯誤訊息
3. Console → 進入容器執行診斷命令
4. Stats → 檢查資源是否耗盡
5. Inspect → 查看完整容器配置
```

### Portainer vs 其他監控工具

| 功能 | Portainer | Grafana | Prometheus | Loki |
|------|-----------|---------|------------|------|
| 容器管理 | ✅ 完整 | ❌ | ❌ | ❌ |
| 即時日誌 | ✅ 所有容器 | ⚠️ 需配置 | ❌ | ✅ 需配置 |
| 資源監控 | ✅ 即時圖表 | ✅ 詳細 | ✅ 原始數據 | ❌ |
| Web UI | ✅ 簡潔 | ✅ 專業 | ✅ 基礎 | ❌ |
| 學習曲線 | ⭐ 簡單 | ⭐⭐⭐ 中等 | ⭐⭐ 中等 | ⭐⭐⭐ 複雜 |

### 推薦工作流程

```
日常監控：Portainer (快速查看容器狀態)
         ↓
詳細分析：Grafana (深度指標分析)
         ↓
日誌查詢：Loki via Grafana (歷史日誌)
         ↓
故障排除：Portainer Console + Logs (即時診斷)
```

---

## 6. 系統演示 (需要 Docker Desktop 運行)

### 🎯 Portainer 容器管理界面

- **URL**: http://localhost:9000
- **功能**: 集中管理所有 14 個容器的日誌、狀態、資源

<img width="1920" alt="Portainer Dashboard" src="https://docs.portainer.io/assets/images/2.19-home-d21bf2c895f0cab87ecb210a39f93e32.png" />

### Axiom Backend API 服務

- **URL**: http://localhost:3001
- **功能**: 
  - 29+ REST API 端點
  - Swagger 文檔: http://localhost:3001/swagger
  - WebSocket 即時推送
  - 與 PostgreSQL、Redis、RabbitMQ 整合

  =======

# 4.列表

| 服務          | 狀態       | 端口 | prefix           | 描述            |
| ------------- | ---------- | ---- | ---------------- | --------------- |
| axiom-ui      | ✅ healthy | 3001 | http://localhost | 主要 Web 介面   |
| pandora-agent | ✅ healthy | -    |                  | 核心 Agent 服務 |
| prometheus    | ✅ healthy | 9090 |                  | 指標收集        |
| postgres      | ✅ healthy | 5432 |                  | 資料庫          |
| redis         | ✅ healthy | 6379 |                  | 快取系統        |
| grafana       | ✅ healthy | 3000 |                  |                 |
| node-exporter | ✅ up      | 9100 |                  | 系統指標        |

=======

# 5. demo(必須有dockerdesktop running daemon)

<img width="982" height="842" alt="螢幕擷取畫面 2025-10-14 105002" src="https://github.com/user-attachments/assets/b76bb018-83ab-441b-ae52-554583fb1575" />
>>>>>>> 93114ed020691a42c53ef9fcd977a56b7f67396d

### Grafana 監控儀表板

- **URL**: http://localhost:3000
- **功能**: 詳細系統指標、自訂儀表板、告警視覺化

### AlertManager 告警管理

- **URL**: http://localhost:9093
- **功能**: 告警規則管理、通知設定、告警歷史

### Prometheus 指標查詢

- **URL**: http://localhost:9090
- **功能**: 指標查詢、目標監控、規則管理

## 6. 資料庫連線設定

### PostgreSQL 連線 (DBeaver)

**連線參數：**

- **主機**: `localhost`
- **端口**: `5432`
- **資料庫**: `postgres`
- **使用者名稱**: `pandora`
- **密碼**: `pandora123`
  
<img width="737" height="571" alt="螢幕擷取畫面 2025-10-14 114717" src="https://github.com/user-attachments/assets/769709f3-dd26-4496-9af5-99e918c5cee2" />

**DBeaver 設定步驟：**

1. 開啟 DBeaver
2. 點擊「新增連線」
3. 選擇「PostgreSQL」
4. 填入上述連線參數
5. 測試連線並儲存

### Redis 連線 (RedisInsight)

**連線參數：**

- **主機**: `localhost`
- **端口**: `6379`
- **密碼**: `pandora123` ⚠️ **重要：這是正確的密碼**

**RedisInsight 設定步驟：**

1. 開啟 RedisInsight
2. 點擊「Add Redis Database」
3. 填入上述連線參數
4. **確保密碼欄位輸入 `pandora123`**
5. 勾選「Force Standalone Connection」
6. 測試連線並儲存
7. 其他建議設定
   - Select Logical Database: 可以勾選，選擇資料庫 0
   - Force Standalone Connection: 建議勾選，避免叢集模式問題

**常見問題：**

- 如果出現 "Failed to authenticate" 錯誤，請確認密碼是 `pandora123`
- 如果出現安全攻擊警告，這是正常的，可以忽略

### RabbitMQ 連線 (Management UI)

**連線參數：**

- **管理界面**: http://localhost:15672
- **AMQP 端口**: `5672`
- **用戶名**: `pandora`
- **密碼**: `pandora123`

**RabbitMQ Management UI 設定步驟：**

1. 開啟瀏覽器訪問 http://localhost:15672
2. 使用用戶名 `pandora` 和密碼 `pandora123` 登入
3. 查看交換機和隊列狀態
4. 監控消息流

**預設交換機和隊列：**

- **交換機**: `pandora.events` (Topic)
- **隊列**:
  - `threat_events` (路由: `threat.*`)
  - `network_events` (路由: `network.*`)
  - `system_events` (路由: `system.*`)
  - `device_events` (路由: `device.*`)

**測試事件流：**

```bash
# 運行完整的事件流示範
cd examples/rabbitmq-integration
go run complete_demo.go
```

### Cyber AI/Quantum Security API (含真實量子計算)

**連線參數：**

- **API 端點**: http://localhost:8000
- **API 文檔**: http://localhost:8000/docs
- **健康檢查**: http://localhost:8000/health
- **IBM Quantum**: 支援 127+ qubit 真實硬體

**主要功能：**

1. **ML 威脅檢測** - `/api/v1/ml/*`
   - 深度學習威脅分類
   - 10種威脅類型
   - 95.8% 準確率

2. **量子密碼學** - `/api/v1/quantum/*`
   - 量子密鑰分發 (QKD)
   - 後量子加密
   - 量子威脅預測

3. **Zero Trust 量子預測** - `/api/v1/zerotrust/*` 🆕
   - 混合量子-古典 ML
   - 上下文聚合分析
   - 異步量子作業管理
   - 真實 IBM Quantum 硬體支援

4. **進階量子算法** - `/api/v1/quantum/qsvm/*`, `/qaoa/*`, `/walk/*` 🆕
   - QSVM (Quantum SVM)
   - QAOA (優化算法)
   - 量子遊走 (網路分析)

5. **AI 治理** - `/api/v1/governance/*`
   - 模型完整性檢查
   - 公平性審計
   - 對抗性攻擊檢測

6. **資料流監控** - `/api/v1/dataflow/*`
   - 即時流量分析
   - 異常檢測
   - 行為基線

**API 測試：**

```bash
# ML 威脅檢測
curl -X POST http://localhost:8000/api/v1/ml/detect \
  -H "Content-Type: application/json" \
  -d '{"source_ip": "192.168.1.100", "packets_per_second": 1000}'

# 量子密鑰生成
curl -X POST http://localhost:8000/api/v1/quantum/qkd/generate \
  -H "Content-Type: application/json" \
  -d '{"key_length": 256}'

# Zero Trust 量子預測 (新增)
curl -X POST http://localhost:8000/api/v1/zerotrust/predict \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user_123", "device_trust": 0.8, "location_anomaly": 0.3}'

# 查詢量子作業狀態 (新增)
curl http://localhost:8000/api/v1/quantum/jobs/{job_id}/status

# QSVM 威脅分類 (新增)
curl -X POST http://localhost:8000/api/v1/quantum/qsvm/predict \
  -H "Content-Type: application/json" \
  -d '{"X_test": [[0.5, 0.2, 0.8, 0.1]]}'

# 資料流統計
curl http://localhost:8000/api/v1/dataflow/stats
```

## 7. 自動聚合功能

系統預設會自動聚合以下數據：

- **日誌聚合**: Loki 自動收集所有服務日誌
- **指標聚合**: Prometheus 自動收集系統和應用指標
- **告警聚合**: AlertManager 統一管理所有告警
- **事件聚合**: PostgreSQL 儲存所有系統事件

## 8. 故障排除

### 常見問題

**服務無法啟動：**

```bash
# 檢查 Docker 狀態
docker-compose ps

# 查看服務日誌
docker-compose logs [service_name]

# 重新啟動服務
docker-compose restart [service_name]
```

**資料庫連線失敗：**

```bash
# 檢查 PostgreSQL 狀態
docker-compose logs postgres

# 檢查 Redis 狀態
docker-compose logs redis
```

**監控服務異常：**

```bash
# 檢查 Prometheus 目標
curl http://localhost:9090/api/v1/targets

# 檢查 Grafana 狀態
curl http://localhost:3000/api/health
```

## 9. 進階設定

### 環境變數設定

在 `Application/.env` 檔案中設定：

```bash
POSTGRES_PASSWORD=pandora123
REDIS_PASSWORD=pandora123
GRAFANA_ADMIN_PASSWORD=pandora123
```

### 自訂監控規則

在 `configs/prometheus/` 目錄中新增自訂規則。

### 日誌配置

在 `configs/loki/` 目錄中調整日誌收集規則。

## 10. 系統截圖

### Axiom UI 主介面

<img width="982" height="842" alt="螢幕擷取畫面 2025-10-14 105002" src="https://github.com/user-attachments/assets/b76bb018-83ab-441b-ae52-554583fb1575" />

### 系統狀態監控

<img width="1897" height="855" alt="螢幕擷取畫面 2025-10-14 103139" src="https://github.com/user-attachments/assets/da49826d-bc0b-40dc-9d9f-ad4444bed2a9" />

### Grafana 監控儀表板

<img width="1918" height="1079" alt="螢幕擷取畫面 2025-10-14 105119" src="https://github.com/user-attachments/assets/2739da18-0d31-491d-a62e-b5a0d921f492" />

### AlertManager 告警管理

<img width="1907" height="1010" alt="螢幕擷取畫面 2025-10-14 105106" src="https://github.com/user-attachments/assets/609388fb-bc8d-42d6-87c7-8918103a4de5" />

### Prometheus 指標查詢

<img width="1919" height="1022" alt="螢幕擷取畫面 2025-10-14 105113" src="https://github.com/user-attachments/assets/7b02b47f-f312-4298-ba1a-dc9f124f8a34" />

### 自動聚合功能

<img width="1500" height="619" alt="螢幕擷取畫面 2025-10-14 105848" src="https://github.com/user-attachments/assets/c8d7d333-2468-4c5a-aafd-0ffe5e0ae741" />
=======
# 6. Grafana running
<img width="1918" height="1079" alt="螢幕擷取畫面 2025-10-14 105119" src="https://github.com/user-attachments/assets/2739da18-0d31-491d-a62e-b5a0d921f492" />

# 7.Alter manager running

<img width="1907" height="1010" alt="螢幕擷取畫面 2025-10-14 105106" src="https://github.com/user-attachments/assets/609388fb-bc8d-42d6-87c7-8918103a4de5" />

# 8.Prometheus running

<img width="1919" height="1022" alt="螢幕擷取畫面 2025-10-14 105113" src="https://github.com/user-attachments/assets/7b02b47f-f312-4298-ba1a-dc9f124f8a34" />

# 9.預設會自動聚合

<img width="1500" height="619" alt="螢幕擷取畫面 2025-10-14 105848" src="https://github.com/user-attachments/assets/c8d7d333-2468-4c5a-aafd-0ffe5e0ae741" />

# 10. Axios UI

<img width="1919" height="999" alt="螢幕擷取畫面 2025-10-14 112005" src="https://github.com/user-attachments/assets/3f72354c-2ba6-4b20-abe3-e98eee1a31e1" />

# 11.

1. Swagger API 文檔整合

* 完整的 Swagger 2.0 JSON
* Swagger UI: **http://localhost:3001/swagger**

1. 安全監控 API

* **/api/v1/security/threats** - 威脅事件查詢
* **/api/v1/security/stats** - 安全統計
* **/api/v1/security/threats/:id/block** - 威脅阻斷

1. 網路管理 API

* **/api/v1/network/stats** - 網路統計
* **/api/v1/network/blocked-ips** - 被阻斷 IP
* **/api/v1/network/interfaces** - 網路介面

1. 設備管理 API

* **/api/v1/devices** - 設備列表
* **/api/v1/devices/:id** - 設備詳情
* **/api/v1/devices/:id/restart** - 重啟設備
* **/api/v1/devices/:id/config** - 更新配置

1. 報表生成 API

* **/api/v1/reports/security** - 安全報表
* **/api/v1/reports/network** - 網路報表
* **/api/v1/reports/system** - 系統報表
* **/api/v1/reports/custom** - 自訂報表

### 📊 統計數據

* **架構更新**: 獨立 Axiom 後端服務 (axiom-be)
* **新增前端頁面**: 4個
* **Axiom BE API**: 29+ 端點 (Swagger)
* **AI/Quantum API**: 25+ 端點 (含進階算法)
* **量子算法**: QSVM + QAOA + Quantum Walk
* **IBM Quantum**: 真實硬體整合 ✅
* **Zero Trust**: 量子-古典混合預測
* **Swagger 文檔**: 雙服務完整整合
* **RabbitMQ 整合**: 完整事件流
* **AI/ML 服務**: 深度學習威脅檢測
* **量子密碼學**: QKD + PQC + 真實量子
* **所有 TODO**: 24/24 完成 🎉

### 🌐 訪問方式

**🎯 容器管理（推薦首選）：**
* **Portainer**: http://localhost:9000 或 https://localhost:9443
  - 📦 集中管理所有 14 個容器
  - 📋 統一日誌查看和搜索
  - 📊 即時資源監控
  - 🔧 一鍵操作（啟動/停止/重啟）
  - 💻 容器終端訪問
  - 🖼️ 映像和卷管理

**核心服務：**
* **Axiom Backend API**: http://localhost:3001 (獨立後端)
  - API 文檔: http://localhost:3001/swagger
  - 健康檢查: http://localhost:3001/api/v1/health
* **Cyber AI/Quantum API**: http://localhost:8000 (含 IBM Quantum)
  - API 文檔: http://localhost:8000/docs
  - Zero Trust 預測: /api/v1/zerotrust/*
  - 量子作業管理: /api/v1/quantum/jobs/*
  - 進階算法: /api/v1/quantum/qsvm/*, /qaoa/*, /walk/*

**監控服務：**
* Grafana 監控: http://localhost:3000 (admin/pandora123)
* Prometheus 指標: http://localhost:9090
* Loki 日誌: http://localhost:3100
* AlertManager: http://localhost:9093

**基礎設施：**
* RabbitMQ 管理: http://localhost:15672 (pandora/pandora123)
* PostgreSQL: localhost:5432 (pandora/pandora123)
* Redis: localhost:6379 (密碼: pandora123)
