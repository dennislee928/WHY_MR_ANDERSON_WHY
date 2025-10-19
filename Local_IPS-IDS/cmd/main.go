package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pandora_box_console_ids_ips/internal/device"
	"pandora_box_console_ids_ips/internal/network"
	"pandora_box_console_ids_ips/internal/pin"
	"pandora_box_console_ids_ips/internal/token"
	"pandora_box_console_ids_ips/internal/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	logger  *logrus.Logger
)

func main() {
	// 初始化日誌
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// 設定命令列參數
	rootCmd := &cobra.Command{
		Use:   "internet-blocker",
		Short: "網路阻斷器主程式",
		Long:  `一個基於USB-SERIAL CH340的智慧網路阻斷器系統`,
		Run:   run,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "設定檔路徑 (預設: ./config.yaml)")
	rootCmd.PersistentFlags().String("device-port", "/dev/ttyUSB0", "USB-SERIAL 裝置埠號")
	rootCmd.PersistentFlags().String("log-level", "info", "日誌等級 (debug, info, warn, error)")
	rootCmd.PersistentFlags().Duration("timeout", 30*time.Minute, "無訊息超時時間")
	rootCmd.PersistentFlags().String("block-time", "20:00", "每日阻斷時間 (HH:MM)")
	rootCmd.PersistentFlags().String("unlock-time", "08:00", "每日解鎖時間 (HH:MM)")

	// 綁定環境變數
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		logger.Fatalf("綁定命令列參數失敗: %v", err)
	}
	viper.SetEnvPrefix("IB")
	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	// 載入設定檔
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Warnf("無法讀取設定檔: %v", err)
	}

	// 設定日誌等級
	logLevel, err := logrus.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	logger.Info("啟動網路阻斷器系統...")

	// 建立上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 設定信號處理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("收到終止信號，正在關閉系統...")
		cancel()
	}()

	// 初始化各個模組
	if err := initializeModules(ctx); err != nil {
		logger.Fatalf("模組初始化失敗: %v", err)
	}

	// 等待上下文取消
	<-ctx.Done()
	logger.Info("系統已關閉")
}

func initializeModules(ctx context.Context) error {
	// 1. 初始化裝置管理器
	deviceManager := device.NewManager(logger)
	if err := deviceManager.Initialize(viper.GetString("device-port")); err != nil {
		return fmt.Errorf("裝置初始化失敗: %v", err)
	}

	// 2. 確認產品類型
	detector := device.NewDetector(deviceManager)
	productType, err := detector.DetectProductType()
	if err != nil {
		logger.Warnf("產品類型偵測失敗: %v", err)
	} else {
		logger.Infof("偵測到產品類型: %s", productType.String())
	}

	// 3. 初始化網路管理器
	networkManager := network.NewManager(logger)
	if err := networkManager.Initialize(); err != nil {
		return fmt.Errorf("網路管理器初始化失敗: %v", err)
	}

	// 4. 初始化PIN碼系統
	pinSystem := pin.NewSystem(logger)
	if err := pinSystem.Initialize(); err != nil {
		return fmt.Errorf("PIN碼系統初始化失敗: %v", err)
	}

	// 5. 初始化USB Token認證系統
	tokenAuth := token.NewAuth(logger)
	if err := tokenAuth.Initialize(); err != nil {
		logger.Warnf("USB Token認證系統初始化失敗: %v", err)
	}

	// 6. 啟動主控制迴圈
	go runMainLoop(ctx, deviceManager, networkManager, pinSystem, tokenAuth)

	return nil
}

func runMainLoop(ctx context.Context, deviceMgr *device.Manager, networkMgr *network.Manager, pinSys *pin.System, tokenAuth *token.Auth) {
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
			return
		case <-ticker.C:
			currentTime := time.Now()

			// 檢查是否有新訊息
			if deviceMgr.HasNewMessage() {
				lastMessageTime = currentTime
				logger.Debug("收到裝置訊息，重置計時器")
			}

			// 3.1 檢查超時或阻斷時間
			if shouldBlockNetwork(currentTime, lastMessageTime, timeout, blockTime) {
				if !networkMgr.IsBlocked() {
					if err := networkMgr.BlockEthernet(); err != nil {
						logger.Errorf("阻斷網路失敗: %v", err)
					} else {
						logger.Info("網路已阻斷")
					}
				}
			}

			// 3.2 檢查解鎖時間
			if timeUtils.IsTimeReached(unlockTime) {
				if networkMgr.IsBlocked() {
					if err := handleUnlockProcess(deviceMgr, networkMgr, pinSys, tokenAuth); err != nil {
						logger.Errorf("解鎖程序失敗: %v", err)
					}
				}
			}

			// 定期清理過期Token
			if tokenAuth.IsEnabled() {
				tokenAuth.CleanExpiredTokens()
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

func handleUnlockProcess(deviceMgr *device.Manager, networkMgr *network.Manager, pinSys *pin.System, tokenAuth *token.Auth) error {
	logger.Info("開始解鎖程序...")

	// 生成並顯示PIN碼
	pinCode := pinSys.GeneratePinCode()
	if err := deviceMgr.DisplayPinCode(pinCode); err != nil {
		return fmt.Errorf("顯示PIN碼失敗: %v", err)
	}

	// 等待用戶輸入PIN碼
	userPin, err := pinSys.WaitForPinInput()
	if err != nil {
		return fmt.Errorf("等待PIN碼輸入失敗: %v", err)
	}

	// 驗證PIN碼
	if !pinSys.ValidatePinCode(userPin, pinCode) {
		logger.Warn("PIN碼驗證失敗")
		return fmt.Errorf("PIN碼驗證失敗")
	}

	// 可選：驗證USB Token
	if tokenAuth.IsEnabled() {
		if !tokenAuth.ValidateToken() {
			logger.Warn("USB Token驗證失敗")
			return fmt.Errorf("USB Token驗證失敗")
		}
	}

	// 啟用網路
	if err := networkMgr.EnableEthernet(); err != nil {
		return fmt.Errorf("啟用網路失敗: %v", err)
	}

	logger.Info("解鎖成功，網路已啟用")
	return nil
}
