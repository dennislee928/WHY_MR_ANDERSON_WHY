package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// CommandValidator 提供安全的命令執行驗證
type CommandValidator struct {
	allowedCommands map[string]bool
	allowedArgs     map[string]*regexp.Regexp
}

// NewCommandValidator 創建新的命令驗證器
func NewCommandValidator() *CommandValidator {
	return &CommandValidator{
		allowedCommands: map[string]bool{
			// 系統信息命令
			"ls":       true,
			"ps":       true,
			"df":       true,
			"free":     true,
			"uptime":   true,
			"hostname": true,
			"whoami":   true,
			
			// 網路命令
			"ping":     true,
			"netstat":  true,
			"ss":       true,
			"ip":       true,
			"ifconfig": true,
			
			// 日誌命令
			"tail":     true,
			"head":     true,
			"cat":      true,
			"grep":     true,
			
			// Docker 命令
			"docker":   true,
		},
		allowedArgs: map[string]*regexp.Regexp{
			// 只允許安全的參數模式
			"ls":     regexp.MustCompile(`^[-altrh]+$`),
			"ps":     regexp.MustCompile(`^[-aux]+$`),
			"tail":   regexp.MustCompile(`^-[nf]?\d*$`),
			"head":   regexp.MustCompile(`^-n?\d+$`),
			"ping":   regexp.MustCompile(`^-[c]\d+$`),
			"docker": regexp.MustCompile(`^(ps|logs|stats|inspect)$`),
		},
	}
}

// ValidateCommand 驗證命令是否安全
func (cv *CommandValidator) ValidateCommand(cmd string, args []string) error {
	// 檢查命令是否在白名單中
	if !cv.allowedCommands[cmd] {
		return fmt.Errorf("command not allowed: %s", cmd)
	}

	// 檢查參數是否安全
	if pattern, exists := cv.allowedArgs[cmd]; exists {
		for _, arg := range args {
			// 跳過文件路徑（以 / 或 . 開頭）
			if strings.HasPrefix(arg, "/") || strings.HasPrefix(arg, ".") {
				continue
			}
			
			// 驗證其他參數
			if !pattern.MatchString(arg) {
				return fmt.Errorf("argument not allowed for command %s: %s", cmd, arg)
			}
		}
	}

	// 檢查是否包含危險字符
	dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "\n", "\r"}
	for _, arg := range args {
		for _, char := range dangerousChars {
			if strings.Contains(arg, char) {
				return fmt.Errorf("dangerous character detected in argument: %s", char)
			}
		}
	}

	return nil
}

// ExecuteCommand 安全地執行命令
func (cv *CommandValidator) ExecuteCommand(cmd string, args ...string) ([]byte, error) {
	// 驗證命令
	if err := cv.ValidateCommand(cmd, args); err != nil {
		return nil, err
	}

	// 執行命令
	output, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("command execution failed: %w", err)
	}

	return output, nil
}

// ExecuteCommandString 安全地執行命令（字符串版本）
func (cv *CommandValidator) ExecuteCommandString(cmdStr string) ([]byte, error) {
	// 解析命令字符串
	parts := strings.Fields(cmdStr)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	cmd := parts[0]
	args := parts[1:]

	return cv.ExecuteCommand(cmd, args...)
}

// AddAllowedCommand 添加允許的命令（用於擴展）
func (cv *CommandValidator) AddAllowedCommand(cmd string, argPattern string) error {
	cv.allowedCommands[cmd] = true
	
	if argPattern != "" {
		pattern, err := regexp.Compile(argPattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}
		cv.allowedArgs[cmd] = pattern
	}
	
	return nil
}

// RemoveAllowedCommand 移除允許的命令
func (cv *CommandValidator) RemoveAllowedCommand(cmd string) {
	delete(cv.allowedCommands, cmd)
	delete(cv.allowedArgs, cmd)
}

// IsCommandAllowed 檢查命令是否被允許
func (cv *CommandValidator) IsCommandAllowed(cmd string) bool {
	return cv.allowedCommands[cmd]
}

// GetAllowedCommands 獲取所有允許的命令列表
func (cv *CommandValidator) GetAllowedCommands() []string {
	commands := make([]string, 0, len(cv.allowedCommands))
	for cmd := range cv.allowedCommands {
		commands = append(commands, cmd)
	}
	return commands
}

