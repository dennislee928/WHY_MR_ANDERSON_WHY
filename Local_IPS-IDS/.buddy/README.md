# Buddy Works 管道檔案
## 單獨導入每個管道

> 📅 **創建日期**: 2025-10-09  
> 🎯 **目標**: 解決 "Definition must contain only one pipeline" 錯誤

---

## 🚨 重要說明

Buddy Works 一次只能導入一個管道。請按照以下順序分別導入每個檔案：

---

## 📁 管道檔案列表

### 1. 01-build-installers.yml
- **管道名稱**: Build On-Premise Installers
- **觸發方式**: Push (main/dev)
- **優先級**: HIGH
- **用途**: 構建所有平台的安裝檔

### 2. 02-ci-pipeline.yml
- **管道名稱**: CI Pipeline
- **觸發方式**: Push (main/dev)
- **優先級**: NORMAL
- **用途**: 持續集成檢查

### 3. 03-kubernetes-deployment.yml
- **管道名稱**: Kubernetes Deployment
- **觸發方式**: Manual (Click)
- **優先級**: HIGH
- **用途**: 部署到 Kubernetes 集群

### 4. 04-performance-testing.yml
- **管道名稱**: Performance Testing
- **觸發方式**: Manual (Click)
- **優先級**: NORMAL
- **用途**: 性能測試和基準測試

### 5. 05-security-audit.yml
- **管道名稱**: Security Audit
- **觸發方式**: Manual (Click)
- **優先級**: HIGH
- **用途**: 安全掃描和審計

### 6. 06-chaos-engineering.yml
- **管道名稱**: Chaos Engineering
- **觸發方式**: Manual (Click)
- **優先級**: NORMAL
- **用途**: 混沌工程測試

---

## 🚀 導入步驟

### 在 Buddy Works 中：

1. 點擊 "Pipelines" → "Add new"
2. 選擇 "Import YAML" → "From Git"
3. 設置：
   - **PROJECT**: `Local_IPS-IDS (This project)`
   - **BRANCH**: `main` 或 `dev`
   - **YAML PATH**: `.buddy/pipelines/01-build-installers.yml`
4. 點擊 "Import pipeline"
5. 重複步驟 1-4，導入其他管道

---

## 🔧 配置需求

### 環境變數

在 Buddy 項目設置中添加：

| 變數名稱 | 類型 | 描述 |
|----------|------|------|
| `GITHUB_TOKEN` | Secret | GitHub Personal Access Token |
| `BUDDY_REPO_SLUG` | 自動 | 倉庫 slug |
| `BUDDY_EXECUTION_BRANCH` | 自動 | 當前分支 |

### GitHub Token 權限

- ✅ `repo` - 完整倉庫訪問
- ✅ `write:packages` - 上傳 artifacts
- ✅ `read:org` - 讀取組織資訊

---

## 📊 導入順序建議

1. **Build On-Premise Installers** (最重要)
2. **CI Pipeline** (持續集成)
3. **Kubernetes Deployment** (部署)
4. **Performance Testing** (測試)
5. **Security Audit** (安全)
6. **Chaos Engineering** (彈性)

---

## 🎯 導入完成後

您將擁有完整的 CI/CD 流程：

- ✅ 自動構建和發布
- ✅ 持續集成檢查
- ✅ Kubernetes 部署
- ✅ 性能測試
- ✅ 安全審計
- ✅ 混沌工程

---

## 📚 相關文檔

- [導入指南](../../docs/BUDDY-PIPELINE-IMPORT-GUIDE.md)
- [設置指南](../../docs/BUDDY-WORKS-SETUP.md)
- [YAML 修復](../../docs/BUDDY-YAML-FIX.md)

---

**狀態**: ✅ 管道檔案已創建  
**下一步**: 開始導入管道  
**預計時間**: 10-15 分鐘

**🎉 準備好導入 Buddy Works 管道了！**
