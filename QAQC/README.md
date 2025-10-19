# Security Platform API Tests - Robot Framework

## 📋 **專案概述**

本專案使用 Robot Framework 對 Security Platform 的 Cloudflare Workers API 端點進行全面的自動化測試。

## 🚀 **功能特色**

### **測試覆蓋範圍**
- ✅ **健康檢查測試** - API 可用性和基本功能
- ✅ **安全性測試** - 威脅檢測和防護功能
- ✅ **網路測試** - 網路統計和監控
- ✅ **設備管理測試** - 設備列表和管理
- ✅ **效能測試** - 回應時間和吞吐量
- ✅ **錯誤處理測試** - 異常情況處理
- ✅ **CORS 測試** - 跨域請求支援
- ✅ **速率限制測試** - API 限制機制

### **Cloudflare Workers 特定測試**
- ✅ **D1 Database** - 資料庫連線和操作
- ✅ **KV Namespaces** - 快取和會話管理
- ✅ **R2 Buckets** - 檔案儲存功能
- ✅ **Workers AI** - AI 威脅檢測
- ✅ **Durable Objects** - WebSocket 管理
- ✅ **Analytics Engine** - 事件分析
- ✅ **Rate Limiter** - 速率限制
- ✅ **Secrets Store** - 安全憑證管理
- ✅ **Queue** - 非同步處理
- ✅ **Vectorize** - 向量搜尋

## 📁 **檔案結構**

```
QAQC/
├── api_tests.robot                    # 主要 API 測試
├── cloudflare_workers_tests.robot     # Cloudflare Workers 特定測試
├── test_suite_config.robot            # 測試套件配置
├── run_tests.sh                       # Linux/Mac 執行腳本
├── run_tests.ps1                      # Windows PowerShell 執行腳本
├── requirements.txt                   # Python 依賴套件
└── README.md                          # 本檔案
```

## 🛠️ **安裝與設定**

### **1. 安裝 Python 依賴**
```bash
# 安裝 Robot Framework 和相關套件
pip install -r requirements.txt

# 或手動安裝
pip install robotframework robotframework-requests
```

### **2. 驗證安裝**
```bash
# 檢查 Robot Framework 版本
robot --version

# 檢查 Python 版本
python --version
```

## 🎯 **執行測試**

### **Linux/Mac 使用 Bash**
```bash
# 執行所有測試
./run_tests.sh

# 執行特定測試類型
./run_tests.sh smoke          # 煙霧測試
./run_tests.sh regression     # 回歸測試
./run_tests.sh cloudflare     # Cloudflare Workers 測試
./run_tests.sh performance    # 效能測試
./run_tests.sh integration    # 整合測試
```

### **Windows 使用 PowerShell**
```powershell
# 執行所有測試
.\run_tests.ps1

# 執行特定測試類型
.\run_tests.ps1 smoke          # 煙霧測試
.\run_tests.ps1 regression     # 回歸測試
.\run_tests.ps1 cloudflare     # Cloudflare Workers 測試
.\run_tests.ps1 performance    # 效能測試
.\run_tests.ps1 integration    # 整合測試
```

### **直接使用 Robot Framework**
```bash
# 執行特定測試檔案
robot api_tests.robot
robot cloudflare_workers_tests.robot

# 執行特定標籤的測試
robot --include smoke api_tests.robot
robot --include performance api_tests.robot

# 執行測試並生成報告
robot --outputdir results --logdir logs *.robot
```

## 📊 **測試類型說明**

### **1. Smoke Tests (煙霧測試)**
- 基本功能驗證
- API 可用性檢查
- 快速驗證部署狀態

### **2. Regression Tests (回歸測試)**
- 完整功能測試
- 資料驗證
- 錯誤處理測試

### **3. Cloudflare Workers Tests**
- Workers 特定功能
- Bindings 狀態檢查
- 服務整合測試

### **4. Performance Tests (效能測試)**
- 回應時間測試
- 吞吐量測試
- 速率限制測試

### **5. Integration Tests (整合測試)**
- 端到端測試
- 服務間整合
- 資料一致性測試

## 🔧 **配置選項**

### **環境變數**
```bash
# API 端點配置
export BASE_URL="https://security-platform-worker.workers.dev"
export API_VERSION="v1"

# 測試配置
export TIMEOUT=30
export RETRY_COUNT=3
```

### **測試標籤**
- `smoke` - 煙霧測試
- `regression` - 回歸測試
- `performance` - 效能測試
- `integration` - 整合測試
- `health` - 健康檢查
- `security` - 安全性測試
- `network` - 網路測試
- `devices` - 設備測試

## 📈 **測試報告**

### **HTML 報告**
測試執行後會生成詳細的 HTML 報告：
- `report.html` - 測試摘要報告
- `log.html` - 詳細測試日誌
- `output.xml` - XML 格式結果

### **報告內容**
- 測試執行統計
- 通過/失敗率
- 執行時間分析
- 錯誤詳情
- 截圖和日誌

## 🚨 **故障排除**

### **常見問題**

#### **1. API 無法存取**
```bash
# 檢查 API 端點
curl https://security-platform-worker.workers.dev/api/v1/health

# 檢查網路連線
ping security-platform-worker.workers.dev
```

#### **2. Robot Framework 未安裝**
```bash
# 重新安裝
pip install --upgrade robotframework robotframework-requests
```

#### **3. 權限問題 (Windows)**
```powershell
# 以管理員身份執行 PowerShell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

#### **4. Python 版本問題**
```bash
# 檢查 Python 版本 (需要 3.7+)
python --version

# 使用 Python 3 明確執行
python3 -m pip install robotframework
```

### **除錯模式**
```bash
# 啟用詳細日誌
robot --loglevel DEBUG *.robot

# 只執行特定測試
robot --test "Test Health Check Endpoint" api_tests.robot
```

## 📝 **自訂測試**

### **新增測試案例**
```robot
*** Test Cases ***
My Custom Test
    [Documentation]    My custom test case
    [Tags]    custom
    ${response}=    GET    ${BASE_URL}/api/${API_VERSION}/my-endpoint
    Should Be Equal As Strings    ${response.status_code}    200
```

### **新增關鍵字**
```robot
*** Keywords ***
My Custom Keyword
    [Documentation]    My custom keyword
    [Arguments]    ${param1}    ${param2}
    Log    Executing custom keyword with ${param1} and ${param2}
```

## 🔄 **CI/CD 整合**

### **GitHub Actions**
```yaml
name: API Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Set up Python
        uses: actions/setup-python@v2
        with:
          python-version: '3.9'
      - name: Install dependencies
        run: pip install -r QAQC/requirements.txt
      - name: Run tests
        run: cd QAQC && ./run_tests.sh smoke
```

### **Jenkins Pipeline**
```groovy
pipeline {
    agent any
    stages {
        stage('Test') {
            steps {
                sh 'cd QAQC && ./run_tests.sh all'
            }
        }
    }
    post {
        always {
            publishHTML([
                allowMissing: false,
                alwaysLinkToLastBuild: true,
                keepAll: true,
                reportDir: 'QAQC/results',
                reportFiles: 'report.html',
                reportName: 'API Test Report'
            ])
        }
    }
}
```

## 📚 **相關資源**

- [Robot Framework 官方文件](https://robotframework.org/)
- [Robot Framework Requests Library](https://github.com/MarketSquare/robotframework-requests)
- [Cloudflare Workers 文件](https://developers.cloudflare.com/workers/)
- [API 測試最佳實踐](https://martinfowler.com/articles/practical-test-pyramid.html)

## 🤝 **貢獻指南**

1. Fork 專案
2. 建立功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交變更 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 開啟 Pull Request

## 📄 **授權**

本專案採用 MIT 授權 - 詳見 [LICENSE](LICENSE) 檔案

---

**注意**: 請確保在執行測試前，Cloudflare Workers 已正確部署且 API 端點可正常存取。
