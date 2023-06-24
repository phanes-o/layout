package trace

import (
	"context"

	"go-micro.dev/v4/metadata"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	// instrumentationName is the name of this instrumentation package.
	instrumentationName = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	// GRPCStatusCodeKey is convention for numeric status code of a gRPC request.
	GRPCStatusCodeKey = attribute.Key("rpc.grpc.status_code")
)

// config is a group of options for this instrumentation.
type traceConfig struct {
	Propagators    propagation.TextMapPropagator
	TracerProvider trace.TracerProvider
}

// Option applies an option value for a config.
type Option interface {
	apply(*traceConfig)
}

// newConfig returns a config configured with all the passed Options.
func newConfig(opts []Option) *traceConfig {
	c := &traceConfig{
		Propagators:    otel.GetTextMapPropagator(),
		TracerProvider: otel.GetTracerProvider(),
	}
	for _, o := range opts {
		o.apply(c)
	}
	return c
}

type propagatorsOption struct{ p propagation.TextMapPropagator }

func (o propagatorsOption) apply(c *traceConfig) {
	if o.p != nil {
		c.Propagators = o.p
	}
}

// WithPropagators returns an Option to use the Propagators when extracting
// and injecting trace context from requests.
func WithPropagators(p propagation.TextMapPropagator) Option {
	return propagatorsOption{p: p}
}

type tracerProviderOption struct{ tp trace.TracerProvider }

func (o tracerProviderOption) apply(c *traceConfig) {
	if o.tp != nil {
		c.TracerProvider = o.tp
	}
}

// WithTracerProvider returns an Option to use the TracerProvider when
// creating a Tracer.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return tracerProviderOption{tp: tp}
}

type metadataSupplier struct {
	metadata metadata.Metadata
}

// assert that metadataSupplier implements the TextMapCarrier interface.
var _ propagation.TextMapCarrier = &metadataSupplier{}

func (s *metadataSupplier) Get(key string) string {
	values, ok := s.metadata.Get(key)
	if !ok {
		return ""
	}
	return values
}

func (s *metadataSupplier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(s.metadata))
	for key := range s.metadata {
		out = append(out, key)
	}
	return out
}

// Inject injects correlation context and span context into the gRPC
// metadata object. This function is meant to be used on outgoing
// requests.
func Inject(ctx context.Context, md metadata.Metadata, opts ...Option) {
	c := newConfig(opts)
	c.Propagators.Inject(ctx, &metadataSupplier{
		metadata: md,
	})
}

// Extract returns the correlation context and span context that
// another service encoded in the gRPC metadata object with Inject.
// This function is meant to be used on incoming requests.
func Extract(ctx context.Context, md metadata.Metadata, opts ...Option) (baggage.Baggage, trace.SpanContext) {
	c := newConfig(opts)
	ctx = c.Propagators.Extract(ctx, &metadataSupplier{
		metadata: md,
	})

	return baggage.FromContext(ctx), trace.SpanContextFromContext(ctx)
}

func SpanIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}

	return ""
}

func TraceIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}
