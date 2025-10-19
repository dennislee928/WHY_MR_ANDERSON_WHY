package device

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Manager 管理Arduino/ESP IoT裝置
type Manager struct {
	logger      *logrus.Logger
	port        *SerialDevice
	portName    string
	isConnected bool
	lastMessage time.Time
	lastPinCode string
	mu          sync.RWMutex
}

// NewManager 建立新的裝置管理器
func NewManager(logger *logrus.Logger) *Manager {
	return &Manager{
		logger: logger,
	}
}

// Initialize 初始化Arduino/ESP裝置連線
func (m *Manager) Initialize(portName string) error {
	m.portName = portName

	// 建立序列埠裝置
	m.port = NewSerialDevice(portName, 115200, time.Second*1, m.logger)

	// 連接到裝置
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.port.Connect(ctx); err != nil {
		return fmt.Errorf("無法開啟串列埠 %s: %v", portName, err)
	}

	m.isConnected = true
	m.lastMessage = time.Now()

	m.logger.Infof("成功連線到Arduino/ESP IoT裝置: %s", portName)

	// 啟動訊息監聽
	go m.listenForMessages()

	return nil
}

// listenForMessages 監聽Arduino/ESP裝置訊息
func (m *Manager) listenForMessages() {
	ctx := context.Background()

	for m.isConnected {
		data, err := m.port.Read(ctx)
		if err != nil {
			if err == context.Canceled || err == context.DeadlineExceeded {
				continue
			}
			m.logger.WithError(err).Error("讀取序列埠資料時發生錯誤")
			time.Sleep(100 * time.Millisecond)
			continue
		}

		message := strings.TrimSpace(string(data))
		if message != "" {
			m.mu.Lock()
			m.lastMessage = time.Now()
			m.mu.Unlock()

			m.logger.Debugf("收到IoT裝置訊息: %s", message)

			// 解析特殊訊息
			m.parseSpecialMessage(message)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// parseSpecialMessage 解析特殊訊息
func (m *Manager) parseSpecialMessage(message string) {
	// 解析PIN碼回應: PIN:123
	if strings.HasPrefix(message, "PIN:") {
		pinCode := strings.TrimPrefix(message, "PIN:")
		m.mu.Lock()
		m.lastPinCode = pinCode
		m.mu.Unlock()
		m.logger.Infof("收到PIN碼: %s", pinCode)
	}

	// 解析狀態回應: STATUS:READY
	if strings.HasPrefix(message, "STATUS:") {
		status := strings.TrimPrefix(message, "STATUS:")
		m.logger.Infof("裝置狀態: %s", status)
	}

	// 解析錯誤回應: ERROR:NO_DISPLAY
	if strings.HasPrefix(message, "ERROR:") {
		errorMsg := strings.TrimPrefix(message, "ERROR:")
		m.logger.Errorf("裝置錯誤: %s", errorMsg)
	}

	// 解析心跳回應: HEARTBEAT:1234567890
	if strings.HasPrefix(message, "HEARTBEAT:") {
		timestamp := strings.TrimPrefix(message, "HEARTBEAT:")
		m.logger.Debugf("收到心跳: %s", timestamp)
	}

	// 解析顯示回應: DISPLAY:PIN_CODE:123
	if strings.HasPrefix(message, "DISPLAY:") {
		displayMsg := strings.TrimPrefix(message, "DISPLAY:")
		m.logger.Infof("顯示訊息: %s", displayMsg)
	}

	// 解析資訊回應: INFO:ESP32-001
	if strings.HasPrefix(message, "INFO:") {
		info := strings.TrimPrefix(message, "INFO:")
		m.logger.Debugf("裝置資訊: %s", info)
	}
}

// HasNewMessage 檢查是否有新訊息
func (m *Manager) HasNewMessage() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 如果最後訊息時間在最近5秒內，認為有新訊息
	return time.Since(m.lastMessage) < 5*time.Second
}

// GetLastMessageTime 取得最後訊息時間
func (m *Manager) GetLastMessageTime() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastMessage
}

// GetLastPinCode 取得最後收到的PIN碼
func (m *Manager) GetLastPinCode() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastPinCode
}

// DisplayPinCode 在Arduino/ESP螢幕上顯示PIN碼
func (m *Manager) DisplayPinCode(pinCode string) error {
	if !m.isConnected {
		return fmt.Errorf("裝置未連線")
	}

	// Arduino/ESP指令格式: DISPLAY_PIN:123
	command := fmt.Sprintf("DISPLAY_PIN:%s\n", pinCode)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.port.Write(ctx, []byte(command))
	if err != nil {
		return fmt.Errorf("發送PIN碼顯示指令失敗: %v", err)
	}

	m.logger.Infof("已發送PIN碼到Arduino/ESP裝置: %s", pinCode)
	return nil
}

// SendCommand 發送指令到Arduino/ESP裝置
func (m *Manager) SendCommand(command string) error {
	if !m.isConnected {
		return fmt.Errorf("裝置未連線")
	}

	fullCommand := command + "\n"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.port.Write(ctx, []byte(fullCommand))
	if err != nil {
		return fmt.Errorf("發送指令失敗: %v", err)
	}

	m.logger.Debugf("已發送指令到Arduino/ESP: %s", command)
	return nil
}

// GetDeviceInfo 取得Arduino/ESP裝置資訊
func (m *Manager) GetDeviceInfo() (map[string]string, error) {
	if !m.isConnected {
		return nil, fmt.Errorf("裝置未連線")
	}

	// 發送裝置資訊查詢指令
	if err := m.SendCommand("GET_INFO"); err != nil {
		return nil, err
	}

	// 等待回應
	time.Sleep(1 * time.Second)

	// 這裡應該實作讀取Arduino/ESP回應的邏輯
	// 暫時返回模擬資料
	info := map[string]string{
		"device_type":      "Arduino/ESP IoT Device",
		"firmware_version": "2.0.0",
		"serial_number":    "ESP32-001",
		"status":           "connected",
		"board_type":       "ESP32",
	}

	return info, nil
}

// RequestPinGeneration 請求Arduino/ESP生成PIN碼
func (m *Manager) RequestPinGeneration() error {
	if !m.isConnected {
		return fmt.Errorf("裝置未連線")
	}

	// 發送PIN碼生成請求
	if err := m.SendCommand("GENERATE_PIN"); err != nil {
		return fmt.Errorf("發送PIN碼生成請求失敗: %v", err)
	}

	m.logger.Info("已請求Arduino/ESP生成PIN碼")
	return nil
}

// DisplaySuccess 在Arduino/ESP螢幕上顯示成功訊息
func (m *Manager) DisplaySuccess() error {
	return m.SendCommand("DISPLAY_SUCCESS")
}

// DisplayError 在Arduino/ESP螢幕上顯示錯誤訊息
func (m *Manager) DisplayError() error {
	return m.SendCommand("DISPLAY_ERROR")
}

// Close 關閉Arduino/ESP裝置連線
func (m *Manager) Close() error {
	if m.port != nil {
		m.isConnected = false
		return m.port.Close()
	}
	return nil
}

// IsConnected 檢查Arduino/ESP裝置是否連線
func (m *Manager) IsConnected() bool {
	return m.isConnected
}

// ValidateArduinoESPResponse 驗證Arduino/ESP回應格式
func (m *Manager) ValidateArduinoESPResponse(response string) bool {
	// Arduino/ESP回應格式驗證
	patterns := []string{
		`^PIN:\d{3}$`,  // PIN:123
		`^STATUS:\w+$`, // STATUS:READY
		`^ERROR:\w+$`,  // ERROR:NO_DISPLAY
		`^INFO:.+$`,    // INFO:ESP32-001
	}

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, response)
		if matched {
			return true
		}
	}

	return false
}
