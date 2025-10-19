#!/bin/bash

# Fly.io Volume 調整腳本
# 用於降低 pandora-monitoring 的 volume 大小以減少費用

set -e

echo "🔍 檢查當前 Fly.io volumes..."

# 檢查當前 volumes
echo "📋 當前 volumes 列表："
flyctl volumes list --app pandora-monitoring

echo ""
echo "⚠️  WARNING: 調整 volume 大小需要停機時間！"
echo "📌 建議步驟："
echo "1. 備份重要數據"
echo "2. 創建新的較小 volume (3GB)"
echo "3. 停止應用"
echo "4. 將數據遷移到新 volume"
echo "5. 刪除舊的大 volume"
echo "6. 重新啟動應用"
echo ""

read -p "是否繼續執行 volume 調整？ (y/N): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 取消操作"
    exit 1
fi

echo "🔧 開始 volume 調整流程..."

# 1. 創建新的 3GB volume
echo "📦 創建新的 3GB volume..."
flyctl volumes create monitoring_data_new \
    --app pandora-monitoring \
    --region nrt \
    --size 3

# 2. 停止應用以進行數據遷移
echo "⏹️  停止應用進行維護..."
flyctl apps suspend pandora-monitoring

echo ""
echo "✅ Volume 調整準備完成！"
echo ""
echo "📝 下一步手動操作："
echo "1. 使用 flyctl ssh console 連接到機器"
echo "2. 手動複製重要資料從 /data 到新 volume"
echo "3. 更新應用配置使用新 volume"
echo "4. 刪除舊 volume: flyctl volumes delete <OLD_VOLUME_ID>"
echo "5. 重新啟動應用: flyctl apps resume pandora-monitoring"
echo ""
echo "💡 或者使用簡單方式：重新部署應用讓它使用新的 3GB volume 配置"
