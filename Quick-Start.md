# Quick Start Guide - Unified Security Platform

**Deploy to Cloudflare Workers, Oracle Cloud (OCI), or IBM Cloud in under 30 minutes!**

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Project Setup](#project-setup)
- [Option 1: Cloudflare Workers (Serverless)](#option-1-cloudflare-workers-serverless)
- [Option 2: Oracle Cloud - OCI (Always Free)](#option-2-oracle-cloud---oci-always-free)
- [Option 3: IBM Cloud (Lite Tier)](#option-3-ibm-cloud-lite-tier)
- [Hybrid Multi-Cloud Setup](#hybrid-multi-cloud-setup)
- [Verification](#verification)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Required Tools

Install these tools before starting:

```bash
# 1. Git
git --version

# 2. Node.js 18+ (for Cloudflare)
node --version
npm --version

# 3. Go 1.24+ (for backend development)
go version

# 4. Python 3.11+ (for AI/Quantum services)
python --version

# 5. Docker & Docker Compose
docker --version
docker-compose --version
```

### Cloud Accounts

Sign up for free accounts:

1. **Cloudflare**: https://dash.cloudflare.com/sign-up
2. **Oracle Cloud**: https://www.oracle.com/cloud/free/
3. **IBM Cloud**: https://cloud.ibm.com/registration

---

## Project Setup

### 1. Clone and Initialize

```bash
# Clone the repository
git clone <your-repository-url>
cd WHY_MR_ANDERSON_WHY

# Verify structure
ls -la
# You should see: src/, infrastructure/, cicd/, docs/, etc.
```

### 2. Environment Configuration

Create your environment file:

```bash
# Copy example
cp .env.example .env

# Edit with your values
nano .env  # or vim, code, etc.
```

**Minimum required variables**:
```bash
# Database
DB_PASSWORD=your_strong_password_here

# IBM Quantum (optional but recommended)
IBM_QUANTUM_TOKEN=your_ibm_quantum_token
```

---

## Option 1: Cloudflare Workers (Serverless)

**Best for**: Global API, low latency, zero maintenance  
**Cost**: $0/month (up to 10M requests)  
**Setup time**: ~10 minutes

### Step 1: Install Wrangler CLI

```bash
npm install -g wrangler

# Login to Cloudflare
wrangler login
# This opens a browser - authorize the CLI
```

### Step 2: Get Your Account Details

```bash
# Get your account ID
wrangler whoami

# Note down your account ID
# Example: Account ID: a1b2c3d4e5f6...
```

Or via Dashboard:
1. Go to https://dash.cloudflare.com
2. Copy the Account ID from the URL or right sidebar

### Step 3: Configure Wrangler

Edit `infrastructure/cloud-configs/cloudflare/wrangler.toml`:

```toml
account_id = "YOUR_ACCOUNT_ID_HERE"  # Paste your account ID

# Everything else can stay as default for now
```

### Step 4: Create D1 Database

```bash
cd infrastructure/cloud-configs/cloudflare

# Install dependencies
npm install

# Create D1 database
wrangler d1 create security_platform_db

# Output will show:
# [[d1_databases]]
# binding = "DB"
# database_name = "security_platform_db"
# database_id = "xxxx-xxxx-xxxx-xxxx"

# Copy the database_id and update wrangler.toml
```

Update `wrangler.toml`:
```toml
[[d1_databases]]
binding = "DB"
database_name = "security_platform_db"
database_id = "PASTE_DATABASE_ID_HERE"  # From previous step
```

### Step 5: Initialize Database

```bash
# Run the schema initialization
wrangler d1 execute security_platform_db --file=schema.sql

# Verify tables were created
wrangler d1 execute security_platform_db --command="SELECT name FROM sqlite_master WHERE type='table'"
```

### Step 6: Create KV Namespaces

```bash
# Create CACHE namespace
wrangler kv:namespace create "CACHE"
# Output: id = "xxxx..."

# Create CACHE preview namespace
wrangler kv:namespace create "CACHE" --preview
# Output: preview_id = "yyyy..."

# Create SESSIONS namespace
wrangler kv:namespace create "SESSIONS"
wrangler kv:namespace create "SESSIONS" --preview
```

Update `wrangler.toml` with the IDs:
```toml
[[kv_namespaces]]
binding = "CACHE"
id = "xxxx..."  # From CACHE create
preview_id = "yyyy..."  # From CACHE preview

[[kv_namespaces]]
binding = "SESSIONS"
id = "zzzz..."  # From SESSIONS create
preview_id = "aaaa..."  # From SESSIONS preview
```

### Step 7: Deploy to Cloudflare

```bash
# Development deployment
npm run dev
# Test locally at http://localhost:8787

# Production deployment
wrangler deploy --env production

# You'll get a URL like:
# https://security-platform-worker.YOUR-SUBDOMAIN.workers.dev
```

### Step 8: Test Your Deployment

```bash
# Test health endpoint
curl https://YOUR-WORKER-URL.workers.dev/api/v1/health

# Expected response:
# {"healthy":true,"checks":{"database":true,"kv":true}}

# Test status endpoint
curl https://YOUR-WORKER-URL.workers.dev/api/v1/status
```

### Step 9: Monitor Logs

```bash
# Tail logs in real-time
npm run tail

# Or
wrangler tail
```

### Cloudflare Custom Domain (Optional)

```bash
# In Cloudflare Dashboard:
# 1. Workers & Pages > YOUR_WORKER > Triggers
# 2. Add Custom Domain
# 3. Enter: api.yourdomain.com
# 4. Cloudflare automatically configures DNS
```

---

## Option 2: Oracle Cloud - OCI (Always Free)

**Best for**: Full control, maximum resources, production workloads  
**Cost**: $0/month (Always Free tier)  
**Setup time**: ~30 minutes

### Step 1: Install OCI CLI (Optional but Recommended)

```bash
# macOS/Linux
bash -c "$(curl -L https://raw.githubusercontent.com/oracle/oci-cli/master/scripts/install/install.sh)"

# Windows (PowerShell)
# Download from: https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm

# Verify installation
oci --version
```

### Step 2: Install Terraform

```bash
# macOS
brew install terraform

# Linux
wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
unzip terraform_1.6.0_linux_amd64.zip
sudo mv terraform /usr/local/bin/

# Windows
choco install terraform

# Verify
terraform --version
```

### Step 3: Generate OCI API Keys

```bash
# Create OCI config directory
mkdir -p ~/.oci

# Generate API key pair
openssl genrsa -out ~/.oci/oci_api_key.pem 2048
chmod 600 ~/.oci/oci_api_key.pem

# Generate public key
openssl rsa -pubout -in ~/.oci/oci_api_key.pem -out ~/.oci/oci_api_key_public.pem

# Display public key (you'll need this)
cat ~/.oci/oci_api_key_public.pem
```

### Step 4: Add API Key to OCI Console

1. Login to https://cloud.oracle.com
2. Click Profile Icon (top right) > **User Settings**
3. Under **Resources**, click **API Keys**
4. Click **Add API Key**
5. Select **Paste Public Key**
6. Paste contents of `~/.oci/oci_api_key_public.pem`
7. Click **Add**

**Important**: Copy the Configuration File Preview shown - you'll need these values!

### Step 5: Get Required OCIDs

Collect these from OCI Console:

```bash
# Tenancy OCID
# Profile > Tenancy > OCID (click "Copy")

# User OCID  
# Profile > User Settings > OCID

# Compartment OCID
# Identity & Security > Compartments > (root or create new) > OCID

# Fingerprint
# Shown when you added API key
```

### Step 6: Configure Terraform

```bash
cd infrastructure/cloud-configs/oci/terraform

# Copy example
cp terraform.tfvars.example terraform.tfvars

# Edit with your values
nano terraform.tfvars
```

Fill in your details:
```hcl
# OCI Authentication
tenancy_ocid     = "ocid1.tenancy.oc1..aaaaaa..."
user_ocid        = "ocid1.user.oc1..aaaaaa..."
fingerprint      = "xx:xx:xx:xx:xx:xx:xx:xx..."
private_key_path = "~/.oci/oci_api_key.pem"

# Region and Compartment
region          = "us-ashburn-1"  # or your preferred region
compartment_id  = "ocid1.compartment.oc1..aaaaaa..."

# Instance Configuration (Always Free - ARM)
instance_shape     = "VM.Standard.A1.Flex"
instance_ocpus     = 2
instance_memory_gb = 12

# Storage (max 200GB free)
data_volume_size_gb = 100

# SSH Key
ssh_public_key_path = "~/.ssh/id_rsa.pub"
```

### Step 7: Generate SSH Key (if you don't have one)

```bash
# Generate new SSH key
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa_oci

# Use this path in terraform.tfvars
ssh_public_key_path = "~/.ssh/id_rsa_oci.pub"
```

### Step 8: Initialize and Deploy with Terraform

```bash
# Initialize Terraform
terraform init

# Preview what will be created
terraform plan

# Review the output - you should see:
# - VCN (Virtual Cloud Network)
# - 2 Compute Instances
# - Block Volume
# - Object Storage Bucket
# - Security Lists

# Deploy (takes 5-10 minutes)
terraform apply

# Type 'yes' when prompted
```

### Step 9: Get Server IPs

```bash
# After deployment completes
terraform output

# You'll see:
# app_server_public_ip = "xxx.xxx.xxx.xxx"
# db_server_public_ip = "xxx.xxx.xxx.xxx"
# ssh_connection_app = "ssh ubuntu@xxx.xxx.xxx.xxx"
# ssh_connection_db = "ssh ubuntu@xxx.xxx.xxx.xxx"
```

### Step 10: Configure Servers

**A. Connect to DB Server**

```bash
# SSH to DB server
ssh ubuntu@DB_SERVER_IP

# Mount block volume (first time only)
sudo lsblk  # Find the device (usually /dev/sdb)
sudo mkfs.ext4 /dev/sdb
sudo mount /dev/sdb /mnt/data

# Add to fstab for auto-mount
echo "/dev/sdb /mnt/data ext4 defaults 0 0" | sudo tee -a /etc/fstab

# Create data directories
sudo mkdir -p /mnt/data/{postgres,redis,prometheus,grafana,backups}
sudo chown -R ubuntu:ubuntu /mnt/data
```

**B. Setup Database Environment**

```bash
# Still on DB server
cd /opt/database

# Create .env file
cat > .env <<EOF
DB_PASSWORD=your_strong_password
GRAFANA_PASSWORD=your_grafana_password
EOF

# Start database services
docker-compose up -d

# Verify
docker-compose ps
```

**C. Configure App Server**

```bash
# SSH to App server
ssh ubuntu@APP_SERVER_IP

cd /opt/security-platform

# Create .env file
cat > .env <<EOF
DB_HOST=DB_SERVER_PRIVATE_IP  # Use private IP, not public!
DB_PASSWORD=your_strong_password
REDIS_HOST=DB_SERVER_PRIVATE_IP
IBM_QUANTUM_TOKEN=your_token_if_you_have_one
EOF

# Get DB server private IP
# From OCI Console: Compute > Instances > db-server > Primary VNIC > Private IP
# Or from DB server: ip addr show | grep inet

# Start application services
docker-compose up -d

# Verify
docker-compose ps
```

### Step 11: Configure Firewall (if needed)

OCI Security Lists should already allow traffic, but verify:

```bash
# On both servers, check if needed
sudo ufw status

# If firewall is active, allow necessary ports
sudo ufw allow 22    # SSH
sudo ufw allow 80    # HTTP
sudo ufw allow 443   # HTTPS
sudo ufw allow 3000  # Grafana
sudo ufw allow 3001  # Backend API
sudo ufw allow 8000  # AI/Quantum API
```

### Step 12: Access Your Application

```bash
# Frontend
http://APP_SERVER_PUBLIC_IP

# Backend API
http://APP_SERVER_PUBLIC_IP:3001/api/v1/status

# AI/Quantum API
http://APP_SERVER_PUBLIC_IP:8000/api/v1/status

# Grafana
http://DB_SERVER_PUBLIC_IP:3000
# Default: admin / your_grafana_password

# Prometheus
http://DB_SERVER_PUBLIC_IP:9090
```

---

## Option 3: IBM Cloud (Lite Tier)

**Best for**: Quick prototypes, managed services, PaaS  
**Cost**: $0/month (Lite tier)  
**Setup time**: ~15 minutes

### Step 1: Install IBM Cloud CLI

```bash
# macOS
curl -fsSL https://clis.cloud.ibm.com/install/osx | sh

# Linux
curl -fsSL https://clis.cloud.ibm.com/install/linux | sh

# Windows (PowerShell as Administrator)
iex (New-Object Net.WebClient).DownloadString('https://clis.cloud.ibm.com/install/powershell')

# Verify
ibmcloud --version
```

### Step 2: Install Cloud Foundry Plugin

```bash
# Install CF plugin
ibmcloud cf install

# Verify
ibmcloud cf --version
```

### Step 3: Login to IBM Cloud

```bash
# Login
ibmcloud login

# Enter your IBM Cloud email and password

# Or with Single Sign-On
ibmcloud login --sso

# Select region
ibmcloud target -r us-south  # or your preferred region
```

### Step 4: Create an Organization and Space (if needed)

```bash
# Create organization (if you don't have one)
ibmcloud account org-create my-org

# Create space
ibmcloud account space-create dev -o my-org

# Target the org and space
ibmcloud target --cf
# Select your org and space from the list
```

### Step 5: Create Database Services

**PostgreSQL**:
```bash
# Create Databases for PostgreSQL (Lite/Standard plan)
ibmcloud resource service-instance-create security-platform-db \
  databases-for-postgresql standard us-south \
  -p '{
    "members_memory_allocation_mb": 1024,
    "members_disk_allocation_mb": 5120
  }'

# Wait for provisioning (can take 5-10 minutes)
ibmcloud resource service-instance security-platform-db

# Create credentials
ibmcloud resource service-key-create security-platform-db-key \
  Manager --instance-name security-platform-db

# View credentials (save these!)
ibmcloud resource service-key security-platform-db-key
```

**Redis** (optional):
```bash
# Create Databases for Redis
ibmcloud resource service-instance-create security-platform-redis \
  databases-for-redis standard us-south \
  -p '{
    "members_memory_allocation_mb": 256,
    "members_disk_allocation_mb": 1024
  }'

# Create credentials
ibmcloud resource service-key-create security-platform-redis-key \
  Manager --instance-name security-platform-redis
```

### Step 6: Configure Environment Variables

From the credentials you got above, note:
- PostgreSQL host, port, username, password, database name
- Redis host, port, password (if using)

Create `.env` files for each component:

**Backend** (`src/backend/.env`):
```bash
DB_HOST=xxxxx.databases.appdomain.cloud
DB_PORT=30000
DB_USER=ibm_cloud_xxxxx
DB_PASSWORD=xxxxx
DB_NAME=ibmclouddb

REDIS_HOST=xxxxx.databases.appdomain.cloud
REDIS_PORT=30001
REDIS_PASSWORD=xxxxx
```

**AI/Quantum** (`src/ai-quantum/.env`):
```bash
# Same as backend
DB_HOST=xxxxx.databases.appdomain.cloud
DB_PORT=30000
# ... etc
```

**Frontend** (`src/frontend/.env.production`):
```bash
NEXT_PUBLIC_API_URL=https://security-platform-api.mybluemix.net
NEXT_PUBLIC_WS_URL=wss://security-platform-api.mybluemix.net
```

### Step 7: Update Cloud Foundry Manifest

Edit `infrastructure/cloud-configs/ibm/manifest.yml`:

```yaml
applications:
  - name: security-platform-api
    memory: 256M
    instances: 1
    buildpacks:
      - go_buildpack
    path: ../../../src/backend
    # ... rest stays the same
    
  # Update routes with your preferred subdomain
    routes:
      - route: YOUR-APP-NAME-api.mybluemix.net
```

### Step 8: Deploy Applications

```bash
cd infrastructure/cloud-configs/ibm

# Deploy all applications
ibmcloud cf push

# This will deploy:
# - security-platform-api (Backend)
# - security-platform-ai (AI/Quantum)
# - security-platform-frontend (Frontend)

# Deployment takes 5-10 minutes
```

### Step 9: Bind Services (if not auto-bound)

```bash
# Bind database to backend
ibmcloud cf bind-service security-platform-api security-platform-db

# Bind database to AI service
ibmcloud cf bind-service security-platform-ai security-platform-db

# Restage applications
ibmcloud cf restage security-platform-api
ibmcloud cf restage security-platform-ai
```

### Step 10: Verify Deployment

```bash
# Check application status
ibmcloud cf apps

# You should see all apps as "started"

# Get URLs
ibmcloud cf app security-platform-frontend

# Access your applications
# Frontend: https://security-platform.mybluemix.net
# API: https://security-platform-api.mybluemix.net/api/v1/health
```

### Step 11: View Logs

```bash
# Tail logs
ibmcloud cf logs security-platform-api --recent

# Stream logs
ibmcloud cf logs security-platform-api

# In another terminal
ibmcloud cf logs security-platform-ai
```

---

## Hybrid Multi-Cloud Setup

**Recommended Architecture**: Cloudflare (edge) + OCI (compute) + IBM (database)

### Why Hybrid?

- **Cloudflare**: Global CDN, DDoS protection, API gateway
- **OCI**: Heavy compute (AI/ML), free forever
- **IBM**: Managed databases, Watson AI integration

### Setup Steps

**1. Deploy Backend to OCI**:
```bash
# Follow OCI steps above for full backend deployment
```

**2. Deploy Database to IBM Cloud**:
```bash
# Use IBM Cloud managed PostgreSQL (from Step 5 above)
# More reliable than self-hosting
```

**3. Deploy API Gateway to Cloudflare**:
```bash
# Update Cloudflare Worker to proxy to OCI
# Edit infrastructure/cloud-configs/cloudflare/src/api.js

# Add in handleApiRequest function:
async function handleApiRequest(request, env, ctx) {
  const url = new URL(request.url);
  
  // Proxy to OCI backend
  const backendUrl = `http://YOUR_OCI_IP:3001${url.pathname}`;
  
  return fetch(backendUrl, {
    method: request.method,
    headers: request.headers,
    body: request.body
  });
}
```

**4. Configure Services**:
```bash
# On OCI App Server, use IBM database
cat > /opt/security-platform/.env <<EOF
DB_HOST=xxxxx.databases.appdomain.cloud  # IBM Cloud
DB_PORT=30000
DB_USER=ibm_cloud_xxxxx
DB_PASSWORD=xxxxx
EOF
```

**Benefits**:
- **Cost**: Still $0/month!
- **Performance**: Best of all clouds
- **Reliability**: Distributed across providers
- **Scalability**: Easy to scale each component independently

---

## Verification

### Health Checks

**Cloudflare**:
```bash
curl https://YOUR-WORKER.workers.dev/api/v1/health
# Expected: {"healthy":true,...}
```

**OCI**:
```bash
curl http://YOUR_OCI_IP:3001/api/v1/health
curl http://YOUR_OCI_IP:8000/api/v1/status
```

**IBM Cloud**:
```bash
curl https://YOUR-APP-api.mybluemix.net/api/v1/health
```

### Database Connectivity

**Cloudflare D1**:
```bash
wrangler d1 execute security_platform_db \
  --command="SELECT COUNT(*) FROM threats"
```

**OCI PostgreSQL**:
```bash
ssh ubuntu@OCI_DB_IP
docker exec postgres psql -U sectools -d security -c "SELECT version()"
```

**IBM PostgreSQL**:
```bash
# Get connection string from credentials
psql "host=xxxxx.databases.appdomain.cloud port=30000 dbname=ibmclouddb user=ibm_cloud_xxxxx sslmode=require"
```

### End-to-End Test

```bash
# Test API endpoint
curl -X POST https://YOUR_DEPLOYMENT/api/v1/ml/detect \
  -H "Content-Type: application/json" \
  -d '{
    "source_ip": "192.168.1.100",
    "packets_per_second": 1000,
    "syn_count": 50
  }'

# Expected: Threat detection response
```

---

## Troubleshooting

### Cloudflare Issues

**Problem**: "Error: Could not find database"
```bash
# Solution: Verify database ID in wrangler.toml
wrangler d1 list
# Copy the correct database_id to wrangler.toml
```

**Problem**: "Rate limit exceeded"
```bash
# Solution: You hit free tier limits
# Check usage: https://dash.cloudflare.com > Workers > Analytics
# Implement caching to reduce requests
```

**Problem**: "Binding not found"
```bash
# Solution: KV namespace not configured
wrangler kv:namespace list
# Update wrangler.toml with correct namespace IDs
```

### OCI Issues

**Problem**: "Terraform apply fails with authorization error"
```bash
# Solution: Check API key configuration
# 1. Verify fingerprint matches in OCI Console
# 2. Ensure API key file path is correct
# 3. Check key permissions: chmod 600 ~/.oci/oci_api_key.pem
```

**Problem**: "Out of capacity for shape VM.Standard.A1.Flex"
```bash
# Solution: Try different availability domain or region
# Edit main.tf:
availability_domain = data.oci_identity_availability_domains.ads.availability_domains[1].name

# Or change region in terraform.tfvars:
region = "us-phoenix-1"  # Try different region
```

**Problem**: "Cannot SSH to instances"
```bash
# Solution: Check security list
# 1. OCI Console > Networking > VCNs > Your VCN > Security Lists
# 2. Ensure Ingress Rule for port 22 exists
# 3. Check SSH key: ssh -i ~/.ssh/id_rsa_oci ubuntu@IP
```

**Problem**: "Services not starting on OCI"
```bash
# SSH to server
ssh ubuntu@OCI_IP

# Check cloud-init status
cloud-init status

# View cloud-init logs
sudo cat /var/log/cloud-init-output.log

# Manually install Docker if needed
sudo apt-get update
sudo apt-get install -y docker.io docker-compose
```

### IBM Cloud Issues

**Problem**: "CF push fails with insufficient memory"
```bash
# Solution: Reduce memory in manifest.yml
# Or scale up (paid):
ibmcloud cf scale security-platform-api -m 512M
```

**Problem**: "Database connection timeout"
```bash
# Solution: Check credentials
ibmcloud resource service-key security-platform-db-key

# Verify SSL mode is required
# Add to connection string: ?sslmode=require
```

**Problem**: "Application crashes on startup"
```bash
# Check logs
ibmcloud cf logs security-platform-api --recent

# Common issues:
# 1. Missing environment variables
# 2. Database not bound
# 3. Wrong buildpack

# Rebind database
ibmcloud cf bind-service security-platform-api security-platform-db
ibmcloud cf restage security-platform-api
```

### General Issues

**Problem**: "CORS errors in browser"
```bash
# Solution: Check backend CORS configuration
# Ensure frontend URL is in allowed origins
```

**Problem**: "WebSocket connection failed"
```bash
# Cloudflare: Durable Objects must be configured
# OCI: Check firewall allows WebSocket upgrade
# IBM: Cloud Foundry should support WebSockets by default
```

**Problem**: "High latency"
```bash
# Cloudflare: Check edge location (should be close to user)
# OCI: Use correct region close to users
# IBM: Consider using CDN in front
# Hybrid: Route through Cloudflare to nearest OCI/IBM backend
```

---

## Next Steps

1. **Configure CI/CD**:
   - See `cicd/buddy/README.md` for Buddy CI
   - See `cicd/argocd/README.md` for Argo CD
   - See `cicd/harness/README.md` for Harness

2. **Add Monitoring**:
   ```bash
   # OCI: Access Grafana
   http://OCI_DB_IP:3000
   # Default: admin / your_grafana_password
   ```

3. **Setup Backups**:
   ```bash
   # OCI: Daily backups configured via cloud-init
   # IBM: Use IBM Cloud Backup service
   # Cloudflare: D1 automatic backups
   ```

4. **Custom Domain**:
   - Cloudflare: Workers > Custom Domains
   - OCI: Configure DNS A record to public IP
   - IBM: Cloud Foundry > Custom Domains

5. **Security Hardening**:
   - Enable HTTPS (all platforms)
   - Configure firewall rules
   - Rotate secrets regularly
   - Enable MFA on cloud accounts

---

## Cost Monitoring

### Cloudflare
```bash
# Check usage
wrangler dash
# Or: https://dash.cloudflare.com > Workers > Analytics
```

### OCI
```bash
# Always Free resources have no cost
# But monitor: https://cloud.oracle.com > Billing > Cost Analysis
```

### IBM Cloud
```bash
# Check usage
ibmcloud billing account-usage

# Or: https://cloud.ibm.com > Manage > Billing and usage
```

---

## Support

- **Documentation**: Full docs in `docs/` directory
- **Issues**: GitHub Issues
- **Community**: GitHub Discussions
- **Email**: support@example.com

---

**Congratulations!** 🎉 You've deployed the Unified Security Platform to the cloud!

For detailed configuration and advanced topics, see:
- [Cloudflare Guide](infrastructure/cloud-configs/cloudflare/README.md)
- [OCI Guide](infrastructure/cloud-configs/oci/README.md)
- [IBM Cloud Guide](infrastructure/cloud-configs/ibm/README.md)
- [Cost Comparison](docs/deployment/cost-comparison.md)

