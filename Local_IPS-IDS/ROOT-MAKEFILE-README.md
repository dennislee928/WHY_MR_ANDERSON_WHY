# 根目錄 Makefile 說明

> **重要**: 專案現在有兩個 Makefile，用途不同

---

## 📋 Makefile 用途區分

### 1. **根目錄 Makefile**（./Makefile）

**用途**: 整體專案管理

**主要功能**:
- Docker 操作（build, push, deploy）
- 服務管理（start, stop, restart）
- 整合測試
- 文檔生成
- 監控和健康檢查

**使用場景**:
```bash
# 在專案根目錄
make help           # 查看所有命令
make docker-build   # 構建 Docker 映像
make deploy         # 部署服務
make test           # 執行測試
```

**不適用**:
- 編譯單個 Go 程式
- 後端開發構建

---

### 2. **Application/be/Makefile**

**用途**: 後端程式構建（專用）

**主要功能**:
- 編譯 Go 二進位檔案
- 跨平台構建
- 後端開發和測試
- 打包發行版

**使用場景**:
```bash
# 在 Application/be/ 目錄
cd Application/be
make all          # 編譯所有程式
make agent        # 只編譯 Agent
make run-agent    # 編譯並執行
make info         # 顯示配置
```

**適用**:
- 後端開發
- 程式編譯
- 單元測試

---

## 🎯 推薦工作流程

### 場景 1: 開發後端

```bash
cd Application/be
make all
make run-agent
```

### 場景 2: 開發前端

```bash
cd Application/Fe
npm run dev
```

### 場景 3: 完整構建（地端部署）

```bash
cd Application
# Windows
.\build-local.ps1

# Linux/macOS
./build-local.sh
```

### 場景 4: Docker 部署

```bash
# 在根目錄
make docker-build
make deploy
```

### 場景 5: 整合測試

```bash
# 在根目錄
make full-test
```

---

## 📊 Makefile 對照表

| 功能 | 根目錄 | Application/be/ |
|------|--------|-----------------|
| 編譯 Go | ❌ (委派) | ✅ 主要 |
| Docker 操作 | ✅ 主要 | ❌ |
| 服務管理 | ✅ 主要 | ❌ |
| 跨平台編譯 | ❌ | ✅ 主要 |
| 打包發行 | ❌ | ✅ 主要 |
| 文檔生成 | ✅ 主要 | ❌ |
| 監控操作 | ✅ 主要 | ❌ |

---

## 💡 最佳實踐

### ✅ 推薦

```bash
# 後端開發
cd Application/be && make run-agent

# 完整構建
cd Application && ./build-local.sh

# Docker 部署
make docker-build && make deploy
```

### ❌ 不推薦

```bash
# 不要在根目錄直接編譯（應該用 Application/be/）
make build-agent

# 不要混用兩個 Makefile 的命令
```

---

## 🔄 未來改進

可能會考慮：
1. 統一為單一 Makefile
2. 使用 make 的 include 機制
3. 或保持當前雙 Makefile 結構

目前建議：**保持雙 Makefile**，職責分明。

---

---

## 🎯 Phase 1-3 完成後的 Makefile 用途

### 當前系統狀態（v3.0.0）

系統現在是完整的雲原生微服務架構，有多種部署方式：

**1. 本地開發（Application/be/Makefile）**
```bash
cd Application/be
make all          # 編譯所有服務
make run-agent    # 運行 Agent
```

**2. Docker Compose 部署（根目錄 Makefile）**
```bash
cd deployments/onpremise
docker-compose up -d
```

**3. Kubernetes 部署（Helm）**
```bash
cd deployments/helm
helm install pandora-box ./pandora-box
```

**4. GitOps 部署（ArgoCD）**
```bash
kubectl apply -f deployments/argocd/application.yaml
```

---

## 📚 相關文檔

- 📖 [完整專案結構](docs/COMPLETE-PROJECT-STRUCTURE.md) ⭐ 最新
- 📖 [Kubernetes 部署指南](docs/KUBERNETES-DEPLOYMENT.md)
- 📖 [GitOps 指南](docs/GITOPS-ARGOCD.md)
- 📖 [微服務快速啟動](docs/MICROSERVICES-QUICKSTART.md)

---

**維護**: Pandora Security Team  
**版本**: 3.0.0 (AI 驅動智能安全平台)  
**狀態**: 🏆 世界級生產就緒  
**最後更新**: 2025-10-09

