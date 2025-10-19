// +build windows

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	
	"pandora_box_console_ids_ips/internal/windows"
)

func main() {
	// 命令行參數
	axiomURL := flag.String("axiom-url", "http://localhost:3001", "Axiom Backend URL")
	agentID := flag.String("agent-id", "", "Agent ID")
	pollInterval := flag.Duration("poll-interval", 30, "Poll interval in seconds")
	batchSize := flag.Int("batch-size", 100, "Batch size for log collection")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// 初始化 Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 如果沒有指定 Agent ID，使用機器名
	if *agentID == "" {
		hostname, _ := os.Hostname()
		*agentID = hostname
	}

	logger.Infof("Starting Windows Event Log Agent...")
	logger.Infof("Agent ID: %s", *agentID)
	logger.Infof("Axiom Backend URL: %s", *axiomURL)
	logger.Infof("Poll Interval: %v", *pollInterval)

	// 創建收集器
	collector := windows.NewModernEventLogCollector(logger)
	collector.SetPollInterval(*pollInterval)
	collector.SetBatchSize(*batchSize)

	// 創建上傳器
	uploader := windows.NewEventLogUploader(logger, *axiomURL, *agentID)

	// 開始收集
	go collector.StartCollection(func(logs []windows.WindowsEventLog) error {
		return uploader.UploadLogs(logs)
	})

	logger.Info("Windows Event Log Agent started successfully")

	// 等待中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down agent...")
}



package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	
	"pandora_box_console_ids_ips/internal/windows"
)

func main() {
	// 命令行參數
	axiomURL := flag.String("axiom-url", "http://localhost:3001", "Axiom Backend URL")
	agentID := flag.String("agent-id", "", "Agent ID")
	pollInterval := flag.Duration("poll-interval", 30, "Poll interval in seconds")
	batchSize := flag.Int("batch-size", 100, "Batch size for log collection")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// 初始化 Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 如果沒有指定 Agent ID，使用機器名
	if *agentID == "" {
		hostname, _ := os.Hostname()
		*agentID = hostname
	}

	logger.Infof("Starting Windows Event Log Agent...")
	logger.Infof("Agent ID: %s", *agentID)
	logger.Infof("Axiom Backend URL: %s", *axiomURL)
	logger.Infof("Poll Interval: %v", *pollInterval)

	// 創建收集器
	collector := windows.NewModernEventLogCollector(logger)
	collector.SetPollInterval(*pollInterval)
	collector.SetBatchSize(*batchSize)

	// 創建上傳器
	uploader := windows.NewEventLogUploader(logger, *axiomURL, *agentID)

	// 開始收集
	go collector.StartCollection(func(logs []windows.WindowsEventLog) error {
		return uploader.UploadLogs(logs)
	})

	logger.Info("Windows Event Log Agent started successfully")

	// 等待中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down agent...")
}


package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	
	"pandora_box_console_ids_ips/internal/windows"
)

func main() {
	// 命令行參數
	axiomURL := flag.String("axiom-url", "http://localhost:3001", "Axiom Backend URL")
	agentID := flag.String("agent-id", "", "Agent ID")
	pollInterval := flag.Duration("poll-interval", 30, "Poll interval in seconds")
	batchSize := flag.Int("batch-size", 100, "Batch size for log collection")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// 初始化 Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 如果沒有指定 Agent ID，使用機器名
	if *agentID == "" {
		hostname, _ := os.Hostname()
		*agentID = hostname
	}

	logger.Infof("Starting Windows Event Log Agent...")
	logger.Infof("Agent ID: %s", *agentID)
	logger.Infof("Axiom Backend URL: %s", *axiomURL)
	logger.Infof("Poll Interval: %v", *pollInterval)

	// 創建收集器
	collector := windows.NewModernEventLogCollector(logger)
	collector.SetPollInterval(*pollInterval)
	collector.SetBatchSize(*batchSize)

	// 創建上傳器
	uploader := windows.NewEventLogUploader(logger, *axiomURL, *agentID)

	// 開始收集
	go collector.StartCollection(func(logs []windows.WindowsEventLog) error {
		return uploader.UploadLogs(logs)
	})

	logger.Info("Windows Event Log Agent started successfully")

	// 等待中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down agent...")
}



package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	
	"pandora_box_console_ids_ips/internal/windows"
)

func main() {
	// 命令行參數
	axiomURL := flag.String("axiom-url", "http://localhost:3001", "Axiom Backend URL")
	agentID := flag.String("agent-id", "", "Agent ID")
	pollInterval := flag.Duration("poll-interval", 30, "Poll interval in seconds")
	batchSize := flag.Int("batch-size", 100, "Batch size for log collection")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	// 初始化 Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// 如果沒有指定 Agent ID，使用機器名
	if *agentID == "" {
		hostname, _ := os.Hostname()
		*agentID = hostname
	}

	logger.Infof("Starting Windows Event Log Agent...")
	logger.Infof("Agent ID: %s", *agentID)
	logger.Infof("Axiom Backend URL: %s", *axiomURL)
	logger.Infof("Poll Interval: %v", *pollInterval)

	// 創建收集器
	collector := windows.NewModernEventLogCollector(logger)
	collector.SetPollInterval(*pollInterval)
	collector.SetBatchSize(*batchSize)

	// 創建上傳器
	uploader := windows.NewEventLogUploader(logger, *axiomURL, *agentID)

	// 開始收集
	go collector.StartCollection(func(logs []windows.WindowsEventLog) error {
		return uploader.UploadLogs(logs)
	})

	logger.Info("Windows Event Log Agent started successfully")

	// 等待中斷信號
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down agent...")
}

