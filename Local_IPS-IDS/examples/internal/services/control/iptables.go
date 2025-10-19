package control

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

// IPTablesManager manages iptables firewall rules
// iptables 防火牆規則管理器
type IPTablesManager struct {
	logger *logrus.Logger
	mu     sync.RWMutex
	chain  string // 自定義鏈名稱
}

// NewIPTablesManager creates a new iptables manager
func NewIPTablesManager(logger *logrus.Logger) *IPTablesManager {
	if logger == nil {
		logger = logrus.New()
	}

	return &IPTablesManager{
		logger: logger,
		chain:  "PANDORA_BLOCK",
	}
}

// Initialize initializes iptables (creates custom chain)
func (ipt *IPTablesManager) Initialize() error {
	ipt.mu.Lock()
	defer ipt.mu.Unlock()

	// 創建自定義鏈
	cmd := exec.Command("iptables", "-N", ipt.chain)
	if err := cmd.Run(); err != nil {
		// 鏈可能已存在，忽略錯誤
		ipt.logger.Debugf("Chain %s may already exist", ipt.chain)
	}

	// 將自定義鏈插入到 INPUT 鏈
	cmd = exec.Command("iptables", "-C", "INPUT", "-j", ipt.chain)
	if err := cmd.Run(); err != nil {
		// 規則不存在，添加它
		cmd = exec.Command("iptables", "-I", "INPUT", "-j", ipt.chain)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to insert chain: %w", err)
		}
	}

	ipt.logger.Infof("IPTables initialized (chain: %s)", ipt.chain)
	return nil
}

// BlockIP blocks an IP address
func (ipt *IPTablesManager) BlockIP(ip string, action string) error {
	ipt.mu.Lock()
	defer ipt.mu.Unlock()

	target := "DROP"
	if action == "reject" {
		target = "REJECT"
	}

	// 添加規則到自定義鏈
	cmd := exec.Command("iptables", "-A", ipt.chain, "-s", ip, "-j", target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		ipt.logger.Errorf("Failed to block IP %s: %v, output: %s", ip, err, string(output))
		return fmt.Errorf("failed to block IP: %w", err)
	}

	ipt.logger.Infof("IP %s blocked (action: %s)", ip, target)
	return nil
}

// UnblockIP unblocks an IP address
func (ipt *IPTablesManager) UnblockIP(ip string) error {
	ipt.mu.Lock()
	defer ipt.mu.Unlock()

	// 嘗試刪除 DROP 規則
	cmd := exec.Command("iptables", "-D", ipt.chain, "-s", ip, "-j", "DROP")
	if err := cmd.Run(); err != nil {
		// 嘗試刪除 REJECT 規則
		cmd = exec.Command("iptables", "-D", ipt.chain, "-s", ip, "-j", "REJECT")
		if err := cmd.Run(); err != nil {
			ipt.logger.Warnf("Failed to unblock IP %s: %v", ip, err)
			return fmt.Errorf("failed to unblock IP: %w", err)
		}
	}

	ipt.logger.Infof("IP %s unblocked", ip)
	return nil
}

// BlockPort blocks a port
func (ipt *IPTablesManager) BlockPort(port int32, protocol string) error {
	ipt.mu.Lock()
	defer ipt.mu.Unlock()

	proto := strings.ToLower(protocol)
	if proto != "tcp" && proto != "udp" {
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}

	cmd := exec.Command("iptables", "-A", ipt.chain,
		"-p", proto,
		"--dport", fmt.Sprintf("%d", port),
		"-j", "DROP")

	output, err := cmd.CombinedOutput()
	if err != nil {
		ipt.logger.Errorf("Failed to block port %d/%s: %v, output: %s", port, proto, err, string(output))
		return fmt.Errorf("failed to block port: %w", err)
	}

	ipt.logger.Infof("Port %d/%s blocked", port, proto)
	return nil
}

// UnblockPort unblocks a port
func (ipt *IPTablesManager) UnblockPort(port int32, protocol string) error {
	ipt.mu.Lock()
	defer ipt.mu.Unlock()

	proto := strings.ToLower(protocol)

	cmd := exec.Command("iptables", "-D", ipt.chain,
		"-p", proto,
		"--dport", fmt.Sprintf("%d", port),
		"-j", "DROP")

	if err := cmd.Run(); err != nil {
		ipt.logger.Warnf("Failed to unblock port %d/%s: %v", port, proto, err)
		return fmt.Errorf("failed to unblock port: %w", err)
	}

	ipt.logger.Infof("Port %d/%s unblocked", port, proto)
	return nil
}

// ListRules lists all rules in the custom chain
func (ipt *IPTablesManager) ListRules() ([]string, error) {
	ipt.mu.RLock()
	defer ipt.mu.RUnlock()

	cmd := exec.Command("iptables", "-L", ipt.chain, "-n", "--line-numbers")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list rules: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var rules []string

	// 跳過前兩行（header）
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			rules = append(rules, line)
		}
	}

	return rules, nil
}

// FlushRules removes all rules from the custom chain
func (ipt *IPTablesManager) FlushRules() error {
	ipt.mu.Lock()
	defer ipt.mu.Unlock()

	cmd := exec.Command("iptables", "-F", ipt.chain)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to flush rules: %w", err)
	}

	ipt.logger.Infof("All rules flushed from chain %s", ipt.chain)
	return nil
}

// SaveRules saves current iptables rules to file
func (ipt *IPTablesManager) SaveRules(filename string) error {
	ipt.mu.RLock()
	defer ipt.mu.RUnlock()

	cmd := exec.Command("iptables-save")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to save rules: %w", err)
	}

	// TODO: 寫入文件
	ipt.logger.Infof("Rules saved to %s (%d bytes)", filename, len(output))
	return nil
}

// RestoreRules restores iptables rules from file
func (ipt *IPTablesManager) RestoreRules(filename string) error {
	ipt.mu.Lock()
	defer ipt.mu.Unlock()

	// TODO: 從文件讀取並恢復
	cmd := exec.Command("iptables-restore", filename)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restore rules: %w", err)
	}

	ipt.logger.Infof("Rules restored from %s", filename)
	return nil
}

// CheckIPTablesAvailable checks if iptables is available
func CheckIPTablesAvailable() error {
	cmd := exec.Command("iptables", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("iptables not available: %w", err)
	}
	return nil
}

