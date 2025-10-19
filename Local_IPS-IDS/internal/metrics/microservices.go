package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MicroserviceMetrics contains Prometheus metrics for microservices
// 微服務 Prometheus 指標
type MicroserviceMetrics struct {
	// gRPC 指標
	GRPCRequestsTotal   *prometheus.CounterVec
	GRPCRequestDuration *prometheus.HistogramVec
	GRPCRequestsInFlight *prometheus.GaugeVec

	// 業務指標 - Device Service
	DeviceConnectionsTotal *prometheus.CounterVec
	DeviceReadOperations   *prometheus.CounterVec
	DeviceWriteOperations  *prometheus.CounterVec
	DeviceBytesRead        *prometheus.CounterVec
	DeviceBytesWritten     *prometheus.CounterVec
	DeviceErrors           *prometheus.CounterVec
	DeviceActiveConnections prometheus.Gauge

	// 業務指標 - Network Service
	NetworkPacketsTotal    *prometheus.CounterVec
	NetworkBytesTotal      *prometheus.CounterVec
	NetworkFlowsTotal      prometheus.Counter
	NetworkAnomaliesTotal  *prometheus.CounterVec
	NetworkAttacksTotal    *prometheus.CounterVec
	NetworkActiveSessions  prometheus.Gauge

	// 業務指標 - Control Service
	ControlBlocksTotal     *prometheus.CounterVec
	ControlUnblocksTotal   *prometheus.CounterVec
	ControlRulesTotal      prometheus.Gauge
	ControlActiveBlocks    prometheus.Gauge
	ControlPacketsBlocked  prometheus.Counter

	// RabbitMQ 指標
	RabbitMQPublishTotal   *prometheus.CounterVec
	RabbitMQPublishErrors  *prometheus.CounterVec
	RabbitMQSubscribeTotal *prometheus.CounterVec
	RabbitMQConsumeErrors  *prometheus.CounterVec
}

// NewMicroserviceMetrics creates new microservice metrics
func NewMicroserviceMetrics(serviceName string) *MicroserviceMetrics {
	return &MicroserviceMetrics{
		// gRPC 指標
		GRPCRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_requests_total",
				Help: "Total number of gRPC requests",
			},
			[]string{"service", "method", "status"},
		),
		GRPCRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_request_duration_seconds",
				Help:    "Duration of gRPC requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "method"},
		),
		GRPCRequestsInFlight: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "grpc_requests_in_flight",
				Help: "Number of gRPC requests currently being processed",
			},
			[]string{"service", "method"},
		),

		// Device Service 指標
		DeviceConnectionsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "device_connections_total",
				Help: "Total number of device connections",
			},
			[]string{"device_id", "status"},
		),
		DeviceReadOperations: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "device_read_operations_total",
				Help: "Total number of device read operations",
			},
			[]string{"device_id"},
		),
		DeviceWriteOperations: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "device_write_operations_total",
				Help: "Total number of device write operations",
			},
			[]string{"device_id"},
		),
		DeviceBytesRead: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "device_bytes_read_total",
				Help: "Total bytes read from devices",
			},
			[]string{"device_id"},
		),
		DeviceBytesWritten: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "device_bytes_written_total",
				Help: "Total bytes written to devices",
			},
			[]string{"device_id"},
		),
		DeviceErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "device_errors_total",
				Help: "Total number of device errors",
			},
			[]string{"device_id", "error_type"},
		),
		DeviceActiveConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "device_active_connections",
				Help: "Number of active device connections",
			},
		),

		// Network Service 指標
		NetworkPacketsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "network_packets_total",
				Help: "Total number of network packets captured",
			},
			[]string{"protocol", "direction"},
		),
		NetworkBytesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "network_bytes_total",
				Help: "Total bytes of network traffic",
			},
			[]string{"protocol", "direction"},
		),
		NetworkFlowsTotal: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "network_flows_total",
				Help: "Total number of network flows",
			},
		),
		NetworkAnomaliesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "network_anomalies_total",
				Help: "Total number of network anomalies detected",
			},
			[]string{"anomaly_type"},
		),
		NetworkAttacksTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "network_attacks_total",
				Help: "Total number of network attacks detected",
			},
			[]string{"attack_type"},
		),
		NetworkActiveSessions: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "network_active_sessions",
				Help: "Number of active monitoring sessions",
			},
		),

		// Control Service 指標
		ControlBlocksTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "control_blocks_total",
				Help: "Total number of blocks applied",
			},
			[]string{"type", "action"},
		),
		ControlUnblocksTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "control_unblocks_total",
				Help: "Total number of unblocks",
			},
			[]string{"type"},
		),
		ControlRulesTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "control_rules_total",
				Help: "Total number of firewall rules",
			},
		),
		ControlActiveBlocks: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "control_active_blocks",
				Help: "Number of active blocks",
			},
		),
		ControlPacketsBlocked: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "control_packets_blocked_total",
				Help: "Total number of packets blocked",
			},
		),

		// RabbitMQ 指標
		RabbitMQPublishTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rabbitmq_publish_total",
				Help: "Total number of messages published to RabbitMQ",
			},
			[]string{"exchange", "routing_key", "status"},
		),
		RabbitMQPublishErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rabbitmq_publish_errors_total",
				Help: "Total number of RabbitMQ publish errors",
			},
			[]string{"exchange", "routing_key", "error_type"},
		),
		RabbitMQSubscribeTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rabbitmq_subscribe_total",
				Help: "Total number of messages consumed from RabbitMQ",
			},
			[]string{"queue", "status"},
		),
		RabbitMQConsumeErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rabbitmq_consume_errors_total",
				Help: "Total number of RabbitMQ consume errors",
			},
			[]string{"queue", "error_type"},
		),
	}
}

// RecordGRPCRequest records a gRPC request
func (m *MicroserviceMetrics) RecordGRPCRequest(service, method, status string, duration float64) {
	m.GRPCRequestsTotal.WithLabelValues(service, method, status).Inc()
	m.GRPCRequestDuration.WithLabelValues(service, method).Observe(duration)
}

// RecordDeviceConnection records a device connection
func (m *MicroserviceMetrics) RecordDeviceConnection(deviceID, status string) {
	m.DeviceConnectionsTotal.WithLabelValues(deviceID, status).Inc()
	if status == "connected" {
		m.DeviceActiveConnections.Inc()
	} else if status == "disconnected" {
		m.DeviceActiveConnections.Dec()
	}
}

// RecordDeviceRead records a device read operation
func (m *MicroserviceMetrics) RecordDeviceRead(deviceID string, bytes int64) {
	m.DeviceReadOperations.WithLabelValues(deviceID).Inc()
	m.DeviceBytesRead.WithLabelValues(deviceID).Add(float64(bytes))
}

// RecordNetworkPacket records a network packet
func (m *MicroserviceMetrics) RecordNetworkPacket(protocol, direction string, bytes int64) {
	m.NetworkPacketsTotal.WithLabelValues(protocol, direction).Inc()
	m.NetworkBytesTotal.WithLabelValues(protocol, direction).Add(float64(bytes))
}

// RecordNetworkAnomaly records a network anomaly
func (m *MicroserviceMetrics) RecordNetworkAnomaly(anomalyType string) {
	m.NetworkAnomaliesTotal.WithLabelValues(anomalyType).Inc()
}

// RecordControlBlock records a control block operation
func (m *MicroserviceMetrics) RecordControlBlock(blockType, action string) {
	m.ControlBlocksTotal.WithLabelValues(blockType, action).Inc()
	m.ControlActiveBlocks.Inc()
}

// RecordControlUnblock records a control unblock operation
func (m *MicroserviceMetrics) RecordControlUnblock(blockType string) {
	m.ControlUnblocksTotal.WithLabelValues(blockType).Inc()
	m.ControlActiveBlocks.Dec()
}

// RecordRabbitMQPublish records a RabbitMQ publish operation
func (m *MicroserviceMetrics) RecordRabbitMQPublish(exchange, routingKey, status string) {
	m.RabbitMQPublishTotal.WithLabelValues(exchange, routingKey, status).Inc()
}

