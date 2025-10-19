# Buddy Works 管道導入指南
## 單獨導入每個管道

> 📅 **創建日期**: 2025-10-09  
> 🎯 **目標**: 解決 "Definition must contain only one pipeline" 錯誤  
> ✅ **狀態**: 已創建單獨的管道檔案

---

## 🐛 問題解決

### 原始錯誤

```
ERROR Parsing YAML failed. Definition must contain only one pipeline.
```

### 原因

Buddy Works 一次只能導入一個管道，但我們的 `buddy.yml` 包含了 17 個管道。

### 解決方案

將每個管道拆分成單獨的 YAML 檔案，分別導入。

---

## 📁 管道檔案結構

### 創建的檔案

```
.buddy/pipelines/
├── 01-build-installers.yml     # Build On-Premise Installers
├── 02-ci-pipeline.yml          # CI Pipeline  
├── 03-kubernetes-deployment.yml # Kubernetes Deployment
├── 04-performance-testing.yml  # Performance Testing
├── 05-security-audit.yml       # Security Audit
├── 06-chaos-engineering.yml    # Chaos Engineering
└── README.md                   # 導入說明
```

### 每個檔案包含

- ✅ 單一管道定義
- ✅ 完整的 action 配置
- ✅ 正確的 Docker 鏡像
- ✅ 修復後的 action types

---

## 🚀 導入步驟

### Step 1: 導入 Build On-Premise Installers

1. 在 Buddy Works 中點擊 "Pipelines" → "Add new"
2. 選擇 "Import YAML" → "From Git"
3. 設置：
   - **PROJECT**: `Local_IPS-IDS (This project)`
   - **BRANCH**: `main` 或 `dev`
   - **YAML PATH**: `.buddy/pipelines/01-build-installers.yml`
4. 點擊 "Import pipeline"

**注意**: 如果路徑錯誤，嘗試以下格式：
- `.buddy/pipelines/01-build-installers.yml` (推薦)
- `buddy/pipelines/01-build-installers.yml`
- `.buddy/pipelines/01-build-installers`

### Step 2: 導入 CI Pipeline

1. 重複步驟 1
2. 設置：
   - **YAML PATH**: `.buddy/pipelines/02-ci-pipeline.yml`
3. 點擊 "Import pipeline"

### Step 3: 導入 Kubernetes Deployment

1. 重複步驟 1
2. 設置：
   - **YAML PATH**: `.buddy/pipelines/03-kubernetes-deployment.yml`
3. 點擊 "Import pipeline"

### Step 4: 導入 Performance Testing

1. 重複步驟 1
2. 設置：
   - **YAML PATH**: `.buddy/pipelines/04-performance-testing.yml`
3. 點擊 "Import pipeline"

### Step 5: 導入 Security Audit

1. 重複步驟 1
2. 設置：
   - **YAML PATH**: `.buddy/pipelines/05-security-audit.yml`
3. 點擊 "Import pipeline"

### Step 6: 導入 Chaos Engineering

1. 重複步驟 1
2. 設置：
   - **YAML PATH**: `.buddy/pipelines/06-chaos-engineering.yml`
3. 點擊 "Import pipeline"

---

## 📊 管道對照表

| 檔案 | 管道名稱 | 觸發方式 | 優先級 | 狀態 |
|------|----------|----------|--------|------|
| `01-build-installers.yml` | Build On-Premise Installers | Push (main/dev) | HIGH | ✅ 準備導入 |
| `02-ci-pipeline.yml` | CI Pipeline | Push (main/dev) | NORMAL | ✅ 準備導入 |
| `03-kubernetes-deployment.yml` | Kubernetes Deployment | Manual | HIGH | ✅ 準備導入 |
| `04-performance-testing.yml` | Performance Testing | Manual | NORMAL | ✅ 準備導入 |
| `05-security-audit.yml` | Security Audit | Manual | HIGH | ✅ 準備導入 |
| `06-chaos-engineering.yml` | Chaos Engineering | Manual | NORMAL | ✅ 準備導入 |

---

## 🔧 配置需求

### 環境變數（每個管道都需要）

在 Buddy Works 項目設置中添加：

| 變數名稱 | 類型 | 描述 | 範例 |
|----------|------|------|------|
| `GITHUB_TOKEN` | Secret | GitHub Personal Access Token | `ghp_xxxx...` |
| `BUDDY_REPO_SLUG` | 自動 | 倉庫 slug（owner/repo） | `your-org/Local_IPS-IDS` |
| `BUDDY_EXECUTION_BRANCH` | 自動 | 當前執行分支 | `main` 或 `dev` |

### GitHub Token 權限

需要以下權限：
- ✅ `repo` - 完整倉庫訪問
- ✅ `write:packages` - 上傳 artifacts
- ✅ `read:org` - 讀取組織資訊

---

## 🎯 導入順序建議

### 第一階段：核心管道

1. **Build On-Premise Installers** - 最重要的管道
2. **CI Pipeline** - 持續集成

### 第二階段：部署管道

3. **Kubernetes Deployment** - 容器化部署

### 第三階段：測試管道

4. **Performance Testing** - 性能驗證
5. **Security Audit** - 安全檢查
6. **Chaos Engineering** - 彈性測試

---

## 📝 導入檢查清單

### 導入前準備

- [ ] 配置 `GITHUB_TOKEN` 環境變數
- [ ] 確認 GitHub Token 有適當權限
- [ ] 檢查倉庫分支（main/dev）

### 導入過程

- [ ] 導入 01-build-installers.yml
- [ ] 導入 02-ci-pipeline.yml
- [ ] 導入 03-kubernetes-deployment.yml
- [ ] 導入 04-performance-testing.yml
- [ ] 導入 05-security-audit.yml
- [ ] 導入 06-chaos-engineering.yml

### 導入後驗證

- [ ] 檢查管道配置
- [ ] 驗證觸發條件
- [ ] 測試手動觸發
- [ ] 檢查環境變數

---

## 🚨 常見問題

### Q1: 導入失敗怎麼辦？

**A**: 檢查 YAML 語法，確保：
- 沒有多餘的空格
- 縮排正確
- 沒有不支持的 action types

### Q2: 環境變數未設置怎麼辦？

**A**: 在項目設置中添加：
- Settings → Variables → Add variable
- 設置為 Secret（如果是敏感資訊）

### Q3: GitHub Token 權限不足？

**A**: 重新生成 Token 並確保有：
- `repo` 權限
- `write:packages` 權限
- `read:org` 權限

### Q4: 管道不觸發？

**A**: 檢查：
- 觸發條件（Push 分支）
- 分支名稱是否正確
- 管道是否啟用

---

## 🎊 導入完成後

### 管道狀態

導入完成後，您將擁有：

- ✅ 6 個獨立的管道
- ✅ 完整的 CI/CD 流程
- ✅ 自動和手動觸發選項
- ✅ 多平台構建支持
- ✅ 安全掃描和測試

### 下一步

1. 配置 Kubernetes 集群連接
2. 設置 Slack 通知
3. 運行第一個管道
4. 監控執行結果

---

## 📚 相關文檔

- [Buddy Works 設置指南](BUDDY-WORKS-SETUP.md)
- [YAML 修復說明](BUDDY-YAML-FIX.md)
- [Phase 4 路線圖](../PHASE4-ROADMAP.md)

---

**狀態**: ✅ 管道檔案已創建  
**下一步**: 開始導入管道  
**預計時間**: 10-15 分鐘

**🎉 準備好導入 Buddy Works 管道了！**
