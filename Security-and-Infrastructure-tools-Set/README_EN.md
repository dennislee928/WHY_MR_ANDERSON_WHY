# Security & Infrastructure Tools Set

<div align="center">

![Version](https://img.shields.io/badge/version-1.0.0-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)
![Docker](https://img.shields.io/badge/Docker-20.10+-blue.svg)
![Docker Compose](https://img.shields.io/badge/Docker%20Compose-2.0+-blue.svg)

**A Docker-based Open Source Security Scanning & Infrastructure Management Platform**

[繁體中文文檔](./README.md) | [Architecture](./ARCHITECTURE.md) | [Tools Reference](./TOOLS.md)

</div>

---

## 📋 Table of Contents

- [Introduction](#introduction)
- [Core Features](#core-features)
- [Quick Start](#quick-start)
- [System Architecture](#system-architecture)
- [Prerequisites](#prerequisites)
- [Installation & Configuration](#installation--configuration)
- [Usage Guide](#usage-guide)
- [Integrated Tools](#integrated-tools)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
- [Development Guide](#development-guide)
- [FAQ](#faq)
- [Contributing](#contributing)
- [License](#license)

---

## Introduction

**Security & Infrastructure Tools Set** is a comprehensive containerized security scanning platform that integrates industry-leading open-source security tools with a unified deployment, management, and query interface. This project follows Docker best practices and adopts a microservices architecture, suitable for personal learning, team usage, and production deployment.

### What Problems Does It Solve?

- ✅ **Environment Consistency**: Eliminate "works on my machine" issues
- ✅ **Rapid Deployment**: Launch a complete security scanning platform with one command
- ✅ **Tool Integration**: Unified management of multiple security scanning tools
- ✅ **Result Aggregation**: Centralized database for storing and querying scan results
- ✅ **Scalability**: Easily add new scanners or services
- ✅ **Best Practices**: Built-in security configuration, health checks, resource limits

### Use Cases

| Scenario | Description |
|----------|-------------|
| 🎓 **Security Learning** | Experience industry-standard tools and understand security scanning workflows |
| 👨‍💻 **Personal Use** | Quickly set up local security testing environment |
| 👥 **Team Collaboration** | Unified scanning platform with shared results and tracking |
| 🏢 **Enterprise Deployment** | Scalable security scanning infrastructure for production |
| 🔬 **Security Research** | Rapidly validate vulnerabilities and test POCs |

---

## Core Features

### 🎯 Technical Features

- **🐳 Container-First**: All services run in Docker containers, one-click deployment
- **🔧 Microservices Architecture**: Each tool runs independently without interference
- **💾 Centralized Storage**: PostgreSQL manages scan results and metadata uniformly
- **🔐 Secret Management**: HashiCorp Vault centralizes sensitive data management
- **🌐 Reverse Proxy**: Traefik with automatic HTTPS and service discovery
- **📊 GitOps Support**: ArgoCD for declarative deployment
- **🏥 Health Checks**: Automatic service status monitoring and dependency management
- **📈 Observability**: Complete logging and metrics collection ready

### 🛡️ Security Features

- **Multi-layer Defense**: Network isolation, authentication, least privilege
- **Secret Rotation**: Automated secret lifecycle management
- **Audit Logs**: Complete record of all operations and access
- **Resource Limits**: Prevention of resource exhaustion attacks
- **Security Updates**: Fixed version tags with controlled updates

### 🚀 Integrated Tools

| Category | Tool | Status |
|----------|------|--------|
| Vulnerability Scan | Nuclei | ✅ Integrated |
| Network Scan | Nmap | ✅ Integrated |
| Asset Discovery | AMASS | ✅ Integrated |
| Database | PostgreSQL | ✅ Integrated |
| Secret Management | Vault | ✅ Integrated |
| Reverse Proxy | Traefik | ✅ Integrated |
| CI/CD | ArgoCD | ✅ Integrated |
| Orchestration | SecureCodeBox | ✅ Integrated |

For more optional tools, see [TOOLS.md](./TOOLS.md)

---

## Quick Start

### ⚡ 3-Minute Quick Deployment

```bash
# 1. Clone the project
git clone https://github.com/your-username/Security-and-Infrastructure-tools-Set.git
cd Security-and-Infrastructure-tools-Set

# 2. Configure environment variables (optional, skip to use defaults)
cp .env.template .env
# Edit .env to modify sensitive information like database password

# 3. Start all services
make up

# 4. Check service status
make health

# 5. Run your first scan
make scan-nuclei TARGET=https://example.com
```

### 🎉 Success!

Access the following services:
- **Traefik Dashboard**: <http://localhost:8080>
- **Vault UI**: <http://localhost:8200>
- **ArgoCD UI**: <http://localhost:8081>
- **Web UI**: <http://localhost:8082>

---

## System Architecture

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         External Users                           │
│              (Developers, Security Teams, Automation)            │
└────────────────────────────┬────────────────────────────────────┘
                             │
                    ┌────────▼────────┐
                    │    Traefik      │ ◄── 🌐 Reverse Proxy & SSL
                    │  (Port 80/443)  │     Load Balancing, Service Discovery
                    └────────┬────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
   ┌────▼─────┐      ┌──────▼──────┐      ┌─────▼──────┐
   │  Vault   │      │   ArgoCD    │      │  Web UI    │
   │  :8200   │      │   :8081     │      │  :8082     │
   └────┬─────┘      └──────┬──────┘      └─────┬──────┘
        │ 🔐 Secrets        │ 🚀 GitOps         │ 📊 Query Interface
        │                   │                    │
        └───────────────────┼────────────────────┘
                            │
                    ┌───────▼────────┐
                    │   PostgreSQL   │ ◄── 💾 Central Database
                    │     :5432      │     Scan Results, Metadata
                    └───────┬────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
   ┌────▼─────┐      ┌─────▼──────┐      ┌────▼──────┐
   │ Scanner  │      │  Operator  │      │  Parsers  │
   │ Nuclei   │      │ SecCodeBox │      │ N/A/N     │
   │ Nmap     │      │            │      │           │
   │ AMASS    │      │ 🔧 Orchestration │  │ 📋 Parsers │
   └──────────┘      └────────────┘      └───────────┘
        │                   │                   │
        └───────────────────┼───────────────────┘
                            │
                    ┌───────▼────────┐
                    │  Scan Results  │ ◄── 📁 Shared Storage Volume
                    │    Volume      │     Scan Output Files
                    └────────────────┘
```

For detailed architecture design, see [ARCHITECTURE.md](./ARCHITECTURE.md)

---

## Prerequisites

### Hardware Requirements

| Environment | CPU | Memory | Disk Space |
|-------------|-----|--------|-----------|
| Minimum | 2 cores | 4GB | 20GB |
| Recommended | 4 cores | 8GB | 50GB |
| Production | 8+ cores | 16GB+ | 100GB+ |

### Software Requirements

- **OS**: Linux, macOS, Windows (WSL2)
- **Docker**: 20.10 or higher
- **Docker Compose**: 2.0 or higher
- **Git**: For cloning the project
- **Make**: For running commands (Windows requires additional installation)

### Environment Check

```bash
# Check Docker version
docker --version
# Output: Docker version 20.10.x

# Check Docker Compose version
docker-compose --version
# Output: Docker Compose version 2.x.x

# Check if Docker is running
docker ps
# Should display container list (even if empty)

# Check Make
make --version
# Windows users can use Git Bash or install Make for Windows
```

---

## Installation & Configuration

### Step 1: Clone the Project

```bash
git clone https://github.com/your-username/Security-and-Infrastructure-tools-Set.git
cd Security-and-Infrastructure-tools-Set
```

### Step 2: Configure Environment Variables

The project provides an environment variable template. Copy and modify:

```bash
# Copy template
cp .env.template .env

# Edit environment variables
nano .env  # or use your preferred editor
```

**Key Configuration Items**:

```bash
# 🔴 Must modify (Production)
DB_PASSWORD=<strong_password>      # Database password
VAULT_TOKEN=<vault_root_token>    # Vault root token

# 🟡 Recommended to modify
SCAN_CONCURRENCY=10               # Scan concurrency
NUCLEI_RATE_LIMIT=150            # Rate limit
NMAP_TIMING=T4                    # Nmap scan speed

# 🟢 Optional configuration
TZ=Asia/Taipei                    # Timezone
DEBUG=false                       # Debug mode
```

For complete environment variable descriptions, see the README.md

### Step 3: Initial Deployment

```bash
# Use Makefile for one-click startup
cd Make_Files
make up

# Or use Docker Compose directly
cd Docker/compose
docker-compose up -d

# Wait for services to start (about 30 seconds)
sleep 30

# Check service status
make health
# or
docker-compose ps
```

---

## Usage Guide

### Makefile Command Reference

The project provides convenient Makefile commands:

```bash
# Enter Make_Files directory
cd Make_Files

# Display all available commands
make help
```

#### Service Management

```bash
# Start all services
make up

# Stop all services
make down

# Restart all services
make restart

# View service status
make ps

# View health status
make health

# View real-time logs
make logs
```

#### Scan Operations

##### Nuclei Scan

```bash
# Scan single target
make scan-nuclei TARGET=https://example.com

# Scan with specific templates
docker-compose run --rm scanner-nuclei \
    nuclei -u https://example.com -t /templates/cves/ -o /results/cve-scan.json

# Specify severity
docker-compose run --rm scanner-nuclei \
    nuclei -u https://example.com -severity critical,high -o /results/critical.json
```

##### Nmap Scan

```bash
# Basic scan
make scan-nmap TARGET=192.168.1.1

# Scan entire subnet
make scan-nmap TARGET=192.168.1.0/24

# Service version detection
docker-compose run --rm nmap nmap -sV 192.168.1.1 -oX /results/nmap-version.xml
```

##### AMASS Scan

```bash
# Subdomain enumeration
docker-compose run --rm scanner-amass amass enum -d example.com -o /results/amass-subs.txt

# Passive mode (no direct target probing)
docker-compose run --rm scanner-amass amass enum -passive -d example.com -o /results/amass-passive.txt
```

#### Database Operations

```bash
# Backup database
make backup

# Manual backup
docker-compose exec -T postgres pg_dump -U sectools security > backup-$(date +%Y%m%d).sql

# Restore database
docker-compose exec -T postgres psql -U sectools security < backup-20251017.sql

# Enter PostgreSQL CLI
docker exec -it postgres psql -U sectools -d security
```

---

## Integrated Tools

### Detailed Tool Introduction

#### 🔍 Nuclei - Fast Vulnerability Scanner

**Overview**: Template-based vulnerability scanner developed by ProjectDiscovery

**Features**:
- 🚀 Extremely fast (written in Go)
- 📝 YAML templates, easy to customize
- 🔄 Community-driven, rapid template updates
- 📊 Low false positive rate

**Usage Example**:
```bash
# Basic scan
make scan-nuclei TARGET=https://example.com

# Custom template directory
docker-compose run --rm -v ./custom-templates:/custom scanner-nuclei \
    nuclei -u https://example.com -t /custom -o /results/custom.json
```

#### 🌐 Nmap - The King of Network Scanning

**Overview**: Classic network exploration and security auditing tool

**Features**:
- 🎯 Precise port scanning
- 🔬 Service version detection
- 🖥️ OS fingerprinting
- 📜 NSE script engine extension

**Scan Types**:
```bash
# TCP SYN scan (default, requires root)
docker-compose run --rm nmap nmap -sS 192.168.1.1

# TCP Connect scan (no root required)
docker-compose run --rm nmap nmap -sT 192.168.1.1

# UDP scan
docker-compose run --rm nmap nmap -sU 192.168.1.1
```

#### 🗺️ AMASS - Asset Discovery Expert

**Overview**: OWASP project for deep subdomain discovery and external attack surface management

**Features**:
- 🔍 Multi-source integration
- 🤫 Passive/active modes
- 🌐 DNS enumeration
- 📊 Relationship graph visualization

---

## Best Practices

### Security Configuration Recommendations

#### 1. Secret Management

**❌ Don't do this**:
```yaml
environment:
  DB_PASSWORD: "plaintext_password"  # Plain text password
```

**✅ Should do this**:
```yaml
environment:
  DB_PASSWORD: ${DB_PASSWORD}  # Read from environment variable or .env
```

**🔒 Best Practice**:
```yaml
# Use Docker secrets
secrets:
  db_password:
    file: ./secrets/db_password.txt

services:
  postgres:
    secrets:
      - db_password
    environment:
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
```

#### 2. Network Isolation

```yaml
# Create multiple networks to isolate different tiers
networks:
  frontend:  # User-facing services
  backend:   # Internal services
  database:  # Database-only

services:
  traefik:
    networks:
      - frontend
  
  api:
    networks:
      - frontend
      - backend
  
  postgres:
    networks:
      - backend
      - database
```

#### 3. Resource Limits

```yaml
services:
  scanner-nuclei:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
```

### Performance Tuning

#### 1. PostgreSQL Optimization

```sql
-- Adjust shared buffers (25% of container memory)
ALTER SYSTEM SET shared_buffers = '512MB';

-- Adjust work memory
ALTER SYSTEM SET work_mem = '16MB';

-- Enable parallel queries
ALTER SYSTEM SET max_parallel_workers_per_gather = 2;

-- Reload configuration
SELECT pg_reload_conf();
```

#### 2. Nuclei Tuning

```bash
# Adjust concurrency and rate
nuclei -u https://example.com \
  -c 50 \                    # Concurrency
  -rate-limit 150 \          # Requests per second
  -timeout 5 \               # Timeout
  -retries 1 \               # Retries
  -bulk-size 25              # Bulk size
```

---

## Troubleshooting

### Common Issues

#### Issue 1: Service Fails to Start

**Symptoms**: Service status shows `Exited` after `docker-compose up -d`

**Diagnosis**:
```bash
# View service logs
docker-compose logs service-name

# Check container status
docker-compose ps

# Check resource usage
docker stats
```

**Possible Causes**:
- Port conflict: Modify port mapping in docker-compose.yml
- Insufficient memory: Increase Docker memory limit or reduce services
- Configuration error: Check environment variables and config files

#### Issue 2: PostgreSQL Health Check Fails

**Symptoms**: `postgres` service status shows `unhealthy`

**Solution**:
```bash
# 1. View PostgreSQL logs
docker-compose logs postgres

# 2. Manually test health check command
docker exec -it postgres pg_isready -U sectools

# 3. Check database connectivity
docker exec -it postgres psql -U sectools -d security -c "SELECT 1;"

# 4. Restart PostgreSQL
docker-compose restart postgres
```

---

## Development Guide

### Adding New Scan Tools

#### Step 1: Select Tool

Refer to [TOOLS.md](./TOOLS.md) to choose a tool to integrate, e.g., `trivy`

#### Step 2: Add to docker-compose.yml

```yaml
services:
  scanner-trivy:
    image: aquasec/trivy:latest
    volumes:
      - scan_results:/results
      - trivy_cache:/root/.cache/trivy
    networks:
      - security_net
    command: ["--help"]  # Default command
```

#### Step 3: Add Makefile Command

```makefile
# Make_Files/Makefile
scan-trivy:
	docker-compose run --rm scanner-trivy \
		trivy image --format json --output /results/trivy-$(shell date +%Y%m%d-%H%M%S).json $(TARGET)
```

---

## FAQ

### General Questions

**Q: Can it be used in production?**
A: Yes, but the following adjustments are needed:
- Change all default passwords
- Enable Traefik SSL
- Configure firewall rules
- Set up monitoring and alerting
- Regular database backups

**Q: Which operating systems are supported?**
A: Any system supporting Docker:
- Linux (recommended)
- macOS
- Windows 10/11 with WSL2

**Q: How much resources are needed?**
A:
- Minimum: 2-core CPU, 4GB RAM
- Recommended: 4-core CPU, 8GB RAM
- Production: 8+ core CPU, 16+ GB RAM

---

## Contributing

We welcome contributions of all kinds!

### How to Contribute

1. **Fork the project**
2. **Create a feature branch** (`git checkout -b feature/AmazingFeature`)
3. **Commit changes** (`git commit -m 'Add some AmazingFeature'`)
4. **Push the branch** (`git push origin feature/AmazingFeature`)
5. **Open a Pull Request**

For detailed guidelines, see [CONTRIBUTING.md](./CONTRIBUTING.md)

---

## Roadmap

### v1.1 (2025 Q2)

- [ ] Web UI Dashboard
- [ ] Complete REST API
- [ ] Redis Task Queue
- [ ] Automated Report Generation

### v1.2 (2025 Q3)

- [ ] Prometheus + Grafana Monitoring
- [ ] ELK Stack Log Aggregation
- [ ] N8N Workflow Automation
- [ ] Trivy Container Scanning Integration

### v2.0 (2025 Q4)

- [ ] AI-Assisted Analysis (Ollama + ChromaDB)
- [ ] Kubernetes Helm Charts
- [ ] Multi-tenancy Support
- [ ] MISP Threat Intelligence Integration

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

---

## Acknowledgments

Thanks to the following open-source projects:

- [ProjectDiscovery Nuclei](https://github.com/projectdiscovery/nuclei)
- [Nmap](https://nmap.org/)
- [OWASP AMASS](https://github.com/OWASP/Amass)
- [HashiCorp Vault](https://www.vaultproject.io/)
- [Traefik](https://traefik.io/)
- [SecureCodeBox](https://www.securecodebox.io/)

---

## Contact

- **Project Issues**: [GitHub Issues](https://github.com/your-username/Security-and-Infrastructure-tools-Set/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-username/Security-and-Infrastructure-tools-Set/discussions)
- **Email**: security-tools@example.com

---

## Disclaimer

This tool is for legal security testing and research only. Users must comply with local laws and regulations and obtain explicit authorization from target system owners. Project maintainers are not responsible for any misuse.

**Please use this tool responsibly!**

---

<div align="center">

**If this project helps you, please give it a ⭐ Star!**

Made with ❤️ by Security Community

</div>

