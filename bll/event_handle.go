package bll

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	log "phanes/collector/logger"
	"phanes/event"
)

func init() {
	Register(&eventHandle{})
}

type eventHandle struct {
}

func (h *eventHandle) init() func() {
	// register event handler
	event.Subscribe(event.ExampleEvent, h.ExampleEventHandler)

	return func() {}
}

func (h *eventHandle) ExampleEventHandler(ctx context.Context, e event.Event, data interface{}) {
	ctx, span := otel.GetTracerProvider().Tracer("event").Start(ctx, "ExampleEventHandler")
	defer func() {
		span.End()
	}()

	log.InfoCtx(ctx, "Got event:", zap.String("event", string(e)), zap.Any("data", data))
	// handle event
}
