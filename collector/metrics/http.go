package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"phanes/config"
)

var Http = &httpRecorder{}

type httpRecorder struct {
	httpRequestTimeDuration prometheus.Histogram
	httpRequestCounter      prometheus.Counter
	errorCounter            prometheus.Counter
	concurrentConnections   prometheus.Gauge
	responseCodeCounter     *prometheus.CounterVec
}

func (hr *httpRecorder) Init() []prometheus.Collector {
	var collectors = make([]prometheus.Collector, 0, 0)

	// HTTP response time
	hr.httpRequestTimeDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: config.Conf.Name,
		Name:      "http_response_time_seconds",
		Help:      "HTTP response time",
		Buckets:   prometheus.DefBuckets,
	})
	collectors = append(collectors, hr.httpRequestTimeDuration)

	// Total number of HTTP requests
	hr.httpRequestCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: config.Conf.Name,
		Name:      "http_requests_total",
		Help:      "Total number of HTTP requests",
	})
	collectors = append(collectors, hr.httpRequestCounter)

	// Total number of HTTP errors
	hr.errorCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "Total number of HTTP errors",
	})
	collectors = append(collectors, hr.errorCounter)

	// Number of concurrent HTTP connections
	hr.concurrentConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "http_concurrent_connections",
		Help: "Number of concurrent HTTP connections",
	})
	collectors = append(collectors, hr.concurrentConnections)

	// Total number of HTTP response codes
	hr.responseCodeCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_response_codes_total",
		Help: "Total number of HTTP response codes",
	}, []string{"StatusCode", "TraceID"})
	collectors = append(collectors, hr.responseCodeCounter)

	return collectors
}

func (hr *httpRecorder) RequestTimeDurationRecord(val float64, labels prometheus.Labels) {
	if exemplarObserver, ok := hr.httpRequestTimeDuration.(prometheus.ExemplarObserver); ok {
		exemplarObserver.ObserveWithExemplar(val, labels)
	} else {
		hr.httpRequestTimeDuration.Observe(val)
	}
}

func (hr *httpRecorder) RequestCountInc(labels prometheus.Labels) {
	if exemplarObserver, ok := hr.httpRequestCounter.(prometheus.ExemplarObserver); ok {
		exemplarObserver.ObserveWithExemplar(1, labels)
	} else {
		hr.httpRequestCounter.Inc()
	}
}

func (hr *httpRecorder) ConcurrentConnectionSet(val float64, labels prometheus.Labels) {
	if exemplarObserver, ok := hr.concurrentConnections.(prometheus.ExemplarObserver); ok {
		exemplarObserver.ObserveWithExemplar(val, labels)
	} else {
		hr.concurrentConnections.Inc()
	}
}

func (hr *httpRecorder) RequestErrorInc(labels prometheus.Labels) {
	if exemplarObserver, ok := hr.errorCounter.(prometheus.ExemplarObserver); ok {
		exemplarObserver.ObserveWithExemplar(1, labels)
	} else {
		hr.errorCounter.Inc()
	}
}

func (hr *httpRecorder) ResponseCodeCounterInc(labels prometheus.Labels) {
	hr.responseCodeCounter.With(labels).Inc()
}
