package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	api "go.opentelemetry.io/otel/metric"
	"phanes/collector/metrics"
)

// BaseMetric is a example metrics middleware
func BaseMetric() gin.HandlerFunc {
	return func(c *gin.Context) {
		// request count
		counter, _ := metrics.Meter.Int64Counter(fmt.Sprintf("http_request_count_%s", c.Request.URL.String()), api.WithDescription("a simple counter"))
		counter.Add(c.Request.Context(), 1)
		start := time.Now()
		c.Next()
		histogram, _ := metrics.Meter.Float64Histogram(fmt.Sprintf("http_request_seconds_%s", c.Request.URL.String()), api.WithDescription("a very nice histogram"))
		histogram.Record(c.Request.Context(), time.Now().Sub(start).Seconds())
	}
}
