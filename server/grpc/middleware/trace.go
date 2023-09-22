package middleware

import (
	"context"
	"fmt"

	"go-micro.dev/v4/client"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
	libTrace "phanes/lib/trace"
	"phanes/utils"
)

func ClientTraceWrapper() client.CallWrapper {
	return func(callFunc client.CallFunc) client.CallFunc {
		return func(ctx context.Context, node *registry.Node, req client.Request, rsp interface{}, opts client.CallOptions) error {
			var err error
			md, _ := metadata.FromContext(ctx)
			m := metadata.Copy(md)

			tracer := otel.Tracer("go.micro.client")
			spanName := fmt.Sprintf("%s:%s", req.Service(), req.Endpoint())

			RPCSystemGRPC := semconv.RPCSystemKey.String("grpc")

			attrs := []attribute.KeyValue{
				RPCSystemGRPC,
				attribute.String("rpc.system", "grpc"),
				attribute.String("rpc.service", req.Service()),
				attribute.String("rpc.method", req.Method()),
				attribute.String("rpc.method", utils.ToJsonString(req.Body())),
			}
			ctx, span := tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(attrs...))

			libTrace.Inject(ctx, m)
			ctx = metadata.NewContext(ctx, m)


			defer func() {
				span.SetAttributes(
					attribute.String("rpc.req", utils.ToJsonString(rsp)))
				if err != nil {
					span.SetAttributes(
						// 设置事件为异常
						attribute.String("event", "error"),
						// 设置 message 为 err.Error().
						attribute.String("message", err.Error()),
					)
					span.SetStatus(codes.Error, err.Error())
				} else {
					// 如果没有发生异常，span 状态则为 ok
					span.SetStatus(codes.Ok, "OK")
				}
				span.End()
			}()

			err = callFunc(ctx, node, req, rsp, opts)
			return err
		}
	}
}

func ServerTraceWrapper() server.HandlerWrapper {
	return func(handlerFunc server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			var (
				err error
			)
			requestMetadata, ok := metadata.FromContext(ctx)
			if !ok {
				fmt.Println("no metadata")
			}
			metadataCopy := metadata.Copy(requestMetadata)
			bags, spanCtx := libTrace.Extract(ctx, metadataCopy)
			ctx = baggage.ContextWithBaggage(ctx, bags)

			tracer := otel.Tracer(
				"go.micro.server",
				trace.WithInstrumentationVersion("1.0"),
			)
			spanName := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
			RPCSystemGRPC := semconv.RPCSystemKey.String("grpc")
			attrs := []attribute.KeyValue{RPCSystemGRPC}
			attrs = append(attrs, attribute.String("rpc.system", "grpc-manager-server"))
			attrs = append(attrs, attribute.String("rpc.service", req.Service()))
			attrs = append(attrs, attribute.String("rpc.method", req.Method()))

			ctx, span := tracer.Start(
				trace.ContextWithRemoteSpanContext(ctx, spanCtx),
				spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(attrs...),
			)
			defer func() {
				span.SetAttributes(
					attribute.String("rpc.req", utils.ToJsonString(req.Body())))
				if err != nil {
					span.SetAttributes(
						// 设置事件为异常
						attribute.String("event", "error"),
						// 设置 message 为 err.Error().
						attribute.String("message", err.Error()),
					)
					span.SetStatus(codes.Error, err.Error())
				} else {
					// 如果没有发生异常，span 状态则为 ok
					span.SetStatus(codes.Ok, "OK")
				}
				span.End()

			}()
			err = handlerFunc(ctx, req, rsp)
			return err

		}
	}
}
