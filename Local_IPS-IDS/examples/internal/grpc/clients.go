package grpc

import (
	"context"
	"fmt"
	"time"

	devicepb "pandora_box_console_ids_ips/api/proto/device"
	networkpb "pandora_box_console_ids_ips/api/proto/network"
	controlpb "pandora_box_console_ids_ips/api/proto/control"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// ClientConfig contains configuration for gRPC clients
type ClientConfig struct {
	Address        string
	Timeout        time.Duration
	RetryAttempts  int
	KeepAlive      time.Duration
	KeepAliveTime  time.Duration
}

// DefaultClientConfig returns default client configuration
func DefaultClientConfig(address string) *ClientConfig {
	return &ClientConfig{
		Address:        address,
		Timeout:        10 * time.Second,
		RetryAttempts:  3,
		KeepAlive:      30 * time.Second,
		KeepAliveTime:  10 * time.Second,
	}
}

// DeviceClient wraps the Device Service gRPC client
type DeviceClient struct {
	conn   *grpc.ClientConn
	client devicepb.DeviceServiceClient
	config *ClientConfig
	logger *logrus.Logger
}

// NewDeviceClient creates a new Device Service client
func NewDeviceClient(config *ClientConfig, logger *logrus.Logger) (*DeviceClient, error) {
	if logger == nil {
		logger = logrus.New()
	}

	// 創建 gRPC 連接
	conn, err := grpc.Dial(
		config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                config.KeepAlive,
			Timeout:             config.KeepAliveTime,
			PermitWithoutStream: true,
		}),
		grpc.WithDefaultCallOptions(
			grpc.WaitForReady(true),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial device service: %w", err)
	}

	client := devicepb.NewDeviceServiceClient(conn)

	logger.Infof("Connected to Device Service at %s", config.Address)

	return &DeviceClient{
		conn:   conn,
		client: client,
		config: config,
		logger: logger,
	}, nil
}

// Connect connects to a device
func (c *DeviceClient) Connect(ctx context.Context, deviceID, port string, baudRate int32) (*devicepb.ConnectResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &devicepb.ConnectRequest{
		DeviceId: deviceID,
		Port:     port,
		BaudRate: baudRate,
	}

	return c.client.Connect(ctx, req)
}

// GetStatus gets device status
func (c *DeviceClient) GetStatus(ctx context.Context, deviceID string) (*devicepb.StatusResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &devicepb.StatusRequest{
		DeviceId: deviceID,
	}

	return c.client.GetStatus(ctx, req)
}

// ListDevices lists all devices
func (c *DeviceClient) ListDevices(ctx context.Context, includeDisconnected bool) (*devicepb.ListDevicesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &devicepb.ListDevicesRequest{
		IncludeDisconnected: includeDisconnected,
	}

	return c.client.ListDevices(ctx, req)
}

// Close closes the client connection
func (c *DeviceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// NetworkClient wraps the Network Service gRPC client
type NetworkClient struct {
	conn   *grpc.ClientConn
	client networkpb.NetworkServiceClient
	config *ClientConfig
	logger *logrus.Logger
}

// NewNetworkClient creates a new Network Service client
func NewNetworkClient(config *ClientConfig, logger *logrus.Logger) (*NetworkClient, error) {
	if logger == nil {
		logger = logrus.New()
	}

	conn, err := grpc.Dial(
		config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                config.KeepAlive,
			Timeout:             config.KeepAliveTime,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial network service: %w", err)
	}

	client := networkpb.NewNetworkServiceClient(conn)

	logger.Infof("Connected to Network Service at %s", config.Address)

	return &NetworkClient{
		conn:   conn,
		client: client,
		config: config,
		logger: logger,
	}, nil
}

// StartMonitoring starts network monitoring
func (c *NetworkClient) StartMonitoring(ctx context.Context, interfaceName string) (*networkpb.MonitorResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &networkpb.MonitorRequest{
		InterfaceName:   interfaceName,
		PromiscuousMode: true,
		SnapshotLength:  1600,
	}

	return c.client.StartMonitoring(ctx, req)
}

// GetStatistics gets network statistics
func (c *NetworkClient) GetStatistics(ctx context.Context, sessionID string) (*networkpb.StatsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &networkpb.StatsRequest{
		SessionId: sessionID,
	}

	return c.client.GetStatistics(ctx, req)
}

// Close closes the client connection
func (c *NetworkClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ControlClient wraps the Control Service gRPC client
type ControlClient struct {
	conn   *grpc.ClientConn
	client controlpb.ControlServiceClient
	config *ClientConfig
	logger *logrus.Logger
}

// NewControlClient creates a new Control Service client
func NewControlClient(config *ClientConfig, logger *logrus.Logger) (*ControlClient, error) {
	if logger == nil {
		logger = logrus.New()
	}

	conn, err := grpc.Dial(
		config.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                config.KeepAlive,
			Timeout:             config.KeepAliveTime,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial control service: %w", err)
	}

	client := controlpb.NewControlServiceClient(conn)

	logger.Infof("Connected to Control Service at %s", config.Address)

	return &ControlClient{
		conn:   conn,
		client: client,
		config: config,
		logger: logger,
	}, nil
}

// BlockIP blocks an IP address
func (c *ControlClient) BlockIP(ctx context.Context, ipAddress, reason string, duration int32) (*controlpb.BlockIPResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &controlpb.BlockIPRequest{
		IpAddress:       ipAddress,
		Reason:          reason,
		DurationSeconds: duration,
		Action:          controlpb.BlockAction_BLOCK_ACTION_DROP,
	}

	return c.client.BlockIP(ctx, req)
}

// UnblockIP unblocks an IP address
func (c *ControlClient) UnblockIP(ctx context.Context, ipAddress, reason string) (*controlpb.UnblockIPResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &controlpb.UnblockIPRequest{
		IpAddress: ipAddress,
		Reason:    reason,
	}

	return c.client.UnblockIP(ctx, req)
}

// BlockPort blocks a port
func (c *ControlClient) BlockPort(ctx context.Context, port int32, protocol, reason string, duration int32) (*controlpb.BlockPortResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &controlpb.BlockPortRequest{
		Port:            port,
		Protocol:        protocol,
		Reason:          reason,
		DurationSeconds: duration,
	}

	return c.client.BlockPort(ctx, req)
}

// GetBlockList gets the current block list
func (c *ControlClient) GetBlockList(ctx context.Context) (*controlpb.BlockListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &controlpb.BlockListRequest{
		Type:           controlpb.BlockType_BLOCK_TYPE_UNKNOWN,
		IncludeExpired: false,
	}

	return c.client.GetBlockList(ctx, req)
}

// Close closes the client connection
func (c *ControlClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ServiceClients contains all gRPC service clients
type ServiceClients struct {
	Device  *DeviceClient
	Network *NetworkClient
	Control *ControlClient
}

// NewServiceClients creates all service clients
func NewServiceClients(deviceAddr, networkAddr, controlAddr string, logger *logrus.Logger) (*ServiceClients, error) {
	deviceClient, err := NewDeviceClient(DefaultClientConfig(deviceAddr), logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create device client: %w", err)
	}

	networkClient, err := NewNetworkClient(DefaultClientConfig(networkAddr), logger)
	if err != nil {
		deviceClient.Close()
		return nil, fmt.Errorf("failed to create network client: %w", err)
	}

	controlClient, err := NewControlClient(DefaultClientConfig(controlAddr), logger)
	if err != nil {
		deviceClient.Close()
		networkClient.Close()
		return nil, fmt.Errorf("failed to create control client: %w", err)
	}

	return &ServiceClients{
		Device:  deviceClient,
		Network: networkClient,
		Control: controlClient,
	}, nil
}

// Close closes all client connections
func (sc *ServiceClients) Close() error {
	var errs []error

	if err := sc.Device.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := sc.Network.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := sc.Control.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing clients: %v", errs)
	}

	return nil
}

