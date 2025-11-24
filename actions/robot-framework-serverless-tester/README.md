# Serverless E2E Tester (Robot Framework)

一個開箱即用的 GitHub Action，用於在 CI 環境中執行 Robot Framework API 測試，無需手動安裝 Python、Java 或任何依賴套件。

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
        uses: ./actions/robot-framework-serverless-tester
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
  uses: ./actions/robot-framework-serverless-tester
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
  uses: ./actions/robot-framework-serverless-tester
  with:
    test_dir: 'QAQC'
    target_url: 'https://your-api.workers.dev'
    include_tags: 'smoke'
    exclude_tags: 'slow,integration'
```

### 使用變數檔案

```yaml
- name: Run Tests with Custom Variables
  uses: ./actions/robot-framework-serverless-tester
  with:
    test_dir: 'QAQC'
    target_url: 'https://your-api.workers.dev'
    variable_file: 'QAQC/variables.py,QAQC/config.py'
```

### 並行執行

```yaml
- name: Run Tests in Parallel
  uses: ./actions/robot-framework-serverless-tester
  with:
    test_dir: 'QAQC'
    target_url: 'https://your-api.workers.dev'
    processes: '4'
```

### 完整範例（包含所有功能）

```yaml
name: Comprehensive E2E Tests

on:
  push:
    branches: [ main ]
  workflow_dispatch:

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Robot Framework E2E Tests
        id: robot-tests
        uses: ./actions/robot-framework-serverless-tester
        env:
          API_KEY: ${{ secrets.API_KEY }}
          DATABASE_URL: ${{ secrets.DATABASE_URL }}
        with:
          test_dir: 'QAQC'
          target_url: 'https://your-api.workers.dev'
          report_dir: 'test-reports'
          include_tags: 'smoke,regression'
          exclude_tags: 'slow'
          variable_file: 'QAQC/test_variables.py'
          processes: '2'
          timeout: '5 minutes'
          log_level: 'DEBUG'
      
      - name: Upload Test Reports
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: robot-framework-reports
          path: test-reports/
          retention-days: 7
      
      - name: Check Test Results
        if: steps.robot-tests.outputs.exit_code != '0'
        run: |
          echo "Tests failed with exit code: ${{ steps.robot-tests.outputs.exit_code }}"
          exit 1
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

## 🔐 GitHub Secrets 使用方式

這個 Action 會自動將所有環境變數注入為 Robot Framework 變數。在 workflow 中使用 `env` 區塊傳遞 secrets：

```yaml
- name: Run Tests
  uses: ./actions/robot-framework-serverless-tester
  env:
    SECRET_KEY: ${{ secrets.SECRET_KEY }}
    API_TOKEN: ${{ secrets.API_TOKEN }}
  with:
    test_dir: 'QAQC'
    target_url: 'https://your-api.workers.dev'
```

在 Robot Framework 測試中，這些變數可以直接使用：

```robot
*** Test Cases ***
Test with Secret
    ${response}=    GET    ${BASE_URL}/api/endpoint
    ...    headers={"Authorization": "Bearer ${API_TOKEN}"}
    Should Be Equal As Strings    ${response.status_code}    200
```

**安全提示**：
- Secrets 不會出現在日誌中
- 只處理大寫環境變數（避免洩漏系統變數）
- 自動排除 GitHub Actions 內部變數

## 📊 查看測試報告

測試執行後，報告會上傳為 GitHub Artifacts：

1. 前往 GitHub Actions 頁面
2. 選擇對應的 workflow run
3. 在 Artifacts 區塊下載 `robot-framework-reports`
4. 解壓縮後開啟 `report.html` 查看詳細報告

報告包含：
- `report.html` - 測試摘要報告
- `log.html` - 詳細測試日誌
- `output.xml` - XML 格式結果（用於 CI 整合）

## 🎯 適用場景

### Cloudflare Workers 部署驗證

```yaml
- name: Deploy to Cloudflare Workers
  run: wrangler deploy

- name: Verify Deployment with E2E Tests
  uses: ./actions/robot-framework-serverless-tester
  with:
    test_dir: 'QAQC'
    target_url: 'https://your-worker.workers.dev'
    include_tags: 'smoke'
```

### 多環境測試

```yaml
strategy:
  matrix:
    environment: [staging, production]
    include:
      - environment: staging
        url: https://staging-api.example.com
      - environment: production
        url: https://api.example.com

- name: Test ${{ matrix.environment }}
  uses: ./actions/robot-framework-serverless-tester
  with:
    test_dir: 'QAQC'
    target_url: ${{ matrix.url }}
```

## 🚨 故障排除

### 測試目錄不存在

```
ERROR: Test directory does not exist: QAQC
```

**解決方案**：確認 `test_dir` 路徑正確，相對於 repository 根目錄。

### 變數檔案找不到

```
WARNING: Variable file not found: variables.py (skipping)
```

**解決方案**：確認變數檔案路徑正確，或移除 `variable_file` input。

### 測試超時

如果測試執行時間過長，可以：

1. 增加超時時間：
```yaml
timeout: '10 minutes'
```

2. 使用標籤過濾，只執行快速測試：
```yaml
include_tags: 'smoke'
exclude_tags: 'slow'
```

### 報告未生成

如果報告未生成，檢查：

1. 測試是否成功執行（即使失敗也會生成報告）
2. `report_dir` 路徑是否正確
3. 檢查 workflow 日誌中的錯誤訊息

## 🔧 進階配置

### 自訂變數檔案範例

建立 `variables.py`：

```python
def get_variables():
    return {
        'API_VERSION': 'v1',
        'TIMEOUT': 30,
        'RETRY_COUNT': 3
    }
```

在 workflow 中使用：

```yaml
variable_file: 'QAQC/variables.py'
```

### 標籤使用建議

在 Robot Framework 測試中使用標籤：

```robot
*** Test Cases ***
Quick Smoke Test
    [Tags]    smoke    quick
    ${response}=    GET    ${BASE_URL}/health
    Should Be Equal As Strings    ${response.status_code}    200

Slow Integration Test
    [Tags]    integration    slow
    # ... 長時間執行的測試
```

然後在 workflow 中過濾：

```yaml
include_tags: 'smoke'      # 只執行快速測試
exclude_tags: 'slow'       # 排除慢速測試
```

## 📚 相關資源

- [Robot Framework 官方文件](https://robotframework.org/)
- [Robot Framework Requests Library](https://github.com/MarketSquare/robotframework-requests)
- [GitHub Actions 文件](https://docs.github.com/en/actions)

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

## 📄 授權

MIT License

---

**注意**：此 Action 專為 Serverless API 測試設計，特別適合 Cloudflare Workers、Vercel Functions、AWS Lambda 等部署後的驗證測試。

