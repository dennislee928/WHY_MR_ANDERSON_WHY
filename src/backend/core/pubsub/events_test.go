package pubsub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewThreatEvent(t *testing.T) {
	event := NewThreatEvent("ddos", "192.168.1.100", "DDoS attack detected", "blocked", 8)

	assert.Equal(t, "ddos", event.ThreatType)
	assert.Equal(t, "192.168.1.100", event.SourceIP)
	assert.Equal(t, "DDoS attack detected", event.Description)
	assert.Equal(t, "blocked", event.Action)
	assert.Equal(t, 8, event.ThreatLevel)
	assert.Equal(t, "high", event.Severity)
	assert.Equal(t, "pandora-agent", event.Source)
	assert.NotEmpty(t, event.ID)
	assert.WithinDuration(t, time.Now(), event.Timestamp, 1*time.Second)
}

func TestNewNetworkEvent(t *testing.T) {
	event := NewNetworkEvent("port_scan", "192.168.1.100", "10.0.0.1", "tcp")

	assert.Equal(t, "port_scan", event.EventSubType)
	assert.Equal(t, "192.168.1.100", event.SourceIP)
	assert.Equal(t, "10.0.0.1", event.DestIP)
	assert.Equal(t, "tcp", event.Protocol)
	assert.Equal(t, "medium", event.Severity)
	assert.NotEmpty(t, event.ID)
}

func TestNewSystemEvent(t *testing.T) {
	event := NewSystemEvent("pandora-agent", "running", "Agent started successfully")

	assert.Equal(t, "pandora-agent", event.Component)
	assert.Equal(t, "running", event.Status)
	assert.Equal(t, "Agent started successfully", event.Message)
	assert.Equal(t, "info", event.Severity)
	assert.NotEmpty(t, event.ID)
}

func TestNewDeviceEvent(t *testing.T) {
	event := NewDeviceEvent("usb-001", "usb-serial", "connected")

	assert.Equal(t, "usb-001", event.DeviceID)
	assert.Equal(t, "usb-serial", event.DeviceType)
	assert.Equal(t, "connected", event.Status)
	assert.Equal(t, "info", event.Severity)
	assert.NotEmpty(t, event.ID)
}

func TestToJSON(t *testing.T) {
	event := NewThreatEvent("ddos", "192.168.1.100", "Test", "blocked", 8)

	data, err := ToJSON(event)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
	assert.Contains(t, string(data), "ddos")
	assert.Contains(t, string(data), "192.168.1.100")
}

func TestFromJSON(t *testing.T) {
	original := NewThreatEvent("ddos", "192.168.1.100", "Test", "blocked", 8)
	original.TargetIP = "10.0.0.1"
	original.TargetPort = 80

	data, err := ToJSON(original)
	require.NoError(t, err)

	var decoded ThreatEvent
	err = FromJSON(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original.ThreatType, decoded.ThreatType)
	assert.Equal(t, original.SourceIP, decoded.SourceIP)
	assert.Equal(t, original.TargetIP, decoded.TargetIP)
	assert.Equal(t, original.TargetPort, decoded.TargetPort)
	assert.Equal(t, original.ThreatLevel, decoded.ThreatLevel)
}

func TestSeverityFromThreatLevel(t *testing.T) {
	tests := []struct {
		level    int
		expected string
	}{
		{10, "critical"},
		{9, "critical"},
		{8, "high"},
		{7, "high"},
		{6, "medium"},
		{4, "medium"},
		{3, "low"},
		{1, "low"},
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.level)), func(t *testing.T) {
			result := severityFromThreatLevel(tt.level)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRoutingKey(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  string
	}{
		{EventTypeThreatDetected, "threat.detected"},
		{EventTypeNetworkAttack, "network.attack"},
		{EventTypeSystemStarted, "system.started"},
		{EventTypeDeviceConnected, "device.connected"},
	}

	for _, tt := range tests {
		t.Run(string(tt.eventType), func(t *testing.T) {
			result := GetRoutingKey(tt.eventType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEventIDUniqueness(t *testing.T) {
	ids := make(map[string]bool)

	// 生成 1000 個事件 ID，確保唯一性
	for i := 0; i < 1000; i++ {
		event := NewThreatEvent("test", "1.1.1.1", "test", "test", 1)
		assert.False(t, ids[event.ID], "Duplicate event ID: %s", event.ID)
		ids[event.ID] = true
	}
}

