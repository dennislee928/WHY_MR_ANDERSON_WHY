# 地端部署配置

此目錄包含地端（On-Premise）部署所需的所有配置檔案。

## 📁 目錄內容

```
onpremise/
├── docker-compose.yml         # Docker Compose 主配置
├── docker-compose.test.yml    # 測試環境配置
└── README.md                  # 本檔案
```

## 🚀 使用 Docker Compose 部署

### 前置需求

- Docker 20.10+
- Docker Compose 2.0+

### 快速開始

```bash
# 1. 進入此目錄
cd deployments/onpremise

# 2. 啟動所有服務
docker-compose up -d

# 3. 檢查服務狀態
docker-compose ps

# 4. 查看日誌
docker-compose logs -f
```

### 停止服務

```bash
docker-compose down
```

### 重新構建

```bash
docker-compose build --no-cache
docker-compose up -d
```

## 📊 服務端口

| 服務 | 端口 | 說明 |
|------|------|------|
| Frontend | 3001 | Web UI |
| Grafana | 3000 | 監控儀表板 |
| Prometheus | 9090 | 指標收集 |
| Loki | 3100 | 日誌聚合 |
| Agent API | 8080 | Agent API |
| PostgreSQL | 5432 | 資料庫 |
| Redis | 6379 | 快取 |

## 🔧 自訂配置

編輯 `docker-compose.yml` 修改：
- 服務端口
- 環境變數
- 資源限制
- 儲存卷位置

## 📝 環境變數

創建 `.env` 檔案：

```bash
# 資料庫
POSTGRES_PASSWORD=your_secure_password
POSTGRES_DB=pandora

# Redis
REDIS_PASSWORD=your_redis_password

# 應用程式
LOG_LEVEL=info
DEVICE_PORT=/dev/ttyUSB0
```

---

**維護者**: Pandora Security Team  
**最後更新**: 2025-10-09

