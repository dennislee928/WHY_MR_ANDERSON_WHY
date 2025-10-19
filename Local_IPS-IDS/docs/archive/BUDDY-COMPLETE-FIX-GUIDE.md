# Buddy Works 完整修復指南
## 解決所有導入問題

> 📅 **修復日期**: 2025-10-09  
> 🎯 **目標**: 解決所有 Buddy Works 導入問題  
> ✅ **狀態**: 提供完整解決方案

---

## 🐛 問題總結

### 問題 1: 多管道檔案
```
ERROR Parsing YAML failed. Definition must contain only one pipeline.
```
**檔案**: `.buddy/pipeline.fixed.yml`  
**原因**: 包含 17 個管道，但 Buddy 一次只能導入一個

### 問題 2: 路徑錯誤
```
Yaml path: Wrong location
```
**檔案**: 所有單獨的管道檔案  
**原因**: `cached_dirs` 路徑格式錯誤

### 問題 3: Cached Directories 路徑
```
At least one value of 'cached_dirs' array doesn't start with '/'
```
**原因**: Buddy Works 要求所有 `cached_dirs` 路徑必須以 `/` 開頭

---

## ✅ 完整解決方案

### 方案 A: 使用根目錄檔案（推薦）

**最簡單的解決方案**，使用已經修復的根目錄檔案：

1. **在 Buddy Works 中**:
   - 點擊 "Pipelines" → "Add new"
   - 選擇 "Import YAML" → "From Git"
   - 設置：
     - **PROJECT**: `Local_IPS-IDS (This project)`
     - **BRANCH**: `main` 或 `dev`
     - **YAML PATH**: `build-installers.yml`
   - 點擊 "Import pipeline"

2. **配置環境變數**:
   - 添加 `GITHUB_TOKEN` (Secret)

### 方案 B: 使用修復後的單獨檔案

**已修復的檔案**:
- ✅ `.buddy/01-build-installers.yml` - Build On-Premise Installers
- ✅ `.buddy/02-ci-pipeline.yml` - CI Pipeline
- ✅ `.buddy/03-kubernetes-deployment.yml` - Kubernetes Deployment
- ✅ `.buddy/04-performance-testing.yml` - Performance Testing
- ✅ `.buddy/05-security-audit.yml` - Security Audit
- ✅ `.buddy/06-chaos-engineering.yml` - Chaos Engineering

**導入步驟**:

1. **導入 Build On-Premise Installers**:
   ```
   YAML PATH: .buddy/01-build-installers.yml
   ```

2. **導入 CI Pipeline**:
   ```
   YAML PATH: .buddy/02-ci-pipeline.yml
   ```

3. **導入 Kubernetes Deployment**:
   ```
   YAML PATH: .buddy/03-kubernetes-deployment.yml
   ```

4. **導入 Performance Testing**:
   ```
   YAML PATH: .buddy/04-performance-testing.yml
   ```

5. **導入 Security Audit**:
   ```
   YAML PATH: .buddy/05-security-audit.yml
   ```

6. **導入 Chaos Engineering**:
   ```
   YAML PATH: .buddy/06-chaos-engineering.yml
   ```

### 方案 C: 使用 Inline YAML

如果路徑問題持續，使用 "Inline YAML" 選項：

1. 在 Buddy 中選擇 "Inline YAML"
2. 複製任何一個修復後的檔案內容
3. 直接貼上到 Buddy 編輯器

---

## 🔧 已修復的問題

### 1. Cached Directories 路徑修復

**修復前**:
```yaml
cached_dirs:
- "node_modules"      # ❌ 錯誤格式
- ".next/cache"       # ❌ 錯誤格式
```

**修復後**:
```yaml
cached_dirs:
- "/node_modules"     # ✅ 正確格式
- "/.next/cache"      # ✅ 正確格式
```

### 2. GitHub Release Action 修復

**修復前**:
```yaml
- action: "Create GitHub Release"
  type: "GITHUB_RELEASE"  # ❌ 不支持的類型
```

**修復後**:
```yaml
- action: "Create GitHub Release"
  type: "BUILD"           # ✅ 使用 BUILD + GitHub API
  docker_image_name: "library/alpine"
  execute_commands:
  - "curl -X POST -H 'Authorization: token $GITHUB_TOKEN' ..."
```

### 3. Kubernetes Apply Action 修復

**修復前**:
```yaml
- action: "Deploy to Kubernetes"
  type: "KUBERNETES_APPLY"  # ❌ 不支持的類型
```

**修復後**:
```yaml
- action: "Deploy to Kubernetes"
  type: "BUILD"             # ✅ 使用 BUILD + kubectl
  docker_image_name: "bitnami/kubectl"
  execute_commands:
  - "kubectl apply -f deployments/kubernetes/"
```

---

## 📊 檔案狀態對照表

| 檔案 | 狀態 | 問題 | 修復 |
|------|------|------|------|
| `build-installers.yml` | ✅ 可用 | 無 | 已修復所有問題 |
| `.buddy/01-build-installers.yml` | ✅ 可用 | cached_dirs | 已修復 |
| `.buddy/02-ci-pipeline.yml` | ✅ 可用 | cached_dirs | 已修復 |
| `.buddy/03-kubernetes-deployment.yml` | ✅ 可用 | 無 | 無需修復 |
| `.buddy/04-performance-testing.yml` | ✅ 可用 | 無 | 無需修復 |
| `.buddy/05-security-audit.yml` | ✅ 可用 | 無 | 無需修復 |
| `.buddy/06-chaos-engineering.yml` | ✅ 可用 | 無 | 無需修復 |
| `.buddy/pipeline.fixed.yml` | ❌ 不可用 | 多管道 + cached_dirs | 已修復但包含多管道 |

---

## 🚀 推薦導入順序

### 第一階段：核心管道

1. **Build On-Premise Installers** (最重要)
   - 使用: `build-installers.yml`
   - 觸發: Push (main/dev)
   - 用途: 構建所有平台安裝檔

2. **CI Pipeline** (持續集成)
   - 使用: `.buddy/02-ci-pipeline.yml`
   - 觸發: Push (main/dev)
   - 用途: 代碼檢查和測試

### 第二階段：部署管道

3. **Kubernetes Deployment**
   - 使用: `.buddy/03-kubernetes-deployment.yml`
   - 觸發: Manual
   - 用途: K8s 集群部署

### 第三階段：測試管道

4. **Performance Testing**
   - 使用: `.buddy/04-performance-testing.yml`
   - 觸發: Manual
   - 用途: 性能驗證

5. **Security Audit**
   - 使用: `.buddy/05-security-audit.yml`
   - 觸發: Manual
   - 用途: 安全檢查

6. **Chaos Engineering**
   - 使用: `.buddy/06-chaos-engineering.yml`
   - 觸發: Manual
   - 用途: 彈性測試

---

## 🔧 配置需求

### 環境變數

在 Buddy Works 項目設置中添加：

| 變數名稱 | 類型 | 描述 | 範例 |
|----------|------|------|------|
| `GITHUB_TOKEN` | Secret | GitHub Personal Access Token | `ghp_xxxx...` |
| `BUDDY_REPO_SLUG` | 自動 | 倉庫 slug | `your-org/Local_IPS-IDS` |
| `BUDDY_EXECUTION_BRANCH` | 自動 | 當前分支 | `main` 或 `dev` |

### GitHub Token 權限

需要以下權限：
- ✅ `repo` - 完整倉庫訪問
- ✅ `write:packages` - 上傳 artifacts
- ✅ `read:org` - 讀取組織資訊

---

## 🎯 故障排除

### 檢查清單

- [ ] 檔案存在於指定分支
- [ ] 路徑格式正確
- [ ] YAML 語法有效
- [ ] Buddy 有倉庫訪問權限
- [ ] 分支名稱正確
- [ ] 環境變數已設置

### 常見錯誤解決

1. **檔案不存在**:
   - 檢查分支
   - 確認檔案已推送

2. **權限問題**:
   - 檢查 GitHub 倉庫權限
   - 確認 Buddy 整合設置

3. **YAML 語法錯誤**:
   - 驗證 YAML 格式
   - 檢查縮排

4. **環境變數未設置**:
   - 在項目設置中添加變數
   - 設置為 Secret（如果是敏感資訊）

---

## 🎊 成功指標

導入成功後，您應該看到：

- ✅ 管道出現在 Buddy 項目中
- ✅ 管道配置正確顯示
- ✅ 觸發條件設置正確
- ✅ Actions 列表完整
- ✅ 環境變數正確設置

---

## 📚 相關文檔

- [YAML 修復說明](BUDDY-YAML-FIX.md) - 詳細修復記錄
- [路徑修復指南](BUDDY-PATH-FIX.md) - 路徑問題解決
- [管道導入指南](BUDDY-PIPELINE-IMPORT-GUIDE.md) - 完整導入流程
- [Buddy Works 設置](BUDDY-WORKS-SETUP.md) - 環境配置

---

**狀態**: ✅ 所有問題已修復  
**推薦方案**: 使用 `build-installers.yml`  
**下一步**: 開始導入管道  
**預計時間**: 5-10 分鐘

**🎉 現在應該可以成功導入 Buddy Works 管道了！**
