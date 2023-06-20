package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
	"phanes/errors"
	"phanes/utils"

	log "phanes/collector/logger"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func LogAndTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err           error
			tracer        = otel.GetTracerProvider().Tracer("http-request")
			spanName      = fmt.Sprintf("%s-%s", c.Request.URL, c.Request.Method)
			requestParams = make(map[string]interface{})
		)

		ctx, span := tracer.Start(c.Request.Context(), spanName)
		spanCtx := span.SpanContext()
		traceID := spanCtx.TraceID()
		spanID := spanCtx.SpanID()
		traceFlags := spanCtx.TraceFlags()

		span.SetAttributes(attribute.String("url", c.Request.URL.String()))
		span.SetAttributes(attribute.String("remote_addr", c.Request.RemoteAddr))
		c.Request = c.Request.WithContext(ctx)

		defer func() {
			if err != nil {
				span.RecordError(err)
			}
			span.End()
		}()

		if c.Request.Method == http.MethodPost {
			buffer, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				log.ErrorCtx(c, "read request body error", zap.String("err", err.Error()))
			}

			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))
			if err = json.Unmarshal(buffer, &requestParams); err != nil {
				log.ErrorCtx(c, "unmarshal request body error", zap.String("err", err.Error()))
			}
		}

		if c.Request.Method == http.MethodGet {
			params := strings.Split(c.Request.URL.String(), "?")
			if len(params) > 1 {
				kvs := strings.Split(params[1], "&")
				for _, kv := range kvs {
					kvs := strings.Split(kv, "=")
					if len(kvs) > 1 {
						requestParams[kvs[0]] = kvs[1]
					}
				}
			}
		}

		newWriter := customWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = newWriter

		span.SetAttributes(attribute.String("request_params", utils.ToJsonString(requestParams)))

		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				err = e
				errType := errors.GetType(e.Err)
				if errType == errors.None {
					// check request params
					if errs, ok := e.Err.(validator.ValidationErrors); ok {
						c.JSON(400, gin.H{
							"trace_id": traceID,
							"code":     errType,
							"message":  removeTopStruct(errs.Translate(trans)),
						})
					} else {
						// some can't show error
						c.JSON(500, gin.H{
							"trace_id": traceID,
							"code":     500,
							"msg":      "Server Internal Error",
						})
					}
				} else if errType == 1000 {
					c.JSON(http.StatusUnauthorized, nil)
					// phanes Common errors handle
				} else if errType > 1000 && errType < 2000 {
					c.JSON(http.StatusOK, gin.H{
						"trace_id": traceID,
						"code":     errType,
						"message":  e.Error(),
					})
					// customer error handle here
				} else if errType >= 2000 && errType < 3000 {
					c.JSON(http.StatusOK, gin.H{
						"trace_id": traceID,
						"code":     errType,
						"message":  e.Error(),
					})
				}
			}
		}
		l := log.WithFields(
			zap.String("url", c.Request.URL.String()),
			zap.String("method", c.Request.Method),
			zap.String("trace_id", traceID.String()),
			zap.String("span_id", spanID.String()),
			zap.String("trace_flag", traceFlags.String()),
			zap.String("request", newWriter.body.String()),
			zap.String("response", newWriter.body.String()),
			zap.Int64("timestamp", time.Now().UnixNano()),
			zap.Int("response-status", c.Writer.Status()),
		)

		if err != nil {
			l.ErrorCtx(c, "[http] request failed ", zap.String("err_info", err.Error()))
		} else {
			l.Info("[http] request success")
		}

		span.SetAttributes(attribute.String("response_status", fmt.Sprintf("%d", c.Writer.Status())))
		resp := newWriter.body.String()
		span.SetAttributes(attribute.String("response_body", resp))
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
