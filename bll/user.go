package bll

import (
	"context"
	"github.com/phanes-o/proto/base"
	"github.com/phanes-o/proto/dto"
	log "go-micro.dev/v4/logger"
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

func (a *user) onEvent(ed *event.Data) {

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
	_, err = a.user.Create(u)
	if err != nil {
		log.Error(err)
		return errors.Wrap(err, "user create failed")
	}

	return nil
}

func (a *user) Delete(ctx context.Context, p *base.Int64) error {
	return a.user.Delete(p.Value)
}
