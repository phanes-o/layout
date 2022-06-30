package trace

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"phanes/utils"
)

func Init() {
	tp, err := JaegerTraceProvider()
	if err != nil {
		utils.Throw(err)
	}
	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
}

func JaegerTraceProvider() (*trace.TracerProvider, error) {
	// todo: use config to set service name, service ID, jaeger url etc
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("phanes"),
			attribute.String("environment", "production"),
			attribute.Int64("ID", 1),
		)),
	)
	return tp, nil
}
