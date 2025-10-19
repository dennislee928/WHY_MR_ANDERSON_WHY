package control

import (
	"context"
	"fmt"
	"sync"
	"time"

	"pandora_box_console_ids_ips/internal/pubsub"
	pb "pandora_box_console_ids_ips/api/proto/control"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Config contains configuration for the control service
type Config struct {
	FirewallBackend    string        // "iptables", "nftables", or "mock"
	DefaultBlockAction string        // "drop" or "reject"
	MaxBlockDuration   time.Duration
}

// Service implements the ControlService gRPC service
type Service struct {
	pb.UnimplementedControlServiceServer

	config       *Config
	mq           pubsub.MessageQueue
	logger       *logrus.Logger
	blockList    map[string]*BlockEntry
	firewallRules map[string]*FirewallRule
	mu           sync.RWMutex
	startTime    time.Time
	metrics      *ServiceMetrics
}

// BlockEntry represents a blocked item
type BlockEntry struct {
	EntryID   string
	Type      string // "ip" or "port"
	Value     string
	Reason    string
	Action    string
	CreatedAt time.Time
	ExpiresAt time.Time
	IsActive  bool
	CreatedBy string
	Metadata  map[string]string
}

// FirewallRule represents a firewall rule
type FirewallRule struct {
	RuleID      string
	Name        string
	Action      string
	SourceIP    string
	DestIP      string
	SourcePort  int32
	DestPort    int32
	Protocol    string
	Direction   string
	Priority    int32
	Enabled     bool
	Description string
	CreatedAt   time.Time
	CreatedBy   string
}

// ServiceMetrics tracks service-level metrics
type ServiceMetrics struct {
	ActiveBlocks   int64
	TotalBlocks    int64
	TotalUnblocks  int64
	ActiveRules    int64
	PacketsBlocked int64
	mu             sync.RWMutex
}

// NewService creates a new control service
func NewService(config *Config, mq pubsub.MessageQueue, logger *logrus.Logger) *Service {
	if logger == nil {
		logger = logrus.New()
	}

	service := &Service{
		config:        config,
		mq:            mq,
		logger:        logger,
		blockList:     make(map[string]*BlockEntry),
		firewallRules: make(map[string]*FirewallRule),
		startTime:     time.Now(),
		metrics:       &ServiceMetrics{},
	}

	// 啟動過期清理 goroutine
	go service.cleanupExpiredBlocks()

	return service
}

// BlockIP blocks an IP address
func (s *Service) BlockIP(ctx context.Context, req *pb.BlockIPRequest) (*pb.BlockIPResponse, error) {
	s.logger.Infof("Blocking IP: %s (reason: %s)", req.IpAddress, req.Reason)

	ruleID := fmt.Sprintf("block_ip_%s_%d", req.IpAddress, time.Now().Unix())

	expiresAt := time.Time{}
	if req.DurationSeconds > 0 {
		expiresAt = time.Now().Add(time.Duration(req.DurationSeconds) * time.Second)
	}

	entry := &BlockEntry{
		EntryID:   ruleID,
		Type:      "ip",
		Value:     req.IpAddress,
		Reason:    req.Reason,
		Action:    req.Action.String(),
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		IsActive:  true,
		CreatedBy: "system",
		Metadata:  req.Metadata,
	}

	s.mu.Lock()
	s.blockList[ruleID] = entry
	s.metrics.ActiveBlocks++
	s.metrics.TotalBlocks++
	s.mu.Unlock()

	// TODO: 實際的 iptables/nftables 規則應用
	if err := s.applyIPBlock(req.IpAddress, req.Action); err != nil {
		s.logger.Errorf("Failed to apply IP block: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to apply block: %v", err)
	}

	// 發布阻斷事件
	event := pubsub.NewNetworkEvent("ip_blocked", req.IpAddress, "", "")
	event.Type = pubsub.EventTypeNetworkBlocked
	event.Metadata["reason"] = req.Reason
	event.Metadata["duration"] = fmt.Sprintf("%d", req.DurationSeconds)
	if message, err := pubsub.ToJSON(event); err == nil {
		s.mq.Publish(ctx, "pandora.events", "network.blocked", message)
	}

	s.logger.Infof("IP %s blocked successfully (rule: %s)", req.IpAddress, ruleID)

	return &pb.BlockIPResponse{
		Success:   true,
		Message:   "IP blocked successfully",
		RuleId:    ruleID,
		ExpiresAt: timestamppb.New(expiresAt),
	}, nil
}

// UnblockIP unblocks an IP address
func (s *Service) UnblockIP(ctx context.Context, req *pb.UnblockIPRequest) (*pb.UnblockIPResponse, error) {
	s.logger.Infof("Unblocking IP: %s (reason: %s)", req.IpAddress, req.Reason)

	s.mu.Lock()
	defer s.mu.Unlock()

	rulesRemoved := 0
	for ruleID, entry := range s.blockList {
		if entry.Type == "ip" && entry.Value == req.IpAddress && entry.IsActive {
			entry.IsActive = false
			rulesRemoved++
			s.metrics.ActiveBlocks--
			s.metrics.TotalUnblocks++

			// TODO: 實際移除 iptables/nftables 規則
			if err := s.removeIPBlock(req.IpAddress); err != nil {
				s.logger.Errorf("Failed to remove IP block: %v", err)
			}

			delete(s.blockList, ruleID)
		}
	}

	if rulesRemoved == 0 {
		return nil, status.Errorf(codes.NotFound, "no active blocks found for IP: %s", req.IpAddress)
	}

	// 發布解除阻斷事件
	event := pubsub.NewNetworkEvent("ip_unblocked", req.IpAddress, "", "")
	event.Metadata["reason"] = req.Reason
	if message, err := pubsub.ToJSON(event); err == nil {
		s.mq.Publish(ctx, "pandora.events", "network.unblocked", message)
	}

	s.logger.Infof("IP %s unblocked (%d rules removed)", req.IpAddress, rulesRemoved)

	return &pb.UnblockIPResponse{
		Success:      true,
		Message:      "IP unblocked successfully",
		RulesRemoved: int32(rulesRemoved),
	}, nil
}

// BlockPort blocks a port
func (s *Service) BlockPort(ctx context.Context, req *pb.BlockPortRequest) (*pb.BlockPortResponse, error) {
	s.logger.Infof("Blocking port: %d/%s (reason: %s)", req.Port, req.Protocol, req.Reason)

	ruleID := fmt.Sprintf("block_port_%d_%s_%d", req.Port, req.Protocol, time.Now().Unix())

	expiresAt := time.Time{}
	if req.DurationSeconds > 0 {
		expiresAt = time.Now().Add(time.Duration(req.DurationSeconds) * time.Second)
	}

	entry := &BlockEntry{
		EntryID:   ruleID,
		Type:      "port",
		Value:     fmt.Sprintf("%d/%s", req.Port, req.Protocol),
		Reason:    req.Reason,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		IsActive:  true,
		CreatedBy: "system",
	}

	s.mu.Lock()
	s.blockList[ruleID] = entry
	s.metrics.ActiveBlocks++
	s.metrics.TotalBlocks++
	s.mu.Unlock()

	// TODO: 實際的端口阻斷邏輯
	if err := s.applyPortBlock(req.Port, req.Protocol); err != nil {
		s.logger.Errorf("Failed to apply port block: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to apply block: %v", err)
	}

	s.logger.Infof("Port %d/%s blocked successfully (rule: %s)", req.Port, req.Protocol, ruleID)

	return &pb.BlockPortResponse{
		Success: true,
		Message: "Port blocked successfully",
		RuleId:  ruleID,
	}, nil
}

// UnblockPort unblocks a port
func (s *Service) UnblockPort(ctx context.Context, req *pb.UnblockPortRequest) (*pb.UnblockPortResponse, error) {
	s.logger.Infof("Unblocking port: %d/%s", req.Port, req.Protocol)

	s.mu.Lock()
	defer s.mu.Unlock()

	portValue := fmt.Sprintf("%d/%s", req.Port, req.Protocol)
	rulesRemoved := 0

	for ruleID, entry := range s.blockList {
		if entry.Type == "port" && entry.Value == portValue && entry.IsActive {
			entry.IsActive = false
			rulesRemoved++
			s.metrics.ActiveBlocks--
			s.metrics.TotalUnblocks++

			// TODO: 實際移除端口阻斷
			if err := s.removePortBlock(req.Port, req.Protocol); err != nil {
				s.logger.Errorf("Failed to remove port block: %v", err)
			}

			delete(s.blockList, ruleID)
		}
	}

	if rulesRemoved == 0 {
		return nil, status.Errorf(codes.NotFound, "no active blocks found for port: %d/%s", req.Port, req.Protocol)
	}

	s.logger.Infof("Port %d/%s unblocked (%d rules removed)", req.Port, req.Protocol, rulesRemoved)

	return &pb.UnblockPortResponse{
		Success:      true,
		Message:      "Port unblocked successfully",
		RulesRemoved: int32(rulesRemoved),
	}, nil
}

// GetBlockList gets the current block list
func (s *Service) GetBlockList(ctx context.Context, req *pb.BlockListRequest) (*pb.BlockListResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var entries []*pb.BlockEntry
	activeCount := int32(0)

	for _, entry := range s.blockList {
		// 過濾類型
		if req.Type != pb.BlockType_BLOCK_TYPE_UNKNOWN {
			if (req.Type == pb.BlockType_BLOCK_TYPE_IP && entry.Type != "ip") ||
				(req.Type == pb.BlockType_BLOCK_TYPE_PORT && entry.Type != "port") {
				continue
			}
		}

		// 過濾過期項目
		if !req.IncludeExpired && !entry.ExpiresAt.IsZero() && time.Now().After(entry.ExpiresAt) {
			continue
		}

		if entry.IsActive {
			activeCount++
		}

		blockType := pb.BlockType_BLOCK_TYPE_IP
		if entry.Type == "port" {
			blockType = pb.BlockType_BLOCK_TYPE_PORT
		}

		entries = append(entries, &pb.BlockEntry{
			EntryId:   entry.EntryID,
			Type:      blockType,
			Value:     entry.Value,
			Reason:    entry.Reason,
			CreatedAt: timestamppb.New(entry.CreatedAt),
			ExpiresAt: timestamppb.New(entry.ExpiresAt),
			IsActive:  entry.IsActive,
			CreatedBy: entry.CreatedBy,
			Metadata:  entry.Metadata,
		})
	}

	// TODO: 實現分頁
	return &pb.BlockListResponse{
		Entries:     entries,
		TotalCount:  int32(len(entries)),
		ActiveCount: activeCount,
	}, nil
}

// ApplyFirewallRule applies a firewall rule
func (s *Service) ApplyFirewallRule(ctx context.Context, req *pb.FirewallRuleRequest) (*pb.FirewallRuleResponse, error) {
	s.logger.Infof("Applying firewall rule: %s", req.Rule.Name)

	ruleID := fmt.Sprintf("fw_rule_%d", time.Now().UnixNano())

	rule := &FirewallRule{
		RuleID:      ruleID,
		Name:        req.Rule.Name,
		Action:      req.Rule.Action.String(),
		SourceIP:    req.Rule.SourceIp,
		DestIP:      req.Rule.DestIp,
		SourcePort:  req.Rule.SourcePort,
		DestPort:    req.Rule.DestPort,
		Protocol:    req.Rule.Protocol,
		Direction:   req.Rule.Direction,
		Priority:    req.Rule.Priority,
		Enabled:     req.Rule.Enabled,
		Description: req.Rule.Description,
		CreatedAt:   time.Now(),
		CreatedBy:   req.Rule.CreatedBy,
	}

	s.mu.Lock()
	s.firewallRules[ruleID] = rule
	s.metrics.ActiveRules++
	s.mu.Unlock()

	// TODO: 實際應用防火牆規則
	if err := s.applyFirewallRule(rule); err != nil {
		s.logger.Errorf("Failed to apply firewall rule: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to apply rule: %v", err)
	}

	// 發布規則應用事件
	event := pubsub.NewSystemEvent("control-service", "rule_applied",
		fmt.Sprintf("Firewall rule %s applied", req.Rule.Name))
	event.Metadata["rule_id"] = ruleID
	event.Metadata["rule_name"] = req.Rule.Name
	if message, err := pubsub.ToJSON(event); err == nil {
		s.mq.Publish(ctx, "pandora.events", "control.rule_applied", message)
	}

	s.logger.Infof("Firewall rule %s applied successfully (ID: %s)", req.Rule.Name, ruleID)

	return &pb.FirewallRuleResponse{
		Success: true,
		Message: "Firewall rule applied successfully",
		RuleId:  ruleID,
	}, nil
}

// RemoveFirewallRule removes a firewall rule
func (s *Service) RemoveFirewallRule(ctx context.Context, req *pb.RemoveRuleRequest) (*pb.RemoveRuleResponse, error) {
	s.logger.Infof("Removing firewall rule: %s", req.RuleId)

	s.mu.Lock()
	defer s.mu.Unlock()

	rule, exists := s.firewallRules[req.RuleId]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "rule not found: %s", req.RuleId)
	}

	// TODO: 實際移除防火牆規則
	if err := s.removeFirewallRule(rule); err != nil {
		s.logger.Errorf("Failed to remove firewall rule: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to remove rule: %v", err)
	}

	delete(s.firewallRules, req.RuleId)
	s.metrics.ActiveRules--

	s.logger.Infof("Firewall rule %s removed successfully", req.RuleId)

	return &pb.RemoveRuleResponse{
		Success: true,
		Message: "Firewall rule removed successfully",
	}, nil
}

// GetFirewallRules gets all firewall rules
func (s *Service) GetFirewallRules(ctx context.Context, _ *emptypb.Empty) (*pb.FirewallRulesResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var rules []*pb.FirewallRule

	for _, rule := range s.firewallRules {
		action := pb.RuleAction_RULE_ACTION_ALLOW
		if rule.Action == "RULE_ACTION_DENY" {
			action = pb.RuleAction_RULE_ACTION_DENY
		}

		rules = append(rules, &pb.FirewallRule{
			RuleId:      rule.RuleID,
			Name:        rule.Name,
			Action:      action,
			SourceIp:    rule.SourceIP,
			DestIp:      rule.DestIP,
			SourcePort:  rule.SourcePort,
			DestPort:    rule.DestPort,
			Protocol:    rule.Protocol,
			Direction:   rule.Direction,
			Priority:    rule.Priority,
			Enabled:     rule.Enabled,
			Description: rule.Description,
			CreatedAt:   timestamppb.New(rule.CreatedAt),
			CreatedBy:   rule.CreatedBy,
		})
	}

	return &pb.FirewallRulesResponse{
		Rules:      rules,
		TotalCount: int32(len(rules)),
	}, nil
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
	metrics := &pb.ControlServiceMetrics{
		ActiveBlocks:   s.metrics.ActiveBlocks,
		TotalBlocks:    s.metrics.TotalBlocks,
		TotalUnblocks:  s.metrics.TotalUnblocks,
		ActiveRules:    s.metrics.ActiveRules,
		PacketsBlocked: s.metrics.PacketsBlocked,
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

// cleanupExpiredBlocks periodically removes expired blocks
func (s *Service) cleanupExpiredBlocks() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for ruleID, entry := range s.blockList {
			if entry.IsActive && !entry.ExpiresAt.IsZero() && now.After(entry.ExpiresAt) {
				s.logger.Infof("Block expired: %s (%s)", entry.Value, ruleID)
				entry.IsActive = false
				s.metrics.ActiveBlocks--

				// TODO: 實際移除規則
				if entry.Type == "ip" {
					s.removeIPBlock(entry.Value)
				} else if entry.Type == "port" {
					// 解析端口和協議
					// s.removePortBlock(port, protocol)
				}

				delete(s.blockList, ruleID)
			}
		}
		s.mu.Unlock()
	}
}

// applyIPBlock applies an IP block rule
func (s *Service) applyIPBlock(ip string, action pb.BlockAction) error {
	// TODO: 實際的 iptables/nftables 命令
	// 示例：
	// iptables -A INPUT -s {ip} -j DROP
	// nft add rule ip filter input ip saddr {ip} drop

	s.logger.Debugf("Applying IP block: %s (action: %s)", ip, action)

	// 模擬實現
	switch s.config.FirewallBackend {
	case "iptables":
		// cmd := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP")
		// return cmd.Run()
		s.logger.Debugf("[MOCK] iptables -A INPUT -s %s -j DROP", ip)
		return nil

	case "nftables":
		// cmd := exec.Command("nft", "add", "rule", "ip", "filter", "input", "ip", "saddr", ip, "drop")
		// return cmd.Run()
		s.logger.Debugf("[MOCK] nft add rule ip filter input ip saddr %s drop", ip)
		return nil

	case "mock":
		s.logger.Debugf("[MOCK] IP %s blocked", ip)
		return nil

	default:
		return fmt.Errorf("unsupported firewall backend: %s", s.config.FirewallBackend)
	}
}

// removeIPBlock removes an IP block rule
func (s *Service) removeIPBlock(ip string) error {
	// TODO: 實際的 iptables/nftables 命令
	// 示例：
	// iptables -D INPUT -s {ip} -j DROP
	// nft delete rule ip filter input handle {handle}

	s.logger.Debugf("Removing IP block: %s", ip)

	switch s.config.FirewallBackend {
	case "iptables":
		s.logger.Debugf("[MOCK] iptables -D INPUT -s %s -j DROP", ip)
		return nil

	case "nftables":
		s.logger.Debugf("[MOCK] nft delete rule for IP %s", ip)
		return nil

	case "mock":
		s.logger.Debugf("[MOCK] IP %s unblocked", ip)
		return nil

	default:
		return fmt.Errorf("unsupported firewall backend: %s", s.config.FirewallBackend)
	}
}

// applyPortBlock applies a port block rule
func (s *Service) applyPortBlock(port int32, protocol string) error {
	s.logger.Debugf("Applying port block: %d/%s", port, protocol)

	// TODO: 實際的防火牆命令
	s.logger.Debugf("[MOCK] Port %d/%s blocked", port, protocol)
	return nil
}

// removePortBlock removes a port block rule
func (s *Service) removePortBlock(port int32, protocol string) error {
	s.logger.Debugf("Removing port block: %d/%s", port, protocol)

	// TODO: 實際的防火牆命令
	s.logger.Debugf("[MOCK] Port %d/%s unblocked", port, protocol)
	return nil
}

// applyFirewallRule applies a firewall rule
func (s *Service) applyFirewallRule(rule *FirewallRule) error {
	s.logger.Debugf("Applying firewall rule: %s", rule.Name)

	// TODO: 實際的防火牆規則應用
	s.logger.Debugf("[MOCK] Firewall rule %s applied", rule.Name)
	return nil
}

// removeFirewallRule removes a firewall rule
func (s *Service) removeFirewallRule(rule *FirewallRule) error {
	s.logger.Debugf("Removing firewall rule: %s", rule.Name)

	// TODO: 實際的防火牆規則移除
	s.logger.Debugf("[MOCK] Firewall rule %s removed", rule.Name)
	return nil
}

