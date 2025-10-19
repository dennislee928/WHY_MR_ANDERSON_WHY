// +build performance

package performance

import (
	"context"
	"fmt"
	"testing"
	"time"

	"pandora_box_console_ids_ips/internal/grpc"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

// BenchmarkDeviceServiceGetStatus benchmarks Device Service GetStatus RPC
func BenchmarkDeviceServiceGetStatus(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	config := grpc.DefaultClientConfig("localhost:50051")
	client, err := grpc.NewDeviceClient(config, logger)
	require.NoError(b, err)
	defer client.Close()

	ctx := context.Background()

	// 先連接設備
	_, err = client.Connect(ctx, "bench-device", "/dev/ttyUSB0", 115200)
	require.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.GetStatus(ctx, "bench-device")
			if err != nil {
				b.Errorf("GetStatus failed: %v", err)
			}
		}
	})
}

// BenchmarkDeviceServiceListDevices benchmarks Device Service ListDevices RPC
func BenchmarkDeviceServiceListDevices(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	config := grpc.DefaultClientConfig("localhost:50051")
	client, err := grpc.NewDeviceClient(config, logger)
	require.NoError(b, err)
	defer client.Close()

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.ListDevices(ctx, false)
			if err != nil {
				b.Errorf("ListDevices failed: %v", err)
			}
		}
	})
}

// BenchmarkNetworkServiceGetStatistics benchmarks Network Service GetStatistics RPC
func BenchmarkNetworkServiceGetStatistics(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	config := grpc.DefaultClientConfig("localhost:50052")
	client, err := grpc.NewNetworkClient(config, logger)
	require.NoError(b, err)
	defer client.Close()

	ctx := context.Background()

	// 先開始監控
	monitorResp, err := client.StartMonitoring(ctx, "eth0")
	require.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.GetStatistics(ctx, monitorResp.SessionId)
			if err != nil {
				b.Errorf("GetStatistics failed: %v", err)
			}
		}
	})
}

// BenchmarkControlServiceBlockIP benchmarks Control Service BlockIP RPC
func BenchmarkControlServiceBlockIP(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	config := grpc.DefaultClientConfig("localhost:50053")
	client, err := grpc.NewControlClient(config, logger)
	require.NoError(b, err)
	defer client.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ip := fmt.Sprintf("192.168.%d.%d", i/256, i%256)
		_, err := client.BlockIP(ctx, ip, "benchmark test", 60)
		if err != nil {
			b.Errorf("BlockIP failed: %v", err)
		}
	}
}

// BenchmarkControlServiceGetBlockList benchmarks Control Service GetBlockList RPC
func BenchmarkControlServiceGetBlockList(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	config := grpc.DefaultClientConfig("localhost:50053")
	client, err := grpc.NewControlClient(config, logger)
	require.NoError(b, err)
	defer client.Close()

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := client.GetBlockList(ctx)
			if err != nil {
				b.Errorf("GetBlockList failed: %v", err)
			}
		}
	})
}

// BenchmarkEndToEndFlow benchmarks a complete end-to-end flow
func BenchmarkEndToEndFlow(b *testing.B) {
	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)

	clients, err := grpc.NewServiceClients(
		"localhost:50051",
		"localhost:50052",
		"localhost:50053",
		logger,
	)
	require.NoError(b, err)
	defer clients.Close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// 1. 查詢設備狀態
		_, err := clients.Device.GetStatus(ctx, "bench-device")
		if err != nil {
			b.Errorf("GetStatus failed: %v", err)
		}

		// 2. 獲取網路統計
		_, err = clients.Network.GetStatistics(ctx, "test-session")
		if err != nil {
			// 預期可能失敗（session 不存在）
		}

		// 3. 查詢阻斷列表
		_, err = clients.Control.GetBlockList(ctx)
		if err != nil {
			b.Errorf("GetBlockList failed: %v", err)
		}
	}
}

// TestLoadTest performs a load test on all services
func TestLoadTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	clients, err := grpc.NewServiceClients(
		"localhost:50051",
		"localhost:50052",
		"localhost:50053",
		logger,
	)
	require.NoError(t, err)
	defer clients.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 測試參數
	duration := 30 * time.Second
	concurrency := 100
	requestsPerGoroutine := 100

	t.Logf("Starting load test: duration=%v, concurrency=%d", duration, concurrency)

	// 統計
	type stats struct {
		requests int64
		errors   int64
		latency  []time.Duration
	}

	deviceStats := &stats{}
	networkStats := &stats{}
	controlStats := &stats{}

	// 啟動並發請求
	done := make(chan bool, concurrency*3)

	// Device Service 負載
	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < requestsPerGoroutine; j++ {
				start := time.Now()
				_, err := clients.Device.GetStatus(ctx, "bench-device")
				latency := time.Since(start)

				deviceStats.requests++
				deviceStats.latency = append(deviceStats.latency, latency)
				if err != nil {
					deviceStats.errors++
				}
			}
			done <- true
		}()
	}

	// Network Service 負載
	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < requestsPerGoroutine; j++ {
				start := time.Now()
				_, err := clients.Network.GetStatistics(ctx, "test-session")
				latency := time.Since(start)

				networkStats.requests++
				networkStats.latency = append(networkStats.latency, latency)
				if err != nil {
					networkStats.errors++
				}
			}
			done <- true
		}()
	}

	// Control Service 負載
	for i := 0; i < concurrency; i++ {
		go func() {
			for j := 0; j < requestsPerGoroutine; j++ {
				start := time.Now()
				_, err := clients.Control.GetBlockList(ctx)
				latency := time.Since(start)

				controlStats.requests++
				controlStats.latency = append(controlStats.latency, latency)
				if err != nil {
					controlStats.errors++
				}
			}
			done <- true
		}()
	}

	// 等待所有請求完成
	for i := 0; i < concurrency*3; i++ {
		<-done
	}

	// 輸出結果
	t.Logf("\n========== Load Test Results ==========")
	t.Logf("Device Service:")
	t.Logf("  Requests: %d, Errors: %d (%.2f%%)", 
		deviceStats.requests, deviceStats.errors,
		float64(deviceStats.errors)/float64(deviceStats.requests)*100)

	t.Logf("Network Service:")
	t.Logf("  Requests: %d, Errors: %d (%.2f%%)",
		networkStats.requests, networkStats.errors,
		float64(networkStats.errors)/float64(networkStats.requests)*100)

	t.Logf("Control Service:")
	t.Logf("  Requests: %d, Errors: %d (%.2f%%)",
		controlStats.requests, controlStats.errors,
		float64(controlStats.errors)/float64(controlStats.requests)*100)
}

