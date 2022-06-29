package middleware

import (
	"context"
	"go-micro.dev/v4/server"
	log "phanes/collector/logger"
)

func Log() server.HandlerWrapper {
	return func(next server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			log.WithFields(log.Fields{
				"method":       req.Method(),
				"service":      req.Service(),
				"endpoint":     req.Endpoint(),
				"content-type": req.ContentType(),
				"header":       req.Header(),
				"body":         req.Body(),
			}).Info("request")
			return next(ctx, req, rsp)
		}
	}
}
