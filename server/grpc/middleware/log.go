package middleware

import (
	"context"
	"go-micro.dev/v4/server"
	log "phanes/collector/logger"
	"time"
)

func Log() server.HandlerWrapper {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			start := time.Now()
			log.WithContext(ctx).WithFields(log.Fields{
				"method":       req.Method(),
				"service":      req.Service(),
				"endpoint":     req.Endpoint(),
				"content-type": req.ContentType(),
				"header":       req.Header(),
				"body":         req.Body(),
			}).Info("request")
			if err := next(ctx, req, rsp); err != nil {
				log.WithContext(ctx).WithFields(log.Fields{
					"method":       req.Method(),
					"service":      req.Service(),
					"endpoint":     req.Endpoint(),
					"content-type": req.ContentType(),
					"header":       req.Header(),
					"body":         req.Body(),
					"time_span":    time.Now().Sub(start).String(),
				}).Error(err.Error())
				return err
			} else {
				log.WithContext(ctx).WithFields(log.Fields{
					"method":       req.Method(),
					"service":      req.Service(),
					"endpoint":     req.Endpoint(),
					"content-type": req.ContentType(),
					"header":       req.Header(),
					"body":         req.Body(),
					"time_span":    time.Now().Sub(start).String(),
				}).Info("response")
			}
			return nil
		}
	}
}
