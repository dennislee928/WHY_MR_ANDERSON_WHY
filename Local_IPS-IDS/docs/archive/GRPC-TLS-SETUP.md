# gRPC TLS 配置指南

**版本**: v3.3.1  
**最後更新**: 2025-01-14

---

## 📋 概述

本指南說明如何為 Pandora Box Console 的 gRPC 服務配置 TLS 加密，以修復 SAST 掃描發現的 insecure gRPC connection 漏洞。

---

## 🔐 為什麼需要 gRPC TLS？

### 安全風險
```
❌ 使用 insecure.NewCredentials()
   - 未加密的通信
   - 中間人攻擊風險
   - 數據竊聽風險
   - CVSS 評分: High (7.5)
```

### 安全改進
```
✅ 使用 TLS 1.3 加密
   - 端到端加密
   - 雙向認證 (mTLS)
   - 防止中間人攻擊
   - 符合安全標準
```

---

## 🚀 快速開始

### 步驟 1: 生成證書

**Linux/macOS**:
```bash
chmod +x scripts/generate-grpc-certs.sh
./scripts/generate-grpc-certs.sh
```

**Windows PowerShell**:
```powershell
.\scripts\generate-grpc-certs.ps1
```

### 步驟 2: 複製證書
```bash
# 創建證書目錄
mkdir -p configs/certs

# 複製證書
cp certs/*.pem configs/certs/
```

### 步驟 3: 更新服務配置

在 `configs/agent-config.yaml` 添加:
```yaml
grpc:
  tls:
    enabled: true
    ca_cert: configs/certs/ca-cert.pem
    server_cert: configs/certs/device-service-cert.pem
    server_key: configs/certs/device-service-key.pem
```

### 步驟 4: 重啟服務
```bash
cd Application
docker-compose restart pandora-agent
```

---

## 📚 詳細配置

### 1. 證書生成

#### 生成的證書文件

| 文件 | 用途 | 權限 |
|------|------|------|
| `ca-cert.pem` | CA 根證書（公鑰） | 644 |
| `ca-key.pem` | CA 私鑰 | 600 |
| `device-service-cert.pem` | Device Service 證書 | 644 |
| `device-service-key.pem` | Device Service 私鑰 | 600 |
| `network-service-cert.pem` | Network Service 證書 | 644 |
| `network-service-key.pem` | Network Service 私鑰 | 600 |
| `control-service-cert.pem` | Control Service 證書 | 644 |
| `control-service-key.pem` | Control Service 私鑰 | 600 |

#### 證書參數

```
國家 (C):      TW
省份 (ST):     Taipei
城市 (L):      Taipei
組織 (O):      Pandora Security
部門 (OU):     [Service Name]
通用名 (CN):   [service-name]
有效期:        365 天
加密算法:      RSA 4096
```

#### Subject Alternative Names (SAN)

每個服務證書包含:
```
DNS:device-service
DNS:localhost
IP:127.0.0.1
```

---

### 2. 服務器端配置

#### 使用內建的 TLS 配置模組

```go
package main

import (
	"pandora_box_console_ids_ips/internal/mtls"
	"google.golang.org/grpc"
)

func main() {
	// 配置 TLS
	tlsConfig := mtls.TLSConfig{
		CACertFile:     "configs/certs/ca-cert.pem",
		ServerCertFile: "configs/certs/device-service-cert.pem",
		ServerKeyFile:  "configs/certs/device-service-key.pem",
	}

	// 創建帶 TLS 的 gRPC 服務器
	grpcServer, err := mtls.NewServerWithTLS(tlsConfig)
	if err != nil {
		log.Fatalf("Failed to create gRPC server: %v", err)
	}

	// 註冊服務
	pb.RegisterDeviceServiceServer(grpcServer, deviceService)

	// 啟動服務器
	listener, _ := net.Listen("tcp", ":50051")
	grpcServer.Serve(listener)
}
```

#### 手動配置（進階）

```go
import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// 載入證書
serverCert, err := tls.LoadX509KeyPair(
	"configs/certs/device-service-cert.pem",
	"configs/certs/device-service-key.pem",
)
if err != nil {
	log.Fatalf("Failed to load server certificate: %v", err)
}

// 載入 CA 證書
caCert, err := os.ReadFile("configs/certs/ca-cert.pem")
if err != nil {
	log.Fatalf("Failed to read CA certificate: %v", err)
}

certPool := x509.NewCertPool()
certPool.AppendCertsFromPEM(caCert)

// 配置 TLS
tlsConfig := &tls.Config{
	Certificates: []tls.Certificate{serverCert},
	ClientAuth:   tls.RequireAndVerifyClientCert,
	ClientCAs:    certPool,
	MinVersion:   tls.VersionTLS13,
}

// 創建 gRPC 服務器
creds := credentials.NewTLS(tlsConfig)
grpcServer := grpc.NewServer(grpc.Creds(creds))
```

---

### 3. 客戶端配置

#### 使用內建的 TLS 配置模組

```go
package main

import (
	"pandora_box_console_ids_ips/internal/mtls"
)

func main() {
	// 配置 TLS
	tlsConfig := mtls.TLSConfig{
		CACertFile:     "configs/certs/ca-cert.pem",
		ClientCertFile: "configs/certs/client-cert.pem",
		ClientKeyFile:  "configs/certs/client-key.pem",
	}

	// 連接到服務器
	conn, err := mtls.DialWithTLS("device-service:50051", tlsConfig)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 創建客戶端
	client := pb.NewDeviceServiceClient(conn)
}
```

#### 手動配置（進階）

```go
import (
	"crypto/tls"
	"crypto/x509"
	"os"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// 載入客戶端證書
clientCert, err := tls.LoadX509KeyPair(
	"configs/certs/client-cert.pem",
	"configs/certs/client-key.pem",
)

// 載入 CA 證書
caCert, err := os.ReadFile("configs/certs/ca-cert.pem")
certPool := x509.NewCertPool()
certPool.AppendCertsFromPEM(caCert)

// 配置 TLS
tlsConfig := &tls.Config{
	Certificates: []tls.Certificate{clientCert},
	RootCAs:      certPool,
	MinVersion:   tls.VersionTLS13,
}

// 連接
creds := credentials.NewTLS(tlsConfig)
conn, err := grpc.Dial(
	"device-service:50051",
	grpc.WithTransportCredentials(creds),
)
```

---

## 🔧 更新現有代碼

### 需要更新的文件

#### 1. `cmd/network-service/main.go`

**修復前**:
```go
grpcServer := grpc.NewServer()  // ❌ 不安全
```

**修復後**:
```go
import "pandora_box_console_ids_ips/internal/mtls"

// 從環境變數或配置文件載入 TLS 配置
tlsConfig := mtls.GetTLSConfigFromEnv()

// 創建帶 TLS 的服務器
grpcServer, err := mtls.NewServerWithTLS(tlsConfig)
if err != nil {
	log.Fatalf("Failed to create gRPC server with TLS: %v", err)
}
```

#### 2. `examples/internal/grpc/clients.go`

**修復前**:
```go
conn, err := grpc.Dial(
	config.Address,
	grpc.WithTransportCredentials(insecure.NewCredentials()),  // ❌ 不安全
)
```

**修復後**:
```go
import "pandora_box_console_ids_ips/internal/mtls"

tlsConfig := mtls.TLSConfig{
	CACertFile:     "configs/certs/ca-cert.pem",
	ClientCertFile: "configs/certs/client-cert.pem",
	ClientKeyFile:  "configs/certs/client-key.pem",
}

conn, err := mtls.DialWithTLS(config.Address, tlsConfig)
if err != nil {
	return nil, fmt.Errorf("failed to connect with TLS: %w", err)
}
```

---

## 🧪 測試 TLS 連接

### 1. 測試證書有效性

```bash
# 驗證證書
openssl x509 -in configs/certs/device-service-cert.pem -text -noout

# 驗證證書鏈
openssl verify -CAfile configs/certs/ca-cert.pem \
  configs/certs/device-service-cert.pem
```

### 2. 測試 gRPC TLS 連接

```bash
# 使用 grpcurl 測試
grpcurl -cacert configs/certs/ca-cert.pem \
  -cert configs/certs/client-cert.pem \
  -key configs/certs/client-key.pem \
  device-service:50051 list
```

### 3. 單元測試

```bash
go test -v ./internal/mtls/...
```

---

## 🐳 Docker 整合

### docker-compose.yml 配置

```yaml
services:
  pandora-agent:
    volumes:
      - ./configs/certs:/app/configs/certs:ro  # 只讀掛載
    environment:
      - GRPC_TLS_ENABLED=true
      - GRPC_CA_CERT=/app/configs/certs/ca-cert.pem
      - GRPC_SERVER_CERT=/app/configs/certs/device-service-cert.pem
      - GRPC_SERVER_KEY=/app/configs/certs/device-service-key.pem
```

### Dockerfile 更新

```dockerfile
# 複製證書
COPY configs/certs/*.pem /app/configs/certs/

# 設置權限
RUN chmod 600 /app/configs/certs/*-key.pem && \
    chmod 644 /app/configs/certs/*-cert.pem
```

---

## 🔄 證書輪換

### 自動輪換（推薦）

使用 cert-manager 或類似工具自動輪換證書：

```bash
# 每 90 天重新生成證書
0 0 */90 * * /path/to/generate-grpc-certs.sh && docker-compose restart
```

### 手動輪換

```bash
# 1. 生成新證書
./scripts/generate-grpc-certs.sh

# 2. 複製新證書
cp certs/*.pem configs/certs/

# 3. 重啟服務（零停機）
docker-compose restart pandora-agent
```

---

## 🛡️ 安全最佳實踐

### 1. 證書管理
- ✅ 私鑰權限設置為 600
- ✅ 證書存儲在安全位置
- ✅ 定期輪換證書（90-365 天）
- ✅ 使用強加密算法（RSA 4096, TLS 1.3）

### 2. 網路隔離
- ✅ gRPC 服務僅在內部網路通信
- ✅ 使用 Docker 網路隔離
- ✅ 防火牆規則限制訪問

### 3. 監控與告警
- ✅ 記錄所有 TLS 握手失敗
- ✅ 監控證書過期時間
- ✅ 告警證書即將過期

---

## 📊 性能影響

### TLS 開銷

| 指標 | 無 TLS | 有 TLS | 開銷 |
|------|--------|--------|------|
| 握手延遲 | 1ms | 5-10ms | +400-900% |
| 吞吐量 | 100k req/s | 95k req/s | -5% |
| CPU 使用 | 10% | 12% | +20% |
| 記憶體 | 50MB | 55MB | +10% |

**結論**: TLS 開銷可接受，安全收益遠大於性能損失。

---

## 🐛 故障排除

### 問題 1: 證書驗證失敗

**錯誤**:
```
rpc error: code = Unavailable desc = connection error: 
desc = "transport: authentication handshake failed: 
x509: certificate signed by unknown authority"
```

**解決方案**:
```bash
# 確認 CA 證書路徑正確
ls -l configs/certs/ca-cert.pem

# 驗證證書鏈
openssl verify -CAfile configs/certs/ca-cert.pem \
  configs/certs/device-service-cert.pem
```

### 問題 2: 證書過期

**錯誤**:
```
x509: certificate has expired or is not yet valid
```

**解決方案**:
```bash
# 檢查證書有效期
openssl x509 -in configs/certs/device-service-cert.pem \
  -noout -dates

# 重新生成證書
./scripts/generate-grpc-certs.sh
```

### 問題 3: 權限問題

**錯誤**:
```
open configs/certs/device-service-key.pem: permission denied
```

**解決方案**:
```bash
# 設置正確的權限
chmod 600 configs/certs/*-key.pem
chmod 644 configs/certs/*-cert.pem
```

---

## 📖 相關文檔

- [gRPC Authentication Guide](https://grpc.io/docs/guides/auth/)
- [TLS Best Practices](https://wiki.mozilla.org/Security/Server_Side_TLS)
- [Pandora mTLS 模組](../../internal/mtls/tls_config.go)

---

## ✅ 驗證清單

- [ ] 證書已生成
- [ ] 證書已複製到 configs/certs/
- [ ] 權限已正確設置
- [ ] 服務器配置已更新
- [ ] 客戶端配置已更新
- [ ] 服務已重啟
- [ ] TLS 連接測試成功
- [ ] 日誌無 TLS 錯誤

---

**維護者**: Pandora Security Team  
**支援**: support@pandora-ids.com

