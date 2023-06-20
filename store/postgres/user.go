package postgres

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	log "phanes/collector/logger"
	"phanes/errors"
	"phanes/model/entity"
)

type user struct{}

func NewUser() *user {
	return &user{}
}

func (a *user) Create(ctx context.Context, u *entity.User) (id int64, err error) {
	p := otel.GetTracerProvider()
	tracer := p.Tracer("store")
	ctx, span := tracer.Start(ctx, "Store.User.Create")
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}()
	//err = db.Model(&entity.User{}).Create(u).Error
	err = errors.New("test")
	log.ErrorCtx(ctx, "[store] create user error", zap.String("err_info", err.Error()))
	return u.ID, err
}

func (a *user) Find(id int64) (user *entity.User, err error) {
	u := new(entity.User)
	err = db.Model(&entity.User{}).Find(u, id).Error
	return
}

func (a *user) Update(id int64, updates map[string]interface{}) (err error) {
	return db.Model(&entity.User{}).Where("id = ?", id).Updates(updates).Error
}

func (a *user) Delete(id int64) (err error) {
	return db.Delete(&entity.User{}, id).Error
}

func (a *user) List(opts map[string]interface{}) (users []*entity.User, err error) {
	// todo: implement
	return nil, nil
}
