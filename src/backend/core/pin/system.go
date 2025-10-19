package pin

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// System PIN碼系統
type System struct {
	logger     *logrus.Logger
	currentPin string
	pinExpiry  time.Time
	mu         sync.RWMutex
	inputChan  chan string
	isWaiting  bool
}

// NewSystem 建立新的PIN碼系統
func NewSystem(logger *logrus.Logger) *System {
	return &System{
		logger:    logger,
		inputChan: make(chan string, 1),
	}
}

// Initialize 初始化PIN碼系統
func (s *System) Initialize() error {
	s.logger.Info("PIN碼系統初始化完成")
	return nil
}

// GeneratePinCode 生成3位數PIN碼
func (s *System) GeneratePinCode() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 生成100-999之間的隨機數
	max := big.NewInt(900)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		s.logger.Errorf("生成隨機數失敗: %v", err)
		// 使用時間戳作為備用方案
		n = big.NewInt(int64(time.Now().Nanosecond() % 900))
	}

	pinCode := fmt.Sprintf("%03d", n.Int64()+100)
	s.currentPin = pinCode
	s.pinExpiry = time.Now().Add(10 * time.Minute) // PIN碼10分鐘後過期

	s.logger.Infof("生成新PIN碼: %s (過期時間: %s)", pinCode, s.pinExpiry.Format("15:04:05"))
	return pinCode
}

// ValidatePinCode 驗證PIN碼
func (s *System) ValidatePinCode(inputPin, expectedPin string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 檢查PIN碼是否過期
	if time.Now().After(s.pinExpiry) {
		s.logger.Warn("PIN碼已過期")
		return false
	}

	// 簡單的字串比較（實際應用中可能需要更複雜的驗證）
	isValid := inputPin == expectedPin

	if isValid {
		s.logger.Info("PIN碼驗證成功")
	} else {
		s.logger.Warnf("PIN碼驗證失敗: 輸入=%s, 期望=%s", inputPin, expectedPin)
	}

	return isValid
}

// WaitForPinInput 等待用戶輸入PIN碼
func (s *System) WaitForPinInput() (string, error) {
	s.mu.Lock()
	s.isWaiting = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.isWaiting = false
		s.mu.Unlock()
	}()

	s.logger.Info("等待用戶輸入PIN碼...")
	fmt.Print("請輸入3位數PIN碼: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("讀取輸入失敗: %v", err)
	}

	// 清理輸入（移除換行符和空白）
	pinCode := strings.TrimSpace(input)

	// 驗證PIN碼格式
	if !s.isValidPinFormat(pinCode) {
		return "", fmt.Errorf("PIN碼格式無效: %s", pinCode)
	}

	s.logger.Infof("用戶輸入PIN碼: %s", pinCode)
	return pinCode, nil
}

// isValidPinFormat 驗證PIN碼格式
func (s *System) isValidPinFormat(pinCode string) bool {
	// 檢查是否為3位數字
	if len(pinCode) != 3 {
		return false
	}

	for _, char := range pinCode {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// GetCurrentPin 取得當前PIN碼
func (s *System) GetCurrentPin() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentPin
}

// IsPinExpired 檢查PIN碼是否過期
func (s *System) IsPinExpired() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return time.Now().After(s.pinExpiry)
}

// IsWaitingForInput 檢查是否正在等待輸入
func (s *System) IsWaitingForInput() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isWaiting
}

// ClearPin 清除當前PIN碼
func (s *System) ClearPin() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentPin = ""
	s.pinExpiry = time.Time{}
	s.logger.Info("PIN碼已清除")
}

// GenerateSecurePin 生成安全PIN碼（使用加密隨機數）
func (s *System) GenerateSecurePin() (string, error) {
	// 生成加密安全的隨機數
	bytes := make([]byte, 1)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("生成加密隨機數失敗: %v", err)
	}

	// 確保在100-999範圍內
	pinValue := int(bytes[0])%900 + 100
	pinCode := fmt.Sprintf("%03d", pinValue)

	s.logger.Infof("生成安全PIN碼: %s", pinCode)
	return pinCode, nil
}

// HashPin 對PIN碼進行雜湊（用於安全儲存）
func (s *System) HashPin(pinCode string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pinCode), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("PIN碼雜湊失敗: %v", err)
	}
	return string(hashedBytes), nil
}

// VerifyHashedPin 驗證雜湊後的PIN碼
func (s *System) VerifyHashedPin(pinCode, hashedPin string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPin), []byte(pinCode))
	return err == nil
}

// SetPinExpiry 設定PIN碼過期時間
func (s *System) SetPinExpiry(expiry time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pinExpiry = time.Now().Add(expiry)
	s.logger.Infof("PIN碼過期時間設定為: %s", s.pinExpiry.Format("15:04:05"))
}
