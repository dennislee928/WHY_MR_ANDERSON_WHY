這是一個關於 **Pandora Box Console IDS-IPS 系統** 的詳細介紹，特別著重於其微服務架構和多層次的資安防護措施。

## 系統介紹：Pandora Box Console IDS-IPS

Pandora Box Console IDS-IPS 系統採用微服務（Microservice）架構，旨在提供完整的入侵偵測與預防系統（IDS-IPS）主控台。在部署方面，系統提供多種方案，其中一套是基於五個不同的 PaaS 平台、實現零基礎設施成本的部署策略。

這個部署方案旨在實現高可用性、完整的監控系統，以及高度的自動化部署能力。

---

## I. 系統架構與部署策略（PaaS 微服務模型）

本系統的核心部署策略是將 10 個微服務組件分散到 5 個不同的 PaaS 平台上，以達到分散風險並利用各平台免費資源的目標。

### 1. 核心元件與平台分佈

| 微服務 (Microservice) | 建議平台 | 部署方式 / 理由 | 核心功能 | 來源 |
| :--- | :--- | :--- | :--- | :--- |
| **pandora-agent** (核心後端) | **Koyeb** | Web Service (Docker) / 提供 2 個永不休眠 (Always-on) 的 Nano 容器，適合 24/7 運行 | 核心 API 服務、Go 後端 | |
| **promtail** (日誌收集) | **Koyeb** | Sidecar 隨 Agent 部署 / 將日誌推送到 Loki | 日誌收集 | |
| **axiom-ui** (UI 伺服器/前端) | **Patr.io** | Container Deployment / 分離使用者介面流量與後端服務 | 前端使用者介面 (Next.js) | |
| **nginx** (反向代理) | **Render** | Web Service / 作為統一流量入口，路由請求至 UI 和 Grafana | 統一路由、流量分離 | |
| **postgres** (資料庫) | **Railway.app** | Managed Database / 提供穩定的核心數據基礎 | 核心數據儲存 | |
| **redis** (快取) | **Render** | Managed Redis / 提供快取層，分離部署以分散風險 | 快取、Pub/Sub 事件系統 |
| **Prometheus** (監控指標) | **Fly.io** | App (Docker) / 需要永久儲存空間 (Persistent Volume) 存放時間序列數據 | 系統指標收集 | |
| **Loki** (日誌聚合) | **Fly.io** | App (Docker) / 需要永久儲存空間 | 日誌聚合與儲存 | |
| **Grafana** (視覺化) | **Fly.io** | App (Docker) / 需頻繁與 Prometheus 和 Loki 通訊，部署在同一平台可降低網路延遲 | 儀表板視覺化 | |
| **Alertmanager** (告警) | **Fly.io** | App (Docker) / 必須全時運行 | 告警通知系統 | |

### 2. 關鍵架構特點

*   **多平台整合**：系統成功整合了 Railway、Render、Koyeb、Patr.io 和 Fly.io 共 5 個不同的 PaaS 平台，並配置了 10 個微服務組件，且所有服務使用免費方案實現零基礎設施成本 ($0/月)。
*   **監控系統集中**：完整的監控與告警系統（Prometheus, Loki, Grafana, Alertmanager）集中部署在 **Fly.io** 上，以確保它們可以頻繁通訊並降低網路延遲。
*   **持久化儲存**：Fly.io 是少數在免費方案中提供數 GB 永久儲存空間的平台，因此 Prometheus 和 Loki 利用 Fly.io 的 Volumes 配置來存放時間序列數據和日誌。為了符合 Fly.io 免費方案每個應用只能有 1 個 Volume 的限制，系統使用單一 Volume 存放所有監控數據（例如 10GB）。
*   **容器多進程管理**：在 Koyeb（Agent + Promtail）和 Fly.io（監控系統）中，採用 **Supervisor** 來管理單一容器內的多個服務進程。
*   **流量統一入口**：Nginx 作為反向代理部署在 Render 上，統一外部流量入口，並將請求導向 Patr.io 上的 UI 服務和 Fly.io 上的 Grafana。
*   **事件驅動**：系統內建 Pub/Sub 事件系統，支援 Redis Pub/Sub 實作和 In-Memory Pub/Sub，用於處理認證、安全、設備和系統事件。
*   **IoT 支援**：核心後端支援 **MQTT 協定**，提供完整的 MQTT v3.1.1 協定支援，包括 QoS 0/1/2、TLS/SSL 加密通訊和自動重連機制。

---

## II. 視覺化架構圖描述

由於無法直接提供圖形，以下為 Pandora Box Console IDS-IPS 系統的 PaaS 微服務架構的邏輯描述：

1.  **使用者存取層 (Entry Point)**：
    *   外部使用者首先透過 **Render 上的 Nginx 反向代理** 進入系統。
    *   Nginx 處理路由，將前端 UI 請求導向 **Patr.io 上的 Axiom UI**。
    *   Nginx 也將監控儀表板的請求導向 **Fly.io 上的 Grafana**。

2.  **核心業務層 (Backend/Cache/DB)**：
    *   核心後端 **Pandora Agent** 運行在 **Koyeb** 的永不休眠容器中，處理 API 請求、業務邏輯、速率限制和安全事件。
    *   Agent 依賴 **Railway.app 上的 PostgreSQL 資料庫** 進行核心數據持久化儲存。
    *   Agent 依賴 **Render 上的 Redis** 提供快取和 Pub/Sub 事件處理。

3.  **可觀察性層 (Observability)**：
    *   **Promtail** 作為 Sidecar 隨 **Koyeb 上的 Agent** 一同運行，負責收集 Agent 的日誌，並透過公網 URL 將日誌推送到 **Fly.io 上的 Loki**。
    *   **Prometheus** 也部署在 **Fly.io** 上，負責收集系統和服務的指標數據。
    *   **Grafana** 部署在 **Fly.io** 上，從 Prometheus 取得指標數據，從 Loki 取得日誌數據，進行視覺化呈現。
    *   **Alertmanager** 與監控系統一同部署在 **Fly.io** 上，負責處理告警通知。

---

## III. 資安與防護措施

本系統實施了多層次的資安考量，涵蓋網路、資料、存取控制以及應用程式層的防護。

### 1. 應用程式層安全 (Agent 核心功能)

*   **速率限制與暴力攻擊防護 (Rate Limiting & Brute Force Protection)**：
    *   採用 **Token Bucket** 演算法實施速率限制，例如設定每秒 10 個請求，突發容量 20。
    *   實施帳號鎖定機制：**5 次失敗嘗試後鎖定 10 分鐘**。
    *   實施 **IP 封鎖機制**：觸發鎖定後，封鎖該 IP 地址 1 小時。
    *   支援動態速率調整和白名單功能。
*   **API 安全**：
    *   實施 CORS 保護和 Security Headers。
    *   支援 PSK（Pre-Shared Key）認證和時間戳驗證。
    *   具備 TPM（Trusted Platform Module）硬體認證功能。
*   **MQTT 安全**：
    *   支援 **TLS/SSL 加密**所有通訊。
    *   使用用戶名/密碼和客戶端證書進行認證。

### 2. 基礎設施與網路安全

*   **通訊加密**：強制使用 **TLS 加密**所有通訊 (HTTPS 強制加密)，確保資料傳輸安全。
*   ** Secrets 管理**：環境變數與 Secrets 嚴格分離，以保護敏感配置信息。GitHub Actions 流程中也需要設定多個 GCP/OCI 或 Vercel 配置密鑰。
*   **資料安全**：
    *   加密敏感資料。
    *   定期備份資料庫。
    *   建議使用強密碼並定期輪換密鑰。
*   **存取控制**：
    *   實施 **最小權限原則**。
    *   在 Kubernetes 環境（GCP/OCI 部署）中，使用 **RBAC**（Role-Based Access Control）控制存取。
    *   限制管理端口的訪問，並設定適當的防火牆規則。
    *   Grafana 設定適當的使用者權限，預設登入資訊為 `admin / pandora123`（但在生產環境中應使用強密碼）。

總體而言，系統在設計時就將安全性納入考量，並確保在生產環境中使用強密碼和適當的安全配置。
```mermaid
graph TD
    subgraph PaaS服務分佈
        direction LR
        subgraph Render_Patrio [流量入口與前端 Render / Patr.io]
            Nginx(Nginx Render 反向代理)
            UI(Axiom UI Patr.io 前端)
        end

        subgraph Koyeb_Core [核心業務層 Koyeb / Railway / Render]
            Agent(Pandora Agent Koyeb Go 後端)
            Promtail(Promtail Koyeb Sidecar)
            DB(PostgreSQL Railway.app 資料庫)
            Cache(Redis Render 快取/PubSub)
        end

        subgraph Flyio_Monitoring [集中監控與告警 Fly.io]
            Prometheus(Prometheus 指標收集)
            Loki(Loki 日誌聚合)
            Grafana(Grafana 儀表板視覺化)
            AlertManager(Alertmanager 告警系統)
        end
        
        Note(["*PV: Fly.io 提供持久化儲存 Volume"])
    end
    
    User(外部使用者/IoT 裝置)
    
    %% 流量流向
    User -->|HTTPS/MQTT 流量| Nginx
    Nginx -->|路由 UI 請求| UI
    Nginx -->|路由 Grafana| Grafana
    UI -->|API 請求| Agent
    
    %% 核心業務流
    Agent -->|數據持久化| DB
    Agent -->|快取/事件| Cache
    
    %% 可觀察性流
    Agent -->|指標抓取| Prometheus
    Promtail -->|日誌推送| Loki
    
    %% 監控內部流
    Grafana -->|查詢 Metrics| Prometheus
    Grafana -->|查詢 Logs| Loki
    
    Prometheus -->|發送告警| AlertManager

    %% 顏色提示
    style Nginx fill:#E6F0FF, stroke:#0066FF
    style UI fill:#E6F0FF, stroke:#0066FF
    style Agent fill:#D6FFE6, stroke:#00AA00
    style Promtail fill:#D6FFE6, stroke:#00AA00
    style DB fill:#FFF0CC, stroke:#FFAA00
    style Cache fill:#FFF0CC, stroke:#FFAA00
    style Prometheus fill:#FFD6D6, stroke:#FF0000
    style Loki fill:#FFD6D6, stroke:#FF0000
    style Grafana fill:#FFD6D6, stroke:#FF0000
    style AlertManager fill:#FFD6D6, stroke:#FF0000
    style User fill:#f9f,stroke:#333
    ```