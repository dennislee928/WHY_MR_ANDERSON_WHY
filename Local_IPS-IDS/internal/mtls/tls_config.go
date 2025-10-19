package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TLSConfig 包含 TLS 配置
type TLSConfig struct {
	CACertFile     string
	ServerCertFile string
	ServerKeyFile  string
	ClientCertFile string
	ClientKeyFile  string
}

// LoadServerTLSCredentials 載入服務器 TLS 憑證
func LoadServerTLSCredentials(config TLSConfig) (credentials.TransportCredentials, error) {
	// 載入服務器證書和私鑰
	serverCert, err := tls.LoadX509KeyPair(config.ServerCertFile, config.ServerKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}

	// 載入 CA 證書
	caCert, err := os.ReadFile(config.CACertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	// 創建證書池
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add CA certificate to pool")
	}

	// 配置 TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS13, // 強制使用 TLS 1.3
		CipherSuites: []uint16{
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}

	return credentials.NewTLS(tlsConfig), nil
}

// LoadClientTLSCredentials 載入客戶端 TLS 憑證
func LoadClientTLSCredentials(config TLSConfig) (credentials.TransportCredentials, error) {
	// 載入客戶端證書和私鑰
	clientCert, err := tls.LoadX509KeyPair(config.ClientCertFile, config.ClientKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %w", err)
	}

	// 載入 CA 證書
	caCert, err := os.ReadFile(config.CACertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %w", err)
	}

	// 創建證書池
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to add CA certificate to pool")
	}

	// 配置 TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS13, // 強制使用 TLS 1.3
		CipherSuites: []uint16{
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}

	return credentials.NewTLS(tlsConfig), nil
}

// NewServerWithTLS 創建帶 TLS 的 gRPC 服務器
func NewServerWithTLS(config TLSConfig) (*grpc.Server, error) {
	creds, err := LoadServerTLSCredentials(config)
	if err != nil {
		return nil, err
	}

	return grpc.NewServer(grpc.Creds(creds)), nil
}

// DialWithTLS 使用 TLS 連接到 gRPC 服務器
func DialWithTLS(address string, config TLSConfig) (*grpc.ClientConn, error) {
	creds, err := LoadClientTLSCredentials(config)
	if err != nil {
		return nil, err
	}

	return grpc.Dial(address, grpc.WithTransportCredentials(creds))
}

// GetTLSConfigFromEnv 從環境變數獲取 TLS 配置
func GetTLSConfigFromEnv() TLSConfig {
	return TLSConfig{
		CACertFile:     getEnv("GRPC_CA_CERT", "configs/certs/ca-cert.pem"),
		ServerCertFile: getEnv("GRPC_SERVER_CERT", "configs/certs/server-cert.pem"),
		ServerKeyFile:  getEnv("GRPC_SERVER_KEY", "configs/certs/server-key.pem"),
		ClientCertFile: getEnv("GRPC_CLIENT_CERT", "configs/certs/client-cert.pem"),
		ClientKeyFile:  getEnv("GRPC_CLIENT_KEY", "configs/certs/client-key.pem"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

