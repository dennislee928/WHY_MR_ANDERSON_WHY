package device

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ProductType 產品類型
type ProductType int

const (
	Unknown ProductType = iota
	InternetBlocker
	NetworkController
	SecurityDevice
)

// String 回傳產品類型字串
func (pt ProductType) String() string {
	switch pt {
	case InternetBlocker:
		return "Internet Blocker"
	case NetworkController:
		return "Network Controller"
	case SecurityDevice:
		return "Security Device"
	default:
		return "Unknown"
	}
}

// Detector 產品類型偵測器
type Detector struct {
	deviceManager *Manager
}

// NewDetector 建立新的產品類型偵測器
func NewDetector(deviceManager *Manager) *Detector {
	return &Detector{
		deviceManager: deviceManager,
	}
}

// DetectProductType 偵測產品類型
func (d *Detector) DetectProductType() (ProductType, error) {
	if !d.deviceManager.IsConnected() {
		return Unknown, fmt.Errorf("裝置未連線")
	}

	// 發送產品識別指令
	if err := d.deviceManager.SendCommand("IDENTIFY_PRODUCT"); err != nil {
		return Unknown, fmt.Errorf("發送識別指令失敗: %v", err)
	}

	// 等待回應並分析
	time.Sleep(1 * time.Second)

	// 取得裝置資訊
	info, err := d.deviceManager.GetDeviceInfo()
	if err != nil {
		return Unknown, fmt.Errorf("取得裝置資訊失敗: %v", err)
	}

	// 根據裝置資訊判斷產品類型
	return d.analyzeProductType(info), nil
}

// analyzeProductType 分析產品類型
func (d *Detector) analyzeProductType(info map[string]string) ProductType {
	deviceType := strings.ToLower(info["device_type"])
	firmwareVersion := info["firmware_version"]
	serialNumber := info["serial_number"]

	// 檢查裝置類型字串
	if strings.Contains(deviceType, "ch340") {
		// 檢查序號模式
		if d.isInternetBlockerSerial(serialNumber) {
			return InternetBlocker
		}

		// 檢查韌體版本
		if d.isInternetBlockerFirmware(firmwareVersion) {
			return InternetBlocker
		}

		// 檢查序號前綴
		if strings.HasPrefix(serialNumber, "IB-") {
			return InternetBlocker
		}

		if strings.HasPrefix(serialNumber, "NC-") {
			return NetworkController
		}

		if strings.HasPrefix(serialNumber, "SD-") {
			return SecurityDevice
		}
	}

	return Unknown
}

// isInternetBlockerSerial 檢查是否為網路阻斷器序號
func (d *Detector) isInternetBlockerSerial(serialNumber string) bool {
	// 網路阻斷器序號模式: IB-XXXX-XXXX
	pattern := `^IB-\d{4}-\d{4}$`
	matched, _ := regexp.MatchString(pattern, serialNumber)
	return matched
}

// isInternetBlockerFirmware 檢查是否為網路阻斷器韌體
func (d *Detector) isInternetBlockerFirmware(version string) bool {
	// 網路阻斷器韌體版本模式: IB-vX.X.X
	pattern := `^IB-v\d+\.\d+\.\d+$`
	matched, _ := regexp.MatchString(pattern, version)
	return matched
}

// ValidateProductType 驗證產品類型
func (d *Detector) ValidateProductType(expectedType ProductType) error {
	detectedType, err := d.DetectProductType()
	if err != nil {
		return fmt.Errorf("偵測產品類型失敗: %v", err)
	}

	if detectedType != expectedType {
		return fmt.Errorf("產品類型不符: 期望 %s, 實際 %s",
			expectedType.String(), detectedType.String())
	}

	return nil
}

// GetProductCapabilities 取得產品功能
func (d *Detector) GetProductCapabilities() (map[string]bool, error) {
	productType, err := d.DetectProductType()
	if err != nil {
		return nil, err
	}

	capabilities := make(map[string]bool)

	switch productType {
	case InternetBlocker:
		capabilities["network_blocking"] = true
		capabilities["pin_display"] = true
		capabilities["time_control"] = true
		capabilities["usb_token_auth"] = true
		capabilities["ethernet_management"] = true

	case NetworkController:
		capabilities["network_blocking"] = true
		capabilities["pin_display"] = false
		capabilities["time_control"] = true
		capabilities["usb_token_auth"] = false
		capabilities["ethernet_management"] = true

	case SecurityDevice:
		capabilities["network_blocking"] = false
		capabilities["pin_display"] = true
		capabilities["time_control"] = false
		capabilities["usb_token_auth"] = true
		capabilities["ethernet_management"] = false

	default:
		capabilities["network_blocking"] = false
		capabilities["pin_display"] = false
		capabilities["time_control"] = false
		capabilities["usb_token_auth"] = false
		capabilities["ethernet_management"] = false
	}

	return capabilities, nil
}
