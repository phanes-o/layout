package metrics

import (
	"net/http"

	client "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	metricApi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/zap"
	log "phanes/collector/logger"
	"phanes/config"
)

var defaultListenAddr = ":2223"
var Meter metricApi.Meter

func Init() func() {
	// Create non-global registry.
	registry := client.NewRegistry()
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	reader, err := prometheus.New(prometheus.WithRegisterer(registry))
	if err != nil {
		log.Fatal(err.Error())
	}

	provider := metric.NewMeterProvider(metric.WithReader(reader))
	Meter = provider.Meter(
		config.Conf.Name,
		metricApi.WithInstrumentationVersion(config.Conf.Version),
	)

	if config.Conf.Collect.Metric.Listen != "" {
		defaultListenAddr = config.Conf.Collect.Metric.Listen
	}
	go serveMetrics(defaultListenAddr, registry)
	return func() {}
}

func serveMetrics(addr string, registry *client.Registry) {
	log.Info("serving metrics at /metrics", zap.String("addr", addr))
	http.Handle("/metrics", promhttp.HandlerFor(
		registry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("error serving http: ", zap.String("err", err.Error()))
		return
	}
}
