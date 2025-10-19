package network

import (
	"context"
	"fmt"
	"sync"
	"time"

	"pandora_box_console_ids_ips/internal/pubsub"
	pb "pandora_box_console_ids_ips/api/proto/network"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Config contains configuration for the network service
type Config struct {
	DefaultInterface string
	SnapshotLength   int
	Promiscuous      bool
	Timeout          time.Duration
}

// Service implements the NetworkService gRPC service
type Service struct {
	pb.UnimplementedNetworkServiceServer

	config    *Config
	mq        pubsub.MessageQueue
	logger    *logrus.Logger
	sessions  map[string]*MonitorSession
	mu        sync.RWMutex
	startTime time.Time
	metrics   *ServiceMetrics
}

// MonitorSession represents an active monitoring session
type MonitorSession struct {
	SessionID     string
	InterfaceName string
	StartTime     time.Time
	Active        bool
	Statistics    *NetworkStatistics
	Flows         map[string]*FlowInfo
	mu            sync.RWMutex
}

// NetworkStatistics tracks network statistics
type NetworkStatistics struct {
	TotalPackets   int64
	TotalBytes     int64
	TCPPackets     int64
	UDPPackets     int64
	ICMPPackets    int64
	OtherPackets   int64
	DroppedPackets int64
	StartTime      time.Time
	LastUpdate     time.Time
}

// FlowInfo contains information about a network flow
type FlowInfo struct {
	FlowID      string
	SourceIP    string
	DestIP      string
	SourcePort  int32
	DestPort    int32
	Protocol    string
	PacketCount int64
	ByteCount   int64
	FirstSeen   time.Time
	LastSeen    time.Time
	State       pb.FlowState
}

// ServiceMetrics tracks service-level metrics
type ServiceMetrics struct {
	ActiveSessions    int64
	TotalFlows        int64
	AnomaliesDetected int64
	AttacksDetected   int64
	mu                sync.RWMutex
}

// NewService creates a new network service
func NewService(config *Config, mq pubsub.MessageQueue, logger *logrus.Logger) *Service {
	if logger == nil {
		logger = logrus.New()
	}

	return &Service{
		config:    config,
		mq:        mq,
		logger:    logger,
		sessions:  make(map[string]*MonitorSession),
		startTime: time.Now(),
		metrics:   &ServiceMetrics{},
	}
}

// StartMonitoring starts monitoring network traffic
func (s *Service) StartMonitoring(ctx context.Context, req *pb.MonitorRequest) (*pb.MonitorResponse, error) {
	s.logger.Infof("Starting monitoring on interface: %s", req.InterfaceName)

	sessionID := fmt.Sprintf("session_%d", time.Now().UnixNano())

	session := &MonitorSession{
		SessionID:     sessionID,
		InterfaceName: req.InterfaceName,
		StartTime:     time.Now(),
		Active:        true,
		Statistics: &NetworkStatistics{
			StartTime:  time.Now(),
			LastUpdate: time.Now(),
		},
		Flows: make(map[string]*FlowInfo),
	}

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.metrics.ActiveSessions++
	s.mu.Unlock()

	// 啟動監控 goroutine
	go s.monitorTraffic(ctx, session)

	// 發布監控開始事件
	event := pubsub.NewSystemEvent("network-service", "monitoring_started", 
		fmt.Sprintf("Monitoring started on %s", req.InterfaceName))
	event.Metadata["session_id"] = sessionID
	event.Metadata["interface"] = req.InterfaceName
	if message, err := pubsub.ToJSON(event); err == nil {
		s.mq.Publish(ctx, "pandora.events", "system.started", message)
	}

	s.logger.Infof("Monitoring session %s started", sessionID)

	return &pb.MonitorResponse{
		Success:   true,
		Message:   "Monitoring started successfully",
		SessionId: sessionID,
	}, nil
}

// StopMonitoring stops monitoring network traffic
func (s *Service) StopMonitoring(ctx context.Context, req *pb.StopRequest) (*pb.StopResponse, error) {
	s.logger.Infof("Stopping monitoring session: %s", req.SessionId)

	s.mu.Lock()
	session, exists := s.sessions[req.SessionId]
	if !exists {
		s.mu.Unlock()
		return nil, status.Errorf(codes.NotFound, "session not found: %s", req.SessionId)
	}

	session.Active = false
	s.metrics.ActiveSessions--
	s.mu.Unlock()

	// 獲取最終統計
	stats := s.buildStatisticsResponse(session)

	// 發布監控停止事件
	event := pubsub.NewSystemEvent("network-service", "monitoring_stopped",
		fmt.Sprintf("Monitoring stopped for session %s", req.SessionId))
	event.Metadata["session_id"] = req.SessionId
	if message, err := pubsub.ToJSON(event); err == nil {
		s.mq.Publish(ctx, "pandora.events", "system.stopped", message)
	}

	s.logger.Infof("Monitoring session %s stopped", req.SessionId)

	return &pb.StopResponse{
		Success:    true,
		Message:    "Monitoring stopped successfully",
		FinalStats: stats,
	}, nil
}

// GetStatistics gets network statistics
func (s *Service) GetStatistics(ctx context.Context, req *pb.StatsRequest) (*pb.StatsResponse, error) {
	s.mu.RLock()
	session, exists := s.sessions[req.SessionId]
	s.mu.RUnlock()

	if !exists {
		return nil, status.Errorf(codes.NotFound, "session not found: %s", req.SessionId)
	}

	stats := s.buildStatisticsResponse(session)

	// 獲取 top flows
	topFlows := s.getTopFlows(session, 10)

	// 獲取協議統計
	protocolStats := s.getProtocolStatistics(session)

	return &pb.StatsResponse{
		Statistics:    stats,
		TopFlows:      topFlows,
		ProtocolStats: protocolStats,
	}, nil
}

// AnalyzeTraffic analyzes network traffic (streaming)
func (s *Service) AnalyzeTraffic(req *pb.AnalyzeRequest, stream pb.NetworkService_AnalyzeTrafficServer) error {
	s.logger.Debugf("Starting traffic analysis for session: %s", req.SessionId)

	s.mu.RLock()
	session, exists := s.sessions[req.SessionId]
	s.mu.RUnlock()

	if !exists {
		return status.Errorf(codes.NotFound, "session not found: %s", req.SessionId)
	}

	ticker := time.NewTicker(time.Duration(req.WindowSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			// 分析流量模式
			analyses := s.analyzeTrafficPatterns(session, req.AnalysisTypes)

			for _, analysis := range analyses {
				if err := stream.Send(analysis); err != nil {
					s.logger.Errorf("Failed to send analysis: %v", err)
					return err
				}

				// 如果檢測到攻擊，發布事件
				if analysis.Confidence > 0.8 {
					s.publishAttackEvent(stream.Context(), analysis)
				}
			}
		}
	}
}

// DetectAnomalies detects network anomalies (streaming)
func (s *Service) DetectAnomalies(req *pb.AnomalyRequest, stream pb.NetworkService_DetectAnomaliesServer) error {
	s.logger.Debugf("Starting anomaly detection for session: %s", req.SessionId)

	s.mu.RLock()
	session, exists := s.sessions[req.SessionId]
	s.mu.RUnlock()

	if !exists {
		return status.Errorf(codes.NotFound, "session not found: %s", req.SessionId)
	}

	ticker := time.NewTicker(time.Duration(req.WindowSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			// 檢測異常
			anomalies := s.detectAnomalies(session, req.Threshold)

			for _, anomaly := range anomalies {
				if err := stream.Send(anomaly); err != nil {
					s.logger.Errorf("Failed to send anomaly: %v", err)
					return err
				}

				// 發布異常事件
				s.publishAnomalyEvent(stream.Context(), anomaly)
				s.metrics.mu.Lock()
				s.metrics.AnomaliesDetected++
				s.metrics.mu.Unlock()
			}
		}
	}
}

// GetFlowInfo gets information about a specific flow
func (s *Service) GetFlowInfo(ctx context.Context, req *pb.FlowInfoRequest) (*pb.FlowInfoResponse, error) {
	flowKey := fmt.Sprintf("%s:%d-%s:%d-%s",
		req.SourceIp, req.SourcePort,
		req.DestIp, req.DestPort,
		req.Protocol)

	s.mu.RLock()
	defer s.mu.RUnlock()

	// 在所有 session 中查找 flow
	for _, session := range s.sessions {
		session.mu.RLock()
		if flow, exists := session.Flows[flowKey]; exists {
			session.mu.RUnlock()
			return &pb.FlowInfoResponse{
				Flow: &pb.FlowInfo{
					FlowId:      flow.FlowID,
					SourceIp:    flow.SourceIP,
					DestIp:      flow.DestIP,
					SourcePort:  flow.SourcePort,
					DestPort:    flow.DestPort,
					Protocol:    flow.Protocol,
					PacketCount: flow.PacketCount,
					ByteCount:   flow.ByteCount,
					StartTime:   timestamppb.New(flow.FirstSeen),
					EndTime:     timestamppb.New(flow.LastSeen),
					State:       flow.State,
				},
				RecentPackets: []*pb.PacketInfo{}, // TODO: 實現封包歷史
			}, nil
		}
		session.mu.RUnlock()
	}

	return nil, status.Errorf(codes.NotFound, "flow not found")
}

// GetHealth checks the health of the service
func (s *Service) GetHealth(ctx context.Context, _ *emptypb.Empty) (*pb.HealthResponse, error) {
	healthy := s.Health(ctx) == nil

	statusStr := "healthy"
	if !healthy {
		statusStr = "unhealthy"
	}

	dependencies := make(map[string]string)
	if s.mq != nil {
		if err := s.mq.Health(ctx); err != nil {
			dependencies["rabbitmq"] = "unhealthy"
		} else {
			dependencies["rabbitmq"] = "healthy"
		}
	}

	s.metrics.mu.RLock()
	metrics := &pb.NetworkServiceMetrics{
		ActiveSessions:    s.metrics.ActiveSessions,
		TotalFlows:        s.metrics.TotalFlows,
		AnomaliesDetected: s.metrics.AnomaliesDetected,
		AttacksDetected:   s.metrics.AttacksDetected,
	}
	s.metrics.mu.RUnlock()

	return &pb.HealthResponse{
		Healthy:        healthy,
		Status:         statusStr,
		Version:        "1.0.0",
		UptimeSeconds:  int64(time.Since(s.startTime).Seconds()),
		Dependencies:   dependencies,
		Metrics:        metrics,
	}, nil
}

// Health checks if the service is healthy
func (s *Service) Health(ctx context.Context) error {
	if s.mq != nil {
		if err := s.mq.Health(ctx); err != nil {
			return fmt.Errorf("rabbitmq unhealthy: %w", err)
		}
	}
	return nil
}

// Private helper methods

// monitorTraffic monitors network traffic for a session
func (s *Service) monitorTraffic(ctx context.Context, session *MonitorSession) {
	s.logger.Infof("Starting traffic monitoring for session %s", session.SessionID)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !session.Active {
				return
			}

			// TODO: 實際的封包捕獲邏輯（使用 gopacket/libpcap）
			// 這裡是模擬實現
			s.simulateTrafficCapture(session)
		}
	}
}

// simulateTrafficCapture simulates traffic capture (for development)
func (s *Service) simulateTrafficCapture(session *MonitorSession) {
	session.mu.Lock()
	defer session.mu.Unlock()

	// 模擬捕獲一些封包
	session.Statistics.TotalPackets += 100
	session.Statistics.TotalBytes += 65000
	session.Statistics.TCPPackets += 80
	session.Statistics.UDPPackets += 15
	session.Statistics.ICMPPackets += 5
	session.Statistics.LastUpdate = time.Now()

	// 模擬一些流量
	flowKey := "192.168.1.100:54321-10.0.0.1:80-tcp"
	if flow, exists := session.Flows[flowKey]; exists {
		flow.PacketCount += 10
		flow.ByteCount += 6500
		flow.LastSeen = time.Now()
	} else {
		session.Flows[flowKey] = &FlowInfo{
			FlowID:      flowKey,
			SourceIP:    "192.168.1.100",
			DestIP:      "10.0.0.1",
			SourcePort:  54321,
			DestPort:    80,
			Protocol:    "tcp",
			PacketCount: 10,
			ByteCount:   6500,
			FirstSeen:   time.Now(),
			LastSeen:    time.Now(),
			State:       pb.FlowState_FLOW_STATE_ESTABLISHED,
		}
		s.metrics.mu.Lock()
		s.metrics.TotalFlows++
		s.metrics.mu.Unlock()
	}
}

// buildStatisticsResponse builds a statistics response
func (s *Service) buildStatisticsResponse(session *MonitorSession) *pb.NetworkStatistics {
	session.mu.RLock()
	defer session.mu.RUnlock()

	stats := session.Statistics
	duration := time.Since(stats.StartTime).Seconds()
	if duration == 0 {
		duration = 1
	}

	return &pb.NetworkStatistics{
		TotalPackets:     stats.TotalPackets,
		TotalBytes:       stats.TotalBytes,
		TcpPackets:       stats.TCPPackets,
		UdpPackets:       stats.UDPPackets,
		IcmpPackets:      stats.ICMPPackets,
		OtherPackets:     stats.OtherPackets,
		DroppedPackets:   stats.DroppedPackets,
		PacketsPerSecond: float64(stats.TotalPackets) / duration,
		BytesPerSecond:   float64(stats.TotalBytes) / duration,
		StartTime:        timestamppb.New(stats.StartTime),
		LastUpdate:       timestamppb.New(stats.LastUpdate),
	}
}

// getTopFlows returns the top N flows by byte count
func (s *Service) getTopFlows(session *MonitorSession, limit int) []*pb.FlowStatistics {
	session.mu.RLock()
	defer session.mu.RUnlock()

	// TODO: 實現實際的排序邏輯
	var flows []*pb.FlowStatistics

	count := 0
	for _, flow := range session.Flows {
		if count >= limit {
			break
		}

		flows = append(flows, &pb.FlowStatistics{
			SourceIp:    flow.SourceIP,
			DestIp:      flow.DestIP,
			SourcePort:  flow.SourcePort,
			DestPort:    flow.DestPort,
			Protocol:    flow.Protocol,
			PacketCount: flow.PacketCount,
			ByteCount:   flow.ByteCount,
			FirstSeen:   timestamppb.New(flow.FirstSeen),
			LastSeen:    timestamppb.New(flow.LastSeen),
		})
		count++
	}

	return flows
}

// getProtocolStatistics returns statistics per protocol
func (s *Service) getProtocolStatistics(session *MonitorSession) []*pb.ProtocolStatistics {
	session.mu.RLock()
	defer session.mu.RUnlock()

	stats := session.Statistics
	total := float64(stats.TotalPackets)
	if total == 0 {
		total = 1
	}

	return []*pb.ProtocolStatistics{
		{
			Protocol:    "tcp",
			PacketCount: stats.TCPPackets,
			ByteCount:   stats.TCPPackets * 800, // 估算
			Percentage:  float64(stats.TCPPackets) / total * 100,
		},
		{
			Protocol:    "udp",
			PacketCount: stats.UDPPackets,
			ByteCount:   stats.UDPPackets * 600,
			Percentage:  float64(stats.UDPPackets) / total * 100,
		},
		{
			Protocol:    "icmp",
			PacketCount: stats.ICMPPackets,
			ByteCount:   stats.ICMPPackets * 100,
			Percentage:  float64(stats.ICMPPackets) / total * 100,
		},
	}
}

// analyzeTrafficPatterns analyzes traffic for attack patterns
func (s *Service) analyzeTrafficPatterns(session *MonitorSession, analysisTypes []string) []*pb.AnalysisResponse {
	var results []*pb.AnalysisResponse

	session.mu.RLock()
	defer session.mu.RUnlock()

	for _, analysisType := range analysisTypes {
		switch analysisType {
		case "ddos":
			// 檢測 DDoS 攻擊
			if session.Statistics.TotalPackets > 100000 {
				results = append(results, &pb.AnalysisResponse{
					AnalysisType: "ddos",
					Confidence:   0.85,
					Description:  "Potential DDoS attack detected",
					Details: map[string]string{
						"packet_rate": fmt.Sprintf("%.0f pps", 
							float64(session.Statistics.TotalPackets)/time.Since(session.StartTime).Seconds()),
					},
					Timestamp: timestamppb.Now(),
				})
				s.metrics.mu.Lock()
				s.metrics.AttacksDetected++
				s.metrics.mu.Unlock()
			}

		case "port_scan":
			// 檢測端口掃描
			if len(session.Flows) > 100 {
				results = append(results, &pb.AnalysisResponse{
					AnalysisType: "port_scan",
					Confidence:   0.75,
					Description:  "Potential port scan detected",
					Details: map[string]string{
						"unique_flows": fmt.Sprintf("%d", len(session.Flows)),
					},
					Timestamp: timestamppb.Now(),
				})
			}

		case "malware":
			// 檢測惡意軟體通訊
			// TODO: 實現實際的惡意軟體檢測邏輯
		}
	}

	return results
}

// detectAnomalies detects network anomalies
func (s *Service) detectAnomalies(session *MonitorSession, threshold float64) []*pb.AnomalyResponse {
	var anomalies []*pb.AnomalyResponse

	session.mu.RLock()
	defer session.mu.RUnlock()

	// 檢測異常流量
	for _, flow := range session.Flows {
		// 計算異常分數（簡化版）
		anomalyScore := s.calculateAnomalyScore(flow)

		if anomalyScore > threshold {
			anomalies = append(anomalies, &pb.AnomalyResponse{
				AnomalyType:  "unusual_traffic",
				AnomalyScore: anomalyScore,
				SourceIp:     flow.SourceIP,
				DestIp:       flow.DestIP,
				SourcePort:   flow.SourcePort,
				DestPort:     flow.DestPort,
				Protocol:     flow.Protocol,
				Description:  fmt.Sprintf("Unusual traffic pattern detected (score: %.2f)", anomalyScore),
				DetectedAt:   timestamppb.Now(),
			})
		}
	}

	return anomalies
}

// calculateAnomalyScore calculates an anomaly score for a flow
func (s *Service) calculateAnomalyScore(flow *FlowInfo) float64 {
	// 簡化的異常分數計算
	// TODO: 實現基於 ML 的異常檢測

	score := 0.0

	// 高封包率
	duration := flow.LastSeen.Sub(flow.FirstSeen).Seconds()
	if duration > 0 {
		pps := float64(flow.PacketCount) / duration
		if pps > 1000 {
			score += 0.3
		}
	}

	// 大量數據傳輸
	if flow.ByteCount > 10000000 { // > 10MB
		score += 0.2
	}

	// 可疑端口
	suspiciousPorts := map[int32]bool{
		22: true, 23: true, 3389: true, // SSH, Telnet, RDP
	}
	if suspiciousPorts[flow.DestPort] {
		score += 0.2
	}

	return score
}

// publishAttackEvent publishes an attack event to RabbitMQ
func (s *Service) publishAttackEvent(ctx context.Context, analysis *pb.AnalysisResponse) {
	event := pubsub.NewNetworkEvent(analysis.AnalysisType, "", "", "")
	event.Type = pubsub.EventTypeNetworkAttack
	event.EventSubType = analysis.AnalysisType

	if message, err := pubsub.ToJSON(event); err == nil {
		if err := s.mq.Publish(ctx, "pandora.events", "network.attack", message); err != nil {
			s.logger.Errorf("Failed to publish attack event: %v", err)
		}
	}
}

// publishAnomalyEvent publishes an anomaly event to RabbitMQ
func (s *Service) publishAnomalyEvent(ctx context.Context, anomaly *pb.AnomalyResponse) {
	event := pubsub.NewNetworkEvent("anomaly", anomaly.SourceIp, anomaly.DestIp, anomaly.Protocol)
	event.Type = pubsub.EventTypeNetworkAnomaly
	event.SourcePort = int(anomaly.SourcePort)
	event.DestPort = int(anomaly.DestPort)

	if message, err := pubsub.ToJSON(event); err == nil {
		if err := s.mq.Publish(ctx, "pandora.events", "network.anomaly", message); err != nil {
			s.logger.Errorf("Failed to publish anomaly event: %v", err)
		}
	}
}

