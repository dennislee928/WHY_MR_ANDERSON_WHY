package device

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

// SerialDevice represents a USB-SERIAL device connection
// USB-SERIAL 設備連接
type SerialDevice struct {
	port       serial.Port
	portName   string
	baudRate   int
	config     *serial.Mode
	logger     *logrus.Logger
	mu         sync.RWMutex
	connected  bool
	lastRead   time.Time
	lastWrite  time.Time
	metrics    *DeviceMetrics
}

// NewSerialDevice creates a new serial device connection
func NewSerialDevice(portName string, baudRate int, logger *logrus.Logger) *SerialDevice {
	if logger == nil {
		logger = logrus.New()
	}

	return &SerialDevice{
		portName: portName,
		baudRate: baudRate,
		config: &serial.Mode{
			BaudRate: baudRate,
			DataBits: 8,
			Parity:   serial.NoParity,
			StopBits: serial.OneStopBit,
		},
		logger:  logger,
		metrics: &DeviceMetrics{},
	}
}

// Open opens the serial port
func (sd *SerialDevice) Open() error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	if sd.connected {
		return fmt.Errorf("device already connected")
	}

	port, err := serial.Open(sd.portName, sd.config)
	if err != nil {
		sd.logger.Errorf("Failed to open serial port %s: %v", sd.portName, err)
		return fmt.Errorf("failed to open serial port: %w", err)
	}

	// 設置讀取超時
	if err := port.SetReadTimeout(30 * time.Second); err != nil {
		port.Close()
		return fmt.Errorf("failed to set read timeout: %w", err)
	}

	sd.port = port
	sd.connected = true

	sd.logger.Infof("Serial port %s opened successfully (baud rate: %d)", sd.portName, sd.baudRate)
	return nil
}

// Close closes the serial port
func (sd *SerialDevice) Close() error {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	if !sd.connected || sd.port == nil {
		return nil
	}

	if err := sd.port.Close(); err != nil {
		sd.logger.Errorf("Failed to close serial port: %v", err)
		return err
	}

	sd.connected = false
	sd.port = nil

	sd.logger.Infof("Serial port %s closed", sd.portName)
	return nil
}

// Read reads data from the serial port
func (sd *SerialDevice) Read(buffer []byte) (int, error) {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	if !sd.connected || sd.port == nil {
		return 0, fmt.Errorf("device not connected")
	}

	n, err := sd.port.Read(buffer)
	if err != nil && err != io.EOF {
		sd.metrics.Errors++
		sd.logger.Errorf("Failed to read from serial port: %v", err)
		return n, err
	}

	sd.lastRead = time.Now()
	sd.metrics.BytesRead += int64(n)
	sd.metrics.ReadOperations++

	return n, err
}

// Write writes data to the serial port
func (sd *SerialDevice) Write(data []byte) (int, error) {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	if !sd.connected || sd.port == nil {
		return 0, fmt.Errorf("device not connected")
	}

	n, err := sd.port.Write(data)
	if err != nil {
		sd.metrics.Errors++
		sd.logger.Errorf("Failed to write to serial port: %v", err)
		return n, err
	}

	sd.lastWrite = time.Now()
	sd.metrics.BytesWritten += int64(n)
	sd.metrics.WriteOperations++

	return n, nil
}

// ReadLine reads a line from the serial port (until \n)
func (sd *SerialDevice) ReadLine(timeout time.Duration) (string, error) {
	buffer := make([]byte, 1024)
	line := make([]byte, 0, 1024)

	deadline := time.Now().Add(timeout)

	for {
		if time.Now().After(deadline) {
			return "", fmt.Errorf("read timeout")
		}

		n, err := sd.Read(buffer)
		if err != nil && err != io.EOF {
			return "", err
		}

		if n > 0 {
			for i := 0; i < n; i++ {
				if buffer[i] == '\n' {
					return string(line), nil
				}
				line = append(line, buffer[i])
			}
		}

		// 短暫休眠避免 CPU 100%
		time.Sleep(10 * time.Millisecond)
	}
}

// WriteLine writes a line to the serial port (adds \n)
func (sd *SerialDevice) WriteLine(line string) error {
	data := []byte(line + "\n")
	_, err := sd.Write(data)
	return err
}

// IsConnected checks if the device is connected
func (sd *SerialDevice) IsConnected() bool {
	sd.mu.RLock()
	defer sd.mu.RUnlock()
	return sd.connected
}

// GetMetrics returns device metrics
func (sd *SerialDevice) GetMetrics() *DeviceMetrics {
	sd.mu.RLock()
	defer sd.mu.RUnlock()

	return &DeviceMetrics{
		BytesRead:       sd.metrics.BytesRead,
		BytesWritten:    sd.metrics.BytesWritten,
		ReadOperations:  sd.metrics.ReadOperations,
		WriteOperations: sd.metrics.WriteOperations,
		Errors:          sd.metrics.Errors,
	}
}

// ResetMetrics resets device metrics
func (sd *SerialDevice) ResetMetrics() {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	sd.metrics = &DeviceMetrics{}
}

// ListAvailablePorts lists all available serial ports
func ListAvailablePorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, fmt.Errorf("failed to list ports: %w", err)
	}

	return ports, nil
}

// DetectCH340Devices detects CH340 USB-SERIAL devices
func DetectCH340Devices(logger *logrus.Logger) ([]string, error) {
	ports, err := ListAvailablePorts()
	if err != nil {
		return nil, err
	}

	var ch340Ports []string

	for _, port := range ports {
		// CH340 設備通常在 Linux 上顯示為 /dev/ttyUSB*
		// 在 Windows 上顯示為 COM*
		// 可以通過嘗試打開來驗證
		if logger != nil {
			logger.Debugf("Found serial port: %s", port)
		}

		// TODO: 更精確的 CH340 檢測邏輯
		// 可以通過 USB VID/PID 識別
		ch340Ports = append(ch340Ports, port)
	}

	return ch340Ports, nil
}

