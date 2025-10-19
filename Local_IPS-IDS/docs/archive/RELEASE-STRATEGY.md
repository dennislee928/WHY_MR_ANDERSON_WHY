# GitHub Release 策略說明

## 概述

本專案實現了自動化的 GitHub Release 策略，支援三種不同的發布方式：

1. **手動 Tag 發布**（正式版本）
2. **Main 分支自動發布**（穩定版本）
3. **Dev 分支自動發布**（開發版本）

## Release 觸發條件

### 1. 手動 Tag 發布（推薦用於正式版本）

**觸發方式**：
```bash
git tag -a v1.2.3 -m "Release version 1.2.3"
git push origin v1.2.3
```

**特點**：
- ✅ Tag 名稱：使用推送的 tag 名稱（如 `v1.2.3`）
- ✅ Release 類型：正式 Release
- ✅ 版本號：使用 tag 中的版本號
- ✅ 適用場景：生產環境發布、重要里程碑版本

### 2. Main 分支自動發布

**觸發方式**：
```bash
git push origin main
```

**特點**：
- ✅ Tag 名稱：自動創建 `v{version}` tag（如 `v0.1.0+main.f6d3712`）
- ✅ Release 類型：正式 Release
- ✅ 版本號：基於 Git 歷史自動生成
- ✅ 適用場景：穩定版本的持續發布

**自動生成的版本號格式**：
- 有 tag 的情況：`{tag}.{commit_count}+{hash}` → `1.2.3.15+f6d3712`
- 無 tag 的情況：`0.1.0+main.{hash}` → `0.1.0+main.f6d3712`

### 3. Dev 分支自動發布

**觸發方式**：
```bash
git push origin dev
```

**特點**：
- ✅ Tag 名稱：自動創建 `dev-{version}` tag（如 `dev-0.1.0+dev.e698947`）
- ✅ Release 類型：Prerelease（預發布版本）
- ✅ 版本號：基於 Git 歷史自動生成
- ✅ 適用場景：開發測試、功能驗證

**自動生成的版本號格式**：
- 有 tag 的情況：`{tag}.{commit_count}+{hash}` → `1.2.3.15+e698947`
- 無 tag 的情況：`0.1.0+dev.{hash}` → `0.1.0+dev.e698947`

## Tag 命名規則

| 觸發方式 | Tag 格式 | 範例 | Release 類型 |
|---------|---------|------|-------------|
| 手動 tag | `v{major}.{minor}.{patch}` | `v1.2.3` | 正式 Release |
| Main 分支 | `v{version}` | `v0.1.0+main.f6d3712` | 正式 Release |
| Dev 分支 | `dev-{version}` | `dev-0.1.0+dev.e698947` | Prerelease |

## 版本號生成邏輯

### 優先級

1. **手動指定版本**（workflow_dispatch 輸入）
2. **Git Tag 版本**（推送 tag 時）
3. **自動生成版本**（基於 Git 歷史）

### 自動生成規則

```bash
# 如果有最近的 tag
LATEST_TAG=$(git describe --tags --abbrev=0)
TAG_VERSION="${LATEST_TAG#v}"  # 移除 v 前綴
COMMIT_COUNT=$(git rev-list --count HEAD ^$(git rev-list --max-parents=0 HEAD))
SHORT_HASH=$(git rev-parse --short HEAD)
VERSION="${TAG_VERSION}.${COMMIT_COUNT}+${SHORT_HASH}"
# 結果：1.2.3.15+f6d3712

# 如果沒有 tag
BRANCH_NAME="${{ github.ref_name }}"
SHORT_HASH=$(git rev-parse --short HEAD)
VERSION="0.1.0+${BRANCH_NAME}.${SHORT_HASH}"
# 結果：0.1.0+main.f6d3712 或 0.1.0+dev.e698947
```

## Release 內容

每個 Release 包含以下內容：

### 構建產物

#### Windows
- `pandora-box-console-{version}-windows-amd64-setup.exe` - Windows 安裝程式

#### Linux
- `pandora-box-console-{version}_amd64.deb` - Debian/Ubuntu 套件
- `pandora-box-console-{version}.rpm` - RedHat/CentOS 套件

#### 虛擬機/光碟
- `pandora-box-console-{version}-amd64.iso` - ISO 安裝光碟
- `packer-config.pkr.hcl` - OVA 虛擬機配置

### Release 資訊

- 版本號
- 構建日期
- Git Commit Hash
- 安裝說明
- 系統需求
- 更新日誌連結

## 工作流程圖

```
┌─────────────────────────────────────────────────────────────┐
│                     推送代碼或 Tag                            │
└────────────┬────────────────────────────────────────────────┘
             │
             ├─── 手動 Tag (v1.2.3)
             │    └─→ 使用 tag 版本 → 創建正式 Release
             │
             ├─── Main 分支推送
             │    └─→ 自動創建 v{version} tag → 創建正式 Release
             │
             └─── Dev 分支推送
                  └─→ 自動創建 dev-{version} tag → 創建 Prerelease
```

## 使用建議

### 開發階段
```bash
# 在 dev 分支開發
git checkout dev
git add .
git commit -m "feat: 新功能"
git push origin dev
# ✅ 自動創建 prerelease，方便測試
```

### 穩定版本
```bash
# 合併到 main 分支
git checkout main
git merge dev
git push origin main
# ✅ 自動創建正式 release
```

### 正式發布
```bash
# 創建正式版本 tag
git tag -a v1.2.3 -m "Release version 1.2.3"
git push origin v1.2.3
# ✅ 創建帶有明確版本號的正式 release
```

## Tag 衝突處理

如果自動生成的 tag 已經存在，workflow 會：

1. 檢測到 tag 已存在
2. 刪除舊的 tag
3. 創建新的 tag
4. 推送到遠端

這確保了每次推送都能成功創建 release，不會因為 tag 衝突而失敗。

## 注意事項

1. **權限要求**：workflow 需要 `contents: write` 權限來創建 tag 和 release
2. **Tag 命名**：避免手動創建與自動生成格式相同的 tag
3. **版本號格式**：自動生成的版本號符合 Semantic Versioning 和 DEB/RPM 格式要求
4. **Prerelease**：dev 分支的 release 會被標記為 prerelease，不會出現在正式版本列表中

## 故障排除

### 問題：Release 沒有創建

**檢查項目**：
1. 確認所有前置 jobs 都成功完成
2. 檢查觸發條件是否符合（main/dev 分支或 v* tag）
3. 確認 GITHUB_TOKEN 有足夠權限

### 問題：Tag 推送失敗

**解決方法**：
1. 檢查是否有 tag 衝突
2. 確認 workflow 有推送權限
3. 查看 workflow 日誌中的詳細錯誤訊息

### 問題：版本號格式錯誤

**解決方法**：
1. 確認 Git 歷史記錄完整（fetch-depth: 0）
2. 檢查版本號生成邏輯
3. 驗證版本號是否以數字開頭

## 相關文檔

- [Workflow 修正報告](WORKFLOW-FIX-REPORT.md)
- [Workflow 修正摘要](WORKFLOW-FIXES-APPLIED.md)
- [GitHub Actions 文檔](https://docs.github.com/en/actions)
- [Semantic Versioning](https://semver.org/)

