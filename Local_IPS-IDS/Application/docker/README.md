# Application Docker 映像檔

此目錄包含所有用於地端部署的 Docker 映像定義檔（Dockerfiles）。

## 📋 Dockerfile 清單

### 核心服務

| Dockerfile | 服務 | 說明 |
|------------|------|------|
| `agent.dockerfile` | Pandora Agent | 主要 Agent 程式 |
| `agent.koyeb.dockerfile` | Pandora Agent (Koyeb) | Koyeb 優化版本 |
| `server-be.dockerfile` | Backend API | 後端 API 伺服器 |
| `ui.patr.dockerfile` | UI Server | UI 伺服器 |

### 監控服務

| Dockerfile | 服務 | 說明 |
|------------|------|------|
| `monitoring.dockerfile` | Monitoring Stack | Prometheus+Grafana+Loki 整合 |

### 輔助服務

| Dockerfile | 服務 | 說明 |
|------------|------|------|
| `nginx.dockerfile` | Nginx | 反向代理 |
| `server-fe.dockerfile` | Frontend | 前端靜態伺服器 |
| `test.dockerfile` | Test Runner | 測試環境 |

**總計**: 8 個 Dockerfiles

---

## 🚀 使用方式

這些 Dockerfiles 由 `docker-compose.yml` 自動使用：

```bash
# 在 Application/ 目錄
docker-compose up -d
```

或使用啟動腳本：

```bash
# Windows
.\docker-start.ps1

# Linux/macOS
./docker-start.sh
```

---

## 🔧 服務架構

```
Application/
├── docker-compose.yml        # 服務編排
├── docker/                   # Dockerfiles（本目錄）
│   ├── agent.dockerfile      → pandora-agent 服務
│   ├── ui.patr.dockerfile    → axiom-ui 服務
│   └── ...
├── .env                      # 環境變數
└── [其他目錄]
```

---

## 📊 服務列表（docker-compose.yml）

| 服務 | 映像來源 | 端口 |
|------|----------|------|
| pandora-agent | agent.dockerfile | 8080 |
| axiom-ui | ui.patr.dockerfile | 3001 |
| prometheus | prom/prometheus:v2.47.0 | 9090 |
| grafana | grafana/grafana:10.2.0 | 3000 |
| loki | grafana/loki:2.9.2 | 3100 |
| alertmanager | prom/alertmanager:v0.26.0 | 9093 |
| postgres | postgres:15-alpine | 5432 |
| redis | redis:7.2-alpine | 6379 |
| nginx | nginx:1.25-alpine | 80, 443 |
| promtail | grafana/promtail:2.9.2 | - |
| node-exporter | prom/node-exporter:v1.6.1 | 9100 |

**總計**: 11 個服務

---

**維護**: Pandora Security Team
**最後更新**: 2025-10-09

