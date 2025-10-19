# ArgoCD GitOps Configuration

GitOps-based continuous deployment using ArgoCD for Kubernetes clusters.

## Overview

ArgoCD provides declarative, GitOps continuous delivery for Kubernetes with:
- Git as source of truth
- Automatic synchronization
- Multi-cluster management
- Rollback capabilities
- Health monitoring

## Prerequisites

1. Kubernetes cluster (1.20+)
2. kubectl configured
3. Helm 3 (optional, for Helm charts)

## Installation

### 1. Install ArgoCD

```bash
# Create namespace
kubectl create namespace argocd

# Install ArgoCD
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for pods to be ready
kubectl wait --for=condition=Ready pods --all -n argocd --timeout=300s
```

### 2. Access ArgoCD UI

```bash
# Port forward to ArgoCD server
kubectl port-forward svc/argocd-server -n argocd 8080:443

# Get initial admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d

# Access at: https://localhost:8080
# Username: admin
# Password: (from above command)
```

### 3. Install ArgoCD CLI (Optional)

```bash
# macOS
brew install argocd

# Linux
curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
chmod +x /usr/local/bin/argocd

# Login
argocd login localhost:8080
```

## Configuration

### 1. Apply Application Configuration

```bash
# Apply main application
kubectl apply -f cicd/argocd/application.yaml

# Apply multi-environment ApplicationSet
kubectl apply -f cicd/argocd/applicationset.yaml

# Apply notifications
kubectl apply -f cicd/argocd/notifications.yaml
```

### 2. Configure Repository Access

```bash
# Add Git repository
argocd repo add https://github.com/your-org/security-platform.git \
  --username your-username \
  --password your-token

# Or via SSH
argocd repo add git@github.com:your-org/security-platform.git \
  --ssh-private-key-path ~/.ssh/id_rsa
```

### 3. Configure Image Registry

```bash
# Add Docker registry credentials
kubectl create secret docker-registry regcred \
  --docker-server=https://index.docker.io/v1/ \
  --docker-username=your-username \
  --docker-password=your-password \
  --docker-email=your-email \
  -n security-platform
```

### 4. Configure Notifications

Edit `notifications.yaml` and add tokens:

```yaml
stringData:
  slack-token: "xoxb-your-slack-bot-token"
  github-token: "ghp_your-github-token"
```

Apply:
```bash
kubectl apply -f cicd/argocd/notifications.yaml
```

## Usage

### Deploy Application

#### Option 1: Via UI
1. Open ArgoCD UI (https://localhost:8080)
2. Click "New App"
3. Fill in details or import from application.yaml
4. Click "Create"
5. Click "Sync"

#### Option 2: Via CLI
```bash
# Create application
argocd app create security-platform \
  --repo https://github.com/your-org/security-platform.git \
  --path infrastructure/kubernetes \
  --dest-server https://kubernetes.default.svc \
  --dest-namespace security-platform

# Sync application
argocd app sync security-platform

# Watch sync progress
argocd app wait security-platform
```

#### Option 3: Via kubectl
```bash
kubectl apply -f cicd/argocd/application.yaml
```

### Monitor Deployment

```bash
# Get application status
argocd app get security-platform

# View sync history
argocd app history security-platform

# View application logs
argocd app logs security-platform

# Watch application
argocd app watch security-platform
```

### Rollback

```bash
# List revisions
argocd app history security-platform

# Rollback to previous revision
argocd app rollback security-platform

# Rollback to specific revision
argocd app rollback security-platform 5
```

## Multi-Environment Deployment

The ApplicationSet configuration automatically creates applications for multiple environments:

```bash
# View all applications
argocd app list

# You should see:
# - security-platform-development
# - security-platform-staging
# - security-platform-production

# Sync specific environment
argocd app sync security-platform-production
```

## Multi-Cloud Deployment

Deploy to multiple Kubernetes clusters:

### 1. Add Clusters

```bash
# Add OCI cluster
argocd cluster add oci-cluster \
  --name oci-cluster \
  --kubeconfig ~/.kube/oci-config

# Add IBM cluster
argocd cluster add ibm-cluster \
  --name ibm-cluster \
  --kubeconfig ~/.kube/ibm-config

# List clusters
argocd cluster list
```

### 2. Deploy to Multiple Clouds

```bash
# Apply multi-cloud ApplicationSet
kubectl apply -f cicd/argocd/applicationset.yaml

# This creates:
# - security-platform-oci-production
# - security-platform-ibm-production
```

## Sync Strategies

### Automatic Sync

Enabled in `application.yaml`:
```yaml
syncPolicy:
  automated:
    prune: true      # Delete resources not in Git
    selfHeal: true   # Override manual changes
```

### Manual Sync

Disable automatic sync:
```yaml
syncPolicy: {}
```

Then sync manually:
```bash
argocd app sync security-platform
```

### Sync Waves

Control deployment order with annotations:

```yaml
metadata:
  annotations:
    argocd.argoproj.io/sync-wave: "1"  # Deploy first
```

Order:
1. Wave 0: Namespaces, ConfigMaps, Secrets
2. Wave 1: Databases, Redis
3. Wave 2: Backend services
4. Wave 3: Frontend services

## Health Checks

ArgoCD automatically monitors health. Custom health checks:

```yaml
# In application.yaml
health:
  check:
    - kind: Deployment
      jsonPath: .status.conditions[?(@.type=="Progressing")].status
      value: "True"
```

## Secrets Management

### Option 1: Sealed Secrets

```bash
# Install Sealed Secrets
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml

# Seal a secret
kubeseal --format yaml < secret.yaml > sealed-secret.yaml

# Commit sealed-secret.yaml to Git
```

### Option 2: External Secrets Operator

```bash
# Install External Secrets Operator
helm repo add external-secrets https://charts.external-secrets.io
helm install external-secrets external-secrets/external-secrets -n external-secrets --create-namespace

# Configure to use cloud secret managers (AWS Secrets Manager, GCP Secret Manager, etc.)
```

### Option 3: Vault Integration

```bash
# Install Vault
helm repo add hashicorp https://helm.releases.hashicorp.com
helm install vault hashicorp/vault

# Configure Vault ArgoCD plugin
```

## Monitoring and Alerts

### Prometheus Metrics

ArgoCD exposes Prometheus metrics:
```yaml
# ServiceMonitor for Prometheus
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: argocd-metrics
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: argocd-metrics
  endpoints:
    - port: metrics
```

### Grafana Dashboards

Import ArgoCD dashboard:
1. Grafana UI > Dashboards > Import
2. ID: 14584 (ArgoCD dashboard)
3. Select Prometheus data source

### Slack Notifications

Already configured in `notifications.yaml`. Test:

```bash
# Trigger a sync
argocd app sync security-platform

# Check Slack channel for notification
```

## Troubleshooting

### Application OutOfSync

```bash
# Check diff
argocd app diff security-platform

# Force sync
argocd app sync security-platform --force
```

### Sync Failed

```bash
# View sync errors
argocd app get security-platform

# View detailed logs
kubectl logs -n argocd -l app.kubernetes.io/name=argocd-application-controller

# Retry sync
argocd app sync security-platform --retry-limit 5
```

### Health Degraded

```bash
# Check resource health
argocd app get security-platform --show-operation

# View pod status
kubectl get pods -n security-platform

# Describe problematic pods
kubectl describe pod <pod-name> -n security-platform
```

### Resource Stuck Deleting

```bash
# Remove finalizers
kubectl patch app security-platform -n argocd \
  -p '{"metadata":{"finalizers":[]}}' \
  --type merge

# Force delete
argocd app delete security-platform --cascade=false
```

## Best Practices

1. **Use Git Branches**: Separate branches for environments
   ```
   - main → production
   - staging → staging
   - develop → development
   ```

2. **Structured Directories**:
   ```
   infrastructure/kubernetes/
   ├── base/           # Base manifests
   ├── overlays/
   │   ├── development/
   │   ├── staging/
   │   └── production/
   ```

3. **Sync Waves**: Order deployments
4. **Health Checks**: Define custom health checks
5. **Notifications**: Configure Slack/email alerts
6. **RBAC**: Restrict access per environment
7. **Secrets**: Use Sealed Secrets or external managers

## Performance Tuning

```yaml
# In argocd-cm ConfigMap
data:
  # Increase concurrent syncs
  application.resourceTrackingMethod: annotation
  
  # Optimize resource compare
  resource.compareoptions: |
    ignoreAggregatedRoles: true
  
  # Adjust sync timeout
  timeout.reconciliation: 180s
```

## Cost

ArgoCD is **100% free** and open-source!

- No licensing costs
- Community support
- Enterprise support available (optional)

## Further Reading

- [ArgoCD Documentation](https://argo-cd.readthedocs.io/)
- [ApplicationSet Documentation](https://argo-cd.readthedocs.io/en/stable/user-guide/application-set/)
- [GitOps Principles](https://opengitops.dev/)
- [Best Practices](https://argo-cd.readthedocs.io/en/stable/user-guide/best_practices/)

