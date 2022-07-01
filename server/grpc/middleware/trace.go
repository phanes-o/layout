package middleware

import (
	"context"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel"
)

type trace struct {
	client.Client
}

func (t *trace) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {

	return nil
}

func ClientTraceWrapper() client.Client {
	return &trace{}
}

func ServerTraceWrapper() server.HandlerWrapper {
	return func(handlerFunc server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {

			ctx, span := otel.Tracer("go.micro.server").Start(ctx, req.Method())
			defer span.End()

			return nil
		}
	}
}
