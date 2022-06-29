package v1

import (
	"context"
	"github.com/phanes-o/proto/dto"
	"github.com/phanes-o/proto/primitive"
	"phanes/bll"
)

type User struct{}

func (u *User) Create(ctx context.Context, request *dto.CreateUserRequest, empty *primitive.Empty) error {
	return bll.User.Create(ctx, request)
}

func (u *User) Delete(ctx context.Context, p *primitive.Int64, empty *primitive.Empty) error {
	return bll.User.Delete(ctx, p)
}
