package bll

import (
	"phanes/event"
	"phanes/store"
	"phanes/store/postgres"
)

var User = &user{}

type user struct {
	user store.IUser
}

func (u *user) onEvent(ed *event.Data) {

}

func (u *user) init() func() {
	u.user = postgres.NewUser()

	return func() {}
}
