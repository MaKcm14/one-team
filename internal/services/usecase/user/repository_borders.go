package user

import (
	"context"

	entity "github.com/MaKcm14/one-team/internal/entity/user"
)

type IAuthRepo interface {
	GetUser(ctx context.Context, login string) (entity.User, error)
}
