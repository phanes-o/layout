package store

import (
	"phanes/model/entity"
)

type IUser interface {
	Create(u *entity.User) (id int64, err error)
	Find(id int64) (user *entity.User, err error)
	Update(id int64, updates map[string]interface{}) (err error)
	Delete(id int64) (err error)
	List(opts map[string]interface{}) (users []*entity.User, err error)
}
