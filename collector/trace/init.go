package trace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"phanes/config"
	"phanes/utils"
)

func Init() func() {
	tp, err := JaegerTraceProvider()
	if err != nil {
		utils.Throw(err)
	}
	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	return func() {}
}

var defaultAddr = "http://localhost:14268/api/traces"

func JaegerTraceProvider() (*trace.TracerProvider, error) {
	if config.Conf.Collect.Trace.Addr != "" {
		defaultAddr = config.Conf.Collect.Trace.Addr
	}
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(defaultAddr)))
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.Conf.Base.Name),
			attribute.String("environment", config.Conf.Base.Env),
			attribute.String("version", config.Conf.Base.Version),
		)),
	)
	return tp, nil
}
