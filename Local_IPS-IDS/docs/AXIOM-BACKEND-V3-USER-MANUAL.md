
# Axiom Backend V3 使用手冊

> **版本**: 3.0.0  
> **日期**: 2025-10-16  
> **目標用戶**: 系統管理員、運維工程師

---

## 📚 目錄

1. [快速開始](#快速開始)
2. [核心功能](#核心功能)
3. [使用場景](#使用場景)
4. [最佳實踐](#最佳實踐)
5. [常見問題](#常見問題)

---

## 快速開始

### 訪問 Web UI

1. 打開瀏覽器訪問：`http://localhost:3001`
2. 登入（默認：admin / admin123）
3. 開始使用

### 主要功能入口

- **服務管理**: `/services-management`
- **量子控制**: `/quantum-control`
- **Windows 日誌**: `/windows-logs`
- **Nginx 配置**: `/nginx-config`

---

## 核心功能

### 1. 服務管理

**位置**: 服務管理 → 總覽

**功能**:
- 查看所有 13 個服務的健康狀態
- 即時健康檢查
- 服務啟動/停止（透過 Portainer）

**使用步驟**:
1. 進入服務管理頁面
2. 查看服務健康狀態
3. 點擊「刷新狀態」進行即時檢查
4. 點擊個別服務查看詳情

### 2. 量子功能

#### 2.1 生成量子密鑰 (QKD)

**位置**: 量子控制 → 生成密鑰

**步驟**:
1. 選擇密鑰長度 (64-2048 bits)
2. 選擇後端 (模擬器 或 IBM Quantum)
3. 點擊「生成密鑰」
4. 等待作業完成
5. 查看生成的密鑰

**建議**:
- 測試時使用模擬器
- 生產環境使用 IBM Quantum
- 密鑰長度建議 256 bits

#### 2.2 Zero Trust 預測

**位置**: 量子控制 → Zero Trust

**步驟**:
1. 輸入用戶 ID
2. 輸入 IP 地址
3. 選擇是否使用量子算法
4. 查看信任分數和決策

**理解結果**:
- trust_score > 0.8: 高信任，允許訪問
- trust_score 0.5-0.8: 中等信任，可能需要 MFA
- trust_score < 0.5: 低信任，拒絕或嚴格限制

### 3. Windows 日誌查看

**位置**: Windows 日誌

**功能**:
- 搜索和過濾日誌
- 按類型/級別/時間範圍過濾
- 查看日誌統計

**使用步驟**:
1. 選擇日誌類型 (System/Security/Application)
2. 選擇日誌級別
3. 輸入關鍵字搜索
4. 點擊「搜索」
5. 查看結果並分頁

**重點關注**:
- **Critical** 和 **Error** 級別的日誌
- 安全相關的事件 ID (如 4625 登入失敗)

### 4. Nginx 配置管理

**位置**: Nginx 配置

**功能**:
- 查看當前配置
- 編輯配置
- 驗證配置
- 重載 Nginx

**使用步驟**:
1. 點擊「編輯配置」
2. 修改配置內容
3. 點擊「保存配置」（自動驗證）
4. 點擊「重載 Nginx」使配置生效

**注意事項**:
⚠️ 配置錯誤可能導致服務中斷
✅ 保存時會自動驗證配置
✅ 保存前會自動備份舊配置

---

## 使用場景

### 場景 1: 監控系統健康

**目標**: 快速了解所有服務的健康狀態

**步驟**:
1. 進入「服務管理」頁面
2. 查看總覽卡片（總服務數、健康服務、健康率）
3. 查看個別服務狀態
4. 對不健康的服務進行診斷

### 場景 2: 調查安全事件

**目標**: 追蹤可疑的登入嘗試

**步驟**:
1. 進入「Windows 日誌」
2. 選擇日誌類型: Security
3. 選擇級別: Error 或 Warning
4. 搜索事件 ID: 4625 (登入失敗)
5. 分析日誌訊息和來源 IP
6. 採取相應措施（封鎖 IP等）

### 場景 3: 執行量子加密

**目標**: 為敏感數據生成量子密鑰

**步驟**:
1. 進入「量子控制」
2. 選擇密鑰長度: 256 bits
3. 選擇後端: ibm_quantum（生產）或 simulator（測試）
4. 點擊「生成密鑰」
5. 等待作業完成
6. 複製生成的密鑰用於加密

### 場景 4: 優化 Nginx 配置

**目標**: 添加新的上游服務

**步驟**:
1. 進入「Nginx 配置」
2. 點擊「編輯配置」
3. 添加新的 upstream 配置
4. 保存（自動驗證）
5. 重載 Nginx
6. 驗證新服務可訪問

---

## 最佳實踐

### 1. 定期檢查服務健康

建議每天至少檢查一次所有服務的健康狀態。

### 2. 監控量子作業成功率

量子作業成功率應保持在 90% 以上。如果低於此值，檢查：
- IBM Quantum 連接
- 網路狀況
- 作業參數是否合理

### 3. 及時處理 Critical 日誌

Windows 日誌中的 Critical 級別應立即處理。

### 4. 配置變更記錄

每次修改 Nginx 配置時：
- 記錄變更原因
- 測試配置有效性
- 備份舊配置

### 5. 定期清理舊日誌

設置日誌保留策略，避免資料庫過大：
- Windows 日誌: 保留 30 天
- API 日誌: 保留 7 天
- 量子作業: 保留 90 天

---

## 常見問題

### Q1: 為什麼某個服務顯示為不健康？

**A**: 可能原因：
1. 服務未啟動
2. 網路連接問題
3. 配置錯誤
4. 資源不足

檢查方法：
```bash
docker ps | grep service-name
docker logs service-name
```

### Q2: 量子作業一直處於 pending 狀態？

**A**: 可能原因：
1. IBM Quantum 隊列繁忙
2. 網路連接問題
3. Token 過期

解決方案：
- 檢查 cyber-ai-quantum 服務日誌
- 驗證 IBM Quantum Token
- 嘗試使用模擬器

### Q3: Windows 日誌沒有新數據？

**A**: 檢查：
1. Windows Log Agent 是否運行
2. Agent 配置是否正確
3. 網路連接是否正常

驗證 Agent：
```powershell
Get-Service | Where-Object {$_.Name -like "*pandora*"}
```

### Q4: Nginx 配置保存失敗？

**A**: 常見原因：
1. 配置語法錯誤
2. 缺少必要的指令
3. 路徑錯誤

建議：
- 使用在線 Nginx 配置驗證工具
- 檢查配置文件語法
- 查看錯誤訊息詳情

### Q5: 如何查看 API 使用情況？

**A**: 查詢 API 日誌表：
```sql
SELECT 
    path, 
    COUNT(*) as request_count,
    AVG(duration) as avg_duration_ms
FROM api_logs
WHERE created_at >= NOW() - INTERVAL '24 hours'
GROUP BY path
ORDER BY request_count DESC
LIMIT 10;
```

---

## 進階功能

### 使用 API 直接調用

所有 Web UI 功能都可以通過 REST API 調用：

```bash
# 查詢 Prometheus
curl -X POST http://localhost:3001/api/v2/prometheus/query \
  -H "Content-Type: application/json" \
  -d '{"query": "up"}'

# 生成量子密鑰
curl -X POST http://localhost:3001/api/v2/quantum/qkd/generate \
  -H "Content-Type: application/json" \
  -d '{"key_length": 256, "backend": "simulator", "shots": 1024}'

# 查詢 Windows 日誌
curl "http://localhost:3001/api/v2/logs/windows?log_type=Security&page=1&page_size=10"
```

### 自動化腳本

使用 API 創建自動化腳本：

```python
import requests

# 定期檢查服務健康
def check_services():
    services = ['prometheus', 'loki', 'quantum']
    for service in services:
        url = f'http://localhost:3001/api/v2/{service}/health'
        response = requests.get(url)
        if not response.ok:
            send_alert(f'{service} is unhealthy!')

# 自動生成量子密鑰
def generate_daily_key():
    url = 'http://localhost:3001/api/v2/quantum/qkd/generate'
    data = {"key_length": 256, "backend": "ibm_quantum", "shots": 2048}
    response = requests.post(url, json=data)
    return response.json()['data']['key']
```

---

## 支援和幫助

### 文檔
- [API 文檔](./AXIOM-BACKEND-V3-API-DOCUMENTATION.md)
- [部署指南](./AXIOM-BACKEND-V3-DEPLOYMENT-GUIDE.md)
- [架構文檔](./AXIOM-BACKEND-V3-COMPLETE-PLAN.md)

### 故障排除
- 查看應用日誌
- 檢查系統資源
- 聯繫技術支援

---

**手冊版本**: 3.0.0  
**最後更新**: 2025-10-16  
**維護者**: Axiom Backend Team



# Axiom Backend V3 使用手冊

> **版本**: 3.0.0  
> **日期**: 2025-10-16  
> **目標用戶**: 系統管理員、運維工程師

---

## 📚 目錄

1. [快速開始](#快速開始)
2. [核心功能](#核心功能)
3. [使用場景](#使用場景)
4. [最佳實踐](#最佳實踐)
5. [常見問題](#常見問題)

---

## 快速開始

### 訪問 Web UI

1. 打開瀏覽器訪問：`http://localhost:3001`
2. 登入（默認：admin / admin123）
3. 開始使用

### 主要功能入口

- **服務管理**: `/services-management`
- **量子控制**: `/quantum-control`
- **Windows 日誌**: `/windows-logs`
- **Nginx 配置**: `/nginx-config`

---

## 核心功能

### 1. 服務管理

**位置**: 服務管理 → 總覽

**功能**:
- 查看所有 13 個服務的健康狀態
- 即時健康檢查
- 服務啟動/停止（透過 Portainer）

**使用步驟**:
1. 進入服務管理頁面
2. 查看服務健康狀態
3. 點擊「刷新狀態」進行即時檢查
4. 點擊個別服務查看詳情

### 2. 量子功能

#### 2.1 生成量子密鑰 (QKD)

**位置**: 量子控制 → 生成密鑰

**步驟**:
1. 選擇密鑰長度 (64-2048 bits)
2. 選擇後端 (模擬器 或 IBM Quantum)
3. 點擊「生成密鑰」
4. 等待作業完成
5. 查看生成的密鑰

**建議**:
- 測試時使用模擬器
- 生產環境使用 IBM Quantum
- 密鑰長度建議 256 bits

#### 2.2 Zero Trust 預測

**位置**: 量子控制 → Zero Trust

**步驟**:
1. 輸入用戶 ID
2. 輸入 IP 地址
3. 選擇是否使用量子算法
4. 查看信任分數和決策

**理解結果**:
- trust_score > 0.8: 高信任，允許訪問
- trust_score 0.5-0.8: 中等信任，可能需要 MFA
- trust_score < 0.5: 低信任，拒絕或嚴格限制

### 3. Windows 日誌查看

**位置**: Windows 日誌

**功能**:
- 搜索和過濾日誌
- 按類型/級別/時間範圍過濾
- 查看日誌統計

**使用步驟**:
1. 選擇日誌類型 (System/Security/Application)
2. 選擇日誌級別
3. 輸入關鍵字搜索
4. 點擊「搜索」
5. 查看結果並分頁

**重點關注**:
- **Critical** 和 **Error** 級別的日誌
- 安全相關的事件 ID (如 4625 登入失敗)

### 4. Nginx 配置管理

**位置**: Nginx 配置

**功能**:
- 查看當前配置
- 編輯配置
- 驗證配置
- 重載 Nginx

**使用步驟**:
1. 點擊「編輯配置」
2. 修改配置內容
3. 點擊「保存配置」（自動驗證）
4. 點擊「重載 Nginx」使配置生效

**注意事項**:
⚠️ 配置錯誤可能導致服務中斷
✅ 保存時會自動驗證配置
✅ 保存前會自動備份舊配置

---

## 使用場景

### 場景 1: 監控系統健康

**目標**: 快速了解所有服務的健康狀態

**步驟**:
1. 進入「服務管理」頁面
2. 查看總覽卡片（總服務數、健康服務、健康率）
3. 查看個別服務狀態
4. 對不健康的服務進行診斷

### 場景 2: 調查安全事件

**目標**: 追蹤可疑的登入嘗試

**步驟**:
1. 進入「Windows 日誌」
2. 選擇日誌類型: Security
3. 選擇級別: Error 或 Warning
4. 搜索事件 ID: 4625 (登入失敗)
5. 分析日誌訊息和來源 IP
6. 採取相應措施（封鎖 IP等）

### 場景 3: 執行量子加密

**目標**: 為敏感數據生成量子密鑰

**步驟**:
1. 進入「量子控制」
2. 選擇密鑰長度: 256 bits
3. 選擇後端: ibm_quantum（生產）或 simulator（測試）
4. 點擊「生成密鑰」
5. 等待作業完成
6. 複製生成的密鑰用於加密

### 場景 4: 優化 Nginx 配置

**目標**: 添加新的上游服務

**步驟**:
1. 進入「Nginx 配置」
2. 點擊「編輯配置」
3. 添加新的 upstream 配置
4. 保存（自動驗證）
5. 重載 Nginx
6. 驗證新服務可訪問

---

## 最佳實踐

### 1. 定期檢查服務健康

建議每天至少檢查一次所有服務的健康狀態。

### 2. 監控量子作業成功率

量子作業成功率應保持在 90% 以上。如果低於此值，檢查：
- IBM Quantum 連接
- 網路狀況
- 作業參數是否合理

### 3. 及時處理 Critical 日誌

Windows 日誌中的 Critical 級別應立即處理。

### 4. 配置變更記錄

每次修改 Nginx 配置時：
- 記錄變更原因
- 測試配置有效性
- 備份舊配置

### 5. 定期清理舊日誌

設置日誌保留策略，避免資料庫過大：
- Windows 日誌: 保留 30 天
- API 日誌: 保留 7 天
- 量子作業: 保留 90 天

---

## 常見問題

### Q1: 為什麼某個服務顯示為不健康？

**A**: 可能原因：
1. 服務未啟動
2. 網路連接問題
3. 配置錯誤
4. 資源不足

檢查方法：
```bash
docker ps | grep service-name
docker logs service-name
```

### Q2: 量子作業一直處於 pending 狀態？

**A**: 可能原因：
1. IBM Quantum 隊列繁忙
2. 網路連接問題
3. Token 過期

解決方案：
- 檢查 cyber-ai-quantum 服務日誌
- 驗證 IBM Quantum Token
- 嘗試使用模擬器

### Q3: Windows 日誌沒有新數據？

**A**: 檢查：
1. Windows Log Agent 是否運行
2. Agent 配置是否正確
3. 網路連接是否正常

驗證 Agent：
```powershell
Get-Service | Where-Object {$_.Name -like "*pandora*"}
```

### Q4: Nginx 配置保存失敗？

**A**: 常見原因：
1. 配置語法錯誤
2. 缺少必要的指令
3. 路徑錯誤

建議：
- 使用在線 Nginx 配置驗證工具
- 檢查配置文件語法
- 查看錯誤訊息詳情

### Q5: 如何查看 API 使用情況？

**A**: 查詢 API 日誌表：
```sql
SELECT 
    path, 
    COUNT(*) as request_count,
    AVG(duration) as avg_duration_ms
FROM api_logs
WHERE created_at >= NOW() - INTERVAL '24 hours'
GROUP BY path
ORDER BY request_count DESC
LIMIT 10;
```

---

## 進階功能

### 使用 API 直接調用

所有 Web UI 功能都可以通過 REST API 調用：

```bash
# 查詢 Prometheus
curl -X POST http://localhost:3001/api/v2/prometheus/query \
  -H "Content-Type: application/json" \
  -d '{"query": "up"}'

# 生成量子密鑰
curl -X POST http://localhost:3001/api/v2/quantum/qkd/generate \
  -H "Content-Type: application/json" \
  -d '{"key_length": 256, "backend": "simulator", "shots": 1024}'

# 查詢 Windows 日誌
curl "http://localhost:3001/api/v2/logs/windows?log_type=Security&page=1&page_size=10"
```

### 自動化腳本

使用 API 創建自動化腳本：

```python
import requests

# 定期檢查服務健康
def check_services():
    services = ['prometheus', 'loki', 'quantum']
    for service in services:
        url = f'http://localhost:3001/api/v2/{service}/health'
        response = requests.get(url)
        if not response.ok:
            send_alert(f'{service} is unhealthy!')

# 自動生成量子密鑰
def generate_daily_key():
    url = 'http://localhost:3001/api/v2/quantum/qkd/generate'
    data = {"key_length": 256, "backend": "ibm_quantum", "shots": 2048}
    response = requests.post(url, json=data)
    return response.json()['data']['key']
```

---

## 支援和幫助

### 文檔
- [API 文檔](./AXIOM-BACKEND-V3-API-DOCUMENTATION.md)
- [部署指南](./AXIOM-BACKEND-V3-DEPLOYMENT-GUIDE.md)
- [架構文檔](./AXIOM-BACKEND-V3-COMPLETE-PLAN.md)

### 故障排除
- 查看應用日誌
- 檢查系統資源
- 聯繫技術支援

---

**手冊版本**: 3.0.0  
**最後更新**: 2025-10-16  
**維護者**: Axiom Backend Team


# Axiom Backend V3 使用手冊

> **版本**: 3.0.0  
> **日期**: 2025-10-16  
> **目標用戶**: 系統管理員、運維工程師

---

## 📚 目錄

1. [快速開始](#快速開始)
2. [核心功能](#核心功能)
3. [使用場景](#使用場景)
4. [最佳實踐](#最佳實踐)
5. [常見問題](#常見問題)

---

## 快速開始

### 訪問 Web UI

1. 打開瀏覽器訪問：`http://localhost:3001`
2. 登入（默認：admin / admin123）
3. 開始使用

### 主要功能入口

- **服務管理**: `/services-management`
- **量子控制**: `/quantum-control`
- **Windows 日誌**: `/windows-logs`
- **Nginx 配置**: `/nginx-config`

---

## 核心功能

### 1. 服務管理

**位置**: 服務管理 → 總覽

**功能**:
- 查看所有 13 個服務的健康狀態
- 即時健康檢查
- 服務啟動/停止（透過 Portainer）

**使用步驟**:
1. 進入服務管理頁面
2. 查看服務健康狀態
3. 點擊「刷新狀態」進行即時檢查
4. 點擊個別服務查看詳情

### 2. 量子功能

#### 2.1 生成量子密鑰 (QKD)

**位置**: 量子控制 → 生成密鑰

**步驟**:
1. 選擇密鑰長度 (64-2048 bits)
2. 選擇後端 (模擬器 或 IBM Quantum)
3. 點擊「生成密鑰」
4. 等待作業完成
5. 查看生成的密鑰

**建議**:
- 測試時使用模擬器
- 生產環境使用 IBM Quantum
- 密鑰長度建議 256 bits

#### 2.2 Zero Trust 預測

**位置**: 量子控制 → Zero Trust

**步驟**:
1. 輸入用戶 ID
2. 輸入 IP 地址
3. 選擇是否使用量子算法
4. 查看信任分數和決策

**理解結果**:
- trust_score > 0.8: 高信任，允許訪問
- trust_score 0.5-0.8: 中等信任，可能需要 MFA
- trust_score < 0.5: 低信任，拒絕或嚴格限制

### 3. Windows 日誌查看

**位置**: Windows 日誌

**功能**:
- 搜索和過濾日誌
- 按類型/級別/時間範圍過濾
- 查看日誌統計

**使用步驟**:
1. 選擇日誌類型 (System/Security/Application)
2. 選擇日誌級別
3. 輸入關鍵字搜索
4. 點擊「搜索」
5. 查看結果並分頁

**重點關注**:
- **Critical** 和 **Error** 級別的日誌
- 安全相關的事件 ID (如 4625 登入失敗)

### 4. Nginx 配置管理

**位置**: Nginx 配置

**功能**:
- 查看當前配置
- 編輯配置
- 驗證配置
- 重載 Nginx

**使用步驟**:
1. 點擊「編輯配置」
2. 修改配置內容
3. 點擊「保存配置」（自動驗證）
4. 點擊「重載 Nginx」使配置生效

**注意事項**:
⚠️ 配置錯誤可能導致服務中斷
✅ 保存時會自動驗證配置
✅ 保存前會自動備份舊配置

---

## 使用場景

### 場景 1: 監控系統健康

**目標**: 快速了解所有服務的健康狀態

**步驟**:
1. 進入「服務管理」頁面
2. 查看總覽卡片（總服務數、健康服務、健康率）
3. 查看個別服務狀態
4. 對不健康的服務進行診斷

### 場景 2: 調查安全事件

**目標**: 追蹤可疑的登入嘗試

**步驟**:
1. 進入「Windows 日誌」
2. 選擇日誌類型: Security
3. 選擇級別: Error 或 Warning
4. 搜索事件 ID: 4625 (登入失敗)
5. 分析日誌訊息和來源 IP
6. 採取相應措施（封鎖 IP等）

### 場景 3: 執行量子加密

**目標**: 為敏感數據生成量子密鑰

**步驟**:
1. 進入「量子控制」
2. 選擇密鑰長度: 256 bits
3. 選擇後端: ibm_quantum（生產）或 simulator（測試）
4. 點擊「生成密鑰」
5. 等待作業完成
6. 複製生成的密鑰用於加密

### 場景 4: 優化 Nginx 配置

**目標**: 添加新的上游服務

**步驟**:
1. 進入「Nginx 配置」
2. 點擊「編輯配置」
3. 添加新的 upstream 配置
4. 保存（自動驗證）
5. 重載 Nginx
6. 驗證新服務可訪問

---

## 最佳實踐

### 1. 定期檢查服務健康

建議每天至少檢查一次所有服務的健康狀態。

### 2. 監控量子作業成功率

量子作業成功率應保持在 90% 以上。如果低於此值，檢查：
- IBM Quantum 連接
- 網路狀況
- 作業參數是否合理

### 3. 及時處理 Critical 日誌

Windows 日誌中的 Critical 級別應立即處理。

### 4. 配置變更記錄

每次修改 Nginx 配置時：
- 記錄變更原因
- 測試配置有效性
- 備份舊配置

### 5. 定期清理舊日誌

設置日誌保留策略，避免資料庫過大：
- Windows 日誌: 保留 30 天
- API 日誌: 保留 7 天
- 量子作業: 保留 90 天

---

## 常見問題

### Q1: 為什麼某個服務顯示為不健康？

**A**: 可能原因：
1. 服務未啟動
2. 網路連接問題
3. 配置錯誤
4. 資源不足

檢查方法：
```bash
docker ps | grep service-name
docker logs service-name
```

### Q2: 量子作業一直處於 pending 狀態？

**A**: 可能原因：
1. IBM Quantum 隊列繁忙
2. 網路連接問題
3. Token 過期

解決方案：
- 檢查 cyber-ai-quantum 服務日誌
- 驗證 IBM Quantum Token
- 嘗試使用模擬器

### Q3: Windows 日誌沒有新數據？

**A**: 檢查：
1. Windows Log Agent 是否運行
2. Agent 配置是否正確
3. 網路連接是否正常

驗證 Agent：
```powershell
Get-Service | Where-Object {$_.Name -like "*pandora*"}
```

### Q4: Nginx 配置保存失敗？

**A**: 常見原因：
1. 配置語法錯誤
2. 缺少必要的指令
3. 路徑錯誤

建議：
- 使用在線 Nginx 配置驗證工具
- 檢查配置文件語法
- 查看錯誤訊息詳情

### Q5: 如何查看 API 使用情況？

**A**: 查詢 API 日誌表：
```sql
SELECT 
    path, 
    COUNT(*) as request_count,
    AVG(duration) as avg_duration_ms
FROM api_logs
WHERE created_at >= NOW() - INTERVAL '24 hours'
GROUP BY path
ORDER BY request_count DESC
LIMIT 10;
```

---

## 進階功能

### 使用 API 直接調用

所有 Web UI 功能都可以通過 REST API 調用：

```bash
# 查詢 Prometheus
curl -X POST http://localhost:3001/api/v2/prometheus/query \
  -H "Content-Type: application/json" \
  -d '{"query": "up"}'

# 生成量子密鑰
curl -X POST http://localhost:3001/api/v2/quantum/qkd/generate \
  -H "Content-Type: application/json" \
  -d '{"key_length": 256, "backend": "simulator", "shots": 1024}'

# 查詢 Windows 日誌
curl "http://localhost:3001/api/v2/logs/windows?log_type=Security&page=1&page_size=10"
```

### 自動化腳本

使用 API 創建自動化腳本：

```python
import requests

# 定期檢查服務健康
def check_services():
    services = ['prometheus', 'loki', 'quantum']
    for service in services:
        url = f'http://localhost:3001/api/v2/{service}/health'
        response = requests.get(url)
        if not response.ok:
            send_alert(f'{service} is unhealthy!')

# 自動生成量子密鑰
def generate_daily_key():
    url = 'http://localhost:3001/api/v2/quantum/qkd/generate'
    data = {"key_length": 256, "backend": "ibm_quantum", "shots": 2048}
    response = requests.post(url, json=data)
    return response.json()['data']['key']
```

---

## 支援和幫助

### 文檔
- [API 文檔](./AXIOM-BACKEND-V3-API-DOCUMENTATION.md)
- [部署指南](./AXIOM-BACKEND-V3-DEPLOYMENT-GUIDE.md)
- [架構文檔](./AXIOM-BACKEND-V3-COMPLETE-PLAN.md)

### 故障排除
- 查看應用日誌
- 檢查系統資源
- 聯繫技術支援

---

**手冊版本**: 3.0.0  
**最後更新**: 2025-10-16  
**維護者**: Axiom Backend Team



# Axiom Backend V3 使用手冊

> **版本**: 3.0.0  
> **日期**: 2025-10-16  
> **目標用戶**: 系統管理員、運維工程師

---

## 📚 目錄

1. [快速開始](#快速開始)
2. [核心功能](#核心功能)
3. [使用場景](#使用場景)
4. [最佳實踐](#最佳實踐)
5. [常見問題](#常見問題)

---

## 快速開始

### 訪問 Web UI

1. 打開瀏覽器訪問：`http://localhost:3001`
2. 登入（默認：admin / admin123）
3. 開始使用

### 主要功能入口

- **服務管理**: `/services-management`
- **量子控制**: `/quantum-control`
- **Windows 日誌**: `/windows-logs`
- **Nginx 配置**: `/nginx-config`

---

## 核心功能

### 1. 服務管理

**位置**: 服務管理 → 總覽

**功能**:
- 查看所有 13 個服務的健康狀態
- 即時健康檢查
- 服務啟動/停止（透過 Portainer）

**使用步驟**:
1. 進入服務管理頁面
2. 查看服務健康狀態
3. 點擊「刷新狀態」進行即時檢查
4. 點擊個別服務查看詳情

### 2. 量子功能

#### 2.1 生成量子密鑰 (QKD)

**位置**: 量子控制 → 生成密鑰

**步驟**:
1. 選擇密鑰長度 (64-2048 bits)
2. 選擇後端 (模擬器 或 IBM Quantum)
3. 點擊「生成密鑰」
4. 等待作業完成
5. 查看生成的密鑰

**建議**:
- 測試時使用模擬器
- 生產環境使用 IBM Quantum
- 密鑰長度建議 256 bits

#### 2.2 Zero Trust 預測

**位置**: 量子控制 → Zero Trust

**步驟**:
1. 輸入用戶 ID
2. 輸入 IP 地址
3. 選擇是否使用量子算法
4. 查看信任分數和決策

**理解結果**:
- trust_score > 0.8: 高信任，允許訪問
- trust_score 0.5-0.8: 中等信任，可能需要 MFA
- trust_score < 0.5: 低信任，拒絕或嚴格限制

### 3. Windows 日誌查看

**位置**: Windows 日誌

**功能**:
- 搜索和過濾日誌
- 按類型/級別/時間範圍過濾
- 查看日誌統計

**使用步驟**:
1. 選擇日誌類型 (System/Security/Application)
2. 選擇日誌級別
3. 輸入關鍵字搜索
4. 點擊「搜索」
5. 查看結果並分頁

**重點關注**:
- **Critical** 和 **Error** 級別的日誌
- 安全相關的事件 ID (如 4625 登入失敗)

### 4. Nginx 配置管理

**位置**: Nginx 配置

**功能**:
- 查看當前配置
- 編輯配置
- 驗證配置
- 重載 Nginx

**使用步驟**:
1. 點擊「編輯配置」
2. 修改配置內容
3. 點擊「保存配置」（自動驗證）
4. 點擊「重載 Nginx」使配置生效

**注意事項**:
⚠️ 配置錯誤可能導致服務中斷
✅ 保存時會自動驗證配置
✅ 保存前會自動備份舊配置

---

## 使用場景

### 場景 1: 監控系統健康

**目標**: 快速了解所有服務的健康狀態

**步驟**:
1. 進入「服務管理」頁面
2. 查看總覽卡片（總服務數、健康服務、健康率）
3. 查看個別服務狀態
4. 對不健康的服務進行診斷

### 場景 2: 調查安全事件

**目標**: 追蹤可疑的登入嘗試

**步驟**:
1. 進入「Windows 日誌」
2. 選擇日誌類型: Security
3. 選擇級別: Error 或 Warning
4. 搜索事件 ID: 4625 (登入失敗)
5. 分析日誌訊息和來源 IP
6. 採取相應措施（封鎖 IP等）

### 場景 3: 執行量子加密

**目標**: 為敏感數據生成量子密鑰

**步驟**:
1. 進入「量子控制」
2. 選擇密鑰長度: 256 bits
3. 選擇後端: ibm_quantum（生產）或 simulator（測試）
4. 點擊「生成密鑰」
5. 等待作業完成
6. 複製生成的密鑰用於加密

### 場景 4: 優化 Nginx 配置

**目標**: 添加新的上游服務

**步驟**:
1. 進入「Nginx 配置」
2. 點擊「編輯配置」
3. 添加新的 upstream 配置
4. 保存（自動驗證）
5. 重載 Nginx
6. 驗證新服務可訪問

---

## 最佳實踐

### 1. 定期檢查服務健康

建議每天至少檢查一次所有服務的健康狀態。

### 2. 監控量子作業成功率

量子作業成功率應保持在 90% 以上。如果低於此值，檢查：
- IBM Quantum 連接
- 網路狀況
- 作業參數是否合理

### 3. 及時處理 Critical 日誌

Windows 日誌中的 Critical 級別應立即處理。

### 4. 配置變更記錄

每次修改 Nginx 配置時：
- 記錄變更原因
- 測試配置有效性
- 備份舊配置

### 5. 定期清理舊日誌

設置日誌保留策略，避免資料庫過大：
- Windows 日誌: 保留 30 天
- API 日誌: 保留 7 天
- 量子作業: 保留 90 天

---

## 常見問題

### Q1: 為什麼某個服務顯示為不健康？

**A**: 可能原因：
1. 服務未啟動
2. 網路連接問題
3. 配置錯誤
4. 資源不足

檢查方法：
```bash
docker ps | grep service-name
docker logs service-name
```

### Q2: 量子作業一直處於 pending 狀態？

**A**: 可能原因：
1. IBM Quantum 隊列繁忙
2. 網路連接問題
3. Token 過期

解決方案：
- 檢查 cyber-ai-quantum 服務日誌
- 驗證 IBM Quantum Token
- 嘗試使用模擬器

### Q3: Windows 日誌沒有新數據？

**A**: 檢查：
1. Windows Log Agent 是否運行
2. Agent 配置是否正確
3. 網路連接是否正常

驗證 Agent：
```powershell
Get-Service | Where-Object {$_.Name -like "*pandora*"}
```

### Q4: Nginx 配置保存失敗？

**A**: 常見原因：
1. 配置語法錯誤
2. 缺少必要的指令
3. 路徑錯誤

建議：
- 使用在線 Nginx 配置驗證工具
- 檢查配置文件語法
- 查看錯誤訊息詳情

### Q5: 如何查看 API 使用情況？

**A**: 查詢 API 日誌表：
```sql
SELECT 
    path, 
    COUNT(*) as request_count,
    AVG(duration) as avg_duration_ms
FROM api_logs
WHERE created_at >= NOW() - INTERVAL '24 hours'
GROUP BY path
ORDER BY request_count DESC
LIMIT 10;
```

---

## 進階功能

### 使用 API 直接調用

所有 Web UI 功能都可以通過 REST API 調用：

```bash
# 查詢 Prometheus
curl -X POST http://localhost:3001/api/v2/prometheus/query \
  -H "Content-Type: application/json" \
  -d '{"query": "up"}'

# 生成量子密鑰
curl -X POST http://localhost:3001/api/v2/quantum/qkd/generate \
  -H "Content-Type: application/json" \
  -d '{"key_length": 256, "backend": "simulator", "shots": 1024}'

# 查詢 Windows 日誌
curl "http://localhost:3001/api/v2/logs/windows?log_type=Security&page=1&page_size=10"
```

### 自動化腳本

使用 API 創建自動化腳本：

```python
import requests

# 定期檢查服務健康
def check_services():
    services = ['prometheus', 'loki', 'quantum']
    for service in services:
        url = f'http://localhost:3001/api/v2/{service}/health'
        response = requests.get(url)
        if not response.ok:
            send_alert(f'{service} is unhealthy!')

# 自動生成量子密鑰
def generate_daily_key():
    url = 'http://localhost:3001/api/v2/quantum/qkd/generate'
    data = {"key_length": 256, "backend": "ibm_quantum", "shots": 2048}
    response = requests.post(url, json=data)
    return response.json()['data']['key']
```

---

## 支援和幫助

### 文檔
- [API 文檔](./AXIOM-BACKEND-V3-API-DOCUMENTATION.md)
- [部署指南](./AXIOM-BACKEND-V3-DEPLOYMENT-GUIDE.md)
- [架構文檔](./AXIOM-BACKEND-V3-COMPLETE-PLAN.md)

### 故障排除
- 查看應用日誌
- 檢查系統資源
- 聯繫技術支援

---

**手冊版本**: 3.0.0  
**最後更新**: 2025-10-16  
**維護者**: Axiom Backend Team

