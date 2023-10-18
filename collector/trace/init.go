package trace

import (
	"context"

	libtrace "github.com/phanes-o/lib/otel/trace"
	"github.com/phanes-o/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"phanes/config"
)

var (
	defaultJaegerAddr   = "http://localhost:14268/api/traces"
	defaultZipkinAddr   = "http://localhost:9411/api/v2/spans"
	defaultOTLPHttpAddr = "http://localhost:4317/v1/traces"
	defaultOTLPGrpcAddr = "http://localhost:4318/v1/traces"
)

func Init() func() {
	if !config.Conf.Collect.Trace.Enabled {
		return func() {}
	}
	
	var (
		url          string
		providerType libtrace.ProviderType
		options      []libtrace.ProviderOption
	)

	switch config.Conf.Collect.Trace.Type {
	case "jaeger":
		url = defaultJaegerAddr
		providerType = libtrace.ProviderTypeJaeger
	case "zipkin":
		url = defaultZipkinAddr
		providerType = libtrace.ProviderTypeZipkin
	case "otlp":
		switch config.Conf.Collect.Trace.Protocol {
		case "grpc":
			url = defaultOTLPGrpcAddr
			options = append(options, libtrace.WithInsecure(true))
			options = append(options, libtrace.WithOtelProtocol(libtrace.OtelGRPC))
		case "http":
			url = defaultOTLPHttpAddr
			options = append(options, libtrace.WithOtelProtocol(libtrace.OtelHTTP))

		}
		providerType = libtrace.ProviderTypeOTLP
	}

	if config.Conf.Collect.Trace.Addr == "" {
		url = config.Conf.Collect.Trace.Addr
	}

	options = append(options, libtrace.WithURL(url))
	options = append(options, libtrace.WithName(config.Conf.Name))
	options = append(options, libtrace.WithEnv(config.Conf.Env))
	options = append(options, libtrace.WithVersion(config.Conf.Version))

	provider, err := libtrace.NewTraceProvider(providerType, options...)
	utils.Throw(err)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetTracerProvider(provider)

	return func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			utils.Throw(err)
		}
	}
}
