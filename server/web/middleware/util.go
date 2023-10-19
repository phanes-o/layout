package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"phanes/config"
	"phanes/lib/translation"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/phanes-o/lib/otel/trace"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	log "phanes/collector/logger"
	"phanes/collector/metrics"
	"phanes/errors"
)

var (
	defaultValidateTrans = "en"
	translate            ut.Translator
	once                 sync.Once
)

func GetRequestParams(c *gin.Context) map[string]interface{} {
	var params = make(map[string]interface{})
	switch c.Request.Method {
	case http.MethodGet:
		splits := strings.Split(c.Request.URL.String(), "?")
		if len(params) > 1 {
			kvs := strings.Split(splits[1], "&")
			for _, kv := range kvs {
				kvs := strings.Split(kv, "=")
				if len(kvs) > 1 {
					params[kvs[0]] = kvs[1]
				}
			}
		}
	case http.MethodPost:
		buffer, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.ErrorCtx(c, "read request body error", zap.String("err", err.Error()))
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))
		if err = json.Unmarshal(buffer, &params); err != nil {
			log.ErrorCtx(c, "unmarshal request body error", zap.String("err", err.Error()))
		}
	}

	return params
}

func HandleResponse(c *gin.Context) error {
	var (
		err     error
		conf    = config.Conf.Server.Http
		traceID = trace.TraceIDFromContext(c.Request.Context())
	)

	if conf.ValidateTrans != "" {
		defaultValidateTrans = conf.ValidateTrans
	}

	// initialize translation
	once.Do(func() {
		if translate, err = translation.InitTrans(defaultValidateTrans); err != nil {
			log.ErrorCtx(c, "translation InitTrans error", zap.String("err", err.Error()))
			return
		}
	})

	if len(c.Errors) == 0 {
		return nil
	}

	err = c.Errors[0].Err
	errType := errors.GetType(err)

	errHandler := NewHttpErrorHandler(err)

	if errType == errors.None {
		errHandler.Unexportable(c)

		if config.Conf.Collect.Metric.Enabled {
			traceLabel := prometheus.Labels{"TraceID": traceID, "StatusCode": "500"}
			metrics.Http.ResponseCodeCounterInc(traceLabel)
		}
	} else {
		errHandler.Exportable(c)

		if config.Conf.Collect.Metric.Enabled {
			traceLabel := prometheus.Labels{"TraceID": traceID, "StatusCode": "400"}
			metrics.Http.ResponseCodeCounterInc(traceLabel)
		}
	}

	return err
}

func RemoveTopStruct(fields map[string]string) string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	var str string
	for _, v := range res {
		return v
	}
	return str
}
