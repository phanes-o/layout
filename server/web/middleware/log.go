package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"phanes/errors"
	"phanes/utils"
	
	log "phanes/collector/logger"


	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err           error
			token         = c.GetHeader("Authorization")
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
		span.SetAttributes(attribute.String("token", token))
		c.Request = c.Request.WithContext(ctx)

		defer func() {
			if err != nil {
				span.SetAttributes(attribute.String("err:", err.Error()))
			}
			span.End()
		}()

		if c.Request.Method == http.MethodPost {
			buffer, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				log.Errorf("read request body error [%v]", err)
			}

			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))
			if err = json.Unmarshal(buffer, &requestParams); err != nil {
				log.Errorf("unmarshal request body error [%v]", err)
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
						c.JSON(500, gin.H{
							"trace_id": traceID,
							"code":     500,
							"msg":      "Server Internal Error",
						})
					}
				} else if errType == 1000 {
					c.JSON(http.StatusUnauthorized, nil)
				} else if errType > 1000 && errType < 5000 {
					c.JSON(http.StatusOK, gin.H{
						"trace_id": traceID,
						"code":     errType,
						"message":  errType.String(),
					})
				}
			}
		}

		l := log.WithFields(log.Fields{
			"url":             c.Request.URL.String(),
			"token":           token,
			"method":          c.Request.Method,
			"span_id":         spanID,
			"trace_id":        traceID,
			"request":         utils.ToJsonString(requestParams),
			"response":        newWriter.body.String(),
			"timestamp":       time.Now().UnixNano(),
			"trace_flags":     traceFlags,
			"response-status": c.Writer.Status(),
		})

		if err != nil {
			l.Error("request failed ", err)
		} else {
			l.Info("request success")
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
