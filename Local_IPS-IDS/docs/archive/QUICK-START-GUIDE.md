# Pandora Box Console - 快速入門指南

> **版本**: v3.0.0 (On-Premise)  
> **分支**: dev  
> **更新**: 2025-10-09

---

## 🚀 3分鐘快速開始

### 方式 1: 本地構建（開發/測試）⭐

**Windows**:
```powershell
# 1. 進入 Application 目錄
cd Application

# 2. 執行構建（需要 Go 1.24+ 和 Node.js 18+）
.\build-local.ps1 -Version "3.0.0"

# 3. 啟動所有服務
cd dist
.\start.bat

# 4. 訪問 Web 介面
# http://localhost:3001
```

**Linux/macOS**:
```bash
# 1. 進入 Application 目錄
cd Application

# 2. 執行構建（需要 Go 1.24+ 和 Node.js 18+）
chmod +x build-local.sh
./build-local.sh all "3.0.0"

# 3. 啟動所有服務
cd dist
chmod +x start.sh
./start.sh

# 4. 訪問 Web 介面
# http://localhost:3001
```

---

### 方式 2: 使用安裝檔（生產環境）

#### 步驟 1: 下載安裝檔

從 [GitHub Releases](https://github.com/your-org/pandora_box_console_IDS-IPS/releases) 下載：

| 作業系統 | 檔案 |
|----------|------|
| Windows | `pandora-box-console-*-windows-amd64-setup.exe` |
| Ubuntu/Debian | `pandora-box-console_*_amd64.deb` |
| RHEL/CentOS | `pandora-box-console-*.rpm` |
| 虛擬機 | `pandora-box-console-*.ova` |
| 通用 | `pandora-box-console-*-amd64.iso` |

#### 步驟 2: 安裝

**Windows**:
```powershell
# 雙擊執行安裝程式
pandora-box-console-3.0.0-windows-amd64-setup.exe

# 按照安裝精靈操作
# 完成後從開始選單啟動
```

**Ubuntu/Debian**:
```bash
sudo dpkg -i pandora-box-console_3.0.0_amd64.deb
sudo systemctl start pandora-agent
```

**RHEL/CentOS**:
```bash
sudo rpm -i pandora-box-console-3.0.0.rpm
sudo systemctl start pandora-agent
```

#### 步驟 3: 訪問

開啟瀏覽器訪問: **http://localhost:3001**

---

## 📋 系統需求

### 最低配置

- **CPU**: 2 核心
- **RAM**: 4 GB
- **儲存**: 20 GB
- **OS**: Windows 10+, Ubuntu 20.04+, Debian 11+, RHEL 8+

### 軟體依賴

- **PostgreSQL** 14+（必須）
- **Redis** 7+（必須）
- **Go** 1.24+（僅開發）
- **Node.js** 18+（僅開發）

---

## 🎯 常用命令

### 開發命令

```bash
# 前端開發
cd Application/Fe
npm install
npm run dev           # http://localhost:3001

# 後端開發  
cd Application/be
make all
make run-agent
```

### 構建命令

```powershell
# Windows完整構建
cd Application
.\build-local.ps1

# Linux/macOS完整構建
cd Application
./build-local.sh
```

### 服務管理（Linux）

```bash
# 啟動
sudo systemctl start pandora-agent

# 停止
sudo systemctl stop pandora-agent

# 狀態
sudo systemctl status pandora-agent

# 日誌
sudo journalctl -u pandora-agent -f
```

---

## 🔍 驗證安裝

### 檢查服務

```bash
# 檢查服務狀態
curl http://localhost:3001/api/v1/status

# 檢查 Prometheus
curl http://localhost:9090/-/healthy

# 檢查 Grafana
curl http://localhost:3000/api/health
```

### 訪問介面

- **主介面**: http://localhost:3001
- **Grafana**: http://localhost:3000
- **Prometheus**: http://localhost:9090
- **API 文檔**: http://localhost:8080/swagger

---

## 🆘 快速故障排除

### 問題 1: 端口被占用

```bash
# Windows
netstat -ano | findstr ":3001"
taskkill /PID <PID> /F

# Linux/macOS
lsof -ti:3001 | xargs kill -9
```

### 問題 2: 無法連接資料庫

```bash
# 啟動 PostgreSQL
sudo systemctl start postgresql

# 測試連接
psql -U postgres -h localhost
```

### 問題 3: 前端無法啟動

```bash
cd Application/Fe
rm -rf node_modules package-lock.json .next
npm cache clean --force
npm install
npm run dev
```

### 問題 4: 後端編譯失敗

```bash
cd <專案根目錄>
go clean -modcache
go mod download
go mod tidy

cd Application/be
make clean
make all
```

---

## 📚 延伸閱讀

- [完整部署指南](ONPREMISE-DEPLOYMENT-GUIDE.md)
- [專案結構說明](README-PROJECT-STRUCTURE.md)
- [前端開發指南](Application/Fe/README.md)
- [後端開發指南](Application/be/README.md)
- [測試清單](TESTING-CHECKLIST.md)
- [重構報告](RESTRUCTURE-FINAL-REPORT.md)

---

## ✅ 成功安裝的標誌

如果您看到以下內容，表示安裝成功：

1. ✅ 瀏覽器可以訪問 http://localhost:3001
2. ✅ 看到 "Pandora Box Console" 儀表板
3. ✅ 系統狀態顯示為"線上"
4. ✅ 無錯誤訊息

---

**需要幫助？**
- 📧 support@pandora-ids.com
- 🐛 [GitHub Issues](https://github.com/your-org/pandora_box_console_IDS-IPS/issues)
- 💬 [Discussions](https://github.com/your-org/pandora_box_console_IDS-IPS/discussions)

---

**版本**: v3.0.0  
**最後更新**: 2025-10-09  
🎉 **歡迎使用 Pandora Box Console！**

