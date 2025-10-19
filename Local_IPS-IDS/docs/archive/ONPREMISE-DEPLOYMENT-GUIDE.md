# Pandora Box Console IDS-IPS - 地端部署指南

> **版本**: v3.0 (dev 分支)  
> **更新日期**: 2025-10-09  
> **狀態**: ✅ 完成

---

## 📋 目錄

1. [概述](#概述)
2. [系統需求](#系統需求)
3. [快速開始](#快速開始)
4. [安裝方法](#安裝方法)
5. [本地構建](#本地構建)
6. [CI/CD 自動化](#cicd-自動化)
7. [故障排除](#故障排除)
8. [常見問題](#常見問題)

---

## 概述

本指南說明如何在**地端（On-Premise）環境**部署 Pandora Box Console IDS-IPS 系統。

### 🆕 新功能（v3.0）

- ✅ **Application/** 目錄結構 - 統一的應用程式管理
- ✅ **本地構建腳本** - Windows 和 Linux/macOS 一鍵構建
- ✅ **安裝檔生成** - 支援 .exe, .deb, .rpm, .iso, .ova 格式
- ✅ **CI/CD 自動化** - GitHub Actions 自動構建和發布
- ✅ **啟動腳本** - 自動生成的服務啟動/停止腳本

### 部署架構

```
本地伺服器
├── Frontend (Port 3001)     - Next.js Web UI
├── Backend (Port 8080)      - Go 後端服務
│   ├── Pandora Agent
│   ├── Console API
│   └── UI Server
├── Monitoring               - 監控系統
│   ├── Prometheus (9090)
│   ├── Grafana (3000)
│   └── Loki (3100)
└── Storage                  - 資料儲存
    ├── PostgreSQL (5432)
    └── Redis (6379)
```

---

## 系統需求

### 最低配置

| 項目 | 需求 |
|------|------|
| **CPU** | 2 核心 |
| **記憶體** | 4 GB RAM |
| **儲存空間** | 20 GB |
| **作業系統** | Windows 10+, Ubuntu 20.04+, Debian 11+, RHEL 8+ |

### 建議配置

| 項目 | 建議 |
|------|------|
| **CPU** | 4 核心以上 |
| **記憶體** | 8 GB RAM 以上 |
| **儲存空間** | 50 GB 以上（含日誌和資料） |
| **網路** | 千兆網卡 |

### 軟體依賴

#### 必要依賴
- **PostgreSQL** 14+
- **Redis** 7+

#### 開發環境額外需求
- **Go** 1.24+
- **Node.js** 18+
- **Git**
- **Make** (Linux/macOS) 或 **PowerShell** (Windows)

---

## 快速開始

### 方式 1: 使用預建安裝檔（最簡單）⭐

1. 從 [GitHub Releases](https://github.com/your-org/pandora_box_console_IDS-IPS/releases) 下載安裝檔

2. 依照您的作業系統選擇：

   **Windows**:
   ```powershell
   # 執行安裝程式
   pandora-box-console-1.0.0-windows-amd64-setup.exe
   ```

   **Ubuntu/Debian**:
   ```bash
   sudo dpkg -i pandora-box-console_1.0.0_amd64.deb
   sudo systemctl start pandora-agent
   ```

   **RedHat/CentOS**:
   ```bash
   sudo rpm -i pandora-box-console-1.0.0-1.x86_64.rpm
   sudo systemctl start pandora-agent
   ```

3. 訪問 Web 介面: http://localhost:3001

### 方式 2: 使用本地構建腳本

1. 克隆專案
   ```bash
   git clone https://github.com/your-org/pandora_box_console_IDS-IPS.git
   cd pandora_box_console_IDS-IPS
   git checkout dev
   ```

2. 執行構建腳本

   **Windows**:
   ```powershell
   cd Application
   .\build-local.ps1
   cd dist
   .\start.bat
   ```

   **Linux/macOS**:
   ```bash
   cd Application
   chmod +x build-local.sh
   ./build-local.sh
   cd dist
   ./start.sh
   ```

3. 訪問 Web 介面: http://localhost:3001

---

## 安裝方法

### Windows 安裝程式 (.exe)

#### 特點
- Inno Setup 精美安裝精靈
- 自動配置系統服務
- 開始選單捷徑
- 自動卸載功能

#### 安裝步驟

1. **下載安裝程式**
   - 檔案名稱: `pandora-box-console-{version}-windows-amd64-setup.exe`
   - 來源: [GitHub Releases](https://github.com/your-org/pandora_box_console_IDS-IPS/releases)

2. **執行安裝**
   - 雙擊執行安裝程式
   - 按照安裝精靈指示操作
   - 選擇安裝路徑（預設: `C:\Program Files\PandoraBox`）
   - 完成安裝

3. **啟動服務**
   - 從開始選單找到 "Pandora Box Console"
   - 或執行: `C:\Program Files\PandoraBox\pandora-agent.exe`

4. **訪問 Web 介面**
   - 開啟瀏覽器訪問: http://localhost:3001

### Linux 套件 (.deb)

#### 適用系統
- Ubuntu 20.04, 22.04
- Debian 11, 12

#### 安裝步驟

```bash
# 1. 下載 .deb 套件
wget https://github.com/your-org/pandora_box_console_IDS-IPS/releases/download/v1.0.0/pandora-box-console_1.0.0_amd64.deb

# 2. 安裝套件
sudo dpkg -i pandora-box-console_1.0.0_amd64.deb

# 3. 安裝依賴（如有缺失）
sudo apt-get install -f

# 4. 啟動服務
sudo systemctl start pandora-agent
sudo systemctl enable pandora-agent

# 5. 檢查狀態
sudo systemctl status pandora-agent

# 6. 查看日誌
sudo journalctl -u pandora-agent -f
```

### Linux 套件 (.rpm)

#### 適用系統
- Red Hat Enterprise Linux 8, 9
- CentOS 8, 9
- Fedora 36+

#### 安裝步驟

```bash
# 1. 下載 .rpm 套件
wget https://github.com/your-org/pandora_box_console_IDS-IPS/releases/download/v1.0.0/pandora-box-console-1.0.0-1.x86_64.rpm

# 2. 安裝套件
sudo rpm -i pandora-box-console-1.0.0-1.x86_64.rpm

# 3. 啟動服務
sudo systemctl start pandora-agent
sudo systemctl enable pandora-agent

# 4. 檢查狀態
sudo systemctl status pandora-agent
```

### ISO 安裝光碟

#### 特點
- 可開機安裝光碟
- 包含所有必要檔案
- 支援離線安裝

#### 使用步驟

```bash
# 1. 掛載 ISO 映像
sudo mkdir -p /mnt/pandora
sudo mount -o loop pandora-box-console-1.0.0-amd64.iso /mnt/pandora

# 2. 執行安裝腳本
cd /mnt/pandora
sudo ./install.sh

# 3. 卸載 ISO
cd ~
sudo umount /mnt/pandora

# 4. 啟動服務
sudo systemctl start pandora-agent
```

### OVA 虛擬機映像

#### 特點
- 預配置的完整虛擬機
- 開箱即用
- 支援 VirtualBox 和 VMware

#### 使用步驟

1. **匯入 OVA 到虛擬化平台**

   **VirtualBox**:
   - 檔案 → 匯入應用裝置
   - 選擇 `.ova` 檔案
   - 調整虛擬機設定（記憶體、CPU）
   - 匯入

   **VMware**:
   - 檔案 → 開啟
   - 選擇 `.ova` 檔案
   - 完成匯入

2. **啟動虛擬機**
   - 預設使用者名稱: `pandora`
   - 預設密碼: `pandora`

3. **服務已自動啟動**
   - 訪問 Web 介面: http://{VM_IP}:3001
   - 訪問 Grafana: http://{VM_IP}:3000

---

## 本地構建

### 使用自動構建腳本（推薦）

#### Windows (PowerShell)

```powershell
cd Application

# 基本構建
.\build-local.ps1

# 指定版本
.\build-local.ps1 -Version "1.0.0"

# 只構建後端
.\build-local.ps1 -SkipFrontend

# 只構建前端
.\build-local.ps1 -SkipBackend

# 清理後重新構建
.\build-local.ps1 -Clean
```

#### Linux/macOS (Bash)

```bash
cd Application

# 基本構建
./build-local.sh

# 指定版本
./build-local.sh all "1.0.0"

# 只構建後端
SKIP_FRONTEND=true ./build-local.sh

# 只構建前端
SKIP_BACKEND=true ./build-local.sh

# 清理後重新構建
CLEAN=true ./build-local.sh
```

### 手動構建

#### 構建後端

```bash
cd Application/be

# 下載依賴
make deps

# 構建所有程式
make all

# 或分別構建
make agent
make console
make ui

# 跨平台構建
make build-windows  # 構建 Windows 版本
make build-linux    # 構建 Linux 版本

# 打包發行版
make package
```

#### 構建前端

```bash
cd Application/Fe

# 安裝依賴
npm install

# 開發模式
npm run dev

# 生產構建
npm run build

# 執行生產版本
npm run start

# 型別檢查
npm run type-check

# Linting
npm run lint
```

---

## CI/CD 自動化

### GitHub Actions Workflows

本專案包含完整的 CI/CD 自動化流程：

#### 1. CI Pipeline (`.github/workflows/ci.yml`)

- **觸發條件**: 推送到 `dev` 或 `main` 分支，或 Pull Request
- **執行內容**:
  - Go 程式碼檢查（vet, fmt, test）
  - 前端檢查（type-check, lint, test）
  - Docker 映像構建
  - 安全掃描（Trivy）

#### 2. 安裝檔構建 (`.github/workflows/build-onpremise-installers.yml`)

- **觸發條件**: 
  - 推送到 `dev` 或 `main` 分支
  - 創建版本標籤（`v*`）
  - 手動觸發

- **構建產物**:
  - Windows 安裝程式 (.exe)
  - Linux 套件 (.deb, .rpm)
  - ISO 安裝光碟
  - OVA 虛擬機映像

- **自動發布**:
  - 版本標籤觸發時自動創建 GitHub Release
  - 上傳所有構建產物

### 手動觸發 CI/CD

#### 方式 1: GitHub Web 介面

1. 進入專案的 Actions 頁面
2. 選擇 "Build On-Premise Installers" workflow
3. 點擊 "Run workflow"
4. 選擇分支並輸入版本號（可選）
5. 點擊 "Run workflow" 確認

#### 方式 2: 創建版本標籤

```bash
# 創建版本標籤
git tag -a v1.0.0 -m "Release v1.0.0"

# 推送標籤到遠端
git push origin v1.0.0

# 自動觸發構建和發布
```

### 下載構建產物

#### 從 GitHub Actions

1. 進入 Actions 頁面
2. 選擇對應的 workflow run
3. 滾動到底部的 "Artifacts" 區域
4. 下載需要的產物

#### 從 GitHub Releases

1. 進入 Releases 頁面
2. 選擇對應的版本
3. 下載 Assets 中的安裝檔

---

## 故障排除

### 常見問題

#### 1. 構建失敗

**症狀**: `build-local.ps1` 或 `build-local.sh` 執行失敗

**解決方案**:

```bash
# 檢查 Go 版本
go version  # 需要 1.24+

# 檢查 Node.js 版本
node --version  # 需要 18+

# 清理並重新下載依賴
cd Application/be
go clean -modcache
go mod download

cd ../Fe
rm -rf node_modules package-lock.json
npm install
```

#### 2. 服務無法啟動

**症狀**: 執行 `start.bat` 或 `start.sh` 後服務無法啟動

**解決方案**:

```bash
# 檢查端口是否被占用
netstat -an | findstr ":3001"  # Windows
netstat -an | grep ":3001"     # Linux

# 檢查 PostgreSQL 和 Redis 是否運行
systemctl status postgresql
systemctl status redis

# 查看日誌
tail -f Application/dist/backend/logs/*.log  # Linux
# 或檢查 Windows 事件檢視器
```

#### 3. 無法連接資料庫

**症狀**: 日誌顯示資料庫連接錯誤

**解決方案**:

```bash
# 檢查 PostgreSQL 是否運行
sudo systemctl start postgresql

# 測試連接
psql -h localhost -U postgres -d pandora

# 檢查配置檔案
cat Application/dist/backend/configs/agent-config.yaml
```

#### 4. 前端頁面無法載入

**症狀**: 訪問 http://localhost:3001 無響應

**解決方案**:

```bash
# 檢查 UI Server 是否運行
ps aux | grep axiom-ui  # Linux
tasklist | findstr "axiom-ui"  # Windows

# 檢查防火牆規則
sudo ufw allow 3001  # Linux
netsh advfirewall firewall add rule name="Pandora UI" dir=in action=allow protocol=TCP localport=3001  # Windows
```

### 日誌位置

#### Linux/macOS
- Agent: `Application/dist/backend/logs/agent.log`
- Console: `Application/dist/backend/logs/console.log`
- UI: `Application/dist/backend/logs/ui.log`

#### Windows
- 安裝版: `C:\Program Files\PandoraBox\logs\`
- 本地構建: `Application\dist\backend\logs\`

#### Systemd (Linux)
```bash
sudo journalctl -u pandora-agent -f
```

---

## 常見問題

### Q: 支援哪些作業系統？

**A**: 
- **Windows**: 10, 11, Server 2019, Server 2022
- **Linux**: Ubuntu 20.04+, Debian 11+, RHEL 8+, CentOS 8+

### Q: 可以在虛擬機中運行嗎？

**A**: 可以！我們提供 OVA 虛擬機映像，支援 VirtualBox 和 VMware。

### Q: 如何升級到新版本？

**A**: 
- **Windows**: 執行新版本的安裝程式，會自動升級
- **Linux**: 
  ```bash
  sudo dpkg -i pandora-box-console_NEW_VERSION_amd64.deb  # Debian/Ubuntu
  sudo rpm -U pandora-box-console-NEW_VERSION.rpm         # RHEL/CentOS
  ```

### Q: 如何備份資料？

**A**:
```bash
# 備份 PostgreSQL
pg_dump -U postgres pandora > backup.sql

# 備份配置檔案
tar -czf config-backup.tar.gz /opt/pandora-box/configs
```

### Q: 如何變更監聽端口？

**A**: 編輯配置檔案 `configs/ui-config.yaml`:
```yaml
server:
  port: 3001  # 變更為您想要的端口
```

### Q: 是否需要網際網路連接？

**A**: 不需要。系統可以完全離線運行（使用 ISO 或 OVA 安裝）。

---

## 技術支援

- **文檔**: [完整文檔](README.md)
- **問題回報**: [GitHub Issues](https://github.com/your-org/pandora_box_console_IDS-IPS/issues)
- **討論區**: [GitHub Discussions](https://github.com/your-org/pandora_box_console_IDS-IPS/discussions)
- **電子郵件**: support@pandora-ids.com

---

## 下一步

✅ 已完成安裝？參考以下文檔繼續：

- [使用說明](README.md#使用說明)
- [設定指南](README.md#設定說明)
- [API 文檔](README.md#api-使用)
- [監控指南](README.md#監控與告警)

---

**版本**: v3.0 (dev)  
**維護者**: Pandora Security Team  
**最後更新**: 2025-10-09

