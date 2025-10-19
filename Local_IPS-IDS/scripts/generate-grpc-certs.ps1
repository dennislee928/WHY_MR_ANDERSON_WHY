# Pandora Box Console - gRPC TLS 證書生成腳本 (PowerShell)
# 為所有 gRPC 服務生成 TLS 證書

$ErrorActionPreference = "Stop"

$CERT_DIR = "certs"
$DAYS_VALID = 365

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "  🔐 gRPC TLS 證書生成工具" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

# 創建證書目錄
if (-not (Test-Path $CERT_DIR)) {
    New-Item -ItemType Directory -Path $CERT_DIR | Out-Null
}
Set-Location $CERT_DIR

# 1. 生成 CA 證書（根證書）
Write-Host "`n📜 步驟 1/4: 生成 CA 根證書..." -ForegroundColor Yellow
openssl req -x509 -newkey rsa:4096 -days $DAYS_VALID -nodes `
  -keyout ca-key.pem `
  -out ca-cert.pem `
  -subj "/C=TW/ST=Taipei/L=Taipei/O=Pandora Security/OU=Security/CN=Pandora CA"

Write-Host "✅ CA 證書已生成: ca-cert.pem, ca-key.pem" -ForegroundColor Green

# 2. 生成服務證書（Device Service）
Write-Host "`n📜 步驟 2/4: 生成 Device Service 證書..." -ForegroundColor Yellow
openssl req -newkey rsa:4096 -nodes `
  -keyout device-service-key.pem `
  -out device-service-req.pem `
  -subj "/C=TW/ST=Taipei/L=Taipei/O=Pandora Security/OU=Device Service/CN=device-service"

# 創建 SAN 配置文件
"subjectAltName=DNS:device-service,DNS:localhost,IP:127.0.0.1" | Out-File -FilePath device-san.cnf -Encoding ASCII

openssl x509 -req -in device-service-req.pem -days $DAYS_VALID `
  -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial `
  -out device-service-cert.pem `
  -extfile device-san.cnf

Write-Host "✅ Device Service 證書已生成" -ForegroundColor Green

# 3. 生成服務證書（Network Service）
Write-Host "`n📜 步驟 3/4: 生成 Network Service 證書..." -ForegroundColor Yellow
openssl req -newkey rsa:4096 -nodes `
  -keyout network-service-key.pem `
  -out network-service-req.pem `
  -subj "/C=TW/ST=Taipei/L=Taipei/O=Pandora Security/OU=Network Service/CN=network-service"

"subjectAltName=DNS:network-service,DNS:localhost,IP:127.0.0.1" | Out-File -FilePath network-san.cnf -Encoding ASCII

openssl x509 -req -in network-service-req.pem -days $DAYS_VALID `
  -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial `
  -out network-service-cert.pem `
  -extfile network-san.cnf

Write-Host "✅ Network Service 證書已生成" -ForegroundColor Green

# 4. 生成服務證書（Control Service）
Write-Host "`n📜 步驟 4/4: 生成 Control Service 證書..." -ForegroundColor Yellow
openssl req -newkey rsa:4096 -nodes `
  -keyout control-service-key.pem `
  -out control-service-req.pem `
  -subj "/C=TW/ST=Taipei/L=Taipei/O=Pandora Security/OU=Control Service/CN=control-service"

"subjectAltName=DNS:control-service,DNS:localhost,IP:127.0.0.1" | Out-File -FilePath control-san.cnf -Encoding ASCII

openssl x509 -req -in control-service-req.pem -days $DAYS_VALID `
  -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial `
  -out control-service-cert.pem `
  -extfile control-san.cnf

Write-Host "✅ Control Service 證書已生成" -ForegroundColor Green

# 清理臨時文件
Remove-Item *.req.pem, *.cnf -ErrorAction SilentlyContinue

Write-Host "`n=========================================" -ForegroundColor Cyan
Write-Host "  ✅ 所有 gRPC TLS 證書已生成！" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "`n📁 證書位置: $(Get-Location)" -ForegroundColor Yellow
Write-Host "`n📋 生成的證書:" -ForegroundColor Yellow
Get-ChildItem *.pem | Format-Table Name, Length, LastWriteTime

Write-Host "`n🔐 使用方式:" -ForegroundColor Cyan
Write-Host "  1. 複製證書到 configs/certs/"
Write-Host "  2. 更新 gRPC 服務器配置"
Write-Host "  3. 更新 gRPC 客戶端配置"
Write-Host "  4. 重新構建並重啟服務"
Write-Host "`n📖 詳細文檔: docs/GRPC-TLS-SETUP.md`n"

Set-Location ..

