# CI/CD Workflows 指南

> **平台**: GitHub Actions

---

## 📋 Workflow 列表

### 主要 Workflows（dev 分支）

1. **ci.yml** - CI Pipeline
   - 觸發: push/PR to dev
   - 功能: 測試、構建、安全掃描

2. **build-onpremise-installers.yml** - 安裝檔構建
   - 觸發: push to dev, tags, 手動
   - 功能: 生成 .exe/.deb/.rpm/.iso/.ova

### 停用 Workflows（僅 main 分支）

3. **deploy-gcp.yml** - GCP部署（僅手動）
4. **deploy-oci.yml** - OCI部署（僅手動）
5. **deploy-paas.yml** - PaaS部署（僅手動）
6. **terraform-deploy.yml** - Terraform部署（僅手動）

---

## 🚀 觸發方式

### 自動觸發

```bash
git push origin dev
# 觸發: ci.yml, build-onpremise-installers.yml
```

### 標籤觸發

```bash
git tag -a v3.0.0 -m "Release v3.0.0"
git push origin v3.0.0
# 觸發: 所有 workflows + Release
```

### 手動觸發

GitHub → Actions → 選擇 workflow → Run workflow

---

**維護**: DevOps Team

