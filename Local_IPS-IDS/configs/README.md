# Pandora Box Console - 配置檔案

此目錄包含所有服務的配置檔案。

## 📁 目錄結構

```
configs/
├── agent-config.yaml              # Agent 主配置
├── agent-config.yaml.template     # Agent 配置範本
├── console-config.yaml            # Console 配置
├── console-config.yaml.template   # Console 配置範本
├── ui-config.yaml.template        # UI Server 配置範本
├── grafana/                       # Grafana 配置
│   ├── dashboards/
│   ├── grafana.ini
│   └── provisioning/
├── prometheus/                    # Prometheus 配置
│   ├── prometheus.yml
│   └── rules/
├── loki.yaml                      # Loki 配置
├── alertmanager.yml               # AlertManager 配置
├── nginx/                         # Nginx 配置
│   ├── nginx-paas.conf
│   └── default-paas.conf
├── postgres/                      # PostgreSQL 初始化
│   └── init.sql
├── promtail-paas.yaml            # Promtail 配置
├── supervisord-*.conf             # Supervisord 配置（PaaS）
└── README.md                      # 本檔案
```

## 🎯 地端部署使用的配置

### 核心配置

- `agent-config.yaml` - Agent 主程式配置
- `console-config.yaml` - Console API 配置
- `ui-config.yaml.template` - UI Server 配置（需複製並自訂）

### 監控配置

- `grafana/` - Grafana 儀表板和數據源
- `prometheus/prometheus.yml` - Prometheus 抓取配置
- `loki.yaml` - Loki 日誌聚合配置
- `alertmanager.yml` - 告警管理配置

### 資料庫配置

- `postgres/init.sql` - PostgreSQL 初始化腳本

## 🚀 快速配置

### 1. 複製範本

```bash
cp agent-config.yaml.template agent-config.yaml
cp console-config.yaml.template console-config.yaml
cp ui-config.yaml.template ui-config.yaml
```

### 2. 編輯配置

編輯複製的配置檔案，設定：
- 資料庫連線字串
- Redis 連線資訊
- 裝置端口
- 日誌等級
- 監控端點

### 3. 環境變數

也可以使用環境變數覆蓋配置：

```bash
export LOG_LEVEL=debug
export DEVICE_PORT=/dev/ttyUSB0
export DATABASE_URL=postgresql://user:pass@localhost:5432/pandora
```

## 📝 配置說明

詳細的配置說明請參考各配置檔案中的註解。

---

**維護者**: Pandora Security Team  
**最後更新**: 2025-10-09

