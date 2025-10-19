# GitOps with ArgoCD 部署指南

## 概述

本指南說明如何使用 ArgoCD 實現 GitOps 自動化部署流程。

---

## 什麼是 GitOps？

GitOps 是一種持續交付的方式，使用 Git 作為唯一的真實來源：

- **聲明式**: 所有配置都以聲明式方式存儲在 Git 中
- **版本控制**: 所有變更都有完整的歷史記錄
- **自動化**: 自動同步 Git 狀態到集群
- **可審計**: 所有操作都可追蹤

---

## ArgoCD 安裝

### 1. 安裝 ArgoCD

```bash
# 創建命名空間
kubectl create namespace argocd

# 安裝 ArgoCD
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# 等待所有 Pod 就緒
kubectl wait --for=condition=Ready pods --all -n argocd --timeout=300s
```

### 2. 訪問 ArgoCD UI

```bash
# 獲取初始密碼
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

# 端口轉發
kubectl port-forward svc/argocd-server -n argocd 8080:443

# 瀏覽器訪問
open https://localhost:8080

# 登入
# 用戶名: admin
# 密碼: <上面獲取的密碼>
```

### 3. 安裝 ArgoCD CLI

```bash
# macOS
brew install argocd

# Linux
curl -sSL -o argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
chmod +x argocd
sudo mv argocd /usr/local/bin/

# Windows
choco install argocd-cli
```

### 4. CLI 登入

```bash
argocd login localhost:8080 --username admin --password <password> --insecure
```

---

## 部署 Pandora Box

### 方法 1: 使用 ArgoCD UI

1. 登入 ArgoCD UI
2. 點擊 "NEW APP"
3. 填寫應用資訊：
   - **Application Name**: pandora-box
   - **Project**: default
   - **Sync Policy**: Automatic
   - **Repository URL**: https://github.com/your-org/pandora-box.git
   - **Revision**: main
   - **Path**: deployments/helm/pandora-box
   - **Destination Cluster**: https://kubernetes.default.svc
   - **Namespace**: pandora-system
4. 點擊 "CREATE"

### 方法 2: 使用 ArgoCD CLI

```bash
argocd app create pandora-box \
  --repo https://github.com/your-org/pandora-box.git \
  --path deployments/helm/pandora-box \
  --dest-server https://kubernetes.default.svc \
  --dest-namespace pandora-system \
  --sync-policy automated \
  --auto-prune \
  --self-heal
```

### 方法 3: 使用 Kubernetes Manifest

```bash
kubectl apply -f deployments/argocd/application.yaml
```

---

## 自動同步策略

### 啟用自動同步

```yaml
syncPolicy:
  automated:
    prune: true        # 自動刪除不在 Git 中的資源
    selfHeal: true     # 自動修復手動更改
    allowEmpty: false  # 不允許空應用
```

### 同步選項

```yaml
syncOptions:
  - CreateNamespace=true           # 自動創建命名空間
  - PrunePropagationPolicy=foreground  # 刪除策略
  - PruneLast=true                 # 最後刪除資源
```

### 重試策略

```yaml
retry:
  limit: 5
  backoff:
    duration: 5s
    factor: 2
    maxDuration: 3m
```

---

## 多環境管理

### 環境結構

```
deployments/
├── helm/
│   └── pandora-box/
│       ├── Chart.yaml
│       ├── values.yaml              # 預設值
│       └── values/
│           ├── dev.yaml             # 開發環境
│           ├── staging.yaml         # 測試環境
│           └── production.yaml      # 生產環境
└── argocd/
    ├── dev-application.yaml
    ├── staging-application.yaml
    └── production-application.yaml
```

### 開發環境

```yaml
# deployments/argocd/dev-application.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: pandora-box-dev
  namespace: argocd
spec:
  source:
    repoURL: https://github.com/your-org/pandora-box.git
    targetRevision: dev
    path: deployments/helm/pandora-box
    helm:
      valueFiles:
        - values.yaml
        - values/dev.yaml
  destination:
    server: https://kubernetes.default.svc
    namespace: pandora-dev
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

### 生產環境

```yaml
# deployments/argocd/production-application.yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: pandora-box-prod
  namespace: argocd
spec:
  source:
    repoURL: https://github.com/your-org/pandora-box.git
    targetRevision: main
    path: deployments/helm/pandora-box
    helm:
      valueFiles:
        - values.yaml
        - values/production.yaml
  destination:
    server: https://kubernetes.default.svc
    namespace: pandora-prod
  syncPolicy:
    automated:
      prune: false  # 生產環境不自動刪除
      selfHeal: false  # 生產環境不自動修復
```

---

## 部署工作流程

### 1. 開發流程

```bash
# 1. 創建功能分支
git checkout -b feature/new-feature

# 2. 修改配置
vim deployments/helm/pandora-box/values/dev.yaml

# 3. 提交變更
git add .
git commit -m "feat: add new feature configuration"

# 4. 推送到遠端
git push origin feature/new-feature

# 5. 創建 Pull Request
# ArgoCD 會自動檢測變更並同步到開發環境
```

### 2. 測試流程

```bash
# 1. 合併到 staging 分支
git checkout staging
git merge feature/new-feature

# 2. 推送到遠端
git push origin staging

# ArgoCD 自動同步到測試環境
```

### 3. 生產發布流程

```bash
# 1. 合併到 main 分支
git checkout main
git merge staging

# 2. 創建標籤
git tag -a v2.1.0 -m "Release v2.1.0"

# 3. 推送到遠端
git push origin main --tags

# 4. 手動同步到生產環境（安全起見）
argocd app sync pandora-box-prod
```

---

## 監控和告警

### 查看應用狀態

```bash
# 列出所有應用
argocd app list

# 查看應用詳情
argocd app get pandora-box

# 查看同步歷史
argocd app history pandora-box
```

### 查看應用健康狀態

```bash
# 健康狀態
argocd app get pandora-box --show-operation

# 資源樹
argocd app tree pandora-box
```

### 配置告警

```yaml
# argocd-notifications-cm ConfigMap
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-notifications-cm
  namespace: argocd
data:
  service.slack: |
    token: $slack-token
  
  trigger.on-deployed: |
    - when: app.status.operationState.phase in ['Succeeded']
      send: [app-deployed]
  
  trigger.on-health-degraded: |
    - when: app.status.health.status == 'Degraded'
      send: [app-health-degraded]
  
  template.app-deployed: |
    message: |
      Application {{.app.metadata.name}} has been deployed.
      Sync Status: {{.app.status.sync.status}}
    slack:
      attachments: |
        [{
          "title": "{{.app.metadata.name}}",
          "color": "good"
        }]
```

---

## 回滾操作

### 使用 ArgoCD 回滾

```bash
# 查看歷史
argocd app history pandora-box

# 回滾到特定版本
argocd app rollback pandora-box 5

# 回滾到上一個版本
argocd app rollback pandora-box
```

### 使用 Git 回滾

```bash
# 回滾 Git commit
git revert HEAD
git push origin main

# ArgoCD 會自動檢測並同步
```

---

## 進階功能

### 1. ApplicationSet

管理多個相似的應用：

```yaml
apiVersion: argoproj.io/v1alpha1
kind: ApplicationSet
metadata:
  name: pandora-box-environments
  namespace: argocd
spec:
  generators:
  - list:
      elements:
      - env: dev
        namespace: pandora-dev
        branch: dev
      - env: staging
        namespace: pandora-staging
        branch: staging
      - env: prod
        namespace: pandora-prod
        branch: main
  template:
    metadata:
      name: 'pandora-box-{{env}}'
    spec:
      source:
        repoURL: https://github.com/your-org/pandora-box.git
        targetRevision: '{{branch}}'
        path: deployments/helm/pandora-box
        helm:
          valueFiles:
          - values/{{env}}.yaml
      destination:
        server: https://kubernetes.default.svc
        namespace: '{{namespace}}'
```

### 2. Sync Waves

控制資源部署順序：

```yaml
apiVersion: v1
kind: Service
metadata:
  name: postgresql
  annotations:
    argocd.argoproj.io/sync-wave: "0"  # 先部署資料庫
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: device-service
  annotations:
    argocd.argoproj.io/sync-wave: "1"  # 再部署服務
```

### 3. Resource Hooks

在同步前後執行操作：

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: db-migration
  annotations:
    argocd.argoproj.io/hook: PreSync
    argocd.argoproj.io/hook-delete-policy: HookSucceeded
spec:
  template:
    spec:
      containers:
      - name: migrate
        image: pandora-box/migrator:latest
        command: ["./migrate.sh"]
```

---

## 安全最佳實踐

### 1. RBAC 配置

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-rbac-cm
  namespace: argocd
data:
  policy.csv: |
    p, role:developer, applications, get, */*, allow
    p, role:developer, applications, sync, */dev/*, allow
    p, role:admin, applications, *, */*, allow
    g, dev-team, role:developer
    g, ops-team, role:admin
```

### 2. 使用私有倉庫

```bash
# 添加私有 Git 倉庫
argocd repo add https://github.com/your-org/pandora-box.git \
  --username <username> \
  --password <password>

# 或使用 SSH
argocd repo add git@github.com:your-org/pandora-box.git \
  --ssh-private-key-path ~/.ssh/id_rsa
```

### 3. Sealed Secrets

加密敏感資訊：

```bash
# 安裝 Sealed Secrets
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.18.0/controller.yaml

# 創建 sealed secret
echo -n 'pandora123' | kubectl create secret generic db-password \
  --dry-run=client \
  --from-file=password=/dev/stdin \
  -o yaml | \
  kubeseal -o yaml > sealed-db-password.yaml

# 提交到 Git
git add sealed-db-password.yaml
git commit -m "Add sealed database password"
```

---

## 故障排除

### 應用無法同步

```bash
# 查看詳細錯誤
argocd app get pandora-box --show-operation

# 強制刷新
argocd app get pandora-box --refresh

# 硬同步（忽略差異）
argocd app sync pandora-box --force
```

### 資源卡在刪除狀態

```bash
# 手動刪除 finalizers
kubectl patch app pandora-box -n argocd -p '{"metadata":{"finalizers":null}}' --type=merge
```

---

## 總結

使用 ArgoCD 實現 GitOps 的優勢：

✅ **自動化**: Git push 即自動部署  
✅ **可追蹤**: 所有變更都有記錄  
✅ **可回滾**: 輕鬆回滾到任何版本  
✅ **一致性**: 確保集群狀態與 Git 一致  
✅ **多環境**: 輕鬆管理多個環境  

---

**現在你已經掌握了 GitOps 自動化部署！** 🚀

