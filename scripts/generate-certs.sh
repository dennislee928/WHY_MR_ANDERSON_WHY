#!/bin/bash
# 生成微服務 mTLS 證書腳本
# Generate mTLS certificates for microservices

set -e

CERTS_DIR="deployments/onpremise/certs"
VALIDITY_DAYS=90

echo "========================================="
echo "  Pandora Box mTLS Certificate Generator"
echo "========================================="
echo ""

# 創建證書目錄
mkdir -p $CERTS_DIR/{ca,device,network,control,prometheus,grafana,loki}

# 1. 生成 CA 證書
echo "1. Generating CA certificate..."
openssl genrsa -out $CERTS_DIR/ca/ca.key 4096
openssl req -new -x509 \
  -key $CERTS_DIR/ca/ca.key \
  -sha256 \
  -subj "//C=TW//ST=Taipei//O=Pandora//CN=Pandora Root CA" \
  -days $VALIDITY_DAYS \
  -out $CERTS_DIR/ca/ca.crt

echo "✅ CA certificate generated"

# 2. 生成 Device Service 證書
echo ""
echo "2. Generating Device Service certificate..."
openssl genrsa -out $CERTS_DIR/device/server.key 4096
openssl req -new \
  -key $CERTS_DIR/device/server.key \
  -out $CERTS_DIR/device/server.csr \
  -subj "//C=TW//ST=Taipei//O=Pandora//CN=device-service"

openssl x509 -req \
  -in $CERTS_DIR/device/server.csr \
  -CA $CERTS_DIR/ca/ca.crt \
  -CAkey $CERTS_DIR/ca/ca.key \
  -CAcreateserial \
  -out $CERTS_DIR/device/server.crt \
  -days $VALIDITY_DAYS \
  -sha256

echo "✅ Device Service certificate generated"

# 3. 生成 Network Service 證書
echo ""
echo "3. Generating Network Service certificate..."
openssl genrsa -out $CERTS_DIR/network/server.key 4096
openssl req -new \
  -key $CERTS_DIR/network/server.key \
  -out $CERTS_DIR/network/server.csr \
  -subj "//C=TW//ST=Taipei//O=Pandora//CN=network-service"

openssl x509 -req \
  -in $CERTS_DIR/network/server.csr \
  -CA $CERTS_DIR/ca/ca.crt \
  -CAkey $CERTS_DIR/ca/ca.key \
  -CAcreateserial \
  -out $CERTS_DIR/network/server.crt \
  -days $VALIDITY_DAYS \
  -sha256

echo "✅ Network Service certificate generated"

# 4. 生成 Control Service 證書
echo ""
echo "4. Generating Control Service certificate..."
openssl genrsa -out $CERTS_DIR/control/server.key 4096
openssl req -new \
  -key $CERTS_DIR/control/server.key \
  -out $CERTS_DIR/control/server.csr \
  -subj "//C=TW//ST=Taipei//O=Pandora//CN=control-service"

openssl x509 -req \
  -in $CERTS_DIR/control/server.csr \
  -CA $CERTS_DIR/ca/ca.crt \
  -CAkey $CERTS_DIR/ca/ca.key \
  -CAcreateserial \
  -out $CERTS_DIR/control/server.crt \
  -days $VALIDITY_DAYS \
  -sha256

echo "✅ Control Service certificate generated"

# 5. 生成客戶端證書（用於 Engine 調用微服務）
echo ""
echo "5. Generating client certificate..."
openssl genrsa -out $CERTS_DIR/ca/client.key 4096
openssl req -new \
  -key $CERTS_DIR/ca/client.key \
  -out $CERTS_DIR/ca/client.csr \
  -subj "//C=TW//ST=Taipei//O=Pandora//CN=pandora-client"

openssl x509 -req \
  -in $CERTS_DIR/ca/client.csr \
  -CA $CERTS_DIR/ca/ca.crt \
  -CAkey $CERTS_DIR/ca/ca.key \
  -CAcreateserial \
  -out $CERTS_DIR/ca/client.crt \
  -days $VALIDITY_DAYS \
  -sha256

echo "✅ Client certificate generated"

# 6. 複製 CA 證書到各服務目錄
echo ""
echo "6. Copying CA certificate to service directories..."
for service in device network control prometheus grafana loki; do
  cp $CERTS_DIR/ca/ca.crt $CERTS_DIR/$service/
done

echo "✅ CA certificate copied"

# 7. 設置權限
echo ""
echo "7. Setting permissions..."
chmod 600 $CERTS_DIR/*/*.key
chmod 644 $CERTS_DIR/*/*.crt

echo "✅ Permissions set"

# 8. 顯示證書資訊
echo ""
echo "========================================="
echo "  Certificate Information"
echo "========================================="
echo ""
echo "CA Certificate:"
openssl x509 -in $CERTS_DIR/ca/ca.crt -noout -subject -dates

echo ""
echo "Device Service Certificate:"
openssl x509 -in $CERTS_DIR/device/server.crt -noout -subject -dates

echo ""
echo "Network Service Certificate:"
openssl x509 -in $CERTS_DIR/network/server.crt -noout -subject -dates

echo ""
echo "Control Service Certificate:"
openssl x509 -in $CERTS_DIR/control/server.crt -noout -subject -dates

echo ""
echo "========================================="
echo "  ✅ All certificates generated!"
echo "========================================="
echo ""
echo "Certificates location: $CERTS_DIR"
echo "Validity: $VALIDITY_DAYS days"
echo ""
echo "⚠️  Remember to rotate certificates before expiration!"
echo "    Use: ./scripts/rotate-certs.sh"
echo ""

