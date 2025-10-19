#!/bin/bash
# deploy-flyio.sh

echo "🚀 開始部署 Fly.io 監控系統..."

# 檢查是否已登入
if ! fly auth whoami &> /dev/null; then
    echo "請先登入 Fly.io: fly auth login"
    exit 1
fi

# 建立 Volumes（如果不存在）
echo "📦 建立持久化儲存..."
fly volumes create prometheus_data --size 3 --region nrt || true
fly volumes create loki_data --size 3 --region nrt || true
fly volumes create grafana_data --size 1 --region nrt || true
fly volumes create alertmanager_data --size 1 --region nrt || true

# 設定 Secrets
echo "🔐 設定環境變數..."
fly secrets set GRAFANA_ADMIN_PASSWORD=pandora123 || true

# 部署
echo "🚀 部署應用..."
fly deploy --config fly.toml --dockerfile Dockerfile.monitoring --remote-only

echo "✅ 部署完成！"
echo "🌐 Grafana: https://pandora-monitoring.fly.dev"
echo "📊 Prometheus: https://pandora-monitoring.fly.dev/prometheus"
echo "📝 Loki: https://pandora-monitoring.fly.dev/loki"
echo "🚨 AlertManager: https://pandora-monitoring.fly.dev/alertmanager"