package bll

import (
	log "go-micro.dev/v4/logger"
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

func (a *user) Create(u *entity.User) (err error) {
	log.Debugf("create new user", u)

	_, err = a.user.Create(u)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
