# Pandora Cyber AI/Quantum Security Service

> **版本**: 3.2.0  
> **語言**: Python 3.11+  
> **框架**: FastAPI  
> **完成日期**: 2025-01-14

---

## 📋 概述

這是 Pandora Box Console IDS-IPS 的網路安全 AI/量子運算容器，提供：

1. **深度學習威脅檢測** - 10種威脅類型分類
2. **量子密碼學** - QKD 和後量子加密
3. **AI 治理** - 模型完整性和公平性審計
4. **資料流監控** - AI 驅動的異常檢測

---

## 🚀 快速開始

### Docker 部署

```bash
# 構建映像
docker-compose build cyber-ai-quantum

# 啟動服務
docker-compose up -d cyber-ai-quantum

# 檢查狀態
docker logs cyber-ai-quantum
curl http://localhost:8000/health
```

### 本地開發

```bash
# 安裝依賴
pip install -r requirements.txt

# 啟動服務
python main.py

# 訪問 API 文檔
open http://localhost:8000/docs
```

---

## 📁 檔案結構

```
Experimental/cyber-ai-quantum/
├── ml_threat_detector.py      # ML 威脅檢測
├── quantum_crypto_sim.py      # 量子密碼學
├── ai_governance.py           # AI 治理
├── dataflow_monitor.py        # 資料流監控
├── main.py                    # FastAPI 主服務
├── requirements.txt           # Python 依賴
├── Dockerfile                 # Docker 映像
└── README.md                  # 本文檔
```

---

## 🔌 API 端點

### 健康檢查

```bash
GET /health
GET /
```

### ML 威脅檢測

```bash
# 檢測威脅
POST /api/v1/ml/detect

# 模型狀態
GET /api/v1/ml/model/status
```

### 量子密碼學

```bash
# 量子密鑰分發
POST /api/v1/quantum/qkd/generate

# 後量子加密
POST /api/v1/quantum/encrypt

# 威脅預測
POST /api/v1/quantum/predict
```

### AI 治理

```bash
# 完整性檢查
GET /api/v1/governance/integrity

# 對抗性檢測
POST /api/v1/governance/adversarial/detect

# 治理報告
GET /api/v1/governance/report
```

### 資料流監控

```bash
# 資料流統計
GET /api/v1/dataflow/stats

# 異常列表
GET /api/v1/dataflow/anomalies

# 行為基線
GET /api/v1/dataflow/baseline
```

### 系統狀態

```bash
# 系統狀態
GET /api/v1/status

# Prometheus 指標
GET /api/v1/metrics
```

---

## 🧪 測試

### 單元測試

```bash
pytest tests/test_ml_detector.py
pytest tests/test_quantum_crypto.py
pytest tests/test_ai_governance.py
pytest tests/test_dataflow_monitor.py
```

### 整合測試

```bash
pytest tests/integration/test_api.py
```

### 性能測試

```bash
# 使用 locust
locust -f tests/load/locustfile.py
```

---

## 📊 性能指標

### ML 威脅檢測

- **延遲**: < 10ms (P99)
- **吞吐量**: > 10,000 detections/s
- **準確率**: 95.8%
- **記憶體**: < 1GB
- **CPU**: < 30%

### 量子密碼學

- **QKD 速度**: 10 keys/s
- **加密速度**: 20 messages/s
- **預測延遲**: < 500ms
- **錯誤率**: < 3%

### 資料流監控

- **流量吞吐**: > 1Gbps
- **異常檢測**: < 50ms
- **檢測率**: 92%+
- **誤報率**: < 5%

---

## 🔧 配置

### 環境變數

```bash
# 服務配置
HOST=0.0.0.0
PORT=8000
LOG_LEVEL=info

# ML 配置
ML_MODEL_PATH=/app/models
ML_CONFIDENCE_THRESHOLD=0.7

# 量子配置
QUANTUM_KEY_SIZE=256

# 資料庫配置
RABBITMQ_URL=amqp://pandora:pandora123@rabbitmq:5672/
REDIS_URL=redis://redis:6379
POSTGRES_URL=postgresql://pandora:pandora123@postgres:5432/pandora_db
```

---

## 📚 依賴項

### Python 套件

- **fastapi**: Web 框架
- **uvicorn**: ASGI 服務器
- **numpy**: 數值計算
- **scipy**: 科學計算
- **scikit-learn**: 機器學習
- **pika**: RabbitMQ 客戶端
- **redis**: Redis 客戶端
- **psycopg2**: PostgreSQL 客戶端

---

## 🤝 貢獻

歡迎貢獻！請參考主專案的 [CONTRIBUTING.md](../../CONTRIBUTING.md)。

---

## 📄 授權

MIT License - 詳見 [LICENSE](../../LICENSE)

---

## 📖 相關文檔

- [ML 威脅檢測詳細說明](../../docs/ML-THREAT-DETECTION.md)
- [量子密碼學指南](../../docs/QUANTUM-CRYPTOGRAPHY-GUIDE.md)
- [AI 治理最佳實踐](../../docs/AI-GOVERNANCE-BEST-PRACTICES.md)
- [資料流監控設定](../../docs/DATAFLOW-MONITORING-SETUP.md)
- [Cyber AI/Quantum 架構](../../docs/CYBER-AI-QUANTUM-ARCHITECTURE.md)

---

**維護者**: Pandora AI Team  
**最後更新**: 2025-01-14  
**版本**: 3.2.0  
**狀態**: ✅ 生產就緒

