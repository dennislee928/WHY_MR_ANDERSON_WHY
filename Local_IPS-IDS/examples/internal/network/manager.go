package network

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Manager 網路管理器
type Manager struct {
	logger        *logrus.Logger
	isBlocked     bool
	blockTime     time.Time
	mu            sync.RWMutex
	interfaceName string
}

// NewManager 建立新的網路管理器
func NewManager(logger *logrus.Logger) *Manager {
	return &Manager{
		logger:        logger,
		interfaceName: "eth0", // 預設乙太網路介面
	}
}

// Initialize 初始化網路管理器
func (m *Manager) Initialize() error {
	// 偵測乙太網路介面
	if err := m.detectEthernetInterface(); err != nil {
		m.logger.Warnf("無法偵測乙太網路介面: %v", err)
	}

	m.logger.Info("網路管理器初始化完成")
	return nil
}

// detectEthernetInterface 偵測乙太網路介面
func (m *Manager) detectEthernetInterface() error {
	switch runtime.GOOS {
	case "linux":
		return m.detectLinuxInterface()
	case "windows":
		return m.detectWindowsInterface()
	case "darwin":
		return m.detectMacOSInterface()
	default:
		return fmt.Errorf("不支援的作業系統: %s", runtime.GOOS)
	}
}

// detectLinuxInterface 偵測Linux乙太網路介面
func (m *Manager) detectLinuxInterface() error {
	cmd := exec.Command("ip", "link", "show")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("執行ip link show失敗: %v", err)
	}

	lines := strings.Split(string(output), "\n")

	// 優先尋找eth0, eth1等乙太網路介面
	ethernetInterfaces := []string{"eth0", "eth1", "eth2", "enp0s3", "enp0s8"}

	for _, iface := range ethernetInterfaces {
		for _, line := range lines {
			if strings.Contains(line, iface) && !strings.Contains(line, "state DOWN") {
				m.interfaceName = iface
				m.logger.Infof("偵測到乙太網路介面: %s", m.interfaceName)
				return nil
			}
		}
	}

	// 如果沒找到，嘗試其他方式
	for _, line := range lines {
		if strings.Contains(line, "eth") && strings.Contains(line, "state UP") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				m.interfaceName = strings.TrimSpace(parts[1])
				m.logger.Infof("偵測到乙太網路介面: %s", m.interfaceName)
				return nil
			}
		}
	}

	return fmt.Errorf("未找到可用的乙太網路介面")
}

// detectWindowsInterface 偵測Windows乙太網路介面
func (m *Manager) detectWindowsInterface() error {
	cmd := exec.Command("netsh", "interface", "show", "interface")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("執行netsh失敗: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Ethernet") && strings.Contains(line, "Enabled") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				m.interfaceName = parts[3]
				m.logger.Infof("偵測到乙太網路介面: %s", m.interfaceName)
				return nil
			}
		}
	}

	return fmt.Errorf("未找到可用的乙太網路介面")
}

// detectMacOSInterface 偵測macOS乙太網路介面
func (m *Manager) detectMacOSInterface() error {
	cmd := exec.Command("ifconfig")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("執行ifconfig失敗: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "en") && strings.Contains(line, "flags=") {
			parts := strings.Split(line, ":")
			if len(parts) >= 1 {
				m.interfaceName = strings.TrimSpace(parts[0])
				m.logger.Infof("偵測到乙太網路介面: %s", m.interfaceName)
				return nil
			}
		}
	}

	return fmt.Errorf("未找到可用的乙太網路介面")
}

// BlockEthernet 阻斷乙太網路連線
func (m *Manager) BlockEthernet() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.isBlocked {
		return nil // 已經阻斷
	}

	switch runtime.GOOS {
	case "linux":
		if err := m.blockLinuxEthernet(); err != nil {
			return err
		}
	case "windows":
		if err := m.blockWindowsEthernet(); err != nil {
			return err
		}
	case "darwin":
		if err := m.blockMacOSEthernet(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("不支援的作業系統: %s", runtime.GOOS)
	}

	m.isBlocked = true
	m.blockTime = time.Now()
	m.logger.Info("乙太網路已阻斷")
	return nil
}

// blockLinuxEthernet Linux阻斷乙太網路
func (m *Manager) blockLinuxEthernet() error {
	validInterface := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	if !validInterface.MatchString(m.interfaceName) {
		return fmt.Errorf("invalid interface name")
	}
	// 方法1: 使用ip命令停用介面
	cmd := exec.Command("sudo", "ip", "link", "set", m.interfaceName, "down")
	if err := cmd.Run(); err != nil {
		m.logger.Warnf("使用ip命令停用介面失敗: %v", err)

		// 方法2: 使用ifconfig命令
		cmd = exec.Command("sudo", "ifconfig", m.interfaceName, "down")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("停用乙太網路介面失敗: %v", err)
		}
	}

	return nil
}

// blockWindowsEthernet Windows阻斷乙太網路
func (m *Manager) blockWindowsEthernet() error {
	validInterface := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	if !validInterface.MatchString(m.interfaceName) {
		return fmt.Errorf("invalid interface name")
	}
	cmd := exec.Command("netsh", "interface", "set", "interface", m.interfaceName, "disable")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停用乙太網路介面失敗: %v", err)
	}

	return nil
}

// blockMacOSEthernet macOS阻斷乙太網路
func (m *Manager) blockMacOSEthernet() error {
	validInterface := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	if !validInterface.MatchString(m.interfaceName) {
		return fmt.Errorf("invalid interface name")
	}
	cmd := exec.Command("sudo", "ifconfig", m.interfaceName, "down")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停用乙太網路介面失敗: %v", err)
	}

	return nil
}

// EnableEthernet 啟用乙太網路連線
func (m *Manager) EnableEthernet() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.isBlocked {
		return nil // 已經啟用
	}

	switch runtime.GOOS {
	case "linux":
		if err := m.enableLinuxEthernet(); err != nil {
			return err
		}
	case "windows":
		if err := m.enableWindowsEthernet(); err != nil {
			return err
		}
	case "darwin":
		if err := m.enableMacOSEthernet(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("不支援的作業系統: %s", runtime.GOOS)
	}

	m.isBlocked = false
	m.logger.Info("乙太網路已啟用")
	return nil
}

// enableLinuxEthernet Linux啟用乙太網路
func (m *Manager) enableLinuxEthernet() error {
	validInterface := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	if !validInterface.MatchString(m.interfaceName) {
		return fmt.Errorf("invalid interface name")
	}
	// 方法1: 使用ip命令啟用介面
	cmd := exec.Command("sudo", "ip", "link", "set", m.interfaceName, "up")
	if err := cmd.Run(); err != nil {
		m.logger.Warnf("使用ip命令啟用介面失敗: %v", err)

		// 方法2: 使用ifconfig命令
		cmd = exec.Command("sudo", "ifconfig", m.interfaceName, "up")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("啟用乙太網路介面失敗: %v", err)
		}
	}

	return nil
}

// enableWindowsEthernet Windows啟用乙太網路
func (m *Manager) enableWindowsEthernet() error {
	validInterface := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	if !validInterface.MatchString(m.interfaceName) {
		return fmt.Errorf("invalid interface name")
	}
	cmd := exec.Command("netsh", "interface", "set", "interface", m.interfaceName, "enable")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("啟用乙太網路介面失敗: %v", err)
	}

	return nil
}

// enableMacOSEthernet macOS啟用乙太網路
func (m *Manager) enableMacOSEthernet() error {
	validInterface := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	if !validInterface.MatchString(m.interfaceName) {
		return fmt.Errorf("invalid interface name")
	}
	cmd := exec.Command("sudo", "ifconfig", m.interfaceName, "up")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("啟用乙太網路介面失敗: %v", err)
	}

	return nil
}

// IsBlocked 檢查網路是否被阻斷
func (m *Manager) IsBlocked() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isBlocked
}

// GetBlockTime 取得阻斷時間
func (m *Manager) GetBlockTime() time.Time {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.blockTime
}

// GetInterfaceName 取得介面名稱
func (m *Manager) GetInterfaceName() string {
	return m.interfaceName
}

// TestConnection 測試網路連線
func (m *Manager) TestConnection() bool {
	// 測試連線到Google DNS
	cmd := exec.Command("ping", "-c", "1", "8.8.8.8")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "1", "8.8.8.8")
	}

	err := cmd.Run()
	return err == nil
}
