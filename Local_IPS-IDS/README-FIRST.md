# Box控制台：量子增強智能運維平台

## Pandora Box Console IDS-IPS

 **版本** ：v3.4.1 (穩定版) / Axiom Backend V3 (開發藍圖)
 **狀態** ：🏆 **世界級生產就緒** + 量子增強 + 工作流自動化
 **更新日期** ：2025-10-15



---

# 完整啟動:

cd Application
   docker-compose up -d

---

## 🚀 專案核心與架構概覽 (v3.4.1)

本專案是一個基於微服務架構的 IDS-IPS 系統，其核心為  **Axiom Backend V2/V3** ，它作為統一的 API Gateway 和控制中心，統一管理所有服務。

最新版本 **v3.4.1** 強化了統一入口點 (Nginx) 和工作流自動化 (n8n) 能力。

### 🌐 核心架構流 (C4 Model - System Context)

所有外部請求都必須通過 **Nginx** (API Gateway) 進行路由和安全驗證。

```
mermaid

```

```
graph TD
    A[外部使用者/Agent] -->|Port 80/443| B(Nginx: API Gateway)
    B -->|/api/| C(Axiom Backend: 核心控制中心)
    B -->|/grafana/| D(Grafana: 視覺化儀表板)
    B -->|/prometheus/| E(Prometheus: 指標收集)
    B -->|其他服務/靜態資源| F(其他內部服務/UI)
    C -->|API 呼叫| G(Cyber AI/Quantum: 量子安全服務)
    C -->|DB 連接| H(PostgreSQL: 資料持久化)
    C -->|事件發布| I(RabbitMQ: 消息隊列)
    C -->|快取/限流| J(Redis: 快取系統)
    D --> E
    I --> K(n8n: 工作流自動化)
    K --> L(外部服務, e.g., Slack/Email)
    M(Portainer: 容器管理) --> B
    M --> K
    subgraph 核心服務
        C
        G
    end
    subgraph 監控與基礎設施
        D
        E
        H
        I
        J
        K
    end
    subgraph 閘道層
        B
    end

    style B fill:#00CCAA,stroke:#333
    style C fill:#ADD8E6,stroke:#333
    style G fill:#E6CCFF,stroke:#333
    style K fill:#FFDDDD,stroke:#333
```

* **Nginx** 作為統一入口點，提供反向代理、安全閘道和效能優化（Gzip 壓縮、連接池）功能。
* **Axiom Backend (Port 3001)** 負責統一管理 13 個服務。
* **n8n (Port 5678)** 是新增的工作流自動化平台，用於告警通知、API 串接和複雜業務流程編排。

---

## 📊 系統服務組件清單 (v3.4.1)

系統總共包含  **15 個容器服務** 。Axiom Backend V2/V3 負責對這些服務進行集中管理.

```
graph TD
    subgraph 閘道與管理層 (Ports 80/443, 9000, 5678)
        N[Nginx: API Gateway] -->|路由| A
        N -->|路由| B
        N -->|路由| C
        P(Portainer: 容器管理) -->|控制| A
        P -->|控制| B
        O(n8n: 工作流自動化) -->|Webhook/API| A
        O -->|發布事件| H
    end

    subgraph 核心服務層 (Port 3001, 8000)
        A(Axiom Backend: 統一控制中心)
        B(Cyber AI/Quantum: 量子安全服務)
        F(pandora-agent: 核心 Agent)
    end

    subgraph 數據與基礎設施 (Ports 5432, 6379, 5672)
        D(PostgreSQL: 資料庫)
        E(Redis: 快取/限流)
        H(RabbitMQ: 消息隊列)
        I(RabbitMQ Mgmt: 15672)
    end

    subgraph 可觀測性與監控 (Ports 9090, 3000, 3100, 9093)
        C(Prometheus: 指標收集)
        J(Grafana: 儀表板)
        K(Loki: 日誌聚合)
        L(AlertManager: 告警管理)
        M(Node Exporter: 系統指標)
        R(Promtail: Log Collector)
    end

    N --> C
    N --> J
    N --> K
    N --> L
    F --> A
    M --> C
    R --> K
    A --> D
    A --> E
    A --> H
    A --> B
    C --> J
    C --> L
    L --> O
    O --> D
```

| 服務                       | 端口      | 描述            | 核心功能                   |
| -------------------------- | --------- | --------------- | -------------------------- |
| **Nginx**🌐          | 80/443    | 統一入口閘道    | 反向代理、安全、效能優化   |
| **Axiom Backend**    | 3001      | 核心 API 服務   | 統一管理 13 個服務         |
| **Cyber AI/Quantum** | 8000      | AI/量子安全服務 | Zero Trust 預測、QKD、QSVM |
| **n8n**🔄            | 5678      | 工作流自動化    | 視覺化工作流編輯器         |
| **Portainer**📦      | 9000/9443 | 容器管理平台    | 集中管理 15 個容器         |
| **Prometheus**       | 9090      | 指標收集        | 系統及應用指標             |
| **Grafana**          | 3000      | 監控儀表板      | 視覺化、告警通知           |
| **Loki**             | 3100      | 日誌聚合        | LogQL 查詢                 |
| **PostgreSQL**       | 5432      | 資料庫          | 儲存 9 個 GORM 模型        |
| **Redis**            | 6379      | 快取系統        | 15+ Key 模式、速率限制     |
| **RabbitMQ**         | 5672      | 消息隊列        | 事件驅動架構               |

---

## 🏗️ Axiom Backend V3 藍圖：統一控制中心

Axiom Backend V3 計劃實現  **300+ API 端點** ，旨在成為世界級的統一 API Gateway 和控制中心。它將 13 個核心服務整合為一個單一的、智能的控制層。

### 📊 Phase 1 資料庫設計 (GORM Models)

Phase 1 架構設計已完成，定義了 9 個核心 GORM 模型來追蹤系統狀態、安全事件和量子作業。

```
classDiagram
    direction LR
    class SystemModel {
        <<Abstract>>
        +int ID
        +time CreatedAt
        +time UpdatedAt
    }

    class User{
        -string Username
        -string Role (RBAC)
        -string APIKey
    }

    class Session{
        -string SessionToken
        -time ExpiresAt
    }

    class Service{
        -string ServiceName (e.g., prometheus)
        -string Status (healthy/unhealthy)
        -string Endpoint
    }

    class Alert{
        -string AlertName
        -string Severity
        -string Status
        -JSONB Context
    }

    class MetricSnapshot{
        -string MetricName
        -float Value
        -time Timestamp
    }

    class WindowsLog{
        -string EventID
        -string LogType
        -string Message
    }

    class QuantumJob{
        -string JobID
        -string Algorithm (QKD/QSVM/QAOA)
        -string Status (pending/completed)
        -JSONB Result
    }

    class ConfigHistory{
        -string Service
        -string ConfigVersion
        -string ChangeSummary
    }

    Service <|-- SystemModel
    Alert <|-- SystemModel
    MetricSnapshot <|-- SystemModel
    WindowsLog <|-- SystemModel
    QuantumJob <|-- SystemModel
    ConfigHistory <|-- SystemModel
    User <|-- SystemModel
    Session <|-- SystemModel

    User "1" -- "N" Session : 登入會話管理
    Service "1" -- "N" ConfigHistory : 配置變更追蹤
    Alert "1" -- "N" QuantumJob : 觸發量子分析 (Optional)
    WindowsLog "1" -- "N" Alert : 日誌觸發告警

    note for Service "追蹤 13 個核心服務狀態"
    note for QuantumJob "支援 10+ 種量子算法"
```

* PostgreSQL 採用完整的索引設計、外鍵約束、JSONB 欄位支援和軟刪除特性。
* Redis 快取架構包含 15+ 種 Key 模式，用於服務健康狀態、即時指標和 API 速率限制。

---

## 🌟 核心創新亮點 (P0 級功能)

Axiom Backend V3 的核心價值在於其跨服務協同的組合功能 (Combined APIs)，以及多項業界首創的創新功能。

| 創新功能                       | 描述                                                          | 組合服務                                      | 狀態 (V3 進度)    |
| ------------------------------ | ------------------------------------------------------------- | --------------------------------------------- | ----------------- |
| **時間旅行調試**⭐⭐⭐   | 業界首創，捕獲完整系統狀態快照，進行時間點恢復和 What-If 分析 | Loki + Prometheus + PG + Redis                | ✅ 已實現核心功能 |
| **智能自癒系統**⭐⭐⭐   | AI 自動診斷根因，並執行修復動作 (如重啟、擴容、配置回滾)      | AlertManager + AI + Agent + Portainer + n8n   | ✅ 已實現核心編排 |
| **自適應安全策略**⭐⭐⭐ | 實時風險評分，動態調整訪問控制，自動蜜罐部署                  | AI + Nginx + Agent                            | ✅ 已實現核心服務 |
| **零信任自動驗證**⭐⭐⭐ | Agent 持續收集狀態，AI 計算實時信任分數，動態調整權限         | Agent + AI + AlertManager + Loki              | ⏳ 待實施 (P0)    |
| **一鍵事件調查**         | 組合所有可觀測數據，自動生成事件調查報告                      | Loki + Prometheus + AlertManager + Agent + AI | ✅ 已實現         |

### 🧠 組合功能範例：零信任自動驗證流水線 (P0)

此 P0 級組合 API 旨在實現  **下一代安全架構** 。

```
sequenceDiagram
    participant A as Pandora Agent (設備端)
    participant B as Axiom BE (控制中心)
    participant C as Cyber AI/Quantum (信任分數計算)
    participant D as Loki/Prometheus (狀態數據)
    participant E as AlertManager/n8n (響應)

    A ->> B: 1. 上報設備健康狀態/行為日誌
    B ->> C: 2. 請求計算[User/Device]實時信任分數
    C ->> D: 3. 提取歷史數據和實時指標
    C -->> B: 4. 返回 Trust Score (0-100)
    alt Score < 0.5 (低信任)
        B ->> E: 5. 觸發[Low Trust]告警/n8n工作流
        E ->> E: 6. 執行: (a) 隔離主機; (b) 要求 MFA
        E -->> B: 7. 報告調整結果
    else Score > 0.8 (高信任)
        B ->> A: 5. 維持/提升訪問權限
    end
    B ->> H(PostgreSQL): 8. 記錄所有驗證決策 (審計日誌)
```

---

## 📈 實施路線圖與進度追蹤 (Axiom Backend V3)

Axiom Backend V3 項目總預計時間約為  **40-50 天** 。截至目前 (2025-10-16)，基礎和部分創新功能已完成。

### 📊 總體進度

| 階段                                 | 預計時間           | 狀態      | 完成度         | API 數量       | 優先級 |
| ------------------------------------ | ------------------ | --------- | -------------- | -------------- | ------ |
| Phase 1: 架構設計                    | 1 天               | ✅ 完成   | 100%           | -              | P0     |
| Phase 2.1-2.4: 基礎 API (13服務)     | 3 天               | ✅ 完成   | 100%           | 30+            | P0     |
| Phase 2.6: 核心組合 APIs (5個)       | 2 天               | ✅ 完成   | 100%           | 5              | P0     |
| **Phase 7 (部分): 創新功能**   | 7-10 天            | ✅ 完成   | ~30%           | 15+            | P0/P1  |
| Phase 2.5: 實用功能 APIs             | 2 天               | ⏳ 待實施 | 0%             | 40+            | P1     |
| **Phase 11-13 (P0): 企業架構** | 10+ 天             | ⏳ 待實施 | 0%             | -              | P0     |
| Phase 6, 8, 9 (餘): 實驗與高級組合   | 15+ 天             | ⏳ 待實施 | 0%             | 80+            | P1-P3  |
| **總計**                       | **40-50 天** | 🚧 進行中 | **~45%** | **300+** | -      |

### 📅 實施階段規劃 (Gantt Chart)

P0 核心和 P1 基礎架構是短期焦點。

```
gantt
    dateFormat YYYY-MM-DD
    title Axiom Backend V3 實施計劃 (43天)

    section 核心基礎 (已完成/P0)
    Phase 1: 架構設計 :done, ID1, 2025-10-16, 1d
    Phase 2.1-2.4: 基礎服務API :done, ID2, after ID1, 3d
    Phase 2.6: 核心組合API :done, ID3, after ID2, 2d
    Phase 7.1/7.3: 時間旅行/自適應安全 :done, ID4, after ID3, 3d

    section 企業核心 (P0 - 立即開始)
    Phase 11: Agent 進階架構 :active, ID5, after ID4, 3d
    Phase 12: 四層儲存架構 :active, ID6, after ID5, 7d
    Phase 13: 合規性引擎 :active, ID7, after ID6, 5d

    section 增值功能 (P1)
    Phase 2.5: 實用功能API : ID8, after ID4, 5d
    Phase 3: Agent 增強 (Log Collector) : ID9, after ID8, 2d
    Phase 4/5: Frontend/文檔測試 : ID10, after ID9, 5d
    Phase 7.2: 數字孿生 : ID11, after ID7, 3d

    section 實驗與創新 (P2/P3)
    Phase 6/8/9: 實驗性/高級組合 : ID12, after ID11, 17d

    section 里程碑
    Milestone 1: 核心功能就緒 (Day 5) : milestone, 2025-10-21
    Milestone 2: 增值功能完成 (Day 16) : milestone, 2025-10-28
    Milestone 3: 生產就緒 (Day 23) : milestone, 2025-11-04
    Milestone 4: 完全體 (Day 50) : milestone, 2025-11-28
```

---

## 🛠️ 部署與啟動指南 (Docker Compose)

### 1. 服務訪問

本地開發環境使用  **Nginx 統一入口** ，或直接透過服務端口訪問。

```
stateDiagram-v2
    direction LR
    state Nginx_Entry <<choice>>
    state Axiom_BE <<choice>>
    state Monitoring <<choice>>
    state Automation <<choice>>

    [*] --> Nginx_Entry : http://localhost:80/443

    Nginx_Entry --> Axiom_BE: /api/ (3001)
    Nginx_Entry --> Monitoring: /grafana/, /prometheus/

    Axiom_BE --> API_DOCS : /swagger (3001)
    Axiom_BE --> Quantum : :8000

    Monitoring --> Grafana : :3000
    Monitoring --> Prometheus : :9090
    Monitoring --> Loki : :3100
    Monitoring --> AlertManager : :9093

    Automation --> n8n : :5678 (admin/pandora123)

    Database_Infra --> PostgreSQL : :5432 (pandora/pandora123)
    Database_Infra --> Redis : :6379 (pandora123)
    Database_Infra --> RabbitMQ_Mgmt : :15672 (pandora/pandora123)

    Nginx_Entry --> Container_Management : /portainer/ (9000/9443)
```

### 2. 部署環境要求

| 項目                 | 要求            | 說明                                 |
| -------------------- | --------------- | ------------------------------------ |
| **CPU**        | 4 核心以上      | 建議運行 15 個容器                   |
| **記憶體**     | 8GB 以上        | 建議                                 |
| **Go**         | 1.21 或更高版本 | Axiom Backend 技術棧                 |
| **PostgreSQL** | 15 或更高版本   | 核心資料庫                           |
| **Docker**     | 24.0 或更高版本 | 推薦使用 Docker Compose              |
| **依賴服務**   | 15 個容器       | Prometheus, Grafana, Loki, n8n, etc. |

### 3. 資料庫遷移 (PostgreSQL)

Axiom Backend V3 使用 PostgreSQL 15+。初始部署時，應用程式啟動會自動執行 `AutoMigrate()` 創建所有 9 個必要的表。

生產環境建議使用 **SQL 遷移腳本** 方式，以獲得完全控制、可審查和支援回滾的能力。

| 檔案 (示例)                  | 版本  | 說明         |
| ---------------------------- | ----- | ------------ |
| `001_initial_schema.sql`   | 3.0.0 | 初始 Schema  |
| `002_add_indexes.sql`      | 3.0.1 | 添加性能索引 |
| `003_add_audit_fields.sql` | 3.1.0 | 添加審計欄位 |

### 4. 故障排除 (範例)

| 問題                          | 症狀                                | 解決方案                                                                       |
| ----------------------------- | ----------------------------------- | ------------------------------------------------------------------------------ |
| **資料庫連接失敗**      | Failed to connect to PostgreSQL     | 檢查 POSTGRES_HOST/PORT/USER/PASSWORD 環境變數；檢查 PostgreSQL 服務是否運行。 |
| **Prometheus 查詢失敗** | prometheus health check failed      | 確認 Prometheus 服務運行 (curl `http://localhost:9090/-/healthy`)。          |
| **量子服務不可用**      | quantum service health check failed | 檢查 `cyber-ai-quantum`服務運行狀態，或確認 `QUANTUM_URL`環境變數。        |
| **Nginx 配置失敗**      | 配置語法錯誤                        | 配置保存時會自動驗證；檢查配置路徑是否正確。                                   |

---

## 🔒 系統安全性與技術債務

| 項目                     | 狀態                 | 詳情                                                        | 來源 |
| ------------------------ | -------------------- | ----------------------------------------------------------- | ---- |
| **SAST 安全評分**  | **A (95/100)** | 已修復 67 個安全漏洞 (Critical/High/Medium)。               |      |
| **安全性**         | 零信任 + 量子增強    | 實施 mTLS 雙向認證；Zero Trust 量子-古典混合預測。          |      |
| **微服務複雜度**   | P0 追蹤中            | 多元件 (gRPC, RabbitMQ, 3+ 微服務) 協同，故障分析複雜度高。 |      |
| **單元測試覆蓋率** | 偏低                 | 文件顯示測試代碼僅 180+ 行，應提高至 80%+。                 |      |
| **性能驗證**       | P0 待實施            | 聲稱 500K req/s, < 2ms 延遲，需進行負載測試驗證。           |      |
