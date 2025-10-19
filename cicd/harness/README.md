# Harness CI/CD Configuration

Enterprise-grade CI/CD pipeline using Harness for multi-cloud deployment.

## Overview

Harness provides:
- Visual pipeline builder
- Multi-cloud deployment orchestration
- Canary and blue-green deployments
- Automated rollback
- Cost optimization
- Enterprise governance

## Prerequisites

1. Harness account (sign up at https://harness.io)
2. GitHub/GitLab repository connected
3. Kubernetes cluster (for OCI deployment)
4. Cloud provider credentials (Cloudflare, IBM)

## Setup

### 1. Create Harness Project

```bash
# Via UI:
# 1. Login to Harness
# 2. Create Organization (if needed)
# 3. Create Project: "security-platform"
```

### 2. Configure Connectors

#### GitHub Connector
1. Project Setup > Connectors > New Connector
2. Select "GitHub"
3. Name: `github_connector`
4. URL: `https://github.com/your-org/security-platform`
5. Credentials: Personal Access Token or SSH Key
6. Test connection

#### Docker Registry Connector
1. Connectors > New Connector > Docker Registry
2. Name: `docker_hub`
3. Registry URL: `https://index.docker.io/v1/`
4. Username & Password
5. Test connection

#### Kubernetes Connector (for OCI)
1. Connectors > New Connector > Kubernetes Cluster
2. Name: `kubernetes_connector`
3. Select "Delegate" or "Direct K8s"
4. Upload kubeconfig or use delegate
5. Test connection

### 3. Configure Secrets

Add these secrets in Project Settings > Secrets:

```bash
# Cloudflare
- cloudflare_api_token
- cloudflare_account_id (as text secret)

# IBM Cloud
- ibm_cloud_api_key

# Slack (for notifications)
- slack_webhook

# Docker Hub
- docker_hub_username
- docker_hub_password
```

### 4. Create Services

#### Cloudflare Workers Service
1. Services > New Service
2. Name: `cloudflare_workers`
3. Deployment Type: Custom
4. Add manifest (wrangler config)

#### OCI Kubernetes Service
1. Services > New Service
2. Name: `security_platform_oci`
3. Deployment Type: Kubernetes
4. Manifest: K8s manifests from repo

#### IBM Cloud Foundry Service
1. Services > New Service
2. Name: `ibm_cloud_foundry`
3. Deployment Type: Custom
4. Add manifest.yml

### 5. Create Environments

```bash
# Production Environments
- cloudflare_production
- oci_production
- ibm_production

# Staging/Development
- cloudflare_staging
- oci_staging
- ibm_staging
```

### 6. Create Infrastructure Definitions

For each environment, define infrastructure:

**OCI Kubernetes**:
- Type: Kubernetes
- Cluster: kubernetes_connector
- Namespace: security-platform

**Cloudflare**:
- Type: Custom
- Script-based deployment

**IBM Cloud**:
- Type: Custom
- Cloud Foundry deployment

### 7. Import Pipeline

```bash
# Option 1: Via UI
1. Pipelines > Create Pipeline
2. Import from Git
3. Select repository: github_connector
4. Path: cicd/harness/.harness/pipelines/multi-cloud-deployment.yaml

# Option 2: Via CLI
harness pipeline create \
  --file cicd/harness/.harness/pipelines/multi-cloud-deployment.yaml \
  --project security-platform
```

## Usage

### Run Pipeline

#### Via UI
1. Pipelines > multi-cloud-deployment
2. Click "Run"
3. Select branch (e.g., `main`)
4. Choose deployment strategy (rolling/canary/bluegreen)
5. Click "Run Pipeline"

#### Via CLI
```bash
# Install Harness CLI
curl -sSfL https://github.com/harness/harness-cli/releases/latest/download/harness-cli-linux-amd64 -o harness
chmod +x harness
sudo mv harness /usr/local/bin/

# Login
harness login --api-key <your-api-key>

# Run pipeline
harness pipeline execute \
  --project security-platform \
  --org default \
  --pipeline multi-cloud-deployment \
  --branch main
```

#### Via API
```bash
curl -X POST \
  https://app.harness.io/gateway/pipeline/api/pipeline/execute/security-platform/multi-cloud-deployment \
  -H "x-api-key: ${HARNESS_API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "runtimeInputYaml": "pipeline:\n  variables:\n  - name: deployment_strategy\n    value: rolling"
  }'
```

## Deployment Strategies

### Rolling Deployment (Default)

Gradually replaces old instances:
```yaml
deployment_strategy: rolling
```

- Zero downtime
- Gradual rollout
- Easy rollback

### Canary Deployment

Deploys to subset of instances first:
```yaml
deployment_strategy: canary
```

- Test with small traffic percentage
- Gradual increase
- Automatic rollback on failure

### Blue-Green Deployment

Maintains two environments:
```yaml
deployment_strategy: bluegreen
```

- Zero downtime
- Instant rollback
- Double resources temporarily

## Monitoring

### Pipeline Dashboard

View in Harness UI:
- Execution history
- Success/failure rates
- Deployment frequency
- Mean time to recovery (MTTR)

### Step Logs

```bash
# Via UI
1. Select pipeline execution
2. Click on any step
3. View detailed logs

# Via CLI
harness pipeline logs \
  --execution-id <execution-id>
```

### Metrics

Harness collects:
- Deployment frequency
- Lead time for changes
- Change failure rate
- MTTR

## Approvals

Manual approval required before production deployment:

```yaml
- step:
    type: HarnessApproval
    spec:
      approvers:
        userGroups:
          - account.admin
      minimumCount: 1
```

Approvers can:
- Approve deployment
- Reject deployment
- Add comments

## Notifications

### Slack Notifications

Configured in pipeline:
```yaml
notificationRules:
  - name: Pipeline Success
    pipelineEvents:
      - type: PipelineSuccess
    notificationMethod:
      type: Slack
      spec:
        slackWebhookUrl: <+secrets.getValue("slack_webhook")>
```

### Email Notifications

Add in Harness UI:
1. Pipeline > Notifications
2. Add Notification Rule
3. Select "Email"
4. Add recipients

### Custom Webhooks

```yaml
notificationMethod:
  type: Webhook
  spec:
    webhookUrl: https://your-webhook.com/notify
    method: POST
    body: |
      {
        "pipeline": "<+pipeline.name>",
        "status": "<+pipeline.status>"
      }
```

## Rollback

### Automatic Rollback

Enabled by default on deployment failure.

### Manual Rollback

```bash
# Via UI
1. Deployments > Select deployment
2. Click "Rollback"

# Via CLI
harness deployment rollback \
  --execution-id <execution-id>
```

## Cost Management

### Cloud Cost Dashboard

Harness tracks cloud costs:
1. Cloud Costs > Dashboard
2. View costs by service, environment, cloud
3. Set budgets and alerts

### Recommendations

Harness provides cost optimization recommendations:
- Right-sizing
- Spot instances
- Reserved instances
- Idle resource cleanup

## Governance

### OPA Policies

Define policies for deployments:

```rego
# policy.rego
package deployment

deny["Must have manual approval for production"] {
  input.environment == "production"
  input.approval == false
}

deny["Cannot deploy on Friday"] {
  time.weekday(time.now_ns()) == "Friday"
}
```

Apply in Harness:
1. Policies > New Policy
2. Upload policy.rego
3. Apply to pipelines

### RBAC

Configure role-based access:
1. Access Control > Roles
2. Define custom roles
3. Assign to users/groups

Example roles:
- **Developer**: View pipelines, trigger builds
- **DevOps**: Full pipeline management
- **Admin**: Full access

## Troubleshooting

### Pipeline Failure

```bash
# Check step logs
1. Select failed execution
2. Click on failed step
3. View error logs

# Common issues:
- Missing secrets
- Invalid connector
- Resource quota exceeded
- Timeout
```

### Deployment Stuck

```bash
# Check infrastructure
kubectl get pods -n security-platform

# Check Harness delegate
kubectl logs -n harness-delegate-ng <delegate-pod>

# Retry deployment
# UI: Click "Retry Failed Pipeline"
```

### Connector Issues

```bash
# Test connectors
1. Connectors > Select connector
2. Click "Test Connection"
3. View error details

# Common fixes:
- Update credentials
- Check network connectivity
- Verify firewall rules
```

## Performance Tips

1. **Use Delegates**: For faster deployments
2. **Enable Caching**: Docker layer caching
3. **Parallel Steps**: Run tests in parallel
4. **Optimize Images**: Use multi-stage builds
5. **Limit History**: Keep 10-20 deployments

## Cost

Harness pricing:
- **Free Tier**: 
  - 5 services
  - Unlimited users
  - Community support

- **Paid Tiers**:
  - Team: $100/month (25 services)
  - Enterprise: Custom pricing

For our use case:
- **Recommended**: Free tier (5 services covers all clouds)

## Integration with Other CI/CD

### Trigger from GitHub Actions

```yaml
# .github/workflows/trigger-harness.yml
name: Trigger Harness
on:
  push:
    branches: [main]
jobs:
  trigger:
    runs-on: ubuntu-latest
    steps:
      - name: Trigger Harness Pipeline
        run: |
          curl -X POST \
            https://app.harness.io/gateway/pipeline/api/webhook/execute \
            -H "x-api-key: ${{ secrets.HARNESS_API_KEY }}" \
            -d '{"branch": "main"}'
```

### Webhook from Buddy

In Buddy pipeline:
```yaml
- action: "Trigger Harness"
  type: "BUILD"
  execute_commands:
    - |
      curl -X POST https://app.harness.io/api/webhook/...
```

## Further Reading

- [Harness Documentation](https://docs.harness.io/)
- [Pipeline Studio](https://docs.harness.io/article/2chyf1acil-add-a-stage)
- [CD Strategies](https://docs.harness.io/article/0zvnn5s1ph-deployment-concepts-and-strategies)
- [Cloud Cost Management](https://docs.harness.io/category/0u951ug1es-cloud-cost-management)
- [Best Practices](https://docs.harness.io/article/9u8jk9p3yg-harness-continuous-delivery-best-practices)

