package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pandora_box_console_ids_ips/internal/axiom"
	"pandora_box_console_ids_ips/internal/metrics"

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
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// 設定命令列參數
	rootCmd := &cobra.Command{
		Use:   "axiom-ui",
		Short: "Axiom UI 伺服器",
		Long:  `Pandora Box Console IDS-IPS Axiom UI 前端伺服器`,
		Run:   runUIServer,
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "設定檔路徑 (預設: ./ui-config.yaml)")
	rootCmd.PersistentFlags().String("listen-port", "3001", "UI伺服器監聽端口")
	rootCmd.PersistentFlags().String("metrics-port", "8081", "Metrics伺服器監聽端口")
	rootCmd.PersistentFlags().String("log-level", "info", "日誌等級 (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("prometheus-url", "http://prometheus:9090", "Prometheus伺服器URL")
	rootCmd.PersistentFlags().String("grafana-url", "http://grafana:3000", "Grafana伺服器URL")

	// 綁定環境變數
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		logger.Fatalf("綁定命令列參數失敗: %v", err)
	}
	viper.SetEnvPrefix("AXIOM_UI")
	viper.AutomaticEnv()

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

func runUIServer(cmd *cobra.Command, args []string) {
	// 載入設定檔
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("ui-config")
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

	logger.Info("啟動 Axiom UI 伺服器...")

	// 建立上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 設定信號處理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("收到終止信號，正在關閉 UI 伺服器...")
		cancel()
	}()

	// 初始化 Prometheus 指標
	metricsClient := metrics.NewPrometheusMetrics(logger)

	// 啟動 Metrics 伺服器
	go func() {
		metricsPort := viper.GetString("metrics-port")
		logger.Infof("啟動 Metrics 伺服器於端口: %s", metricsPort)
		if err := metricsClient.StartMetricsServer(metricsPort); err != nil {
			logger.Errorf("Metrics 伺服器啟動失敗: %v", err)
		}
	}()

	// 初始化 UI 伺服器
	uiServer := axiom.NewUIServer(logger, metricsClient)

	// 啟動 UI 伺服器
	go func() {
		listenPort := viper.GetString("listen-port")
		if err := uiServer.StartUIServer(listenPort); err != nil {
			logger.Errorf("UI 伺服器啟動失敗: %v", err)
			cancel()
		}
	}()

	// 等待上下文取消
	<-ctx.Done()
	logger.Info("Axiom UI 伺服器已關閉")
}
