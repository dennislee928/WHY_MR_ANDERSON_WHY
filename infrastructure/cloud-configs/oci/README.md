# Oracle Cloud Infrastructure (OCI) Deployment

Deploy the Security Platform to OCI using the Always Free tier.

## Always Free Tier Resources

Oracle Cloud offers generous Always Free tier resources:

### Compute
- **2x AMD VMs**: VM.Standard.E2.1.Micro (1/8 OCPU, 1 GB RAM each)
- **4x ARM vCPUs**: VM.Standard.A1.Flex (Ampere A1, up to 24 GB RAM total)

### Storage
- **200 GB** Block Volume storage
- **10 GB** Object Storage (50,000 API requests/month)

### Networking
- Flexible Network Load Balancer
- Outbound Data Transfer (10 TB/month)

**Recommended**: Use ARM-based VMs for better performance within free tier.

## Prerequisites

1. OCI Account (sign up at https://www.oracle.com/cloud/free/)
2. Terraform installed:
   ```bash
   # macOS
   brew install terraform

   # Linux
   wget https://releases.hashicorp.com/terraform/1.6.0/terraform_1.6.0_linux_amd64.zip
   unzip terraform_1.6.0_linux_amd64.zip
   sudo mv terraform /usr/local/bin/
   ```
3. OCI CLI (optional but recommended):
   ```bash
   bash -c "$(curl -L https://raw.githubusercontent.com/oracle/oci-cli/master/scripts/install/install.sh)"
   ```

## Setup

### 1. Configure OCI API Access

```bash
# Create OCI config directory
mkdir -p ~/.oci

# Generate API key pair
openssl genrsa -out ~/.oci/oci_api_key.pem 2048
chmod 600 ~/.oci/oci_api_key.pem
openssl rsa -pubout -in ~/.oci/oci_api_key.pem -out ~/.oci/oci_api_key_public.pem
```

### 2. Add API Key to OCI Console

1. Log in to OCI Console
2. Go to Profile > API Keys
3. Click "Add API Key"
4. Upload `~/.oci/oci_api_key_public.pem`
5. Note the configuration details shown

### 3. Get Required OCIDs

You'll need:
- **Tenancy OCID**: Profile > Tenancy > OCID
- **User OCID**: Profile > User Settings > OCID
- **Compartment OCID**: Identity > Compartments > (select compartment) > OCID
- **API Key Fingerprint**: Shown when you added the API key

### 4. Configure Terraform

```bash
cd infrastructure/cloud-configs/oci/terraform

# Copy example variables
cp terraform.tfvars.example terraform.tfvars

# Edit with your values
vim terraform.tfvars
```

Fill in:
```hcl
tenancy_ocid     = "ocid1.tenancy.oc1..aaaaaa..."
user_ocid        = "ocid1.user.oc1..aaaaaa..."
fingerprint      = "xx:xx:xx:xx..."
compartment_id   = "ocid1.compartment.oc1..aaaaaa..."
region           = "us-ashburn-1"
```

### 5. Generate SSH Key

```bash
# Generate SSH key pair if you don't have one
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa_oci

# Update terraform.tfvars
ssh_public_key_path = "~/.ssh/id_rsa_oci.pub"
```

## Deployment

### Initialize Terraform

```bash
cd infrastructure/cloud-configs/oci/terraform
terraform init
```

### Plan Deployment

```bash
terraform plan
```

Review the resources to be created:
- VCN (Virtual Cloud Network)
- Internet Gateway
- Security Lists
- 2x Compute Instances (ARM-based)
- Block Volume (100 GB)
- Object Storage Bucket

### Deploy

```bash
terraform apply
```

Type `yes` when prompted. Deployment takes ~5-10 minutes.

### Get Connection Details

```bash
terraform output
```

You'll see:
```
app_server_public_ip = "xxx.xxx.xxx.xxx"
db_server_public_ip = "xxx.xxx.xxx.xxx"
ssh_connection_app = "ssh ubuntu@xxx.xxx.xxx.xxx"
ssh_connection_db = "ssh ubuntu@xxx.xxx.xxx.xxx"
```

## Post-Deployment Configuration

### 1. Mount Block Volume

SSH to DB server:
```bash
ssh ubuntu@<db_server_ip>

# Find block device
lsblk

# Format and mount (only first time)
sudo mkfs.ext4 /dev/sdb
sudo mount /dev/sdb /mnt/data

# Add to fstab for auto-mount
echo "/dev/sdb /mnt/data ext4 defaults 0 0" | sudo tee -a /etc/fstab

# Create data directories
sudo mkdir -p /mnt/data/{postgres,redis,prometheus,grafana,backups}
sudo chown -R ubuntu:ubuntu /mnt/data
```

### 2. Configure Environment Variables

On DB server:
```bash
cd /opt/database

# Create .env file
cat > .env <<EOF
DB_PASSWORD=<strong-password>
GRAFANA_PASSWORD=<grafana-password>
EOF

# Update prometheus.yml with app server IP
APP_SERVER_IP=<app-server-ip>
sed -i "s/APP_SERVER_IP/$APP_SERVER_IP/g" prometheus.yml
```

On App server:
```bash
ssh ubuntu@<app_server_ip>
cd /opt/security-platform

# Create .env file
cat > .env <<EOF
DB_HOST=<db-server-internal-ip>
DB_PASSWORD=<same-as-db-server>
REDIS_HOST=<db-server-internal-ip>
IBM_QUANTUM_TOKEN=<your-token>
EOF

# Update docker-compose override
DB_SERVER_IP=<db-server-internal-ip>
sed -i "s/\${DB_SERVER_IP}/$DB_SERVER_IP/g" docker-compose.override.yml
```

### 3. Start Services

On DB server:
```bash
cd /opt/database
docker-compose up -d

# Verify
docker-compose ps
```

On App server:
```bash
cd /opt/security-platform
docker-compose up -d

# Verify
docker-compose ps
```

## Access Services

Once deployed, access services at:

- **Frontend**: http://<app_server_ip>
- **Backend API**: http://<app_server_ip>:3001
- **AI/Quantum API**: http://<app_server_ip>:8000
- **Grafana**: http://<db_server_ip>:3000
- **Prometheus**: http://<db_server_ip>:9090

## Monitoring

### Check Instance Status

```bash
# SSH to instances
ssh ubuntu@<ip>

# Check running containers
docker ps

# View logs
docker-compose logs -f

# System resources
htop
df -h
```

### OCI Console

Monitor in OCI Console:
- Compute > Instances > (select instance)
- View CPU, Memory, Network metrics
- Access console connection

## Backups

### Automatic Database Backups

Backups run daily at 2 AM (configured in cloud-init):
```bash
# Manual backup
ssh ubuntu@<db_server_ip>
sudo /usr/local/bin/backup-databases.sh

# View backups
ls -lh /mnt/data/backups/
```

### Object Storage Sync

Upload backups to Object Storage:
```bash
# Install OCI CLI on DB server
ssh ubuntu@<db_server_ip>
bash -c "$(curl -L https://raw.githubusercontent.com/oracle/oci-cli/master/scripts/install/install.sh)"

# Configure OCI CLI
oci setup config

# Sync backups
oci os object bulk-upload \
  --bucket-name security-platform-backups-production \
  --src-dir /mnt/data/backups/
```

### Block Volume Backups

Create backups via OCI Console or CLI:
```bash
oci bv backup create \
  --volume-id <volume-ocid> \
  --display-name "security-platform-data-$(date +%Y%m%d)"
```

## Scaling

### Vertical Scaling (Within Free Tier)

Modify instance resources:
```hcl
# terraform.tfvars
instance_ocpus     = 4  # Max for free tier
instance_memory_gb = 24 # Max for free tier total
```

Apply changes:
```bash
terraform apply
```

### Horizontal Scaling (Requires Paid)

Add more instances by duplicating resources in `main.tf`.

## Cost Management

### Stay Within Free Tier

✅ **Always Free Resources**:
- 2x AMD VMs or 4x ARM vCPUs: **$0**
- 200 GB Block Volume: **$0**
- 10 GB Object Storage: **$0**
- 10 TB outbound transfer: **$0**

⚠️ **Monitor Usage**:
- Check OCI Console > Billing > Cost Analysis
- Set up budget alerts

### Optimize Costs

1. **Use ARM instances**: Better performance in free tier
2. **Clean up old backups**: Keep only 7 days locally
3. **Archive to Object Storage**: Cheaper long-term storage
4. **Monitor block volume**: Stay under 200 GB

## Troubleshooting

### Instance Won't Start

```bash
# Check OCI console for errors
# Common issues:
# - Out of capacity (try different availability domain)
# - Service limits reached
# - Insufficient free tier resources

# Via OCI CLI
oci compute instance list \
  --compartment-id <compartment-ocid> \
  --lifecycle-state RUNNING
```

### Can't SSH

```bash
# Check security list allows port 22
# Verify public IP assigned
terraform output app_server_public_ip

# Use console connection
# OCI Console > Compute > Instance > Console Connection
```

### Services Not Starting

```bash
ssh ubuntu@<ip>

# Check cloud-init status
cloud-init status

# View cloud-init logs
sudo cat /var/log/cloud-init-output.log

# Manually run setup
sudo docker-compose up -d
```

## Cleanup

### Destroy Resources

```bash
cd infrastructure/cloud-configs/oci/terraform
terraform destroy
```

Type `yes` to confirm. This will:
- Terminate all instances
- Delete block volumes
- Remove VCN and networking
- Delete object storage bucket (if empty)

⚠️ **Warning**: This is irreversible. Backup data first!

### Keep Free Tier Resources

To keep some resources:
```bash
# Remove specific resources
terraform state rm oci_core_instance.app_server

# Then destroy the rest
terraform destroy
```

## Security Best Practices

1. **Change Default Passwords**: Update all default credentials
2. **Enable Firewall**: Configure security lists properly
3. **Use Private IPs**: Communicate between instances via private network
4. **Regular Updates**: Keep software updated
5. **Monitor Logs**: Review security logs regularly
6. **Enable MFA**: Use multi-factor authentication for OCI account

## Performance Tuning

### ARM Instances

Optimize for ARM architecture:
```dockerfile
# Use ARM-compatible images
FROM --platform=linux/arm64 ubuntu:22.04
```

### Docker Optimization

```bash
# Prune unused resources
docker system prune -a

# Limit container resources
# Add to docker-compose.yml:
deploy:
  resources:
    limits:
      cpus: '1.0'
      memory: 4G
```

## Further Reading

- [OCI Always Free Tier](https://www.oracle.com/cloud/free/)
- [OCI Terraform Provider](https://registry.terraform.io/providers/oracle/oci/latest/docs)
- [OCI Documentation](https://docs.oracle.com/en-us/iaas/Content/home.htm)
- [OCI CLI Reference](https://docs.oracle.com/en-us/iaas/tools/oci-cli/latest/oci_cli_docs/)

