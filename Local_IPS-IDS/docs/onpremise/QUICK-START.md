# 快速入門指南 - 地端部署

> **版本**: v3.0.0  
> **分支**: dev  
> **更新**: 2025-10-09

---

## 🚀 3分鐘快速開始

### Windows 用戶

```powershell
# 1. 進入 Application 目錄
cd Application

# 2. 執行構建
.\build-local.ps1

# 3. 啟動服務
cd dist
.\start.bat

# 4. 訪問 http://localhost:3001
```

### Linux/macOS 用戶

```bash
# 1. 進入 Application 目錄
cd Application

# 2. 執行構建
chmod +x build-local.sh
./build-local.sh

# 3. 啟動服務
cd dist
chmod +x start.sh
./start.sh

# 4. 訪問 http://localhost:3001
```

---

## 📋 前置需求

- **Go** 1.24+
- **Node.js** 18+
- **PostgreSQL** 14+
- **Redis** 7+

---

## ✅ 驗證安裝

```bash
# 檢查服務
curl http://localhost:3001/api/v1/status

# 訪問介面
# http://localhost:3001        - 主介面
# http://localhost:3000        - Grafana
# http://localhost:9090        - Prometheus
```

---

**詳細說明**: 請參考 [Application/README.md](../../Application/README.md)

