package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"phanes/config"
	"phanes/lib/translation"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	log "phanes/collector/logger"
	"phanes/collector/metrics"
	"phanes/errors"
	"phanes/lib/trace"
)

var (
	defaultValidateTrans = "en"
	translate            ut.Translator
	lock                 sync.Mutex
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
		traceID = trace.TraceIDFromContext(c.Request.Context())
	)

	if config.Conf.Http.ValidateTrans != "" {
		defaultValidateTrans = config.Conf.Http.ValidateTrans
	}

	if translate == nil {
		lock.Lock()
		defer lock.Unlock()
		if translate, err = translation.InitTrans(defaultValidateTrans); err != nil {
			return err
		}
	}

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
						"msg":      RemoveTopStruct(errs.Translate(translate)),
					})
					traceLabel := prometheus.Labels{"TraceID": traceID, "StatusCode": "400"}
					metrics.Http.ResponseCodeCounterInc(traceLabel)
				} else {
					// some can't show error
					c.JSON(500, gin.H{
						"trace_id": traceID,
						"code":     500,
						"msg":      "Server Internal Error",
					})
					traceLabel := prometheus.Labels{"TraceID": traceID, "StatusCode": "500"}
					metrics.Http.ResponseCodeCounterInc(traceLabel)
				}
			} else if errType == 1000 {
				c.JSON(http.StatusUnauthorized, nil)
				traceLabel := prometheus.Labels{"TraceID": traceID, "StatusCode": "401"}
				metrics.Http.ResponseCodeCounterInc(traceLabel)
				// hello Common errors handle
			} else if errType > 1000 && errType < 2000 {
				c.JSON(http.StatusOK, gin.H{
					"trace_id": traceID,
					"code":     errType,
					"msg":      e.Error(),
				})
				traceLabel := prometheus.Labels{"TraceID": traceID, "StatusCode": "200"}
				metrics.Http.ResponseCodeCounterInc(traceLabel)
				// customer error handle here
			} else if errType >= 2000 && errType < 3000 {
				c.JSON(http.StatusOK, gin.H{
					"trace_id": traceID,
					"code":     errType,
					"msg":      e.Error(),
				})
				traceLabel := prometheus.Labels{"TraceID": traceID, "StatusCode": "200"}
				metrics.Http.ResponseCodeCounterInc(traceLabel)
			}
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
