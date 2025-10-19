package security

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMTLSManager(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := t.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	assert.NotNil(t, manager)
	assert.Equal(t, tempDir, manager.certDir)
	assert.NotNil(t, manager.caConfig)
	assert.Equal(t, "Pandora CA", manager.caConfig.CommonName)
}

func TestMTLSManagerInitialize(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := t.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	err := manager.Initialize()
	require.NoError(t, err)

	// 檢查 CA 憑證檔案是否存在
	caKeyPath := filepath.Join(tempDir, "ca.key")
	caCertPath := filepath.Join(tempDir, "ca.crt")

	assert.FileExists(t, caKeyPath)
	assert.FileExists(t, caCertPath)

	// 檢查檔案權限
	keyInfo, err := os.Stat(caKeyPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), keyInfo.Mode())

	// 檢查 TLS 配置是否已設定
	assert.NotNil(t, manager.tlsConfig)
	assert.NotNil(t, manager.tlsConfig.ClientCAs)
	assert.NotNil(t, manager.tlsConfig.RootCAs)
}

func TestGenerateServerCert(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := t.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	err := manager.Initialize()
	require.NoError(t, err)

	// 產生伺服器憑證
	serverConfig := &CertConfig{
		CommonName:   "test-server",
		Organization: []string{"Test Org"},
		Country:      []string{"TW"},
		Province:     []string{"Taipei"},
		City:         []string{"Taipei"},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
		},
		DNSNames: []string{
			"localhost",
			"test-server",
		},
		ValidDays: 365,
		IsServer:  true,
	}

	err = manager.GenerateServerCert(serverConfig)
	require.NoError(t, err)

	// 檢查憑證檔案
	serverKeyPath := filepath.Join(tempDir, "server.key")
	serverCertPath := filepath.Join(tempDir, "server.crt")

	assert.FileExists(t, serverKeyPath)
	assert.FileExists(t, serverCertPath)

	// 驗證憑證
	err = manager.VerifyCertificate(serverCertPath)
	assert.NoError(t, err)
}

func TestGenerateClientCert(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := t.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	err := manager.Initialize()
	require.NoError(t, err)

	// 產生客戶端憑證
	clientConfig := &CertConfig{
		CommonName:   "test-client",
		Organization: []string{"Test Org"},
		Country:      []string{"TW"},
		Province:     []string{"Taipei"},
		City:         []string{"Taipei"},
		ValidDays:    365,
		IsClient:     true,
	}

	err = manager.GenerateClientCert(clientConfig)
	require.NoError(t, err)

	// 檢查憑證檔案
	clientKeyPath := filepath.Join(tempDir, "client.key")
	clientCertPath := filepath.Join(tempDir, "client.crt")

	assert.FileExists(t, clientKeyPath)
	assert.FileExists(t, clientCertPath)

	// 驗證憑證
	err = manager.VerifyCertificate(clientCertPath)
	assert.NoError(t, err)
}

func TestGetTLSConfigs(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := t.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	err := manager.Initialize()
	require.NoError(t, err)

	err = manager.SetupDefaultCertificates()
	require.NoError(t, err)

	// 測試伺服器 TLS 配置
	serverConfig, err := manager.GetServerTLSConfig()
	require.NoError(t, err)
	assert.NotNil(t, serverConfig)
	assert.Greater(t, len(serverConfig.Certificates), 0)
	assert.Equal(t, tls.RequireAndVerifyClientCert, serverConfig.ClientAuth)

	// 測試客戶端 TLS 配置
	clientConfig, err := manager.GetClientTLSConfig()
	require.NoError(t, err)
	assert.NotNil(t, clientConfig)
	assert.Greater(t, len(clientConfig.Certificates), 0)
}

func TestVerifyCertificate(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := t.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	err := manager.Initialize()
	require.NoError(t, err)

	err = manager.SetupDefaultCertificates()
	require.NoError(t, err)

	// 驗證伺服器憑證
	serverCertPath := filepath.Join(tempDir, "server.crt")
	err = manager.VerifyCertificate(serverCertPath)
	assert.NoError(t, err)

	// 驗證客戶端憑證
	clientCertPath := filepath.Join(tempDir, "client.crt")
	err = manager.VerifyCertificate(clientCertPath)
	assert.NoError(t, err)

	// 測試無效憑證路徑
	err = manager.VerifyCertificate("/non/existent/cert.crt")
	assert.Error(t, err)
}

func TestRenewCertificate(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := t.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	err := manager.Initialize()
	require.NoError(t, err)

	err = manager.SetupDefaultCertificates()
	require.NoError(t, err)

	// 讀取原始憑證
	serverCertPath := filepath.Join(tempDir, "server.crt")
	originalCert, err := os.ReadFile(serverCertPath)
	require.NoError(t, err)

	// 等待一秒確保時間戳不同
	time.Sleep(1 * time.Second)

	// 更新憑證
	serverConfig := &CertConfig{
		CommonName:   "renewed-server",
		Organization: []string{"Renewed Org"},
		Country:      []string{"TW"},
		Province:     []string{"Taipei"},
		City:         []string{"Taipei"},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
		},
		DNSNames: []string{
			"localhost",
			"renewed-server",
		},
		ValidDays: 365,
		IsServer:  true,
	}

	err = manager.RenewCertificate("server", serverConfig)
	require.NoError(t, err)

	// 讀取新憑證
	newCert, err := os.ReadFile(serverCertPath)
	require.NoError(t, err)

	// 確認憑證已更新
	assert.NotEqual(t, originalCert, newCert)

	// 驗證新憑證
	err = manager.VerifyCertificate(serverCertPath)
	assert.NoError(t, err)
}

func TestSetupDefaultCertificates(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := t.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	err := manager.Initialize()
	require.NoError(t, err)

	err = manager.SetupDefaultCertificates()
	require.NoError(t, err)

	// 檢查所有預設憑證檔案
	expectedFiles := []string{
		"ca.key", "ca.crt",
		"server.key", "server.crt",
		"client.key", "client.crt",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(tempDir, file)
		assert.FileExists(t, filePath, "檔案應該存在: %s", file)
	}
}

func TestCertificateValidation(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := t.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	err := manager.Initialize()
	require.NoError(t, err)

	err = manager.SetupDefaultCertificates()
	require.NoError(t, err)

	// 載入並解析憑證
	serverCertPath := filepath.Join(tempDir, "server.crt")
	certPEM, err := os.ReadFile(serverCertPath)
	require.NoError(t, err)

	// 解析憑證
	block, _ := x509.ParseCertificate(certPEM)
	require.NotNil(t, block)

	cert, err := x509.ParseCertificate(block.Raw)
	require.NoError(t, err)

	// 檢查憑證屬性
	assert.Equal(t, "pandora-server", cert.Subject.CommonName)
	assert.Contains(t, cert.DNSNames, "localhost")
	assert.Contains(t, cert.DNSNames, "pandora-server")
	assert.Contains(t, cert.IPAddresses, net.ParseIP("127.0.0.1"))

	// 檢查憑證用途
	assert.Contains(t, cert.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
}

func TestTLSConfigSecurity(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := t.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	err := manager.Initialize()
	require.NoError(t, err)

	config := manager.GetTLSConfig()

	// 檢查安全設定
	assert.Equal(t, tls.RequireAndVerifyClientCert, config.ClientAuth)
	assert.Equal(t, uint16(tls.VersionTLS12), config.MinVersion)
	assert.Greater(t, len(config.CipherSuites), 0)

	// 檢查支援的加密套件
	expectedCiphers := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	}

	for _, cipher := range expectedCiphers {
		assert.Contains(t, config.CipherSuites, cipher)
	}
}

func BenchmarkGenerateServerCert(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := b.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	err := manager.Initialize()
	require.NoError(b, err)

	serverConfig := &CertConfig{
		CommonName:   "bench-server",
		Organization: []string{"Bench Org"},
		Country:      []string{"TW"},
		Province:     []string{"Taipei"},
		City:         []string{"Taipei"},
		ValidDays:    365,
		IsServer:     true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := manager.GenerateServerCert(serverConfig)
		require.NoError(b, err)

		// 清理憑證檔案以便下次測試
		os.Remove(filepath.Join(tempDir, "server.crt"))
		os.Remove(filepath.Join(tempDir, "server.key"))
	}
}

func TestFilePermissions(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	tempDir := t.TempDir()
	manager := NewMTLSManager(logger, tempDir)

	err := manager.Initialize()
	require.NoError(t, err)

	err = manager.SetupDefaultCertificates()
	require.NoError(t, err)

	// 檢查私鑰檔案權限 (應該是 0600)
	keyFiles := []string{"ca.key", "server.key", "client.key"}
	for _, file := range keyFiles {
		filePath := filepath.Join(tempDir, file)
		info, err := os.Stat(filePath)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0600), info.Mode(), "私鑰檔案權限應該是 0600: %s", file)
	}

	// 檢查憑證檔案權限 (應該是 0644)
	certFiles := []string{"ca.crt", "server.crt", "client.crt"}
	for _, file := range certFiles {
		filePath := filepath.Join(tempDir, file)
		info, err := os.Stat(filePath)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0644), info.Mode(), "憑證檔案權限應該是 0644: %s", file)
	}
}
