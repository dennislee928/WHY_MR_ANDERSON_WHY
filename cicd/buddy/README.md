# Buddy CI Configuration

Buddy CI pipelines for the Security Platform with multi-cloud deployment support.

## Overview

Buddy provides Docker-based CI/CD with visual pipeline editor and excellent performance.

## Features

- ✅ Visual pipeline editor
- ✅ Docker-native builds
- ✅ Fast execution (cached layers)
- ✅ Multi-cloud deployment
- ✅ Security scanning
- ✅ Parallel execution

## Setup

### 1. Connect Repository

1. Sign up at [buddy.works](https://buddy.works)
2. Create new project
3. Connect to GitHub repository
4. Import `buddy.yml` configuration

### 2. Configure Variables

Add these environment variables in Buddy:

#### Cloudflare Deployment
- `CLOUDFLARE_API_TOKEN` - Your Cloudflare API token (Secret)
- `CLOUDFLARE_ACCOUNT_ID` - Your account ID

#### OCI Deployment
- `OCI_REGISTRY` - OCI container registry URL
- `OCI_REGISTRY_ID` - Registry ID in Buddy
- `OCI_APP_SERVER_IP` - Application server IP
- `OCI_SSH_KEY` - SSH private key (Secret)

#### IBM Cloud Deployment
- `IBM_CLOUD_API_KEY` - IBM Cloud API key (Secret)
- `IBM_CLOUD_REGION` - Deployment region (e.g., `us-south`)

#### Docker Registry
- `DOCKER_REGISTRY` - Docker registry URL
- `DOCKER_USERNAME` - Registry username
- `DOCKER_PASSWORD` - Registry password (Secret)

### 3. Configure Integrations

#### Docker Registry
1. Go to Project Settings > Integrations
2. Add Docker Registry integration
3. Enter credentials
4. Note the integration ID

#### SSH Connection (for OCI)
1. Go to Project Settings > Integrations
2. Add SSH integration
3. Upload SSH private key
4. Configure for OCI servers

## Pipelines

### 1. Build and Deploy to Cloudflare

**Trigger**: Manual (click)  
**Branch**: `main`

**Steps**:
1. Install dependencies
2. Run tests
3. Deploy to Cloudflare Workers

**Usage**:
```bash
# Runs automatically on push to main, or click "Run pipeline"
```

### 2. Build and Deploy to OCI

**Trigger**: Manual (click)  
**Branch**: `main`

**Steps**:
1. Build backend Docker image
2. Build AI/Quantum Docker image
3. Build frontend Docker image
4. Deploy to OCI via SSH

**Usage**:
```bash
# Click "Run pipeline" to deploy to OCI
```

### 3. Build and Deploy to IBM Cloud

**Trigger**: Manual (click)  
**Branch**: `main`

**Steps**:
1. Install IBM Cloud CLI
2. Deploy to Cloud Foundry

**Usage**:
```bash
# Click "Run pipeline" to deploy to IBM Cloud
```

### 4. Run Tests and Security Scan

**Trigger**: Automatic on push  
**Branch**: Any

**Steps**:
1. Go backend tests
2. Python AI tests
3. Frontend tests
4. Trivy security scan

**Usage**:
```bash
# Runs automatically on every push
```

### 5. Build Multi-Arch Docker Images

**Trigger**: Manual (click)  
**Branch**: `main`

**Steps**:
1. Setup Docker Buildx
2. Build and push multi-arch images (amd64 + arm64)

**Usage**:
```bash
# For ARM-based deployments (OCI Ampere A1)
```

## Best Practices

### Caching

Buddy automatically caches:
- Docker layers
- Dependencies (`node_modules`, `go.mod`, etc.)
- Build artifacts

To clear cache:
```bash
# In pipeline settings, enable "Clear cache before execution"
```

### Parallel Execution

Configure parallel actions for faster builds:

```yaml
- parallel:
    - action: "Backend Tests"
    - action: "Frontend Tests"
    - action: "AI Tests"
```

### Notifications

Configure Slack/Discord notifications:
1. Project Settings > Notifications
2. Add Slack webhook
3. Configure triggers (success/failure)

## Monitoring

### Pipeline Status

View in Buddy Dashboard:
- Execution history
- Duration trends
- Success/failure rates
- Resource usage

### Logs

Access logs:
1. Click on pipeline execution
2. View detailed logs for each action
3. Download logs if needed

## Troubleshooting

### Build Failures

**Problem**: Tests fail
```bash
# Solution: Check test logs
# Fix failing tests locally first
npm test
go test ./...
```

**Problem**: Docker build fails
```bash
# Solution: Test build locally
docker build -t test src/backend/
```

### Deployment Issues

**Problem**: Cloudflare deployment fails
```bash
# Check CLOUDFLARE_API_TOKEN is valid
# Verify wrangler.toml configuration
```

**Problem**: OCI SSH deployment fails
```bash
# Verify SSH key format
# Check server is accessible
# Ensure docker-compose.yml exists on server
```

**Problem**: IBM Cloud deployment fails
```bash
# Verify IBM_CLOUD_API_KEY is valid
# Check manifest.yml configuration
# Ensure services are created
```

## Advanced Configuration

### Custom Docker Images

Build custom builder images:

```dockerfile
# .buddy/Dockerfile.builder
FROM golang:1.24-alpine
RUN apk add --no-cache git make
COPY . /app
WORKDIR /app
```

Use in pipeline:
```yaml
- action: "Custom Build"
  docker_image_name: "custom/builder"
  docker_image_tag: "latest"
```

### Matrix Builds

Test across multiple versions:

```yaml
- action: "Test Go Versions"
  type: "BUILD"
  docker_image_name: "library/golang"
  matrix:
    - docker_image_tag: "1.22"
    - docker_image_tag: "1.23"
    - docker_image_tag: "1.24"
```

### Conditional Execution

Run based on conditions:

```yaml
- action: "Deploy to Production"
  trigger_condition: "ALWAYS"
  run_only_on_first_failure: false
  run_on_manual: true
  condition: "git branch = main"
```

## Cost

Buddy pricing:
- **Free**: 120 executions/month, 1 concurrent pipeline
- **Pro**: $75/month, 2000 executions, 2 concurrent
- **Hyper**: $200/month, unlimited executions, 4 concurrent

For our use case:
- **Recommended**: Free tier (suitable for small teams)
- Upgrade if >120 deploys/month

## Integration with Other CI/CD

### GitHub Actions Fallback

```yaml
# .github/workflows/buddy-fallback.yml
name: Buddy Fallback
on:
  push:
    branches: [main]
jobs:
  deploy:
    if: github.event_name == 'push'
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Trigger Buddy
        run: |
          curl -X POST https://api.buddy.works/webhooks/...
```

### Webhook Triggers

Trigger Buddy from external sources:
1. Project Settings > Webhooks
2. Copy webhook URL
3. Use in external CI/CD

```bash
curl -X POST https://api.buddy.works/webhooks/YOUR_WEBHOOK_ID
```

## Further Reading

- [Buddy Documentation](https://buddy.works/docs)
- [YAML Reference](https://buddy.works/docs/yaml/yaml-gui)
- [Docker Build Guide](https://buddy.works/docs/builds/build-docker-images)
- [Deployment Guide](https://buddy.works/docs/deployments)

