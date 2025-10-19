package pubsub

import (
	"encoding/json"
	"fmt"
	"time"
)

// EventType defines the type of event
// 事件類型枚舉
type EventType string

const (
	// Threat Events - 威脅事件
	EventTypeThreatDetected   EventType = "threat.detected"
	EventTypeThreatBlocked    EventType = "threat.blocked"
	EventTypeThreatAnalyzed   EventType = "threat.analyzed"
	EventTypeThreatResolved   EventType = "threat.resolved"

	// Network Events - 網路事件
	EventTypeNetworkAttack    EventType = "network.attack"
	EventTypeNetworkScan      EventType = "network.scan"
	EventTypeNetworkAnomaly   EventType = "network.anomaly"
	EventTypeNetworkBlocked   EventType = "network.blocked"

	// System Events - 系統事件
	EventTypeSystemStarted    EventType = "system.started"
	EventTypeSystemStopped    EventType = "system.stopped"
	EventTypeSystemError      EventType = "system.error"
	EventTypeSystemHealthy    EventType = "system.healthy"

	// Device Events - 設備事件
	EventTypeDeviceConnected  EventType = "device.connected"
	EventTypeDeviceDisconnect EventType = "device.disconnected"
	EventTypeDeviceData       EventType = "device.data"
	EventTypeDeviceError      EventType = "device.error"
)

// BaseEvent contains common fields for all events
// 基礎事件結構，所有事件的共同字段
type BaseEvent struct {
	// ID is the unique identifier for the event
	ID string `json:"id"`

	// Type is the event type
	Type EventType `json:"type"`

	// Timestamp is when the event occurred
	Timestamp time.Time `json:"timestamp"`

	// Source is the component that generated the event
	Source string `json:"source"`

	// Severity indicates the severity level (low, medium, high, critical)
	Severity string `json:"severity"`

	// Tags for categorization and filtering
	Tags []string `json:"tags,omitempty"`

	// Metadata contains additional information
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ThreatEvent represents a security threat event
// 威脅事件，用於安全威脅檢測
type ThreatEvent struct {
	BaseEvent

	// ThreatType describes the type of threat (e.g., "ddos", "malware", "intrusion")
	ThreatType string `json:"threat_type"`

	// ThreatLevel indicates the threat level (1-10)
	ThreatLevel int `json:"threat_level"`

	// SourceIP is the IP address of the threat source
	SourceIP string `json:"source_ip"`

	// TargetIP is the IP address of the target
	TargetIP string `json:"target_ip,omitempty"`

	// TargetPort is the port number of the target
	TargetPort int `json:"target_port,omitempty"`

	// Protocol is the network protocol (tcp, udp, icmp, etc.)
	Protocol string `json:"protocol,omitempty"`

	// Description provides details about the threat
	Description string `json:"description"`

	// Action taken (e.g., "blocked", "logged", "alerted")
	Action string `json:"action"`

	// Evidence contains supporting data
	Evidence map[string]interface{} `json:"evidence,omitempty"`
}

// NetworkEvent represents a network-related event
// 網路事件，用於網路流量監控
type NetworkEvent struct {
	BaseEvent

	// EventSubType specifies the network event subtype
	EventSubType string `json:"event_subtype"`

	// SourceIP is the source IP address
	SourceIP string `json:"source_ip"`

	// SourcePort is the source port number
	SourcePort int `json:"source_port,omitempty"`

	// DestIP is the destination IP address
	DestIP string `json:"dest_ip"`

	// DestPort is the destination port number
	DestPort int `json:"dest_port,omitempty"`

	// Protocol is the network protocol
	Protocol string `json:"protocol"`

	// BytesSent is the number of bytes sent
	BytesSent int64 `json:"bytes_sent,omitempty"`

	// BytesReceived is the number of bytes received
	BytesReceived int64 `json:"bytes_received,omitempty"`

	// PacketCount is the number of packets
	PacketCount int64 `json:"packet_count,omitempty"`

	// Duration is the duration of the connection
	Duration time.Duration `json:"duration,omitempty"`

	// Flags contains TCP flags or other protocol-specific flags
	Flags []string `json:"flags,omitempty"`

	// Payload contains sample payload data (if applicable)
	Payload string `json:"payload,omitempty"`
}

// SystemEvent represents a system-level event
// 系統事件，用於系統狀態監控
type SystemEvent struct {
	BaseEvent

	// Component is the system component (e.g., "agent", "engine", "ui")
	Component string `json:"component"`

	// Status is the current status (e.g., "running", "stopped", "error")
	Status string `json:"status"`

	// Message provides additional information
	Message string `json:"message"`

	// ErrorCode is the error code (if applicable)
	ErrorCode string `json:"error_code,omitempty"`

	// ErrorDetails contains error details
	ErrorDetails string `json:"error_details,omitempty"`

	// Metrics contains system metrics
	Metrics map[string]interface{} `json:"metrics,omitempty"`
}

// DeviceEvent represents a device-related event
// 設備事件，用於 IoT 設備監控
type DeviceEvent struct {
	BaseEvent

	// DeviceID is the unique identifier for the device
	DeviceID string `json:"device_id"`

	// DeviceType is the type of device (e.g., "usb-serial", "sensor")
	DeviceType string `json:"device_type"`

	// DeviceName is the human-readable name of the device
	DeviceName string `json:"device_name,omitempty"`

	// Port is the device port (e.g., "/dev/ttyUSB0")
	Port string `json:"port,omitempty"`

	// Status is the device status (e.g., "connected", "disconnected", "error")
	Status string `json:"status"`

	// Data contains device-specific data
	Data map[string]interface{} `json:"data,omitempty"`

	// ErrorMessage contains error information (if applicable)
	ErrorMessage string `json:"error_message,omitempty"`
}

// NewThreatEvent creates a new threat event
// 創建威脅事件的輔助函數
func NewThreatEvent(threatType, sourceIP, description, action string, threatLevel int) *ThreatEvent {
	return &ThreatEvent{
		BaseEvent: BaseEvent{
			ID:        generateEventID(),
			Type:      EventTypeThreatDetected,
			Timestamp: time.Now(),
			Source:    "pandora-agent",
			Severity:  severityFromThreatLevel(threatLevel),
			Tags:      []string{"threat", threatType},
			Metadata:  make(map[string]interface{}),
		},
		ThreatType:  threatType,
		ThreatLevel: threatLevel,
		SourceIP:    sourceIP,
		Description: description,
		Action:      action,
		Evidence:    make(map[string]interface{}),
	}
}

// NewNetworkEvent creates a new network event
// 創建網路事件的輔助函數
func NewNetworkEvent(eventSubType, sourceIP, destIP, protocol string) *NetworkEvent {
	return &NetworkEvent{
		BaseEvent: BaseEvent{
			ID:        generateEventID(),
			Type:      EventTypeNetworkAttack,
			Timestamp: time.Now(),
			Source:    "pandora-agent",
			Severity:  "medium",
			Tags:      []string{"network", eventSubType},
			Metadata:  make(map[string]interface{}),
		},
		EventSubType: eventSubType,
		SourceIP:     sourceIP,
		DestIP:       destIP,
		Protocol:     protocol,
	}
}

// NewSystemEvent creates a new system event
// 創建系統事件的輔助函數
func NewSystemEvent(component, status, message string) *SystemEvent {
	return &SystemEvent{
		BaseEvent: BaseEvent{
			ID:        generateEventID(),
			Type:      EventTypeSystemStarted,
			Timestamp: time.Now(),
			Source:    component,
			Severity:  "info",
			Tags:      []string{"system", component},
			Metadata:  make(map[string]interface{}),
		},
		Component: component,
		Status:    status,
		Message:   message,
		Metrics:   make(map[string]interface{}),
	}
}

// NewDeviceEvent creates a new device event
// 創建設備事件的輔助函數
func NewDeviceEvent(deviceID, deviceType, status string) *DeviceEvent {
	return &DeviceEvent{
		BaseEvent: BaseEvent{
			ID:        generateEventID(),
			Type:      EventTypeDeviceConnected,
			Timestamp: time.Now(),
			Source:    "pandora-agent",
			Severity:  "info",
			Tags:      []string{"device", deviceType},
			Metadata:  make(map[string]interface{}),
		},
		DeviceID:   deviceID,
		DeviceType: deviceType,
		Status:     status,
		Data:       make(map[string]interface{}),
	}
}

// ToJSON converts an event to JSON
// 將事件轉換為 JSON
func ToJSON(event interface{}) ([]byte, error) {
	return json.Marshal(event)
}

// FromJSON parses JSON into an event
// 從 JSON 解析事件
func FromJSON(data []byte, event interface{}) error {
	return json.Unmarshal(data, event)
}

// generateEventID generates a unique event ID
// 生成唯一的事件 ID
func generateEventID() string {
	return fmt.Sprintf("evt_%d", time.Now().UnixNano())
}

// severityFromThreatLevel converts threat level to severity string
// 將威脅等級轉換為嚴重程度字符串
func severityFromThreatLevel(level int) string {
	switch {
	case level >= 9:
		return "critical"
	case level >= 7:
		return "high"
	case level >= 4:
		return "medium"
	default:
		return "low"
	}
}

// GetRoutingKey returns the routing key for an event type
// 獲取事件類型對應的路由鍵
func GetRoutingKey(eventType EventType) string {
	return string(eventType)
}

