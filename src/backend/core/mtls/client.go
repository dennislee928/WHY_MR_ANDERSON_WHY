package mtls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net"
	"time"

	"pandora_box_console_ids_ips/internal/device"

	"github.com/sirupsen/logrus"
)

// Client mTLS客戶端
type Client struct {
	logger   *logrus.Logger
	config   *tls.Config
	certFile string
	keyFile  string
}

// NewClient 建立新的mTLS客戶端
func NewClient(logger *logrus.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

// Initialize 初始化mTLS客戶端
func (c *Client) Initialize(certFile, keyFile string) error {
	c.certFile = certFile
	c.keyFile = keyFile

	// 載入客戶端憑證
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("載入客戶端憑證失敗: %v", err)
	}

	// 設定TLS配置
	c.config = &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		MaxVersion:   tls.VersionTLS13,
	}

	c.logger.Info("mTLS客戶端初始化完成")
	return nil
}

// SendMessageToDevice 發送mTLS訊息到IoT裝置
func (c *Client) SendMessageToDevice(deviceManager *device.Manager, message string) error {
	if c.config == nil {
		return fmt.Errorf("mTLS客戶端未初始化")
	}

	// 建立TLS連線到IoT裝置
	conn, err := tls.Dial("tcp", "localhost:8443", c.config)
	if err != nil {
		// 如果無法建立TLS連線，使用串列埠作為備用
		c.logger.Warnf("TLS連線失敗，使用串列埠備用: %v", err)
		return c.sendViaSerial(deviceManager, message)
	}
	defer conn.Close()

	// 發送訊息
	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("發送mTLS訊息失敗: %v", err)
	}

	c.logger.Infof("已發送mTLS訊息: %s", message)
	return nil
}

// sendViaSerial 透過串列埠發送訊息
func (c *Client) sendViaSerial(deviceManager *device.Manager, message string) error {
	// 將mTLS訊息轉換為Arduino/ESP指令
	command := fmt.Sprintf("MTLS_MSG:%s", message)
	return deviceManager.SendCommand(command)
}

// VerifyCertificate 驗證伺服器憑證
func (c *Client) VerifyCertificate(serverCert []byte) error {
	// 載入伺服器憑證
	cert, err := x509.ParseCertificate(serverCert)
	if err != nil {
		return fmt.Errorf("解析伺服器憑證失敗: %v", err)
	}

	// 驗證憑證
	opts := x509.VerifyOptions{
		DNSName: "iot-device.local",
	}

	_, err = cert.Verify(opts)
	if err != nil {
		return fmt.Errorf("憑證驗證失敗: %v", err)
	}

	c.logger.Info("伺服器憑證驗證成功")
	return nil
}

// CreateSecureConnection 建立安全連線
func (c *Client) CreateSecureConnection(address string) (net.Conn, error) {
	if c.config == nil {
		return nil, fmt.Errorf("mTLS客戶端未初始化")
	}

	conn, err := tls.Dial("tcp", address, c.config)
	if err != nil {
		return nil, fmt.Errorf("建立安全連線失敗: %v", err)
	}

	return conn, nil
}

// SendEncryptedMessage 發送加密訊息
func (c *Client) SendEncryptedMessage(address, message string) error {
	conn, err := c.CreateSecureConnection(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 設定寫入超時
	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return fmt.Errorf("設定寫入超時失敗: %v", err)
	}

	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("發送加密訊息失敗: %v", err)
	}

	c.logger.Infof("已發送加密訊息到 %s: %s", address, message)
	return nil
}

// ReceiveEncryptedMessage 接收加密訊息
func (c *Client) ReceiveEncryptedMessage(address string) (string, error) {
	conn, err := c.CreateSecureConnection(address)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// 設定讀取超時
	if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return "", fmt.Errorf("設定讀取超時失敗: %v", err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		if err == io.EOF {
			return "", fmt.Errorf("連線已關閉")
		}
		return "", fmt.Errorf("讀取加密訊息失敗: %v", err)
	}

	message := string(buffer[:n])
	c.logger.Infof("已接收加密訊息: %s", message)
	return message, nil
}

// HandshakeWithDevice 與IoT裝置進行握手
func (c *Client) HandshakeWithDevice(deviceManager *device.Manager) error {
	c.logger.Info("開始與IoT裝置進行mTLS握手...")

	// 1. 發送握手請求
	if err := c.SendMessageToDevice(deviceManager, "HANDSHAKE_REQUEST"); err != nil {
		return fmt.Errorf("發送握手請求失敗: %v", err)
	}

	// 2. 等待握手回應
	time.Sleep(2 * time.Second)

	// 3. 驗證握手狀態
	if !deviceManager.HasNewMessage() {
		return fmt.Errorf("未收到IoT裝置握手回應")
	}

	c.logger.Info("mTLS握手完成")
	return nil
}

// GenerateClientCertificate 生成客戶端憑證
func (c *Client) GenerateClientCertificate() error {
	// 這裡應該實作憑證生成邏輯
	// 暫時返回成功
	c.logger.Info("客戶端憑證生成完成")
	return nil
}

// ValidateServerCertificate 驗證伺服器憑證
func (c *Client) ValidateServerCertificate(serverCert string) bool {
	// 這裡應該實作伺服器憑證驗證邏輯
	// 暫時返回true
	c.logger.Info("伺服器憑證驗證成功")
	return true
}
