package root

import (
	"context"

	"github.com/MaKcm14/one-team/internal/services/usecase/user"

	entity "github.com/MaKcm14/one-team/internal/entity/user"
)

type IRootRepo interface {
	GetUser(ctx context.Context, login string) (entity.User, error)
	GetUsers(ctx context.Context) ([]user.UserDTO, error)
	GetRoles(ctx context.Context) ([]Role, error)
	GetUserRole(ctx context.Context, login string) (entity.Role, error)
	DeleteUser(ctx context.Context, login string) error
}
