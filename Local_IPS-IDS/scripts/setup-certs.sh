#!/bin/bash
# Pandora Box Console SSL/TLS 憑證設定腳本

set -e

# 設定變數
CERT_DIR="./certs"
COUNTRY="TW"
STATE="Taipei"
CITY="Taipei"
ORG="Pandora Box Console"
ORG_UNIT="Security Team"
CA_COMMON_NAME="Pandora CA"
SERVER_COMMON_NAME="pandora-server"
CLIENT_COMMON_NAME="pandora-client"
DAYS=365
CA_DAYS=3650

# 顏色定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日誌函數
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 檢查 OpenSSL
check_openssl() {
    if ! command -v openssl &> /dev/null; then
        log_error "OpenSSL 未安裝，請先安裝 OpenSSL"
        exit 1
    fi
    log_info "OpenSSL 版本: $(openssl version)"
}

# 建立憑證目錄
create_cert_dir() {
    log_step "建立憑證目錄..."
    mkdir -p "$CERT_DIR"
    chmod 755 "$CERT_DIR"
    log_info "憑證目錄已建立: $CERT_DIR"
}

# 產生 CA 私鑰
generate_ca_key() {
    log_step "產生 CA 私鑰..."
    openssl genrsa -out "$CERT_DIR/ca.key" 4096
    chmod 600 "$CERT_DIR/ca.key"
    log_info "CA 私鑰已產生: $CERT_DIR/ca.key"
}

# 產生 CA 憑證
generate_ca_cert() {
    log_step "產生 CA 憑證..."
    openssl req -new -x509 \
        -key "$CERT_DIR/ca.key" \
        -sha256 \
        -subj "/C=$COUNTRY/ST=$STATE/L=$CITY/O=$ORG/OU=$ORG_UNIT/CN=$CA_COMMON_NAME" \
        -days $CA_DAYS \
        -out "$CERT_DIR/ca.crt"
    
    chmod 644 "$CERT_DIR/ca.crt"
    log_info "CA 憑證已產生: $CERT_DIR/ca.crt"
}

# 產生伺服器私鑰
generate_server_key() {
    log_step "產生伺服器私鑰..."
    openssl genrsa -out "$CERT_DIR/server.key" 4096
    chmod 600 "$CERT_DIR/server.key"
    log_info "伺服器私鑰已產生: $CERT_DIR/server.key"
}

# 產生伺服器憑證簽名請求
generate_server_csr() {
    log_step "產生伺服器憑證簽名請求..."
    openssl req -new \
        -key "$CERT_DIR/server.key" \
        -out "$CERT_DIR/server.csr" \
        -subj "/C=$COUNTRY/ST=$STATE/L=$CITY/O=$ORG/OU=$ORG_UNIT/CN=$SERVER_COMMON_NAME"
    log_info "伺服器 CSR 已產生: $CERT_DIR/server.csr"
}

# 建立伺服器憑證擴展檔案
create_server_ext() {
    log_step "建立伺服器憑證擴展檔案..."
    cat > "$CERT_DIR/server.ext" << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = pandora-server
DNS.3 = pandora-agent
DNS.4 = axiom-ui
DNS.5 = grafana
DNS.6 = prometheus
DNS.7 = loki
DNS.8 = alertmanager
IP.1 = 127.0.0.1
IP.2 = ::1
IP.3 = 172.20.0.2
IP.4 = 172.20.0.3
IP.5 = 172.20.0.4
EOF
    log_info "伺服器擴展檔案已建立: $CERT_DIR/server.ext"
}

# 簽署伺服器憑證
sign_server_cert() {
    log_step "簽署伺服器憑證..."
    openssl x509 -req \
        -in "$CERT_DIR/server.csr" \
        -CA "$CERT_DIR/ca.crt" \
        -CAkey "$CERT_DIR/ca.key" \
        -CAcreateserial \
        -out "$CERT_DIR/server.crt" \
        -days $DAYS \
        -sha256 \
        -extfile "$CERT_DIR/server.ext"
    
    chmod 644 "$CERT_DIR/server.crt"
    log_info "伺服器憑證已簽署: $CERT_DIR/server.crt"
}

# 產生客戶端私鑰
generate_client_key() {
    log_step "產生客戶端私鑰..."
    openssl genrsa -out "$CERT_DIR/client.key" 4096
    chmod 600 "$CERT_DIR/client.key"
    log_info "客戶端私鑰已產生: $CERT_DIR/client.key"
}

# 產生客戶端憑證簽名請求
generate_client_csr() {
    log_step "產生客戶端憑證簽名請求..."
    openssl req -new \
        -key "$CERT_DIR/client.key" \
        -out "$CERT_DIR/client.csr" \
        -subj "/C=$COUNTRY/ST=$STATE/L=$CITY/O=$ORG/OU=$ORG_UNIT/CN=$CLIENT_COMMON_NAME"
    log_info "客戶端 CSR 已產生: $CERT_DIR/client.csr"
}

# 建立客戶端憑證擴展檔案
create_client_ext() {
    log_step "建立客戶端憑證擴展檔案..."
    cat > "$CERT_DIR/client.ext" << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
extendedKeyUsage = clientAuth
EOF
    log_info "客戶端擴展檔案已建立: $CERT_DIR/client.ext"
}

# 簽署客戶端憑證
sign_client_cert() {
    log_step "簽署客戶端憑證..."
    openssl x509 -req \
        -in "$CERT_DIR/client.csr" \
        -CA "$CERT_DIR/ca.crt" \
        -CAkey "$CERT_DIR/ca.key" \
        -CAcreateserial \
        -out "$CERT_DIR/client.crt" \
        -days $DAYS \
        -sha256 \
        -extfile "$CERT_DIR/client.ext"
    
    chmod 644 "$CERT_DIR/client.crt"
    log_info "客戶端憑證已簽署: $CERT_DIR/client.crt"
}

# 驗證憑證
verify_certificates() {
    log_step "驗證憑證..."
    
    # 驗證伺服器憑證
    if openssl verify -CAfile "$CERT_DIR/ca.crt" "$CERT_DIR/server.crt" > /dev/null 2>&1; then
        log_info "伺服器憑證驗證成功"
    else
        log_error "伺服器憑證驗證失敗"
        exit 1
    fi
    
    # 驗證客戶端憑證
    if openssl verify -CAfile "$CERT_DIR/ca.crt" "$CERT_DIR/client.crt" > /dev/null 2>&1; then
        log_info "客戶端憑證驗證成功"
    else
        log_error "客戶端憑證驗證失敗"
        exit 1
    fi
}

# 產生 PKCS#12 格式檔案 (可選)
generate_pkcs12() {
    log_step "產生 PKCS#12 格式檔案..."
    
    # 伺服器 PKCS#12
    openssl pkcs12 -export \
        -out "$CERT_DIR/server.p12" \
        -inkey "$CERT_DIR/server.key" \
        -in "$CERT_DIR/server.crt" \
        -certfile "$CERT_DIR/ca.crt" \
        -password pass:pandora123
    
    # 客戶端 PKCS#12
    openssl pkcs12 -export \
        -out "$CERT_DIR/client.p12" \
        -inkey "$CERT_DIR/client.key" \
        -in "$CERT_DIR/client.crt" \
        -certfile "$CERT_DIR/ca.crt" \
        -password pass:pandora123
    
    log_info "PKCS#12 檔案已產生 (密碼: pandora123)"
}

# 顯示憑證資訊
show_certificate_info() {
    log_step "憑證資訊摘要..."
    
    echo ""
    echo "CA 憑證:"
    openssl x509 -in "$CERT_DIR/ca.crt" -noout -subject -issuer -dates
    
    echo ""
    echo "伺服器憑證:"
    openssl x509 -in "$CERT_DIR/server.crt" -noout -subject -issuer -dates
    
    echo ""
    echo "客戶端憑證:"
    openssl x509 -in "$CERT_DIR/client.crt" -noout -subject -issuer -dates
    
    echo ""
    echo "憑證檔案清單:"
    ls -la "$CERT_DIR"/*.{crt,key,p12} 2>/dev/null || true
}

# 清理暫存檔案
cleanup() {
    log_step "清理暫存檔案..."
    rm -f "$CERT_DIR"/*.csr
    rm -f "$CERT_DIR"/*.ext
    rm -f "$CERT_DIR"/*.srl
    log_info "暫存檔案已清理"
}

# 設定 Docker 權限
setup_docker_permissions() {
    log_step "設定 Docker 權限..."
    
    # 建立 Docker 使用的憑證目錄
    mkdir -p "$CERT_DIR/docker"
    
    # 複製憑證到 Docker 目錄
    cp "$CERT_DIR/ca.crt" "$CERT_DIR/docker/"
    cp "$CERT_DIR/server.crt" "$CERT_DIR/docker/"
    cp "$CERT_DIR/server.key" "$CERT_DIR/docker/"
    cp "$CERT_DIR/client.crt" "$CERT_DIR/docker/"
    cp "$CERT_DIR/client.key" "$CERT_DIR/docker/"
    
    # 設定適當權限
    chmod 644 "$CERT_DIR/docker"/*.crt
    chmod 600 "$CERT_DIR/docker"/*.key
    
    log_info "Docker 憑證權限已設定"
}

# 產生憑證指紋
generate_fingerprints() {
    log_step "產生憑證指紋..."
    
    echo ""
    echo "憑證指紋 (SHA256):"
    echo "CA 憑證: $(openssl x509 -fingerprint -sha256 -noout -in "$CERT_DIR/ca.crt" | cut -d= -f2)"
    echo "伺服器憑證: $(openssl x509 -fingerprint -sha256 -noout -in "$CERT_DIR/server.crt" | cut -d= -f2)"
    echo "客戶端憑證: $(openssl x509 -fingerprint -sha256 -noout -in "$CERT_DIR/client.crt" | cut -d= -f2)"
}

# 主要函數
main() {
    echo "====================================="
    echo "Pandora Box Console SSL/TLS 憑證設定"
    echo "====================================="
    echo ""
    
    # 檢查是否已存在憑證
    if [[ -f "$CERT_DIR/ca.crt" && -f "$CERT_DIR/server.crt" && -f "$CERT_DIR/client.crt" ]]; then
        log_warn "憑證檔案已存在"
        read -p "是否要重新產生憑證? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "跳過憑證產生"
            show_certificate_info
            exit 0
        fi
        
        log_warn "刪除現有憑證..."
        rm -rf "$CERT_DIR"
    fi
    
    # 執行憑證產生流程
    check_openssl
    create_cert_dir
    
    # CA 憑證
    generate_ca_key
    generate_ca_cert
    
    # 伺服器憑證
    generate_server_key
    generate_server_csr
    create_server_ext
    sign_server_cert
    
    # 客戶端憑證
    generate_client_key
    generate_client_csr
    create_client_ext
    sign_client_cert
    
    # 驗證和清理
    verify_certificates
    generate_pkcs12
    setup_docker_permissions
    cleanup
    
    # 顯示結果
    show_certificate_info
    generate_fingerprints
    
    echo ""
    log_info "SSL/TLS 憑證設定完成！"
    log_info "憑證檔案位於: $CERT_DIR"
    echo ""
    echo "使用說明:"
    echo "1. CA 憑證: $CERT_DIR/ca.crt (信任此憑證)"
    echo "2. 伺服器憑證: $CERT_DIR/server.crt + $CERT_DIR/server.key"
    echo "3. 客戶端憑證: $CERT_DIR/client.crt + $CERT_DIR/client.key"
    echo "4. PKCS#12 檔案: $CERT_DIR/server.p12 + $CERT_DIR/client.p12 (密碼: pandora123)"
    echo ""
}

# 執行主函數
main "$@"
