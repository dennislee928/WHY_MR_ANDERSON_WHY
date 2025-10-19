package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

// TLSConfig contains TLS configuration for gRPC
type TLSConfig struct {
	// Server configuration
	ServerCertFile string
	ServerKeyFile  string
	ClientCAFile   string

	// Client configuration
	ClientCertFile string
	ClientKeyFile  string
	ServerCAFile   string

	// Server name for verification
	ServerName string
}

// LoadServerTLSCredentials loads TLS credentials for gRPC server
// 載入 gRPC 服務端 TLS 憑證
func LoadServerTLSCredentials(config *TLSConfig) (credentials.TransportCredentials, error) {
	// 載入服務端憑證和私鑰
	serverCert, err := tls.LoadX509KeyPair(config.ServerCertFile, config.ServerKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load server cert: %w", err)
	}

	// 載入客戶端 CA 憑證（用於驗證客戶端）
	clientCAPool := x509.NewCertPool()
	clientCA, err := os.ReadFile(config.ClientCAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client CA: %w", err)
	}

	if !clientCAPool.AppendCertsFromPEM(clientCA) {
		return nil, fmt.Errorf("failed to append client CA")
	}

	// 配置 TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCAPool,
		MinVersion:   tls.VersionTLS13,
	}

	return credentials.NewTLS(tlsConfig), nil
}

// LoadClientTLSCredentials loads TLS credentials for gRPC client
// 載入 gRPC 客戶端 TLS 憑證
func LoadClientTLSCredentials(config *TLSConfig) (credentials.TransportCredentials, error) {
	// 載入客戶端憑證和私鑰
	clientCert, err := tls.LoadX509KeyPair(config.ClientCertFile, config.ClientKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert: %w", err)
	}

	// 載入服務端 CA 憑證（用於驗證服務端）
	serverCAPool := x509.NewCertPool()
	serverCA, err := os.ReadFile(config.ServerCAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read server CA: %w", err)
	}

	if !serverCAPool.AppendCertsFromPEM(serverCA) {
		return nil, fmt.Errorf("failed to append server CA")
	}

	// 配置 TLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      serverCAPool,
		ServerName:   config.ServerName,
		MinVersion:   tls.VersionTLS13,
	}

	return credentials.NewTLS(tlsConfig), nil
}

// GenerateSelfSignedCert generates a self-signed certificate (for development)
// 生成自簽名憑證（僅用於開發環境）
func GenerateSelfSignedCert(certFile, keyFile string) error {
	// TODO: 實現自簽名憑證生成
	// 使用 crypto/x509 和 crypto/rsa 生成憑證
	return fmt.Errorf("not implemented")
}

