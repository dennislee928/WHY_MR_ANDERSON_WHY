package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pandora_box_console_ids_ips/internal/grpc"

	"github.com/sirupsen/logrus"
)

// Orchestrator demonstrates microservices communication
// 編排器示例：展示微服務間的通訊
func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	logger.Info("Starting Microservices Orchestrator Example")

	// 創建所有服務的客戶端
	clients, err := grpc.NewServiceClients(
		getEnv("DEVICE_SERVICE_URL", "localhost:50051"),
		getEnv("NETWORK_SERVICE_URL", "localhost:50052"),
		getEnv("CONTROL_SERVICE_URL", "localhost:50053"),
		logger,
	)
	if err != nil {
		log.Fatalf("Failed to create service clients: %v", err)
	}
	defer clients.Close()

	logger.Info("All service clients connected")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ========== 場景 1: 設備管理 ==========
	logger.Info("\n========== Scenario 1: Device Management ==========")

	// 連接設備
	connectResp, err := clients.Device.Connect(ctx, "usb-001", "/dev/ttyUSB0", 115200)
	if err != nil {
		logger.Errorf("Failed to connect device: %v", err)
	} else {
		logger.Infof("✅ Device connected: %s", connectResp.Message)
	}

	// 查詢設備狀態
	statusResp, err := clients.Device.GetStatus(ctx, "usb-001")
	if err != nil {
		logger.Errorf("Failed to get device status: %v", err)
	} else {
		logger.Infof("✅ Device status: %s (uptime: %.0fs)",
			statusResp.Status, statusResp.Metrics.UptimeSeconds)
	}

	// 列出所有設備
	listResp, err := clients.Device.ListDevices(ctx, false)
	if err != nil {
		logger.Errorf("Failed to list devices: %v", err)
	} else {
		logger.Infof("✅ Total devices: %d", listResp.TotalCount)
		for _, device := range listResp.Devices {
			logger.Infof("   - %s: %s (%s)", device.DeviceId, device.DeviceName, device.Port)
		}
	}

	// ========== 場景 2: 網路監控 ==========
	logger.Info("\n========== Scenario 2: Network Monitoring ==========")

	// 開始監控
	monitorResp, err := clients.Network.StartMonitoring(ctx, "eth0")
	if err != nil {
		logger.Errorf("Failed to start monitoring: %v", err)
	} else {
		logger.Infof("✅ Monitoring started: session %s", monitorResp.SessionId)

		// 等待一段時間收集數據
		time.Sleep(5 * time.Second)

		// 獲取統計數據
		statsResp, err := clients.Network.GetStatistics(ctx, monitorResp.SessionId)
		if err != nil {
			logger.Errorf("Failed to get statistics: %v", err)
		} else {
			stats := statsResp.Statistics
			logger.Infof("✅ Network Statistics:")
			logger.Infof("   - Total Packets: %d", stats.TotalPackets)
			logger.Infof("   - Total Bytes: %d", stats.TotalBytes)
			logger.Infof("   - TCP: %d, UDP: %d, ICMP: %d",
				stats.TcpPackets, stats.UdpPackets, stats.IcmpPackets)
			logger.Infof("   - Rate: %.2f pps, %.2f Bps",
				stats.PacketsPerSecond, stats.BytesPerSecond)

			logger.Infof("✅ Top Flows:")
			for i, flow := range statsResp.TopFlows {
				logger.Infof("   %d. %s:%d -> %s:%d (%s) - %d packets, %d bytes",
					i+1, flow.SourceIp, flow.SourcePort, flow.DestIp, flow.DestPort,
					flow.Protocol, flow.PacketCount, flow.ByteCount)
			}
		}
	}

	// ========== 場景 3: 威脅檢測和阻斷 ==========
	logger.Info("\n========== Scenario 3: Threat Detection and Blocking ==========")

	// 模擬檢測到 DDoS 攻擊
	attackerIP := "192.168.1.100"
	logger.Infof("🚨 DDoS attack detected from %s", attackerIP)

	// 調用 Control Service 阻斷 IP
	blockResp, err := clients.Control.BlockIP(ctx, attackerIP, "DDoS attack detected", 3600)
	if err != nil {
		logger.Errorf("Failed to block IP: %v", err)
	} else {
		logger.Infof("✅ IP blocked: %s (rule: %s)", attackerIP, blockResp.RuleId)
		logger.Infof("   Expires at: %s", blockResp.ExpiresAt.AsTime().Format(time.RFC3339))
	}

	// 查詢阻斷列表
	blockListResp, err := clients.Control.GetBlockList(ctx)
	if err != nil {
		logger.Errorf("Failed to get block list: %v", err)
	} else {
		logger.Infof("✅ Block List (Total: %d, Active: %d):",
			blockListResp.TotalCount, blockListResp.ActiveCount)
		for _, entry := range blockListResp.Entries {
			logger.Infof("   - %s: %s (reason: %s)",
				entry.Type, entry.Value, entry.Reason)
		}
	}

	// ========== 場景 4: 端口控制 ==========
	logger.Info("\n========== Scenario 4: Port Control ==========")

	// 阻斷可疑端口
	suspiciousPort := int32(23) // Telnet
	portBlockResp, err := clients.Control.BlockPort(ctx, suspiciousPort, "tcp", "Suspicious port", 7200)
	if err != nil {
		logger.Errorf("Failed to block port: %v", err)
	} else {
		logger.Infof("✅ Port blocked: %d/tcp (rule: %s)", suspiciousPort, portBlockResp.RuleId)
	}

	// ========== 場景 5: 服務健康檢查 ==========
	logger.Info("\n========== Scenario 5: Service Health Checks ==========")

	// 檢查所有服務健康狀態
	services := []struct {
		name   string
		check  func() error
	}{
		{"Device Service", func() error {
			_, err := clients.Device.ListDevices(ctx, false)
			return err
		}},
		{"Network Service", func() error {
			_, err := clients.Network.GetStatistics(ctx, "test")
			return err
		}},
		{"Control Service", func() error {
			_, err := clients.Control.GetBlockList(ctx)
			return err
		}},
	}

	for _, svc := range services {
		if err := svc.check(); err != nil {
			logger.Warnf("⚠️  %s: unhealthy (%v)", svc.name, err)
		} else {
			logger.Infof("✅ %s: healthy", svc.name)
		}
	}

	// ========== 等待中斷信號 ==========
	logger.Info("\n========== Orchestrator Running ==========")
	logger.Info("Press Ctrl+C to stop...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("\nShutting down...")

	// 清理：解除阻斷
	if _, err := clients.Control.UnblockIP(ctx, attackerIP, "Manual cleanup"); err != nil {
		logger.Errorf("Failed to unblock IP: %v", err)
	} else {
		logger.Infof("✅ IP %s unblocked", attackerIP)
	}

	logger.Info("Orchestrator stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

