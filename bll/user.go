package bll

import (
	"context"

	"github.com/phanes-o/proto/base"
	"github.com/phanes-o/proto/dto"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	log "phanes/collector/logger"
	"phanes/errors"
	"phanes/event"
	"phanes/model/entity"
	"phanes/store"
	"phanes/store/postgres"
)

var User = &user{}

type user struct {
	user store.IUser
}

func (a *user) init() func() {
	a.user = postgres.NewUser()
	return func() {}
}

func (a *user) Create(ctx context.Context, in *dto.CreateUserRequest) (err error) {
	u := &entity.User{
		Username: in.Username,
		Password: in.Password,
	}

	p := otel.GetTracerProvider()
	tracer := p.Tracer("bll")
	ctx, span := tracer.Start(ctx, "Bll.User.Create")
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}()

	_, err = a.user.Create(ctx, u)
	if err != nil {
		log.ErrorCtx(ctx, "[bll] create user failed", zap.String("err_info", err.Error()))
		return errors.Wrap(err, "user create failed")
	}

	// publish event to event bus
	event.Bus.PublishAsync(event.ExampleEvent, u.Username)
	return nil
}

func (a *user) Delete(ctx context.Context, p *base.Int64) error {
	return a.user.Delete(p.Value)
}
