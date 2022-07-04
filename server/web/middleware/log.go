package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	log "phanes/collector/logger"
	"phanes/errors"
	"time"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err      error
			token    = c.GetHeader("Authorization")
			tracer   = otel.GetTracerProvider().Tracer("http-request")
			spanName = fmt.Sprintf("%s-%s", c.Request.URL, c.Request.Method)
		)

		ctx, span := tracer.Start(c.Request.Context(), spanName)
		spanCtx := span.SpanContext()
		traceID := spanCtx.TraceID()
		spanID := spanCtx.SpanID()
		traceFlags := spanCtx.TraceFlags()

		span.SetAttributes(attribute.String("url", c.Request.URL.String()))
		span.SetAttributes(attribute.String("remote_addr", c.Request.RemoteAddr))
		span.SetAttributes(attribute.String("token", token))
		c.Request = c.Request.WithContext(ctx)

		defer func() {
			if err != nil {
				span.SetAttributes(attribute.String("err:", err.Error()))
			}
			span.End()
		}()

		start := time.Now()
		log.WithContext(c.Request.Context()).WithFields(log.Fields{
			"timestamp":   time.Now().UnixNano(),
			"trace_id":    traceID,
			"span_id":     spanID,
			"trace_flags": traceFlags,
			"url":         c.Request.URL.String(),
			"method":      c.Request.Method,
			"ip":          c.ClientIP(),
		}).Info("request")

		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				err = e
				errType := errors.GetType(e.Err)
				if errType == errors.None {
					c.JSON(500, gin.H{
						"trace_id": traceID,
						"code":     500,
						"message":  "server internal error, more info in logs",
					})
				} else if errType > 3000 && errType < 5000 {
					c.JSON(400, gin.H{
						"trace_id": traceID,
						"code":     errType,
						"msg":      err.Error(),
					})
				}

				log.WithContext(c.Request.Context()).WithFields(log.Fields{
					"timestamp":   time.Now().UnixNano(),
					"trace_id":    traceID,
					"span_id":     spanID,
					"trace_flags": traceFlags,
					"url":         c.Request.URL.String(),
					"method":      c.Request.Method,
					"ip":          c.ClientIP(),
					"error_type":  errType.String(),
					"time_span":   time.Now().Sub(start).String(),
				}).Error(err.Error())
			}
		} else {
			log.WithContext(c.Request.Context()).WithFields(log.Fields{
				"timestamp":   time.Now().UnixNano(),
				"trace_id":    traceID,
				"span_id":     spanID,
				"trace_flags": traceFlags,
				"rul":         c.Request.URL.String(),
				"method":      c.Request.Method,
				"ip":          c.ClientIP(),
				"time_span":   time.Now().Sub(start).String(),
			}).Info("response")
		}
	}
}
