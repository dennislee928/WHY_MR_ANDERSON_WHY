# IBM Cloud Deployment

Deploy the Security Platform to IBM Cloud using the Lite (free) tier.

## IBM Cloud Lite Tier

IBM Cloud offers a generous free tier with:

### Compute
- **Cloud Foundry Apps**: 256 MB memory per app (multiple apps allowed)
- **Functions**: 400,000 GB-seconds of execution time/month
- **Container Registry**: 500 MB storage

### Databases
- **Cloudant NoSQL**: 1 GB storage
- **Db2**: 200 MB storage
- **PostgreSQL** (Compose): Limited free tier
- **Redis**: 256 MB memory

### Storage
- **Object Storage**: 25 GB/month
- **Block Storage**: Limited free tier

### Services
- **Watson AI**: Various free tier limits
- **API Connect**: 50,000 API calls/month

## Prerequisites

1. IBM Cloud account (sign up at https://cloud.ibm.com/registration)
2. IBM Cloud CLI:
   ```bash
   # macOS
   curl -fsSL https://clis.cloud.ibm.com/install/osx | sh

   # Linux
   curl -fsSL https://clis.cloud.ibm.com/install/linux | sh

   # Windows (PowerShell)
   iex (New-Object Net.WebClient).DownloadString('https://clis.cloud.ibm.com/install/powershell')
   ```
3. Cloud Foundry CLI plugin:
   ```bash
   ibmcloud cf install
   ```

## Setup

### 1. Login to IBM Cloud

```bash
ibmcloud login

# Or with SSO
ibmcloud login --sso

# Select your region
ibmcloud target -r us-south

# Target Cloud Foundry org and space
ibmcloud target --cf
```

### 2. Create Services

#### PostgreSQL Database

```bash
# Create Databases for PostgreSQL instance (Lite plan)
ibmcloud resource service-instance-create security-platform-db \
  databases-for-postgresql standard us-south \
  -p '{
    "members_memory_allocation_mb": 1024,
    "members_disk_allocation_mb": 5120
  }'

# Get connection details
ibmcloud resource service-key-create security-platform-db-key \
  Manager --instance-name security-platform-db
```

#### Redis

```bash
# Create Databases for Redis instance (Lite plan)
ibmcloud resource service-instance-create security-platform-redis \
  databases-for-redis standard us-south \
  -p '{
    "members_memory_allocation_mb": 256,
    "members_disk_allocation_mb": 1024
  }'

# Get connection details
ibmcloud resource service-key-create security-platform-redis-key \
  Manager --instance-name security-platform-redis
```

#### Object Storage (for backups)

```bash
# Create Object Storage instance
ibmcloud resource service-instance-create security-platform-storage \
  cloud-object-storage lite global

# Create service credentials
ibmcloud resource service-key-create security-platform-storage-key \
  Writer --instance-name security-platform-storage
```

### 3. Configure Environment

Create `.env` files in each application directory:

**Backend (`src/backend/.env`)**:
```bash
DB_HOST=<postgres-host-from-credentials>
DB_PORT=<postgres-port>
DB_USER=<postgres-user>
DB_PASSWORD=<postgres-password>
DB_NAME=<postgres-database>

REDIS_HOST=<redis-host>
REDIS_PORT=<redis-port>
REDIS_PASSWORD=<redis-password>

IBM_QUANTUM_TOKEN=<your-token>
```

**AI/Quantum (`src/ai-quantum/.env`)**:
```bash
DB_HOST=<postgres-host>
DB_PORT=<postgres-port>
DB_USER=<postgres-user>
DB_PASSWORD=<postgres-password>

REDIS_HOST=<redis-host>
REDIS_PORT=<redis-port>

IBM_QUANTUM_TOKEN=<your-token>
```

**Frontend (`src/frontend/.env.production`)**:
```bash
NEXT_PUBLIC_API_URL=https://security-platform-api.mybluemix.net
NEXT_PUBLIC_WS_URL=wss://security-platform-api.mybluemix.net
```

## Deployment

### Option 1: Cloud Foundry (Recommended for Lite Tier)

```bash
cd infrastructure/cloud-configs/ibm

# Deploy all applications
ibmcloud cf push

# Or deploy individually
ibmcloud cf push security-platform-api
ibmcloud cf push security-platform-ai
ibmcloud cf push security-platform-frontend
```

### Option 2: Kubernetes (Code Engine)

```bash
# Create Code Engine project
ibmcloud ce project create --name security-platform

# Target the project
ibmcloud ce project select --name security-platform

# Build and deploy backend
cd src/backend
ibmcloud ce application create --name backend \
  --image us.icr.io/security-platform/backend:latest \
  --port 3001 \
  --min-scale 0 --max-scale 1 \
  --cpu 0.25 --memory 0.5G

# Build and deploy AI service
cd ../ai-quantum
ibmcloud ce application create --name ai-quantum \
  --image us.icr.io/security-platform/ai-quantum:latest \
  --port 8000 \
  --min-scale 0 --max-scale 1 \
  --cpu 0.5 --memory 1G

# Build and deploy frontend
cd ../frontend
ibmcloud ce application create --name frontend \
  --image us.icr.io/security-platform/frontend:latest \
  --port 3000 \
  --min-scale 0 --max-scale 1 \
  --cpu 0.25 --memory 0.5G
```

### Option 3: Container Registry + Virtual Servers

```bash
# Build Docker images
docker build -t us.icr.io/security-platform/backend:latest src/backend
docker build -t us.icr.io/security-platform/ai-quantum:latest src/ai-quantum
docker build -t us.icr.io/security-platform/frontend:latest src/frontend

# Login to IBM Container Registry
ibmcloud cr login

# Push images
docker push us.icr.io/security-platform/backend:latest
docker push us.icr.io/security-platform/ai-quantum:latest
docker push us.icr.io/security-platform/frontend:latest
```

## Access Services

After deployment, access your applications:

- **Frontend**: https://security-platform.mybluemix.net
- **Backend API**: https://security-platform-api.mybluemix.net
- **AI/Quantum API**: https://security-platform-ai.mybluemix.net

## Monitoring

### Cloud Foundry Logs

```bash
# Tail logs
ibmcloud cf logs security-platform-api --recent
ibmcloud cf logs security-platform-ai --recent

# Stream logs
ibmcloud cf logs security-platform-api
```

### Metrics Dashboard

Access in IBM Cloud Console:
- Resource List > Cloud Foundry Apps > (select app)
- View CPU, Memory, Network metrics

### Health Checks

```bash
# Check application health
ibmcloud cf app security-platform-api
ibmcloud cf app security-platform-ai

# Check service bindings
ibmcloud cf services
```

## Scaling

### Vertical Scaling

Update memory/instances in `manifest.yml`:
```yaml
applications:
  - name: security-platform-api
    memory: 512M  # Increase memory
    instances: 2  # Increase instances
```

Redeploy:
```bash
ibmcloud cf push security-platform-api
```

### Horizontal Scaling

```bash
# Scale instances
ibmcloud cf scale security-platform-api -i 2

# Scale memory
ibmcloud cf scale security-platform-api -m 512M
```

## Cost Management

### Stay Within Lite Tier

✅ **Free Resources**:
- Cloud Foundry: 256 MB × multiple apps
- Databases: Lite plans available
- Object Storage: 25 GB/month
- Container Registry: 500 MB

⚠️ **Monitor Usage**:
- Check IBM Cloud Console > Billing
- Set up spending notifications
- Review monthly usage reports

### Optimize Costs

1. **Use Lite plan services**: PostgreSQL, Redis Lite plans
2. **Minimize instances**: Start with 1 instance per app
3. **Use autoscaling**: Scale down during low traffic
4. **Clean up unused resources**: Remove old apps and services

## Backups

### Database Backups

```bash
# PostgreSQL backup (manual)
ibmcloud cdb deployment-backups-list <deployment-id>

# Trigger on-demand backup
ibmcloud cdb deployment-backup-now <deployment-id>
```

### Object Storage Backups

```bash
# Install rclone
brew install rclone  # macOS
# or apt-get install rclone  # Linux

# Configure IBM COS
rclone config

# Sync backups
rclone sync /local/backups ibmcos:security-platform-backups
```

## Troubleshooting

### Application Won't Start

```bash
# Check logs
ibmcloud cf logs security-platform-api --recent

# Check events
ibmcloud cf events security-platform-api

# Restart application
ibmcloud cf restart security-platform-api
```

### Service Binding Issues

```bash
# List services
ibmcloud cf services

# Unbind and rebind
ibmcloud cf unbind-service security-platform-api security-platform-db
ibmcloud cf bind-service security-platform-api security-platform-db
ibmcloud cf restage security-platform-api
```

### Memory/Quota Issues

```bash
# Check quota
ibmcloud cf quotas

# Check app resource usage
ibmcloud cf app security-platform-api

# Scale down if needed
ibmcloud cf scale security-platform-api -m 128M
```

## CI/CD Integration

### Using IBM Continuous Delivery

```bash
# Create toolchain
ibmcloud dev toolchain-create

# Or use the UI:
# IBM Cloud Console > DevOps > Toolchains > Create Toolchain
```

### GitHub Actions

See `cicd/` directory for GitHub Actions workflows that deploy to IBM Cloud.

## Cleanup

### Remove Applications

```bash
# Delete apps
ibmcloud cf delete security-platform-api -f
ibmcloud cf delete security-platform-ai -f
ibmcloud cf delete security-platform-frontend -f
```

### Remove Services

```bash
# Delete service instances
ibmcloud resource service-instance-delete security-platform-db -f
ibmcloud resource service-instance-delete security-platform-redis -f
ibmcloud resource service-instance-delete security-platform-storage -f
```

## Security Best Practices

1. **Use service credentials**: Don't hardcode passwords
2. **Enable TLS**: All connections use HTTPS by default
3. **Rotate keys**: Regularly rotate service credentials
4. **Use IAM**: Configure proper access controls
5. **Monitor logs**: Review security logs regularly

## Performance Tuning

### Cloud Foundry Apps

```yaml
# manifest.yml optimizations
applications:
  - name: security-platform-api
    memory: 256M
    disk_quota: 1G
    instances: 1
    health-check-type: http
    health-check-http-endpoint: /health
    timeout: 180
```

### Database Optimization

- Use connection pooling
- Implement caching with Redis
- Index frequently queried columns
- Monitor slow queries

## Further Reading

- [IBM Cloud Documentation](https://cloud.ibm.com/docs)
- [Cloud Foundry Documentation](https://docs.cloudfoundry.org/)
- [IBM Cloud CLI Reference](https://cloud.ibm.com/docs/cli)
- [Databases for PostgreSQL](https://cloud.ibm.com/docs/databases-for-postgresql)
- [Code Engine](https://cloud.ibm.com/docs/codeengine)

