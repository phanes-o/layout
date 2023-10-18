package middleware

import (
	"context"
	"time"

	"github.com/phanes-o/utils"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	log "phanes/collector/logger"
)

func Log() server.HandlerWrapper {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			start := time.Now()
			log.WithFields(
				zap.String("method", req.Method()),
				zap.String("service", req.Service()),
				zap.String("endpoint", req.Endpoint()),
				zap.String("content-type", req.ContentType()),
				zap.String("header", utils.ToJsonString(req.Header())),
				zap.String("body", utils.ToJsonString(req.Body())),
			)
			log.Info("request")

			if err := next(ctx, req, rsp); err != nil {
				log.WithFields(zap.String("time_span", time.Now().Sub(start).String())).Error("request error", zap.String("err", err.Error()))
				return err
			} else {
				log.WithFields(zap.String("time_span", time.Now().Sub(start).String())).Info("response")
			}
			return nil
		}
	}
}
