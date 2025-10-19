# GitHub Workflow 修正報告

## 修正日期
2025-10-09 (Updated: Final Fix)

## 檔案
`.github/workflows/build-onpremise-installers.yml`

## 問題描述

### 1. 原始錯誤
```
(Line: 100, Col: 16): Unrecognized named-value: 'matrix'
(Line: 128, Col: 16): Unrecognized named-value: 'matrix'
```

### 2. PowerShell 執行錯誤
```
Unexpected token '}' in expression or statement.
```

## 根本原因

1. **Matrix 表達式語法錯誤**
   - 多層條件運算符沒有正確使用括號
   - 嵌套的 `${{ }}` 表達式在 `run` 區塊中無效

2. **Emoji 字符編碼問題**
   - Windows PowerShell runner 無法正確處理 emoji 字符（✅、❌、⚠️）
   - 導致 PowerShell 腳本解析錯誤

3. **中文字符編碼問題（最終根本原因）**
   - YAML 多行字串與 PowerShell 的中文字符組合導致解析錯誤
   - GitHub Actions 在將 YAML 轉換為臨時 PowerShell 腳本時出現編碼問題
   - 特別是在多行 if-else 語句中包含中文字符時，會導致 "Unexpected token '}'" 錯誤

## 修正措施

### 1. 分離構建步驟

**修正前：**
```yaml
- name: 構建後端二進位檔案
  shell: ${{ matrix.os == 'windows-latest' && 'powershell' || 'bash' }}
  env:
    GOOS: ${{ matrix.os == 'ubuntu-latest' && 'linux' || matrix.os == 'windows-latest' && 'windows' || 'darwin' }}
  run: |
    ${{ matrix.os == 'windows-latest' && 'New-Item...' || 'mkdir...' }}
```

**修正後：**
```yaml
- name: 構建後端二進位檔案 (Windows)
  if: matrix.os == 'windows-latest'
  shell: powershell
  env:
    GOOS: windows
  run: |
    New-Item -ItemType Directory -Force -Path "dist/backend"
    # ... Windows 特定腳本

- name: 構建後端二進位檔案 (Linux/Mac)
  if: matrix.os != 'windows-latest'
  shell: bash
  env:
    GOOS: ${{ matrix.os == 'ubuntu-latest' && 'linux' || 'darwin' }}
  run: |
    mkdir -p dist/backend
    # ... Linux/Mac 特定腳本
```

### 2. 移除 Emoji 字符

**替換規則：**
- `✅` → `[SUCCESS]`
- `❌` → `[ERROR]`
- `⚠️` → `[WARNING]`
- `🏷️` → 移除
- `📅` → 移除
- `📝` → 移除

**修正範圍：**
- 版本資訊輸出（prepare job）
- 後端構建步驟（Windows 和 Linux/Mac）
- 配置檔案複製步驟（Windows 和 Linux/Mac）
- 前端構建步驟

### 3. 將所有訊息改為英文（最終解決方案）

**問題：**即使移除 emoji，中文訊息仍然導致 PowerShell 解析錯誤。

**解決方案：**
- 將所有步驟名稱和訊息改為英文
- 使用單行 PowerShell 命令，以分號分隔
- 簡化腳本結構，移除不必要的註釋和空行

**範例：**
```powershell
# 修正前（多行 + 中文）
if (Test-Path "configs") {
  Copy-Item -Path "configs" -Destination "dist/backend/" -Recurse -Force
  Write-Host "[SUCCESS] configs 目錄已複製" -ForegroundColor Green
} else {
  Write-Host "[WARNING] configs 目錄不存在" -ForegroundColor Yellow
}

# 修正後（單行 + 英文）
if (Test-Path "configs") { Copy-Item -Path "configs" -Destination "dist/backend/" -Recurse -Force; Write-Host "[SUCCESS] configs directory copied" -ForegroundColor Green } else { Write-Host "[WARNING] configs directory not found" -ForegroundColor Yellow }
```

### 4. 修正表達式括號

**修正前：**
```yaml
shell: ${{ matrix.os == 'windows-latest' && 'powershell' || 'bash' }}
```

**修正後：**
```yaml
shell: ${{ (matrix.os == 'windows-latest' && 'powershell') || 'bash' }}
```

## 額外修正（測試後發現的問題）

### 5. Linux DEB 套件構建目錄問題

**錯誤**：`dpkg-deb: error: unable to create 'dist/...' : No such file or directory`

**修正**：在構建 DEB 套件前添加 `mkdir -p dist`

### 6. Inno Setup 語言檔案不存在

**錯誤**：`Couldn't open include file "ChineseSimplified.isl"`

**修正**：
- 移除不存在的中文語言檔案引用
- 只保留英文語言 (Default.isl)
- 將所有中文描述改為英文

### 7. Inno Setup 檔案路徑錯誤

**錯誤**：`No files found matching "installer\installer\backend\*"`（雙重路徑）

**修正**：
- 在 installer 目錄中執行 ISCC（使用 Push-Location/Pop-Location）
- Source 路徑使用相對於腳本的正確路徑：`backend\*` 和 `frontend\*`
- 輸出到 `installer\dist`，然後移動到根目錄的 `dist`

## 驗證結果

✅ YAML 語法驗證通過
```bash
python -c "import yaml; yaml.safe_load(open('.github/workflows/build-onpremise-installers.yml'))"
```

✅ 無嵌套表達式問題
```bash
grep -n '\${{.*\${{' .github/workflows/build-onpremise-installers.yml
# 無結果
```

✅ Matrix 引用正確
- 所有 `matrix.os` 和 `matrix.arch` 引用都在正確的上下文中
- 條件表達式使用正確的括號

✅ 目錄結構正確
- Linux DEB: 輸出到 `dist/` 目錄
- Windows Installer: 從 `installer/dist/` 移動到 `dist/` 目錄

✅ Inno Setup 配置正確
- 只使用英文語言檔案
- 路徑相對於腳本位置正確

### 8. GitHub Release 被跳過問題

**問題**：`創建 GitHub Release` 作業被跳過

**原因**：
- Release 條件只允許 tag 觸發：`if: startsWith(github.ref, 'refs/tags/v')`
- 但 workflow 從 `dev` 分支觸發，不滿足條件

**修正**：
1. 添加 `dev` 分支到 workflow 觸發條件
2. 修改 Release 條件，允許 `dev` 分支觸發
3. 設定 `dev` 分支的 Release 為 prerelease
4. 在 Release 名稱中標記 dev 版本

**修正後**：
```yaml
# 觸發條件
on:
  push:
    branches:
      - main
      - dev  # ← 新增

# Release 條件
if: startsWith(github.ref, 'refs/tags/v') || 
    (github.ref == 'refs/heads/main' && github.event_name == 'push') || 
    (github.ref == 'refs/heads/dev' && github.event_name == 'push')

# Release 設定
name: Pandora Box Console v${{ needs.prepare.outputs.version }}${{ github.ref == 'refs/heads/dev' && '-dev' || '' }}
prerelease: ${{ github.ref == 'refs/heads/dev' }}
```

### 9. Chocolatey 服務不可用問題

**問題**：Inno Setup 安裝失敗，Chocolatey 返回 503 錯誤

**原因**：Chocolatey 服務器暫時不可用（Service Unavailable: Back-end server is at capacity）

**修正**：添加備用安裝方法
1. 首先嘗試 Chocolatey 安裝
2. 如果失敗，直接從官方網站下載 Inno Setup 安裝程式
3. 使用靜默安裝參數自動安裝

**修正後**：
```powershell
try {
  choco install innosetup -y
  if ($LASTEXITCODE -ne 0) { throw "Chocolatey installation failed" }
} catch {
  Write-Host "Chocolatey failed, trying direct download..." -ForegroundColor Yellow
  $innosetupUrl = "https://files.jrsoftware.org/is/6/innosetup-6.5.4.exe"
  $installerPath = "$env:TEMP\innosetup-installer.exe"
  Invoke-WebRequest -Uri $innosetupUrl -OutFile $installerPath
  Start-Process -FilePath $installerPath -ArgumentList "/SILENT", "/DIR=C:\Program Files (x86)\Inno Setup 6" -Wait
}
```

### 10. DEB 套件版本號格式錯誤

**問題**：`dpkg-deb: error: parsing file 'debian/DEBIAN/control': 'Version' field value 'e698947': version number does not start with digit`

**原因**：
- DEB/RPM 版本號必須以數字開頭
- `git describe --tags --always` 可能返回 Git commit hash（如 `e698947`）
- 這不符合套件管理器的版本號格式要求

**修正**：改善版本號生成邏輯
1. 優先使用正式 tag 版本
2. 如果有 tag，使用 `tag.commit_count+hash` 格式
3. 如果沒有 tag，使用 `0.1.0+branch.hash` 格式
4. 確保所有版本號都以數字開頭

**修正後**：
```bash
# 版本號生成邏輯
if [[ "${{ github.event_name }}" == "workflow_dispatch" && -n "${{ github.event.inputs.version }}" ]]; then
  VERSION="${{ github.event.inputs.version }}"
elif [[ "${{ github.ref }}" == refs/tags/* ]]; then
  VERSION="${GITHUB_REF#refs/tags/v}"
else
  # 嘗試取得最近的 tag
  LATEST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
  if [[ -n "$LATEST_TAG" ]]; then
    TAG_VERSION="${LATEST_TAG#v}"
    COMMIT_COUNT=$(git rev-list --count HEAD ^$(git rev-list --max-parents=0 HEAD))
    SHORT_HASH=$(git rev-parse --short HEAD)
    VERSION="${TAG_VERSION}.${COMMIT_COUNT}+${SHORT_HASH}"
  else
    BRANCH_NAME="${{ github.ref_name }}"
    SHORT_HASH=$(git rev-parse --short HEAD)
    VERSION="0.1.0+${BRANCH_NAME}.${SHORT_HASH}"
  fi
fi

# 確保版本號以數字開頭
if [[ ! "$VERSION" =~ ^[0-9] ]]; then
  VERSION="0.1.0+${VERSION}"
fi
```

**版本號格式範例**：
- 正式 tag：`v1.2.3` → `1.2.3`
- 有 tag 的開發版本：`1.2.3.15+e698947`
- 無 tag 的分支：`0.1.0+dev.e698947`

### 11. ISO Volume ID 字串過長問題

**問題**：`genisoimage: Volume ID string too long`

**原因**：
- `genisoimage` 的 Volume ID 必須不超過 32 個字符
- 原始 Volume ID：`"Pandora Box Console ${VERSION}"`
- 當 VERSION 為 `0.1.0+main.f6d3712` 時，總長度超過 32 字符

**修正**：添加 Volume ID 長度檢查和截斷邏輯
1. 首先嘗試使用 `PandoraBox-${VERSION}` 格式
2. 如果超過 32 字符，截斷版本號並使用 `PandoraBox-v${VERSION:0:20}` 格式
3. 確保 Volume ID 始終符合 genisoimage 的 32 字符限制

**修正後**：
```bash
# 創建簡短的 Volume ID（限制 32 字符）
VOLUME_ID="PandoraBox-${VERSION}"
if [[ ${#VOLUME_ID} -gt 32 ]]; then
  VOLUME_ID="PandoraBox-v${VERSION:0:20}"
fi

genisoimage \
  -o dist/pandora-box-console-${VERSION}-amd64.iso \
  -b isolinux/isolinux.bin \
  -c isolinux/boot.cat \
  -no-emul-boot \
  -boot-load-size 4 \
  -boot-info-table \
  -J -R -V "${VOLUME_ID}" \
  iso/
```

**Volume ID 範例**：
- 短版本號：`PandoraBox-1.2.3` (17 字符)
- 長版本號：`PandoraBox-v0.1.0+main.f6d3712` (31 字符)
- 超長版本號：`PandoraBox-v0.1.0+main.f6d3712` → 截斷為 `PandoraBox-v0.1.0+main.f6d37` (32 字符)

### 12. GitHub Releases 需要 tag 問題

**問題**：`⚠️ GitHub Releases requires a tag`

**原因**：
- GitHub Releases API 要求必須基於 Git tag 創建
- 我們嘗試從 `dev` 分支創建 release，但沒有對應的 tag
- `softprops/action-gh-release` 需要 tag 才能創建 release

**修正**：為 dev 分支自動創建 tag
1. 在創建 release 之前，為 dev 分支自動創建一個 tag
2. 使用格式 `dev-{version}` 作為 tag 名稱
3. 確保 checkout 有完整的 git 歷史記錄和推送權限

**修正後**：
```yaml
# 觸發條件（支援 main、dev 分支和正式 tag）
if: startsWith(github.ref, 'refs/tags/v') || 
    (github.ref == 'refs/heads/main' && github.event_name == 'push') || 
    (github.ref == 'refs/heads/dev' && github.event_name == 'push')

# Checkout 配置
- name: Checkout 程式碼
  uses: actions/checkout@v4
  with:
    fetch-depth: 0
    token: ${{ secrets.GITHUB_TOKEN }}

# 創建 tag 步驟（支援 main 和 dev 分支）
- name: 為分支創建 tag
  if: github.ref == 'refs/heads/main' || github.ref == 'refs/heads/dev'
  run: |
    if [[ "${{ github.ref }}" == "refs/heads/main" ]]; then
      TAG_NAME="v${{ needs.prepare.outputs.version }}"
    else
      TAG_NAME="dev-${{ needs.prepare.outputs.version }}"
    fi
    
    git config --local user.email "action@github.com"
    git config --local user.name "GitHub Action"
    
    # 檢查 tag 是否已存在，如果存在則刪除後重新創建
    if git rev-parse "$TAG_NAME" >/dev/null 2>&1; then
      echo "Tag $TAG_NAME already exists, deleting it first"
      git tag -d "$TAG_NAME"
      git push origin ":refs/tags/$TAG_NAME" || true
    fi
    
    git tag -a "$TAG_NAME" -m "Release ${{ needs.prepare.outputs.version }} from ${{ github.ref_name }}"
    git push origin "$TAG_NAME"

# Release 配置
- name: 創建 Release
  uses: softprops/action-gh-release@v1
  with:
    tag_name: ${{ github.ref == 'refs/heads/main' && format('v{0}', needs.prepare.outputs.version) || 
                   (github.ref == 'refs/heads/dev' && format('dev-{0}', needs.prepare.outputs.version)) || 
                   github.ref_name }}
    name: Pandora Box Console v${{ needs.prepare.outputs.version }}${{ github.ref == 'refs/heads/dev' && '-dev' || '' }}
    prerelease: ${{ github.ref == 'refs/heads/dev' }}
```

**Tag 命名策略**：
- **main 分支**：`v{version}` → `v0.1.0+main.f6d3712`（正式 release）
- **dev 分支**：`dev-{version}` → `dev-0.1.0+dev.e698947`（prerelease）
- **手動 tag**：使用原始 tag 名稱 → `v1.2.3`（正式 release）

### 13. Release 資產上傳失敗問題

**問題**：`Error: Failed to upload release asset index.js. received status code 404`

**原因**：
- `actions/download-artifact@v4` 下載的檔案結構包含多層目錄
- 使用 `release-artifacts/**/*` glob 模式可能導致路徑問題
- 多個 artifact 中可能有同名檔案（如 `index.js`）導致衝突
- Tag 創建和 Release 創建之間可能存在時間差

**修正**：整理檔案結構並改善上傳邏輯
1. 下載所有 artifacts 到 `release-artifacts/` 目錄
2. 將所有檔案扁平化複製到 `release-files/` 目錄
3. 使用簡單的 glob 模式 `release-files/*` 上傳
4. 添加等待時間確保 tag 完全同步

**修正後**：
```yaml
- name: 下載所有構建產物
  uses: actions/download-artifact@v4
  with:
    path: release-artifacts

- name: 整理構建產物
  run: |
    # 列出下載的檔案結構
    echo "Downloaded artifacts structure:"
    ls -R release-artifacts/
    
    # 將所有檔案移到根目錄，方便上傳
    mkdir -p release-files
    find release-artifacts -type f -exec cp {} release-files/ \;
    
    echo "Files to upload:"
    ls -lh release-files/

- name: 為分支創建 tag
  if: github.ref == 'refs/heads/main' || github.ref == 'refs/heads/dev'
  run: |
    # ... 創建 tag ...
    git push origin "$TAG_NAME"
    
    # 等待 tag 同步
    echo "Waiting for tag to be available..."
    sleep 5

- name: 創建 Release
  uses: softprops/action-gh-release@v1
  with:
    files: release-files/*  # 使用扁平化的檔案目錄
```

**改進效果**：
- 避免檔案路徑問題
- 解決同名檔案衝突
- 確保 tag 完全同步後再創建 release
- 更清晰的檔案上傳日誌

### 14. 上傳不必要的檔案導致失敗

**問題**：`Error: Failed to upload release asset base.js. received status code 404`

**原因**：
- `find release-artifacts -type f -exec cp {} release-files/ \;` 複製了**所有檔案**
- 包括 frontend node_modules 中的數千個 `.js` 檔案
- 檔案數量過多導致上傳超時或失敗
- 同名檔案被覆蓋導致檔案丟失

**修正**：只複製實際的構建產物
1. 使用特定的檔案模式篩選（`.exe`, `.deb`, `.rpm`, `.iso` 等）
2. 為 macOS 二進制檔案創建 `.tar.gz` 壓縮包
3. 添加錯誤抑制（`2>/dev/null || true`）避免找不到檔案時失敗
4. 顯示實際上傳的檔案列表和數量

**修正後**：
```bash
# 只複製實際的構建產物
mkdir -p release-files

# Windows 安裝程式（只從 dist 目錄）
find release-artifacts -path "*/dist/*.exe" -type f -exec cp {} release-files/ \; 2>/dev/null || true

# Linux 套件
find release-artifacts -name "*.deb" -type f -exec cp {} release-files/ \; 2>/dev/null || true
find release-artifacts -name "*.rpm" -type f -exec cp {} release-files/ \; 2>/dev/null || true

# ISO 映像
find release-artifacts -name "*.iso" -type f -exec cp {} release-files/ \; 2>/dev/null || true
find release-artifacts -name "*.iso.md5" -type f -exec cp {} release-files/ \; 2>/dev/null || true

# macOS 二進制檔案（創建壓縮包）
if [ -d "release-artifacts/backend-macos-latest-amd64" ]; then
  cd release-artifacts/backend-macos-latest-amd64
  tar -czf ../../release-files/pandora-box-console-${VERSION}-darwin-amd64.tar.gz *
  cd ../..
fi

if [ -d "release-artifacts/backend-macos-latest-arm64" ]; then
  cd release-artifacts/backend-macos-latest-arm64
  tar -czf ../../release-files/pandora-box-console-${VERSION}-darwin-arm64.tar.gz *
  cd ../..
fi
```

**上傳的檔案清單**：
- ✅ `pandora-box-console-{version}-windows-amd64-setup.exe` - Windows 安裝程式
- ✅ `pandora-box-console_{version}_amd64.deb` - Debian/Ubuntu 套件
- ✅ `pandora-box-console-{version}.rpm` - RedHat/CentOS 套件
- ✅ `pandora-box-console-{version}-darwin-amd64.tar.gz` - macOS Intel 二進制檔案
- ✅ `pandora-box-console-{version}-darwin-arm64.tar.gz` - macOS Apple Silicon 二進制檔案
- ✅ `pandora-box-console-{version}-amd64.iso` - ISO 安裝光碟
- ✅ `pandora-box-console-{version}-amd64.iso.md5` - ISO 校驗和
- ✅ `packer-config.pkr.hcl` - OVA 虛擬機配置

**改進效果**：
- 大幅減少上傳檔案數量（從數千個減少到 ~8 個）
- 避免上傳不必要的 node_modules 檔案
- 包含所有平台的構建產物（Windows、Linux、macOS）
- 更快的上傳速度和更高的成功率

**進一步優化**：只上傳實際的構建產物

**問題**：上傳了過多不需要的檔案（如 frontend node_modules 中的數千個 .js 檔案）

**最終修正**：
```bash
# 只複製實際的構建產物
mkdir -p release-files

# Windows 安裝程式
find release-artifacts -name "*.exe" -type f -exec cp {} release-files/ \;

# Linux 套件
find release-artifacts -name "*.deb" -type f -exec cp {} release-files/ \;
find release-artifacts -name "*.rpm" -type f -exec cp {} release-files/ \;

# ISO 映像
find release-artifacts -name "*.iso" -type f -exec cp {} release-files/ \;
find release-artifacts -name "*.iso.md5" -type f -exec cp {} release-files/ \;

# Packer 配置
find release-artifacts -name "*.hcl" -type f -exec cp {} release-files/ \;

# OVA 檔案（可選）
find release-artifacts -name "*.ova" -type f -exec cp {} release-files/ \; || true
```

**上傳的檔案類型**：
- ✅ `.exe` - Windows 安裝程式
- ✅ `.deb` - Debian/Ubuntu 套件
- ✅ `.rpm` - RedHat/CentOS 套件
- ✅ `.iso` - ISO 安裝光碟
- ✅ `.iso.md5` - ISO 校驗和
- ✅ `.hcl` - Packer 配置
- ✅ `.ova` - OVA 虛擬機（如果有）

**不再上傳**：
- ❌ `.js` 檔案（frontend 源碼）
- ❌ `node_modules` 內容
- ❌ 其他開發檔案

## 改進效果

### 1. 可讀性提升
- 分離的步驟更容易閱讀和維護
- 每個平台的腳本獨立，邏輯清晰

### 2. 可維護性提升
- 減少複雜的條件表達式
- 更容易添加平台特定邏輯

### 3. 可靠性提升
- 消除編碼相關的執行錯誤
- 減少 GitHub Actions 解析錯誤的風險

## 建議

### 後續測試
1. 觸發 workflow 進行完整測試
2. 驗證各平台構建產物
3. 檢查上傳的 artifacts

### 最佳實踐
1. 避免在腳本中使用 emoji 字符
2. 為不同平台創建獨立步驟
3. 使用 `if` 條件而非複雜的三元表達式
4. 保持 YAML 表達式簡潔

## 相關文件
- [GitHub Actions 文檔 - 表達式](https://docs.github.com/en/actions/learn-github-actions/expressions)
- [GitHub Actions 文檔 - Matrix](https://docs.github.com/en/actions/using-jobs/using-a-matrix-for-your-jobs)
- [PowerShell 編碼最佳實踐](https://docs.microsoft.com/en-us/powershell/scripting/dev-cross-plat/vscode/understanding-file-encoding)

## 提交建議

```bash
git add .github/workflows/build-onpremise-installers.yml
git commit -m "fix: 修正 build-onpremise-installers workflow 語法錯誤

主要修正：
- 分離 Windows 和 Linux/Mac 構建步驟，避免複雜的條件表達式
- 移除所有 emoji 字符，使用純文字標記（[SUCCESS], [ERROR], [WARNING]）
- 修正 matrix 表達式的括號問題
- 改善腳本可讀性和可維護性

解決問題：
- Unrecognized named-value: 'matrix' 錯誤
- PowerShell 腳本 'Unexpected token' 解析錯誤
- Windows runner 的 emoji 編碼問題
"
```

