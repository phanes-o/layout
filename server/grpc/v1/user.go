package v1

import (
	"context"
	"github.com/phanes-o/proto/base"
	"github.com/phanes-o/proto/dto"
	"phanes/bll"
)

type User struct{}

func (u *User) Create(ctx context.Context, request *dto.CreateUserRequest, empty *base.Empty) error {
	return bll.User.Create(ctx, request)
}

func (u *User) Delete(ctx context.Context, p *base.Int64, empty *base.Empty) error {
	return bll.User.Delete(ctx, p)
}
