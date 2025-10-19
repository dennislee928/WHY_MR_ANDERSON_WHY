package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// MTLSManager mTLS 管理器
type MTLSManager struct {
	logger    *logrus.Logger
	certDir   string
	caConfig  *CAConfig
	tlsConfig *tls.Config
}

// CAConfig CA 憑證配置
type CAConfig struct {
	Country      []string
	Province     []string
	City         []string
	Organization []string
	OrgUnit      []string
	CommonName   string
	ValidDays    int
}

// CertConfig 憑證配置
type CertConfig struct {
	CommonName   string
	Organization []string
	Country      []string
	Province     []string
	City         []string
	IPAddresses  []net.IP
	DNSNames     []string
	ValidDays    int
	IsServer     bool
	IsClient     bool
}

// NewMTLSManager 建立新的 mTLS 管理器
func NewMTLSManager(logger *logrus.Logger, certDir string) *MTLSManager {
	return &MTLSManager{
		logger:  logger,
		certDir: certDir,
		caConfig: &CAConfig{
			Country:      []string{"TW"},
			Province:     []string{"Taipei"},
			City:         []string{"Taipei"},
			Organization: []string{"Pandora Box Console"},
			OrgUnit:      []string{"Security Team"},
			CommonName:   "Pandora CA",
			ValidDays:    3650, // 10年
		},
	}
}

// Initialize 初始化 mTLS 管理器
func (m *MTLSManager) Initialize() error {
	// 建立憑證目錄
	if err := os.MkdirAll(m.certDir, 0755); err != nil {
		return fmt.Errorf("建立憑證目錄失敗: %v", err)
	}

	// 檢查並建立 CA 憑證
	caKeyPath := filepath.Join(m.certDir, "ca.key")
	caCertPath := filepath.Join(m.certDir, "ca.crt")

	if !m.fileExists(caKeyPath) || !m.fileExists(caCertPath) {
		m.logger.Info("建立 CA 憑證...")
		if err := m.generateCA(); err != nil {
			return fmt.Errorf("建立 CA 憑證失敗: %v", err)
		}
	}

	// 載入 CA 憑證
	if err := m.loadCA(); err != nil {
		return fmt.Errorf("載入 CA 憑證失敗: %v", err)
	}

	m.logger.Info("mTLS 管理器初始化完成")
	return nil
}

// generateCA 產生 CA 憑證
func (m *MTLSManager) generateCA() error {
	// 產生私鑰
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("產生 CA 私鑰失敗: %v", err)
	}

	// 建立憑證範本
	caTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:            m.caConfig.Country,
			Province:           m.caConfig.Province,
			Locality:           m.caConfig.City,
			Organization:       m.caConfig.Organization,
			OrganizationalUnit: m.caConfig.OrgUnit,
			CommonName:         m.caConfig.CommonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, m.caConfig.ValidDays),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
	}

	// 自簽名 CA 憑證
	caCertDER, err := x509.CreateCertificate(rand.Reader, &caTemplate, &caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("建立 CA 憑證失敗: %v", err)
	}

	// 儲存 CA 私鑰
	caKeyPath := filepath.Join(m.certDir, "ca.key")
	caKeyFile, err := os.Create(caKeyPath)
	if err != nil {
		return fmt.Errorf("建立 CA 私鑰檔案失敗: %v", err)
	}
	defer caKeyFile.Close()

	caKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caKey),
	}
	if err := pem.Encode(caKeyFile, caKeyPEM); err != nil {
		return fmt.Errorf("編碼 CA 私鑰失敗: %v", err)
	}

	// 設定私鑰檔案權限
	if err := os.Chmod(caKeyPath, 0600); err != nil {
		return fmt.Errorf("設定 CA 私鑰權限失敗: %v", err)
	}

	// 儲存 CA 憑證
	caCertPath := filepath.Join(m.certDir, "ca.crt")
	caCertFile, err := os.Create(caCertPath)
	if err != nil {
		return fmt.Errorf("建立 CA 憑證檔案失敗: %v", err)
	}
	defer caCertFile.Close()

	caCertPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caCertDER,
	}
	if err := pem.Encode(caCertFile, caCertPEM); err != nil {
		return fmt.Errorf("編碼 CA 憑證失敗: %v", err)
	}

	m.logger.Info("CA 憑證建立完成")
	return nil
}

// loadCA 載入 CA 憑證
func (m *MTLSManager) loadCA() error {
	caCertPath := filepath.Join(m.certDir, "ca.crt")
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("讀取 CA 憑證失敗: %v", err)
	}

	caCertBlock, _ := pem.Decode(caCertPEM)
	if caCertBlock == nil {
		return fmt.Errorf("解析 CA 憑證 PEM 失敗")
	}

	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return fmt.Errorf("解析 CA 憑證失敗: %v", err)
	}

	// 建立憑證池
	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(caCert)

	// 設定 TLS 配置
	m.tlsConfig = &tls.Config{
		ClientCAs:  caCertPool,
		RootCAs:    caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		},
	}

	return nil
}

// GenerateServerCert 產生伺服器憑證
func (m *MTLSManager) GenerateServerCert(config *CertConfig) error {
	return m.generateCert(config, "server")
}

// GenerateClientCert 產生客戶端憑證
func (m *MTLSManager) GenerateClientCert(config *CertConfig) error {
	return m.generateCert(config, "client")
}

// generateCert 產生憑證
func (m *MTLSManager) generateCert(config *CertConfig, certType string) error {
	// 載入 CA 私鑰
	caKeyPath := filepath.Join(m.certDir, "ca.key")
	caKeyPEM, err := os.ReadFile(caKeyPath)
	if err != nil {
		return fmt.Errorf("讀取 CA 私鑰失敗: %v", err)
	}

	caKeyBlock, _ := pem.Decode(caKeyPEM)
	if caKeyBlock == nil {
		return fmt.Errorf("解析 CA 私鑰 PEM 失敗")
	}

	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("解析 CA 私鑰失敗: %v", err)
	}

	// 載入 CA 憑證
	caCertPath := filepath.Join(m.certDir, "ca.crt")
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("讀取 CA 憑證失敗: %v", err)
	}

	caCertBlock, _ := pem.Decode(caCertPEM)
	if caCertBlock == nil {
		return fmt.Errorf("解析 CA 憑證 PEM 失敗")
	}

	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return fmt.Errorf("解析 CA 憑證失敗: %v", err)
	}

	// 產生憑證私鑰
	certKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("產生憑證私鑰失敗: %v", err)
	}

	// 建立憑證範本
	serialNumber, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:      config.Country,
			Province:     config.Province,
			Locality:     config.City,
			Organization: config.Organization,
			CommonName:   config.CommonName,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(0, 0, config.ValidDays),
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		IPAddresses: config.IPAddresses,
		DNSNames:    config.DNSNames,
	}

	// 設定憑證用途
	if config.IsServer {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
	}
	if config.IsClient {
		certTemplate.ExtKeyUsage = append(certTemplate.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
	}

	// 建立憑證
	certDER, err := x509.CreateCertificate(rand.Reader, &certTemplate, caCert, &certKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("建立憑證失敗: %v", err)
	}

	// 儲存憑證私鑰
	keyPath := filepath.Join(m.certDir, fmt.Sprintf("%s.key", certType))
	keyFile, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("建立憑證私鑰檔案失敗: %v", err)
	}
	defer keyFile.Close()

	keyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certKey),
	}
	if err := pem.Encode(keyFile, keyPEM); err != nil {
		return fmt.Errorf("編碼憑證私鑰失敗: %v", err)
	}

	// 設定私鑰檔案權限
	if err := os.Chmod(keyPath, 0600); err != nil {
		return fmt.Errorf("設定憑證私鑰權限失敗: %v", err)
	}

	// 儲存憑證
	certPath := filepath.Join(m.certDir, fmt.Sprintf("%s.crt", certType))
	certFile, err := os.Create(certPath)
	if err != nil {
		return fmt.Errorf("建立憑證檔案失敗: %v", err)
	}
	defer certFile.Close()

	certPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	}
	if err := pem.Encode(certFile, certPEM); err != nil {
		return fmt.Errorf("編碼憑證失敗: %v", err)
	}

	m.logger.Infof("%s 憑證建立完成: %s", certType, config.CommonName)
	return nil
}

// GetTLSConfig 取得 TLS 配置
func (m *MTLSManager) GetTLSConfig() *tls.Config {
	return m.tlsConfig.Clone()
}

// GetServerTLSConfig 取得伺服器 TLS 配置
func (m *MTLSManager) GetServerTLSConfig() (*tls.Config, error) {
	serverCert, err := tls.LoadX509KeyPair(
		filepath.Join(m.certDir, "server.crt"),
		filepath.Join(m.certDir, "server.key"),
	)
	if err != nil {
		return nil, fmt.Errorf("載入伺服器憑證失敗: %v", err)
	}

	config := m.tlsConfig.Clone()
	config.Certificates = []tls.Certificate{serverCert}
	return config, nil
}

// GetClientTLSConfig 取得客戶端 TLS 配置
func (m *MTLSManager) GetClientTLSConfig() (*tls.Config, error) {
	clientCert, err := tls.LoadX509KeyPair(
		filepath.Join(m.certDir, "client.crt"),
		filepath.Join(m.certDir, "client.key"),
	)
	if err != nil {
		return nil, fmt.Errorf("載入客戶端憑證失敗: %v", err)
	}

	config := m.tlsConfig.Clone()
	config.Certificates = []tls.Certificate{clientCert}
	return config, nil
}

// VerifyCertificate 驗證憑證
func (m *MTLSManager) VerifyCertificate(certPath string) error {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("讀取憑證失敗: %v", err)
	}

	certBlock, _ := pem.Decode(certPEM)
	if certBlock == nil {
		return fmt.Errorf("解析憑證 PEM 失敗")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return fmt.Errorf("解析憑證失敗: %v", err)
	}

	// 載入 CA 憑證池
	caCertPath := filepath.Join(m.certDir, "ca.crt")
	caCertPEM, err := os.ReadFile(caCertPath)
	if err != nil {
		return fmt.Errorf("讀取 CA 憑證失敗: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertPEM)

	// 驗證憑證
	opts := x509.VerifyOptions{
		Roots: caCertPool,
	}

	_, err = cert.Verify(opts)
	if err != nil {
		return fmt.Errorf("憑證驗證失敗: %v", err)
	}

	// 檢查憑證有效期
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return fmt.Errorf("憑證尚未生效")
	}
	if now.After(cert.NotAfter) {
		return fmt.Errorf("憑證已過期")
	}

	m.logger.Infof("憑證驗證成功: %s", cert.Subject.CommonName)
	return nil
}

// RenewCertificate 更新憑證
func (m *MTLSManager) RenewCertificate(certType string, config *CertConfig) error {
	// 備份舊憑證
	oldCertPath := filepath.Join(m.certDir, fmt.Sprintf("%s.crt", certType))
	oldKeyPath := filepath.Join(m.certDir, fmt.Sprintf("%s.key", certType))
	backupCertPath := fmt.Sprintf("%s.backup.%d", oldCertPath, time.Now().Unix())
	backupKeyPath := fmt.Sprintf("%s.backup.%d", oldKeyPath, time.Now().Unix())

	if m.fileExists(oldCertPath) {
		if err := os.Rename(oldCertPath, backupCertPath); err != nil {
			return fmt.Errorf("備份舊憑證失敗: %v", err)
		}
	}

	if m.fileExists(oldKeyPath) {
		if err := os.Rename(oldKeyPath, backupKeyPath); err != nil {
			return fmt.Errorf("備份舊私鑰失敗: %v", err)
		}
	}

	// 產生新憑證
	if err := m.generateCert(config, certType); err != nil {
		// 恢復備份
		if renameErr := os.Rename(backupCertPath, oldCertPath); renameErr != nil {
			m.logger.Errorf("恢復憑證備份失敗: %v", renameErr)
		}
		if renameErr := os.Rename(backupKeyPath, oldKeyPath); renameErr != nil {
			m.logger.Errorf("恢復金鑰備份失敗: %v", renameErr)
		}
		return fmt.Errorf("產生新憑證失敗: %v", err)
	}

	// 刪除備份
	os.Remove(backupCertPath)
	os.Remove(backupKeyPath)

	m.logger.Infof("憑證更新完成: %s", certType)
	return nil
}

// fileExists 檢查檔案是否存在
func (m *MTLSManager) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// SetupDefaultCertificates 設定預設憑證
func (m *MTLSManager) SetupDefaultCertificates() error {
	// 伺服器憑證配置
	serverConfig := &CertConfig{
		CommonName:   "pandora-server",
		Organization: []string{"Pandora Box Console"},
		Country:      []string{"TW"},
		Province:     []string{"Taipei"},
		City:         []string{"Taipei"},
		IPAddresses: []net.IP{
			net.ParseIP("127.0.0.1"),
			net.ParseIP("::1"),
		},
		DNSNames: []string{
			"localhost",
			"pandora-server",
			"pandora-agent",
			"axiom-ui",
		},
		ValidDays: 365,
		IsServer:  true,
	}

	// 客戶端憑證配置
	clientConfig := &CertConfig{
		CommonName:   "pandora-client",
		Organization: []string{"Pandora Box Console"},
		Country:      []string{"TW"},
		Province:     []string{"Taipei"},
		City:         []string{"Taipei"},
		ValidDays:    365,
		IsClient:     true,
	}

	// 產生憑證
	if err := m.GenerateServerCert(serverConfig); err != nil {
		return fmt.Errorf("產生伺服器憑證失敗: %v", err)
	}

	if err := m.GenerateClientCert(clientConfig); err != nil {
		return fmt.Errorf("產生客戶端憑證失敗: %v", err)
	}

	return nil
}
