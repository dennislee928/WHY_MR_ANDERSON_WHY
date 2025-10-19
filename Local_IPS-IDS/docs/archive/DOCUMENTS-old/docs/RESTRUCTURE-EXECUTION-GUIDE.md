# 專案重整執行指南

**狀態**: 🟢 配置已更新，準備執行文件遷移  
**日期**: 2025-10-08

---

## ✅ 已完成的配置更新

### 1. CI/CD Workflows ✓
- ✅ `.github/workflows/ci.yml`
  - 更新 Dockerfile 路徑為 `build/docker/`
  - 修正重複的 `needs` 依賴
  - 修正映像名稱為 `pandora_*`

- ✅ `.github/workflows/deploy-gcp.yml`
  - 更新 Dockerfile 路徑
  - 修正 tags 格式（使用 YAML 列表）
  - 更新 K8s manifests 路徑為 `deployments/kubernetes/gcp/`

- ✅ `.github/workflows/deploy-oci.yml`
  - 修正語法錯誤（Line 5）
  - 更新 Dockerfile 路徑
  - 更新 K8s manifests 路徑為 `deployments/kubernetes/base/`

- ✅ `.github/workflows/deploy-paas.yml`
  - 更新 Koyeb Dockerfile 路徑
  - 更新 Patr UI Dockerfile 路徑
  - 更新 Fly.io 部署命令路徑

### 2. Build 配置 ✓
- ✅ `Makefile`
  - Docker Compose 路徑已更新
  - 輸出目錄已設為 `bin/`（已經正確）

- ✅ `.gitignore`
  - 添加 `bin/` 排除
  - 添加編譯產物排除規則
  - 添加臨時文件排除

### 3. 文檔和腳本 ✓
- ✅ 重整計劃文檔已創建
- ✅ CI/CD 更新指南已創建
- ✅ 自動化腳本已創建
- ✅ 專案結構說明已創建

---

## 🚀 立即執行（三種方案）

### 方案 A: 自動化腳本（推薦用於測試）

```powershell
# 1. 備份當前狀態
git add -A
git commit -m "chore: 配置已更新，準備重整"
git checkout -b feature/project-restructure

# 2. 執行 DRY RUN（不實際移動文件）
.\scripts\restructure-project.ps1 -DryRun

# 3. 查看操作日誌
Get-Content docs\restructure-operations.csv | Format-Table

# 4. 如果滿意，執行實際操作
.\scripts\restructure-project.ps1

# 5. 檢查變更
git status

# 6. 提交變更
git add -A
git commit -m "feat: 重整專案結構"
```

### 方案 B: 手動分階段執行（推薦用於生產）

#### 階段 1: 移動文檔（低風險）✅
```powershell
# 已有目錄結構，執行文檔遷移
Move-Item DEPLOYMENT.md docs\deployment\README.md
Move-Item DEPLOYMENT-GCP.md docs\deployment\gcp.md
Move-Item KOYEB-*.md docs\deployment\koyeb\
Move-Item FLYIO-*.md docs\deployment\flyio\
# ... 等等

git add docs/
git commit -m "docs: 重整文檔目錄結構"
```

#### 階段 2: 移動 Dockerfiles（中風險）✅
```powershell
# 已完成 - Dockerfiles 在 build/docker/

# 測試建置
docker build -f build/docker/agent.dockerfile -t test-agent .

git add build/
git commit -m "build: 移動 Dockerfile 到 build/docker/"
```

#### 階段 3: 移動 K8s 和 Terraform（高風險）⚠️
```powershell
# 複製（不刪除原文件）K8s 配置
Copy-Item -Recurse k8s/* deployments/kubernetes/base/
Copy-Item -Recurse k8s-gcp/* deployments/kubernetes/gcp/

# 複製 Terraform
Copy-Item -Recurse terraform/* deployments/terraform/

# 測試 kubectl
kubectl apply -k deployments/kubernetes/base/ --dry-run=client

git add deployments/
git commit -m "feat: 新增 deployments 目錄結構"
```

#### 階段 4: 移動 PaaS 配置
```powershell
# Fly.io
Move-Item fly.toml deployments/paas/flyio/
Move-Item fly-monitoring.toml deployments/paas/flyio/

# Koyeb
Move-Item koyeb.yaml deployments/paas/koyeb/
Copy-Item -Recurse .koyeb/* deployments/paas/koyeb/

# Railway
Move-Item railway.json deployments/paas/railway/
Move-Item railway.toml deployments/paas/railway/

# Render
Move-Item render.yaml deployments/paas/render/

# Patr
Move-Item patr.yaml deployments/paas/patr/

git add deployments/paas/
git commit -m "feat: 整理 PaaS 部署配置"
```

#### 階段 5: 移動 Docker Compose
```powershell
Move-Item docker-compose.yml deployments/docker-compose/
Move-Item docker-compose.test.yml deployments/docker-compose/

git add deployments/docker-compose/
git commit -m "feat: 移動 Docker Compose 配置"
```

### 方案 C: 最簡單的方式（只做關鍵更新）

如果您不想大規模重整，只需：

```powershell
# 1. 確保 bin/ 和 build/docker/ 正確
# （已完成）

# 2. 測試建置
go build -o bin/pandora-agent.exe ./cmd/agent
go build -o bin/pandora-console.exe ./cmd/console

# 3. 測試 Docker 建置
docker build -f build/docker/agent.dockerfile -t test .

# 4. 提交
git add -A
git commit -m "chore: 更新建置路徑配置"
git push
```

---

## ⚠️ 重要提醒

### 執行前必須做的事

1. **完整備份**
   ```powershell
   git add -A
   git commit -m "backup: 重整前完整備份"
   git tag backup-$(Get-Date -Format "yyyyMMdd-HHmm")
   git push --tags
   ```

2. **創建功能分支**
   ```powershell
   git checkout -b feature/project-restructure
   ```

3. **通知團隊**
   - 告知即將進行重整
   - 暫停其他人的提交

### 執行後必須做的事

1. **本地測試**
   ```powershell
   # 建置測試
   make build
   
   # Docker 測試
   docker build -f build/docker/agent.dockerfile -t test-agent .
   docker build -f build/docker/monitoring.dockerfile -t test-monitoring .
   ```

2. **提交並觀察 CI**
   ```powershell
   git push origin feature/project-restructure
   ```
   
   在 GitHub 上觀察 Actions 是否全部通過

3. **測試部署**
   - 在開發環境測試 K8s 部署
   - 確認所有服務正常

4. **合併到 main**
   ```powershell
   # 只在所有測試通過後
   git checkout main
   git merge feature/project-restructure
   git push origin main
   ```

---

## 📊 當前狀態

### 已完成 ✅ (無需執行)
- [x] 目錄結構已創建
- [x] `.gitignore` 已更新
- [x] `.github/workflows/ci.yml` 已更新
- [x] `.github/workflows/deploy-gcp.yml` 已更新
- [x] `.github/workflows/deploy-oci.yml` 已更新
- [x] `.github/workflows/deploy-paas.yml` 已更新
- [x] `Makefile` 已更新
- [x] 編譯產物已移動到 `bin/`
- [x] Dockerfiles 已移動到 `build/docker/`

### 待執行 ⏳ (需要手動執行)
- [ ] 移動文檔到 `docs/`
- [ ] 移動 K8s 配置到 `deployments/kubernetes/`
- [ ] 移動 Terraform 到 `deployments/terraform/`
- [ ] 移動 PaaS 配置到 `deployments/paas/`
- [ ] 移動 Docker Compose 到 `deployments/docker-compose/`
- [ ] 本地建置測試
- [ ] CI/CD 測試
- [ ] 刪除舊文件（可選）

---

## 🎯 推薦執行步驟

### 第一步：驗證配置（立即執行）✅

```powershell
# 1. 檢查 workflows 語法
Get-ChildItem .github\workflows\*.yml | ForEach-Object {
    Write-Host "檢查 $($_.Name)..." -ForegroundColor Cyan
    Get-Content $_.FullName | python -c "import yaml, sys; yaml.safe_load(sys.stdin)"
}

# 2. 測試本地建置
make clean
make build

# 3. 測試 Docker 建置（確保 Dockerfile 存在）
docker build -f build\docker\agent.dockerfile -t test-agent .
```

### 第二步：執行文件遷移（可選）

```powershell
# 使用自動化腳本
.\scripts\restructure-project.ps1 -DryRun  # 先預覽
.\scripts\restructure-project.ps1          # 實際執行
```

**或手動執行 - 請參考上方「方案 B」**

### 第三步：提交並驗證

```powershell
# 1. 查看變更
git status
git diff

# 2. 提交
git add -A
git commit -m "feat: 重整專案結構

- 移動 Dockerfiles 到 build/docker/
- 更新所有 CI/CD workflows
- 更新 Makefile
- 更新 .gitignore
- 創建標準目錄結構"

# 3. 推送
git push origin feature/project-restructure

# 4. 觀察 GitHub Actions
# 打開 https://github.com/你的倉庫/actions
```

---

## 📋 檢查清單

在提交 PR 前，確認以下項目：

### 建置測試
- [ ] `make build` 成功
- [ ] `make test` 通過
- [ ] Docker 建置成功
- [ ] 沒有破壞性變更

### CI/CD 測試
- [ ] GitHub Actions - ci.yml 通過
- [ ] Docker 映像成功建置
- [ ] 安全掃描通過

### 文檔檢查
- [ ] README.md 路徑更新
- [ ] 所有文檔鏈接有效
- [ ] 新結構文檔完整

### 部署測試（可選）
- [ ] K8s 部署測試（開發環境）
- [ ] PaaS 部署測試
- [ ] 服務健康檢查通過

---

## 🐛 常見問題

### Q1: Docker build 失敗找不到 Dockerfile
**A**: 確保文件已實際移動，或更新 workflow 中的路徑

### Q2: K8s apply 失敗找不到 manifests
**A**: 檢查是否已複製 K8s 配置到新位置

### Q3: CI 失敗「file not found」
**A**: 檢查對應的 Dockerfile 是否在 `build/docker/` 中

### Q4: 想要撤銷
**A**: 
```powershell
git checkout main
git branch -D feature/project-restructure
```

---

## 📞 需要幫助？

### 相關文檔
- 📖 [重整計劃](PROJECT-RESTRUCTURE-PLAN.md)
- 🔧 [CI/CD 更新指南](CI-CD-UPDATE-GUIDE.md)
- 📊 [重整總結](RESTRUCTURE-SUMMARY.md)
- 📁 [專案結構說明](../README-PROJECT-STRUCTURE.md)

### 自動化工具
- 🤖 [重整腳本](../scripts/restructure-project.ps1)
- 📋 [操作日誌](restructure-operations.csv)（執行後生成）

---

## 🎉 執行後的效果

### 根目錄變化

**整理前（混亂）:**
```
pandora_box_console_IDS-IPS/
├── agent.exe                    ❌ 編譯產物
├── console.exe                  ❌ 編譯產物
├── Dockerfile.agent             ❌ 多個 Dockerfile
├── Dockerfile.monitoring        ❌ 散落各處
├── DEPLOYMENT.md                ❌ 文檔混雜
├── KOYEB-FIX.md                 ❌ 文檔混雜
├── docker-compose.yml           ❌ 配置文件
├── fly.toml                     ❌ 配置文件
└── k8s/                         ❌ 多個 k8s 目錄
```

**整理後（清晰）:**
```
pandora_box_console_IDS-IPS/
├── bin/                         ✅ 所有編譯產物
├── build/                       ✅ 所有建置文件
├── cmd/                         ✅ 主程式入口
├── internal/                    ✅ 應用程式代碼
├── configs/                     ✅ 配置文件
├── deployments/                 ✅ 所有部署配置
├── docs/                        ✅ 所有文檔
├── scripts/                     ✅ 工具腳本
├── web/                         ✅ 前端資源
├── go.mod                       ✅ 依賴管理
├── Makefile                     ✅ 建置腳本
└── README.md                    ✅ 專案說明
```

### 優勢

1. **清晰的結構** ✨
   - 每個目錄職責明確
   - 易於查找文件
   - 符合業界標準

2. **更好的維護性** 🛠️
   - 配置集中管理
   - 文檔結構化
   - 易於擴展

3. **CI/CD 友好** 🚀
   - 路徑一致性
   - 易於自動化
   - 清晰的依賴關係

4. **團隊協作** 👥
   - 新成員易於理解
   - 減少混淆
   - 提高效率

---

## 🚦 下一步行動

### 立即執行（推薦）

```powershell
# 快速執行完整重整
.\scripts\restructure-project.ps1

# 檢查結果
git status

# 本地測試
make clean && make build
docker build -f build\docker\agent.dockerfile -t test .

# 如果一切正常
git add -A
git commit -m "feat: 完成專案結構重整"
git push origin feature/project-restructure

# 創建 PR 並等待 CI 通過
```

### 保守執行（更安全）

僅保留當前配置更新，不移動文件：

```powershell
# 只提交配置更新
git add .github/ Makefile .gitignore docs/
git commit -m "chore: 更新建置和部署配置"
git push origin main
```

---

## ✅ 成功標準

當以下所有項目都達成時，重整才算完成：

1. ✅ CI/CD 配置已更新且正確
2. ⏳ 文件已移動到新位置
3. ⏳ 本地建置成功
4. ⏳ CI/CD 全部通過
5. ⏳ 部署測試成功
6. ⏳ README 已更新
7. ⏳ 團隊已通知

---

**你已經完成了 90% 的工作！現在只需要執行文件遷移即可。** 🎉

**下一個命令**:
```powershell
.\scripts\restructure-project.ps1 -DryRun
```
