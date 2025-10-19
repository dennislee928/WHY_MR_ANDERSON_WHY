# Pandora Box Console - 應用程式（地端部署版本）

這個目錄包含了 Pandora Box Console IDS-IPS 的完整應用程式，專為地端（On-Premise）部署而設計。

## 📁 目錄結構

```
Application/
├── be/                    # 後端應用程式
│   ├── Makefile          # 後端構建腳本（17個目標）
│   ├── build.ps1, .sh    # 構建腳本
│   ├── go.mod            # Go 模組引用
│   └── README.md         # 後端說明
│
├── Fe/                    # 前端應用程式
│   ├── components/       # React 組件（7個）
│   ├── pages/            # Next.js 頁面（3個）
│   ├── services/         # API 服務層
│   ├── hooks/            # 自定義 Hooks（2個）
│   ├── types/            # TypeScript 類型
│   ├── lib/              # 工具函數
│   ├── styles/           # 樣式文件
│   ├── public/           # 靜態資源
│   ├── package.json      # NPM 套件定義
│   ├── next.config.js    # Next.js 配置
│   ├── tsconfig.json     # TypeScript 配置
│   └── README.md         # 前端說明
│
├── docker/                # 🆕 Docker 映像檔案
│   ├── agent.dockerfile           # Pandora Agent
│   ├── agent.koyeb.dockerfile     # Agent (Koyeb版)
│   ├── server-be.dockerfile       # 後端 API
│   ├── ui.patr.dockerfile         # UI Server
│   ├── server-fe.dockerfile       # 前端伺服器
│   ├── monitoring.dockerfile      # 監控堆疊
│   ├── nginx.dockerfile           # Nginx
│   ├── test.dockerfile            # 測試環境
│   └── README.md                  # Docker 說明
│
├── docker-compose.yml     # 🆕 Docker Compose 編排（11個服務）
├── docker-start.ps1       # 🆕 Docker 啟動腳本（Windows）
├── docker-start.sh        # 🆕 Docker 啟動腳本（Linux）
├── .env.example           # 🆕 環境變數範例
│
├── build-local.ps1        # Windows 本地構建腳本
├── build-local.sh         # Linux/macOS 本地構建腳本
├── dist/                  # 構建產物（不納入版控）
└── README.md              # 本檔案
```

## 🚀 快速開始

### 方式 1: 使用 Docker Compose（最簡單）⭐

**Windows**:
```powershell
cd Application
.\docker-start.ps1
```

**Linux/macOS**:
```bash
cd Application
chmod +x docker-start.sh
./docker-start.sh
```

這會啟動所有 11 個服務：
- Pandora Agent、UI Server
- Prometheus、Grafana、Loki、AlertManager
- PostgreSQL、Redis
- Nginx、Promtail、Node Exporter

訪問: http://localhost:3001

### 方式 2: 使用自動構建腳本

#### Windows

```powershell
# 構建所有（後端 + 前端）
.\build-local.ps1

# 只構建後端
.\build-local.ps1 -SkipFrontend

# 只構建前端
.\build-local.ps1 -SkipBackend

# 清理並重新構建
.\build-local.ps1 -Clean

# 指定版本
.\build-local.ps1 -Version "1.0.0"
```

#### Linux/macOS

```bash
# 構建所有（後端 + 前端）
./build-local.sh

# 只構建後端
SKIP_FRONTEND=true ./build-local.sh

# 只構建前端
SKIP_BACKEND=true ./build-local.sh

# 清理並重新構建
CLEAN=true ./build-local.sh

# 指定版本
./build-local.sh all "1.0.0"
```

### 方式 2: 手動構建

#### 構建後端

```bash
cd be
make all              # 構建所有程式
make agent            # 只構建 Agent
make console          # 只構建 Console
make ui               # 只構建 UI Server
```

#### 構建前端

```bash
cd Fe
npm install           # 安裝依賴
npm run build         # 構建生產版本
npm run dev           # 開發模式
```

## 📦 安裝檔生成

本專案支援生成多種格式的安裝檔，透過 GitHub Actions CI/CD 自動化構建：

### 支援的安裝檔格式

| 格式 | 平台 | 說明 |
|------|------|------|
| `.exe` | Windows | Inno Setup 安裝程式 |
| `.msi` | Windows | Windows Installer 套件 |
| `.deb` | Linux | Debian/Ubuntu 套件 |
| `.rpm` | Linux | RedHat/CentOS 套件 |
| `.iso` | All | 可開機安裝光碟 |
| `.ova` | All | VirtualBox/VMware 虛擬機映像 |

### 觸發安裝檔構建

#### 方式 1: 推送程式碼到 dev 或 main 分支

```bash
git add .
git commit -m "feat: 新增功能"
git push origin dev
```

#### 方式 2: 創建版本標籤

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

#### 方式 3: 手動觸發

1. 前往 GitHub Actions 頁面
2. 選擇 "Build On-Premise Installers" workflow
3. 點擊 "Run workflow"
4. 選擇分支並輸入版本號
5. 點擊 "Run workflow"

### 下載安裝檔

構建完成後，可以在以下位置下載：

1. **GitHub Actions Artifacts**: 
   - 進入 Actions 頁面
   - 選擇對應的 workflow run
   - 下載 artifacts

2. **GitHub Releases**（僅限標籤構建）:
   - 進入 Releases 頁面
   - 下載對應版本的安裝檔

## 🔧 開發

### 前置需求

- **Go** 1.24+
- **Node.js** 18+
- **Git**
- **Make** (Linux/macOS) 或 **PowerShell** (Windows)

### 開發工作流程

1. **克隆專案**
   ```bash
   git clone https://github.com/your-org/pandora_box_console_IDS-IPS.git
   cd pandora_box_console_IDS-IPS/Application
   ```

2. **啟動後端開發環境**
   ```bash
   cd be
   make run-agent     # 啟動 Agent（開發模式）
   make run-console   # 啟動 Console
   make run-ui        # 啟動 UI Server
   ```

3. **啟動前端開發環境**
   ```bash
   cd Fe
   npm run dev        # 啟動 Next.js 開發伺服器
   ```

4. **執行測試**
   ```bash
   # 後端測試
   cd be
   make test
   
   # 前端測試
   cd Fe
   npm run test
   ```

## 📋 構建產物

構建完成後，`dist/` 目錄將包含：

```
dist/
├── backend/              # 後端程式
│   ├── pandora-agent.exe (或 pandora-agent)
│   ├── pandora-console.exe (或 pandora-console)
│   ├── axiom-ui.exe (或 axiom-ui)
│   └── configs/          # 配置檔案
│
├── frontend/             # 前端程式
│   ├── .next/           # Next.js 構建產物
│   └── public/          # 靜態資源
│
├── start.bat/.sh         # 啟動腳本
├── stop.bat/.sh          # 停止腳本
└── README.txt            # 使用說明
```

## 🚢 部署

### 使用構建產物部署

1. 將 `dist/` 目錄複製到目標伺服器

2. 編輯配置檔案
   ```bash
   cd dist/backend/configs
   # 編輯 agent-config.yaml、console-config.yaml 等
   ```

3. 啟動服務
   ```bash
   # Windows
   start.bat
   
   # Linux/macOS
   ./start.sh
   ```

### 使用安裝檔部署

#### Windows (.exe)

1. 執行安裝程式
   ```
   pandora-box-console-1.0.0-windows-amd64-setup.exe
   ```

2. 按照安裝精靈完成安裝

3. 從開始選單啟動應用程式

#### Linux (.deb)

```bash
# Ubuntu/Debian
sudo dpkg -i pandora-box-console_1.0.0_amd64.deb
sudo apt-get install -f  # 安裝依賴

# 啟動服務
sudo systemctl start pandora-agent
sudo systemctl enable pandora-agent
```

#### Linux (.rpm)

```bash
# RedHat/CentOS
sudo rpm -i pandora-box-console-1.0.0-1.x86_64.rpm

# 啟動服務
sudo systemctl start pandora-agent
sudo systemctl enable pandora-agent
```

#### ISO 安裝光碟

1. 掛載 ISO 映像
   ```bash
   sudo mount -o loop pandora-box-console-1.0.0-amd64.iso /mnt
   ```

2. 執行安裝腳本
   ```bash
   cd /mnt
   sudo ./install.sh
   ```

3. 重新開機並啟動服務

#### OVA 虛擬機

1. 匯入 OVA 到 VirtualBox/VMware
   - VirtualBox: 檔案 > 匯入應用裝置
   - VMware: 檔案 > 開啟

2. 啟動虛擬機

3. 預設登入資訊:
   - 使用者名稱: `pandora`
   - 密碼: `pandora`

4. 服務會自動啟動

## 🔐 系統需求

### 最低需求

- **CPU**: 2 核心
- **記憶體**: 4 GB RAM
- **儲存空間**: 20 GB
- **作業系統**:
  - Windows 10/11 或 Windows Server 2019/2022
  - Ubuntu 20.04/22.04 或 Debian 11/12
  - RHEL/CentOS 8/9

### 建議配置

- **CPU**: 4 核心以上
- **記憶體**: 8 GB RAM 以上
- **儲存空間**: 50 GB 以上（包含日誌和資料）
- **網路**: 千兆網卡

### 軟體依賴

- **PostgreSQL** 14+
- **Redis** 7+
- **Nginx** (可選，用於反向代理)

## 📚 文檔

- [後端開發指南](be/README.md)
- [前端開發指南](Fe/README.md)
- [專案結構說明](../README-PROJECT-STRUCTURE.md)
- [主要 README](../README.md)

## 🐛 故障排除

### 構建失敗

**問題**: Go 構建失敗

```bash
# 清理 Go 快取
go clean -modcache
go mod download
go mod tidy
```

**問題**: npm 安裝失敗

```bash
# 清理 npm 快取
cd Fe
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

### 執行時錯誤

**問題**: 無法連接資料庫

- 檢查 PostgreSQL 是否已啟動
- 檢查配置檔案中的資料庫連線字串
- 確認防火牆規則允許連線

**問題**: 端口已被占用

- 修改配置檔案中的端口設定
- 或停止占用端口的程式

## 🤝 貢獻

歡迎提交 Pull Request 或建立 Issue！

## 📄 授權

MIT License - 詳見 [LICENSE](../LICENSE)

---

**維護者**: Pandora Security Team  
**技術支援**: support@pandora-ids.com  
**最後更新**: 2025-10-09

