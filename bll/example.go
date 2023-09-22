package bll

import (
	"context"

	"go.opentelemetry.io/otel"
	log "phanes/collector/logger"
	"phanes/event"
)

func init() {
	Register(Example)
}

var Example = &example{}

type example struct{}

func (e *example) init() func() {
	return func() {}
}

func (e *example) ExamplePublishEvent(ctx context.Context) {
	ctx, span := otel.GetTracerProvider().Tracer("example").Start(ctx, "ExamplePublishEvent")
	span.AddEvent("ExamplePublishEvent")
	defer func() {
		span.End()
	}()

	log.InfoCtx(ctx, "publish event")
	// publish event
	event.Publish(ctx, event.ExampleEvent, "example payload")
}
