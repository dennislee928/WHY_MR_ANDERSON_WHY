# Axiom Backend V3

> **企業級 IDS/IPS 統一管理平台**  
> **版本**: 3.1.0  
> **狀態**: 核心功能完成

---

## 🌟 特色功能

### 企業級核心功能 ✅
- ✅ **雙模式 Agent 架構** - External (mTLS) + Internal (直連)
- ✅ **四層智能儲存** - Hot (Redis) → Cold (PostgreSQL) → Archive (S3)
- ✅ **多法規合規引擎** - GDPR/PCI-DSS/HIPAA/SOX/ISO27001
- ✅ **PII 自動檢測與匿名化** - 6種類型，4種方法
- ✅ **GDPR 刪除權** - 完整工作流
- ✅ **完整性驗證** - SHA-256 Hash Chain

### 創新功能 ✅
- ✅ **時間旅行調試** - 系統快照、What-If 分析
- ✅ **自適應安全** - 動態風險評分、自動蜜罐
- ✅ **自癒系統** - AI 診斷、自動修復
- ✅ **統一可觀測性** - 跨服務整合、智能降噪
- ✅ **API 治理** - 健康評分、使用分析
- ✅ **技術債務追蹤** - 自動掃描、修復路線圖

---

## 🚀 快速開始

### 1. 前置需求

- Go 1.20+
- PostgreSQL 14+
- Redis 7+
- Docker (可選)

### 2. 安裝依賴

```bash
cd Application/be
go mod download
```

### 3. 配置環境變量

```bash
# 複製範例配置
cp .env.example .env

# 編輯配置
vim .env
```

### 4. 初始化資料庫

```bash
# 執行 Migrations
psql -U pandora -d pandora_db -f ../../database/migrations/001_initial_schema.sql
psql -U pandora -d pandora_db -f ../../database/migrations/002_agent_and_compliance_schema.sql
```

### 5. 啟動服務

```bash
go run cmd/server/main.go
```

服務將在 `http://localhost:3001` 啟動。

### 6. 測試 API

```bash
# 健康檢查
curl http://localhost:3001/health

# Agent 註冊
curl -X POST http://localhost:3001/api/v2/agent/register \
  -H "Content-Type: application/json" \
  -d '{
    "mode": "internal",
    "hostname": "test-server",
    "ip_address": "127.0.0.1",
    "capabilities": ["windows_logs"]
  }'
```

---

## 📚 API 文檔

### 完整 API 清單

**70+ 個端點** 分為以下類別：

1. **基礎服務** (24個)
   - Prometheus, Loki, Quantum, Nginx, Windows Logs

2. **Agent 管理** (11個)
   - 註冊、心跳、資產發現、合規檢查、遠端執行

3. **Storage 管理** (2個)
   - 統計、手動轉移

4. **Compliance** (10個)
   - PII 檢測、匿名化、GDPR 刪除、資料匯出

5. **組合功能** (7個)
   - 事件調查、性能分析、可觀測性、告警降噪

6. **創新功能** (16+個)
   - 時間旅行、自適應安全、自癒、API 治理等

詳細文檔請參閱: `docs/AXIOM-BACKEND-V3-API-DOCUMENTATION.md`

---

## 🏗️ 架構

```
┌─────────────────────────────────────────────────────────┐
│                    Axiom Backend V3                     │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐            │
│  │ Handlers │→ │ Services │→ │ Storage  │            │
│  └──────────┘  └──────────┘  └──────────┘            │
│       ↓              ↓              ↓                  │
│  ┌──────────────────────────────────────────┐         │
│  │         四層儲存架構                      │         │
│  │  Hot → Warm → Cold → Archive             │         │
│  │  Redis → Loki → PostgreSQL → S3/MinIO    │         │
│  └──────────────────────────────────────────┘         │
│                                                         │
│  ┌──────────────────────────────────────────┐         │
│  │         合規性引擎                        │         │
│  │  • PII 檢測   • 匿名化   • GDPR          │         │
│  │  • 審計追蹤   • 完整性驗證                │         │
│  └──────────────────────────────────────────┘         │
│                                                         │
└─────────────────────────────────────────────────────────┘
         ↑                                    ↑
    External Agents                     Internal Agents
    (via Nginx/mTLS)                    (Direct Connect)
```

---

## 📁 項目結構

```
Application/be/
├── cmd/server/              # 主程序
│   ├── main.go
│   └── routes.go
├── internal/
│   ├── agent/               # Phase 11: Agent 管理
│   ├── storage/             # Phase 12: 儲存層
│   ├── compliance/          # Phase 13: 合規性
│   ├── service/             # 業務邏輯層
│   ├── handler/             # HTTP 處理層
│   ├── model/               # 資料模型
│   ├── dto/                 # 請求結構
│   ├── vo/                  # 響應結構
│   ├── client/              # HTTP Client
│   ├── database/            # 資料庫管理
│   ├── cache/               # 快取管理
│   └── errors/              # 錯誤處理
├── go.mod
└── Makefile
```

---

## 🔧 開發指南

### 添加新的 API

1. **定義 Model** (如需要)
```go
// internal/model/your_model.go
type YourModel struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"type:varchar(255)"`
    CreatedAt time.Time
}
```

2. **定義 DTO/VO**
```go
// internal/dto/your_dto.go
type YourRequest struct {
    Name string `json:"name" binding:"required"`
}

// internal/vo/your_vo.go
type YourResponse struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}
```

3. **實現 Service**
```go
// internal/service/your_service.go
type YourService struct {
    db *database.Database
}

func (s *YourService) YourMethod(ctx context.Context) (*YourResponse, error) {
    // 業務邏輯
}
```

4. **實現 Handler**
```go
// internal/handler/your_handler.go
func (h *YourHandler) YourEndpoint(c *gin.Context) {
    // HTTP 處理
}
```

5. **註冊路由**
```go
// cmd/server/routes.go
v2.POST("/your-endpoint", yourHandler.YourEndpoint)
```

---

## 🧪 測試

### 單元測試
```bash
go test ./...
```

### API 測試
```bash
# 使用 Postman 或 curl
curl http://localhost:3001/api/v2/...
```

---

## 📦 部署

### Docker

```bash
# 構建鏡像
docker build -f ../docker/axiom-be-v3.dockerfile -t axiom-backend:v3 .

# 運行容器
docker run -p 3001:3001 \
  -e POSTGRES_HOST=postgres \
  -e REDIS_HOST=redis \
  axiom-backend:v3
```

### Docker Compose

```bash
cd ../..
docker-compose up axiom-be
```

---

## 📊 監控

### 健康檢查端點

```bash
GET /health
```

### 儲存統計

```bash
GET /api/v2/storage/tiers/stats
```

### Agent 健康

```bash
GET /api/v2/agent/health
```

---

## 🔐 安全考慮

### 已實施
1. ✅ mTLS 認證 (External Agents)
2. ✅ API Key 認證
3. ✅ CORS 配置
4. ✅ SQL 注入防護 (GORM)
5. ✅ 完整性驗證
6. ✅ 審計日誌

### 建議
1. 使用真實的 TLS 證書（生產環境）
2. 啟用 Rate Limiting
3. 配置 WAF
4. 定期安全掃描

---

## 📞 支援

### 文檔
- [API 文檔](../../docs/AXIOM-BACKEND-V3-API-DOCUMENTATION.md)
- [部署指南](../../docs/AXIOM-BACKEND-V3-DEPLOYMENT-GUIDE.md)
- [用戶手冊](../../docs/AXIOM-BACKEND-V3-USER-MANUAL.md)
- [Migration 指南](../../docs/AXIOM-BACKEND-V3-MIGRATION-GUIDE.md)

### 問題回報
請在項目 Issues 中提交問題。

---

## 📝 變更日誌

### v3.1.0 (2025-10-16)
- ✅ 實施 Phase 11: Agent 進階架構
- ✅ 實施 Phase 12: 四層儲存架構 (70%)
- ✅ 實施 Phase 13: 合規性引擎
- ✅ 新增 30+ API 端點
- ✅ 新增完整性驗證機制

### v3.0.0 (2025-10-15)
- ✅ 基礎架構完成
- ✅ 核心 API 實施
- ✅ Frontend 整合
- ✅ 文檔生成

---

## 📄 授權

本項目為內部專案。

---

**專案狀態**: 🟢 核心功能完成，進入測試階段  
**維護者**: Axiom Development Team  
**最後更新**: 2025-10-16
