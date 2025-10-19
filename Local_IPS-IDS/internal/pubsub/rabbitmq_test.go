// +build integration

package pubsub

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRabbitMQConnection tests basic connection to RabbitMQ
func TestRabbitMQConnection(t *testing.T) {
	config := DefaultConfig()
	mq, err := NewRabbitMQ(config)
	require.NoError(t, err)
	defer mq.Close()

	// Test health check
	ctx := context.Background()
	err = mq.Health(ctx)
	assert.NoError(t, err)
}

// TestPublishAndSubscribe tests publishing and subscribing to messages
func TestPublishAndSubscribe(t *testing.T) {
	config := DefaultConfig()
	mq, err := NewRabbitMQ(config)
	require.NoError(t, err)
	defer mq.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a channel to receive messages
	received := make(chan *ThreatEvent, 1)

	// Subscribe to threat events
	handler := func(ctx context.Context, msg *Message) error {
		var event ThreatEvent
		if err := FromJSON(msg.Body, &event); err != nil {
			return err
		}
		received <- &event
		return nil
	}

	err = mq.Subscribe(ctx, "threat_events", handler)
	require.NoError(t, err)

	// Give subscriber time to start
	time.Sleep(1 * time.Second)

	// Publish a test event
	event := NewThreatEvent("test_ddos", "192.168.1.100", "Test DDoS", "blocked", 8)
	message, err := ToJSON(event)
	require.NoError(t, err)

	err = mq.Publish(ctx, "pandora.events", "threat.detected", message)
	require.NoError(t, err)

	// Wait for message
	select {
	case receivedEvent := <-received:
		assert.Equal(t, event.ThreatType, receivedEvent.ThreatType)
		assert.Equal(t, event.SourceIP, receivedEvent.SourceIP)
		assert.Equal(t, event.ThreatLevel, receivedEvent.ThreatLevel)
	case <-ctx.Done():
		t.Fatal("Timeout waiting for message")
	}
}

// TestMultipleEvents tests publishing multiple events
func TestMultipleEvents(t *testing.T) {
	config := DefaultConfig()
	mq, err := NewRabbitMQ(config)
	require.NoError(t, err)
	defer mq.Close()

	ctx := context.Background()

	// Publish multiple events
	eventCount := 10
	for i := 0; i < eventCount; i++ {
		event := NewThreatEvent("test", "192.168.1.100", "Test", "blocked", 5)
		message, _ := ToJSON(event)
		err := mq.Publish(ctx, "pandora.events", "threat.detected", message)
		require.NoError(t, err)
	}

	// All events should be published successfully
	assert.NoError(t, err)
}

// TestPublishJSON tests the PublishJSON helper function
func TestPublishJSON(t *testing.T) {
	config := DefaultConfig()
	mq, err := NewRabbitMQ(config)
	require.NoError(t, err)
	defer mq.Close()

	ctx := context.Background()

	event := NewThreatEvent("test", "192.168.1.100", "Test", "blocked", 5)
	err = mq.PublishJSON(ctx, "pandora.events", "threat.detected", event)
	assert.NoError(t, err)
}

// TestReconnect tests automatic reconnection
func TestReconnect(t *testing.T) {
	t.Skip("Requires manual RabbitMQ restart")

	config := DefaultConfig()
	config.ReconnectDelay = 1 * time.Second
	config.MaxReconnectAttempts = 5

	mq, err := NewRabbitMQ(config)
	require.NoError(t, err)
	defer mq.Close()

	ctx := context.Background()

	// Publish initial message
	event := NewThreatEvent("test", "192.168.1.100", "Test", "blocked", 5)
	message, _ := ToJSON(event)
	err = mq.Publish(ctx, "pandora.events", "threat.detected", message)
	require.NoError(t, err)

	// Manually restart RabbitMQ here
	t.Log("Restart RabbitMQ now...")
	time.Sleep(10 * time.Second)

	// Try to publish again (should reconnect automatically)
	err = mq.Publish(ctx, "pandora.events", "threat.detected", message)
	assert.NoError(t, err)
}

// BenchmarkPublish benchmarks message publishing
func BenchmarkPublish(b *testing.B) {
	config := DefaultConfig()
	mq, err := NewRabbitMQ(config)
	require.NoError(b, err)
	defer mq.Close()

	ctx := context.Background()
	event := NewThreatEvent("bench", "192.168.1.100", "Benchmark", "logged", 5)
	message, _ := ToJSON(event)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mq.Publish(ctx, "pandora.events", "threat.detected", message)
	}
}

// BenchmarkPublishJSON benchmarks JSON publishing
func BenchmarkPublishJSON(b *testing.B) {
	config := DefaultConfig()
	mq, err := NewRabbitMQ(config)
	require.NoError(b, err)
	defer mq.Close()

	ctx := context.Background()
	event := NewThreatEvent("bench", "192.168.1.100", "Benchmark", "logged", 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mq.PublishJSON(ctx, "pandora.events", "threat.detected", event)
	}
}

