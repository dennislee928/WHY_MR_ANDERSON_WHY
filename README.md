# Serverless E2E Tester (Robot Framework)

A ready-to-use GitHub Action for running Robot Framework API tests in CI, without manually installing Python, Java, or any dependencies.

[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub Marketplace](https://img.shields.io/badge/Marketplace-Serverless%20E2E%20Tester-blue.svg)](https://github.com/marketplace/actions/serverless-e2e-tester-robot-framework-by-dennis)

## Core Value

**Remove the pain of setting up a serverless test environment.** Robot Framework works well, but installing Java, Python, and dependencies in CI is tedious. This Action lets you focus on writing tests instead of configuring the environment.

## Features

- **Ready to use** - Includes Robot Framework, Requests Library, and JSON Library
- **Environment injection** - Automatically injects GitHub Secrets as test variables
- **Report output** - Automatically generates HTML test reports
- **Tag filtering** - Supports `--include` and `--exclude` tags
- **Variable files** - Supports custom variable files
- **Parallel execution** - Supports multi-process parallel tests
- **Zero configuration** - No dependencies to install

## Usage Examples

### Basic usage

```yaml
name: API E2E Tests

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Run Robot Framework Tests
        uses: dennislee928/WHY_MR_ANDERSON_WHY@v0.0.2
        with:
          test_dir: 'QAQC'
          target_url: 'https://your-api.workers.dev'
      
      - name: Upload Test Reports
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: robot-framework-reports
          path: reports/
          retention-days: 7
```

### Using GitHub Secrets

```yaml
- name: Run Robot Framework Tests
  uses: dennislee928/WHY_MR_ANDERSON_WHY@v0.0.2
  env:
    API_KEY: ${{ secrets.API_KEY }}
    AUTH_TOKEN: ${{ secrets.AUTH_TOKEN }}
  with:
    test_dir: 'QAQC'
    target_url: 'https://your-api.workers.dev'
```

These secrets are automatically injected as variables in Robot Framework tests:

```robot
*** Test Cases ***
Test API with Authentication
    ${headers}=    Create Dictionary
    ...    Authorization=Bearer ${AUTH_TOKEN}
    ...    X-API-Key=${API_KEY}
    ${response}=    GET    ${BASE_URL}/api/v1/endpoint    headers=${headers}
    Should Be Equal As Strings    ${response.status_code}    200
```

### Filtering by tags

```yaml
- name: Run Smoke Tests Only
  uses: dennislee928/WHY_MR_ANDERSON_WHY@v0.0.2
  with:
    test_dir: 'QAQC'
    target_url: 'https://your-api.workers.dev'
    include_tags: 'smoke'
    exclude_tags: 'slow,integration'
```

### Parallel execution

```yaml
- name: Run Tests in Parallel
  uses: dennislee928/WHY_MR_ANDERSON_WHY@v0.0.2
  with:
    test_dir: 'QAQC'
    target_url: 'https://your-api.workers.dev'
    processes: '4'
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `test_dir` | Directory containing `.robot` test files | Yes | - |
| `target_url` | Deployed API base URL (for example, a Cloudflare Worker URL) | Yes | - |
| `report_dir` | Directory where test reports are stored | No | `reports` |
| `include_tags` | Test tags to include (comma-separated, for example: `smoke,regression`) | No | - |
| `exclude_tags` | Test tags to exclude (comma-separated) | No | - |
| `variable_file` | Path to a Robot Framework variable file (comma-separated for multiple files) | No | - |
| `processes` | Number of parallel processes (`1` runs tests sequentially) | No | `1` |
| `timeout` | Test timeout in Robot Framework format (for example: `5 minutes`) | No | - |
| `log_level` | Log level (`TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`, `NONE`) | No | `INFO` |
| `artifact_name` | Name of the uploaded artifact | No | `robot-framework-reports` |
| `artifact_retention_days` | Artifact retention in days | No | `7` |

## Outputs

| Output | Description |
|--------|-------------|
| `exit_code` | Robot Framework exit code (`0` means success; non-zero means failure) |
| `report_path` | Path to the generated test report directory |

## Use Cases

- **Cloudflare Workers** deployment verification
- **AWS Lambda / Vercel** serverless API testing
- **CI/CD pipeline** smoke tests

---

# Unified Security & Infrastructure Platform (Project Context)

> The sections below describe the full platform this Action comes from.

[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![Python Version](https://img.shields.io/badge/Python-3.11+-blue.svg)](https://python.org)
[![Docker](https://img.shields.io/badge/Docker-20.10+-blue.svg)](https://docker.com)

[Traditional Chinese](README.zh-TW.md) | English

## Overview

A comprehensive, cloud-native security and infrastructure management platform combining:
- **IDS/IPS System** - Real-time intrusion detection and prevention
- **AI/ML Threat Detection** - Deep learning-based security analysis
- **Quantum Computing Integration** - IBM Quantum for advanced cryptography
- **Security Scanning Tools** - Integrated Nuclei, Nmap, AMASS scanners
- **Multi-Cloud Deployment** - Support for Cloudflare Workers, OCI, IBM Cloud

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│           Unified Security & Infrastructure Platform         │
└─────────────────────────────────────────────────────────────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
    ┌────▼─────┐      ┌──────▼──────┐      ┌─────▼──────┐
    │ Frontend │      │   Backend   │      │ AI/Quantum │
    │  (React) │      │    (Go)     │      │  (Python)  │
    └────┬─────┘      └──────┬──────┘      └─────┬──────┘
         │                   │                    │
         └───────────────────┼────────────────────┘
                             │
                    ┌────────▼────────┐
                    │  Infrastructure  │
                    │   - Docker       │
                    │   - Kubernetes   │
                    │   - Multi-Cloud  │
                    └──────────────────┘
```

## Quick Start

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- Go 1.24+ (for local development)
- Python 3.11+ (for AI/Quantum features)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/dennislee928/WHY_MR_ANDERSON_WHY.git
   cd WHY_MR_ANDERSON_WHY
   ```

2. **Setup Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start Services**
   ```bash
   docker-compose up -d
   ```

## Documentation

- [Architecture Details](docs/architecture/system-design.md)
- [API Reference](docs/development/api-reference.md)
- [Security Guide](docs/security/)
- [Deployment Guide](docs/deployment/)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

- **Author**: Dennislee928
- **Project**: [GitHub Repository](https://github.com/dennislee928/WHY_MR_ANDERSON_WHY)
