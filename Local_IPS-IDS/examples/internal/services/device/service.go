package device

import (
	"context"
	"fmt"
	"sync"
	"time"

	"pandora_box_console_ids_ips/internal/pubsub"
	pb "pandora_box_console_ids_ips/api/proto/device"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Config contains configuration for the device service
type Config struct {
	DefaultPort     string
	DefaultBaudRate int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

// Service implements the DeviceService gRPC service
type Service struct {
	pb.UnimplementedDeviceServiceServer

	config    *Config
	mq        pubsub.MessageQueue
	logger    *logrus.Logger
	devices   map[string]*DeviceConnection
	mu        sync.RWMutex
	startTime time.Time
}

// DeviceConnection represents a connection to a device
type DeviceConnection struct {
	DeviceID    string
	Port        string
	BaudRate    int
	Connected   bool
	ConnectedAt time.Time
	Metrics     *DeviceMetrics
}

// DeviceMetrics tracks device metrics
type DeviceMetrics struct {
	BytesRead       int64
	BytesWritten    int64
	ReadOperations  int64
	WriteOperations int64
	Errors          int64
}

// NewService creates a new device service
func NewService(config *Config, mq pubsub.MessageQueue, logger *logrus.Logger) *Service {
	if logger == nil {
		logger = logrus.New()
	}

	return &Service{
		config:    config,
		mq:        mq,
		logger:    logger,
		devices:   make(map[string]*DeviceConnection),
		startTime: time.Now(),
	}
}

// Connect establishes a connection to a device
func (s *Service) Connect(ctx context.Context, req *pb.ConnectRequest) (*pb.ConnectResponse, error) {
	s.logger.Infof("Connecting to device: %s on port %s", req.DeviceId, req.Port)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if already connected
	if conn, exists := s.devices[req.DeviceId]; exists && conn.Connected {
		return &pb.ConnectResponse{
			Success: false,
			Message: "Device already connected",
		}, nil
	}

	// Create device connection
	conn := &DeviceConnection{
		DeviceID:    req.DeviceId,
		Port:        req.Port,
		BaudRate:    int(req.BaudRate),
		Connected:   true,
		ConnectedAt: time.Now(),
		Metrics:     &DeviceMetrics{},
	}

	s.devices[req.DeviceId] = conn

	// Publish device connected event
	event := pubsub.NewDeviceEvent(req.DeviceId, "usb-serial", "connected")
	event.Port = req.Port
	event.Data["baud_rate"] = req.BaudRate
	if message, err := pubsub.ToJSON(event); err == nil {
		s.mq.Publish(ctx, "pandora.events", "device.connected", message)
	}

	s.logger.Infof("Device %s connected successfully", req.DeviceId)

	return &pb.ConnectResponse{
		Success: true,
		Message: "Device connected successfully",
		Device: &pb.DeviceInfo{
			DeviceId:    req.DeviceId,
			DeviceType:  "usb-serial",
			DeviceName:  "CH340 Serial",
			Port:        req.Port,
			ConnectedAt: timestamppb.New(conn.ConnectedAt),
		},
	}, nil
}

// Disconnect closes the connection to a device
func (s *Service) Disconnect(ctx context.Context, req *pb.DisconnectRequest) (*pb.DisconnectResponse, error) {
	s.logger.Infof("Disconnecting device: %s", req.DeviceId)

	s.mu.Lock()
	defer s.mu.Unlock()

	conn, exists := s.devices[req.DeviceId]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "device not found: %s", req.DeviceId)
	}

	conn.Connected = false

	// Publish device disconnected event
	event := pubsub.NewDeviceEvent(req.DeviceId, "usb-serial", "disconnected")
	event.Port = conn.Port
	if message, err := pubsub.ToJSON(event); err == nil {
		s.mq.Publish(ctx, "pandora.events", "device.disconnected", message)
	}

	delete(s.devices, req.DeviceId)

	s.logger.Infof("Device %s disconnected", req.DeviceId)

	return &pb.DisconnectResponse{
		Success: true,
		Message: "Device disconnected successfully",
	}, nil
}

// ReadData reads data from the device (streaming)
func (s *Service) ReadData(req *pb.ReadDataRequest, stream pb.DeviceService_ReadDataServer) error {
	s.logger.Debugf("Reading data from device: %s", req.DeviceId)

	s.mu.RLock()
	conn, exists := s.devices[req.DeviceId]
	s.mu.RUnlock()

	if !exists || !conn.Connected {
		return status.Errorf(codes.NotFound, "device not connected: %s", req.DeviceId)
	}

	// TODO: 實際的設備讀取邏輯
	// 這裡是模擬實現
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			// 模擬讀取數據
			data := []byte(fmt.Sprintf("Device %s data at %s", req.DeviceId, time.Now().Format(time.RFC3339)))

			resp := &pb.DataResponse{
				DeviceId:  req.DeviceId,
				Data:      data,
				Timestamp: timestamppb.Now(),
				BytesRead: int32(len(data)),
			}

			if err := stream.Send(resp); err != nil {
				s.logger.Errorf("Failed to send data: %v", err)
				return err
			}

			conn.Metrics.BytesRead += int64(len(data))
			conn.Metrics.ReadOperations++
		}
	}
}

// WriteData writes data to the device
func (s *Service) WriteData(ctx context.Context, req *pb.WriteDataRequest) (*pb.WriteDataResponse, error) {
	s.logger.Debugf("Writing %d bytes to device: %s", len(req.Data), req.DeviceId)

	s.mu.RLock()
	conn, exists := s.devices[req.DeviceId]
	s.mu.RUnlock()

	if !exists || !conn.Connected {
		return nil, status.Errorf(codes.NotFound, "device not connected: %s", req.DeviceId)
	}

	// TODO: 實際的設備寫入邏輯

	conn.Metrics.BytesWritten += int64(len(req.Data))
	conn.Metrics.WriteOperations++

	return &pb.WriteDataResponse{
		Success:      true,
		BytesWritten: int32(len(req.Data)),
		Message:      "Data written successfully",
	}, nil
}

// GetStatus gets the current status of a device
func (s *Service) GetStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conn, exists := s.devices[req.DeviceId]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "device not found: %s", req.DeviceId)
	}

	deviceStatus := pb.DeviceStatus_DEVICE_STATUS_DISCONNECTED
	if conn.Connected {
		deviceStatus = pb.DeviceStatus_DEVICE_STATUS_CONNECTED
	}

	return &pb.StatusResponse{
		DeviceId: req.DeviceId,
		Status:   deviceStatus,
		Info: &pb.DeviceInfo{
			DeviceId:    conn.DeviceID,
			DeviceType:  "usb-serial",
			DeviceName:  "CH340 Serial",
			Port:        conn.Port,
			ConnectedAt: timestamppb.New(conn.ConnectedAt),
		},
		Metrics: &pb.DeviceMetrics{
			BytesRead:       conn.Metrics.BytesRead,
			BytesWritten:    conn.Metrics.BytesWritten,
			ReadOperations:  conn.Metrics.ReadOperations,
			WriteOperations: conn.Metrics.WriteOperations,
			Errors:          conn.Metrics.Errors,
			UptimeSeconds:   time.Since(conn.ConnectedAt).Seconds(),
		},
	}, nil
}

// ListDevices lists all available devices
func (s *Service) ListDevices(ctx context.Context, req *pb.ListDevicesRequest) (*pb.ListDevicesResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var devices []*pb.DeviceInfo

	for _, conn := range s.devices {
		if !req.IncludeDisconnected && !conn.Connected {
			continue
		}

		devices = append(devices, &pb.DeviceInfo{
			DeviceId:    conn.DeviceID,
			DeviceType:  "usb-serial",
			DeviceName:  "CH340 Serial",
			Port:        conn.Port,
			ConnectedAt: timestamppb.New(conn.ConnectedAt),
		})
	}

	return &pb.ListDevicesResponse{
		Devices:    devices,
		TotalCount: int32(len(devices)),
	}, nil
}

// GetHealth checks the health of the service
func (s *Service) GetHealth(ctx context.Context, _ *emptypb.Empty) (*pb.HealthResponse, error) {
	healthy := s.Health(ctx) == nil

	status := "healthy"
	if !healthy {
		status = "unhealthy"
	}

	dependencies := make(map[string]string)
	if s.mq != nil {
		if err := s.mq.Health(ctx); err != nil {
			dependencies["rabbitmq"] = "unhealthy"
		} else {
			dependencies["rabbitmq"] = "healthy"
		}
	}

	return &pb.HealthResponse{
		Healthy:        healthy,
		Status:         status,
		Version:        "1.0.0",
		UptimeSeconds:  int64(time.Since(s.startTime).Seconds()),
		Dependencies:   dependencies,
	}, nil
}

// Health checks if the service is healthy
func (s *Service) Health(ctx context.Context) error {
	// Check RabbitMQ connection
	if s.mq != nil {
		if err := s.mq.Health(ctx); err != nil {
			return fmt.Errorf("rabbitmq unhealthy: %w", err)
		}
	}

	return nil
}

