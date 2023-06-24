package middleware

import (
	"bytes"
	"fmt"
	"time"

	"go.uber.org/zap"
	"phanes/utils"

	log "phanes/collector/logger"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func LogAndTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err      error
			tracer   = otel.GetTracerProvider().Tracer("http-request")
			spanName = fmt.Sprintf("%s-%s", c.Request.URL, c.Request.Method)
			params   = GetRequestParams(c)
		)

		ctx, span := tracer.Start(c.Request.Context(), spanName)
		traceID := span.SpanContext().TraceID().String()

		span.SetAttributes(attribute.String("url", c.Request.URL.String()))
		span.SetAttributes(attribute.String("remote_addr", c.Request.RemoteAddr))
		span.SetAttributes(attribute.String("request_params", utils.ToJsonString(params)))
		c.Request = c.Request.WithContext(ctx)

		defer func() {
			if err != nil {
				span.RecordError(err)
			}

			span.End()
		}()

		newWriter := customWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = newWriter
		c.Next()
		err = HandleResponse(c)

		logger := log.WithFields(
			zap.String("url", c.Request.URL.String()),
			zap.String("method", c.Request.Method),
			zap.String("trace_id", traceID),
			zap.String("span_id", span.SpanContext().SpanID().String()),
			zap.String("trace_flag", span.SpanContext().TraceFlags().String()),
			zap.String("request", utils.ToJsonString(params)),
			zap.String("response", newWriter.body.String()),
			zap.Int64("timestamp", time.Now().UnixNano()),
			zap.Int("response-status", c.Writer.Status()),
		)

		if err != nil {
			logger.ErrorCtx(c.Request.Context(), "[http] request failed ", zap.String("err_info", err.Error()))
		} else {
			logger.InfoCtx(c.Request.Context(), "[http] request success")
		}

		resp := newWriter.body.String()
		span.SetAttributes(attribute.String("response_body", resp))
		span.SetAttributes(attribute.String("response_status", fmt.Sprintf("%d", c.Writer.Status())))

	}
}

type customWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (c customWriter) Write(p []byte) (int, error) {
	if _, err := c.ResponseWriter.Write(p); err != nil {
		return 0, err
	}
	return c.body.Write(p)
}
