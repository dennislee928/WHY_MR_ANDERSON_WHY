# CI/CD 配置更新指南

執行專案重整後需要更新的 CI/CD 配置文件。

---

## 📋 更新清單

### 1. `.github/workflows/ci.yml`

#### 當前問題
- ❌ Line 107: 重複的 Dockerfile 路徑
- ❌ Line 141: 重複的 needs 依賴

#### 修正方案

**第一處修正 (Line 93-114):**
```yaml
# 當前（第 93-114 行）
- name: Login to GitHub Container Registry
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}

- name: Build and Push Docker Image to GHCR
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./Dockerfile.${{ matrix.image }}  # ❌ 錯誤路徑
    push: ${{ github.event_name == 'push' }}
    tags: |
      ghcr.io/${{ github.repository_owner }}/mitake_${{ matrix.image }}:latest
      ghcr.io/${{ github.repository_owner }}/mitake_${{ matrix.image }}:${{ github.sha }}
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**改為:**
```yaml
- name: Login to GitHub Container Registry
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}

- name: Build and Push Docker Image to GHCR
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./build/docker/${{ matrix.image }}.dockerfile  # ✅ 新路徑
    push: ${{ github.event_name == 'push' }}
    tags: |
      ghcr.io/${{ github.repository_owner }}/pandora_${{ matrix.image }}:latest
      ghcr.io/${{ github.repository_owner }}/pandora_${{ matrix.image }}:${{ github.sha }}
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**第二處修正 (Line 138-151):**
```yaml
# 當前（第 138-151 行）
final-check:
  runs-on: ubuntu-latest
  needs: [basic-check, frontend-check, docker-build-test, security-scan]  # ❌ 第一個 needs
  needs: [basic-check,  docker-build-and-push]  # ❌ 重複的 needs (且名稱錯誤)
  if: always()
  steps:
  - name: All checks completed
    run: |
      echo "✅ CI Pipeline 完成！"
      echo "- 基本檢查: ${{ needs.basic-check.result }}"
      echo "- 前端檢查: ${{ needs.frontend-check.result }}"
      echo "- Docker 建置: ${{ needs.docker-build-test.result }}"
      echo "- 安全掃描: ${{ needs.security-scan.result }}"
```

**改為:**
```yaml
final-check:
  runs-on: ubuntu-latest
  needs: [basic-check, frontend-check, docker-build-test, security-scan]  # ✅ 只保留一個正確的 needs
  if: always()
  steps:
  - name: All checks completed
    run: |
      echo "✅ CI Pipeline 完成！"
      echo "- 基本檢查: ${{ needs.basic-check.result }}"
      echo "- 前端檢查: ${{ needs.frontend-check.result }}"
      echo "- Docker 建置: ${{ needs.docker-build-test.result }}"
      echo "- 安全掃描: ${{ needs.security-scan.result }}"
```

---

### 2. `.github/workflows/deploy-gcp.yml`

#### 修正方案

**Line 61-70 (Agent Image):**
```yaml
# 當前
- name: Build and push Agent image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./Dockerfile.agent  # ❌ 舊路徑
    push: true
    tags: ${{ env.GCP_REGISTRY }}/pandora-agent:${{ github.sha }}
    tags: ${{ env.GCP_REGISTRY }}/pandora-agent:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**改為:**
```yaml
- name: Build and push Agent image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./build/docker/agent.dockerfile  # ✅ 新路徑
    push: true
    tags: |
      ${{ env.GCP_REGISTRY }}/pandora-agent:${{ github.sha }}
      ${{ env.GCP_REGISTRY }}/pandora-agent:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**Line 72-81 (Console Image):**
```yaml
# 當前
- name: Build and push Console image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./Dockerfile  # ❌ 舊路徑
    push: true
    tags: ${{ env.GCP_REGISTRY }}/pandora-console:${{ github.sha }}
    tags: ${{ env.GCP_REGISTRY }}/pandora-console:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**改為:**
```yaml
- name: Build and push Console image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./build/docker/console.dockerfile  # ✅ 新路徑
    push: true
    tags: |
      ${{ env.GCP_REGISTRY }}/pandora-console:${{ github.sha }}
      ${{ env.GCP_REGISTRY }}/pandora-console:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**Line 109-118 (K8s Deployment):**
```yaml
# 當前
- name: Update image tags in GCP manifests
  run: |
    find k8s-gcp/ -name "*.yaml" -exec sed -i "s|gcr.io/YOUR_PROJECT_ID|${{ env.GCP_REGISTRY }}|g" {} \;
    find k8s-gcp/ -name "*.yaml" -exec sed -i "s|:latest|:${{ github.sha }}|g" {} \;

- name: Deploy to GKE
  run: |
    kubectl apply -k k8s-gcp/  # ❌ 舊路徑
    kubectl rollout status deployment/pandora-agent -n pandora-box --timeout=300s
    kubectl rollout status deployment/pandora-console -n pandora-box --timeout=300s
```

**改為:**
```yaml
- name: Update image tags in GCP manifests
  run: |
    find deployments/kubernetes/gcp/ -name "*.yaml" -exec sed -i "s|gcr.io/YOUR_PROJECT_ID|${{ env.GCP_REGISTRY }}|g" {} \;
    find deployments/kubernetes/gcp/ -name "*.yaml" -exec sed -i "s|:latest|:${{ github.sha }}|g" {} \;

- name: Deploy to GKE
  run: |
    kubectl apply -k deployments/kubernetes/gcp/  # ✅ 新路徑
    kubectl rollout status deployment/pandora-agent -n pandora-box --timeout=300s
    kubectl rollout status deployment/pandora-console -n pandora-box --timeout=300s
```

---

### 3. `.github/workflows/deploy-oci.yml`

#### 語法錯誤修正

**Line 5:**
```yaml
# 當前
on:
  push:
    branches: [ temp_locked" ]  # ❌ 語法錯誤：缺少開頭引號
```

**改為:**
```yaml
on:
  push:
    branches: [ "temp_locked" ]  # ✅ 正確語法
```

#### Dockerfile 路徑修正

**Line 58-67:**
```yaml
# 當前
- name: Build and push Agent image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./Dockerfile.agent  # ❌ 舊路徑
    push: true
    tags: ${{ env.OCI_REGISTRY }}/pandora-agent:${{ github.sha }}
    tags: ${{ env.OCI_REGISTRY }}/pandora-agent:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**改為:**
```yaml
- name: Build and push Agent image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./build/docker/agent.dockerfile  # ✅ 新路徑
    push: true
    tags: |
      ${{ env.OCI_REGISTRY }}/pandora-agent:${{ github.sha }}
      ${{ env.OCI_REGISTRY }}/pandora-agent:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**Line 69-78:**
```yaml
# 當前
- name: Build and push Console image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./Dockerfile  # ❌ 舊路徑
    push: true
    tags: ${{ env.OCI_REGISTRY }}/pandora-console:${{ github.sha }}
    tags: ${{ env.OCI_REGISTRY }}/pandora-console:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**改為:**
```yaml
- name: Build and push Console image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./build/docker/console.dockerfile  # ✅ 新路徑
    push: true
    tags: |
      ${{ env.OCI_REGISTRY }}/pandora-console:${{ github.sha }}
      ${{ env.OCI_REGISTRY }}/pandora-console:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**Line 104-112:**
```yaml
# 當前
- name: Update image tags in manifests
  run: |
    find k8s/ -name "*.yaml" -exec sed -i "s|iad.ocir.io/YOUR_NAMESPACE|${{ env.OCI_REGISTRY }}|g" {} \;
    find k8s/ -name "*.yaml" -exec sed -i "s|:latest|:${{ github.sha }}|g" {} \;

- name: Deploy to Kubernetes
  run: |
    kubectl apply -k k8s/  # ❌ 舊路徑
    kubectl rollout status deployment/pandora-agent -n pandora-box --timeout=300s
    kubectl rollout status deployment/pandora-console -n pandora-box --timeout=300s
```

**改為:**
```yaml
- name: Update image tags in manifests
  run: |
    find deployments/kubernetes/base/ -name "*.yaml" -exec sed -i "s|iad.ocir.io/YOUR_NAMESPACE|${{ env.OCI_REGISTRY }}|g" {} \;
    find deployments/kubernetes/base/ -name "*.yaml" -exec sed -i "s|:latest|:${{ github.sha }}|g" {} \;

- name: Deploy to Kubernetes
  run: |
    kubectl apply -k deployments/kubernetes/base/  # ✅ 新路徑
    kubectl rollout status deployment/pandora-agent -n pandora-box --timeout=300s
    kubectl rollout status deployment/pandora-console -n pandora-box --timeout=300s
```

---

### 4. `.github/workflows/deploy-paas.yml`

#### Dockerfile 路徑修正

**Line 131-139:**
```yaml
# 當前
- name: Build and push Agent image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./Dockerfile.agent.koyeb  # ❌ 舊路徑
    push: true
    tags: ${{ secrets.DOCKER_USERNAME }}/pandora-agent:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**改為:**
```yaml
- name: Build and push Agent image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./build/docker/agent.koyeb.dockerfile  # ✅ 新路徑
    push: true
    tags: ${{ secrets.DOCKER_USERNAME }}/pandora-agent:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**Line 174-182:**
```yaml
# 當前
- name: Build and push UI image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./Dockerfile.ui.patr  # ❌ 舊路徑
    push: true
    tags: ${{ secrets.DOCKER_USERNAME }}/axiom-ui:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**改為:**
```yaml
- name: Build and push UI image
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./build/docker/ui.patr.dockerfile  # ✅ 新路徑
    push: true
    tags: ${{ secrets.DOCKER_USERNAME }}/axiom-ui:latest
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

**Line 208-212:**
```yaml
# 當前
- name: Deploy to Fly.io
  env:
    FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
  run: |
    flyctl deploy --config fly.toml --dockerfile Dockerfile.monitoring --remote-only  # ❌ 舊路徑
```

**改為:**
```yaml
- name: Deploy to Fly.io
  env:
    FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
  run: |
    flyctl deploy --config deployments/paas/flyio/fly.toml --dockerfile build/docker/monitoring.dockerfile --remote-only  # ✅ 新路徑
```

---

## 🚀 快速更新腳本

創建一個腳本來自動應用這些更改：

```powershell
# update-ci-workflows.ps1

# 備份
Copy-Item .github/workflows/ci.yml .github/workflows/ci.yml.backup
Copy-Item .github/workflows/deploy-gcp.yml .github/workflows/deploy-gcp.yml.backup
Copy-Item .github/workflows/deploy-oci.yml .github/workflows/deploy-oci.yml.backup
Copy-Item .github/workflows/deploy-paas.yml .github/workflows/deploy-paas.yml.backup

# 更新路徑
(Get-Content .github/workflows/ci.yml) `
    -replace 'file: ./Dockerfile\.', 'file: ./build/docker/' `
    -replace '\$\{\{ matrix\.image \}\}$', '${{ matrix.image }}.dockerfile' |
    Set-Content .github/workflows/ci.yml

Write-Host "✅ CI workflows 已更新"
Write-Host "📝 請手動檢查並移除重複的 needs 行"
```

---

## ✅ 驗證清單

更新後請驗證：

- [ ] 所有 Dockerfile 路徑已更新
- [ ] 所有 K8s manifests 路徑已更新
- [ ] 移除重複的 needs 依賴
- [ ] 修正語法錯誤
- [ ] 本地測試建置
- [ ] 提交前運行 `yamllint` 檢查

---

## 📞 遇到問題？

如果 CI/CD 失敗：
1. 檢查文件是否實際移動到新位置
2. 檢查路徑分隔符（`/` vs `\`）
3. 確認 context 設置正確
4. 查看 GitHub Actions 日誌

**記得**: 先在測試分支驗證，確認無誤後再合併到 main！
