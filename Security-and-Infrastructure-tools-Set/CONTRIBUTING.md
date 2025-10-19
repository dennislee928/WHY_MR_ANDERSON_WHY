# 貢獻指南 | Contributing Guide

[English](#english) | [繁體中文](#繁體中文)

---

## English

Thank you for your interest in contributing to the Security & Infrastructure Tools Set! This document provides guidelines for contributing to the project.

### Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

### How to Contribute

#### 1. Reporting Bugs

Before submitting a bug report:
- Check if the bug has already been reported in [Issues](https://github.com/your-username/Security-and-Infrastructure-tools-Set/issues)
- Test with the latest version
- Collect diagnostic information

**Bug Report Template**:
```markdown
**Description**: Clear description of the bug

**Steps to Reproduce**:
1. Step one
2. Step two
3. Expected vs actual behavior

**Environment**:
- OS: Ubuntu 22.04
- Docker version: 20.10.21
- Docker Compose version: 2.15.1

**Logs**:
```
paste relevant logs here
```
```

#### 2. Suggesting Features

Feature requests are welcome! Please:
- Check if the feature has been suggested
- Explain the use case
- Describe the expected behavior

**Feature Request Template**:
```markdown
**Feature**: Brief title

**Problem**: What problem does this solve?

**Proposed Solution**: How should it work?

**Alternatives**: Other solutions you've considered
```

#### 3. Contributing Code

**Workflow**:

1. **Fork the repository**
   ```bash
   # Clone your fork
   git clone https://github.com/your-username/Security-and-Infrastructure-tools-Set.git
   cd Security-and-Infrastructure-tools-Set
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

3. **Make your changes**
   - Follow existing code style
   - Add tests if applicable
   - Update documentation

4. **Test your changes**
   ```bash
   # Run tests
   cd Make_Files
   make test
   
   # Test deployment
   make up
   make health
   ```

5. **Commit with clear messages**
   ```bash
   git add .
   git commit -m "feat: add ZAP scanner integration"
   ```

   **Commit Message Format**:
   ```
   type(scope): subject

   body (optional)

   footer (optional)
   ```

   **Types**:
   - `feat`: New feature
   - `fix`: Bug fix
   - `docs`: Documentation changes
   - `style`: Code style changes (formatting, etc.)
   - `refactor`: Code refactoring
   - `test`: Adding or updating tests
   - `chore`: Maintenance tasks

6. **Push and create Pull Request**
   ```bash
   git push origin feature/amazing-feature
   ```

   Then create a PR on GitHub with:
   - Clear description of changes
   - Link to related issues
   - Screenshots/examples if applicable

### Code Style

#### Docker Compose
```yaml
# Use lowercase with underscores
services:
  scanner_nuclei:  # ✅ Good
    image: projectdiscovery/nuclei:latest
    
  ScannerNuclei:   # ❌ Avoid
```

#### Python
- Follow PEP 8
- Use type hints
- Add docstrings

```python
def parse_results(file_path: str) -> dict:
    """
    Parse scan results from JSON file.
    
    Args:
        file_path: Path to the JSON file
        
    Returns:
        Dictionary containing parsed results
    """
    pass
```

#### Shell Scripts
- Use `#!/bin/bash`
- Add error handling with `set -e`
- Add help messages

```bash
#!/bin/bash
set -e

# Script description
# Usage: ./script.sh [options]

if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <target>"
    exit 1
fi
```

### Documentation

- Update README.md for major changes
- Add docstrings to new functions
- Update relevant documentation in `docs/`
- Keep Chinese and English docs in sync

### Testing

#### Unit Tests
```bash
# Run unit tests
python -m pytest tests/unit/

# With coverage
pytest --cov=scripts tests/
```

#### Integration Tests
```bash
# Run integration tests
bash tests/integration/test_full_scan.sh
```

### Review Process

1. **Automated Checks**: CI/CD will run tests and linters
2. **Code Review**: Maintainers will review your code
3. **Feedback**: Address review comments
4. **Approval**: Minimum 1 approval required
5. **Merge**: Maintainers will merge the PR

### Getting Help

- **Questions**: [GitHub Discussions](https://github.com/your-username/Security-and-Infrastructure-tools-Set/discussions)
- **Chat**: [Discord](https://discord.gg/xxx)
- **Email**: security-tools@example.com

---

## 繁體中文

感謝您有興趣為安全與基礎設施工具集貢獻！本文件提供專案貢獻指南。

### 行為準則

參與此專案，即表示您同意維護對所有人都尊重和包容的環境。

### 如何貢獻

#### 1. 回報 Bug

提交 bug 報告前：
- 檢查該 bug 是否已在 [Issues](https://github.com/your-username/Security-and-Infrastructure-tools-Set/issues) 中回報
- 使用最新版本測試
- 收集診斷資訊

**Bug 報告範本**:
```markdown
**描述**: 清楚描述 bug

**重現步驟**:
1. 步驟一
2. 步驟二
3. 預期結果 vs 實際結果

**環境**:
- 作業系統: Ubuntu 22.04
- Docker 版本: 20.10.21
- Docker Compose 版本: 2.15.1

**日誌**:
```
貼上相關日誌
```
```

#### 2. 提出新功能

歡迎功能請求！請：
- 檢查該功能是否已被提出
- 說明使用場景
- 描述預期行為

**功能請求範本**:
```markdown
**功能**: 簡短標題

**問題**: 這解決了什麼問題？

**建議方案**: 應該如何運作？

**替代方案**: 您考慮過的其他解決方案
```

#### 3. 貢獻程式碼

**工作流程**:

1. **Fork 儲存庫**
   ```bash
   # 克隆您的 fork
   git clone https://github.com/your-username/Security-and-Infrastructure-tools-Set.git
   cd Security-and-Infrastructure-tools-Set
   ```

2. **創建功能分支**
   ```bash
   git checkout -b feature/amazing-feature
   ```

3. **進行修改**
   - 遵循現有代碼風格
   - 如適用，添加測試
   - 更新文件

4. **測試您的修改**
   ```bash
   # 執行測試
   cd Make_Files
   make test
   
   # 測試部署
   make up
   make health
   ```

5. **使用清晰的訊息提交**
   ```bash
   git add .
   git commit -m "feat: 添加 ZAP 掃描器整合"
   ```

   **Commit 訊息格式**:
   ```
   類型(範圍): 主題

   正文 (可選)

   頁腳 (可選)
   ```

   **類型**:
   - `feat`: 新功能
   - `fix`: Bug 修復
   - `docs`: 文件更新
   - `style`: 代碼風格更改（格式化等）
   - `refactor`: 代碼重構
   - `test`: 添加或更新測試
   - `chore`: 維護任務

6. **推送並創建 Pull Request**
   ```bash
   git push origin feature/amazing-feature
   ```

   然後在 GitHub 上創建 PR，包含：
   - 清楚的變更描述
   - 連結到相關 issue
   - 如適用，提供截圖/範例

### 代碼風格

#### Docker Compose
```yaml
# 使用小寫加底線
services:
  scanner_nuclei:  # ✅ 正確
    image: projectdiscovery/nuclei:latest
    
  ScannerNuclei:   # ❌ 避免
```

#### Python
- 遵循 PEP 8
- 使用型別提示
- 添加 docstrings

```python
def parse_results(file_path: str) -> dict:
    """
    從 JSON 檔案解析掃描結果。
    
    參數:
        file_path: JSON 檔案路徑
        
    返回:
        包含解析結果的字典
    """
    pass
```

#### Shell 腳本
- 使用 `#!/bin/bash`
- 使用 `set -e` 添加錯誤處理
- 添加幫助訊息

```bash
#!/bin/bash
set -e

# 腳本描述
# 用法: ./script.sh [選項]

if [ "$#" -lt 1 ]; then
    echo "用法: $0 <目標>"
    exit 1
fi
```

### 文件

- 重大變更需更新 README.md
- 為新函數添加 docstrings
- 更新 `docs/` 中的相關文件
- 保持中英文文件同步

### 測試

#### 單元測試
```bash
# 執行單元測試
python -m pytest tests/unit/

# 帶覆蓋率
pytest --cov=scripts tests/
```

#### 整合測試
```bash
# 執行整合測試
bash tests/integration/test_full_scan.sh
```

### 審查流程

1. **自動檢查**: CI/CD 會執行測試和 linter
2. **代碼審查**: 維護者會審查您的代碼
3. **反饋**: 處理審查意見
4. **批准**: 需要至少 1 個批准
5. **合併**: 維護者會合併 PR

### 獲取幫助

- **問題**: [GitHub Discussions](https://github.com/your-username/Security-and-Infrastructure-tools-Set/discussions)
- **聊天**: [Discord](https://discord.gg/xxx)
- **Email**: security-tools@example.com

---

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Thank You! 感謝您！

Your contributions make this project better for everyone. We appreciate your time and effort!

您的貢獻讓這個專案變得更好。我們感謝您的時間和努力！

