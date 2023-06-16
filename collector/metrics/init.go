package metrics

import (
	"net/http"
	log "phanes/collector/logger"
	"phanes/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	metricApi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

var defaultListenAddr = ":2223"
var Meter metricApi.Meter

func Init() func() {
	reader, err := prometheus.New()
	if config.Conf.Collect.Metric.Listen != "" {
		defaultListenAddr = config.Conf.Collect.Metric.Listen
	}

	if err != nil {
		log.Fatal(err)
	}

	Meter = metric.NewMeterProvider(metric.WithReader(reader)).Meter("phanes")
	go serveMetrics(defaultListenAddr)
	return func() {}
}

func serveMetrics(addr string) {
	log.Info("serving metrics at 0.0.0.0:%s/metrics", addr)
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("error serving http: %v", err)
		return
	}
}
