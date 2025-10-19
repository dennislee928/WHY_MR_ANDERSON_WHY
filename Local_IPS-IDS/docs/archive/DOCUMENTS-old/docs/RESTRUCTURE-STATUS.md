# 專案重整狀態報告

**日期**: 2025-10-08  
**狀態**: 🟡 準備階段完成，等待執行

---

## ✅ 已完成的準備工作

### 1. 目錄結構創建
已創建以下新目錄結構：

```
✓ bin/                              # 編譯產物
✓ build/docker/                     # Dockerfile 集中管理
✓ build/package/                    # 打包腳本
✓ docs/architecture/                # 架構文檔
✓ docs/deployment/                  # 部署指南
✓ docs/development/                 # 開發指南
✓ docs/operations/                  # 運維文檔
✓ deployments/kubernetes/base/      # K8s 基礎配置
✓ deployments/kubernetes/gcp/       # GCP K8s
✓ deployments/kubernetes/oci/       # OCI K8s
✓ deployments/terraform/            # Terraform
✓ deployments/paas/                 # PaaS 配置
✓ deployments/docker-compose/       # Docker Compose
```

### 2. 腳本工具
創建了自動化重整腳本：
- ✅ `scripts/restructure-project.ps1` - 自動化遷移腳本
- ✅ `docs/PROJECT-RESTRUCTURE-PLAN.md` - 詳細重整計劃

### 3. 文檔
- ✅ 重整計劃文檔
- ✅ 本狀態報告

---

## 📋 待執行的操作

### 執行選項

#### 選項 A: 完全自動化（推薦用於測試）

```powershell
# DRY RUN - 查看將會執行什麼操作（不實際移動文件）
.\scripts\restructure-project.ps1 -DryRun

# 實際執行（請先備份！）
.\scripts\restructure-project.ps1
```

#### 選項 B: 手動執行（更安全，推薦）

1. **備份專案**
   ```powershell
   git add -A
   git commit -m "backup: 重整前備份"
   git branch restructure-backup
   ```

2. **創建新分支**
   ```powershell
   git checkout -b feature/project-restructure
   ```

3. **分階段執行**
   - 階段 1: 移動編譯產物 (低風險)
   - 階段 2: 移動 Dockerfile (需更新 CI)
   - 階段 3: 移動文檔 (低風險)
   - 階段 4-7: 移動配置文件 (需更新多個配置)

---

## 🔄 需要更新的配置文件

### 高優先級（必須更新）

#### 1. CI/CD Workflows

**`.github/workflows/ci.yml`**
```yaml
# 需要更新的路徑:
- file: ./Dockerfile.${{ matrix.image }}
# 改為:
+ file: ./build/docker/${{ matrix.image }}.dockerfile
```

**`.github/workflows/deploy-gcp.yml`**
```yaml
# 需要更新的路徑:
- file: ./Dockerfile.agent
+ file: ./build/docker/agent.dockerfile

- file: ./Dockerfile
+ file: ./build/docker/console.dockerfile

- kubectl apply -k k8s-gcp/
+ kubectl apply -k deployments/kubernetes/gcp/
```

**`.github/workflows/deploy-oci.yml`**
```yaml
# 需要更新的路徑:
- file: ./Dockerfile.agent
+ file: ./build/docker/agent.dockerfile

- kubectl apply -k k8s/
+ kubectl apply -k deployments/kubernetes/base/
```

**`.github/workflows/deploy-paas.yml`**
```yaml
# 需要更新的路徑:
- file: ./Dockerfile.agent.koyeb
+ file: ./build/docker/agent.koyeb.dockerfile

- flyctl deploy --config fly.toml --dockerfile Dockerfile.monitoring
+ flyctl deploy --config deployments/paas/flyio/fly.toml --dockerfile build/docker/monitoring.dockerfile
```

#### 2. Makefile

```makefile
# 需要更新:
-  go build -o pandora-agent ./cmd/agent
+  go build -o bin/pandora-agent ./cmd/agent

-  docker build -f Dockerfile.agent -t pandora-agent .
+  docker build -f build/docker/agent.dockerfile -t pandora-agent .
```

#### 3. .gitignore

```gitignore
# 添加:
+ bin/
+ build/temp/
+ *.exe
+ *.dll
```

### 中優先級

#### 4. PaaS 配置文件

**`deployments/paas/koyeb/koyeb.yaml`**
```yaml
docker:
-  dockerfile: Dockerfile.agent.koyeb
+  dockerfile: build/docker/agent.koyeb.dockerfile
```

**`deployments/paas/flyio/fly.toml`**
```toml
[build]
-  dockerfile = "Dockerfile.monitoring"
+  dockerfile = "build/docker/monitoring.dockerfile"
```

#### 5. K8s Kustomization

**`deployments/kubernetes/base/kustomization.yaml`**
```yaml
# 確保所有資源路徑正確
resources:
- namespace.yaml
- configmap.yaml
- secrets.yaml
# ... 等等
```

---

## ⚠️ 風險評估

### 高風險操作
- ❌ 移動 K8s 配置（影響生產部署）
- ❌ 移動 Terraform 配置（影響基礎設施）

### 中風險操作
- ⚠️ 移動 Dockerfile（影響 CI/CD）
- ⚠️ 移動 PaaS 配置（影響部署）

### 低風險操作
- ✅ 移動文檔（不影響運行）
- ✅ 移動編譯產物（可重新編譯）

---

## 🚀 推薦執行步驟

### 第一階段：低風險操作
1. 移動編譯產物到 `bin/`
2. 移動文檔到 `docs/`
3. 更新 README.md 中的文檔鏈接
4. 提交第一個 PR

### 第二階段：中風險操作
1. 移動 Dockerfile 到 `build/docker/`
2. 更新 .gitignore
3. 更新 Makefile
4. 測試本地建置
5. 提交第二個 PR

### 第三階段：高風險操作
1. 更新 CI/CD workflows
2. 測試 CI/CD 流程
3. 移動 Docker Compose
4. 提交第三個 PR

### 第四階段：K8s 和 Terraform
1. 複製（不移動）K8s 配置到新位置
2. 複製（不移動）Terraform 到新位置
3. 更新 workflows 指向新位置
4. 測試部署
5. 驗證後再刪除舊位置

---

## 📊 進度追蹤

- [x] 創建目錄結構
- [x] 創建重整腳本
- [x] 創建文檔
- [ ] 執行文件遷移
- [ ] 更新 .gitignore
- [ ] 更新 Makefile
- [ ] 更新 CI/CD workflows
- [ ] 測試本地建置
- [ ] 測試 CI/CD
- [ ] 清理舊文件

---

## 💡 建議

### 立即執行（低風險）
```powershell
# 1. 先測試看看會發生什麼
.\scripts\restructure-project.ps1 -DryRun

# 2. 查看操作日誌
cat docs/restructure-operations.csv

# 3. 如果滿意，執行實際操作
.\scripts\restructure-project.ps1
```

### 分步執行（更安全）
1. 手動移動編譯產物和文檔
2. 提交並測試
3. 再移動 Dockerfile 和配置
4. 逐步更新 CI/CD

---

## 📞 需要協助？

如果遇到問題：
1. 檢查 `docs/restructure-operations.csv` 日誌
2. 使用 `git status` 查看變更
3. 隨時可以 `git checkout .` 撤銷變更

---

**下一步建議**: 
1. 執行 DRY RUN 查看效果
2. 備份專案
3. 在新分支中執行重整
4. 測試建置流程
5. 更新 CI/CD 配置

