#!/bin/bash
# 證書輪換腳本
# Certificate rotation script

set -e

CERTS_DIR="deployments/onpremise/certs"
BACKUP_DIR="$CERTS_DIR/backup_$(date +%Y%m%d_%H%M%S)"
VALIDITY_DAYS=90

echo "========================================="
echo "  Pandora Box Certificate Rotation"
echo "========================================="
echo ""

# 1. 備份現有證書
echo "1. Backing up existing certificates..."
mkdir -p $BACKUP_DIR
cp -r $CERTS_DIR/* $BACKUP_DIR/ 2>/dev/null || true
echo "✅ Certificates backed up to: $BACKUP_DIR"

# 2. 檢查證書過期時間
echo ""
echo "2. Checking certificate expiration..."

check_cert_expiry() {
    local cert_file=$1
    local service_name=$2
    
    if [ -f "$cert_file" ]; then
        expiry_date=$(openssl x509 -in "$cert_file" -noout -enddate | cut -d= -f2)
        expiry_epoch=$(date -d "$expiry_date" +%s 2>/dev/null || date -j -f "%b %d %T %Y %Z" "$expiry_date" +%s)
        now_epoch=$(date +%s)
        days_left=$(( ($expiry_epoch - $now_epoch) / 86400 ))
        
        echo "   $service_name: $days_left days remaining"
        
        if [ $days_left -lt 7 ]; then
            echo "   ⚠️  WARNING: Certificate expires in less than 7 days!"
        fi
    fi
}

check_cert_expiry "$CERTS_DIR/ca/ca.crt" "CA"
check_cert_expiry "$CERTS_DIR/device/server.crt" "Device Service"
check_cert_expiry "$CERTS_DIR/network/server.crt" "Network Service"
check_cert_expiry "$CERTS_DIR/control/server.crt" "Control Service"

# 3. 生成新證書
echo ""
echo "3. Generating new certificates..."
./scripts/generate-certs.sh

# 4. 重新載入服務（零停機時間）
echo ""
echo "4. Reloading services..."

# 發送 SIGHUP 信號讓服務重新載入證書
docker-compose -f deployments/onpremise/docker-compose.yml kill -s SIGHUP device-service
docker-compose -f deployments/onpremise/docker-compose.yml kill -s SIGHUP network-service
docker-compose -f deployments/onpremise/docker-compose.yml kill -s SIGHUP control-service

echo "✅ Services reloaded"

# 5. 驗證新證書
echo ""
echo "5. Verifying new certificates..."

# 等待服務重新載入
sleep 5

# 檢查服務健康狀態
for port in 8081 8082 8083; do
    if curl -s -f http://localhost:$port/health > /dev/null; then
        echo "✅ Service on port $port is healthy"
    else
        echo "❌ Service on port $port is unhealthy"
    fi
done

echo ""
echo "========================================="
echo "  ✅ Certificate rotation completed!"
echo "========================================="
echo ""
echo "Backup location: $BACKUP_DIR"
echo "New certificates valid for: $VALIDITY_DAYS days"
echo ""
echo "Next rotation due: $(date -d "+$VALIDITY_DAYS days" +%Y-%m-%d 2>/dev/null || date -v +${VALIDITY_DAYS}d +%Y-%m-%d)"
echo ""

