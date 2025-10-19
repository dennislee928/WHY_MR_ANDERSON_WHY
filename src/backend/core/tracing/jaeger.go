package tracing

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// TracingConfig contains tracing configuration
type TracingConfig struct {
	ServiceName     string
	AgentHost       string
	AgentPort       string
	SamplerType     string
	SamplerParam    float64
	LogSpans        bool
	ReporterLogSpans bool
}

// Tracer wraps the Jaeger tracer
type Tracer struct {
	tracer opentracing.Tracer
	closer io.Closer
	logger *logrus.Logger
}

// NewTracer creates a new Jaeger tracer
func NewTracer(cfg *TracingConfig, logger *logrus.Logger) (*Tracer, error) {
	if logger == nil {
		logger = logrus.New()
	}

	// 配置 Jaeger
	jcfg := config.Configuration{
		ServiceName: cfg.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  cfg.SamplerType,
			Param: cfg.SamplerParam,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            cfg.ReporterLogSpans,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  fmt.Sprintf("%s:%s", cfg.AgentHost, cfg.AgentPort),
		},
	}

	// 創建 tracer
	tracer, closer, err := jcfg.NewTracer(
		config.Logger(jaeger.StdLogger),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tracer: %w", err)
	}

	// 設置為全局 tracer
	opentracing.SetGlobalTracer(tracer)

	logger.Infof("Jaeger tracer initialized for service: %s", cfg.ServiceName)

	return &Tracer{
		tracer: tracer,
		closer: closer,
		logger: logger,
	}, nil
}

// StartSpan starts a new span
func (t *Tracer) StartSpan(operationName string) opentracing.Span {
	return t.tracer.StartSpan(operationName)
}

// StartSpanFromContext starts a new span from context
func (t *Tracer) StartSpanFromContext(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, operationName)
	return span, ctx
}

// InjectSpan injects span context into carrier
func (t *Tracer) InjectSpan(span opentracing.Span, carrier interface{}) error {
	return t.tracer.Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		carrier,
	)
}

// ExtractSpan extracts span context from carrier
func (t *Tracer) ExtractSpan(carrier interface{}) (opentracing.SpanContext, error) {
	return t.tracer.Extract(
		opentracing.HTTPHeaders,
		carrier,
	)
}

// Close closes the tracer
func (t *Tracer) Close() error {
	if t.closer != nil {
		return t.closer.Close()
	}
	return nil
}

// SpanHelper provides helper methods for span management
type SpanHelper struct {
	tracer *Tracer
	logger *logrus.Logger
}

// NewSpanHelper creates a new span helper
func NewSpanHelper(tracer *Tracer, logger *logrus.Logger) *SpanHelper {
	return &SpanHelper{
		tracer: tracer,
		logger: logger,
	}
}

// TraceGRPCCall traces a gRPC call
func (sh *SpanHelper) TraceGRPCCall(ctx context.Context, method string) (opentracing.Span, context.Context) {
	span, ctx := sh.tracer.StartSpanFromContext(ctx, method)
	ext.SpanKindRPCClient.Set(span)
	ext.Component.Set(span, "grpc")
	return span, ctx
}

// TraceHTTPRequest traces an HTTP request
func (sh *SpanHelper) TraceHTTPRequest(ctx context.Context, method, url string) (opentracing.Span, context.Context) {
	span, ctx := sh.tracer.StartSpanFromContext(ctx, fmt.Sprintf("HTTP %s", method))
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPMethod.Set(span, method)
	ext.HTTPUrl.Set(span, url)
	ext.Component.Set(span, "http")
	return span, ctx
}

// TraceDBQuery traces a database query
func (sh *SpanHelper) TraceDBQuery(ctx context.Context, dbType, query string) (opentracing.Span, context.Context) {
	span, ctx := sh.tracer.StartSpanFromContext(ctx, "DB Query")
	ext.DBType.Set(span, dbType)
	ext.DBStatement.Set(span, query)
	ext.Component.Set(span, "database")
	return span, ctx
}

// TraceMessageQueue traces a message queue operation
func (sh *SpanHelper) TraceMessageQueue(ctx context.Context, operation, queue string) (opentracing.Span, context.Context) {
	span, ctx := sh.tracer.StartSpanFromContext(ctx, fmt.Sprintf("MQ %s", operation))
	span.SetTag("mq.operation", operation)
	span.SetTag("mq.queue", queue)
	ext.Component.Set(span, "message_queue")
	return span, ctx
}

// TraceCache traces a cache operation
func (sh *SpanHelper) TraceCache(ctx context.Context, operation, key string) (opentracing.Span, context.Context) {
	span, ctx := sh.tracer.StartSpanFromContext(ctx, fmt.Sprintf("Cache %s", operation))
	span.SetTag("cache.operation", operation)
	span.SetTag("cache.key", key)
	ext.Component.Set(span, "cache")
	return span, ctx
}

// FinishSpanWithError finishes a span and records an error if present
func (sh *SpanHelper) FinishSpanWithError(span opentracing.Span, err error) {
	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV(
			"event", "error",
			"message", err.Error(),
		)
	}
	span.Finish()
}

// AddSpanTags adds multiple tags to a span
func (sh *SpanHelper) AddSpanTags(span opentracing.Span, tags map[string]interface{}) {
	for key, value := range tags {
		span.SetTag(key, value)
	}
}

// LogSpanEvent logs an event in a span
func (sh *SpanHelper) LogSpanEvent(span opentracing.Span, event string, fields map[string]interface{}) {
	kvs := make([]interface{}, 0, len(fields)*2+2)
	kvs = append(kvs, "event", event)
	
	for key, value := range fields {
		kvs = append(kvs, key, value)
	}
	
	span.LogKV(kvs...)
}

// Middleware provides tracing middleware for various frameworks
type Middleware struct {
	tracer *Tracer
	helper *SpanHelper
	logger *logrus.Logger
}

// NewMiddleware creates a new tracing middleware
func NewMiddleware(tracer *Tracer, logger *logrus.Logger) *Middleware {
	return &Middleware{
		tracer: tracer,
		helper: NewSpanHelper(tracer, logger),
		logger: logger,
	}
}

// WrapGRPCUnaryServer wraps a gRPC unary server with tracing
func (m *Middleware) WrapGRPCUnaryServer(ctx context.Context, req interface{}, info interface{}, handler func(context.Context, interface{}) (interface{}, error)) (interface{}, error) {
	// 從 metadata 中提取 span context
	span := m.tracer.StartSpan(fmt.Sprintf("gRPC Server: %v", info))
	defer span.Finish()

	ctx = opentracing.ContextWithSpan(ctx, span)

	ext.SpanKindRPCServer.Set(span)
	ext.Component.Set(span, "grpc")

	// 執行 handler
	resp, err := handler(ctx, req)

	if err != nil {
		ext.Error.Set(span, true)
		span.LogKV("event", "error", "message", err.Error())
	}

	return resp, err
}

// TracingMetrics contains tracing metrics
type TracingMetrics struct {
	TotalSpans      int64
	ActiveSpans     int64
	AvgSpanDuration time.Duration
	ErrorSpans      int64
}

// GetMetrics returns tracing metrics
func (t *Tracer) GetMetrics() *TracingMetrics {
	// 實際實現需要從 Jaeger 收集指標
	return &TracingMetrics{
		TotalSpans:      0,
		ActiveSpans:     0,
		AvgSpanDuration: 0,
		ErrorSpans:      0,
	}
}

// DefaultConfig returns default tracing configuration
func DefaultConfig(serviceName string) *TracingConfig {
	return &TracingConfig{
		ServiceName:      serviceName,
		AgentHost:        "localhost",
		AgentPort:        "6831",
		SamplerType:      "const",
		SamplerParam:     1.0, // 100% sampling for development
		LogSpans:         false,
		ReporterLogSpans: false,
	}
}

// ProductionConfig returns production tracing configuration
func ProductionConfig(serviceName string) *TracingConfig {
	return &TracingConfig{
		ServiceName:      serviceName,
		AgentHost:        "jaeger-agent",
		AgentPort:        "6831",
		SamplerType:      "probabilistic",
		SamplerParam:     0.1, // 10% sampling for production
		LogSpans:         false,
		ReporterLogSpans: false,
	}
}

