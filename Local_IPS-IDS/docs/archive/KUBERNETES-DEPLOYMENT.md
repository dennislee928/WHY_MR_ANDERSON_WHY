# Kubernetes 部署指南

## 概述

本指南說明如何在 Kubernetes 集群上部署 Pandora Box Console IDS-IPS 微服務架構。

---

## 前置需求

### 軟體需求

- Kubernetes 1.24+
- Helm 3.10+
- kubectl 1.24+
- Docker 20.10+

### 硬體需求（最小）

| 資源 | 最小值 | 推薦值 |
|------|--------|--------|
| CPU 核心 | 4 | 8+ |
| 記憶體 | 8GB | 16GB+ |
| 儲存空間 | 100GB | 200GB+ |
| 節點數量 | 3 | 5+ |

---

## 快速開始

### 1. 創建命名空間

```bash
kubectl create namespace pandora-system
```

### 2. 使用 Helm 部署

```bash
# 添加 Helm 倉庫
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# 部署 Pandora Box
cd deployments/helm
helm install pandora-box ./pandora-box \
  --namespace pandora-system \
  --create-namespace \
  --wait
```

### 3. 驗證部署

```bash
# 檢查所有 Pod 狀態
kubectl get pods -n pandora-system

# 檢查服務
kubectl get svc -n pandora-system

# 檢查 HPA
kubectl get hpa -n pandora-system
```

---

## 詳細配置

### 自訂 values.yaml

創建自訂配置文件 `custom-values.yaml`：

```yaml
# 自訂微服務副本數
deviceService:
  replicaCount: 3
  resources:
    limits:
      memory: "512Mi"
      cpu: "1000m"

networkService:
  replicaCount: 5
  autoscaling:
    maxReplicas: 30

# 自訂資料庫配置
postgresql:
  primary:
    persistence:
      size: 100Gi
    resources:
      limits:
        memory: "2Gi"
        cpu: "2000m"

# 啟用 Ingress
ingress:
  enabled: true
  hosts:
    - host: pandora.yourdomain.com
      paths:
        - path: /
          pathType: Prefix
```

使用自訂配置部署：

```bash
helm install pandora-box ./pandora-box \
  --namespace pandora-system \
  -f custom-values.yaml
```

---

## 手動部署（不使用 Helm）

### 1. 部署 PostgreSQL

```bash
kubectl apply -f deployments/kubernetes/postgresql.yaml
```

### 2. 部署微服務

```bash
# Device Service
kubectl apply -f deployments/kubernetes/device-service.yaml

# Network Service
kubectl apply -f deployments/kubernetes/network-service.yaml

# Control Service
kubectl apply -f deployments/kubernetes/control-service.yaml
```

### 3. 創建 Secrets

```bash
# RabbitMQ Secret
kubectl create secret generic rabbitmq-secret \
  --from-literal=url='amqp://pandora:pandora123@rabbitmq:5672/' \
  -n pandora-system

# mTLS Certificates
kubectl create secret tls device-service-certs \
  --cert=certs/device-service.crt \
  --key=certs/device-service.key \
  -n pandora-system
```

---

## 自動擴展配置

### HorizontalPodAutoscaler (HPA)

系統已預配置 HPA，基於 CPU 和記憶體使用率自動擴展：

```yaml
# Device Service HPA
minReplicas: 2
maxReplicas: 10
targetCPUUtilizationPercentage: 70
targetMemoryUtilizationPercentage: 80
```

### 查看 HPA 狀態

```bash
kubectl get hpa -n pandora-system
kubectl describe hpa device-service-hpa -n pandora-system
```

### 自訂 HPA

```bash
kubectl autoscale deployment network-service \
  --cpu-percent=60 \
  --min=3 \
  --max=30 \
  -n pandora-system
```

---

## 監控和日誌

### Prometheus 監控

```bash
# 訪問 Prometheus UI
kubectl port-forward -n pandora-system svc/prometheus-server 9090:9090

# 瀏覽器訪問
open http://localhost:9090
```

### Grafana 儀表板

```bash
# 訪問 Grafana UI
kubectl port-forward -n pandora-system svc/grafana 3000:3000

# 預設登入
# 用戶名: admin
# 密碼: pandora123
```

### 查看日誌

```bash
# 查看特定服務日誌
kubectl logs -f deployment/device-service -n pandora-system

# 查看所有微服務日誌
kubectl logs -f -l component=microservice -n pandora-system

# 使用 stern（推薦）
stern -n pandora-system device-service
```

---

## 網路策略

### 啟用 NetworkPolicy

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: device-service-netpol
  namespace: pandora-system
spec:
  podSelector:
    matchLabels:
      app: device-service
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
      - podSelector:
          matchLabels:
            component: microservice
      ports:
      - protocol: TCP
        port: 50051
  egress:
    - to:
      - podSelector:
          matchLabels:
            app: rabbitmq
      ports:
      - protocol: TCP
        port: 5672
```

應用網路策略：

```bash
kubectl apply -f deployments/kubernetes/network-policies.yaml
```

---

## 故障排除

### Pod 無法啟動

```bash
# 檢查 Pod 狀態
kubectl describe pod <pod-name> -n pandora-system

# 檢查事件
kubectl get events -n pandora-system --sort-by='.lastTimestamp'

# 檢查日誌
kubectl logs <pod-name> -n pandora-system --previous
```

### 服務無法連接

```bash
# 檢查服務端點
kubectl get endpoints -n pandora-system

# 測試服務連接
kubectl run -it --rm debug --image=nicolaka/netshoot --restart=Never -n pandora-system -- bash
# 在容器內
nslookup device-service
curl http://device-service:8081/health
```

### HPA 不工作

```bash
# 檢查 metrics-server
kubectl top nodes
kubectl top pods -n pandora-system

# 如果 metrics-server 未安裝
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
```

---

## 升級和回滾

### 使用 Helm 升級

```bash
# 升級到新版本
helm upgrade pandora-box ./pandora-box \
  --namespace pandora-system \
  -f custom-values.yaml

# 查看升級歷史
helm history pandora-box -n pandora-system

# 回滾到上一個版本
helm rollback pandora-box -n pandora-system

# 回滾到特定版本
helm rollback pandora-box 2 -n pandora-system
```

### 使用 kubectl 升級

```bash
# 更新映像
kubectl set image deployment/device-service \
  device-service=pandora-box/device-service:2.1.0 \
  -n pandora-system

# 查看滾動更新狀態
kubectl rollout status deployment/device-service -n pandora-system

# 回滾部署
kubectl rollout undo deployment/device-service -n pandora-system
```

---

## 備份和恢復

### 備份 PostgreSQL

```bash
# 創建備份
kubectl exec -n pandora-system postgresql-0 -- \
  pg_dump -U pandora pandora > backup-$(date +%Y%m%d).sql

# 恢復備份
kubectl exec -i -n pandora-system postgresql-0 -- \
  psql -U pandora pandora < backup-20251009.sql
```

### 使用 Velero 備份整個命名空間

```bash
# 安裝 Velero
velero install --provider aws --bucket pandora-backups

# 備份命名空間
velero backup create pandora-backup --include-namespaces pandora-system

# 恢復備份
velero restore create --from-backup pandora-backup
```

---

## 生產環境最佳實踐

### 1. 資源限制

始終設置資源請求和限制：

```yaml
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "512Mi"
    cpu: "1000m"
```

### 2. 健康檢查

配置適當的 liveness 和 readiness 探針：

```yaml
livenessProbe:
  httpGet:
    path: /live
    port: 8081
  initialDelaySeconds: 30
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /ready
    port: 8081
  initialDelaySeconds: 5
  periodSeconds: 5
```

### 3. PodDisruptionBudget

確保高可用性：

```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: device-service-pdb
spec:
  minAvailable: 1
  selector:
    matchLabels:
      app: device-service
```

### 4. 安全性

- 使用 NetworkPolicy 限制流量
- 啟用 RBAC
- 使用 Secrets 管理敏感資訊
- 定期更新映像

### 5. 監控和告警

- 配置 Prometheus 告警規則
- 設置 Grafana 儀表板
- 啟用日誌聚合（如 Loki）

---

## 參考資源

- [Kubernetes 官方文檔](https://kubernetes.io/docs/)
- [Helm 官方文檔](https://helm.sh/docs/)
- [Prometheus Operator](https://prometheus-operator.dev/)
- [ArgoCD 文檔](https://argo-cd.readthedocs.io/)

---

**部署完成後，系統將在 Kubernetes 上以高可用、可擴展的方式運行！** 🚀

