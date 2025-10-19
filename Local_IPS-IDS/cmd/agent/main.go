package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pandora_box_console_ids_ips/internal/device"
	"pandora_box_console_ids_ips/internal/grafana"
	"pandora_box_console_ids_ips/internal/mtls"
	"pandora_box_console_ids_ips/internal/network"
	"pandora_box_console_ids_ips/internal/pin"
	"pandora_box_console_ids_ips/internal/token"
	"pandora_box_console_ids_ips/internal/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// main 主函數
func main() {
	// 設定命令列參數
	rootCmd := &cobra.Command{
		Use:   "internet-blocker-agent",
		Short: "網路阻斷器代理程式",
		Long:  `一個基於USB-SERIAL CH340的智慧網路阻斷器代理系統`,
		Run:   runAgent,
	}

	rootCmd.PersistentFlags().String("config", "", "設定檔路徑 (預設: ./agent-config.yaml)")
	rootCmd.PersistentFlags().String("device-port", "/dev/ttyUSB0", "USB-SERIAL 裝置埠號")
	rootCmd.PersistentFlags().String("log-level", "info", "日誌等級 (debug, info, warn, error)")
	rootCmd.PersistentFlags().Duration("timeout", 30*time.Minute, "無訊息超時時間")
	rootCmd.PersistentFlags().String("block-time", "20:00", "每日阻斷時間 (HH:MM)")
	rootCmd.PersistentFlags().String("unlock-time", "08:00", "每日解鎖時間 (HH:MM)")
	rootCmd.PersistentFlags().String("grafana-url", "http://localhost:3000", "Grafana伺服器URL")
	rootCmd.PersistentFlags().String("mtls-cert", "", "mTLS憑證檔案路徑")
	rootCmd.PersistentFlags().String("mtls-key", "", "mTLS私鑰檔案路徑")

	// 綁定環境變數
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		fmt.Printf("綁定命令列參數失敗: %v\n", err)
		os.Exit(1)
	}
	viper.SetEnvPrefix("IB_AGENT")
	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("執行命令失敗: %v\n", err)
		os.Exit(1)
	}
}

// runAgent 執行代理程式
func runAgent(cmd *cobra.Command, args []string) {
	// 初始化日誌
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// 讀取配置
	cfgFile, _ := cmd.Flags().GetString("config")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("agent-config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("/app/configs")
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Warnf("讀取配置檔案失敗: %v，使用預設值", err)
	} else {
		logger.Infof("已載入配置檔案: %s", viper.ConfigFileUsed())
	}

	// 建立代理程式
	agent := NewAgent()

	// 初始化代理程式（在雲端環境可能沒有實體設備）
	if err := agent.Initialize(); err != nil {
		logger.Errorf("初始化代理程式失敗: %v", err)
		// 在雲端環境，即使初始化失敗也繼續運行健康檢查服務器
	}

	// 啟動 HTTP 健康檢查服務器
	go startHealthCheckServer(logger)

	// 創建 context 用於優雅關閉
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 設定信號處理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 在 goroutine 中運行 agent
	go func() {
		if err := agent.Run(ctx); err != nil {
			logger.Errorf("Agent 運行錯誤: %v", err)
		}
	}()

	// 等待中斷信號
	sig := <-sigChan
	logger.Infof("收到信號 %v，正在關閉...", sig)
	cancel()

	// 等待清理完成
	time.Sleep(2 * time.Second)
	logger.Info("Agent 已安全關閉")
}

// startHealthCheckServer 啟動健康檢查 HTTP 服務器
func startHealthCheckServer(logger *logrus.Logger) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// 健康檢查端點
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"status":    "healthy",
			"service":   "pandora-agent",
			"timestamp": time.Now().Format(time.RFC3339),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"%s","service":"%s","timestamp":"%s"}`,
			response["status"], response["service"], response["timestamp"])
	})

	// 根路徑
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"service":"pandora-agent","version":"1.0.0"}`)
	})

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Infof("健康檢查服務器啟動於端口 %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Errorf("健康檢查服務器錯誤: %v", err)
	}
}

// Agent 網路阻斷器代理程式
type Agent struct {
	deviceManager  *device.Manager
	networkManager *network.Manager
	pinSystem      *pin.System
	tokenAuth      *token.Auth
	mtlsClient     *mtls.Client
	grafanaClient  *grafana.Client
	logger         *logrus.Logger
}

// NewAgent 建立新的代理程式
func NewAgent() *Agent {
	return &Agent{}
}

// Initialize 初始化Agent
func (a *Agent) Initialize() error {
	// 初始化日誌
	a.logger = logrus.New()
	a.logger.SetLevel(logrus.InfoLevel)

	// 1. 初始化裝置管理器 (Arduino/ESP)
	deviceManager := device.NewManager(a.logger)
	if err := deviceManager.Initialize(viper.GetString("device-port")); err != nil {
		return fmt.Errorf("裝置初始化失敗: %v", err)
	}
	a.deviceManager = deviceManager

	// 2. 初始化網路管理器
	networkManager := network.NewManager(a.logger)
	if err := networkManager.Initialize(); err != nil {
		return fmt.Errorf("網路管理器初始化失敗: %v", err)
	}
	a.networkManager = networkManager

	// 3. 初始化PIN碼系統
	pinSystem := pin.NewSystem(a.logger)
	if err := pinSystem.Initialize(); err != nil {
		return fmt.Errorf("PIN碼系統初始化失敗: %v", err)
	}
	a.pinSystem = pinSystem

	// 4. 初始化USB Token認證系統
	tokenAuth := token.NewAuth(a.logger)
	if err := tokenAuth.Initialize(); err != nil {
		a.logger.Warnf("USB Token認證系統初始化失敗: %v", err)
	}
	a.tokenAuth = tokenAuth

	// 5. 初始化mTLS客戶端
	mtlsClient := mtls.NewClient(a.logger)
	if err := mtlsClient.Initialize(viper.GetString("mtls-cert"), viper.GetString("mtls-key")); err != nil {
		return fmt.Errorf("mTLS客戶端初始化失敗: %v", err)
	}
	a.mtlsClient = mtlsClient

	// 6. 初始化Grafana客戶端
	grafanaClient := grafana.NewClient(a.logger)
	if err := grafanaClient.Initialize(viper.GetString("grafana-url")); err != nil {
		return fmt.Errorf(" grafana客戶端初始化失敗: %v", err)
	}
	a.grafanaClient = grafanaClient

	a.logger.Info("Agent初始化完成")
	return nil
}

// Run 執行Agent主迴圈
func (a *Agent) Run(ctx context.Context) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastMessageTime := time.Now()
	timeout := viper.GetDuration("timeout")
	blockTime := viper.GetString("block-time")
	unlockTime := viper.GetString("unlock-time")

	// 時間工具
	timeUtils := &utils.TimeUtils{}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			currentTime := time.Now()

			// 檢查是否有新訊息
			if a.deviceManager.HasNewMessage() {
				lastMessageTime = currentTime
				a.logger.Debug("收到IoT裝置訊息，重置計時器")
			}

			// 3.1 檢查超時或阻斷時間
			if shouldBlockNetwork(currentTime, lastMessageTime, timeout, blockTime) {
				if !a.networkManager.IsBlocked() {
					if err := a.handleNetworkBlock(); err != nil {
						a.logger.Errorf("阻斷網路失敗: %v", err)
					}
				}
			}

			// 3.2 檢查解鎖時間
			if timeUtils.IsTimeReached(unlockTime) {
				if a.networkManager.IsBlocked() {
					if err := a.handleUnlockProcess(); err != nil {
						a.logger.Errorf("解鎖程序失敗: %v", err)
					}
				}
			}

			// 定期清理過期Token
			if a.tokenAuth.IsEnabled() {
				a.tokenAuth.CleanExpiredTokens()
			}
		}
	}
}

// handleNetworkBlock 處理網路阻斷
func (a *Agent) handleNetworkBlock() error {
	a.logger.Info("開始網路阻斷程序...")

	// 阻斷網路
	if err := a.networkManager.BlockEthernet(); err != nil {
		return fmt.Errorf("阻斷網路失敗: %v", err)
	}

	// 記錄到Grafana
	if err := a.grafanaClient.LogEvent("network_blocked", "網路已阻斷", map[string]interface{}{
		"reason": "timeout_or_schedule",
		"time":   time.Now().Format(time.RFC3339),
	}); err != nil {
		a.logger.Warnf("記錄到Grafana失敗: %v", err)
	}

	// 發送訊息到IoT裝置顯示"ERROR"
	if err := a.deviceManager.SendCommand("DISPLAY_ERROR"); err != nil {
		a.logger.Warnf("發送錯誤訊息到裝置失敗: %v", err)
	}

	a.logger.Info("網路阻斷完成")
	return nil
}

// handleUnlockProcess 處理解鎖程序
func (a *Agent) handleUnlockProcess() error {
	a.logger.Info("開始解鎖程序...")

	// 1. 發送mTLS訊息到IoT裝置
	if err := a.mtlsClient.SendMessageToDevice(a.deviceManager, "REQUEST_PIN"); err != nil {
		return fmt.Errorf("發送mTLS訊息失敗: %v", err)
	}

	// 2. 等待IoT裝置生成並顯示PIN碼
	pinCode, err := a.waitForPinFromDevice()
	if err != nil {
		return fmt.Errorf("等待PIN碼失敗: %v", err)
	}

	// 3. 等待用戶輸入PIN碼
	userPin, err := a.pinSystem.WaitForPinInput()
	if err != nil {
		return fmt.Errorf("等待PIN碼輸入失敗: %v", err)
	}

	// 4. 驗證PIN碼
	if !a.pinSystem.ValidatePinCode(userPin, pinCode) {
		a.logger.Warn("PIN碼驗證失敗")
		return fmt.Errorf("PIN碼驗證失敗")
	}

	// 5. 可選：驗證USB Token
	if a.tokenAuth.IsEnabled() {
		if !a.tokenAuth.ValidateToken() {
			a.logger.Warn("USB Token驗證失敗")
			return fmt.Errorf("USB Token驗證失敗")
		}
	}

	// 6. 向Grafana伺服器驗證
	if err := a.grafanaClient.VerifyAgent(a.tokenAuth); err != nil {
		a.logger.Warnf("Grafana驗證失敗: %v", err)
		return fmt.Errorf(" grafana驗證失敗: %v", err)
	}

	// 7. 啟用網路
	if err := a.networkManager.EnableEthernet(); err != nil {
		return fmt.Errorf("啟用網路失敗: %v", err)
	}

	// 8. 記錄成功到Grafana
	if err := a.grafanaClient.LogEvent("network_unlocked", "網路已解鎖", map[string]interface{}{
		"success": true,
		"time":    time.Now().Format(time.RFC3339),
	}); err != nil {
		a.logger.Warnf("記錄到Grafana失敗: %v", err)
	}

	// 9. 發送成功訊息到IoT裝置
	if err := a.deviceManager.SendCommand("DISPLAY_SUCCESS"); err != nil {
		a.logger.Warnf("發送成功訊息到裝置失敗: %v", err)
	}

	a.logger.Info("解鎖成功，網路已啟用")
	return nil
}

// waitForPinFromDevice 等待從IoT裝置接收PIN碼
func (a *Agent) waitForPinFromDevice() (string, error) {
	a.logger.Info("等待IoT裝置生成PIN碼...")

	// 發送PIN碼生成請求
	if err := a.deviceManager.SendCommand("GENERATE_PIN"); err != nil {
		return "", fmt.Errorf("發送PIN碼生成請求失敗: %v", err)
	}

	// 等待最多30秒
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return "", fmt.Errorf("等待PIN碼超時")
		case <-ticker.C:
			// 檢查是否有PIN碼回應
			if pinCode := a.deviceManager.GetLastPinCode(); pinCode != "" {
				a.logger.Infof("收到PIN碼: %s", pinCode)
				return pinCode, nil
			}
		}
	}
}

func shouldBlockNetwork(currentTime, lastMessageTime time.Time, timeout time.Duration, blockTime string) bool {
	// 檢查超時
	if currentTime.Sub(lastMessageTime) > timeout {
		return true
	}

	// 檢查阻斷時間
	currentTimeStr := currentTime.Format("15:04")
	return currentTimeStr >= blockTime
}
