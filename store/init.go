package store

import "phanes/store/postgres"

var (
	User IUser
)

func Init() func() {
	postgres.Init()

	return func() {}
}
