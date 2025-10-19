package utils

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// TimeUtils 時間工具
type TimeUtils struct{}

// ParseTime 解析時間字串 (HH:MM格式)
func (tu *TimeUtils) ParseTime(timeStr string) (hour, minute int, err error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("時間格式錯誤: %s", timeStr)
	}

	hour, err = parseInt(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("小時格式錯誤: %s", parts[0])
	}

	minute, err = parseInt(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("分鐘格式錯誤: %s", parts[1])
	}

	if hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("小時超出範圍: %d", hour)
	}

	if minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("分鐘超出範圍: %d", minute)
	}

	return hour, minute, nil
}

// IsTimeReached 檢查是否到達指定時間
func (tu *TimeUtils) IsTimeReached(targetTime string) bool {
	hour, minute, err := tu.ParseTime(targetTime)
	if err != nil {
		return false
	}

	now := time.Now()
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())

	// 檢查是否在目標時間的前後1分鐘內
	diff := now.Sub(target)
	return diff >= -time.Minute && diff <= time.Minute
}

// FormatDuration 格式化時間間隔
func (tu *TimeUtils) FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f秒", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0f分鐘", d.Minutes())
	} else {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		return fmt.Sprintf("%d小時%d分鐘", hours, minutes)
	}
}

// SystemUtils 系統工具
type SystemUtils struct{}

// IsRoot 檢查是否為root權限
func (su *SystemUtils) IsRoot() bool {
	return os.Geteuid() == 0
}

// RequireRoot 要求root權限
func (su *SystemUtils) RequireRoot() error {
	if !su.IsRoot() {
		return fmt.Errorf("此操作需要root權限")
	}
	return nil
}

// GetOS 取得作業系統類型
func (su *SystemUtils) GetOS() string {
	return runtime.GOOS
}

// ExecuteCommand 執行系統命令（使用安全驗證器）
func (su *SystemUtils) ExecuteCommand(command string, args ...string) (string, error) {
	// 使用命令驗證器進行安全檢查
	validator := &CommandValidator{
		allowedCommands: map[string]bool{
			"ls": true, "ps": true, "df": true, "free": true,
			"ping": true, "netstat": true, "docker": true,
			"tail": true, "cat": true, "grep": true,
		},
	}
	
	// 驗證命令
	if !validator.allowedCommands[command] {
		return "", fmt.Errorf("執行命令失敗: command not allowed: %s", command)
	}
	
	// 基本的輸入驗證
	validInput := regexp.MustCompile(`^[a-zA-Z0-9_\-\./\\]+$`)
	if !validInput.MatchString(command) {
		return "", fmt.Errorf("執行命令失敗: invalid command input")
	}
	for _, arg := range args {
		if !validInput.MatchString(arg) {
			return "", fmt.Errorf("執行命令失敗: invalid argument input")
		}
	}
	
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("執行命令失敗: %v", err)
	}
	return string(output), nil
}

// CommandValidator 簡單的命令驗證器
type CommandValidator struct {
	allowedCommands map[string]bool
}

// ExecuteCommandWithSudo 使用sudo執行命令
func (su *SystemUtils) ExecuteCommandWithSudo(command string, args ...string) (string, error) {
	if su.GetOS() == "windows" {
		return "", fmt.Errorf("Windows不支援sudo")
	}

	sudoArgs := append([]string{command}, args...)
	return su.ExecuteCommand("sudo", sudoArgs...)
}

// ValidationUtils 驗證工具
type ValidationUtils struct{}

// ValidatePinCode 驗證PIN碼格式
func (vu *ValidationUtils) ValidatePinCode(pinCode string) error {
	if len(pinCode) != 3 {
		return fmt.Errorf("PIN碼必須為3位數")
	}

	for _, char := range pinCode {
		if char < '0' || char > '9' {
			return fmt.Errorf("PIN碼只能包含數字")
		}
	}

	return nil
}

// ValidateTimeFormat 驗證時間格式
func (vu *ValidationUtils) ValidateTimeFormat(timeStr string) error {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return fmt.Errorf("時間格式必須為HH:MM")
	}

	hour, err := parseInt(parts[0])
	if err != nil {
		return fmt.Errorf("小時格式錯誤")
	}

	minute, err := parseInt(parts[1])
	if err != nil {
		return fmt.Errorf("分鐘格式錯誤")
	}

	if hour < 0 || hour > 23 {
		return fmt.Errorf("小時必須在0-23之間")
	}

	if minute < 0 || minute > 59 {
		return fmt.Errorf("分鐘必須在0-59之間")
	}

	return nil
}

// ValidateDuration 驗證時間間隔
func (vu *ValidationUtils) ValidateDuration(duration time.Duration) error {
	if duration <= 0 {
		return fmt.Errorf("時間間隔必須大於0")
	}

	if duration > 24*time.Hour {
		return fmt.Errorf("時間間隔不能超過24小時")
	}

	return nil
}

// FileUtils 檔案工具
type FileUtils struct{}

// FileExists 檢查檔案是否存在
func (fu *FileUtils) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// CreateDirIfNotExists 如果目錄不存在則建立
func (fu *FileUtils) CreateDirIfNotExists(path string) error {
	if !fu.FileExists(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// WriteFile 寫入檔案
func (fu *FileUtils) WriteFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("建立檔案失敗: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("寫入檔案失敗: %v", err)
	}

	return nil
}

// ReadFile 讀取檔案
func (fu *FileUtils) ReadFile(path string) (string, error) {
	if strings.Contains(path, "../") || strings.Contains(path, "..\\") {
		return "", fmt.Errorf("Invalid file path")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("讀取檔案失敗: %v", err)
	}
	return string(content), nil
}

// 輔助函數
func parseInt(s string) (int, error) {
	var result int
	for _, char := range s {
		if char < '0' || char > '9' {
			return 0, fmt.Errorf("無效的數字: %c", char)
		}
		result = result*10 + int(char-'0')
	}
	return result, nil
}

// ConfigUtils 設定工具
type ConfigUtils struct{}

// GetEnvWithDefault 取得環境變數，如果不存在則使用預設值
func (cu *ConfigUtils) GetEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetBoolEnvWithDefault 取得布林環境變數
func (cu *ConfigUtils) GetBoolEnvWithDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}

// GetIntEnvWithDefault 取得整數環境變數
func (cu *ConfigUtils) GetIntEnvWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := parseInt(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
