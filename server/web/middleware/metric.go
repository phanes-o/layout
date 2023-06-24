package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	api "go.opentelemetry.io/otel/metric"
	"phanes/collector/metrics"
	"phanes/config"
	"phanes/lib/trace"
	"time"
)

// BaseMetric is a example metrics middleware
func BaseMetric() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request count
		trace_id := trace.TraceIDFromContext(c.Request.Context())
		value := attribute.String("trace_id", trace_id)

		counter, _ := metrics.Meter.Int64Counter(fmt.Sprintf("%s_request_count", config.Conf.Name), api.WithDescription("Number of requests"))
		counter.Add(c.Request.Context(), 1, api.WithAttributes(value))

		start := time.Now()
		histogram, _ := metrics.Meter.Float64Histogram(fmt.Sprintf("%s_request_time_seconds", config.Conf.Name), api.WithDescription("Response time in seconds"))

		c.Next()
		histogram.Record(c.Request.Context(), time.Since(start).Seconds(), api.WithAttributes(value))

	}
}
