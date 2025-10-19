package device

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

// SerialDevice 代表序列埠裝置
type SerialDevice struct {
	Port     string
	BaudRate int
	Timeout  time.Duration
	logger   *logrus.Logger
}

// NewSerialDevice 建立新的序列埠裝置
func NewSerialDevice(port string, baudRate int, timeout time.Duration, logger *logrus.Logger) *SerialDevice {
	return &SerialDevice{
		Port:     port,
		BaudRate: baudRate,
		Timeout:  timeout,
		logger:   logger,
	}
}

// Connect 連接到序列埠裝置
func (s *SerialDevice) Connect(ctx context.Context) error {
	s.logger.WithFields(logrus.Fields{
		"port":      s.Port,
		"baud_rate": s.BaudRate,
	}).Info("嘗試連接序列埠裝置")

	// 模擬連接過程
	select {
	case <-time.After(1 * time.Second):
		s.logger.Info("序列埠裝置連接成功")
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Read 從序列埠讀取資料
func (s *SerialDevice) Read(ctx context.Context) ([]byte, error) {
	// 模擬讀取資料
	select {
	case <-time.After(100 * time.Millisecond):
		data := []byte("模擬序列埠資料")
		s.logger.WithField("data_length", len(data)).Debug("從序列埠讀取資料")
		return data, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Write 向序列埠寫入資料
func (s *SerialDevice) Write(ctx context.Context, data []byte) error {
	s.logger.WithField("data_length", len(data)).Debug("向序列埠寫入資料")

	// 模擬寫入過程
	select {
	case <-time.After(50 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close 關閉序列埠連接
func (s *SerialDevice) Close() error {
	s.logger.Info("關閉序列埠連接")
	return nil
}

// IsConnected 檢查是否已連接
func (s *SerialDevice) IsConnected() bool {
	// 模擬連接狀態檢查
	return true
}

// GetDeviceInfo 取得裝置資訊
func (s *SerialDevice) GetDeviceInfo() map[string]interface{} {
	return map[string]interface{}{
		"port":      s.Port,
		"baud_rate": s.BaudRate,
		"timeout":   s.Timeout.String(),
		"connected": s.IsConnected(),
	}
}

// SerialDeviceManager 管理多個序列埠裝置
type SerialDeviceManager struct {
	devices map[string]*SerialDevice
	logger  *logrus.Logger
}

// NewSerialDeviceManager 建立新的序列埠裝置管理器
func NewSerialDeviceManager(logger *logrus.Logger) *SerialDeviceManager {
	return &SerialDeviceManager{
		devices: make(map[string]*SerialDevice),
		logger:  logger,
	}
}

// AddDevice 新增裝置
func (m *SerialDeviceManager) AddDevice(id string, device *SerialDevice) {
	m.devices[id] = device
	m.logger.WithField("device_id", id).Info("新增序列埠裝置")
}

// GetDevice 取得裝置
func (m *SerialDeviceManager) GetDevice(id string) (*SerialDevice, bool) {
	device, exists := m.devices[id]
	return device, exists
}

// RemoveDevice 移除裝置
func (m *SerialDeviceManager) RemoveDevice(id string) error {
	if device, exists := m.devices[id]; exists {
		if err := device.Close(); err != nil {
			return err
		}
		delete(m.devices, id)
		m.logger.WithField("device_id", id).Info("移除序列埠裝置")
	}
	return nil
}

// ListDevices 列出所有裝置
func (m *SerialDeviceManager) ListDevices() map[string]map[string]interface{} {
	result := make(map[string]map[string]interface{})
	for id, device := range m.devices {
		result[id] = device.GetDeviceInfo()
	}
	return result
}

// CloseAll 關閉所有裝置
func (m *SerialDeviceManager) CloseAll() error {
	for id, device := range m.devices {
		if err := device.Close(); err != nil {
			m.logger.WithField("device_id", id).WithError(err).Error("關閉裝置時發生錯誤")
			return err
		}
	}
	m.logger.Info("所有序列埠裝置已關閉")
	return nil
}
