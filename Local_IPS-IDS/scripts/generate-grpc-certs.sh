#!/bin/bash
# Pandora Box Console - gRPC TLS 證書生成腳本
# 為所有 gRPC 服務生成 TLS 證書

set -e

CERT_DIR="certs"
DAYS_VALID=365

echo "========================================="
echo "  🔐 gRPC TLS 證書生成工具"
echo "========================================="

# 創建證書目錄
mkdir -p "$CERT_DIR"
cd "$CERT_DIR"

# 1. 生成 CA 證書（根證書）
echo -e "\n📜 步驟 1/4: 生成 CA 根證書..."

# Windows Git Bash 兼容性修復
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
  # Windows: 使用 MSYS_NO_PATHCONV 防止路徑轉換
  MSYS_NO_PATHCONV=1 openssl req -x509 -newkey rsa:4096 -days $DAYS_VALID -nodes \
    -keyout ca-key.pem \
    -out ca-cert.pem \
    -subj "//C=TW\ST=Taipei\L=Taipei\O=Pandora Security\OU=Security\CN=Pandora CA"
else
  # Linux/macOS
  openssl req -x509 -newkey rsa:4096 -days $DAYS_VALID -nodes \
    -keyout ca-key.pem \
    -out ca-cert.pem \
    -subj "/C=TW/ST=Taipei/L=Taipei/O=Pandora Security/OU=Security/CN=Pandora CA"
fi

echo "✅ CA 證書已生成: ca-cert.pem, ca-key.pem"

# 2. 生成服務證書（Device Service）
echo -e "\n📜 步驟 2/4: 生成 Device Service 證書..."

if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
  MSYS_NO_PATHCONV=1 openssl req -newkey rsa:4096 -nodes \
    -keyout device-service-key.pem \
    -out device-service-req.pem \
    -subj "//C=TW\ST=Taipei\L=Taipei\O=Pandora Security\OU=Device Service\CN=device-service"
else
  openssl req -newkey rsa:4096 -nodes \
    -keyout device-service-key.pem \
    -out device-service-req.pem \
    -subj "/C=TW/ST=Taipei/L=Taipei/O=Pandora Security/OU=Device Service/CN=device-service"
fi

openssl x509 -req -in device-service-req.pem -days $DAYS_VALID \
  -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial \
  -out device-service-cert.pem \
  -extfile <(printf "subjectAltName=DNS:device-service,DNS:localhost,IP:127.0.0.1")

echo "✅ Device Service 證書已生成"

# 3. 生成服務證書（Network Service）
echo -e "\n📜 步驟 3/4: 生成 Network Service 證書..."

if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
  MSYS_NO_PATHCONV=1 openssl req -newkey rsa:4096 -nodes \
    -keyout network-service-key.pem \
    -out network-service-req.pem \
    -subj "//C=TW\ST=Taipei\L=Taipei\O=Pandora Security\OU=Network Service\CN=network-service"
else
  openssl req -newkey rsa:4096 -nodes \
    -keyout network-service-key.pem \
    -out network-service-req.pem \
    -subj "/C=TW/ST=Taipei/L=Taipei/O=Pandora Security/OU=Network Service/CN=network-service"
fi

openssl x509 -req -in network-service-req.pem -days $DAYS_VALID \
  -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial \
  -out network-service-cert.pem \
  -extfile <(printf "subjectAltName=DNS:network-service,DNS:localhost,IP:127.0.0.1")

echo "✅ Network Service 證書已生成"

# 4. 生成服務證書（Control Service）
echo -e "\n📜 步驟 4/4: 生成 Control Service 證書..."

if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
  MSYS_NO_PATHCONV=1 openssl req -newkey rsa:4096 -nodes \
    -keyout control-service-key.pem \
    -out control-service-req.pem \
    -subj "//C=TW\ST=Taipei\L=Taipei\O=Pandora Security\OU=Control Service\CN=control-service"
else
  openssl req -newkey rsa:4096 -nodes \
    -keyout control-service-key.pem \
    -out control-service-req.pem \
    -subj "/C=TW/ST=Taipei/L=Taipei/O=Pandora Security/OU=Control Service/CN=control-service"
fi

openssl x509 -req -in control-service-req.pem -days $DAYS_VALID \
  -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial \
  -out control-service-cert.pem \
  -extfile <(printf "subjectAltName=DNS:control-service,DNS:localhost,IP:127.0.0.1")

echo "✅ Control Service 證書已生成"

# 清理臨時文件
rm -f *.req.pem

# 設置權限
chmod 600 *-key.pem
chmod 644 *-cert.pem

echo -e "\n========================================="
echo "  ✅ 所有 gRPC TLS 證書已生成！"
echo "========================================="
echo -e "\n📁 證書位置: $(pwd)"
echo -e "\n📋 生成的證書:"
ls -lh *.pem

echo -e "\n🔐 使用方式:"
echo "  1. 複製證書到 configs/certs/"
echo "  2. 更新 gRPC 服務器配置"
echo "  3. 更新 gRPC 客戶端配置"
echo "  4. 重新構建並重啟服務"
echo -e "\n📖 詳細文檔: docs/GRPC-TLS-SETUP.md\n"

