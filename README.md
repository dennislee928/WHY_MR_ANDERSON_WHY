# Serverless E2E Tester (Robot Framework)

一個開箱即用的 GitHub Action，用於在 CI 環境中執行 Robot Framework API 測試，無需手動安裝 Python、Java 或任何依賴套件。

[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Serverless%20E2E%20Tester-blue.svg)](https://github.com/marketplace/actions/serverless-e2e-tester-robot-framework-by-dennis)

## 🎯 核心價值

**解決 Serverless 測試環境建置的痛苦** - Robot Framework 很好用，但在 CI 環境安裝 Java/Python/Dependencies 很煩。這個 Action 讓您專注於撰寫測試，而不是配置環境。

## ✨ 功能特色

- ✅ **開箱即用** - 內建 Robot Framework + Requests Library + JSON Library
- ✅ **環境注入** - 自動將 GitHub Secrets 注入為測試變數
- ✅ **報告輸出** - 自動生成 HTML 測試報告
- ✅ **標籤過濾** - 支援 `--include` 和 `--exclude` 標籤
- ✅ **變數檔案** - 支援自訂變數檔案
- ✅ **並行執行** - 支援多進程並行測試
- ✅ **零配置** - 不需要安裝任何依賴

## 📋 使用範例

### 基本使用

```yaml
name: API E2E Tests

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Robot Framework Tests
        uses: dennislee928/WHY_MR_ANDERSON_WHY@v0.0.2
        with:
          test_dir: 'QAQC'
          target_url: 'https://your-api.workers.dev'
      
      - name: Upload Test Reports
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: robot-framework-reports
          path: reports/
          retention-days: 7
```

### 使用 GitHub Secrets

```yaml
- name: Run Robot Framework Tests
  uses: dennislee928/WHY_MR_ANDERSON_WHY@v0.0.2
  env:
    API_KEY: ${{ secrets.API_KEY }}
    AUTH_TOKEN: ${{ secrets.AUTH_TOKEN }}
  with:
    test_dir: 'QAQC'
    target_url: 'https://your-api.workers.dev'
```

在 Robot Framework 測試中，這些 secrets 會自動注入為變數：

```robot
*** Test Cases ***
Test API with Authentication
    ${headers}=    Create Dictionary
    ...    Authorization=Bearer ${AUTH_TOKEN}
    ...    X-API-Key=${API_KEY}
    ${response}=    GET    ${BASE_URL}/api/v1/endpoint    headers=${headers}
    Should Be Equal As Strings    ${response.status_code}    200
```

### 使用標籤過濾

```yaml
- name: Run Smoke Tests Only
  uses: dennislee928/WHY_MR_ANDERSON_WHY@v0.0.2
  with:
    test_dir: 'QAQC'
    target_url: 'https://your-api.workers.dev'
    include_tags: 'smoke'
    exclude_tags: 'slow,integration'
```

### 並行執行

```yaml
- name: Run Tests in Parallel
  uses: dennislee928/WHY_MR_ANDERSON_WHY@v0.0.2
  with:
    test_dir: 'QAQC'
    target_url: 'https://your-api.workers.dev'
    processes: '4'
```

## 📖 Inputs 說明

| Input | 說明 | 必填 | 預設值 |
|-------|------|------|--------|
| `test_dir` | 包含 `.robot` 測試檔案的目錄 | ✅ | - |
| `target_url` | 部署的 API 基礎 URL（例如 Cloudflare Worker URL） | ✅ | - |
| `report_dir` | 測試報告儲存目錄 | ❌ | `reports` |
| `include_tags` | 要包含的測試標籤（逗號分隔，例如：`smoke,regression`） | ❌ | - |
| `exclude_tags` | 要排除的測試標籤（逗號分隔） | ❌ | - |
| `variable_file` | Robot Framework 變數檔案路徑（多個檔案用逗號分隔） | ❌ | - |
| `processes` | 並行執行進程數（預設 1 為順序執行） | ❌ | `1` |
| `timeout` | 測試超時時間（Robot Framework 格式，例如：`5 minutes`） | ❌ | - |
| `log_level` | 日誌級別（TRACE, DEBUG, INFO, WARN, ERROR, NONE） | ❌ | `INFO` |
| `artifact_name` | 上傳的 artifact 名稱 | ❌ | `robot-framework-reports` |
| `artifact_retention_days` | Artifact 保留天數 | ❌ | `7` |

## 📤 Outputs 說明

| Output | 說明 |
|--------|------|
| `exit_code` | Robot Framework 執行退出碼（0 = 成功，非零 = 失敗） |
| `report_path` | 生成的測試報告目錄路徑 |

## 🎯 適用場景

- **Cloudflare Workers** 部署驗證
- **AWS Lambda / Vercel** Serverless API 測試
- **CI/CD Pipeline** 中的 Smoke Tests

---
---

# Unified Security & Infrastructure Platform (Project Context)

> 以下為本專案完整平台的說明文件，包含此 Action 的來源專案背景。

[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Python Version](https://img.shields.io/badge/Python-3.11+-blue.svg)](https://python.org)
[![Docker](https://img.shields.io/badge/Docker-20.10+-blue.svg)](https://docker.com)

[繁體中文](README.zh-TW.md) | English

## Overview

A comprehensive, cloud-native security and infrastructure management platform combining:
- **IDS/IPS System** - Real-time intrusion detection and prevention
- **AI/ML Threat Detection** - Deep learning-based security analysis
- **Quantum Computing Integration** - IBM Quantum for advanced cryptography
- **Security Scanning Tools** - Integrated Nuclei, Nmap, AMASS scanners
- **Multi-Cloud Deployment** - Support for Cloudflare Workers, OCI, IBM Cloud

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│           Unified Security & Infrastructure Platform         │
└─────────────────────────────────────────────────────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
    ┌────▼─────┐      ┌──────▼──────┐      ┌─────▼──────┐
    │ Frontend │      │   Backend   │      │ AI/Quantum │
    │  (React) │      │    (Go)     │      │  (Python)  │
    └────┬─────┘      └──────┬──────┘      └─────┬──────┘
         │                   │                    │
         └───────────────────┼────────────────────┘
                             │
                    ┌────────▼────────┐
                    │  Infrastructure  │
                    │   - Docker       │
                    │   - Kubernetes   │
                    │   - Multi-Cloud  │
                    └──────────────────┘
```

## Quick Start

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.24+ (for local development)
- Python 3.11+ (for AI/Quantum features)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/dennislee928/WHY_MR_ANDERSON_WHY.git
   cd WHY_MR_ANDERSON_WHY
   ```

2. **Setup Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start Services**
   ```bash
   docker-compose up -d
   ```

## Documentation

- [Architecture Details](docs/architecture/system-design.md)
- [API Reference](docs/development/api-reference.md)
- [Security Guide](docs/security/)
- [Deployment Guide](docs/deployment/)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

- **Author**: Dennislee928
- **Project**: [GitHub Repository](https://github.com/dennislee928/WHY_MR_ANDERSON_WHY)
